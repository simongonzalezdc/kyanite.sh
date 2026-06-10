package ai

import (
	"context"
	"errors"
	"testing"
	"time"
)

// ── Config Tests ─────────────────────────────────────────────────────

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig("focus")
	if cfg.App != "focus" {
		t.Errorf("expected app 'focus', got %q", cfg.App)
	}
	if cfg.OllamaURL != "http://nucbox:11434" {
		t.Errorf("expected default ollama URL, got %q", cfg.OllamaURL)
	}
	if cfg.Model != "gemma4:12b" {
		t.Errorf("expected default model, got %q", cfg.Model)
	}
	if cfg.DBHost != "nucbox" {
		t.Errorf("expected default db host, got %q", cfg.DBHost)
	}
	if cfg.DBPort != 5432 {
		t.Errorf("expected default db port 5432, got %d", cfg.DBPort)
	}
}

func TestDefaultConfigEachApp(t *testing.T) {
	for _, app := range []string{"focus", "noise", "syntax", "prism"} {
		cfg := DefaultConfig(app)
		if cfg.App != app {
			t.Errorf("expected app %q, got %q", app, cfg.App)
		}
	}
}

// ── GenerateOptions Tests ────────────────────────────────────────────

func TestGenerateOptions(t *testing.T) {
	opts := &GenerateOptions{}
	WithTemperature(0.5)(opts)
	WithMaxTokens(2048)(opts)
	WithSystemPrompt("test prompt")(opts)
	WithJSONMode()(opts)
	WithStopWords("END", "STOP")(opts)
	WithTimeout(30 * time.Second)(opts)

	if opts.Temperature != 0.5 {
		t.Errorf("expected temperature 0.5, got %f", opts.Temperature)
	}
	if opts.MaxTokens != 2048 {
		t.Errorf("expected max tokens 2048, got %d", opts.MaxTokens)
	}
	if opts.SystemPrompt != "test prompt" {
		t.Errorf("expected system prompt, got %q", opts.SystemPrompt)
	}
	if !opts.JSONMode {
		t.Error("expected JSON mode")
	}
	if len(opts.StopWords) != 2 {
		t.Errorf("expected 2 stop words, got %d", len(opts.StopWords))
	}
	if opts.Timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %v", opts.Timeout)
	}
}

// ── Brain Lifecycle Tests ────────────────────────────────────────────

func TestNewBrainWithoutServices(t *testing.T) {
	cfg := Config{
		OllamaURL:   "http://localhost:19999",
		Model:       "test",
		Timeout:     2 * time.Second,
		WhisperBin:  "nonexistent-whisper-binary",
		DBHost:      "nonexistent-host",
		DBPort:      5432,
		DBName:      "nonexistent",
		DBUser:      "nonexistent",
		DBPassword:  "nonexistent",
		App:         "test",
	}

	brain, err := New(cfg)
	if err != nil {
		t.Fatalf("New() should not fail even with unreachable services: %v", err)
	}
	defer brain.Close()

	if brain.IsLLMAvailable(context.Background()) {
		t.Error("LLM should not be available with unreachable URL")
	}
	if brain.IsSTTAvailable() {
		t.Error("STT should not be available with nonexistent binary")
	}
	if brain.IsMemoryAvailable(context.Background()) {
		t.Error("Memory should not be available with nonexistent host")
	}
	if brain.App() != "test" {
		t.Errorf("expected app 'test', got %q", brain.App())
	}
}

// ── Typed Error Tests ────────────────────────────────────────────────

func TestBrainGenerateReturnsTypedErrors(t *testing.T) {
	cfg := Config{
		OllamaURL:  "http://localhost:19999",
		Model:      "test",
		Timeout:    1 * time.Second,
		WhisperBin: "nonexistent",
		DBHost:     "nonexistent",
		App:        "test",
	}

	brain, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer brain.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = brain.Generate(ctx, "test prompt")
	if err == nil {
		t.Fatal("Generate should fail with unreachable LLM")
	}
	if !errors.Is(err, ErrLLMUnavailable) && !errors.Is(err, ErrLLMRequestFailed) {
		t.Errorf("expected ErrLLMUnavailable or ErrLLMRequestFailed, got: %v", err)
	}
}

func TestBrainSTTReturnsTypedErrors(t *testing.T) {
	cfg := Config{
		OllamaURL:  "http://localhost:19999",
		Model:      "test",
		Timeout:    1 * time.Second,
		WhisperBin: "nonexistent",
		DBHost:     "nonexistent",
		App:        "test",
	}

	brain, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer brain.Close()

	ctx := context.Background()

	_, err = brain.TranscribeFile(ctx, "nonexistent.wav")
	if !errors.Is(err, ErrSTTNotInstalled) {
		t.Errorf("expected ErrSTTNotInstalled, got: %v", err)
	}

	_, err = brain.TranscribePCM(ctx, []byte{})
	if !errors.Is(err, ErrSTTNotInstalled) {
		t.Errorf("expected ErrSTTNotInstalled, got: %v", err)
	}
}

func TestBrainMemoryReturnsTypedErrors(t *testing.T) {
	cfg := Config{
		OllamaURL:  "http://localhost:19999",
		Model:      "test",
		Timeout:    1 * time.Second,
		WhisperBin: "nonexistent",
		DBHost:     "nonexistent",
		App:        "test",
	}

	brain, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer brain.Close()

	ctx := context.Background()

	err = brain.SaveContext(ctx, "session1", "key1", "value1")
	if !errors.Is(err, ErrMemoryUnreachable) {
		t.Errorf("expected ErrMemoryUnreachable, got: %v", err)
	}

	var result string
	err = brain.LoadContext(ctx, "session1", "key1", &result)
	if !errors.Is(err, ErrMemoryUnreachable) {
		t.Errorf("expected ErrMemoryUnreachable, got: %v", err)
	}
}

func TestBrainNotInitializedErrors(t *testing.T) {
	// Create a brain with nil llm by directly constructing
	brain := &Brain{app: "test"}

	ctx := context.Background()

	_, err := brain.Generate(ctx, "test")
	if !errors.Is(err, ErrBrainNotInitialized) {
		t.Errorf("expected ErrBrainNotInitialized, got: %v", err)
	}

	_, err = brain.Chat(ctx, nil)
	if !errors.Is(err, ErrBrainNotInitialized) {
		t.Errorf("expected ErrBrainNotInitialized, got: %v", err)
	}

	_, err = brain.StreamChat(ctx, nil)
	if !errors.Is(err, ErrBrainNotInitialized) {
		t.Errorf("expected ErrBrainNotInitialized, got: %v", err)
	}
}

// ── Prompt Tests ─────────────────────────────────────────────────────


func TestSyntaxSuggestPrompt(t *testing.T) {
	for _, typ := range []string{"continue", "improve", "dialogue", "description", "character", "unknown"} {
		p := SyntaxSuggestPrompt(typ, "The dark forest loomed.", "A scary scene")
		if p == "" {
			t.Errorf("prompt for type %q should not be empty", typ)
		}
	}
}

func TestPrismPalettePrompt(t *testing.T) {
	p := PrismPalettePrompt("sunset over the ocean")
	if p == "" {
		t.Error("prompt should not be empty")
	}
}

func TestPrismContrastPrompt(t *testing.T) {
	p := PrismContrastPrompt("#ffffff", "#000000", 21.0)
	if p == "" {
		t.Error("prompt should not be empty for passing contrast")
	}
	p = PrismContrastPrompt("#777777", "#888888", 1.5)
	if p == "" {
		t.Error("prompt should not be empty for failing contrast")
	}
}

// ── LLM Client Unit Tests ────────────────────────────────────────────

func TestLLMClientCreation(t *testing.T) {
	client := NewLLMClient("http://localhost:11434", "gemma4:12b", 30*time.Second)
	if client == nil {
		t.Fatal("client should not be nil")
	}
	if client.baseURL != "http://localhost:11434" {
		t.Errorf("expected baseURL, got %q", client.baseURL)
	}
	if client.model != "gemma4:12b" {
		t.Errorf("expected model, got %q", client.model)
	}
}

func TestLLMHealthCache(t *testing.T) {
	client := NewLLMClient("http://localhost:19999", "test", 1*time.Second)
	ctx := context.Background()

	// First call should check
	avail := client.IsAvailable(ctx)
	if avail {
		t.Error("should not be available")
	}

	// Second call should return cached result (same false)
	avail = client.IsAvailable(ctx)
	if avail {
		t.Error("cached result should also be false")
	}
}

// ── STT Client Tests ─────────────────────────────────────────────────

func TestSTTClientCreation(t *testing.T) {
	client := NewSTTClient("whisper-cli", "/path/to/model.bin", "en")
	if client == nil {
		t.Fatal("client should not be nil")
	}
	if client.bin != "whisper-cli" {
		t.Errorf("expected bin, got %q", client.bin)
	}
}

// ── Whisper Output Parsing ──────────────────────────────────────────

func TestParseWhisperOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"plain text", "Hello world", "Hello world"},
		{"multiline", "Hello\nWorld", "Hello World"},
		{"with timestamps", "[00:00:00.000 --> 00:00:05.000]  Hello world", "Hello world"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseWhisperOutput(tt.input)
			if got != tt.expected {
				t.Errorf("parseWhisperOutput(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// ── Message Building ────────────────────────────────────────────────

func TestBuildMessages(t *testing.T) {
	msgs := buildMessages("hello", &GenerateOptions{})
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if msgs[0].Role != "user" || msgs[0].Content != "hello" {
		t.Errorf("unexpected message: %+v", msgs[0])
	}

	msgs = buildMessages("hello", &GenerateOptions{SystemPrompt: "system"})
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Role != "system" {
		t.Error("first message should be system")
	}
}

// ── Helper Tests ────────────────────────────────────────────────────

func TestEnvOr(t *testing.T) {
	got := envOr("NONEXISTENT_VAR_12345", "fallback")
	if got != "fallback" {
		t.Errorf("expected fallback, got %q", got)
	}
}

func TestDefaultWhisperModel(t *testing.T) {
	model := defaultWhisperModel()
	if model == "" {
		t.Error("should return a non-empty default")
	}
}

func TestIsError(t *testing.T) {
	if !IsError(ErrLLMUnavailable, ErrLLMUnavailable) {
		t.Error("IsError should match same error")
	}
	if IsError(ErrLLMUnavailable, ErrSTTNotInstalled) {
		t.Error("IsError should not match different errors")
	}
}

// ── Sentinel Error Tests ────────────────────────────────────────────

func TestSentinelErrors(t *testing.T) {
	sentinelErrors := []error{
		ErrLLMUnavailable,
		ErrLLMRequestFailed,
		ErrSTTNotInstalled,
		ErrSTTFailed,
		ErrMemoryUnreachable,
		ErrBrainNotInitialized,
		ErrRateLimited,
		ErrContextCanceled,
	}
	for _, e := range sentinelErrors {
		if e.Error() == "" {
			t.Errorf("sentinel error %v should have a message", e)
		}
	}
}

// ── Transient Error Detection ────────────────────────────────────────

func TestIsTransientError(t *testing.T) {
	if isTransientError(context.DeadlineExceeded) {
		t.Error("deadline exceeded should not be transient")
	}
	if isTransientError(context.Canceled) {
		t.Error("canceled should not be transient")
	}
}
