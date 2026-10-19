package api

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
	once   sync.Once
	router *gin.Engine
)

// Inisialisasi router
func setupRouter() {
	router = gin.Default()
	router.Use(cors.Default())

	// Koneksi database
	config.ConnectDB()

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

	protected.GET("/auth/me", controllers.GetCurrentUser)

	protected.POST("/events", controllers.CreateEvent)
	protected.PUT("/events/:id", controllers.UpdateEventById)
	protected.DELETE("/events/:id", controllers.DeleteEventById)

	// Booking
	protected.POST("/booking", controllers.CreateBookingEvent)
	protected.GET("/booking/user", controllers.GetBookingByUser)
	protected.DELETE("/booking/:id", controllers.DeleteBooking)
}

// Digunakan oleh main.go
func SetupRouter() *gin.Engine {
	once.Do(setupRouter)
	return router
}

// Digunakan oleh Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	SetupRouter().ServeHTTP(w, r)
}
