package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateGroupsTable() *gormigrate.Migration {
	type Class struct {
		ID uint `gorm:"primarykey"`
	}

	type Group struct {
		ID        uint   `gorm:"primarykey"`
		Name      string `gorm:"type:varchar(255);not null"`
		ClassID   uint   `gorm:"not null"`
		Class     Class  `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020114024",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Group{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Group{})
		},
	}
}
