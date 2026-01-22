package plugins

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPluginSettingsModel_NewPluginSettingsModel(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	if model == nil {
		t.Fatal("Failed to create plugin settings model")
	}

	if model.manager != manager {
		t.Error("Manager not set correctly")
	}

	if model.selectedPlugin != 0 {
		t.Errorf("Expected selectedPlugin to be 0, got %d", model.selectedPlugin)
	}

	if model.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to be 0, got %d", model.scrollOffset)
	}

	if model.showDetails {
		t.Error("Expected showDetails to be false")
	}

	if model.detailsPlugin != "" {
		t.Errorf("Expected detailsPlugin to be empty, got '%s'", model.detailsPlugin)
	}
}

func TestPluginSettingsModel_Init(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestPluginSettingsModel_Update_KeyMessages(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	// Add some test plugins
	plugin1 := CreateMockPlugin("plugin1", "Plugin 1", true)
	plugin2 := CreateMockPlugin("plugin2", "Plugin 2", false)

	RegisterTestPlugin(manager, plugin1)
	RegisterTestPlugin(manager, plugin2)

	tests := []struct {
		name          string
		keyMsg        string
		initialState  func(*PluginSettingsModel)
		expectedState func(*PluginSettingsModel) bool
	}{
		{
			name:   "Navigate down",
			keyMsg: "down",
			initialState: func(m *PluginSettingsModel) {
				m.selectedPlugin = 0
				m.showDetails = false
			},
			expectedState: func(m *PluginSettingsModel) bool {
				return m.selectedPlugin == 1 && !m.showDetails
			},
		},
		{
			name:   "Navigate up",
			keyMsg: "up",
			initialState: func(m *PluginSettingsModel) {
				m.selectedPlugin = 1
				m.showDetails = false
			},
			expectedState: func(m *PluginSettingsModel) bool {
				return m.selectedPlugin == 0 && !m.showDetails
			},
		},
		{
			name:   "Navigate down with k key",
			keyMsg: "k",
			initialState: func(m *PluginSettingsModel) {
				m.selectedPlugin = 1
				m.showDetails = false
			},
			expectedState: func(m *PluginSettingsModel) bool {
				return m.selectedPlugin == 0 && !m.showDetails
			},
		},
		{
			name:   "Navigate up with j key",
			keyMsg: "j",
			initialState: func(m *PluginSettingsModel) {
				m.selectedPlugin = 0
				m.showDetails = false
			},
			expectedState: func(m *PluginSettingsModel) bool {
				return m.selectedPlugin == 1 && !m.showDetails
			},
		},
		{
			name:   "Enter to show details",
			keyMsg: "enter",
			initialState: func(m *PluginSettingsModel) {
				m.selectedPlugin = 0
				m.showDetails = false
				m.detailsPlugin = ""
			},
			expectedState: func(m *PluginSettingsModel) bool {
				return m.showDetails && m.detailsPlugin != ""
			},
		},
		{
			name:   "D key to hide details",
			keyMsg: "d",
			initialState: func(m *PluginSettingsModel) {
				m.showDetails = true
				m.detailsPlugin = "plugin1"
			},
			expectedState: func(m *PluginSettingsModel) bool {
				return !m.showDetails && m.detailsPlugin == ""
			},
		},
		{
			name:   "Escape to hide details",
			keyMsg: "esc",
			initialState: func(m *PluginSettingsModel) {
				m.showDetails = true
				m.detailsPlugin = "plugin1"
			},
			expectedState: func(m *PluginSettingsModel) bool {
				return !m.showDetails && m.detailsPlugin == ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initialState(model)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.keyMsg)}
			newModel, _ := model.Update(msg)

			newSettingsModel := newModel.(*PluginSettingsModel)
			if !tt.expectedState(newSettingsModel) {
				t.Errorf("State validation failed for key '%s'", tt.keyMsg)
			}
		})
	}
}

func TestPluginSettingsModel_Update_WindowSizeMsg(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, _ := model.Update(msg)

	newSettingsModel := newModel.(*PluginSettingsModel)
	if newSettingsModel.width != 100 || newSettingsModel.height != 50 {
		t.Errorf("Expected dimensions to be (100, 50), got (%d, %d)",
			newSettingsModel.width, newSettingsModel.height)
	}
}

func TestPluginSettingsModel_View_EmptyState(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)
	model.SetDimensions(80, 24)

	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	if !contains(view, "No plugins loaded") {
		t.Error("View should contain 'No plugins loaded' message")
	}
}

func TestPluginSettingsModel_View_PluginList(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)
	model.SetDimensions(80, 24)

	// Add test plugins
	plugin1 := CreateMockPlugin("plugin1", "Plugin 1", true)
	plugin2 := CreateMockPlugin("plugin2", "Plugin 2", false)

	RegisterTestPlugin(manager, plugin1)
	RegisterTestPlugin(manager, plugin2)

	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	if !contains(view, "Plugin Manager") {
		t.Error("View should contain 'Plugin Manager' title")
	}

	if !contains(view, "Plugin 1") {
		t.Error("View should contain 'Plugin 1'")
	}

	if !contains(view, "Plugin 2") {
		t.Error("View should contain 'Plugin 2'")
	}

	if !contains(view, "Enabled") {
		t.Error("View should contain 'Enabled' status")
	}

	if !contains(view, "Disabled") {
		t.Error("View should contain 'Disabled' status")
	}
}

func TestPluginSettingsModel_View_PluginDetails(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)
	model.SetDimensions(80, 24)

	// Add test plugin
	plugin := CreateMockPlugin("test_plugin", "Test Plugin", true)
	plugin.metadata.Author = "Test Author"
	plugin.metadata.License = "MIT"
	plugin.metadata.Description = "A test plugin for testing"

	RegisterTestPlugin(manager, plugin)

	// Show details for the plugin
	model.showDetails = true
	model.detailsPlugin = "test_plugin"

	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	if !contains(view, "Plugin: Test Plugin") {
		t.Error("View should contain plugin name")
	}

	if !contains(view, "test_plugin") {
		t.Error("View should contain plugin ID")
	}

	if !contains(view, "Test Author") {
		t.Error("View should contain plugin author")
	}

	if !contains(view, "MIT") {
		t.Error("View should contain plugin license")
	}

	if !contains(view, "A test plugin for testing") {
		t.Error("View should contain plugin description")
	}
}

func TestPluginSettingsModel_renderEmptyState(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)
	model.SetDimensions(80, 24)

	view := model.renderEmptyState()
	if view == "" {
		t.Error("Empty state view should not be empty")
	}

	if !contains(view, "No plugins loaded") {
		t.Error("Empty state should contain 'No plugins loaded' message")
	}

	if !contains(view, "~/.noise/plugins/") {
		t.Error("Empty state should contain plugin directory paths")
	}
}

func TestPluginSettingsModel_renderPluginList(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)
	model.SetDimensions(80, 24)

	// Add test plugins
	plugin1 := CreateMockPlugin("plugin1", "Plugin 1", true)
	plugin2 := CreateMockPlugin("plugin2", "Plugin 2", false)

	RegisterTestPlugin(manager, plugin1)
	RegisterTestPlugin(manager, plugin2)

	view := model.renderPluginList()
	if view == "" {
		t.Error("Plugin list view should not be empty")
	}

	if !contains(view, "Plugin Manager") {
		t.Error("Plugin list should contain title")
	}

	if !contains(view, "Plugin 1") {
		t.Error("Plugin list should contain first plugin")
	}

	if !contains(view, "Plugin 2") {
		t.Error("Plugin list should contain second plugin")
	}
}

func TestPluginSettingsModel_renderPluginDetails(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)
	model.SetDimensions(80, 24)

	// Add test plugin with capabilities
	plugin := CreateMockPlugin("test_plugin", "Test Plugin", true)
	plugin.metadata.Author = "Test Author"
	plugin.metadata.License = "MIT"
	plugin.metadata.Description = "A test plugin for testing"
	plugin.metadata.Capabilities = []Capability{
		CapabilityExportFormat,
		CapabilityEditorTool,
	}

	RegisterTestPlugin(manager, plugin)

	model.detailsPlugin = "test_plugin"

	view := model.renderPluginDetails()
	if view == "" {
		t.Error("Plugin details view should not be empty")
	}

	if !contains(view, "Plugin: Test Plugin") {
		t.Error("Plugin details should contain plugin name")
	}

	if !contains(view, "Capabilities:") {
		t.Error("Plugin details should contain capabilities section")
	}

	if !contains(view, "export_format") {
		t.Error("Plugin details should contain export_format capability")
	}

	if !contains(view, "editor_tool") {
		t.Error("Plugin details should contain editor_tool capability")
	}
}

func TestPluginSettingsModel_renderPluginActions(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	// Test with enabled plugin
	enabledPlugin := CreateMockPlugin("enabled_plugin", "Enabled Plugin", true)
	actions := model.renderPluginActions(enabledPlugin)
	if !contains(actions, "Disable Plugin") {
		t.Error("Enabled plugin should show 'Disable Plugin' action")
	}

	// Test with disabled plugin
	disabledPlugin := CreateMockPlugin("disabled_plugin", "Disabled Plugin", false)
	actions = model.renderPluginActions(disabledPlugin)
	if !contains(actions, "Enable Plugin") {
		t.Error("Disabled plugin should show 'Enable Plugin' action")
	}
}

func TestPluginSettingsModel_formatCapabilities(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	capabilities := []Capability{
		CapabilityExportFormat,
		CapabilityEditorTool,
		CapabilityTheoryTool,
	}

	formatted := model.formatCapabilities(capabilities)

	if len(formatted) != 3 {
		t.Errorf("Expected 3 formatted capabilities, got %d", len(formatted))
	}

	if !containsSlice(formatted, "export_format") {
		t.Error("Formatted capabilities should contain 'export_format'")
	}

	if !containsSlice(formatted, "editor_tool") {
		t.Error("Formatted capabilities should contain 'editor_tool'")
	}

	if !containsSlice(formatted, "theory_tool") {
		t.Error("Formatted capabilities should contain 'theory_tool'")
	}
}

func TestPluginSettingsModel_getPluginIDs(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	// Add test plugins
	plugin1 := CreateMockPlugin("plugin1", "Plugin 1", true)
	plugin2 := CreateMockPlugin("plugin2", "Plugin 2", false)
	plugin3 := CreateMockPlugin("plugin3", "Plugin 3", true)

	RegisterTestPlugin(manager, plugin1)
	RegisterTestPlugin(manager, plugin2)
	RegisterTestPlugin(manager, plugin3)

	ids := model.getPluginIDs()

	if len(ids) != 3 {
		t.Errorf("Expected 3 plugin IDs, got %d", len(ids))
	}

	// Check that all expected IDs are present
	expectedIDs := []string{"plugin1", "plugin2", "plugin3"}
	for _, expectedID := range expectedIDs {
		if !containsSlice(ids, expectedID) {
			t.Errorf("Plugin IDs should contain '%s'", expectedID)
		}
	}
}

func TestPluginSettingsModel_SetDimensions(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	model.SetDimensions(100, 50)

	if model.width != 100 || model.height != 50 {
		t.Errorf("Expected dimensions to be (100, 50), got (%d, %d)",
			model.width, model.height)
	}
}

func TestPluginSettingsModel_GetDimensions(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	model.width = 120
	model.height = 60

	width, height := model.GetDimensions()

	if width != 120 || height != 60 {
		t.Errorf("Expected dimensions to be (120, 60), got (%d, %d)", width, height)
	}
}

func TestPluginSettingsModel_FocusBlur(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	// These methods should not panic
	model.Focus()
	model.Blur()
}

func TestPluginSettingsModel_TogglePlugin(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)

	// Add test plugin
	plugin := CreateMockPlugin("toggle_plugin", "Toggle Plugin", true)
	RegisterTestPlugin(manager, plugin)

	// Show details for the plugin
	model.showDetails = true
	model.detailsPlugin = "toggle_plugin"
	model.selectedPlugin = 0 // Select the plugin

	// Toggle plugin (should disable it)
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(msg)

	newSettingsModel := newModel.(*PluginSettingsModel)

	// Check that plugin is now disabled
	toggledPlugin, err := manager.GetPlugin("toggle_plugin")
	if err != nil {
		t.Fatalf("Failed to get plugin: %v", err)
	}

	if toggledPlugin.IsEnabled() {
		t.Error("Plugin should be disabled after toggle")
	}

	// Toggle again (should enable it)
	newModel, _ = newSettingsModel.Update(msg)
	newSettingsModel = newModel.(*PluginSettingsModel)

	toggledPlugin, err = manager.GetPlugin("toggle_plugin")
	if err != nil {
		t.Fatalf("Failed to get plugin: %v", err)
	}

	if !toggledPlugin.IsEnabled() {
		t.Error("Plugin should be enabled after second toggle")
	}
}

func TestPluginSettingsModel_Scrolling(t *testing.T) {
	manager := CreateTestPluginManager(t)
	model := NewPluginSettingsModel(manager)
	model.SetDimensions(80, 10) // Small height to test scrolling

	// Add many plugins to test scrolling
	for i := 0; i < 20; i++ {
		plugin := CreateMockPlugin(
			fmt.Sprintf("plugin%d", i),
			fmt.Sprintf("Plugin %d", i),
			true,
		)
		RegisterTestPlugin(manager, plugin)
	}

	// Test navigation down - selectedPlugin should increase
	initialSelected := model.selectedPlugin
	for i := 0; i < 15; i++ {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
		newModel, _ := model.Update(msg)
		model = newModel.(*PluginSettingsModel)
	}

	if model.selectedPlugin <= initialSelected {
		t.Error("Selected plugin index should have increased")
	}

	// Test navigation up - selectedPlugin should decrease
	highSelected := model.selectedPlugin
	for i := 0; i < 10; i++ {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
		newModel, _ := model.Update(msg)
		model = newModel.(*PluginSettingsModel)
	}

	if model.selectedPlugin >= highSelected {
		t.Error("Selected plugin index should have decreased")
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 1; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func containsSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
