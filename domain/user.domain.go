package domain

import (
	"context"
	"time"

	"github.com/ipincamp/go-edsa-api/domain/dto"
	"gorm.io/gorm"
)

type Role struct {
	ID          string       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string       `gorm:"type:varchar(50);uniqueIndex;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Permission struct {
	ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name      string `gorm:"type:varchar(50);uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID              string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name            string `gorm:"type:varchar(100);not null"`
	Email           string `gorm:"type:varchar(255);uniqueIndex;not null"`
	EmailVerifiedAt *time.Time
	Password        string `gorm:"type:varchar(255);not null" json:"-"`
	RoleID          string `gorm:"type:uuid"`
	Role            Role
	Permissions     []Permission `gorm:"many2many:user_permissions;"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	FindByID(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	Save(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

type UserService interface {
	GetAll(ctx context.Context) ([]dto.UserData, error)
	Profile(ctx context.Context, userID string) (dto.UserData, error)
	Update(ctx context.Context, userID string, request dto.UpdateUserRequest) error
}

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) (dto.AuthResponse, error)
	Login(ctx context.Context, request dto.LoginRequest) (dto.AuthResponse, error)
	Logout(ctx context.Context, user User) error
}
