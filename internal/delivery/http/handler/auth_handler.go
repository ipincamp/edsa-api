package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/dto"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/service"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"github.com/ipincamp/go-edsa-api/pkg/validator"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *validator.CustomValidator
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.NewValidator(),
	}
}

// Register - User registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	registeredUser, err := h.authService.Register(user)
	if err != nil {
		status, msg := mapRegisterError(err)
		return util.SendError(c, status, msg, nil)
	}

	res := &dto.AuthResponse{
		User: dto.ToUserResponse(*registeredUser.User),
		Token: dto.Token{
			Access:  registeredUser.Token.Access,
			Refresh: registeredUser.Token.Refresh,
		},
	}
	return util.SendSuccess(c, fiber.StatusCreated, "user registered successfully", res)
}

// mapRegisterError memetakan error pada proses register ke status dan pesan yang sesuai
func mapRegisterError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrEmailExists):
		return fiber.StatusConflict, err.Error()
	case errors.Is(err, service.ErrDefaultRoleNotFound):
		fallthrough
	case errors.Is(err, service.ErrGenerateAccessToken):
		fallthrough
	case errors.Is(err, service.ErrGenerateRefreshToken):
		fallthrough
	case errors.Is(err, service.ErrHashingPassword):
		fallthrough
	case errors.Is(err, service.ErrUserCreation):
		return fiber.StatusInternalServerError, err.Error()
	default:
		return fiber.StatusInternalServerError, err.Error()
	}
}

// Login - User login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	loggedUser, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		status, msg := mapLoginError(err)
		return util.SendError(c, status, msg, nil)
	}

	res := &dto.AuthResponse{
		User: dto.ToUserResponse(*loggedUser.User),
		Token: dto.Token{
			Access:  loggedUser.Token.Access,
			Refresh: loggedUser.Token.Refresh,
		},
	}
	return util.SendSuccess(c, fiber.StatusOK, "login successful", res)
}

// mapLoginError memetakan error pada proses login ke status dan pesan yang sesuai
func mapLoginError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		fallthrough
	case errors.Is(err, service.ErrInvalidCredentials):
		return fiber.StatusUnauthorized, "invalid credentials"
	case errors.Is(err, service.ErrCheckCredentials):
		fallthrough
	case errors.Is(err, service.ErrGenerateAccessToken):
		fallthrough
	case errors.Is(err, service.ErrGenerateRefreshToken):
		return fiber.StatusInternalServerError, err.Error()
	default:
		return fiber.StatusUnauthorized, err.Error()
	}
}

// RefreshToken - Refresh access token using refresh token
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	token, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		status, msg := mapRefreshTokenError(err)
		return util.SendError(c, status, msg, nil)
	}

	res := &dto.Token{
		Access:  token.Access,
		Refresh: token.Refresh,
	}
	return util.SendSuccess(c, fiber.StatusOK, "token refreshed successfully", res)
}

// mapRefreshTokenError memetakan error pada proses refresh token ke status dan pesan yang sesuai
func mapRefreshTokenError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrInvalidToken):
		return fiber.StatusUnauthorized, err.Error()
	case errors.Is(err, service.ErrGenerateAccessToken):
		fallthrough
	case errors.Is(err, service.ErrGenerateRefreshToken):
		return fiber.StatusInternalServerError, err.Error()
	default:
		return fiber.StatusInternalServerError, err.Error()
	}
}

// Logout - User logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	if user == nil {
		return util.SendError(c, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	// logout
	if err := h.authService.Logout(user.ID); err != nil {
		return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// response
	return util.SendSuccess(c, fiber.StatusOK, "logged out successfully", nil)
}
