package batcher

import (
	"log"
	"sync"
	"time"

	"github.com/ipincamp/edsa/internal/models"
	"github.com/ipincamp/edsa/internal/repositories"
	"github.com/ipincamp/edsa/internal/utils"
)

type RegistrationRequest struct {
	Name     string
	Email    string
	Password string
}

type Processor struct {
	queue      []RegistrationRequest
	userRepo   repositories.UserRepository
	emailCache repositories.EmailCache
	mu         sync.Mutex
	interval   time.Duration
}

func NewProcessor(userRepo repositories.UserRepository, emailCache repositories.EmailCache, interval time.Duration) *Processor {
	return &Processor{
		queue:      make([]RegistrationRequest, 0),
		userRepo:   userRepo,
		emailCache: emailCache,
		interval:   interval,
	}
}

func (p *Processor) AddToQueue(req RegistrationRequest) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.queue = append(p.queue, req)
}

func (p *Processor) Start() {
	log.Printf("Starting registration batch processor with a %v interval...", p.interval)
	ticker := time.NewTicker(p.interval)

	go func() {
		for range ticker.C {
			p.processBatch()
		}
	}()
}

func (p *Processor) processBatch() {
	p.mu.Lock()
	if len(p.queue) == 0 {
		p.mu.Unlock()
		return
	}

	processingQueue := make([]RegistrationRequest, len(p.queue))
	copy(processingQueue, p.queue)
	p.queue = make([]RegistrationRequest, 0)
	p.mu.Unlock()

	log.Printf("Processing batch of %d registration requests...", len(processingQueue))

	usersToCreate := make([]models.User, 0, len(processingQueue))
	emailsToUpdate := make([]string, 0, len(processingQueue))

	for _, req := range processingQueue {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			log.Printf("Failed to hash password for %s, removing from cache.", req.Email)
			p.emailCache.RemoveEmail(req.Email)
			continue
		}

		usersToCreate = append(usersToCreate, models.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: hashedPassword,
			Status:   models.StatusActive,
		})
		emailsToUpdate = append(emailsToUpdate, req.Email)
	}

	if len(usersToCreate) > 0 {
		if err := p.userRepo.CreateUsersInBatch(usersToCreate); err != nil {
			log.Printf("ERROR: Failed to save user batch to database: %v", err)
			for _, email := range emailsToUpdate {
				p.emailCache.RemoveEmail(email)
			}
		} else {
			for _, email := range emailsToUpdate {
				p.emailCache.SetEmailStatus(email, repositories.StatusRegistered)
			}
			log.Printf("Successfully saved %d new users and updated cache.", len(usersToCreate))
		}
	}
}
