package seeders

import (
	"log"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"gorm.io/gorm"
)

// Fungsi untuk menjalankan semua seeders
func RunAllSeeders(db *gorm.DB) error {
	// 1. Jalankan seeder penting yang harus ada di semua environment
	if err := RoleSeeder(db); err != nil {
		return err
	}
	if err := UserAdminSeeder(db); err != nil {
		return err
	}
	if err := BookSeeder(db); err != nil {
		return err
	}
	if err := GameSeeder(db); err != nil {
		return err
	}

	// 2. Cek environment dari config
	env := config.AppConfig.App.Env

	// 3. Hanya jalankan seeder tambahan jika BUKAN production
	if env != "production" {
		log.Printf("Running additional seeders for '%s' environment...", env)
		if err := UserTeacherSeeder(db); err != nil {
			return err
		}
		if err := UserStudentSeeder(db); err != nil {
			return err
		}
		if err := UserPublicSeeder(db); err != nil {
			return err
		}
		// Panggil seeder lain di sini
	} else {
		log.Println("Production environment detected. Only Role and Admin seeders were run.")
	}

	return nil
}
