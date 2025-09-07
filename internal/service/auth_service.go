package service

import (
	"errors"
	"time"

	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/repository"
	"github.com/ipincamp/go-edsa-api/pkg/bloom"
	"github.com/ipincamp/go-edsa-api/pkg/cache"
	"github.com/ipincamp/go-edsa-api/pkg/hash"
	"github.com/ipincamp/go-edsa-api/pkg/token"
	"gorm.io/gorm"
)

var (
	ErrEmailExists          = errors.New("email already exists")
	ErrDefaultRoleNotFound  = errors.New("default role not found in cache")
	ErrHashingPassword      = errors.New("error hashing password")
	ErrGenerateAccessToken  = errors.New("error generating access token")
	ErrGenerateRefreshToken = errors.New("error generating refresh token")
	ErrUserCreation         = errors.New("error creating user")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrCheckCredentials     = errors.New("failed check credentials")
	ErrInvalidToken         = errors.New("invalid or expired refresh token")
)

type authResponse struct {
	User  *domain.User
	Token *tokenResponse
}

type tokenResponse struct {
	Access  string
	Refresh string
}

type AuthService interface {
	Register(user *domain.User) (authResponse, error)
	Login(email, password string) (authResponse, error)
	RefreshToken(token string) (tokenResponse, error)
	Logout(userID string) error
}

type authService struct {
	db                   *gorm.DB
	userRepo             repository.UserRepository
	roleRepo             repository.RoleRepository
	tokenMaker           token.PasetoMaker
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	bloomFilter          *bloom.BloomFilterManager
	hashParams           hash.Argon2Params
}

func NewAuthService(db *gorm.DB, userRepo repository.UserRepository, roleRepo repository.RoleRepository, tokenMaker token.PasetoMaker, accessTokenDuration time.Duration, refreshTokenDuration time.Duration, bloomFilter *bloom.BloomFilterManager) AuthService {
	return &authService{
		db:                   db,
		userRepo:             userRepo,
		roleRepo:             roleRepo,
		tokenMaker:           tokenMaker,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		bloomFilter:          bloomFilter,
		hashParams: hash.Argon2Params{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}

func (s *authService) Register(user *domain.User) (authResponse, error) {
	if s.bloomFilter.Test(user.Email) {
		return authResponse{}, ErrEmailExists
	}

	guestRole, found := cache.GetRoleByName(constant.RoleGuest.String())
	if !found {
		return authResponse{}, ErrDefaultRoleNotFound
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := hash.CreateHash(user.Password, &s.hashParams)
		if err != nil {
			return ErrHashingPassword
		}
		user.Password = hashedPassword
		user.RoleID = guestRole.ID

		txUserRepo := repository.NewUserRepository(tx)
		if err := txUserRepo.Create(user); err != nil {
			return ErrUserCreation
		}

		return nil
	})

	if err != nil {
		return authResponse{}, err
	}

	s.bloomFilter.Add(user.Email)
	go s.bloomFilter.Save()

	user.Role = guestRole
	cache.AddUserToCache(*user)

	// generate token
	accessToken, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Role.ID,
		"access",
		s.accessTokenDuration,
	)
	if err != nil {
		return authResponse{}, ErrGenerateAccessToken
	}
	refreshToken, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Role.ID,
		"refresh",
		s.refreshTokenDuration,
	)
	if err != nil {
		return authResponse{}, ErrGenerateRefreshToken
	}

	response := authResponse{
		User: user,
		Token: &tokenResponse{
			Access:  accessToken,
			Refresh: refreshToken,
		},
	}

	return response, nil
}

func (s *authService) Login(email, password string) (authResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return authResponse{}, ErrUserNotFound
	}

	match, err := hash.ComparePasswordAndHash(password, user.Password)
	if err != nil {
		return authResponse{}, ErrCheckCredentials
	}
	if !match {
		return authResponse{}, ErrInvalidCredentials
	}

	accessToken, err := s.tokenMaker.CreateToken(
		user.ID,
		user.RoleID,
		"access",
		s.accessTokenDuration,
	)
	if err != nil {
		return authResponse{}, ErrGenerateAccessToken
	}
	refreshToken, err := s.tokenMaker.CreateToken(
		user.ID,
		user.RoleID,
		"refresh",
		s.refreshTokenDuration,
	)
	if err != nil {
		return authResponse{}, ErrGenerateRefreshToken
	}

	return authResponse{
		User: user,
		Token: &tokenResponse{
			Access:  accessToken,
			Refresh: refreshToken,
		},
	}, nil
}

func (s *authService) RefreshToken(token string) (tokenResponse, error) {
	payload, err := s.tokenMaker.VerifyToken(token)
	if err != nil {
		return tokenResponse{}, ErrInvalidToken
	}
	if payload.TokenType != "refresh" {
		return tokenResponse{}, ErrInvalidToken
	}

	accessToken, err := s.tokenMaker.CreateToken(
		payload.UserID,
		payload.RoleID,
		"access",
		s.accessTokenDuration,
	)
	if err != nil {
		return tokenResponse{}, ErrGenerateAccessToken
	}
	refreshToken, err := s.tokenMaker.CreateToken(
		payload.UserID,
		payload.RoleID,
		"refresh",
		s.refreshTokenDuration,
	)
	if err != nil {
		return tokenResponse{}, ErrGenerateRefreshToken
	}

	return tokenResponse{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (s *authService) Logout(userID string) error {
	// TODO: Implement token revocation or blacklisting if necessary
	return nil
}
