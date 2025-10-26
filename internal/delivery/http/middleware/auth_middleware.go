package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/pkg/applogger"
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

		// 1. Validasi token DULU untuk dapat userID, sessionID, dan email
		userID, sessionID, userEmail, err := tokenSvc.ValidateToken(tokenString)
		if err != nil {
			// Sertakan detail error dari ValidateToken
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid or expired token", err.Error())
		}

		// 2. Cek apakah SESSION-nya di-blacklist
		isBlacklisted, err := blacklistSvc.IsSessionBlacklisted(c.Context(), sessionID)
		if err != nil {
			// Log error internal ini
			applogger.ErrorLogger.Printf("AuthMiddleware: Error checking blacklist for session %s: %v", sessionID, err)
			return utils.SendSimpleError(c, fiber.StatusInternalServerError, "Session check error", "Could not verify session status")
		}
		if isBlacklisted {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Session expired", "This session has been logged out")
		}

		// Simpan user ID, session ID, dan email di context
		c.Locals("userID", userID)
		c.Locals("sessionID", sessionID)
		c.Locals("userEmail", userEmail)

		return c.Next()
	}
}
