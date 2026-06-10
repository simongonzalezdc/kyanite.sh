package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyanite/syntax/internal/ai"
	"gopkg.in/yaml.v3"
)

// AppConfig represents the application configuration
type AppConfig struct {
	AI ai.Config `yaml:"ai"`
}

// LoadConfig loads the application configuration
func LoadConfig() (*AppConfig, error) {
	configDir := GetConfigDir()
	configPath := filepath.Join(configDir, "config.yaml")

	// Return default config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &AppConfig{
			AI: ai.DefaultConfig(),
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the application configuration
func SaveConfig(config *AppConfig) error {
	configDir := GetConfigDir()
	if err := EnsureDir(configDir); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	return AtomicWriteFile(configPath, data, 0600)
}
