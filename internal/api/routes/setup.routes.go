package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/response"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello EDSA!")
	})

	api := app.Group("/api")
	UserRoutes(api)

	app.Use(func(c *fiber.Ctx) error {
		return response.Error(c, fiber.StatusNotFound, "Not Found")
	})
}
