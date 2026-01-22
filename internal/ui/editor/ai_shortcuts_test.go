package editor

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func TestEditorPane_AIUnstickShortcut(t *testing.T) {
	// Create editor pane
	ta := textarea.New()
	ta.SetValue("First line\nSecond line")
	model := NewEditorPaneModel(ta, nil)

	// Create AI unstick key message
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Alt: true, Runes: []rune{'g'}}

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel
	if !editorPane.IsContinueMode() {
		t.Error("Expected continue mode to be active after Ctrl+G")
	}

	if len(editorPane.GetContinueSuggestions()) == 0 {
		t.Error("Expected continue suggestions to be populated")
	}
}

func TestEditorPane_AISparkShortcut(t *testing.T) {
	// Create editor pane
	ta := textarea.New()
	ta.SetValue("Some content")
	model := NewEditorPaneModel(ta, nil)

	// Create AI spark key message
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Alt: true, Runes: []rune{'r'}}

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel
	if !editorPane.IsRapidBrainstorm() {
		t.Error("Expected rapid brainstorm mode to be active after Ctrl+R")
	}

	if editorPane.GetBrainstormTheme() == "" {
		t.Error("Expected brainstorm theme to be set")
	}
}

func TestEditorPane_AITweakShortcut(t *testing.T) {
	// Create editor pane with content
	ta := textarea.New()
	ta.SetValue("First line\nSecond line to tweak")
	model := NewEditorPaneModel(ta, nil)

	// Create AI tweak key message
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Alt: true, Runes: []rune{'v'}}

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel
	if !editorPane.IsVariationMode() {
		t.Error("Expected variation mode to be active after Ctrl+V")
	}

	if editorPane.GetVariationOriginal() == "" {
		t.Error("Expected variation original to be set")
	}

	if editorPane.GetVariationOriginal() != "Second line to tweak" {
		t.Errorf("Expected variation original to be 'Second line to tweak', got '%s'", editorPane.GetVariationOriginal())
	}
}

func TestEditorPane_AICheckShortcut(t *testing.T) {
	// Create editor pane with content
	ta := textarea.New()
	ta.SetValue("Line to check for quality")
	model := NewEditorPaneModel(ta, nil)

	// Create AI check key message
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Alt: true, Runes: []rune{'c'}}

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel

	// Check that quality check comment was added
	content := editorPane.GetText()
	// Debug: Print the content to see what's happening
	t.Logf("Content after AI check: %q", content)

	if !contains(content, "<!-- Quality Check:") {
		t.Errorf("Expected quality check comment to be added to content. Got: %q", content)
	}

	if !contains(content, "OKAY") {
		t.Errorf("Expected quality rating to be in content. Got: %q", content)
	}
}

func TestEditorPane_AIUnstickWithEmptyContent(t *testing.T) {
	// Create editor pane with empty content
	ta := textarea.New()
	ta.SetValue("")
	model := NewEditorPaneModel(ta, nil)

	// Create AI unstick key message
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Alt: true, Runes: []rune{'g'}}

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel
	if !editorPane.IsContinueMode() {
		t.Error("Expected continue mode to be active even with empty content")
	}
}

func TestEditorPane_AITweakWithEmptyContent(t *testing.T) {
	// Create editor pane with empty content
	ta := textarea.New()
	ta.SetValue("")
	model := NewEditorPaneModel(ta, nil)

	// Create AI tweak key message
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Alt: true, Runes: []rune{'v'}}

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel
	// Should not activate variation mode with empty content
	if editorPane.IsVariationMode() {
		t.Error("Expected variation mode to not be active with empty content")
	}
}

func TestEditorPane_AICheckWithEmptyContent(t *testing.T) {
	// Create editor pane with empty content
	ta := textarea.New()
	ta.SetValue("")
	model := NewEditorPaneModel(ta, nil)

	// Create AI check key message
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Alt: true, Runes: []rune{'c'}}

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel

	// Check that no quality check comment was added to empty content
	content := editorPane.GetText()
	if contains(content, "Quality Check:") {
		t.Error("Expected no quality check comment to be added to empty content")
	}
}

func TestEditorPane_ContinueModeSelection(t *testing.T) {
	// Create editor pane
	ta := textarea.New()
	ta.SetValue("First line")
	model := NewEditorPaneModel(ta, nil)

	// Start continue mode
	model.StartContinueMode()

	// Select first suggestion (press '1')
	keyMsg := tea.KeyMsg(tea.Key{
		Type:  tea.KeyRunes,
		Runes: []rune{'1'},
	})

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel

	// Continue mode should be cancelled after selection
	if editorPane.IsContinueMode() {
		t.Error("Expected continue mode to be cancelled after selection")
	}

	// Content should be updated
	content := editorPane.GetText()
	if !contains(content, "Continue with this line...") {
		t.Error("Expected content to be updated with selected suggestion")
	}
}

func TestEditorPane_VariationModeSelection(t *testing.T) {
	// Create editor pane
	ta := textarea.New()
	ta.SetValue("Original line")
	model := NewEditorPaneModel(ta, nil)

	// Start variation mode
	model.StartVariationMode("Original line")

	// Select first variation (press '1')
	keyMsg := tea.KeyMsg(tea.Key{
		Type:  tea.KeyRunes,
		Runes: []rune{'1'},
	})

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel

	// Variation mode should be cancelled after selection
	if editorPane.IsVariationMode() {
		t.Error("Expected variation mode to be cancelled after selection")
	}

	// Content should be updated
	content := editorPane.GetText()
	if !contains(content, "Variation 1 of: Original line") {
		t.Error("Expected content to be updated with selected variation")
	}
}

func TestEditorPane_BrainstormModeSelection(t *testing.T) {
	// Create editor pane
	ta := textarea.New()
	ta.SetValue("Some content")
	model := NewEditorPaneModel(ta, nil)

	// Start brainstorm mode
	model.StartRapidBrainstorm("love")

	// Select first angle (press '1')
	keyMsg := tea.KeyMsg(tea.Key{
		Type:  tea.KeyRunes,
		Runes: []rune{'1'},
	})

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel

	// Brainstorm mode should be cancelled after selection
	if editorPane.IsRapidBrainstorm() {
		t.Error("Expected brainstorm mode to be cancelled after selection")
	}

	// Content should be updated with opening line
	content := editorPane.GetText()
	if !contains(content, "Opening line for:") {
		t.Error("Expected content to be updated with opening line")
	}
}

func TestEditorPane_CancelAIModes(t *testing.T) {
	// Test cancelling continue mode
	ta := textarea.New()
	ta.SetValue("Some content")
	model := NewEditorPaneModel(ta, nil)

	// Start continue mode
	model.StartContinueMode()

	// Press Escape
	keyMsg := tea.KeyMsg(tea.Key{
		Type: tea.KeyEsc,
	})

	// Handle the key
	updatedModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel
	if editorPane.IsContinueMode() {
		t.Error("Expected continue mode to be cancelled after Escape")
	}

	// Test cancelling variation mode
	model.StartVariationMode("Test line")

	// Press Escape again
	updatedModel, cmd = model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane = updatedModel
	if editorPane.IsVariationMode() {
		t.Error("Expected variation mode to be cancelled after Escape")
	}

	// Test cancelling brainstorm mode
	model.StartRapidBrainstorm("test")

	// Press Escape again
	updatedModel, cmd = model.Update(keyMsg)

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane = updatedModel
	if editorPane.IsRapidBrainstorm() {
		t.Error("Expected brainstorm mode to be cancelled after Escape")
	}
}

func TestEditorPane_AIResponseTime(t *testing.T) {
	// Test that AI operations complete quickly
	ta := textarea.New()
	ta.SetValue("Test content")
	model := NewEditorPaneModel(ta, nil)

	start := time.Now()

	// Trigger AI unstick
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Alt: true, Runes: []rune{'g'}}

	updatedModel, cmd := model.Update(keyMsg)

	duration := time.Since(start)

	// Should complete quickly (under 100ms for fallback)
	if duration > 100*time.Millisecond {
		t.Errorf("AI operation took too long: %v", duration)
	}

	if cmd != nil {
		t.Error("Unexpected command returned")
	}

	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	editorPane := updatedModel
	if !editorPane.IsContinueMode() {
		t.Error("Expected continue mode to be active")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
