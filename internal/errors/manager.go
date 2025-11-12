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

// ErrorContext contains contextual information for an error
type ErrorContext struct {
	RequestID string            `json:"request_id,omitempty"`
	UserID    string            `json:"user_id,omitempty"`
	SessionID string            `json:"session_id,omitempty"`
	TraceID   string            `json:"trace_id,omitempty"`
	Component string            `json:"component,omitempty"`
	Operation string            `json:"operation,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// RecoveryStats tracks statistics for recovery attempts
type RecoveryStats struct {
	Strategy     ErrorRecovery `json:"strategy"`
	SuccessCount int           `json:"success_count"`
	FailureCount int           `json:"failure_count"`
	LastAttempt  time.Time     `json:"last_attempt"`
}

// ErrorManager manages error handling, logging, and recovery
type ErrorManager struct {
	logger    *logging.Logger
	config    *ErrorConfig
	reporters []ErrorReporter
	mu        sync.RWMutex

	// Enhanced error tracking
	errorStats   map[string]int
	errorHistory []ErrorReport
	historyLimit int
	contextStore map[string]*ErrorContext

	// Performance monitoring
	recoveryStats map[ErrorRecovery]*RecoveryStats
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
		logger:        logger,
		config:        config,
		reporters:     make([]ErrorReporter, 0),
		errorStats:    make(map[string]int),
		errorHistory:  make([]ErrorReport, 0),
		historyLimit:  1000,
		contextStore:  make(map[string]*ErrorContext),
		recoveryStats: make(map[ErrorRecovery]*RecoveryStats),
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
	appErr := em.wrapError(err)

	// If error is nil, return a report with nil error
	if appErr == nil {
		return &ErrorReport{
			Error:     nil,
			Handled:   true,
			Recovered: false,
			Context:   make(map[string]string),
			Timestamp: time.Now(),
		}
	}

	// Update error statistics
	em.updateErrorStats(appErr.Category)

	// Extract context information
	errorContext := em.extractContext(ctx, appErr)

	report := &ErrorReport{
		Error:     appErr,
		Handled:   true,
		Recovered: false,
		Context:   errorContext,
		Timestamp: time.Now(),
	}

	// Store error in history
	em.storeErrorInHistory(*report)

	// Log the error with structured fields
	em.logErrorWithContext(report.Error, errorContext)

	// Report the error if configured
	if em.config.EnableReporting {
		em.reportError(ctx, report)
	}

	return report
}

// HandleErrorWithRecovery handles an error with recovery attempt
func (em *ErrorManager) HandleErrorWithRecovery(ctx context.Context, err error, recoveryFunc RecoveryFunc) *ErrorReport {
	report := em.HandleError(ctx, err)

	if recoveryFunc != nil && report.Error != nil && report.Error.CanRecover(report.Error.Recovery) {
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
	em.mu.RLock()
	defer em.mu.RUnlock()

	stats := map[string]interface{}{
		"total_errors":      0,
		"critical_errors":   0,
		"high_errors":       0,
		"medium_errors":     0,
		"low_errors":        0,
		"recovery_success":  0,
		"recovery_failures": 0,
		"error_categories":  make(map[string]int),
		"recent_errors":     0,
		"error_rate":        0.0,
	}

	// Count errors by category and severity
	categoryCounts := make(map[string]int)
	severityCounts := map[string]int{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
	}

	for category, count := range em.errorStats {
		stats["total_errors"] = stats["total_errors"].(int) + count
		categoryCounts[string(category)] = count
		stats["error_categories"].(map[string]int)[string(category)] = count
	}

	// Count recent errors (last hour)
	recentCount := 0
	now := time.Now()
	oneHourAgo := now.Add(-time.Hour)
	for _, report := range em.errorHistory {
		if report.Timestamp.After(oneHourAgo) {
			recentCount++
		}

		// Count by severity (skip if error is nil)
		if report.Error != nil {
			switch report.Error.Severity {
			case SeverityCritical:
				severityCounts["critical"]++
				stats["critical_errors"] = stats["critical_errors"].(int) + 1
			case SeverityHigh:
				severityCounts["high"]++
				stats["high_errors"] = stats["high_errors"].(int) + 1
			case SeverityMedium:
				severityCounts["medium"]++
				stats["medium_errors"] = stats["medium_errors"].(int) + 1
			case SeverityLow:
				severityCounts["low"]++
				stats["low_errors"] = stats["low_errors"].(int) + 1
			}
		}
	}

	stats["recent_errors"] = recentCount

	// Calculate error rate (errors per hour)
	hours := 24.0 // Default to 24 hours for now
	if len(em.errorHistory) > 0 {
		earliest := em.errorHistory[0].Timestamp
		elapsed := now.Sub(earliest).Hours()
		if elapsed > 0 {
			hours = elapsed
		}
	}

	if hours > 0 {
		stats["error_rate"] = float64(stats["total_errors"].(int)) / hours
	}

	// Add recovery statistics
	for _, recStats := range em.recoveryStats {
		stats["recovery_success"] = stats["recovery_success"].(int) + recStats.SuccessCount
		stats["recovery_failures"] = stats["recovery_failures"].(int) + recStats.FailureCount
	}

	return stats
}

// updateErrorStats updates error statistics
func (em *ErrorManager) updateErrorStats(category ErrorCategory) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.errorStats[string(category)]++
}

// storeErrorInHistory stores an error in the history
func (em *ErrorManager) storeErrorInHistory(report ErrorReport) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.errorHistory = append(em.errorHistory, report)

	// Limit history size
	if len(em.errorHistory) > em.historyLimit {
		em.errorHistory = em.errorHistory[1:]
	}
}

// extractContext extracts context information for an error
func (em *ErrorManager) extractContext(ctx context.Context, err *AppError) map[string]string {
	context := make(map[string]string)

	// Extract basic error information
	context["error_code"] = err.Code
	context["error_category"] = string(err.Category)
	context["error_severity"] = string(err.Severity)
	context["error_recovery"] = string(err.Recovery)

	if err.Operation != "" {
		context["operation"] = err.Operation
	}

	// Extract context values from the context if available
	if ctx != nil {
		if reqID := ctx.Value("request_id"); reqID != nil {
			if id, ok := reqID.(string); ok {
				context["request_id"] = id
			}
		}

		if userID := ctx.Value("user_id"); userID != nil {
			if id, ok := userID.(string); ok {
				context["user_id"] = id
			}
		}

		if sessionID := ctx.Value("session_id"); sessionID != nil {
			if id, ok := sessionID.(string); ok {
				context["session_id"] = id
			}
		}

		if traceID := ctx.Value("trace_id"); traceID != nil {
			if id, ok := traceID.(string); ok {
				context["trace_id"] = id
			}
		}
	}

	return context
}

// logErrorWithContext logs an error with structured context fields
func (em *ErrorManager) logErrorWithContext(err *AppError, context map[string]string) {
	if err == nil {
		return
	}

	// Create a slice of key-value pairs for structured logging
	fields := make([]interface{}, 0)
	for k, v := range context {
		fields = append(fields, k, v)
	}

	// Add error-specific fields
	fields = append(fields, "error", err.Error())
	fields = append(fields, "error_code", err.Code)
	fields = append(fields, "error_category", err.Category)
	fields = append(fields, "error_severity", err.Severity)

	if err.Operation != "" {
		fields = append(fields, "operation", err.Operation)
	}

	if err.StackTrace != "" {
		fields = append(fields, "stack_trace", err.StackTrace)
	}

	// Create log message
	logMsg := fmt.Sprintf("[%s] %s", err.Category, err.Message)
	if err.Operation != "" {
		logMsg = fmt.Sprintf("%s (operation: %s)", logMsg, err.Operation)
	}

	// Log based on severity with structured fields
	switch err.Severity {
	case SeverityCritical:
		em.logger.Error(logMsg, fields)
	case SeverityHigh:
		em.logger.Warn(logMsg, fields)
	case SeverityMedium:
		em.logger.Info(logMsg, fields)
	case SeverityLow:
		em.logger.Debug(logMsg, fields)
	default:
		em.logger.Debug(logMsg, fields)
	}
}

// UpdateRecoveryStats updates recovery statistics
func (em *ErrorManager) UpdateRecoveryStats(recovery ErrorRecovery, success bool) {
	em.mu.Lock()
	defer em.mu.Unlock()

	stats, exists := em.recoveryStats[recovery]
	if !exists {
		stats = &RecoveryStats{
			Strategy:     recovery,
			SuccessCount: 0,
			FailureCount: 0,
			LastAttempt:  time.Time{},
		}
		em.recoveryStats[recovery] = stats
	}

	stats.LastAttempt = time.Now()
	if success {
		stats.SuccessCount++
	} else {
		stats.FailureCount++
	}
}

// GetRecentErrors returns recent errors within the specified duration
func (em *ErrorManager) GetRecentErrors(duration time.Duration) []ErrorReport {
	em.mu.RLock()
	defer em.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var recent []ErrorReport

	for _, report := range em.errorHistory {
		if report.Timestamp.After(cutoff) {
			recent = append(recent, report)
		}
	}

	return recent
}

// GetErrorHistory returns the full error history
func (em *ErrorManager) GetErrorHistory() []ErrorReport {
	em.mu.RLock()
	defer em.mu.RUnlock()

	// Return a copy to prevent external modification
	history := make([]ErrorReport, len(em.errorHistory))
	copy(history, em.errorHistory)

	return history
}

// ClearErrorHistory clears the error history
func (em *ErrorManager) ClearErrorHistory() {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.errorHistory = make([]ErrorReport, 0)
	em.errorStats = make(map[string]int)
}

// SetErrorContext stores context information for future errors
func (em *ErrorManager) SetErrorContext(key string, context *ErrorContext) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.contextStore[key] = context
}

// GetErrorContext retrieves stored context information
func (em *ErrorManager) GetErrorContext(key string) *ErrorContext {
	em.mu.RLock()
	defer em.mu.RUnlock()

	return em.contextStore[key]
}

// ClearErrorContext clears stored context information
func (em *ErrorManager) ClearErrorContext(key string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	delete(em.contextStore, key)
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
