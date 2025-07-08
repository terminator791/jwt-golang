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

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.Use(middleware.RateLimitAuth())
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
		}

		// Routes yang memerlukan autentikasi (middleware)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/user/profile", authController.GetUserProfile)

			protected.POST("/auth/logout", authController.Logout)

			terminal := protected.Group("/terminal")
			{
				terminal.POST("/create", terminalController.CreateTerminal)
			}
		}
	}
}
