package user

import (
	"context"
	"errors"
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
}

func NewUserService(
	userRepo usecase.UserRepository,
	roleRepo usecase.RoleRepository,
	passSvc usecase.PasswordService,
	tokenSvc usecase.TokenService,
	cfg *config.Config,
	logger usecase.ActivityLoggerService,
	blacklistSvc usecase.SessionBlacklistService,
) usecase.UserService {
	return &userService{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		passSvc:      passSvc,
		tokenSvc:     tokenSvc,
		cfg:          cfg,
		logger:       logger,
		blacklistSvc: blacklistSvc,
	}
}

// --- Helper Mapper ---
func toUserResponse(user *domain.User) *domain.UserResponse {
	return &domain.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		RoleName: user.Role.Name,
		JoinedAt: user.CreatedAt,
		// ProfilePictureURL: user.ProfilePictureURL,
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
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		RoleID:   defaultRole.ID,
		// ProfilePictureURL akan default ke string kosong
	}

	// 5. Simpan ke database
	if err := s.userRepo.Create(ctx, user); err != nil {
		applogger.ErrorLogger.Printf("Register: failed to create user %s: %v", user.Email, err)
		return nil, errors.New("failed to create user")
	}
	user.Role = *defaultRole
	// Ambil CreatedAt yang di-generate DB (meskipun mapper di bawah akan menggunakannya)
	// user.CreatedAt = ... (userRepo.Create seharusnya meng-update ID, kita asumsikan CreatedAt juga)

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

	// 8. Kembalikan respons
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

	// 2. Kembalikan respons
	return toUserResponse(user), nil
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
