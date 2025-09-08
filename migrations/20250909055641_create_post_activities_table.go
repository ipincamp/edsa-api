package migrations

import (
	"encoding/json"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreatePostActivitiesTable() *gormigrate.Migration {
	type Book struct {
		ID string `gorm:"type:uuid;primary_key"`
	}

	type PostActivity struct {
		ID            string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		BookID        string          `gorm:"type:uuid;not null;index"`
		Type          string          `gorm:"type:varchar(100);not null"`
		ConfigData    json.RawMessage `gorm:"type:jsonb"`
		PointsAwarded int             `gorm:"not null"`
		Book          Book            `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	return &gormigrate.Migration{
		ID: "20250909055641",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&PostActivity{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("post_activities")
		},
	}
}
