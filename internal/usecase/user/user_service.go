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
	passSvc  usecase.PasswordService
	tokenSvc usecase.TokenService
	cfg      *config.Config
}

func NewUserService(
	userRepo usecase.UserRepository,
	passSvc usecase.PasswordService,
	tokenSvc usecase.TokenService,
	cfg *config.Config,
) usecase.UserService {
	return &userService{
		userRepo: userRepo,
		passSvc:  passSvc,
		tokenSvc: tokenSvc,
		cfg:      cfg,
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

	// 3. Buat domain user baru
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		// TODO: Ini harusnya mengambil ID "user" dari database secara dinamis
		//       (misal: 2), bukan di-hardcode.
		//       Ini memerlukan RoleRepository.
		RoleID: 2,
	}

	// 4. Simpan ke database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// 5. Buat Access Token
	accessTTL := time.Duration(s.cfg.Security.AccessTokenTTLMin) * time.Minute
	accessToken, err := s.tokenSvc.CreateToken(user, accessTTL)
	if err != nil {
		return nil, errors.New("failed to create access token")
	}

	// 6. Buat Refresh Token
	refreshTTL := time.Duration(s.cfg.Security.RefreshTokenTTLMin) * time.Minute
	refreshToken, err := s.tokenSvc.CreateToken(user, refreshTTL)
	if err != nil {
		return nil, errors.New("failed to create refresh token")
	}

	// 7. Kembalikan respons
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
