package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/routes"
	"github.com/ipincamp/edsa/internal/config"
)

func main() {
	env := config.LoadEnv()
	app := fiber.New()

	routes.SetupRoutes(app)

	port := env.AppPort
	if port == "" {
		port = "3000"
	}
	serverAddr := fmt.Sprintf(":%s", port)

	log.Printf("Server is starting on port %s", port)
	log.Fatal(app.Listen(serverAddr))
}
