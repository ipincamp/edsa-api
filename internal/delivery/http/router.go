package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/middleware"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

func SetupRoutes(
	app *fiber.App,
	userHandler *UserHandler,
	adminHandler *AdminHandler,
	tokenSvc usecase.TokenService,
	userRepo usecase.UserRepository,
) {
	app.Use(logger.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Grup API v1
	api := app.Group("/api/v1")

	// Rute Autentikasi
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// Rute yang dilindungi
	protected := api.Group("/users")
	protected.Use(middleware.AuthMiddleware(tokenSvc))
	protected.Get("/me", userHandler.GetMe)

	// --- Rute Administrasi ---
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(tokenSvc))
	admin.Use(middleware.AdminMiddleware(userRepo))

	// Rute Subjects
	subjects := admin.Group("/subjects")
	subjects.Post("/", adminHandler.CreateSubject)
	subjects.Get("/", adminHandler.GetAllSubjects)
	subjects.Get("/:id", adminHandler.GetSubjectByID)
	subjects.Put("/:id", adminHandler.UpdateSubject)
	subjects.Delete("/:id", adminHandler.DeleteSubject)

	// Rute Classes
	classes := admin.Group("/classes")
	classes.Post("/", adminHandler.CreateClass)
	classes.Get("/", adminHandler.GetAllClasses)
	classes.Get("/:id", adminHandler.GetClassByID)
	classes.Put("/:id", adminHandler.UpdateClass)
	classes.Delete("/:id", adminHandler.DeleteClass)

	// Rute Groups
	groups := admin.Group("/groups")
	groups.Post("/", adminHandler.CreateGroup)
	groups.Get("/", adminHandler.GetAllGroups)
	groups.Get("/:id", adminHandler.GetGroupByID)
	groups.Put("/:id", adminHandler.UpdateGroup)
	groups.Delete("/:id", adminHandler.DeleteGroup)

	// TODO: Rute manajemen penugasan (Assign/Unassign)
}
