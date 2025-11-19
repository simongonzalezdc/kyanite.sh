package errors

import (
	"errors"
	"fmt"
)

// Common error types
var (
	// Task errors
	ErrTaskNotFound        = errors.New("task not found")
	ErrTaskAlreadyExists   = errors.New("task already exists")
	ErrEmptyDescription    = errors.New("task description cannot be empty")
	ErrInvalidPriority     = errors.New("invalid priority level")
	ErrInvalidStatus       = errors.New("invalid task status")
	ErrInvalidDeadline     = errors.New("invalid deadline format")

	// Template errors
	ErrTemplateNotFound      = errors.New("template not found")
	ErrTemplateAlreadyExists = errors.New("template already exists")
	ErrEmptyTemplateName     = errors.New("template name cannot be empty")

	// Recurrence errors
	ErrInvalidRecurrence = errors.New("invalid recurrence pattern")
	ErrRecurrenceInPast  = errors.New("recurrence end date is in the past")

	// Storage errors
	ErrStorageRead  = errors.New("failed to read from storage")
	ErrStorageWrite = errors.New("failed to write to storage")
	ErrStorageInit  = errors.New("failed to initialize storage")

	// AI errors
	ErrAINotAvailable = errors.New("AI service not available")
	ErrAIParseFailed  = errors.New("AI failed to parse task")
	ErrAITimeout      = errors.New("AI request timed out")

	// Command errors
	ErrNothingToUndo = errors.New("nothing to undo")
	ErrNothingToRedo = errors.New("nothing to redo")
)

// TaskError wraps task-related errors with context
type TaskError struct {
	Op      string // Operation that failed
	TaskID  string // Task ID if applicable
	Err     error  // Underlying error
	Message string // Human-readable message
}

func (e *TaskError) Error() string {
	if e.TaskID != "" {
		return fmt.Sprintf("%s failed for task %s: %v", e.Op, e.TaskID, e.Err)
	}
	return fmt.Sprintf("%s failed: %v", e.Op, e.Err)
}

func (e *TaskError) Unwrap() error {
	return e.Err
}

// NewTaskError creates a new task error
func NewTaskError(op, taskID string, err error) *TaskError {
	return &TaskError{
		Op:     op,
		TaskID: taskID,
		Err:    err,
	}
}

// ValidationError represents input validation errors
type ValidationError struct {
	Field   string // Field that failed validation
	Value   string // Invalid value
	Message string // Human-readable error message
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
	}
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, value, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// StorageError wraps storage-related errors
type StorageError struct {
	Op   string // Operation that failed
	Path string // File path if applicable
	Err  error  // Underlying error
}

func (e *StorageError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s failed for %s: %v", e.Op, e.Path, e.Err)
	}
	return fmt.Sprintf("%s failed: %v", e.Op, e.Err)
}

func (e *StorageError) Unwrap() error {
	return e.Err
}

// NewStorageError creates a new storage error
func NewStorageError(op, path string, err error) *StorageError {
	return &StorageError{
		Op:   op,
		Path: path,
		Err:  err,
	}
}

// IsNotFound returns true if the error is a "not found" error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrTaskNotFound) ||
		errors.Is(err, ErrTemplateNotFound)
}

// IsValidationError returns true if the error is a validation error
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsStorageError returns true if the error is a storage error
func IsStorageError(err error) bool {
	var storageErr *StorageError
	return errors.As(err, &storageErr)
}
