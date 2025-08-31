package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/api/handler"
	"github.com/ipincamp/go-edsa-api/internal/api/middleware"
	"github.com/ipincamp/go-edsa-api/internal/api/router"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/database/connection"
	"github.com/ipincamp/go-edsa-api/internal/database/migration"
	"github.com/ipincamp/go-edsa-api/internal/database/seeder"
	"github.com/ipincamp/go-edsa-api/internal/repository"
	"github.com/ipincamp/go-edsa-api/internal/service"
	"github.com/ipincamp/go-edsa-api/internal/util"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	rollbackFlag := flag.Bool("rollback", false, "Rollback the last database migration")
	seedFlag := flag.Bool("seed", false, "Run database seeders")
	flag.Parse()

	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	if *migrateFlag {
		migration.Migrate(dbConnection)
		return
	}
	if *rollbackFlag {
		migration.Rollback(dbConnection)
		return
	}
	if *seedFlag {
		if cnf.Env == "production" {
			log.Println("Seeder is disabled in production environment.")
			return
		}
		seeder.Seed(dbConnection, cnf)
		return
	}

	util.LoadRolesAndPermissions(dbConnection)
	roles := util.GetAllRoles()

	util.PrintPermissionTree(roles)
	log.Printf(
		"├── %s%sStarting server...%s",
		constant.Color("bold"),
		constant.Color("yellow"),
		constant.Color("reset"),
	)

	log.Printf(
		"├── %s%sLoading email bloom filter...%s",
		constant.Color("bold"),
		constant.Color("yellow"),
		constant.Color("reset"),
	)
	bloomFilter, err := util.NewBloomFilterManager(cnf.BloomFilterPath, 10000, 0.01)
	if err != nil {
		log.Fatalf("Failed to initialize bloom filter: %v", err)
	}

	userRepository := repository.NewUser(dbConnection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	allUsers, err := userRepository.FindAll(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch users to populate bloom filter: %v", err)
	}

	for _, user := range allUsers {
		bloomFilter.Add(user.Email)
	}
	if err := bloomFilter.Save(); err != nil {
		log.Printf("Warning: Failed to save initial bloom filter: %v", err)
	}
	log.Printf("│   └── %sEmail bloom filter loaded with %d entries.%s", constant.Color("green"), len(allUsers), constant.Color("reset"))

	validator := util.NewValidator()
	userService := service.NewUser(userRepository)
	authService := service.NewAuth(userRepository, dbConnection, cnf, bloomFilter)

	authMiddleware := middleware.NewAuth(cnf)
	permissionMiddleware := middleware.NewPermission()
	limiterMiddleware := middleware.NewLimiter()

	authHandler := handler.NewAuth(authService, validator)
	userHandler := handler.NewUser(userService, validator)

	app := fiber.New()
	router.Setup(
		app,
		authHandler,
		userHandler,
		authMiddleware,
		permissionMiddleware,
		limiterMiddleware,
	)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Resource not found",
		})
	})

	router.PrintRoutes(app)
	log.Printf(
		"└── %s%sServer is listening on %s:%s%s",
		constant.Color("bold"),
		constant.Color("green"),
		cnf.Server.Host,
		cnf.Server.Port,
		constant.Color("reset"),
	)
	err = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
