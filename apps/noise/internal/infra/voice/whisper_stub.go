//go:build !cgo || !whisper

package voice

import (
	"errors"

	"github.com/kyanite/noise/internal/logging"
)

// Common errors for whisper engine
var (
	ErrModelNotLoaded   = errors.New("whisper model not loaded")
	ErrTranscribeFailed = errors.New("transcription failed")
	ErrInvalidSamples   = errors.New("invalid audio samples")
	ErrContextCreation  = errors.New("failed to create whisper context")
	ErrWhisperNotBuilt  = errors.New("whisper support not compiled - build with: CGO_ENABLED=1 go build -tags whisper")
)

// WhisperConfig configures the whisper engine
type WhisperConfig struct {
	ModelPath     string
	Language      string
	Translate     bool
	NoContext     bool
	SingleSegment bool
	PrintProgress bool
	Threads       int
	SpeedUp       bool
}

// DefaultWhisperConfig returns the default configuration
func DefaultWhisperConfig() WhisperConfig {
	return WhisperConfig{
		Language:  "en",
		NoContext: true,
	}
}

// Segment represents a transcribed segment with timing
type Segment struct {
	Text     string
	Start    int64
	End      int64
	NoSpeech bool
}

// TranscriptionResult contains the full transcription output
type TranscriptionResult struct {
	Text     string
	Segments []Segment
	Language string
}

// WhisperEngine is a stub when whisper.cpp is not available
type WhisperEngine struct {
	config WhisperConfig
	logger *logging.Logger
}

// NewWhisperEngine returns an error when whisper support is not compiled
func NewWhisperEngine(modelPath string, cfg WhisperConfig, logger *logging.Logger) (*WhisperEngine, error) {
	return nil, ErrWhisperNotBuilt
}

// LoadModel is not available without whisper support
func (we *WhisperEngine) LoadModel(modelPath string) error {
	return ErrWhisperNotBuilt
}

// IsLoaded always returns false without whisper support
func (we *WhisperEngine) IsLoaded() bool {
	return false
}

// Transcribe is not available without whisper support
func (we *WhisperEngine) Transcribe(samples []float32) (*TranscriptionResult, error) {
	return nil, ErrWhisperNotBuilt
}

// TranscribeSimple is not available without whisper support
func (we *WhisperEngine) TranscribeSimple(samples []float32) (string, error) {
	return "", ErrWhisperNotBuilt
}

// GetModelInfo returns information about the stub
func (we *WhisperEngine) GetModelInfo() map[string]interface{} {
	return map[string]interface{}{
		"loaded":  false,
		"message": "whisper support not compiled",
	}
}

// Close is a no-op for the stub
func (we *WhisperEngine) Close() error {
	return nil
}

// SupportedLanguages returns a list of supported language codes
func SupportedLanguages() []string {
	return []string{"en"}
}
