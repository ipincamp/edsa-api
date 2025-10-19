package main

import (
	"log"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/database/migrations"
	"github.com/ipincamp/go-edsa-api/internal/pkg/database"
)

func main() {
	config.LoadConfig()
	db := database.NewPostgresConnection(config.GetDatabaseDSN())

	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetAllMigrations())

	if len(os.Args) < 2 {
		log.Fatal("Missing command. Usage: go run cmd/migrate/main.go [up|down]")
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
	default:
		log.Fatalf("Unknown command: %s. Use 'up' or 'down'.", command)
	}
}
