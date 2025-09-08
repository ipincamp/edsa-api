package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

// CourseGroupRepository adalah kontrak untuk akses data grup kursus/kelas
type CourseGroupRepository interface {
	Create(ctx context.Context, courseGroup *domain.CourseGroup) error
	FindByCode(ctx context.Context, code string) (*domain.CourseGroup, error)
	FindByID(ctx context.Context, id string) (*domain.CourseGroup, error)
	GetByTeacherID(ctx context.Context, teacherID string) ([]domain.CourseGroup, error)
}

type courseGroupRepository struct {
	db *gorm.DB
}

// NewCourseGroupRepository membuat instance baru courseGroupRepository
func NewCourseGroupRepository(db *gorm.DB) CourseGroupRepository {
	return &courseGroupRepository{db: db}
}

// Create membuat kelas baru
func (r *courseGroupRepository) Create(ctx context.Context, courseGroup *domain.CourseGroup) error {
	return r.db.WithContext(ctx).Create(courseGroup).Error
}

// FindByCode untuk mencari group berdasarkan kode join
func (r *courseGroupRepository) FindByCode(ctx context.Context, code string) (*domain.CourseGroup, error) {
	var courseGroup domain.CourseGroup
	err := r.db.WithContext(ctx).Where("group_code = ?", code).First(&courseGroup).Error
	return &courseGroup, err
}

// FindByID mencari kelas berdasarkan ID
func (r *courseGroupRepository) FindByID(ctx context.Context, id string) (*domain.CourseGroup, error) {
	var courseGroup domain.CourseGroup
	// Preload ClassTeachers untuk mendapatkan info guru dan kapasitas
	err := r.db.WithContext(ctx).Preload("ClassTeachers").Where("id = ?", id).First(&courseGroup).Error
	return &courseGroup, err
}

// GetByTeacherID mendapatkan semua kelas yang diajar oleh seorang guru
func (r *courseGroupRepository) GetByTeacherID(ctx context.Context, teacherID string) ([]domain.CourseGroup, error) {
	var groups []domain.CourseGroup
	err := r.db.WithContext(ctx).
		Joins("JOIN class_teachers ON class_teachers.course_group_id = course_groups.id").
		Where("class_teachers.teacher_id = ?", teacherID).
		Find(&groups).Error
	return groups, err
}
