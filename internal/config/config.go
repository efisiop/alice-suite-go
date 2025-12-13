package config

import (
	"log"
	"os"
)

// Config holds application configuration
type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
	AIAPIKey  string
	Env       string // "development" or "production"
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		Port:   getEnvOrDefault("PORT", "8080"),
		DBPath: getEnvOrDefault("DB_PATH", "data/alice-suite.db"),
		AIAPIKey: getEnvOrDefault("AI_API_KEY", ""),
		Env:    getEnvOrDefault("ENV", "development"),
	}

	// JWT Secret - required in production, optional in development
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		if cfg.Env == "production" {
			log.Fatal("JWT_SECRET environment variable is required for production")
		}
		log.Println("WARNING: JWT_SECRET not set, using development fallback")
		cfg.JWTSecret = "alice-suite-go-secret-key-change-in-production"
	}

	return cfg
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Validate checks that required configuration is present
func (c *Config) Validate() error {
	if c.Port == "" {
		log.Fatal("PORT configuration is required")
	}

	if c.DBPath == "" {
		log.Fatal("DB_PATH configuration is required")
	}

	// In production, validate critical settings
	if c.Env == "production" {
		if c.JWTSecret == "" || c.JWTSecret == "alice-suite-go-secret-key-change-in-production" {
			log.Fatal("JWT_SECRET must be set to a secure value in production")
		}
	}

	return nil
}

