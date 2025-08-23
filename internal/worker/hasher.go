package worker

import (
	"log"

	"github.com/ipincamp/edsa/internal/utils"
)

type Job struct {
	Password   string
	ResultChan chan Result
}

type Result struct {
	HashedPassword string
	Err            error
}

var jobQueue chan Job

func StartHasherWorkers(numWorkers int) {
	jobQueue = make(chan Job, 100)

	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobQueue)
	}
	log.Printf("Started %d password hasher workers.", numWorkers)
}

func worker(_ int, jobs <-chan Job) {
	for job := range jobs {
		hashedPassword, err := utils.HashPassword(job.Password)
		job.ResultChan <- Result{
			HashedPassword: hashedPassword,
			Err:            err,
		}
	}
}

func HashPasswordAsync(password string) (string, error) {
	resultChan := make(chan Result)
	job := Job{
		Password:   password,
		ResultChan: resultChan,
	}

	jobQueue <- job

	result := <-resultChan
	return result.HashedPassword, result.Err
}
