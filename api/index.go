package handler

import (
	"net/http"
	"sync"

	"backend/eventApp/config"
	"backend/eventApp/controllers"
	"backend/eventApp/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	app  *gin.Engine
	once sync.Once
)

// Setup dipakai oleh main.go (lokal) dan Handler (Vercel)
func Setup() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://domain-frontend-kamu.vercel.app", "http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	api := router.Group("/api")

	// Auth
	api.POST("/auth/register", controllers.RegisterUser)
	api.POST("/auth/login", controllers.LoginUser)

	// Events
	api.GET("/events", controllers.GetAllEvents)
	api.GET("/events/:id", controllers.GetEventById)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.RequiredAuth())
	{
		protected.GET("/auth/me", controllers.GetCurrentUser)

		protected.POST("/events", controllers.CreateEvent)
		protected.PUT("/events/:id", controllers.UpdateEventById)
		protected.DELETE("/events/:id", controllers.DeleteEventById)

		// Booking
		protected.POST("/booking", controllers.CreateBookingEvent)
		protected.GET("/booking/user", controllers.GetBookingByUser)
		protected.DELETE("/booking/:id", controllers.DeleteBooking)
	}

	return router
}

// Handler adalah entry point untuk Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		gin.SetMode(gin.ReleaseMode)
		config.ConnectDB()
		app = Setup()
	})
	app.ServeHTTP(w, r)
}
