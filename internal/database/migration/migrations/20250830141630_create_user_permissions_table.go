package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateUserPermissionsTable() *gormigrate.Migration {
	type UserPermission struct {
		UserID       string `gorm:"type:uuid;primaryKey"`
		PermissionID string `gorm:"type:uuid;primaryKey"`
	}

	return &gormigrate.Migration{
		ID: "20250830141630",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&UserPermission{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("user_permissions")
		},
	}
}
