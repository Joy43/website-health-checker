package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration parameters for the application.
type Config struct {
	AppName string
	AppEnv  string
	AppPort int

	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
}

// Load reads config from environment variables (loading from .env file if it exists).
func Load() (*Config, error) {
	// Attempt to load .env file if it exists. Ignore failure if .env is missing.
	_ = godotenv.Load()

	cfg := &Config{
		AppName:       getEnv("APP_NAME", "Health Checker"),
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnvAsInt("APP_PORT", 8000),
		DBHost:        getEnv("DB_HOST", "mysql"),
		DBPort:        getEnvAsInt("DB_PORT", 3306),
		DBName:        getEnv("DB_NAME", "health_checker"),
		DBUser:        getEnv("DB_USER", "health_user"),
		DBPassword:    getEnv("DB_PASSWORD", "health_password"),
		RedisHost:     getEnv("REDIS_HOST", "redis"),
		RedisPort:     getEnvAsInt("REDIS_PORT", 6379),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate ensures all required configuration options are present and valid.
func (c *Config) Validate() error {
	if c.AppName == "" {
		return fmt.Errorf("APP_NAME is required")
	}
	if c.AppPort <= 0 {
		return fmt.Errorf("APP_PORT must be greater than 0")
	}
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBPort <= 0 {
		return fmt.Errorf("DB_PORT must be greater than 0")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.RedisHost == "" {
		return fmt.Errorf("REDIS_HOST is required")
	}
	if c.RedisPort <= 0 {
		return fmt.Errorf("REDIS_PORT must be greater than 0")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}
