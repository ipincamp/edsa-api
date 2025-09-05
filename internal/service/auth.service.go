package service

import (
	"context"
	"errors"
	"log"

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
	bloomFilter    *util.BloomFilterManager
}

func NewAuth(userRepository domain.UserRepository, db *gorm.DB, cfg *config.Config, bloomFilter *util.BloomFilterManager) domain.AuthService {
	return &authService{
		userRepository: userRepository,
		db:             db,
		config:         cfg,
		bloomFilter:    bloomFilter,
	}
}

func (s *authService) Register(ctx context.Context, request dto.RegisterRequest) (dto.AuthResponse, error) {
	if s.bloomFilter.Test(request.Email) {
		return dto.AuthResponse{}, errors.New("email already exists")
	}

	defaultRole, found := util.GetRoleByName(constant.RoleGuest.String())
	if !found {
		return dto.AuthResponse{}, errors.New("default role not found in cache")
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

	s.bloomFilter.Add(newUser.Email)
	go func() {
		if err := s.bloomFilter.Save(); err != nil {
			log.Printf("Error saving bloom filter after registration: %v", err)
		}
	}()

	newUser.Role = defaultRole
	return createAuthResponse(newUser, s)
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

	return createAuthResponse(user, s)
}

func (s *authService) Logout(ctx context.Context, user domain.User) error {
	return nil
}
