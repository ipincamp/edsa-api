package services

import (
	"errors"

	"github.com/ipincamp/edsa/internal/api/dto"
	"github.com/ipincamp/edsa/internal/models"
	"github.com/ipincamp/edsa/internal/repositories"
	"github.com/ipincamp/edsa/internal/utils"
	"github.com/ipincamp/edsa/internal/worker"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*models.User, error)
	Login(req *dto.LoginRequest) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(req *dto.RegisterRequest) (*models.User, error) {
	existingUser, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := worker.HashPasswordAsync(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) Login(req *dto.LoginRequest) (*models.User, error) {
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	match, err := utils.CheckPasswordHash(req.Password, user.Password)
	if err != nil || !match {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
