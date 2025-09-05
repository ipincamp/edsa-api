package handler

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/service"
	"github.com/ipincamp/go-edsa-api/internal/util"
)

type UserHandler struct {
	UserService service.UserService
	Validator   *util.GoValidator
}

func NewUser(userService service.UserService, validator *util.GoValidator) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Validator:   validator,
	}
}

func (h *UserHandler) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	page, limit := util.GetPaginationParams(ctx)

	res, err := h.UserService.All(c, page, limit)
	if err != nil {
		return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to retrieve users", err.Error())
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "Users retrieved successfully", res)
}

func (h *UserHandler) Profile(ctx *fiber.Ctx) error {
	loggedInUser, ok := ctx.Locals("user").(domain.User)
	if !ok {
		return dto.SendError(ctx, fiber.StatusUnauthorized, "Invalid user data in token")
	}

	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	res, err := h.UserService.Profile(c, loggedInUser.ID)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return dto.SendError(ctx, fiber.StatusNotFound, err.Error())
		}
		return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to retrieve user profile")
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "Profile retrieved successfully", res)
}

func (h *UserHandler) UpdateProfile(ctx *fiber.Ctx) error {
	loggedInUser, ok := ctx.Locals("user").(domain.User)
	if !ok {
		return dto.SendError(ctx, fiber.StatusUnauthorized, "Invalid user data in token")
	}

	var request dto.UpdateProfileUserRequest
	c, cancel, err := parseAndValidate(ctx, h.Validator, &request)
	if err != nil {
		return err
	}
	defer cancel()

	err = h.UserService.UpdateProfile(c, loggedInUser.ID, request)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return dto.SendError(ctx, fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, constant.ErrUnauthorized) {
			return dto.SendError(ctx, fiber.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, constant.ErrInvalidInput) {
			return dto.SendError(ctx, fiber.StatusBadRequest, err.Error())
		}
		return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to update profile")
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "Profile updated successfully", nil)
}

func (h *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")
	if userID == "" {
		return dto.SendError(ctx, fiber.StatusBadRequest, "User ID is required")
	}

	var request dto.UpdateProfileUserRequest
	c, cancel, err := parseAndValidate(ctx, h.Validator, &request)
	if err != nil {
		return err
	}
	defer cancel()

	err = h.UserService.UpdateProfile(c, userID, request)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return dto.SendError(ctx, fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, constant.ErrInvalidInput) {
			return dto.SendError(ctx, fiber.StatusBadRequest, err.Error())
		}
		return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to update user")
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "User updated successfully", nil)
}
