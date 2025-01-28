package config

import (
	"fmt"
	"os"
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
	config := &Config{
		Server: ServerConfig{
			Port: os.Getenv("APP_SERVER_PORT"),
			Host: os.Getenv("APP_SERVER_HOST"),
		},
		DB: DBConfig{
			Host:     os.Getenv("APP_DB_HOST"),
			Port:     os.Getenv("APP_DB_PORT"),
			User:     os.Getenv("APP_DB_USER"),
			Password: os.Getenv("APP_DB_PASSWORD"),
			Name:     os.Getenv("APP_DB_NAME"),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(cfg *Config) error {
	if cfg.DB.Password == "" {
		return fmt.Errorf("database password is required")
	}
	return nil
}
