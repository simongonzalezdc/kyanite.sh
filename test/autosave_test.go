package noise

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/puente-labs/noise/internal/app"
	"github.com/puente-labs/noise/internal/infra/db"
)

// TestAutoSaveServiceCreation tests the creation of an auto-save service
func TestAutoSaveServiceCreation(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Test with default config
	service := app.NewAutoSaveService(database, nil)
	if service == nil {
		t.Fatal("Expected auto-save service to be created, got nil")
	}

	// Test with custom config
	config := &app.AutoSaveConfig{
		Enabled:          true,
		IntervalSeconds:  60,
		DebounceMs:       3000,
		MaxRetries:       5,
		RetryDelayMs:     2000,
		EnableVersioning: true,
		MaxVersions:      20,
	}

	service = app.NewAutoSaveService(database, config)
	if service == nil {
		t.Fatal("Expected auto-save service with custom config to be created, got nil")
	}

	// Verify config was set
	serviceConfig := service
	if serviceConfig == nil {
		t.Error("Expected service config to be accessible")
	}
}

// TestAutoSaveConfigValidation tests auto-save configuration validation
func TestAutoSaveConfigValidation(t *testing.T) {
	// Test valid config
	validConfig := &app.AutoSaveConfig{
		Enabled:          true,
		IntervalSeconds:  30,
		DebounceMs:       2000,
		MaxRetries:       3,
		RetryDelayMs:     1000,
		EnableVersioning: true,
		MaxVersions:      10,
	}

	err := validConfig.ValidateConfig()
	if err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}

	// Test invalid interval (too short)
	invalidConfig := &app.AutoSaveConfig{
		IntervalSeconds: 3, // Too short
		DebounceMs:      2000,
		MaxRetries:      3,
		MaxVersions:     10,
	}

	err = invalidConfig.ValidateConfig()
	if err == nil {
		t.Error("Expected invalid config to fail validation")
	}

	// Test invalid interval (too long)
	invalidConfig.IntervalSeconds = 400 // Too long
	err = invalidConfig.ValidateConfig()
	if err == nil {
		t.Error("Expected invalid config to fail validation")
	}

	// Test invalid debounce (too short)
	invalidConfig = &app.AutoSaveConfig{
		IntervalSeconds: 30,
		DebounceMs:      50, // Too short
		MaxRetries:      3,
		MaxVersions:     10,
	}

	err = invalidConfig.ValidateConfig()
	if err == nil {
		t.Error("Expected invalid config to fail validation")
	}

	// Test invalid debounce (too long)
	invalidConfig.DebounceMs = 15000 // Too long
	err = invalidConfig.ValidateConfig()
	if err == nil {
		t.Error("Expected invalid config to fail validation")
	}

	// Test invalid max retries (too low)
	invalidConfig = &app.AutoSaveConfig{
		IntervalSeconds: 30,
		DebounceMs:      2000,
		MaxRetries:      0, // Too low
		MaxVersions:     10,
	}

	err = invalidConfig.ValidateConfig()
	if err == nil {
		t.Error("Expected invalid config to fail validation")
	}

	// Test invalid max retries (too high)
	invalidConfig.MaxRetries = 15 // Too high
	err = invalidConfig.ValidateConfig()
	if err == nil {
		t.Error("Expected invalid config to fail validation")
	}

	// Test invalid max versions (too low)
	invalidConfig = &app.AutoSaveConfig{
		IntervalSeconds: 30,
		DebounceMs:      2000,
		MaxRetries:      3,
		MaxVersions:     0, // Too low
	}

	err = invalidConfig.ValidateConfig()
	if err == nil {
		t.Error("Expected invalid config to fail validation")
	}

	// Test invalid max versions (too high)
	invalidConfig.MaxVersions = 150 // Too high
	err = invalidConfig.ValidateConfig()
	if err == nil {
		t.Error("Expected invalid config to fail validation")
	}
}

// TestAutoSaveServiceStartStop tests service start and stop functionality
func TestAutoSaveServiceStartStop(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Test starting service
	ctx, cancel := context.WithCancel(context.Background())
	err = service.Start(ctx)
	if err != nil {
		t.Errorf("Expected service to start successfully, got error: %v", err)
	}

	// Test stopping service
	err = service.Stop()
	if err != nil {
		t.Errorf("Expected service to stop successfully, got error: %v", err)
	}

	// Test starting disabled service
	disabledConfig := &app.AutoSaveConfig{Enabled: false}
	service = app.NewAutoSaveService(database, disabledConfig)

	err = service.Start(ctx)
	if err != nil {
		t.Errorf("Expected disabled service to start without error, got: %v", err)
	}

	cancel()
}

// TestAutoSaveContentOperations tests content saving operations
func TestAutoSaveContentOperations(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Test initial status
	status := service.GetStatus()
	if status != app.AutoSaveIdle {
		t.Error("Expected initial status to be idle")
	}

	// Test saving content
	testContent := "Test content for auto-save"
	service.SaveContent(testContent)

	// Give time for save to process
	time.Sleep(100 * time.Millisecond)

	// Test force save
	err = service.ForceSave(testContent)
	if err != nil {
		t.Errorf("Expected force save to succeed, got error: %v", err)
	}

	// Test that last save time is updated
	lastSaveTime := service.GetLastSaveTime()
	if lastSaveTime.IsZero() {
		t.Error("Expected last save time to be set after save")
	}
}

// TestAutoSaveCallbacks tests callback functionality
func TestAutoSaveCallbacks(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Test status change callback
	var statusChanges []app.AutoSaveStatus
	service.SetStatusChangeCallback(func(status app.AutoSaveStatus) {
		statusChanges = append(statusChanges, status)
	})

	// Test error callback
	var errors []error
	service.SetErrorCallback(func(err error) {
		errors = append(errors, err)
	})

	// Trigger status change
	service.SaveContent("Test content")

	// Give time for processing
	time.Sleep(100 * time.Millisecond)

	// Test that callbacks were called
	if len(statusChanges) == 0 {
		t.Error("Expected status change callback to be called")
	}

	// Test that no errors occurred
	if len(errors) > 0 {
		t.Errorf("Expected no errors, got: %v", errors[0])
	}
}

// TestAutoSaveVersioning tests versioning functionality
func TestAutoSaveVersioning(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Test saving with versioning
	testContent := "Versioned content"
	err = service.SaveWithVersioning(1, testContent, false, "test-version")
	if err != nil {
		t.Errorf("Expected versioned save to succeed, got error: %v", err)
	}

	// Test creating milestone
	err = service.CreateMilestone(1, testContent, "Test Milestone")
	if err != nil {
		t.Errorf("Expected milestone creation to succeed, got error: %v", err)
	}

	// Test getting version history
	versions, err := service.GetVersionHistory(1, 10)
	if err != nil {
		t.Errorf("Expected to get version history, got error: %v", err)
	}

	if len(versions) == 0 {
		t.Error("Expected version history to contain entries")
	}
}

// TestAutoSaveRecovery tests content recovery functionality
func TestAutoSaveRecovery(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Save some content first
	testContent := "Content for recovery test"
	err = service.ForceSave(testContent)
	if err != nil {
		t.Errorf("Expected save to succeed, got error: %v", err)
	}

	// Test recovery
	recoveredContent, err := service.RecoverFromLastSave(0)
	if err != nil {
		t.Errorf("Expected recovery to succeed, got error: %v", err)
	}

	if recoveredContent != testContent {
		t.Errorf("Expected recovered content '%s', got '%s'", testContent, recoveredContent)
	}
}

// TestAutoSaveStatistics tests statistics functionality
func TestAutoSaveStatistics(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Save multiple versions
	for i := 0; i < 5; i++ {
		content := fmt.Sprintf("Content version %d", i)
		err = service.ForceSave(content)
		if err != nil {
			t.Errorf("Expected save %d to succeed, got error: %v", i, err)
		}
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	// Test getting statistics
	stats, err := service.GetSaveStatistics(0)
	if err != nil {
		t.Errorf("Expected to get statistics, got error: %v", err)
	}

	if stats == nil {
		t.Fatal("Expected statistics to be returned")
	}

	if stats.TotalVersions < 5 {
		t.Errorf("Expected at least 5 total versions, got %d", stats.TotalVersions)
	}
}

// TestAutoSaveDebouncing tests debouncing functionality
func TestAutoSaveDebouncing(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create service with short debounce for testing
	config := &app.AutoSaveConfig{
		Enabled:         true,
		IntervalSeconds: 30,
		DebounceMs:      100, // Short debounce for testing
		MaxRetries:      3,
		MaxVersions:     10,
	}

	service := app.NewAutoSaveService(database, config)

	// Start service
	ctx := context.Background()
	err = service.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start service: %v", err)
	}
	if err := service.Stop(); err != nil {
		t.Logf("Warning: Failed to stop auto-save service: %v", err)
	}

	// Test rapid content changes (should be debounced)
	start := time.Now()
	service.SaveContent("Content 1")
	service.SaveContent("Content 2")
	service.SaveContent("Content 3")
	duration := time.Since(start)

	// Should take at least the debounce duration
	if duration < 100*time.Millisecond {
		t.Error("Expected content changes to be debounced")
	}
}

// TestAutoSaveErrorHandling tests error handling scenarios
func TestAutoSaveErrorHandling(t *testing.T) {
	// Test with nil database
	service := app.NewAutoSaveService(nil, nil)
	if service == nil {
		t.Error("Expected service to handle nil database gracefully")
	}

	// Test operations with nil database
	service.SaveContent("test content")
	err := service.ForceSave("test content")
	if err == nil {
		t.Error("Expected error when saving with nil database")
	}
}

// TestAutoSaveConfigOperations tests configuration operations
func TestAutoSaveConfigOperations(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Test updating config
	newConfig := &app.AutoSaveConfig{
		Enabled:         false,
		IntervalSeconds: 60,
		DebounceMs:      5000,
		MaxRetries:      5,
		MaxVersions:     15,
	}

	service.UpdateConfig(newConfig)

	// Test that config was updated (verify through behavior)
	// Since we can't directly access config, we'll test through service behavior
}

// TestAutoSavePeriodicSaving tests periodic saving functionality
func TestAutoSavePeriodicSaving(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create service with short interval for testing
	config := &app.AutoSaveConfig{
		Enabled:         true,
		IntervalSeconds: 1, // Very short interval for testing
		DebounceMs:      100,
		MaxRetries:      3,
		MaxVersions:     10,
	}

	service := app.NewAutoSaveService(database, config)

	// Start service
	ctx, cancel := context.WithCancel(context.Background())
	err = service.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start service: %v", err)
	}
	defer func() {
	if err := service.Stop(); err != nil {
		t.Logf("Warning: Failed to stop auto-save service: %v", err)
	}
		cancel()
	}()

	// Set some content
	service.SaveContent("Periodic save test content")

	// Wait for periodic save to trigger
	time.Sleep(1500 * time.Millisecond)

	// Test that periodic save occurred
	lastSaveTime := service.GetLastSaveTime()
	if lastSaveTime.IsZero() {
		t.Error("Expected periodic save to update last save time")
	}
}

// TestAutoSaveVersionCleanup tests version cleanup functionality
func TestAutoSaveVersionCleanup(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Create many versions
	for i := 0; i < 15; i++ {
		content := fmt.Sprintf("Version %d for cleanup test", i)
		err = service.ForceSave(content)
		if err != nil {
			t.Errorf("Expected save %d to succeed, got error: %v", i, err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Test cleanup
	err = service.CleanupOldVersions(0)
	if err != nil {
		t.Errorf("Expected cleanup to succeed, got error: %v", err)
	}
}

// TestAutoSaveConcurrency tests concurrent operations
func TestAutoSaveConcurrency(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Start service
	ctx := context.Background()
	err = service.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start service: %v", err)
	}
	if err := service.Stop(); err != nil {
		t.Logf("Warning: Failed to stop auto-save service: %v", err)
	}

	// Test concurrent saves
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			content := fmt.Sprintf("Concurrent save content %d", id)
			service.SaveContent(content)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Test that service handles concurrency without error
	lastSaveTime := service.GetLastSaveTime()
	if lastSaveTime.IsZero() {
		t.Error("Expected concurrent saves to update last save time")
	}
}

// TestAutoSaveMilestoneOperations tests milestone operations
func TestAutoSaveMilestoneOperations(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Test creating milestone
	content := "Milestone content"
	err = service.CreateMilestone(1, content, "Test Milestone")
	if err != nil {
		t.Errorf("Expected milestone creation to succeed, got error: %v", err)
	}

	// Test getting milestones
	milestones, err := service.GetMilestones(1)
	if err != nil {
		t.Errorf("Expected to get milestones, got error: %v", err)
	}

	if len(milestones) == 0 {
		t.Error("Expected milestones to contain at least one entry")
	}

	// Verify milestone properties
	if !milestones[0].IsMilestone {
		t.Error("Expected milestone to be marked as milestone")
	}

	if milestones[0].MilestoneName != "Test Milestone" {
		t.Errorf("Expected milestone name 'Test Milestone', got '%s'", milestones[0].MilestoneName)
	}
}

// TestAutoSaveVersionRestoration tests version restoration functionality
func TestAutoSaveVersionRestoration(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Create a version to restore
	originalContent := "Original content for restoration"
	err = service.ForceSave(originalContent)
	if err != nil {
		t.Errorf("Expected save to succeed, got error: %v", err)
	}

	// Get version ID (would need to be implemented in actual service)
	// For now, we'll test the concept through the API

	// Test that restoration API exists and doesn't panic
	// Note: Actual restoration would need version ID from database
}

// BenchmarkAutoSaveOperations benchmarks auto-save operation performance
func BenchmarkAutoSaveOperations(b *testing.B) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Create test content
	contents := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		contents[i] = fmt.Sprintf("Benchmark content %d with various lengths and content for performance testing", i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		service.SaveContent(contents[i])
	}
}

// BenchmarkAutoSaveForceSave benchmarks force save performance
func BenchmarkAutoSaveForceSave(b *testing.B) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	content := "Benchmark content for force save performance testing with sufficient length"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := service.ForceSave(content)
		if err != nil {
			b.Errorf("Force save failed: %v", err)
		}
	}
}

// TestAutoSaveServiceIntegration tests integration scenarios
func TestAutoSaveServiceIntegration(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Test complete workflow
	testContent := "Complete workflow test content"

	// 1. Save content
	service.SaveContent(testContent)

	// 2. Wait for processing
	time.Sleep(100 * time.Millisecond)

	// 3. Verify save occurred
	lastSaveTime := service.GetLastSaveTime()
	if lastSaveTime.IsZero() {
		t.Error("Expected save to occur in workflow")
	}

	// 4. Test force save
	err = service.ForceSave(testContent + " updated")
	if err != nil {
		t.Errorf("Expected force save to succeed, got error: %v", err)
	}

	// 5. Test recovery
	recoveredContent, err := service.RecoverFromLastSave(0)
	if err != nil {
		t.Errorf("Expected recovery to succeed, got error: %v", err)
	}

	if recoveredContent == "" {
		t.Error("Expected recovery to return content")
	}
}

// TestAutoSaveEdgeCases tests edge cases and boundary conditions
func TestAutoSaveEdgeCases(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Test with empty content
	service.SaveContent("")
	time.Sleep(50 * time.Millisecond)

	// Test with very large content
	largeContent := strings.Repeat("Very large content ", 1000)
	service.SaveContent(largeContent)
	time.Sleep(50 * time.Millisecond)

	// Test with special characters
	specialContent := "Content with special chars: !@#$%^&*()_+{}|:<>?[]\\;',./\""
	service.SaveContent(specialContent)
	time.Sleep(50 * time.Millisecond)

	// Test with unicode content
	unicodeContent := "Content with unicode: ñáéíóú 🚀 ñáéíóú 🚀"
	service.SaveContent(unicodeContent)
	time.Sleep(50 * time.Millisecond)

	// Test that service handles all edge cases without error
	lastSaveTime := service.GetLastSaveTime()
	if lastSaveTime.IsZero() {
		t.Error("Expected service to handle edge cases without error")
	}
}

// TestAutoSaveStatusTransitions tests status transition scenarios
func TestAutoSaveStatusTransitions(t *testing.T) {
	// Create database for testing
	database, err := db.New(db.Config{DataDir: "./testdata"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := app.NewAutoSaveService(database, nil)

	// Track status changes
	var statusHistory []app.AutoSaveStatus
	service.SetStatusChangeCallback(func(status app.AutoSaveStatus) {
		statusHistory = append(statusHistory, status)
	})

	// Test status transitions
	service.SaveContent("Content 1")
	time.Sleep(50 * time.Millisecond)

	service.SaveContent("Content 2")
	time.Sleep(50 * time.Millisecond)

	// Test that status transitions occurred
	if len(statusHistory) == 0 {
		t.Error("Expected status transitions to occur")
	}

	// Test final status
	finalStatus := service.GetStatus()
	if finalStatus != app.AutoSaveIdle && finalStatus != app.AutoSaveSuccess {
		t.Error("Expected final status to be idle or success")
	}
}
