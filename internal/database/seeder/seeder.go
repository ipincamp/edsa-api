package seeder

import (
	"log"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/database/seeder/seeders"
	"gorm.io/gorm"
)

type Seeder func(db *gorm.DB) error

func getSeeders(cnf *config.Config) []Seeder {
	return []Seeder{
		seeders.RoleSeeder,
		seeders.UserSeeder(cnf),
		// Other seeders can be added here
	}
}

func Seed(db *gorm.DB, cnf *config.Config) {
	err := db.Transaction(func(tx *gorm.DB) error {
		log.Println("Running seeders inside a transaction...")

		for _, seeder := range getSeeders(cnf) {
			if err := seeder(tx); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Seeding failed, transaction rolled back: %v", err)
	}
	log.Println("Seeding completed successfully")
}
