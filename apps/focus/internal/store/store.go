package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kyanite/focus/pkg/models"
)

// Store handles persistence of tasks
type Store struct {
	filePath string
	mu       sync.Mutex
}

// New creates a new store instance
func New(filePath string) *Store {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		// If we can't create the directory, fall back to current directory
		// Log the error for debugging (will be added with logging framework)
		// For now, use a safe fallback
		return &Store{
			filePath: "./tasks.json",
		}
	}

	return &Store{
		filePath: filePath,
	}
}

// Load retrieves all tasks from storage
func (s *Store) Load() ([]models.Task, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return empty slice
			return []models.Task{}, nil
		}
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	if len(data) == 0 {
		// Empty file, return empty slice
		return []models.Task{}, nil
	}

	var tasks []models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks file: %w", err)
	}

	return tasks, nil
}

// Save persists tasks to storage
func (s *Store) Save(tasks []models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize tasks: %w", err)
	}

	// Write to temporary file first
	tempPath := s.filePath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Atomically rename temp file to target file
	if err := os.Rename(tempPath, s.filePath); err != nil {
		// Clean up temp file if rename failed
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to save tasks file: %w", err)
	}

	return nil
}
