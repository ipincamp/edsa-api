package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreatePagesTable() *gormigrate.Migration {
	type Book struct {
		ID uint `gorm:"primarykey"`
	}

	type Page struct {
		ID              uint   `gorm:"primarykey"`
		BookID          uint   `gorm:"not null;index"`
		Book            Book   `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		PageNumber      int    `gorm:"not null;index"`
		NarrativeText   string `gorm:"type:text"`
		InstructionText string `gorm:"type:text"`
		CreatedAt       time.Time
		UpdatedAt       time.Time
		DeletedAt       gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020135041",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Page{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("pages")
		},
	}
}
