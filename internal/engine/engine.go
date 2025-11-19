package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kyanite/focus/internal/repository"
	pkgerrors "github.com/kyanite/focus/pkg/errors"
	"github.com/kyanite/focus/pkg/models"
)

// Engine handles task management operations with in-memory caching
type Engine struct {
	repo repository.Repository

	// In-memory cache for performance
	cache      []models.Task
	cacheIndex map[string]int // ID -> index mapping for O(1) lookups
	cacheDirty bool           // Track if cache needs persisting
	mu         sync.RWMutex   // Protect cache access
	cacheValid bool           // Track if cache is loaded
}

// New creates a new engine instance
func New(repo repository.Repository) *Engine {
	e := &Engine{
		repo:       repo,
		cacheIndex: make(map[string]int),
		cacheValid: false,
	}
	// Eagerly load cache on creation
	_ = e.loadCache()
	return e
}

// loadCache loads tasks from storage into memory
func (e *Engine) loadCache() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	tasks, err := e.repo.Load()
	if err != nil {
		return err
	}

	e.cache = tasks
	e.cacheIndex = make(map[string]int, len(tasks))
	for i, task := range tasks {
		e.cacheIndex[task.ID] = i
	}
	e.cacheValid = true
	e.cacheDirty = false

	return nil
}

// flushCache saves cached tasks to storage if dirty
func (e *Engine) flushCache() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.cacheDirty {
		return nil // Nothing to save
	}

	if err := e.repo.Save(e.cache); err != nil {
		return err
	}

	e.cacheDirty = false
	return nil
}

// invalidateCache marks cache as invalid, forcing reload on next access
func (e *Engine) invalidateCache() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cacheValid = false
}

// AddTask creates and stores a new task
func (e *Engine) AddTask(parsedTask models.ParsedTask) (models.Task, error) {
	// Validate input
	if parsedTask.Description == "" {
		return models.Task{}, pkgerrors.ErrEmptyDescription
	}

	now := time.Now()
	task := models.Task{
		ID:                 generateID(),
		Description:        parsedTask.Description,
		Status:             "pending",
		Priority:           parsedTask.Priority,
		Deadline:           parsedTask.Deadline,
		Categories:         parsedTask.Categories,
		Notes:              parsedTask.Notes,
		CreatedAt:          now,
		UpdatedAt:          now,
		RecurrencePattern:  parsedTask.RecurrencePattern,
		RecurrenceInterval: parsedTask.RecurrenceInterval,
		RecurrenceEndDate:  parsedTask.RecurrenceEndDate,
	}

	// Validate priority
	if task.Priority == "" {
		task.Priority = "medium"
	}

	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}

	if !validPriorities[task.Priority] {
		task.Priority = "medium" // Normalize invalid priority
	}

	// Add to cache
	e.mu.Lock()
	e.cache = append(e.cache, task)
	e.cacheIndex[task.ID] = len(e.cache) - 1
	e.cacheDirty = true
	e.mu.Unlock()

	// Persist to disk
	if err := e.flushCache(); err != nil {
		return models.Task{}, fmt.Errorf("failed to save task: %w", err)
	}

	return task, nil
}

// AddSubtask creates a new task as a subtask of an existing parent task
func (e *Engine) AddSubtask(parentID, description, priority string, categories []string, deadline string) (string, error) {
	if description == "" {
		return "", fmt.Errorf("subtask description cannot be empty")
	}

	// Get parent task
	parent, err := e.GetTask(parentID)
	if err != nil {
		return "", fmt.Errorf("parent task not found: %w", err)
	}

	// Create subtask
	now := time.Now()
	subtask := models.Task{
		ID:          generateID(),
		Description: description,
		Status:      "pending",
		Priority:    priority,
		Categories:  categories,
		CreatedAt:   now,
		UpdatedAt:   now,
		ParentID:    parentID,
	}

	// Parse deadline if provided
	if deadline != "" {
		deadlineTime, err := time.Parse("2006-01-02", deadline)
		if err == nil {
			subtask.Deadline = deadlineTime
		}
	}

	// Validate and normalize priority
	if subtask.Priority == "" {
		subtask.Priority = parent.Priority // Inherit from parent if not specified
	}
	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}
	if !validPriorities[subtask.Priority] {
		subtask.Priority = "medium"
	}

	// Add subtask to cache
	e.mu.Lock()
	e.cache = append(e.cache, subtask)
	e.cacheIndex[subtask.ID] = len(e.cache) - 1

	// Update parent task to include this subtask
	if parentIdx, exists := e.cacheIndex[parentID]; exists {
		e.cache[parentIdx].AddSubtask(subtask.ID)
		e.cache[parentIdx].UpdatedAt = now
	}
	e.cacheDirty = true
	e.mu.Unlock()

	// Persist to disk
	if err := e.flushCache(); err != nil {
		return "", fmt.Errorf("failed to save subtask: %w", err)
	}

	return subtask.ID, nil
}

// ListTasks returns all tasks, optionally filtered
func (e *Engine) ListTasks(filter string) ([]models.Task, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if filter == "" || filter == "all" {
		// Return copy of cache to prevent external modification
		result := make([]models.Task, len(e.cache))
		copy(result, e.cache)
		return result, nil
	}

	// Pre-allocate with estimated capacity
	filtered := make([]models.Task, 0, len(e.cache)/2)
	for _, task := range e.cache {
		if filter == "active" && task.Status == "pending" {
			filtered = append(filtered, task)
		} else if filter == "completed" && task.Status == "completed" {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

// CompleteTask marks a task as completed
func (e *Engine) CompleteTask(id string) error {
	return e.updateTaskStatus(id, "completed")
}

// DeleteTask removes a task
func (e *Engine) DeleteTask(id string) error {
	e.mu.Lock()

	// Find task in cache using index
	idx, exists := e.cacheIndex[id]
	if !exists {
		e.mu.Unlock()
		return fmt.Errorf("task with ID %s not found", id)
	}

	// Remove from cache
	e.cache = append(e.cache[:idx], e.cache[idx+1:]...)

	// Rebuild index for all tasks after deleted one
	delete(e.cacheIndex, id)
	for i := idx; i < len(e.cache); i++ {
		e.cacheIndex[e.cache[i].ID] = i
	}

	e.cacheDirty = true
	e.mu.Unlock()

	// Persist to disk
	return e.flushCache()
}

// GetTask retrieves a specific task by ID
func (e *Engine) GetTask(id string) (models.Task, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// O(1) lookup using index
	idx, exists := e.cacheIndex[id]
	if !exists {
		return models.Task{}, pkgerrors.NewTaskError("GetTask", id, pkgerrors.ErrTaskNotFound)
	}

	return e.cache[idx], nil
}

// updateTaskStatus updates the status of a task
func (e *Engine) updateTaskStatus(id, status string) error {
	e.mu.Lock()

	// O(1) lookup using index
	idx, exists := e.cacheIndex[id]
	if !exists {
		e.mu.Unlock()
		return pkgerrors.NewTaskError("updateTaskStatus", id, pkgerrors.ErrTaskNotFound)
	}

	e.cache[idx].Status = status
	e.cache[idx].UpdatedAt = time.Now()
	e.cacheDirty = true
	e.mu.Unlock()

	// Persist to disk
	return e.flushCache()
}

// UpdateTask updates a task with new values
func (e *Engine) UpdateTask(updatedTask models.Task) error {
	e.mu.Lock()

	// O(1) lookup using index
	idx, exists := e.cacheIndex[updatedTask.ID]
	if !exists {
		e.mu.Unlock()
		return fmt.Errorf("task with ID %s not found", updatedTask.ID)
	}

	// Preserve original creation time
	updatedTask.CreatedAt = e.cache[idx].CreatedAt
	updatedTask.UpdatedAt = time.Now()
	e.cache[idx] = updatedTask
	e.cacheDirty = true
	e.mu.Unlock()

	// Persist to disk
	return e.flushCache()
}

// UpdateTaskStatus updates only the status of a task
func (e *Engine) UpdateTaskStatus(id, status string) error {
	// Reuse the private updateTaskStatus method
	return e.updateTaskStatus(id, status)
}

// generateID creates a unique ID for tasks using UUID
func generateID() string {
	return uuid.New().String()
}
