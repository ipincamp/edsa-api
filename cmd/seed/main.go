package main

import (
	"log"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/database/seeders"
	"github.com/ipincamp/go-edsa-api/internal/pkg/database"
)

func main() {
	config.LoadConfig()
	db := database.NewPostgresConnection(config.GetDatabaseDSN())

	log.Println("Starting database seeding transaction...")

	// 1. Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Failed to start transaction: %v", tx.Error)
	}

	// 2. Defer rollback jika terjadi panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("Seeding panicked, transaction rolled back. Panic: %v", r)
		}
	}()

	// 3. Jalankan semua seeder di dalam transaksi
	log.Println("Running seeders...")
	if err := seeders.RunAllSeeders(tx); err != nil {
		// 4. Jika ada error, rollback
		log.Printf("ERROR: Seeding failed, rolling back transaction... Error: %v", err)
		if rbErr := tx.Rollback().Error; rbErr != nil {
			log.Fatalf("Failed to rollback transaction: %v", rbErr)
		}
		log.Println("Transaction rolled back successfully.")
		return // Keluar setelah rollback
	}

	// 5. Jika sukses, commit
	log.Println("Seeding successful, committing transaction...")
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
	log.Println("Seeding completed and committed.")
}
