package config

import (
	"backend/eventApp/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	// host := os.Getenv("DB_HOST")
	// user := os.Getenv("DB_USER")
	// dbName := os.Getenv("DB_NAME")
	// dbPass := os.Getenv("DB_PASSWORD")
	// port := os.Getenv("DB_PORT")

	// if host == "" || user == "" || dbName == "" || port == "" {
	// 	log.Fatal("Konfigurasi database belum lengkap")
	// }

	dsn := "postgresql://neondb_owner:npg_mUu9SNEMeFW0@ep-bold-smoke-ao71bs75-pooler.c-2.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Makassar", host, user, dbPass, dbName, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("gagal terhubung ke database:", err)
	}

	// sqlDB, err := database.DB()
	// if err != nil {
	// 	log.Fatal("gagal terkoneksi ...")
	// }
	// // Memastikan database benar-benar bisa diakses
	// if err := sqlDB.Ping(); err != nil {
	// 	log.Fatal("Database tidak dapat diakses:", err)
	// }

	//lakukan automigrate dari model event
	err = database.AutoMigrate(&models.Event{}, &models.User{}, &models.Booking{})
	if err != nil {
		log.Fatal("Gagal melakukan migrate database", err)
	}

	DB = database
	log.Println("Berhasil tekoneksi ke database! :)")

}
