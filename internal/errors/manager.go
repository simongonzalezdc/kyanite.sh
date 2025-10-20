package errors

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/logging"
)

// ErrorHandler defines the interface for error handling
type ErrorHandler interface {
	HandleError(ctx context.Context, err error) *ErrorReport
	HandleErrorWithRecovery(ctx context.Context, err error, recoveryFunc RecoveryFunc) *ErrorReport
	LogError(err error)
	ReportError(err error) *ErrorReport
}

// RecoveryFunc represents a function that can attempt to recover from an error
type RecoveryFunc func(ctx context.Context, err error) error

// ErrorReport contains information about an error occurrence
type ErrorReport struct {
	Error     *AppError         `json:"error"`
	Handled   bool              `json:"handled"`
	Recovered bool              `json:"recovered"`
	Recovery  *RecoveryAttempt  `json:"recovery,omitempty"`
	Context   map[string]string `json:"context,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// RecoveryAttempt contains information about a recovery attempt
type RecoveryAttempt struct {
	Strategy    ErrorRecovery `json:"strategy"`
	Attempts    int           `json:"attempts"`
	MaxAttempts int           `json:"max_attempts"`
	Success     bool          `json:"success"`
	Duration    time.Duration `json:"duration"`
	Result      string        `json:"result,omitempty"`
}

// ErrorManager manages error handling, logging, and recovery
type ErrorManager struct {
	logger    *logging.Logger
	config    *ErrorConfig
	reporters []ErrorReporter
	mu        sync.RWMutex
}

// ErrorConfig contains configuration for error handling
type ErrorConfig struct {
	// Logging configuration
	LogLevel     logging.LogLevel `json:"log_level"`
	LogToFile    bool             `json:"log_to_file"`
	LogDirectory string           `json:"log_directory"`
	MaxLogSize   int64            `json:"max_log_size"` // in bytes
	MaxLogFiles  int              `json:"max_log_files"`

	// Error reporting configuration
	EnableReporting bool     `json:"enable_reporting"`
	ReportURL       string   `json:"report_url,omitempty"`
	ReportFilters   []string `json:"report_filters,omitempty"`

	// Recovery configuration
	DefaultRetries  int           `json:"default_retries"`
	RetryDelay      time.Duration `json:"retry_delay"`
	MaxRecoveryTime time.Duration `json:"max_recovery_time"`

	// User notification configuration
	NotifyUser     bool `json:"notify_user"`
	ShowStackTrace bool `json:"show_stack_trace"`
}

// ErrorReporter defines the interface for error reporting
type ErrorReporter interface {
	Report(ctx context.Context, report *ErrorReport) error
	Name() string
}

// DefaultErrorConfig returns a default error configuration
func DefaultErrorConfig() *ErrorConfig {
	return &ErrorConfig{
		LogLevel:        logging.INFO,
		LogToFile:       true,
		LogDirectory:    "./logs",
		MaxLogSize:      10 * 1024 * 1024, // 10MB
		MaxLogFiles:     5,
		EnableReporting: false,
		DefaultRetries:  3,
		RetryDelay:      time.Second,
		MaxRecoveryTime: time.Minute,
		NotifyUser:      true,
		ShowStackTrace:  false,
	}
}

// NewErrorManager creates a new error manager
func NewErrorManager(logger *logging.Logger, config *ErrorConfig) *ErrorManager {
	if config == nil {
		config = DefaultErrorConfig()
	}

	manager := &ErrorManager{
		logger:    logger,
		config:    config,
		reporters: make([]ErrorReporter, 0),
	}

	// Create log directory if needed
	if config.LogToFile {
		if err := os.MkdirAll(config.LogDirectory, 0755); err != nil {
			logger.Errorf("Failed to create log directory: %v", err)
		}
	}

	return manager
}

// AddReporter adds an error reporter to the manager
func (em *ErrorManager) AddReporter(reporter ErrorReporter) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.reporters = append(em.reporters, reporter)
}

// HandleError handles an error with logging and reporting
func (em *ErrorManager) HandleError(ctx context.Context, err error) *ErrorReport {
	report := &ErrorReport{
		Error:     em.wrapError(err),
		Handled:   true,
		Recovered: false,
		Timestamp: time.Now(),
	}

	// Log the error
	em.logError(report.Error)

	// Report the error if configured
	if em.config.EnableReporting {
		em.reportError(ctx, report)
	}

	return report
}

// HandleErrorWithRecovery handles an error with recovery attempt
func (em *ErrorManager) HandleErrorWithRecovery(ctx context.Context, err error, recoveryFunc RecoveryFunc) *ErrorReport {
	report := em.HandleError(ctx, err)

	if recoveryFunc != nil && report.Error.CanRecover(report.Error.Recovery) {
		recovery := em.attemptRecovery(ctx, report.Error, recoveryFunc)
		report.Recovery = recovery
		report.Recovered = recovery.Success
	}

	return report
}

// LogError logs an error without creating a full report
func (em *ErrorManager) LogError(err error) {
	em.logError(em.wrapError(err))
}

// ReportError creates a report for an error without handling it
func (em *ErrorManager) ReportError(err error) *ErrorReport {
	report := &ErrorReport{
		Error:     em.wrapError(err),
		Handled:   false,
		Recovered: false,
		Timestamp: time.Now(),
	}

	if em.config.EnableReporting {
		em.reportError(context.Background(), report)
	}

	return report
}

// wrapError converts any error to an AppError
func (em *ErrorManager) wrapError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	// Try to categorize the error based on its type and message
	return em.categorizeError(err)
}

// categorizeError attempts to categorize an unknown error
func (em *ErrorManager) categorizeError(err error) *AppError {
	if err == nil {
		return nil
	}

	msg := err.Error()

	// Categorize based on error message patterns
	switch {
	case contains(msg, "database", "sql", "connection"):
		return NewDatabaseError("unknown", err)
	case contains(msg, "file", "open", "read", "write", "permission"):
		return NewFileError("unknown", "unknown", err)
	case contains(msg, "parse", "syntax", "invalid", "malformed"):
		return NewParsingError("unknown", err)
	case contains(msg, "network", "connection", "timeout"):
		return NewAppError("NETWORK_ERROR", "Network operation failed", err, CategoryNetwork, SeverityMedium, RecoveryRetry)
	default:
		return NewAppError("UNKNOWN_ERROR", "An unexpected error occurred", err, CategoryUnknown, SeverityMedium, RecoveryNone)
	}
}

// logError logs an error based on its severity
func (em *ErrorManager) logError(err *AppError) {
	if err == nil {
		return
	}

	logMsg := fmt.Sprintf("[%s] %s", err.Category, err.Message)
	if err.Operation != "" {
		logMsg = fmt.Sprintf("%s (operation: %s)", logMsg, err.Operation)
	}

	// Log based on severity
	switch err.Severity {
	case SeverityCritical:
		em.logger.Error(logMsg, "error", err.Error(), "category", err.Category, "severity", err.Severity)
	case SeverityHigh:
		em.logger.Warn(logMsg, "error", err.Error(), "category", err.Category, "severity", err.Severity)
	case SeverityMedium:
		em.logger.Info(logMsg, "error", err.Error(), "category", err.Category, "severity", err.Severity)
	case SeverityLow:
		em.logger.Debug(logMsg, "error", err.Error(), "category", err.Category, "severity", err.Severity)
	default:
		em.logger.Debug(logMsg, "error", err.Error(), "category", err.Category)
	}

	// Log stack trace if configured and available
	if em.config.ShowStackTrace && err.StackTrace != "" {
		em.logger.Debug("Stack trace", "stack", err.StackTrace)
	}
}

// reportError sends the error to all configured reporters
func (em *ErrorManager) reportError(ctx context.Context, report *ErrorReport) {
	em.mu.RLock()
	reporters := make([]ErrorReporter, len(em.reporters))
	copy(reporters, em.reporters)
	em.mu.RUnlock()

	for _, reporter := range reporters {
		if em.shouldReport(report.Error) {
			if err := reporter.Report(ctx, report); err != nil {
				em.logger.Errorf("Failed to report error to %s: %v", reporter.Name(), err)
			}
		}
	}
}

// shouldReport determines if an error should be reported based on filters
func (em *ErrorManager) shouldReport(err *AppError) bool {
	if len(em.config.ReportFilters) == 0 {
		return true
	}

	for _, filter := range em.config.ReportFilters {
		if err.Category == ErrorCategory(filter) {
			return true
		}
	}
	return false
}

// attemptRecovery attempts to recover from an error using the provided function
func (em *ErrorManager) attemptRecovery(ctx context.Context, err *AppError, recoveryFunc RecoveryFunc) *RecoveryAttempt {
	start := time.Now()
	attempt := &RecoveryAttempt{
		Strategy:    err.Recovery,
		MaxAttempts: err.MaxRetries,
	}

	// Check if we've exceeded maximum recovery time
	ctx, cancel := context.WithTimeout(ctx, em.config.MaxRecoveryTime)
	defer cancel()

	for attempt.Attempts < attempt.MaxAttempts {
		attempt.Attempts++

		select {
		case <-ctx.Done():
			attempt.Result = "Recovery timeout"
			attempt.Duration = time.Since(start)
			return attempt
		default:
		}

		// Attempt recovery
		recoveryErr := recoveryFunc(ctx, err)
		if recoveryErr == nil {
			attempt.Success = true
			attempt.Duration = time.Since(start)
			attempt.Result = "Recovery successful"
			return attempt
		}

		// Log failed recovery attempt
		em.logger.Debugf("Recovery attempt %d failed: %v", attempt.Attempts, recoveryErr)

		// Wait before retry (except on last attempt)
		if attempt.Attempts < attempt.MaxAttempts {
			select {
			case <-ctx.Done():
				attempt.Result = "Recovery timeout during retry delay"
				attempt.Duration = time.Since(start)
				return attempt
			case <-time.After(em.config.RetryDelay):
			}
		}
	}

	attempt.Duration = time.Since(start)
	attempt.Result = "Max recovery attempts exceeded"
	return attempt
}

// contains checks if a string contains any of the given substrings (case-insensitive)
func contains(s string, substrings ...string) bool {
	// Convert to lowercase for case-insensitive comparison
	lowerS := strings.ToLower(s)
	for _, substr := range substrings {
		if substr == "" {
			continue
		}
		lowerSubstr := strings.ToLower(substr)
		// Simple contains check (in production, consider using strings.Contains with strings.ToLower)
		// For now, we'll do a basic check
		if len(lowerS) >= len(lowerSubstr) {
			for i := 0; i <= len(lowerS)-len(lowerSubstr); i++ {
				match := true
				for j := 0; j < len(lowerSubstr); j++ {
					if lowerS[i+j] != lowerSubstr[j] {
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
		}
	}
	return false
}

// GetErrorStats returns statistics about error handling
func (em *ErrorManager) GetErrorStats() map[string]interface{} {
	// In a full implementation, this would track and return error statistics
	return map[string]interface{}{
		"total_errors":      0,
		"critical_errors":   0,
		"recovery_success":  0,
		"recovery_failures": 0,
	}
}

// Close gracefully shuts down the error manager
func (em *ErrorManager) Close() error {
	em.mu.Lock()
	defer em.mu.Unlock()

	// Close all reporters
	for _, reporter := range em.reporters {
		if closer, ok := reporter.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				em.logger.Errorf("Error closing reporter %s: %v", reporter.Name(), err)
			}
		}
	}

	em.reporters = nil
	return nil
}
