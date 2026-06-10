package ai

import (
	"context"
	"errors"
)

// Brain is the unified inference brain for kyanite.sh apps.
//
// It provides three capabilities:
//   - LLM: text generation via Ollama on the NUCBox inference server
//   - STT: speech-to-text via local whisper.cpp
//   - Memory: persistent context via PostgreSQL on the NUCBox
//
// Apps create a Brain with their name and a Config, then call Generate,
// Transcribe, and memory methods through a single interface.
type Brain struct {
	llm    *LLMClient
	stt    *STTClient
	memory *MemoryStore
	app    string
}

// New creates a new Brain from the given config.
// Memory (PostgreSQL) and STT are optional — if unavailable, the Brain
// degrades gracefully (LLM always works if Ollama is reachable).
func New(cfg Config) (*Brain, error) {
	b := &Brain{app: cfg.App}

	// LLM is required
	b.llm = NewLLMClient(cfg.OllamaURL, cfg.Model, cfg.Timeout)

	// STT is optional
	stt := NewSTTClient(cfg.WhisperBin, cfg.WhisperModel, cfg.WhisperLang)
	if err := stt.IsAvailable(); err != nil {
		b.stt = nil // STT not available
	} else {
		b.stt = stt
	}

	// Memory is optional
	mem, err := NewMemoryStore(cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword, cfg.DBSSLMode)
	if err != nil {
		b.memory = nil // Memory not available
	} else {
		b.memory = mem
		// Best-effort schema init
		_ = mem.InitSchema(context.Background())
	}

	return b, nil
}

// ── LLM Methods ────────────────────────────────────────────────────

// Generate sends a prompt to the LLM and returns the response text.
func (b *Brain) Generate(ctx context.Context, prompt string, opts ...GenerateOption) (string, error) {
	if b.llm == nil {
		return "", ErrBrainNotInitialized
	}
	return b.llm.Generate(ctx, prompt, opts...)
}

// Chat sends a conversation to the LLM and returns the assistant's response.
func (b *Brain) Chat(ctx context.Context, messages []Message, opts ...GenerateOption) (string, error) {
	if b.llm == nil {
		return "", ErrBrainNotInitialized
	}
	return b.llm.Chat(ctx, messages, opts...)
}

// StreamChat sends a conversation and returns a channel of response chunks.
func (b *Brain) StreamChat(ctx context.Context, messages []Message, opts ...GenerateOption) (<-chan StreamChunk, error) {
	if b.llm == nil {
		return nil, ErrBrainNotInitialized
	}
	return b.llm.StreamChat(ctx, messages, opts...)
}

// IsLLMAvailable checks if the Ollama server is reachable.
func (b *Brain) IsLLMAvailable(ctx context.Context) bool {
	return b.llm != nil && b.llm.IsAvailable(ctx)
}

// ── STT Methods ─────────────────────────────────────────────────────

// TranscribeFile transcribes an audio file using local whisper.cpp.
func (b *Brain) TranscribeFile(ctx context.Context, audioPath string) (*Transcription, error) {
	if b.stt == nil {
		return nil, ErrSTTNotInstalled
	}
	return b.stt.TranscribeFile(ctx, audioPath)
}

// TranscribePCM transcribes raw PCM audio data.
func (b *Brain) TranscribePCM(ctx context.Context, pcmData []byte) (*Transcription, error) {
	if b.stt == nil {
		return nil, ErrSTTNotInstalled
	}
	return b.stt.TranscribePCM(ctx, pcmData)
}

// IsSTTAvailable checks if whisper.cpp is available.
func (b *Brain) IsSTTAvailable() bool {
	return b.stt != nil
}

// ── Memory Methods ──────────────────────────────────────────────────

// SaveContext stores a value in the app's session context.
func (b *Brain) SaveContext(ctx context.Context, sessionID, key string, value any) error {
	if b.memory == nil {
		return ErrMemoryUnreachable
	}
	return b.memory.SaveContext(ctx, b.app, sessionID, key, value)
}

// LoadContext retrieves a value from the app's session context.
func (b *Brain) LoadContext(ctx context.Context, sessionID, key string, dest any) error {
	if b.memory == nil {
		return ErrMemoryUnreachable
	}
	return b.memory.LoadContext(ctx, b.app, sessionID, key, dest)
}

// SaveMessage appends a message to the conversation history.
func (b *Brain) SaveMessage(ctx context.Context, sessionID, role, content string) error {
	if b.memory == nil {
		return nil // Silently skip if memory unavailable
	}
	return b.memory.SaveMessage(ctx, b.app, sessionID, role, content)
}

// LoadConversation loads recent conversation history.
func (b *Brain) LoadConversation(ctx context.Context, sessionID string, limit int) ([]Message, error) {
	if b.memory == nil {
		return nil, ErrMemoryUnreachable
	}
	return b.memory.LoadConversation(ctx, b.app, sessionID, limit)
}

// IsMemoryAvailable checks if the memory store is reachable.
func (b *Brain) IsMemoryAvailable(ctx context.Context) bool {
	return b.memory != nil && b.memory.IsAvailable(ctx)
}

// Close releases all resources.
func (b *Brain) Close() {
	if b.memory != nil {
		b.memory.Close()
	}
}

// App returns the app name this brain was created for.
func (b *Brain) App() string {
	return b.app
}

// IsError reports whether an error matches a specific brain error type.
// This is a convenience wrapper around errors.Is for callers that don't
// want to import the errors package.
func IsError(err, target error) bool {
	return errors.Is(err, target)
}

