package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

type BookRepository interface {
	FindAllOrderedByLevel(ctx context.Context) ([]domain.Book, error)
	FindByID(ctx context.Context, id string) (*domain.Book, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) FindAllOrderedByLevel(ctx context.Context) ([]domain.Book, error) {
	var books []domain.Book
	err := r.db.WithContext(ctx).Order("level asc").Find(&books).Error
	return books, err
}

func (r *bookRepository) FindByID(ctx context.Context, id string) (*domain.Book, error) {
	var book domain.Book
	err := r.db.WithContext(ctx).Preload("Pages").Where("id = ?", id).First(&book).Error
	return &book, err
}
