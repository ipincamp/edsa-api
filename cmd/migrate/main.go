package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/database/migrations"
	"github.com/ipincamp/go-edsa-api/internal/pkg/database"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"gorm.io/gorm"
)

// dropAllTables dynamically drops all tables based on GORM models.
func dropAllTables(db *gorm.DB) error {
	log.Println("🧹 Dropping tables dynamically...") // Changed log message

	models := []interface{}{
		&repo.ApprovalRequestGORM{},
		&repo.UserInteractionAttemptGORM{},
		&repo.GroupBookSettingGORM{},
		&repo.MediaAssetGORM{},
		&repo.ActivityLogGORM{},
		&repo.UserGameScoreGORM{},
		&repo.UserBookProgressGORM{},
		&repo.InteractionGORM{},
		&repo.UserGroupGORM{},
		&repo.PageGORM{},
		&repo.GroupGORM{},
		&repo.ClassGORM{},
		&repo.UserGORM{},
		&repo.GameGORM{},
		&repo.BookGORM{},
		&repo.SubjectGORM{},
		&repo.RoleGORM{},
	}

	if err := db.Migrator().DropTable(models...); err != nil {
		log.Printf("❌ Error dropping GORM model tables: %v", err) // Changed log message
		return err
	}
	log.Println("✅ Dynamically dropped GORM model tables.") // Changed log message

	migrationsTableName := gormigrate.DefaultOptions.TableName
	if db.Migrator().HasTable(migrationsTableName) {
		if err := db.Migrator().DropTable(migrationsTableName); err != nil {
			log.Printf("⚠️ Could not drop migrations table '%s': %v", migrationsTableName, err) // Changed log message
			// Consider if this should be a fatal error depending on your needs
		} else {
			log.Printf("🗑️ Dropped migrations table: %s", migrationsTableName) // Changed log message
		}
	} else {
		log.Printf("ℹ️ Migrations table '%s' does not exist, skipping drop.", migrationsTableName) // Changed log message
	}

	return nil
}

// Function to ask for confirmation
func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/N]: ", prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading confirmation: %v", err)
			return false
		}
		response = strings.ToLower(strings.TrimSpace(response))
		switch response {
		case "y", "yes":
			return true
		case "", "n", "no":
			return false
		}
	}
}

func main() {
	// Load config first
	config.LoadConfig()
	cfg := config.AppConfig

	log.Println("🔧 Initializing database connection...")
	db, err := database.NewPostgresConnection(cfg) // Pass config
	if err != nil {                                // Handle error
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("🔗 Database connection established.") // Log success after check

	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetAllMigrations())

	if len(os.Args) < 2 {
		log.Fatal("❌ Missing command. Usage: go run cmd/migrate/main.go [up|down|reset|drop-all]")
	}

	command := os.Args[1]
	log.Printf("🚀 Executing command: '%s'", command) // Log command execution start

	switch command {
	case "up":
		log.Println("⬆️ Running pending migrations...")
		if err := m.Migrate(); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}
		log.Println("✅ Migrations completed successfully.")
	case "down":
		log.Println("⬇️ Rolling back the last migration...")
		if err := m.RollbackLast(); err != nil {
			log.Fatalf("❌ Rollback failed: %v", err)
		}
		log.Println("✅ Rollback successful.")
	case "reset":
		log.Println("🔄 Preparing database reset...")
		if !askForConfirmation("⚠️ WARNING: This will DROP ALL tables and re-run ALL migrations. Are you sure?") {
			log.Println("🛑 Database reset cancelled by user.")
			os.Exit(0)
		}

		log.Println("🔥 Resetting database (dropping all tables)...")
		if err := dropAllTables(db); err != nil {
			log.Fatalf("❌ Database reset failed during table drop: %v", err)
		}
		log.Println("🗑️ All application tables dropped.")

		log.Println("⬆️ Running all migrations after reset...")
		if err := m.Migrate(); err != nil {
			log.Fatalf("❌ Could not run migrations after reset: %v", err)
		}
		log.Println("✅ Migrations ran successfully after reset.")
		log.Println("✨ Database reset complete.")

	case "drop-all":
		log.Println("💣 Preparing to drop all tables...")
		if !askForConfirmation("🔥🔥 DANGER ZONE: This will permanently DROP ALL application and migration tables. Data will be lost! Are you absolutely sure?") {
			log.Println("🛑 Drop all tables cancelled by user.")
			os.Exit(0)
		}

		log.Println("🔥🔥 DANGER: Dropping all application tables (including migrations table)...")
		if err := dropAllTables(db); err != nil {
			log.Fatalf("❌ Could not drop all tables: %v", err)
		}
		log.Println("✅ All application tables dropped successfully. Database is now empty of these tables.")

	default:
		log.Fatalf("❓ Unknown command: '%s'. Use 'up', 'down', 'reset', or 'drop-all'.", command)
	}
}
