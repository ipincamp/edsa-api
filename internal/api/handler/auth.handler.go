package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/service"
	"github.com/ipincamp/go-edsa-api/internal/util"
)

type AuthHandler struct {
	AuthService service.AuthService
	Validator   *util.GoValidator
}

func NewAuth(authService service.AuthService, validator *util.GoValidator) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Validator:   validator,
	}
}

func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	var request dto.RegisterRequest
	c, cancel, err := parseAndValidateBody(ctx, h.Validator, &request)
	if err != nil {
		return handleValidationError(ctx, err)
	}
	defer cancel()

	res, err := h.AuthService.Register(c, request)
	if err != nil {
		return handleServiceError(ctx, err, "Failed to register user")
	}

	return dto.SendSuccess(ctx, fiber.StatusCreated, "User registered successfully", res)
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var request dto.LoginRequest
	c, cancel, err := parseAndValidateBody(ctx, h.Validator, &request)
	if err != nil {
		return handleValidationError(ctx, err)
	}
	defer cancel()

	res, err := h.AuthService.Login(c, request)
	if err != nil {
		return handleServiceError(ctx, err, "Failed to log in")
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "Login successful", res)
}

func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	c, cancel := createContextWithTimeout(ctx)
	defer cancel()

	err := h.AuthService.Logout(c, ctx.Locals("user").(domain.User))
	if err != nil {
		return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to log out", err)
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "Logout successful", nil)
}
