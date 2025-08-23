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
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	existingUser, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	job := worker.RegistrationJob{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	worker.QueueRegistrationJob(job)

	return nil
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
