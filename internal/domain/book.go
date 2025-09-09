package domain

import (
	"time"
)

// Book adalah entitas untuk data buku cerita
type Book struct {
	ID                 string  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title              string  `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description        string  `gorm:"type:text"`
	CoverImageURL      string  `gorm:"type:varchar(255)"`
	Level              int     `gorm:"type:smallint;not null"`
	UnlockDependencyID *string `gorm:"type:uuid"` // Pointer karena bisa NULL untuk buku level 1

	// Relasi
	Pages          []Page         `gorm:"foreignKey:BookID"`
	PostActivities []PostActivity `gorm:"foreignKey:BookID"`

	// Foreign key constraint (opsional, tapi baik untuk kejelasan)
	UnlockDependency *Book `gorm:"foreignKey:UnlockDependencyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
