package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

// RoleRepository adalah kontrak untuk akses data role
type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*domain.Role, error)
}

// roleRepository adalah implementasi RoleRepository menggunakan GORM
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository membuat instance baru roleRepository
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// FindByName mencari role berdasarkan nama
func (r *roleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).Take(&role).Error
	return &role, err
}
