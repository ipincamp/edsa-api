package migrations

import (
	"encoding/json"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateActivityLogsTable() *gormigrate.Migration {
	type User struct {
		ID uuid.UUID `gorm:"type:uuid;primarykey"`
	}

	type ActivityLog struct {
		ID             uint            `gorm:"primarykey"`
		UserID         uuid.UUID       `gorm:"not null;index"`
		User           User            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		SessionID      uuid.UUID       `gorm:"not null;index"`
		Action         string          `gorm:"type:varchar(100);not null;index"`
		TimestampStart time.Time       `gorm:"not null"`
		DurationMs     *int            // Pointer int agar bisa nullable
		Details        json.RawMessage `gorm:"type:jsonb"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
		DeletedAt      gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020155835",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ActivityLog{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("activity_logs")
		},
	}
}
