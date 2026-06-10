package ai

import "context"

// Provider is the interface that both Client and BrainClient satisfy.
// Use this in app code so the concrete AI backend is swappable.
type Provider interface {
	IsEnabled() bool
	GetSuggestion(ctx context.Context, suggestionType SuggestionType, content string, context string) (*Suggestion, error)
	CheckConnection(ctx context.Context) error
}
