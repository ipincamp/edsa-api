package migrations

import (
	"encoding/json"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateInteractionsTable() *gormigrate.Migration {
	type Page struct {
		ID string `gorm:"type:uuid;primary_key"`
	}

	type Interaction struct {
		ID            string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		PageID        string          `gorm:"type:uuid;not null;index"`
		Type          string          `gorm:"type:varchar(100);not null"`
		ConfigData    json.RawMessage `gorm:"type:jsonb"`
		PointsAwarded int             `gorm:"not null"`
		Page          Page            `gorm:"foreignKey:PageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	return &gormigrate.Migration{
		ID: "20250909055537",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Interaction{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("interactions")
		},
	}
}
