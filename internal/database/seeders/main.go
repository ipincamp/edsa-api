package seeders

import (
	"log"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"gorm.io/gorm"
)

// Fungsi untuk menjalankan semua seeders
func RunAllSeeders(db *gorm.DB) {
	// 1. Jalankan seeder penting yang harus ada di semua environment
	RoleSeeder(db)
	UserAdminSeeder(db)
	BookSeeder(db)
	GameSeeder(db)

	// 2. Cek environment dari config
	env := config.AppConfig.App.Env

	// 3. Hanya jalankan seeder tambahan jika BUKAN production
	if env != "production" {
		log.Printf("Running additional seeders for '%s' environment...", env)
		UserTeacherSeeder(db)
		UserStudentSeeder(db)
		UserPublicSeeder(db)
		// Panggil seeder lain di sini
	} else {
		log.Println("Production environment detected. Only Role and Admin seeders were run.")
	}
}
