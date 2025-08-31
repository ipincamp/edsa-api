package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain"
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

	roles, err := middleware.LoadAndCachePermissions(dbConnection)
	if err != nil {
		log.Fatalf("Failed to load and cache permissions: %v", err)
	}

	log.Printf(
		"%s%sStarting server...%s",
		constant.Color("bold"),
		constant.Color("yellow"),
		constant.Color("reset"),
	)
	printPermissionTree(roles)

	app := fiber.New()
	validator := util.NewValidator()

	userRepository := repository.NewUser(dbConnection)
	userService := service.NewUser(userRepository)
	authService := service.NewAuth(userRepository, dbConnection, cnf)

	authMiddleware := middleware.NewAuth(cnf)
	permissionMiddleware := middleware.NewPermission()

	authHandler := handler.NewAuth(authService, validator)
	userHandler := handler.NewUser(userService, validator)

	router.Setup(app, authHandler, userHandler, authMiddleware, permissionMiddleware)
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

func printPermissionTree(roles []domain.Role) {
	log.Printf(
		"├── %s%sCaching roles and permissions...%s",
		constant.Color("bold"),
		constant.Color("yellow"),
		constant.Color("reset"),
	)

	basePrefix := "│   "
	for i, role := range roles {
		isLastRole := i == len(roles)-1
		roleConnector := "├──"
		if isLastRole {
			roleConnector = "└──"
		}

		log.Printf(
			"%s%s%s %s%s%s%s",
			basePrefix, roleConnector, constant.Color("reset"),
			constant.Color("bold"), constant.Color("yellow"), role.Name, constant.Color("reset"),
		)

		permParentPrefix := basePrefix + "│   "
		if isLastRole {
			permParentPrefix = basePrefix + "    "
		}

		for j, p := range role.Permissions {
			isLastPerm := j == len(role.Permissions)-1
			permConnector := "├──"
			if isLastPerm {
				permConnector = "└──"
			}
			log.Printf(
				"%s%s%s %s%s%s%s",
				permParentPrefix, constant.Color("gray"), permConnector, constant.Color("reset"),
				constant.Color("green"), p.Name, constant.Color("reset"),
			)
		}
	}

	log.Printf(
		"├── %s%sCached permissions for %d roles.%s",
		constant.Color("bold"), constant.Color("green"), len(roles), constant.Color("reset"),
	)
}
