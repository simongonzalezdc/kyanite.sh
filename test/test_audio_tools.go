package lyricforge_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/puente-labs/lyricforge/internal/ui"
)

// TestAudioModelCreation tests that the audio model can be created successfully
func TestAudioModelCreation(t *testing.T) {
	audioModel := ui.NewAudioModel()
	if audioModel == nil {
		t.Fatal("NewAudioModel() returned nil")
	}

	// Test initial state
	if audioModel.Init() != nil {
		t.Error("AudioModel.Init() should return nil initially")
	}
}

// TestAudioModelUpdate tests that the audio model handles messages correctly
func TestAudioModelUpdate(t *testing.T) {
	audioModel := ui.NewAudioModel()

	// Test window size message
	width, height := 120, 40
	msg := tea.WindowSizeMsg{Width: width, Height: height}

	updatedModel, cmd := audioModel.Update(msg)
	if updatedModel == nil {
		t.Fatal("Update returned nil model")
	}

	if cmd != nil {
		t.Error("Update should not return command for window size message")
	}
}

// TestAudioModelView tests that the audio model renders correctly
func TestAudioModelView(t *testing.T) {
	audioModel := ui.NewAudioModel()

	// Set dimensions first
	width, height := 100, 30
	msg := tea.WindowSizeMsg{Width: width, Height: height}
	audioModel.Update(msg)

	view := audioModel.View()
	if view == "" {
		t.Error("AudioModel.View() returned empty string")
	}

	// Check that view contains expected content
	expectedContent := "🎵 Audio Tools"
	if !strings.Contains(view, expectedContent) {
		t.Errorf("View should contain '%s', got: %s", expectedContent, view)
	}
}

// TestAudioToolsNavigation tests tool switching functionality
func TestAudioToolsNavigation(t *testing.T) {
	audioModel := ui.NewAudioModel()

	// Test tab navigation
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	audioModel.Update(tabMsg)

	// Test direct tool selection
	tool1Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	audioModel.Update(tool1Msg)

	// Test spacebar for play/pause
	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	audioModel.Update(spaceMsg)

	// Test stop
	stopMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	audioModel.Update(stopMsg)
}

// TestMetronomeFunctionality tests metronome-specific features
func TestMetronomeFunctionality(t *testing.T) {
	audioModel := ui.NewAudioModel()

	// Test tempo adjustment
	plusMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}}
	audioModel.Update(plusMsg)

	minusMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}}
	audioModel.Update(minusMsg)

	// Test time signature toggle
	timeSigMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
	audioModel.Update(timeSigMsg)
}

// TestHelpSystem tests the help functionality
func TestHelpSystem(t *testing.T) {
	audioModel := ui.NewAudioModel()

	// Test help toggle
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	audioModel.Update(helpMsg)

	// Test escape to exit help
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	audioModel.Update(escMsg)
}

// BenchmarkAudioModelView benchmarks the view rendering performance
func BenchmarkAudioModelView(b *testing.B) {
	audioModel := ui.NewAudioModel()

	// Set dimensions
	width, height := 120, 40
	msg := tea.WindowSizeMsg{Width: width, Height: height}
	audioModel.Update(msg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = audioModel.View()
	}
}

// TestAudioToolsIntegration tests that audio tools integrate properly with the existing system
func TestAudioToolsIntegration(t *testing.T) {
	// Test that audio model can be created through the same pattern as other models
	audioModel := ui.NewAudioModel()

	// Test that it handles the same message types as other UI components
	width, height := 100, 30
	windowMsg := tea.WindowSizeMsg{Width: width, Height: height}
	_, cmd := audioModel.Update(windowMsg)
	if cmd != nil {
		t.Error("Window size message should not produce commands")
	}

	// Test keyboard navigation
	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	_, cmd = audioModel.Update(keyMsg)
	if cmd != nil {
		t.Error("Tab key should not produce commands")
	}
}
