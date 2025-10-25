package migrations

import (
	"encoding/json"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUserInteractionAttemptsTable() *gormigrate.Migration {
	type User struct {
		ID uuid.UUID `gorm:"type:uuid;primarykey"`
	}
	type Interaction struct {
		ID uint `gorm:"primarykey"`
	}

	type UserInteractionAttempt struct {
		ID              uint            `gorm:"primarykey"`
		UserID          uuid.UUID       `gorm:"type:uuid;not null;index"`
		User            User            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		InteractionID   uint            `gorm:"not null;index"`
		Interaction     Interaction     `gorm:"foreignKey:InteractionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		Timestamp       time.Time       `gorm:"not null"`      // Waktu percobaan
		UserAnswer      json.RawMessage `gorm:"type:jsonb"`    // Jawaban user (teks, JSON, dll)
		IsCorrect       bool            `gorm:"default:false"` // Apakah jawaban benar
		ScoreAwarded    float64         `gorm:"default:0"`     // Skor yang didapat
		DurationSeconds int             `gorm:"default:0"`     // Durasi pengerjaan (detik)
		CreatedAt       time.Time
		UpdatedAt       time.Time
		DeletedAt       gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251025190637",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&UserInteractionAttempt{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("user_interaction_attempts")
		},
	}
}
