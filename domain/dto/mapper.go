package dto

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain"
)

func SendSuccess(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(ResponseJSON{
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

	return c.Status(statusCode).JSON(ResponseJSON{
		Status:  false,
		Message: message,
		Error:   details,
	})
}

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role.Name,
		JoinedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
