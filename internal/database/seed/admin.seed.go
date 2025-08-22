package seeds

import (
	"log"

	"github.com/ipincamp/edsa/internal/config"
	"github.com/ipincamp/edsa/internal/models"
	"github.com/ipincamp/edsa/internal/utils"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB, env *config.Env) {
	hashedPassword, err := utils.HashPassword(env.AdminPassword)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	adminUser := models.User{
		Name:     env.AdminName,
		Password: hashedPassword,
	}

	result := db.Where(models.User{Email: env.AdminEmail}).
		Assign(adminUser).
		FirstOrCreate(&models.User{})

	if result.Error != nil {
		log.Fatalf("Cannot seed admin user: %v", result.Error)
	}

	if result.RowsAffected > 0 {
		log.Println("Seeded or updated admin user record.")
	} else {
		log.Println("Admin user record already up to date.")
	}
}
