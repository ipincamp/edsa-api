package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type userService struct {
	userRepo usecase.UserRepository
	roleRepo usecase.RoleRepository
	passSvc  usecase.PasswordService
	tokenSvc usecase.TokenService
	cfg      *config.Config
	logger   usecase.ActivityLoggerService
}

func NewUserService(
	userRepo usecase.UserRepository,
	roleRepo usecase.RoleRepository,
	passSvc usecase.PasswordService,
	tokenSvc usecase.TokenService,
	cfg *config.Config,
	logger usecase.ActivityLoggerService,
) usecase.UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		passSvc:  passSvc,
		tokenSvc: tokenSvc,
		cfg:      cfg,
		logger:   logger,
	}
}

func (s *userService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// 1. Cek apakah email sudah ada
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("database error")
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// 2. Hash password
	hashedPassword, err := s.passSvc.Hash(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 3. Dapatkan role default (public)
	defaultRole, err := s.roleRepo.FindByName(ctx, domain.RoleNamePublic)
	if err != nil {
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
	}

	// 5. Simpan ke database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// 6. Buat Access Token
	accessTTL := time.Duration(s.cfg.Security.AccessTokenTTLMin) * time.Minute
	accessToken, err := s.tokenSvc.CreateToken(user, accessTTL)
	if err != nil {
		return nil, errors.New("failed to create access token")
	}

	// 7. Buat Refresh Token
	refreshTTL := time.Duration(s.cfg.Security.RefreshTokenTTLMin) * time.Minute
	refreshToken, err := s.tokenSvc.CreateToken(user, refreshTTL)
	if err != nil {
		return nil, errors.New("failed to create refresh token")
	}

	// 8. Kembalikan respons
	return &domain.AuthResponse{
		User: domain.UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			RoleID: user.RoleID,
		},
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

	// Log aktivitas login
	s.logger.Log(ctx, domain.ActivityLog{
		UserID: user.ID,
		Action: domain.ActionLogin,
	})

	// 3. Buat Access Token
	accessTTL := time.Duration(s.cfg.Security.AccessTokenTTLMin) * time.Minute
	accessToken, err := s.tokenSvc.CreateToken(user, accessTTL)
	if err != nil {
		return nil, errors.New("failed to create access token")
	}

	// 4. Buat Refresh Token
	refreshTTL := time.Duration(s.cfg.Security.RefreshTokenTTLMin) * time.Minute
	refreshToken, err := s.tokenSvc.CreateToken(user, refreshTTL)
	if err != nil {
		return nil, errors.New("failed to create refresh token")
	}

	// 5. Kembalikan respons
	return &domain.AuthResponse{
		User: domain.UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			RoleID: user.RoleID,
		},
		Token: domain.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *userService) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.TokenResponse, error) {
	// 1. Validasi refresh token
	userID, err := s.tokenSvc.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// 2. Dapatkan data user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found for this token")
	}

	// 3. Buat Access Token baru
	accessTTL := time.Duration(s.cfg.Security.AccessTokenTTLMin) * time.Minute
	accessToken, err := s.tokenSvc.CreateToken(user, accessTTL)
	if err != nil {
		return nil, errors.New("failed to create new access token")
	}

	// 4. Buat Refresh Token baru (Best practice: rotasi refresh token)
	refreshTTL := time.Duration(s.cfg.Security.RefreshTokenTTLMin) * time.Minute
	refreshToken, err := s.tokenSvc.CreateToken(user, refreshTTL)
	if err != nil {
		return nil, errors.New("failed to create new refresh token")
	}

	// 5. Kembalikan token baru
	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Logout(ctx context.Context, userID uuid.UUID) error {
	// Karena Paseto stateless, "logout" di sisi server berarti mencatat aktivitas.
	// Klien bertanggung jawab untuk menghapus token.
	// Jika ada blocklist (cth: Redis), token bisa ditambahkan di sini.

	// Log aktivitas logout
	s.logger.Log(ctx, domain.ActivityLog{
		UserID: userID,
		Action: domain.ActionLogout,
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
	return &domain.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		RoleID: user.RoleID,
	}, nil
}
