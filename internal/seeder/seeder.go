package seeder

import (
	"log"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/seeder/seeders"
	"gorm.io/gorm"
)

// Seeder adalah tipe fungsi untuk menjalankan proses seeding
type Seeder func(db *gorm.DB) error

// getSeeders mengembalikan daftar seeder yang akan dijalankan
func getSeeders(cnf *config.Config) []Seeder {
	return []Seeder{
		seeders.RoleSeeder,
		seeders.UserSeeder(cnf),
	}
}

// Seed menjalankan semua seeder dalam satu transaksi
func Seed(db *gorm.DB, cnf *config.Config) {
	log.Println("Running seeders inside a transaction...")
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, seeder := range getSeeders(cnf) {
			if err := seeder(tx); err != nil {
				return err
			}
		}
		return nil
	})
	handleSeedResult(err)
}

// handleSeedResult menangani hasil akhir proses seeding
func handleSeedResult(err error) {
	if err != nil {
		log.Fatalf("Seeding failed, transaction rolled back: %v", err)
	}
	log.Println("Seeding completed successfully")
}
