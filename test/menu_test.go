package noise

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/puente-labs/noise/internal/ui"
)

// TestMenuModelCreation tests menu model creation
func TestMenuModelCreation(t *testing.T) {
	model := ui.NewMenuModel()

	if model == nil {
		t.Fatal("Expected menu model to be created, got nil")
	}

	// Test initial state
	// Since selected and focused are private, we'll test through behavior
	view := model.View()
	if view == "" {
		t.Error("Expected view to render")
	}

	// Verify view contains expected content
	if !contains(view, "noise.sh") {
		t.Error("Expected view to contain noise.sh title")
	}
}

// TestMenuModelWindowSize tests window size handling
func TestMenuModelWindowSize(t *testing.T) {
	model := ui.NewMenuModel()

	// Test window size message
	width, height := 100, 30
	msg := tea.WindowSizeMsg{Width: width, Height: height}
	updatedModel, _ := model.Update(msg)

	if updatedModel == nil {
		t.Error("Expected model to be updated")
	}

	// Verify view renders with new dimensions
	view := updatedModel.View()
	if view == "" {
		t.Error("Expected view to render with dimensions set")
	}
}

// TestMenuModelKeyHandling tests key handling
func TestMenuModelKeyHandling(t *testing.T) {
	model := ui.NewMenuModel()

	// Set dimensions
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Test up key
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := model.Update(upMsg)

	if updatedModel == nil {
		t.Error("Expected model to be updated with up key")
	}

	view := updatedModel.View()
	if view == "" {
		t.Error("Expected view to render after up key")
	}

	// Test down key
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = model.Update(downMsg)

	if updatedModel == nil {
		t.Error("Expected model to be updated with down key")
	}

	view = updatedModel.View()
	if view == "" {
		t.Error("Expected view to render after down key")
	}

	// Test enter key
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = model.Update(enterMsg)

	if updatedModel == nil {
		t.Error("Expected model to be updated with enter key")
	}

	view = updatedModel.View()
	if view == "" {
		t.Error("Expected view to render after enter key")
	}
}

// TestMenuModelResponsiveMode tests responsive mode behavior
func TestMenuModelResponsiveMode(t *testing.T) {
	model := ui.NewMenuModel()

	// Test full mode (large terminal)
	model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	view := model.View()

	if view == "" {
		t.Error("Expected view to render in full mode")
	}

	// Verify full mode content
	if !contains(view, "🎵 noise.sh") {
		t.Error("Expected full mode title")
	}

	// Test compact mode (medium terminal)
	model.Update(tea.WindowSizeMsg{Width: 85, Height: 25})
	view = model.View()

	if view == "" {
		t.Error("Expected view to render in compact mode")
	}

	// Test minimal mode (small terminal)
	model.Update(tea.WindowSizeMsg{Width: 65, Height: 20})
	view = model.View()

	if view == "" {
		t.Error("Expected view to render in minimal mode")
	}

	// Verify minimal mode content
	if !contains(view, "🎵 LF") {
		t.Error("Expected minimal mode title")
	}
}

// TestMenuModelAnimation tests animation behavior
func TestMenuModelAnimation(t *testing.T) {
	model := ui.NewMenuModel()

	// Set dimensions
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Test selection animation (enter key)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(enterMsg)

	if updatedModel == nil {
		t.Error("Expected model to be updated with enter key")
	}

	// Test selection change animation (down key)
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = model.Update(downMsg)

	if updatedModel == nil {
		t.Error("Expected model to be updated with down key")
	}

	// Test animation tick message
	tickMsg := ui.AnimationTickMsg{}
	updatedModel, _ = model.Update(tickMsg)

	if updatedModel == nil {
		t.Error("Expected model to be updated with animation tick")
	}
}

// TestMenuModelMenuItems tests menu items
func TestMenuModelMenuItems(t *testing.T) {
	model := ui.NewMenuModel()

	// Set dimensions
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Get view
	view := model.View()
	if view == "" {
		t.Error("Expected view to render")
	}

	// Verify expected menu items are present
	expectedItems := []string{
		"New Song",
		"Open Song",
		"Theory Tools",
		"Audio Tools",
		"Project Manager",
		"Settings",
		"Exit",
	}

	for _, item := range expectedItems {
		if !contains(view, item) {
			t.Errorf("Expected view to contain menu item: %s", item)
		}
	}
}

// TestMenuModelDescriptions tests menu item descriptions
func TestMenuModelDescriptions(t *testing.T) {
	model := ui.NewMenuModel()

	// Set dimensions
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Get view
	view := model.View()
	if view == "" {
		t.Error("Expected view to render")
	}

	// Verify expected descriptions are present
	expectedDescs := []string{
		"Create a new song",
		"Open an existing song",
		"Music theory and rhyme tools",
		"Metronome and chord playback",
		"Manage songs and projects",
		"Application settings",
		"Exit noise.sh",
	}

	for _, desc := range expectedDescs {
		if !contains(view, desc) {
			t.Errorf("Expected view to contain description: %s", desc)
		}
	}
}

// TestMenuModelNavigation tests menu navigation
func TestMenuModelNavigation(t *testing.T) {
	model := ui.NewMenuModel()

	// Set dimensions
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Test navigation down
	for i := 0; i < 5; i++ {
		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		updatedModel, _ := model.Update(downMsg)

		if updatedModel == nil {
			t.Errorf("Expected model to be updated with down key on iteration %d", i)
		}

		view := updatedModel.View()
		if view == "" {
			t.Errorf("Expected view to render after down key on iteration %d", i)
		}
	}

	// Test navigation up
	for i := 0; i < 5; i++ {
		upMsg := tea.KeyMsg{Type: tea.KeyUp}
		updatedModel, _ := model.Update(upMsg)

		if updatedModel == nil {
			t.Errorf("Expected model to be updated with up key on iteration %d", i)
		}

		view := updatedModel.View()
		if view == "" {
			t.Errorf("Expected view to render after up key on iteration %d", i)
		}
	}
}

// TestMenuModelEdgeCases tests edge cases and boundary conditions
func TestMenuModelEdgeCases(t *testing.T) {
	model := ui.NewMenuModel()

	// Test with zero dimensions
	model.Update(tea.WindowSizeMsg{Width: 0, Height: 0})
	view := model.View()
	if view == "" {
		t.Error("Expected view to render with zero dimensions")
	}

	// Test with very small dimensions
	model.Update(tea.WindowSizeMsg{Width: 10, Height: 5})
	view = model.View()
	if view == "" {
		t.Error("Expected view to render with very small dimensions")
	}

	// Test with very large dimensions
	model.Update(tea.WindowSizeMsg{Width: 500, Height: 200})
	view = model.View()
	if view == "" {
		t.Error("Expected view to render with very large dimensions")
	}

	// Test with unrecognized key
	unrecognizedMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	updatedModel, _ := model.Update(unrecognizedMsg)

	if updatedModel == nil {
		t.Error("Expected model to handle unrecognized key")
	}

	view = updatedModel.View()
	if view == "" {
		t.Error("Expected view to render after unrecognized key")
	}
}

// TestMenuModelConsistency tests model consistency
func TestMenuModelConsistency(t *testing.T) {
	model := ui.NewMenuModel()

	// Set dimensions
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Test multiple updates
	for i := 0; i < 10; i++ {
		// Navigate
		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		model, _ = model.Update(downMsg)

		// Resize
		resizeMsg := tea.WindowSizeMsg{Width: 100 + i, Height: 30 + i}
		model, _ = model.Update(resizeMsg)

		// Verify view still renders
		view := model.View()
		if view == "" {
			t.Errorf("Expected view to render on iteration %d", i)
		}
	}
}

// TestMenuModelPerformance tests performance with many updates
func TestMenuModelPerformance(t *testing.T) {
	model := ui.NewMenuModel()

	// Set dimensions
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Test many key updates
	keys := []tea.KeyMsg{
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
	}

	for i := 0; i < 100; i++ {
		key := keys[i%len(keys)]
		updatedModel, _ := model.Update(key)

		if updatedModel == nil {
			t.Errorf("Expected model to be updated on iteration %d", i)
		}

		// Don't check view on every iteration for performance
		if i%10 == 0 {
			view := updatedModel.View()
			if view == "" {
				t.Errorf("Expected view to render on iteration %d", i)
			}
		}
	}
}

// TestMenuModelResponsiveModeTransitions tests responsive mode transitions
func TestMenuModelResponsiveModeTransitions(t *testing.T) {
	model := ui.NewMenuModel()

	// Test transitions between modes
	dimensions := []struct {
		width  int
		height int
	}{
		{120, 40}, // Full
		{95, 30},  // Full
		{85, 25},  // Compact
		{75, 20},  // Compact
		{65, 20},  // Minimal
		{55, 15},  // Minimal
	}

	for _, dim := range dimensions {
		model.Update(tea.WindowSizeMsg{Width: dim.width, Height: dim.height})
		view := model.View()

		if view == "" {
			t.Errorf("Expected view to render with dimensions %dx%d", dim.width, dim.height)
		}
	}
}

// TestMenuModelStyling tests menu styling
func TestMenuModelStyling(t *testing.T) {
	model := ui.NewMenuModel()

	// Set dimensions
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Get view
	view := model.View()
	if view == "" {
		t.Error("Expected view to render")
	}

	// Verify view contains styled elements
	if !contains(view, "🎵") {
		t.Error("Expected view to contain styled elements")
	}
}

// BenchmarkMenuModelViewRendering benchmarks view rendering performance
func BenchmarkMenuModelViewRendering(b *testing.B) {
	model := ui.NewMenuModel()
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

// BenchmarkMenuModelKeyHandling benchmarks key handling performance
func BenchmarkMenuModelKeyHandling(b *testing.B) {
	model := ui.NewMenuModel()
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Test key messages
	keys := []tea.KeyMsg{
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
		{Type: tea.KeyEnter},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := keys[i%len(keys)]
		_, _ = model.Update(key)
	}
}

// BenchmarkMenuModelResponsiveModeChanges benchmarks responsive mode changes
func BenchmarkMenuModelResponsiveModeChanges(b *testing.B) {
	model := ui.NewMenuModel()

	// Test various dimensions
	dimensions := []struct {
		width  int
		height int
	}{
		{120, 40}, // Full
		{85, 25},  // Compact
		{65, 20},  // Minimal
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dim := dimensions[i%len(dimensions)]
		model.Update(tea.WindowSizeMsg{Width: dim.width, Height: dim.height})
		_ = model.View()
	}
}
