package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	// Id          int otomatis generate dr db
	gorm.Model
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Image       string `json:"image"`
	ImageID     string `json:"imageId"`
	UserID      int    `json:"userId"`
	// User        User      `gorm:"foreignKey:UserID" json:"-"` "-" artinya sembunyikan field ini dalam json
	User     User      `gorm:"foreignKey:UserID" json:"user"`
	DateTime time.Time `json:"datetime"`
	Booking  []Booking `gorm:"foreignKey:EventID" json:"listBooking"`
}

//simpan di lokal bukan database
// var events []Event = []Event{}

// // function utk simpan event
// func (e Event) Save() {
// 	events = append(events, e)
// }

// // show all event
// func GetAllEvents() []Event {
// 	return events
// }
