package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/story"
)

func TestHandleEditorKeys(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		expectedScreen Screen
	}{
		{"c navigates to characters", "c", ScreenCharacters},
		{"s navigates to scenes", "s", ScreenScenes},
		{"l navigates to locations", "l", ScreenLocations},
		{"b navigates to backups", "b", ScreenBackups},
		{"t navigates to stats", "t", ScreenStats},
		{"e navigates to export", "e", ScreenExport},
		{"h navigates to help", "h", ScreenHelp},
		{"esc returns to welcome", "esc", ScreenWelcome},
		{"unknown key stays on editor", "x", ScreenEditor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			m.CurrentScreen = ScreenEditor
			m.CurrentProject = &story.Project{
				Title:  "Test Project",
				Author: "Test Author",
				Genre:  "Test Genre",
			}

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			result, _ := m.handleEditorKeys(msg)
			resultModel := result.(Model)

			if resultModel.CurrentScreen != tt.expectedScreen {
				t.Errorf("CurrentScreen = %v, expected %v", resultModel.CurrentScreen, tt.expectedScreen)
			}
		})
	}
}

func TestHandleEditorKeys_HelpSetsPreviousScreen(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenEditor
	m.CurrentProject = &story.Project{Title: "Test"}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")}
	result, _ := m.handleEditorKeys(msg)
	resultModel := result.(Model)

	if resultModel.PreviousScreen != ScreenEditor {
		t.Errorf("PreviousScreen = %v, expected ScreenEditor", resultModel.PreviousScreen)
	}
}

func TestHandleEditorKeys_BackupsResetState(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenEditor
	m.CurrentProject = &story.Project{Title: "Test"}
	m.SelectedIndex = 5
	m.Message = "Test message"
	m.Error = &testError{"test error"}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("b")}
	result, _ := m.handleEditorKeys(msg)
	resultModel := result.(Model)

	if resultModel.SelectedIndex != 0 {
		t.Errorf("SelectedIndex = %d, expected 0", resultModel.SelectedIndex)
	}
	if resultModel.Message != "" {
		t.Errorf("Message = %q, expected empty", resultModel.Message)
	}
	if resultModel.Error != nil {
		t.Error("Error should be nil after navigating to backups")
	}
}

func TestHandleEditorKeys_ExportResetState(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenEditor
	m.CurrentProject = &story.Project{Title: "Test"}
	m.SelectedIndex = 5
	m.Message = "Test message"
	m.Error = &testError{"test error"}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")}
	result, _ := m.handleEditorKeys(msg)
	resultModel := result.(Model)

	if resultModel.SelectedIndex != 0 {
		t.Errorf("SelectedIndex = %d, expected 0", resultModel.SelectedIndex)
	}
	if resultModel.Message != "" {
		t.Errorf("Message = %q, expected empty", resultModel.Message)
	}
	if resultModel.Error != nil {
		t.Error("Error should be nil after navigating to export")
	}
}

func TestViewEditor_NoProject(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenEditor

	view := m.viewEditor()

	if view != "No project loaded" {
		t.Errorf("viewEditor() = %q, expected 'No project loaded'", view)
	}
}

func TestViewEditor_WithProject(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenEditor
	m.Width = 80
	m.Height = 24
	m.CurrentProject = &story.Project{
		Title:           "My Novel",
		Author:          "John Doe",
		Genre:           "Fantasy",
		Status:          "draft",
		TotalWords:      5000,
		TotalScenes:     10,
		TotalCharacters: 5,
	}

	view := m.viewEditor()

	expectedContent := []string{
		"My Novel",
		"John Doe",
		"Fantasy",
		"draft",
		"5000",
		"10",
		"5",
		"Characters",
		"Scenes",
		"Locations",
		"Statistics",
		"Export",
	}

	for _, content := range expectedContent {
		if !contains(view, content) {
			t.Errorf("viewEditor should contain %q", content)
		}
	}
}

func TestViewEditor_WithMessage(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenEditor
	m.Width = 80
	m.Height = 24
	m.CurrentProject = &story.Project{
		Title: "Test Project",
	}
	m.Message = "Operation successful"

	view := m.viewEditor()

	if !contains(view, "Operation successful") {
		t.Error("viewEditor should include message")
	}
}
