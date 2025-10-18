package errors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/puente-labs/noise/internal/logging"
)

// FileErrorReporter reports errors to a file
type FileErrorReporter struct {
	filePath string
	file     *os.File
	logger   *logging.Logger
	mu       sync.Mutex
}

// NewFileErrorReporter creates a new file-based error reporter
func NewFileErrorReporter(filePath string, logger *logging.Logger) (*FileErrorReporter, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create error log directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open error log file: %w", err)
	}

	return &FileErrorReporter{
		filePath: filePath,
		file:     file,
		logger:   logger,
	}, nil
}

// Report writes an error report to the file
func (r *FileErrorReporter) Report(ctx context.Context, report *ErrorReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Convert report to JSON
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal error report: %w", err)
	}

	// Write timestamp and report
	timestamp := time.Now().Format(time.RFC3339)
	if _, err := fmt.Fprintf(r.file, "[%s] Error Report:\n%s\n\n", timestamp, jsonData); err != nil {
		return fmt.Errorf("failed to write error report: %w", err)
	}

	// Flush to ensure it's written immediately
	if err := r.file.Sync(); err != nil {
		r.logger.Warn("Failed to sync error log file", "error", err)
	}

	return nil
}

// Name returns the reporter name
func (r *FileErrorReporter) Name() string {
	return "file"
}

// Close closes the file
func (r *FileErrorReporter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.file != nil {
		err := r.file.Close()
		r.file = nil
		return err
	}
	return nil
}

// ConsoleErrorReporter reports errors to the console
type ConsoleErrorReporter struct {
	writer io.Writer
	logger *logging.Logger
}

// NewConsoleErrorReporter creates a new console-based error reporter
func NewConsoleErrorReporter(writer io.Writer, logger *logging.Logger) *ConsoleErrorReporter {
	return &ConsoleErrorReporter{
		writer: writer,
		logger: logger,
	}
}

// Report writes an error report to the console
func (r *ConsoleErrorReporter) Report(ctx context.Context, report *ErrorReport) error {
	// Only report high severity errors to console to avoid spam
	if report.Error.Severity == SeverityLow {
		return nil
	}

	message := fmt.Sprintf("[%s] %s: %s", report.Error.Category, report.Error.Severity, report.Error.Message)

	if report.Error.Operation != "" {
		message += fmt.Sprintf(" (operation: %s)", report.Error.Operation)
	}

	if report.Error.Component != "" {
		message += fmt.Sprintf(" (component: %s)", report.Error.Component)
	}

	fmt.Fprintf(r.writer, "ERROR: %s\n", message)

	// Include recovery information if available
	if report.Recovery != nil {
		fmt.Fprintf(r.writer, "  Recovery: %s (attempts: %d/%d, success: %v)\n",
			report.Recovery.Strategy,
			report.Recovery.Attempts,
			report.Recovery.MaxAttempts,
			report.Recovery.Success,
		)
	}

	// Include context information
	if len(report.Context) > 0 {
		fmt.Fprintf(r.writer, "  Context: %v\n", report.Context)
	}

	fmt.Fprintf(r.writer, "\n")

	return nil
}

// Name returns the reporter name
func (r *ConsoleErrorReporter) Name() string {
	return "console"
}

// ExternalErrorReporter reports errors to an external service
type ExternalErrorReporter struct {
	endpoint string
	apiKey   string
	logger   *logging.Logger
	client   *http.Client
}

// NewExternalErrorReporter creates a new external service error reporter
func NewExternalErrorReporter(endpoint, apiKey string, logger *logging.Logger) *ExternalErrorReporter {
	return &ExternalErrorReporter{
		endpoint: endpoint,
		apiKey:   apiKey,
		logger:   logger,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Report sends an error report to an external service
func (r *ExternalErrorReporter) Report(ctx context.Context, report *ErrorReport) error {
	// Only report critical and high severity errors to external services
	if report.Error.Severity != SeverityCritical && report.Error.Severity != SeverityHigh {
		return nil
	}

	// Prepare payload
	payload := map[string]interface{}{
		"error":     report.Error,
		"handled":   report.Handled,
		"recovered": report.Recovered,
		"recovery":  report.Recovery,
		"context":   report.Context,
		"timestamp": report.Timestamp,
		"version":   "1.0.0", // Application version
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal error report for external service: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("User-Agent", "noise.sh-ErrorReporter/1.0")

	// Send request
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send error report: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("external service returned error: %s (status: %d)", string(body), resp.StatusCode)
	}

	r.logger.Debug("Error report sent to external service successfully", "status", resp.StatusCode)
	return nil
}

// Name returns the reporter name
func (r *ExternalErrorReporter) Name() string {
	return "external"
}

// Close closes the reporter
func (r *ExternalErrorReporter) Close() error {
	// Nothing to close for HTTP client
	return nil
}

// CompositeErrorReporter combines multiple reporters
type CompositeErrorReporter struct {
	reporters []ErrorReporter
	logger    *logging.Logger
}

// NewCompositeErrorReporter creates a new composite error reporter
func NewCompositeErrorReporter(reporters []ErrorReporter, logger *logging.Logger) *CompositeErrorReporter {
	return &CompositeErrorReporter{
		reporters: reporters,
		logger:    logger,
	}
}

// Report sends the error report to all configured reporters
func (r *CompositeErrorReporter) Report(ctx context.Context, report *ErrorReport) error {
	var errors []error

	for _, reporter := range r.reporters {
		if err := reporter.Report(ctx, report); err != nil {
			r.logger.Warn("Reporter failed to report error", "reporter", reporter.Name(), "error", err)
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple reporters failed: %v", errors)
	}

	return nil
}

// Name returns the reporter name
func (r *CompositeErrorReporter) Name() string {
	return "composite"
}

// Close closes all reporters
func (r *CompositeErrorReporter) Close() error {
	var errors []error

	for _, reporter := range r.reporters {
		if closer, ok := reporter.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				r.logger.Warn("Failed to close reporter", "reporter", reporter.Name(), "error", err)
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple reporters failed to close: %v", errors)
	}

	return nil
}

// AddReporter adds a reporter to the composite
func (r *CompositeErrorReporter) AddReporter(reporter ErrorReporter) {
	r.reporters = append(r.reporters, reporter)
}

// RemoveReporter removes a reporter from the composite
func (r *CompositeErrorReporter) RemoveReporter(name string) {
	for i, reporter := range r.reporters {
		if reporter.Name() == name {
			r.reporters = append(r.reporters[:i], r.reporters[i+1:]...)
			break
		}
	}
}

// GetReporters returns all configured reporters
func (r *CompositeErrorReporter) GetReporters() []ErrorReporter {
	return r.reporters
}

// ErrorMetrics tracks error statistics
type ErrorMetrics struct {
	TotalErrors      int64            `json:"total_errors"`
	ErrorsByCategory map[string]int64 `json:"errors_by_category"`
	ErrorsBySeverity map[string]int64 `json:"errors_by_severity"`
	RecoverySuccess  int64            `json:"recovery_success"`
	RecoveryFailure  int64            `json:"recovery_failure"`
	LastErrorTime    time.Time        `json:"last_error_time"`
	mu               sync.RWMutex
}

// NewErrorMetrics creates a new error metrics tracker
func NewErrorMetrics() *ErrorMetrics {
	return &ErrorMetrics{
		ErrorsByCategory: make(map[string]int64),
		ErrorsBySeverity: make(map[string]int64),
	}
}

// RecordError records an error in the metrics
func (m *ErrorMetrics) RecordError(err *AppError) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalErrors++
	m.ErrorsByCategory[string(err.Category)]++
	m.ErrorsBySeverity[string(err.Severity)]++
	m.LastErrorTime = time.Now()
}

// RecordRecovery records a recovery attempt in the metrics
func (m *ErrorMetrics) RecordRecovery(success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if success {
		m.RecoverySuccess++
	} else {
		m.RecoveryFailure++
	}
}

// GetMetrics returns a copy of the current metrics
func (m *ErrorMetrics) GetMetrics() ErrorMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return ErrorMetrics{
		TotalErrors:      m.TotalErrors,
		ErrorsByCategory: m.copyMap(m.ErrorsByCategory),
		ErrorsBySeverity: m.copyMap(m.ErrorsBySeverity),
		RecoverySuccess:  m.RecoverySuccess,
		RecoveryFailure:  m.RecoveryFailure,
		LastErrorTime:    m.LastErrorTime,
	}
}

// copyMap creates a copy of a string to int64 map
func (m *ErrorMetrics) copyMap(original map[string]int64) map[string]int64 {
	copy := make(map[string]int64)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

// Global error metrics instance
var globalErrorMetrics *ErrorMetrics

// init initializes the global error metrics
func init() {
	globalErrorMetrics = NewErrorMetrics()
}

// GetGlobalErrorMetrics returns the global error metrics
func GetGlobalErrorMetrics() *ErrorMetrics {
	return globalErrorMetrics
}
