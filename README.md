# EDSA API

## Description

This is a Laravel 12 API project for the EDSA (English Digital Storytelling Application) project. It is designed to provide a backend for the EDSA application, which is a platform for creating and sharing digital stories in English.
The API is built using Laravel 12 and follows RESTful principles. It provides endpoints for managing users, stories, and other resources related to the EDSA project.

## Features
- User authentication and authorization
- CRUD operations for stories
- CRUD operations for glosaries
- CRUD operations for users
- File uploads for audio files
- Search functionality for stories
- Pagination for large datasets
- Rate limiting for API requests
- CORS support for cross-origin requests
- API versioning
- Scramble for API documentation

## Requirements

- **Composer** `v2.8.12`
- **PHP** `v8.4.1-NTS`
- **Laravel Installer** `v5.18.0`
- **Git** `v2.51.1`
- **MariaDB** `v10.11.13`
- **NodeJS** `v22.20.0`
- **NPM** `v10.9.3`

## Running on Docker

- Select stage of `docker-compose.yml` inside *docker/docker-compose*
```js
development.yml // for development stage

production.yml // for production stage
```

- Copy and paste in root project directory, then rename to `docker-compose.yml`

- Build the image:
```bash
docker compose up -d --build
```

## Installation

> **Note:** This project is a work in progress and may not be fully functional.

- Clone the repository:

```bash
git clone https://github.com/ipincamp/laravel-12-edsa-api.git
cd laravel-12-edsa-api
```

- Copy the `.env.example` file to `.env` and update the database configuration:

```bash
cp .env.example .env
```

- Fill in the `.env` file with your database credentials and other environment variables.

- For development stage:
```bash
composer run dev-setup
```

- For production stage:
```bash
composer run pro-setup
```

- Endpoints:
  - API: `http://localhost:8000/api/v1`
  - Web: `http://localhost:8000`

## Additional Setup for Docker
- If you are running the application in a Docker container, you may need to set the correct permissions for the storage and bootstrap/cache directories:

```bash
sudo chown -R $USER:$USER api
cd api
sudo chown -R www-data:www-data storage bootstrap/cache
sudo chmod -R 755 storage bootstrap/cache
```

## License

This project is closed-sourced software licensed under the [MIT license](LICENSE).
