// Deprecated: Use brain_provider.go with pkg/ai instead.
package ai

import (
	"context"
)

// OpenRouterProvider implements the Provider interface for OpenRouter
type OpenRouterProvider struct {
	manager *Manager
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(manager *Manager) *OpenRouterProvider {
	return &OpenRouterProvider{manager: manager}
}

// ParseTask parses natural language into a structured task using OpenRouter
func (p *OpenRouterProvider) ParseTask(ctx context.Context, input string) (*ParsedTask, error) {
	return p.manager.parseWithOpenRouter(ctx, input)
}

// ChatAssistant provides conversational AI assistance using OpenRouter
func (p *OpenRouterProvider) ChatAssistant(ctx context.Context, message string, taskContext []string) (string, error) {
	return p.manager.chatWithOpenRouter(ctx, message, taskContext)
}

// IsAvailable checks if OpenRouter is available (requires API key)
func (p *OpenRouterProvider) IsAvailable() bool {
	return p.manager.openRouterKey != ""
}

// GetName returns the provider name
func (p *OpenRouterProvider) GetName() string {
	return "OpenRouter"
}
