package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ai "github.com/kyanite/ai"
)

// BrainProvider implements the Provider interface using the shared pkg/ai Brain.
// It delegates LLM calls to the NUCBox Ollama server via the unified Brain client.
type BrainProvider struct {
	brain         *ai.Brain
	promptBuilder *PromptBuilder
}

// NewBrainProvider creates a new Brain-based provider.
// It initialises a Brain with DefaultConfig("focus"), which reads
// KYANITE_OLLAMA_URL, KYANITE_MODEL, etc. from the environment.
// If the Brain cannot be created (e.g. NUCBox unreachable), the provider
// is still created but IsAvailable() returns false — the fallback chain
// will degrade gracefully.
func NewBrainProvider() *BrainProvider {
	cfg := ai.DefaultConfig("focus")
	brain, err := ai.New(cfg)
	if err != nil {
		return &BrainProvider{brain: nil, promptBuilder: NewPromptBuilder()}
	}
	return &BrainProvider{brain: brain, promptBuilder: NewPromptBuilder()}
}

// ParseTask sends a parse prompt through Brain.Generate and decodes the
// JSON response into a ParsedTask.
func (p *BrainProvider) ParseTask(ctx context.Context, input string) (*ParsedTask, error) {
	if p.brain == nil {
		return nil, fmt.Errorf("%w: focus brain provider", ai.ErrBrainNotInitialized)
	}

	prompt := p.promptBuilder.BuildParsePrompt(input)

	resp, err := p.brain.Generate(ctx, prompt, ai.WithJSONMode())
	if err != nil {
		return nil, fmt.Errorf("brain generate: %w", err)
	}

	// Extract JSON from potential markdown-wrapped response.
	cleaned := extractJSONFromResponse(resp)

	var task ParsedTask
	if err := json.Unmarshal([]byte(cleaned), &task); err != nil {
		return nil, fmt.Errorf("brain parse json: %w (raw: %s)", err, resp)
	}

	return &task, nil
}

// ChatAssistant sends a chat prompt through Brain.Generate.
func (p *BrainProvider) ChatAssistant(ctx context.Context, message string, taskContext []string) (string, error) {
	if p.brain == nil {
		return "", fmt.Errorf("%w: focus brain provider", ai.ErrBrainNotInitialized)
	}

	prompt := p.promptBuilder.BuildChatPrompt(message, taskContext)

	resp, err := p.brain.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("brain generate: %w", err)
	}

	return strings.TrimSpace(resp), nil
}

// IsAvailable checks if the Brain's LLM backend (Ollama on NUCBox) is reachable.
func (p *BrainProvider) IsAvailable() bool {
	if p.brain == nil {
		return false
	}
	return p.brain.IsLLMAvailable(context.Background())
}

// GetName returns the provider name.
func (p *BrainProvider) GetName() string {
	return "Brain (NUCBox)"
}

// Brain returns the underlying *ai.Brain for callers that need direct access
// (e.g. memory or STT features). Returns nil if unavailable.
func (p *BrainProvider) Brain() *ai.Brain {
	return p.brain
}

// Close releases Brain resources.
func (p *BrainProvider) Close() {
	if p.brain != nil {
		p.brain.Close()
	}
}

// extractJSONFromResponse extracts JSON from a response that might be wrapped
// in markdown code fences or prose.
func extractJSONFromResponse(response string) string {
	// Try to find JSON object in the response.
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start != -1 && end != -1 && end > start {
		return response[start : end+1]
	}

	return response
}
