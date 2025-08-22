package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello EDSA!")
	})

	api := app.Group("/api")
	UserRoutes(api)
}
