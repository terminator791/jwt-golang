package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter mengatur rate limiting berdasarkan IP address
type RateLimiter struct {
	// ipMap menyimpan IP dan timestamp request terakhir
	ipMap map[string][]time.Time
	// mutex untuk thread safety
	mu sync.RWMutex
	// durasi window untuk rate limiting
	windowSize time.Duration
	// maksimum request dalam 1 window
	maxRequests int
}

// NewRateLimiter membuat instance RateLimiter baru
func NewRateLimiter(windowSize time.Duration, maxRequests int) *RateLimiter {
	limiter := &RateLimiter{
		ipMap:       make(map[string][]time.Time),
		windowSize:  windowSize,
		maxRequests: maxRequests,
	}

	// Membersihkan data lama secara periodik
	go func() {
		for {
			time.Sleep(time.Minute * 5) // bersihkan setiap 5 menit
			limiter.mu.Lock()
			now := time.Now()
			for ip, timestamps := range limiter.ipMap {
				var validTimestamps []time.Time
				for _, ts := range timestamps {
					if now.Sub(ts) < limiter.windowSize {
						validTimestamps = append(validTimestamps, ts)
					}
				}
				if len(validTimestamps) > 0 {
					limiter.ipMap[ip] = validTimestamps
				} else {
					delete(limiter.ipMap, ip)
				}
			}
			limiter.mu.Unlock()
		}
	}()

	return limiter
}

// RateLimitMiddleware menciptakan middleware Gin untuk rate limiting
func RateLimitMiddleware(windowSize time.Duration, maxRequests int) gin.HandlerFunc {
	limiter := NewRateLimiter(windowSize, maxRequests)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mu.Lock()
		now := time.Now()

		// Inisialisasi jika IP belum ada di map
		if _, exists := limiter.ipMap[ip]; !exists {
			limiter.ipMap[ip] = []time.Time{}
		}

		// Filter timestamps yang masih dalam window
		var validTimestamps []time.Time
		for _, ts := range limiter.ipMap[ip] {
			if now.Sub(ts) < limiter.windowSize {
				validTimestamps = append(validTimestamps, ts)
			}
		}

		// Cek apakah sudah melewati batas
		if len(validTimestamps) >= limiter.maxRequests {
			limiter.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":              "error",
				"message":             "Terlalu banyak permintaan. Coba lagi setelah beberapa saat.",
				"retry_after_seconds": int(limiter.windowSize.Seconds()),
			})
			c.Abort()
			return
		}

		// Tambahkan timestamp request saat ini
		validTimestamps = append(validTimestamps, now)
		limiter.ipMap[ip] = validTimestamps

		limiter.mu.Unlock()
		c.Next()
	}
}

// (5 request per menit)
func RateLimitAuth() gin.HandlerFunc {
	return RateLimitMiddleware(time.Minute, 5)
}

// (30 request per menit)
func RateLimitGeneral() gin.HandlerFunc {
	return RateLimitMiddleware(time.Minute, 30)
}
