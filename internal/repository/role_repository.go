package repository

import (
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByName(name string) (*domain.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) FindByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return &role, err
}
