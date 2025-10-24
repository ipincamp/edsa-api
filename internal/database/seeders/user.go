package seeders

import (
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/database/factories"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
	"gorm.io/gorm"
)

func UserAdminSeeder(db *gorm.DB) error {
	log.Println("Seeding admin user...")

	// 1. Dapatkan ID role "admin" dari database
	var userRole repo.RoleGORM
	if err := db.Where("name = ?", "admin").First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find 'admin' role. Did you run SeedRoles? Error: %w", err)
	}

	// 2. Ambil kredensial admin dari config (menggunakan global AppConfig)
	adminCfg := config.AppConfig.Seeder
	adminEmail := adminCfg.AdminEmail
	adminPassword := adminCfg.AdminPassword
	adminName := adminCfg.AdminName

	// Ambil URL avatar default
	baseURL := config.AppConfig.Storage.StoragePublicBaseURL
	defaultAvatarURL := baseURL + "/public/uploads/avatar.png"

	// 3. Cek dulu agar email unik
	var existing repo.UserGORM
	if db.Where("email = ?", adminEmail).First(&existing).Error == nil {
		log.Printf("User with email %s already exists, skipping.\n", adminEmail)
		return nil
	}

	// 4. Hash password admin
	passSvc := argon2id.NewPasswordService()
	hashedPassword, err := passSvc.Hash(adminPassword)
	if err != nil {
		return fmt.Errorf("failed to hash admin password for seeder: %w", err)
	}

	// 5. Buat user admin baru dari config (tidak pakai factory)
	user := &repo.UserGORM{
		Name:              adminName,
		Email:             adminEmail,
		Password:          hashedPassword,
		RoleID:            userRole.ID,
		ProfilePictureURL: defaultAvatarURL,
	}

	// 6. Simpan ke database
	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}
	fmt.Println("Seeded admin user successfully from config.")
	return nil
}

func UserTeacherSeeder(db *gorm.DB) error {
	log.Println("Seeding teachers...")

	// 1. Dapatkan ID role "teacher" dari database
	var userRole repo.RoleGORM
	if err := db.Where("name = ?", "teacher").First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find 'teacher' role. Did you run SeedRoles? Error: %w", err)
	}

	count := 3
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
			return fmt.Errorf("failed to seed teacher user: %w", err)
		}
	}
	fmt.Printf("Seeded %d users successfully.\n", count)
	return nil
}

func UserStudentSeeder(db *gorm.DB) error {
	log.Println("Seeding students...")

	// 1. Dapatkan ID role "student" dari database
	var userRole repo.RoleGORM
	if err := db.Where("name = ?", "student").First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find 'student' role. Did you run SeedRoles? Error: %w", err)
	}

	count := 5
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
			return fmt.Errorf("failed to seed student user: %w", err)
		}
	}
	fmt.Printf("Seeded %d users successfully.\n", count)
	return nil
}

func UserPublicSeeder(db *gorm.DB) error {
	log.Println("Seeding public users...")

	// 1. Dapatkan ID role "public" dari database
	var userRole repo.RoleGORM
	if err := db.Where("name = ?", "public").First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find 'public' role. Did you run SeedRoles? Error: %w", err)
	}

	count := 2
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
			return fmt.Errorf("failed to seed public user: %w", err)
		}
	}
	fmt.Printf("Seeded %d users successfully.\n", count)
	return nil
}
