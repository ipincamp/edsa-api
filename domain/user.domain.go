package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

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
	List(ctx context.Context, limit, offset int, roleName string) ([]User, int64, error)
	FindByID(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
