package handler

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/util"
)

// ValidationError represents validation errors with HTTP status code
type ValidationError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(code int, message string, details interface{}) *ValidationError {
	return &ValidationError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// handleValidationError handles validation error and sends appropriate HTTP response
func handleValidationError(ctx *fiber.Ctx, err error) error {
	if validationErr, ok := err.(*ValidationError); ok {
		return dto.SendError(ctx, validationErr.Code, validationErr.Message, validationErr.Details)
	}
	return dto.SendError(ctx, fiber.StatusInternalServerError, "Internal server error", err.Error())
}

// handleServiceError handles service layer errors and sends appropriate HTTP response
func handleServiceError(ctx *fiber.Ctx, err error, defaultMessage string) error {
	switch {
	case strings.Contains(err.Error(), "not found"):
		return dto.SendError(ctx, fiber.StatusNotFound, err.Error())
	case strings.Contains(err.Error(), "already exists"):
		return dto.SendError(ctx, fiber.StatusConflict, err.Error())
	case strings.Contains(err.Error(), "invalid credentials"):
		return dto.SendError(ctx, fiber.StatusUnauthorized, err.Error())
	case strings.Contains(err.Error(), "invalid input"):
		return dto.SendError(ctx, fiber.StatusBadRequest, err.Error())
	case strings.Contains(err.Error(), "unauthorized"):
		return dto.SendError(ctx, fiber.StatusUnauthorized, err.Error())
	case strings.Contains(err.Error(), "forbidden"):
		return dto.SendError(ctx, fiber.StatusForbidden, err.Error())
	default:
		return dto.SendError(ctx, fiber.StatusInternalServerError, defaultMessage, err.Error())
	}
}

// validateRequest validates a request struct and returns ValidationError if invalid
func validateRequest(validator *util.GoValidator, request interface{}) error {
	if errs := validator.Validate(request); len(errs) > 0 {
		return NewValidationError(fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}
	return nil
}

// parseAndValidateBody parses request body, validates it, and returns context with timeout
func parseAndValidateBody[T any](ctx *fiber.Ctx, validator *util.GoValidator, request *T) (context.Context, context.CancelFunc, error) {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)

	if err := ctx.BodyParser(request); err != nil {
		cancel()
		return nil, nil, NewValidationError(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := validateRequest(validator, request); err != nil {
		cancel()
		return nil, nil, err
	}

	return c, cancel, nil
}

// parseAndValidateQuery parses query parameters and validates them
func parseAndValidateQuery[T any](ctx *fiber.Ctx, validator *util.GoValidator, request *T) error {
	if err := ctx.QueryParser(request); err != nil {
		return NewValidationError(fiber.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	return validateRequest(validator, request)
}

// parseAndValidateParams parses path parameters and validates them
func parseAndValidateParams[T any](ctx *fiber.Ctx, validator *util.GoValidator, request *T) error {
	if err := ctx.ParamsParser(request); err != nil {
		return NewValidationError(fiber.StatusBadRequest, "Invalid path parameters", err.Error())
	}

	return validateRequest(validator, request)
}

// createContextWithTimeout creates a context with 10 second timeout
func createContextWithTimeout(ctx *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx.Context(), 10*time.Second)
}
