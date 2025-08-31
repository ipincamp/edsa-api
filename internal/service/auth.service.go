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
	var res dto.AuthResponse
	var newUser domain.User
	var defaultRole domain.Role

	err := s.db.Transaction(func(tx *gorm.DB) error {
		_, err := s.userRepository.FindByEmail(ctx, tx, request.Email)
		if err == nil {
			return errors.New("email already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		hashedPassword, err := util.HashPassword(request.Password)
		if err != nil {
			return err
		}

		if err := tx.WithContext(ctx).
			Where("name = ?", constant.RoleStudent).
			First(&defaultRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("default user role not found")
			}
			return err
		}

		newUser = domain.User{
			Name:     request.Name,
			Email:    request.Email,
			Password: hashedPassword,
			RoleID:   defaultRole.ID,
		}

		if err := s.userRepository.Save(ctx, tx, &newUser); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.AuthResponse{}, err
	}

	newUser.Role = defaultRole
	tokenTTL := time.Duration(s.config.Paseto.TokenTTLMin) * time.Minute
	pasetoMaker, err := util.NewPasetoMaker(s.config.Paseto.SecretKey)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	token, err := pasetoMaker.CreateToken(newUser.ID, newUser.Role.ID, tokenTTL)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	res = dto.AuthResponse{
		Token: token,
		User: dto.UserData{
			ID:        newUser.ID,
			Name:      newUser.Name,
			Email:     newUser.Email,
			Role:      newUser.Role.Name,
			JoinedAt:  newUser.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: newUser.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	return res, nil
}

func (s *authService) Login(ctx context.Context, request dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := s.userRepository.FindByEmail(ctx, s.db, request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AuthResponse{}, errors.New("invalid credentials")
		}
		return dto.AuthResponse{}, err
	}

	match, err := util.CheckPasswordHash(request.Password, user.Password)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if !match {
		return dto.AuthResponse{}, errors.New("invalid credentials")
	}

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

func (s *authService) Logout(ctx context.Context, user domain.User) error {
	return nil
}
