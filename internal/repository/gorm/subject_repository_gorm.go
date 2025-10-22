package gorm

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type subjectRepositoryGORM struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) usecase.SubjectRepository {
	return &subjectRepositoryGORM{db: db}
}

func (r *subjectRepositoryGORM) Create(ctx context.Context, subject *domain.Subject) error {
	gormSubject := SubjectFromDomain(subject)
	result := r.db.WithContext(ctx).Create(gormSubject)
	if result.Error != nil {
		return result.Error
	}
	subject.ID = gormSubject.ID // Kembalikan ID yang baru dibuat
	return nil
}

func (r *subjectRepositoryGORM) FindAll(ctx context.Context) ([]domain.Subject, error) {
	var gormSubjects []SubjectGORM
	if err := r.db.WithContext(ctx).Find(&gormSubjects).Error; err != nil {
		return nil, err
	}

	var domainSubjects []domain.Subject
	for _, s := range gormSubjects {
		domainSubjects = append(domainSubjects, *s.ToDomain())
	}
	return domainSubjects, nil
}

func (r *subjectRepositoryGORM) FindByID(ctx context.Context, id uint) (*domain.Subject, error) {
	var gormSubject SubjectGORM
	result := r.db.WithContext(ctx).First(&gormSubject, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormSubject.ToDomain(), nil
}

func (r *subjectRepositoryGORM) Update(ctx context.Context, subject *domain.Subject) error {
	gormSubject := SubjectFromDomain(subject)
	return r.db.WithContext(ctx).Save(gormSubject).Error
}

func (r *subjectRepositoryGORM) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&SubjectGORM{}, id).Error
}
