package controllers

import (
	"backend/eventApp/config"
	"backend/eventApp/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthInputRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthInputLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func RegisterUser(c *gin.Context) {
	var input AuthInputRegister

	//validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//hash password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to encrypt password",
		})
		return
	}

	//simpan ke DB
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	userCreated := config.DB.Create(&user).Error
	if userCreated != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email sudah terdaftar",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil register",
		"data": gin.H{
			"id":     user.ID,
			"email":  user.Email,
			"name":   user.Name,
			"events": user.Events,
		},
	})

}

func LoginUser(c *gin.Context) {

	var input AuthInputLogin
	//validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//cek apakah user sudah ada dlm database
	var user models.User
	userData := config.DB.Where("email=?", input.Email).First(&user).Error
	if userData != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email doesn't exists",
		})
		return
	}

	//compare password
	matchPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if matchPassword != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Password doesn't matched",
		})
		return
	}

	//buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,                                   //buat jwt dr ID user
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), //masa expired
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
		"token": tokenString,
	})

}

func GetCurrentUser(c *gin.Context) {
	//ambil userId
	userId, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})

		return
	}

	var user models.User

	userData := config.DB.Select("id", "name", "email").First(&user, userId).Error
	if userData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"ID":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"events": user.Events,
		},
	})

}
