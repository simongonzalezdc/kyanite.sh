package engine

import (
	"fmt"
	"time"
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/models"
)

// Engine handles task management operations
type Engine struct {
	store *store.Store
}

// New creates a new engine instance
func New(store *store.Store) *Engine {
	return &Engine{
		store: store,
	}
}

// AddTask creates and stores a new task
func (e *Engine) AddTask(parsedTask models.ParsedTask) (models.Task, error) {
	// Validate input
	if parsedTask.Description == "" {
		return models.Task{}, fmt.Errorf("task description cannot be empty")
	}
	
	tasks, err := e.store.Load()
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to load tasks: %w", err)
	}
	
	now := time.Now()
	task := models.Task{
		ID:          generateID(),
		Description: parsedTask.Description,
		Status:      "pending",
		Priority:    parsedTask.Priority,
		Deadline:    parsedTask.Deadline,
		Categories:  parsedTask.Categories,
		CreatedAt:   now,
		UpdatedAt:   now,
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
	
	tasks = append(tasks, task)
	err = e.store.Save(tasks)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to save task: %w", err)
	}
	
	return task, nil
}

// ListTasks returns all tasks, optionally filtered
func (e *Engine) ListTasks(filter string) ([]models.Task, error) {
	tasks, err := e.store.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}
	
	if filter == "" || filter == "all" {
		return tasks, nil
	}
	
	var filtered []models.Task
	for _, task := range tasks {
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
	tasks, err := e.store.Load()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}
	
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return e.store.Save(tasks)
		}
	}
	
	return fmt.Errorf("task with ID %s not found", id)
}

// GetTask retrieves a specific task by ID
func (e *Engine) GetTask(id string) (models.Task, error) {
	tasks, err := e.store.Load()
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to load tasks: %w", err)
	}
	
	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}
	
	return models.Task{}, fmt.Errorf("task with ID %s not found", id)
}

// updateTaskStatus updates the status of a task
func (e *Engine) updateTaskStatus(id, status string) error {
	tasks, err := e.store.Load()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}
	
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			return e.store.Save(tasks)
		}
	}
	
	return fmt.Errorf("task with ID %s not found", id)
}

// UpdateTask updates a task with new values
func (e *Engine) UpdateTask(updatedTask models.Task) error {
	tasks, err := e.store.Load()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}
	
	for i, task := range tasks {
		if task.ID == updatedTask.ID {
			// Preserve original creation time
			updatedTask.CreatedAt = task.CreatedAt
			updatedTask.UpdatedAt = time.Now()
			tasks[i] = updatedTask
			return e.store.Save(tasks)
		}
	}
	
	return fmt.Errorf("task with ID %s not found", updatedTask.ID)
}

// UpdateTaskStatus updates only the status of a task
func (e *Engine) UpdateTaskStatus(id, status string) error {
	tasks, err := e.store.Load()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}
	
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			return e.store.Save(tasks)
		}
	}
	
	return fmt.Errorf("task with ID %s not found", id)
}

// generateID creates a simple ID for tasks
func generateID() string {
	// In a real implementation, this would use a proper UUID generator
	// For now, using a simple timestamp-based approach
	return time.Now().Format("20060102150405")
}
