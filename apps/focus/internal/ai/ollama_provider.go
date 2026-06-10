package ai

import (
	"context"
)

// OllamaProvider implements the Provider interface for Ollama
type OllamaProvider struct {
	manager *Manager
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(manager *Manager) *OllamaProvider {
	return &OllamaProvider{manager: manager}
}

// ParseTask parses natural language into a structured task using Ollama
func (p *OllamaProvider) ParseTask(ctx context.Context, input string) (*ParsedTask, error) {
	return p.manager.parseWithOllama(ctx, input)
}

// ChatAssistant provides conversational AI assistance using Ollama
func (p *OllamaProvider) ChatAssistant(ctx context.Context, message string, taskContext []string) (string, error) {
	return p.manager.chatWithOllama(ctx, message, taskContext)
}

// IsAvailable checks if Ollama is available
func (p *OllamaProvider) IsAvailable() bool {
	return p.manager.IsOllamaAvailable()
}

// GetName returns the provider name
func (p *OllamaProvider) GetName() string {
	return "Ollama"
}
