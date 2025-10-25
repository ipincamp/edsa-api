package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type UserHandler struct {
	userService usecase.UserService
	validate    *validator.GoPlaygroundValidator
}

func NewUserHandler(us usecase.UserService, v *validator.GoPlaygroundValidator) *UserHandler {
	return &UserHandler{
		userService: us,
		validate:    v,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest

	// Parse & Validasi
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// Panggil Usecase
	authResponse, err := h.userService.Register(c.Context(), &req)
	if err != nil {
		// Seharusnya ada error handling yang lebih baik di sini
		return utils.SendSimpleError(c, fiber.StatusConflict, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "User registered successfully", authResponse)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest

	// Parse & Validasi
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// Panggil Usecase
	authResponse, err := h.userService.Login(c.Context(), &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Login successful", authResponse)
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var req domain.RefreshTokenRequest

	// Parse & Validasi
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// Panggil Usecase
	tokenResponse, err := h.userService.RefreshToken(c.Context(), &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Token refreshed successfully", tokenResponse)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// Ambil user ID dari middleware
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid user ID in token")
	}

	// Ambil session ID dari middleware
	sessionID, ok := c.Locals("sessionID").(uuid.UUID)
	if !ok {
		// Jika sessionID tidak ada, ini adalah error krusial untuk logout
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid session ID in token")
	}

	// Panggil Usecase
	if err := h.userService.Logout(c.Context(), userID, sessionID); err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Logged out successfully", nil)
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	// Ambil user ID dari middleware
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid user ID in token")
	}

	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User profile retrieved successfully", user)
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	// 1. Ambil user ID dari middleware
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid user ID in token")
	}

	var req domain.ChangePasswordRequest

	// 2. Parse & Validasi
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// 3. Panggil Usecase
	if err := h.userService.ChangePassword(c.Context(), userID, &req); err != nil {
		if err.Error() == "invalid current password" {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, err.Error(), err.Error())
		}
		if err.Error() == "user not found" {
			return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
		}
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	// 4. Kembalikan sukses
	return utils.SendSuccess(c, fiber.StatusOK, "Password changed successfully", nil)
}

func (h *UserHandler) UpdateUserDetails(c *fiber.Ctx) error {
	// 1. Ambil user ID dari middleware
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid user ID in token")
	}

	var req domain.UpdateDetailsRequest

	// 2. Parse & Validasi
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// 3. Panggil Usecase
	updatedUser, err := h.userService.UpdateUserDetails(c.Context(), userID, &req)
	if err != nil {
		if err.Error() == "invalid password confirmation" {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, err.Error(), err.Error())
		}
		if err.Error() == "user not found" {
			return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
		}
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	// 4. Kembalikan sukses dengan data user yang diperbarui
	return utils.SendSuccess(c, fiber.StatusOK, "User details updated successfully", updatedUser)
}

// DeleteAccount menangani 'DELETE /api/v1/users/me'
// Endpoint ini memiliki dua status:
// 1. Jika 'confirmation_token' tidak ada: Memulai proses, mengirim email.
// 2. Jika 'confirmation_token' ada: Mengkonfirmasi proses, menghapus akun.
func (h *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	// Ambil user ID dari middleware
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid user ID in token")
	}

	// 1. Coba parse DTO "parsial" (dengan field omitempty)
	var req domain.DeletionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// 2. Periksa apakah ini langkah pertama (meminta token)
	if req.ConfirmationToken == "" {
		if err := h.userService.RequestAccountDeletion(c.Context(), userID); err != nil {
			// Jika gagal (cth: user tidak ada, email gagal kirim)
			return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
		}
		// Beri tahu klien bahwa email telah dikirim
		return utils.SendSimpleError(c,
			fiber.StatusUnprocessableEntity,
			"Confirmation required",
			"A confirmation email has been sent to you. Please provide the token from the email to delete your account.",
		)
	}

	// 3. Jika token ada, ini adalah langkah kedua (konfirmasi)
	// Kita parse ulang body menggunakan DTO "ketat" (dengan validasi "required")
	var confirmReq domain.ConfirmDeletionRequest
	if err := c.BodyParser(&confirmReq); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	// Lakukan validasi ketat
	if errs := h.validate.ValidateStruct(confirmReq); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// 4. Panggil usecase konfirmasi
	if err := h.userService.ConfirmAccountDeletion(c.Context(), userID, &confirmReq); err != nil {
		// Tangani error spesifik (cth: password salah, token expired)
		if err.Error() == "invalid current password" {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, err.Error(), err.Error())
		}
		if err.Error() == "token has expired" || strings.Contains(err.Error(), "invalid token") {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, err.Error(), err.Error())
		}
		if err.Error() == "confirmation token does not match authenticated user" {
			return utils.SendSimpleError(c, fiber.StatusForbidden, err.Error(), err.Error())
		}
		// Error umum
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	// 5. Sukses
	return utils.SendSuccess(c, fiber.StatusOK, "Account deleted successfully", nil)
}

// ResendVerification menangani 'POST /auth/resend-verification'
func (h *UserHandler) ResendVerification(c *fiber.Ctx) error {
	var req domain.ResendVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	err := h.userService.SendVerificationEmail(c.Context(), req.Email)
	if err != nil {
		// Handle error spesifik
		if err.Error() == "email already verified" {
			return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
		}
		// Jangan ekspos error "user not found"
		if err.Error() == "user not found" {
			return utils.SendSuccess(c, fiber.StatusOK, "If your email is registered and not verified, a verification link has been sent.", nil)
		}
		// Error internal lainnya
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Verification email sent successfully.", nil)
}

// VerifyEmail menangani 'POST /auth/verify-email'
func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	var req domain.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	err := h.userService.VerifyEmail(c.Context(), req.Token)
	if err != nil {
		if strings.Contains(err.Error(), "invalid or expired token") {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, err.Error(), err.Error())
		}
		if err.Error() == "email already verified" {
			return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
		}
		if err.Error() == "user associated with token not found" {
			return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
		}
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Email verified successfully.", nil)
}

// ForgotPassword menangani 'POST /auth/forgot-password'
func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	var req domain.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	// Panggil usecase. Jangan ekspos error internal ke user.
	_ = h.userService.SendPasswordResetEmail(c.Context(), req.Email)
	// Selalu kembalikan respons sukses untuk mencegah email enumeration
	return utils.SendSuccess(c, fiber.StatusOK, "If your email is registered, a password reset link has been sent.", nil)
}

// ResetPassword menangani 'POST /auth/reset-password'
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	var req domain.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	err := h.userService.ResetPassword(c.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid or expired token") {
			return utils.SendSimpleError(c, fiber.StatusUnauthorized, err.Error(), err.Error())
		}
		if strings.Contains(err.Error(), "not found or inactive") {
			return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
		}
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Password reset successfully.", nil)
}
