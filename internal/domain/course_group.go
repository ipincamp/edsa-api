package domain

import (
	"time"

	"gorm.io/gorm"
)

type CourseGroup struct {
	ID               string             `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name             string             `gorm:"type:varchar(255);not null"`
	GroupCode        string             `gorm:"type:varchar(50);uniqueIndex;not null"`
	GroupDescription string             `gorm:"type:text"`
	CourseID         string             `gorm:"type:uuid;not null"`
	Course           Course             // Relasi many-to-one ke Course
	Enrollments      []Enrollment       // Relasi one-to-many ke Enrollment
	ClassTeachers    []ClassTeacher     // Relasi one-to-many ke ClassTeacher
	JoinRequests     []JoinGroupRequest `gorm:"foreignKey:TargetGroupID"` // Relasi one-to-many
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
