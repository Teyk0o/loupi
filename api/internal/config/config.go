// Package config handles application configuration loaded from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration values for the application.
type Config struct {
	Env            string
	Port           string
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPassword     string
	DBSSL          string
	RedisHost      string
	RedisPort      string
	RedisPassword  string
	JWTSecret      string
	EncryptionKey  string
	PhotosDir      string
	BcryptCost     int
	AllowedOrigins string
}

// Load reads configuration from environment variables.
// In development, it loads from a .env file if present.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Env:            getEnv("LOUPI_ENV", "development"),
		Port:           getEnv("LOUPI_PORT", "8080"),
		DBHost:         getEnv("LOUPI_DB_HOST", "localhost"),
		DBPort:         getEnv("LOUPI_DB_PORT", "5432"),
		DBName:         getEnv("LOUPI_DB_NAME", "loupi"),
		DBUser:         getEnv("LOUPI_DB_USER", "loupi"),
		DBPassword:     getEnv("LOUPI_DB_PASSWORD", "loupi_dev"),
		DBSSL:          getEnv("LOUPI_DB_SSL", "disable"),
		RedisHost:      getEnv("LOUPI_REDIS_HOST", "localhost"),
		RedisPort:      getEnv("LOUPI_REDIS_PORT", "6379"),
		RedisPassword:  getEnv("LOUPI_REDIS_PASSWORD", ""),
		JWTSecret:      getEnv("LOUPI_JWT_SECRET", ""),
		EncryptionKey:  getEnv("LOUPI_ENCRYPTION_KEY", ""),
		PhotosDir:      getEnv("LOUPI_PHOTOS_DIR", "./photos"),
		AllowedOrigins: getEnv("LOUPI_ALLOWED_ORIGINS", "http://localhost:3000"),
	}

	bcryptCost, err := strconv.Atoi(getEnv("LOUPI_BCRYPT_COST", "12"))
	if err != nil {
		return nil, fmt.Errorf("invalid LOUPI_BCRYPT_COST: %w", err)
	}
	cfg.BcryptCost = bcryptCost

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// DatabaseURL returns the PostgreSQL connection string.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSL,
	)
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// validate checks that required configuration values are set.
func (c *Config) validate() error {
	if c.JWTSecret == "" && !c.IsDevelopment() {
		return fmt.Errorf("LOUPI_JWT_SECRET is required in production")
	}
	if c.JWTSecret == "" {
		c.JWTSecret = "dev-secret-do-not-use-in-production"
	}

	if c.EncryptionKey == "" && !c.IsDevelopment() {
		return fmt.Errorf("LOUPI_ENCRYPTION_KEY is required in production")
	}
	if c.EncryptionKey == "" {
		c.EncryptionKey = "0000000000000000000000000000000000000000000000000000000000000000"
	}

	return nil
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
