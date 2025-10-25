package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddVerificationToUsersTable() *gormigrate.Migration {
	type User struct {
		ID uuid.UUID `gorm:"type:uuid;primarykey"`

		EmailVerifiedAt *time.Time `gorm:"index"`                       // Pointer ke time.Time akan membuatnya nullable
		IsActive        bool       `gorm:"default:true;not null;index"` // Default true
	}

	return &gormigrate.Migration{
		ID: "20251025192427",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&User{})
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn("users", "email_verified_at"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn("users", "is_active"); err != nil {
				return err
			}
			return nil
		},
	}
}
