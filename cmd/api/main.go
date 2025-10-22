package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http"
	"github.com/ipincamp/go-edsa-api/internal/pkg/database"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
	"github.com/ipincamp/go-edsa-api/internal/service/paseto"
	"github.com/ipincamp/go-edsa-api/internal/service/storage"
	"github.com/ipincamp/go-edsa-api/internal/usecase/admin"
	"github.com/ipincamp/go-edsa-api/internal/usecase/app"
	"github.com/ipincamp/go-edsa-api/internal/usecase/dashboard"
	"github.com/ipincamp/go-edsa-api/internal/usecase/logger"
	"github.com/ipincamp/go-edsa-api/internal/usecase/media"
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
	fileStorageService := storage.NewLocalStorageService(cfg)

	// 5. Init Repositories
	userRepository := repo.NewUserRepository(db)
	roleRepository := repo.NewRoleRepository(db)
	// Admin
	subjectRepository := repo.NewSubjectRepository(db)
	classRepository := repo.NewClassRepository(db)
	groupRepository := repo.NewGroupRepository(db)
	// Konten
	bookRepository := repo.NewBookRepository(db)
	pageRepository := repo.NewPageRepository(db)
	interactionRepository := repo.NewInteractionRepository(db)
	// Progres
	progressRepository := repo.NewUserBookProgressRepository(db)
	// Game
	gameRepository := repo.NewGameRepository(db)
	userGameScoreRepository := repo.NewUserGameScoreRepository(db)
	// Log Aktivitas
	activityLogRepository := repo.NewActivityLogRepository(db)
	// Repo Media
	mediaAssetRepository := repo.NewMediaAssetRepository(db)

	// 6. Init Usecases
	loggerService := logger.NewActivityLoggerService(activityLogRepository)
	userService := user.NewUserService(
		userRepository,
		roleRepository,
		passwordService,
		tokenService,
		cfg,
		loggerService,
	)
	adminService := admin.NewAdminService(
		subjectRepository,
		classRepository,
		groupRepository,
		bookRepository,
		pageRepository,
		interactionRepository,
	)
	appService := app.NewAppService(
		bookRepository,
		progressRepository,
		gameRepository,
		userGameScoreRepository,
		loggerService,
	)
	dashboardService := dashboard.NewDashboardService(
		activityLogRepository,
		userRepository,
	)
	mediaService := media.NewMediaService(
		fileStorageService,
		mediaAssetRepository,
		cfg,
	)

	// 7. Init Handlers
	userHandler := http.NewUserHandler(userService, validate)
	adminHandler := http.NewAdminHandler(adminService, validate)
	appHandler := http.NewAppHandler(appService, validate)
	dashboardHandler := http.NewDashboardHandler(dashboardService, validate)
	mediaHandler := http.NewMediaHandler(mediaService, validate)

	// 8. Init Fiber App
	app := fiber.New(fiber.Config{
		// Error handling kustom
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Failed" // Pesan default

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message // Gunakan pesan dari Fiber jika ada
			}

			// Gunakan helper response error yang baru
			// 'message' adalah pesan utama, err.Error() adalah detail di 'value'
			return utils.SendSimpleError(c, code, message, err.Error())
		},
	})

	// 9. Setup Routes
	http.SetupRoutes(
		app,
		userHandler,
		adminHandler,
		appHandler,
		dashboardHandler,
		mediaHandler,
		tokenService,
		userRepository,
		cfg,
	)

	// 10. Start Server
	log.Printf("Starting server on port %d...", cfg.App.Port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.App.Port)))
}
