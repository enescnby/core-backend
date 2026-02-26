package main

import (
	"core-backend/internal/middleware"
	"log"

	"github.com/gofiber/fiber/v2"

	"core-backend/internal/config"
	"core-backend/internal/database"
	"core-backend/internal/handlers"
	"core-backend/internal/repositories"
	"core-backend/internal/services"
	"core-backend/pkg/logger"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()

	config.LoadConfig()

	database.Connect()
	defer database.Close()

	database.Migrate()

	app := fiber.New()

	keyRepo := repositories.NewKeyRepository(database.DB)
	userRepo := repositories.NewUserRepository(database.DB)
	auditRepo := repositories.NewAuditRepository(database.DB)

	authService := services.NewAuthService(userRepo, auditRepo)
	keyService := services.NewKeyService(keyRepo)

	authHandler := handlers.NewAuthHandler(authService)
	keyHandler := handlers.NewKeyHandler(keyService)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	authGroup := v1.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login/init", authHandler.LoginInit)
	authGroup.Post("/login/verify", authHandler.LoginVerify)

	keys := v1.Group("/keys", middleware.Protected())
	keys.Get("/:id", keyHandler.GetPublicKey)

	logger.Log.Info("CoreGuard API is starting...")
	log.Fatal(app.Listen(config.AppConfig.AppPort))
}
