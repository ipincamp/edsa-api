package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

// UserProgressRepository adalah kontrak untuk akses data progres user
type UserProgressRepository interface {
	GetProgressByUser(ctx context.Context, userID string) ([]domain.UserProgress, error)
	FindOrCreate(ctx context.Context, progress *domain.UserProgress) error
	FindByUserAndBook(ctx context.Context, userID, bookID string) (*domain.UserProgress, error)
	Update(ctx context.Context, progress *domain.UserProgress) error
}

type userProgressRepository struct {
	db *gorm.DB
}

// NewUserProgressRepository membuat instance baru
func NewUserProgressRepository(db *gorm.DB) UserProgressRepository {
	return &userProgressRepository{db: db}
}

// GetProgressByUser mengambil semua data progres milik seorang user
func (r *userProgressRepository) GetProgressByUser(ctx context.Context, userID string) ([]domain.UserProgress, error) {
	var progressList []domain.UserProgress
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&progressList).Error
	return progressList, err
}

// FindOrCreate mencari progres berdasarkan userID dan bookID. Jika tidak ada, maka akan dibuatkan yang baru.
// Ini sangat berguna saat user pertama kali membuka sebuah buku.
func (r *userProgressRepository) FindOrCreate(ctx context.Context, progress *domain.UserProgress) error {
	// GORM akan mencari record yang cocok dengan kondisi di Where().
	// Jika tidak ditemukan, GORM akan membuat record baru menggunakan seluruh data dari objek 'progress'.
	return r.db.WithContext(ctx).
		Where(domain.UserProgress{UserID: progress.UserID, BookID: progress.BookID}).
		FirstOrCreate(progress).Error
}

// FindByUserAndBook mencari progres berdasarkan userID dan bookID.
func (r *userProgressRepository) FindByUserAndBook(ctx context.Context, userID, bookID string) (*domain.UserProgress, error) {
	var progress domain.UserProgress
	err := r.db.WithContext(ctx).Where("user_id = ? AND book_id = ?", userID, bookID).First(&progress).Error
	return &progress, err
}

// Update memperbarui data progres user di database.
func (r *userProgressRepository) Update(ctx context.Context, progress *domain.UserProgress) error {
	return r.db.WithContext(ctx).Save(progress).Error
}
