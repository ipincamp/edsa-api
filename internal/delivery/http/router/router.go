package router

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/handler"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/middleware"
	"github.com/ipincamp/go-edsa-api/internal/repository"
	"github.com/ipincamp/go-edsa-api/internal/service"
	"github.com/ipincamp/go-edsa-api/internal/util"
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
	courseRepo := repository.NewCourseRepository(db)
	courseGroupRepo := repository.NewCourseGroupRepository(db)
	joinRequestRepo := repository.NewJoinRequestRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	bookRepo := repository.NewBookRepository(db)
	userProgressRepo := repository.NewUserProgressRepository(db)
	interactionRepo := repository.NewInteractionRepository(db)
	pageRepo := repository.NewPageRepository(db)

	// Populate Bloom Filter on startup
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		emails, err := userRepo.FindAllEmails(ctx)
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
	courseService := service.NewCourseService(
		db,
		courseRepo,
		courseGroupRepo,
		joinRequestRepo,
		enrollmentRepo,
		userRepo,
		roleRepo,
	)
	learningService := service.NewLearningService(
		db,
		bookRepo,
		userProgressRepo,
		interactionRepo,
		pageRepo,
	)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	courseHandler := handler.NewCourseHandler(courseService)
	learningHandler := handler.NewLearningHandler(learningService)

	// Middleware
	authRequired := middleware.AuthMiddleware(tokenMaker)
	teacherOnly := middleware.RequireRole(constant.RoleTeacher.String())
	guestOnly := middleware.RequireRole(constant.RoleGuest.String())

	// Add custom error handler for Method Not Allowed
	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			if e, ok := err.(*fiber.Error); ok && e.Code == fiber.StatusMethodNotAllowed {
				return util.SendError(c, fiber.StatusMethodNotAllowed, "The requested method is not allowed for this resource")
			}
		}
		return err
	})

	// Routes
	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the API")
	})
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)           // DONE
	auth.Post("/login", authHandler.Login)                 // DONE
	auth.Post("/refresh", authHandler.RefreshToken)        // DONE
	auth.Post("/logout", authRequired, authHandler.Logout) // DONE

	// Protected user routes
	users := v1.Group("/users").Use(authRequired)
	users.Get("/profile", userHandler.GetProfile)      // DONE
	users.Patch("/profile", userHandler.UpdateProfile) // DONE

	// users.Post("/", middleware.RequirePermission(constant.UsersCreate), userHandler.CreateUser) // Uncomment jika diperlukan
	users.Get("/", middleware.RequirePermission(constant.UsersListAll), userHandler.GetAllUsers)                                      // DONE
	users.Get("/:userId", middleware.RequirePermissionOrOwnership(constant.UsersViewOther, "userId"), userHandler.GetUserByID)        // DONE
	users.Patch("/:userId", middleware.RequirePermissionOrOwnership(constant.UsersUpdateOther, "userId"), userHandler.UpdateUserByID) // DONE
	// users.Delete("/:userId", middleware.RequirePermission(constant.UsersDelete), userHandler.DeleteUser) // Uncomment jika diperlukan

	// Course Management & Learning Routes
	classes := v1.Group("/classes").Use(authRequired)
	classes.Post("/join", guestOnly, courseHandler.ApplyToJoinGroup)                 // Guest mendaftar kelas
	classes.Get("/my", teacherOnly, courseHandler.GetMyClasses)                      // Guru melihat kelasnya
	classes.Get("/:groupID/students", teacherOnly, courseHandler.GetStudentsByClass) // Guru melihat murid di kelasnya

	joinRequests := v1.Group("/join-requests").Use(authRequired, teacherOnly)
	joinRequests.Post("/:requestID/handle", courseHandler.HandleJoinRequest) // Guru memproses permintaan

	learning := v1.Group("/learning").Use(authRequired)
	learning.Get("/books", learningHandler.GetAvailableBooks)
	learning.Get("/books/:bookID", learningHandler.GetBookDetail)
	learning.Post("/interactions/submit", learningHandler.SubmitInteraction)

	// Admin routes
	admin := v1.Group("/admin", authRequired, middleware.RequireRole(constant.RoleAdmin.String()))
	admin.Get("/users", userHandler.GetAllUsers)

	// Course management for admin
	admin.Post("/courses", courseHandler.CreateCourse)            // DONE
	admin.Get("/courses", courseHandler.GetAllCourses)            // DONE
	admin.Get("/courses/:courseId", courseHandler.GetCourseByID)  // DONE
	admin.Patch("/courses/:courseId", courseHandler.UpdateCourse) // DONE
	admin.Delete("/courses/:courseId", courseHandler.DeleteCourse)

	// Teacher routes
	teacher := v1.Group("/teacher", authRequired, middleware.RequireRole(constant.RoleTeacher.String()))
	teacher.Post("/handle-request/:requestID", courseHandler.HandleJoinRequest)

	// Fallback routes
	api.Use("*", func(c *fiber.Ctx) error {
		return util.SendError(c, fiber.StatusNotFound, "API endpoint not found")
	})
	app.Use("*", func(c *fiber.Ctx) error {
		return util.SendError(c, fiber.StatusNotFound, "Endpoint not found")
	})
}
