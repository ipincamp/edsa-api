package main

import (
	"log"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/migrations"
	"github.com/ipincamp/go-edsa-api/pkg/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.CreateUsersTable(),
		migrations.CreateRolesTable(),
		migrations.AddRoleIdToUsersTable(),
	})

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "up":
			if err = m.Migrate(); err != nil {
				log.Fatalf("Could not migrate: %v", err)
			}
			log.Printf("Migration run successfully")
			return
		case "down":
			if err = m.RollbackLast(); err != nil {
				log.Fatalf("Could not rollback: %v", err)
			}
			log.Printf("Rollback run successfully")
			return
		default:
			log.Printf("Usage: go run cmd/migrate/main.go [up|down]")
			return
		}
	}

	// Default action is to migrate up
	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration run successfully")
}
