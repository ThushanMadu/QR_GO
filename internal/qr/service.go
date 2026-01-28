package qr

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

// Service defines the business logic for QR codes.
type Service interface {
	Generate(data []byte, size int) ([]byte, error)
}

// service implements the Service interface.
type service struct{}

// NewService creates a new instance of the QR service.
func NewService() Service {
	return &service{}
}

// Generate creates a PNG QR code for the given data and size.
func (s *service) Generate(data []byte, size int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	// Generate the QR code
	// qrcode.Medium is the error recovery level
	png, err := qrcode.Encode(string(data), qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("failed to encode QR code: %w", err)
	}

	return png, nil
}
