package services

import (
	"errors"
	"fmt"

	"github.com/ipincamp/edsa/internal/api/dto"
	"github.com/ipincamp/edsa/internal/batcher"
	"github.com/ipincamp/edsa/internal/models"
	"github.com/ipincamp/edsa/internal/repositories"
	"github.com/ipincamp/edsa/internal/utils"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*models.User, error)
}

type authService struct {
	userRepo   repositories.UserRepository
	emailCache repositories.EmailCache
	processor  *batcher.Processor
}

func NewAuthService(userRepo repositories.UserRepository, emailCache repositories.EmailCache, processor *batcher.Processor) AuthService {
	return &authService{userRepo, emailCache, processor}
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	// Cek status email langsung di cache.
	status, exists := s.emailCache.GetEmailStatus(req.Email)

	if exists {
		if status == repositories.StatusProcessing {
			return errors.New("registration for this email is already in progress")
		}
		if status == repositories.StatusRegistered {
			return errors.New("email already exists")
		}
	}

	// Jika tidak ada di cache, tambahkan ke antrian dan set status ke "processing".
	s.processor.AddToQueue(batcher.RegistrationRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	s.emailCache.SetEmailStatus(req.Email, repositories.StatusProcessing)

	return nil
}

func (s *authService) Login(req *dto.LoginRequest) (*models.User, error) {
	status, exists := s.emailCache.GetEmailStatus(req.Email)

	if !exists {
		return nil, errors.New("invalid credentials")
	}

	if status == repositories.StatusProcessing {
		return nil, errors.New("your account registration is still being processed, please try again in a few minutes")
	}

	// Jika status "registered", baru ambil data dari database.
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
