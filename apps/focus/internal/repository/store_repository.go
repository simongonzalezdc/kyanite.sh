package repository

import (
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/models"
)

// StoreRepository implements Repository using the file-based Store
type StoreRepository struct {
	store *store.Store
}

// NewStoreRepository creates a new StoreRepository
func NewStoreRepository(storagePath string) *StoreRepository {
	return &StoreRepository{
		store: store.New(storagePath),
	}
}

// Load retrieves all tasks from storage
func (r *StoreRepository) Load() ([]models.Task, error) {
	return r.store.Load()
}

// Save persists all tasks to storage
func (r *StoreRepository) Save(tasks []models.Task) error {
	return r.store.Save(tasks)
}

// Close releases any resources (no-op for file store)
func (r *StoreRepository) Close() error {
	// File-based store doesn't need cleanup
	return nil
}
