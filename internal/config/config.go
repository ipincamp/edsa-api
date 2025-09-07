package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config adalah struktur utama untuk konfigurasi aplikasi
type Config struct {
	Database Database `mapstructure:",squash"` // Konfigurasi database
	Server   Server   `mapstructure:",squash"` // Konfigurasi server
	Paseto   Paseto   `mapstructure:",squash"` // Konfigurasi token Paseto
	Seeder   Seeder   `mapstructure:",squash"` // Konfigurasi seeder
	Bloom    Bloom    `mapstructure:",squash"` // Konfigurasi bloom filter
}

// Database adalah konfigurasi koneksi database
type Database struct {
	Host string `mapstructure:"DB_HOST"`
	Port string `mapstructure:"DB_PORT"`
	Name string `mapstructure:"DB_NAME"`
	User string `mapstructure:"DB_USERNAME"`
	Pass string `mapstructure:"DB_PASSWORD"`
	Tz   string `mapstructure:"DB_TIMEZONE"`
}

// Server adalah konfigurasi server aplikasi
type Server struct {
	Env  string `mapstructure:"SERVER_ENV"`
	Host string `mapstructure:"SERVER_HOST"`
	Port string `mapstructure:"SERVER_PORT"`
}

// Paseto adalah konfigurasi token Paseto
type Paseto struct {
	SecretKey            string        `mapstructure:"PASETO_SECRET_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"PASETO_ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"PASETO_REFRESH_TOKEN_DURATION"`
}

// Seeder adalah konfigurasi data seeder
type Seeder struct {
	Admin Admin `mapstructure:",squash"`
}

// Admin adalah konfigurasi data admin untuk seeder
type Admin struct {
	Name     string `mapstructure:"ADMIN_NAME"`
	Email    string `mapstructure:"ADMIN_EMAIL"`
	Password string `mapstructure:"ADMIN_PASSWORD"`
}

// Bloom adalah konfigurasi bloom filter
type Bloom struct {
	EmailFilterPath      string        `mapstructure:"BLOOM_EMAIL_FILTER_PATH"`
	IntervalRegeneration time.Duration `mapstructure:"BLOOM_INTERVAL_REGENERATION"`
}

// LoadConfig membaca konfigurasi dari file .env dan environment variable
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
