package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

// TeacherMiddleware memeriksa apakah pengguna yang diautentikasi adalah guru
// Middleware ini HARUS dijalankan SETELAH AuthMiddleware
func TeacherMiddleware(userRepo usecase.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil userID dari context
		userID, ok := c.Locals("userID").(uuid.UUID)
		if !ok || userID == uuid.Nil {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid session", "Authentication context not found")
		}

		// 2. Dapatkan data pengguna dari repository
		user, err := userRepo.FindByID(c.Context(), userID)
		if err != nil {
			return utils.SendSimpleError(c, fiber.StatusInternalServerError, "Session error", err.Error())
		}
		if user == nil {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid session", "User not found")
		}

		// 3. Periksa peran pengguna
		if user.Role.Name != domain.RoleNameTeacher {
			return utils.SendSimpleError(c, fiber.StatusForbidden, "Forbidden", "Access restricted to teachers only")
		}

		// 4. Lolos, lanjutkan ke handler berikutnya
		return c.Next()
	}
}
