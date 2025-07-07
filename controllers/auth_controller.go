package controllers

import (
	"net/http"

	"github.com/terminator791/jwt-golang/services"
	"github.com/gin-gonic/gin"
)

// AuthController - Interface untuk controller autentikasi
type AuthController interface {
	Login(c *gin.Context)
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