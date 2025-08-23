package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/dto"
	"github.com/ipincamp/edsa/internal/api/middleware"
	"github.com/ipincamp/edsa/internal/api/response"
	"github.com/ipincamp/edsa/internal/repositories"
	"github.com/ipincamp/edsa/internal/utils"
)

type UserHandler struct {
	userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	payload := c.Locals(middleware.AuthorizationPayloadKey).(*utils.PasetoPayload)

	user, err := h.userRepo.FindUserByID(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}

	formattedUser := dto.FormatUser(user)

	return response.Success(c, fiber.StatusOK, "Profile retrieved successfully", formattedUser)
}
