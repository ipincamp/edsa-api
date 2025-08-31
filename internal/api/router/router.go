package router

import (
	"log"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/api/handler"
	"github.com/ipincamp/go-edsa-api/internal/api/middleware"
	"github.com/ipincamp/go-edsa-api/internal/constant"
)

func Setup(
	app *fiber.App,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
	permissionMiddleware *middleware.PermissionMiddleware,
	limiterMiddleware *middleware.LimiterMiddleware,
) {
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", limiterMiddleware.Make(5, 1), authHandler.Register)
	auth.Post("/login", limiterMiddleware.Make(5, 1), authHandler.Login)
	auth.Post("/logout", authMiddleware.Auth(), authHandler.Logout)

	// User routes
	user := api.Group("/users", authMiddleware.Auth())
	user.Get("/profile", userHandler.Profile)
	user.Patch("/:userID", userHandler.Update)

	// Admin routes
	admin := api.Group(
		"/admin",
		authMiddleware.Auth(),
		permissionMiddleware.CheckRole(constant.RoleAdmin.String()),
	)
	admin.Get("/users", userHandler.Index)

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the API")
	})
}

func PrintRoutes(app *fiber.App) {
	methodColors := map[string]string{
		"GET":     constant.Color("green"),
		"POST":    constant.Color("yellow"),
		"PUT":     constant.Color("blue"),
		"PATCH":   constant.Color("blue"),
		"DELETE":  constant.Color("red"),
		"HEAD":    constant.Color("cyan"),
		"OPTIONS": constant.Color("gray"),
	}
	log.Printf("├── %s%sRegistered API Routes%s", constant.Color("bold"), constant.Color("yellow"), constant.Color("reset"))

	basePrefix := "│   "
	routesMap := make(map[string][]string)
	for _, route := range app.GetRoutes(true) {
		routesMap[route.Path] = append(routesMap[route.Path], route.Method)
	}
	paths := make([]string, 0, len(routesMap))
	for path := range routesMap {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var groupNames []string
	tempSet := make(map[string]bool)
	for _, path := range paths {
		prefix := getRoutePrefix(path)
		if !tempSet[prefix] {
			tempSet[prefix] = true
			groupNames = append(groupNames, prefix)
		}
	}

	for i, groupName := range groupNames {
		isLastGroup := i == len(groupNames)-1
		groupConnector := "├──"
		routePrefix := basePrefix + "│   "
		if isLastGroup {
			groupConnector = "└──"
			routePrefix = basePrefix + "    "
		}
		log.Printf(
			"%s %s%s%s%s",
			basePrefix+groupConnector,
			constant.Color("bold"),
			constant.Color("cyan"),
			strings.ToUpper(groupName),
			constant.Color("reset"),
		)

		var groupPaths []string
		for _, path := range paths {
			if getRoutePrefix(path) == groupName {
				groupPaths = append(groupPaths, path)
			}
		}

		for j, path := range groupPaths {
			isLastRouteInGroup := j == len(groupPaths)-1
			routeConnector := "├──"
			if isLastRouteInGroup {
				routeConnector = "└──"
			}

			var coloredMethods []string
			for _, method := range routesMap[path] {
				color := methodColors[method]
				coloredMethods = append(coloredMethods, color+method+constant.Color("reset"))
			}

			log.Printf("%s [%s] %s",
				routePrefix+routeConnector,
				strings.Join(coloredMethods, ", "),
				path,
			)
		}
	}
}

func getRoutePrefix(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) > 1 {
		if parts[0] == "api" && len(parts) > 1 {
			return parts[1]
		}
		return parts[0]
	}
	return "general"
}
