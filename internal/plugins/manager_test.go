package plugins

import (
	"os"
	"testing"

	"github.com/puente-labs/lyricforge/internal/config"
	"github.com/puente-labs/lyricforge/internal/logging"
)

func TestPluginManager(t *testing.T) {
	// Create a test configuration
	cfg := config.DefaultConfig()

	// Create a test logger
	logger, err := logging.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Create plugin manager
	manager := NewManager(cfg, logger)

	// Test loading plugins
	err = manager.LoadPlugins()
	if err != nil {
		t.Errorf("Failed to load plugins: %v", err)
	}

	// Check that plugins were loaded (should find the example plugin)
	plugins := manager.GetPlugins()
	if len(plugins) == 0 {
		t.Log("No plugins loaded - this might be expected if no plugin files exist")
	} else {
		t.Logf("Loaded %d plugins", len(plugins))
		for _, plugin := range plugins {
			t.Logf("Plugin: %s v%s - %s", plugin.Metadata().Name, plugin.Metadata().Version, plugin.Metadata().Description)
		}
	}

	// Test plugin capabilities
	exportPlugins := manager.GetPluginsByCapability(CapabilityExportFormat)
	if len(exportPlugins) > 0 {
		t.Logf("Found %d export format plugins", len(exportPlugins))
	}

	theoryPlugins := manager.GetPluginsByCapability(CapabilityTheoryTool)
	if len(theoryPlugins) > 0 {
		t.Logf("Found %d theory tool plugins", len(theoryPlugins))
	}
}

func TestPluginSettingsModel(t *testing.T) {
	// Create a test configuration
	cfg := config.DefaultConfig()

	// Create a test logger
	logger, err := logging.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Create plugin manager
	manager := NewManager(cfg, logger)

	// Create plugin settings model
	settingsModel := NewPluginSettingsModel(manager)

	// Test basic functionality
	if settingsModel == nil {
		t.Fatal("Failed to create plugin settings model")
	}

	// Test dimensions setting
	settingsModel.SetDimensions(100, 50)

	// Test view rendering (should not panic)
	view := settingsModel.View()
	if view == "" {
		t.Log("Empty view - no plugins loaded")
	} else {
		t.Logf("Settings view rendered successfully, length: %d", len(view))
	}
}

func TestPluginMetadata(t *testing.T) {
	// Test plugin metadata structure
	metadata := &PluginMetadata{
		ID:          "test_plugin",
		Name:        "Test Plugin",
		Version:     "1.0.0",
		Description: "A test plugin",
		Author:      "Test Author",
		License:     "MIT",
		Capabilities: []Capability{
			CapabilityExportFormat,
			CapabilityTheoryTool,
		},
		Enabled: true,
	}

	if metadata.ID != "test_plugin" {
		t.Errorf("Expected ID 'test_plugin', got '%s'", metadata.ID)
	}

	if len(metadata.Capabilities) != 2 {
		t.Errorf("Expected 2 capabilities, got %d", len(metadata.Capabilities))
	}

	// Test capability checking
	found := false
	for _, cap := range metadata.Capabilities {
		if cap == CapabilityExportFormat {
			found = true
			break
		}
	}
	if !found {
		t.Error("Export format capability should be present")
	}
}

func TestPluginDirectoryCreation(t *testing.T) {
	// Test that plugin directories are created correctly
	cfg := config.DefaultConfig()
	pluginDirs := getDefaultPluginDirs(cfg)

	for _, dir := range pluginDirs {
		// Check if directory exists or can be created
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// Try to create it
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Errorf("Failed to create plugin directory %s: %v", dir, err)
			} else {
				t.Logf("Created plugin directory: %s", dir)
				// Clean up
				os.RemoveAll(dir)
			}
		} else {
			t.Logf("Plugin directory exists: %s", dir)
		}
	}
}
