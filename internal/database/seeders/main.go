package seeders

import (
	"log"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"gorm.io/gorm"
)

// Update RunAllSeeders to accept a logger
func RunAllSeeders(db *gorm.DB, logger *log.Logger) error {
	// Pass the logger to each seeder function
	logger.Println("Seeding Roles...")
	if err := RoleSeeder(db, logger); err != nil {
		return err
	}
	logger.Println("Seeding Default Avatars...")
	if err := DefaultAvatarsSeeder(db, logger); err != nil {
		return err
	}
	logger.Println("Seeding Admin User...")
	if err := UserAdminSeeder(db, logger); err != nil {
		return err
	}
	logger.Println("Seeding Books and Covers...")
	if err := BookSeeder(db, logger); err != nil {
		return err
	}
	logger.Println("Seeding Games...")
	if err := GameSeeder(db, logger); err != nil {
		return err
	}

	env := config.AppConfig.App.Env
	if env != "production" {
		logger.Printf("Running additional seeders for '%s' environment...", env)
		logger.Println("Seeding Teacher Users...")
		if err := UserTeacherSeeder(db, logger); err != nil {
			return err
		}
		logger.Println("Seeding Student Users...")
		if err := UserStudentSeeder(db, logger); err != nil {
			return err
		}
		logger.Println("Seeding Public Users...")
		if err := UserPublicSeeder(db, logger); err != nil {
			return err
		}
		logger.Println("Seeding Classes and Groups...")
		if err := ClassGroupSeeder(db, logger); err != nil {
			return err
		}
	} else {
		logger.Println("Production environment detected. Only essential seeders were run.")
	}

	return nil
}
