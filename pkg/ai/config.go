package ai

import (
	"os"
	"strconv"
	"time"

	"github.com/kyanite/config"
)

// Config holds the configuration for the inference brain.
//
// Defaults are designed for the kyanite.sh tailnet setup:
//   - LLM: Ollama on NUCBox over tailnet
//   - STT: local whisper.cpp on this machine
//   - Memory: PostgreSQL on NUCBox
type Config struct {
	// LLM configuration
	OllamaURL string // Ollama API endpoint (default: "http://nucbox:11434")
	Model     string // Model name (default: "gemma4:12b")
	Timeout   time.Duration

	// STT configuration
	WhisperBin   string // Path to whisper-cli or whisper-stream binary
	WhisperModel string // Path to GGML model file (e.g. ggml-base.en.bin)
	WhisperLang  string // Language code (default: "en")

	// Memory configuration
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	// App identity (used for memory namespacing)
	App string // "focus", "noise", "syntax", or "prism"
}

// DefaultConfig returns a config with sensible defaults for the kyanite.sh tailnet.
func DefaultConfig(app string) Config {
	return Config{
		OllamaURL: envOr("KYANITE_OLLAMA_URL", "http://nucbox:11434"),
		Model:     envOr("KYANITE_MODEL", "gemma4:12b"),
		Timeout:   60 * time.Second,

		WhisperBin:   envOr("KYANITE_WHISPER_BIN", "whisper-stream"),
		WhisperModel: envOr("KYANITE_WHISPER_MODEL", defaultWhisperModel()),
		WhisperLang:  envOr("KYANITE_WHISPER_LANG", "en"),

		DBHost:     envOr("KYANITE_DB_HOST", "nucbox"),
		DBPort:     envOrInt("KYANITE_DB_PORT", 5432),
		DBName:     envOr("KYANITE_DB_NAME", "kyanite"),
		DBUser:     envOr("KYANITE_DB_USER", "kyanite"),
		DBPassword: envOr("KYANITE_DB_PASSWORD", ""),
		DBSSLMode:  envOr("KYANITE_DB_SSLMODE", "disable"),

		App: app,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func defaultWhisperModel() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		home + "/.local/share/achiote-voice/models/ggml-large-v3-turbo.bin",
		home + "/.local/share/whisper/models/ggml-base.en.bin",
		"ggml-base.en.bin",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return "ggml-base.en.bin"
}

// ConfigFromRoot creates an ai.Config from the unified config Root.
// Falls back to DefaultConfig if root is nil (no config file loaded).
func ConfigFromRoot(root *config.Root, app string) Config {
	if root == nil {
		return DefaultConfig(app)
	}
	return Config{
		OllamaURL: root.Brain.OllamaURL,
		Model:     root.Brain.Model,
		Timeout:   root.Brain.Timeout,
		WhisperBin:   root.Brain.WhisperBin,
		WhisperModel: root.Brain.WhisperModel,
		WhisperLang:  envOr("KYANITE_WHISPER_LANG", "en"),
		DBHost:     root.Brain.DBHost,
		DBPort:     root.Brain.DBPort,
		DBName:     root.Brain.DBName,
		DBUser:     root.Brain.DBUser,
		DBPassword: root.Brain.DBPassword,
		DBSSLMode:  envOr("KYANITE_DB_SSLMODE", "disable"),
		App: app,
	}
}
