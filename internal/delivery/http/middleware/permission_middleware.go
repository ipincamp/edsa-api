package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"github.com/ipincamp/go-edsa-api/pkg/cache"
)

// RequireRole - Middleware to check if user has specific role(s)
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*domain.User)
		if user == nil {
			return util.SendError(c, fiber.StatusUnauthorized, "unauthorized")
		}

		// Get user's role from cache to ensure we have the latest data
		userRole, found := cache.GetRoleByID(user.RoleID)
		if !found {
			return util.SendError(c, fiber.StatusForbidden, "user role not found")
		}

		// Check if user's role is in the allowed roles
		for _, allowedRole := range roles {
			if strings.EqualFold(userRole.Name, allowedRole) {
				return c.Next()
			}
		}

		return util.SendError(c, fiber.StatusForbidden, "insufficient role permissions")
	}
}

// RequirePermission - Middleware to check if user has specific permission(s)
func RequirePermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*domain.User)
		if user == nil {
			return util.SendError(c, fiber.StatusUnauthorized, "unauthorized")
		}

		// Get user's role from cache to get permissions
		userRole, found := cache.GetRoleByID(user.RoleID)
		if !found {
			return util.SendError(c, fiber.StatusForbidden, "user role not found")
		}

		// Parse role permissions from JSON
		rolePermissions, err := userRole.GetPermissions()
		if err != nil {
			return util.SendError(c, fiber.StatusInternalServerError, "failed to parse role permissions")
		}

		// Check if user has all required permissions
		for _, requiredPermission := range permissions {
			if hasPermission, exists := rolePermissions[requiredPermission]; !exists || !hasPermission {
				return util.SendError(c, fiber.StatusForbidden, "insufficient permissions")
			}
		}

		return c.Next()
	}
}

// RequireAdmin - Shorthand middleware for admin role
func RequireAdmin() fiber.Handler {
	return RequireRole(constant.RoleAdmin.String())
}

// RequireTeacher - Shorthand middleware for teacher role
func RequireTeacher() fiber.Handler {
	return RequireRole(constant.RoleTeacher.String())
}

// RequireTeacherOrAdmin - Shorthand middleware for teacher or admin role
func RequireTeacherOrAdmin() fiber.Handler {
	return RequireRole(constant.RoleTeacher.String(), constant.RoleAdmin.String())
}

// RequireOwnershipOrAdmin - Middleware to check if user owns the resource or is admin
// This checks if the userId parameter matches the authenticated user's ID or if user is admin
func RequireOwnershipOrAdmin(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*domain.User)
		if user == nil {
			return util.SendError(c, fiber.StatusUnauthorized, "unauthorized")
		}

		// Get user's role from cache
		userRole, found := cache.GetRoleByID(user.RoleID)
		if !found {
			return util.SendError(c, fiber.StatusForbidden, "user role not found")
		}

		// If user is admin, allow access
		if strings.EqualFold(userRole.Name, constant.RoleAdmin.String()) {
			return c.Next()
		}

		// Check if user owns the resource
		resourceUserID := c.Params(paramName)
		if resourceUserID == user.ID {
			return c.Next()
		}

		return util.SendError(c, fiber.StatusForbidden, "access denied: you can only access your own resources")
	}
}
