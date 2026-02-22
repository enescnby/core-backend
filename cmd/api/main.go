package main

import (
	"core-backend/internal/config"
	"core-backend/internal/database"
	"core-backend/internal/handlers"
	"core-backend/internal/repositories"
	"core-backend/internal/services"
	"core-backend/pkg/logger"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	//Logger
	logger.InitLogger()
	defer logger.Log.Sync()

	config.LoadConfig()

	database.Connect()
	defer database.Close()

	database.Migrate()

	userRepo := repositories.NewUserRepository()
	auditRepo := repositories.NewAuditRepository()

	authService := services.NewAuthService(userRepo, auditRepo)

	authHandler := handlers.NewAuthHandler(authService)

	app := fiber.New(fiber.Config{
		AppName: "CoreGuard E2EE Backend v0.1.0",
	})

	api := app.Group("/api/v1")
	api.Post("/register", authHandler.Register)
	api.Post("/login/init", authHandler.LoginInit)
	api.Post("/login/verify", authHandler.LoginVerify)
	
	logger.Log.Info(fmt.Sprintf("starting CoreGuard server on port %s", config.AppConfig.AppPort))

	if err := app.Listen(config.AppConfig.AppPort); err != nil {
		logger.Log.Fatal("server failed to start", zap.Error(err))
	}
}
