package controllers

import (
	"backend/eventApp/config"
	"backend/eventApp/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingInput struct {
	Phone   string `json:"phone" binding:"required"`
	EventID int    `json:"eventId" binding:"required"`
}

func CreateBookingEvent(c *gin.Context) {
	userId, _ := c.Get("UserID")

	var input BookingInput
	var booking models.Booking

	//validation
	if errValidation := c.ShouldBindJSON(&input); errValidation != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Error(),
		})
		return
	}

	//kondisi jika user sudah pernah booking eventId
	Booking := config.DB.Where("user_id = ? AND event_id = ?", userId.(int), input.EventID).First(&booking).Error
	if Booking == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "anda sudah booking event ini!",
		})
		return
	}

	//Generate Booking code
	CodeBooking := fmt.Sprintf("BK-%sE%dU%d", time.Now().Format("20060102"), input.EventID, userId.(int))

	//simpan kedalam database
	bookingData := models.Booking{
		Phone:       input.Phone,
		EventID:     input.EventID,
		BookingCode: CodeBooking,
		UserID:      userId.(int),
	}

	errCreateBooking := config.DB.Create(&bookingData).Error
	if errCreateBooking != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "create booking gagal!",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "berhasil daftar event",
		"data":    bookingData,
	})

}

func GetBookingByUser(c *gin.Context) {
	var booking []models.Booking
	userId, _ := c.Get("UserID") //ambil berdasarkan userId

	errBookingData := config.DB.Preload("Event").Preload("Event.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Where("user_id = ?", userId).Find(&booking).Error

	if errBookingData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "event tidak ditemukan",
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "data booking berhasil ditampilkan",
		"data":    booking,
	})

}

func DeleteBooking(c *gin.Context) {
	var booking models.Booking
	paramsId := c.Param("id")

	userId, _ := c.Get("UserID")

	bookingData := config.DB.First(&booking, paramsId).Error
	if bookingData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "booking event tidak ditemukan",
		})
		return
	}

	if booking.UserID != userId.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "kamu bukan pemiliknya, booking event tidak bisa dihapus",
		})
		return
	}

	config.DB.Unscoped().Delete(&booking)
	c.JSON(http.StatusOK, gin.H{
		"message": "data booking event ini berhasil dihapus!",
	})
}
