package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(ctx context.Context, course *domain.Course) error
	FindByID(ctx context.Context, id string) (domain.Course, error)
	FindAll(ctx context.Context, limit, offset int) ([]domain.Course, int64, error)
	Update(ctx context.Context, course *domain.Course) error
	Delete(ctx context.Context, id string) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(ctx context.Context, course *domain.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *courseRepository) FindByID(ctx context.Context, id string) (domain.Course, error) {
	var course domain.Course
	err := r.db.WithContext(ctx).First(&course, "id = ?", id).Error
	return course, err
}

func (r *courseRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Course, int64, error) {
	var courses []domain.Course
	var total int64

	err := r.db.WithContext(ctx).Model(&domain.Course{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&courses).Error
	if err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

func (r *courseRepository) Update(ctx context.Context, course *domain.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

func (r *courseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Course{}, "id = ?", id).Error
}
