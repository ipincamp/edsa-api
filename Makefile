.PHONY: build run run-prod new-migration migrate-up migrate-down seed clean

# Variabel yang digunakan dalam Makefile
NAME ?= new_migration
MIGRATE_SCRIPT = ./scripts/migrate.sh
SEEDER_APP = ./cmd/seeder/main.go
MAIN_APP = ./cmd/app/main.go

# Build: Membuat executable aplikasi
build:
	@echo "Building application..."
	go build -o ./bin/server $(MAIN_APP)
	@echo "Build successful! Executable is at ./bin/server"

# Run: Menjalankan aplikasi di mode pengembangan
run:
	@echo "Starting application..."
	go run $(MAIN_APP)

# Run-prod: Menjalankan executable yang sudah di-build
run-prod: build
	@echo "Running application in production mode..."
	APP_ENV=production ./bin/server

# Migrasi: Membuat file migrasi baru
# Contoh: make new-migration name=create_products_table
new-migration:
	@echo "Creating new migration: $(NAME)"
	migrate create -ext sql -dir migrations -seq $(NAME)

# Migrasi: Menjalankan migrasi database
migrate-up:
	@echo "Running database migrations UP..."
	$(MIGRATE_SCRIPT) up

# Migrasi: Rollback satu langkah migrasi
migrate-down:
	@echo "Running database migrations DOWN..."
	$(MIGRATE_SCRIPT) down -1

# Seeder: Menjalankan seeder database
seed:
	@echo "Running database seeder..."
	go run $(SEEDER_APP)

# Clean: Membersihkan file sementara dan cache
clean:
	@echo "Cleaning temporary files and cache..."
	go clean -modcache
	rm -rf ./bin
	@echo "Clean complete!"

.DEFAULT_GOAL := run