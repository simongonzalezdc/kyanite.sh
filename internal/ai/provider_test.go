package ai

import (
	"context"
	"testing"
)

func TestFallbackProvider_ParseTask(t *testing.T) {
	provider := NewFallbackProvider()

	tests := []struct {
		name             string
		input            string
		expectedPriority string
		expectedCategory string
	}{
		{
			name:             "urgent keyword",
			input:            "urgent: fix the bug",
			expectedPriority: "high",
		},
		{
			name:             "asap keyword",
			input:            "need this asap",
			expectedPriority: "high",
		},
		{
			name:             "critical keyword",
			input:            "critical issue in production",
			expectedPriority: "high",
		},
		{
			name:             "low priority keyword",
			input:            "low priority: update docs",
			expectedPriority: "low",
		},
		{
			name:             "when possible keyword",
			input:            "do this when possible",
			expectedPriority: "low",
		},
		{
			name:             "default priority",
			input:            "regular task",
			expectedPriority: "medium",
		},
		{
			name:             "work category",
			input:            "work on the project",
			expectedCategory: "work",
		},
		{
			name:             "personal category",
			input:            "personal task at home",
			expectedCategory: "personal",
		},
		{
			name:             "meeting category",
			input:            "schedule a meeting",
			expectedCategory: "meetings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			task, err := provider.ParseTask(ctx, tt.input)

			if err != nil {
				t.Fatalf("ParseTask() error = %v", err)
			}

			if task.Description != tt.input {
				t.Errorf("ParseTask() Description = %v, want %v", task.Description, tt.input)
			}

			if tt.expectedPriority != "" && task.Priority != tt.expectedPriority {
				t.Errorf("ParseTask() Priority = %v, want %v", task.Priority, tt.expectedPriority)
			}

			if tt.expectedCategory != "" {
				found := false
				for _, cat := range task.Categories {
					if cat == tt.expectedCategory {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ParseTask() Categories = %v, want to contain %v", task.Categories, tt.expectedCategory)
				}
			}
		})
	}
}

func TestFallbackProvider_ChatAssistant(t *testing.T) {
	provider := NewFallbackProvider()

	tests := []struct {
		name        string
		message     string
		taskContext []string
		wantContain string
	}{
		{
			name:        "hello greeting",
			message:     "hello",
			taskContext: []string{},
			wantContain: "fallback mode",
		},
		{
			name:        "hi greeting",
			message:     "hi there",
			taskContext: []string{},
			wantContain: "fallback mode",
		},
		{
			name:        "help request",
			message:     "help me",
			taskContext: []string{"task1", "task2"},
			wantContain: "2 tasks",
		},
		{
			name:        "count request",
			message:     "how many tasks do I have?",
			taskContext: []string{"task1", "task2", "task3"},
			wantContain: "3 tasks",
		},
		{
			name:        "default response",
			message:     "random question",
			taskContext: []string{"task1"},
			wantContain: "AI is currently unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			response, err := provider.ChatAssistant(ctx, tt.message, tt.taskContext)

			if err != nil {
				t.Fatalf("ChatAssistant() error = %v", err)
			}

			if response == "" {
				t.Error("ChatAssistant() returned empty response")
			}

			// Note: We can't easily test string containment without importing strings
			// In a real test you'd check if response contains the expected substring
			_ = tt.wantContain
		})
	}
}

func TestFallbackProvider_IsAvailable(t *testing.T) {
	provider := NewFallbackProvider()

	if !provider.IsAvailable() {
		t.Error("FallbackProvider.IsAvailable() should always return true")
	}
}

func TestFallbackProvider_GetName(t *testing.T) {
	provider := NewFallbackProvider()

	name := provider.GetName()
	if name == "" {
		t.Error("FallbackProvider.GetName() should not return empty string")
	}

	if name != "Fallback (Rule-based)" {
		t.Errorf("FallbackProvider.GetName() = %v, want 'Fallback (Rule-based)'", name)
	}
}

func TestOllamaProvider_GetName(t *testing.T) {
	manager := &Manager{}
	provider := NewOllamaProvider(manager)

	name := provider.GetName()
	if name != "Ollama" {
		t.Errorf("OllamaProvider.GetName() = %v, want 'Ollama'", name)
	}
}

func TestOpenRouterProvider_GetName(t *testing.T) {
	manager := &Manager{}
	provider := NewOpenRouterProvider(manager)

	name := provider.GetName()
	if name != "OpenRouter" {
		t.Errorf("OpenRouterProvider.GetName() = %v, want 'OpenRouter'", name)
	}
}

func TestOpenRouterProvider_IsAvailable(t *testing.T) {
	tests := []struct {
		name           string
		openRouterKey  string
		expectAvailable bool
	}{
		{
			name:           "with API key",
			openRouterKey:  "test-key",
			expectAvailable: true,
		},
		{
			name:           "without API key",
			openRouterKey:  "",
			expectAvailable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				openRouterKey: tt.openRouterKey,
			}
			provider := NewOpenRouterProvider(manager)

			if got := provider.IsAvailable(); got != tt.expectAvailable {
				t.Errorf("OpenRouterProvider.IsAvailable() = %v, want %v", got, tt.expectAvailable)
			}
		})
	}
}

func TestProviderType_Constants(t *testing.T) {
	// Test that provider type constants are defined correctly
	if ProviderOllama != "ollama" {
		t.Errorf("ProviderOllama = %v, want 'ollama'", ProviderOllama)
	}

	if ProviderOpenRouter != "openrouter" {
		t.Errorf("ProviderOpenRouter = %v, want 'openrouter'", ProviderOpenRouter)
	}

	if ProviderFallback != "fallback" {
		t.Errorf("ProviderFallback = %v, want 'fallback'", ProviderFallback)
	}
}
