package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateSubjectsTable() *gormigrate.Migration {
	type Subject struct {
		ID        uint   `gorm:"primarykey"`
		Name      string `gorm:"type:varchar(255);not null"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020113230",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Subject{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("subjects")
		},
	}
}
