package noise

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Kyanite/noise/internal/ui/editor"
)

// TestHelpPaneModelCreation tests help pane model creation
func TestHelpPaneModelCreation(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	if model == nil {
		t.Fatal("Expected help pane model to be created, got nil")
	}

	// Test initial state
	if model.GetShortcutManager() != shortcutManager {
		t.Error("Expected shortcut manager to be set")
	}

	// Test initial dimensions (can't directly access private fields)
	// We'll verify through view rendering behavior

	// Test initial focus state
	// Since focused is private, we'll test through behavior
	view := model.View()
	if view == "" {
		t.Error("Expected view to render even with zero dimensions")
	}
}

// TestHelpPaneDimensions tests dimension setting
func TestHelpPaneDimensions(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Test setting dimensions
	model.SetDimensions(100, 30)

	// Verify dimensions were set (through view rendering)
	view := model.View()
	if view == "" {
		t.Error("Expected view to render with dimensions set")
	}

	// Test setting different dimensions
	model.SetDimensions(80, 24)
	view = model.View()
	if view == "" {
		t.Error("Expected view to render with new dimensions")
	}
}

// TestHelpPaneFocusBlur tests focus and blur functionality
func TestHelpPaneFocusBlur(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Test focus
	model.Focus()
	// Since focused is private, we'll test through behavior
	view := model.View()
	if view == "" {
		t.Error("Expected view to render when focused")
	}

	// Test blur
	model.Blur()
	view = model.View()
	if view == "" {
		t.Error("Expected view to render when blurred")
	}
}

// TestHelpPaneKeyHandling tests key handling
func TestHelpPaneKeyHandling(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)
	model.SetDimensions(100, 30)

	// Enable help mode
	shortcutManager.SetHelpMode(true)
	if !shortcutManager.IsHelpMode() {
		t.Error("Expected help mode to be enabled")
	}

	// Test ESC key
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(escMsg)

	// Verify help mode was disabled
	if shortcutManager.IsHelpMode() {
		t.Error("Expected help mode to be disabled after ESC")
	}

	// Verify model was updated
	if updatedModel == nil {
		t.Error("Expected model to be updated")
	}

	// Re-enable help mode for other tests
	shortcutManager.SetHelpMode(true)

	// Test Q key
	qMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, _ = model.Update(qMsg) // Ignore the returned model

	// Verify help mode was disabled
	if shortcutManager.IsHelpMode() {
		t.Error("Expected help mode to be disabled after Q")
	}

	// Re-enable help mode for other tests
	shortcutManager.SetHelpMode(true)

	// Test Enter key
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _ = model.Update(enterMsg) // Ignore the returned model

	// Verify help mode was disabled
	if shortcutManager.IsHelpMode() {
		t.Error("Expected help mode to be disabled after Enter")
	}
}

// TestHelpPaneViewRendering tests view rendering
func TestHelpPaneViewRendering(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Test view rendering with dimensions
	model.SetDimensions(100, 30)
	view := model.View()

	if view == "" {
		t.Error("Expected view to render with dimensions set")
	}

	// Verify view contains expected content
	if !contains(view, "Keyboard Shortcuts") {
		t.Error("Expected view to contain keyboard shortcuts title")
	}

	if !contains(view, "Navigation") {
		t.Error("Expected view to contain navigation category")
	}

	if !contains(view, "Edit") {
		t.Error("Expected view to contain edit category")
	}
}

// TestHelpPaneResponsiveMode tests responsive mode behavior
func TestHelpPaneResponsiveMode(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Test full mode (large terminal)
	model.SetDimensions(120, 40)
	view := model.View()

	if view == "" {
		t.Error("Expected view to render in full mode")
	}

	// Verify full mode content
	if !contains(view, "Keyboard Shortcuts Reference") {
		t.Error("Expected full mode title")
	}

	// Test compact mode (medium terminal)
	model.SetDimensions(90, 25)
	view = model.View()

	if view == "" {
		t.Error("Expected view to render in compact mode")
	}

	// Verify compact mode content
	if !contains(view, "Shortcuts") {
		t.Error("Expected compact mode title")
	}

	// Test minimal mode (small terminal)
	model.SetDimensions(70, 20)
	view = model.View()

	if view == "" {
		t.Error("Expected view to render in minimal mode")
	}

	// Verify minimal mode content
	if !contains(view, "Help") {
		t.Error("Expected minimal mode title")
	}
}

// TestHelpPaneShortcutManager tests shortcut manager integration
func TestHelpPaneShortcutManager(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Test getting shortcut manager
	retrievedManager := model.GetShortcutManager()
	if retrievedManager != shortcutManager {
		t.Error("Expected to get the same shortcut manager")
	}

	// Test setting new shortcut manager
	newManager := editor.NewShortcutManager()
	model.SetShortcutManager(newManager)

	retrievedManager = model.GetShortcutManager()
	if retrievedManager != newManager {
		t.Error("Expected to get the new shortcut manager")
	}

	// Test with nil shortcut manager
	model.SetShortcutManager(nil)
	retrievedManager = model.GetShortcutManager()
	if retrievedManager != nil {
		t.Error("Expected nil shortcut manager to be returned")
	}

	// Test view rendering with nil manager
	model.SetDimensions(100, 30)
	view := model.View()
	if view == "" {
		t.Error("Expected view to render even with nil shortcut manager")
	}
}

// TestHelpPaneCategoryRendering tests category rendering
func TestHelpPaneCategoryRendering(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Set dimensions for full mode
	model.SetDimensions(120, 40)

	// Set context to editor
	shortcutManager.SetContext(editor.ContextEditor)

	view := model.View()
	if view == "" {
		t.Error("Expected view to render")
	}

	// Verify expected categories are present
	expectedCategories := []string{
		"Navigation",
		"Edit",
		"Search",
		"File",
		"View",
		"Application",
		"Tools",
	}

	for _, category := range expectedCategories {
		if !contains(view, category) {
			t.Errorf("Expected view to contain category: %s", category)
		}
	}
}

// TestHelpPaneCompactCategoryRendering tests compact category rendering
func TestHelpPaneCompactCategoryRendering(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Set dimensions for compact mode
	model.SetDimensions(90, 25)

	// Set context to editor
	shortcutManager.SetContext(editor.ContextEditor)

	view := model.View()
	if view == "" {
		t.Error("Expected view to render in compact mode")
	}

	// Verify only essential categories are present
	essentialCategories := []string{
		"Navigation",
		"Edit",
		"File",
	}

	for _, category := range essentialCategories {
		if !contains(view, category) {
			t.Errorf("Expected compact view to contain category: %s", category)
		}
	}
}

// TestHelpPaneMinimalRendering tests minimal rendering
func TestHelpPaneMinimalRendering(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Set dimensions for minimal mode
	model.SetDimensions(70, 20)

	// Set context to editor
	shortcutManager.SetContext(editor.ContextEditor)

	view := model.View()
	if view == "" {
		t.Error("Expected view to render in minimal mode")
	}

	// Verify minimal content is present
	if !contains(view, "Help") {
		t.Error("Expected minimal view to contain help title")
	}

	// Verify only essential shortcuts are shown
	// (This is a basic check, actual implementation might vary)
	if !contains(view, "Navigation") {
		t.Error("Expected minimal view to contain navigation shortcuts")
	}
}

// TestHelpPaneContextAwareRendering tests context-aware rendering
func TestHelpPaneContextAwareRendering(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Set dimensions for full mode
	model.SetDimensions(120, 40)

	// Test with editor context
	shortcutManager.SetContext(editor.ContextEditor)
	editorView := model.View()

	// Test with preview context
	shortcutManager.SetContext(editor.ContextPreview)
	previewView := model.View()

	// Views should be different based on context
	if editorView == previewView {
		t.Error("Expected views to be different for different contexts")
	}

	// Test with global context
	shortcutManager.SetContext(editor.ContextGlobal)
	globalView := model.View()

	// Global view should also be different
	if editorView == globalView || previewView == globalView {
		t.Error("Expected global view to be different from other contexts")
	}
}

// TestHelpPaneFooterRendering tests footer rendering
func TestHelpPaneFooterRendering(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Test full mode footer
	model.SetDimensions(120, 40)
	view := model.View()

	if !contains(view, "Press ESC, Q, or Enter to return to editor") {
		t.Error("Expected full mode footer")
	}

	// Test compact mode footer
	model.SetDimensions(90, 25)
	view = model.View()

	if !contains(view, "ESC/Q/Enter: back") {
		t.Error("Expected compact mode footer")
	}

	// Test minimal mode footer
	model.SetDimensions(70, 20)
	view = model.View()

	if !contains(view, "ESC: back") {
		t.Error("Expected minimal mode footer")
	}
}

// TestHelpPaneEdgeCases tests edge cases and boundary conditions
func TestHelpPaneEdgeCases(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Test with zero dimensions
	model.SetDimensions(0, 0)
	view := model.View()
	if view == "" {
		t.Error("Expected view to render with zero dimensions")
	}

	// Test with very small dimensions
	model.SetDimensions(10, 5)
	view = model.View()
	if view == "" {
		t.Error("Expected view to render with very small dimensions")
	}

	// Test with very large dimensions
	model.SetDimensions(500, 200)
	view = model.View()
	if view == "" {
		t.Error("Expected view to render with very large dimensions")
	}

	// Test with nil shortcut manager
	model.SetShortcutManager(nil)
	model.SetDimensions(100, 30)
	view = model.View()
	if view == "" {
		t.Error("Expected view to render with nil shortcut manager")
	}
}

// TestHelpPaneStyleConsistency tests style consistency
func TestHelpPaneStyleConsistency(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Set dimensions for full mode
	model.SetDimensions(120, 40)

	// Render view
	view := model.View()
	if view == "" {
		t.Error("Expected view to render")
	}

	// View should contain styled elements
	// (This is a basic check, actual implementation might vary)
	if !contains(view, "ðŸŽ¹") {
		t.Error("Expected view to contain styled elements")
	}

	if !contains(view, "ðŸ“‚") {
		t.Error("Expected view to contain category icons")
	}
}

// TestHelpPanePerformance tests performance with large content
func TestHelpPanePerformance(t *testing.T) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Set dimensions
	model.SetDimensions(120, 40)

	// Test multiple view renders
	for i := 0; i < 10; i++ {
		view := model.View()
		if view == "" {
			t.Errorf("Expected view to render on iteration %d", i)
		}
	}

	// Test multiple dimension changes
	for i := 0; i < 10; i++ {
		width := 80 + (i * 5)
		height := 24 + (i * 2)
		model.SetDimensions(width, height)

		view := model.View()
		if view == "" {
			t.Errorf("Expected view to render with dimensions %dx%d", width, height)
		}
	}
}

// BenchmarkHelpPaneViewRendering benchmarks view rendering performance
func BenchmarkHelpPaneViewRendering(b *testing.B) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)
	model.SetDimensions(120, 40)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

// BenchmarkHelpPaneUpdate benchmarks update performance
func BenchmarkHelpPaneUpdate(b *testing.B) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)
	model.SetDimensions(120, 40)

	// Test key messages
	keys := []tea.KeyMsg{
		{Type: tea.KeyEsc},
		{Type: tea.KeyRunes, Runes: []rune{'q'}},
		{Type: tea.KeyEnter},
		{Type: tea.KeyRunes, Runes: []rune{'x'}}, // Unhandled key
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := keys[i%len(keys)]
		_, _ = model.Update(key)
	}
}

// BenchmarkHelpPaneResponsiveMode benchmarks responsive mode changes
func BenchmarkHelpPaneResponsiveMode(b *testing.B) {
	shortcutManager := editor.NewShortcutManager()
	model := editor.NewHelpPaneModel(shortcutManager)

	// Test various dimensions
	dimensions := []struct {
		width  int
		height int
	}{
		{70, 20},  // Minimal
		{90, 25},  // Compact
		{120, 40}, // Full
		{150, 50}, // Large
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dim := dimensions[i%len(dimensions)]
		model.SetDimensions(dim.width, dim.height)
		_ = model.View()
	}
}
