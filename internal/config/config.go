package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database Database `mapstructure:",squash"`
	Server   Server   `mapstructure:",squash"`
	Paseto   Paseto   `mapstructure:",squash"`
	Seeder   Seeder   `mapstructure:",squash"`
	Bloom    Bloom    `mapstructure:",squash"`
}

type Database struct {
	Host string `mapstructure:"DB_HOST"`
	Port string `mapstructure:"DB_PORT"`
	Name string `mapstructure:"DB_NAME"`
	User string `mapstructure:"DB_USERNAME"`
	Pass string `mapstructure:"DB_PASSWORD"`
	Tz   string `mapstructure:"DB_TIMEZONE"`
}

type Server struct {
	Env  string `mapstructure:"SERVER_ENV"`
	Host string `mapstructure:"SERVER_HOST"`
	Port string `mapstructure:"SERVER_PORT"`
}

type Paseto struct {
	SecretKey            string        `mapstructure:"PASETO_SECRET_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"PASETO_ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"PASETO_REFRESH_TOKEN_DURATION"`
}

type Seeder struct {
	Admin Admin `mapstructure:",squash"`
}

type Admin struct {
	Name     string `mapstructure:"ADMIN_NAME"`
	Email    string `mapstructure:"ADMIN_EMAIL"`
	Password string `mapstructure:"ADMIN_PASSWORD"`
}

type Bloom struct {
	EmailFilterPath      string        `mapstructure:"BLOOM_EMAIL_FILTER_PATH"`
	IntervalRegeneration time.Duration `mapstructure:"BLOOM_INTERVAL_REGENERATION"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
