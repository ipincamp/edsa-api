package handler

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/util"
)

type AuthHandler struct {
	AuthService domain.AuthService
	Validator   *util.GoValidator
}

func NewAuth(authService domain.AuthService, validator *util.GoValidator) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Validator:   validator,
	}
}

func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	var request dto.RegisterRequest
	c, cancel, err := parseAndValidate(ctx, h.Validator, &request)
	if err != nil {
		return err
	}
	defer cancel()

	res, err := h.AuthService.Register(c, request)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return dto.SendError(ctx, fiber.StatusConflict, err.Error(), nil)
		}
		return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to register user", err)
	}

	return dto.SendSuccess(ctx, fiber.StatusCreated, "User registered successfully", res)
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var request dto.LoginRequest
	c, cancel, err := parseAndValidate(ctx, h.Validator, &request)
	if err != nil {
		return err
	}
	defer cancel()

	res, err := h.AuthService.Login(c, request)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			return dto.SendError(ctx, fiber.StatusUnauthorized, err.Error(), nil)
		}
		return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to log in", err)
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "Login successful", res)
}

func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	err := h.AuthService.Logout(c, ctx.Locals("user").(domain.User))
	if err != nil {
		return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to log out", err)
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "Logout successful", nil)
}
