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
	// validation - body
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	// registration
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	registeredUser, err := h.authService.Register(user)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			return util.SendError(c, fiber.StatusConflict, err.Error(), nil)
		}
		if errors.Is(err, service.ErrDefaultRoleNotFound) {
			return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		if errors.Is(err, service.ErrGenerateAccessToken) {
			return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		if errors.Is(err, service.ErrGenerateRefreshToken) {
			return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		if errors.Is(err, service.ErrHashingPassword) {
			return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		if errors.Is(err, service.ErrUserCreation) {
			return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// response
	res := &dto.AuthResponse{
		User: dto.ToUserResponse(*registeredUser.User),
		Token: dto.Token{
			Access:  registeredUser.Token.Access,
			Refresh: registeredUser.Token.Refresh,
		},
	}
	return util.SendSuccess(c, fiber.StatusCreated, "user registered successfully", res)
}

// Login - User login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// validation - body
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	// login
	loggedUser, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return util.SendError(c, fiber.StatusUnauthorized, "invalid credentials", nil)
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			return util.SendError(c, fiber.StatusUnauthorized, "invalid credentials", nil)
		}
		if errors.Is(err, service.ErrCheckCredentials) {
			return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		if errors.Is(err, service.ErrGenerateAccessToken) {
			return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		if errors.Is(err, service.ErrGenerateRefreshToken) {
			return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		return util.SendError(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	// response
	res := &dto.AuthResponse{
		User: dto.ToUserResponse(*loggedUser.User),
		Token: dto.Token{
			Access:  loggedUser.Token.Access,
			Refresh: loggedUser.Token.Refresh,
		},
	}
	return util.SendSuccess(c, fiber.StatusOK, "login successful", res)
}

// RefreshToken - Refresh access token using refresh token
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// validation - body
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	// refresh token
	token, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			return util.SendError(c, fiber.StatusUnauthorized, err.Error(), nil)
		}
		if errors.Is(err, service.ErrGenerateAccessToken) {
			return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		if errors.Is(err, service.ErrGenerateRefreshToken) {
			return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		return util.SendError(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// response
	res := &dto.Token{
		Access:  token.Access,
		Refresh: token.Refresh,
	}
	return util.SendSuccess(c, fiber.StatusOK, "token refreshed successfully", res)
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
