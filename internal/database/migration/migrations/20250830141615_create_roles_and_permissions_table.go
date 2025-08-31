package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateRolesAndPermissionsTable() *gormigrate.Migration {
	type Role struct {
		ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Name      string `gorm:"type:varchar(50);uniqueIndex;not null"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	type Permission struct {
		ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Name      string `gorm:"type:varchar(50);uniqueIndex;not null"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	type RolePermission struct {
		RoleID       string `gorm:"type:uuid;primaryKey"`
		PermissionID string `gorm:"type:uuid;primaryKey"`
	}

	return &gormigrate.Migration{
		ID: "20250830141615",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Role{}, &Permission{}, &RolePermission{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("roles", "permissions", "role_permissions")
		},
	}
}
