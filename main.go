package main

import (
	"fmt"
	"log"
	"os"

	"github.com/IqbalBPH/golang-e-ticketing/config"
	"github.com/IqbalBPH/golang-e-ticketing/models"
	"github.com/IqbalBPH/golang-e-ticketing/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	// Inisialisasi database
	config.InitDB()
	db := config.GetDB()

	// Seed admin user untuk testing jika belum ada
	seedAdminUser(db)

	// Setup Gin
	r := gin.Default()

	// Middleware global
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(r)

	// Ambil port dari environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	// Jalankan server
	log.Printf("Server running on port %s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// seedAdminUser - Membuat admin user jika belum ada
func seedAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)

	// Jika belum ada user, buat admin default
	if count == 0 {
		// Hash password untuk admin
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password for admin user")
		}

		admin := models.User{
			FullName: "Admin System",
			Email:    "admin@e-ticketing.com",
			Password: string(hashedPassword),
			Phone:    "08123456789",
			UserType: models.UserTypeAdmin,
		}

		result := db.Create(&admin)
		if result.Error != nil {
			log.Printf("Failed to seed admin user: %v", result.Error)
		} else {
			log.Println("Admin user created successfully")
		}
	}
}