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

type UserHandler struct {
	userService service.UserService
	validator   *validator.CustomValidator
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.NewValidator(),
	}
}

// GetProfile - Get current user profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	if user == nil {
		return util.SendError(c, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	// profile
	profile, err := h.userService.GetProfile(user.ID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return util.SendError(c, fiber.StatusNotFound, "user not found", nil)
		}
		return util.SendError(c, fiber.StatusInternalServerError, "failed to get profile", nil)
	}

	// response
	return util.SendSuccess(c, fiber.StatusOK, "profile retrieved successfully", profile)
}

// UpdateProfile - Update current user profile (self)
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	if user == nil {
		return util.SendError(c, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	// validation - body
	var req dto.UpdateProfileUserRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	// update profile
	err := h.userService.UpdateProfile(user.ID, req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return util.SendError(c, fiber.StatusNotFound, "user not found", nil)
		}
		if errors.Is(err, service.ErrInvalidOldPassword) {
			return util.SendError(c, fiber.StatusBadRequest, "old password is incorrect", nil)
		}
		return util.SendError(c, fiber.StatusInternalServerError, "failed to update profile", nil)
	}

	// response
	return util.SendSuccess(c, fiber.StatusOK, "profile updated successfully", nil)
}

// GetAllUsers - Get all users with pagination (admin only)
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// validation - query
	var req dto.UserFilterRequest
	if err := c.QueryParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid query parameters", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// get users
	users, err := h.userService.GetAllUsers(req.Page, req.Limit, req.Role)
	if err != nil {
		return util.SendError(c, fiber.StatusInternalServerError, "failed to get users", nil)
	}

	// response
	return util.SendSuccess(c, fiber.StatusOK, "users retrieved successfully", users)
}

// GetUserByID - Get user by ID (admin only)
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	// validation - params
	var req dto.UserIDRequest
	if err := c.ParamsParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid user ID", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	// profile
	profile, err := h.userService.GetProfile(req.UserId)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return util.SendError(c, fiber.StatusNotFound, "user not found", nil)
		}
		return util.SendError(c, fiber.StatusInternalServerError, "failed to get user", nil)
	}

	// response
	return util.SendSuccess(c, fiber.StatusOK, "user retrieved successfully", profile)
}

// UpdateUserByID - Update user by ID (admin only)
func (h *UserHandler) UpdateUserByID(c *fiber.Ctx) error {
	// validation - params
	var userIDReq dto.UserIDRequest
	if err := c.ParamsParser(&userIDReq); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid user ID", nil)
	}
	if validationErrors := h.validator.Validate(userIDReq); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	// validation - body
	var req dto.UpdateProfileUserRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	// update profile
	err := h.userService.UpdateProfile(userIDReq.UserId, req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return util.SendError(c, fiber.StatusNotFound, "user not found", nil)
		}
		if errors.Is(err, service.ErrInvalidOldPassword) {
			return util.SendError(c, fiber.StatusBadRequest, "old password is incorrect", nil)
		}
		return util.SendError(c, fiber.StatusInternalServerError, "failed to update user", nil)
	}

	// response
	return util.SendSuccess(c, fiber.StatusOK, "user updated successfully", nil)
}
