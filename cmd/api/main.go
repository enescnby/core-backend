package main

import (
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
	keyService := services.NewKeyService(keyRepo)
	keyHandler := handlers.NewKeyHandler(keyService)

	app.Get("/keys/:id", keyHandler.GetPublicKey)

	logger.Log.Info("CoreGuard API başlatılıyor...")
	log.Fatal(app.Listen(":8080"))
}
