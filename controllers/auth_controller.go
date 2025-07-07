package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/terminator791/jwt-golang/models"
	"github.com/terminator791/jwt-golang/services"
	"github.com/terminator791/jwt-golang/utils"
)

// AuthController - Interface untuk controller autentikasi
type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	GetUserProfile(c *gin.Context)
	Logout(c *gin.Context)
}

// authController - Implementasi AuthController
type authController struct {
	authService services.AuthService
}

// NewAuthController - Membuat instance baru AuthController
func NewAuthController() AuthController {
	return &authController{
		authService: services.NewAuthService(),
	}
}

// LoginRequest - Struktur request untuk login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest - Struktur request untuk register
type RegisterRequest struct {
	FullName    string    `json:"full_name" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	Password    string    `json:"password" binding:"required,min=6"`
	Phone       string    `json:"phone"`
	DateOfBirth time.Time `json:"date_of_birth"`
	UserType    string    `json:"user_type"`
}

// Login - Handler untuk endpoint login
func (ctrl *authController) Login(c *gin.Context) {
	var request LoginRequest

	// Binding request body ke struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format request tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	// Panggil service untuk login
	token, user, err := ctrl.authService.Login(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Gagal login",
			"error":   err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login berhasil",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"user_id":   user.UserID,
				"full_name": user.FullName,
				"email":     user.Email,
				"user_type": user.UserType,
			},
		},
	})
}

// Register - Handler untuk endpoint register
func (ctrl *authController) Register(c *gin.Context) {
	var request RegisterRequest

	// Binding request body ke struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format request tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	// Konversi request ke model User
	user := models.User{
		FullName:    request.FullName,
		Email:       request.Email,
		Phone:       request.Phone,
		DateOfBirth: request.DateOfBirth,
	}

	// Konversi string user_type ke UserType jika ada
	if request.UserType != "" {
		switch request.UserType {
		case string(models.UserTypeAdmin):
			user.UserType = models.UserTypeAdmin
		case string(models.UserTypeStaff):
			user.UserType = models.UserTypeStaff
		default:
			user.UserType = models.UserTypeCustomer
		}
	} else {
		user.UserType = models.UserTypeCustomer
	}

	// Panggil service untuk register
	registeredUser, err := ctrl.authService.Register(user, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Gagal registrasi",
			"error":   err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data": gin.H{
			"user_id":   registeredUser.UserID,
			"full_name": registeredUser.FullName,
			"email":     registeredUser.Email,
			"phone":     registeredUser.Phone,
			"user_type": registeredUser.UserType,
		},
	})
}

// GetUserProfile - Handler untuk mendapatkan profil user
func (ctrl *authController) GetUserProfile(c *gin.Context) {
	// Dapatkan user_id dari context (set oleh middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID tidak ditemukan di context",
		})
		return
	}

	// Konversi ke UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Tipe data user ID tidak valid",
		})
		return
	}

	// Dapatkan data user dari database
	user, err := ctrl.authService.GetUserByID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Gagal mendapatkan profil user",
			"error":   err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profil user berhasil didapatkan",
		"data": gin.H{
			"user_id":       user.UserID,
			"full_name":     user.FullName,
			"email":         user.Email,
			"phone":         user.Phone,
			"date_of_birth": user.DateOfBirth,
			"user_type":     user.UserType,
			"created_at":    user.CreatedAt,
		},
	})
}

// Logout - Handler untuk logout (minimal implementation)
// Logout - Handler untuk logout dengan blacklist token
func (ctrl *authController) Logout(c *gin.Context) {
	// Ambil token dari context
	tokenString, exists := c.Get("tokenString")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Token tidak ditemukan",
		})
		return
	}

	// Masukkan token ke blacklist
	err := utils.BlacklistToken(tokenString.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Gagal logout",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout berhasil",
	})
}
