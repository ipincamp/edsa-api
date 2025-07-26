FROM php:8.2-apache

RUN apt-get update && apt-get install -y \
    zip unzip git curl libpng-dev libonig-dev libxml2-dev libicu-dev \
    && apt-get install -y libzip-dev libfreetype6-dev libjpeg-dev libwebp-dev libsodium-dev \
    && docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp \
    && docker-php-ext-install pdo_mysql mbstring exif pcntl bcmath gd intl sodium zip \
    && docker-php-ext-enable opcache

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

COPY --from=node:lts-jod /usr/local /usr/local

RUN a2enmod rewrite

COPY docker/apache/default.conf /etc/apache2/sites-available/000-default.conf

COPY docker/init.sh /usr/local/bin/init.sh

RUN chmod +x /usr/local/bin/init.sh

ENTRYPOINT ["/usr/local/bin/init.sh"]
