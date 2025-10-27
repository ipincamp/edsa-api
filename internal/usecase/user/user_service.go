package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/applogger"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type userService struct {
	userRepo     usecase.UserRepository
	roleRepo     usecase.RoleRepository
	passSvc      usecase.PasswordService
	tokenSvc     usecase.TokenService
	cfg          *config.Config
	logger       usecase.ActivityLoggerService
	blacklistSvc usecase.SessionBlacklistService
	emailSvc     usecase.EmailService
	progressRepo usecase.UserBookProgressRepository
	mediaRepo    usecase.MediaAssetRepository
}

func NewUserService(
	userRepo usecase.UserRepository,
	roleRepo usecase.RoleRepository,
	passSvc usecase.PasswordService,
	tokenSvc usecase.TokenService,
	cfg *config.Config,
	logger usecase.ActivityLoggerService,
	blacklistSvc usecase.SessionBlacklistService,
	emailSvc usecase.EmailService,
	progressRepo usecase.UserBookProgressRepository,
	mediaRepo usecase.MediaAssetRepository,
) usecase.UserService {
	return &userService{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		passSvc:      passSvc,
		tokenSvc:     tokenSvc,
		cfg:          cfg,
		logger:       logger,
		blacklistSvc: blacklistSvc,
		emailSvc:     emailSvc,
		progressRepo: progressRepo,
		mediaRepo:    mediaRepo,
	}
}

// --- Helper Mapper ---
func (s *userService) toUserResponse(user *domain.User) *domain.UserResponse {
	var avatarURL string

	// 1. Cek apakah ada avatar KUSTOM yang terpasang
	if user.ProfilePictureID != nil && user.ProfilePicture.ID != uuid.Nil {
		avatarURL = user.ProfilePicture.PublicURL
	} else {
		// 2. Jika TIDAK ada avatar kustom, CARI URL avatar default 'avatar1.png' DARI DB
		defaultAvatarFileName := "avatar1.png"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // Timeout pendek
		defer cancel()

		defaultAsset, err := s.mediaRepo.FindByFileName(ctx, defaultAvatarFileName)
		if err == nil && defaultAsset != nil {
			avatarURL = defaultAsset.PublicURL // Gunakan URL dari DB
		} else {
			// Fallback jika asset default tidak ditemukan di DB atau error query
			if err != nil {
				applogger.ErrorLogger.Printf("toUserResponse: DB error finding default avatar '%s': %v", defaultAvatarFileName, err)
			} else { // err == nil && defaultAsset == nil
				applogger.ErrorLogger.Printf("toUserResponse: CRITICAL! Default avatar asset '%s' not found in DB.", defaultAvatarFileName)
			}
			// Gunakan fallback URL hardcoded
			avatarURL = s.cfg.Storage.StoragePublicBaseURL + "/cdn/" + defaultAvatarFileName // Fallback
		}
	}

	return &domain.UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		RoleName:          user.Role.Name,
		JoinedAt:          user.CreatedAt,
		ProfilePictureURL: avatarURL,
		IsActive:          user.IsActive,
		EmailVerified:     user.EmailVerifiedAt != nil,
		// OverallScore dihitung terpisah
	}
}

// --- Auth Services --
func (s *userService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// 1. Cek apakah email sudah ada
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		applogger.ErrorLogger.Printf("Register: database error checking email %s: %v", req.Email, err)
		return nil, errors.New("database error")
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// 2. Hash password
	hashedPassword, err := s.passSvc.Hash(req.Password)
	if err != nil {
		applogger.ErrorLogger.Printf("Register: failed to hash password: %v", err)
		return nil, errors.New("failed to hash password")
	}

	// 3. Dapatkan role default (public)
	defaultRole, err := s.roleRepo.FindByName(ctx, domain.RoleNamePublic)
	if err != nil {
		applogger.ErrorLogger.Printf("Register: database error fetching role %s: %v", domain.RoleNamePublic, err)
		return nil, errors.New("database error while fetching role")
	}
	if defaultRole == nil {
		// Log error ini karena seharusnya role default selalu ada setelah seeder
		applogger.ErrorLogger.Printf("Register: CRITICAL - Default role '%s' not found in database!", domain.RoleNamePublic)
		return nil, errors.New("default role configuration error")
	}

	var defaultAvatarID *uuid.UUID // Gunakan pointer karena bisa jadi nil
	defaultAvatarFileName := "avatar1.png"
	// Gunakan context baru dengan timeout pendek agar tidak memblokir registrasi terlalu lama
	ctxForAvatar, cancelAvatar := context.WithTimeout(ctx, 3*time.Second)
	defer cancelAvatar() // Pastikan context dibatalkan

	defaultAsset, errAsset := s.mediaRepo.FindByFileName(ctxForAvatar, defaultAvatarFileName)

	if errAsset != nil {
		// Jika ada error saat mencari avatar, log sebagai warning tapi JANGAN gagalkan registrasi
		applogger.ErrorLogger.Printf("Register: WARNING - DB error finding default avatar '%s': %v. User will have NULL avatar ID.", defaultAvatarFileName, errAsset)
	} else if defaultAsset == nil {
		// Jika avatar default tidak ditemukan di DB (seeder mungkin gagal?), log sebagai warning
		applogger.ErrorLogger.Printf("Register: WARNING - Default avatar asset '%s' not found in DB (Check Seeders?). User will have NULL avatar ID.", defaultAvatarFileName)
	} else {
		// Jika avatar ditemukan, gunakan ID-nya
		defaultAvatarID = &defaultAsset.ID
		applogger.ErrorLogger.Printf("Register: INFO - Found default avatar '%s' with ID: %s", defaultAvatarFileName, defaultAsset.ID.String()) // Log info (opsional)
	}

	// 4. Buat domain user baru
	user := &domain.User{
		Name:             req.Name,
		Email:            req.Email,
		Password:         hashedPassword,
		RoleID:           defaultRole.ID,
		ProfilePictureID: defaultAvatarID,
		IsActive:         true,
		// EmailVerifiedAt akan otomatis nil
	}

	// 5. Simpan ke database
	if err := s.userRepo.Create(ctx, user); err != nil {
		applogger.ErrorLogger.Printf("Register: failed to create user %s: %v", user.Email, err)
		return nil, errors.New("failed to create user")
	}
	// Setelah Create berhasil, user.ID sudah terisi
	user.Role = *defaultRole // Attach role domain object for token creation

	// 6. Buat Session ID baru
	sessionID := uuid.New()

	// 7. Buat PasetoPayload dari domain
	payloadAccess := domain.PasetoPayload{
		UserID:    user.ID.String(),
		Email:     user.Email,
		SessionID: sessionID.String(),
		RoleName:  user.Role.Name,
		TokenType: domain.TokenTypeAccess,
	}
	payloadRefresh := domain.PasetoPayload{
		UserID:    user.ID.String(),
		Email:     user.Email,
		SessionID: sessionID.String(),
		RoleName:  user.Role.Name,
		TokenType: domain.TokenTypeRefresh,
	}

	// 8. Buat Access Token
	accessTTL := time.Duration(s.cfg.Security.AccessTokenTTLMin) * time.Minute
	accessToken, err := s.tokenSvc.CreateToken(payloadAccess, accessTTL)
	if err != nil {
		applogger.ErrorLogger.Printf("Register: failed to create access token for %s: %v", user.Email, err)
		// Sebaiknya tidak mengembalikan error internal ke user
		return nil, errors.New("failed to generate session token")
	}

	// 9. Buat Refresh Token
	refreshTTL := time.Duration(s.cfg.Security.RefreshTokenTTLMin) * time.Minute
	refreshToken, err := s.tokenSvc.CreateToken(payloadRefresh, refreshTTL)
	if err != nil {
		applogger.ErrorLogger.Printf("Register: failed to create refresh token for %s: %v", user.Email, err)
		// Sebaiknya tidak mengembalikan error internal ke user
		return nil, errors.New("failed to generate refresh token")
	}

	// 9. Log aktivitas registrasi
	detailsRegister, _ := json.Marshal(map[string]interface{}{
		"email": req.Email, // Ambil dari request
	})
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    user.ID,
		Action:    domain.ActionRegister,
		SessionID: sessionID,
		Details:   detailsRegister, // Gunakan details baru
	})

	// 10. Kembalikan respons
	// Panggil toUserResponse untuk mendapatkan URL avatar yang benar (termasuk fallback jika ID null)
	userResponse := s.toUserResponse(user)

	return &domain.AuthResponse{
		User: *userResponse,
		Token: domain.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *userService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// 1. Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		applogger.ErrorLogger.Printf("Login: database error checking email %s: %v", req.Email, err)
		return nil, errors.New("database error")
	}
	if user == nil {
		return nil, errors.New("invalid email or password") // Pesan generik
	}
	if user.DeletedAt.Valid || !user.IsActive {
		// Tolak jika user sudah soft deleted (DeletedAt terisi), ATAU tidak aktif
		return nil, errors.New("invalid email or password") // Pesan generik
	}

	// 2. Bandingkan password
	match, err := s.passSvc.Compare(req.Password, user.Password)
	if err != nil || !match {
		return nil, errors.New("invalid email or password")
	}

	// 3. Buat Session ID baru
	sessionID := uuid.New()

	// 4. Log aktivitas login
	detailsLogin, _ := json.Marshal(map[string]interface{}{
		"email": user.Email, // Ambil dari user object
	})
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    user.ID,
		Action:    domain.ActionLogin,
		SessionID: sessionID,
		Details:   detailsLogin, // Gunakan details baru
	})

	// 5. Buat PasetoPayload dari domain
	payloadAccess := domain.PasetoPayload{
		UserID:    user.ID.String(),
		Email:     user.Email,
		SessionID: sessionID.String(),
		RoleName:  user.Role.Name,
		TokenType: domain.TokenTypeAccess,
	}
	payloadRefresh := domain.PasetoPayload{
		UserID:    user.ID.String(),
		Email:     user.Email,
		SessionID: sessionID.String(),
		RoleName:  user.Role.Name,
		TokenType: domain.TokenTypeRefresh,
	}

	// 6. Buat Access Token
	accessTTL := time.Duration(s.cfg.Security.AccessTokenTTLMin) * time.Minute
	accessToken, err := s.tokenSvc.CreateToken(payloadAccess, accessTTL)
	if err != nil {
		applogger.ErrorLogger.Printf("Login: failed to create access token for %s: %v", user.Email, err)
		return nil, errors.New("failed to create access token")
	}

	// 7. Buat Refresh Token
	refreshTTL := time.Duration(s.cfg.Security.RefreshTokenTTLMin) * time.Minute
	refreshToken, err := s.tokenSvc.CreateToken(payloadRefresh, refreshTTL)
	if err != nil {
		applogger.ErrorLogger.Printf("Login: failed to create refresh token for %s: %v", user.Email, err)
		return nil, errors.New("failed to create refresh token")
	}

	// 8. Kembalikan respons
	return &domain.AuthResponse{
		User: *s.toUserResponse(user),
		Token: domain.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *userService) SendVerificationEmail(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	// 1. Find user by ID
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		applogger.ErrorLogger.Printf("SendVerificationEmail: DB error finding user %s: %v", userID, err)
		return errors.New("database error")
	}
	if user == nil {
		// Should not happen if called after AuthMiddleware, but handle defensively
		return errors.New("user not found")
	}
	if user.EmailVerifiedAt != nil {
		return errors.New("email already verified")
	}

	// 2. Generate Verification Attempt ID (ID unik untuk link ini)
	attemptID := uuid.New()

	// 3. Buat PasetoPayload untuk token verifikasi dari domain
	verificationPayload := domain.PasetoPayload{
		UserID:                user.ID.String(),
		Email:                 user.Email,
		SessionID:             sessionID.String(),
		RoleName:              user.Role.Name,
		VerificationAttemptID: attemptID.String(),
		TokenType:             domain.TokenTypeEmailVerification,
	}

	// 4. Create verification token (1 jam)
	verificationToken, err := s.tokenSvc.CreateToken(verificationPayload, 1*time.Hour)
	if err != nil {
		applogger.ErrorLogger.Printf("SendVerificationEmail: Failed to create verification token for %s: %v", user.Email, err)
		return errors.New("failed to create verification token")
	}

	// 5. Create verification link
	backendBaseURL := s.cfg.Storage.StoragePublicBaseURL
	verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s", backendBaseURL, verificationToken)

	// 4. Send email
	subject := "Verifikasi Alamat Email Anda - EDSA"
	body := fmt.Sprintf(
		"Halo %s,<br><br>"+
			"Terima kasih telah mendaftar. Silakan klik tautan di bawah ini untuk memverifikasi alamat email Anda:<br><br>"+
			"<a href=\"%s\" style=\"background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;\">Verifikasi Email</a><br><br>"+
			"Tautan ini berlaku selama 1 jam.<br><br>"+
			"Jika Anda tidak mendaftar, abaikan email ini.<br><br>"+
			"Terima kasih,<br>Tim EDSA",
		user.Name, verificationLink,
	)

	if err := s.emailSvc.SendEmail(ctx, user.Email, subject, body); err != nil {
		applogger.ErrorLogger.Printf("SendVerificationEmail: Failed to send email to %s: %v", user.Email, err)
		return errors.New("failed to send verification email")
	}

	// 7. Log aktivitas resend verification email
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    user.ID,
		Action:    domain.ActionResendVerificationEmail,
		SessionID: sessionID,
		Details: func() json.RawMessage {
			detailMap := map[string]interface{}{"email": user.Email}
			jsonData, _ := json.Marshal(detailMap)
			return jsonData
		}(),
	})

	return nil
}

func (s *userService) VerifyEmail(ctx context.Context, token string) error {
	// 1. Validasi token (dapatkan seluruh payload)
	payload, err := s.tokenSvc.ValidateToken(token, domain.TokenTypeEmailVerification)
	if err != nil {
		// Token tidak valid atau kedaluwarsa
		return fmt.Errorf("token invalid or expired: %w", err)
	}

	// Parse UUIDs dari payload
	userID, _ := uuid.Parse(payload.UserID)
	sessionID, _ := uuid.Parse(payload.SessionID) // Session ID utama pengguna
	attemptIDStr := payload.VerificationAttemptID
	if attemptIDStr == "" {
		return errors.New("invalid verification token: missing attempt ID")
	}
	attemptID, err := uuid.Parse(attemptIDStr)
	if err != nil {
		return errors.New("invalid verification token: invalid attempt ID format")
	}

	// 2. Cek apakah link/attempt ini sudah pernah digunakan (di-blacklist)
	isBlacklisted, err := s.blacklistSvc.IsSessionBlacklisted(ctx, attemptID)
	if err != nil {
		applogger.ErrorLogger.Printf("VerifyEmail: Error checking blacklist for attempt %s: %v", attemptID, err)
		return errors.New("internal server error checking token status")
	}
	if isBlacklisted {
		return errors.New("token already used")
	}

	// 3. Cari user (gunakan userID dari payload)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		applogger.ErrorLogger.Printf("VerifyEmail: DB error finding user %s: %v", userID, err)
		return errors.New("database error finding user")
	}
	if user == nil {
		// User tidak ditemukan -> token dianggap tidak valid
		return errors.New("user associated with token not found") // Error spesifik
	}
	if user.EmailVerifiedAt != nil {
		// Sudah terverifikasi sebelumnya
		return errors.New("email already verified") // Error spesifik
	}

	// 4. Update status verifikasi
	now := time.Now()
	user.EmailVerifiedAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		applogger.ErrorLogger.Printf("VerifyEmail: Failed to update user %s: %v", userID, err)
		return errors.New("failed to update verification status")
	}

	// 5. Blacklist attempt ID ini agar tidak bisa dipakai lagi
	blacklistDuration := 1*time.Hour + 5*time.Minute
	if err := s.blacklistSvc.BlacklistSession(ctx, attemptID, blacklistDuration); err != nil {
		applogger.ErrorLogger.Printf("VerifyEmail: WARNING - Failed to blacklist verification attempt %s after successful verification for user %s: %v", attemptID, userID, err)
	} else {
		applogger.ErrorLogger.Printf("VerifyEmail: INFO - Successfully blacklisted verification attempt %s for user %s", attemptID, userID)
	}

	// 6. Log aktivitas verifikasi email
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    userID,
		Action:    domain.ActionVerifyEmail,
		SessionID: sessionID,
		Details: func() json.RawMessage {
			detailMap := map[string]interface{}{"email": user.Email}
			jsonData, _ := json.Marshal(detailMap)
			return jsonData
		}(),
	})

	return nil
}

func (s *userService) SendPasswordResetEmail(ctx context.Context, email string) error {
	// 1. Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		applogger.ErrorLogger.Printf("SendPasswordResetEmail: DB error finding user %s: %v", email, err)
		return errors.New("database error")
	}
	// Jangan beri tahu jika email tidak ada (untuk keamanan)
	if user == nil || user.DeletedAt.Valid || !user.IsActive {
		applogger.ErrorLogger.Printf("SendPasswordResetEmail: Attempt to reset password for non-existent, inactive or deleted user: %s", email)
		return nil // Kembalikan nil agar attacker tidak tahu email mana yang terdaftar/aktif
	}

	// 2. Buat PasetoPayload untuk reset token
	// SessionID tidak terlalu relevan di sini, bisa buat baru atau gunakan yang acak
	resetSessionID := uuid.New()
	resetPayload := domain.PasetoPayload{
		UserID:    user.ID.String(),
		Email:     user.Email,
		SessionID: resetSessionID.String(),
		RoleName:  user.Role.Name,
		TokenType: domain.TokenTypePasswordReset,
	}

	// 3. Buat token reset password (15 menit)
	resetToken, err := s.tokenSvc.CreateToken(resetPayload, 15*time.Minute)
	if err != nil {
		applogger.ErrorLogger.Printf("SendPasswordResetEmail: Failed to create reset token for %s: %v", email, err)
		return errors.New("failed to create reset token")
	}

	// 4. Buat link reset
	frontendURL := s.cfg.App.FrontendURL
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, resetToken)

	// 5. Kirim email
	subject := "Reset Password Akun EDSA Anda"
	body := fmt.Sprintf(
		"Halo %s,<br><br>"+
			"Kami menerima permintaan untuk mereset password akun Anda. Silakan klik tautan di bawah ini:<br><br>"+
			"<a href=\"%s\" style=\"background-color: #2196F3; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;\">Reset Password</a><br><br>"+
			"Tautan ini hanya berlaku selama <strong>15 menit</strong>.<br><br>"+
			"Jika Anda tidak meminta reset password, abaikan email ini.<br><br>"+
			"Terima kasih,<br>Tim EDSA",
		user.Name, resetLink,
	)

	if err := s.emailSvc.SendEmail(ctx, user.Email, subject, body); err != nil {
		applogger.ErrorLogger.Printf("SendPasswordResetEmail: Failed to send email to %s: %v", user.Email, err)
		return errors.New("failed to send password reset email")
	}

	return nil
}

func (s *userService) ResetPassword(ctx context.Context, req *domain.ResetPasswordRequest) error {
	// 1. Validasi token reset (dapatkan payload)
	payload, err := s.tokenSvc.ValidateToken(req.Token, domain.TokenTypePasswordReset)
	if err != nil {
		return fmt.Errorf("invalid or expired token: %w", err)
	}

	// Ekstrak userID dari payload
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		// Seharusnya tidak terjadi jika CreateToken benar, tapi handle defensively
		applogger.ErrorLogger.Printf("ResetPassword: Invalid UserID format in token payload: %s", payload.UserID)
		return errors.New("invalid user data in token")
	}

	// 2. Cari user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		applogger.ErrorLogger.Printf("ResetPassword: DB error finding user %s: %v", userID, err)
		return errors.New("database error")
	}
	// Pastikan user masih ada dan aktif
	if user == nil || user.DeletedAt.Valid || !user.IsActive {
		return errors.New("user associated with token not found or inactive")
	}

	// 3. Hash password baru
	hashedPassword, err := s.passSvc.Hash(req.NewPassword)
	if err != nil {
		applogger.ErrorLogger.Printf("ResetPassword: Failed to hash new password for user %s: %v", userID, err)
		return errors.New("failed to hash new password")
	}

	// 4. Update password user
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		applogger.ErrorLogger.Printf("ResetPassword: Failed to update password for user %s: %v", userID, err)
		return errors.New("failed to update password")
	}

	// TODO: Idealnya, blacklist semua sesi aktif user ini setelah reset password

	return nil
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenResponse, error) {
	// 1. Validasi refresh token
	payload, err := s.tokenSvc.ValidateToken(refreshToken, domain.TokenTypeRefresh)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token: %w", err)
	}

	// Parse UUIDs
	userID, _ := uuid.Parse(payload.UserID)
	sessionID, _ := uuid.Parse(payload.SessionID)

	// 2. Cek apakah sesi sudah di-blacklist
	isBlacklisted, err := s.blacklistSvc.IsSessionBlacklisted(ctx, sessionID)
	if err != nil {
		return nil, errors.New("session check error")
	}
	if isBlacklisted {
		return nil, errors.New("session has been logged out")
	}

	// 3. Dapatkan data user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		applogger.ErrorLogger.Printf("RefreshToken: user not found for token with UserID %s: %v", userID, err)
		return nil, errors.New("user not found for this token")
	}

	// 4. Buat PasetoPayload baru dari domain
	newPayloadAccess := domain.PasetoPayload{
		UserID:    payload.UserID,
		Email:     payload.Email,
		SessionID: payload.SessionID,
		RoleName:  payload.RoleName,
		TokenType: domain.TokenTypeAccess,
	}
	newPayloadRefresh := domain.PasetoPayload{
		UserID:    payload.UserID,
		Email:     payload.Email,
		SessionID: payload.SessionID,
		RoleName:  payload.RoleName,
		TokenType: domain.TokenTypeRefresh,
	}

	// 5. Buat Access Token baru
	accessTTL := time.Duration(s.cfg.Security.AccessTokenTTLMin) * time.Minute
	accessToken, err := s.tokenSvc.CreateToken(newPayloadAccess, accessTTL)
	if err != nil {
		return nil, errors.New("failed to create new access token")
	}

	// 6. Buat Refresh Token baru
	refreshTTL := time.Duration(s.cfg.Security.RefreshTokenTTLMin) * time.Minute
	newRefreshToken, err := s.tokenSvc.CreateToken(newPayloadRefresh, refreshTTL)
	if err != nil {
		return nil, errors.New("failed to create new refresh token")
	}

	// 7. Kembalikan token baru
	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *userService) Logout(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, email string) error {
	// 0. Ambil data user untuk logging email (opsional, tapi bagus untuk detail log)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		// Log error tapi jangan gagalkan logout
		applogger.ErrorLogger.Printf("Logout: Failed to fetch user %s for logging email: %v", userID, err)
	}

	// 1. Tambahkan SESSION ID ke blacklist
	duration := time.Duration(s.cfg.Security.BlacklistTTLHour) * time.Hour
	if err := s.blacklistSvc.BlacklistSession(ctx, sessionID, duration); err != nil {
		applogger.ErrorLogger.Printf("Logout: failed to blacklist session %s for user %s: %v", sessionID, userID, err)
		return errors.New("failed to invalidate session")
	}

	// 2. Log aktivitas logout
	var detailsLogout json.RawMessage
	if user != nil {
		detailsLogout, _ = json.Marshal(map[string]interface{}{
			"email": user.Email,
		})
	}
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    userID,
		Action:    domain.ActionLogout,
		SessionID: sessionID,
		Details:   detailsLogout, // Gunakan details baru (bisa nil jika user fetch gagal)
	})
	return nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	// 1. Cari user berdasarkan ID
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("database error")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 2. Konversi ke DTO dasar
	userResponse := s.toUserResponse(user)

	// 3. Hitung overall score
	progresses, err := s.progressRepo.FindAllByUserID(ctx, id)
	if err != nil {
		// Log error, tapi jangan gagalkan permintaan. Skor akan 0.
		applogger.ErrorLogger.Printf("GetUserByID: Failed to fetch progresses for user %s: %v", id, err)
	}

	var totalScore float64 = 0.0
	completedCount := 0
	for _, p := range progresses {
		if p.Status == domain.BookProgressStatusCompleted {
			totalScore += p.HighestScore
			completedCount++
		}
	}

	overallScore := 0
	if completedCount > 0 {
		// Lakukan pembagian float, lalu bulatkan ke integer terdekat
		overallScore = int(math.Round(totalScore / float64(completedCount)))
	}

	userResponse.OverallScore = overallScore // Tambahkan skor ke DTO

	// 4. Set EmailVerified di respons DTO
	userResponse.EmailVerified = user.EmailVerifiedAt != nil // Jika EmailVerifiedAt tidak NULL, berarti sudah terverifikasi

	// 5. Kembalikan respons
	return userResponse, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, req *domain.ChangePasswordRequest) error {
	// 1. Ambil user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		applogger.ErrorLogger.Printf("ChangePassword: DB error checking user %s: %v", userID, err)
		return errors.New("database error")
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 2. Validasi password saat ini
	match, err := s.passSvc.Compare(req.CurrentPassword, user.Password)
	if err != nil {
		applogger.ErrorLogger.Printf("ChangePassword: Error comparing password for user %s: %v", userID, err)
		return errors.New("password comparison failed")
	}
	if !match {
		return errors.New("invalid current password")
	}

	// 3. Hash password baru
	hashedPassword, err := s.passSvc.Hash(req.NewPassword)
	if err != nil {
		applogger.ErrorLogger.Printf("ChangePassword: failed to hash new password for user %s: %v", userID, err)
		return errors.New("failed to hash new password")
	}

	// 4. Update password di struct domain
	user.Password = hashedPassword

	// 5. Simpan ke database
	if err := s.userRepo.Update(ctx, user); err != nil {
		applogger.ErrorLogger.Printf("ChangePassword: failed to update user %s in DB: %v", userID, err)
		return errors.New("failed to save new password")
	}

	// TODO: Sebaiknya, semua sesi lain di-blacklist setelah ganti password
	// (Ini di luar scope, tapi penting untuk keamanan)

	return nil
}

func (s *userService) UpdateUserDetails(ctx context.Context, userID uuid.UUID, req *domain.UpdateDetailsRequest) (*domain.UserResponse, error) {
	// 1. Ambil user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		applogger.ErrorLogger.Printf("UpdateUserDetails: DB error checking user %s: %v", userID, err)
		return nil, errors.New("database error")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 2. Validasi password (konfirmasi identitas)
	match, err := s.passSvc.Compare(req.Password, user.Password)
	if err != nil {
		applogger.ErrorLogger.Printf("UpdateUserDetails: Error comparing password for user %s: %v", userID, err)
		return nil, errors.New("password comparison failed")
	}
	if !match {
		return nil, errors.New("invalid password confirmation")
	}

	// 3. Update nama di struct domain
	user.Name = req.Name

	// 4. Simpan ke database
	if err := s.userRepo.Update(ctx, user); err != nil {
		applogger.ErrorLogger.Printf("UpdateUserDetails: failed to update user %s in DB: %v", userID, err)
		return nil, errors.New("failed to save user details")
	}

	// 5. Kembalikan respons DTO
	return s.toUserResponse(user), nil
}

func (s *userService) UpdateAvatar(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) (*domain.UserResponse, error) {
	// 1. Ambil user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		applogger.ErrorLogger.Printf("UpdateAvatar: DB error checking user %s: %v", userID, err)
		return nil, errors.New("database error")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 2. Ambil media asset
	media, err := s.mediaRepo.FindByID(ctx, mediaID)
	if err != nil {
		applogger.ErrorLogger.Printf("UpdateAvatar: DB error checking media %s: %v", mediaID, err)
		return nil, errors.New("database error")
	}
	if media == nil {
		return nil, errors.New("media asset not found")
	}

	// 3. Validasi Keamanan: Pastikan media ini memang diupload oleh user ini
	//    dan tipenya adalah 'user_avatar'.
	if media.UploadedByUserID == nil || *media.UploadedByUserID != userID {
		return nil, errors.New("you do not own this media asset")
	}
	if media.OwnerType != domain.OwnerTypeUserAvatar {
		return nil, errors.New("this media asset is not designated for user avatars")
	}

	// 4. Update ID di struct domain
	user.ProfilePictureID = &mediaID
	user.ProfilePicture = *media // Update relasi yang di-memori

	// 5. Simpan ke database
	if err := s.userRepo.Update(ctx, user); err != nil {
		applogger.ErrorLogger.Printf("UpdateAvatar: failed to update user %s in DB: %v", userID, err)
		return nil, errors.New("failed to save avatar")
	}

	// 6. Kembalikan respons DTO
	return s.toUserResponse(user), nil
}

// DEPRECATED: Gunakan ConfirmAccountDeletion dan RequestAccountDeletion
func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, req *domain.DeleteAccountRequest) error {
	// 1. Ambil user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		applogger.ErrorLogger.Printf("DeleteUser: DB error checking user %s: %v", userID, err)
		return errors.New("database error")
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 2. Validasi password saat ini
	match, err := s.passSvc.Compare(req.CurrentPassword, user.Password)
	if err != nil {
		applogger.ErrorLogger.Printf("DeleteUser: Error comparing password for user %s: %v", userID, err)
		return errors.New("password comparison failed")
	}
	if !match {
		return errors.New("invalid current password")
	}

	// 3. Log aktivitas SEBELUM menghapus
	details, _ := json.Marshal(map[string]interface{}{
		"reason": req.DeletionReason,
	})
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    userID,
		Action:    domain.ActionDeleteAccount,
		SessionID: sessionID,
		Details:   details,
	})

	// 4. Hapus user dari database
	// DB constraint (OnDelete:CASCADE) akan menangani data terkait
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		applogger.ErrorLogger.Printf("DeleteUser: failed to delete user %s from DB: %v", userID, err)
		return errors.New("failed to delete user")
	}

	// 5. Blacklist sesi ini
	duration := time.Duration(s.cfg.Security.BlacklistTTLHour) * time.Hour
	if err := s.blacklistSvc.BlacklistSession(ctx, sessionID, duration); err != nil {
		applogger.ErrorLogger.Printf("DeleteUser: failed to blacklist session %s post-deletion: %v", sessionID, err)
		// Jangan kembalikan error, karena penghapusan user sudah berhasil
	}

	return nil
}

func (s *userService) RequestAccountDeletion(ctx context.Context, userID uuid.UUID) error {
	// 1. Ambil data user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// 2. Buat PasetoPayload untuk token konfirmasi
	// SessionID tidak terlalu relevan di sini, buat acak
	deletionSessionID := uuid.New()
	deletionPayload := domain.PasetoPayload{
		UserID:    user.ID.String(),
		Email:     user.Email,
		SessionID: deletionSessionID.String(),
		RoleName:  user.Role.Name,
		TokenType: domain.TokenTypeAccountDeletion,
	}

	// 3. Buat token konfirmasi (5 menit)
	deletionToken, err := s.tokenSvc.CreateToken(deletionPayload, 5*time.Minute)
	if err != nil {
		applogger.ErrorLogger.Printf("RequestAccountDeletion: Failed to create deletion token for %s: %v", userID, err)
		return errors.New("failed to create confirmation token")
	}

	// 4. Buat link konfirmasi
	// Cth: http://localhost:3000/confirm-delete?token=...
	frontendURL := s.cfg.App.FrontendURL
	confirmationLink := fmt.Sprintf("%s/confirm-delete?token=%s", frontendURL, deletionToken)

	// 5. Buat body email
	subject := "Konfirmasi Penghapusan Akun EDSA"
	body := fmt.Sprintf(
		"Halo %s,<br><br>"+
			"Kami menerima permintaan untuk menghapus akun Anda. Untuk mengkonfirmasi, silakan klik tautan di bawah ini:<br><br>"+
			"<a href=\"%s\" style=\"background-color: #f44336; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;\">Konfirmasi Hapus Akun</a><br><br>"+
			"Tautan ini hanya berlaku selama <strong>5 menit</strong>.<br><br>"+
			"Jika Anda tidak meminta penghapusan ini, abaikan saja email ini.<br><br>"+
			"Terima kasih,<br>Tim EDSA",
		user.Name, confirmationLink,
	)

	// 6. Kirim email
	if err := s.emailSvc.SendEmail(ctx, user.Email, subject, body); err != nil {
		applogger.ErrorLogger.Printf("RequestAccountDeletion: Failed to send email to %s: %v", user.Email, err)
		return errors.New("failed to send confirmation email")
	}

	return nil
}

func (s *userService) ConfirmAccountDeletion(ctx context.Context, userID uuid.UUID, req *domain.ConfirmDeletionRequest) error {
	// 1. Ambil data user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// 2. Validasi Password Saat Ini
	match, err := s.passSvc.Compare(req.CurrentPassword, user.Password)
	if err != nil {
		applogger.ErrorLogger.Printf("ConfirmAccountDeletion: Error comparing password for user %s: %v", userID, err)
		return errors.New("password comparison failed")
	}
	if !match {
		return errors.New("invalid current password")
	}

	// 3. Validasi Token Konfirmasi (dapatkan payload)
	payload, err := s.tokenSvc.ValidateToken(req.ConfirmationToken, domain.TokenTypeAccountDeletion)
	if err != nil {
		return err // Error sudah termasuk "invalid token type"
	}

	// Ekstrak tokenUserID dari payload
	tokenUserID, err := uuid.Parse(payload.UserID)
	if err != nil {
		applogger.ErrorLogger.Printf("ConfirmAccountDeletion: Invalid UserID format in token payload: %s", payload.UserID)
		return errors.New("invalid user data in token")
	}

	// 4. Pastikan token tersebut milik pengguna yang sedang login
	if tokenUserID != userID {
		return errors.New("confirmation token does not match authenticated user")
	}

	// 5. Log aktivitas penghapusan akun
	// Gabungkan reason dan email
	detailsDelete, _ := json.Marshal(map[string]interface{}{
		"reason": req.DeletionReason,
		"email":  user.Email,
	})
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    userID,
		Action:    domain.ActionDeleteAccount,
		SessionID: uuid.New(),
		Details:   detailsDelete,
	})

	// 6. Lakukan Soft Delete
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		applogger.ErrorLogger.Printf("ConfirmAccountDeletion: Failed to soft delete user %s from DB: %v", userID, err)
		return errors.New("failed to deactivate account")
	}

	// 7. Blacklist sesi ini agar token tidak bisa dipakai lagi (Opsional tapi direkomendasikan)
	// Ambil sessionID dari token konfirmasi (jika diperlukan) atau blacklist semua sesi user
	// Untuk simple, kita skip blacklist semua sesi di sini, karena login akan dicegah

	return nil
}
