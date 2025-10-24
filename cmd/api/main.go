package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http"
	"github.com/ipincamp/go-edsa-api/internal/pkg/applogger"
	"github.com/ipincamp/go-edsa-api/internal/pkg/database"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	"github.com/ipincamp/go-edsa-api/internal/repository/cache"
	"github.com/ipincamp/go-edsa-api/internal/repository/gorm"
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
	// 0. Init Error Logger
	applogger.InitErrorLogger()

	// 1. Load Config
	config.LoadConfig()
	cfg := config.AppConfig

	// 2. Init Database (PostgreSQL + GORM)
	db := database.NewPostgresConnection(config.GetDatabaseDSN(), cfg.App.Env)

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
	// Buat GORM Role Repo (untuk di-pass ke cache)
	roleRepositoryGORM := gorm.NewRoleRepository(db)

	// Buat Cache Role Repo
	roleRepositoryCACHE, err := cache.NewRoleRepositoryCACHE(roleRepositoryGORM)
	if err != nil {
		log.Fatalf("Failed to create role cache: %v", err)
	}

	// Inject cache repo ke UserRepository
	// (Struct userRepositoryGORM sudah konsisten)
	userRepositoryGORM := gorm.NewUserRepository(db, roleRepositoryCACHE)

	// Admin
	subjectRepositoryGORM := gorm.NewSubjectRepository(db)
	classRepositoryGORM := gorm.NewClassRepository(db)
	groupRepositoryGORM := gorm.NewGroupRepository(db)
	// Konten
	bookRepositoryGORM := gorm.NewBookRepository(db)
	pageRepositoryGORM := gorm.NewPageRepository(db)
	interactionRepositoryGORM := gorm.NewInteractionRepository(db)
	// Progres
	progressRepositoryGORM := gorm.NewUserBookProgressRepository(db)
	// Game
	gameRepositoryGORM := gorm.NewGameRepository(db)
	userGameScoreRepositoryGORM := gorm.NewUserGameScoreRepository(db)
	// Log Aktivitas
	activityLogRepositoryGORM := gorm.NewActivityLogRepository(db)
	// Repo Media
	mediaAssetRepositoryGORM := gorm.NewMediaAssetRepository(db)
	// Repo Pengaturan Buku Grup
	groupBookSettingRepositoryGORM := gorm.NewGroupBookSettingRepository(db)

	// 6. Init Usecases
	loggerService := logger.NewActivityLoggerService(activityLogRepositoryGORM)
	sessionBlacklistService := cache.NewSessionBlacklistCACHE()
	userService := user.NewUserService(
		userRepositoryGORM,
		roleRepositoryCACHE,
		passwordService,
		tokenService,
		cfg,
		loggerService,
		sessionBlacklistService,
	)
	adminService := admin.NewAdminService(
		subjectRepositoryGORM,
		classRepositoryGORM,
		groupRepositoryGORM,
		bookRepositoryGORM,
		pageRepositoryGORM,
		interactionRepositoryGORM,
	)
	appService := app.NewAppService(
		bookRepositoryGORM,
		progressRepositoryGORM,
		gameRepositoryGORM,
		userGameScoreRepositoryGORM,
		loggerService,
		userRepositoryGORM,
		groupRepositoryGORM,
		groupBookSettingRepositoryGORM,
	)
	dashboardService := dashboard.NewDashboardService(
		activityLogRepositoryGORM,
		userRepositoryGORM,
		groupRepositoryGORM,
		bookRepositoryGORM,
		groupBookSettingRepositoryGORM,
	)
	mediaService := media.NewMediaService(
		fileStorageService,
		mediaAssetRepositoryGORM,
		cfg,
	)

	// 7. Init Handlers
	userHandler := http.NewUserHandler(userService, validate)
	adminHandler := http.NewAdminHandler(adminService, validate)
	appHandler := http.NewAppHandler(appService, dashboardService, validate)
	dashboardHandler := http.NewDashboardHandler(dashboardService, validate)
	mediaHandler := http.NewMediaHandler(mediaService, validate, loggerService)

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
		userRepositoryGORM,
		sessionBlacklistService,
		cfg,
	)

	// 10. Start Server with Graceful Shutdown

	// Channel untuk mendengarkan sinyal OS (cth: Ctrl+C, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Jalankan server di goroutine agar tidak memblokir
	go func() {
		log.Printf("Starting server on port %d...", cfg.App.Port)
		if err := app.Listen(fmt.Sprintf(":%d", cfg.App.Port)); err != nil {
			// Kita gunakan log.Printf agar shutdown bisa berjalan
			log.Printf("Server failed to start: %v", err)
			quit <- syscall.Signal(0) // Kirim sinyal dummy untuk memicu shutdown
		}
	}()

	// Blokir main goroutine sampai sinyal diterima
	<-quit

	log.Println("Received shutdown signal. Gracefully shutting down...")

	// Beri waktu 5 detik untuk server menyelesaikan request yang sedang berjalan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown Fiber server
	if err := app.Shutdown(); err != nil {
		log.Printf("Fiber server shutdown failed: %v", err)
	}

	// Shutdown Logger Service (ini akan memproses sisa batch)
	if err := loggerService.Shutdown(ctx); err != nil {
		log.Printf("Logger service shutdown failed: %v", err)
	}

	log.Println("Server gracefully shut down.")
}
