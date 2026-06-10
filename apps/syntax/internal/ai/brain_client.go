package ai

import (
	"context"
	"fmt"
	"sync"
	"time"

	ai "github.com/kyanite/ai"
)

// BrainClient wraps pkg/ai.Brain to provide the same interface as Client.
// It uses the shared kyanite AI module with centralized prompts and
// NUCBox-hosted Ollama for inference.
type BrainClient struct {
	brain *ai.Brain
	mu    sync.Mutex
}

// NewBrainClient creates a new AI client backed by pkg/ai.Brain.
// Falls back to offline mode if the brain cannot be initialized.
func NewBrainClient() *BrainClient {
	cfg := ai.DefaultConfig("syntax")
	brain, err := ai.New(cfg)
	if err != nil {
		// Brain failed to initialize — return a disabled client.
		// The app still works offline.
		return &BrainClient{}
	}
	return &BrainClient{brain: brain}
}

// IsEnabled returns whether the AI assistant is available.
func (b *BrainClient) IsEnabled() bool {
	return b.brain != nil
}

// GetSuggestion gets an AI suggestion based on the context.
// Uses pkg/ai.SyntaxSuggestPrompt for prompt construction.
func (b *BrainClient) GetSuggestion(ctx context.Context, suggestionType SuggestionType, content string, contextStr string) (*Suggestion, error) {
	if b.brain == nil {
		return nil, fmt.Errorf("%w: syntax client", ai.ErrBrainNotInitialized)
	}

	typeName := suggestionTypeName(suggestionType)
	prompt := ai.SyntaxSuggestPrompt(typeName, content, contextStr)

	response, err := b.brain.Generate(ctx, prompt,
		ai.WithTemperature(0.7),
		ai.WithMaxTokens(500),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate suggestion: %w", err)
	}

	return &Suggestion{
		Type:      suggestionType,
		Content:   response,
		Timestamp: time.Now(),
	}, nil
}

// CheckConnection verifies that the Ollama server is reachable via Brain.
func (b *BrainClient) CheckConnection(ctx context.Context) error {
	if b.brain == nil {
		return fmt.Errorf("%w: syntax client", ai.ErrBrainNotInitialized)
	}
	if !b.brain.IsLLMAvailable(ctx) {
		return fmt.Errorf("Ollama server is not reachable")
	}
	return nil
}

// Close releases Brain resources if the client owns them.
func (b *BrainClient) Close() {
	if b.brain != nil {
		b.brain.Close()
	}
}

// suggestionTypeName maps a SuggestionType to the string keys used by
// pkg/ai.SyntaxSuggestPrompt.
func suggestionTypeName(t SuggestionType) string {
	switch t {
	case SuggestionContinue:
		return "continue"
	case SuggestionImprove:
		return "improve"
	case SuggestionDialogue:
		return "dialogue"
	case SuggestionDescription:
		return "description"
	case SuggestionCharacter:
		return "character"
	default:
		return ""
	}
}
