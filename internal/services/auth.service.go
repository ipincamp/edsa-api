package services

import (
	"errors"

	"github.com/ipincamp/edsa/internal/api/dto"
	"github.com/ipincamp/edsa/internal/database"
	"github.com/ipincamp/edsa/internal/models"
	"github.com/ipincamp/edsa/internal/repositories"
	"github.com/ipincamp/edsa/internal/utils"
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
	var newUser *models.User

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		txUserRepo := s.userRepo.WithTx(tx)

		existingUser, err := txUserRepo.FindUserByEmail(req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existingUser != nil {
			return errors.New("email already exists")
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}

		userToCreate := &models.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: hashedPassword,
			Status:   models.StatusActive,
		}

		if err := txUserRepo.CreateUser(userToCreate); err != nil {
			return err
		}

		newUser = userToCreate
		return nil
	})

	if err != nil {
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

	if user.Status != models.StatusActive {
		return nil, errors.New("your account is not active")
	}

	match, err := utils.CheckPasswordHash(req.Password, user.Password)
	if err != nil || !match {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
