// Package testutil provides shared test helpers for kyanite.sh apps.
package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/kyanite/ai"
)

// MockBrain wraps a real Brain backed by an httptest.Server that mimics Ollama.
// Use this in tests to verify AI behavior without needing NUCBox.
type MockBrain struct {
	Brain       *ai.Brain
	Server      *httptest.Server
	mu          sync.Mutex
	Responses   map[string]string // prompt substring -> response
	LastPrompt  string
	RequestCount int
}

// NewMockBrain creates a mock Brain backed by an httptest Ollama server.
func NewMockBrain(app string) *MockBrain {
	mb := &MockBrain{
		Responses: make(map[string]string),
	}

	mux := http.NewServeMux()
	// /api/tags is the LLMClient.IsAvailable health probe (T7-05).
	// MockBrain needs to answer 200 here or the brain reports itself
	// as unavailable and the manager falls back to basicParse.
	mux.HandleFunc("/api/tags", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"models":[]}`)
	})
	mux.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			Stream bool `json:"stream"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		mb.mu.Lock()
		mb.RequestCount++
		// Get the last user message as the prompt
		prompt := ""
		for _, m := range req.Messages {
			if m.Role == "user" {
				prompt = m.Content
			}
		}
		mb.LastPrompt = prompt
		response := ""
		for key, resp := range mb.Responses {
			if contains(prompt, key) {
				response = resp
				break
			}
		}
		if response == "" {
			response = "mock response"
		}
		// Wrap in Ollama response format so LLMClient can extract message.content.
		wrapped, _ := json.Marshal(map[string]any{
			"message": map[string]string{
				"content": response,
				"role":    "assistant",
			},
			"done": true,
		})
		mb.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(wrapped))
	})

	server := httptest.NewServer(mux)
	mb.Server = server

	// Create a Brain pointing at the mock server
	cfg := ai.Config{
		OllamaURL: server.URL,
		Model:     "mock",
		App:       app,
	}
	brain, _ := ai.New(cfg)
	mb.Brain = brain

	return mb
}

// Close shuts down the mock server and Brain.
func (mb *MockBrain) Close() {
	mb.Server.Close()
	if mb.Brain != nil {
		mb.Brain.Close()
	}
}

// SetResponse sets a mock response for prompts containing the given substring.
func (mb *MockBrain) SetResponse(promptSubstring, response string) {
	mb.mu.Lock()
	mb.Responses[promptSubstring] = response
	mb.mu.Unlock()
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsSubstr(s, substr)
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
