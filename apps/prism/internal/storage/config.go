package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents application configuration
type Config struct {
	Theme            string `toml:"theme"`
	AutoSave         bool   `toml:"auto_save"`
	AutoSaveInterval int    `toml:"auto_save_interval"`
}

// GetConfigDir returns the platform-specific config directory
func GetConfigDir() (string, error) {
	// Use os.UserConfigDir() with fallback to ~/.config/prism
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		if filepath.IsAbs(xdgConfigHome) {
			return filepath.Join(xdgConfigHome, "prism"), nil
		}
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".config", "prism"), nil
	}

	return filepath.Join(configDir, "prism"), nil
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return DefaultConfig(), err
	}

	configPath := filepath.Join(configDir, "config.toml")

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("failed to read config: %w", err)
	}

	// Parse TOML
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves configuration to file
func SaveConfig(cfg *Config) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.toml")

	// Marshal config to TOML
	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write atomically
	if err := AtomicWrite(configPath, data); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Theme:            "amber-night",
		AutoSave:         true,
		AutoSaveInterval: 30,
	}
}
