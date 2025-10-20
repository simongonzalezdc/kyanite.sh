package cli

import (
	"strings"
	"testing"

	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/validation"
)

// TestAddCommand tests the add command with various inputs
func TestAddCommand(t *testing.T) {
	// Test validation
	t.Run("EmptyDescription", func(t *testing.T) {
		err := validation.ValidateTaskDescription("")
		if err == nil {
			t.Error("Expected validation error for empty description")
		}
	})

	t.Run("SQLInjection", func(t *testing.T) {
		err := validation.ValidateTaskDescription("'; DROP TABLE tasks; --")
		if err == nil {
			t.Error("Expected validation error for SQL injection attempt")
		}
	})

	t.Run("ValidDescription", func(t *testing.T) {
		err := validation.ValidateTaskDescription("Valid task description")
		if err != nil {
			t.Errorf("Unexpected validation error: %v", err)
		}
	})

	t.Run("TooLongDescription", func(t *testing.T) {
		longDesc := strings.Repeat("a", 1001)
		err := validation.ValidateTaskDescription(longDesc)
		if err == nil {
			t.Error("Expected validation error for too long description")
		}
	})
}

// TestStoreOperations tests store functionality
func TestStoreOperations(t *testing.T) {
	// Create temporary store
	tempFile := t.TempDir() + "/test_tasks.json"
	store := store.New(tempFile)

	// Test store creation
	if store == nil {
		t.Error("Failed to create store")
	}
}

// TestValidationFunctions tests validation utilities
func TestValidationFunctions(t *testing.T) {
	t.Run("ValidPriority", func(t *testing.T) {
		validPriorities := []string{"low", "medium", "high", "urgent"}
		for _, priority := range validPriorities {
			err := validation.ValidateTaskPriority(priority)
			if err != nil {
				t.Errorf("Expected no error for priority %s, got: %v", priority, err)
			}
		}
	})

	t.Run("InvalidPriority", func(t *testing.T) {
		err := validation.ValidateTaskPriority("invalid")
		if err == nil {
			t.Error("Expected validation error for invalid priority")
		}
	})

	t.Run("ValidTheme", func(t *testing.T) {
		validThemes := []string{"synthwave", "light", "plain"}
		for _, theme := range validThemes {
			err := validation.ValidateTheme(theme)
			if err != nil {
				t.Errorf("Expected no error for theme %s, got: %v", theme, err)
			}
		}
	})

	t.Run("InvalidTheme", func(t *testing.T) {
		err := validation.ValidateTheme("invalid")
		if err == nil {
			t.Error("Expected validation error for invalid theme")
		}
	})

	t.Run("ValidTimeFormat", func(t *testing.T) {
		validFormats := []string{"12h", "24h"}
		for _, format := range validFormats {
			err := validation.ValidateTimeFormat(format)
			if err != nil {
				t.Errorf("Expected no error for time format %s, got: %v", format, err)
			}
		}
	})

	t.Run("InvalidTimeFormat", func(t *testing.T) {
		err := validation.ValidateTimeFormat("invalid")
		if err == nil {
			t.Error("Expected validation error for invalid time format")
		}
	})
}

// TestSanitization tests input sanitization
func TestSanitization(t *testing.T) {
	t.Run("RemoveNullBytes", func(t *testing.T) {
		input := "test\x00injection"
		sanitized := validation.SanitizeInput(input)
		if strings.Contains(sanitized, "\x00") {
			t.Error("Expected null bytes to be removed")
		}
	})

	t.Run("LimitLength", func(t *testing.T) {
		longInput := strings.Repeat("a", 2000)
		sanitized := validation.SanitizeInput(longInput)
		if len(sanitized) > 1000 {
			t.Errorf("Expected sanitized input to be limited to 1000 chars, got %d", len(sanitized))
		}
	})

	t.Run("TrimWhitespace", func(t *testing.T) {
		input := "  test input  "
		sanitized := validation.SanitizeInput(input)
		if strings.TrimSpace(sanitized) != sanitized {
			t.Error("Expected whitespace to be trimmed")
		}
	})
}

// TestConfigurationCommands tests config command functionality
func TestConfigurationCommands(t *testing.T) {
	t.Run("ThemeValidation", func(t *testing.T) {
		err := validation.ValidateTheme("synthwave")
		if err != nil {
			t.Errorf("Expected no error for valid theme, got: %v", err)
		}

		err = validation.ValidateTheme("invalid_theme")
		if err == nil {
			t.Error("Expected error for invalid theme")
		}
	})

	t.Run("PriorityValidation", func(t *testing.T) {
		err := validation.ValidateTaskPriority("high")
		if err != nil {
			t.Errorf("Expected no error for valid priority, got: %v", err)
		}

		err = validation.ValidateTaskPriority("invalid")
		if err == nil {
			t.Error("Expected error for invalid priority")
		}
	})
}

// BenchmarkValidation benchmarks validation performance
func BenchmarkValidation(b *testing.B) {
	description := "This is a test task description for benchmarking validation performance"

	b.Run("ValidateDescription", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			validation.ValidateTaskDescription(description)
		}
	})

	b.Run("SanitizeInput", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			validation.SanitizeInput(description)
		}
	})
}
