package sync

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPairingManager(t *testing.T) {
	t.Run("generates unique codes", func(t *testing.T) {
		pm := NewPairingManager(nil)

		code1 := pm.GenerateCode()
		code2 := pm.GenerateCode()

		if code1 == code2 {
			t.Error("generated codes should be unique")
		}

		if len(code1) != PairingCodeLength {
			t.Errorf("code length should be %d, got %d", PairingCodeLength, len(code1))
		}
	})

	t.Run("validates codes correctly", func(t *testing.T) {
		pm := NewPairingManager(nil)

		code := pm.GenerateCode()

		// Code should be valid
		if !pm.ValidateCode(code) {
			t.Error("valid code should be validated")
		}

		// Code should not be reusable (one-time use)
		if pm.ValidateCode(code) {
			t.Error("code should only be valid once")
		}
	})

	t.Run("rejects invalid codes", func(t *testing.T) {
		pm := NewPairingManager(nil)

		// Invalid code should be rejected
		if pm.ValidateCode("000000") {
			t.Error("invalid code should be rejected")
		}

		if pm.ValidateCode("") {
			t.Error("empty code should be rejected")
		}
	})

	t.Run("pairs devices", func(t *testing.T) {
		pm := NewPairingManager(nil)

		device := pm.PairDevice("Test Device")

		if device.Name != "Test Device" {
			t.Errorf("expected name 'Test Device', got %s", device.Name)
		}

		if device.ID == "" {
			t.Error("device ID should not be empty")
		}

		if !pm.IsDevicePaired(device.ID) {
			t.Error("device should be paired")
		}
	})

	t.Run("gets paired devices", func(t *testing.T) {
		pm := NewPairingManager(nil)

		pm.PairDevice("Device 1")
		pm.PairDevice("Device 2")

		devices := pm.GetPairedDevices()
		if len(devices) != 2 {
			t.Errorf("expected 2 devices, got %d", len(devices))
		}
	})

	t.Run("unpairs devices", func(t *testing.T) {
		pm := NewPairingManager(nil)

		device := pm.PairDevice("Test Device")

		if !pm.UnpairDevice(device.ID) {
			t.Error("unpair should succeed")
		}

		if pm.IsDevicePaired(device.ID) {
			t.Error("device should not be paired after unpair")
		}
	})

	t.Run("updates last seen", func(t *testing.T) {
		pm := NewPairingManager(nil)

		device := pm.PairDevice("Test Device")
		originalLastSeen := device.LastSeen

		time.Sleep(10 * time.Millisecond)
		pm.UpdateLastSeen(device.ID)

		updatedDevice := pm.GetDevice(device.ID)
		if !updatedDevice.LastSeen.After(originalLastSeen) {
			t.Error("last seen should be updated")
		}
	})
}

func TestMediaStore(t *testing.T) {
	t.Run("creates directories", func(t *testing.T) {
		tempDir := t.TempDir()
		basePath := filepath.Join(tempDir, "media")

		ms, err := NewMediaStore(basePath)
		if err != nil {
			t.Fatalf("failed to create media store: %v", err)
		}

		// Check directories exist
		if _, err := os.Stat(ms.GetVoicePath()); os.IsNotExist(err) {
			t.Error("voice directory should exist")
		}

		if _, err := os.Stat(ms.GetPhotosPath()); os.IsNotExist(err) {
			t.Error("photos directory should exist")
		}
	})

	t.Run("saves voice memo", func(t *testing.T) {
		tempDir := t.TempDir()
		ms, _ := NewMediaStore(filepath.Join(tempDir, "media"))

		data := []byte("fake audio data")
		path, err := ms.SaveVoiceMemo("device123", data)
		if err != nil {
			t.Fatalf("failed to save voice memo: %v", err)
		}

		if path == "" {
			t.Error("path should not be empty")
		}

		// Verify file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("voice memo file should exist")
		}
	})

	t.Run("saves photo", func(t *testing.T) {
		tempDir := t.TempDir()
		ms, _ := NewMediaStore(filepath.Join(tempDir, "media"))

		data := []byte("fake photo data")
		path, err := ms.SavePhoto("device123", data)
		if err != nil {
			t.Fatalf("failed to save photo: %v", err)
		}

		if path == "" {
			t.Error("path should not be empty")
		}

		// Verify file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("photo file should exist")
		}
	})

	t.Run("deletes media", func(t *testing.T) {
		tempDir := t.TempDir()
		ms, _ := NewMediaStore(filepath.Join(tempDir, "media"))

		data := []byte("fake data")
		path, _ := ms.SaveVoiceMemo("device123", data)
		filename := filepath.Base(path)

		err := ms.Delete(filename)
		if err != nil {
			t.Fatalf("failed to delete: %v", err)
		}

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Error("file should be deleted")
		}
	})

	t.Run("lists files", func(t *testing.T) {
		tempDir := t.TempDir()
		ms, _ := NewMediaStore(filepath.Join(tempDir, "media"))

		// Save some files
		ms.SaveVoiceMemo("device1", []byte("data1"))
		ms.SaveVoiceMemo("device2", []byte("data2"))

		files, err := ms.ListVoiceMemos()
		if err != nil {
			t.Fatalf("failed to list files: %v", err)
		}

		if len(files) != 2 {
			t.Errorf("expected 2 files, got %d", len(files))
		}
	})
}

func TestWebSocketHub(t *testing.T) {
	t.Run("tracks client count", func(t *testing.T) {
		hub := NewWebSocketHub(nil)

		// Initially no clients
		if hub.ClientCount() != 0 {
			t.Errorf("expected 0 clients, got %d", hub.ClientCount())
		}
	})

	t.Run("starts and stops", func(t *testing.T) {
		hub := NewWebSocketHub(nil)

		// Start hub in background
		go hub.Run()

		// Give it time to start
		time.Sleep(10 * time.Millisecond)

		// Stop should not panic
		hub.Stop()
	})
}

func TestMessage(t *testing.T) {
	t.Run("creates message with type", func(t *testing.T) {
		msg := Message{
			Type:    "test",
			Payload: []byte(`{"key": "value"}`),
		}

		if msg.Type != "test" {
			t.Errorf("expected type 'test', got %s", msg.Type)
		}
	})
}

func TestCapturedIdea(t *testing.T) {
	t.Run("creates text idea", func(t *testing.T) {
		idea := CapturedIdea{
			ID:        "test-id",
			Type:      IdeaTypeText,
			Content:   "Test content",
			DeviceID:  "device-1",
			CreatedAt: time.Now(),
			SyncedAt:  time.Now(),
		}

		if idea.Type != IdeaTypeText {
			t.Errorf("expected type 'text', got %s", idea.Type)
		}
	})

	t.Run("creates tempo idea", func(t *testing.T) {
		idea := CapturedIdea{
			ID:       "test-id",
			Type:     IdeaTypeTempo,
			BPM:      120,
			DeviceID: "device-1",
		}

		if idea.BPM != 120 {
			t.Errorf("expected BPM 120, got %d", idea.BPM)
		}
	})
}

func TestGenerateNumericCode(t *testing.T) {
	t.Run("generates correct length", func(t *testing.T) {
		code := generateNumericCode(6)
		if len(code) != 6 {
			t.Errorf("expected length 6, got %d", len(code))
		}
	})

	t.Run("contains only digits", func(t *testing.T) {
		code := generateNumericCode(10)
		for _, c := range code {
			if c < '0' || c > '9' {
				t.Errorf("code should contain only digits, found %c", c)
			}
		}
	})
}

func TestGenerateDeviceID(t *testing.T) {
	t.Run("generates unique IDs", func(t *testing.T) {
		id1 := generateDeviceID()
		id2 := generateDeviceID()

		if id1 == id2 {
			t.Error("device IDs should be unique")
		}
	})

	t.Run("generates non-empty IDs", func(t *testing.T) {
		id := generateDeviceID()
		if id == "" {
			t.Error("device ID should not be empty")
		}
	})
}

func TestSyncStatus(t *testing.T) {
	status := SyncStatus{
		Connected:        true,
		LastSync:         time.Now(),
		PendingIdeas:     5,
		ConnectedDevices: 2,
	}

	if !status.Connected {
		t.Error("expected connected to be true")
	}

	if status.PendingIdeas != 5 {
		t.Errorf("expected 5 pending ideas, got %d", status.PendingIdeas)
	}
}

// Integration tests
func TestSyncServerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	t.Run("server starts and stops", func(t *testing.T) {
		tempDir := t.TempDir()
		cfg := DefaultSyncServerConfig()
		cfg.MediaPath = filepath.Join(tempDir, "media")
		cfg.Port = 0 // Use random port

		server, err := NewSyncServer(cfg)
		if err != nil {
			t.Fatalf("failed to create server: %v", err)
		}

		// Start should succeed
		if err := server.Start(); err != nil {
			t.Fatalf("failed to start server: %v", err)
		}

		if !server.IsRunning() {
			t.Error("server should be running")
		}

		// Stop should succeed
		if err := server.Stop(); err != nil {
			t.Fatalf("failed to stop server: %v", err)
		}

		if server.IsRunning() {
			t.Error("server should not be running after stop")
		}
	})
}
