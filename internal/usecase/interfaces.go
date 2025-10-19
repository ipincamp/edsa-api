package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
)

// UserService mendefinisikan logika bisnis untuk pengguna
type UserService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error)
}

// UserRepository mendefinisikan kontrak untuk persistensi data pengguna
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

// PasswordService mendefinisikan kontrak untuk hashing password
type PasswordService interface {
	Hash(password string) (string, error)
	Compare(password, hash string) (bool, error)
}

// TokenService mendefinisikan kontrak untuk pembuatan & validasi token
type TokenService interface {
	CreateToken(user *domain.User) (string, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
}
