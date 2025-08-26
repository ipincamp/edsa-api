package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ipincamp/edsa/internal/api/dto"
	"github.com/ipincamp/edsa/internal/batcher"
	"github.com/ipincamp/edsa/internal/mailer"
	"github.com/ipincamp/edsa/internal/models"
	"github.com/ipincamp/edsa/internal/repositories"
	"github.com/ipincamp/edsa/internal/utils"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	VerifyEmail(token string) error
	Login(req *dto.LoginRequest) (*models.User, error)
}

type authService struct {
	userRepo   repositories.UserRepository
	emailCache repositories.EmailCache
	processor  *batcher.Processor
	mailer     mailer.Mailer
}

func NewAuthService(userRepo repositories.UserRepository, emailCache repositories.EmailCache, processor *batcher.Processor, mailer mailer.Mailer) AuthService {
	return &authService{userRepo, emailCache, processor, mailer}
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	existingUser, err := s.userRepo.FindUserByEmail(req.Email)
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

	// Buat token verifikasi
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	verificationToken := hex.EncodeToString(tokenBytes)

	newUser := &models.User{
		Name:              req.Name,
		Email:             req.Email,
		Password:          hashedPassword,
		Status:            models.StatusPending,
		VerificationToken: &verificationToken, // Simpan token
	}

	if err := s.userRepo.CreateUser(newUser); err != nil {
		return err
	}

	// Kirim email verifikasi
	go s.mailer.Send(newUser.Email, "templates/verification_email.tmpl", map[string]interface{}{
		"Name":             newUser.Name,
		"VerificationLink": fmt.Sprintf("http://127.0.0.1:5000/api/auth/verify?token=%s", verificationToken),
	})

	return nil
}

func (s *authService) VerifyEmail(token string) error {
	user, err := s.userRepo.FindUserByVerificationToken(token)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	user.Status = models.StatusActive
	user.VerificationToken = nil
	return s.userRepo.UpdateUser(user)
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
