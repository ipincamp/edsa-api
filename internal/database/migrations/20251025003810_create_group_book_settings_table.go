package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateGroupBookSettingsTable() *gormigrate.Migration {
	type Group struct {
		ID uint `gorm:"primarykey"`
	}
	type Book struct {
		ID uint `gorm:"primarykey"`
	}

	type GroupBookSetting struct {
		ID         uint  `gorm:"primarykey"`
		GroupID    uint  `gorm:"not null;uniqueIndex:idx_group_book"`
		Group      Group `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		BookID     uint  `gorm:"not null;uniqueIndex:idx_group_book"`
		Book       Book  `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		IsUnlocked bool  `gorm:"default:false"`
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeletedAt  gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251025003810",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&GroupBookSetting{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("group_book_settings")
		},
	}
}
