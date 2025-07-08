package services

import (
	"errors"

	"time"

	"github.com/google/uuid"
	"github.com/terminator791/jwt-golang/config"
	"github.com/terminator791/jwt-golang/models"
	"github.com/terminator791/jwt-golang/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService - Interface untuk service autentikasi
type AuthService interface {
	Login(email, password string) (string, *models.User, error)
	Register(user models.User, password string) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
}

type authService struct{}

// NewAuthService - Membuat instance baru AuthService
func NewAuthService() AuthService {
	return &authService{}
}

// Metode untuk login user
func (s *authService) Login(email, password string) (string, *models.User, error) {
	var user models.User
	db := config.GetDB()

	// Cari user berdasarkan email
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", nil, errors.New("pengguna tidak ditemukan")
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("password salah")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		return "", nil, errors.New("gagal membuat token autentikasi")
	}

	return token, &user, nil
}

// Metode untuk register user baru
func (s *authService) Register(user models.User, password string) (*models.User, error) {
	db := config.GetDB()

	// Validasi email unik
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email sudah terdaftar")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("terjadi kesalahan saat memeriksa email")
	}

	// Validasi nomor telepon unik
	if user.Phone != "" {
		if err := db.Where("phone = ?", user.Phone).First(&existingUser).Error; err == nil {
			return nil, errors.New("nomor telepon sudah terdaftar")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("terjadi kesalahan saat memeriksa nomor telepon")
		}
	}

	// Validasi data user
	if user.FullName == "" {
		return nil, errors.New("nama lengkap tidak boleh kosong")
	}

	if user.Email == "" {
		return nil, errors.New("email tidak boleh kosong")
	}

	// Set default user type jika tidak diisi
	if user.UserType == "" {
		user.UserType = models.UserTypeCustomer
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}
	user.Password = string(hashedPassword)

	// Set waktu pembuatan dan update
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Buat user baru
	if err := db.Create(&user).Error; err != nil {
		return nil, errors.New("gagal membuat user baru: " + err.Error())
	}

	// Jangan kembalikan password
	user.Password = ""

	return &user, nil
}

// GetUserByID - Metode untuk mendapatkan data user berdasarkan ID
func (s *authService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	db := config.GetDB()

	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengguna tidak ditemukan")
		}
		return nil, errors.New("terjadi kesalahan saat mencari pengguna")
	}

	// Jangan kembalikan password
	user.Password = ""

	return &user, nil
}
