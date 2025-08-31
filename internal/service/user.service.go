package service

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"gorm.io/gorm"
)

type userService struct {
	userRepository domain.UserRepository
}

func NewUser(userRepository domain.UserRepository) domain.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetAll(ctx context.Context) ([]dto.UserData, error) {
	users, err := s.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var userData []dto.UserData
	for _, v := range users {
		userData = append(userData, dto.UserData{
			ID:       v.ID,
			Name:     v.Name,
			Email:    v.Email,
			JoinedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return userData, nil
}

func (s *userService) Profile(ctx context.Context, userID string) (dto.UserData, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserData{}, errors.New("user not found")
		}
		return dto.UserData{}, err
	}

	return dto.UserData{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		JoinedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *userService) Update(ctx context.Context, userID string, request dto.UpdateUserRequest) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	isModified := false

	if request.Name != "" && request.Name != user.Name {
		user.Name = request.Name
		isModified = true
	}

	if request.NewPassword != "" {
		if request.OldPassword == "" {
			return errors.New("old password is required to set a new password")
		}

		match, err := util.CheckPasswordHash(request.OldPassword, user.Password)
		if err != nil {
			return err
		}
		if !match {
			return errors.New("invalid old password")
		}

		newHashedPassword, err := util.HashPassword(request.NewPassword)
		if err != nil {
			return err
		}
		user.Password = newHashedPassword
		isModified = true
	}

	if isModified {
		return s.userRepository.Update(ctx, &user)
	}

	return nil
}
