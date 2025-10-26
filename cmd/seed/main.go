package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/database/seeders"
	"github.com/ipincamp/go-edsa-api/internal/pkg/database"
)

// setupSeederLogger initializes and returns a logger for the seeder.
func setupSeederLogger() (*log.Logger, *os.File) {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("❌ Failed to create log directory: %v", err)
	}

	// Create a log file name with date, e.g., seeder_2025-10-26.log
	logFileName := fmt.Sprintf("seeder_%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logDir, logFileName)

	// Open the log file in append mode, create if it doesn't exist
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("❌ Failed to open seeder log file %s: %v", logFilePath, err)
	}

	// Create a multi-writer to log to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Create the logger instance
	seederLogger := log.New(multiWriter, "SEEDER: ", log.Ldate|log.Ltime|log.LUTC)

	seederLogger.Printf("📝 Seeder logging initialized. Writing to %s", logFilePath)
	return seederLogger, logFile
}

func main() {
	// Initialize the dedicated seeder logger
	seederLogger, logFile := setupSeederLogger()
	defer logFile.Close() // Ensure the log file is closed when main exits

	// Use seederLogger for initial messages
	seederLogger.Println("🌱 Loading configuration...")
	// Load config first
	config.LoadConfig()
	cfg := config.AppConfig

	log.Println("🔧 Initializing database connection...")
	db, err := database.NewPostgresConnection(cfg) // Pass config
	if err != nil {                                // Handle error
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("🔗 Database connection established.") // Log success after check

	seederLogger.Println("⏳ Starting database seeding transaction...")

	// 1. Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		seederLogger.Fatalf("❌ Failed to start transaction: %v", tx.Error)
	}

	// 2. Defer rollback jika terjadi panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			seederLogger.Fatalf("💥 Seeding panicked, transaction rolled back. Panic: %v", r)
		}
	}()

	// 3. Jalankan semua seeder di dalam transaksi
	seederLogger.Println("🏃 Running seeders...")
	// Pass the seederLogger to RunAllSeeders
	if err := seeders.RunAllSeeders(tx, seederLogger); err != nil {
		seederLogger.Printf("❌ ERROR: Seeding failed, rolling back transaction... Error: %v", err)
		// 4. Jika ada error, rollback
		if rbErr := tx.Rollback().Error; rbErr != nil {
			seederLogger.Fatalf("❌ Failed to rollback transaction: %v", rbErr)
		}
		seederLogger.Println("⏪ Transaction rolled back successfully.")
		os.Exit(1) // Exit with error code
	}

	// 5. Jika sukses, commit
	seederLogger.Println("✅ Seeding successful, committing transaction...")
	if err := tx.Commit().Error; err != nil {
		seederLogger.Fatalf("❌ Failed to commit transaction: %v", err)
	}
	seederLogger.Println("🎉 Seeding completed and committed.")
}
