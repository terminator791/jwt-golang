package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/terminator791/jwt-golang/models"
	"github.com/terminator791/jwt-golang/services"
)

// TerminalController - Interface untuk controller terminal
type TerminalController interface {
	CreateTerminal(c *gin.Context)
}

// terminalController - Implementasi TerminalController
type terminalController struct {
	terminalService services.TerminalService
}

// NewTerminalController - Membuat instance baru TerminalController
func NewTerminalController() TerminalController {
	return &terminalController{
		terminalService: services.NewTerminalService(),
	}
}

// CreateTerminalRequest - Struktur request untuk membuat terminal
type CreateTerminalRequest struct {
	TerminalName string  `json:"terminal_name" binding:"required"`
	TerminalCode string  `json:"terminal_code" binding:"required"`
	Location     string  `json:"location"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	IsActive     bool    `json:"is_active"`
}

// CreateTerminal - Handler untuk endpoint create terminal
func (ctrl *terminalController) CreateTerminal(c *gin.Context) {
	var request CreateTerminalRequest

	// Binding request body ke struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format request tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	// Konversi request ke model Terminal
	terminal := models.Terminal{
		TerminalName: request.TerminalName,
		TerminalCode: request.TerminalCode,
		Location:     request.Location,
		Latitude:     request.Latitude,
		Longitude:    request.Longitude,
		IsActive:     request.IsActive,
	}

	// Panggil service untuk membuat terminal
	result, err := ctrl.terminalService.CreateTerminal(terminal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Gagal membuat terminal",
			"error":   err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Terminal berhasil dibuat",
		"data":    result,
	})
}
