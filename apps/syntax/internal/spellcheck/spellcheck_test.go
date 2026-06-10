package spellcheck

import (
	"testing"
)

func TestNewChecker(t *testing.T) {
	checker := NewChecker()

	if checker == nil {
		t.Fatal("NewChecker() returned nil")
	}

	// Should be disabled by default
	if checker.IsEnabled() {
		t.Error("NewChecker() should create a disabled checker")
	}

	if checker.language != "en_US" {
		t.Errorf("NewChecker() language = %q, expected %q", checker.language, "en_US")
	}
}

func TestToggle(t *testing.T) {
	checker := NewChecker()

	// Start disabled
	if checker.IsEnabled() {
		t.Error("Checker should start disabled")
	}

	// Toggle to enabled
	checker.Toggle()
	if !checker.IsEnabled() {
		t.Error("After Toggle(), checker should be enabled")
	}

	// Toggle back to disabled
	checker.Toggle()
	if checker.IsEnabled() {
		t.Error("After second Toggle(), checker should be disabled")
	}
}

func TestSetEnabled(t *testing.T) {
	checker := NewChecker()

	checker.SetEnabled(true)
	if !checker.IsEnabled() {
		t.Error("SetEnabled(true) should enable checker")
	}

	checker.SetEnabled(false)
	if checker.IsEnabled() {
		t.Error("SetEnabled(false) should disable checker")
	}
}

func TestCheckWord(t *testing.T) {
	checker := NewChecker()

	// When disabled, all words should pass
	t.Run("disabled checker", func(t *testing.T) {
		if !checker.CheckWord("thisisnotaword") {
			t.Error("Disabled checker should return true for any word")
		}
		if !checker.CheckWord("correct") {
			t.Error("Disabled checker should return true for any word")
		}
	})

	// When enabled, check if aspell is available
	checker.SetEnabled(true)
	if !IsAvailable() {
		t.Skip("Aspell not available, skipping enabled tests")
	}

	t.Run("enabled checker with correct word", func(t *testing.T) {
		// Common words that should be in the dictionary
		correctWords := []string{"hello", "world", "test", "computer"}
		for _, word := range correctWords {
			if !checker.CheckWord(word) {
				t.Errorf("CheckWord(%q) = false, expected true for common word", word)
			}
		}
	})

	t.Run("enabled checker with empty string", func(t *testing.T) {
		// Empty strings should pass (handled specially)
		if !checker.CheckWord("") {
			t.Error("CheckWord(\"\") should return true for empty string")
		}
	})

	t.Run("enabled checker with numbers", func(t *testing.T) {
		// Numbers might be considered incorrect by aspell
		// This test documents the behavior but doesn't enforce it
		checker.CheckWord("12345")
	})
}

func TestGetSuggestions(t *testing.T) {
	checker := NewChecker()

	// When disabled, should return empty slice
	t.Run("disabled checker", func(t *testing.T) {
		suggestions := checker.GetSuggestions("mispeled")
		if len(suggestions) != 0 {
			t.Errorf("Disabled checker returned %d suggestions, expected 0", len(suggestions))
		}
	})

	// When enabled, check if aspell is available
	checker.SetEnabled(true)
	if !IsAvailable() {
		t.Skip("Aspell not available, skipping enabled tests")
	}

	t.Run("enabled checker with misspelled word", func(t *testing.T) {
		suggestions := checker.GetSuggestions("helo")
		// Should get at least one suggestion for "helo" (likely "hello")
		if len(suggestions) == 0 {
			t.Error("GetSuggestions(\"helo\") returned no suggestions")
		}
		// Should not return more than 5 suggestions
		if len(suggestions) > 5 {
			t.Errorf("GetSuggestions() returned %d suggestions, expected max 5", len(suggestions))
		}
	})

	t.Run("enabled checker with correct word", func(t *testing.T) {
		suggestions := checker.GetSuggestions("hello")
		// Correct words typically return no suggestions
		if len(suggestions) > 0 {
			t.Logf("GetSuggestions(\"hello\") returned %d suggestions (this is acceptable)", len(suggestions))
		}
	})

	t.Run("enabled checker with empty string", func(t *testing.T) {
		suggestions := checker.GetSuggestions("")
		if len(suggestions) != 0 {
			t.Errorf("GetSuggestions(\"\") returned %d suggestions, expected 0", len(suggestions))
		}
	})
}

func TestIsAvailable(t *testing.T) {
	available := IsAvailable()
	t.Logf("Aspell available: %v", available)

	// This test just checks that IsAvailable doesn't panic
	// The actual result depends on whether aspell is installed
}

func TestSetLanguage(t *testing.T) {
	checker := NewChecker()

	originalLang := checker.language

	checker.SetLanguage("fr")
	if checker.language != "fr" {
		t.Errorf("SetLanguage(\"fr\") failed, language = %q", checker.language)
	}

	// Restore original
	checker.SetLanguage(originalLang)
	if checker.language != originalLang {
		t.Errorf("Failed to restore language to %q, got %q", originalLang, checker.language)
	}
}
