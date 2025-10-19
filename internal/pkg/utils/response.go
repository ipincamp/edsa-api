package utils

import "github.com/gofiber/fiber/v2"

func SendSuccess(c *fiber.Ctx, code int, data interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func SendError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}

func SendValidationErrors(c *fiber.Ctx, errors interface{}) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"status":  "error",
		"message": "Validation failed",
		"errors":  errors,
	})
}
