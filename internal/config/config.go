package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiration time.Duration
}

func Load() (*Config, error) {
	// Load environment variables from .env
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	jwtExpirationString := os.Getenv("JWT_EXPIRATION")
	if jwtExpirationString == "" {
		jwtExpirationString = "24h"
	}

	jwtExpiration, err := time.ParseDuration(jwtExpirationString)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION: %w", err)
	}

	return &Config{
		Port:          port,
		DatabaseURL:   databaseURL,
		JWTSecret:     jwtSecret,
		JWTExpiration: jwtExpiration,
	}, nil
}
