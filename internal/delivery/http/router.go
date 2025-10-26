package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/middleware"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

func SetupRoutes(
	app *fiber.App,
	userHandler *UserHandler,
	adminHandler *AdminHandler,
	appHandler *AppHandler,
	dashboardHandler *DashboardHandler,
	mediaHandler *MediaHandler,
	tokenSvc usecase.TokenService,
	userRepo usecase.UserRepository,
	blacklistSvc usecase.SessionBlacklistService,
	cfg *config.Config,
) {
	app.Use(logger.New())

	// Static files untuk media (gambar, video, dsb.)
	// Cth: GET /public/uploads/file.png akan disajikan dari ./public/uploads/file.png
	app.Static(cfg.Storage.StoragePublicURL, cfg.Storage.StoragePath, fiber.Static{
		CacheDuration: 365 * 24 * time.Hour,
		Compress:      true, // Mengaktifkan kompresi gzip/brotli
		ByteRange:     true, // Mengizinkan browser me-resume download
	}) // PASSED

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	}) // PASSED

	// Favicon
	app.Static("/favicon.ico", cfg.Storage.StoragePath+"/favicon.svg", fiber.Static{
		CacheDuration: 365 * 24 * time.Hour,
		Compress:      true,
		ByteRange:     true,
	}) // PASSED

	// Grup API v1
	api := app.Group("/api/v1")

	// --- Rute Autentikasi ---

	// Rute Auth (tidak dilindungi)
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)    // PASSED
	auth.Post("/login", userHandler.Login)          // PASSED
	auth.Post("/refresh", userHandler.RefreshToken) // PASSED
	auth.Post(
		"/logout",
		middleware.AuthMiddleware(tokenSvc, blacklistSvc),
		userHandler.Logout,
	) // PASSED
	auth.Post("/resend-verification", userHandler.ResendVerification) // PASSED
	auth.Post("/verify-email", userHandler.VerifyEmail)               // PASSED
	auth.Post("/forgot-password", userHandler.ForgotPassword)         // PASSED
	auth.Post("/reset-password", userHandler.ResetPassword)           // PASSED

	// Rute User (dilindungi)
	protected := api.Group("/users")
	protected.Use(middleware.AuthMiddleware(tokenSvc, blacklistSvc))
	protected.Get("/me", userHandler.GetMe)                       // PASSED
	protected.Patch("/me/password", userHandler.ChangePassword)   // PASSED
	protected.Patch("/me/details", userHandler.UpdateUserDetails) // PASSED
	protected.Patch("/me/avatar", userHandler.UpdateAvatar)       //PASSED
	protected.Delete("/me", userHandler.DeleteAccount)            // PASSED

	// --- Rute Administrasi ---

	// Rute Admin (dilindungi & hanya untuk admin)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(tokenSvc, blacklistSvc))
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

	// --- Rute Manajemen Konten ---

	// Rute Books
	books := admin.Group("/books")
	books.Post("/", adminHandler.CreateBook)          // PASSED
	books.Get("/", adminHandler.GetAllBooks)          // PASSED
	books.Get("/:bookId", adminHandler.GetBookByID)   // PASSED
	books.Patch("/:bookId", adminHandler.UpdateBook)  // PASSED
	books.Delete("/:bookId", adminHandler.DeleteBook) // PASSED

	// Rute Pages (nested under books)
	books.Post("/:bookId/pages", adminHandler.CreatePage)
	books.Get("/:bookId/pages", adminHandler.GetAllPagesForBook)

	// Rute Pages (standalone for single resource mgmt)
	pages := admin.Group("/pages")
	pages.Get("/:pageId", adminHandler.GetPageByID)
	pages.Put("/:pageId", adminHandler.UpdatePage)
	pages.Delete("/:pageId", adminHandler.DeletePage)

	// Rute Interactions (nested under pages)
	pages.Post("/:pageId/interactions", adminHandler.CreateInteraction)
	pages.Get("/:pageId/interactions", adminHandler.GetAllInteractionsForPage)

	// Rute Interactions (standalone for single resource mgmt)
	interactions := admin.Group("/interactions")
	interactions.Get("/:interactionId", adminHandler.GetInteractionByID) // Implied
	interactions.Put("/:interactionId", adminHandler.UpdateInteraction)
	interactions.Delete("/:interactionId", adminHandler.DeleteInteraction)

	// --- Rute Manajemen Progres & Game ---

	// Rute Aplikasi (Siswa & Guru)
	appRoutes := api.Group("/app")
	appRoutes.Use(middleware.AuthMiddleware(tokenSvc, blacklistSvc))

	// Rute Modul "Read"
	appRoutes.Get("/books", appHandler.GetBooksWithProgress)                 // PASSED
	appRoutes.Get("/books/:bookId/restore", appHandler.GetProgressToRestore) // PASSED
	appRoutes.Post("/progress/update", appHandler.UpdatePageProgress)        // PASSED
	appRoutes.Post("/progress/complete", appHandler.CompleteBookProgress)    // PASSED

	// Rute Modul "Game"
	appRoutes.Get("/games", appHandler.GetAllGames)                    // PASSED
	appRoutes.Post("/games/:gameId/score", appHandler.SubmitGameScore) // PASSED

	// Rute Aktivitas Pengguna
	appRoutes.Get("/activity", appHandler.GetMyActivity) // PASSED

	// --- Rute Dashboard (Guru) ---

	// Rute Dashboard Guru
	dashboard := api.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware(tokenSvc, blacklistSvc))
	dashboard.Use(middleware.TeacherMiddleware(userRepo))

	// Rute Laporan Aktivitas
	dashboard.Get("/students/:studentId/activity", dashboardHandler.GetStudentActivity) // PASSED

	// Rute Manajemen Buku Grup
	dashboard.Post("/groups/:groupId/books/:bookId/unlock", dashboardHandler.UnlockBookForGroup)
	// TODO: Rute dashboard lainnya

	// --- Rute Media ---

	// Rute Media
	media := api.Group("/media")
	media.Use(middleware.AuthMiddleware(tokenSvc, blacklistSvc))

	// Rute Upload File
	media.Post("/upload", mediaHandler.UploadFile) // PASSED
	media.Delete("/:id", mediaHandler.DeleteFile)  // PASSED

	// TODO: Rute manajemen penugasan (Assign/Unassign)
	// TODO: Rute untuk Speaking

	// Handle 404 - Not Found
	app.Use(func(c *fiber.Ctx) error {
		return utils.SendSimpleError(
			c,
			fiber.StatusNotFound,
			"Resource Not Found",
			"Route '"+c.Path()+"' not found on this server",
		)
	}) // PASSED
}
