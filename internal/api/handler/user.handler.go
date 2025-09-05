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
	var filter dto.UserFilterRequest
	if err := parseAndValidateQuery(ctx, h.Validator, &filter); err != nil {
		return err
	}

	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	res, err := h.UserService.All(c, filter.Page, filter.Limit, filter.Role)
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
	c, cancel, err := parseAndValidateBody(ctx, h.Validator, &request)
	if err != nil {
		return err
	}
	defer cancel()

	err = h.UserService.UpdateProfile(c, loggedInUser.ID, request)
	if err != nil {
		switch {
		case errors.Is(err, constant.ErrNotFound):
			return dto.SendError(ctx, fiber.StatusNotFound, err.Error())
		case errors.Is(err, constant.ErrUnauthorized):
			return dto.SendError(ctx, fiber.StatusUnauthorized, err.Error())
		case errors.Is(err, constant.ErrInvalidInput):
			return dto.SendError(ctx, fiber.StatusBadRequest, err.Error())
		default:
			return dto.SendError(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "Profile updated successfully", nil)
}

func (h *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	var params dto.UserIDRequest
	if err := ctx.ParamsParser(&params); err != nil {
		return dto.SendError(ctx, fiber.StatusBadRequest, "Invalid request parameters", err.Error())
	}

	if validationErrs := h.Validator.Validate(&params); validationErrs != nil {
		return dto.SendError(ctx, fiber.StatusUnprocessableEntity, "Validation failed", validationErrs)
	}

	var request dto.UpdateProfileUserRequest
	if err := ctx.BodyParser(&request); err != nil {
		return dto.SendError(ctx, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrs := h.Validator.Validate(&request); validationErrs != nil {
		return dto.SendError(ctx, fiber.StatusUnprocessableEntity, "Validation failed", validationErrs)
	}

	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	err := h.UserService.UpdateProfile(c, params.UserId, request)
	if err != nil {
		switch {
		case errors.Is(err, constant.ErrNotFound):
			return dto.SendError(ctx, fiber.StatusNotFound, err.Error())
		case errors.Is(err, constant.ErrUnauthorized):
			return dto.SendError(ctx, fiber.StatusUnauthorized, err.Error())
		case errors.Is(err, constant.ErrInvalidInput):
			return dto.SendError(ctx, fiber.StatusBadRequest, err.Error())
		default:
			return dto.SendError(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}

	return dto.SendSuccess(ctx, fiber.StatusOK, "User updated successfully", nil)
}
