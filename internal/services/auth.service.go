package services

import (
	"errors"
	"fmt"

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
	userRepo   repositories.UserRepository
	emailCache repositories.EmailCache
}

func NewAuthService(userRepo repositories.UserRepository, emailCache repositories.EmailCache) AuthService {
	return &authService{userRepo, emailCache}
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	if s.emailCache.EmailExists(req.Email) {
		return errors.New("email already exists")
	}
	if worker.IsEmailBeingProcessed(req.Email) {
		return errors.New("registration for this email is already in progress")
	}

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
	if worker.IsEmailBeingProcessed(req.Email) {
		return nil, errors.New("your account is still being processed, please try again in a moment")
	}

	if !s.emailCache.EmailExists(req.Email) {
		if reason, failed := worker.GetRegistrationFailureReason(req.Email); failed {
			return nil, fmt.Errorf("your account registration failed: %s Please try to register again", reason)
		}
		return nil, errors.New("invalid credentials")
	}

	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Status != models.StatusActive {
		return nil, fmt.Errorf("your account is not active (status: %s)", user.Status)
	}

	match, err := utils.CheckPasswordHash(req.Password, user.Password)
	if err != nil || !match {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
