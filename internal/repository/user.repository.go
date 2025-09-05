package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) LoadAll(ctx context.Context) (result []domain.User, err error) {
	err = r.db.WithContext(ctx).Find(&result).Error
	return
}

func (r *userRepository) FindAll(ctx context.Context, limit int, offset int) (result []domain.User, err error) {
	err = r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&result).Error
	return
}

func (r *userRepository) Count(ctx context.Context) (total int64, err error) {
	err = r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error
	return
}

func (r *userRepository) FindByID(ctx context.Context, id string) (result domain.User, err error) {
	err = r.db.WithContext(ctx).
		Preload("Role").
		First(&result, "id = ?", id).Error
	return
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (result domain.User, err error) {
	err = r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&result).Error
	return
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Model(user).Updates(user).Error
	})
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&domain.User{}, "id = ?", id).Error
	})
}
