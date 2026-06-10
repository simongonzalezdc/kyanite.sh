package ai

import (
	"context"
	"testing"
)

func TestNewBrainClient(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewBrainClient()
	if client == nil {
		t.Fatal("NewBrainClient() returned nil")
	}
}

func TestBrainClient_IsEnabled_Check(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewBrainClient()
	// Brain may or may not be nil depending on local whisper.cpp availability.
	// Just verify IsEnabled() doesn't panic and returns a bool.
	_ = client.IsEnabled()
}

func TestBrainClient_CheckConnection_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewBrainClient()
	err := client.CheckConnection(context.Background())
	if err == nil {
		t.Error("CheckConnection() should return error when brain is nil")
	}
}

func TestBrainClient_GetSuggestion_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewBrainClient()
	_, err := client.GetSuggestion(context.Background(), SuggestionContinue, "some content", "some context")
	if err == nil {
		t.Error("GetSuggestion() should return error when brain is nil")
	}
}

func TestBrainClient_Close_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewBrainClient()
	// Must not panic.
	client.Close()
}

func TestBrainClient_Close_Idempotent(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewBrainClient()
	client.Close()
	client.Close() // second close should also be safe
}

func TestBrainClient_ImplementsProvider(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	// Compile-time check that BrainClient satisfies the Provider interface.
	var _ Provider = NewBrainClient()
}

func TestSuggestionTypeName(t *testing.T) {
	tests := []struct {
		input    SuggestionType
		expected string
	}{
		{SuggestionContinue, "continue"},
		{SuggestionImprove, "improve"},
		{SuggestionDialogue, "dialogue"},
		{SuggestionDescription, "description"},
		{SuggestionCharacter, "character"},
		{SuggestionType(99), ""}, // unknown type
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := suggestionTypeName(tt.input)
			if got != tt.expected {
				t.Errorf("suggestionTypeName(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
