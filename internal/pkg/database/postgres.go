package database

import (
	"fmt"
	"log"
	"os" // Import os
	"time"

	"github.com/ipincamp/go-edsa-api/internal/config" // Import config
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresConnection initializes GORM, configures pooling, and returns the DB instance or an error.
func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) { // Accept full config, return error
	dsn := config.GetDatabaseDSN() // Get DSN using the config method

	var logLevel logger.LogLevel
	if cfg.App.Env == "production" {
		logLevel = logger.Error // Log only errors in production for cleaner logs
	} else {
		logLevel = logger.Info // Log SQL statements in development/staging
	}

	// Configure a more detailed GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Use standard log output
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Log queries slower than 200ms
			LogLevel:                  logLevel,               // Set log level dynamically
			IgnoreRecordNotFoundError: true,                   // Don't log 'record not found' errors
			Colorful:                  false,                  // Disable color coding in logs unless preferred
		},
	)

	// Open database connection with configuration
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      newLogger, // Use the configured logger
		PrepareStmt: true,      // Enable prepared statement caching for performance
		// Consider SkipDefaultTransaction: true if you manually handle all transactions
	})
	if err != nil {
		// Return error instead of Fatalf for better control in caller (main.go)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying sql.DB instance to configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		// Return error if we can't get the underlying DB instance
		return nil, fmt.Errorf("failed to get underlying sql.DB for pool configuration: %w", err)
	}

	// Set connection pool parameters using values from config
	sqlDB.SetMaxIdleConns(cfg.Database.PoolMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.PoolMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.PoolConnMaxLifetime) * time.Minute)

	// Optional: Ping the database to verify connection immediately
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database after connection: %w", err)
	}

	return db, nil // Return the configured DB instance and nil error
}
