package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateJoinGroupRequestsTable() *gormigrate.Migration {
	type User struct {
		ID string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	}
	type CourseGroup struct {
		ID string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	}

	type JoinGroupRequest struct {
		ID              string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		ApplicantUserID string `gorm:"type:uuid;not null"`
		TargetGroupID   string `gorm:"type:uuid;not null"`
		Status          string `gorm:"type:varchar(50);default:'pending';check:status IN ('pending', 'approved', 'rejected')"`
		RequestDate     time.Time
		Applicant       User        `gorm:"foreignKey:ApplicantUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		TargetGroup     CourseGroup `gorm:"foreignKey:TargetGroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}

	return &gormigrate.Migration{
		ID: "20250908013515",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&JoinGroupRequest{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("join_group_requests")
		},
	}
}
