#!/bin/bash
set -e

if [ ! -f /var/www/html/composer.json ]; then
    echo "Laravel not found, starting installation..."
    cd /var/www/html
    composer create-project --prefer-dist laravel/laravel .
else
    echo "Laravel already installed, skipping installation."
fi

# Replace DB_HOST, DB_DATABASE, etc. values from Docker ENV
echo "Updating .env with environment variables. Using connection ${DB_CONNECTION}"
sed -i "s/^DB_CONNECTION=.*/DB_CONNECTION=${DB_CONNECTION}/" .env
sed -i "s/^#* *DB_HOST=.*/DB_HOST=${DB_HOST}/" .env
sed -i "s/^#* *DB_PORT=.*/DB_PORT=${DB_PORT}/" .env
sed -i "s/^#* *DB_DATABASE=.*/DB_DATABASE=${DB_DATABASE}/" .env
sed -i "s/^#* *DB_USERNAME=.*/DB_USERNAME=${DB_USERNAME}/" .env
sed -i "s/^#* *DB_PASSWORD=.*/DB_PASSWORD=${DB_PASSWORD}/" .env

cat .env | grep DB_

# Set ownership and permissions for directories
chown -R $USER_ID:$GROUP_ID /var/www/html
echo "Setting permissions for storage and bootstrap/cache directories..."
# Create folder and log file if not exist
mkdir -p /var/www/html/storage/logs
touch /var/www/html/storage/logs/laravel.log
# Set permissions
chown -R www-data:www-data /var/www/html/storage /var/www/html/bootstrap/cache
chmod -R 775 /var/www/html/storage /var/www/html/bootstrap/cache

cd /var/www/html
echo "Reading environment from .env file..."

# Get APP_ENV from .env file
APP_ENV=$(grep ^APP_ENV= .env | cut -d '=' -f2 | tr -d '\r')

if [[ -z "$APP_ENV" ]]; then
    echo "APP_ENV not found in .env file. Using 'local' as default."
    APP_ENV=local
fi

echo "Laravel environment: $APP_ENV"

# echo "Running Laravel optimization..."
# composer validate --strict
# composer install --optimize-autoloader --no-dev

# Migrate database if needed
php artisan migrate
# Run artisan commands based on environment
if [ "$APP_ENV" = "production" ]; then
    # echo "Production mode: running caching commands..."
    echo "Production mode: running setup scripts..."
    # php artisan config:clear
    # php artisan cache:clear
    # php artisan view:clear
    # php artisan route:clear

    # php artisan config:cache
    # php artisan route:cache
    # php artisan view:cache
    composer run-script pro-setup
else
    # echo "Development mode: clearing caches..."
    echo "Development mode: running setup scripts..."
    # php artisan config:clear
    # php artisan cache:clear
    # php artisan view:clear
    # php artisan route:clear
    composer run-script dev-setup
fi

# Run Apache to keep the container alive
exec apache2-foreground