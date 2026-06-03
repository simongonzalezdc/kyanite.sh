package app

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/config"
	"github.com/Kyanite/noise/internal/infra/voice"
	"github.com/Kyanite/noise/internal/logging"
)

// VoiceService errors
var (
	ErrVoiceNotAvailable  = errors.New("voice-to-text is not available")
	ErrVoiceNotEnabled    = errors.New("voice-to-text is not enabled")
	ErrModelNotReady      = errors.New("whisper model not ready")
	ErrAlreadyDictating   = errors.New("already dictating")
	ErrNotDictating       = errors.New("not dictating")
	ErrDictationCancelled = errors.New("dictation cancelled")
)

// TranscriptionResult contains the result of a voice transcription
type TranscriptionResult struct {
	Text      string
	Duration  time.Duration
	Segments  []voice.Segment
	Cancelled bool
}

// VoiceState represents the current state of the voice service
type VoiceState int

const (
	VoiceStateIdle VoiceState = iota
	VoiceStateRecording
	VoiceStateProcessing
	VoiceStateError
)

// String returns the string representation of the voice state
func (s VoiceState) String() string {
	switch s {
	case VoiceStateIdle:
		return "idle"
	case VoiceStateRecording:
		return "recording"
	case VoiceStateProcessing:
		return "processing"
	case VoiceStateError:
		return "error"
	default:
		return "unknown"
	}
}

// VoiceService provides voice-to-text functionality
type VoiceService struct {
	capture      *voice.AudioCapture
	engine       *voice.WhisperEngine
	modelManager *voice.ModelManager
	config       *config.VoiceConfig
	logger       *logging.Logger

	// State
	state     VoiceState
	startTime time.Time
	lastError error
	mu        sync.Mutex

	// Event callbacks
	onStateChange   func(VoiceState)
	onTranscription func(TranscriptionResult)
	onLevelChange   func(float32)
	onProgress      func(int64, int64) // Model download progress
}

// NewVoiceService creates a new voice service
// The service auto-creates necessary directories and downloads models on first use
func NewVoiceService(cfg *config.Config, logger *logging.Logger) (*VoiceService, error) {
	if logger == nil {
		logger = logging.GetDefaultLogger()
	}

	if !cfg.Voice.Enabled {
		return nil, ErrVoiceNotEnabled
	}

	// Auto-create models directory in data dir
	modelsDir := cfg.GetDataDir() + "/models"
	modelManager, err := voice.NewModelManager(modelsDir, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create model manager: %w", err)
	}

	// Create audio capture
	captureConfig := voice.DefaultAudioCaptureConfig()
	captureConfig.Logger = logger
	capture, err := voice.NewAudioCapture(captureConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio capture: %w", err)
	}

	vs := &VoiceService{
		capture:      capture,
		modelManager: modelManager,
		config:       &cfg.Voice,
		logger:       logger,
		state:        VoiceStateIdle,
	}

	// Set up audio level callback
	capture.OnLevelChange(func(level float32) {
		vs.mu.Lock()
		callback := vs.onLevelChange
		vs.mu.Unlock()
		if callback != nil {
			callback(level)
		}
	})

	return vs, nil
}

// NewVoiceServiceWithAutoSetup creates a voice service that auto-downloads models if needed
// This is the recommended way to create a voice service for end-user apps
func NewVoiceServiceWithAutoSetup(cfg *config.Config, logger *logging.Logger, onProgress func(status string, progress float64)) (*VoiceService, error) {
	vs, err := NewVoiceService(cfg, logger)
	if err != nil {
		return nil, err
	}

	// Check if model needs to be downloaded
	modelName := cfg.Voice.Model
	if modelName == "" {
		modelName = voice.ModelBaseEN
	}

	if vs.modelManager.IsModelAvailable(modelName) {
		return vs, nil
	}

	if onProgress != nil {
		onProgress("Downloading voice model (first-time setup)...", 0)
	}

	// Download model with progress updates
	_, err = vs.modelManager.EnsureModel(modelName, func(downloaded, total int64) {
		if onProgress != nil && total > 0 {
			progress := float64(downloaded) / float64(total)
			onProgress(fmt.Sprintf("Downloading %s...", modelName), progress)
		}
	})

	if err != nil {
		if onProgress != nil {
			onProgress("Model download failed - voice will be unavailable", 0)
		}
		logger.Warnf("Failed to download voice model: %v", err)
		// Don't fail - voice just won't be available until model is downloaded
	} else {
		if onProgress != nil {
			onProgress("Voice model ready", 1.0)
		}
	}

	return vs, nil
}

// Initialize loads the whisper model (can be called asynchronously)
func (vs *VoiceService) Initialize() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if vs.engine != nil {
		return nil // Already initialized
	}

	modelName := vs.config.Model
	if modelName == "" {
		modelName = voice.ModelBaseEN // Default
	}

	// Check if model is available, download if not
	modelPath, err := vs.modelManager.EnsureModel(modelName, func(downloaded, total int64) {
		callback := vs.onProgress
		if callback != nil {
			callback(downloaded, total)
		}
	})
	if err != nil {
		vs.lastError = err
		return fmt.Errorf("failed to ensure model: %w", err)
	}

	// Create whisper engine
	whisperConfig := voice.DefaultWhisperConfig()
	whisperConfig.Language = vs.config.Language
	if whisperConfig.Language == "" {
		whisperConfig.Language = "en"
	}

	engine, err := voice.NewWhisperEngine(modelPath, whisperConfig, vs.logger)
	if err != nil {
		vs.lastError = err
		return fmt.Errorf("failed to create whisper engine: %w", err)
	}

	vs.engine = engine
	vs.logger.Info("Voice service initialized successfully")
	return nil
}

// IsAvailable returns whether voice-to-text is available
func (vs *VoiceService) IsAvailable() bool {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return vs.config.Enabled && vs.engine != nil && vs.engine.IsLoaded()
}

// IsModelReady returns whether the whisper model is loaded and ready
func (vs *VoiceService) IsModelReady() bool {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return vs.engine != nil && vs.engine.IsLoaded()
}

// GetState returns the current voice service state
func (vs *VoiceService) GetState() VoiceState {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return vs.state
}

// setState updates the state and notifies listeners
func (vs *VoiceService) setState(state VoiceState) {
	vs.state = state
	callback := vs.onStateChange
	if callback != nil {
		// Call outside of lock
		go callback(state)
	}
}

// StartDictation begins recording audio for transcription
func (vs *VoiceService) StartDictation() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if !vs.config.Enabled {
		return ErrVoiceNotEnabled
	}

	if vs.engine == nil || !vs.engine.IsLoaded() {
		return ErrModelNotReady
	}

	if vs.state == VoiceStateRecording {
		return ErrAlreadyDictating
	}

	// Start recording
	if err := vs.capture.StartRecording(); err != nil {
		vs.lastError = err
		vs.setState(VoiceStateError)
		return fmt.Errorf("failed to start recording: %w", err)
	}

	vs.startTime = time.Now()
	vs.setState(VoiceStateRecording)
	vs.logger.Debug("Dictation started")

	return nil
}

// StopDictation stops recording and returns the transcribed text
func (vs *VoiceService) StopDictation() (string, error) {
	vs.mu.Lock()

	if vs.state != VoiceStateRecording {
		vs.mu.Unlock()
		return "", ErrNotDictating
	}

	// Stop recording
	samples, err := vs.capture.StopRecording()
	if err != nil {
		vs.lastError = err
		vs.setState(VoiceStateError)
		vs.mu.Unlock()
		return "", fmt.Errorf("failed to stop recording: %w", err)
	}

	duration := time.Since(vs.startTime)
	vs.setState(VoiceStateProcessing)
	engine := vs.engine
	vs.mu.Unlock()

	vs.logger.Debugf("Processing %d samples (%.1fs of audio)", len(samples), duration.Seconds())

	// Transcribe
	result, err := engine.Transcribe(samples)
	if err != nil {
		vs.mu.Lock()
		vs.lastError = err
		vs.setState(VoiceStateError)
		vs.mu.Unlock()
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	vs.mu.Lock()
	vs.setState(VoiceStateIdle)
	callback := vs.onTranscription
	vs.mu.Unlock()

	// Notify listeners
	if callback != nil {
		callback(TranscriptionResult{
			Text:     result.Text,
			Duration: duration,
			Segments: result.Segments,
		})
	}

	vs.logger.Debugf("Transcription complete: %q", result.Text)
	return result.Text, nil
}

// CancelDictation stops recording without transcribing
func (vs *VoiceService) CancelDictation() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if vs.state != VoiceStateRecording {
		return ErrNotDictating
	}

	// Stop recording and discard samples
	_, err := vs.capture.StopRecording()
	if err != nil {
		vs.lastError = err
	}

	vs.setState(VoiceStateIdle)
	vs.logger.Debug("Dictation cancelled")

	// Notify listeners
	callback := vs.onTranscription
	if callback != nil {
		go callback(TranscriptionResult{Cancelled: true})
	}

	return nil
}

// GetDuration returns the current recording duration
func (vs *VoiceService) GetDuration() time.Duration {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if vs.state != VoiceStateRecording {
		return 0
	}
	return time.Since(vs.startTime)
}

// GetAudioLevel returns the current audio input level (0.0 to 1.0)
func (vs *VoiceService) GetAudioLevel() float32 {
	return vs.capture.PeakLevel()
}

// GetLastError returns the last error that occurred
func (vs *VoiceService) GetLastError() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return vs.lastError
}

// OnStateChange sets a callback for state changes
func (vs *VoiceService) OnStateChange(fn func(VoiceState)) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.onStateChange = fn
}

// OnTranscription sets a callback for completed transcriptions
func (vs *VoiceService) OnTranscription(fn func(TranscriptionResult)) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.onTranscription = fn
}

// OnLevelChange sets a callback for audio level changes
func (vs *VoiceService) OnLevelChange(fn func(float32)) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.onLevelChange = fn
}

// OnProgress sets a callback for model download progress
func (vs *VoiceService) OnProgress(fn func(downloaded, total int64)) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.onProgress = fn
}

// GetModelManager returns the model manager for UI access
func (vs *VoiceService) GetModelManager() *voice.ModelManager {
	return vs.modelManager
}

// GetConfig returns the voice configuration
func (vs *VoiceService) GetConfig() *config.VoiceConfig {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return vs.config
}

// ListMicrophones returns available microphone devices
func (vs *VoiceService) ListMicrophones() ([]string, error) {
	devices, err := vs.capture.ListDevices()
	if err != nil {
		return nil, err
	}

	names := make([]string, len(devices))
	for i, d := range devices {
		names[i] = d.Name()
	}
	return names, nil
}

// Close releases all resources
func (vs *VoiceService) Close() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	var errs []error

	if vs.capture != nil {
		if err := vs.capture.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if vs.engine != nil {
		if err := vs.engine.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing voice service: %v", errs)
	}

	vs.logger.Info("Voice service closed")
	return nil
}

// CheckMaxDuration returns whether the current recording has exceeded max duration
func (vs *VoiceService) CheckMaxDuration() bool {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if vs.state != VoiceStateRecording {
		return false
	}

	maxDuration := time.Duration(vs.config.MaxDuration) * time.Second
	if maxDuration <= 0 {
		maxDuration = 60 * time.Second // Default 60 seconds
	}

	return time.Since(vs.startTime) >= maxDuration
}
