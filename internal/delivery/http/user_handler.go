package http

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/applogger"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type UserHandler struct {
	userService   usecase.UserService
	mediaService  usecase.MediaService
	loggerService usecase.ActivityLoggerService
	validate      *validator.GoPlaygroundValidator
}

func NewUserHandler(
	us usecase.UserService,
	ms usecase.MediaService,
	logger usecase.ActivityLoggerService,
	v *validator.GoPlaygroundValidator,
) *UserHandler {
	return &UserHandler{
		userService:   us,
		mediaService:  ms,
		loggerService: logger,
		validate:      v,
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
	// Ambil token dari header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Missing Authorization Header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid Authorization Header format")
	}

	refreshToken := parts[1]

	// Panggil Usecase
	tokenResponse, err := h.userService.RefreshToken(c.Context(), refreshToken)
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

	// Ambil email dari middleware
	userEmail, ok := c.Locals("userEmail").(string)
	if !ok || userEmail == "" {
		// Ini seharusnya tidak terjadi jika AuthMiddleware benar
		applogger.ErrorLogger.Printf("Logout Handler: userEmail not found or empty in context for userID %s", userID)
		// Tetap lanjutkan logout, tapi log tanpa email
		userEmail = "[email not found in context]"
	}

	// Panggil Usecase dengan email
	if err := h.userService.Logout(c.Context(), userID, sessionID, userEmail); err != nil { // Teruskan email
		// userService.Logout sekarang mengembalikan error jika blacklist gagal
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, "Logout failed", err.Error())
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

// UpdateAvatar menangani 'PATCH /api/v1/users/me/avatar'
func (h *UserHandler) UpdateAvatar(c *fiber.Ctx) error {
	// 1. Ambil user ID dari middleware
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid user ID in token")
	}
	sessionID, _ := c.Locals("sessionID").(uuid.UUID)

	// 2. Ambil file dari form field "picture"
	file, err := c.FormFile("picture")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Missing 'picture' file in form-data", err.Error())
	}

	// 3. Validasi ukuran (sesuai media_handler)
	const maxFileSize = 1 * 1024 * 1024 // 1 MB
	if file.Size > maxFileSize {
		return utils.SendSimpleError(c, fiber.StatusRequestEntityTooLarge, "File is too large", "File size must be no more than 1 MB")
	}

	// 4. Panggil MediaService.UploadFile
	// OwnerID adalah user_id, OwnerType adalah 'user_avatar'
	uploaderIDPtr := &userID
	asset, err := h.mediaService.UploadFile(c.Context(), file, userID.String(), domain.OwnerTypeUserAvatar, uploaderIDPtr)
	if err != nil {
		// mediaService.UploadFile sudah menangani error (mime type, dll)
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	// 5. Log aktivitas upload
	details, _ := json.Marshal(map[string]interface{}{
		"asset_id":   asset.ID,
		"file_name":  asset.FileName,
		"owner_type": domain.OwnerTypeUserAvatar,
	})
	h.loggerService.Log(c.Context(), domain.ActivityLog{
		UserID:    userID,
		SessionID: sessionID,
		Action:    domain.ActionMediaUpload,
		Details:   details,
	})

	// 6. Panggil Usecase UserService.UpdateAvatar dengan MediaID baru
	updatedUser, err := h.userService.UpdateAvatar(c.Context(), userID, asset.ID)
	if err != nil {
		if err.Error() == "media asset not found" || err.Error() == "user not found" {
			return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
		}
		if err.Error() == "you do not own this media asset" || err.Error() == "this media asset is not designated for user avatars" {
			return utils.SendSimpleError(c, fiber.StatusForbidden, err.Error(), err.Error())
		}
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	// 7. Kembalikan sukses dengan data user yang diperbarui
	return utils.SendSuccess(c, fiber.StatusOK, "Avatar updated successfully", updatedUser)
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
	// --- Get userID from context ---
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid user ID in token")
	}

	// --- Get sessionID from context ---
	sessionID, ok := c.Locals("sessionID").(uuid.UUID)
	if !ok || sessionID == uuid.Nil {
		// Seharusnya sessionID selalu ada jika AuthMiddleware berhasil
		applogger.ErrorLogger.Printf("ResendVerification Handler: sessionID not found in context for userID %s", userID)
		return utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid session", "Session identifier missing")
	}

	// Panggil service dengan userID dan sessionID
	err := h.userService.SendVerificationEmail(c.Context(), userID, sessionID)
	if err != nil {
		// Handle specific errors
		if err.Error() == "email already verified" {
			return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
		}
		if err.Error() == "user not found" {
			// This shouldn't happen if the token is valid, but handle defensively
			return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
		}
		// Other internal errors
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, "Failed to send verification email", err.Error())
	}

	// Adjusted success message
	return utils.SendSuccess(c, fiber.StatusOK, "Verification email sent successfully to your registered email address.", nil)
}

// VerifyEmailDirect menangani GET /auth/verify-email?token=...
func (h *UserHandler) VerifyEmailDirect(c *fiber.Ctx) error {
	// 1. Ambil token dari query parameter
	token := c.Query("token")
	if token == "" {
		// Tidak ada token, kirim respons error sederhana
		// Gunakan SendString untuk respons teks
		return c.Status(fiber.StatusBadRequest).SendString("Verification token missing.")
	}

	// 2. Panggil service VerifyEmail
	err := h.userService.VerifyEmail(c.Context(), token)

	// 3. Kirim respons berdasarkan hasil service
	if err == nil {
		// Sukses
		return c.Status(fiber.StatusOK).SendString("Verification successful.")
	}

	// Cek error spesifik dari service
	errMsg := err.Error()
	if errMsg == "email already verified" {
		return c.Status(fiber.StatusBadRequest).SendString("Email already verified.")
	}
	if strings.Contains(errMsg, "token invalid") || strings.Contains(errMsg, "token has expired") || errMsg == "token already used" || strings.Contains(errMsg, "not found") {
		return c.Status(fiber.StatusBadRequest).SendString("Token is no longer valid.")
	}

	// Error internal lainnya
	applogger.ErrorLogger.Printf("VerifyEmailDirect Handler: Internal error during verification: %v", err) // Log error internal
	// Berikan pesan generik ke pengguna untuk error internal
	return c.Status(fiber.StatusInternalServerError).SendString("Verification failed due to an internal error. Please try again later or contact support.")
}

// VerifyEmail menangani 'POST /auth/verify-email'
// DEPRECATED: Gunakan VerifyEmailDirect sebagai gantinya
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
