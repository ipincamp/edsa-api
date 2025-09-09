package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

type PageRepository interface {
	IsLastPage(ctx context.Context, pageID, bookID string) (bool, error)
}

type pageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &pageRepository{db: db}
}

// IsLastPage mengecek apakah sebuah halaman adalah halaman terakhir dari sebuah buku
func (r *pageRepository) IsLastPage(ctx context.Context, pageID, bookID string) (bool, error) {
	var lastPage domain.Page
	err := r.db.WithContext(ctx).
		Where("book_id = ?", bookID).
		Order("page_number desc").
		First(&lastPage).Error

	if err != nil {
		return false, err
	}
	return lastPage.ID == pageID, nil
}
