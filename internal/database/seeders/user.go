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

// Fungsi helper getExistingEmails (baru)
func getExistingEmails(db *gorm.DB) (map[string]bool, error) {
	existingEmails := make(map[string]bool)
	var emails []string
	if err := db.Model(&repo.UserGORM{}).Pluck("email", &emails).Error; err != nil {
		return nil, fmt.Errorf("failed to load existing emails: %w", err)
	}
	for _, email := range emails {
		existingEmails[email] = true
	}
	return existingEmails, nil
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

// seedUsersBatch sekarang hash password sekali saja
func seedUsersBatch(db *gorm.DB, roleName, avatarFileName string, count int) error {
	log.Printf("Seeding %d %s users...", count, roleName)

	// 1. Dapatkan Role ID
	var userRole repo.RoleGORM
	if err := db.Where("name = ?", roleName).First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find '%s' role: %w", roleName, err)
	}

	// 2. Hash password default SEKALI di sini
	passSvc := argon2id.NewPasswordService()
	defaultPassword := "password" // Password default untuk user dummy
	hashedDefaultPassword, err := passSvc.Hash(defaultPassword)
	if err != nil {
		return fmt.Errorf("failed to hash default password for %s seeder: %w", roleName, err)
	}
	log.Printf("   Hashed default password for %s users.", roleName)

	// 3. Ambil email yang sudah ada
	existingEmails, err := getExistingEmails(db)
	if err != nil {
		return err
	}
	log.Printf("   Loaded %d existing emails.", len(existingEmails))

	// 4. Buat user baru dalam batch
	newUsers := make([]*repo.UserGORM, 0, count)
	generatedEmails := make(map[string]bool)
	attempts := 0
	maxAttempts := count * 5

	for len(newUsers) < count && attempts < maxAttempts {
		attempts++
		// Panggil factory dengan hash yang sudah dibuat
		user := factories.UserFactory(db, userRole.ID, avatarFileName, hashedDefaultPassword)

		// 5. Cek duplikasi email
		if existingEmails[user.Email] || generatedEmails[user.Email] {
			// log.Printf("   Email %s already exists, regenerating...", user.Email) // Optional logging
			continue
		}
		newUsers = append(newUsers, user)
		generatedEmails[user.Email] = true
	}

	if len(newUsers) < count {
		log.Printf("   WARNING: Could only generate %d unique users after %d attempts.", len(newUsers), attempts)
	}

	// 6. Batch Insert
	if len(newUsers) > 0 {
		log.Printf("   Inserting %d new %s users into database...", len(newUsers), roleName)
		if err := db.Create(&newUsers).Error; err != nil {
			return fmt.Errorf("failed to batch insert %s users: %w", roleName, err)
		}
		log.Printf("   Successfully inserted %d users.", len(newUsers))
	} else {
		log.Printf("   No new unique users generated.")
	}

	return nil
}

func UserTeacherSeeder(db *gorm.DB) error {
	return seedUsersBatch(db, "teacher", "avatar3.png", 3)
}

func UserStudentSeeder(db *gorm.DB) error {
	return seedUsersBatch(db, "student", "avatar2.png", 5)
}

func UserPublicSeeder(db *gorm.DB) error {
	return seedUsersBatch(db, "public", "avatar1.png", 2)
}
