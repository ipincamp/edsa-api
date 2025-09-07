package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/router"
	"github.com/ipincamp/go-edsa-api/pkg/cache"
	"github.com/ipincamp/go-edsa-api/pkg/database"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Load cache
	cache.LoadCache(db)

	// Create fiber app
	app := fiber.New()

	// Setup router
	router.Setup(app, db)

	// Start server
	log.Fatal(app.Listen(cfg.Server.Host + ":" + cfg.Server.Port))
}
