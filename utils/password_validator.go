package utils

import (
	"errors"
	"regexp"
	"unicode"
)

// PasswordStrengthValidator memvalidasi kekuatan password
func ValidatePasswordStrength(password string) error {
	// Cek panjang minimum (sudah ada di binding, tapi ditambahkan lagi untuk lengkapnya)
	if len(password) < 6 {
		return errors.New("password harus minimal 6 karakter")
	}

	// Cek apakah ada huruf kecil
	hasLowerCase := false
	// Cek apakah ada angka
	hasNumber := false

	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowerCase = true
		} else if unicode.IsDigit(char) {
			hasNumber = true
		}
	}

	if !hasLowerCase {
		return errors.New("password harus mengandung minimal 1 huruf kecil")
	}

	if !hasNumber {
		return errors.New("password harus mengandung minimal 1 angka")
	}

	return nil
}

// SanitizeInput membersihkan input dari potensi XSS
func SanitizeInput(input string) string {
	// Menghapus tag HTML/JavaScript yang berbahaya
	// Ini adalah implementasi sederhana, untuk produksi gunakan library yang lebih robust
	stripTagsRegex := regexp.MustCompile("<[^>]*>")
	return stripTagsRegex.ReplaceAllString(input, "")
}

// SanitizeUserInput membersihkan berbagai field input user dari XSS
func SanitizeUserInput(fullName, email, phone string) (string, string, string) {
	return SanitizeInput(fullName),
		SanitizeInput(email),
		SanitizeInput(phone)
}
