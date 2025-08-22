package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AppPort string
}

func LoadEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	return &Env{
		AppPort: os.Getenv("APP_PORT"),
	}
}
