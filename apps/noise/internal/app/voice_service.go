package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kyanite/noise/internal/config"
	"github.com/kyanite/noise/internal/infra/brain"
	"github.com/kyanite/noise/internal/infra/voice"
	"github.com/kyanite/noise/internal/logging"
)

// VoiceService errors
var (
	ErrVoiceNotAvailable  = errors.New("voice-to-text is not available")
	ErrVoiceNotEnabled    = errors.New("voice-to-text is not enabled")
	ErrModelNotReady      = errors.New("STT not ready — brain STT unavailable")
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

// VoiceService provides voice-to-text functionality via brain STT.
type VoiceService struct {
	capture     *voice.AudioCapture
	brainClient *brain.Client
	config      *config.VoiceConfig
	logger      *logging.Logger

	// State
	state     VoiceState
	startTime time.Time
	lastError error
	mu        sync.Mutex

	// Event callbacks
	onStateChange   func(VoiceState)
	onTranscription func(TranscriptionResult)
	onLevelChange   func(float32)
	onProgress      func(int64, int64)
}

// NewVoiceService creates a new voice service backed by brain STT.
func NewVoiceService(cfg *config.Config, logger *logging.Logger) (*VoiceService, error) {
	if logger == nil {
		logger = logging.GetDefaultLogger()
	}

	if !cfg.Voice.Enabled {
		return nil, ErrVoiceNotEnabled
	}

	// Create audio capture
	captureConfig := voice.DefaultAudioCaptureConfig()
	captureConfig.Logger = logger
	capture, err := voice.NewAudioCapture(captureConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio capture: %w", err)
	}

	vs := &VoiceService{
		capture:     capture,
		brainClient: brain.NewClient(),
		config:      &cfg.Voice,
		logger:      logger,
		state:       VoiceStateIdle,
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

// NewVoiceServiceWithAutoSetup creates a voice service with brain-backed STT.
func NewVoiceServiceWithAutoSetup(cfg *config.Config, logger *logging.Logger, onProgress func(status string, progress float64)) (*VoiceService, error) {
	vs, err := NewVoiceService(cfg, logger)
	if err != nil {
		return nil, err
	}

	if onProgress != nil {
		if vs.brainClient != nil && vs.brainClient.IsSTTAvailable() {
			onProgress("Voice STT ready (via brain)", 1.0)
		} else {
			onProgress("Voice STT unavailable — brain STT not reachable", 0)
		}
	}

	return vs, nil
}

// Initialize prepares the voice service for transcription via brain STT.
func (vs *VoiceService) Initialize() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if vs.brainClient != nil && vs.brainClient.IsSTTAvailable() {
		vs.logger.Info("Voice service initialized (using brain STT)")
		return nil
	}

	vs.lastError = ErrModelNotReady
	return ErrModelNotReady
}

// IsAvailable returns whether voice-to-text is available
func (vs *VoiceService) IsAvailable() bool {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	if !vs.config.Enabled {
		return false
	}
	return vs.brainClient.IsSTTAvailable()
}

// IsModelReady returns whether brain STT is ready
func (vs *VoiceService) IsModelReady() bool {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return vs.brainClient.IsSTTAvailable()
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

	if !vs.brainClient.IsSTTAvailable() {
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
	brainClient := vs.brainClient
	vs.mu.Unlock()

	vs.logger.Debugf("Processing %d samples (%.1fs of audio)", len(samples), duration.Seconds())

	// Transcribe via brain STT
	if brainClient == nil || !brainClient.IsSTTAvailable() {
		vs.mu.Lock()
		vs.setState(VoiceStateError)
		vs.mu.Unlock()
		return "", ErrModelNotReady
	}

	transcribed, err := brainClient.TranscribePCM(context.Background(), samples)
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
			Text:     transcribed,
			Duration: duration,
		})
	}

	vs.logger.Debugf("Transcription complete: %q", transcribed)
	return transcribed, nil
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

	if vs.brainClient != nil {
		vs.brainClient.Close()
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
