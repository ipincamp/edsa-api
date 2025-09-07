package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateClassTeachersTable() *gormigrate.Migration {
	type User struct {
		ID string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	}
	type CourseGroup struct {
		ID string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	}

	// Junction table for User (as Teacher) and CourseGroup
	type ClassTeacher struct {
		ID               string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		TeacherID        string `gorm:"type:uuid;not null"`
		CourseGroupID    string `gorm:"type:uuid;not null"`
		TeachingCapacity int
		User             User        `gorm:"foreignKey:TeacherID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CourseGroup      CourseGroup `gorm:"foreignKey:CourseGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}

	return &gormigrate.Migration{
		ID: "20250908013230",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ClassTeacher{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("class_teachers")
		},
	}
}
