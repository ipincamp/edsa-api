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

	log.Println("Running seeders...")
	seeders.RunAllSeeders(db)
	log.Println("Seeding completed.")
}
