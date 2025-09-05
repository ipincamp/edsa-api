package middleware

import (
	"encoding/json"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/util"
)

type PermissionMiddleware struct{}

func NewPermission() *PermissionMiddleware {
	return &PermissionMiddleware{}
}

func (m *PermissionMiddleware) CheckRole(requiredRoles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(domain.User)
		if !ok {
			return dto.SendError(ctx, fiber.StatusForbidden, "User data not found in context")
		}

		userRole, found := util.GetRoleByID(user.Role.ID)
		if !found {
			return dto.SendError(ctx, fiber.StatusForbidden, "Invalid user role")
		}

		if slices.Contains(requiredRoles, userRole.Name) {
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

		userRole, found := util.GetRoleByID(user.Role.ID)
		if !found {
			return dto.SendError(ctx, fiber.StatusForbidden, "Invalid user role")
		}

		var permissions map[string]bool
		if err := json.Unmarshal(userRole.Permissions, &permissions); err != nil {
			return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to parse user permissions")
		}

		if hasPermission, ok := permissions[requiredPermission]; ok && hasPermission {
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

		userRole, found := util.GetRoleByID(user.Role.ID)
		if !found {
			return dto.SendError(ctx, fiber.StatusForbidden, "Invalid user role")
		}

		if userRole.Name == roleName {
			return ctx.Next()
		}

		var permissions map[string]bool
		if err := json.Unmarshal(userRole.Permissions, &permissions); err != nil {
			return dto.SendError(ctx, fiber.StatusInternalServerError, "Failed to parse user permissions")
		}

		if hasPermission, ok := permissions[permissionName]; ok && hasPermission {
			return ctx.Next()
		}

		return dto.SendError(ctx, fiber.StatusForbidden, "You don't have the required role or permission")
	}
}
