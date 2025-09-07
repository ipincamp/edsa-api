package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"github.com/ipincamp/go-edsa-api/pkg/cache"
	"github.com/ipincamp/go-edsa-api/pkg/token"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

var (
	ErrWrongTokenType = errors.New("wrong token type")
)

func AuthMiddleware(tokenMaker *token.PasetoMaker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return util.SendError(c, fiber.StatusUnauthorized, "authorization header is not provided", nil)
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return util.SendError(c, fiber.StatusUnauthorized, "invalid authorization header format", nil)
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			return util.SendError(c, fiber.StatusUnauthorized, "unsupported authorization type", nil)
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return util.SendError(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		if payload.TokenType != "access" {
			return util.SendError(c, fiber.StatusUnauthorized, ErrWrongTokenType.Error(), nil)
		}

		user, found := cache.GetUserFromCacheByID(payload.UserID)
		if !found {
			return util.SendError(c, fiber.StatusUnauthorized, "user for this token not found", nil)
		}

		c.Locals("user", &user)
		return c.Next()
	}
}
