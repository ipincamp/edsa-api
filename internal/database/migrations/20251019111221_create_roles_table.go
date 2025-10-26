package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateRolesTable() *gormigrate.Migration {
	type Role struct {
		ID        uint   `gorm:"primarykey"`
		Name      string `gorm:"type:varchar(50);uniqueIndex;not null"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251019111221",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Role{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Role{})
		},
	}
}
