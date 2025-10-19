package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type userService struct {
	userRepo usecase.UserRepository
	passSvc  usecase.PasswordService
	tokenSvc usecase.TokenService
}

func NewUserService(
	userRepo usecase.UserRepository,
	passSvc usecase.PasswordService,
	tokenSvc usecase.TokenService,
) usecase.UserService {
	return &userService{
		userRepo: userRepo,
		passSvc:  passSvc,
		tokenSvc: tokenSvc,
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

	// 2. Hash password (Requirement 9)
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

	// 5. Buat token (Requirement 8)
	token, err := s.tokenSvc.CreateToken(user)
	if err != nil {
		return nil, errors.New("failed to create token")
	}

	// 6. Kembalikan respons
	return &domain.AuthResponse{
		User: domain.UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			RoleID: user.RoleID,
		},
		Token: token,
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

	// 2. Bandingkan password (Requirement 9)
	match, err := s.passSvc.Compare(req.Password, user.Password)
	if err != nil || !match {
		return nil, errors.New("invalid email or password")
	}

	// 3. Buat token (Requirement 8)
	token, err := s.tokenSvc.CreateToken(user)
	if err != nil {
		return nil, errors.New("failed to create token")
	}

	// 4. Kembalikan respons
	return &domain.AuthResponse{
		User: domain.UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			RoleID: user.RoleID,
		},
		Token: token,
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("database error")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &domain.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		RoleID: user.RoleID,
	}, nil
}
