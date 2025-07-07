package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/terminator791/jwt-golang/controllers"
	"github.com/terminator791/jwt-golang/middleware"
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
		auth.Use(middleware.RateLimitAuth())
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
		}

		// Routes yang memerlukan autentikasi
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/user/profile", authController.GetUserProfile)

			// Logout
			protected.POST("/auth/logout", authController.Logout)

			// Terminal routes
			terminal := protected.Group("/terminal")
			{
				terminal.POST("/create", terminalController.CreateTerminal)
			}
		}
	}
}
