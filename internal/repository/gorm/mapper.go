package gorm

import "github.com/ipincamp/go-edsa-api/internal/domain"

// Map GORM model ke Domain entity
func (u *UserGORM) ToDomain() *domain.User {
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		RoleID:    u.RoleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}

// Map Domain entity ke GORM model
func UserFromDomain(u *domain.User) *UserGORM {
	return &UserGORM{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		RoleID:    u.RoleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}

// Map GORM model ke Domain entity
func (r *RoleGORM) ToDomain() *domain.Role {
	return &domain.Role{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}
}

// Map Domain entity ke GORM model
func RoleFromDomain(r *domain.Role) *RoleGORM {
	return &RoleGORM{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}
}
