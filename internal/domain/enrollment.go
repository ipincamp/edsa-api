package domain

import "time"

type Enrollment struct {
	ID                string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	StudentID         string `gorm:"type:uuid;not null"`
	CourseGroupID     string `gorm:"type:uuid;not null"`
	HomeroomTeacherID string `gorm:"type:uuid;not null"`
	JoinDate          time.Time
	Student           User        `gorm:"foreignKey:StudentID"`
	CourseGroup       CourseGroup `gorm:"foreignKey:CourseGroupID"`
	HomeroomTeacher   User        `gorm:"foreignKey:HomeroomTeacherID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
