package errors

import (
	"fmt"
	"time"
)

// ErrorCategory represents different categories of errors
type ErrorCategory string

const (
	// CategoryValidation represents input validation errors
	CategoryValidation ErrorCategory = "validation"
	// CategoryDatabase represents database operation errors
	CategoryDatabase ErrorCategory = "database"
	// CategoryFile represents file I/O operation errors
	CategoryFile ErrorCategory = "file"
	// CategoryNetwork represents network operation errors
	CategoryNetwork ErrorCategory = "network"
	// CategoryUI represents user interface errors
	CategoryUI ErrorCategory = "ui"
	// CategoryConfiguration represents configuration errors
	CategoryConfiguration ErrorCategory = "configuration"
	// CategoryPermission represents permission/access errors
	CategoryPermission ErrorCategory = "permission"
	// CategoryResource represents resource limit errors
	CategoryResource ErrorCategory = "resource"
	// CategoryParsing represents parsing/serialization errors
	CategoryParsing ErrorCategory = "parsing"
	// CategoryUnknown represents unknown/unexpected errors
	CategoryUnknown ErrorCategory = "unknown"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity string

const (
	// SeverityLow represents minor errors that don't affect functionality
	SeverityLow ErrorSeverity = "low"
	// SeverityMedium represents errors that affect some functionality
	SeverityMedium ErrorSeverity = "medium"
	// SeverityHigh represents errors that affect core functionality
	SeverityHigh ErrorSeverity = "high"
	// SeverityCritical represents errors that prevent application from running
	SeverityCritical ErrorSeverity = "critical"
)

// ErrorRecovery represents recovery strategies for errors
type ErrorRecovery string

const (
	// RecoveryNone indicates no recovery is possible
	RecoveryNone ErrorRecovery = "none"
	// RecoveryRetry indicates the operation can be retried
	RecoveryRetry ErrorRecovery = "retry"
	// RecoveryFallback indicates a fallback option is available
	RecoveryFallback ErrorRecovery = "fallback"
	// RecoveryGraceful indicates graceful degradation is possible
	RecoveryGraceful ErrorRecovery = "graceful"
	// RecoveryManual indicates manual intervention is required
	RecoveryManual ErrorRecovery = "manual"
)

// AppError represents a structured application error
type AppError struct {
	// Core error information
	Code     string        `json:"code"`
	Message  string        `json:"message"`
	Cause    error         `json:"-"` // Original error (not serialized)
	Category ErrorCategory `json:"category"`
	Severity ErrorSeverity `json:"severity"`
	Recovery ErrorRecovery `json:"recovery"`

	// Context information
	Operation string            `json:"operation,omitempty"`
	Component string            `json:"component,omitempty"`
	UserID    string            `json:"user_id,omitempty"`
	SessionID string            `json:"session_id,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`

	// Timing information
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration,omitempty"`

	// Recovery information
	RecoveryAttempts int `json:"recovery_attempts"`
	MaxRetries       int `json:"max_retries"`

	// Stack trace (for debugging)
	StackTrace string `json:"stack_trace,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// Is implements error comparison for errors.Is()
func (e *AppError) Is(target error) bool {
	if appErr, ok := target.(*AppError); ok {
		return e.Code == appErr.Code && e.Category == appErr.Category
	}
	return false
}

// HasCategory checks if the error belongs to a specific category
func (e *AppError) HasCategory(category ErrorCategory) bool {
	return e.Category == category
}

// HasSeverity checks if the error has at least the specified severity
func (e *AppError) HasSeverity(severity ErrorSeverity) bool {
	return getSeverityLevel(e.Severity) >= getSeverityLevel(severity)
}

// CanRecover checks if the error supports the specified recovery strategy
func (e *AppError) CanRecover(strategy ErrorRecovery) bool {
	return e.Recovery == strategy || e.Recovery == RecoveryFallback || e.Recovery == RecoveryGraceful
}

// ShouldRetry checks if the error can be retried
func (e *AppError) ShouldRetry() bool {
	return e.CanRecover(RecoveryRetry) && e.RecoveryAttempts < e.MaxRetries
}

// IncrementRetry increments the recovery attempt counter
func (e *AppError) IncrementRetry() {
	e.RecoveryAttempts++
}

// IsRetryable checks if the error can be retried based on current attempts
func (e *AppError) IsRetryable() bool {
	return e.ShouldRetry() && e.RecoveryAttempts < e.MaxRetries
}

// WithContext adds context information to the error
func (e *AppError) WithContext(key, value string) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}

// WithOperation sets the operation context
func (e *AppError) WithOperation(operation string) *AppError {
	e.Operation = operation
	return e
}

// WithComponent sets the component context
func (e *AppError) WithComponent(component string) *AppError {
	e.Component = component
	return e
}

// WithDuration sets the operation duration
func (e *AppError) WithDuration(duration time.Duration) *AppError {
	e.Duration = duration
	return e
}

// getSeverityLevel converts severity to numeric level for comparison
func getSeverityLevel(severity ErrorSeverity) int {
	switch severity {
	case SeverityLow:
		return 1
	case SeverityMedium:
		return 2
	case SeverityHigh:
		return 3
	case SeverityCritical:
		return 4
	default:
		return 0
	}
}

// NewAppError creates a new application error
func NewAppError(code, message string, cause error, category ErrorCategory, severity ErrorSeverity, recovery ErrorRecovery) *AppError {
	return &AppError{
		Code:             code,
		Message:          message,
		Cause:            cause,
		Category:         category,
		Severity:         severity,
		Recovery:         recovery,
		Timestamp:        time.Now(),
		RecoveryAttempts: 0,
		MaxRetries:       getDefaultMaxRetries(recovery),
	}
}

// getDefaultMaxRetries returns default retry count based on recovery strategy
func getDefaultMaxRetries(recovery ErrorRecovery) int {
	switch recovery {
	case RecoveryRetry:
		return 3
	case RecoveryFallback:
		return 1
	case RecoveryGraceful:
		return 2
	default:
		return 0
	}
}

// Common error constructors for different categories

// NewValidationError creates a validation error
func NewValidationError(message string, cause error) *AppError {
	return NewAppError("VALIDATION_FAILED", message, cause, CategoryValidation, SeverityMedium, RecoveryManual)
}

// NewDatabaseError creates a database operation error
func NewDatabaseError(operation string, cause error) *AppError {
	return NewAppError("DATABASE_ERROR", fmt.Sprintf("Database operation '%s' failed", operation), cause, CategoryDatabase, SeverityHigh, RecoveryRetry)
}

// NewFileError creates a file I/O error
func NewFileError(operation, filepath string, cause error) *AppError {
	return NewAppError("FILE_ERROR", fmt.Sprintf("File operation '%s' failed for '%s'", operation, filepath), cause, CategoryFile, SeverityHigh, RecoveryRetry).
		WithOperation(operation).
		WithContext("filepath", filepath)
}

// NewPermissionError creates a permission/access error
func NewPermissionError(resource string, cause error) *AppError {
	return NewAppError("PERMISSION_DENIED", fmt.Sprintf("Permission denied for resource: %s", resource), cause, CategoryPermission, SeverityMedium, RecoveryManual).
		WithContext("resource", resource)
}

// NewParsingError creates a parsing/serialization error
func NewParsingError(operation string, cause error) *AppError {
	return NewAppError("PARSING_ERROR", fmt.Sprintf("Failed to parse %s", operation), cause, CategoryParsing, SeverityHigh, RecoveryManual).
		WithOperation(operation)
}

// NewResourceError creates a resource limit error
func NewResourceError(resource string, cause error) *AppError {
	return NewAppError("RESOURCE_LIMIT", fmt.Sprintf("Resource limit exceeded for: %s", resource), cause, CategoryResource, SeverityMedium, RecoveryGraceful).
		WithContext("resource", resource)
}

// NewConfigurationError creates a configuration error
func NewConfigurationError(setting string, cause error) *AppError {
	return NewAppError("CONFIG_ERROR", fmt.Sprintf("Configuration error for setting: %s", setting), cause, CategoryConfiguration, SeverityCritical, RecoveryManual).
		WithContext("setting", setting)
}

// NewUIError creates a user interface error
func NewUIError(component string, cause error) *AppError {
	return NewAppError("UI_ERROR", fmt.Sprintf("UI error in component: %s", component), cause, CategoryUI, SeverityMedium, RecoveryGraceful).
		WithComponent(component)
}
