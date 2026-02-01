package main

import (
	"log"

	"auth-go/internal/config"
	"auth-go/internal/database"
	"auth-go/internal/handler"
	"auth-go/internal/middleware"
	"auth-go/internal/repository"
	"auth-go/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Connect Database
	db := database.ConnectDB(cfg)

	// 3. Init Repositories
	userRepo := repository.NewUserRepository(db)
	resetRepo := repository.NewPasswordResetRepository(db)

	// 4. Init Services
	emailService := service.NewEmailService(cfg)
	authService := service.NewAuthService(userRepo, resetRepo, emailService, cfg)
	userService := service.NewUserService(userRepo)

	// 5. Init Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// 6. Init Router
	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 7. Setup Middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())

	// 8. Define Routes
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)

			// Protected Auth Route (e.g., Get Current User)
			auth.GET("/me", middleware.AuthMiddleware(cfg), userHandler.GetProfile)
		}

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg))
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/profile", userHandler.GetProfile)
		}
	}

	// 9. Start Server
	log.Printf("Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
