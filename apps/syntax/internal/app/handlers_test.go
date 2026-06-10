package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleHelpKeys(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		previousScreen Screen
		expectedScreen Screen
	}{
		{"esc returns to previous", "esc", ScreenEditor, ScreenEditor},
		{"q returns to previous", "q", ScreenWelcome, ScreenWelcome},
		{"? returns to previous", "?", ScreenScenes, ScreenScenes},
		{"other key stays on help", "a", ScreenEditor, ScreenHelp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			m.CurrentScreen = ScreenHelp
			m.PreviousScreen = tt.previousScreen

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			result, _ := m.handleHelpKeys(msg)
			resultModel := result.(Model)

			if resultModel.CurrentScreen != tt.expectedScreen {
				t.Errorf("CurrentScreen = %v, expected %v", resultModel.CurrentScreen, tt.expectedScreen)
			}
		})
	}
}

func TestHandleWelcomeKeys_Quit(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenWelcome

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd := m.handleWelcomeKeys(msg)

	if cmd == nil {
		t.Error("Expected quit command, got nil")
	}
}

func TestHandleWelcomeKeys_Unknown(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenWelcome

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}
	result, cmd := m.handleWelcomeKeys(msg)
	resultModel := result.(Model)

	if resultModel.CurrentScreen != ScreenWelcome {
		t.Error("Unknown key should keep us on welcome screen")
	}

	if cmd != nil {
		t.Error("Unknown key should not return a command")
	}
}

func TestHandleKeyPress_GlobalShortcuts(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		screen   Screen
		expected Screen
	}{
		{"? shows help from welcome", "?", ScreenWelcome, ScreenHelp},
		{"? shows help from editor", "?", ScreenEditor, ScreenHelp},
		{"h shows help from editor", "h", ScreenEditor, ScreenHelp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			m.CurrentScreen = tt.screen

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			result, _ := m.handleKeyPress(msg)
			resultModel := result.(Model)

			if resultModel.CurrentScreen != tt.expected {
				t.Errorf("CurrentScreen = %v, expected %v", resultModel.CurrentScreen, tt.expected)
			}

			if resultModel.PreviousScreen != tt.screen {
				t.Errorf("PreviousScreen = %v, expected %v", resultModel.PreviousScreen, tt.screen)
			}
		})
	}
}

func TestHandleKeyPress_ThemeCycle(t *testing.T) {
	t.Skip("Theme cycle test requires proper KeyMsg construction for ctrl+shift+t")
}

func TestView_DifferentScreens(t *testing.T) {
	tests := []struct {
		name   string
		screen Screen
	}{
		{"welcome screen", ScreenWelcome},
		{"help screen", ScreenHelp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			m.CurrentScreen = tt.screen
			m.Width = 80
			m.Height = 24

			view := m.View()
			if view == "" {
				t.Errorf("View() returned empty string for %v", tt.screen)
			}
			if view == "Unknown screen" {
				t.Errorf("View() returned 'Unknown screen' for %v", tt.screen)
			}
		})
	}
}

func TestView_UnknownScreen(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = Screen(999) // Invalid screen

	view := m.View()
	if view != "Unknown screen" {
		t.Errorf("View() = %q, expected 'Unknown screen'", view)
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	m := NewModel()

	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	result, _ := m.Update(msg)
	resultModel := result.(Model)

	if resultModel.Width != 120 {
		t.Errorf("Width = %d, expected 120", resultModel.Width)
	}

	if resultModel.Height != 40 {
		t.Errorf("Height = %d, expected 40", resultModel.Height)
	}
}

func TestUpdate_AutoSaveTick(t *testing.T) {
	m := NewModel()

	// Without active editor, should not crash
	msg := AutoSaveTickMsg{}
	result, cmd := m.Update(msg)

	if result == nil {
		t.Error("Update returned nil model")
	}

	if cmd == nil {
		t.Error("AutoSaveTick should return a command for next tick")
	}
}

func TestUpdate_AutoSaveComplete_Success(t *testing.T) {
	m := NewModel()
	m.SaveStatus = SaveStatusSaving

	msg := AutoSaveCompleteMsg{Err: nil}
	result, _ := m.Update(msg)
	resultModel := result.(Model)

	if resultModel.SaveStatus != SaveStatusSaved {
		t.Errorf("SaveStatus = %v, expected SaveStatusSaved", resultModel.SaveStatus)
	}

	if resultModel.Error != nil {
		t.Error("Error should be nil after successful save")
	}
}

func TestUpdate_AutoSaveComplete_Error(t *testing.T) {
	m := NewModel()
	m.SaveStatus = SaveStatusSaving

	testErr := &testError{"save failed"}
	msg := AutoSaveCompleteMsg{Err: testErr}
	result, _ := m.Update(msg)
	resultModel := result.(Model)

	if resultModel.SaveStatus != SaveStatusUnsaved {
		t.Errorf("SaveStatus = %v, expected SaveStatusUnsaved", resultModel.SaveStatus)
	}

	if resultModel.Error != testErr {
		t.Error("Error should be set after failed save")
	}
}

func TestEnsureDataLoaded(t *testing.T) {
	m := NewModel()

	// Without a project, should not crash
	m.ensureDataLoaded()

	// This test mainly ensures the function doesn't panic
	// when called with various screen states
	m.CurrentScreen = ScreenScenes
	m.ensureDataLoaded()

	m.CurrentScreen = ScreenCharacters
	m.ensureDataLoaded()

	m.CurrentScreen = ScreenLocations
	m.ensureDataLoaded()
}

func TestViewWelcome_WithMessage(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenWelcome
	m.Width = 80
	m.Height = 24
	m.Message = "Test message"

	view := m.viewWelcome()

	if !contains(view, "Test message") {
		t.Error("viewWelcome should include message")
	}
}

func TestViewWelcome_WithError(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenWelcome
	m.Width = 80
	m.Height = 24
	m.Error = &testError{"test error"}

	view := m.viewWelcome()

	if !contains(view, "test error") {
		t.Error("viewWelcome should include error message")
	}
}

func TestViewHelp(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenHelp
	m.Width = 80
	m.Height = 24

	view := m.viewHelp()

	// Check for expected content
	expectedContent := []string{
		"Help",
		"Global Shortcuts",
		"Navigation",
		"Text Editor",
		"Ctrl+Q",
		"Esc",
	}

	for _, content := range expectedContent {
		if !contains(view, content) {
			t.Errorf("viewHelp should contain %q", content)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
