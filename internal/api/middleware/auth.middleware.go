package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/response"
	"github.com/ipincamp/edsa/internal/config"
	"github.com/ipincamp/edsa/internal/utils"
)

const (
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationTypeBearer = "Bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func Protected(env *config.Env) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(AuthorizationHeaderKey)
		if len(authHeader) == 0 {
			return response.Error(c, fiber.StatusUnauthorized, "Authorization header is not provided")
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		authType := strings.ToLower(fields[0])
		if authType != strings.ToLower(AuthorizationTypeBearer) {
			return response.Error(c, fiber.StatusUnauthorized, "Unsupported authorization type")
		}

		accessToken := fields[1]
		payload, err := utils.VerifyToken(accessToken, env.PasetoSymmetricKey)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		c.Locals(AuthorizationPayloadKey, payload)
		return c.Next()
	}
}
