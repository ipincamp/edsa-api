package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateClassesTable() *gormigrate.Migration {
	type Subject struct {
		ID uint `gorm:"primarykey"`
	}

	type Class struct {
		ID        uint    `gorm:"primarykey"`
		Name      string  `gorm:"type:varchar(255);not null"`
		SubjectID uint    `gorm:"not null"`
		Subject   Subject `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020113624",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Class{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Class{})
		},
	}
}
