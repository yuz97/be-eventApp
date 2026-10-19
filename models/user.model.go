package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email" gorm:"uniqueIndex"`
	Password string `json:"password" binding:"required,min=6"`
	//relasi One to Many (1 user bisa memiliki banyak event)
	Events []Event `gorm:"foreignKey:UserID" json:"-"`
}
