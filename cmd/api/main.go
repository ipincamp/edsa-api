package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http"
	"github.com/ipincamp/go-edsa-api/internal/pkg/database"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
	"github.com/ipincamp/go-edsa-api/internal/service/paseto"
	"github.com/ipincamp/go-edsa-api/internal/usecase/admin"
	"github.com/ipincamp/go-edsa-api/internal/usecase/user"
)

func main() {
	// 1. Load Config
	config.LoadConfig()
	cfg := config.AppConfig

	// 2. Init Database (PostgreSQL + GORM)
	db := database.NewPostgresConnection(config.GetDatabaseDSN())

	// 3. Init Validator
	validate := validator.NewValidator()

	// 4. Init Services
	passwordService := argon2id.NewPasswordService()
	tokenService, err := paseto.NewPasetoService(cfg.Security.PasetoSymmetricKey)
	if err != nil {
		log.Fatalf("Failed to init Paseto service: %v", err)
	}

	// 5. Init Repositories
	userRepository := repo.NewUserRepository(db)
	roleRepository := repo.NewRoleRepository(db)
	// Repositori Manajemen User
	subjectRepository := repo.NewSubjectRepository(db)
	classRepository := repo.NewClassRepository(db)
	groupRepository := repo.NewGroupRepository(db)
	// Repositori Konten
	bookRepository := repo.NewBookRepository(db)
	pageRepository := repo.NewPageRepository(db)
	interactionRepository := repo.NewInteractionRepository(db)

	// 6. Init Usecases
	userService := user.NewUserService(
		userRepository,
		roleRepository,
		passwordService,
		tokenService,
		cfg,
	)
	adminService := admin.NewAdminService(
		subjectRepository,
		classRepository,
		groupRepository,
		bookRepository,
		pageRepository,
		interactionRepository,
	)

	// 7. Init Handlers
	userHandler := http.NewUserHandler(userService, validate)
	adminHandler := http.NewAdminHandler(adminService, validate)

	// 8. Init Fiber App
	app := fiber.New(fiber.Config{
		// Error handling kustom
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		},
	})

	// 9. Setup Routes
	http.SetupRoutes(
		app,
		userHandler,
		adminHandler,
		tokenService,
		userRepository,
	)

	// 10. Start Server
	log.Printf("Starting server on port %d...", cfg.App.Port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.App.Port)))
}
