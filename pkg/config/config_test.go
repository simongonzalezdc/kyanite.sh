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
