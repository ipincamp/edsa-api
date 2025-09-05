package service

import (
	"context"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/domain/dto"
)

type UserService interface {
	All(ctx context.Context, page int, limit int, roleName string) (*dto.PaginatedResponse, error)
	Profile(ctx context.Context, userID string) (dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userID string, request dto.UpdateProfileUserRequest) error
}

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) (dto.AuthResponse, error)
	Login(ctx context.Context, request dto.LoginRequest) (dto.AuthResponse, error)
	Logout(ctx context.Context, user domain.User) error
}
