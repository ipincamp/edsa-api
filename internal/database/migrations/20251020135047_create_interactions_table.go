package migrations

import (
	"encoding/json"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateInteractionsTable() *gormigrate.Migration {
	type Page struct {
		ID uint `gorm:"primarykey"`
	}

	type Interaction struct {
		ID        uint            `gorm:"primarykey"`
		PageID    uint            `gorm:"not null;index"`
		Page      Page            `gorm:"foreignKey:PageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		Type      string          `gorm:"type:varchar(50);not null"` // "Tap", "DragDrop", "Speaking", "Quiz"
		Config    json.RawMessage `gorm:"type:jsonb"`                // Menyimpan parameter (jawaban benar, skor, dll)
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020135047",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Interaction{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Interaction{})
		},
	}
}
