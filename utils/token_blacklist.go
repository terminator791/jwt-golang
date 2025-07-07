package utils

import (
	"sync"
	"time"
)

// TokenBlacklist - Struct untuk menyimpan token yang sudah logout
type TokenBlacklist struct {
	blacklistedTokens map[string]time.Time
	mutex             sync.RWMutex
	cleanupInterval   time.Duration
}

// Global instance dari TokenBlacklist
var tokenBlacklist *TokenBlacklist
var once sync.Once

// GetTokenBlacklist - Mendapatkan instance singleton dari TokenBlacklist
func GetTokenBlacklist() *TokenBlacklist {
	once.Do(func() {
		tokenBlacklist = &TokenBlacklist{
			blacklistedTokens: make(map[string]time.Time),
			cleanupInterval:   1 * time.Hour,
		}
		go tokenBlacklist.periodicCleanup()
	})
	return tokenBlacklist
}

// AddToBlacklist - Menambahkan token ke blacklist
func (tb *TokenBlacklist) AddToBlacklist(token string, expiry time.Time) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.blacklistedTokens[token] = expiry
}

// IsBlacklisted - Memeriksa apakah token masuk dalam blacklist
func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	tb.mutex.RLock()
	defer tb.mutex.RUnlock()
	_, exists := tb.blacklistedTokens[token]
	return exists
}

// periodicCleanup - Membersihkan token yang sudah expired secara periodik
func (tb *TokenBlacklist) periodicCleanup() {
	ticker := time.NewTicker(tb.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		tb.cleanup()
	}
}

// cleanup - Menghapus token yang sudah expired dari blacklist
func (tb *TokenBlacklist) cleanup() {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	for token, expiry := range tb.blacklistedTokens {
		if now.After(expiry) {
			delete(tb.blacklistedTokens, token)
		}
	}
}