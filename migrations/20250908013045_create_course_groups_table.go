package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateCourseGroupsTable() *gormigrate.Migration {
	type User struct {
		ID string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	}
	type Course struct {
		ID string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	}

	type CourseGroup struct {
		ID               string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Name             string `gorm:"type:varchar(255);not null"`
		GroupCode        string `gorm:"type:varchar(50);uniqueIndex;not null"`
		GroupDescription string `gorm:"type:text"`
		CourseID         string `gorm:"type:uuid;not null"`
		Course           Course `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt        time.Time
		UpdatedAt        time.Time
		DeletedAt        gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20250908013045",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&CourseGroup{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("course_groups")
		},
	}
}
