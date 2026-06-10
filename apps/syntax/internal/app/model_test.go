package app

import (
	"testing"
	"time"

	"github.com/kyanite/syntax/internal/editor"
	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/story"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.CurrentScreen != ScreenWelcome {
		t.Errorf("NewModel() CurrentScreen = %v, expected ScreenWelcome", m.CurrentScreen)
	}

	if m.SaveStatus != SaveStatusSaved {
		t.Errorf("NewModel() SaveStatus = %v, expected SaveStatusSaved", m.SaveStatus)
	}

	if m.SpellChecker == nil {
		t.Error("NewModel() SpellChecker is nil")
	}

	if m.ThemeManager == nil {
		t.Error("NewModel() ThemeManager is nil")
	}
}

func TestHasProject(t *testing.T) {
	tests := []struct {
		name     string
		project  *story.Project
		expected bool
	}{
		{"no project", nil, false},
		{"with project", &story.Project{Title: "Test"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{CurrentProject: tt.project}
			result := m.HasProject()
			if result != tt.expected {
				t.Errorf("HasProject() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHasUnsavedChanges(t *testing.T) {
	tests := []struct {
		name     string
		buffer   *editor.Buffer
		modified bool
		expected bool
	}{
		{"no buffer", nil, false, false},
		{"buffer not modified", editor.NewBuffer("test"), false, false},
		{"buffer modified", editor.NewBuffer("test"), true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{Buffer: tt.buffer}
			if tt.buffer != nil && tt.modified {
				tt.buffer.SetModified(true)
			}
			result := m.HasUnsavedChanges()
			if result != tt.expected {
				t.Errorf("HasUnsavedChanges() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsEditorActive(t *testing.T) {
	tests := []struct {
		name    string
		screen  Screen
		buffer  *editor.Buffer
		scene   *scene.Scene
		expected bool
	}{
		{"not in editor screen", ScreenWelcome, editor.NewBuffer(""), &scene.Scene{}, false},
		{"in editor but no buffer", ScreenTextEditor, nil, &scene.Scene{}, false},
		{"in editor but no scene", ScreenTextEditor, editor.NewBuffer(""), nil, false},
		{"editor fully active", ScreenTextEditor, editor.NewBuffer(""), &scene.Scene{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				CurrentScreen: tt.screen,
				Buffer:        tt.buffer,
				CurrentScene:  tt.scene,
			}
			result := m.IsEditorActive()
			if result != tt.expected {
				t.Errorf("IsEditorActive() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCanAutoSave(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		screen       Screen
		buffer       *editor.Buffer
		scene        *scene.Scene
		project      *story.Project
		modified     bool
		lastEditTime time.Time
		expected     bool
	}{
		{
			name:         "all conditions met",
			screen:       ScreenTextEditor,
			buffer:       editor.NewBuffer("test"),
			scene:        &scene.Scene{},
			project:      &story.Project{},
			modified:     true,
			lastEditTime: now.Add(-5 * time.Second), // 5 seconds ago
			expected:     true,
		},
		{
			name:         "not in editor",
			screen:       ScreenWelcome,
			buffer:       editor.NewBuffer("test"),
			scene:        &scene.Scene{},
			project:      &story.Project{},
			modified:     true,
			lastEditTime: now.Add(-5 * time.Second),
			expected:     false,
		},
		{
			name:         "no buffer",
			screen:       ScreenTextEditor,
			buffer:       nil,
			scene:        &scene.Scene{},
			project:      &story.Project{},
			modified:     false,
			lastEditTime: now.Add(-5 * time.Second),
			expected:     false,
		},
		{
			name:         "not modified",
			screen:       ScreenTextEditor,
			buffer:       editor.NewBuffer("test"),
			scene:        &scene.Scene{},
			project:      &story.Project{},
			modified:     false,
			lastEditTime: now.Add(-5 * time.Second),
			expected:     false,
		},
		{
			name:         "edited too recently",
			screen:       ScreenTextEditor,
			buffer:       editor.NewBuffer("test"),
			scene:        &scene.Scene{},
			project:      &story.Project{},
			modified:     true,
			lastEditTime: now.Add(-1 * time.Second), // Only 1 second ago
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				CurrentScreen: tt.screen,
				Buffer:        tt.buffer,
				CurrentScene:  tt.scene,
				CurrentProject: tt.project,
				LastEditTime:  tt.lastEditTime,
			}
			if tt.buffer != nil && tt.modified {
				tt.buffer.SetModified(true)
			}
			result := m.CanAutoSave()
			if result != tt.expected {
				t.Errorf("CanAutoSave() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestResetInputState(t *testing.T) {
	m := Model{
		InputMode:    true,
		InputValue:   "test",
		ReplaceValue: "replace",
	}

	m.ResetInputState()

	if m.InputMode {
		t.Error("ResetInputState() did not set InputMode to false")
	}
	if m.InputValue != "" {
		t.Errorf("ResetInputState() InputValue = %q, expected empty string", m.InputValue)
	}
	if m.ReplaceValue != "" {
		t.Errorf("ResetInputState() ReplaceValue = %q, expected empty string", m.ReplaceValue)
	}
}

func TestGetContentDimensions(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		height         int
		expectedWidth  int
		expectedHeight int
	}{
		{"standard size", 80, 24, 80, 20},
		{"large size", 120, 40, 120, 36},
		{"small size", 40, 10, 40, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				Width:  tt.width,
				Height: tt.height,
			}
			w, h := m.GetContentDimensions()
			if w != tt.expectedWidth {
				t.Errorf("GetContentDimensions() width = %d, expected %d", w, tt.expectedWidth)
			}
			if h != tt.expectedHeight {
				t.Errorf("GetContentDimensions() height = %d, expected %d", h, tt.expectedHeight)
			}
		})
	}
}

func TestSetMessage(t *testing.T) {
	m := Model{}
	msg := "Test message"

	m.SetMessage(msg)

	if m.Message != msg {
		t.Errorf("SetMessage() Message = %q, expected %q", m.Message, msg)
	}
}

func TestSetError(t *testing.T) {
	m := Model{Message: "Some message"}
	err := &testError{"test error"}

	m.SetError(err)

	if m.Error != err {
		t.Errorf("SetError() Error = %v, expected %v", m.Error, err)
	}
	if m.Message != "" {
		t.Errorf("SetError() should clear Message, got %q", m.Message)
	}

	// Test with nil error
	m.Message = "Another message"
	m.SetError(nil)
	if m.Error != nil {
		t.Errorf("SetError(nil) Error = %v, expected nil", m.Error)
	}
	if m.Message != "Another message" {
		t.Error("SetError(nil) should not clear Message")
	}
}

func TestClearFeedback(t *testing.T) {
	m := Model{
		Message: "Test message",
		Error:   &testError{"test error"},
	}

	m.ClearFeedback()

	if m.Message != "" {
		t.Errorf("ClearFeedback() Message = %q, expected empty string", m.Message)
	}
	if m.Error != nil {
		t.Errorf("ClearFeedback() Error = %v, expected nil", m.Error)
	}
}

func TestExitEditor(t *testing.T) {
	tests := []struct {
		name           string
		saveIfModified bool
		modified       bool
	}{
		{"save on exit with modifications", true, true},
		{"save on exit without modifications", true, false},
		{"don't save on exit with modifications", false, true},
		{"don't save on exit without modifications", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := editor.NewBuffer("test content")
			if tt.modified {
				buffer.SetModified(true)
			}

			m := Model{
				CurrentScreen: ScreenTextEditor,
				CurrentScene:  &scene.Scene{ID: "scene1"},
				Buffer:        buffer,
			}

			m.ExitEditor(tt.saveIfModified)

			if m.CurrentScene != nil {
				t.Error("ExitEditor() should clear CurrentScene")
			}
			if m.Buffer != nil {
				t.Error("ExitEditor() should clear Buffer")
			}
			if m.CurrentScreen != ScreenScenes {
				t.Errorf("ExitEditor() CurrentScreen = %v, expected ScreenScenes", m.CurrentScreen)
			}
		})
	}
}

// Helper type for testing errors
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestModelInitialization(t *testing.T) {
	// Test that all helper types are properly initialized
	m := NewModel()

	// Check theme-related initialization
	if m.ThemeManager == nil {
		t.Error("ThemeManager not initialized")
	}
	if m.CurrentTheme.Name == "" {
		t.Error("CurrentTheme not initialized")
	}

	// Check spell checker initialization
	if m.SpellChecker == nil {
		t.Error("SpellChecker not initialized")
	}
	if m.SpellChecker.IsEnabled() {
		t.Error("SpellChecker should be disabled by default")
	}

	// Check time fields
	if m.LastSaveTime.IsZero() {
		t.Error("LastSaveTime should be initialized to current time")
	}
}

func TestModelHelperMethodsIntegration(t *testing.T) {
	// Integration test: Create a realistic scenario
	m := NewModel()

	// Initially should have no project
	if m.HasProject() {
		t.Error("New model should not have a project")
	}

	// Add a project
	m.CurrentProject = &story.Project{
		Title:     "Test Project",
		Directory: "/tmp/test",
	}

	if !m.HasProject() {
		t.Error("Model should have a project after setting CurrentProject")
	}

	// Initially should not have unsaved changes
	if m.HasUnsavedChanges() {
		t.Error("New model should not have unsaved changes")
	}

	// Add a buffer and scene (but don't modify)
	m.Buffer = editor.NewBuffer("test content")
	m.CurrentScene = &scene.Scene{ID: "scene1", Name: "Opening"}
	m.CurrentScreen = ScreenTextEditor

	if m.HasUnsavedChanges() {
		t.Error("Should not have unsaved changes immediately after creating buffer")
	}

	// Editor should be active now
	if !m.IsEditorActive() {
		t.Error("Editor should be active with buffer, scene, and correct screen")
	}

	// Modify the buffer
	m.Buffer.SetModified(true)
	if !m.HasUnsavedChanges() {
		t.Error("Should have unsaved changes after modifying buffer")
	}

	// Should not be able to auto-save immediately
	m.LastEditTime = time.Now()
	if m.CanAutoSave() {
		t.Error("Should not be able to auto-save immediately after edit")
	}

	// Should be able to auto-save after waiting
	m.LastEditTime = time.Now().Add(-5 * time.Second)
	if !m.CanAutoSave() {
		t.Error("Should be able to auto-save 5 seconds after edit")
	}
}
