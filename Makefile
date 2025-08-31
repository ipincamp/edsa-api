# ==============================================================================
# Makefile for Go Project
# ==============================================================================

# Your application binary name
BINARY_NAME=edsa
MAIN_GO=./main.go
TIMEZONE=Asia/Jakarta

# Default command
help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## --------------------------------------
## Build & Run Commands
## --------------------------------------

clean: ## Clean up build artifacts
	@echo "Cleaning up..."
	@rm -rf ./bin/*

build: ## Compile Go code into binary
	@echo "Building binary..."
	@mkdir -p ./bin
	@go build -o ./bin/$(BINARY_NAME) $(MAIN_GO)

run: build ## Run the application in production mode
	@echo "Running in production mode..."
	@ENV=production TZ=$(TIMEZONE) ./bin/$(BINARY_NAME)

run-dev:  ## Run the application in development mode
	@echo "Running in development mode..."
	@ENV=development TZ=$(TIMEZONE) go run $(MAIN_GO)

debug: ## Run the application with Delve debugger
	@echo "Starting debugger (Delve)..."
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@dlv debug $(MAIN_GO)

## --------------------------------------
## Migration Commands
## --------------------------------------

migrate-create: ## Create a new migration file with template. Example: make migrate-create name=create_users_table
	@echo "Creating migration file with template..."
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	func_name=$$(echo "$(name)" | sed -e 's/_\([a-z]\)/\u\1/g' -e 's/^\([a-z]\)/\u\1/g'); \
	filepath=internal/migration/migrations/$${timestamp}_$(name).go; \
	printf 'package migrations\n\nimport (\n\t"github.com/go-gormigrate/gormigrate/v2"\n\t"gorm.io/gorm"\n)\n\nfunc %s() *gormigrate.Migration {\n\t// TODO: Define struct here\n\t// Example: type YourStruct struct {}\n\treturn &gormigrate.Migration{\n\t\tID: "%s",\n\t\tMigrate: func(tx *gorm.DB) error {\n\t\t\t// TODO: Implement table or column creation here\n\t\t\t// Example: return tx.AutoMigrate(&YourStruct{})\n\t\t\treturn nil\n\t\t},\n\t\tRollback: func(tx *gorm.DB) error {\n\t\t\t// TODO: Implement table or column deletion here\n\t\t\t// Example: return tx.Migrator().DropTable("your_table")\n\t\t\treturn nil\n\t\t},\n\t}\n}\n' "$$func_name" "$$timestamp" > $$filepath; \
	echo "Successfully created: $$filepath"

migrate-up: ## Run all pending migrations
	@echo "Running migrations..."
	@go run $(MAIN_GO) --migrate

migrate-down: ## Rollback the last migration
	@echo "Rolling back last migration..."
	@go run $(MAIN_GO) --rollback

## --------------------------------------
## Seeder Commands
## --------------------------------------

seed-create: ## Create a new seeder file with template. Example: make seed-create name=product
	@echo "Creating seeder file with template..."
	@if [ -z "$(name)" ]; then \
		echo "Usage: make seed-create name=<seeder_name>"; \
		exit 1; \
	fi
	@func_name=$$(echo "$(name)" | sed -e 's/_\([a-z]\)/\u\1/g' -e 's/^\([a-z]\)/\u\1/g')Seeder; \
	filepath=internal/database/seeder/seeders/$(name).seeder.go; \
	printf 'package seeders\n\nimport (\n\t"log"\n\n\t"gorm.io/gorm"\n)\n\nfunc %s(db *gorm.DB) {\n\t// TODO: Implement your seeder logic here\n\t// Use db.FirstOrCreate() to avoid duplicates\n\tlog.Println("%s ran successfully")\n}\n' "$$func_name" "$$func_name" > $$filepath; \
	echo "Successfully created: $$filepath"

db-seed: ## Run all registered seeders
	@echo "Running database seeders..."
	@go run $(MAIN_GO) --seed

.PHONY: build clean db-seed debug help migrate-create migrate-down migrate-up run run-dev seed-create
