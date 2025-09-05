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

func (r *userRepository) List(ctx context.Context, limit, offset int, roleName string) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&domain.User{})

		if roleName != "" {
			query = query.Joins("JOIN roles ON roles.id = users.role_id").Where("roles.name = ?", roleName)
		}

		if err := query.Count(&total).Error; err != nil {
			return err
		}

		paginatedQuery := query
		if limit > 0 {
			paginatedQuery = paginatedQuery.Limit(limit).Offset(offset)
		}

		if err := paginatedQuery.Preload("Role").Find(&users).Error; err != nil {
			return err
		}
		return nil
	})

	return users, total, err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("Role").First(&user, "id = ?", id).Error
	return user, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}
