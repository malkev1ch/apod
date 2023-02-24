package middleware

import (
	"github.com/malkev1ch/apod/pkg/logger"
)

// Manager is middleware manager.
type Manager struct {
	logger logger.Logger
}

// NewManager creates a new Manager.
func NewManager(logger logger.Logger) *Manager {
	return &Manager{logger: logger}
}
