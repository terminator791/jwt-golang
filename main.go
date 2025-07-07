package main

import (
	"fmt"
	"log"
	"os"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/terminator791/jwt-golang/config"
	"github.com/terminator791/jwt-golang/models"
	"github.com/terminator791/jwt-golang/routes"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		log.Println("Database kosong. Membuat akun admin default...")

		// Hash password untuk admin
		password := "admin123" // Password default admin
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Gagal melakukan hash password untuk akun admin")
		}

		// Buat admin dengan data lengkap
		admin := models.User{
			UserID:      uuid.New(),
			FullName:    "Admin Sistem E-Ticketing",
			Email:       "admin@e-ticketing.com",
			Password:    string(hashedPassword),
			Phone:       "081234567890",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), // Tanggal lahir dummy
			UserType:    models.UserTypeAdmin,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		result := db.Create(&admin)
		if result.Error != nil {
			log.Printf("Gagal membuat akun admin default: %v", result.Error)
		} else {
			log.Println("✅ Akun admin berhasil dibuat dengan detail:")
			log.Println("   Email   : admin@e-ticketing.com")
			log.Println("   Password: admin123")
		}
	} else {
		log.Println("Database sudah berisi data. Melewati pembuatan akun admin default.")
	}
}
