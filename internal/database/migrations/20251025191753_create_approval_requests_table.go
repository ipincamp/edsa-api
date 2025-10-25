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
		// HARUS ada, tidak boleh NULL.
		// Jika user dihapus, request ini juga dihapus (CASCADE).
		UserID uuid.UUID `gorm:"type:uuid;not null;index"`
		User   User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

		// RequestType (VARCHAR)
		RequestType string `gorm:"type:varchar(100);not null;index"` // Cth: "DELETE_ACCOUNT", "REGISTER_USER"

		// Status (VARCHAR)
		Status string `gorm:"type:varchar(50);default:'pending';not null;index"` // "pending", "approved", "rejected"

		// Reason (TEXT)
		Reason string `gorm:"type:text"` // Alasan dari UserID

		// ReviewerID (FK - admin yang review)
		// Boleh NULL (jika belum di-assign).
		// Jika reviewer akan dihapus, database akan MENCEGAH (RESTRICT).
		ReviewerID *uuid.UUID `gorm:"type:uuid;index"`
		Reviewer   User       `gorm:"foreignKey:ReviewerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

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
