package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ipincamp/go-edsa-api/domain/dto"
	"gorm.io/gorm"
)

type Role struct {
	ID          string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string          `gorm:"type:varchar(50);uniqueIndex;not null"`
	Permissions json.RawMessage `gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type User struct {
	ID              string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name            string `gorm:"type:varchar(100);not null"`
	Email           string `gorm:"type:varchar(255);uniqueIndex;not null"`
	EmailVerifiedAt *time.Time
	Password        string `gorm:"type:varchar(255);not null" json:"-"`
	RoleID          string `gorm:"type:uuid"`
	Role            Role
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type UserRepository interface {
	LoadAll(ctx context.Context) ([]User, error)
	FindAll(ctx context.Context, limit int, offset int) ([]User, error)
	Count(ctx context.Context) (int64, error)
	FindByID(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	Save(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

type UserService interface {
	GetAll(ctx context.Context, page, limit int) (*dto.PaginatedResponse, error)
	Profile(ctx context.Context, userID string) (dto.UserData, error)
	Update(ctx context.Context, userID string, request dto.UpdateUserRequest) error
}

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) (dto.AuthResponse, error)
	// TODO send email verification
	// TODO verify redirect email verification
	Login(ctx context.Context, request dto.LoginRequest) (dto.AuthResponse, error)
	// TODO refresh token
	Logout(ctx context.Context, user User) error
}
