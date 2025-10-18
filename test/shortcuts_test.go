package lyricforge_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/puente-labs/lyricforge/internal/ui/editor"
)

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestShortcutManagerCreation tests shortcut manager creation
func TestShortcutManagerCreation(t *testing.T) {
	manager := editor.NewShortcutManager()

	if manager == nil {
		t.Fatal("Expected shortcut manager to be created, got nil")
	}

	// Test initial context
	context := manager.GetContext()
	if context != editor.ContextGlobal {
		t.Errorf("Expected initial context to be ContextGlobal, got %v", context)
	}

	// Test initial help mode
	if manager.IsHelpMode() {
		t.Error("Expected help mode to be disabled initially")
	}
}

// TestShortcutManagerContext tests context switching
func TestShortcutManagerContext(t *testing.T) {
	manager := editor.NewShortcutManager()

	// Test setting context
	manager.SetContext(editor.ContextEditor)
	context := manager.GetContext()
	if context != editor.ContextEditor {
		t.Errorf("Expected context to be ContextEditor, got %v", context)
	}

	// Test setting another context
	manager.SetContext(editor.ContextPreview)
	context = manager.GetContext()
	if context != editor.ContextPreview {
		t.Errorf("Expected context to be ContextPreview, got %v", context)
	}

	// Test setting all contexts
	contexts := []editor.KeyContext{
		editor.ContextGlobal,
		editor.ContextEditor,
		editor.ContextPreview,
		editor.ContextSearch,
		editor.ContextMenu,
		editor.ContextHelp,
	}

	for _, ctx := range contexts {
		manager.SetContext(ctx)
		if manager.GetContext() != ctx {
			t.Errorf("Expected context to be %v, got %v", ctx, manager.GetContext())
		}
	}
}

// TestShortcutManagerHelpMode tests help mode functionality
func TestShortcutManagerHelpMode(t *testing.T) {
	manager := editor.NewShortcutManager()

	// Test enabling help mode
	manager.SetHelpMode(true)
	if !manager.IsHelpMode() {
		t.Error("Expected help mode to be enabled")
	}

	// Test disabling help mode
	manager.SetHelpMode(false)
	if manager.IsHelpMode() {
		t.Error("Expected help mode to be disabled")
	}
}

// TestShortcutKeyHandling tests key handling functionality
func TestShortcutKeyHandling(t *testing.T) {
	manager := editor.NewShortcutManager()

	// Test global shortcut (tab)
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	action, handled := manager.HandleKey(tabMsg)

	if !handled {
		t.Error("Expected tab key to be handled")
	}

	if action.Type != editor.ActionNextPane {
		t.Errorf("Expected ActionNextPane for tab key, got %v", action.Type)
	}

	// Test global shortcut (ctrl+q)
	ctrlQMsg := tea.KeyMsg{Type: tea.KeyCtrlC, Runes: []rune{'q'}}
	action, handled = manager.HandleKey(ctrlQMsg)

	if !handled {
		t.Error("Expected ctrl+q key to be handled")
	}

	if action.Type != editor.ActionQuit {
		t.Errorf("Expected ActionQuit for ctrl+q key, got %v", action.Type)
	}

	// Test help mode toggle (F1)
	f1Msg := tea.KeyMsg{Type: tea.KeyF1}
	action, handled = manager.HandleKey(f1Msg)

	if !handled {
		t.Error("Expected F1 key to be handled")
	}

	if action.Type != editor.ActionToggleHelp {
		t.Errorf("Expected ActionToggleHelp for F1 key, got %v", action.Type)
	}

	// Verify help mode was toggled
	if !manager.IsHelpMode() {
		t.Error("Expected help mode to be enabled after F1 press")
	}
}

// TestShortcutContextSpecificKeyHandling tests context-specific key handling
func TestShortcutContextSpecificKeyHandling(t *testing.T) {
	manager := editor.NewShortcutManager()

	// Test editor context shortcuts
	manager.SetContext(editor.ContextEditor)

	// Test editor navigation (home)
	homeMsg := tea.KeyMsg{Type: tea.KeyHome}
	action, handled := manager.HandleKey(homeMsg)

	if !handled {
		t.Error("Expected home key to be handled in editor context")
	}

	if action.Type != editor.ActionStartOfLine {
		t.Errorf("Expected ActionStartOfLine for home key, got %v", action.Type)
	}

	// Test editor editing (ctrl+c)
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	action, handled = manager.HandleKey(ctrlCMsg)

	if !handled {
		t.Error("Expected ctrl+c key to be handled in editor context")
	}

	if action.Type != editor.ActionCopy {
		t.Errorf("Expected ActionCopy for ctrl+c key, got %v", action.Type)
	}

	// Test preview context shortcuts
	manager.SetContext(editor.ContextPreview)

	// Test preview navigation (up)
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	action, handled = manager.HandleKey(upMsg)

	if !handled {
		t.Error("Expected up key to be handled in preview context")
	}

	if action.Type != editor.ActionPreviewUp {
		t.Errorf("Expected ActionPreviewUp for up key, got %v", action.Type)
	}
}

// TestShortcutBindingsByContext tests getting bindings by context
func TestShortcutBindingsByContext(t *testing.T) {
	manager := editor.NewShortcutManager()

	// Test getting global bindings
	globalBindings := manager.GetBindingsForContext(editor.ContextGlobal)
	if len(globalBindings) == 0 {
		t.Error("Expected global bindings to be available")
	}

	// Test getting editor bindings
	editorBindings := manager.GetBindingsForContext(editor.ContextEditor)
	if len(editorBindings) == 0 {
		t.Error("Expected editor bindings to be available")
	}

	// Verify global bindings are included in editor bindings
	if len(editorBindings) <= len(globalBindings) {
		t.Error("Expected editor bindings to include global bindings")
	}

	// Test getting preview bindings
	previewBindings := manager.GetBindingsForContext(editor.ContextPreview)
	if len(previewBindings) == 0 {
		t.Error("Expected preview bindings to be available")
	}

	// Verify global bindings are included in preview bindings
	if len(previewBindings) <= len(globalBindings) {
		t.Error("Expected preview bindings to include global bindings")
	}
}

// TestShortcutBindingsByCategory tests getting bindings by category
func TestShortcutBindingsByCategory(t *testing.T) {
	manager := editor.NewShortcutManager()
	manager.SetContext(editor.ContextEditor)

	// Test getting navigation bindings
	navBindings := manager.GetBindingsByCategory("Navigation")
	if len(navBindings) == 0 {
		t.Error("Expected navigation bindings to be available")
	}

	// Test getting edit bindings
	editBindings := manager.GetBindingsByCategory("Edit")
	if len(editBindings) == 0 {
		t.Error("Expected edit bindings to be available")
	}

	// Test getting file bindings
	fileBindings := manager.GetBindingsByCategory("File")
	if len(fileBindings) == 0 {
		t.Error("Expected file bindings to be available")
	}

	// Test getting view bindings
	viewBindings := manager.GetBindingsByCategory("View")
	if len(viewBindings) == 0 {
		t.Error("Expected view bindings to be available")
	}
}

// TestShortcutAllBindings tests getting all bindings
func TestShortcutAllBindings(t *testing.T) {
	manager := editor.NewShortcutManager()

	// Test getting all bindings
	allBindings := manager.GetAllBindings()
	if len(allBindings) == 0 {
		t.Error("Expected bindings to be available")
	}

	// Verify bindings have required properties
	for _, binding := range allBindings {
		if binding.Key.Help().Key == "" {
			t.Error("Expected binding to have key help")
		}

		if binding.Description == "" {
			t.Error("Expected binding to have description")
		}

		if binding.Category == "" {
			t.Error("Expected binding to have category")
		}
	}
}

// TestShortcutHelpText tests help text generation
func TestShortcutHelpText(t *testing.T) {
	manager := editor.NewShortcutManager()
	manager.SetContext(editor.ContextEditor)

	// Test getting help text
	helpText := manager.GetHelpText()
	if helpText == "" {
		t.Error("Expected help text to be generated")
	}

	// Verify help text contains expected categories
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
		if !containsString(helpText, category) {
			t.Errorf("Expected help text to contain category: %s", category)
		}
	}
}

// TestShortcutStatusBarHints tests status bar hints generation
func TestShortcutStatusBarHints(t *testing.T) {
	manager := editor.NewShortcutManager()

	// Test getting hints for editor context
	manager.SetContext(editor.ContextEditor)
	hints := manager.GetStatusBarHints()
	if hints == "" {
		t.Error("Expected status bar hints to be generated")
	}

	// Verify hints contain expected shortcuts
	if !containsString(hints, "Ctrl+F") {
		t.Error("Expected hints to contain Ctrl+F")
	}

	if !containsString(hints, "Ctrl+S") {
		t.Error("Expected hints to contain Ctrl+S")
	}

	if !containsString(hints, "Tab") {
		t.Error("Expected hints to contain Tab")
	}

	// Test getting hints for preview context
	manager.SetContext(editor.ContextPreview)
	hints = manager.GetStatusBarHints()
	if hints == "" {
		t.Error("Expected status bar hints to be generated for preview")
	}

	// Verify hints contain expected shortcuts
	if !containsString(hints, "↑↓") {
		t.Error("Expected hints to contain arrow keys")
	}

	// Test getting hints for global context
	manager.SetContext(editor.ContextGlobal)
	hints = manager.GetStatusBarHints()
	if hints == "" {
		t.Error("Expected status bar hints to be generated for global")
	}

	// Verify hints contain expected shortcuts
	if !containsString(hints, "F1") {
		t.Error("Expected hints to contain F1")
	}

	if !containsString(hints, "Esc") {
		t.Error("Expected hints to contain Esc")
	}
}

// TestShortcutActionString tests action string representation
func TestShortcutActionString(t *testing.T) {
	// Test known actions
	actions := []struct {
		actionType editor.ShortcutActionType
		expected   string
	}{
		{editor.ActionNextPane, "NextPane"},
		{editor.ActionPrevPane, "PrevPane"},
		{editor.ActionCopy, "Copy"},
		{editor.ActionPaste, "Paste"},
		{editor.ActionSave, "Save"},
		{editor.ActionFind, "Find"},
		{editor.ActionQuit, "Quit"},
		{editor.ActionToggleHelp, "ToggleHelp"},
		{editor.ActionUnknown, "Unknown"},
	}

	for _, test := range actions {
		result := test.actionType.String()
		if result != test.expected {
			t.Errorf("Expected action string '%s', got '%s'", test.expected, result)
		}
	}
}

// TestShortcutComplexKeyBindings tests complex key bindings
func TestShortcutComplexKeyBindings(t *testing.T) {
	manager := editor.NewShortcutManager()
	manager.SetContext(editor.ContextEditor)

	// Test complex key combinations

	// Ctrl+Shift+Home (select to start of file)
	ctrlShiftHome := tea.KeyMsg{
		Type: tea.KeyHome,
		Alt:  true,
	}
	action, handled := manager.HandleKey(ctrlShiftHome)

	if !handled {
		t.Error("Expected ctrl+shift+home to be handled")
	}

	if action.Type != editor.ActionSelectToStartOfFile {
		t.Errorf("Expected ActionSelectToStartOfFile for ctrl+shift+home, got %v", action.Type)
	}

	// Ctrl+Shift+S (save as)
	ctrlShiftS := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Alt:   true,
		Runes: []rune{'s'},
	}
	action, handled = manager.HandleKey(ctrlShiftS)

	if !handled {
		t.Error("Expected ctrl+shift+s to be handled")
	}

	if action.Type != editor.ActionSaveAs {
		t.Errorf("Expected ActionSaveAs for ctrl+shift+s, got %v", action.Type)
	}
}

// TestShortcutContextActivation tests context activation behavior
func TestShortcutContextActivation(t *testing.T) {
	manager := editor.NewShortcutManager()

	// Test that preview shortcuts only work in preview context
	upMsg := tea.KeyMsg{Type: tea.KeyUp}

	// In global context, up should not be handled as preview navigation
	action, handled := manager.HandleKey(upMsg)
	if handled && action.Type == editor.ActionPreviewUp {
		t.Error("Expected up key to not be handled as preview navigation in global context")
	}

	// In preview context, up should be handled as preview navigation
	manager.SetContext(editor.ContextPreview)
	action, handled = manager.HandleKey(upMsg)
	if !handled {
		t.Error("Expected up key to be handled in preview context")
	}
	if action.Type != editor.ActionPreviewUp {
		t.Errorf("Expected ActionPreviewUp for up key in preview context, got %v", action.Type)
	}
}

// TestShortcutActionCreation tests action creation from bindings
func TestShortcutActionCreation(t *testing.T) {
	manager := editor.NewShortcutManager()
	manager.SetContext(editor.ContextEditor)

	// Test that actions are created with correct properties
	ctrlSMsg := tea.KeyMsg{Type: tea.KeyCtrlS}
	action, handled := manager.HandleKey(ctrlSMsg)

	if !handled {
		t.Error("Expected ctrl+s to be handled")
	}

	if action.Type != editor.ActionSave {
		t.Errorf("Expected ActionSave for ctrl+s, got %v", action.Type)
	}

	if action.Description == "" {
		t.Error("Expected action to have description")
	}

	if action.Category == "" {
		t.Error("Expected action to have category")
	}

	// Verify description matches binding
	if action.Description != "Save" {
		t.Errorf("Expected action description 'Save', got '%s'", action.Description)
	}

	// Verify category matches binding
	if action.Category != "File" {
		t.Errorf("Expected action category 'File', got '%s'", action.Category)
	}
}

// TestShortcutUnrecognizedKeys tests handling of unrecognized keys
func TestShortcutUnrecognizedKeys(t *testing.T) {
	manager := editor.NewShortcutManager()

	// Test unrecognized key
	unrecognizedMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	action, handled := manager.HandleKey(unrecognizedMsg)

	if handled {
		t.Error("Expected unrecognized key to not be handled")
	}

	if action.Type != editor.ActionUnknown {
		t.Errorf("Expected ActionUnknown for unrecognized key, got %v", action.Type)
	}
}

// TestShortcutKeyConflicts tests key conflict handling
func TestShortcutKeyConflicts(t *testing.T) {
	manager := editor.NewShortcutManager()
	manager.SetContext(editor.ContextEditor)

	// Test potential key conflicts
	// Ctrl+W is used for both "close file" and "toggle word wrap"
	// In editor context, it should toggle word wrap
	ctrlWMsg := tea.KeyMsg{Type: tea.KeyCtrlW}
	action, handled := manager.HandleKey(ctrlWMsg)

	if !handled {
		t.Error("Expected ctrl+w to be handled in editor context")
	}

	if action.Type != editor.ActionToggleWordWrap {
		t.Errorf("Expected ActionToggleWordWrap for ctrl+w in editor context, got %v", action.Type)
	}

	// In global context, it should close file
	manager.SetContext(editor.ContextGlobal)
	action, handled = manager.HandleKey(ctrlWMsg)

	if !handled {
		t.Error("Expected ctrl+w to be handled in global context")
	}

	if action.Type != editor.ActionCloseFile {
		t.Errorf("Expected ActionCloseFile for ctrl+w in global context, got %v", action.Type)
	}
}

// TestShortcutPerformance tests performance with many key presses
func TestShortcutPerformance(t *testing.T) {
	manager := editor.NewShortcutManager()
	manager.SetContext(editor.ContextEditor)

	// Test performance with many key presses
	keys := []tea.KeyMsg{
		{Type: tea.KeyCtrlC},
		{Type: tea.KeyCtrlV},
		{Type: tea.KeyCtrlZ},
		{Type: tea.KeyCtrlY},
		{Type: tea.KeyCtrlF},
		{Type: tea.KeyCtrlS},
		{Type: tea.KeyTab},
		{Type: tea.KeyHome},
		{Type: tea.KeyEnd},
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
	}

	// Process each key multiple times
	for i := 0; i < 100; i++ {
		for _, key := range keys {
			_, handled := manager.HandleKey(key)
			if !handled {
				t.Errorf("Expected key to be handled: %v", key.Type)
			}
		}
	}
}

// BenchmarkShortcutKeyHandling benchmarks key handling performance
func BenchmarkShortcutKeyHandling(b *testing.B) {
	manager := editor.NewShortcutManager()
	manager.SetContext(editor.ContextEditor)

	// Test common keys
	keys := []tea.KeyMsg{
		{Type: tea.KeyCtrlC},
		{Type: tea.KeyCtrlV},
		{Type: tea.KeyCtrlS},
		{Type: tea.KeyTab},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := keys[i%len(keys)]
		_, _ = manager.HandleKey(key)
	}
}

// BenchmarkShortcutHelpTextGeneration benchmarks help text generation
func BenchmarkShortcutHelpTextGeneration(b *testing.B) {
	manager := editor.NewShortcutManager()
	manager.SetContext(editor.ContextEditor)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = manager.GetHelpText()
	}
}

// BenchmarkShortcutStatusBarHints benchmarks status bar hints generation
func BenchmarkShortcutStatusBarHints(b *testing.B) {
	manager := editor.NewShortcutManager()

	contexts := []editor.KeyContext{
		editor.ContextGlobal,
		editor.ContextEditor,
		editor.ContextPreview,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := contexts[i%len(contexts)]
		manager.SetContext(ctx)
		_ = manager.GetStatusBarHints()
	}
}
