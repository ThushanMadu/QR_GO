package http

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/thushanmadu/qr-go/internal/qr"
)

// HandlerConfig configures the HTTP handler.
type HandlerConfig struct {
	MaxBodySize int64
	MinQRSize   int
	MaxQRSize   int
	DefaultSize int
}

type Handler struct {
	svc    qr.Service
	logger *slog.Logger
	cfg    HandlerConfig
}

func NewHandler(svc qr.Service, logger *slog.Logger, cfg HandlerConfig) *Handler {
	if cfg.DefaultSize <= 0 {
		cfg.DefaultSize = 256
	}
	return &Handler{
		svc:    svc,
		logger: logger,
		cfg:    cfg,
	}
}

func (h *Handler) writeJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Generate handles the QR code generation request.
// POST /generate: body = raw content; optional query: size=N
// GET /generate: query content=... (URL-encoded); optional size=N
func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed. Use GET or POST.")
		return
	}

	var body []byte
	switch r.Method {
	case http.MethodGet:
		content := r.URL.Query().Get("content")
		if content == "" {
			h.writeJSONError(w, http.StatusBadRequest, "Missing query parameter: content")
			return
		}
		decoded, err := url.QueryUnescape(content)
		if err != nil {
			decoded = content
		}
		body = []byte(decoded)
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxBodySize)
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			h.logger.Error("failed to read request body", "error", err)
			if err.Error() == "http: request body too large" {
				h.writeJSONError(w, http.StatusRequestEntityTooLarge, "Request body too large")
				return
			}
			h.writeJSONError(w, http.StatusInternalServerError, "Failed to read request body")
			return
		}
		_ = r.Body.Close()
		if len(body) == 0 {
			h.writeJSONError(w, http.StatusBadRequest, "Request body is empty")
			return
		}
	}

	size := h.parseSize(r.URL.Query().Get("size"))
	if size <= 0 {
		h.writeJSONError(w, http.StatusBadRequest, "Invalid size parameter. Must be between "+strconv.Itoa(h.cfg.MinQRSize)+" and "+strconv.Itoa(h.cfg.MaxQRSize))
		return
	}

	png, err := h.svc.Generate(body, size)
	if err != nil {
		h.logger.Error("failed to generate QR code", "error", err)
		h.writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(png); err != nil {
		h.logger.Error("failed to write response", "error", err)
		return
	}
}

func (h *Handler) parseSize(s string) int {
	if s == "" {
		return h.cfg.DefaultSize
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < h.cfg.MinQRSize || n > h.cfg.MaxQRSize {
		return 0
	}
	return n
}

// HealthCheck is a simple endpoint to verify the service is running.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Root returns API information.
func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"service": "QR Go",
		"version": "1.0",
		"endpoints": map[string]string{
			"GET  /":        "API info",
			"GET  /health":  "Health check",
			"GET  /generate": "Generate QR (query: content, size)",
			"POST /generate": "Generate QR (body: content; query: size)",
		},
	})
}
