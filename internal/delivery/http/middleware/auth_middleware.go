package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

func AuthMiddleware(tokenSvc usecase.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Missing Authorization Header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Invalid Authorization Header format")
		}

		tokenString := parts[1]
		userID, sessionID, err := tokenSvc.ValidateToken(tokenString)
		if err != nil {
			return utils.SendError(c, fiber.StatusUnauthorized, err.Error())
		}

		// Simpan user ID di context untuk handler selanjutnya
		c.Locals("userID", userID)
		c.Locals("sessionID", sessionID)
		return c.Next()
	}
}
