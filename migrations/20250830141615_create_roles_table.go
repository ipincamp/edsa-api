package migrations

import (
	"encoding/json"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateRolesTable() *gormigrate.Migration {
	type Role struct {
		ID          string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Name        string          `gorm:"type:varchar(50);uniqueIndex;not null"`
		Permissions json.RawMessage `gorm:"type:jsonb"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	return &gormigrate.Migration{
		ID: "20250830141615",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Role{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("roles")
		},
	}
}
