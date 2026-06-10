// Package brain wraps github.com/kyanite/ai Brain as the unified AI client
// for the noise app, providing both LLM and STT through a single interface.
package brain

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	ai "github.com/kyanite/ai"
	"github.com/kyanite/noise/internal/logging"
)

// Client wraps a pkg/ai Brain and exposes LLM + STT methods for the noise app.
type Client struct {
	brain *ai.Brain
	cfg   ai.Config
}

// NewClient creates a new Brain-backed client using ai.DefaultConfig("noise").
// If Brain creation fails (e.g. NUCBox unreachable), the client still initializes
// but LLM/STT calls will return errors — the app remains usable offline.
func NewClient() *Client {
	cfg := ai.DefaultConfig("noise")
	brain, err := ai.New(cfg)
	if err != nil {
		// Log but don't fail — app must work offline.
		logging.GetDefaultLogger().Warnf("brain init failed (offline mode): %v", err)
	}
	return &Client{
		brain: brain,
		cfg:   cfg,
	}
}

// NewClientWithConfig creates a new Brain-backed client with explicit config.
func NewClientWithConfig(cfg ai.Config) *Client {
	brain, err := ai.New(cfg)
	if err != nil {
		logging.GetDefaultLogger().Warnf("brain init failed (offline mode): %v", err)
	}
	return &Client{
		brain: brain,
		cfg:   cfg,
	}
}

// Brain returns the underlying pkg/ai Brain for direct access.
func (c *Client) Brain() *ai.Brain {
	return c.brain
}

// Generate sends a prompt to the LLM and returns the completion.
// This satisfies the app/ai.QuickLLMClient interface (model and options are
// ignored — the Brain uses its configured model).
func (c *Client) Generate(ctx context.Context, model, prompt string, options map[string]any) (string, error) {
	if c.brain == nil {
		return "", fmt.Errorf("%w: noise client", ai.ErrBrainNotInitialized)
	}
	return c.brain.Generate(ctx, prompt)
}

// GenerateWithOptions sends a prompt with fine-grained options.
func (c *Client) GenerateWithOptions(ctx context.Context, prompt string, opts ...ai.GenerateOption) (string, error) {
	if c.brain == nil {
		return "", fmt.Errorf("%w: noise client", ai.ErrBrainNotInitialized)
	}
	return c.brain.Generate(ctx, prompt, opts...)
}

// Chat sends a conversation to the LLM and returns the assistant response.
func (c *Client) Chat(ctx context.Context, messages []ai.Message, opts ...ai.GenerateOption) (string, error) {
	if c.brain == nil {
		return "", fmt.Errorf("%w: noise client", ai.ErrBrainNotInitialized)
	}
	return c.brain.Chat(ctx, messages, opts...)
}

// StreamGenerate sends a prompt and streams the response over the returned channel.
// This uses real streaming via Brain.StreamChat when available.
//
// Error reporting:
//   - If brain is nil, returns (nil, ErrBrainNotInitialized).
//   - If StreamChat fails to start, returns (nil, error).
//   - If a mid-stream error occurs, the channel is closed early without
//     sending any data. Callers can detect this by checking whether the
//     channel closed without producing meaningful output.
func (c *Client) StreamGenerate(ctx context.Context, model, prompt string) (<-chan string, error) {
	if c.brain == nil {
		return nil, fmt.Errorf("%w: noise client", ai.ErrBrainNotInitialized)
	}

	messages := []ai.Message{{Role: "user", Content: prompt}}
	stream, err := c.brain.StreamChat(ctx, messages)
	if err != nil {
		return nil, err
	}

	ch := make(chan string, 8)
	go func() {
		defer close(ch)
		for chunk := range stream {
			if chunk.Error != nil {
				// Log the error but don't send it as data — close cleanly.
				logging.GetDefaultLogger().Warnf("stream error: %v", chunk.Error)
				return
			}
			if chunk.Content != "" {
				ch <- chunk.Content
			}
			if chunk.Done {
				return
			}
		}
	}()
	return ch, nil
}

// TranscribePCM transcribes raw float32 audio samples via Brain.STT.
// The samples are converted to 16-bit PCM bytes as required by whisper.cpp.
func (c *Client) TranscribePCM(ctx context.Context, samples []float32) (string, error) {
	if c.brain == nil {
		return "", fmt.Errorf("%w: noise client", ai.ErrBrainNotInitialized)
	}
	if !c.brain.IsSTTAvailable() {
		return "", fmt.Errorf("STT not available: whisper.cpp not found")
	}

	pcmData := float32ToPCM16(samples)
	xfer, err := c.brain.TranscribePCM(ctx, pcmData)
	if err != nil {
		return "", err
	}
	return xfer.Text, nil
}

// TranscribeFile transcribes a WAV file via Brain.STT.
func (c *Client) TranscribeFile(ctx context.Context, audioPath string) (string, error) {
	if c.brain == nil {
		return "", fmt.Errorf("%w: noise client", ai.ErrBrainNotInitialized)
	}
	xfer, err := c.brain.TranscribeFile(ctx, audioPath)
	if err != nil {
		return "", err
	}
	return xfer.Text, nil
}

// IsAvailable checks if the LLM (Ollama) is reachable.
func (c *Client) IsAvailable(ctx context.Context) bool {
	if c.brain == nil {
		return false
	}
	return c.brain.IsLLMAvailable(ctx)
}

// IsSTTAvailable checks if STT (whisper.cpp) is available.
func (c *Client) IsSTTAvailable() bool {
	if c.brain == nil {
		return false
	}
	return c.brain.IsSTTAvailable()
}

// Close releases Brain resources.
func (c *Client) Close() {
	if c.brain != nil {
		c.brain.Close()
	}
}

// GetConfig returns the config used to create this client.
func (c *Client) GetConfig() ai.Config {
	return c.cfg
}

// GetModelStatus returns status information about the AI backend.
func (c *Client) GetModelStatus() map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return map[string]interface{}{
		"provider":      "brain",
		"ollama_url":    c.cfg.OllamaURL,
		"model":         c.cfg.Model,
		"ollama_status": c.IsAvailable(ctx),
		"stt_available": c.IsSTTAvailable(),
		"last_check":    time.Now(),
	}
}

// float32ToPCM16 converts float32 audio samples to 16-bit PCM bytes.
func float32ToPCM16(samples []float32) []byte {
	buf := make([]byte, len(samples)*2)
	for i, s := range samples {
		val := int16(s * 32767.0)
		if s > 1.0 {
			val = 32767
		} else if s < -1.0 {
			val = -32768
		}
		binary.LittleEndian.PutUint16(buf[i*2:], uint16(val))
	}
	return buf
}
