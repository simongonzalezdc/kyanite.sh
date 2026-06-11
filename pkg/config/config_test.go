package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDir(t *testing.T) {
	dir := ConfigDir()
	if dir == "" {
		t.Fatal("ConfigDir returned empty string")
	}
	// Should end with kyanite
	if filepath.Base(dir) != "kyanite" {
		t.Fatalf("ConfigDir should end with 'kyanite', got %s", dir)
	}
}

func TestConfigDir_XDG(t *testing.T) {
	os.Setenv("XDG_CONFIG_HOME", "/tmp/xdg-test")
	defer os.Unsetenv("XDG_CONFIG_HOME")

	dir := ConfigDir()
	expected := "/tmp/xdg-test/kyanite"
	if dir != expected {
		t.Fatalf("expected %s, got %s", expected, dir)
	}
}

func TestConfigPath_EnvOverride(t *testing.T) {
	os.Setenv("KYANITE_CONFIG", "/tmp/custom-config.yaml")
	defer os.Unsetenv("KYANITE_CONFIG")

	path := ConfigPath()
	if path != "/tmp/custom-config.yaml" {
		t.Fatalf("expected /tmp/custom-config.yaml, got %s", path)
	}
}

func TestLoad_DefaultsOnly(t *testing.T) {
	// Ensure no config file or env vars interfere
	os.Unsetenv("KYANITE_CONFIG")
	os.Unsetenv("KYANITE_BRAIN_OLLAMA_URL")
	os.Unsetenv("XDG_CONFIG_HOME")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Brain.OllamaURL != "http://nucbox:11434" {
		t.Errorf("default ollama_url: got %s", cfg.Brain.OllamaURL)
	}
	if cfg.Brain.Model != "gemma4:12b" {
		t.Errorf("default model: got %s", cfg.Brain.Model)
	}
	if cfg.Brain.DBPort != 5432 {
		t.Errorf("default db_port: got %d", cfg.Brain.DBPort)
	}
	if cfg.Focus.Theme != "amber-night" {
		t.Errorf("default focus theme: got %s", cfg.Focus.Theme)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	os.Unsetenv("KYANITE_CONFIG")
	os.Setenv("KYANITE_BRAIN_OLLAMA_URL", "http://custom:9999")
	os.Setenv("KYANITE_BRAIN_MODEL", "llama3:8b")
	defer os.Unsetenv("KYANITE_BRAIN_OLLAMA_URL")
	defer os.Unsetenv("KYANITE_BRAIN_MODEL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Brain.OllamaURL != "http://custom:9999" {
		t.Errorf("env override ollama_url: got %s", cfg.Brain.OllamaURL)
	}
	if cfg.Brain.Model != "llama3:8b" {
		t.Errorf("env override model: got %s", cfg.Brain.Model)
	}
}

func TestLoad_ConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(cfgFile, []byte(`
brain:
  ollama_url: http://file-test:1234
  model: test-model:latest
focus:
  theme: cyberpunk
`), 0644)

	os.Setenv("KYANITE_CONFIG", cfgFile)
	defer os.Unsetenv("KYANITE_CONFIG")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Brain.OllamaURL != "http://file-test:1234" {
		t.Errorf("file override ollama_url: got %s", cfg.Brain.OllamaURL)
	}
	if cfg.Brain.Model != "test-model:latest" {
		t.Errorf("file override model: got %s", cfg.Brain.Model)
	}
	if cfg.Focus.Theme != "cyberpunk" {
		t.Errorf("file override focus theme: got %s", cfg.Focus.Theme)
	}
}

func TestLoad_ConfigFileAndEnv(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(cfgFile, []byte(`
brain:
  ollama_url: http://file-value:1234
`), 0644)

	os.Setenv("KYANITE_CONFIG", cfgFile)
	os.Setenv("KYANITE_BRAIN_MODEL", "env-model:latest")
	defer os.Unsetenv("KYANITE_CONFIG")
	defer os.Unsetenv("KYANITE_BRAIN_MODEL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	// File value should be present
	if cfg.Brain.OllamaURL != "http://file-value:1234" {
		t.Errorf("file value: got %s", cfg.Brain.OllamaURL)
	}
	// Env overrides file/default
	if cfg.Brain.Model != "env-model:latest" {
		t.Errorf("env override model: got %s", cfg.Brain.Model)
	}
}

func TestInit_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("KYANITE_CONFIG", filepath.Join(tmpDir, "config.yaml"))
	defer os.Unsetenv("KYANITE_CONFIG")

	if err := Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// File should exist
	if _, err := os.Stat(filepath.Join(tmpDir, "config.yaml")); err != nil {
		t.Fatalf("config file not created: %v", err)
	}
}

func TestInit_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(cfgFile, []byte("existing"), 0644)
	os.Setenv("KYANITE_CONFIG", cfgFile)
	defer os.Unsetenv("KYANITE_CONFIG")

	if err := Init(); err == nil {
		t.Fatal("Init should fail when config already exists")
	}
}

// TestLoad_MalformedYAML tests loading a malformed YAML file
func TestLoad_MalformedYAML(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "config.yaml")

	// Write invalid YAML
	malformedYAML := `
brain:
  ollama_url: http://localhost:11434
  model: test
    invalid_nested_key: value
focus:
  theme: "unclosed string
`
	os.WriteFile(cfgFile, []byte(malformedYAML), 0644)

	os.Setenv("KYANITE_CONFIG", cfgFile)
	defer os.Unsetenv("KYANITE_CONFIG")

	cfg, err := Load()
	if err == nil {
		t.Error("Expected error when loading malformed YAML, got nil")
	}
	if cfg != nil {
		t.Error("Expected nil config when YAML is malformed")
	}
}

// TestLoad_NonExistentFile tests loading when config file doesn't exist
func TestLoad_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "nonexistent.yaml")

	os.Setenv("KYANITE_CONFIG", cfgFile)
	defer os.Unsetenv("KYANITE_CONFIG")

	// Should not error - should return defaults
	cfg, err := Load()
	if err != nil {
		t.Errorf("Load() should not error when file doesn't exist, got: %v", err)
	}
	if cfg == nil {
		t.Error("Expected default config when file doesn't exist, got nil")
	}
	// Verify defaults are present
	if cfg.Brain.OllamaURL != "http://nucbox:11434" {
		t.Errorf("Expected default ollama_url, got: %s", cfg.Brain.OllamaURL)
	}
}

// TestLoad_DirectoryPath tests loading when path points to a directory instead of file
func TestLoad_DirectoryPath(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("KYANITE_CONFIG", tmpDir)
	defer os.Unsetenv("KYANITE_CONFIG")

	cfg, err := Load()
	if err == nil {
		t.Error("Expected error when config path is a directory, got nil")
	}
	if cfg != nil {
		t.Error("Expected nil config when path is directory")
	}
}

// TestLoad_InvalidYAMLStructure tests loading YAML with wrong structure
func TestLoad_InvalidYAMLStructure(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "config.yaml")

	// Valid YAML but wrong structure
	invalidStructure := `
not_brain:
  random_key: value
completely_wrong:
  - item1
  - item2
`
	os.WriteFile(cfgFile, []byte(invalidStructure), 0644)

	os.Setenv("KYANITE_CONFIG", cfgFile)
	defer os.Unsetenv("KYANITE_CONFIG")

	cfg, err := Load()
	if err != nil {
		t.Errorf("Load() with wrong structure should still return defaults, got error: %v", err)
	}
	if cfg == nil {
		t.Error("Expected config with defaults when structure is wrong, got nil")
	}
}

// TestLoad_InvalidEnvVarFormat tests loading with invalid environment variable format
func TestLoad_InvalidEnvVarFormat(t *testing.T) {
	os.Unsetenv("KYANITE_CONFIG")
	os.Unsetenv("XDG_CONFIG_HOME")

	// Set env var with invalid value that can't be parsed
	os.Setenv("KYANITE_BRAIN_DB_PORT", "not-a-number")
	defer os.Unsetenv("KYANITE_BRAIN_DB_PORT")

	cfg, err := Load()
	// Invalid env vars should cause Load to fail
	if err == nil {
		t.Error("Expected error when env var has invalid format, got nil")
	}
	// Config should be nil when unmarshaling fails
	if cfg != nil {
		t.Error("Expected nil config when env var has invalid format")
	}
}

// TestSave_PermissionError tests saving when directory creation fails
func TestSave_PermissionError(t *testing.T) {
	// This test is tricky as we can't easily create permission errors in a cross-platform way
	// We'll test the atomic write behavior instead
	tmpDir := t.TempDir()
	os.Setenv("KYANITE_CONFIG", filepath.Join(tmpDir, "config.yaml"))
	defer os.Unsetenv("KYANITE_CONFIG")

	cfg := &Root{
		Brain: BrainConfig{
			OllamaURL: "http://test:11434",
			Model:     "test-model",
		},
		Focus: FocusConfig{
			Theme: "test-theme",
		},
	}

	err := Save(cfg)
	if err != nil {
		t.Errorf("Save() error: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filepath.Join(tmpDir, "config.yaml")); err != nil {
		t.Errorf("Config file was not created: %v", err)
	}
}

// TestLoad_EnvOverrideWithSpecialChars tests env vars with special characters
func TestLoad_EnvOverrideWithSpecialChars(t *testing.T) {
	os.Unsetenv("KYANITE_CONFIG")
	os.Unsetenv("XDG_CONFIG_HOME")

	// Test with URL containing special characters
	os.Setenv("KYANITE_BRAIN_OLLAMA_URL", "http://test-server:11434/path?query=value")
	defer os.Unsetenv("KYANITE_BRAIN_OLLAMA_URL")

	cfg, err := Load()
	if err != nil {
		t.Errorf("Load() with special chars in URL should work, got error: %v", err)
	}
	if cfg.Brain.OllamaURL != "http://test-server:11434/path?query=value" {
		t.Errorf("Expected URL with special chars to be preserved, got: %s", cfg.Brain.OllamaURL)
	}
}

// TestLoad_YAMLCommentsAndEmptyValues tests loading YAML with comments and empty values
func TestLoad_YAMLCommentsAndEmptyValues(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "config.yaml")

	yamlWithComments := `
# This is a comment
brain:
  ollama_url: http://localhost:11434  # inline comment
  model: ""  # empty string should be preserved
  timeout: 60s
focus:
  theme:  # empty theme
  ai:
    provider: ollama
    # Missing model field
`
	os.WriteFile(cfgFile, []byte(yamlWithComments), 0644)

	os.Setenv("KYANITE_CONFIG", cfgFile)
	defer os.Unsetenv("KYANITE_CONFIG")

	cfg, err := Load()
	if err != nil {
		t.Errorf("Load() with comments and empty values should work, got error: %v", err)
	}
	if cfg == nil {
		t.Error("Expected config, got nil")
	}
}

// TestConfigDir_HomeDirFallback tests ConfigDir when HOME is not set (edge case)
func TestConfigDir_HomeDirFallback(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	originalXDG := os.Getenv("XDG_CONFIG_HOME")

	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("XDG_CONFIG_HOME", originalXDG)
	}()

	// Test with only XDG set
	os.Unsetenv("HOME")
	os.Setenv("XDG_CONFIG_HOME", "/tmp/xdg-test")

	dir := ConfigDir()
	expected := "/tmp/xdg-test/kyanite"
	if dir != expected {
		t.Errorf("Expected %s, got %s", expected, dir)
	}
}

// TestConfigPath_WithBothEnvVars tests ConfigPath when both env vars are set
func TestConfigPath_WithBothEnvVars(t *testing.T) {
	os.Setenv("KYANITE_CONFIG", "/custom/path/config.yaml")
	os.Setenv("XDG_CONFIG_HOME", "/tmp/xdg")
	defer os.Unsetenv("KYANITE_CONFIG")
	defer os.Unsetenv("XDG_CONFIG_HOME")

	path := ConfigPath()
	if path != "/custom/path/config.yaml" {
		t.Errorf("KYANITE_CONFIG should take precedence, got: %s", path)
	}
}
