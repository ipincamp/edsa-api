package gorm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type userBookProgressRepositoryGORM struct {
	db *gorm.DB
}

func NewUserBookProgressRepository(db *gorm.DB) usecase.UserBookProgressRepository {
	return &userBookProgressRepositoryGORM{db: db}
}

func (r *userBookProgressRepositoryGORM) FindOrCreate(ctx context.Context, progress *domain.UserBookProgress) error {
	gormProgress := UserBookProgressFromDomain(progress)

	// Cari berdasarkan UserID dan BookID. Jika tidak ada, buat baru
	if err := r.db.WithContext(ctx).
		Where(UserBookProgressGORM{UserID: gormProgress.UserID, BookID: gormProgress.BookID}).
		FirstOrCreate(gormProgress).Error; err != nil {
		return err
	}

	// Update domain model dengan ID yang mungkin baru dibuat
	*progress = *gormProgress.ToDomain()
	return nil
}

func (r *userBookProgressRepositoryGORM) FindByUserAndBook(ctx context.Context, userID uuid.UUID, bookID uint) (*domain.UserBookProgress, error) {
	var gormProgress UserBookProgressGORM
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND book_id = ?", userID, bookID).
		First(&gormProgress)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Tidak ditemukan itu bukan error
		}
		return nil, result.Error
	}
	return gormProgress.ToDomain(), nil
}

func (r *userBookProgressRepositoryGORM) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.UserBookProgress, error) {
	var gormProgresses []UserBookProgressGORM
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&gormProgresses).Error; err != nil {
		return nil, err
	}

	var domainProgresses []domain.UserBookProgress
	for _, p := range gormProgresses {
		domainProgresses = append(domainProgresses, *p.ToDomain())
	}
	return domainProgresses, nil
}

func (r *userBookProgressRepositoryGORM) Update(ctx context.Context, progress *domain.UserBookProgress) error {
	gormProgress := UserBookProgressFromDomain(progress)
	return r.db.WithContext(ctx).Save(gormProgress).Error
}
