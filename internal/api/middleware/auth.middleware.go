package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/util"
)

type AuthMiddleware struct {
	config *config.Config
}

func NewAuth(config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
	}
}

func (m *AuthMiddleware) Auth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return dto.SendError(ctx, fiber.StatusUnauthorized, "Authorization header is required", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return dto.SendError(ctx, fiber.StatusUnauthorized, "Invalid authorization header format", nil)
		}

		tokenString := parts[1]
		pasetoMaker, err := util.NewPasetoMaker(m.config.Paseto.SecretKey)
		if err != nil {
			return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to create token maker", err)
		}

		payload, err := pasetoMaker.VerifyToken(tokenString)
		if err != nil {
			return dto.SendError(ctx, fiber.StatusUnauthorized, "Invalid or expired token", err)
		}

		userFromToken := domain.User{
			ID: payload.UserID,
			Role: domain.Role{
				ID: payload.RoleID,
			},
		}

		ctx.Locals("user", userFromToken)

		return ctx.Next()
	}
}
