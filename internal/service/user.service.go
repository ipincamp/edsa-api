package service

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"gorm.io/gorm"
)

type userService struct {
	userRepository domain.UserRepository
}

func NewUser(userRepository domain.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) All(ctx context.Context, page int, limit int, roleName string) (*dto.PaginatedResponse, error) {
	offset := util.CalculateOffset(page, limit)

	users, total, err := s.userRepository.List(ctx, limit, offset, roleName)
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

func (s *userService) Profile(ctx context.Context, userID string) (dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, constant.ErrNotFound
		}
		return dto.UserResponse{}, err
	}

	return dto.ToUserResponse(user), nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, request dto.UpdateProfileUserRequest) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constant.ErrNotFound
		}
		return err
	}

	if request.Name != nil && *request.Name != "" {
		user.Name = *request.Name
	}

	if request.NewPassword != nil && *request.NewPassword != "" {
		if request.OldPassword == nil || *request.OldPassword == "" {
			return constant.ErrInvalidInput
		}

		match, err := util.CheckPasswordHash(*request.OldPassword, user.Password)
		if err != nil {
			return err
		}
		if !match {
			return constant.ErrInvalidInput
		}

		newHashedPassword, err := util.HashPassword(*request.NewPassword)
		if err != nil {
			return err
		}
		user.Password = newHashedPassword
	}

	return s.userRepository.Update(ctx, &user)
}
