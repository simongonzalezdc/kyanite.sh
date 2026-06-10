package ai

import (
	"context"
	"testing"
)

func TestNewBrainProvider(t *testing.T) {
	// Point at a definitely-unreachable Ollama URL so the provider
	// is created but the brain is nil (offline mode).
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	provider := NewBrainProvider()
	if provider == nil {
		t.Fatal("NewBrainProvider() returned nil")
	}
}

func TestBrainProvider_IsAvailable_Unreachable(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	provider := NewBrainProvider()
	if provider.IsAvailable() {
		t.Error("IsAvailable() should be false when NUCBox is unreachable")
	}
}

func TestBrainProvider_ParseTask_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	provider := NewBrainProvider()
	_, err := provider.ParseTask(context.Background(), "buy groceries")
	if err == nil {
		t.Error("ParseTask() should return error when brain is nil")
	}
}

func TestBrainProvider_ChatAssistant_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	provider := NewBrainProvider()
	_, err := provider.ChatAssistant(context.Background(), "hello", nil)
	if err == nil {
		t.Error("ChatAssistant() should return error when brain is nil")
	}
}

func TestBrainProvider_Close_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	provider := NewBrainProvider()
	// Must not panic.
	provider.Close()
}

func TestBrainProvider_Close_NilReceiver(t *testing.T) {
	// Construct one with nil brain, close twice — must not panic.
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	provider := NewBrainProvider()
	provider.Close()
	provider.Close() // second close should also be safe
}

func TestBrainProvider_GetName(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	provider := NewBrainProvider()
	name := provider.GetName()
	if name == "" {
		t.Error("GetName() should return a non-empty string")
	}
	if name != "Brain (NUCBox)" {
		t.Errorf("GetName() = %q, want %q", name, "Brain (NUCBox)")
	}
}

func TestBrainProvider_ImplementsInterface(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	// Compile-time check that BrainProvider satisfies the Provider interface.
	var _ Provider = NewBrainProvider()
}

func TestBrainProvider_Brain_Check(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	provider := NewBrainProvider()
	// Brain() may be nil or not depending on local whisper.cpp availability.
	// Just verify the method doesn't panic and returns a valid value.
	_ = provider.Brain()
}

func TestExtractJSONFromResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain JSON object",
			input:    `{"task": "buy milk"}`,
			expected: `{"task": "buy milk"}`,
		},
		{
			name:     "JSON in markdown fence",
			input:    "```json\n{\"task\": \"walk dog\"}\n```",
			expected: `{"task": "walk dog"}`,
		},
		{
			name:     "JSON wrapped in prose",
			input:    "Here is the result: {\"task\": \"code review\"} and that's it.",
			expected: `{"task": "code review"}`,
		},
		{
			name:     "no JSON at all",
			input:    "no json here",
			expected: "no json here",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractJSONFromResponse(tt.input)
			if got != tt.expected {
				t.Errorf("extractJSONFromResponse(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
