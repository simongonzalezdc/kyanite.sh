package integration

import (
	"context"
	"testing"

	"github.com/kyanite/focus/internal/ai"
)

// TestAIProviderIntegration tests AI provider switching and fallback behavior
func TestAIProviderIntegration(t *testing.T) {
	manager := ai.New()
	ctx := context.Background()

	t.Run("parse_task_with_fallback", func(t *testing.T) {
		input := "urgent: complete the report by tomorrow"

		result, err := manager.ParseTask(ctx, input)
		if err != nil {
			t.Fatalf("Failed to parse task: %v", err)
		}

		if result == nil {
			t.Fatal("Expected parsed task, got nil")
		}

		if result.Description == "" {
			t.Error("Description should not be empty")
		}

		// Fallback provider should detect "urgent" keyword
		if result.Priority != "high" {
			t.Errorf("Expected high priority for urgent task, got %s", result.Priority)
		}
	})

	t.Run("parse_task_with_keywords", func(t *testing.T) {
		testCases := []struct {
			name             string
			input            string
			expectedPriority string
		}{
			{
				name:             "urgent keyword",
				input:            "urgent: fix production bug",
				expectedPriority: "high",
			},
			{
				name:             "asap keyword",
				input:            "need to fix this asap",
				expectedPriority: "high",
			},
			{
				name:             "critical keyword",
				input:            "critical security patch",
				expectedPriority: "high",
			},
			{
				name:             "low priority keyword",
				input:            "low priority: refactor code",
				expectedPriority: "low",
			},
			{
				name:             "when possible keyword",
				input:            "update docs when possible",
				expectedPriority: "low",
			},
			{
				name:             "no keyword",
				input:            "regular task",
				expectedPriority: "medium",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := manager.ParseTask(ctx, tc.input)
				if err != nil {
					t.Fatalf("Failed to parse task: %v", err)
				}

				if result.Priority != tc.expectedPriority {
					t.Errorf("Expected priority %s, got %s", tc.expectedPriority, result.Priority)
				}
			})
		}
	})

	t.Run("parse_task_with_categories", func(t *testing.T) {
		testCases := []struct {
			name             string
			input            string
			expectedCategory string
		}{
			{
				name:             "work category",
				input:            "work meeting at 2pm",
				expectedCategory: "work",
			},
			{
				name:             "personal category",
				input:            "personal: call mom",
				expectedCategory: "personal",
			},
			{
				name:             "meeting category",
				input:            "meeting with team",
				expectedCategory: "meeting",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := manager.ParseTask(ctx, tc.input)
				if err != nil {
					t.Fatalf("Failed to parse task: %v", err)
				}

				found := false
				for _, cat := range result.Categories {
					if cat == tc.expectedCategory {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("Expected category %s to be present in %v", tc.expectedCategory, result.Categories)
				}
			})
		}
	})

	t.Run("suggest_tasks_fallback", func(t *testing.T) {
		existingTasks := []string{
			"Write unit tests",
			"Review pull requests",
		}

		suggestions, err := manager.SuggestTasks(ctx, existingTasks)
		if err != nil {
			t.Fatalf("Failed to suggest tasks: %v", err)
		}

		// Fallback provider should return helpful suggestions
		if len(suggestions) == 0 {
			t.Error("Expected at least one suggestion")
		}
	})

	t.Run("chat_assistant_fallback", func(t *testing.T) {
		testCases := []struct {
			name          string
			question      string
			shouldContain string
		}{
			{
				name:          "greeting",
				question:      "hello",
				shouldContain: "Hello",
			},
			{
				name:          "help request",
				question:      "help",
				shouldContain: "help",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				response, err := manager.ChatAssistant(ctx, tc.question, []string{})
				if err != nil {
					t.Fatalf("Failed to chat: %v", err)
				}

				if response == "" {
					t.Error("Expected non-empty response")
				}
			})
		}
	})

	t.Run("validation_filters_invalid_categories", func(t *testing.T) {
		input := "task with x and abcdefghijklmnopqrstuvwxyz categories"

		result, err := manager.ParseTask(ctx, input)
		if err != nil {
			t.Fatalf("Failed to parse task: %v", err)
		}

		// Check that short categories (< 2 chars) are filtered
		for _, cat := range result.Categories {
			if len(cat) < 2 {
				t.Errorf("Category %s is too short (< 2 chars)", cat)
			}
			if len(cat) > 20 {
				t.Errorf("Category %s is too long (> 20 chars)", cat)
			}
		}
	})

	t.Run("validation_limits_categories", func(t *testing.T) {
		// Create input with many category keywords
		input := "work personal meeting project task todo urgent critical high medium low priority categories test"

		result, err := manager.ParseTask(ctx, input)
		if err != nil {
			t.Fatalf("Failed to parse task: %v", err)
		}

		// Should limit to 10 categories
		if len(result.Categories) > 10 {
			t.Errorf("Expected max 10 categories, got %d", len(result.Categories))
		}
	})
}
