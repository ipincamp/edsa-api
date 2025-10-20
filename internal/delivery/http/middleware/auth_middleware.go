package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		userID, err := tokenSvc.ValidateToken(tokenString)
		if err != nil {
			return utils.SendError(c, fiber.StatusUnauthorized, err.Error())
		}

		// Simpan user ID di context untuk handler selanjutnya
		c.Locals("userID", userID)

		// Ambil Session ID dari header
		sessionIDHeader := c.Get("X-Session-ID")
		sessionID, err := uuid.Parse(sessionIDHeader)
		if err != nil {
			// Jika header tidak ada atau invalid, kita set Nil.
			// Logger akan membuatkan ID baru jika diperlukan (fallback).
			sessionID = uuid.Nil
		}
		c.Locals("sessionID", sessionID)

		return c.Next()
	}
}
