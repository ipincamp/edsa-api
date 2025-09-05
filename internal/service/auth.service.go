package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/repository"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"gorm.io/gorm"
)

type authService struct {
	userRepository domain.UserRepository
	db             *gorm.DB
	config         *config.Config
	bloomFilter    *util.BloomFilterManager
}

func NewAuth(userRepository domain.UserRepository, db *gorm.DB, cfg *config.Config, bloomFilter *util.BloomFilterManager) AuthService {
	return &authService{
		userRepository: userRepository,
		db:             db,
		config:         cfg,
		bloomFilter:    bloomFilter,
	}
}

func (s *authService) Register(ctx context.Context, request dto.RegisterRequest) (dto.AuthResponse, error) {
	if s.bloomFilter.Test(request.Email) {
		return dto.AuthResponse{}, constant.ErrConflict
	}

	hashedPassword, err := util.HashPassword(request.Password)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	newUser := domain.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		defaultRole, found := util.GetRoleByName(constant.RoleGuest.String())
		if !found {
			return errors.New("default role not found in cache")
		}
		newUser.RoleID = defaultRole.ID

		txUserRepo := repository.NewUser(tx)
		return txUserRepo.Create(ctx, &newUser)
	}); err != nil {
		return dto.AuthResponse{}, err
	}

	createdUser, err := s.userRepository.FindByEmail(ctx, newUser.Email)
	if err != nil {
		log.Printf("CRITICAL: Failed to fetch user right after registration: %v", err)
		return dto.AuthResponse{}, errors.New("failed to finalize registration")
	}

	s.bloomFilter.Add(createdUser.Email)
	go func() {
		if err := s.bloomFilter.Save(); err != nil {
			log.Printf("Error saving bloom filter after registration: %v", err)
		}
	}()

	util.AddUserToCache(createdUser)

	return s.createAuthResponse(createdUser)
}

func (s *authService) Login(ctx context.Context, request dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AuthResponse{}, constant.ErrInvalidCredentials
		}
		return dto.AuthResponse{}, err
	}

	match, err := util.CheckPasswordHash(request.Password, user.Password)
	if err != nil || !match {
		return dto.AuthResponse{}, constant.ErrInvalidCredentials
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
		User:  dto.ToUserResponse(user),
	}
	return res, nil
}
