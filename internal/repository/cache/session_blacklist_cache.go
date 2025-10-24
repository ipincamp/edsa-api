package cache

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

// sessionBlacklistCACHE adalah implementasi in-memory dari SessionBlacklistService
type sessionBlacklistCACHE struct {
	store map[uuid.UUID]time.Time
	mtx   sync.RWMutex
}

// NewSessionBlacklistCACHE membuat instance cache blacklist baru.
// Kita ubah return type nya agar sesuai dengan interface
func NewSessionBlacklistCACHE() usecase.SessionBlacklistService {
	cache := &sessionBlacklistCACHE{
		store: make(map[uuid.UUID]time.Time),
	}
	go cache.startCleaner()
	return cache
}

// BlacklistSession menambahkan sessionID ke daftar hitam
func (c *sessionBlacklistCACHE) BlacklistSession(ctx context.Context, sessionID uuid.UUID, duration time.Duration) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	expirationTime := time.Now().Add(duration)
	c.store[sessionID] = expirationTime
	return nil
}

// IsSessionBlacklisted memeriksa apakah sessionID ada di daftar hitam
func (c *sessionBlacklistCACHE) IsSessionBlacklisted(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	expirationTime, exists := c.store[sessionID]
	if !exists {
		return false, nil // Tidak ada di blacklist
	}

	if time.Now().After(expirationTime) {
		return false, nil // Ada di blacklist, tapi sudah kedaluwarsa
	}

	return true, nil // Ada di blacklist dan masih aktif
}

// startCleaner (tetap sama, tapi membersihkan map[uuid.UUID])
func (c *sessionBlacklistCACHE) startCleaner() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.mtx.Lock()
		now := time.Now()
		for sessionID, expirationTime := range c.store {
			if now.After(expirationTime) {
				delete(c.store, sessionID)
			}
		}
		c.mtx.Unlock()
	}
}
