package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/config"
)

func TestDefaultConfigValues(t *testing.T) {
	c := config.DefaultConfig()
	if c == nil {
		t.Fatal("expected non-nil default config")
	}
	if c.App.Name != "noise.sh" {
		t.Fatalf("expected app name noise.sh, got %q", c.App.Name)
	}
	if c.App.AutoSave != true {
		t.Fatalf("expected autosave true")
	}
	if c.App.AutoSaveInterval < 1*time.Second {
		t.Fatalf("autosave interval too small")
	}
	if c.Database.Type != "sqlite" {
		t.Fatalf("expected sqlite DB type")
	}
	if c.UI.Theme == "" {
		t.Fatalf("expected default theme")
	}
	if !c.AI.Enabled {
		t.Fatalf("expected AI enabled by default")
	}
	if c.Audio.PlaybackGain <= 0 {
		t.Fatalf("expected positive playback gain")
	}
}

func TestValidateErrors(t *testing.T) {
	// Empty app name
	c := config.DefaultConfig()
	c.App.Name = ""
	if err := c.Validate(); err == nil {
		t.Fatal("expected validation error for empty app name")
	}

	// AutoSaveInterval too small
	c = config.DefaultConfig()
	c.App.AutoSaveInterval = 0
	if err := c.Validate(); err == nil {
		t.Fatal("expected validation error for autosave interval")
	}

	// AI enabled but no provider
	c = config.DefaultConfig()
	c.AI.Provider = ""
	c.AI.Enabled = true
	if err := c.Validate(); err == nil {
		t.Fatal("expected validation error for AI provider when enabled")
	}

	// Audio enabled but non-positive sample rate
	c = config.DefaultConfig()
	c.Audio.Enabled = true
	c.Audio.ChordSampleRate = 0
	if err := c.Validate(); err == nil {
		t.Fatal("expected validation error for chord sample rate")
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	// Ensure UserHomeDir points to tmp
	prevHome := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", prevHome)
	if err := os.Setenv("USERPROFILE", tmp); err != nil {
		t.Fatal(err)
	}

	cfg := config.DefaultConfig()
	cfg.App.Name = "roundtrip-test"
	cfg.AI.APIKey = "secret-key"
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	path := filepath.Join(tmp, ".config", "noise", "noise.yaml")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("expected config file at %s, read error: %v", path, err)
	}
	if !strings.Contains(string(data), "roundtrip-test") {
		t.Fatalf("written config does not contain updated name")
	}

	loaded, err := config.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.App.Name != "roundtrip-test" {
		t.Fatalf("expected loaded name roundtrip-test, got %q", loaded.App.Name)
	}
	if loaded.AI.APIKey != "secret-key" {
		t.Fatalf("expected AI API key preserved")
	}
}

func TestLoadWithEnvOverrides(t *testing.T) {
	tmp := t.TempDir()
	prevHome := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", prevHome)
	if err := os.Setenv("USERPROFILE", tmp); err != nil {
		t.Fatal(err)
	}

	// Set environment overrides
	os.Setenv("NOISE_APP_DATA_DIR", "/tmp/noise-data")
	os.Setenv("NOISE_AI_API_KEY", "env-key")
	os.Setenv("NOISE_AI_ENABLED", "false")
	os.Setenv("NOISE_DEV_DEBUG", "true")
	defer func() {
		os.Unsetenv("NOISE_APP_DATA_DIR")
		os.Unsetenv("NOISE_AI_API_KEY")
		os.Unsetenv("NOISE_AI_ENABLED")
		os.Unsetenv("NOISE_DEV_DEBUG")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load failed with env overrides: %v", err)
	}
	if cfg.App.DataDir != "/tmp/noise-data" {
		t.Fatalf("expected data dir overridden by env, got %q", cfg.App.DataDir)
	}
	if cfg.AI.APIKey != "env-key" {
		t.Fatalf("expected AI API key from env")
	}
	if cfg.AI.Enabled != false {
		t.Fatalf("expected AI enabled=false from env")
	}
	if !cfg.IsDebug() {
		t.Fatalf("expected debug true from env")
	}
}

func TestLoadInvalidConfigFileReturnsError(t *testing.T) {
	tmp := t.TempDir()
	prevHome := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", prevHome)
	if err := os.Setenv("USERPROFILE", tmp); err != nil {
		t.Fatal(err)
	}

	cfgDir := filepath.Join(tmp, ".config", "noise")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	badPath := filepath.Join(cfgDir, "noise.yaml")
	if err := ioutil.WriteFile(badPath, []byte("::: this is :::: not yaml :::: "), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected Load to fail for invalid config file")
	}
}

func TestGetDataDirCreatesDirectory(t *testing.T) {
	tmp := t.TempDir()
	cfg := config.DefaultConfig()
	cfg.App.DataDir = filepath.Join(tmp, "data-subdir")
	got := cfg.GetDataDir()
	if got != cfg.App.DataDir {
		t.Fatalf("expected GetDataDir to return %q, got %q", cfg.App.DataDir, got)
	}
	if fi, err := os.Stat(got); err != nil || !fi.IsDir() {
		t.Fatalf("expected data dir to be created at %s", got)
	}
}

func TestIsDebugAndIsDevMode(t *testing.T) {
	c := config.DefaultConfig()
	c.Dev.Debug = true
	if !c.IsDebug() {
		t.Fatal("expected IsDebug true when Debug set")
	}
	if !c.IsDevMode() {
		t.Fatal("expected IsDevMode true when Debug set")
	}

	c = config.DefaultConfig()
	c.Dev.Profile = true
	if !c.IsDevMode() {
		t.Fatal("expected IsDevMode true when Profile set")
	}

	c = config.DefaultConfig()
	c.Dev.Trace = true
	if !c.IsDevMode() {
		t.Fatal("expected IsDevMode true when Trace set")
	}
}

func TestSaveFailsWhenHomeIsFile(t *testing.T) {
	tmp := t.TempDir()
	// create a file and point USERPROFILE to it
	file := filepath.Join(tmp, "afile")
	if err := ioutil.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	prevHome := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", prevHome)
	if err := os.Setenv("USERPROFILE", file); err != nil {
		t.Fatal(err)
	}

	cfg := config.DefaultConfig()
	if err := cfg.Save(); err == nil {
		t.Fatal("expected Save to fail when home points to a file")
	}
}
