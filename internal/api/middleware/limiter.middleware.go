package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/ipincamp/go-edsa-api/domain/dto"
)

type LimiterMiddleware struct{}

func NewLimiter() *LimiterMiddleware {
	return &LimiterMiddleware{}
}

func (lm *LimiterMiddleware) Make(max int, expirationMinutes int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: time.Duration(expirationMinutes) * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return dto.SendError(c, fiber.StatusTooManyRequests, "Too many requests, please try again later.")
		},
	})
}
