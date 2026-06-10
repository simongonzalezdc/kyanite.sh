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

// GetRecentSessions returns the most recent sessions for the syntax app.
// Returns nil if the brain or memory is unavailable (best-effort).
func (b *BrainClient) GetRecentSessions(ctx context.Context, limit int) ([]ai.Session, error) {
	if b.brain == nil {
		return nil, nil
	}
	return b.brain.GetRecentSessions(ctx, limit)
}

// SaveSession persists the current session state for resume-on-startup.
// Silently skipped if the brain or memory is unavailable (best-effort).
func (b *BrainClient) SaveSession(ctx context.Context, sessionID, title string, state any) error {
	if b.brain == nil {
		return nil
	}
	return b.brain.SaveSession(ctx, sessionID, title, state)
}

// SaveCrossAppContext stores a context link for other kyanite apps.
// Silently skipped if the brain or memory is unavailable (best-effort).
func (b *BrainClient) SaveCrossAppContext(ctx context.Context, targetApp, contextType, summary string, score float32) error {
	if b.brain == nil {
		return nil
	}
	return b.brain.SaveCrossAppContext(ctx, targetApp, contextType, summary, score)
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

// Brain returns the underlying ai.Brain instance, or nil if unavailable.
// This allows callers (e.g., aipanel) to use StreamChat and other
// brain-level capabilities directly.
func (b *BrainClient) Brain() *ai.Brain {
	return b.brain
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
