package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/terminator791/jwt-golang/models"
)

// JWTClaim - Payload JWT
type JWTClaim struct {
	UserID   uuid.UUID       `json:"user_id"`
	Email    string          `json:"email"`
	UserType models.UserType `json:"user_type"`
	jwt.StandardClaims
}

// GenerateJWT - Membuat token JWT baru
func GenerateJWT(user models.User) (string, error) {
	// Set expiration time (24 jam)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Buat payload JWT
	claims := &JWTClaim{
		UserID:   user.UserID,
		Email:    user.Email,
		UserType: user.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Buat token dengan signing method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ambil secret key dari environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "default_secret_key_for_development" // Default key untuk development
	}

	// Tanda tangani token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT - Validasi token JWT
func ValidateJWT(tokenString string) (*JWTClaim, error) {
	// Ambil secret key dari environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "default_secret_key_for_development" // Default key untuk development
	}

	// Parse token JWT
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Validasi token dan extract claims
	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	// Cek apakah token sudah expired
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token sudah expired")
	}

	return claims, nil
}

// BlacklistToken - Masukkan token ke dalam blacklist
func BlacklistToken(tokenString string) error {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return err
	}

	// Ambil waktu expiry dari claims
	expiryTime := time.Unix(claims.ExpiresAt, 0)

	// Tambahkan token ke blacklist
	tokenBlacklist := GetTokenBlacklist()
	tokenBlacklist.AddToBlacklist(tokenString, expiryTime)

	return nil
}
