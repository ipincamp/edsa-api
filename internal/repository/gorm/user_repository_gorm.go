package gorm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type userRepositoryGORM struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) usecase.UserRepository {
	return &userRepositoryGORM{db: db}
}

func (r *userRepositoryGORM) Create(ctx context.Context, user *domain.User) error {
	gormUser := UserFromDomain(user)
	result := r.db.WithContext(ctx).Create(gormUser)
	if result.Error != nil {
		return result.Error
	}
	user.ID = gormUser.ID
	return nil
}

func (r *userRepositoryGORM) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var gormUser UserGORM
	// Tambahkan .Preload("Role") jika Anda ingin data role ikut ter-load
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&gormUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormUser.ToDomain(), nil
}

func (r *userRepositoryGORM) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var gormUser UserGORM
	// Tambahkan .Preload("Role") jika Anda ingin data role ikut ter-load
	result := r.db.WithContext(ctx).First(&gormUser, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormUser.ToDomain(), nil
}
