package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateGamesTable() *gormigrate.Migration {
	type Game struct {
		ID               uint   `gorm:"primarykey"`
		Name             string `gorm:"type:varchar(255);not null"`
		Type             string `gorm:"type:varchar(100)"`
		RelatedBookTheme string `gorm:"type:varchar(100)"`
		CreatedAt        time.Time
		UpdatedAt        time.Time
		DeletedAt        gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020151914",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Game{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("games")
		},
	}
}
