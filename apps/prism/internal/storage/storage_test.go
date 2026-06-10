package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	color "github.com/kyanite/prism/internal/color"
	palette "github.com/kyanite/prism/internal/palette"
)

func cleanupTempDir(t *testing.T, path string) {
	t.Helper()
	t.Cleanup(func() {
		if err := os.RemoveAll(path); err != nil {
			t.Errorf("Failed to remove temp dir %s: %v", path, err)
		}
	})
}

// TestStorageIntegration tests complete storage workflows
func TestStorageIntegration(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "prism-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	cleanupTempDir(t, tmpDir)

	// Override config directory for testing
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	t.Run("ConfigSaveAndLoad", func(t *testing.T) {
		// Create and save a config
		cfg := &Config{
			Theme:            "custom-theme",
			AutoSave:         true,
			AutoSaveInterval: 60,
		}

		err := SaveConfig(cfg)
		if err != nil {
			t.Fatalf("SaveConfig failed: %v", err)
		}

		// Load the config
		loadedCfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		// Verify loaded config matches saved config
		if loadedCfg.Theme != cfg.Theme {
			t.Errorf("Theme = %s, want %s", loadedCfg.Theme, cfg.Theme)
		}
		if loadedCfg.AutoSave != cfg.AutoSave {
			t.Errorf("AutoSave = %v, want %v", loadedCfg.AutoSave, cfg.AutoSave)
		}
		if loadedCfg.AutoSaveInterval != cfg.AutoSaveInterval {
			t.Errorf("AutoSaveInterval = %d, want %d", loadedCfg.AutoSaveInterval, cfg.AutoSaveInterval)
		}
	})

	t.Run("DefaultConfig", func(t *testing.T) {
		// Remove config file if it exists
		configDir, _ := GetConfigDir()
		if err := os.RemoveAll(configDir); err != nil {
			t.Fatalf("Failed to remove config dir: %v", err)
		}

		// Load config should return default when file doesn't exist
		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		defaultCfg := DefaultConfig()
		if cfg.Theme != defaultCfg.Theme {
			t.Errorf("Theme = %s, want %s", cfg.Theme, defaultCfg.Theme)
		}
	})

	t.Run("PaletteSaveLoadDelete", func(t *testing.T) {
		// Create a test palette
		baseColor, _ := color.ParseHex("#FF5733")
		pal, _ := palette.Generate(baseColor, palette.Triadic)
		pal.Name = "Test Triadic Palette"
		pal.Description = "Integration test palette"
		pal.Tags = []string{"test", "triadic"}

		// Save palette
		err := SavePalette(pal)
		if err != nil {
			t.Fatalf("SavePalette failed: %v", err)
		}

		// Load palette
		loadedPal, err := LoadPalette(pal.ID)
		if err != nil {
			t.Fatalf("LoadPalette failed: %v", err)
		}

		// Verify palette data
		if loadedPal.Name != pal.Name {
			t.Errorf("Name = %s, want %s", loadedPal.Name, pal.Name)
		}
		if loadedPal.Description != pal.Description {
			t.Errorf("Description = %s, want %s", loadedPal.Description, pal.Description)
		}
		if len(loadedPal.Colors) != len(pal.Colors) {
			t.Errorf("Colors count = %d, want %d", len(loadedPal.Colors), len(pal.Colors))
		}
		if len(loadedPal.Tags) != len(pal.Tags) {
			t.Errorf("Tags count = %d, want %d", len(loadedPal.Tags), len(pal.Tags))
		}

		// Delete palette
		err = DeletePalette(pal.ID)
		if err != nil {
			t.Fatalf("DeletePalette failed: %v", err)
		}

		// Verify palette is deleted
		_, err = LoadPalette(pal.ID)
		if err == nil {
			t.Error("LoadPalette should fail after delete")
		}
	})

	t.Run("ListPalettes", func(t *testing.T) {
		// Clean up any existing palettes
		configDir, _ := GetConfigDir()
		palettesDir := filepath.Join(configDir, "palettes")
		if err := os.RemoveAll(palettesDir); err != nil {
			t.Fatalf("Failed to remove palettes dir: %v", err)
		}

		// Create multiple test palettes
		palettes := []palette.Palette{}
		rules := []palette.HarmonyRule{
			palette.Monochromatic,
			palette.Complementary,
			palette.Analogous,
		}

		for i, rule := range rules {
			baseColor, _ := color.ParseHex("#FF5733")
			pal, _ := palette.Generate(baseColor, rule)
			pal.Name = string(rule)
			palettes = append(palettes, pal)

			err := SavePalette(pal)
			if err != nil {
				t.Fatalf("SavePalette %d failed: %v", i, err)
			}
		}

		// List all palettes
		loadedPalettes, err := ListPalettes()
		if err != nil {
			t.Fatalf("ListPalettes failed: %v", err)
		}

		// Verify count
		if len(loadedPalettes) != len(palettes) {
			t.Errorf("ListPalettes count = %d, want %d", len(loadedPalettes), len(palettes))
		}

		// Clean up
		for _, pal := range palettes {
			if err := DeletePalette(pal.ID); err != nil {
				t.Fatalf("DeletePalette cleanup failed: %v", err)
			}
		}
	})

	t.Run("AtomicWrite", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "atomic-test.txt")
		testData := []byte("test data for atomic write")

		err := AtomicWrite(testFile, testData)
		if err != nil {
			t.Fatalf("AtomicWrite failed: %v", err)
		}

		// Verify file exists and has correct content
		data, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}

		if string(data) != string(testData) {
			t.Errorf("File content = %s, want %s", string(data), string(testData))
		}
	})

	t.Run("ConcurrentPaletteSaves", func(t *testing.T) {
		// Test concurrent saves to ensure atomic writes work correctly
		done := make(chan bool)

		for i := 0; i < 5; i++ {
			go func(index int) {
				baseColor, _ := color.ParseHex("#FF5733")
				pal, _ := palette.Generate(baseColor, palette.Triadic)
				pal.Name = "Concurrent Test"
				err := SavePalette(pal)
				if err != nil {
					t.Errorf("Concurrent SavePalette %d failed: %v", index, err)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 5; i++ {
			<-done
		}
	})
}

func TestGetConfigDirIgnoresRelativeXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "relative-config")

	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir failed: %v", err)
	}

	if !filepath.IsAbs(configDir) {
		t.Fatalf("Expected absolute config dir, got %q", configDir)
	}
	if filepath.Clean(configDir) == filepath.Join("relative-config", "prism") {
		t.Fatalf("Expected relative XDG_CONFIG_HOME to be ignored, got %q", configDir)
	}
}

// TestStorageErrorHandling tests error conditions
func TestStorageErrorHandling(t *testing.T) {
	t.Run("LoadNonexistentPalette", func(t *testing.T) {
		_, err := LoadPalette("nonexistent-id-12345")
		if err == nil {
			t.Error("LoadPalette should fail for nonexistent palette")
		}
	})

	t.Run("DeleteNonexistentPalette", func(t *testing.T) {
		err := DeletePalette("nonexistent-id-12345")
		if err == nil {
			t.Error("DeletePalette should fail for nonexistent palette")
		}
	})

	t.Run("ListPalettesEmptyDirectory", func(t *testing.T) {
		// Create temporary test directory
		tmpDir, err := os.MkdirTemp("", "prism-empty-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		cleanupTempDir(t, tmpDir)

		// Override config directory
		t.Setenv("XDG_CONFIG_HOME", tmpDir)

		palettes, err := ListPalettes()
		if err != nil {
			t.Fatalf("ListPalettes should not fail on empty directory: %v", err)
		}

		if len(palettes) != 0 {
			t.Errorf("ListPalettes should return empty slice, got %d palettes", len(palettes))
		}
	})
}

// TestPaletteVersioning tests palette updates
func TestPaletteVersioning(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "prism-version-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	cleanupTempDir(t, tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create and save a palette
	baseColor, _ := color.ParseHex("#FF5733")
	pal, _ := palette.Generate(baseColor, palette.Triadic)
	originalName := "Original Name"
	pal.Name = originalName

	err = SavePalette(pal)
	if err != nil {
		t.Fatalf("SavePalette failed: %v", err)
	}

	// Simulate time passing
	time.Sleep(10 * time.Millisecond)

	// Update the palette
	pal.Name = "Updated Name"
	pal.UpdatedAt = time.Now()
	err = SavePalette(pal)
	if err != nil {
		t.Fatalf("SavePalette (update) failed: %v", err)
	}

	// Load and verify update
	loadedPal, err := LoadPalette(pal.ID)
	if err != nil {
		t.Fatalf("LoadPalette failed: %v", err)
	}

	if loadedPal.Name != "Updated Name" {
		t.Errorf("Updated name = %s, want %s", loadedPal.Name, "Updated Name")
	}

	if !loadedPal.UpdatedAt.After(loadedPal.CreatedAt) {
		t.Error("UpdatedAt should be after CreatedAt")
	}
}
