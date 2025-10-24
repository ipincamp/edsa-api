package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

func AuthMiddleware(tokenSvc usecase.TokenService, blacklistSvc usecase.SessionBlacklistService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Missing Authorization Header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid Authorization Header format")
		}

		tokenString := parts[1]

		// 1. Validasi token DULU untuk dapat sessionID
		userID, sessionID, err := tokenSvc.ValidateToken(tokenString)
		if err != nil {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", err.Error())
		}

		// 2. Cek apakah SESSION-nya di-blacklist
		isBlacklisted, err := blacklistSvc.IsSessionBlacklisted(c.Context(), sessionID)
		if err != nil {
			return utils.SendSimpleError(c, fiber.StatusInternalServerError, "Session error", err.Error())
		}
		if isBlacklisted {
			// Sesi ini sudah logout
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Token has been logged out")
		}

		// Simpan user ID dan session ID di context
		c.Locals("userID", userID)
		c.Locals("sessionID", sessionID)

		return c.Next()
	}
}
