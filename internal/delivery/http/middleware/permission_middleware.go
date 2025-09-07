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
// RequireRole adalah middleware untuk memeriksa apakah user memiliki salah satu role yang diizinkan
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*domain.User)
		if user == nil {
			return sendPermError(c, fiber.StatusUnauthorized, "unauthorized")
		}

		userRole, found := cache.GetRoleByID(user.RoleID)
		if !found {
			return sendPermError(c, fiber.StatusForbidden, "user role not found")
		}

		for _, allowedRole := range roles {
			if strings.EqualFold(userRole.Name, allowedRole) {
				return c.Next()
			}
		}
		return sendPermError(c, fiber.StatusForbidden, "insufficient role permissions")
	}
}

// RequirePermission - Middleware to check if user has specific permission(s)
// RequirePermission adalah middleware untuk memeriksa apakah user memiliki semua permission yang dibutuhkan
func RequirePermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*domain.User)
		if user == nil {
			return sendPermError(c, fiber.StatusUnauthorized, "unauthorized")
		}

		userRole, found := cache.GetRoleByID(user.RoleID)
		if !found {
			return sendPermError(c, fiber.StatusForbidden, "user role not found")
		}

		rolePermissions, err := userRole.GetPermissions()
		if err != nil {
			return sendPermError(c, fiber.StatusInternalServerError, "failed to parse role permissions")
		}

		for _, requiredPermission := range permissions {
			if hasPermission, exists := rolePermissions[requiredPermission]; !exists || !hasPermission {
				return sendPermError(c, fiber.StatusForbidden, "insufficient permissions")
			}
		}
		return c.Next()
	}
}

// RequireAdmin - Shorthand middleware for admin role
// RequireAdmin adalah middleware untuk membatasi akses hanya untuk admin
func RequireAdmin() fiber.Handler {
	return RequireRole(constant.RoleAdmin.String())
}

// RequireTeacher - Shorthand middleware for teacher role
// RequireTeacher adalah middleware untuk membatasi akses hanya untuk teacher
func RequireTeacher() fiber.Handler {
	return RequireRole(constant.RoleTeacher.String())
}

// RequireTeacherOrAdmin - Shorthand middleware for teacher or admin role
// RequireTeacherOrAdmin adalah middleware untuk membatasi akses hanya untuk teacher atau admin
func RequireTeacherOrAdmin() fiber.Handler {
	return RequireRole(constant.RoleTeacher.String(), constant.RoleAdmin.String())
}

// RequireOwnershipOrAdmin adalah middleware untuk membatasi akses hanya untuk owner resource atau admin
func RequireOwnershipOrAdmin(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*domain.User)
		if user == nil {
			return sendPermError(c, fiber.StatusUnauthorized, "unauthorized")
		}

		userRole, found := cache.GetRoleByID(user.RoleID)
		if !found {
			return sendPermError(c, fiber.StatusForbidden, "user role not found")
		}

		if strings.EqualFold(userRole.Name, constant.RoleAdmin.String()) {
			return c.Next()
		}

		resourceUserID := c.Params(paramName)
		if resourceUserID == user.ID {
			return c.Next()
		}

		return sendPermError(c, fiber.StatusForbidden, "access denied: you can only access your own resources")
	}
}

// sendPermError adalah helper untuk mengirim error pada permission/role middleware
func sendPermError(c *fiber.Ctx, status int, message string) error {
	return util.SendError(c, status, message)
}
