package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AppEnv        string
	AppPort       string
	AppTZ         string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	AdminName     string
	AdminEmail    string
	AdminPassword string
}

func LoadEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	return &Env{
		AppEnv:        os.Getenv("APP_ENV"),
		AppPort:       os.Getenv("APP_PORT"),
		AppTZ:         os.Getenv("TZ"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPass:        os.Getenv("DB_PASS"),
		DBName:        os.Getenv("DB_NAME"),
		AdminName:     os.Getenv("USER_ADMIN_NAME"),
		AdminEmail:    os.Getenv("USER_ADMIN_EMAIL"),
		AdminPassword: os.Getenv("USER_ADMIN_PASSWORD"),
	}
}
