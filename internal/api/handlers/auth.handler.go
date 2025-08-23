package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/dto"
	"github.com/ipincamp/edsa/internal/api/response"
	"github.com/ipincamp/edsa/internal/api/validator"
	"github.com/ipincamp/edsa/internal/config"
	"github.com/ipincamp/edsa/internal/services"
	"github.com/ipincamp/edsa/internal/utils"
)

type AuthHandler struct {
	authService services.AuthService
	env         *config.Env
	validator   *validator.CustomValidator
}

func NewAuthHandler(authService services.AuthService, env *config.Env, validator *validator.CustomValidator) *AuthHandler {
	return &AuthHandler{authService, env, validator}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(&req); errs != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation failed", errs)
	}

	err := h.authService.Register(&req)
	if err != nil {
		return response.Error(c, fiber.StatusConflict, err.Error())
	}

	return response.Success(c, fiber.StatusAccepted, "Registration request accepted and is being processed", nil)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(&req); errs != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation failed", errs)
	}

	user, err := h.authService.Login(&req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	payload, err := utils.NewPayload(user, h.env.PasetoExpireInHours)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create token payload")
	}

	token, err := utils.CreateToken(payload, h.env.PasetoSymmetricKey)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create token")
	}

	return response.Success(c, fiber.StatusOK, "Login successful", fiber.Map{"token": token})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, "Logout successful", nil)
}
