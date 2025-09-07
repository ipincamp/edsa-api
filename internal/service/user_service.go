package service

import (
	"context"
	"errors"
	"time"

	"github.com/ipincamp/go-edsa-api/internal/delivery/http/dto"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/repository"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"github.com/ipincamp/go-edsa-api/pkg/cache"
	"github.com/ipincamp/go-edsa-api/pkg/hash"
	"gorm.io/gorm"
)

var ErrInvalidOldPassword = errors.New("old password is incorrect")

// UserService adalah kontrak untuk service user
type UserService interface {
	GetAllUsers(page int, limit int, roleName string) (*dto.PaginatedResponse, error)
	GetProfile(userID string) (*dto.UserResponse, error)
	UpdateProfile(userID string, request dto.UpdateProfileUserRequest) error
	UpdateUserByID(userID string, request dto.UpdateUserRequest) error
}

// userService adalah implementasi UserService
type userService struct {
	db         *gorm.DB
	userRepo   repository.UserRepository
	hashParams hash.Argon2Params
}

// NewUserService membuat instance baru userService
func NewUserService(db *gorm.DB, userRepo repository.UserRepository) UserService {
	return &userService{
		db:         db,
		userRepo:   userRepo,
		hashParams: hash.DefaultArgon2Params,
	}
}

// GetAllUsers mengambil daftar user dengan paginasi dan filter role
func (s *userService) GetAllUsers(page int, limit int, roleName string) (*dto.PaginatedResponse, error) {
	offset := util.CalculateOffset(page, limit)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type result struct {
		resp *dto.PaginatedResponse
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		users, total, err := s.userRepo.List(ctx, limit, offset, roleName)
		if err != nil {
			resultChan <- result{resp: nil, err: err}
			return
		}
		userData := dto.ToUserListResponse(users)
		metaData := util.GeneratePagination(page, limit, total)
		resultChan <- result{resp: &dto.PaginatedResponse{
			Data: userData,
			Meta: metaData,
		}, err: nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

// UpdateUserByID mengubah data user oleh admin
func (s *userService) UpdateUserByID(userID string, request dto.UpdateUserRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type result struct {
		err error
	}
	resultChan := make(chan result, 1)

	go func() {
		user, err := s.userRepo.FindByID(ctx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				resultChan <- result{err: ErrUserNotFound}
				return
			}
			resultChan <- result{err: err}
			return
		}

		// Prepare updates map
		updates := make(map[string]interface{})

		if request.Name != nil && *request.Name != "" && *request.Name != user.Name {
			updates["name"] = *request.Name
		}

		if request.Email != nil && *request.Email != "" && *request.Email != user.Email {
			updates["email"] = *request.Email
		}

		if request.Role != nil && *request.Role != "" && *request.Role != user.Role.Name {
			// Find role by name to get role_id, check cache first
			var role domain.Role
			cachedRole, found := cache.GetRoleByName(*request.Role)
			if found {
				role = cachedRole
			} else {
				err := s.db.WithContext(ctx).Select("id").Where("name = ?", *request.Role).First(&role).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						resultChan <- result{err: errors.New("role not found")}
						return
					}
					resultChan <- result{err: err}
					return
				}
			}
			updates["role_id"] = role.ID
		}

		if request.Password != nil && *request.Password != "" {
			// Check if new password is same as current password
			samePassword, err := hash.ComparePasswordAndHash(*request.Password, user.Password)
			if err != nil {
				resultChan <- result{err: err}
				return
			}
			if !samePassword {
				newHashedPassword, err := hash.CreateHash(*request.Password, &s.hashParams)
				if err != nil {
					resultChan <- result{err: err}
					return
				}
				updates["password"] = newHashedPassword
			}
		}

		if len(updates) == 0 {
			resultChan <- result{err: nil} // Tidak ada perubahan
			return
		}
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return tx.Model(&domain.User{}).Where("id = ?", userID).Updates(updates).Error
		})
		if err != nil {
			resultChan <- result{err: err}
			return
		}
		updatedUser, err := s.userRepo.FindByID(ctx, userID)
		if err == nil {
			cache.AddUserToCache(*updatedUser)
		}
		resultChan <- result{err: nil}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-resultChan:
		return res.err
	}
}

// GetProfile mengambil data profil user berdasarkan userID
func (s *userService) GetProfile(userID string) (*dto.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type result struct {
		resp *dto.UserResponse
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		user, err := s.userRepo.FindByID(ctx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				resultChan <- result{resp: nil, err: ErrUserNotFound}
				return
			}
			resultChan <- result{resp: nil, err: err}
			return
		}
		response := dto.ToUserResponse(*user)
		resultChan <- result{resp: &response, err: nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

// UpdateProfile mengubah data profil user
func (s *userService) UpdateProfile(userID string, request dto.UpdateProfileUserRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type result struct {
		err error
	}
	resultChan := make(chan result, 1)

	go func() {
		user, err := s.userRepo.FindByID(ctx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				resultChan <- result{err: ErrUserNotFound}
				return
			}
			resultChan <- result{err: err}
			return
		}

		// Prepare updates map
		updates := make(map[string]interface{})

		if request.Name != nil && *request.Name != "" && *request.Name != user.Name {
			updates["name"] = *request.Name
		}

		if request.NewPassword != nil && *request.NewPassword != "" {
			if request.OldPassword == nil || *request.OldPassword == "" {
				resultChan <- result{err: ErrInvalidOldPassword}
				return
			}
			match, err := hash.ComparePasswordAndHash(*request.OldPassword, user.Password)
			if err != nil {
				resultChan <- result{err: err}
				return
			}
			if !match {
				resultChan <- result{err: ErrInvalidOldPassword}
				return
			}

			// Check if new password is same as current password
			samePassword, err := hash.ComparePasswordAndHash(*request.NewPassword, user.Password)
			if err != nil {
				resultChan <- result{err: err}
				return
			}
			if !samePassword {
				newHashedPassword, err := hash.CreateHash(*request.NewPassword, &s.hashParams)
				if err != nil {
					resultChan <- result{err: err}
					return
				}
				updates["password"] = newHashedPassword
			}
		}

		if len(updates) == 0 {
			resultChan <- result{err: nil} // Tidak ada perubahan
			return
		}

		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return tx.Model(&domain.User{}).Where("id = ?", userID).Updates(updates).Error
		})
		if err != nil {
			resultChan <- result{err: err}
			return
		}

		updatedUser, err := s.userRepo.FindByID(ctx, userID)
		if err == nil {
			cache.AddUserToCache(*updatedUser)
		}
		resultChan <- result{err: nil}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-resultChan:
		return res.err
	}
}
