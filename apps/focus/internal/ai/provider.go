package ai

import (
	"context"
)

// Provider defines the interface for AI providers
type Provider interface {
	// ParseTask parses natural language into a structured task
	ParseTask(ctx context.Context, input string) (*ParsedTask, error)

	// ChatAssistant provides conversational AI assistance
	ChatAssistant(ctx context.Context, message string, taskContext []string) (string, error)

	// IsAvailable checks if the provider is available
	IsAvailable() bool

	// GetName returns the provider name
	GetName() string
}
