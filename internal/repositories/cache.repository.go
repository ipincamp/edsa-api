package repositories

import (
	"log"
	"sync"

	"github.com/ipincamp/edsa/internal/models"
	"gorm.io/gorm"
)

const (
	StatusRegistered = "registered"
	StatusProcessing = "processing"
)

type EmailCache interface {
	LoadAllEmails() error
	GetEmailStatus(email string) (status string, exists bool)
	SetEmailStatus(email, status string)
	RemoveEmail(email string)
}

type emailCache struct {
	db    *gorm.DB
	cache map[string]string
	mu    sync.RWMutex
}

func NewEmailCache(db *gorm.DB) EmailCache {
	return &emailCache{
		db:    db,
		cache: make(map[string]string),
	}
}

func (c *emailCache) LoadAllEmails() error {
	log.Println("Loading all registered emails into in-memory cache...")
	var emails []string
	if err := c.db.Model(&models.User{}).Pluck("email", &emails).Error; err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, email := range emails {
		c.cache[email] = StatusRegistered
	}
	log.Printf("Successfully loaded %d emails into cache.", len(emails))
	return nil
}

func (c *emailCache) GetEmailStatus(email string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	status, exists := c.cache[email]
	return status, exists
}

func (c *emailCache) SetEmailStatus(email, status string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[email] = status
	log.Printf("Set status for email %s to %s in cache.", email, status)
}

func (c *emailCache) RemoveEmail(email string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, email)
	log.Printf("Removed email %s from cache.", email)
}
