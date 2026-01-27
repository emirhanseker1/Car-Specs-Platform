package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	ApiNinjasKey string
	DatabaseURL  string
}

// LoadConfig reads configuration from .env file or environment variables.
func LoadConfig() (*Config, error) {
	// Attempt to load .env file, but don't fail if it doesn't exist (e.g. production env vars)
	_ = godotenv.Load()

	apiNinjasKey := os.Getenv("API_NINJAS_KEY")
	if apiNinjasKey == "" {
		return nil, fmt.Errorf("API_NINJAS_KEY is not set in environment or .env file")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Default to local sqlite if not set
		databaseURL = "vehicles.db"
	}

	return &Config{
		ApiNinjasKey: apiNinjasKey,
		DatabaseURL:  databaseURL,
	}, nil
}
