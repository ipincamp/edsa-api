package domain

import (
	"time"

	"gorm.io/gorm"
)

// GroupBookSetting adalah entitas domain inti
type GroupBookSetting struct {
	ID         uint
	GroupID    uint
	BookID     uint
	IsUnlocked bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt

	// Relasi
	Group Group
	Book  Book
}
