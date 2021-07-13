package providers

import (
	"github.com/google/uuid"
)

// NewUUID creates a new UUID
func NewUUID() string {
	return uuid.New().String()
}
