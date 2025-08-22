#!/bin/sh
set -e

# Konfigurasi database dari .env
DB_HOST=$(grep DB_HOST .env | cut -d '=' -f2)
DB_PORT=$(grep DB_PORT .env | cut -d '=' -f2)
DB_USER=$(grep DB_USER .env | cut -d '=' -f2)
DB_PASS=$(grep DB_PASS .env | cut -d '=' -f2)
DB_NAME=$(grep DB_NAME .env | cut -d '=' -f2)

# Buat DSN untuk migrate CLI
DB_DSN="mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

# Jalankan migrasi
# $1 adalah argumen yang diberikan (contoh: 'up' atau 'down')
migrate -path migrations -database "${DB_DSN}" "$1"