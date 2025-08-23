package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/handlers"
)

func AuthRoutes(router fiber.Router, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
}
