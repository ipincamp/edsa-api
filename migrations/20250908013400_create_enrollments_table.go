package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateEnrollmentsTable() *gormigrate.Migration {
	type User struct {
		ID string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	}
	type CourseGroup struct {
		ID string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	}

	// Junction table for User (as Student) and CourseGroup
	type Enrollment struct {
		ID                string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		StudentID         string `gorm:"type:uuid;not null"`
		CourseGroupID     string `gorm:"type:uuid;not null"`
		HomeroomTeacherID string `gorm:"type:uuid;not null"` // Foreign key to User table
		JoinDate          time.Time
		Student           User        `gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CourseGroup       CourseGroup `gorm:"foreignKey:CourseGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		HomeroomTeacher   User        `gorm:"foreignKey:HomeroomTeacherID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	return &gormigrate.Migration{
		ID: "20250908013400",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Enrollment{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("enrollments")
		},
	}
}
