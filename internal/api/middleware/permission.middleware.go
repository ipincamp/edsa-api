package middleware

import (
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

		userRole, found := util.GetRoleByID(user.Role.ID)
		if !found {
			return dto.SendError(ctx, fiber.StatusForbidden, "Invalid user role")
		}

		for _, p := range userRole.Permissions {
			if p.Name == requiredPermission {
				return ctx.Next()
			}
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

		for _, p := range userRole.Permissions {
			if p.Name == permissionName {
				return ctx.Next()
			}
		}

		return dto.SendError(ctx, fiber.StatusForbidden, "You don't have the required role or permission")
	}
}
