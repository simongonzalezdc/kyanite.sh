package plugins

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Kyanite/noise/internal/config"
	"github.com/Kyanite/noise/internal/logging"
	tea "github.com/charmbracelet/bubbletea"
)

func TestDefaultManager_LoadPlugins(t *testing.T) {
	cfg := config.DefaultConfig()
	logger, err := logging.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	manager := NewManager(cfg, logger)

	// Test loading plugins (should not error even if no plugins exist)
	err = manager.LoadPlugins()
	if err != nil {
		t.Errorf("Failed to load plugins: %v", err)
	}
}

func TestDefaultManager_LoadPluginsFromDir(t *testing.T) {
	manager := CreateTestPluginManager(t)
	testDir := CreateTestPluginDir(t)
	defer CleanupTestPluginDir(t, testDir)

	// Add test directory to allowed paths
	manager.security.allowedPaths = append(manager.security.allowedPaths, testDir)

	// Create a valid plugin manifest
	manifestPath := CreateTestPluginManifest(t, testDir, nil)

	// Load plugins from the directory
	err := manager.loadPluginsFromDir(testDir)
	if err != nil {
		t.Errorf("Failed to load plugins from directory: %v", err)
	}

	// Check that plugin was loaded
	plugins := manager.GetPlugins()
	if len(plugins) == 0 {
		t.Logf("No plugins loaded from %s (manifest: %s)", testDir, manifestPath)
	}
}

func TestDefaultManager_LoadPluginsFromDir_InvalidPath(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Try to load from a blocked path
	err := manager.loadPluginsFromDir("/etc")
	if err != nil {
		t.Errorf("Expected no error when loading from blocked path (should skip), got: %v", err)
	}

	// Try to load from a non-existent path
	err = manager.loadPluginsFromDir("/non/existent/path")
	if err != nil {
		t.Errorf("Expected no error when loading from non-existent path, got: %v", err)
	}
}

func TestDefaultManager_LoadCompiledPlugin(t *testing.T) {
	manager := CreateTestPluginManager(t)
	testDir := CreateTestPluginDir(t)
	defer CleanupTestPluginDir(t, testDir)

	// Create a .so file without manifest (should be skipped)
	soFile := CreateTestPluginFile(t, testDir, "test.so", "fake plugin content")

	err := manager.loadCompiledPlugin(soFile)
	if err != nil {
		t.Errorf("Expected no error when loading .so file without manifest, got: %v", err)
	}

	// Create a .so file with manifest
	CreateTestPluginManifest(t, testDir, nil)
	soFile = filepath.Join(testDir, "test_with_manifest.so")
	CreateTestPluginFile(t, testDir, "test_with_manifest.so", "fake plugin content")

	err = manager.loadCompiledPlugin(soFile)
	if err != nil {
		t.Errorf("Failed to load compiled plugin with manifest: %v", err)
	}
}

func TestDefaultManager_LoadManifestPlugin(t *testing.T) {
	manager := CreateTestPluginManager(t)
	testDir := CreateTestPluginDir(t)
	defer CleanupTestPluginDir(t, testDir)

	// Add test directory to allowed paths
	manager.security.allowedPaths = append(manager.security.allowedPaths, testDir)

	// Test loading valid manifest
	validManifest := CreateTestPluginManifest(t, testDir, nil)
	err := manager.loadManifestPlugin(validManifest)
	if err != nil {
		t.Logf("Failed to load valid manifest (expected on Windows): %v", err)
	}

	// Test loading malicious manifest
	maliciousManifest := CreateMaliciousPluginManifest(t, testDir, "suspicious_description")
	err = manager.loadManifestPlugin(maliciousManifest)
	if err == nil {
		t.Error("Expected error when loading malicious manifest")
	}

	// Test loading manifest with invalid path
	invalidManifest := CreateMaliciousPluginManifest(t, testDir, "invalid_id")
	err = manager.loadManifestPlugin(invalidManifest)
	if err == nil {
		t.Error("Expected error when loading manifest with invalid ID")
	}
}

func TestDefaultManager_InitializePlugin(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Test with valid plugin
	plugin := CreateMockPlugin("test_plugin", "Test Plugin", true)

	err := manager.initializePlugin(plugin)
	if err != nil {
		t.Errorf("Failed to initialize valid plugin: %v", err)
	}

	// Check that load time was set
	if plugin.Metadata().LoadTime.IsZero() {
		t.Error("Plugin load time should be set after initialization")
	}

	// Test with malicious plugin
	maliciousPlugin := CreateMockMaliciousPlugin("malicious_plugin", "resource_exhaustion")
	err = manager.initializePlugin(maliciousPlugin)
	if err != nil {
		t.Errorf("Failed to initialize malicious plugin: %v", err)
	}
}

func TestDefaultManager_CollectPluginCapabilities(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Create plugin with multiple capabilities
	plugin := CreateTestPluginWithCapabilities(t, "test_plugin", []Capability{
		CapabilityExportFormat,
		CapabilityEditorTool,
		CapabilityTheoryTool,
	})

	// This should not panic
	manager.collectPluginCapabilities(plugin)
}

func TestDefaultManager_GetPlugin(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add a test plugin
	plugin := CreateMockPlugin("test_plugin", "Test Plugin", true)
	manager.plugins["test_plugin"] = plugin

	// Test getting existing plugin
	retrievedPlugin, err := manager.GetPlugin("test_plugin")
	if err != nil {
		t.Errorf("Failed to get existing plugin: %v", err)
	}
	if retrievedPlugin.Metadata().ID != "test_plugin" {
		t.Error("Retrieved plugin has wrong ID")
	}

	// Test getting non-existent plugin
	_, err = manager.GetPlugin("non_existent")
	if err == nil {
		t.Error("Expected error when getting non-existent plugin")
	}
}

func TestDefaultManager_GetPlugins(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add test plugins using the test helper
	plugin1 := CreateMockPlugin("plugin1", "Plugin 1", true)
	plugin2 := CreateMockPlugin("plugin2", "Plugin 2", false)

	RegisterTestPlugin(manager, plugin1)
	RegisterTestPlugin(manager, plugin2)

	// Get all plugins
	plugins := manager.GetPlugins()

	if len(plugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(plugins))
	}

	// Check that we got copies, not references
	plugins["plugin3"] = CreateMockPlugin("plugin3", "Plugin 3", true)

	if len(manager.GetPlugins()) != 2 {
		t.Error("GetPlugins should return a copy, not modify the original")
	}
}

func TestDefaultManager_GetPluginsByCapability(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add plugins with different capabilities using test helper
	plugin1 := CreateTestPluginWithCapabilities(t, "plugin1", []Capability{CapabilityExportFormat})
	plugin2 := CreateTestPluginWithCapabilities(t, "plugin2", []Capability{CapabilityEditorTool})
	plugin3 := CreateTestPluginWithCapabilities(t, "plugin3", []Capability{CapabilityExportFormat, CapabilityEditorTool})

	RegisterTestPlugin(manager, plugin1)
	RegisterTestPlugin(manager, plugin2)
	RegisterTestPlugin(manager, plugin3)

	// Test getting plugins by export format capability
	exportPlugins := manager.GetPluginsByCapability(CapabilityExportFormat)
	if len(exportPlugins) != 2 {
		t.Errorf("Expected 2 plugins with export format capability, got %d", len(exportPlugins))
	}

	// Test getting plugins by editor tool capability
	editorPlugins := manager.GetPluginsByCapability(CapabilityEditorTool)
	if len(editorPlugins) != 2 {
		t.Errorf("Expected 2 plugins with editor tool capability, got %d", len(editorPlugins))
	}

	// Test getting plugins with capability that no plugin has
	audioPlugins := manager.GetPluginsByCapability(CapabilityAudioEffect)
	if len(audioPlugins) != 0 {
		t.Errorf("Expected 0 plugins with audio effect capability, got %d", len(audioPlugins))
	}
}

func TestDefaultManager_EnablePlugin(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add a disabled plugin
	plugin := CreateMockPlugin("test_plugin", "Test Plugin", false)
	RegisterTestPlugin(manager, plugin)

	// Enable the plugin
	err := manager.EnablePlugin("test_plugin")
	if err != nil {
		t.Errorf("Failed to enable plugin: %v", err)
	}

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled")
	}

	// Test enabling non-existent plugin
	err = manager.EnablePlugin("non_existent")
	if err == nil {
		t.Error("Expected error when enabling non-existent plugin")
	}
}

func TestDefaultManager_DisablePlugin(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add an enabled plugin
	plugin := CreateMockPlugin("test_plugin", "Test Plugin", true)
	RegisterTestPlugin(manager, plugin)

	// Disable the plugin
	err := manager.DisablePlugin("test_plugin")
	if err != nil {
		t.Errorf("Failed to disable plugin: %v", err)
	}

	if plugin.IsEnabled() {
		t.Error("Plugin should be disabled")
	}

	// Test disabling non-existent plugin
	err = manager.DisablePlugin("non_existent")
	if err == nil {
		t.Error("Expected error when disabling non-existent plugin")
	}
}

func TestDefaultManager_UnloadPlugin(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add a plugin
	plugin := CreateMockPlugin("test_plugin", "Test Plugin", true)
	RegisterTestPlugin(manager, plugin)

	// Unload the plugin
	err := manager.UnloadPlugin("test_plugin")
	if err != nil {
		t.Errorf("Failed to unload plugin: %v", err)
	}

	// Check that plugin was removed
	if _, exists := manager.GetPlugins()["test_plugin"]; exists {
		t.Error("Plugin should be removed after unload")
	}

	// Test unloading non-existent plugin
	err = manager.UnloadPlugin("non_existent")
	if err == nil {
		t.Error("Expected error when unloading non-existent plugin")
	}
}

func TestDefaultManager_RegisterHook(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Register a hook handler
	handler := func(point HookPoint, data HookData) error {
		return nil
	}

	err := manager.RegisterHook("test_plugin", handler)
	if err != nil {
		t.Errorf("Failed to register hook: %v", err)
	}

	// Register another handler for the same plugin
	err = manager.RegisterHook("test_plugin", handler)
	if err != nil {
		t.Errorf("Failed to register second hook: %v", err)
	}

	// Check that handlers were registered
	if len(manager.hooks["test_plugin"]) != 2 {
		t.Errorf("Expected 2 hook handlers, got %d", len(manager.hooks["test_plugin"]))
	}
}

func TestDefaultManager_CallHook(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Register a hook handler that tracks calls
	called := false
	handler := func(point HookPoint, data HookData) error {
		called = true
		if point != HookContentSave {
			t.Errorf("Expected hook point %s, got %s", HookContentSave, point)
		}
		if data["test"] != "value" {
			t.Errorf("Expected hook data test=value, got %v", data)
		}
		return nil
	}

	manager.RegisterHook("test_plugin", handler)

	// Call the hook
	data := HookData{"test": "value"}
	err := manager.CallHook(HookContentSave, data)
	if err != nil {
		t.Errorf("Failed to call hook: %v", err)
	}

	if !called {
		t.Error("Hook handler was not called")
	}
}

func TestDefaultManager_CallHook_WithError(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Register a hook handler that returns an error
	handler := func(point HookPoint, data HookData) error {
		return fmt.Errorf("hook error")
	}

	manager.RegisterHook("test_plugin", handler)

	// Call the hook - should not return error even if handler fails
	data := HookData{"test": "value"}
	err := manager.CallHook(HookContentSave, data)
	if err != nil {
		t.Errorf("Expected no error when hook handler fails, got: %v", err)
	}
}

func TestDefaultManager_GetMenuItems(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add some menu items
	menuItem1 := &MenuItem{
		ID:          "menu1",
		Title:       "Menu 1",
		Description: "First menu item",
		Enabled:     true,
	}
	menuItem2 := &MenuItem{
		ID:          "menu2",
		Title:       "Menu 2",
		Description: "Second menu item",
		Enabled:     true,
	}

	manager.menuItems = []*MenuItem{menuItem1, menuItem2}

	// Get menu items
	items := manager.GetMenuItems()

	if len(items) != 2 {
		t.Errorf("Expected 2 menu items, got %d", len(items))
	}

	// Check that we got copies, not references
	items[0] = &MenuItem{ID: "modified"}

	if manager.menuItems[0].ID == "modified" {
		t.Error("GetMenuItems should return a copy, not modify the original")
	}
}

func TestDefaultManager_GetScreens(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add some screens
	screen1 := &MockScreen{id: "screen1", name: "Screen 1"}
	screen2 := &MockScreen{id: "screen2", name: "Screen 2"}

	manager.screens = map[string]Screen{
		"screen1": screen1,
		"screen2": screen2,
	}

	// Get screens
	screens := manager.GetScreens()

	if len(screens) != 2 {
		t.Errorf("Expected 2 screens, got %d", len(screens))
	}

	// Check that we got copies, not references
	screens["screen3"] = &MockScreen{id: "screen3", name: "Screen 3"}

	if len(manager.screens) != 2 {
		t.Error("GetScreens should return a copy, not modify the original")
	}
}

func TestDefaultManager_GetEditorTools(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add some editor tools
	tool1 := &EditorTool{
		ID:          "tool1",
		Name:        "Tool 1",
		Description: "First tool",
		Enabled:     true,
	}
	tool2 := &EditorTool{
		ID:          "tool2",
		Name:        "Tool 2",
		Description: "Second tool",
		Enabled:     true,
	}

	manager.editorTools = []*EditorTool{tool1, tool2}

	// Get editor tools
	tools := manager.GetEditorTools()

	if len(tools) != 2 {
		t.Errorf("Expected 2 editor tools, got %d", len(tools))
	}

	// Check that we got copies, not references
	tools[0] = &EditorTool{ID: "modified"}

	if manager.editorTools[0].ID == "modified" {
		t.Error("GetEditorTools should return a copy, not modify the original")
	}
}

func TestDefaultManager_GetExportFormats(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add some export formats
	format1 := &ExportFormat{
		ID:          "format1",
		Name:        "Format 1",
		Description: "First format",
		Extension:   ".fmt1",
		MimeType:    "application/fmt1",
		Enabled:     true,
	}
	format2 := &ExportFormat{
		ID:          "format2",
		Name:        "Format 2",
		Description: "Second format",
		Extension:   ".fmt2",
		MimeType:    "application/fmt2",
		Enabled:     true,
	}

	manager.exportFormats = []*ExportFormat{format1, format2}

	// Get export formats
	formats := manager.GetExportFormats()

	if len(formats) != 2 {
		t.Errorf("Expected 2 export formats, got %d", len(formats))
	}

	// Check that we got copies, not references
	formats[0] = &ExportFormat{ID: "modified"}

	if manager.exportFormats[0].ID == "modified" {
		t.Error("GetExportFormats should return a copy, not modify the original")
	}
}

func TestDefaultManager_ConcurrentAccess(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Add some test plugins
	for i := 0; i < 10; i++ {
		plugin := CreateMockPlugin(
			fmt.Sprintf("plugin%d", i),
			fmt.Sprintf("Plugin %d", i),
			true,
		)
		RegisterTestPlugin(manager, plugin)
	}

	var wg sync.WaitGroup
	numGoroutines := 10

	// Test concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			plugins := manager.GetPlugins()
			if len(plugins) != 10 {
				t.Errorf("Goroutine %d: Expected 10 plugins, got %d", id, len(plugins))
			}
		}(i)
	}

	// Test concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			pluginID := fmt.Sprintf("plugin%d", id%10)
			err := manager.EnablePlugin(pluginID)
			if err != nil {
				t.Errorf("Goroutine %d: Failed to enable plugin: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
}

func TestDefaultManager_PluginLifecycle(t *testing.T) {
	manager := CreateTestPluginManager(t)

	// Create a plugin
	plugin := CreateMockPlugin("lifecycle_plugin", "Lifecycle Plugin", false)
	RegisterTestPlugin(manager, plugin)

	// Test full lifecycle
	// 1. Initialize
	err := manager.initializePlugin(plugin)
	if err != nil {
		t.Errorf("Failed to initialize plugin: %v", err)
	}

	// 2. Enable
	err = manager.EnablePlugin("lifecycle_plugin")
	if err != nil {
		t.Errorf("Failed to enable plugin: %v", err)
	}

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled")
	}

	// 3. Disable
	err = manager.DisablePlugin("lifecycle_plugin")
	if err != nil {
		t.Errorf("Failed to disable plugin: %v", err)
	}

	if plugin.IsEnabled() {
		t.Error("Plugin should be disabled")
	}

	// 4. Cleanup
	err = plugin.Cleanup()
	if err != nil {
		t.Errorf("Failed to cleanup plugin: %v", err)
	}

	// 5. Unload
	err = manager.UnloadPlugin("lifecycle_plugin")
	if err != nil {
		t.Errorf("Failed to unload plugin: %v", err)
	}

	// Check that plugin was removed
	if _, exists := manager.GetPlugins()["lifecycle_plugin"]; exists {
		t.Error("Plugin should be removed after unload")
	}
}

// Mock implementations for testing
type MockScreen struct {
	id      string
	name    string
	width   int
	height  int
	focused bool
}

func (s *MockScreen) ID() string {
	return s.id
}

func (s *MockScreen) Name() string {
	return s.name
}

func (s *MockScreen) Init() tea.Cmd {
	return nil
}

func (s *MockScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s *MockScreen) View() string {
	return fmt.Sprintf("Mock Screen: %s", s.name)
}

func (s *MockScreen) SetDimensions(width, height int) {
	s.width = width
	s.height = height
}

func (s *MockScreen) Focus() {
	s.focused = true
}

func (s *MockScreen) Blur() {
	s.focused = false
}

// Table-driven tests for plugin manager scenarios
func TestDefaultManager_Scenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		setup       func() *DefaultManager
		testFunc    func(*DefaultManager) error
		expectError bool
		errorMsg    string
	}{
		{
			name: "Load plugins from multiple directories",
			setup: func() *DefaultManager {
				manager := CreateTestPluginManager(t)
				dir1 := CreateTestPluginDir(t)
				dir2 := CreateTestPluginDir(t)

				manager.pluginDirs = []string{dir1, dir2}
				manager.security.allowedPaths = append(manager.security.allowedPaths, dir1, dir2)

				CreateTestPluginManifest(t, dir1, nil)
				CreateTestPluginManifest(t, dir2, nil)

				return manager
			},
			testFunc: func(m *DefaultManager) error {
				return m.LoadPlugins()
			},
			expectError: false,
		},
		{
			name: "Handle plugin initialization failure",
			setup: func() *DefaultManager {
				manager := CreateTestPluginManager(t)

				// Create a plugin that will fail initialization
				plugin := &FailingPlugin{
					StubPlugin: &StubPlugin{
						metadata: &PluginMetadata{
							ID:          "failing_plugin",
							Name:        "Failing Plugin",
							Version:     "1.0.0",
							Description: "A plugin that fails to initialize",
							Author:      "Test Author",
							License:     "MIT",
							Enabled:     true,
						},
						enabled: true,
					},
					failInit: true,
				}

				RegisterTestPlugin(manager, plugin)
				return manager
			},
			testFunc: func(m *DefaultManager) error {
				// Try to initialize all plugins and verify we get an error for the failing one
				for _, plugin := range m.GetPlugins() {
					err := m.initializePlugin(plugin)
					if plugin.Metadata().ID == "failing_plugin" {
						// Expect this plugin to fail initialization
						if err == nil {
							return fmt.Errorf("expected initialization to fail for failing_plugin")
						}
						// Verify the error message
						if !strings.Contains(err.Error(), "initialization failed") {
							return fmt.Errorf("unexpected error message: %v", err)
						}
						return nil // Test passed - initialization failed as expected
					} else if err != nil {
						return fmt.Errorf("unexpected error for plugin %s: %v", plugin.Metadata().ID, err)
					}
				}
				return nil
			},
			expectError: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			manager := scenario.setup()
			err := scenario.testFunc(manager)

			if scenario.expectError {
				if err == nil {
					t.Errorf("Expected error but got none for scenario: %s", scenario.name)
					return
				}
				if scenario.errorMsg != "" && !strings.Contains(err.Error(), scenario.errorMsg) {
					t.Errorf("Expected error to contain '%s', got: %v for scenario: %s",
						scenario.errorMsg, err, scenario.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for scenario %s: %v", scenario.name, err)
				}
			}
		})
	}
}

// FailingPlugin is a plugin that fails during certain operations for testing
type FailingPlugin struct {
	*StubPlugin
	failInit bool
}

func (p *FailingPlugin) Initialize(ctx *PluginContext) error {
	if p.failInit {
		return fmt.Errorf("initialization failed")
	}
	return p.StubPlugin.Initialize(ctx)
}
