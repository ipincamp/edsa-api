package util

import "github.com/gofiber/fiber/v2"

// ResponseJSON adalah struktur standar untuk response API
type ResponseJSON struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SendSuccess mengirim response sukses dengan data
func SendSuccess(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(ResponseJSON{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// SendError mengirim response error dengan detail error opsional
func SendError(c *fiber.Ctx, statusCode int, message string, err ...interface{}) error {
	var details interface{}
	if len(err) > 0 {
		details = err[0]
	}
	return c.Status(statusCode).JSON(ResponseJSON{
		Status:  false,
		Message: message,
		Error:   details,
	})
}
