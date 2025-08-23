package response

import "github.com/gofiber/fiber/v2"

type SuccessResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(SuccessResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, statusCode int, message string, errDetails ...interface{}) error {
	var details interface{}
	if len(errDetails) > 0 {
		details = errDetails[0]
	}

	return c.Status(statusCode).JSON(ErrorResponse{
		Status:  false,
		Message: message,
		Errors:  details,
	})
}
