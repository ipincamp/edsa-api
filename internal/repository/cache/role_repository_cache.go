package cache

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

// roleRepositoryCACHE adalah implementasi RoleRepository yang di-cache
type roleRepositoryCACHE struct {
	byID   map[uint]*domain.Role
	byName map[string]*domain.Role
	mtx    sync.RWMutex
}

// NewRoleRepositoryCACHE membuat instance cache repo dan langsung memuat semua role
func NewRoleRepositoryCACHE(dbRepo usecase.RoleRepository) (usecase.RoleRepository, error) {
	cacheRepo := &roleRepositoryCACHE{
		byID:   make(map[uint]*domain.Role),
		byName: make(map[string]*domain.Role),
	}

	if err := cacheRepo.load(context.Background(), dbRepo); err != nil {
		return nil, fmt.Errorf("failed to pre-load role cache: %w", err)
	}

	return cacheRepo, nil
}

// load mengambil semua role dari DB dan menyimpannya di map
func (r *roleRepositoryCACHE) load(ctx context.Context, dbRepo usecase.RoleRepository) error {
	log.Println("Pre-loading role cache...")

	roles, err := dbRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	for i := range roles {
		role := roles[i] // Ambil salinan/pointer
		r.byID[role.ID] = &role
		r.byName[role.Name] = &role
	}

	log.Printf("Successfully loaded %d roles into cache.", len(roles))
	return nil
}

// FindByName mengambil role dari cache
func (r *roleRepositoryCACHE) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	role, ok := r.byName[name]
	if !ok {
		return nil, nil // Tidak ditemukan
	}
	return role, nil
}

// FindByID mengambil role dari cache
func (r *roleRepositoryCACHE) FindByID(ctx context.Context, id uint) (*domain.Role, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	role, ok := r.byID[id]
	if !ok {
		return nil, nil // Tidak ditemukan
	}
	return role, nil
}

// FindAll mengembalikan semua role dari cache
func (r *roleRepositoryCACHE) FindAll(ctx context.Context) ([]domain.Role, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	roles := make([]domain.Role, 0, len(r.byID))
	for _, role := range r.byID {
		roles = append(roles, *role)
	}
	return roles, nil
}
