package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/database/migrations"
	"github.com/ipincamp/go-edsa-api/internal/pkg/database"
	"gorm.io/gorm"
)

// dropAllTables menghapus semua tabel yang dikenal (kecuali gormigrate)
func dropAllTables(db *gorm.DB) error {
	// List tabel berdasarkan file models.go dan migrations
	// Urutkan dari tabel yang memiliki foreign key ke tabel yang direferensikan
	tables := []string{
		"group_book_settings",
		"media_assets",
		"activity_logs",
		"user_game_scores",
		"user_book_progresses",
		"interactions",
		"user_groups",
		"pages",
		"groups",
		"classes",
		"users",
		"games",
		"books",
		"subjects",
		"roles",
		// tabel gormigrate
		"migrations",
	}

	for _, table := range tables {
		// Kita cek dulu apakah tabelnya ada sebelum drop
		if db.Migrator().HasTable(table) {
			// Gunakan "CASCADE" untuk otomatis drop constraint yang bergantung
			if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
				log.Printf("Could not drop table %s: %v", table, err)
			} else {
				log.Printf("Dropped table: %s", table)
			}
		} else {
			log.Printf("Table %s does not exist, skipping.", table)
		}
	}

	return nil
}

func main() {
	config.LoadConfig()
	db := database.NewPostgresConnection(config.GetDatabaseDSN())

	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetAllMigrations())

	if len(os.Args) < 2 {
		log.Fatal("Missing command. Usage: go run cmd/migrate/main.go [up|down|reset|drop-all]")
	}

	command := os.Args[1]
	switch command {
	case "up":
		log.Println("Running migrations...")
		if err := m.Migrate(); err != nil {
			log.Fatalf("Could not migrate: %v", err)
		}
		log.Println("Migrations ran successfully")
	case "down":
		log.Println("Rolling back last migration...")
		if err := m.RollbackLast(); err != nil {
			log.Fatalf("Could not rollback: %v", err)
		}
		log.Println("Rollback successful")
	case "reset":
		log.Println("Resetting database (dropping all tables)...")
		if err := dropAllTables(db); err != nil {
			log.Fatalf("Could not drop all tables: %v", err)
		}
		log.Println("All tables dropped.")

		log.Println("Running all migrations...")
		if err := m.Migrate(); err != nil {
			log.Fatalf("Could not migrate: %v", err)
		}
		log.Println("Migrations ran successfully")
		log.Println("Database reset complete.")

	case "drop-all":
		log.Println("DANGER: Dropping all known tables (including migrations table)...")
		if err := dropAllTables(db); err != nil {
			log.Fatalf("Could not drop all tables: %v", err)
		}
		log.Println("All known tables dropped successfully. Database is now empty of these tables.")

	default:
		log.Fatalf("Unknown command: %s. Use 'up', 'down', 'reset', or 'drop-all'.", command)
	}
}
