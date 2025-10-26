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

// Fungsi helper getAvatarIDFromDB (tetap sama)
func getAvatarIDFromDB(db *gorm.DB, avatarFileName string) *uuid.UUID {
	var avatarAsset repo.MediaAssetGORM
	result := db.Where("file_name = ?", avatarFileName).First(&avatarAsset)
	if result.Error != nil {
		// Gunakan log standar di sini karena logger utama mungkin belum diinisialisasi saat helper dipanggil
		log.Printf("WARNING: Default avatar '%s' not found in database for seeder/factory. Error: %v", avatarFileName, result.Error)
		return nil
	}
	return &avatarAsset.ID
}

// Fungsi helper getExistingEmails (tetap sama)
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

func UserAdminSeeder(db *gorm.DB, logger *log.Logger) error {
	logger.Println("Seeding admin user...")

	var userRole repo.RoleGORM
	if err := db.Where("name = ?", "admin").First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find 'admin' role. Did you run SeedRoles? Error: %w", err)
	}

	adminCfg := config.AppConfig.Seeder
	adminEmail := adminCfg.AdminEmail
	adminPassword := adminCfg.AdminPassword
	adminName := adminCfg.AdminName

	var existing repo.UserGORM
	if db.Where("email = ?", adminEmail).First(&existing).Error == nil {
		logger.Printf("User with email %s already exists, skipping.", adminEmail)
		return nil
	}

	passSvc := argon2id.NewPasswordService()
	hashedPassword, err := passSvc.Hash(adminPassword)
	if err != nil {
		return fmt.Errorf("failed to hash admin password for seeder: %w", err)
	}

	adminAvatarIDPtr := getAvatarIDFromDB(db, "avatar4.png")

	user := &repo.UserGORM{
		Name:             adminName,
		Email:            adminEmail,
		Password:         hashedPassword,
		RoleID:           userRole.ID,
		ProfilePictureID: adminAvatarIDPtr,
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}
	logger.Printf("Seeded admin user successfully from config.")
	return nil
}

func seedUsersBatch(db *gorm.DB, logger *log.Logger, roleName, avatarFileName string, count int) error { // Tambahkan logger
	logger.Printf("Seeding %d %s users...", count, roleName)

	var userRole repo.RoleGORM
	if err := db.Where("name = ?", roleName).First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find '%s' role: %w", roleName, err)
	}

	passSvc := argon2id.NewPasswordService()
	defaultPassword := "password"
	hashedDefaultPassword, err := passSvc.Hash(defaultPassword)
	if err != nil {
		return fmt.Errorf("failed to hash default password for %s seeder: %w", roleName, err)
	}
	logger.Printf("   Hashed default password for %s users.", roleName)

	existingEmails, err := getExistingEmails(db)
	if err != nil {
		return err
	}
	logger.Printf("   Loaded %d existing emails.", len(existingEmails))

	newUsers := make([]*repo.UserGORM, 0, count)
	generatedEmails := make(map[string]bool)
	attempts := 0
	maxAttempts := count * 5

	for len(newUsers) < count && attempts < maxAttempts {
		attempts++
		user := factories.UserFactory(db, userRole.ID, avatarFileName, hashedDefaultPassword)

		if existingEmails[user.Email] || generatedEmails[user.Email] {
			// logger.Printf("   Email %s already exists, regenerating...", user.Email) // Optional
			continue
		}
		newUsers = append(newUsers, user)
		generatedEmails[user.Email] = true
	}

	if len(newUsers) < count {
		logger.Printf("   WARNING: Could only generate %d unique users after %d attempts.", len(newUsers), attempts)
	}

	if len(newUsers) > 0 {
		logger.Printf("   Inserting %d new %s users into database...", len(newUsers), roleName)
		if err := db.Create(&newUsers).Error; err != nil {
			return fmt.Errorf("failed to batch insert %s users: %w", roleName, err)
		}
		logger.Printf("   Successfully inserted %d users.", len(newUsers))
	} else {
		logger.Printf("   No new unique users generated.")
	}

	return nil
}

func UserTeacherSeeder(db *gorm.DB, logger *log.Logger) error {
	return seedUsersBatch(db, logger, "teacher", "avatar3.png", 3)
}

func UserStudentSeeder(db *gorm.DB, logger *log.Logger) error {
	return seedUsersBatch(db, logger, "student", "avatar2.png", 5)
}

func UserPublicSeeder(db *gorm.DB, logger *log.Logger) error {
	return seedUsersBatch(db, logger, "public", "avatar1.png", 2)
}
