package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Email     string `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	Password  string `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
