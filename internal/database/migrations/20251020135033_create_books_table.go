package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateBooksTable() *gormigrate.Migration {
	type Book struct {
		ID            uint   `gorm:"primarykey"`
		Title         string `gorm:"type:varchar(255);not null"`
		Description   string `gorm:"type:text"`
		CoverImageURL string `gorm:"type:varchar(255)"`
		Theme         string `gorm:"type:varchar(100)"`
		BookOrder     int    `gorm:"default:0;uniqueIndex"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
		DeletedAt     gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020135033",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Book{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("books")
		},
	}
}
