package middleware

import (
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/dto"
	"gorm.io/gorm"
)

type PermissionMiddleware struct{}

var permissionCache = make(map[string]map[string]bool)
var roleIdToNameCache = make(map[string]string)

func NewPermission() *PermissionMiddleware {
	return &PermissionMiddleware{}
}

func LoadAndCachePermissions(db *gorm.DB) ([]domain.Role, error) {
	var roles []domain.Role
	if err := db.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}

	permissionCache = make(map[string]map[string]bool)
	roleIdToNameCache = make(map[string]string)

	for _, role := range roles {
		roleIdToNameCache[role.ID] = role.Name
		perms := make(map[string]bool)
		for _, p := range role.Permissions {
			perms[p.Name] = true
		}
		permissionCache[role.Name] = perms
	}

	return roles, nil
}

func (m *PermissionMiddleware) CheckRole(requiredRoles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(domain.User)
		if !ok {
			return dto.SendError(ctx, fiber.StatusForbidden, "User data not found in context")
		}

		if slices.Contains(requiredRoles, user.Role.Name) {
			return ctx.Next()
		}

		return dto.SendError(ctx, fiber.StatusForbidden, "You don't have the required role")
	}
}

func (m *PermissionMiddleware) CheckPermission(requiredPermission string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(domain.User)
		if !ok {
			return dto.SendError(ctx, fiber.StatusForbidden, "User data not found in context")
		}

		userRoleID := user.Role.ID
		userRoleName, ok := roleIdToNameCache[userRoleID]
		if !ok {
			return dto.SendError(ctx, fiber.StatusForbidden, "Invalid user role")
		}

		if rolesPermissions, ok := permissionCache[userRoleName]; ok && rolesPermissions[requiredPermission] {
			return ctx.Next()
		}

		return dto.SendError(ctx, fiber.StatusForbidden, "You don't have the required permission")
	}
}

func (m *PermissionMiddleware) CheckRoleOrPermission(roleName string, permissionName string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(domain.User)
		if !ok {
			return dto.SendError(ctx, fiber.StatusForbidden, "User data not found in context")
		}

		userRoleName, ok := roleIdToNameCache[user.Role.ID]
		if !ok {
			return dto.SendError(ctx, fiber.StatusForbidden, "Invalid user role")
		}

		if userRoleName == roleName {
			return ctx.Next()
		}

		if rolesPermissions, ok := permissionCache[userRoleName]; ok && rolesPermissions[permissionName] {
			return ctx.Next()
		}

		return dto.SendError(ctx, fiber.StatusForbidden, "You don't have the required role or permission")
	}
}
