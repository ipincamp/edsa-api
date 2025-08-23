package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/handlers"
	"github.com/ipincamp/edsa/internal/api/middleware"
	"github.com/ipincamp/edsa/internal/config"
)

func AuthRoutes(router fiber.Router, authHandler *handlers.AuthHandler, env *config.Env) {
	auth := router.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", middleware.Protected(env), authHandler.Logout)
}
