//go:build cgo && whisper

package voice

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Kyanite/noise/internal/logging"
	whisper "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

// Common errors for whisper engine
var (
	ErrModelNotLoaded    = errors.New("whisper model not loaded")
	ErrTranscribeFailed  = errors.New("transcription failed")
	ErrInvalidSamples    = errors.New("invalid audio samples")
	ErrContextCreation   = errors.New("failed to create whisper context")
)

// WhisperConfig configures the whisper engine
type WhisperConfig struct {
	ModelPath   string  // Path to the GGML model file
	Language    string  // Language code (e.g., "en", "auto")
	Translate   bool    // Translate to English
	NoContext   bool    // Don't use past transcription as context
	SingleSegment bool  // Force single segment output
	PrintProgress bool  // Print progress during transcription
	Threads     int     // Number of threads to use (0 = auto)
	SpeedUp     bool    // Speed up audio 2x for faster processing
}

// DefaultWhisperConfig returns the default configuration
func DefaultWhisperConfig() WhisperConfig {
	return WhisperConfig{
		Language:      "en",
		Translate:     false,
		NoContext:     true,  // Better for dictation
		SingleSegment: false,
		PrintProgress: false,
		Threads:       0,     // Auto-detect
		SpeedUp:       false,
	}
}

// Segment represents a transcribed segment with timing
type Segment struct {
	Text       string
	Start      int64  // Start time in milliseconds
	End        int64  // End time in milliseconds
	NoSpeech   bool   // Probability of no speech
}

// TranscriptionResult contains the full transcription output
type TranscriptionResult struct {
	Text     string
	Segments []Segment
	Language string
}

// WhisperEngine wraps whisper.cpp for speech recognition
type WhisperEngine struct {
	model   whisper.Model
	config  WhisperConfig
	logger  *logging.Logger
	mu      sync.Mutex
	loaded  bool
}

// NewWhisperEngine creates a new whisper engine with the given model
func NewWhisperEngine(modelPath string, cfg WhisperConfig, logger *logging.Logger) (*WhisperEngine, error) {
	if logger == nil {
		logger = logging.GetDefaultLogger()
	}

	cfg.ModelPath = modelPath

	we := &WhisperEngine{
		config: cfg,
		logger: logger,
	}

	// Load the model
	if err := we.LoadModel(modelPath); err != nil {
		return nil, err
	}

	return we, nil
}

// LoadModel loads a whisper model from the given path
func (we *WhisperEngine) LoadModel(modelPath string) error {
	we.mu.Lock()
	defer we.mu.Unlock()

	// Close existing model if any
	if we.model != nil {
		we.model.Close()
		we.model = nil
		we.loaded = false
	}

	we.logger.Infof("Loading whisper model from %s", modelPath)

	model, err := whisper.New(modelPath)
	if err != nil {
		return fmt.Errorf("failed to load whisper model: %w", err)
	}

	we.model = model
	we.loaded = true
	we.config.ModelPath = modelPath

	we.logger.Info("Whisper model loaded successfully")
	return nil
}

// IsLoaded returns whether a model is loaded
func (we *WhisperEngine) IsLoaded() bool {
	we.mu.Lock()
	defer we.mu.Unlock()
	return we.loaded
}

// Transcribe converts audio samples to text
func (we *WhisperEngine) Transcribe(samples []float32) (*TranscriptionResult, error) {
	we.mu.Lock()
	defer we.mu.Unlock()

	if !we.loaded || we.model == nil {
		return nil, ErrModelNotLoaded
	}

	if len(samples) == 0 {
		return nil, ErrInvalidSamples
	}

	// Create a new context for this transcription
	ctx, err := we.model.NewContext()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrContextCreation, err)
	}
	defer ctx.Close()

	// Configure the context
	if err := we.configureContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to configure context: %w", err)
	}

	// Process audio
	we.logger.Debugf("Processing %d audio samples", len(samples))
	if err := ctx.Process(samples, nil, nil); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTranscribeFailed, err)
	}

	// Extract results
	result := &TranscriptionResult{
		Segments: make([]Segment, 0),
	}

	var textBuilder strings.Builder
	
	// Iterate through segments
	for {
		segment, err := ctx.NextSegment()
		if err != nil {
			break // No more segments
		}

		segText := segment.Text
		textBuilder.WriteString(segText)

		result.Segments = append(result.Segments, Segment{
			Text:  segText,
			Start: segment.Start.Milliseconds(),
			End:   segment.End.Milliseconds(),
		})
	}

	result.Text = strings.TrimSpace(textBuilder.String())
	we.logger.Debugf("Transcription complete: %q", result.Text)

	return result, nil
}

// TranscribeSimple is a convenience method that returns just the text
func (we *WhisperEngine) TranscribeSimple(samples []float32) (string, error) {
	result, err := we.Transcribe(samples)
	if err != nil {
		return "", err
	}
	return result.Text, nil
}

// configureContext sets up the whisper context with configuration
func (we *WhisperEngine) configureContext(ctx whisper.Context) error {
	// Set language
	if we.config.Language != "" && we.config.Language != "auto" {
		if err := ctx.SetLanguage(we.config.Language); err != nil {
			we.logger.Warnf("Failed to set language %s: %v", we.config.Language, err)
		}
	}

	// Set translate mode
	ctx.SetTranslate(we.config.Translate)

	// Set thread count
	if we.config.Threads > 0 {
		ctx.SetThreads(uint(we.config.Threads))
	}

	// Set speed up
	ctx.SetSpeedup(we.config.SpeedUp)

	return nil
}

// GetModelInfo returns information about the loaded model
func (we *WhisperEngine) GetModelInfo() map[string]interface{} {
	we.mu.Lock()
	defer we.mu.Unlock()

	info := make(map[string]interface{})
	info["loaded"] = we.loaded
	info["model_path"] = we.config.ModelPath
	info["language"] = we.config.Language

	if we.loaded && we.model != nil {
		// Add model-specific info if available
		info["is_multilingual"] = we.model.IsMultilingual()
	}

	return info
}

// Close releases all resources
func (we *WhisperEngine) Close() error {
	we.mu.Lock()
	defer we.mu.Unlock()

	if we.model != nil {
		we.model.Close()
		we.model = nil
		we.loaded = false
		we.logger.Info("Whisper engine closed")
	}

	return nil
}

// SupportedLanguages returns a list of supported language codes
func SupportedLanguages() []string {
	return []string{
		"en", "zh", "de", "es", "ru", "ko", "fr", "ja", "pt", "tr",
		"pl", "ca", "nl", "ar", "sv", "it", "id", "hi", "fi", "vi",
		"he", "uk", "el", "ms", "cs", "ro", "da", "hu", "ta", "no",
		"th", "ur", "hr", "bg", "lt", "la", "mi", "ml", "cy", "sk",
		"te", "fa", "lv", "bn", "sr", "az", "sl", "kn", "et", "mk",
		"br", "eu", "is", "hy", "ne", "mn", "bs", "kk", "sq", "sw",
		"gl", "mr", "pa", "si", "km", "sn", "yo", "so", "af", "oc",
		"ka", "be", "tg", "sd", "gu", "am", "yi", "lo", "uz", "fo",
		"ht", "ps", "tk", "nn", "mt", "sa", "lb", "my", "bo", "tl",
		"mg", "as", "tt", "haw", "ln", "ha", "ba", "jw", "su",
	}
}
