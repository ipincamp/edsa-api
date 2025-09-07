package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateCoursesTable() *gormigrate.Migration {
	type Course struct {
		ID          string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Name        string `gorm:"type:varchar(255);not null;uniqueIndex"`
		Description string `gorm:"type:text"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20250908012930",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Course{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("courses")
		},
	}
}
