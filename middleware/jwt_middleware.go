package middleware

import (
	"net/http"
	"strings"

	"github.com/terminator791/jwt-golang/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware - Middleware untuk validasi token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header tidak ditemukan",
			})
			c.Abort()
			return
		}

		// Cek format Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Format Authorization header tidak valid",
			})
			c.Abort()
			return
		}

		// Ambil token dan validasi
		tokenString := parts[1]
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		// Set user claims ke context untuk digunakan di handler berikutnya
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("userType", claims.UserType)
		c.Next()
	}
}

// AdminMiddleware - Middleware khusus untuk akses admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pastikan AuthMiddleware sudah dijalankan sebelumnya
		userType, exists := c.Get("userType")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Akses ditolak",
			})
			c.Abort()
			return
		}

		// Cek apakah user adalah admin
		if userType != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Anda tidak memiliki akses untuk fitur ini",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}