package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/middleware"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

func SetupRoutes(
	app *fiber.App,
	userHandler *UserHandler,
	tokenSvc usecase.TokenService,
) {
	app.Use(logger.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Grup API v1
	api := app.Group("/api/v1")

	// Rute Autentikasi
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// Rute yang dilindungi
	protected := api.Group("/users")
	protected.Use(middleware.AuthMiddleware(tokenSvc))
	protected.Get("/me", userHandler.GetMe)
}
