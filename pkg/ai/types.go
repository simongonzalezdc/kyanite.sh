package ai

import (
	"time"
)

// Message represents a single message in a conversation.
type Message struct {
	Role    string    `json:"role"`    // "system", "user", "assistant"
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

// Transcription represents the result of speech-to-text.
type Transcription struct {
	Text       string        `json:"text"`
	Language   string        `json:"language"`
	Duration   time.Duration `json:"duration"`
	Confidence float64       `json:"confidence,omitempty"`
}

// GenerateOptions configures an LLM generation request.
type GenerateOptions struct {
	Temperature  float64
	MaxTokens    int
	StopWords    []string
	SystemPrompt string
	JSONMode     bool
	Timeout      time.Duration // Per-operation timeout override (0 = use client default)
}

// GenerateOption is a functional option for GenerateOptions.
type GenerateOption func(*GenerateOptions)

// WithTemperature sets the sampling temperature.
func WithTemperature(t float64) GenerateOption {
	return func(o *GenerateOptions) { o.Temperature = t }
}

// WithMaxTokens sets the maximum tokens to generate.
func WithMaxTokens(n int) GenerateOption {
	return func(o *GenerateOptions) { o.MaxTokens = n }
}

// WithSystemPrompt sets a system prompt.
func WithSystemPrompt(p string) GenerateOption {
	return func(o *GenerateOptions) { o.SystemPrompt = p }
}

// WithJSONMode enables JSON output mode.
func WithJSONMode() GenerateOption {
	return func(o *GenerateOptions) { o.JSONMode = true }
}

// WithStopWords sets stop sequences.
func WithStopWords(words ...string) GenerateOption {
	return func(o *GenerateOptions) { o.StopWords = words }
}

// WithTimeout sets a per-operation timeout override.
func WithTimeout(d time.Duration) GenerateOption {
	return func(o *GenerateOptions) { o.Timeout = d }
}

// ContextEntry is a persisted context entry for an app session.
type ContextEntry struct {
	App       string
	SessionID string
	Key       string
	Value     string
	UpdatedAt time.Time
}
