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
	"github.com/thushanmadu/qr-go/internal/qr"
	transport "github.com/thushanmadu/qr-go/internal/transport/http"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Setup Logger (JSON for production)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 3. Initialize Dependencies
	svc := qr.NewService()
	h := transport.NewHandler(svc, logger, cfg.MaxBodySize)

	// 4. Setup Router
	mux := http.NewServeMux()
	mux.HandleFunc("/generate", h.Generate)
	mux.HandleFunc("/health", h.HealthCheck)

	// 5. Setup Server with Timeouts (Security Best Practice)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: 2 * time.Second, // Protect against Slowloris
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       60 * time.Second,
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
