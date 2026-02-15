package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Config for logger behavior.
type Config struct {
	Level  string // debug, info, warn, error
	Format string // text, json
	Env    string // development, staging, production
}

// New creates a slog.Logger from config.
// Format: "json" for production/live, "text" for development (human-readable).
// Level: debug, info, warn, error (default info).
func New(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)
	opts := &slog.HandlerOptions{Level: level}

	var w io.Writer = os.Stdout
	var handler slog.Handler
	switch strings.ToLower(strings.TrimSpace(cfg.Format)) {
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	default:
		handler = slog.NewTextHandler(w, opts)
	}

	logger := slog.New(handler)
	logger = logger.With("env", cfg.Env)
	return logger
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
