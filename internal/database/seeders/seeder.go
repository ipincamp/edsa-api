package seeders

import (
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/internal/database/factories"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	log.Println("Seeding roles...")
	roles := []repo.RoleGORM{
		{Name: "admin"},
		{Name: "teacher"},
		{Name: "student"},
		{Name: "public"},
	}

	for _, role := range roles {
		result := db.FirstOrCreate(&role, repo.RoleGORM{Name: role.Name})
		if result.Error != nil {
			log.Printf("Failed to seed role '%s': %v\n", role.Name, result.Error)
		}
		if result.RowsAffected > 0 {
			log.Printf("Seeded role: %s\n", role.Name)
		}
	}
}

func SeedUsers(db *gorm.DB, count int) {
	log.Println("Seeding users...")

	// 1. Dapatkan ID role "student" dari database
	var userRole repo.RoleGORM
	if err := db.Where("name = ?", "student").First(&userRole).Error; err != nil {
		log.Fatalf("Failed to find 'student' role. Did you run SeedRoles? Error: %v", err)
		return
	}

	for i := 0; i < count; i++ {
		// 2. Kirim userRole.ID ke factory
		user := factories.UserFactory(userRole.ID)

		// Cek dulu agar email unik
		var existing repo.UserGORM
		if db.Where("email = ?", user.Email).First(&existing).Error == nil {
			log.Printf("User with email %s already exists, skipping.\n", user.Email)
			continue
		}

		if err := db.Create(user).Error; err != nil {
			log.Printf("Failed to seed user: %v\n", err)
		}
	}
	fmt.Printf("Seeded %d users successfully.\n", count)
}

// Fungsi utama untuk menjalankan semua seeder
func RunAllSeeders(db *gorm.DB) {
	SeedRoles(db)
	SeedUsers(db, 10)
	// Panggil seeder lain di sini
}
