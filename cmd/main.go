package main

import (
	"fmt"
	"inv/config"
	"inv/internal/app"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Create application
	app, err := app.New(&config.Config{
		Server: cfg.Server,
		DB:     cfg.DB,
	})
	if err != nil {
		log.Fatal("Failed to create application:", err)
	}

	// Run application
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := app.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
