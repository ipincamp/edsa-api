package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type UserHandler struct {
	userService usecase.UserService
	validate    *validator.GoPlaygroundValidator
}

func NewUserHandler(us usecase.UserService, v *validator.GoPlaygroundValidator) *UserHandler {
	return &UserHandler{
		userService: us,
		validate:    v,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest

	// Parse & Validasi
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// Panggil Usecase
	authResponse, err := h.userService.Register(c.Context(), &req)
	if err != nil {
		// Seharusnya ada error handling yang lebih baik di sini
		return utils.SendSimpleError(c, fiber.StatusConflict, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "User registered successfully", authResponse)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest

	// Parse & Validasi
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// Panggil Usecase
	authResponse, err := h.userService.Login(c.Context(), &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Login successful", authResponse)
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var req domain.RefreshTokenRequest

	// Parse & Validasi
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// Panggil Usecase
	tokenResponse, err := h.userService.RefreshToken(c.Context(), &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Token refreshed successfully", tokenResponse)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// Ambil user ID dari middleware
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid user ID in token")
	}

	// Ambil session ID dari middleware
	sessionID, ok := c.Locals("sessionID").(uuid.UUID)
	if !ok {
		// Jika sessionID tidak ada, ini adalah error krusial untuk logout
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid session ID in token")
	}

	// Panggil Usecase
	if err := h.userService.Logout(c.Context(), userID, sessionID); err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Logged out successfully", nil)
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	// Ambil user ID dari middleware
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid user ID in token")
	}

	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User profile retrieved successfully", user)
}
