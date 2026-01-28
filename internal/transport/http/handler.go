package http

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/thushanmadu/qr-go/internal/qr"
)

type Handler struct {
	svc         qr.Service
	logger      *slog.Logger
	maxBodySize int64
}

func NewHandler(svc qr.Service, logger *slog.Logger, maxBodySize int64) *Handler {
	return &Handler{
		svc:         svc,
		logger:      logger,
		maxBodySize: maxBodySize,
	}
}

// Generate handles the QR code generation request.
func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit request body size to prevent DoS
	r.Body = http.MaxBytesReader(w, r.Body, h.maxBodySize)

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read request body", "error", err)
		if err.Error() == "http: request body too large" {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Parse size query param, default to 256
	size := 256
	sizeStr := r.URL.Query().Get("size")
	if sizeStr != "" {
		parsedSize, err := strconv.Atoi(sizeStr)
		if err != nil || parsedSize <= 0 {
			http.Error(w, "Invalid size parameter", http.StatusBadRequest)
			return
		}
		size = parsedSize
	}

	// Generate QR code
	png, err := h.svc.Generate(body, size)
	if err != nil {
		h.logger.Error("failed to generate QR code", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(png); err != nil {
		h.logger.Error("failed to write response", "error", err)
		return
	}
}

// HealthCheck is a simple endpoint to verify the service is running.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
