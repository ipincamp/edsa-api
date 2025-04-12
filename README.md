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

- PHP 8.0 or higher
- Composer
- MySQL
- Node.js and npm (for frontend assets)

## Installation

> **Note:** This project is a work in progress and may not be fully functional.

- Clone the repository:

```bash
git clone https://github.com/ipincamp/laravel-12-edsa-api.git
cd laravel-12-edsa-api
```

- Install dependencies:

```bash
composer install
npm install
```

- Copy the `.env.example` file to `.env` and update the database configuration:

```bash
cp .env.example .env
```

- Fill in the `.env` file with your database credentials and other environment variables.

- Generate the application key:

```bash
php artisan key:generate
```

- Run the migrations and seed the database:

```bash
php artisan migrate --seed
```

- Start the development server:

```bash
php artisan serve
```

- Endpoints:
  - API: `http://localhost:8000/api/v1`
  - Web: `http://localhost:8000`

## License

This project is closed-sourced software licensed under the [MIT license](LICENSE).
