package repositories

import (
	"log"
	"sync"

	"github.com/ipincamp/edsa/internal/models"
	"gorm.io/gorm"
)

type EmailCache interface {
	LoadAllEmails() error
	EmailExists(email string) bool
	AddEmail(email string)
}

type emailCache struct {
	db    *gorm.DB
	cache map[string]bool
	mu    sync.Mutex
}

func NewEmailCache(db *gorm.DB) EmailCache {
	return &emailCache{
		db:    db,
		cache: make(map[string]bool),
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
		c.cache[email] = true
	}
	log.Printf("Successfully loaded %d emails into cache.", len(emails))
	return nil
}

func (c *emailCache) EmailExists(email string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exists := c.cache[email]
	return exists
}

func (c *emailCache) AddEmail(email string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[email] = true
}
