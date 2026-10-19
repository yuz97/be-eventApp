package main

import (
	"backend/eventApp/config"
	"backend/eventApp/controllers"
	"backend/eventApp/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// router.GET("/", getEvents)
	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	//atur cors
	router.Use(cors.Default())
	//koneksi database
	config.ConnectDB()

	// Route
	api := router.Group("/api")
	{

		// api.GET("/events", getEvents)
		// api.POST("/events", createEvent)

		//auth
		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		//event
		api.GET("/events", controllers.GetAllEvents)
		api.GET("/events/:id", controllers.GetEventById)

		//protected middleware
		protected := api.Group("/")
		protected.Use(middleware.RequiredAuth())
		{
			protected.GET("/auth/me", controllers.GetCurrentUser)

			protected.POST("/events", controllers.CreateEvent)
			protected.PUT("/events/:id", controllers.UpdateEventById)
			protected.DELETE("/events/:id", controllers.DeleteEventById)

			protected.POST("/booking", controllers.CreateBookingEvent)
			protected.GET("/booking/user", controllers.GetBookingByUser)
			protected.DELETE("/booking/:id", controllers.DeleteBooking)

		}

	}

	router.Run()
}

// function handler
// func getEvents(c *gin.Context) {

// 	events := models.GetAllEvents()

// 	c.JSON(200, gin.H{
// 		"message": "wellcome my brother",
// 		"event` ": events,
// 	})
// }

// func createEvent(c *gin.Context) {
// 	var event models.Event
// 	err := c.ShouldBindJSON(&event)
// 	if err != nil {
// 		c.JSON(500, gin.H{
// 			"message": "Could not parse request data",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	//ambil waktu sekarang
// 	// tm := time.Now()

// 	// event.Id = 1
// 	// event.UserId = 1
// 	// event.DateTime = tm

// 	//save inputan
// 	event.Save()

// 	c.JSON(200, gin.H{
// 		"message": "Data event added success",
// 		"event":   event,
// 	})
// }
