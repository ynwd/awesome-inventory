package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("APP_SERVER_PORT", "8080"),
			Host: getEnv("APP_SERVER_HOST", "0.0.0.0"),
		},
		DB: DBConfig{
			Host:     getEnv("APP_DB_HOST", "localhost"),
			Port:     getEnv("APP_DB_PORT", "5432"),
			User:     getEnv("APP_DB_USER", "postgres"),
			Password: os.Getenv("APP_DB_PASSWORD"),
			Name:     getEnv("APP_DB_NAME", "inventory"),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func validateConfig(cfg *Config) error {
	if cfg.DB.Password == "" {
		return fmt.Errorf("database password is required")
	}
	return nil
}
