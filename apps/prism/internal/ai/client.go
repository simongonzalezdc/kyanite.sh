package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ai "github.com/kyanite/ai"
)

// Client wraps the shared kyanite/ai Brain for Prism's AI features.
//
// It provides two capabilities:
//   - Palette generation from natural language descriptions
//   - Accessible color suggestions that meet WCAG contrast requirements
//
// When the NUCBox inference server is unreachable, all methods return
// graceful errors so the app continues to work offline.
type Client struct {
	brain *ai.Brain
}

// NewClient creates a new AI client for Prism.
// If the Brain cannot be initialized (e.g. config error), the client
// is still usable — IsAvailable will return false and calls will
// return errors.
func NewClient() *Client {
	brain, err := ai.New(ai.DefaultConfig("prism"))
	if err != nil {
		return &Client{brain: nil}
	}
	return &Client{brain: brain}
}

// PaletteResult holds the result of an AI palette generation.
type PaletteResult struct {
	Colors []string `json:"colors"`
	Name   string   `json:"name"`
}

// GeneratePalette generates a 5-color palette from a text description
// using the shared PrismPalettePrompt. Returns hex color codes and a
// palette name.
func (c *Client) GeneratePalette(ctx context.Context, description string) ([]string, error) {
	if c.brain == nil {
		return nil, fmt.Errorf("%w: prism client", ai.ErrBrainNotInitialized)
	}

	if !c.brain.IsLLMAvailable(ctx) {
		return nil, fmt.Errorf("ai: NUCBox Ollama server unreachable")
	}

	prompt := ai.PrismPalettePrompt(description)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := c.brain.Generate(ctx, prompt, ai.WithJSONMode())
	if err != nil {
		return nil, fmt.Errorf("ai: palette generation failed: %w", err)
	}

	// The LLM may wrap the JSON in markdown fences; strip them.
	cleaned := strings.TrimSpace(resp)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var result PaletteResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("ai: failed to parse palette response: %w", err)
	}

	if len(result.Colors) == 0 {
		return nil, fmt.Errorf("ai: no colors returned")
	}

	return result.Colors, nil
}

// SuggestAccessibleColor suggests a foreground color adjustment that meets
// WCAG AA contrast requirements. It uses the shared PrismContrastPrompt.
func (c *Client) SuggestAccessibleColor(ctx context.Context, fg, bg string, ratio float64) (string, error) {
	if c.brain == nil {
		return "", fmt.Errorf("%w: prism client", ai.ErrBrainNotInitialized)
	}

	if !c.brain.IsLLMAvailable(ctx) {
		return "", fmt.Errorf("ai: NUCBox Ollama server unreachable")
	}

	prompt := ai.PrismContrastPrompt(fg, bg, ratio)

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	resp, err := c.brain.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("ai: color suggestion failed: %w", err)
	}

	// The response should be a single hex color code.
	hex := strings.TrimSpace(resp)
	hex = strings.TrimPrefix(hex, "`")
	hex = strings.TrimSuffix(hex, "`")
	hex = strings.TrimSpace(hex)

	// Validate it looks like a hex color.
	if !strings.HasPrefix(hex, "#") || len(hex) != 7 {
		return "", fmt.Errorf("ai: unexpected color format: %q", hex)
	}

	return strings.ToUpper(hex), nil
}

// IsAvailable returns true if the AI brain is initialized and the LLM
// server is reachable.
func (c *Client) IsAvailable() bool {
	if c.brain == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return c.brain.IsLLMAvailable(ctx)
}

// Close releases resources held by the AI client.
func (c *Client) Close() {
	if c.brain != nil {
		c.brain.Close()
	}
}
