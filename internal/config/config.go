package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration.
// All values can be set via environment variables; see .env.example.
type Config struct {
	// Server
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	ReadHeaderTimeout  time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration
	// Request/QR limits
	MaxBodySize        int64
	MinQRSize          int
	MaxQRSize          int
	DefaultQRSize      int
	// Environment and logging: development | staging | production
	Env                string
	LogLevel           string
	LogFormat          string
}

// LoadConfig loads configuration from environment variables or usage defaults.
func LoadConfig() *Config {
	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		ReadTimeout:        getEnvDuration("READ_TIMEOUT", 5*time.Second),
		WriteTimeout:       getEnvDuration("WRITE_TIMEOUT", 10*time.Second),
		ReadHeaderTimeout:  getEnvDuration("READ_HEADER_TIMEOUT", 2*time.Second),
		IdleTimeout:        getEnvDuration("IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:    getEnvDuration("SHUTDOWN_TIMEOUT", 5*time.Second),
		MaxBodySize:        getEnvInt64("MAX_BODY_SIZE", 1024*1024),
		MinQRSize:          getEnvInt("MIN_QR_SIZE", 64),
		MaxQRSize:          getEnvInt("MAX_QR_SIZE", 512),
		DefaultQRSize:      getEnvInt("DEFAULT_QR_SIZE", 256),
		Env:                getEnv("ENV", "development"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		LogFormat:          getEnv("LOG_FORMAT", ""), // empty = auto from ENV
	}
	if cfg.MinQRSize > cfg.MaxQRSize {
		cfg.MinQRSize, cfg.MaxQRSize = cfg.MaxQRSize, cfg.MinQRSize
	}
	if cfg.DefaultQRSize <= 0 {
		cfg.DefaultQRSize = 256
	}
	// Auto log format from ENV if not set
	if cfg.LogFormat == "" {
		switch cfg.Env {
		case "production", "prod", "staging", "live":
			cfg.LogFormat = "json"
		default:
			cfg.LogFormat = "text"
		}
	}
	return cfg
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

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
