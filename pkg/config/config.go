package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Root is the top-level config for the kyanite.sh suite.
// Loaded from ~/.config/kyanite/config.yaml with env var overrides.
type Root struct {
	Brain  BrainConfig `koanf:"brain"`
	Focus  AppConfig   `koanf:"focus"`
	Noise  AppConfig   `koanf:"noise"`
	Syntax AppConfig   `koanf:"syntax"`
	Prism  AppConfig   `koanf:"prism"`
}

// BrainConfig holds inference brain settings.
type BrainConfig struct {
	OllamaURL   string        `koanf:"ollama_url"`
	Model       string        `koanf:"model"`
	Timeout     time.Duration `koanf:"timeout"`
	WhisperBin  string        `koanf:"whisper_bin"`
	WhisperModel string       `koanf:"whisper_model"`
	DBHost      string        `koanf:"db_host"`
	DBPort      int           `koanf:"db_port"`
	DBName      string        `koanf:"db_name"`
	DBUser      string        `koanf:"db_user"`
	DBPassword  string        `koanf:"db_password"`
}

// AppConfig holds per-app settings.
type AppConfig struct {
	Theme       string `koanf:"theme"`
	JournalDir  string `koanf:"journal_dir"`
	SamplesDir  string `koanf:"samples_dir"`
	StoriesDir  string `koanf:"stories_dir"`
	PalettesDir string `koanf:"palettes_dir"`
}

// defaults returns the built-in default configuration.
func defaults() *confmap.Confmap {
	return confmap.Provider(map[string]interface{}{
		"brain.ollama_url":    "http://nucbox:11434",
		"brain.model":         "gemma4:12b",
		"brain.timeout":       "60s",
		"brain.whisper_bin":   "whisper-stream",
		"brain.whisper_model": defaultWhisperModel(),
		"brain.db_host":       "nucbox",
		"brain.db_port":       5432,
		"brain.db_name":       "kyanite",
		"brain.db_user":       "kyanite",
		"brain.db_password":   "",
		"focus.theme":         "amber-night",
		"noise.theme":         "amber-night",
		"syntax.theme":        "amber-night",
		"prism.theme":         "amber-night",
	}, ".")
}

func defaultWhisperModel() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "achiote-voice", "models", "ggml-large-v3-turbo.bin")
}

// envKeyTransform converts KYANITE_BRAIN_OLLAMA_URL to brain.ollama_url.
// The koanf env provider passes the FULL env var name (including prefix) to this callback.
func envKeyTransform(s string) string {
	s = strings.TrimPrefix(s, "KYANITE_")
	s = strings.ToLower(s)
	parts := strings.SplitN(s, "_", 2)
	if len(parts) == 2 {
		return parts[0] + "." + parts[1]
	}
	return s
}

// ConfigDir returns the XDG-compliant config directory.
func ConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "kyanite")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "kyanite")
}

// ConfigPath returns the resolved config file path.
func ConfigPath() string {
	if p := os.Getenv("KYANITE_CONFIG"); p != "" {
		return p
	}
	return filepath.Join(ConfigDir(), "config.yaml")
}

// Load reads configuration from defaults -> config file -> env vars.
// Missing config file is not an error; env vars still apply.
func Load() (*Root, error) {
	k := koanf.New(".")

	// Layer 1: defaults
	if err := k.Load(defaults(), nil); err != nil {
		return nil, fmt.Errorf("load defaults: %w", err)
	}

	// Layer 2: config file (optional)
	cfgPath := ConfigPath()
	if _, err := os.Stat(cfgPath); err == nil {
		if err := k.Load(file.Provider(cfgPath), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("parse config file %s: %w", cfgPath, err)
		}
	}

	// Layer 3: env vars (KYANITE_ prefix, double underscore for nesting)
	if err := k.Load(env.Provider("KYANITE_", ".", envKeyTransform), nil); err != nil {
		return nil, fmt.Errorf("load env vars: %w", err)
	}

	var cfg Root
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Init creates a default config file at the standard location.
func Init() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	cfgPath := ConfigPath()
	if _, err := os.Stat(cfgPath); err == nil {
		return fmt.Errorf("config file already exists: %s", cfgPath)
	}

	content := `# kyanite.sh configuration
# Docs: https://github.com/simongonzalezdc/kyanite.sh#configuration

brain:
  ollama_url: http://nucbox:11434
  model: gemma4:12b
  timeout: 60s
  whisper_bin: whisper-stream
  whisper_model: ~/.local/share/achiote-voice/models/ggml-large-v3-turbo.bin
  db_host: nucbox
  db_port: 5432
  db_name: kyanite
  db_user: kyanite

focus:
  theme: amber-night

noise:
  theme: amber-night

syntax:
  theme: amber-night

prism:
  theme: amber-night
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Printf("Created config: %s\n", cfgPath)
	return nil
}

// Show prints the resolved configuration.
func Show() error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	fmt.Printf("Config file: %s\n\n", ConfigPath())
	fmt.Printf("Brain:\n")
	fmt.Printf("  ollama_url:    %s\n", cfg.Brain.OllamaURL)
	fmt.Printf("  model:         %s\n", cfg.Brain.Model)
	fmt.Printf("  timeout:       %s\n", cfg.Brain.Timeout)
	fmt.Printf("  whisper_bin:   %s\n", cfg.Brain.WhisperBin)
	fmt.Printf("  whisper_model: %s\n", cfg.Brain.WhisperModel)
	fmt.Printf("  db_host:       %s\n", cfg.Brain.DBHost)
	fmt.Printf("  db_port:       %d\n", cfg.Brain.DBPort)
	fmt.Printf("  db_name:       %s\n", cfg.Brain.DBName)
	fmt.Printf("  db_user:       %s\n", cfg.Brain.DBUser)
	fmt.Printf("\nApps:\n")
	for _, app := range []struct {
		name string
		cfg  AppConfig
	}{
		{"focus", cfg.Focus},
		{"noise", cfg.Noise},
		{"syntax", cfg.Syntax},
		{"prism", cfg.Prism},
	} {
		fmt.Printf("  %s: theme=%s\n", app.name, app.cfg.Theme)
	}
	return nil
}
