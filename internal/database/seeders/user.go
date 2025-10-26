package seeders

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/database/factories"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
	"gorm.io/gorm"
)

// Fungsi helper getAvatarIDFromDB (sama seperti di factory)
func getAvatarIDFromDB(db *gorm.DB, avatarFileName string) *uuid.UUID {
	var avatarAsset repo.MediaAssetGORM
	result := db.Where("file_name = ?", avatarFileName).First(&avatarAsset)
	if result.Error != nil {
		log.Printf("WARNING: Default avatar '%s' not found in database for seeder. Error: %v", avatarFileName, result.Error)
		return nil
	}
	return &avatarAsset.ID
}

func UserAdminSeeder(db *gorm.DB) error {
	log.Println("Seeding admin user...")

	// 1. Dapatkan ID role "admin"
	var userRole repo.RoleGORM
	if err := db.Where("name = ?", "admin").First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find 'admin' role. Did you run SeedRoles? Error: %w", err)
	}

	// 2. Ambil kredensial admin
	adminCfg := config.AppConfig.Seeder
	adminEmail := adminCfg.AdminEmail
	adminPassword := adminCfg.AdminPassword
	adminName := adminCfg.AdminName

	// 3. Cek email unik
	var existing repo.UserGORM
	if db.Where("email = ?", adminEmail).First(&existing).Error == nil {
		log.Printf("User with email %s already exists, skipping.\n", adminEmail)
		return nil
	}

	// 4. Hash password
	passSvc := argon2id.NewPasswordService()
	hashedPassword, err := passSvc.Hash(adminPassword)
	if err != nil {
		return fmt.Errorf("failed to hash admin password for seeder: %w", err)
	}

	// 5. Cari ID avatar admin ("avatar4.png") langsung dari DB
	adminAvatarIDPtr := getAvatarIDFromDB(db, "avatar4.png")

	// 6. Buat user admin
	user := &repo.UserGORM{
		Name:             adminName,
		Email:            adminEmail,
		Password:         hashedPassword,
		RoleID:           userRole.ID,
		ProfilePictureID: adminAvatarIDPtr, // Gunakan ID dari DB
	}

	// 7. Simpan ke database
	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}
	fmt.Println("Seeded admin user successfully from config.")
	return nil
}

func UserTeacherSeeder(db *gorm.DB) error {
	log.Println("Seeding teachers...")

	var userRole repo.RoleGORM
	if err := db.Where("name = ?", "teacher").First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find 'teacher' role. Did you run SeedRoles? Error: %w", err)
	}

	count := 3
	for i := 0; i < count; i++ {
		// Panggil factory dengan nama file avatar yang diinginkan
		user := factories.UserFactory(db, userRole.ID, "avatar3.png")

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

	var userRole repo.RoleGORM
	if err := db.Where("name = ?", "student").First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find 'student' role. Did you run SeedRoles? Error: %w", err)
	}

	count := 5
	for i := 0; i < count; i++ {
		// Panggil factory dengan nama file avatar yang diinginkan
		user := factories.UserFactory(db, userRole.ID, "avatar2.png")

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

	var userRole repo.RoleGORM
	if err := db.Where("name = ?", "public").First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find 'public' role. Did you run SeedRoles? Error: %w", err)
	}

	count := 2
	for i := 0; i < count; i++ {
		// Panggil factory dengan nama file avatar yang diinginkan
		user := factories.UserFactory(db, userRole.ID, "avatar1.png")

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
