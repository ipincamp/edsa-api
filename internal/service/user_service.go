package service

import (
	"errors"

	"github.com/ipincamp/go-edsa-api/internal/delivery/http/dto"
	"github.com/ipincamp/go-edsa-api/internal/repository"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"github.com/ipincamp/go-edsa-api/pkg/cache"
	"github.com/ipincamp/go-edsa-api/pkg/hash"
	"gorm.io/gorm"
)

var (
	ErrInvalidOldPassword = errors.New("old password is incorrect")
)

type UserService interface {
	GetAllUsers(page int, limit int, roleName string) (*dto.PaginatedResponse, error)
	GetProfile(userID string) (*dto.UserResponse, error)
	UpdateProfile(userID string, request dto.UpdateProfileUserRequest) error
}

type userService struct {
	db         *gorm.DB
	userRepo   repository.UserRepository
	hashParams hash.Argon2Params
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository) UserService {
	return &userService{
		db:       db,
		userRepo: userRepo,
		hashParams: hash.Argon2Params{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}

func (s *userService) GetAllUsers(page int, limit int, roleName string) (*dto.PaginatedResponse, error) {
	offset := util.CalculateOffset(page, limit)

	users, total, err := s.userRepo.List(limit, offset, roleName)
	if err != nil {
		return nil, err
	}

	userData := dto.ToUserListResponse(users)
	metaData := util.GeneratePagination(page, limit, total)

	return &dto.PaginatedResponse{
		Data: userData,
		Meta: metaData,
	}, nil
}

func (s *userService) GetProfile(userID string) (*dto.UserResponse, error) {
	// Try to get from cache first
	user, found := cache.GetUserFromCacheByID(userID)
	if found {
		response := dto.ToUserResponse(user)
		return &response, nil
	}

	// If not in cache, get from database
	dbUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	response := dto.ToUserResponse(*dbUser)
	return &response, nil
}

func (s *userService) UpdateProfile(userID string, request dto.UpdateProfileUserRequest) error {
	// Get user from database to ensure we have the latest data
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Prepare updates
	hasUpdates := false

	if request.Name != nil && *request.Name != "" {
		user.Name = *request.Name
		hasUpdates = true
	}

	if request.NewPassword != nil && *request.NewPassword != "" {
		if request.OldPassword == nil || *request.OldPassword == "" {
			return ErrInvalidOldPassword
		}

		match, err := hash.ComparePasswordAndHash(*request.OldPassword, user.Password)
		if err != nil {
			return err
		}
		if !match {
			return ErrInvalidOldPassword
		}

		newHashedPassword, err := hash.CreateHash(*request.NewPassword, &s.hashParams)
		if err != nil {
			return err
		}
		user.Password = newHashedPassword
		hasUpdates = true
	}

	if !hasUpdates {
		return nil // No changes to save
	}

	// Update in transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txUserRepo := repository.NewUserRepository(tx)
		return txUserRepo.Update(user)
	})

	if err != nil {
		return err
	}

	// Update cache
	updatedUser, err := s.userRepo.FindByID(userID)
	if err == nil {
		cache.AddUserToCache(*updatedUser)
	}

	return nil
}
