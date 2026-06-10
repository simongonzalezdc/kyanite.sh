package ai

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestClient_IsAvailable_Unreachable(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	if client.IsAvailable() {
		t.Error("IsAvailable() should be false when NUCBox is unreachable")
	}
}

func TestClient_GeneratePalette_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	_, err := client.GeneratePalette(context.Background(), "warm sunset colors")
	if err == nil {
		t.Error("GeneratePalette() should return error when brain is nil")
	}
}

func TestClient_SuggestAccessibleColor_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	_, err := client.SuggestAccessibleColor(context.Background(), "#333333", "#FFFFFF", 4.5)
	if err == nil {
		t.Error("SuggestAccessibleColor() should return error when brain is nil")
	}
}

func TestClient_Close_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	// Must not panic.
	client.Close()
}

func TestClient_Close_Idempotent(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	client.Close()
	client.Close() // second close should also be safe
}
