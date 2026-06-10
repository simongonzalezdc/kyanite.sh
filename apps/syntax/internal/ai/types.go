package ai

import "time"

// SuggestionType represents the type of AI suggestion
type SuggestionType int

const (
	SuggestionContinue SuggestionType = iota
	SuggestionImprove
	SuggestionDialogue
	SuggestionDescription
	SuggestionCharacter
)

// Suggestion represents an AI-generated suggestion
type Suggestion struct {
	Type      SuggestionType
	Content   string
	Timestamp time.Time
}

// Config represents AI assistant configuration
type Config struct {
	Enabled     bool
	Model       string
	APIEndpoint string
	Temperature float64
	MaxTokens   int
	Timeout     time.Duration
}

// DefaultConfig returns the default AI configuration
func DefaultConfig() Config {
	return Config{
		Enabled:     false, // Disabled by default for privacy
		Model:       "llama3.2:3b",
		APIEndpoint: "http://localhost:11434",
		Temperature: 0.7,
		MaxTokens:   500,
		Timeout:     30 * time.Second,
	}
}
