package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateUserProgressTable() *gormigrate.Migration {
	type User struct {
		ID string `gorm:"type:uuid;primary_key"`
	}
	type Book struct {
		ID string `gorm:"type:uuid;primary_key"`
	}
	type Page struct {
		ID string `gorm:"type:uuid;primary_key"`
	}

	type UserProgress struct {
		ID                  string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		UserID              string     `gorm:"type:uuid;not null;index"`
		BookID              string     `gorm:"type:uuid;not null;index"`
		LastCompletedPageID *string    `gorm:"type:uuid;index"` // Pointer karena bisa NULL
		Status              string     `gorm:"type:varchar(50);not null;check:status IN ('not_started', 'in_progress', 'completed')"`
		BookScore           int        `gorm:"not null;default:0"`
		CompletedAt         *time.Time // Pointer karena bisa NULL
		User                User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		Book                Book       `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		LastCompletedPage   Page       `gorm:"foreignKey:LastCompletedPageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		CreatedAt           time.Time
		UpdatedAt           time.Time
	}

	return &gormigrate.Migration{
		ID: "20250909055732",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&UserProgress{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("user_progress")
		},
	}
}
