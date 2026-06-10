package repository

import (
	"github.com/kyanite/focus/pkg/models"
)

// Repository defines the contract for task data access
type Repository interface {
	// Load retrieves all tasks from storage
	Load() ([]models.Task, error)

	// Save persists all tasks to storage
	Save(tasks []models.Task) error

	// Close releases any resources held by the repository
	Close() error
}
