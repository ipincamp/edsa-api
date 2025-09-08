package dto

import (
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/domain"
)

// UserIDRequest adalah DTO untuk request berdasarkan userId
type UserIDRequest struct {
	UserId string `params:"userId" validate:"required,uuid4"`
}

// UserFilterRequest adalah DTO untuk filter list user
type UserFilterRequest struct {
	Page  int    `query:"page" validate:"omitempty,min=1"`
	Limit int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Role  string `query:"role" validate:"omitempty,oneof=teacher student guest"`
}

// UpdateProfileUserRequest adalah DTO untuk update profile user
type UpdateProfileUserRequest struct {
	Name                    *string `json:"name" validate:"omitempty,min=3"`
	OldPassword             *string `json:"old_password,omitempty" validate:"required_with=NewPassword,omitempty"`
	NewPassword             *string `json:"new_password,omitempty" validate:"omitempty,gte=8"`
	NewPasswordConfirmation *string `json:"new_password_confirmation,omitempty" validate:"required_with=NewPassword,omitempty,eqfield=NewPassword"`
}

// UpdateUserRequest adalah DTO untuk update user oleh admin
type UpdateUserRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=3"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Role     *string `json:"role" validate:"omitempty,oneof=teacher student guest"`
	Password *string `json:"password,omitempty" validate:"omitempty,gte=8"`
}

// UserResponse adalah DTO untuk response user
type UserResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`
	JoinedAt  string `json:"joined_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// UserListResponse adalah DTO sederhana untuk response user
type UserListResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	JoinedAt string `json:"joined_at"`
}

// ToUserResponse mengubah domain.User menjadi UserResponse
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

// ToUserListResponse mengubah slice domain.User menjadi slice UserListResponse
func ToUserListResponse(users []domain.User) []UserListResponse {
	userResponses := make([]UserListResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserListResponse{
			ID:       user.ID,
			Name:     user.Name,
			JoinedAt: user.CreatedAt.Format(constant.TimeFormat),
		}
	}
	return userResponses
}
