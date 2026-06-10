// Package voice provides voice-to-text functionality using local speech recognition.
package voice

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/kyanite/noise/internal/logging"
	"github.com/gen2brain/malgo"
)

// Common errors for audio capture
var (
	ErrNoMicrophone     = errors.New("no microphone device found")
	ErrAlreadyRecording = errors.New("already recording")
	ErrNotRecording     = errors.New("not recording")
	ErrCaptureInit      = errors.New("failed to initialize audio capture")
	ErrDeviceStart      = errors.New("failed to start audio device")
)

// AudioFormat specifies the audio format for capture
type AudioFormat struct {
	SampleRate uint32
	Channels   uint32
	Format     malgo.FormatType
}

// DefaultAudioFormat returns the default audio format for Whisper (16kHz mono float32)
func DefaultAudioFormat() AudioFormat {
	return AudioFormat{
		SampleRate: 16000,           // Whisper expects 16kHz
		Channels:   1,               // Mono
		Format:     malgo.FormatF32, // 32-bit float
	}
}

// RingBuffer is a thread-safe circular buffer for audio samples
type RingBuffer struct {
	data     []float32
	size     int
	writePos int
	readPos  int
	count    int
	mu       sync.Mutex
}

// NewRingBuffer creates a new ring buffer with the specified capacity
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		data: make([]float32, capacity),
		size: capacity,
	}
}

// Write appends samples to the buffer
func (rb *RingBuffer) Write(samples []float32) int {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	written := 0
	for _, sample := range samples {
		if rb.count >= rb.size {
			// Buffer full, overwrite oldest data
			rb.readPos = (rb.readPos + 1) % rb.size
		} else {
			rb.count++
		}
		rb.data[rb.writePos] = sample
		rb.writePos = (rb.writePos + 1) % rb.size
		written++
	}
	return written
}

// Read reads all available samples from the buffer
func (rb *RingBuffer) Read() []float32 {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.count == 0 || rb.size == 0 {
		return nil
	}

	result := make([]float32, rb.count)
	for i := 0; i < rb.count; i++ {
		result[i] = rb.data[(rb.readPos+i)%rb.size]
	}

	// Clear buffer after reading
	rb.readPos = 0
	rb.writePos = 0
	rb.count = 0

	return result
}

// Count returns the number of samples in the buffer
func (rb *RingBuffer) Count() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.count
}

// Clear clears the buffer
func (rb *RingBuffer) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.readPos = 0
	rb.writePos = 0
	rb.count = 0
}

// AudioCapture handles microphone input capture
type AudioCapture struct {
	context   *malgo.AllocatedContext
	device    *malgo.Device
	buffer    *RingBuffer
	format    AudioFormat
	recording bool
	startTime time.Time
	peakLevel float32
	mu        sync.Mutex
	logger    *logging.Logger

	// Callbacks
	onLevelChange func(level float32)
}

// AudioCaptureConfig configures the audio capture
type AudioCaptureConfig struct {
	Format        AudioFormat
	BufferSeconds int // Buffer size in seconds of audio
	Logger        *logging.Logger
}

// DefaultAudioCaptureConfig returns the default configuration
func DefaultAudioCaptureConfig() AudioCaptureConfig {
	return AudioCaptureConfig{
		Format:        DefaultAudioFormat(),
		BufferSeconds: 60, // 60 seconds max recording
		Logger:        logging.GetDefaultLogger(),
	}
}

// NewAudioCapture creates a new audio capture instance
func NewAudioCapture(cfg AudioCaptureConfig) (*AudioCapture, error) {
	// Initialize miniaudio context
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCaptureInit, err)
	}

	// Calculate buffer size based on seconds of audio
	bufferSize := int(cfg.Format.SampleRate) * cfg.BufferSeconds * int(cfg.Format.Channels)

	ac := &AudioCapture{
		context: ctx,
		buffer:  NewRingBuffer(bufferSize),
		format:  cfg.Format,
		logger:  cfg.Logger,
	}

	return ac, nil
}

// ListDevices returns available capture (microphone) devices
func (ac *AudioCapture) ListDevices() ([]malgo.DeviceInfo, error) {
	infos, err := ac.context.Devices(malgo.Capture)
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate capture devices: %w", err)
	}
	return infos, nil
}

// StartRecording begins capturing audio from the microphone
func (ac *AudioCapture) StartRecording() error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if ac.recording {
		return ErrAlreadyRecording
	}

	// Clear any previous data
	ac.buffer.Clear()

	// Configure device
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = ac.format.Format
	deviceConfig.Capture.Channels = ac.format.Channels
	deviceConfig.SampleRate = ac.format.SampleRate
	deviceConfig.Alsa.NoMMap = 1

	// Data callback - called when audio data is available
	onRecvFrames := func(pOutputSamples, pInputSamples []byte, frameCount uint32) {
		// Convert bytes to float32 samples
		samples := bytesToFloat32(pInputSamples)

		// Write to buffer
		ac.buffer.Write(samples)

		// Calculate peak level for visual feedback
		var peak float32
		for _, s := range samples {
			if s < 0 {
				s = -s
			}
			if s > peak {
				peak = s
			}
		}

		ac.mu.Lock()
		ac.peakLevel = peak
		callback := ac.onLevelChange
		ac.mu.Unlock()

		if callback != nil {
			callback(peak)
		}
	}

	callbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}

	device, err := malgo.InitDevice(ac.context.Context, deviceConfig, callbacks)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCaptureInit, err)
	}

	if err := device.Start(); err != nil {
		device.Uninit()
		return fmt.Errorf("%w: %v", ErrDeviceStart, err)
	}

	ac.device = device
	ac.recording = true
	ac.startTime = time.Now()
	ac.logger.Infof("Started audio recording at %d Hz", ac.format.SampleRate)

	return nil
}

// StopRecording stops capturing and returns the recorded samples
func (ac *AudioCapture) StopRecording() ([]float32, error) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if !ac.recording {
		return nil, ErrNotRecording
	}

	// Stop and cleanup device
	if ac.device != nil {
		_ = ac.device.Stop()
		ac.device.Uninit()
		ac.device = nil
	}

	ac.recording = false
	duration := time.Since(ac.startTime)
	ac.logger.Infof("Stopped audio recording after %v", duration)

	// Return captured samples
	return ac.buffer.Read(), nil
}

// IsRecording returns whether recording is in progress
func (ac *AudioCapture) IsRecording() bool {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.recording
}

// Duration returns the duration of the current recording
func (ac *AudioCapture) Duration() time.Duration {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if !ac.recording {
		return 0
	}
	return time.Since(ac.startTime)
}

// PeakLevel returns the current peak audio level (0.0 to 1.0)
func (ac *AudioCapture) PeakLevel() float32 {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.peakLevel
}

// OnLevelChange sets a callback for audio level changes
func (ac *AudioCapture) OnLevelChange(fn func(level float32)) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.onLevelChange = fn
}

// Close releases all resources
func (ac *AudioCapture) Close() error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if ac.recording && ac.device != nil {
		_ = ac.device.Stop()
		ac.device.Uninit()
		ac.device = nil
		ac.recording = false
	}

	if ac.context != nil {
		if err := ac.context.Uninit(); err != nil {
			return fmt.Errorf("failed to uninitialize audio context: %w", err)
		}
		ac.context.Free()
		ac.context = nil
	}

	return nil
}

// bytesToFloat32 converts raw bytes to float32 samples
func bytesToFloat32(data []byte) []float32 {
	if len(data) == 0 {
		return nil
	}

	// Each float32 is 4 bytes
	numSamples := len(data) / 4
	samples := make([]float32, numSamples)

	for i := 0; i < numSamples; i++ {
		offset := i * 4
		// Little-endian byte order
		bits := uint32(data[offset]) |
			uint32(data[offset+1])<<8 |
			uint32(data[offset+2])<<16 |
			uint32(data[offset+3])<<24
		samples[i] = float32frombits(bits)
	}

	return samples
}

// float32frombits converts a uint32 bit pattern to float32
func float32frombits(bits uint32) float32 {
	return math.Float32frombits(bits)
}
