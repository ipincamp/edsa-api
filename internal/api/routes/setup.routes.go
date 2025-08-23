package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/handlers"
	"github.com/ipincamp/edsa/internal/api/response"
	"github.com/ipincamp/edsa/internal/config"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, env *config.Env) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello EDSA API!")
	})

	api := app.Group("/api")

	AuthRoutes(api, authHandler)
	UserRoutes(api, userHandler, env)

	app.Use(func(c *fiber.Ctx) error {
		return response.Error(c, fiber.StatusNotFound, "The requested resource was not found on this server")
	})
}
