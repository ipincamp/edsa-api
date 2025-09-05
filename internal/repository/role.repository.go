package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/domain"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRole(db *gorm.DB) domain.RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) List(ctx context.Context) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	return roles, err
}

func (r *roleRepository) FindByID(ctx context.Context, id string) (domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error
	return role, err
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).First(&role, "name = ?", name).Error
	return role, err
}

func (r *roleRepository) Update(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}
