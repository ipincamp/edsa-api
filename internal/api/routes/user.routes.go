package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/edsa/internal/api/handlers"
	"github.com/ipincamp/edsa/internal/api/middleware"
	"github.com/ipincamp/edsa/internal/config"
)

func UserRoutes(router fiber.Router, userHandler *handlers.UserHandler, env *config.Env) {
	router.Get("/profile", middleware.Protected(env), userHandler.GetProfile)
}
