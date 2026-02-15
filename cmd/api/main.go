package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thushanmadu/qr-go/internal/config"
	"github.com/thushanmadu/qr-go/internal/logger"
	"github.com/thushanmadu/qr-go/internal/qr"
	transport "github.com/thushanmadu/qr-go/internal/transport/http"
)

func main() {
	// 1. Load Configuration (all from env; see .env.example)
	cfg := config.LoadConfig()

	// 2. Setup Logger: text + debug in dev, JSON in prod/live
	logger := logger.New(logger.Config{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
		Env:    cfg.Env,
	})
	logger.Info("Logger initialized", "env", cfg.Env, "log_level", cfg.LogLevel, "log_format", cfg.LogFormat)

	// 3. Initialize Dependencies
	svc := qr.NewService()
	h := transport.NewHandler(svc, logger, transport.HandlerConfig{
		MaxBodySize: cfg.MaxBodySize,
		MinQRSize:   cfg.MinQRSize,
		MaxQRSize:   cfg.MaxQRSize,
		DefaultSize: cfg.DefaultQRSize,
	})

	// 4. Setup Router
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Root)
	mux.HandleFunc("/generate", h.Generate)
	mux.HandleFunc("/health", h.HealthCheck)

	// 5. Setup Server with Timeouts (all from config)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	// 6. Run Server
	go func() {
		logger.Info("Starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}
