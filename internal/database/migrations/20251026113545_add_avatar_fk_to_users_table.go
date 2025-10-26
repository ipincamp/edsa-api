package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddAvatarFkToUsersTable() *gormigrate.Migration {
	type MediaAsset struct {
		ID uuid.UUID `gorm:"type:uuid;primarykey"`
	}
	type User struct {
		ID               uuid.UUID `gorm:"type:uuid;primarykey"`
		ProfilePictureID *uuid.UUID
		ProfilePicture   MediaAsset `gorm:"foreignKey:ProfilePictureID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

		ProfilePictureURL string `gorm:"type:varchar(255);default:''"`
	}

	return &gormigrate.Migration{
		ID: "20251026113545",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&User{}); err != nil {
				return err
			}
			return tx.Migrator().DropColumn(&User{}, "profile_picture_url")
		},
		Rollback: func(tx *gorm.DB) error {
			// Tambahkan kolom lama ke tabel users
			if !tx.Migrator().HasColumn(&User{}, "profile_picture_url") {
				if err := tx.Migrator().AddColumn(&User{}, "ProfilePictureURL"); err != nil {
					return err
				}
			}
			// Hapus kolom baru dari tabel users
			return tx.Migrator().DropColumn(&User{}, "profile_picture_id")
		},
	}
}
