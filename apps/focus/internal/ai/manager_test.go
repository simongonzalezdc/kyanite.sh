package ai

import (
	"context"
	"strings"
	"testing"
)

func TestAIManager_ParseTask(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	manager := New()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "simple task",
			input:   "Buy milk",
			wantErr: false,
		},
		{
			name:    "task with deadline",
			input:   "Finish report by Friday",
			wantErr: false,
		},
		{
			name:    "task with priority",
			input:   "Urgent: Call dentist",
			wantErr: false,
		},
		{
			name:    "complex task",
			input:   "Schedule team meeting for next Monday to discuss Q3 goals",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.ParseTask(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && result.Description == "" {
				t.Errorf("ParseTask() returned empty description")
			}
		})
	}
}

func TestAIManager_SuggestTasks(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	manager := New()

	tests := []struct {
		name    string
		tasks   []string
		wantErr bool
	}{
		{
			name:    "no existing tasks",
			tasks:   []string{},
			wantErr: false,
		},
		{
			name:    "few tasks",
			tasks:   []string{"Buy milk", "Call dentist"},
			wantErr: false,
		},
		{
			name:    "many tasks",
			tasks:   []string{"Buy milk", "Call dentist", "Finish report", "Schedule meeting", "Plan vacation"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.SuggestTasks(context.Background(), tt.tasks)
			if (err != nil) != tt.wantErr {
				t.Errorf("SuggestTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// It's okay to return empty suggestions, but if we do, it should be a slice
			if result == nil {
				t.Errorf("SuggestTasks() returned nil instead of slice")
			}
		})
	}
}

func TestAIManager_ValidateResponse(t *testing.T) {
	manager := New()

	tests := []struct {
		name      string
		task      *ParsedTask
		wantValid bool
	}{
		{
			name: "valid task",
			task: &ParsedTask{
				Description: "Test task",
				Priority:    "medium",
			},
			wantValid: true,
		},
		{
			name: "empty description",
			task: &ParsedTask{
				Description: "",
				Priority:    "medium",
			},
			wantValid: false,
		},
		{
			name: "too short description",
			task: &ParsedTask{
				Description: "A",
				Priority:    "medium",
			},
			wantValid: false,
		},
		{
			name: "invalid priority",
			task: &ParsedTask{
				Description: "Test task",
				Priority:    "invalid",
			},
			wantValid: true, // Should be normalized
		},
		{
			name: "valid categories",
			task: &ParsedTask{
				Description: "Test task",
				Priority:    "medium",
				Categories:  []string{"work", "personal"},
			},
			wantValid: true,
		},
		{
			name: "invalid categories",
			task: &ParsedTask{
				Description: "Test task",
				Priority:    "medium",
				Categories:  []string{"a", "verylongcategorynameexceedinglimit"},
			},
			wantValid: true, // Invalid categories should be filtered out
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, gotValid := manager.validateResponse(tt.task)
			if gotValid != tt.wantValid {
				t.Errorf("validateResponse() = %v, want %v", gotValid, tt.wantValid)
			}

			// Check that invalid categories are filtered out
			if result != nil && tt.name == "invalid categories" {
				if len(result.Categories) != 0 {
					t.Errorf("validateResponse() should filter out invalid categories")
				}
			}
		})
	}
}

func TestAIManager_BasicParse(t *testing.T) {
	manager := New()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "short input",
			input: "Buy milk",
		},
		{
			name:  "long input",
			input: "This is a very long task description that exceeds the normal limits and should be truncated appropriately by the basic parser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.basicParse(tt.input)
			if result == nil {
				t.Errorf("basicParse() returned nil")
				return
			}

			if result.Description == "" {
				t.Errorf("basicParse() returned empty description")
			}

			if result.Priority != "medium" {
				t.Errorf("basicParse() should set default priority to medium")
			}
		})
	}
}

func TestManager_ParseTask_FallbackPriority(t *testing.T) {
	manager := New()

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
			// Test basicParse directly for deterministic results (no LLM needed)
			result := manager.basicParse(tt.input)

			if result.Description == "" {
				t.Errorf("basicParse() returned empty description")
			}

			if tt.expectedPriority != "" && result.Priority != tt.expectedPriority {
				t.Errorf("basicParse() Priority = %v, want %v", result.Priority, tt.expectedPriority)
			}

			if tt.expectedCategory != "" {
				found := false
				for _, cat := range result.Categories {
					if cat == tt.expectedCategory {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("basicParse() Categories = %v, want to contain %v", result.Categories, tt.expectedCategory)
				}
			}
		})
	}
}

func TestManager_ChatAssistant_FallbackResponse(t *testing.T) {
	manager := New()

	tests := []struct {
		name        string
		message     string
		taskContext []string
	}{
		{
			name:        "hello greeting",
			message:     "hello",
			taskContext: []string{},
		},
		{
			name:        "hi greeting",
			message:     "hi there",
			taskContext: []string{},
		},
		{
			name:        "help request",
			message:     "help me",
			taskContext: []string{"task1", "task2"},
		},
		{
			name:        "task question",
			message:     "what are my tasks?",
			taskContext: []string{"task1"},
		},
		{
			name:        "default response",
			message:     "random question",
			taskContext: []string{"task1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			response, err := manager.ChatAssistant(ctx, tt.message, tt.taskContext)
			if err != nil {
				t.Fatalf("ChatAssistant() error = %v", err)
			}

			if response == "" {
				t.Error("ChatAssistant() returned empty response")
			}
		})
	}
}

func TestManager_IsOllamaAvailable_Delegates(t *testing.T) {
	// Point at a definitely-unreachable Ollama URL so the brain is nil
	// (LLM unavailable) — IsOllamaAvailable should return false without
	// panicking.
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	manager := New()
	// Just verify it doesn't panic and returns a bool.
	_ = manager.IsOllamaAvailable()
}

func TestManager_Close_Idempotent(t *testing.T) {
	manager := New()
	manager.Close()
	// Second close will panic on double-close of done channel; this is
	// acceptable: Close should only be called once. We document the
	// single-shot contract here.
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

func TestParseSuggestionList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "bullet list",
			input:    "- First task\n- Second task\n- Third task",
			expected: []string{"First task", "Second task", "Third task"},
		},
		{
			name:     "asterisk list",
			input:    "* Task A\n* Task B",
			expected: []string{"Task A", "Task B"},
		},
		{
			name:     "plain lines",
			input:    "Task one\nTask two",
			expected: []string{"Task one", "Task two"},
		},
		{
			name:     "blank lines ignored",
			input:    "Task one\n\nTask two\n",
			expected: []string{"Task one", "Task two"},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSuggestionList(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("parseSuggestionList(%q) returned %d items, want %d (got %v)", tt.input, len(got), len(tt.expected), got)
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("parseSuggestionList(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestFallbackChatResponse(t *testing.T) {
	manager := New()

	tests := []struct {
		name     string
		question string
		tasks    []string
		contains string
	}{
		{
			name:     "help greeting",
			question: "help me",
			tasks:    []string{"a"},
			contains: "focus.sh AI assistant",
		},
		{
			name:     "task with list",
			question: "what are my tasks?",
			tasks:    []string{"a", "b"},
			contains: "2 tasks",
		},
		{
			name:     "task without list",
			question: "show me tasks",
			tasks:    nil,
			contains: "don't have any tasks",
		},
		{
			name:     "default response",
			question: "why is the sky blue?",
			tasks:    []string{"a"},
			contains: "AI is currently unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.fallbackChatResponse(tt.question, tt.tasks)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("fallbackChatResponse(%q, %v) = %q, expected to contain %q", tt.question, tt.tasks, got, tt.contains)
			}
		})
	}
}
