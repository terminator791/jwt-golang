package routes

import (
	"github.com/terminator791/jwt-golang/controllers"
	"github.com/terminator791/jwt-golang/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes - Setup semua route aplikasi
func SetupRoutes(r *gin.Engine) {
	// Instansiasi controller
	authController := controllers.NewAuthController()
	terminalController := controllers.NewTerminalController()

	// API Group
	api := r.Group("/api")
	{
		// Auth Routes - Tidak memerlukan autentikasi
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
		}

		// Terminal Routes - Memerlukan autentikasi
		terminal := api.Group("/terminal")
		terminal.Use(middleware.AuthMiddleware()) // Middleware JWT untuk semua route di bawah
		{
			terminal.POST("", terminalController.CreateTerminal)
		}
	}
}