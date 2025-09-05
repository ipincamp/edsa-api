package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/util"
)

func parseAndValidateBody[T any](ctx *fiber.Ctx, validator *util.GoValidator, request *T) (context.Context, context.CancelFunc, error) {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)

	if err := ctx.BodyParser(request); err != nil {
		cancel()
		return nil, nil, dto.SendError(ctx, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if errs := validator.Validate(request); errs != nil {
		cancel()
		return nil, nil, dto.SendError(ctx, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	return c, cancel, nil
}

func parseAndValidateQuery[T any](ctx *fiber.Ctx, validator *util.GoValidator, request *T) error {
	if err := ctx.QueryParser(request); err != nil {
		return dto.SendError(ctx, fiber.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	if errs := validator.Validate(request); errs != nil {
		return dto.SendError(ctx, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	return nil
}

func parseAndValidateParams[T any](ctx *fiber.Ctx, validator *util.GoValidator, request *T) error {
	if err := ctx.ParamsParser(request); err != nil {
		return dto.SendError(ctx, fiber.StatusBadRequest, "Invalid path parameters", err.Error())
	}

	if errs := validator.Validate(request); errs != nil {
		return dto.SendError(ctx, fiber.StatusUnprocessableEntity, "Validation failed", errs)
	}

	return nil
}
