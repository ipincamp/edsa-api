package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddRoleIdToUsersTable() *gormigrate.Migration {
	type User struct {
		RoleID string `gorm:"type:uuid"`
	}

	return &gormigrate.Migration{
		ID: "20250830141645",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&User{}, "role_id")
		},
	}
}
