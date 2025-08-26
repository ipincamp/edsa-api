package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Env struct {
	AppEnv              string
	AppPort             string
	AppTZ               string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPass              string
	DBName              string
	AdminName           string
	AdminEmail          string
	AdminPassword       string
	PasetoSymmetricKey  string
	PasetoExpireInHours time.Duration
	MailHost            string
	MailPort            int
	MailUsername        string
	MailPassword        string
	MailFrom            string
}

func LoadEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	hours, err := strconv.Atoi(os.Getenv("PASETO_EXPIRE_IN_HOURS"))
	if err != nil {
		hours = 24
	}
	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	return &Env{
		AppEnv:              os.Getenv("APP_ENV"),
		AppPort:             os.Getenv("APP_PORT"),
		AppTZ:               os.Getenv("TZ"),
		DBHost:              os.Getenv("DB_HOST"),
		DBPort:              os.Getenv("DB_PORT"),
		DBUser:              os.Getenv("DB_USER"),
		DBPass:              os.Getenv("DB_PASS"),
		DBName:              os.Getenv("DB_NAME"),
		AdminName:           os.Getenv("USER_ADMIN_NAME"),
		AdminEmail:          os.Getenv("USER_ADMIN_EMAIL"),
		AdminPassword:       os.Getenv("USER_ADMIN_PASSWORD"),
		PasetoSymmetricKey:  os.Getenv("PASETO_SYMMETRIC_KEY"),
		PasetoExpireInHours: time.Duration(hours) * time.Hour,
		MailHost:            os.Getenv("MAIL_HOST"),
		MailPort:            mailPort,
		MailUsername:        os.Getenv("MAIL_USERNAME"),
		MailPassword:        os.Getenv("MAIL_PASSWORD"),
		MailFrom:            os.Getenv("MAIL_FROM"),
	}
}
