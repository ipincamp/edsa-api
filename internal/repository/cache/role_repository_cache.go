package cache

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

// roleRepositoryCACHE is the cached implementation of RoleRepository
type roleRepositoryCACHE struct {
	byID   map[uint]*domain.Role
	byName map[string]*domain.Role
	mtx    sync.RWMutex
}

// NewRoleRepositoryCACHE creates a cache repo instance and immediately loads all roles
func NewRoleRepositoryCACHE(dbRepo usecase.RoleRepository) (usecase.RoleRepository, error) {
	cacheRepo := &roleRepositoryCACHE{
		byID:   make(map[uint]*domain.Role),
		byName: make(map[string]*domain.Role),
	}

	if err := cacheRepo.load(context.Background(), dbRepo); err != nil {
		// Keep this error log as it's critical during startup
		return nil, fmt.Errorf("🚨 failed to pre-load role cache: %w", err) // <-- UPDATED ERROR FORMAT
	}

	return cacheRepo, nil
}

// load fetches all roles from the DB and stores them in the maps
func (r *roleRepositoryCACHE) load(ctx context.Context, dbRepo usecase.RoleRepository) error {
	log.Println("💾 Pre-loading role cache...") // <-- UPDATED LOG

	roles, err := dbRepo.FindAll(ctx)
	if err != nil {
		// Return the error to be handled by NewRoleRepositoryCACHE
		return err
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.byID = make(map[uint]*domain.Role)     // Clear existing cache before loading
	r.byName = make(map[string]*domain.Role) // Clear existing cache before loading

	for i := range roles {
		role := roles[i]
		r.byID[role.ID] = &role
		r.byName[role.Name] = &role
	}

	log.Printf("✅ Successfully loaded %d roles into cache.", len(roles)) // <-- UPDATED LOG
	return nil
}

// FindByName retrieves a role from the cache
func (r *roleRepositoryCACHE) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	role, ok := r.byName[name]
	if !ok {
		return nil, nil // Not found
	}
	return role, nil
}

// FindByID retrieves a role from the cache
func (r *roleRepositoryCACHE) FindByID(ctx context.Context, id uint) (*domain.Role, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	role, ok := r.byID[id]
	if !ok {
		return nil, nil // Not found
	}
	return role, nil
}

// FindAll returns all roles from the cache
func (r *roleRepositoryCACHE) FindAll(ctx context.Context) ([]domain.Role, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	roles := make([]domain.Role, 0, len(r.byID))
	for _, role := range r.byID {
		roles = append(roles, *role)
	}
	return roles, nil
}
