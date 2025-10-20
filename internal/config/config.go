package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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
	
	// Rapid prototyping settings
	Temperatures map[string]float64 `mapstructure:"temperatures"`
	RapidBrainstorm RapidBrainstormConfig `mapstructure:"rapid_brainstorm"`
	Continuation  ContinuationConfig  `mapstructure:"continuation"`
	Variation     VariationConfig     `mapstructure:"variation"`
}

// RapidBrainstormConfig contains settings for rapid brainstorming
type RapidBrainstormConfig struct {
	MaxAngles          int  `mapstructure:"max_angles"`
	GenerateFirstLine  bool `mapstructure:"generate_first_line"`
}

// ContinuationConfig contains settings for line continuation
type ContinuationConfig struct {
	Variations        int `mapstructure:"variations"`
	MaxContextLines   int `mapstructure:"max_context_lines"`
}

// VariationConfig contains settings for line variation
type VariationConfig struct {
	Variations        int `mapstructure:"variations"`
	DefaultConstraint string `mapstructure:"default_constraint"`
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
			Name:             "noise.sh",
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
			Database: "noise",
			Username: "postgres",
			Password: "",
			SSLMode:  "disable",
		},
		UI: UIConfig{
			Theme:           "violet-dusk",
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
			
			// Rapid prototyping settings
			Temperatures: map[string]float64{
				"sketch": 0.9,
				"draft":  0.7,
				"polish": 0.5,
			},
			RapidBrainstorm: RapidBrainstormConfig{
				MaxAngles:         3,
				GenerateFirstLine: true,
			},
			Continuation: ContinuationConfig{
				Variations:      3,
				MaxContextLines: 5,
			},
			Variation: VariationConfig{
				Variations:        3,
				DefaultConstraint: "more concrete",
			},
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
	v.SetConfigName("noise")

	// Enable reading from environment variables
	v.SetEnvPrefix("NOISE")
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

	appConfig := map[string]interface{}{
		"name":               c.App.Name,
		"version":            c.App.Version,
		"data_dir":           c.App.DataDir,
		"auto_save":          c.App.AutoSave,
		"auto_save_interval": c.App.AutoSaveInterval.String(),
		"max_recent_files":   c.App.MaxRecentFiles,
	}

	aiConfig := map[string]interface{}{
		"provider":        c.AI.Provider,
		"model":           c.AI.Model,
		"api_key":         c.AI.APIKey,
		"base_url":        c.AI.BaseURL,
		"temperature":     c.AI.Temperature,
		"max_tokens":      c.AI.MaxTokens,
		"timeout":         c.AI.Timeout.String(),
		"enabled":         c.AI.Enabled,
		"local_models":    c.AI.LocalModels,
		"custom_prompts":  c.AI.CustomPrompts,
		"temperatures":    c.AI.Temperatures,
		"rapid_brainstorm": map[string]interface{}{
			"max_angles":          c.AI.RapidBrainstorm.MaxAngles,
			"generate_first_line": c.AI.RapidBrainstorm.GenerateFirstLine,
		},
		"continuation": map[string]interface{}{
			"variations":        c.AI.Continuation.Variations,
			"max_context_lines": c.AI.Continuation.MaxContextLines,
		},
		"variation": map[string]interface{}{
			"variations":         c.AI.Variation.Variations,
			"default_constraint": c.AI.Variation.DefaultConstraint,
		},
	}

	audioConfig := map[string]interface{}{
		"enabled":           c.Audio.Enabled,
		"metronome_sound":   c.Audio.MetronomeSound,
		"chord_sample_rate": c.Audio.ChordSampleRate,
		"audio_buffer_size": c.Audio.AudioBufferSize,
		"midi_device":       c.Audio.MIDIDevice,
		"playback_gain":     c.Audio.PlaybackGain,
	}

	data := map[string]interface{}{
		"app": appConfig,
		"database": map[string]interface{}{
			"type":     c.Database.Type,
			"host":     c.Database.Host,
			"port":     c.Database.Port,
			"database": c.Database.Database,
			"username": c.Database.Username,
			"password": c.Database.Password,
			"ssl_mode": c.Database.SSLMode,
		},
		"ui": map[string]interface{}{
			"theme":            c.UI.Theme,
			"font_size":        c.UI.FontSize,
			"show_line_numbers": c.UI.ShowLineNumbers,
			"word_wrap":        c.UI.WordWrap,
			"animations":       c.UI.Animations,
		},
		"ai":    aiConfig,
		"audio": audioConfig,
		"dev": map[string]interface{}{
			"debug":         c.Dev.Debug,
			"log_level":     c.Dev.LogLevel,
			"profile":       c.Dev.Profile,
			"trace":         c.Dev.Trace,
			"mock_ai":       c.Dev.MockAI,
			"skip_database": c.Dev.SkipDatabase,
		},
	}

	configPath := filepath.Join(configDir, "noise.yaml")
	buf, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, buf, 0644); err != nil {
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
	if val := os.Getenv("NOISE_APP_DATA_DIR"); val != "" {
		c.App.DataDir = val
	}

	if val := os.Getenv("NOISE_APP_AUTO_SAVE"); val != "" {
		c.App.AutoSave = val == "true"
	}

	// Database settings
	if val := os.Getenv("NOISE_DATABASE_HOST"); val != "" {
		c.Database.Host = val
	}

	if val := os.Getenv("NOISE_DATABASE_PORT"); val != "" {
		// TODO: Parse port from environment variable
		_ = val // Suppress unused variable warning
	}

	// AI settings
	if val := os.Getenv("NOISE_AI_API_KEY"); val != "" {
		c.AI.APIKey = val
	}

	if val := os.Getenv("NOISE_AI_BASE_URL"); val != "" {
		c.AI.BaseURL = val
	}

	if val := os.Getenv("NOISE_AI_ENABLED"); val != "" {
		c.AI.Enabled = val == "true"
	}

	// Dev settings
	if val := os.Getenv("NOISE_DEV_DEBUG"); val != "" {
		c.Dev.Debug = val == "true"
	}

	if val := os.Getenv("NOISE_DEV_LOG_LEVEL"); val != "" {
		c.Dev.LogLevel = val
	}
}

// getDefaultDataDir returns the default data directory
func getDefaultDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./data"
	}
	return filepath.Join(homeDir, ".noise")
}

// getConfigDir returns the configuration directory
func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./config"
	}
	return filepath.Join(homeDir, ".config", "noise")
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
