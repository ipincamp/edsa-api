package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUserBookProgressTable() *gormigrate.Migration {
	type User struct {
		ID uuid.UUID `gorm:"type:uuid;primarykey"`
	}
	type Book struct {
		ID uint `gorm:"primarykey"`
	}

	type UserBookProgress struct {
		ID                   uint      `gorm:"primarykey"`
		UserID               uuid.UUID `gorm:"not null;uniqueIndex:idx_user_book"`
		User                 User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		BookID               uint      `gorm:"not null;uniqueIndex:idx_user_book"`
		Book                 Book      `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
		Status               string    `gorm:"type:varchar(50);default:'locked'"` // 'locked', 'unlocked', 'completed'
		HighestScore         int       `gorm:"default:0"`
		LastPageID           uint      `gorm:"default:0"` // ID dari tabel 'pages'
		CurrentSessionPoints int       `gorm:"default:0"`
		CreatedAt            time.Time
		UpdatedAt            time.Time
		DeletedAt            gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020144957",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&UserBookProgress{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("user_book_progress")
		},
	}
}
