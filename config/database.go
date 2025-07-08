package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/terminator791/jwt-golang/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Inisialisasi koneksi database
func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Membuat string koneksi DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Membuka koneksi ke database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	// Drop
	// err = db.Exec("DROP TABLE IF EXISTS sync_logs, card_balance_logs, top_ups, transactions, gates, cards, fare_matrices, terminals, users CASCADE").Error
	// if err != nil {
	// 	log.Printf("Warning: Failed to drop existing tables: %v", err)
	// }

	// Migrasi semua model tanpa foreign key constraints
	log.Println("Creating tables without foreign key constraints...")
	err = db.Set("gorm:association_autocreate", false).Set("gorm:association_autoupdate", false).AutoMigrate(
		&models.User{},
		&models.Terminal{},
		&models.FareMatrix{},
		&models.Card{},
		&models.Gate{},
		&models.Transaction{},
		&models.TopUp{},
		&models.CardBalanceLog{},
		&models.SyncLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}

	// Setup relasi antar tabel
	log.Println("Creating foreign key constraints...")
	err = models.CreateForeignKeys(db)
	if err != nil {
		log.Printf("Warning: Failed to create some foreign keys: %v", err)
	}

	log.Println("Database migration completed successfully")

	DB = db
}

// GetDB
func GetDB() *gorm.DB {
	return DB
}
