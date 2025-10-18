package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ChordProgression represents a chord progression with metadata
type ChordProgression struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Chords      []string `json:"chords"`
	Mood        string   `json:"mood"`
	Description string   `json:"description"`
	Genre       string   `json:"genre"`
}

// ChordProgressionData represents the full chord progressions dataset
type ChordProgressionData struct {
	Progressions []ChordProgression `json:"progressions"`
}

// ChordProgressionLoader handles loading and caching chord progressions
type ChordProgressionLoader struct {
	mu          sync.RWMutex
	data        *ChordProgressionData
	dataPath    string
	initialized bool
}

var (
	// Global instance for chord progression loading
	chordLoader *ChordProgressionLoader
	once        sync.Once
)

// GetChordProgressionLoader returns the singleton chord progression loader
func GetChordProgressionLoader() *ChordProgressionLoader {
	once.Do(func() {
		chordLoader = &ChordProgressionLoader{
			dataPath: "data/chord_progressions.json",
		}
	})
	return chordLoader
}

// LoadProgressions loads chord progressions from the JSON file
func (cpl *ChordProgressionLoader) LoadProgressions() error {
	cpl.mu.Lock()
	defer cpl.mu.Unlock()

	// Check if already loaded
	if cpl.initialized {
		return nil
	}

	// Read the JSON file
	data, err := os.ReadFile(cpl.dataPath)
	if err != nil {
		return fmt.Errorf("failed to read chord progressions file: %w", err)
	}

	// Parse JSON
	var progressionData ChordProgressionData
	if err := json.Unmarshal(data, &progressionData); err != nil {
		return fmt.Errorf("failed to parse chord progressions JSON: %w", err)
	}

	cpl.data = &progressionData
	cpl.initialized = true

	return nil
}

// GetAllProgressions returns all chord progressions
func (cpl *ChordProgressionLoader) GetAllProgressions() ([]ChordProgression, error) {
	cpl.mu.RLock()
	defer cpl.mu.RUnlock()

	if !cpl.initialized {
		if err := cpl.LoadProgressions(); err != nil {
			return nil, err
		}
	}

	return cpl.data.Progressions, nil
}

// GetProgressionsByMood returns chord progressions filtered by mood
func (cpl *ChordProgressionLoader) GetProgressionsByMood(mood string) ([]ChordProgression, error) {
	progressions, err := cpl.GetAllProgressions()
	if err != nil {
		return nil, err
	}

	var filtered []ChordProgression
	for _, prog := range progressions {
		if prog.Mood == mood {
			filtered = append(filtered, prog)
		}
	}

	return filtered, nil
}

// GetProgressionByID returns a specific chord progression by ID
func (cpl *ChordProgressionLoader) GetProgressionByID(id string) (*ChordProgression, error) {
	progressions, err := cpl.GetAllProgressions()
	if err != nil {
		return nil, err
	}

	for _, prog := range progressions {
		if prog.ID == id {
			return &prog, nil
		}
	}

	return nil, fmt.Errorf("chord progression with ID '%s' not found", id)
}

// GetAvailableMoods returns all available mood categories
func (cpl *ChordProgressionLoader) GetAvailableMoods() ([]string, error) {
	progressions, err := cpl.GetAllProgressions()
	if err != nil {
		return nil, err
	}

	moodSet := make(map[string]bool)
	for _, prog := range progressions {
		moodSet[prog.Mood] = true
	}

	var moods []string
	for mood := range moodSet {
		moods = append(moods, mood)
	}

	return moods, nil
}

// SetDataPath sets the path to the chord progressions data file
func (cpl *ChordProgressionLoader) SetDataPath(path string) {
	cpl.mu.Lock()
	defer cpl.mu.Unlock()
	cpl.dataPath = path
	cpl.initialized = false // Reset to force reload
}

// Reload forces a reload of the chord progressions data
func (cpl *ChordProgressionLoader) Reload() error {
	cpl.mu.Lock()
	defer cpl.mu.Unlock()
	cpl.initialized = false
	return cpl.LoadProgressions()
}

// GetDefaultDataPath returns the default path to the chord progressions data file
func GetDefaultDataPath() string {
	// Try to find the data file relative to the executable or current directory
	if path, err := filepath.Abs("data/chord_progressions.json"); err == nil {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	// Fallback to relative path
	return "data/chord_progressions.json"
}

// Convenience functions for global access

// GetAllChordProgressions returns all chord progressions using the global loader
func GetAllChordProgressions() ([]ChordProgression, error) {
	return GetChordProgressionLoader().GetAllProgressions()
}

// GetChordProgressionsByMood returns chord progressions filtered by mood using the global loader
func GetChordProgressionsByMood(mood string) ([]ChordProgression, error) {
	return GetChordProgressionLoader().GetProgressionsByMood(mood)
}

// GetChordProgressionByID returns a specific chord progression by ID using the global loader
func GetChordProgressionByID(id string) (*ChordProgression, error) {
	return GetChordProgressionLoader().GetProgressionByID(id)
}

// GetAvailableChordMoods returns all available mood categories using the global loader
func GetAvailableChordMoods() ([]string, error) {
	return GetChordProgressionLoader().GetAvailableMoods()
}