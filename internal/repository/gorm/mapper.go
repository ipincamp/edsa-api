package gorm

import "github.com/ipincamp/go-edsa-api/internal/domain"

// --- User Mappers ---

// Map GORM model ke Domain entity
func (u *UserGORM) ToDomain() *domain.User {
	domainUser := &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		RoleID:    u.RoleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}

	// Penting: mapping untuk role jika di-preload
	if u.Role.ID != 0 {
		domainUser.Role = *u.Role.ToDomain()
	}
	return domainUser
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

// --- Role Mappers ---

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

// --- Subject Mappers ---

// Map GORM model ke Domain entity
func (s *SubjectGORM) ToDomain() *domain.Subject {
	return &domain.Subject{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		DeletedAt: s.DeletedAt,
	}
}

// Map Domain entity ke GORM model
func SubjectFromDomain(s *domain.Subject) *SubjectGORM {
	return &SubjectGORM{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		DeletedAt: s.DeletedAt,
	}
}

// --- Class Mappers ---

// Map GORM model ke Domain entity
func (c *ClassGORM) ToDomain() *domain.Class {
	domainClass := &domain.Class{
		ID:        c.ID,
		Name:      c.Name,
		SubjectID: c.SubjectID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
	if c.Subject.ID != 0 {
		domainClass.Subject = *c.Subject.ToDomain()
	}
	return domainClass
}

// Map Domain entity ke GORM model
func ClassFromDomain(c *domain.Class) *ClassGORM {
	return &ClassGORM{
		ID:        c.ID,
		Name:      c.Name,
		SubjectID: c.SubjectID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
}

// --- Group Mappers ---

// Map GORM model ke Domain entity
func (g *GroupGORM) ToDomain() *domain.Group {
	domainGroup := &domain.Group{
		ID:        g.ID,
		Name:      g.Name,
		ClassID:   g.ClassID,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
		DeletedAt: g.DeletedAt,
	}
	if g.Class.ID != 0 {
		domainGroup.Class = *g.Class.ToDomain()
	}
	return domainGroup
}

// Map Domain entity ke GORM model
func GroupFromDomain(g *domain.Group) *GroupGORM {
	return &GroupGORM{
		ID:        g.ID,
		Name:      g.Name,
		ClassID:   g.ClassID,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
		DeletedAt: g.DeletedAt,
	}
}
