package gorm

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type pageRepositoryGORM struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) usecase.PageRepository {
	return &pageRepositoryGORM{db: db}
}

func (r *pageRepositoryGORM) Create(ctx context.Context, page *domain.Page) error {
	gormPage := PageFromDomain(page)
	result := r.db.WithContext(ctx).Create(gormPage)
	if result.Error != nil {
		return result.Error
	}
	page.ID = gormPage.ID
	return nil
}

func (r *pageRepositoryGORM) FindAllByBookID(ctx context.Context, bookID uint) ([]domain.Page, error) {
	var gormPages []PageGORM
	if err := r.db.WithContext(ctx).
		Preload("Book").
		Where("book_id = ?", bookID).
		Order("page_number asc").
		Find(&gormPages).Error; err != nil {
		return nil, err
	}

	var domainPages []domain.Page
	for _, p := range gormPages {
		domainPages = append(domainPages, *p.ToDomain())
	}
	return domainPages, nil
}

func (r *pageRepositoryGORM) FindByID(ctx context.Context, id uint) (*domain.Page, error) {
	var gormPage PageGORM
	result := r.db.WithContext(ctx).Preload("Book").First(&gormPage, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormPage.ToDomain(), nil
}

func (r *pageRepositoryGORM) Update(ctx context.Context, page *domain.Page) error {
	gormPage := PageFromDomain(page)
	return r.db.WithContext(ctx).Save(gormPage).Error
}

func (r *pageRepositoryGORM) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&PageGORM{}, id).Error
}
