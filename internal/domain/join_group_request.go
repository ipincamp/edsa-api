package domain

import "time"

type JoinGroupRequest struct {
	ID              string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ApplicantUserID string `gorm:"type:uuid;not null"`
	TargetGroupID   string `gorm:"type:uuid;not null"`
	Status          string `gorm:"type:varchar(50);default:'pending'"`
	RequestDate     time.Time
	Applicant       User        `gorm:"foreignKey:ApplicantUserID"`
	TargetGroup     CourseGroup `gorm:"foreignKey:TargetGroupID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
