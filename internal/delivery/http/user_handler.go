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
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// Panggil Usecase
	authResponse, err := h.userService.Register(c.Context(), &req)
	if err != nil {
		// Seharusnya ada error handling yang lebih baik di sini
		return utils.SendError(c, fiber.StatusConflict, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, authResponse)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest

	// Parse & Validasi
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// Panggil Usecase
	authResponse, err := h.userService.Login(c.Context(), &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, authResponse)
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var req domain.RefreshTokenRequest

	// Parse & Validasi
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// Panggil Usecase
	tokenResponse, err := h.userService.RefreshToken(c.Context(), &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, tokenResponse)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// Ambil user ID dari middleware
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid token")
	}

	// Panggil Usecase
	if err := h.userService.Logout(c.Context(), userID); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, fiber.Map{"message": "Logged out successfully"})
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	// Ambil user ID dari middleware
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid token")
	}

	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, user)
}
