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
) {
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authMiddleware.Auth(), authHandler.Logout)

	// User routes
	user := api.Group("/users", authMiddleware.Auth())
	user.Get("/me", userHandler.Profile)
	user.Patch("/:userID", userHandler.Update)

	// Admin routes
	admin := api.Group(
		"/admin",
		authMiddleware.Auth(),
		permissionMiddleware.CheckRole(constant.RoleAdmin.String()),
	)
	admin.Get("/users", userHandler.Index)
}

func PrintRoutes(app *fiber.App) {
	methodColors := map[string]string{
		"GET":     constant.Color("green"),
		"POST":    constant.Color("blue"),
		"PUT":     constant.Color("yellow"),
		"PATCH":   constant.Color("yellow"),
		"DELETE":  constant.Color("red"),
		"HEAD":    constant.Color("cyan"),
		"OPTIONS": constant.Color("gray"),
	}

	// Cetak root dari pohon rute
	log.Printf("├── %s%sRegistered API Routes%s", constant.Color("bold"), constant.Color("yellow"), constant.Color("reset"))

	// Awalan dasar untuk semua item di bawah "Registered Routes"
	basePrefix := "│   "

	// Logika untuk mengelompokkan rute tetap sama
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

	// Iterasi melalui setiap grup rute
	for i, groupName := range groupNames {
		isLastGroup := i == len(groupNames)-1
		groupConnector := "├──"
		routePrefix := basePrefix + "│   " // Awalan untuk rute di dalam grup
		if isLastGroup {
			groupConnector = "└──"
			routePrefix = basePrefix + "    " // Gunakan spasi jika ini grup terakhir
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

		// Iterasi melalui setiap rute di dalam grup
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
		// Grup berdasarkan bagian kedua setelah "/api"
		if parts[0] == "api" && len(parts) > 1 {
			return parts[1]
		}
		return parts[0]
	}
	return "general"
}
