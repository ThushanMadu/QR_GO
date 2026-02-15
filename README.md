# QR Go

A minimal HTTP API service that generates QR codes from text or URLs. Built with Go, configurable via environment variables, and container-ready.

## Features

- **Generate QR codes** via GET or POST
- **Configurable size** (pixels) with safe min/max bounds
- **Health check** and API info endpoints
- **Graceful shutdown** and request timeouts
- **JSON error responses** for clients
- **Docker** multi-stage build

## Quick Start

```bash
# Run locally
make run

# Or with Go
go run cmd/api/main.go
```

Then open:

- **API info:** http://localhost:8080/
- **Health:** http://localhost:8080/health
- **QR (GET):** http://localhost:8080/generate?content=Hello%20World&size=200
- **QR (POST):** `curl -X POST http://localhost:8080/generate -d "https://example.com" -o qr.png`

## API Reference

| Method | Endpoint   | Description |
|--------|------------|-------------|
| GET    | `/`        | API information and endpoint list |
| GET    | `/health`  | Health check (`{"status":"ok"}`) |
| GET    | `/generate`| Generate QR from `content` (and optional `size`) |
| POST   | `/generate`| Generate QR from request body; optional query `size` |

### Generate (GET)

- **Query parameters**
  - `content` (required): Text or URL to encode (URL-encoded).
  - `size` (optional): Image size in pixels (default `256`). Clamped to server min/max.

**Example:**

```
GET /generate?content=https%3A%2F%2Fexample.com&size=200
```

Response: `image/png`

### Generate (POST)

- **Body:** Raw bytes (e.g. plain text or URL). Max size is configurable (default 1MB).
- **Query:** `size` (optional), same as GET.

**Example:**

```bash
curl -X POST "http://localhost:8080/generate?size=300" \
  -d "https://github.com/thushanmadu/qr-go" \
  -o qr.png
```

Response: `image/png`

### Error responses

Errors are JSON:

```json
{"error": "Invalid size parameter. Must be between 64 and 512"}
```

Common status codes: `400` (bad request), `413` (body too large), `500` (server error).

## Configuration

All settings are optional and use environment variables:

| Variable           | Default   | Description                    |
|--------------------|-----------|--------------------------------|
| `PORT`             | `8080`    | HTTP server port               |
| `READ_TIMEOUT`     | `5s`      | Read timeout                   |
| `WRITE_TIMEOUT`    | `10s`     | Write timeout                  |
| `SHUTDOWN_TIMEOUT` | `5s`      | Graceful shutdown timeout      |
| `MAX_BODY_SIZE`    | `1048576` | Max POST body size (bytes, 1MB)|
| `MIN_QR_SIZE`      | `64`      | Minimum QR image size (px)     |
| `MAX_QR_SIZE`      | `512`     | Maximum QR image size (px)     |

Example:

```bash
PORT=3000 MAX_QR_SIZE=1024 go run cmd/api/main.go
```

## Build & Run

```bash
# Build binary
make build

# Run binary
./qr-microservice

# Run tests
make test

# Docker build and run
make docker-build
make docker-run
```

## Project layout

```
.
├── cmd/api/main.go          # Entrypoint, server, graceful shutdown
├── internal/
│   ├── config/config.go     # Env-based configuration
│   ├── qr/service.go        # QR generation logic
│   └── transport/http/
│       └── handler.go       # HTTP handlers
├── go.mod
├── Makefile
├── Dockerfile
└── README.md
```

## License

MIT
