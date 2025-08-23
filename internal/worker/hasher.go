package worker

import (
	"log"
	"sync"

	"github.com/ipincamp/edsa/internal/models"
	"github.com/ipincamp/edsa/internal/repositories"
	"github.com/ipincamp/edsa/internal/utils"
)

type RegistrationJob struct {
	Name     string
	Email    string
	Password string
}

type RegistrationProcessor struct {
	jobQueue   chan RegistrationJob
	inProgress map[string]bool
	mu         sync.Mutex
	userRepo   repositories.UserRepository
}

var processor *RegistrationProcessor

func StartRegistrationWorkers(numWorkers int, repo repositories.UserRepository) {
	processor = &RegistrationProcessor{
		jobQueue:   make(chan RegistrationJob, 100),
		inProgress: make(map[string]bool),
		userRepo:   repo,
	}

	for i := 1; i <= numWorkers; i++ {
		go processor.worker(i)
	}
	log.Printf("Started %d registration workers.", numWorkers)
}

func (p *RegistrationProcessor) worker(id int) {
	for job := range p.jobQueue {
		log.Printf("Worker %d: Processing registration for %s", id, job.Email)

		hashedPassword, err := utils.HashPassword(job.Password)
		if err != nil {
			log.Printf("Worker %d: Failed to hash password for %s: %v", id, job.Email, err)
			p.markAsDone(job.Email)
			continue
		}

		newUser := &models.User{
			Name:     job.Name,
			Email:    job.Email,
			Password: hashedPassword,
			Status:   models.StatusActive,
		}

		if err := p.userRepo.CreateUser(newUser); err != nil {
			log.Printf("Worker %d: Failed to create user %s: %v", id, job.Email, err)
		} else {
			log.Printf("Worker %d: Successfully registered user %s", id, job.Email)
		}

		p.markAsDone(job.Email)
	}
}

func QueueRegistrationJob(job RegistrationJob) {
	processor.mu.Lock()
	processor.inProgress[job.Email] = true
	processor.mu.Unlock()

	processor.jobQueue <- job
}

func IsEmailBeingProcessed(email string) bool {
	processor.mu.Lock()
	defer processor.mu.Unlock()
	_, exists := processor.inProgress[email]
	return exists
}

func (p *RegistrationProcessor) markAsDone(email string) {
	processor.mu.Lock()
	defer processor.mu.Unlock()
	delete(p.inProgress, email)
}
