package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/terminator791/jwt-golang/utils"
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

		// Ambil token
		tokenString := parts[1]

		// Cek apakah token di blacklist
		tokenBlacklist := utils.GetTokenBlacklist()
		if tokenBlacklist.IsBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid atau sudah logout",
			})
			c.Abort()
			return
		}

		// Validasi token
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

		// Simpan token string untuk digunakan pada saat logout
		c.Set("tokenString", tokenString)

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
