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

// --- Book Mappers ---

// Map GORM model ke Domain entity
func (b *BookGORM) ToDomain() *domain.Book {
	return &domain.Book{
		ID:            b.ID,
		Title:         b.Title,
		Description:   b.Description,
		CoverImageURL: b.CoverImageURL,
		Theme:         b.Theme,
		BookOrder:     b.BookOrder,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		DeletedAt:     b.DeletedAt,
	}
}

// Map Domain entity ke GORM model
func BookFromDomain(b *domain.Book) *BookGORM {
	return &BookGORM{
		ID:            b.ID,
		Title:         b.Title,
		Description:   b.Description,
		CoverImageURL: b.CoverImageURL,
		Theme:         b.Theme,
		BookOrder:     b.BookOrder,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		DeletedAt:     b.DeletedAt,
	}
}

// --- Page Mappers ---

// Map GORM model ke Domain entity
func (p *PageGORM) ToDomain() *domain.Page {
	domainPage := &domain.Page{
		ID:              p.ID,
		BookID:          p.BookID,
		PageNumber:      p.PageNumber,
		NarrativeText:   p.NarrativeText,
		InstructionText: p.InstructionText,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
		DeletedAt:       p.DeletedAt,
	}
	if p.Book.ID != 0 {
		domainPage.Book = *p.Book.ToDomain()
	}
	return domainPage
}

// Map Domain entity ke GORM model
func PageFromDomain(p *domain.Page) *PageGORM {
	return &PageGORM{
		ID:              p.ID,
		BookID:          p.BookID,
		PageNumber:      p.PageNumber,
		NarrativeText:   p.NarrativeText,
		InstructionText: p.InstructionText,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
		DeletedAt:       p.DeletedAt,
	}
}

// --- Interaction Mappers ---

// Map GORM model ke Domain entity
func (i *InteractionGORM) ToDomain() *domain.Interaction {
	domainInteraction := &domain.Interaction{
		ID:        i.ID,
		PageID:    i.PageID,
		Type:      i.Type,
		Config:    i.Config,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		DeletedAt: i.DeletedAt,
	}
	if i.Page.ID != 0 {
		domainInteraction.Page = *i.Page.ToDomain()
	}
	return domainInteraction
}

// Map Domain entity ke GORM model
func InteractionFromDomain(i *domain.Interaction) *InteractionGORM {
	return &InteractionGORM{
		ID:        i.ID,
		PageID:    i.PageID,
		Type:      i.Type,
		Config:    i.Config,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		DeletedAt: i.DeletedAt,
	}
}

// --- UserBookProgress Mappers ---

// Map GORM model ke Domain entity
func (p *UserBookProgressGORM) ToDomain() *domain.UserBookProgress {
	return &domain.UserBookProgress{
		ID:                   p.ID,
		UserID:               p.UserID,
		BookID:               p.BookID,
		Status:               p.Status,
		HighestScore:         p.HighestScore,
		LastPageID:           p.LastPageID,
		CurrentSessionPoints: p.CurrentSessionPoints,
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
		DeletedAt:            p.DeletedAt,
	}
}

// Map Domain entity ke GORM model
func UserBookProgressFromDomain(p *domain.UserBookProgress) *UserBookProgressGORM {
	return &UserBookProgressGORM{
		ID:                   p.ID,
		UserID:               p.UserID,
		BookID:               p.BookID,
		Status:               p.Status,
		HighestScore:         p.HighestScore,
		LastPageID:           p.LastPageID,
		CurrentSessionPoints: p.CurrentSessionPoints,
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
		DeletedAt:            p.DeletedAt,
	}
}
