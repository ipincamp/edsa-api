package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

// AdminMiddleware memeriksa apakah pengguna yang diautentikasi adalah admin
// Middleware ini HARUS dijalankan SETELAH AuthMiddleware
func AdminMiddleware(userRepo usecase.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil userID dari context (ditaruh oleh AuthMiddleware)
		userID, ok := c.Locals("userID").(uuid.UUID)
		if !ok || userID == uuid.Nil {
			// Ini seharusnya tidak terjadi jika AuthMiddleware berjalan dulu
			return utils.SendError(c, fiber.StatusUnauthorized, "Authentication context not found")
		}

		// 2. Dapatkan data pengguna dari repository
		user, err := userRepo.FindByID(c.Context(), userID)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve user data")
		}
		if user == nil {
			return utils.SendError(c, fiber.StatusUnauthorized, "User not found")
		}

		// 3. Periksa peran pengguna
		if user.Role.Name != domain.RoleNameAdmin {
			return utils.SendError(c, fiber.StatusForbidden, "Forbidden: Access restricted to administrators")
		}

		// 4. Lolos, lanjutkan ke handler berikutnya
		return c.Next()
	}
}
