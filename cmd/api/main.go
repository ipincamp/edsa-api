package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http"
	"github.com/ipincamp/go-edsa-api/internal/pkg/applogger"
	"github.com/ipincamp/go-edsa-api/internal/pkg/database"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	"github.com/ipincamp/go-edsa-api/internal/repository/cache"
	"github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
	"github.com/ipincamp/go-edsa-api/internal/service/email"
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
	log.Println("🚀 Starting EDSA API application...")

	// 0. Init Error Logger
	log.Println("📝 Initializing Error Logger...")
	applogger.InitErrorLogger() // Assuming InitErrorLogger already logs its success/failure
	log.Println("✅ Error Logger Initialized.")

	// 1. Load Config
	log.Println("⚙️ Loading Configuration...")
	config.LoadConfig()
	cfg := config.AppConfig
	log.Println("✅ Configuration Loaded.")

	// 2. Init Database
	log.Println("🔧 Initializing Database Connection...")
	db, err := database.NewPostgresConnection(cfg) // Pass cfg
	if err != nil {
		log.Fatalf("🚨 CRITICAL: Failed to initialize database: %v", err) // Use a distinct emoji for critical errors
	}
	log.Println("✅ Database connection established and pool configured.")

	// 3. Init Validator
	log.Println("🔍 Initializing Validator...")
	validate := validator.NewValidator()
	log.Println("✅ Validator Initialized.")

	// 4. Init Services (Grouped Log)
	log.Println("🛠️ Initializing Core Services (Password, Token, Storage, Email)...")
	passwordService := argon2id.NewPasswordService()
	tokenService, err := paseto.NewPasetoService(cfg.Security.PasetoSymmetricKey)
	if err != nil {
		log.Fatalf("🚨 CRITICAL: Failed to init Paseto service: %v", err)
	}
	fileStorageService := storage.NewLocalStorageService(cfg)
	emailService := email.NewSmtpService(cfg)
	log.Println("✅ Core Services Initialized.")

	// 5. Init Repositories (Grouped Log)
	log.Println("📦 Initializing Repositories (GORM & Cache)...")
	roleRepositoryGORM := gorm.NewRoleRepository(db)
	roleRepositoryCACHE, err := cache.NewRoleRepositoryCACHE(roleRepositoryGORM)
	if err != nil {
		log.Fatalf("🚨 CRITICAL: Failed to create and pre-load role cache: %v", err)
	}
	userRepositoryGORM := gorm.NewUserRepository(db, roleRepositoryCACHE)
	subjectRepositoryGORM := gorm.NewSubjectRepository(db)
	classRepositoryGORM := gorm.NewClassRepository(db)
	groupRepositoryGORM := gorm.NewGroupRepository(db)
	bookRepositoryGORM := gorm.NewBookRepository(db)
	pageRepositoryGORM := gorm.NewPageRepository(db)
	interactionRepositoryGORM := gorm.NewInteractionRepository(db)
	progressRepositoryGORM := gorm.NewUserBookProgressRepository(db)
	gameRepositoryGORM := gorm.NewGameRepository(db)
	userGameScoreRepositoryGORM := gorm.NewUserGameScoreRepository(db)
	activityLogRepositoryGORM := gorm.NewActivityLogRepository(db)
	mediaAssetRepositoryGORM := gorm.NewMediaAssetRepository(db)
	groupBookSettingRepositoryGORM := gorm.NewGroupBookSettingRepository(db)
	log.Println("✅ Repositories Initialized.")

	// 6. Init Usecases (Grouped Log)
	log.Println("🧠 Initializing Usecases/Services...")
	loggerService := logger.NewActivityLoggerService(activityLogRepositoryGORM) // Logs its own start internally
	sessionBlacklistService := cache.NewSessionBlacklistCACHE()                 // Simple in-memory, might not need detailed log
	userService := user.NewUserService(
		userRepositoryGORM, roleRepositoryCACHE, passwordService, tokenService, cfg,
		loggerService, sessionBlacklistService, emailService, progressRepositoryGORM, mediaAssetRepositoryGORM,
	)
	adminService := admin.NewAdminService(
		subjectRepositoryGORM, classRepositoryGORM, groupRepositoryGORM,
		bookRepositoryGORM, pageRepositoryGORM, interactionRepositoryGORM,
	)
	appService := app.NewAppService(
		bookRepositoryGORM, progressRepositoryGORM, gameRepositoryGORM, userGameScoreRepositoryGORM,
		loggerService, userRepositoryGORM, groupRepositoryGORM, groupBookSettingRepositoryGORM,
	)
	dashboardService := dashboard.NewDashboardService(
		activityLogRepositoryGORM, userRepositoryGORM, groupRepositoryGORM,
		bookRepositoryGORM, groupBookSettingRepositoryGORM,
	)
	mediaService := media.NewMediaService(
		fileStorageService, mediaAssetRepositoryGORM, cfg,
	)
	log.Println("✅ Usecases/Services Initialized.")

	// 7. Init Handlers (Grouped Log)
	log.Println("🔌 Initializing HTTP Handlers...")
	userHandler := http.NewUserHandler(userService, mediaService, loggerService, validate)
	adminHandler := http.NewAdminHandler(adminService, validate)
	appHandler := http.NewAppHandler(appService, dashboardService, validate)
	dashboardHandler := http.NewDashboardHandler(dashboardService, validate)
	mediaHandler := http.NewMediaHandler(mediaService, validate, loggerService)
	log.Println("✅ HTTP Handlers Initialized.")

	// 8. Init Fiber App
	log.Println("🌐 Initializing Fiber Application...")
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "An unexpected error occurred" // Default internal error message

			var e *fiber.Error
			if errors.As(err, &e) { // Use errors.As for type assertion
				code = e.Code
				message = e.Message
			}

			// Log the actual error internally using the ErrorLogger
			applogger.ErrorLogger.Printf("Fiber ErrorHandler caught: Code=%d, Message=%s, Error=%v", code, message, err)

			// Send structured error response to client
			return utils.SendSimpleError(c, code, message, err.Error())
		},
	})
	log.Println("✅ Fiber Application Initialized.")

	// Setup CORS Middleware
	log.Println("🔗 Configuring CORS Middleware...")
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Security.CorsAllowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	log.Println("✅ CORS Middleware Configured.")

	// 9. Setup Routes
	log.Println("🛣️ Setting up API Routes...")
	http.SetupRoutes(
		app, userHandler, adminHandler, appHandler, dashboardHandler, mediaHandler,
		tokenService, userRepositoryGORM, sessionBlacklistService, cfg,
	)
	log.Println("✅ API Routes Configured.")

	// 10. Start Server with Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		serverAddr := fmt.Sprintf(":%d", cfg.App.Port)
		log.Printf("👂 Starting server, listening on %s...", serverAddr)
		if err := app.Listen(serverAddr); err != nil {
			log.Printf("🔥 Server failed to start or stopped: %v", err)
			// Signal the main goroutine to stop if the server fails immediately
			// Use non-blocking send in case quit channel is already closed or full
			select {
			case quit <- syscall.SIGTERM:
			default:
			}
		}
	}()

	// Block until signal is received
	sig := <-quit
	log.Printf("🛑 Received signal: %s. Initiating graceful shutdown...", sig)

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Increased timeout slightly
	defer cancel()

	// Shutdown Fiber server
	log.Println("⏳ Shutting down Fiber server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("⚠️ Fiber server shutdown failed: %v", err)
	} else {
		log.Println("✅ Fiber server shut down gracefully.")
	}

	// Shutdown Logger Service
	log.Println("⏳ Shutting down Activity Logger service...")
	// Pass the shutdown context to the logger service shutdown
	if err := loggerService.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ Activity Logger service shutdown failed: %v", err)
	}
	// Logger service logs its own success message internally

	log.Println("🏁 Application shut down complete.")
}
