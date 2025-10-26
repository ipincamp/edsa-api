package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config struct to hold application configuration
type Config struct {
	App         Server
	Database    Database
	Security    Security
	Seeder      Seeder
	BloomFilter BloomFilter
	Storage     Storage
	Email       EmailConfig
}

// Server is a struct to hold server configuration
type Server struct {
	Env         string
	Port        int
	Host        string
	FrontendURL string
}

// Database is a struct to hold database configuration
type Database struct {
	Host                string
	Port                int
	User                string
	Password            string
	Name                string
	SSLMode             string
	TimeZone            string
	PoolMaxIdleConns    int
	PoolMaxOpenConns    int
	PoolConnMaxLifetime int // in minutes
}

// Security is a struct to hold security configuration
type Security struct {
	PasetoSymmetricKey string
	AccessTokenTTLMin  int
	RefreshTokenTTLMin int
	BlacklistTTLHour   int
	CorsAllowedOrigins string
}

// Seeder is a struct to hold seeder configuration
type Seeder struct {
	AdminName     string
	AdminEmail    string
	AdminPassword string
}

// BloomFilter is a struct to hold bloom filter configuration
type BloomFilter struct {
	Path                        string
	RegenerationIntervalInHours int
}

// Storage is a struct to hold storage configuration
type Storage struct {
	StoragePath          string
	StorageUploadDir     string
	StoragePublicURL     string
	StoragePublicBaseURL string
}

// EmailConfig is a struct to hold email configuration
type EmailConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	SenderEmail string
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		App: Server{
			Env:         getEnv("APP_ENV", "development"),
			Port:        getEnvAsInt("APP_PORT", 8000),
			Host:        getEnv("APP_HOST", "localhost"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
		Database: Database{
			Host:                getEnv("DB_HOST", "localhost"),
			Port:                getEnvAsInt("DB_PORT", 5432),
			User:                getEnv("DB_USER", "edsa"),
			Password:            getEnv("DB_PASSWORD", "password"),
			Name:                getEnv("DB_NAME", "edsa_staging"),
			SSLMode:             getEnv("DB_SSLMODE", "disable"),
			TimeZone:            getEnv("DB_TIMEZONE", "Asia/Jakarta"),
			PoolMaxIdleConns:    getEnvAsInt("DB_POOL_MAX_IDLE_CONNS", 10),
			PoolMaxOpenConns:    getEnvAsInt("DB_POOL_MAX_OPEN_CONNS", 100),
			PoolConnMaxLifetime: getEnvAsInt("DB_POOL_CONN_MAX_LIFETIME_MIN", 60),
		},
		Security: Security{
			PasetoSymmetricKey: getEnv("PASETO_SYMMETRIC_KEY", ""),
			AccessTokenTTLMin:  getEnvAsInt("ACCESS_TOKEN_TTL_MIN", 15),
			RefreshTokenTTLMin: getEnvAsInt("REFRESH_TOKEN_TTL_MIN", 43200),
			BlacklistTTLHour:   getEnvAsInt("BLACKLIST_TTL_HOUR", 1),
			CorsAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5000"),
		},
		Seeder: Seeder{
			AdminName:     getEnv("ADMIN_NAME", "Admin Edsa"),
			AdminEmail:    getEnv("ADMIN_EMAIL", "admin@edsa.com"),
			AdminPassword: getEnv("ADMIN_PASSWORD", "password"),
		},
		BloomFilter: BloomFilter{
			Path:                        getEnv("BLOOM_FILTER_PATH", "./bloom_filter.bin"),
			RegenerationIntervalInHours: getEnvAsInt("BLOOM_REGENERATION_INTERVAL_IN_HOURS", 24),
		},
		Storage: Storage{
			StoragePath:          getEnv("STORAGE_PATH", "./public"),
			StorageUploadDir:     getEnv("STORAGE_UPLOAD_DIR", "cdn"),
			StoragePublicURL:     getEnv("STORAGE_PUBLIC_URL", "/"),
			StoragePublicBaseURL: getEnv("STORAGE_PUBLIC_BASE_URL", "http://localhost:8000"),
		},
		Email: EmailConfig{
			Host:        getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:        getEnvAsInt("SMTP_PORT", 587),
			User:        getEnv("SMTP_USER", ""),
			Password:    getEnv("SMTP_PASSWORD", ""),
			SenderEmail: getEnv("SMTP_SENDER_EMAIL", ""),
		},
	}
	if AppConfig.Security.PasetoSymmetricKey == "" {
		log.Fatal("PASETO_SYMMETRIC_KEY is not set")
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(name string, fallback int) int {
	if valStr, exists := os.LookupEnv(name); exists {
		if v, err := strconv.Atoi(valStr); err == nil {
			return v
		}
	}
	return fallback
}

func GetDatabaseDSN() string {
	dsn := "host=" + getEnv("DB_HOST", "localhost") +
		" user=" + getEnv("DB_USER", "edsa") +
		" password=" + getEnv("DB_PASSWORD", "password") +
		" dbname=" + getEnv("DB_NAME", "edsa_staging") +
		" port=" + getEnv("DB_PORT", "5432") +
		" sslmode=" + getEnv("DB_SSLMODE", "disable") +
		" TimeZone=Asia/Jakarta"
	return dsn
}
