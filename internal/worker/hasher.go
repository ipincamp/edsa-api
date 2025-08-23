package worker

import (
	"log"

	"github.com/ipincamp/edsa/internal/models"
	"github.com/ipincamp/edsa/internal/repositories"
	"github.com/ipincamp/edsa/internal/utils"
)

type RegistrationJob struct {
	Name     string
	Email    string
	Password string
}

var jobQueue chan RegistrationJob

var userRepo repositories.UserRepository

func StartRegistrationWorkers(numWorkers int, repo repositories.UserRepository) {
	jobQueue = make(chan RegistrationJob, 100)
	userRepo = repo

	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobQueue)
	}
	log.Printf("Started %d registration workers.", numWorkers)
}

func worker(id int, jobs <-chan RegistrationJob) {
	for job := range jobs {
		log.Printf("Worker %d: Processing registration for %s", id, job.Email)

		hashedPassword, err := utils.HashPassword(job.Password)
		if err != nil {
			log.Printf("Worker %d: Failed to hash password for %s: %v", id, job.Email, err)
			continue
		}

		newUser := &models.User{
			Name:     job.Name,
			Email:    job.Email,
			Password: hashedPassword,
			Status:   models.StatusActive,
		}

		if err := userRepo.CreateUser(newUser); err != nil {
			log.Printf("Worker %d: Failed to create user %s: %v", id, job.Email, err)
		} else {
			log.Printf("Worker %d: Successfully registered user %s", id, job.Email)
		}
	}
}

func QueueRegistrationJob(job RegistrationJob) {
	jobQueue <- job
}
