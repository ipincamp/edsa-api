package dto

import (
	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/internal/constant"
)

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role.Name,
		JoinedAt:  user.CreatedAt.Format(constant.TimeFormat),
		UpdatedAt: user.UpdatedAt.Format(constant.TimeFormat),
	}
}

func ToUserListResponse(users []domain.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, user := range users {
		result[i] = ToUserResponse(user)
	}
	return result
}
