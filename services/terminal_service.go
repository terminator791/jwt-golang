package services

import (
	"errors"

	"github.com/terminator791/jwt-golang/config"
	"github.com/terminator791/jwt-golang/models"
)

// TerminalService - Interface untuk service terminal
type TerminalService interface {
	CreateTerminal(terminal models.Terminal) (*models.Terminal, error)
}

// terminalService - Implementasi TerminalService
type terminalService struct{}

// NewTerminalService - Membuat instance baru TerminalService
func NewTerminalService() TerminalService {
	return &terminalService{}
}

// CreateTerminal - Metode untuk membuat terminal baru
func (s *terminalService) CreateTerminal(terminal models.Terminal) (*models.Terminal, error) {
	db := config.GetDB()

	// Cek apakah terminal code sudah digunakan
	var existingTerminal models.Terminal
	if err := db.Where("terminal_code = ?", terminal.TerminalCode).First(&existingTerminal).Error; err == nil {
		return nil, errors.New("kode terminal sudah digunakan")
	}

	// Validasi data terminal
	if terminal.TerminalName == "" {
		return nil, errors.New("nama terminal tidak boleh kosong")
	}

	if terminal.TerminalCode == "" {
		return nil, errors.New("kode terminal tidak boleh kosong")
	}

	// Buat terminal baru
	if err := db.Create(&terminal).Error; err != nil {
		return nil, err
	}

	return &terminal, nil
}
