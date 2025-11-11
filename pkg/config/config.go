package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	AI        AIConfig  `yaml:"ai" mapstructure:"ai"`
	Theme     string    `yaml:"theme" mapstructure:"theme"`
	Dashboard Dashboard `yaml:"dashboard" mapstructure:"dashboard"`
	Notes     Notes     `yaml:"notes" mapstructure:"notes"`
	UI        UI        `yaml:"ui" mapstructure:"ui"`
}

// AIConfig contains AI-related settings
type AIConfig struct {
	Provider    string  `yaml:"provider" mapstructure:"provider"`
	Model       string  `yaml:"model" mapstructure:"model"`
	Temperature float64 `yaml:"temperature" mapstructure:"temperature"`
	MaxTokens   int     `yaml:"max_tokens" mapstructure:"max_tokens"`
	Timeout     int     `yaml:"timeout" mapstructure:"timeout"` // in seconds
}

// Dashboard contains dashboard-specific settings
type Dashboard struct {
	AutoRefresh     bool `yaml:"auto_refresh" mapstructure:"auto_refresh"`
	RefreshInterval int  `yaml:"refresh_interval" mapstructure:"refresh_interval"` // in seconds
	ShowAnimation   bool `yaml:"show_animation" mapstructure:"show_animation"`
	CompactMode     bool `yaml:"compact_mode" mapstructure:"compact_mode"`
}

// Notes contains notes-related settings
type Notes struct {
	DefaultEditor string `yaml:"default_editor" mapstructure:"default_editor"`
	AutoSave      bool   `yaml:"auto_save" mapstructure:"auto_save"`
	SaveInterval  int    `yaml:"save_interval" mapstructure:"save_interval"` // in minutes
	Directory     string `yaml:"directory" mapstructure:"directory"`
}

// UI contains user interface settings
type UI struct {
	TimeFormat    string `yaml:"time_format" mapstructure:"time_format"` // 12h or 24h
	DateFormat    string `yaml:"date_format" mapstructure:"date_format"`
	ShowHelpTips  bool   `yaml:"show_help_tips" mapstructure:"show_help_tips"`
	Notifications bool   `yaml:"notifications" mapstructure:"notifications"`
	SoundEffects  bool   `yaml:"sound_effects" mapstructure:"sound_effects"`
}

var (
	globalConfigMutex sync.RWMutex
	globalConfig      *Config
	configLoaded      = false
)

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	globalConfigMutex.Lock()
	defer globalConfigMutex.Unlock()

	if configLoaded && globalConfig != nil {
		return globalConfig, nil
	}

	// Set config name and paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add config search paths
	homeDir, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(homeDir, ".focus"))
	}

	// Fallback to current directory
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// If file not found, create it with defaults and set viper to that file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createDefaultConfig(); err != nil {
				return nil, err
			}
			// Point viper to the created file and read again
			homeDir2, err := os.UserHomeDir()
			if err != nil {
				homeDir2 = "."
			}
			cfgPath := filepath.Join(homeDir2, ".focus", "config.yaml")
			viper.SetConfigFile(cfgPath)
			if err := viper.ReadInConfig(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Unmarshal config
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	globalConfig = config
	configLoaded = true

	return config, nil
}

// SaveConfig saves the current configuration to file
func SaveConfig(config *Config) error {
	viper.Set("ai", config.AI)
	viper.Set("theme", config.Theme)
	viper.Set("dashboard", config.Dashboard)
	viper.Set("notes", config.Notes)
	viper.Set("ui", config.UI)

	// Ensure config directory exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".focus")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "config.yaml")
	viper.SetConfigFile(configFile)
	if err := viper.WriteConfigAs(configFile); err != nil {
		return err
	}
	// Ensure file exists on disk
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		f, err2 := os.Create(configFile)
		if err2 != nil {
			return err2
		}
		_ = f.Close()
	}
	return nil
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	globalConfigMutex.RLock()
	loaded := configLoaded
	cfg := globalConfig
	globalConfigMutex.RUnlock()

	if !loaded {
		loadedCfg, err := LoadConfig()
		if err != nil {
			log.Printf("Warning: Failed to load config: %v", err)
			// Return a safe default config to prevent nil dereferences
			if loadedCfg == nil {
				return getDefaultConfig()
			}
		}
		return loadedCfg
	}
	return cfg
}

// getDefaultConfig returns a default configuration to prevent nil dereferences
func getDefaultConfig() *Config {
	return &Config{
		AI: AIConfig{
			Provider:    "ollama",
			Model:       "llama3",
			Temperature: 0.7,
			MaxTokens:   1000,
			Timeout:     30,
		},
		Theme: "amber-night",
		Dashboard: Dashboard{
			AutoRefresh:     true,
			RefreshInterval: 5,
			ShowAnimation:   true,
			CompactMode:     false,
		},
		Notes: Notes{
			DefaultEditor: "auto",
			AutoSave:      true,
			SaveInterval:  5,
			Directory:     "",
		},
		UI: UI{
			TimeFormat:    "12h",
			DateFormat:    "2006-01-02",
			ShowHelpTips:  true,
			Notifications: false,
			SoundEffects:  true,
		},
	}
}

// UpdateConfig updates specific configuration values
func UpdateConfig(updates map[string]any) error {
	config := GetConfig()

	// Apply updates
	for key, value := range updates {
		viper.Set(key, value)
	}

	// Reload and unmarshal
	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	globalConfig = config
	return SaveConfig(config)
}

// setDefaults sets default configuration values
func setDefaults() {
	// AI defaults
	viper.SetDefault("ai.provider", "ollama")
	viper.SetDefault("ai.model", "llama3")
	viper.SetDefault("ai.temperature", 0.7)
	viper.SetDefault("ai.max_tokens", 1000)
	viper.SetDefault("ai.timeout", 30)

	// Theme defaults (Kyanite master)
	viper.SetDefault("theme", "amber-night")

	// Dashboard defaults
	viper.SetDefault("dashboard.auto_refresh", true)
	viper.SetDefault("dashboard.refresh_interval", 5)
	viper.SetDefault("dashboard.show_animation", true)
	viper.SetDefault("dashboard.compact_mode", false)

	// Notes defaults
	viper.SetDefault("notes.default_editor", "auto")
	viper.SetDefault("notes.auto_save", true)
	viper.SetDefault("notes.save_interval", 5)
	viper.SetDefault("notes.directory", "")

	// UI defaults
	viper.SetDefault("ui.time_format", "12h")
	viper.SetDefault("ui.date_format", "2006-01-02")
	viper.SetDefault("ui.show_help_tips", true)
	viper.SetDefault("ui.notifications", false)
	viper.SetDefault("ui.sound_effects", true)
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".focus")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "config.yaml")
	return viper.WriteConfigAs(configFile)
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".focus", "config.yaml")
}

// IsConfigLoaded returns whether configuration has been loaded
func IsConfigLoaded() bool {
	return configLoaded
}

// ReloadConfig forces a reload of the configuration
func ReloadConfig() error {
	globalConfigMutex.Lock()
	configLoaded = false
	globalConfig = nil
	globalConfigMutex.Unlock()
	_, err := LoadConfig()
	return err
}
