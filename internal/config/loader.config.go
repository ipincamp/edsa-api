package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Get() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	tokenTTLMin, err := strconv.Atoi(os.Getenv("PASETO_TOKEN_TTL_MIN"))
	if err != nil {
		log.Fatal("Error parsing PASETO_TOKEN_TTL_MIN", err.Error())
	}

	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
			Tz:   os.Getenv("DB_TZ"),
		},
		Paseto: Paseto{
			SecretKey:   os.Getenv("PASETO_SECRET_KEY"),
			TokenTTLMin: tokenTTLMin,
		},
		Seeder: Seeder{
			Admin: Admin{
				Name:     os.Getenv("ADMIN_NAME"),
				Email:    os.Getenv("ADMIN_EMAIL"),
				Password: os.Getenv("ADMIN_PASSWORD"),
			},
		},
	}
}
