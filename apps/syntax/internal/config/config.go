package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Editor EditorConfig `yaml:"editor"`
	UI     UIConfig     `yaml:"ui"`
	AI     AIConfig     `yaml:"ai"`
}

// EditorConfig contains editor-specific settings
type EditorConfig struct {
	AutoSaveInterval  int  `yaml:"auto_save_interval"`  // seconds, 0 to disable
	AutoSaveIdleTime  int  `yaml:"auto_save_idle_time"` // seconds after last edit
	TabSize           int  `yaml:"tab_size"`
	ShowLineNumbers   bool `yaml:"show_line_numbers"`
	WordWrap          bool `yaml:"word_wrap"`
	VimMode           bool `yaml:"vim_mode"`
	MaxUndoStates     int  `yaml:"max_undo_states"`
	PaneSplit         int  `yaml:"pane_split"` // Left pane percentage (1-99)
}

// UIConfig contains UI preferences
type UIConfig struct {
	DefaultTheme      string `yaml:"default_theme"`
	ShowStats         bool   `yaml:"show_stats_in_statusbar"`
	ConfirmOnDelete   bool   `yaml:"confirm_on_delete"`
	RecentFilesCount  int    `yaml:"recent_files_count"`
}

// AIConfig contains AI assistant settings
type AIConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Model       string `yaml:"model"`
	APIEndpoint string `yaml:"api_endpoint"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int    `yaml:"max_tokens"`
	Timeout     int    `yaml:"timeout"` // seconds
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Editor: EditorConfig{
			AutoSaveInterval: 30,
			AutoSaveIdleTime: 3,
			TabSize:          2,
			ShowLineNumbers:  true,
			WordWrap:         true,
			VimMode:          true,
			MaxUndoStates:    100,
			PaneSplit:        50,
		},
		UI: UIConfig{
			DefaultTheme:     "monochrome",
			ShowStats:        true,
			ConfirmOnDelete:  true,
			RecentFilesCount: 10,
		},
		AI: AIConfig{
			Enabled:     false,
			Model:       "llama3.2:3b",
			APIEndpoint: "http://localhost:11434",
			Temperature: 0.7,
			MaxTokens:   500,
			Timeout:     30,
		},
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "syntax", "config.yaml")
}

// Load loads the configuration from disk, or returns default if not found
func Load() (Config, error) {
	configPath := GetConfigPath()

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), err
	}

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	// Validate and apply defaults for missing fields
	cfg = mergeWithDefaults(cfg)

	return cfg, nil
}

// Save saves the configuration to disk
func Save(cfg Config) error {
	configPath := GetConfigPath()
	configDir := filepath.Dir(configPath)

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(configPath, data, 0600)
}

// mergeWithDefaults fills in missing values with defaults
func mergeWithDefaults(cfg Config) Config {
	defaults := DefaultConfig()

	if cfg.Editor.AutoSaveInterval == 0 {
		cfg.Editor.AutoSaveInterval = defaults.Editor.AutoSaveInterval
	}
	if cfg.Editor.AutoSaveIdleTime == 0 {
		cfg.Editor.AutoSaveIdleTime = defaults.Editor.AutoSaveIdleTime
	}
	if cfg.Editor.TabSize == 0 {
		cfg.Editor.TabSize = defaults.Editor.TabSize
	}
	if cfg.Editor.MaxUndoStates == 0 {
		cfg.Editor.MaxUndoStates = defaults.Editor.MaxUndoStates
	}
	if cfg.Editor.PaneSplit == 0 {
		cfg.Editor.PaneSplit = defaults.Editor.PaneSplit
	}
	if cfg.UI.DefaultTheme == "" {
		cfg.UI.DefaultTheme = defaults.UI.DefaultTheme
	}
	if cfg.UI.RecentFilesCount == 0 {
		cfg.UI.RecentFilesCount = defaults.UI.RecentFilesCount
	}
	if cfg.AI.Model == "" {
		cfg.AI.Model = defaults.AI.Model
	}
	if cfg.AI.APIEndpoint == "" {
		cfg.AI.APIEndpoint = defaults.AI.APIEndpoint
	}
	if cfg.AI.Timeout == 0 {
		cfg.AI.Timeout = defaults.AI.Timeout
	}

	return cfg
}

// AutoSaveInterval returns the auto-save interval as a duration
func (c *EditorConfig) AutoSaveIntervalDuration() time.Duration {
	return time.Duration(c.AutoSaveInterval) * time.Second
}

// AutoSaveIdleTimeDuration returns the idle time before auto-save as a duration
func (c *EditorConfig) AutoSaveIdleTimeDuration() time.Duration {
	return time.Duration(c.AutoSaveIdleTime) * time.Second
}

// TimeoutDuration returns the AI timeout as a duration
func (c *AIConfig) TimeoutDuration() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}
