package plugins

import (
	"testing"
	"time"

	"github.com/kyanite/noise/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

func TestPluginMetadata_Validation(t *testing.T) {
	tests := []struct {
		name        string
		metadata    *PluginMetadata
		expectValid bool
	}{
		{
			name: "Valid plugin metadata",
			metadata: &PluginMetadata{
				ID:          "test_plugin",
				Name:        "Test Plugin",
				Version:     "1.0.0",
				Description: "A test plugin",
				Author:      "Test Author",
				License:     "MIT",
				Capabilities: []Capability{
					CapabilityExportFormat,
				},
				Enabled: true,
			},
			expectValid: true,
		},
		{
			name: "Plugin with multiple capabilities",
			metadata: &PluginMetadata{
				ID:          "multi_cap_plugin",
				Name:        "Multi Capability Plugin",
				Version:     "2.1.0",
				Description: "A plugin with multiple capabilities",
				Author:      "Test Author",
				License:     "Apache-2.0",
				Capabilities: []Capability{
					CapabilityExportFormat,
					CapabilityEditorTool,
					CapabilityTheoryTool,
				},
				Enabled: true,
			},
			expectValid: true,
		},
		{
			name: "Plugin with config",
			metadata: &PluginMetadata{
				ID:          "config_plugin",
				Name:        "Config Plugin",
				Version:     "1.5.2",
				Description: "A plugin with configuration",
				Author:      "Test Author",
				License:     "MIT",
				Capabilities: []Capability{
					CapabilityConfig,
				},
				Config: map[string]interface{}{
					"option1": "value1",
					"option2": 42,
					"option3": true,
				},
				Enabled: true,
			},
			expectValid: true,
		},
		{
			name: "Plugin with no capabilities",
			metadata: &PluginMetadata{
				ID:           "no_cap_plugin",
				Name:         "No Capability Plugin",
				Version:      "0.1.0",
				Description:  "A plugin with no capabilities",
				Author:       "Test Author",
				License:      "MIT",
				Capabilities: []Capability{},
				Enabled:      false,
			},
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - check required fields
			if tt.metadata.ID == "" && tt.expectValid {
				t.Error("Valid plugin should have ID")
			}
			if tt.metadata.Name == "" && tt.expectValid {
				t.Error("Valid plugin should have Name")
			}
			if tt.metadata.Version == "" && tt.expectValid {
				t.Error("Valid plugin should have Version")
			}
		})
	}
}

func TestCapability_Constants(t *testing.T) {
	// Test that all capability constants are defined
	capabilities := []Capability{
		CapabilityMenuItem,
		CapabilityScreen,
		CapabilityUIExtension,
		CapabilityEditorTool,
		CapabilityExportFormat,
		CapabilitySyntax,
		CapabilityTheoryTool,
		CapabilityChordLib,
		CapabilityAudioEffect,
		CapabilityMIDIHandler,
		CapabilityDataProvider,
		CapabilityExporter,
		CapabilityHook,
		CapabilityConfig,
	}

	// Check that all capabilities have string values
	for _, cap := range capabilities {
		if string(cap) == "" {
			t.Errorf("Capability %v should have a string value", cap)
		}
	}
}

func TestPluginContext_Creation(t *testing.T) {
	cfg := config.DefaultConfig()
	metadata := &PluginMetadata{
		ID:      "test_plugin",
		Name:    "Test Plugin",
		Version: "1.0.0",
	}

	ctx := &PluginContext{
		Config:   cfg,
		Metadata: metadata,
	}

	if ctx.Config == nil {
		t.Error("Context should have config")
	}

	if ctx.Metadata == nil {
		t.Error("Context should have metadata")
	}

	if ctx.Metadata.ID != "test_plugin" {
		t.Error("Context metadata should have correct ID")
	}
}

func TestStubPlugin_Implementation(t *testing.T) {
	metadata := &PluginMetadata{
		ID:          "stub_plugin",
		Name:        "Stub Plugin",
		Version:     "1.0.0",
		Description: "A stub plugin for testing",
		Author:      "Test Author",
		License:     "MIT",
		Enabled:     true,
	}

	plugin := &StubPlugin{
		metadata: metadata,
		enabled:  true,
	}

	// Test Metadata
	if plugin.Metadata() != metadata {
		t.Error("Metadata should return the same instance")
	}

	// Test Initialize
	ctx := TestPluginContext(t)
	err := plugin.Initialize(ctx)
	if err != nil {
		t.Errorf("Initialize should not error: %v", err)
	}

	if plugin.Metadata().LoadTime.IsZero() {
		t.Error("Load time should be set after initialization")
	}

	// Test IsEnabled
	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled")
	}

	// Test Disable
	err = plugin.Disable()
	if err != nil {
		t.Errorf("Disable should not error: %v", err)
	}

	if plugin.IsEnabled() {
		t.Error("Plugin should be disabled after Disable()")
	}

	// Test Enable
	err = plugin.Enable()
	if err != nil {
		t.Errorf("Enable should not error: %v", err)
	}

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled after Enable()")
	}

	// Test Cleanup
	err = plugin.Cleanup()
	if err != nil {
		t.Errorf("Cleanup should not error: %v", err)
	}
}

func TestMenuItem_Creation(t *testing.T) {
	handler := func() tea.Cmd {
		return nil
	}

	menuItem := &MenuItem{
		ID:          "test_menu",
		Title:       "Test Menu",
		Description: "A test menu item",
		Shortcut:    "ctrl+t",
		Icon:        "🔧",
		Handler:     handler,
		Enabled:     true,
	}

	if menuItem.ID != "test_menu" {
		t.Error("Menu item ID should be set correctly")
	}

	if menuItem.Handler == nil {
		t.Error("Menu item handler should not be nil")
	}

	if !menuItem.Enabled {
		t.Error("Menu item should be enabled")
	}
}

func TestMenuItem_WithChildren(t *testing.T) {
	child1 := &MenuItem{
		ID:      "child1",
		Title:   "Child 1",
		Enabled: true,
	}
	child2 := &MenuItem{
		ID:      "child2",
		Title:   "Child 2",
		Enabled: false,
	}

	parent := &MenuItem{
		ID:          "parent_menu",
		Title:       "Parent Menu",
		Description: "A menu item with children",
		Children:    []*MenuItem{child1, child2},
		Enabled:     true,
	}

	if len(parent.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(parent.Children))
	}

	if parent.Children[0].ID != "child1" {
		t.Error("First child should be child1")
	}

	if parent.Children[1].ID != "child2" {
		t.Error("Second child should be child2")
	}
}

func TestScreen_Interface(t *testing.T) {
	screen := &MockScreen{
		id:   "test_screen",
		name: "Test Screen",
	}

	// Test ID
	if screen.ID() != "test_screen" {
		t.Error("Screen ID should be set correctly")
	}

	// Test Name
	if screen.Name() != "Test Screen" {
		t.Error("Screen name should be set correctly")
	}

	// Test Init
	cmd := screen.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}

	// Test Update
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := screen.Update(msg)
	if model != screen {
		t.Error("Update should return the same model")
	}
	if cmd != nil {
		t.Error("Update should return nil command")
	}

	// Test View
	view := screen.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	// Test SetDimensions
	screen.SetDimensions(100, 50)
	if screen.width != 100 || screen.height != 50 {
		t.Error("Dimensions should be set correctly")
	}

	// Test Focus/Blur
	screen.Focus()
	if !screen.focused {
		t.Error("Screen should be focused after Focus()")
	}

	screen.Blur()
	if screen.focused {
		t.Error("Screen should not be focused after Blur()")
	}
}

func TestEditorTool_Creation(t *testing.T) {
	handler := func(content string) (string, error) {
		return "processed: " + content, nil
	}

	tool := &EditorTool{
		ID:          "test_tool",
		Name:        "Test Tool",
		Description: "A test editor tool",
		Icon:        "🔧",
		Handler:     handler,
		Shortcut:    "ctrl+t",
		Enabled:     true,
	}

	if tool.ID != "test_tool" {
		t.Error("Tool ID should be set correctly")
	}

	if tool.Handler == nil {
		t.Error("Tool handler should not be nil")
	}

	if !tool.Enabled {
		t.Error("Tool should be enabled")
	}

	// Test handler
	result, err := tool.Handler("test content")
	if err != nil {
		t.Errorf("Handler should not error: %v", err)
	}

	if result != "processed: test content" {
		t.Error("Handler should process content correctly")
	}
}

func TestExportFormat_Creation(t *testing.T) {
	handler := func(content string) ([]byte, error) {
		return []byte("exported: " + content), nil
	}

	format := &ExportFormat{
		ID:          "test_format",
		Name:        "Test Format",
		Description: "A test export format",
		Extension:   ".test",
		MimeType:    "application/test",
		Handler:     handler,
		Enabled:     true,
	}

	if format.ID != "test_format" {
		t.Error("Format ID should be set correctly")
	}

	if format.Handler == nil {
		t.Error("Format handler should not be nil")
	}

	if !format.Enabled {
		t.Error("Format should be enabled")
	}

	// Test handler
	result, err := format.Handler("test content")
	if err != nil {
		t.Errorf("Handler should not error: %v", err)
	}

	if string(result) != "exported: test content" {
		t.Error("Handler should export content correctly")
	}
}

func TestHookPoint_Constants(t *testing.T) {
	// Test that all hook point constants are defined
	hookPoints := []HookPoint{
		HookPreInit,
		HookPostInit,
		HookPreShutdown,
		HookScreenChange,
		HookMenuOpen,
		HookMenuClose,
		HookContentChange,
		HookContentSave,
		HookContentLoad,
		HookPreExport,
		HookPostExport,
		HookAnalysisComplete,
		HookPlaybackStart,
		HookPlaybackStop,
	}

	// Check that all hook points have string values
	for _, point := range hookPoints {
		if string(point) == "" {
			t.Errorf("Hook point %v should have a string value", point)
		}
	}
}

func TestHookData_Usage(t *testing.T) {
	// Create hook data
	data := HookData{
		"content": "test content",
		"file":    "test.txt",
		"count":   42,
	}

	// Test accessing data
	if data["content"] != "test content" {
		t.Error("Hook data should store values correctly")
	}

	if data["count"] != 42 {
		t.Error("Hook data should store different types correctly")
	}

	// Test modifying data
	data["new_field"] = "new value"
	if data["new_field"] != "new value" {
		t.Error("Hook data should be modifiable")
	}
}

func TestHookHandler_Usage(t *testing.T) {
	// Create a hook handler that tracks calls
	called := false
	var receivedPoint HookPoint
	var receivedData HookData

	handler := func(point HookPoint, data HookData) error {
		called = true
		receivedPoint = point
		receivedData = data
		return nil
	}

	// Call the handler
	testData := HookData{"test": "value"}
	err := handler(HookContentSave, testData)
	if err != nil {
		t.Errorf("Handler should not error: %v", err)
	}

	if !called {
		t.Error("Handler should have been called")
	}

	if receivedPoint != HookContentSave {
		t.Error("Handler should receive correct hook point")
	}

	if receivedData["test"] != "value" {
		t.Error("Handler should receive correct data")
	}
}

func TestPluginManager_Interface(t *testing.T) {
	// Test that DefaultManager implements PluginManager interface
	var _ PluginManager = &DefaultManager{}
}

func TestPluginInterface_Compliance(t *testing.T) {
	// Test that StubPlugin implements Plugin interface
	var _ Plugin = &StubPlugin{}

	// Test that MockScreen implements Screen interface
	var _ Screen = &MockScreen{}
}

func TestPluginMetadata_LoadTime(t *testing.T) {
	metadata := &PluginMetadata{
		ID:          "test_plugin",
		Name:        "Test Plugin",
		Version:     "1.0.0",
		Description: "A test plugin",
		Author:      "Test Author",
		License:     "MIT",
		LoadTime:    time.Now(),
		Enabled:     true,
	}

	// Test that load time is set
	if metadata.LoadTime.IsZero() {
		t.Error("Load time should be set")
	}

	// Test load time is recent (within last minute)
	if time.Since(metadata.LoadTime) > time.Minute {
		t.Error("Load time should be recent")
	}
}

func TestPluginMetadata_ConfigSerialization(t *testing.T) {
	// Test that config can store various types
	config := map[string]interface{}{
		"string_val": "test string",
		"int_val":    42,
		"float_val":  3.14,
		"bool_val":   true,
		"array_val":  []string{"a", "b", "c"},
		"nested_val": map[string]interface{}{
			"nested_string": "nested value",
			"nested_int":    123,
		},
	}

	metadata := &PluginMetadata{
		ID:          "config_plugin",
		Name:        "Config Plugin",
		Version:     "1.0.0",
		Description: "A plugin with complex config",
		Author:      "Test Author",
		License:     "MIT",
		Config:      config,
		Enabled:     true,
	}

	// Test accessing config values
	if metadata.Config["string_val"] != "test string" {
		t.Error("Config should store string values correctly")
	}

	if metadata.Config["int_val"] != 42 {
		t.Error("Config should store int values correctly")
	}

	nested := metadata.Config["nested_val"].(map[string]interface{})
	if nested["nested_string"] != "nested value" {
		t.Error("Config should store nested values correctly")
	}
}

// Table-driven tests for plugin type validation
func TestPluginType_Validation(t *testing.T) {
	tests := []struct {
		name        string
		pluginType  interface{}
		expectValid bool
	}{
		{
			name:        "Valid PluginMetadata",
			pluginType:  &PluginMetadata{ID: "test", Name: "Test", Version: "1.0"},
			expectValid: true,
		},
		{
			name:        "Valid MenuItem",
			pluginType:  &MenuItem{ID: "test", Title: "Test"},
			expectValid: true,
		},
		{
			name:        "Valid EditorTool",
			pluginType:  &EditorTool{ID: "test", Name: "Test"},
			expectValid: true,
		},
		{
			name:        "Valid ExportFormat",
			pluginType:  &ExportFormat{ID: "test", Name: "Test"},
			expectValid: true,
		},
		{
			name:        "Valid HookPoint",
			pluginType:  HookContentSave,
			expectValid: true,
		},
		{
			name:        "Valid HookData",
			pluginType:  HookData{"test": "value"},
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - check that the type is not nil
			if tt.pluginType == nil && tt.expectValid {
				t.Error("Valid plugin type should not be nil")
			}
		})
	}
}
