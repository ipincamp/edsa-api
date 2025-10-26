package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUserGameScoresTable() *gormigrate.Migration {
	type User struct {
		ID uuid.UUID `gorm:"type:uuid;primarykey"`
	}
	type Game struct {
		ID uint `gorm:"primarykey"`
	}

	type UserGameScore struct {
		ID           uint      `gorm:"primarykey"`
		UserID       uuid.UUID `gorm:"not null;uniqueIndex:idx_user_game"`
		User         User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		GameID       uint      `gorm:"not null;uniqueIndex:idx_user_game"`
		Game         Game      `gorm:"foreignKey:GameID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
		HighestScore int       `gorm:"default:0"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020151922",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&UserGameScore{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&UserGameScore{})
		},
	}
}
