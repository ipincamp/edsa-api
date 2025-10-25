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
	}
}

// --- Helper Mapper ---
func toUserResponse(user *domain.User) *domain.UserResponse {
	return &domain.UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		RoleName:          user.Role.Name,
		JoinedAt:          user.CreatedAt,
		ProfilePictureURL: user.ProfilePictureURL,
		IsActive:          user.IsActive,
		EmailVerified:     user.EmailVerifiedAt != nil,
		// OverallScore dihitung terpisah
	}
}

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
		return nil, errors.New("default role not found in database")
	}

	// 4. Buat domain user baru
	// Ambil URL avatar default dari config
	baseURL := s.cfg.Storage.StoragePublicBaseURL
	defaultAvatarURL := baseURL + "/public/uploads/avatar1.png"

	user := &domain.User{
		Name:              req.Name,
		Email:             req.Email,
		Password:          hashedPassword,
		RoleID:            defaultRole.ID,
		ProfilePictureURL: defaultAvatarURL,
		IsActive:          true,
	}

	// 5. Simpan ke database
	if err := s.userRepo.Create(ctx, user); err != nil {
		applogger.ErrorLogger.Printf("Register: failed to create user %s: %v", user.Email, err)
		return nil, errors.New("failed to create user")
	}
	user.Role = *defaultRole

	// 6. Buat Session ID baru
	sessionID := uuid.New()

	// 7. Buat Access Token
	accessTTL := time.Duration(s.cfg.Security.AccessTokenTTLMin) * time.Minute
	accessToken, err := s.tokenSvc.CreateToken(user, sessionID, accessTTL)
	if err != nil {
		applogger.ErrorLogger.Printf("Register: failed to create access token for %s: %v", user.Email, err)
		return nil, errors.New("failed to create access token")
	}

	// 8. Buat Refresh Token
	refreshTTL := time.Duration(s.cfg.Security.RefreshTokenTTLMin) * time.Minute
	refreshToken, err := s.tokenSvc.CreateToken(user, sessionID, refreshTTL)
	if err != nil {
		applogger.ErrorLogger.Printf("Register: failed to create refresh token for %s: %v", user.Email, err)
		return nil, errors.New("failed to create refresh token")
	}

	// 9. Log aktivitas registrasi
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    user.ID,
		Action:    domain.ActionRegister,
		SessionID: sessionID,
	})

	// 10. Kembalikan respons
	return &domain.AuthResponse{
		User: *toUserResponse(user),
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

	// 2. Bandingkan password
	match, err := s.passSvc.Compare(req.Password, user.Password)
	if err != nil || !match {
		return nil, errors.New("invalid email or password")
	}

	// 3. Buat Session ID baru
	sessionID := uuid.New()

	// 4. Log aktivitas login
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    user.ID,
		Action:    domain.ActionLogin,
		SessionID: sessionID,
	})

	// 5. Buat Access Token
	accessTTL := time.Duration(s.cfg.Security.AccessTokenTTLMin) * time.Minute
	accessToken, err := s.tokenSvc.CreateToken(user, sessionID, accessTTL)
	if err != nil {
		applogger.ErrorLogger.Printf("Login: failed to create access token for %s: %v", user.Email, err)
		return nil, errors.New("failed to create access token")
	}

	// 6. Buat Refresh Token
	refreshTTL := time.Duration(s.cfg.Security.RefreshTokenTTLMin) * time.Minute
	refreshToken, err := s.tokenSvc.CreateToken(user, sessionID, refreshTTL)
	if err != nil {
		applogger.ErrorLogger.Printf("Login: failed to create refresh token for %s: %v", user.Email, err)
		return nil, errors.New("failed to create refresh token")
	}

	// 7. Kembalikan respons
	return &domain.AuthResponse{
		User: *toUserResponse(user),
		Token: domain.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *userService) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.TokenResponse, error) {
	// 1. Validasi refresh token
	userID, sessionID, err := s.tokenSvc.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// 2. Cek apakah sesi sudah di-blacklist (logout)
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

	// 4. Buat Access Token baru
	accessTTL := time.Duration(s.cfg.Security.AccessTokenTTLMin) * time.Minute
	accessToken, err := s.tokenSvc.CreateToken(user, sessionID, accessTTL)
	if err != nil {
		return nil, errors.New("failed to create new access token")
	}

	// 5. Buat Refresh Token baru
	refreshTTL := time.Duration(s.cfg.Security.RefreshTokenTTLMin) * time.Minute
	refreshToken, err := s.tokenSvc.CreateToken(user, sessionID, refreshTTL)
	if err != nil {
		return nil, errors.New("failed to create new refresh token")
	}

	// 6. Kembalikan token baru
	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Logout(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	// 1. Tambahkan SESSION ID ke blacklist
	duration := time.Duration(s.cfg.Security.BlacklistTTLHour) * time.Hour
	// Gunakan sessionID, bukan tokenString
	if err := s.blacklistSvc.BlacklistSession(ctx, sessionID, duration); err != nil {
		applogger.ErrorLogger.Printf("Logout: failed to blacklist session %s for user %s: %v", sessionID, userID, err)
		// Ini adalah error kritis, kita harus mengembalikannya
		return errors.New("failed to invalidate session")
	}

	// 2. Log aktivitas logout
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    userID,
		Action:    domain.ActionLogout,
		SessionID: sessionID,
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
	userResponse := toUserResponse(user)

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
	return toUserResponse(user), nil
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
	// 1. Ambil data user (terutama email)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// 2. Buat token konfirmasi (Paseto) yang berlaku 5 menit
	// Kita bisa gunakan sessionID acak karena tidak relevan untuk flow ini
	deletionToken, err := s.tokenSvc.CreateToken(user, uuid.New(), 5*time.Minute)
	if err != nil {
		applogger.ErrorLogger.Printf("RequestAccountDeletion: Failed to create deletion token for %s: %v", userID, err)
		return errors.New("failed to create confirmation token")
	}

	// 3. Buat link konfirmasi
	// Cth: http://localhost:3000/confirm-delete?token=...
	frontendURL := s.cfg.App.FrontendURL
	confirmationLink := fmt.Sprintf("%s/confirm-delete?token=%s", frontendURL, deletionToken)

	// 4. Buat body email
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

	// 5. Kirim email
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

	// 3. Validasi Token Konfirmasi
	tokenUserID, _, err := s.tokenSvc.ValidateToken(req.ConfirmationToken)
	if err != nil {
		// Cth: "token has expired" atau "invalid token"
		return err
	}

	// 4. Pastikan token tersebut milik pengguna yang sedang login
	if tokenUserID != userID {
		return errors.New("confirmation token does not match authenticated user")
	}

	// 5. Log aktivitas (termasuk alasannya)
	details, _ := json.Marshal(map[string]interface{}{
		"reason": req.DeletionReason,
	})
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    userID,
		Action:    domain.ActionDeleteAccount,
		SessionID: uuid.New(), // Sesi baru untuk tindakan ini
		Details:   details,
	})

	// 6. Lakukan Hard Delete
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		applogger.ErrorLogger.Printf("ConfirmAccountDeletion: Failed to delete user %s from DB: %v", userID, err)
		return errors.New("failed to delete account")
	}

	// TODO: Blacklist semua sesi yang ada untuk user ini?
	// Saat ini, user repo sudah dihapus, jadi token yang ada tidak akan divalidasi
	// oleh AuthMiddleware (karena userRepo.FindByID akan gagal).
	// Jadi, tidak perlu blacklist manual.

	return nil
}
