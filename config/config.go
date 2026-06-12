package config

import (
	"os"
)

// Config holds all configuration parameters for the application.
type Config struct {
	Port        string
	DatabaseURL string
}

// Load loads configuration from environment variables with sensible defaults.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default local PostgreSQL URL for development
		dbURL = "postgres://shantanukale@localhost:5432/user_management?sslmode=disable"
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
	}
}
