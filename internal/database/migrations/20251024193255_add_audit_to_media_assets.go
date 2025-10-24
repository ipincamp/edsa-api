package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddAuditToMediaAssets() *gormigrate.Migration {
	type User struct {
		ID uuid.UUID `gorm:"type:uuid;primarykey"`
	}

	type MediaAsset struct {
		// Kolom baru
		UploadedByUserID *uuid.UUID `gorm:"type:uuid;index"`
		UploadedByUser   User       `gorm:"foreignKey:UploadedByUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

		DeletedByUserID *uuid.UUID `gorm:"type:uuid;index"`
		DeletedByUser   User       `gorm:"foreignKey:DeletedByUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

		// Kolom lama
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251024193255",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&MediaAsset{})
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn(&MediaAsset{}, "uploaded_by_user_id"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&MediaAsset{}, "deleted_by_user_id"); err != nil {
				return err
			}
			// Rollback foreign key constraint biasanya ditangani oleh GORM saat DropColumn
			return nil
		},
	}
}
