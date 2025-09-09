package domain

import (
	"time"
)

// UserProgress melacak progres belajar seorang user pada sebuah buku
type UserProgress struct {
	ID                  string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID              string     `gorm:"type:uuid;not null;index"`
	BookID              string     `gorm:"type:uuid;not null;index"`
	LastCompletedPageID *string    `gorm:"type:uuid;index"` // Pointer karena bisa NULL
	Status              string     `gorm:"type:varchar(50);not null;default:'not_started';check:status IN ('not_started', 'in_progress', 'completed')"`
	BookScore           int        `gorm:"not null;default:0"`
	CompletedAt         *time.Time // Pointer karena bisa NULL

	// Relasi
	User              User `gorm:"foreignKey:UserID"`
	Book              Book `gorm:"foreignKey:BookID"`
	LastCompletedPage Page `gorm:"foreignKey:LastCompletedPageID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
