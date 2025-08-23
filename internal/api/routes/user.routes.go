package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/response"
)

func UserRoutes(app fiber.Router) {
	users := app.Group("/users")

	users.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "List of users", nil)
	})

	users.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return response.Success(c, fiber.StatusOK, "User details", fiber.Map{
			"id": id,
		})
	})
}
