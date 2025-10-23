package gorm

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type roleRepositoryGORM struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) usecase.RoleRepository {
	return &roleRepositoryGORM{db: db}
}

func (r *roleRepositoryGORM) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	var gormRole RoleGORM
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&gormRole)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormRole.ToDomain(), nil
}

func (r *roleRepositoryGORM) FindByID(ctx context.Context, id uint) (*domain.Role, error) {
	var gormRole RoleGORM
	result := r.db.WithContext(ctx).First(&gormRole, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormRole.ToDomain(), nil
}

func (r *roleRepositoryGORM) FindAll(ctx context.Context) ([]domain.Role, error) {
	var gormRoles []RoleGORM
	if err := r.db.WithContext(ctx).Find(&gormRoles).Error; err != nil {
		return nil, err
	}

	var domainRoles []domain.Role
	for _, role := range gormRoles {
		domainRoles = append(domainRoles, *role.ToDomain())
	}
	return domainRoles, nil
}
