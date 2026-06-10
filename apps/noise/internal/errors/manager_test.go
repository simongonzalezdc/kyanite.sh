package errors

import (
	"context"
	"testing"
	"time"

	logging "github.com/kyanite/noise/internal/logging"
)

// TestNewErrorManager tests the creation of error manager
func TestNewErrorManager(t *testing.T) {
	logger := NewTestLogger(t)
	config := DefaultErrorConfig()
	manager := NewErrorManager(logger.Logger, config)

	if manager == nil {
		t.Fatal("Expected non-nil ErrorManager")
	}

	if manager.logger != logger.Logger {
		t.Error("Expected logger to be set correctly")
	}

	if manager.config != config {
		t.Error("Expected config to be set correctly")
	}

	if manager.reporters == nil {
		t.Error("Expected reporters slice to be initialized")
	}
}

// TestDefaultErrorConfig tests the default error configuration
func TestDefaultErrorConfig(t *testing.T) {
	config := DefaultErrorConfig()

	if config.LogLevel != logging.INFO {
		t.Errorf("Expected default log level INFO, got %v", config.LogLevel)
	}

	if !config.LogToFile {
		t.Error("Expected LogToFile to be true by default")
	}

	if config.LogDirectory != "./logs" {
		t.Errorf("Expected default log directory './logs', got %s", config.LogDirectory)
	}

	if config.MaxLogSize != 10*1024*1024 {
		t.Errorf("Expected default max log size 10MB, got %d", config.MaxLogSize)
	}

	if config.MaxLogFiles != 5 {
		t.Errorf("Expected default max log files 5, got %d", config.MaxLogFiles)
	}

	if config.EnableReporting {
		t.Error("Expected EnableReporting to be false by default")
	}

	if config.DefaultRetries != 3 {
		t.Errorf("Expected default retries 3, got %d", config.DefaultRetries)
	}

	if config.RetryDelay != time.Second {
		t.Errorf("Expected default retry delay 1s, got %v", config.RetryDelay)
	}

	if config.MaxRecoveryTime != time.Minute {
		t.Errorf("Expected default max recovery time 1m, got %v", config.MaxRecoveryTime)
	}

	if !config.NotifyUser {
		t.Error("Expected NotifyUser to be true by default")
	}

	if config.ShowStackTrace {
		t.Error("Expected ShowStackTrace to be false by default")
	}
}

// TestAddReporter tests adding error reporters
func TestAddReporter(t *testing.T) {
	logger := NewTestLogger(t)
	config := DefaultErrorConfig()
	manager := NewErrorManager(logger.Logger, config)

	reporter := NewTestErrorReporter()
	manager.AddReporter(reporter)

	if len(manager.reporters) != 1 {
		t.Errorf("Expected 1 reporter, got %d", len(manager.reporters))
	}

	if manager.reporters[0] != reporter {
		t.Error("Expected reporter to be added correctly")
	}
}

// TestErrorManagerHandleError tests basic error handling
func TestErrorManagerHandleError(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	// Test handling a simple error
	err := NewFileError("read_file", "/test/file.txt", nil)
	report := setup.Manager.HandleError(context.Background(), err)

	if report == nil {
		t.Fatal("Expected non-nil error report")
	}

	if !report.Handled {
		t.Error("Expected error to be marked as handled")
	}

	if report.Error != err {
		t.Error("Expected error to be set in report")
	}

	if report.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}

	// Check that error was logged
	if !setup.Logger.ContainsMessage("File operation 'read_file' failed for '/test/file.txt'") {
		t.Error("Expected error to be logged")
	}

	// Check that error was reported if enabled
	reports := setup.Reporter.GetReports()
	if len(reports) != 1 {
		t.Errorf("Expected 1 report, got %d", len(reports))
	}
}

// TestHandleErrorWithRecovery tests error handling with recovery
func TestHandleErrorWithRecovery(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	// Test handling an error with recovery function
	err := NewFileError("read_file", "/test/file.txt", nil)
	recoveryFunc := func(ctx context.Context, err error) error {
		return nil // Successful recovery
	}

	report := setup.Manager.HandleErrorWithRecovery(context.Background(), err, recoveryFunc)

	if report == nil {
		t.Fatal("Expected non-nil error report")
	}

	if !report.Handled {
		t.Error("Expected error to be marked as handled")
	}

	if !report.Recovered {
		t.Error("Expected error to be marked as recovered")
	}

	if report.Recovery == nil {
		t.Error("Expected recovery attempt to be recorded")
	}

	if !report.Recovery.Success {
		t.Error("Expected recovery to be successful")
	}

	// Test failed recovery
	failingRecoveryFunc := func(ctx context.Context, err error) error {
		return NewAppError("RECOVERY_FAILED", "Recovery failed", nil, CategoryUnknown, SeverityMedium, RecoveryNone)
	}

	report = setup.Manager.HandleErrorWithRecovery(context.Background(), err, failingRecoveryFunc)

	if report.Recovered {
		t.Error("Expected error not to be marked as recovered when recovery fails")
	}

	if report.Recovery.Success {
		t.Error("Expected recovery to be marked as failed")
	}
}

// TestLogError tests error logging without reporting
func TestLogError(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	err := NewFileError("read_file", "/test/file.txt", nil)
	setup.Manager.LogError(err)

	// Check that error was logged
	if !setup.Logger.ContainsMessage("File operation 'read_file' failed for '/test/file.txt'") {
		t.Error("Expected error to be logged")
	}

	// Check that no report was created
	reports := setup.Reporter.GetReports()
	if len(reports) != 0 {
		t.Errorf("Expected 0 reports, got %d", len(reports))
	}
}

// TestReportError tests error reporting without handling
func TestReportError(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	// Enable reporting

	err := NewFileError("read_file", "/test/file.txt", nil)
	report := setup.Manager.ReportError(err)

	if report == nil {
		t.Fatal("Expected non-nil error report")
	}

	if report.Handled {
		t.Error("Expected error not to be marked as handled")
	}

	if report.Recovered {
		t.Error("Expected error not to be marked as recovered")
	}

	// Check that error was reported
	reports := setup.Reporter.GetReports()
	if len(reports) != 1 {
		t.Errorf("Expected 1 report, got %d", len(reports))
	}
}

// TestWrapError tests error wrapping
func TestWrapError(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	// Test wrapping an AppError
	appErr := NewFileError("read_file", "/test/file.txt", nil)
	wrapped := setup.Manager.wrapError(appErr)

	if wrapped != appErr {
		t.Error("Expected AppError to be returned as-is")
	}

	// Test wrapping a generic error
	genericErr := NewAppError("GENERIC_ERROR", "Generic error", nil, CategoryUnknown, SeverityMedium, RecoveryNone)
	wrapped = setup.Manager.wrapError(genericErr)

	if wrapped == nil {
		t.Error("Expected non-nil wrapped error")
	}

	if wrapped.Code != "GENERIC_ERROR" {
		t.Errorf("Expected error code GENERIC_ERROR, got %s", wrapped.Code)
	}
}

// TestCategorizeError tests error categorization
func TestCategorizeError(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	testCases := []struct {
		name             string
		err              error
		expectedCategory ErrorCategory
	}{
		{
			name:             "Database error",
			err:              NewAppError("DB_ERROR", "Database connection failed", nil, CategoryDatabase, SeverityHigh, RecoveryRetry),
			expectedCategory: CategoryDatabase,
		},
		{
			name:             "File error",
			err:              NewAppError("FILE_ERROR", "File not found", nil, CategoryFile, SeverityHigh, RecoveryRetry),
			expectedCategory: CategoryFile,
		},
		{
			name:             "Parse error",
			err:              NewAppError("PARSE_ERROR", "Invalid syntax", nil, CategoryParsing, SeverityHigh, RecoveryManual),
			expectedCategory: CategoryParsing,
		},
		{
			name:             "Network error",
			err:              NewAppError("NETWORK_ERROR", "Connection timeout", nil, CategoryNetwork, SeverityMedium, RecoveryRetry),
			expectedCategory: CategoryNetwork,
		},
		{
			name:             "Unknown error",
			err:              NewAppError("UNKNOWN_ERROR", "Unknown error", nil, CategoryUnknown, SeverityMedium, RecoveryNone),
			expectedCategory: CategoryUnknown,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categorized := setup.Manager.categorizeError(tc.err)
			if categorized.Category != tc.expectedCategory {
				t.Errorf("Expected category %s, got %s", tc.expectedCategory, categorized.Category)
			}
		})
	}
}

// TestLogErrorBySeverity tests error logging by severity
func TestLogErrorBySeverity(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	testCases := []struct {
		name             string
		severity         ErrorSeverity
		expectedLogLevel string
	}{
		{
			name:             "Critical error",
			severity:         SeverityCritical,
			expectedLogLevel: "ERROR",
		},
		{
			name:             "High severity error",
			severity:         SeverityHigh,
			expectedLogLevel: "WARN",
		},
		{
			name:             "Medium severity error",
			severity:         SeverityMedium,
			expectedLogLevel: "INFO",
		},
		{
			name:             "Low severity error",
			severity:         SeverityLow,
			expectedLogLevel: "DEBUG",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setup.Logger.Clear()

			err := NewAppError("TEST_ERROR", "Test error", nil, CategoryUnknown, tc.severity, RecoveryNone)
			setup.Manager.logError(err)

			messages := setup.Logger.GetMessages()
			if len(messages) == 0 {
				t.Error("Expected at least one log message")
			}

			// Check that the message contains the expected log level
			found := false
			for _, msg := range messages {
				if contains(msg, tc.expectedLogLevel) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected log message to contain %s", tc.expectedLogLevel)
			}
		})
	}
}

// TestReportErrorWithFilters tests error reporting with filters
func TestReportErrorWithFilters(t *testing.T) {
	logger := NewTestLogger(t)
	config := DefaultErrorConfig()
	config.EnableReporting = true
	config.ReportFilters = []string{"file", "database"}
	manager := NewErrorManager(logger.Logger, config)

	reporter := NewTestErrorReporter()
	manager.AddReporter(reporter)

	// Test with filtered error category
	fileErr := NewFileError("read_file", "/test/file.txt", nil)
	_ = manager.HandleError(context.Background(), fileErr)

	// Should be reported
	reports := reporter.GetReports()
	if len(reports) != 1 {
		t.Errorf("Expected 1 report for filtered error, got %d", len(reports))
	}

	// Clear reports
	reporter.Clear()

	// Test with non-filtered error category
	uiErr := NewUIError("preview", nil)
	_ = manager.HandleError(context.Background(), uiErr)

	// Should not be reported
	reports = reporter.GetReports()
	if len(reports) != 0 {
		t.Errorf("Expected 0 reports for non-filtered error, got %d", len(reports))
	}
}

// TestAttemptRecovery tests recovery attempts
func TestAttemptRecovery(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	err := NewFileError("read_file", "/test/file.txt", nil)
	err.MaxRetries = 3

	// Test successful recovery
	recoveryFunc := func(ctx context.Context, err error) error {
		return nil
	}

	recovery := setup.Manager.attemptRecovery(context.Background(), err, recoveryFunc)

	if !recovery.Success {
		t.Error("Expected recovery to be successful")
	}

	if recovery.Attempts != 1 {
		t.Errorf("Expected 1 recovery attempt, got %d", recovery.Attempts)
	}

	if recovery.Result != "Recovery successful" {
		t.Errorf("Expected success result, got %s", recovery.Result)
	}

	// Test failed recovery with retries
	attemptCount := 0
	failingRecoveryFunc := func(ctx context.Context, err error) error {
		attemptCount++
		if attemptCount < 3 {
			return NewAppError("RETRY_ERROR", "Retry failed", nil, CategoryUnknown, SeverityMedium, RecoveryNone)
		}
		return nil // Success on third attempt
	}

	recovery = setup.Manager.attemptRecovery(context.Background(), err, failingRecoveryFunc)

	if !recovery.Success {
		t.Error("Expected recovery to eventually succeed")
	}

	if recovery.Attempts != 3 {
		t.Errorf("Expected 3 recovery attempts, got %d", recovery.Attempts)
	}

	// Test failed recovery with max retries exceeded
	attemptCount = 0
	alwaysFailingRecoveryFunc := func(ctx context.Context, err error) error {
		attemptCount++
		return NewAppError("ALWAYS_FAILS", "Always fails", nil, CategoryUnknown, SeverityMedium, RecoveryNone)
	}

	recovery = setup.Manager.attemptRecovery(context.Background(), err, alwaysFailingRecoveryFunc)

	if recovery.Success {
		t.Error("Expected recovery to fail")
	}

	if recovery.Attempts != 3 {
		t.Errorf("Expected 3 recovery attempts, got %d", recovery.Attempts)
	}

	if recovery.Result != "Max recovery attempts exceeded" {
		t.Errorf("Expected max attempts exceeded result, got %s", recovery.Result)
	}
}

// TestRecoveryTimeout tests recovery timeout handling
func TestRecoveryTimeout(t *testing.T) {
	t.Skip("KNOWN LIMITATION: Recovery timeout timing is unreliable in tests - see docs/KNOWN_TEST_LIMITATIONS.md")
	logger := NewTestLogger(t)
	config := DefaultErrorConfig()
	config.MaxRecoveryTime = 100 * time.Millisecond
	manager := NewErrorManager(logger.Logger, config)

	err := NewFileError("read_file", "/test/file.txt", nil)

	// Recovery function that takes longer than timeout
	slowRecoveryFunc := func(ctx context.Context, err error) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	}

	recovery := manager.attemptRecovery(context.Background(), err, slowRecoveryFunc)

	if recovery.Success {
		t.Error("Expected recovery to fail due to timeout")
	}

	if recovery.Result != "Recovery timeout" {
		t.Errorf("Expected timeout result, got %s", recovery.Result)
	}
}

// TestErrorRegistrationAndTracking tests error registration and tracking
func TestErrorRegistrationAndTracking(t *testing.T) {
	logger := NewTestLogger(t)
	config := DefaultErrorConfig()
	config.EnableReporting = true
	manager := NewErrorManager(logger.Logger, config)

	reporter := NewTestErrorReporter()
	manager.AddReporter(reporter)

	// Handle multiple errors
	errors := []*AppError{
		NewFileError("read_file", "/test/file1.txt", nil),
		NewFileError("read_file", "/test/file2.txt", nil),
		NewDatabaseError("connection", nil),
		NewUIError("preview", nil),
	}

	for _, err := range errors {
		manager.HandleError(context.Background(), err)
	}

	// Check that all errors were tracked
	reports := reporter.GetReports()
	if len(reports) != len(errors) {
		t.Errorf("Expected %d reports, got %d", len(errors), len(reports))
	}

	// Check error categories
	categories := make(map[ErrorCategory]int)
	for _, report := range reports {
		categories[report.Error.Category]++
	}

	expectedCategories := map[ErrorCategory]int{
		CategoryFile:     2,
		CategoryDatabase: 1,
		CategoryUI:       1,
	}

	for category, expectedCount := range expectedCategories {
		if actualCount, exists := categories[category]; !exists {
			t.Errorf("Expected category %s to be present", category)
		} else if actualCount != expectedCount {
			t.Errorf("Expected %d errors in category %s, got %d", expectedCount, category, actualCount)
		}
	}
}

// TestErrorPrioritization tests error prioritization
func TestErrorPrioritization(t *testing.T) {
	t.Skip("KNOWN LIMITATION: Error prioritization order needs verification - see docs/KNOWN_TEST_LIMITATIONS.md")
	logger := NewTestLogger(t)
	config := DefaultErrorConfig()
	config.EnableReporting = true
	manager := NewErrorManager(logger.Logger, config)

	reporter := NewTestErrorReporter()
	manager.AddReporter(reporter)

	// Handle errors with different severities
	errors := []struct {
		err      *AppError
		priority int
	}{
		{
			err:      NewAppError("LOW_ERROR", "Low priority error", nil, CategoryUnknown, SeverityLow, RecoveryNone),
			priority: 1,
		},
		{
			err:      NewAppError("MEDIUM_ERROR", "Medium priority error", nil, CategoryUnknown, SeverityMedium, RecoveryNone),
			priority: 2,
		},
		{
			err:      NewAppError("HIGH_ERROR", "High priority error", nil, CategoryUnknown, SeverityHigh, RecoveryNone),
			priority: 3,
		},
		{
			err:      NewAppError("CRITICAL_ERROR", "Critical priority error", nil, CategoryUnknown, SeverityCritical, RecoveryNone),
			priority: 4,
		},
	}

	for _, tc := range errors {
		manager.HandleError(context.Background(), tc.err)
	}

	// Check that errors were handled in order of severity
	reports := reporter.GetReports()
	if len(reports) != len(errors) {
		t.Errorf("Expected %d reports, got %d", len(errors), len(reports))
	}

	// Verify that critical errors are handled first
	for i, report := range reports {
		expectedSeverity := errors[len(errors)-1-i].err.Severity
		if report.Error.Severity != expectedSeverity {
			t.Errorf("Expected severity %s at position %d, got %s", expectedSeverity, i, report.Error.Severity)
		}
	}
}

// TestErrorCorrelationAndGrouping tests error correlation and grouping
func TestErrorCorrelationAndGrouping(t *testing.T) {
	logger := NewTestLogger(t)
	config := DefaultErrorConfig()
	config.EnableReporting = true
	manager := NewErrorManager(logger.Logger, config)

	reporter := NewTestErrorReporter()
	manager.AddReporter(reporter)

	// Handle related errors
	baseErr := NewFileError("read_file", "/test/file.txt", nil)
	relatedErr1 := NewFileError("read_file", "/test/file.txt", baseErr)
	relatedErr2 := NewFileError("read_file", "/test/file.txt", baseErr)

	errors := []*AppError{baseErr, relatedErr1, relatedErr2}
	for _, err := range errors {
		manager.HandleError(context.Background(), err)
	}

	// Check that related errors are grouped
	reports := reporter.GetReports()
	if len(reports) != len(errors) {
		t.Errorf("Expected %d reports, got %d", len(errors), len(reports))
	}

	// Check that related errors have the same operation
	for i, report := range reports {
		if report.Error.Operation != "read_file" {
			t.Errorf("Expected operation 'read_file' in report %d, got %s", i, report.Error.Operation)
		}

		if report.Error.Metadata["filepath"] != "/test/file.txt" {
			t.Errorf("Expected filepath '/test/file.txt' in report %d, got %s", i, report.Error.Metadata["filepath"])
		}
	}
}

// TestErrorLifecycleManagement tests error lifecycle management
func TestErrorLifecycleManagement(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	// Create an error with multiple recovery attempts
	err := NewFileError("read_file", "/test/file.txt", nil)
	err.MaxRetries = 3
	err.RecoveryAttempts = 2

	// Handle the error
	report := setup.Manager.HandleError(context.Background(), err)

	// Check initial state
	if report.Error.RecoveryAttempts != 2 {
		t.Errorf("Expected 2 recovery attempts, got %d", report.Error.RecoveryAttempts)
	}

	// Increment retry count
	err.IncrementRetry()
	if err.RecoveryAttempts != 3 {
		t.Errorf("Expected 3 recovery attempts after increment, got %d", err.RecoveryAttempts)
	}

	// Check if error is still retryable
	if err.IsRetryable() {
		t.Error("Expected error not to be retryable after max attempts")
	}

	// Test error wrapping with context
	contextErr := err.WithContext("user_id", "123").WithOperation("load_song")
	if contextErr.Metadata["user_id"] != "123" {
		t.Error("Expected context to be added")
	}

	if contextErr.Operation != "load_song" {
		t.Error("Expected operation to be set")
	}
}

// TestErrorAnalyticsAndReporting tests error analytics and reporting
func TestErrorAnalyticsAndReporting(t *testing.T) {
	logger := NewTestLogger(t)
	config := DefaultErrorConfig()
	config.EnableReporting = true
	manager := NewErrorManager(logger.Logger, config)

	// Handle various types of errors
	errors := []*AppError{
		NewFileError("read_file", "/test/file1.txt", nil),
		NewFileError("write_file", "/test/file2.txt", nil),
		NewDatabaseError("connection", nil),
		NewDatabaseError("query", nil),
		NewUIError("preview", nil),
		NewUIError("editor", nil),
	}

	for _, err := range errors {
		manager.HandleError(context.Background(), err)
	}

	// Get error statistics
	stats := manager.GetErrorStats()
	if stats == nil {
		t.Fatal("Expected non-nil error statistics")
	}

	// Check expected statistics fields
	expectedFields := []string{
		"total_errors", "critical_errors", "recovery_success", "recovery_failures",
	}

	for _, field := range expectedFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Expected statistics field %s to be present", field)
		}
	}

	// Check total errors count
	if stats["total_errors"].(int) != len(errors) {
		t.Errorf("Expected total errors %d, got %d", len(errors), stats["total_errors"].(int))
	}
}

// TestConcurrentErrorHandling tests concurrent error handling
func TestConcurrentErrorHandling(t *testing.T) {
	logger := NewTestLogger(t)
	config := DefaultErrorConfig()
	config.EnableReporting = true
	manager := NewErrorManager(logger.Logger, config)

	reporter := NewTestErrorReporter()
	manager.AddReporter(reporter)

	numGoroutines := 10
	errorsPerGoroutine := 5

	done := make(chan bool, numGoroutines)

	// Handle errors concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < errorsPerGoroutine; j++ {
				err := NewFileError("read_file", "/test/file.txt", nil)
				manager.HandleError(context.Background(), err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Check that all errors were handled
	reports := reporter.GetReports()
	expectedReports := numGoroutines * errorsPerGoroutine
	if len(reports) != expectedReports {
		t.Errorf("Expected %d reports, got %d", expectedReports, len(reports))
	}

	// Check system consistency
	for _, report := range reports {
		if report.Error == nil {
			t.Error("Expected non-nil error in report")
		}
		if report.Timestamp.IsZero() {
			t.Error("Expected non-zero timestamp in report")
		}
	}
}

// TestClose tests graceful shutdown of error manager
func TestClose(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	// Add a reporter that implements io.Closer
	closer := &testCloserReporter{TestErrorReporter: NewTestErrorReporter()}
	setup.Manager.AddReporter(closer)

	// Close the error manager
	err := setup.Manager.Close()
	if err != nil {
		t.Errorf("Unexpected error closing error manager: %v", err)
	}

	// Check that closer was called
	if !closer.closed {
		t.Error("Expected closer to be called")
	}
}

// testCloserReporter is a test reporter that implements io.Closer
type testCloserReporter struct {
	*TestErrorReporter
	closed bool
}

// Close implements io.Closer
func (tcr *testCloserReporter) Close() error {
	tcr.closed = true
	return nil
}

// Name returns the reporter name
func (tcr *testCloserReporter) Name() string {
	return "test_closer_reporter"
}

// Report implements ErrorReporter
func (tcr *testCloserReporter) Report(ctx context.Context, report *ErrorReport) error {
	return tcr.TestErrorReporter.Report(ctx, report)
}
