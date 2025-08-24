package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/handlers"
	"github.com/ipincamp/edsa/internal/api/routes"
	"github.com/ipincamp/edsa/internal/api/validator"
	"github.com/ipincamp/edsa/internal/config"
	"github.com/ipincamp/edsa/internal/database"
	"github.com/ipincamp/edsa/internal/repositories"
	"github.com/ipincamp/edsa/internal/services"
)

func main() {
	env := config.LoadEnv()
	database.ConnectDB(env)

	customValidator := validator.New()

	userRepo := repositories.NewUserRepository(database.DB)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService, env, customValidator)
	userHandler := handlers.NewUserHandler(userRepo)

	app := fiber.New()
	routes.SetupRoutes(app, authHandler, userHandler, env)

	port := env.AppPort
	if port == "" {
		port = "3000"
	}
	serverAddr := fmt.Sprintf(":%s", port)

	log.Printf("Server is starting on port %s", port)
	log.Fatal(app.Listen(serverAddr))
}
