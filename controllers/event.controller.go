package controllers

import (
	"backend/eventApp/config"
	"backend/eventApp/models"
	"context"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
	"gorm.io/gorm"
)

// initial imagekit
func initImagekit() *imagekit.Client {
	client := imagekit.NewClient(
		option.WithPrivateKey(os.Getenv("IMAGEKIT_PRIVATE_KEY")),
	)

	return &client
}

// func GetAllEvents(c *gin.Context) {
// 	var events []models.Event

// 	config.DB.Find(&events)
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "success to get all data events success",
// 		"data":    events,
// 	})
// }

// dengan feature searching
func GetAllEvents(c *gin.Context) {
	var events []models.Event

	// inisasi dasar query di gorm
	query := config.DB.Model(&models.Event{})

	//tangkap function filter by query
	search := c.Query("search")
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE  ?", "%"+search+"%", "%"+search+"%")
	}

	// hitung total data sebelum di limit
	var totalRows int64
	query.Count(&totalRows)

	//tangkap parameter query dan masukkan nilai default
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "2")

	page, errPage := strconv.Atoi(pageStr)
	if errPage != nil || page < 1 {
		page = 1
	}

	limit, errLimit := strconv.Atoi(limitStr)
	if errLimit != nil || limit < 1 {
		limit = 5
	}

	//hitung offset
	offset := (page - 1) * limit

	//hitung data  halaman / page
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	//load semua feature
	if err := query.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Order("name DESC").Limit(limit).Offset(offset).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch data events",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success to get all data events success",
		"data":    events,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalRows":  totalRows,
			"totalPages": totalPages,
		},
	})
}

func CreateEvent(c *gin.Context) {
	userId, _ := c.Get("UserID")

	name := c.PostForm("name")
	description := c.PostForm("description")
	location := c.PostForm("location")
	dateTimeStr := c.PostForm("datetime")

	errs := map[string]string{}

	if name == "" {
		errs["name"] = "nama wajib diisi"
	}
	if description == "" {
		errs["description"] = "deskripsi wajib diisi"
	}
	if location == "" {
		errs["location"] = "lokasi wajib diisi"
	}

	var parsedTime time.Time
	if dateTimeStr == "" {
		errs["datetime"] = "tanggal & waktu wajib diisi"
	} else {
		var errParse error
		parsedTime, errParse = time.Parse(time.RFC3339, dateTimeStr)
		if errParse != nil {
			errs["datetime"] = "format harus RFC3339, contoh: 2026-10-18T10:00:00Z"
		}
	}

	// validasi gambar
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		errs["image"] = "gambar wajib diisi"
	} else {
		defer file.Close()

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/webp": true,
		}
		if !allowedTypes[header.Header.Get("Content-Type")] {
			errs["image"] = "format gambar harus JPEG, PNG, atau WebP"
		}
	}

	// kirim semua error sekaligus
	if len(errs) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "validasi gagal",
			"errors":  errs,
		})
		return
	}

	// semua valid, baru upload ke imagekit
	ik := initImagekit()
	uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: header.Filename,
	})
	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal upload gambar",
		})
		return
	}

	// simpan ke database
	event := models.Event{
		Name:        name,
		Description: description,
		Location:    location,
		DateTime:    parsedTime,
		Image:       uploadRes.URL,
		ImageID:     uploadRes.FileID,
		UserID:      userId.(int),
	}

	if err := config.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal menyimpan event",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "data event insert success",
		"data":    event,
	})
}

func GetEventById(c *gin.Context) {
	var event models.Event
	paramsId := c.Param("id")

	var eventData = config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Preload("Booking").Preload("Booking.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).First(&event, paramsId).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "event tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data by ID berhasil ditampilkan",
		"data":    event,
	})

}

func UpdateEventById(c *gin.Context) {
	var event models.Event
	paramsId := c.Param("id")

	userId, _ := c.Get("UserID") //cek berdasarkan userId

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "event tidak ditemukan",
		})
		return
	}

	if event.UserID != userId.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "kamu tidak bisa mengupdate event user lain",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "gambar wajib diisi",
		})
		return
	}

	defer file.Close() //close jika tidak jadi upload file

	// validasi tipe file dulu sebelum upload
	contentType := header.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format gambar harus JPEG, PNG, atau WebP",
		})
		return
	}

	//upload file
	fileName := header.Filename
	ik := initImagekit()
	uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: fileName,
	})

	if errUpload == nil {
		//hapus gambar lama
		if event.ImageID != "" {
			ik.Files.Delete(context.Background(), event.ImageID)
		}

		// upload gambar baru
		event.Image = uploadRes.URL
		event.ImageID = uploadRes.FileID
	}

	if name := c.PostForm("name"); name != "" {
		event.Name = name
	}
	if description := c.PostForm("description"); description != "" {
		event.Description = description
	}
	if location := c.PostForm("location"); location != "" {
		event.Location = location
	}
	if dateTimeStr := c.PostForm("datetime"); dateTimeStr != "" {
		parseTime, errParse := time.Parse(time.RFC3339, dateTimeStr)
		if errParse == nil {
			event.DateTime = parseTime
		}
	}

	// config.DB.Model(&event).Updates(input)
	config.DB.Save(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data event berhasil diupdate",
		"data":    event,
	})

}

func DeleteEventById(c *gin.Context) {
	var event models.Event
	paramsId := c.Param("id")

	userId, _ := c.Get("UserID") //cek berdasarkan userId

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event tidak ditemukan",
		})
		return
	}

	if event.UserID != userId.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "kamu tidak bisa menghapus event user lain",
		})
		return
	}

	//hapus file gambar
	if event.ImageID != "" {
		ik := initImagekit()
		ik.Files.Delete(context.Background(), event.ImageID)
	}

	// config.DB.Delete(&event) cuman softdelete
	config.DB.Unscoped().Delete(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data event berhasil dihapus",
	})
}
