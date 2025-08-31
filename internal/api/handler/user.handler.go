package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/util"
)

type UserHandler struct {
	UserService domain.UserService
	Validator   *util.GoValidator
}

func NewUser(userService domain.UserService, validator *util.GoValidator) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Validator:   validator,
	}
}

func (h *UserHandler) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	res, err := h.UserService.Index(c)
	if err != nil {
		return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to retrieve users", err)
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "success", res)
}

func (h *UserHandler) Profile(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	loggedInUser, ok := ctx.Locals("user").(domain.User)
	if !ok {
		return dto.SendError(ctx, fiber.StatusUnauthorized, "Failed to identify user from token", nil)
	}

	res, err := h.UserService.Profile(c, loggedInUser.ID)
	if err != nil {
		return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to retrieve user profile", err)
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "success", res)
}

func (h *UserHandler) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	userID := ctx.Params("userID")
	if userID == "" {
		loggedInUser, ok := ctx.Locals("user").(domain.User)
		if !ok {
			return dto.SendError(ctx, fiber.StatusUnauthorized, "Failed to identify user from token", nil)
		}
		userID = loggedInUser.ID
	}

	var request dto.UpdateUserRequest
	if err := ctx.BodyParser(&request); err != nil {
		return dto.SendError(ctx, fiber.StatusBadRequest, "Invalid request body", err)
	}
	if errs := h.Validator.Validate(request); errs != nil {
		return dto.SendError(ctx, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	err := h.UserService.Update(c, userID, request)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return dto.SendError(ctx, fiber.StatusNotFound, err.Error(), nil)
		case "invalid old password":
			return dto.SendError(ctx, fiber.StatusUnauthorized, err.Error(), nil)
		case "old password is required to set a new password":
			return dto.SendError(ctx, fiber.StatusBadRequest, err.Error(), nil)
		default:
			return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to update user data", err)
		}
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "User data updated successfully", nil)
}
