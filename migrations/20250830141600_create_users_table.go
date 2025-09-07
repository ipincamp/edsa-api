package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateUsersTable() *gormigrate.Migration {
	type User struct {
		ID              string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Name            string `gorm:"type:varchar(100);not null"`
		Email           string `gorm:"type:varchar(255);uniqueIndex;not null"`
		EmailVerifiedAt *time.Time
		Password        string `gorm:"type:varchar(255);not null"`
		CreatedAt       time.Time
		UpdatedAt       time.Time
		DeletedAt       gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20250830141600",
		Migrate: func(tx *gorm.DB) error {
			tx.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
			return tx.AutoMigrate(&User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")
		},
	}
}
