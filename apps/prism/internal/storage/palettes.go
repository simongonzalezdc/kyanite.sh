package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyanite/prism/internal/palette"
)

// SavePalette saves a palette to disk
func SavePalette(p palette.Palette) error {

	configDir, err := GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir: %w", err)
	}

	palettesDir := filepath.Join(configDir, "palettes")
	err = os.MkdirAll(palettesDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create palettes dir: %w", err)
	}

	filename := filepath.Join(palettesDir, p.ID+".json")
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal palette: %w", err)
	}

	return AtomicWrite(filename, data)
}

// LoadPalette loads a palette from disk
func LoadPalette(id string) (*palette.Palette, error) {

	configDir, err := GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config dir: %w", err)
	}

	filename := filepath.Join(configDir, "palettes", id+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read palette: %w", err)
	}

	var p palette.Palette
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal palette: %w", err)
	}

	return &p, nil
}

// ListPalettes lists all saved palettes
func ListPalettes() ([]palette.Palette, error) {

	configDir, err := GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config dir: %w", err)
	}

	palettesDir := filepath.Join(configDir, "palettes")
	files, err := os.ReadDir(palettesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []palette.Palette{}, nil
		}
		return nil, fmt.Errorf("failed to read palettes dir: %w", err)
	}

	// Pre-allocate based on file count for better performance
	palettes := make([]palette.Palette, 0, len(files))
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		id := file.Name()[:len(file.Name())-5] // Remove .json
		p, err := LoadPalette(id)
		if err != nil {
			// Log error but don't crash
			continue
		}
		palettes = append(palettes, *p)
	}

	return palettes, nil
}

// DeletePalette deletes a palette from disk
func DeletePalette(id string) error {

	configDir, err := GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir: %w", err)
	}

	filename := filepath.Join(configDir, "palettes", id+".json")
	return os.Remove(filename)
}
