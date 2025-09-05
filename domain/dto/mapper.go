package dto

import (
	"github.com/ipincamp/go-edsa-api/domain"
)

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role.Name,
		JoinedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToUserListResponse(users []domain.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, user := range users {
		result[i] = UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			JoinedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return result
}
