package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresConnection(dsn string, env string) *gorm.DB {
	// Tentukan level log berdasarkan environment
	var logLevel logger.LogLevel
	if env == "production" {
		logLevel = logger.Silent // Jangan log SQL di production
	} else {
		logLevel = logger.Info // Log SQL di development, dll.
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel), // Gunakan logLevel dinamis
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established")
	return db
}
