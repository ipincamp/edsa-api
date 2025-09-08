package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

// EnrollmentRepository adalah kontrak untuk data pendaftaran
type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *domain.Enrollment) error
}

type enrollmentRepository struct {
	db *gorm.DB
}

// NewEnrollmentRepository membuat instance baru
func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(ctx context.Context, enrollment *domain.Enrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}
