package gorm

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type classRepositoryGORM struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) usecase.ClassRepository {
	return &classRepositoryGORM{db: db}
}

func (r *classRepositoryGORM) Create(ctx context.Context, class *domain.Class) error {
	gormClass := ClassFromDomain(class)
	result := r.db.WithContext(ctx).Create(gormClass)
	if result.Error != nil {
		return result.Error
	}
	class.ID = gormClass.ID
	return nil
}

func (r *classRepositoryGORM) FindAll(ctx context.Context) ([]domain.Class, error) {
	var gormClasses []ClassGORM
	if err := r.db.WithContext(ctx).Preload("Subject").Find(&gormClasses).Error; err != nil {
		return nil, err
	}

	var domainClasses []domain.Class
	for _, c := range gormClasses {
		domainClasses = append(domainClasses, *c.ToDomain())
	}
	return domainClasses, nil
}

func (r *classRepositoryGORM) FindByID(ctx context.Context, id uint) (*domain.Class, error) {
	var gormClass ClassGORM
	result := r.db.WithContext(ctx).Preload("Subject").First(&gormClass, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormClass.ToDomain(), nil
}

func (r *classRepositoryGORM) Update(ctx context.Context, class *domain.Class) error {
	gormClass := ClassFromDomain(class)
	return r.db.WithContext(ctx).Save(gormClass).Error
}

func (r *classRepositoryGORM) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&ClassGORM{}, id).Error
}
