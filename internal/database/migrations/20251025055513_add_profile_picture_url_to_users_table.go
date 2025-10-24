package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddProfilePictureUrlToUsersTable() *gormigrate.Migration {
	type User struct {
		ID                uuid.UUID `gorm:"type:uuid;primarykey"`
		ProfilePictureURL string    `gorm:"type:varchar(255);default:''"`
	}

	return &gormigrate.Migration{
		ID: "20251025055513",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&User{})
		},
		Rollback: func(tx *gorm.DB) error {
			type User struct {
				ID uuid.UUID `gorm:"type:uuid;primarykey"`
			}
			return tx.Migrator().DropColumn(&User{}, "profile_picture_url")
		},
	}
}
