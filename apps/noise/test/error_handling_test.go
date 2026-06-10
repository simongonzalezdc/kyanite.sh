package noise

import (
	"errors"
	"testing"
	"time"

	appErrors "github.com/kyanite/noise/internal/errors"
)

// TestAppErrorCreation tests the creation of application errors
func TestAppErrorCreation(t *testing.T) {
	// Test basic error creation
	err := appErrors.NewAppError("TEST_ERROR", "Test error message", nil, appErrors.CategoryValidation, appErrors.SeverityMedium, appErrors.RecoveryManual)

	if err == nil {
		t.Fatal("Expected error to be created, got nil")
	}

	// Verify error properties
	if err.Code != "TEST_ERROR" {
		t.Errorf("Expected error code 'TEST_ERROR', got '%s'", err.Code)
	}

	if err.Message != "Test error message" {
		t.Errorf("Expected error message 'Test error message', got '%s'", err.Message)
	}

	if err.Category != appErrors.CategoryValidation {
		t.Errorf("Expected error category 'validation', got '%s'", err.Category)
	}

	if err.Severity != appErrors.SeverityMedium {
		t.Errorf("Expected error severity 'medium', got '%s'", err.Severity)
	}

	if err.Recovery != appErrors.RecoveryManual {
		t.Errorf("Expected error recovery 'manual', got '%s'", err.Recovery)
	}

	// Verify timestamp is set
	if err.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

// TestAppErrorWithCause tests error creation with underlying cause
func TestAppErrorWithCause(t *testing.T) {
	// Create underlying error
	originalErr := errors.New("original error")

	// Create app error with cause
	appErr := appErrors.NewAppError("WRAPPED_ERROR", "Wrapped error message", originalErr, appErrors.CategoryDatabase, appErrors.SeverityHigh, appErrors.RecoveryRetry)

	if appErr == nil {
		t.Fatal("Expected error to be created, got nil")
	}

	// Verify cause is preserved
	if appErr.Cause != originalErr {
		t.Error("Expected cause to be preserved")
	}

	// Verify error implements error interface
	errorMsg := appErr.Error()
	if errorMsg == "" {
		t.Error("Expected error message to be non-empty")
	}
}

// TestErrorConstructors tests specific error constructors
func TestErrorConstructors(t *testing.T) {
	// Test validation error
	validationErr := appErrors.NewValidationError("Invalid input", nil)
	if validationErr == nil {
		t.Fatal("Expected validation error to be created")
	}

	if validationErr.Category != appErrors.CategoryValidation {
		t.Error("Expected validation error to have validation category")
	}

	// Test database error
	dbErr := appErrors.NewDatabaseError("SELECT", errors.New("connection failed"))
	if dbErr == nil {
		t.Fatal("Expected database error to be created")
	}

	if dbErr.Category != appErrors.CategoryDatabase {
		t.Error("Expected database error to have database category")
	}

	// Test file error
	fileErr := appErrors.NewFileError("read", "/test/file.txt", errors.New("permission denied"))
	if fileErr == nil {
		t.Fatal("Expected file error to be created")
	}

	if fileErr.Category != appErrors.CategoryFile {
		t.Error("Expected file error to have file category")
	}
}

// TestErrorCategories tests error categories
func TestErrorCategories(t *testing.T) {
	// Test validation category
	err := appErrors.NewAppError("TEST", "Test message", nil, appErrors.CategoryValidation, appErrors.SeverityMedium, appErrors.RecoveryManual)

	if !err.HasCategory(appErrors.CategoryValidation) {
		t.Error("Expected error to have validation category")
	}

	if err.HasCategory(appErrors.CategoryDatabase) {
		t.Error("Expected error to not have database category")
	}
}

// TestErrorSeverities tests error severities
func TestErrorSeverities(t *testing.T) {
	// Test medium severity
	err := appErrors.NewAppError("TEST", "Test message", nil, appErrors.CategoryValidation, appErrors.SeverityMedium, appErrors.RecoveryManual)

	if !err.HasSeverity(appErrors.SeverityLow) {
		t.Error("Expected error to have at least low severity")
	}

	if !err.HasSeverity(appErrors.SeverityMedium) {
		t.Error("Expected error to have medium severity")
	}

	if err.HasSeverity(appErrors.SeverityHigh) {
		t.Error("Expected error to not have high severity")
	}
}

// TestErrorRecoveryStrategies tests recovery strategies
func TestErrorRecoveryStrategies(t *testing.T) {
	// Test retry recovery
	retryErr := appErrors.NewAppError("RETRY", "Retry error", nil, appErrors.CategoryNetwork, appErrors.SeverityMedium, appErrors.RecoveryRetry)

	if !retryErr.CanRecover(appErrors.RecoveryRetry) {
		t.Error("Expected retry error to support retry recovery")
	}

	// Test manual recovery
	manualErr := appErrors.NewAppError("MANUAL", "Manual error", nil, appErrors.CategoryPermission, appErrors.SeverityMedium, appErrors.RecoveryManual)

	if !manualErr.CanRecover(appErrors.RecoveryManual) {
		t.Error("Expected manual error to support manual recovery")
	}

	if manualErr.CanRecover(appErrors.RecoveryRetry) {
		t.Error("Expected manual error to not support retry recovery")
	}
}

// TestErrorChaining tests error chaining and unwrapping
func TestErrorChaining(t *testing.T) {
	// Create chain of errors
	rootErr := errors.New("root cause")
	appErr := appErrors.NewAppError("APP_ERROR", "Application error", rootErr, appErrors.CategoryDatabase, appErrors.SeverityHigh, appErrors.RecoveryRetry)

	// Test error unwrapping
	unwrapped := errors.Unwrap(appErr)
	if unwrapped != rootErr {
		t.Error("Expected Unwrap to return root error")
	}

	// Test errors.Is functionality
	if !errors.Is(appErr, rootErr) {
		t.Error("Expected errors.Is to find root cause")
	}
}

// TestErrorStringRepresentation tests error string representation
func TestErrorStringRepresentation(t *testing.T) {
	// Test error without cause
	err1 := appErrors.NewAppError("SIMPLE_ERROR", "Simple message", nil, appErrors.CategoryValidation, appErrors.SeverityMedium, appErrors.RecoveryManual)
	str1 := err1.Error()

	if str1 == "" {
		t.Error("Expected error string to be non-empty")
	}

	// Test error with cause
	cause := errors.New("underlying issue")
	err2 := appErrors.NewAppError("CAUSE_ERROR", "Error with cause", cause, appErrors.CategoryDatabase, appErrors.SeverityHigh, appErrors.RecoveryRetry)
	str2 := err2.Error()

	if str2 == "" {
		t.Error("Expected error string to be non-empty")
	}

	// Error string should be different from cause string
	causeStr := cause.Error()
	if str2 == causeStr {
		t.Error("Expected app error string to be different from cause string")
	}
}

// TestErrorTimestamp tests error timestamp functionality
func TestErrorTimestamp(t *testing.T) {
	before := time.Now()
	err := appErrors.NewAppError("TIMESTAMP_ERROR", "Timestamp test", nil, appErrors.CategoryValidation, appErrors.SeverityMedium, appErrors.RecoveryManual)
	after := time.Now()

	// Verify timestamp is reasonable
	if err.Timestamp.Before(before) || err.Timestamp.After(after) {
		t.Error("Expected timestamp to be within reasonable range")
	}

	// Test that timestamp is set
	if err.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

// TestErrorEdgeCases tests edge cases and boundary conditions
func TestErrorEdgeCases(t *testing.T) {
	// Test with empty code
	err1 := appErrors.NewAppError("", "Empty code", nil, appErrors.CategoryValidation, appErrors.SeverityMedium, appErrors.RecoveryManual)
	if err1.Code != "" {
		t.Error("Expected empty code to be preserved")
	}

	// Test with empty message
	err2 := appErrors.NewAppError("EMPTY_MSG", "", nil, appErrors.CategoryValidation, appErrors.SeverityMedium, appErrors.RecoveryManual)
	if err2.Message != "" {
		t.Error("Expected empty message to be preserved")
	}

	// Test with nil cause (should be allowed)
	err3 := appErrors.NewAppError("NIL_CAUSE", "Nil cause", nil, appErrors.CategoryValidation, appErrors.SeverityMedium, appErrors.RecoveryManual)
	if err3.Cause != nil {
		t.Error("Expected nil cause to be preserved")
	}

	// Test Unwrap with nil cause
	unwrapped := errors.Unwrap(err3)
	if unwrapped != nil {
		t.Error("Expected Unwrap of nil cause to return nil")
	}
}

// BenchmarkErrorCreation benchmarks error creation performance
func BenchmarkErrorCreation(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := appErrors.NewAppError("BENCH_ERROR", "Benchmark message", nil, appErrors.CategoryValidation, appErrors.SeverityMedium, appErrors.RecoveryManual)
		if err == nil {
			b.Error("Expected error to be created")
		}
	}
}
