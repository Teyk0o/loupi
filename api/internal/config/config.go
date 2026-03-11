// Package config handles application configuration loaded from environment variables.
package config

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
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
	CookieDomain   string
	CookieSecure   bool
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
		AllowedOrigins: getEnv("LOUPI_ALLOWED_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000"),
		CookieDomain:   getEnv("LOUPI_COOKIE_DOMAIN", ""),
	}

	bcryptCost, err := strconv.Atoi(getEnv("LOUPI_BCRYPT_COST", "12"))
	if err != nil {
		return nil, fmt.Errorf("invalid LOUPI_BCRYPT_COST: %w", err)
	}
	cfg.BcryptCost = bcryptCost

	cfg.CookieSecure = !cfg.IsDevelopment()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// DatabaseURL returns the PostgreSQL connection string with properly encoded credentials.
func (c *Config) DatabaseURL() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   fmt.Sprintf("%s:%s", c.DBHost, c.DBPort),
		Path:   c.DBName,
	}
	q := u.Query()
	q.Set("sslmode", c.DBSSL)
	u.RawQuery = q.Encode()
	return u.String()
}

// RedisAddr returns the Redis server address.
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// validate checks that required configuration values are set.
func (c *Config) validate() error {
	if c.JWTSecret == "" {
		if c.IsDevelopment() {
			log.Println("WARNING: LOUPI_JWT_SECRET is not set — using insecure dev default. DO NOT use in production.")
			c.JWTSecret = "dev-secret-do-not-use-in-production-change-me"
		} else {
			return fmt.Errorf("LOUPI_JWT_SECRET is required in production")
		}
	}
	if len(c.JWTSecret) < 32 && !c.IsDevelopment() {
		return fmt.Errorf("LOUPI_JWT_SECRET must be at least 32 characters")
	}

	if c.EncryptionKey == "" {
		if c.IsDevelopment() {
			log.Println("WARNING: LOUPI_ENCRYPTION_KEY is not set — using insecure dev default. DO NOT use in production.")
			c.EncryptionKey = "0102030405060708091011121314151617181920212223242526272829303132"
		} else {
			return fmt.Errorf("LOUPI_ENCRYPTION_KEY is required in production")
		}
	}
	if _, err := hex.DecodeString(c.EncryptionKey); err != nil || len(c.EncryptionKey) != 64 {
		return fmt.Errorf("LOUPI_ENCRYPTION_KEY must be exactly 64 hex characters (32 bytes for AES-256)")
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
