package dto

import (
	"github.com/gofiber/fiber/v2"
)

type responseJSON struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

func SendSuccess(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(responseJSON{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func SendError(c *fiber.Ctx, statusCode int, message string, err ...interface{}) error {
	var details interface{}
	if len(err) > 0 {
		details = err[0]
	}

	return c.Status(statusCode).JSON(responseJSON{
		Status:  false,
		Message: message,
		Error:   details,
	})
}
