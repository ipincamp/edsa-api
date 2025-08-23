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
	jobQueue            chan RegistrationJob
	processingEmails    map[string]bool
	failedRegistrations map[string]string
	mu                  sync.Mutex
	userRepo            repositories.UserRepository
	emailCache          repositories.EmailCache
}

var processor *RegistrationProcessor

func StartRegistrationWorkers(numWorkers int, repo repositories.UserRepository, cache repositories.EmailCache) {
	processor = &RegistrationProcessor{
		jobQueue:            make(chan RegistrationJob, 100),
		processingEmails:    make(map[string]bool),
		failedRegistrations: make(map[string]string),
		userRepo:            repo,
		emailCache:          cache,
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
			log.Printf("Worker %d: Hashing failed for %s: %v", id, job.Email, err)
			p.markAsFailed(job.Email, "Internal server error during processing.")
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
			p.markAsFailed(job.Email, "Could not save user data.")
		} else {
			p.emailCache.AddEmail(job.Email)
			log.Printf("Worker %d: Successfully registered user %s", id, job.Email)
			p.markAsDone(job.Email)
		}
	}
}

func QueueRegistrationJob(job RegistrationJob) {
	processor.mu.Lock()
	defer processor.mu.Unlock()
	processor.processingEmails[job.Email] = true
	delete(processor.failedRegistrations, job.Email)
	processor.jobQueue <- job
}

func IsEmailBeingProcessed(email string) bool {
	processor.mu.Lock()
	defer processor.mu.Unlock()
	_, exists := processor.processingEmails[email]
	return exists
}

func GetRegistrationFailureReason(email string) (string, bool) {
	processor.mu.Lock()
	defer processor.mu.Unlock()
	reason, exists := processor.failedRegistrations[email]
	return reason, exists
}

func (p *RegistrationProcessor) markAsDone(email string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.processingEmails, email)
}

func (p *RegistrationProcessor) markAsFailed(email, reason string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.processingEmails, email)
	p.failedRegistrations[email] = reason
}
