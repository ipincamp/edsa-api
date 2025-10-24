package cache

import (
	"context"
	"sync"
	"time"
)

// tokenBlacklistCACHE adalah implementasi in-memory dari TokenBlacklistService
type tokenBlacklistCACHE struct {
	store map[string]time.Time
	mtx   sync.RWMutex
}

// NewTokenBlacklistCACHE membuat instance cache blacklist baru.
func NewTokenBlacklistCACHE() *tokenBlacklistCACHE {
	cache := &tokenBlacklistCACHE{
		store: make(map[string]time.Time),
	}
	// Mulai goroutine pembersih
	go cache.startCleaner()
	return cache
}

// BlacklistToken menambahkan token ke daftar hitam
func (c *tokenBlacklistCACHE) BlacklistToken(ctx context.Context, tokenString string, duration time.Duration) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	expirationTime := time.Now().Add(duration)
	c.store[tokenString] = expirationTime
	return nil
}

// IsTokenBlacklisted memeriksa apakah token ada di daftar hitam dan belum kedaluwarsa
func (c *tokenBlacklistCACHE) IsTokenBlacklisted(ctx context.Context, tokenString string) (bool, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	expirationTime, exists := c.store[tokenString]
	if !exists {
		return false, nil // Tidak ada di blacklist
	}

	if time.Now().After(expirationTime) {
		return false, nil // Ada di blacklist, tapi sudah kedaluwarsa
	}

	return true, nil // Ada di blacklist dan masih aktif
}

// startCleaner adalah worker background untuk membersihkan token yang kedaluwarsa
func (c *tokenBlacklistCACHE) startCleaner() {
	// Jalankan pembersih setiap 30 menit
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.mtx.Lock()
		now := time.Now()
		for token, expirationTime := range c.store {
			if now.After(expirationTime) {
				delete(c.store, token)
			}
		}
		c.mtx.Unlock()
	}
}
