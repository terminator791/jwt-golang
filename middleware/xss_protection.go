package middleware

import (
	"bytes"
	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// XSSProtection adalah middleware untuk mencegah XSS attack
func XSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Hanya proses request dengan content-type application/json
		if c.ContentType() == "application/json" {
			// Baca body
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
				c.Abort()
				return
			}

			// Restore body untuk middleware berikutnya
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Convert body ke string
			bodyString := string(bodyBytes)

			// Cek apakah ada script tag atau javascript event handlers
			scriptPattern := regexp.MustCompile(`(?i)<script|javascript:|on\w+\s*=|alert\s*\(`)
			if scriptPattern.MatchString(bodyString) {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Input mengandung konten berbahaya yang tidak diizinkan",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
