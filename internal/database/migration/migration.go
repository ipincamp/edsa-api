package migration

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/ipincamp/go-edsa-api/internal/database/migration/migrations"
	"gorm.io/gorm"
)

func getMigrations(db *gorm.DB) []*gormigrate.Migration {
	return []*gormigrate.Migration{
		migrations.CreateUsersTable(),
		migrations.CreateRolesAndPermissionsTable(),
		migrations.CreateUserPermissionsTable(),
		migrations.AddRoleIdToUsersTable(),
		// Other migrations can be added here
	}
}

func Migrate(db *gorm.DB) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations(db))

	log.Println("Running pending migrations...")
	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	log.Println("Migrations ran successfully")
}

func Rollback(db *gorm.DB) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations(db))

	log.Println("Rolling back the last migration...")
	if err := m.RollbackLast(); err != nil {
		log.Fatalf("Could not rollback: %v", err)
	}

	log.Println("Rolled back successfully")
}
