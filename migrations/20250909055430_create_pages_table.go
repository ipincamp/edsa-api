package migrations

import (
	"encoding/json"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreatePagesTable() *gormigrate.Migration {
	type Book struct {
		ID string `gorm:"type:uuid;primary_key"`
	}

	type Page struct {
		ID             string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		BookID         string          `gorm:"type:uuid;not null;index"`
		PageNumber     int             `gorm:"not null"`
		ContentType    string          `gorm:"type:varchar(255);not null"`
		ContentData    json.RawMessage `gorm:"type:jsonb"`
		HasInteraction bool            `gorm:"not null;default:false"`
		Book           Book            `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	return &gormigrate.Migration{
		ID: "20250909055430",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Page{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("pages")
		},
	}
}
