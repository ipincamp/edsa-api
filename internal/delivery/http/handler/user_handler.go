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

	profile, err := h.userService.GetProfile(user.ID)
	if err != nil {
		status, msg := mapUserProfileError(err)
		return util.SendError(c, status, msg, nil)
	}
	return util.SendSuccess(c, fiber.StatusOK, "profile retrieved successfully", profile)
}

// mapUserProfileError memetakan error pada proses get profile ke status dan pesan yang sesuai
func mapUserProfileError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		return fiber.StatusNotFound, "user not found"
	default:
		return fiber.StatusInternalServerError, "failed to get profile"
	}
}

// UpdateProfile - Update current user profile (self)
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	if user == nil {
		return util.SendError(c, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	var req dto.UpdateProfileUserRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	err := h.userService.UpdateProfile(user.ID, req)
	if err != nil {
		status, msg := mapUpdateProfileError(err)
		return util.SendError(c, status, msg, nil)
	}
	return util.SendSuccess(c, fiber.StatusOK, "profile updated successfully", nil)
}

// mapUpdateProfileError memetakan error pada proses update profile ke status dan pesan yang sesuai
func mapUpdateProfileError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		return fiber.StatusNotFound, "user not found"
	case errors.Is(err, service.ErrInvalidOldPassword):
		return fiber.StatusBadRequest, "old password is incorrect"
	default:
		return fiber.StatusInternalServerError, "failed to update profile"
	}
}

// GetAllUsers - Get all users with pagination (admin only)
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	var req dto.UserFilterRequest
	if err := c.QueryParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid query parameters", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	users, err := h.userService.GetAllUsers(req.Page, req.Limit, req.Role)
	if err != nil {
		status, msg := mapGetAllUsersError(err)
		return util.SendError(c, status, msg, nil)
	}
	return util.SendSuccess(c, fiber.StatusOK, "users retrieved successfully", users)
}

// mapGetAllUsersError memetakan error pada proses get all users ke status dan pesan yang sesuai
func mapGetAllUsersError(err error) (int, string) {
	return fiber.StatusInternalServerError, "failed to get users"
}

// GetUserByID - Get user by ID (admin only)
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	var req dto.UserIDRequest
	if err := c.ParamsParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid user ID", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	profile, err := h.userService.GetProfile(req.UserId)
	if err != nil {
		status, msg := mapGetUserByIDError(err)
		return util.SendError(c, status, msg, nil)
	}
	return util.SendSuccess(c, fiber.StatusOK, "user retrieved successfully", profile)
}

// mapGetUserByIDError memetakan error pada proses get user by ID ke status dan pesan yang sesuai
func mapGetUserByIDError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		return fiber.StatusNotFound, "user not found"
	default:
		return fiber.StatusInternalServerError, "failed to get user"
	}
}

// UpdateUserByID - Update user by ID (admin only)
func (h *UserHandler) UpdateUserByID(c *fiber.Ctx) error {
	var userIDReq dto.UserIDRequest
	if err := c.ParamsParser(&userIDReq); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid user ID", nil)
	}
	if validationErrors := h.validator.Validate(userIDReq); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	var req dto.UpdateProfileUserRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	err := h.userService.UpdateProfile(userIDReq.UserId, req)
	if err != nil {
		status, msg := mapUpdateUserByIDError(err)
		return util.SendError(c, status, msg, nil)
	}
	return util.SendSuccess(c, fiber.StatusOK, "user updated successfully", nil)
}

// mapUpdateUserByIDError memetakan error pada proses update user by ID ke status dan pesan yang sesuai
func mapUpdateUserByIDError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		return fiber.StatusNotFound, "user not found"
	case errors.Is(err, service.ErrInvalidOldPassword):
		return fiber.StatusBadRequest, "old password is incorrect"
	default:
		return fiber.StatusInternalServerError, "failed to update user"
	}
}
