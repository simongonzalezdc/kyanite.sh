package errors

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// BenchmarkErrorManagerHandleError benchmarks error handling performance
func BenchmarkErrorManagerHandleError(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	config := DefaultErrorConfig()
	errorManager := NewErrorManager(logger.Logger, config)

	// Create test errors of different types
	errors := []error{
		NewValidationError("validation error", nil),
		NewFileError("read_file", "/test/path", nil),
		NewParsingError("parse_json", fmt.Errorf("parse error")),
		NewAppError("TEST_ERROR", "test error", nil, CategoryUnknown, SeverityLow, RecoveryNone),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := errors[i%len(errors)]
		errorManager.HandleError(context.Background(), err)
	}
}

// BenchmarkErrorManagerHandleErrorWithRecovery benchmarks error handling with recovery
func BenchmarkErrorManagerHandleErrorWithRecovery(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	config := DefaultErrorConfig()
	errorManager := NewErrorManager(logger.Logger, config)

	// Create a test error that can be recovered
	testErr := NewFileError("read_file", "/test/path", nil)
	recoveryFunc := func(ctx context.Context, err error) error {
		return nil // Mock recovery function
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		errorManager.HandleErrorWithRecovery(context.Background(), testErr, recoveryFunc)
	}
}

// BenchmarkErrorManagerGetErrorStats benchmarks statistics retrieval
func BenchmarkErrorManagerGetErrorStats(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	config := DefaultErrorConfig()
	errorManager := NewErrorManager(logger.Logger, config)

	// Add some errors to the manager
	for i := 0; i < 100; i++ {
		testErr := NewValidationError("validation error", nil)
		errorManager.HandleError(context.Background(), testErr)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = errorManager.GetErrorStats()
	}
}

// BenchmarkEnhancedNotificationManager benchmarks notification manager operations
func BenchmarkEnhancedNotificationManager(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	notificationManager := NewEnhancedNotificationManager(*logger.Logger)

	// Create test errors
	errors := make([]*AppError, 100)
	for i := 0; i < 100; i++ {
		errors[i] = NewValidationError("validation error", nil)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := errors[i%len(errors)]
		notificationManager.ShowActionableError("Test Error", err)
	}
}

// BenchmarkGracefulDegradation benchmarks graceful degradation operations
func BenchmarkGracefulDegradation(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	gracefulDegradation := NewGracefulDegradation(logger.Logger)

	testErr := NewValidationError("test error", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		featureName := fmt.Sprintf("feature_%d", i%50)
		gracefulDegradation.DisableFeature(featureName, testErr)
		gracefulDegradation.EnableFeature(featureName)
	}
}

// BenchmarkBackupManagerCreateBackup benchmarks backup creation
func BenchmarkBackupManagerCreateBackup(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	tempDir := b.TempDir()
	backupManager, _ := NewBackupManager(tempDir, 100, logger.Logger)

	// Create a test song
	song := CreateTestSong(1, "Benchmark Test Song")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		backupManager.CreateBackup(song, "benchmark")
	}
}

// BenchmarkBackupManagerRestoreBackup benchmarks backup restoration
func BenchmarkBackupManagerRestoreBackup(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	tempDir := b.TempDir()
	backupManager, _ := NewBackupManager(tempDir, 100, logger.Logger)

	// Create a test song and backup
	song := CreateTestSong(1, "Benchmark Test Song")
	backupInfo, _ := backupManager.CreateBackup(song, "benchmark")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		backupManager.RestoreBackup(backupInfo.ID)
	}
}

// BenchmarkBackupManagerListBackups benchmarks backup listing
func BenchmarkBackupManagerListBackups(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	tempDir := b.TempDir()
	backupManager, _ := NewBackupManager(tempDir, 100, logger.Logger)

	// Create multiple backups
	for i := 0; i < 50; i++ {
		song := CreateTestSong(i, fmt.Sprintf("Test Song %d", i))
		backupManager.CreateBackup(song, "benchmark")
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		backupManager.ListBackups(0) // 0 means list all
	}
}

// BenchmarkRecoveryManager benchmarks recovery manager operations
func BenchmarkRecoveryManager(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	tempDir := b.TempDir()
	backupManager, _ := NewBackupManager(tempDir, 100, logger.Logger)
	recoveryManager := NewRecoveryManager(backupManager, logger.Logger)

	// Create test errors
	errors := []error{
		NewValidationError("validation error", nil),
		NewFileError("read_file", "/test/path", nil),
		NewParsingError("parse_json", fmt.Errorf("parse error")),
		NewAppError("TEST_ERROR", "test error", nil, CategoryUnknown, SeverityLow, RecoveryNone),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := errors[i%len(errors)]
		recoveryManager.Recover(context.Background(), err, map[string]interface{}{
			"iteration": i,
		})
	}
}

// BenchmarkErrorRecoveryUI benchmarks error recovery UI operations
func BenchmarkErrorRecoveryUI(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Add some recovery operations
	for i := 0; i < 10; i++ {
		operation := RecoveryOperation{
			ID:          fmt.Sprintf("operation_%d", i),
			Type:        "test_type",
			Description: fmt.Sprintf("Test operation %d", i),
			Status:      StatusPending,
			Progress:    0.0,
			StartTime:   time.Now(),
		}
		errorRecoveryUI.recoveryOperations = append(errorRecoveryUI.recoveryOperations, operation)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		switch i % 4 {
		case 0:
			errorRecoveryUI.GetRecoveryOperations()
		case 1:
			errorRecoveryUI.GetSystemHealth()
		case 2:
			errorRecoveryUI.ToggleRecoveryPanel()
		case 3:
			errorRecoveryUI.View()
		}
	}
}

// BenchmarkErrorHandlingPipeline benchmarks the complete error handling pipeline
func BenchmarkErrorHandlingPipeline(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	tempDir := b.TempDir()

	// Initialize all error handling components
	config := DefaultErrorConfig()
	errorManager := NewErrorManager(logger.Logger, config)
	notificationManager := NewEnhancedNotificationManager(*logger.Logger)
	backupManager, _ := NewBackupManager(tempDir, 100, logger.Logger)
	recoveryManager := NewRecoveryManager(backupManager, logger.Logger)
	gracefulDegradation := NewGracefulDegradation(logger.Logger)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Create a test song for backup operations
	song := CreateTestSong(1, "Benchmark Test Song")

	// Create test errors
	errors := []error{
		NewValidationError("validation error", nil),
		NewFileError("read_file", "/test/path", nil),
		NewParsingError("parse_json", fmt.Errorf("parse error")),
		NewAppError("TEST_ERROR", "test error", nil, CategoryUnknown, SeverityLow, RecoveryNone),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := errors[i%len(errors)]

		// Step 1: Handle error
		errorManager.HandleError(context.Background(), err)

		// Step 2: Create notification
		notificationManager.ShowActionableError("Benchmark Error", err.(*AppError))

		// Step 3: Attempt recovery
		recoveryManager.Recover(context.Background(), err, map[string]interface{}{
			"iteration": i,
		})

		// Step 4: Handle feature degradation
		if i%10 == 0 {
			featureName := fmt.Sprintf("feature_%d", i%5)
			gracefulDegradation.DisableFeature(featureName, err)
		}

		// Step 5: Create backup
		if i%20 == 0 {
			backupManager.CreateBackup(song, "benchmark")
		}

		// Step 6: Update UI
		if i%5 == 0 {
			errorRecoveryUI.GetSystemHealth()
			errorRecoveryUI.GetRecoveryOperations()
		}
	}
}

// BenchmarkConcurrentErrorHandling benchmarks concurrent error handling
func BenchmarkConcurrentErrorHandling(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	config := DefaultErrorConfig()
	errorManager := NewErrorManager(logger.Logger, config)

	// Create test errors
	errors := []error{
		NewValidationError("validation error", nil),
		NewFileError("read_file", "/test/path", nil),
		NewParsingError("parse_json", fmt.Errorf("parse error")),
		NewAppError("TEST_ERROR", "test error", nil, CategoryUnknown, SeverityLow, RecoveryNone),
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			err := errors[i%len(errors)]
			errorManager.HandleError(context.Background(), err)
			i++
		}
	})
}

// BenchmarkMemoryUsage benchmarks memory usage of error handling components
func BenchmarkMemoryUsage(b *testing.B) {
	logger := NewTestLogger(&testing.T{})
	config := DefaultErrorConfig()
	errorManager := NewErrorManager(logger.Logger, config)

	// Create a large number of errors
	errors := make([]error, 1000)
	for i := 0; i < 1000; i++ {
		errors[i] = NewValidationError(
			fmt.Sprintf("validation error %d", i),
			fmt.Errorf("underlying error %d", i),
		)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := errors[i%len(errors)]
		errorManager.HandleError(context.Background(), err)
	}
}

// BenchmarkErrorCreation benchmarks error creation performance
func BenchmarkErrorCreation(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		switch i % 4 {
		case 0:
			_ = NewValidationError("validation error", nil)
		case 1:
			_ = NewFileError("read_file", "/test/path", nil)
		case 2:
			_ = NewParsingError("parse_json", fmt.Errorf("parse error"))
		case 3:
			_ = NewAppError("TEST_ERROR", "test error", nil, CategoryUnknown, SeverityLow, RecoveryNone)
		}
	}
}

// BenchmarkErrorWrapping benchmarks error wrapping performance
func BenchmarkErrorWrapping(b *testing.B) {
	baseErr := fmt.Errorf("base error")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := NewValidationError("validation error", baseErr)
		wrappedErr := NewFileError("read_file", "/test/path", err)
		_ = NewAppError("TEST_ERROR", "test error", wrappedErr, CategoryUnknown, SeverityLow, RecoveryNone)
	}
}

// BenchmarkErrorFormatting benchmarks error formatting performance
func BenchmarkErrorFormatting(b *testing.B) {
	errors := []error{
		NewValidationError("validation error", nil),
		NewFileError("read_file", "/test/path", nil),
		NewParsingError("parse_json", fmt.Errorf("parse error")),
		NewAppError("TEST_ERROR", "test error", nil, CategoryUnknown, SeverityLow, RecoveryNone),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := errors[i%len(errors)]
		_ = err.Error()
	}
}
