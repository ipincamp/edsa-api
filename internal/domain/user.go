package domain

import (
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
