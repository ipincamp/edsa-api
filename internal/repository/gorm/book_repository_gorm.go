package gorm

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type bookRepositoryGORM struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) usecase.BookRepository {
	return &bookRepositoryGORM{db: db}
}

func (r *bookRepositoryGORM) Create(ctx context.Context, book *domain.Book) error {
	gormBook := BookFromDomain(book)
	result := r.db.WithContext(ctx).Create(gormBook)
	if result.Error != nil {
		return result.Error
	}
	book.ID = gormBook.ID
	return nil
}

func (r *bookRepositoryGORM) FindAll(ctx context.Context) ([]domain.Book, error) {
	var gormBooks []BookGORM
	if err := r.db.WithContext(ctx).Order("book_order asc, title asc").Find(&gormBooks).Error; err != nil {
		return nil, err
	}

	var domainBooks []domain.Book
	for _, b := range gormBooks {
		domainBooks = append(domainBooks, *b.ToDomain())
	}
	return domainBooks, nil
}

func (r *bookRepositoryGORM) FindByID(ctx context.Context, id uint) (*domain.Book, error) {
	var gormBook BookGORM
	result := r.db.WithContext(ctx).First(&gormBook, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormBook.ToDomain(), nil
}

func (r *bookRepositoryGORM) Update(ctx context.Context, book *domain.Book) error {
	gormBook := BookFromDomain(book)
	return r.db.WithContext(ctx).Save(gormBook).Error
}

func (r *bookRepositoryGORM) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&BookGORM{}, id).Error
}

func (r *bookRepositoryGORM) FindByOrder(ctx context.Context, order int) (*domain.Book, error) {
	var gormBook BookGORM
	result := r.db.WithContext(ctx).Where("book_order = ?", order).First(&gormBook)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormBook.ToDomain(), nil
}
