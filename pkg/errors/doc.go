// Package errors provides structured error types and error handling utilities.
//
// This package defines common error types used throughout the application
// and provides structured error wrapping with context information.
//
// Error Types:
//   - TaskError: Wraps task-related errors with operation context
//   - ValidationError: Input validation failures
//   - StorageError: File system and persistence errors
//
// Common Errors:
//   - Task errors: ErrTaskNotFound, ErrEmptyDescription, etc.
//   - Template errors: ErrTemplateNotFound, etc.
//   - AI errors: ErrAINotAvailable, ErrAITimeout, etc.
//   - Command errors: ErrNothingToUndo, ErrNothingToRedo
//
// Example usage:
//
//	if err := someOperation(); err != nil {
//	    if errors.IsNotFound(err) {
//	        // Handle not found
//	    }
//	    return NewTaskError("operation", taskID, err)
//	}
package errors
