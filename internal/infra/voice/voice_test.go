package voice

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	t.Run("basic write and read", func(t *testing.T) {
		rb := NewRingBuffer(10)

		samples := []float32{1.0, 2.0, 3.0, 4.0, 5.0}
		written := rb.Write(samples)

		if written != len(samples) {
			t.Errorf("expected %d samples written, got %d", len(samples), written)
		}

		if rb.Count() != len(samples) {
			t.Errorf("expected count %d, got %d", len(samples), rb.Count())
		}

		result := rb.Read()
		if len(result) != len(samples) {
			t.Errorf("expected %d samples read, got %d", len(samples), len(result))
		}

		for i, v := range result {
			if v != samples[i] {
				t.Errorf("sample %d: expected %f, got %f", i, samples[i], v)
			}
		}

		// Buffer should be empty after read
		if rb.Count() != 0 {
			t.Errorf("expected count 0 after read, got %d", rb.Count())
		}
	})

	t.Run("overflow behavior", func(t *testing.T) {
		rb := NewRingBuffer(5)

		// Write more than capacity
		samples := []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}
		rb.Write(samples)

		// Should only have last 5 samples
		if rb.Count() != 5 {
			t.Errorf("expected count 5, got %d", rb.Count())
		}

		result := rb.Read()
		expected := []float32{3.0, 4.0, 5.0, 6.0, 7.0}
		for i, v := range result {
			if v != expected[i] {
				t.Errorf("sample %d: expected %f, got %f", i, expected[i], v)
			}
		}
	})

	t.Run("clear", func(t *testing.T) {
		rb := NewRingBuffer(10)
		rb.Write([]float32{1.0, 2.0, 3.0})
		rb.Clear()

		if rb.Count() != 0 {
			t.Errorf("expected count 0 after clear, got %d", rb.Count())
		}

		result := rb.Read()
		if result != nil {
			t.Errorf("expected nil from empty buffer, got %v", result)
		}
	})
}

func TestDefaultAudioFormat(t *testing.T) {
	format := DefaultAudioFormat()

	if format.SampleRate != 16000 {
		t.Errorf("expected sample rate 16000, got %d", format.SampleRate)
	}

	if format.Channels != 1 {
		t.Errorf("expected 1 channel, got %d", format.Channels)
	}
}

func TestBytesToFloat32(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		result := bytesToFloat32(nil)
		if result != nil {
			t.Errorf("expected nil for empty input, got %v", result)
		}

		result = bytesToFloat32([]byte{})
		if result != nil {
			t.Errorf("expected nil for empty input, got %v", result)
		}
	})

	t.Run("converts bytes correctly", func(t *testing.T) {
		// Test with known float32 value: 1.0 = 0x3F800000 in IEEE 754
		// Little-endian: 00 00 80 3F
		data := []byte{0x00, 0x00, 0x80, 0x3F}
		result := bytesToFloat32(data)

		if len(result) != 1 {
			t.Fatalf("expected 1 sample, got %d", len(result))
		}

		if result[0] != 1.0 {
			t.Errorf("expected 1.0, got %f", result[0])
		}
	})
}

func TestModelManager(t *testing.T) {
	t.Run("creates models directory", func(t *testing.T) {
		tempDir := t.TempDir()
		modelsDir := filepath.Join(tempDir, "models")

		mm, err := NewModelManager(modelsDir, nil)
		if err != nil {
			t.Fatalf("failed to create model manager: %v", err)
		}

		if mm.GetModelsDir() != modelsDir {
			t.Errorf("expected models dir %s, got %s", modelsDir, mm.GetModelsDir())
		}

		// Check directory was created
		if _, err := os.Stat(modelsDir); os.IsNotExist(err) {
			t.Error("models directory was not created")
		}
	})

	t.Run("model path generation", func(t *testing.T) {
		tempDir := t.TempDir()
		mm, _ := NewModelManager(tempDir, nil)

		path := mm.GetModelPath(ModelBaseEN)
		expected := filepath.Join(tempDir, ModelBaseEN)

		if path != expected {
			t.Errorf("expected path %s, got %s", expected, path)
		}
	})

	t.Run("model availability check", func(t *testing.T) {
		tempDir := t.TempDir()
		mm, _ := NewModelManager(tempDir, nil)

		// Model shouldn't exist yet
		if mm.IsModelAvailable(ModelBaseEN) {
			t.Error("model should not be available")
		}

		// Create a fake model file
		modelPath := mm.GetModelPath(ModelBaseEN)
		if err := os.WriteFile(modelPath, []byte("fake model data"), 0o644); err != nil {
			t.Fatalf("failed to create fake model: %v", err)
		}

		// Now it should be available
		if !mm.IsModelAvailable(ModelBaseEN) {
			t.Error("model should be available")
		}
	})

	t.Run("invalid model name", func(t *testing.T) {
		tempDir := t.TempDir()
		mm, _ := NewModelManager(tempDir, nil)

		_, err := mm.EnsureModel("invalid-model.bin", nil)
		if err == nil {
			t.Error("expected error for invalid model name")
		}
	})

	t.Run("get available models", func(t *testing.T) {
		tempDir := t.TempDir()
		mm, _ := NewModelManager(tempDir, nil)

		// No models initially
		models := mm.GetAvailableModels()
		if len(models) != 0 {
			t.Errorf("expected 0 models, got %d", len(models))
		}

		// Create a fake model
		modelPath := mm.GetModelPath(ModelTinyEN)
		if err := os.WriteFile(modelPath, []byte("fake"), 0o644); err != nil {
			t.Fatalf("failed to create fake model: %v", err)
		}

		models = mm.GetAvailableModels()
		if len(models) != 1 {
			t.Errorf("expected 1 model, got %d", len(models))
		}
	})
}

func TestFormatModelSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 bytes"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{142000000, "135.4 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		result := FormatModelSize(tt.bytes)
		if result != tt.expected {
			t.Errorf("FormatModelSize(%d): expected %s, got %s", tt.bytes, tt.expected, result)
		}
	}
}

func TestRecommendedModel(t *testing.T) {
	tests := []struct {
		englishOnly bool
		preferSpeed bool
		expected    string
	}{
		{true, true, ModelTinyEN},
		{true, false, ModelBaseEN},
		{false, true, ModelTiny},
		{false, false, ModelBase},
	}

	for _, tt := range tests {
		result := RecommendedModel(tt.englishOnly, tt.preferSpeed)
		if result != tt.expected {
			t.Errorf("RecommendedModel(%v, %v): expected %s, got %s",
				tt.englishOnly, tt.preferSpeed, tt.expected, result)
		}
	}
}

func TestWhisperConfig(t *testing.T) {
	cfg := DefaultWhisperConfig()

	if cfg.Language != "en" {
		t.Errorf("expected language 'en', got %s", cfg.Language)
	}

	if cfg.Translate != false {
		t.Error("expected translate to be false")
	}

	if cfg.NoContext != true {
		t.Error("expected no_context to be true for dictation")
	}
}

func TestSupportedLanguages(t *testing.T) {
	langs := SupportedLanguages()

	if len(langs) == 0 {
		t.Error("expected non-empty language list")
	}

	// Check for English
	hasEnglish := false
	for _, lang := range langs {
		if lang == "en" {
			hasEnglish = true
			break
		}
	}

	if !hasEnglish {
		t.Error("expected English in supported languages")
	}
}

// Integration tests that require actual hardware/models are skipped by default
func TestAudioCapture_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Skip if no audio device available
	t.Skip("skipping audio capture test - requires microphone")

	// This would test actual audio capture if hardware is available
	cfg := DefaultAudioCaptureConfig()
	capture, err := NewAudioCapture(cfg)
	if err != nil {
		t.Fatalf("failed to create audio capture: %v", err)
	}
	defer capture.Close()

	devices, err := capture.ListDevices()
	if err != nil {
		t.Fatalf("failed to list devices: %v", err)
	}

	t.Logf("Found %d capture devices", len(devices))
}

func TestWhisperEngine_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Skip if model not available
	t.Skip("skipping whisper test - requires model file")

	// This would test actual transcription if model is available
	cfg := DefaultWhisperConfig()
	engine, err := NewWhisperEngine("path/to/model.bin", cfg, nil)
	if err != nil {
		t.Fatalf("failed to create whisper engine: %v", err)
	}
	defer engine.Close()

	if !engine.IsLoaded() {
		t.Error("expected model to be loaded")
	}
}
