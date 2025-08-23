package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus int

const (
	StatusPending  UserStatus = 1
	StatusActive   UserStatus = 2
	StatusInactive UserStatus = 3
	StatusBanned   UserStatus = 4
)

func (s UserStatus) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusInactive:
		return "inactive"
	case StatusBanned:
		return "banned"
	default:
		return "pending"
	}
}

func (s *UserStatus) Scan(value interface{}) error {
	*s = UserStatus(value.(int64))
	return nil
}

func (s UserStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

type User struct {
	ID        string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Status    UserStatus     `gorm:"type:tinyint;not null;default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
