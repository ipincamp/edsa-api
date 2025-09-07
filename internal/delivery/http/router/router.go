package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/handler"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/middleware"
	"github.com/ipincamp/go-edsa-api/internal/repository"
	"github.com/ipincamp/go-edsa-api/internal/service"
	"github.com/ipincamp/go-edsa-api/pkg/bloom"
	"github.com/ipincamp/go-edsa-api/pkg/token"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	// Config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Token
	tokenMaker, err := token.NewPasetoMaker(cfg.Paseto.SecretKey)
	if err != nil {
		panic(err)
	}

	// Bloom Filter
	bloomFilter, err := bloom.NewBloomFilterManager(cfg.Bloom.EmailFilterPath, 100000, 0.01)
	if err != nil {
		panic(err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Populate Bloom Filter on startup
	go func() {
		emails, err := userRepo.FindAllEmails()
		if err != nil {
			log.Printf("Failed to get emails for bloom filter: %v", err)
			return
		}
		if len(emails) > 0 {
			bloomFilter.Regenerate(emails)
			if err := bloomFilter.Save(); err != nil {
				log.Printf("Failed to save bloom filter: %v", err)
			}
			log.Printf("Bloom filter regenerated with %d emails", len(emails))
		}
	}()

	// Services
	authService := service.NewAuthService(
		db,
		userRepo,
		roleRepo,
		*tokenMaker,
		cfg.Paseto.AccessTokenDuration,
		cfg.Paseto.RefreshTokenDuration,
		bloomFilter,
	)

	userService := service.NewUserService(db, userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// Middleware
	authRequired := middleware.AuthMiddleware(tokenMaker)

	// Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", authRequired, authHandler.Logout)

	// Protected user routes
	users := v1.Group("/users").Use(authRequired)

	// Self profile routes (accessible by all authenticated users)
	users.Get("/profile", userHandler.GetProfile)
	users.Put("/profile", userHandler.UpdateProfile)

	// Admin-only routes
	users.Get("/", middleware.RequireAdmin(), userHandler.GetAllUsers) // GET /api/v1/users?page=1&limit=10&role=student

	// Resource ownership or admin access
	users.Get("/:userId", middleware.RequireOwnershipOrAdmin("userId"), userHandler.GetUserByID)    // GET /api/v1/users/{userId}
	users.Put("/:userId", middleware.RequireOwnershipOrAdmin("userId"), userHandler.UpdateUserByID) // PUT /api/v1/users/{userId}

	// Example of permission-based access (uncomment when needed)
	// users.Post("/", middleware.RequirePermission("users.create"), userHandler.CreateUser)
	// users.Delete("/:userId", middleware.RequirePermission("users.delete"), userHandler.DeleteUser)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the API")
	})
}
