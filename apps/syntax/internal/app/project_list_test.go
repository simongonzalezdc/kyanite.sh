package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/story"
)

func TestHandleProjectListKeys_Navigation(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		initialIndex  int
		numProjects   int
		expectedIndex int
	}{
		{"up from first stays at first", "up", 0, 3, 0},
		{"up decrements index", "up", 2, 3, 1},
		{"down increments index", "down", 0, 3, 1},
		{"down from last stays at last", "down", 2, 3, 2},
		{"k moves up", "k", 2, 3, 1},
		{"j moves down", "j", 1, 3, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			m.CurrentScreen = ScreenProjectList
			m.SelectedIndex = tt.initialIndex
			m.Projects = make([]*story.Project, tt.numProjects)
			for i := 0; i < tt.numProjects; i++ {
				m.Projects[i] = &story.Project{Title: "Project " + string(rune('A'+i))}
			}

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			result, _ := m.handleProjectListKeys(msg)
			resultModel := result.(Model)

			if resultModel.SelectedIndex != tt.expectedIndex {
				t.Errorf("SelectedIndex = %d, expected %d", resultModel.SelectedIndex, tt.expectedIndex)
			}
		})
	}
}

func TestHandleProjectListKeys_Escape(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenProjectList

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")}
	result, _ := m.handleProjectListKeys(msg)
	resultModel := result.(Model)

	if resultModel.CurrentScreen != ScreenWelcome {
		t.Errorf("CurrentScreen = %v, expected ScreenWelcome", resultModel.CurrentScreen)
	}
}

func TestHandleProjectListKeys_EnterWithNoProjects(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenProjectList
	m.Projects = []*story.Project{}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("enter")}
	result, _ := m.handleProjectListKeys(msg)
	resultModel := result.(Model)

	if resultModel.CurrentScreen != ScreenProjectList {
		t.Error("Should stay on project list when no projects exist")
	}
}

func TestHandleProjectListKeys_UnknownKey(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenProjectList
	m.SelectedIndex = 0

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}
	result, _ := m.handleProjectListKeys(msg)
	resultModel := result.(Model)

	if resultModel.CurrentScreen != ScreenProjectList {
		t.Error("Unknown key should keep us on project list")
	}
	if resultModel.SelectedIndex != 0 {
		t.Error("Unknown key should not change selection")
	}
}

func TestViewProjectList_Empty(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenProjectList
	m.Width = 80
	m.Height = 24
	m.Projects = []*story.Project{}

	view := m.viewProjectList()

	expectedContent := []string{
		"Select a Project",
		"No projects found",
		"Press 'n' to create one",
	}

	for _, content := range expectedContent {
		if !contains(view, content) {
			t.Errorf("viewProjectList should contain %q", content)
		}
	}
}

func TestViewProjectList_WithProjects(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenProjectList
	m.Width = 80
	m.Height = 24
	m.SelectedIndex = 1
	m.Projects = []*story.Project{
		{Title: "First Novel"},
		{Title: "Second Novel"},
		{Title: "Third Novel"},
	}

	view := m.viewProjectList()

	expectedContent := []string{
		"Select a Project",
		"First Novel",
		"Second Novel",
		"Third Novel",
		"Navigate",
		"Enter",
		"Esc",
	}

	for _, content := range expectedContent {
		if !contains(view, content) {
			t.Errorf("viewProjectList should contain %q", content)
		}
	}
}
