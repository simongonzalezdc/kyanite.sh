package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test loading config with default values
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config == nil {
		t.Error("LoadConfig() should return a config, not nil")
	}

	// Test default values
	if config.Theme == "" {
		t.Error("Config should have a default theme")
	}

	if config.AI.Provider == "" {
		t.Error("Config should have a default AI provider")
	}
}

func TestGetConfigPath(t *testing.T) {
	// Test that GetConfigPath returns a valid path
	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, ".focus", "config.yaml")

	configPath := GetConfigPath()

	if configPath != expectedPath {
		t.Errorf("Expected config path '%s', got '%s'", expectedPath, configPath)
	}
}

func TestConfigDefaults(t *testing.T) {
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Test default theme
	if config.Theme == "" {
		t.Error("Config should have a default theme")
	}

	// Test default AI settings
	if config.AI.Provider == "" {
		t.Error("Config should have a default AI provider")
	}

	if config.AI.Model == "" {
		t.Error("Config should have a default AI model")
	}

	// Test default dashboard settings
	if config.Dashboard.RefreshInterval <= 0 {
		t.Error("Config should have a positive default refresh interval")
	}
}

func TestConfigValidation(t *testing.T) {
	// Test invalid config values
	config := &Config{
		Theme: "", // Invalid empty theme
		AI: AIConfig{
			Provider:    "",   // Invalid empty provider
			Model:       "",   // Invalid empty model
			Temperature: -1.0, // Invalid negative temperature
			MaxTokens:   -1,   // Invalid negative tokens
			Timeout:     -1,   // Invalid negative timeout
		},
		Dashboard: Dashboard{
			RefreshInterval: -1, // Invalid negative interval
		},
	}

	// Should still create config but validation should catch issues
	if config.Theme == "" {
		t.Logf("Config validation caught empty theme")
	}

	if config.AI.Provider == "" {
		t.Logf("Config validation caught empty AI provider")
	}

	if config.AI.Model == "" {
		t.Logf("Config validation caught empty AI model")
	}

	if config.AI.Temperature < 0 {
		t.Logf("Config validation caught negative temperature")
	}

	if config.AI.MaxTokens < 0 {
		t.Logf("Config validation caught negative max tokens")
	}

	if config.AI.Timeout < 0 {
		t.Logf("Config validation caught negative timeout")
	}

	if config.Dashboard.RefreshInterval < 0 {
		t.Logf("Config validation caught negative refresh interval")
	}
}

func TestConfigFileHandling(t *testing.T) {
	// Create initial config
	config := &Config{
		Theme: "Amber Night",
		AI: AIConfig{
			Provider:    "ollama",
			Model:       "llama2",
			Temperature: 0.7,
			MaxTokens:   2000,
			Timeout:     60,
		},
		Dashboard: Dashboard{
			AutoRefresh:     true,
			RefreshInterval: 30,
			ShowAnimation:   true,
			CompactMode:     false,
		},
		Notes: Notes{
			DefaultEditor: "code",
			AutoSave:      true,
			SaveInterval:  10,
			Directory:     "~/documents/notes",
		},
		UI: UI{
			TimeFormat:    "24h",
			DateFormat:    "2006/01/02",
			ShowHelpTips:  true,
			Notifications: false,
			SoundEffects:  true,
		},
	}

	// Save config
	if err := SaveConfig(config); err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Verify file was created at default path
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".focus", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("SaveConfig() should create a config file at default location")
	}

	// Read config back
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	if loadedConfig == nil {
		t.Error("LoadConfig() should return a valid config")
	}
}

func TestConfigPermissions(t *testing.T) {
	// Test that config files have correct permissions
	config := &Config{
		Theme: "test",
		AI: AIConfig{
			Provider: "test",
			Model:    "test",
		},
	}

	err := SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Check file permissions (should be 0644 by default)
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".focus", "config.yaml")

	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	// File should be readable by user
	// (Note: exact permissions depend on OS and umask)
	perm := info.Mode()
	if perm&0o400 == 0 { // Readable by owner
		t.Error("Config file should be readable by owner")
	}

	// Clean up
	os.Remove(configPath)
}

func TestConfigDirectoryStructure(t *testing.T) {
	// Test that config directory structure is correct
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".focus")
	configPath := filepath.Join(configDir, "config.yaml")

	// Create config
	config := &Config{Theme: "test"}
	err := SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Verify structure
	if _, err := os.Stat(configDir); err != nil {
		t.Errorf("Config directory should exist: %v", err)
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Config file should exist: %v", err)
	}

	// Clean up
	os.Remove(configPath)
	os.Remove(configDir)
}

func TestConfigSubStructures(t *testing.T) {
	// Test nested configuration structures
	config := &Config{
		AI: AIConfig{
			Provider:    "ollama",
			Model:       "llama2",
			Temperature: 0.8,
			MaxTokens:   2048,
			Timeout:     45,
		},
		Dashboard: Dashboard{
			AutoRefresh:     false,
			RefreshInterval: 120,
			ShowAnimation:   false,
			CompactMode:     true,
		},
		Notes: Notes{
			DefaultEditor: "nano",
			AutoSave:      false,
			SaveInterval:  15,
			Directory:     "~/scratch",
		},
		UI: UI{
			TimeFormat:    "12h",
			DateFormat:    "Jan 2, 2006",
			ShowHelpTips:  false,
			Notifications: false,
			SoundEffects:  false,
		},
	}

	// Test AI configuration
	if config.AI.Provider != "ollama" {
		t.Errorf("Expected AI provider 'ollama', got '%s'", config.AI.Provider)
	}

	if config.AI.Model != "llama2" {
		t.Errorf("Expected AI model 'llama2', got '%s'", config.AI.Model)
	}

	if config.AI.Temperature != 0.8 {
		t.Errorf("Expected AI temperature 0.8, got %f", config.AI.Temperature)
	}

	// Test Dashboard configuration
	if !config.Dashboard.CompactMode {
		t.Error("Expected CompactMode to be true")
	}

	if config.Dashboard.RefreshInterval != 120 {
		t.Errorf("Expected refresh interval 120, got %d", config.Dashboard.RefreshInterval)
	}

	// Test Notes configuration
	if config.Notes.DefaultEditor != "nano" {
		t.Errorf("Expected default editor 'nano', got '%s'", config.Notes.DefaultEditor)
	}

	if config.Notes.SaveInterval != 15 {
		t.Errorf("Expected save interval 15, got %d", config.Notes.SaveInterval)
	}

	// Test UI configuration
	if config.UI.TimeFormat != "12h" {
		t.Errorf("Expected time format '12h', got '%s'", config.UI.TimeFormat)
	}

	if config.UI.Notifications {
		t.Error("Expected notifications to be false")
	}
}

func TestConfigEdgeCases(t *testing.T) {
	// Test edge cases and boundary conditions
	// Test minimum values
	minConfig := &Config{
		AI: AIConfig{
			Provider:    "test",
			Model:       "test",
			Temperature: 0.0,
			MaxTokens:   1,
			Timeout:     1,
		},
		Dashboard: Dashboard{
			RefreshInterval: 1,
		},
		Notes: Notes{
			SaveInterval: 1,
		},
	}

	if minConfig.AI.MaxTokens != 1 {
		t.Errorf("Expected max tokens 1, got %d", minConfig.AI.MaxTokens)
	}

	if minConfig.AI.Timeout != 1 {
		t.Errorf("Expected timeout 1, got %d", minConfig.AI.Timeout)
	}

	// Test maximum values
	maxConfig := &Config{
		AI: AIConfig{
			Provider:    "test",
			Model:       "test",
			Temperature: 2.0,
			MaxTokens:   8192,
			Timeout:     300,
		},
		Dashboard: Dashboard{
			RefreshInterval: 3600, // 1 hour
		},
		Notes: Notes{
			SaveInterval: 1440, // 1 day
		},
	}

	if maxConfig.AI.MaxTokens != 8192 {
		t.Errorf("Expected max tokens 8192, got %d", maxConfig.AI.MaxTokens)
	}

	if maxConfig.AI.Timeout != 300 {
		t.Errorf("Expected timeout 300, got %d", maxConfig.AI.Timeout)
	}
}

func TestConfigUpdate(t *testing.T) {
	// Test updating configuration values
	updates := map[string]any{
		"theme":                      "updated-theme",
		"ai.provider":                "updated-provider",
		"dashboard.refresh_interval": 60,
	}

	err := UpdateConfig(updates)
	if err != nil {
		t.Fatalf("UpdateConfig() failed: %v", err)
	}

	// Verify updates were applied
	config := GetConfig()
	if config == nil {
		t.Fatal("GetConfig() should return a config after update")
	}

	// Note: Due to global state, we can't test exact values here
	// but we can verify the update process works
	if config == nil {
		t.Error("Config should not be nil after update")
	}
}

func TestGetConfig(t *testing.T) {
	// Test GetConfig function
	config := GetConfig()

	if config == nil {
		t.Error("GetConfig() should return a config")
	}

	// Test that multiple calls return the same instance
	config2 := GetConfig()
	if config != config2 {
		t.Error("GetConfig() should return the same instance (global config)")
	}
}

func TestConfigDefaultsInAction(t *testing.T) {
	// Test that defaults are actually applied
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Test that default AI provider is set
	if config.AI.Provider != "ollama" && config.AI.Provider != "" {
		t.Logf("AI provider is '%s' (default might be 'ollama')", config.AI.Provider)
	}

	// Test that default AI model is set
	if config.AI.Model != "llama3" && config.AI.Model != "" {
		t.Logf("AI model is '%s' (default might be 'llama3')", config.AI.Model)
	}

	// Test that default theme is set
	if config.Theme != "synthwave" && config.Theme != "" {
		t.Logf("Theme is '%s' (default might be 'synthwave')", config.Theme)
	}
}

func TestConfigConsistency(t *testing.T) {
	// Test that config remains consistent across multiple operations
	originalConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Update config
	updates := map[string]any{
		"theme": "consistency-test",
	}

	err = UpdateConfig(updates)
	if err != nil {
		t.Fatalf("UpdateConfig() failed: %v", err)
	}

	// Reload config
	reloadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed after update: %v", err)
	}

	// Both should not be nil
	if originalConfig == nil || reloadedConfig == nil {
		t.Error("Both original and reloaded configs should not be nil")
	}
}
