package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// Application settings
	App AppConfig `mapstructure:"app"`

	// Database settings
	Database DatabaseConfig `mapstructure:"database"`

	// UI settings
	UI UIConfig `mapstructure:"ui"`

	// AI settings
	AI AIConfig `mapstructure:"ai"`

	// Audio settings
	Audio AudioConfig `mapstructure:"audio"`

	// Development settings
	Dev DevConfig `mapstructure:"dev"`
}

// AppConfig contains general application settings
type AppConfig struct {
	Name             string        `mapstructure:"name"`
	Version          string        `mapstructure:"version"`
	DataDir          string        `mapstructure:"data_dir"`
	AutoSave         bool          `mapstructure:"auto_save"`
	AutoSaveInterval time.Duration `mapstructure:"auto_save_interval"`
	MaxRecentFiles   int           `mapstructure:"max_recent_files"`
}

// DatabaseConfig contains database settings
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// UIConfig contains UI settings
type UIConfig struct {
	Theme           string `mapstructure:"theme"`
	FontSize        int    `mapstructure:"font_size"`
	ShowLineNumbers bool   `mapstructure:"show_line_numbers"`
	WordWrap        bool   `mapstructure:"word_wrap"`
	Animations      bool   `mapstructure:"animations"`
}

// AIConfig contains AI service settings
type AIConfig struct {
	Provider      string            `mapstructure:"provider"`
	Model         string            `mapstructure:"model"`
	APIKey        string            `mapstructure:"api_key"`
	BaseURL       string            `mapstructure:"base_url"`
	Temperature   float64           `mapstructure:"temperature"`
	MaxTokens     int               `mapstructure:"max_tokens"`
	Timeout       time.Duration     `mapstructure:"timeout"`
	Enabled       bool              `mapstructure:"enabled"`
	LocalModels   []string          `mapstructure:"local_models"`
	CustomPrompts map[string]string `mapstructure:"custom_prompts"`
}

// AudioConfig contains audio settings
type AudioConfig struct {
	Enabled         bool    `mapstructure:"enabled"`
	MetronomeSound  string  `mapstructure:"metronome_sound"`
	ChordSampleRate int     `mapstructure:"chord_sample_rate"`
	AudioBufferSize int     `mapstructure:"audio_buffer_size"`
	MIDIDevice      string  `mapstructure:"midi_device"`
	PlaybackGain    float64 `mapstructure:"playback_gain"`
}

// DevConfig contains development settings
type DevConfig struct {
	Debug        bool   `mapstructure:"debug"`
	LogLevel     string `mapstructure:"log_level"`
	Profile      bool   `mapstructure:"profile"`
	Trace        bool   `mapstructure:"trace"`
	MockAI       bool   `mapstructure:"mock_ai"`
	SkipDatabase bool   `mapstructure:"skip_database"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:             "LyricForge",
			Version:          "1.0.0",
			DataDir:          getDefaultDataDir(),
			AutoSave:         true,
			AutoSaveInterval: 30 * time.Second,
			MaxRecentFiles:   10,
		},
		Database: DatabaseConfig{
			Type:     "sqlite",
			Host:     "localhost",
			Port:     5432,
			Database: "lyricforge",
			Username: "postgres",
			Password: "",
			SSLMode:  "disable",
		},
		UI: UIConfig{
			Theme:           "midnight_jazz",
			FontSize:        12,
			ShowLineNumbers: true,
			WordWrap:        true,
			Animations:      true,
		},
		AI: AIConfig{
			Provider:      "ollama",
			Model:         "qwen2.5:7b-instruct",
			APIKey:        "",
			BaseURL:       "http://localhost:11434",
			Temperature:   0.7,
			MaxTokens:     2048,
			Timeout:       30 * time.Second,
			Enabled:       true,
			LocalModels:   []string{"qwen2.5:7b-instruct", "llama3.1:8b"},
			CustomPrompts: make(map[string]string),
		},
		Audio: AudioConfig{
			Enabled:         true,
			MetronomeSound:  "sine",
			ChordSampleRate: 44100,
			AudioBufferSize: 1024,
			MIDIDevice:      "default",
			PlaybackGain:    0.8,
		},
		Dev: DevConfig{
			Debug:        false,
			LogLevel:     "info",
			Profile:      false,
			Trace:        false,
			MockAI:       false,
			SkipDatabase: false,
		},
	}
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	// Start with default configuration
	config := DefaultConfig()

	// Set up viper
	v := viper.New()

	// Set configuration file paths
	configPaths := []string{
		".",
		"./config",
		getConfigDir(),
	}

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	// Set configuration file names (without extension)
	v.SetConfigName("lyricforge")

	// Enable reading from environment variables
	v.SetEnvPrefix("LYRICFORGE")
	v.AutomaticEnv()

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		// If config file doesn't exist, that's okay - we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal configuration
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with environment variables if set
	config.overrideFromEnv()

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	// Create config directory if it doesn't exist
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set up viper for writing
	v := viper.New()

	// Set configuration values
	v.Set("app", c.App)
	v.Set("database", c.Database)
	v.Set("ui", c.UI)
	v.Set("ai", c.AI)
	v.Set("audio", c.Audio)
	v.Set("dev", c.Dev)

	// Write configuration file
	configPath := filepath.Join(configDir, "lyricforge.yaml")
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate app settings
	if c.App.Name == "" {
		return fmt.Errorf("app name cannot be empty")
	}

	if c.App.AutoSaveInterval < time.Second {
		return fmt.Errorf("auto save interval must be at least 1 second")
	}

	// Validate AI settings
	if c.AI.Enabled && c.AI.Provider == "" {
		return fmt.Errorf("AI provider must be specified when AI is enabled")
	}

	// Validate audio settings
	if c.Audio.Enabled && c.Audio.ChordSampleRate <= 0 {
		return fmt.Errorf("chord sample rate must be positive")
	}

	return nil
}

// overrideFromEnv overrides configuration values with environment variables
func (c *Config) overrideFromEnv() {
	// App settings
	if val := os.Getenv("LYRICFORGE_APP_DATA_DIR"); val != "" {
		c.App.DataDir = val
	}

	if val := os.Getenv("LYRICFORGE_APP_AUTO_SAVE"); val != "" {
		c.App.AutoSave = val == "true"
	}

	// Database settings
	if val := os.Getenv("LYRICFORGE_DATABASE_HOST"); val != "" {
		c.Database.Host = val
	}

	if val := os.Getenv("LYRICFORGE_DATABASE_PORT"); val != "" {
		// TODO: Parse port from environment variable
		_ = val // Suppress unused variable warning
	}

	// AI settings
	if val := os.Getenv("LYRICFORGE_AI_API_KEY"); val != "" {
		c.AI.APIKey = val
	}

	if val := os.Getenv("LYRICFORGE_AI_BASE_URL"); val != "" {
		c.AI.BaseURL = val
	}

	if val := os.Getenv("LYRICFORGE_AI_ENABLED"); val != "" {
		c.AI.Enabled = val == "true"
	}

	// Dev settings
	if val := os.Getenv("LYRICFORGE_DEV_DEBUG"); val != "" {
		c.Dev.Debug = val == "true"
	}

	if val := os.Getenv("LYRICFORGE_DEV_LOG_LEVEL"); val != "" {
		c.Dev.LogLevel = val
	}
}

// getDefaultDataDir returns the default data directory
func getDefaultDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./data"
	}
	return filepath.Join(homeDir, ".lyricforge")
}

// getConfigDir returns the configuration directory
func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./config"
	}
	return filepath.Join(homeDir, ".config", "lyricforge")
}

// GetDataDir returns the data directory (creates it if it doesn't exist)
func (c *Config) GetDataDir() string {
	if err := os.MkdirAll(c.App.DataDir, 0755); err != nil {
		// Fallback to current directory
		return "./data"
	}
	return c.App.DataDir
}

// IsDebug returns whether debug mode is enabled
func (c *Config) IsDebug() bool {
	return c.Dev.Debug
}

// IsDevMode returns whether development mode is enabled
func (c *Config) IsDevMode() bool {
	return c.Dev.Debug || c.Dev.Profile || c.Dev.Trace
}
