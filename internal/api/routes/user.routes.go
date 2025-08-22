package routes

import "github.com/gofiber/fiber/v2"

func UserRoutes(app fiber.Router) {
	users := app.Group("/users")

	users.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "List of users",
		})
	})

	users.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON(fiber.Map{
			"message": "User details",
			"id":      id,
		})
	})
}
