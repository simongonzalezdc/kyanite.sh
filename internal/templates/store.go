package templates

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/kyanite/focus/pkg/models"
)

// Store handles template persistence
type Store struct {
	filePath string
}

// NewStore creates a new template store
func NewStore(filePath string) *Store {
	return &Store{
		filePath: filePath,
	}
}

// Load reads all templates from disk
func (s *Store) Load() ([]models.TaskTemplate, error) {
	// Ensure directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// Return empty list if file doesn't exist yet
		return []models.TaskTemplate{}, nil
	}

	// Read file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	// Handle empty file
	if len(data) == 0 {
		return []models.TaskTemplate{}, nil
	}

	// Parse JSON
	var templates []models.TaskTemplate
	if err := json.Unmarshal(data, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}

// Save writes all templates to disk
func (s *Store) Save(templates []models.TaskTemplate) error {
	// Ensure directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(templates, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(s.filePath, data, 0644)
}

// Add creates a new template
func (s *Store) Add(template models.TaskTemplate) (models.TaskTemplate, error) {
	templates, err := s.Load()
	if err != nil {
		return models.TaskTemplate{}, err
	}

	// Generate ID and timestamps
	now := time.Now()
	template.ID = uuid.New().String()
	template.CreatedAt = now
	template.UpdatedAt = now

	// Append and save
	templates = append(templates, template)
	if err := s.Save(templates); err != nil {
		return models.TaskTemplate{}, err
	}

	return template, nil
}

// Delete removes a template by ID
func (s *Store) Delete(id string) error {
	templates, err := s.Load()
	if err != nil {
		return err
	}

	// Find and remove template
	for i, template := range templates {
		if template.ID == id {
			templates = append(templates[:i], templates[i+1:]...)
			return s.Save(templates)
		}
	}

	return nil // Not found, but not an error
}

// Get retrieves a template by ID
func (s *Store) Get(id string) (models.TaskTemplate, error) {
	templates, err := s.Load()
	if err != nil {
		return models.TaskTemplate{}, err
	}

	for _, template := range templates {
		if template.ID == id {
			return template, nil
		}
	}

	return models.TaskTemplate{}, os.ErrNotExist
}

// GetByName retrieves a template by name
func (s *Store) GetByName(name string) (models.TaskTemplate, error) {
	templates, err := s.Load()
	if err != nil {
		return models.TaskTemplate{}, err
	}

	for _, template := range templates {
		if template.Name == name {
			return template, nil
		}
	}

	return models.TaskTemplate{}, os.ErrNotExist
}
