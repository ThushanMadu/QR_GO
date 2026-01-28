package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration.
type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	MaxBodySize     int64
}

// LoadConfig loads configuration from environment variables or usage defaults.
func LoadConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		ReadTimeout:     getEnvDuration("READ_TIMEOUT", 5*time.Second),
		WriteTimeout:    getEnvDuration("WRITE_TIMEOUT", 10*time.Second),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 5*time.Second),
		MaxBodySize:     getEnvInt64("MAX_BODY_SIZE", 1024*1024), // Default 1MB
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}
	return fallback
}
