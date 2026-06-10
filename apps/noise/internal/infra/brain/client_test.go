package brain

import (
	"context"
	"encoding/binary"
	"testing"
	"time"

	ai "github.com/kyanite/ai"
)

func TestNewClient(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestClient_IsAvailable_Unreachable(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	if client.IsAvailable(context.Background()) {
		t.Error("IsAvailable() should be false when NUCBox is unreachable")
	}
}

func TestClient_Generate_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	_, err := client.Generate(context.Background(), "model", "hello", nil)
	if err == nil {
		t.Error("Generate() should return error when brain is nil")
	}
}

func TestClient_GenerateWithOptions_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	_, err := client.GenerateWithOptions(context.Background(), "hello")
	if err == nil {
		t.Error("GenerateWithOptions() should return error when brain is nil")
	}
}

func TestClient_Chat_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	_, err := client.Chat(context.Background(), []ai.Message{{Role: "user", Content: "hi"}})
	if err == nil {
		t.Error("Chat() should return error when brain is nil")
	}
}

func TestClient_StreamGenerate_Offline(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	ch, err := client.StreamGenerate(context.Background(), "model", "hello")
	// Brain may be nil (no whisper) or not (whisper installed).
	// If nil: should return ErrBrainNotInitialized with nil channel.
	// If not nil: may return a channel (streaming to unreachable LLM) or error.
	if err != nil {
		// Error path — verify it's a typed error
		if ch != nil {
			t.Error("StreamGenerate() should return nil channel on error")
		}
		t.Logf("StreamGenerate returned error (expected offline): %v", err)
		return
	}
	// Channel path — consume and expect eventual error/close
	for range ch {
	}
}

func TestClient_TranscribePCM_ShortTimeout(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	if !client.IsSTTAvailable() {
		t.Skip("STT not available, skipping")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := client.TranscribePCM(ctx, []float32{0.1, 0.2})
	// With only 2 samples, whisper may timeout or fail — either is acceptable
	if err == nil {
		t.Log("TranscribePCM returned successfully with minimal data")
	}
}

func TestClient_TranscribeFile_Nonexistent(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	if !client.IsSTTAvailable() {
		// Without STT, the brain wraps the error correctly
		_, err := client.TranscribeFile(context.Background(), "nonexistent.wav")
		if err == nil {
			t.Error("TranscribeFile should error on nonexistent file without STT")
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = client.TranscribeFile(ctx, "nonexistent_file_that_does_not_exist.wav")
	// Whisper may or may not error on nonexistent file — just verify no panic.
}

func TestClient_Close_NilBrain(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	// Must not panic.
	client.Close()
}

func TestClient_Close_Idempotent(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	client.Close()
	client.Close() // second close should also be safe
}

func TestClient_IsSTTAvailable_Check(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	defer client.Close()
	// STT availability depends on whether whisper.cpp is installed locally.
	// Just verify the method doesn't panic and returns a bool.
	_ = client.IsSTTAvailable()
}

func TestClient_Brain_Check(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	defer client.Close()
	// Brain() may be nil or not depending on local setup. Just verify no panic.
	_ = client.Brain()
}

func TestClient_GetConfig(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	cfg := client.GetConfig()
	if cfg.OllamaURL != "http://localhost:19999" {
		t.Errorf("GetConfig().OllamaURL = %q, want %q", cfg.OllamaURL, "http://localhost:19999")
	}
}

func TestClient_GetModelStatus(t *testing.T) {
	t.Setenv("KYANITE_OLLAMA_URL", "http://localhost:19999")

	client := NewClient()
	status := client.GetModelStatus()
	if status == nil {
		t.Fatal("GetModelStatus() returned nil")
	}
	if status["provider"] != "brain" {
		t.Errorf("GetModelStatus()[\"provider\"] = %v, want \"brain\"", status["provider"])
	}
}

func TestFloat32ToPCM16(t *testing.T) {
	tests := []struct {
		name     string
		input    []float32
		expected []byte
	}{
		{
			name:     "empty input",
			input:    []float32{},
			expected: []byte{},
		},
		{
			name:     "zero value",
			input:    []float32{0.0},
			expected: []byte{0x00, 0x00},
		},
		{
			name:     "positive one",
			input:    []float32{1.0},
			expected: []byte{0xff, 0x7f}, // 32767 as little-endian uint16
		},
		{
			name:     "negative one",
			input:    []float32{-1.0},
			expected: []byte{0x01, 0x80}, // int16(-32767) = 0x8001 LE
		},
		{
			name:     "clamped above one",
			input:    []float32{2.0},
			expected: []byte{0xff, 0x7f}, // clamped to 32767
		},
		{
			name:     "clamped below negative one",
			input:    []float32{-2.0},
			expected: []byte{0x00, 0x80}, // clamped to -32768 → 0x8000
		},
		{
			name:     "mid-range positive",
			input:    []float32{0.5},
			expected: func() []byte {
				b := make([]byte, 2)
				s := float32(0.5)
				v := int16(s * 32767.0)
				binary.LittleEndian.PutUint16(b, uint16(v))
				return b
			}(),
		},
		{
			name:     "multiple samples",
			input:    []float32{0.0, 1.0, -1.0},
			expected: func() []byte {
				b := make([]byte, 6)
				binary.LittleEndian.PutUint16(b[0:], 0)
				binary.LittleEndian.PutUint16(b[2:], uint16(int16(32767)))
				neg := float32(-1.0)
				binary.LittleEndian.PutUint16(b[4:], uint16(int16(neg*32767.0)))
				return b
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := float32ToPCM16(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("float32ToPCM16() length = %d, want %d", len(got), len(tt.expected))
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("byte %d: got 0x%02x, want 0x%02x", i, got[i], tt.expected[i])
				}
			}
		})
	}
}
