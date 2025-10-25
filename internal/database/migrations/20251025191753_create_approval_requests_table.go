package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateApprovalRequestsTable() *gormigrate.Migration {
	type User struct {
		ID uuid.UUID `gorm:"type:uuid;primarykey"`
	}

	// Definisikan tabel baru
	type ApprovalRequest struct {
		ID uint `gorm:"primarykey"`

		// UserID (FK - yang meminta)
		UserID *uuid.UUID `gorm:"type:uuid;index"`
		User   User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

		// RequestType (VARCHAR)
		RequestType string `gorm:"type:varchar(100);not null;index"` // Cth: "DELETE_ACCOUNT", "REGISTER_USER"

		// Status (VARCHAR)
		Status string `gorm:"type:varchar(50);default:'pending';not null;index"` // "pending", "approved", "rejected"

		// Reason (TEXT)
		Reason string `gorm:"type:text"` // Alasan dari UserID

		// ReviewerID (FK - admin yang review)
		ReviewerID *uuid.UUID `gorm:"type:uuid;index"`
		Reviewer   User       `gorm:"foreignKey:ReviewerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

		// Timestamps
		ReviewTimestamp *time.Time `gorm:"index"` // Kapan direview (bisa NULL)
		CreatedAt       time.Time  // Otomatis berfungsi sebagai RequestTimestamp
		UpdatedAt       time.Time
		DeletedAt       gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251025191753",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ApprovalRequest{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("approval_requests")
		},
	}
}
