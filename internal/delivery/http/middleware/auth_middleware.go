package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"github.com/ipincamp/go-edsa-api/pkg/cache"
	"github.com/ipincamp/go-edsa-api/pkg/token"
)

// AuthorizationHeaderKey adalah nama header untuk otentikasi
const AuthorizationHeaderKey = "authorization"

// AuthorizationTypeBearer adalah tipe otorisasi yang didukung
const AuthorizationTypeBearer = "bearer"

// AuthorizationPayloadKey adalah key untuk payload otorisasi di context
const AuthorizationPayloadKey = "authorization_payload"

// ErrWrongTokenType digunakan jika token bukan tipe akses
var ErrWrongTokenType = errors.New("wrong token type")

// AuthMiddleware adalah middleware untuk otentikasi JWT Bearer
func AuthMiddleware(tokenMaker *token.PasetoMaker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil header otorisasi
		authorizationHeader := c.Get(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return sendAuthError(c, fiber.StatusUnauthorized, "authorization header is not provided")
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return sendAuthError(c, fiber.StatusUnauthorized, "invalid authorization header format")
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			return sendAuthError(c, fiber.StatusUnauthorized, "unsupported authorization type")
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return sendAuthError(c, fiber.StatusUnauthorized, err.Error())
		}

		if payload.TokenType != "access" {
			return sendAuthError(c, fiber.StatusUnauthorized, ErrWrongTokenType.Error())
		}

		user, found := cache.GetUserFromCacheByID(payload.UserID)
		if !found {
			return sendAuthError(c, fiber.StatusUnauthorized, "user for this token not found")
		}

		c.Locals("user", &user)
		return c.Next()
	}
}

// sendAuthError adalah helper untuk mengirim error otentikasi
func sendAuthError(c *fiber.Ctx, status int, message string) error {
	return util.SendError(c, status, message, nil)
}
