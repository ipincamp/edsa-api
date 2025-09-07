package service

import (
	"context"
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

// authResponse adalah response untuk proses autentikasi
type authResponse struct {
	User  *domain.User
	Token *tokenResponse
}

// tokenResponse adalah response untuk token akses dan refresh
type tokenResponse struct {
	Access  string
	Refresh string
}

// AuthService adalah kontrak untuk service autentikasi
type AuthService interface {
	Register(user *domain.User) (authResponse, error)
	Login(email, password string) (authResponse, error)
	RefreshToken(token string) (tokenResponse, error)
	Logout(userID string) error
}

// authService adalah implementasi AuthService
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

// NewAuthService membuat instance baru authService
func NewAuthService(db *gorm.DB, userRepo repository.UserRepository, roleRepo repository.RoleRepository, tokenMaker token.PasetoMaker, accessTokenDuration time.Duration, refreshTokenDuration time.Duration, bloomFilter *bloom.BloomFilterManager) AuthService {
	return &authService{
		db:                   db,
		userRepo:             userRepo,
		roleRepo:             roleRepo,
		tokenMaker:           tokenMaker,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		bloomFilter:          bloomFilter,
		hashParams:           hash.DefaultArgon2Params,
	}
}

func (s *authService) Register(user *domain.User) (authResponse, error) {
	if s.bloomFilter.Test(user.Email) {
		return authResponse{}, ErrEmailExists
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type result struct {
		resp authResponse
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		guestRole, found := cache.GetRoleByName(constant.RoleGuest.String())
		if !found {
			resultChan <- result{resp: authResponse{}, err: ErrDefaultRoleNotFound}
			return
		}

		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			hashedPassword, err := hash.CreateHash(user.Password, &s.hashParams)
			if err != nil {
				return ErrHashingPassword
			}
			user.Password = hashedPassword
			user.RoleID = guestRole.ID

			txUserRepo := repository.NewUserRepository(tx)
			if err := txUserRepo.Create(ctx, user); err != nil {
				return ErrUserCreation
			}
			return nil
		})
		if err != nil {
			resultChan <- result{resp: authResponse{}, err: err}
			return
		}

		s.bloomFilter.Add(user.Email)
		go s.bloomFilter.Save()

		user.Role = guestRole
		cache.AddUserToCache(*user)

		tokenResp, err := s.createTokens(user.ID, user.Role.ID)
		if err != nil {
			resultChan <- result{resp: authResponse{}, err: err}
			return
		}

		resultChan <- result{resp: authResponse{
			User:  user,
			Token: &tokenResp,
		}, err: nil}
	}()

	select {
	case <-ctx.Done():
		return authResponse{}, ctx.Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

func (s *authService) Login(email, password string) (authResponse, error) {
	if !s.bloomFilter.Test(email) {
		return authResponse{}, ErrUserNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use a channel to handle the database operation with context
	type result struct {
		user *domain.User
		err  error
	}

	resultChan := make(chan result, 1)

	go func() {
		user, err := s.userRepo.FindByEmail(ctx, email)
		resultChan <- result{user: user, err: err}
	}()

	select {
	case <-ctx.Done():
		return authResponse{}, ctx.Err()
	case res := <-resultChan:
		if res.err != nil {
			return authResponse{}, ErrUserNotFound
		}

		match, err := hash.ComparePasswordAndHash(password, res.user.Password)
		if err != nil {
			return authResponse{}, ErrCheckCredentials
		}
		if !match {
			return authResponse{}, ErrInvalidCredentials
		}

		tokenResp, err := s.createTokens(res.user.ID, res.user.RoleID)
		if err != nil {
			return authResponse{}, err
		}
		return authResponse{
			User:  res.user,
			Token: &tokenResp,
		}, nil
	}
}

func (s *authService) RefreshToken(token string) (tokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	type result struct {
		resp tokenResponse
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		payload, err := s.tokenMaker.VerifyToken(token)
		if err != nil {
			resultChan <- result{resp: tokenResponse{}, err: ErrInvalidToken}
			return
		}
		if payload.TokenType != "refresh" {
			resultChan <- result{resp: tokenResponse{}, err: ErrInvalidToken}
			return
		}

		tokenResp, err := s.createTokens(payload.UserID, payload.RoleID)
		resultChan <- result{resp: tokenResp, err: err}
	}()

	select {
	case <-ctx.Done():
		return tokenResponse{}, ctx.Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

// Logout melakukan logout user (implementasi token revocation jika diperlukan)
func (s *authService) Logout(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	type result struct {
		err error
	}
	resultChan := make(chan result, 1)

	go func() {
		// TODO: Implement token revocation or blacklisting if necessary
		resultChan <- result{err: nil}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-resultChan:
		return res.err
	}
}

// createTokens adalah helper untuk membuat access dan refresh token
func (s *authService) createTokens(userID, roleID string) (tokenResponse, error) {
	accessToken, err := s.tokenMaker.CreateToken(
		userID,
		roleID,
		"access",
		s.accessTokenDuration,
	)
	if err != nil {
		return tokenResponse{}, ErrGenerateAccessToken
	}
	refreshToken, err := s.tokenMaker.CreateToken(
		userID,
		roleID,
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
