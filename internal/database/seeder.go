package database

import (
	"log"

	"github.com/ipincamp/edsa/internal/config"
	seeds "github.com/ipincamp/edsa/internal/database/seed"
	"gorm.io/gorm"
)

func RunSeeder(db *gorm.DB, env *config.Env) {
	log.Println("Running database seeder...")

	seeds.SeedAdmin(db, env)

	log.Println("Database seeding complete.")
}
