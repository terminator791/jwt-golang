package services

import (
	"errors"

	"github.com/terminator791/jwt-golang/config"
	"github.com/terminator791/jwt-golang/models"
	"github.com/terminator791/jwt-golang/utils"
	"golang.org/x/crypto/bcrypt"
)

// AuthService - Interface untuk service autentikasi
type AuthService interface {
	Login(email, password string) (string, *models.User, error)
}

// authService - Implementasi AuthService
type authService struct{}

// NewAuthService - Membuat instance baru AuthService
func NewAuthService() AuthService {
	return &authService{}
}

// Login - Metode untuk login user
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
