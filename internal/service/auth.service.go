package service

import (
	"context"
	"errors"
	"time"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"gorm.io/gorm"
)

type authService struct {
	userRepository domain.UserRepository
	db             *gorm.DB
	config         *config.Config
}

func NewAuth(userRepository domain.UserRepository, db *gorm.DB, cfg *config.Config) domain.AuthService {
	return &authService{
		userRepository: userRepository,
		db:             db,
		config:         cfg,
	}
}

func (s *authService) Register(ctx context.Context, request dto.RegisterRequest) (dto.AuthResponse, error) {
	_, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err == nil {
		return dto.AuthResponse{}, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.AuthResponse{}, err
	}

	defaultRole, found := util.GetRoleByName(constant.RoleStudent.String())
	if !found {
		return dto.AuthResponse{}, errors.New("default student role not found in cache")
	}

	hashedPassword, err := util.HashPassword(request.Password)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	newUser := domain.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
		RoleID:   defaultRole.ID,
	}

	if err := s.userRepository.Save(ctx, &newUser); err != nil {
		return dto.AuthResponse{}, err
	}

	newUser.Role = defaultRole
	return s.createAuthResponse(newUser)
}

func (s *authService) Login(ctx context.Context, request dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AuthResponse{}, errors.New("invalid credentials")
		}
		return dto.AuthResponse{}, err
	}

	match, err := util.CheckPasswordHash(request.Password, user.Password)
	if err != nil || !match {
		return dto.AuthResponse{}, errors.New("invalid credentials")
	}

	return s.createAuthResponse(user)
}

func (s *authService) Logout(ctx context.Context, user domain.User) error {
	return nil
}

func (s *authService) createAuthResponse(user domain.User) (dto.AuthResponse, error) {
	tokenTTL := time.Duration(s.config.Paseto.TokenTTLMin) * time.Minute
	pasetoMaker, err := util.NewPasetoMaker(s.config.Paseto.SecretKey)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	token, err := pasetoMaker.CreateToken(user.ID, user.Role.ID, tokenTTL)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	res := dto.AuthResponse{
		Token: token,
		User: dto.UserData{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role.Name,
			JoinedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}
	return res, nil
}
