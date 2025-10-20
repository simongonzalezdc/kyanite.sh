package errors

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewErrorRecoveryUI tests the creation of error recovery UI
func TestNewErrorRecoveryUI(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	if errorRecoveryUI == nil {
		t.Fatal("Expected non-nil ErrorRecoveryUI")
	}

	if errorRecoveryUI.logger != logger.Logger {
		t.Error("Expected logger to be set correctly")
	}

	if errorRecoveryUI.showRecoveryPanel {
		t.Error("Expected recovery panel to be hidden initially")
	}

	if errorRecoveryUI.selectedRecovery != 0 {
		t.Error("Expected selected recovery to be 0 initially")
	}

	if errorRecoveryUI.recoveryOperations == nil {
		t.Error("Expected recovery operations slice to be initialized")
	}

	if errorRecoveryUI.spinner.View() == "" {
		t.Error("Expected spinner to be initialized")
	}

	if errorRecoveryUI.progressBar == nil {
		t.Error("Expected progress bar to be initialized")
	}
}

// TestErrorRecoveryUIShowRecoveryPanel tests showing the recovery panel
func TestErrorRecoveryUIShowRecoveryPanel(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Show recovery panel
	errorRecoveryUI.ShowRecoveryPanel()

	if !errorRecoveryUI.showRecoveryPanel {
		t.Error("Expected recovery panel to be shown")
	}
}

// TestErrorRecoveryUIHideRecoveryPanel tests hiding the recovery panel
func TestErrorRecoveryUIHideRecoveryPanel(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Show then hide recovery panel
	errorRecoveryUI.ShowRecoveryPanel()
	errorRecoveryUI.HideRecoveryPanel()

	if errorRecoveryUI.showRecoveryPanel {
		t.Error("Expected recovery panel to be hidden")
	}
}

// TestErrorRecoveryUIToggleRecoveryPanel tests toggling the recovery panel
func TestErrorRecoveryUIToggleRecoveryPanel(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Initially hidden
	if errorRecoveryUI.showRecoveryPanel {
		t.Error("Expected recovery panel to be hidden initially")
	}

	// Toggle to show
	errorRecoveryUI.ToggleRecoveryPanel()
	if !errorRecoveryUI.showRecoveryPanel {
		t.Error("Expected recovery panel to be shown after toggle")
	}

	// Toggle to hide
	errorRecoveryUI.ToggleRecoveryPanel()
	if errorRecoveryUI.showRecoveryPanel {
		t.Error("Expected recovery panel to be hidden after second toggle")
	}
}

// TestErrorRecoveryUIGetRecoveryOperations tests getting recovery operations
func TestErrorRecoveryUIGetRecoveryOperations(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Initially no operations
	operations := errorRecoveryUI.GetRecoveryOperations()
	if len(operations) != 0 {
		t.Errorf("Expected no recovery operations initially, got %d", len(operations))
	}

	// Add a mock operation
	mockOperation := RecoveryOperation{
		ID:          "test_operation",
		Type:        "test_type",
		Description: "Test operation",
		Status:      StatusPending,
		Progress:    0.0,
		StartTime:   time.Now(),
	}
	errorRecoveryUI.recoveryOperations = append(errorRecoveryUI.recoveryOperations, mockOperation)

	// Get operations
	operations = errorRecoveryUI.GetRecoveryOperations()
	if len(operations) != 1 {
		t.Errorf("Expected 1 recovery operation, got %d", len(operations))
	}

	if operations[0].ID != mockOperation.ID {
		t.Error("Expected operation ID to match")
	}

	// Verify it returns a copy (modifying returned operations shouldn't affect original)
	operations[0].ID = "modified_id"
	if errorRecoveryUI.recoveryOperations[0].ID != mockOperation.ID {
		t.Error("Expected original operations to be unchanged")
	}
}

// TestErrorRecoveryUIGetSystemHealth tests getting system health
func TestErrorRecoveryUIGetSystemHealth(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Get system health
	health := errorRecoveryUI.GetSystemHealth()

	if health.OverallScore != 100 {
		t.Errorf("Expected overall health score 100, got %d", health.OverallScore)
	}

	if health.FeatureHealth == nil {
		t.Error("Expected feature health map to be initialized")
	}

	if health.Recommendations == nil {
		t.Error("Expected recommendations slice to be initialized")
	}

	if health.LastUpdated.IsZero() {
		t.Error("Expected last updated time to be set")
	}
}

// TestErrorRecoveryUIUpdate tests UI updates
func TestErrorRecoveryUIUpdate(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Test key press 'r' to toggle panel
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updatedUI, cmd := errorRecoveryUI.Update(msg)

	if updatedUI == nil {
		t.Error("Expected updated UI to be returned")
	}

	if cmd != nil {
		t.Error("Expected no command to be returned for 'r' key")
	}

	// Test key press 'ctrl+r' to refresh operations
	msg = tea.KeyMsg{Type: tea.KeyCtrlR}
	updatedUI, cmd = errorRecoveryUI.Update(msg)

	if updatedUI == nil {
		t.Error("Expected updated UI to be returned")
	}

	if cmd == nil {
		t.Error("Expected spinner tick command to be returned for ctrl+r")
	}
}

// TestErrorRecoveryUIView tests UI rendering
func TestErrorRecoveryUIView(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Test view when panel is hidden (status bar)
	view := errorRecoveryUI.View()
	if view == "" {
		t.Error("Expected non-empty view when panel is hidden")
	}

	// Show panel and test view when panel is shown
	errorRecoveryUI.ShowRecoveryPanel()
	view = errorRecoveryUI.View()
	if view == "" {
		t.Error("Expected non-empty view when panel is shown")
	}

	// Check that view contains expected elements
	if len(view) < 10 {
		t.Error("Expected view to contain substantial content")
	}
}

// TestErrorRecoveryUIExecuteRecoveryOperation tests executing recovery operations
func TestErrorRecoveryUIExecuteRecoveryOperation(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Try to execute non-existent operation
	err := errorRecoveryUI.ExecuteRecoveryOperation("non_existent")
	if err == nil {
		t.Error("Expected error when executing non-existent operation")
	}

	// Add a mock operation
	mockOperation := RecoveryOperation{
		ID:          "test_operation",
		Type:        "test_type",
		Description: "Test operation",
		Status:      StatusPending,
		Progress:    0.0,
		StartTime:   time.Now(),
	}
	errorRecoveryUI.recoveryOperations = append(errorRecoveryUI.recoveryOperations, mockOperation)

	// Execute operation with unknown type
	err = errorRecoveryUI.ExecuteRecoveryOperation("test_operation")
	if err == nil {
		t.Error("Expected error when executing operation with unknown type")
	}

	// Check that operation status was updated to failed
	if errorRecoveryUI.recoveryOperations[0].Status != StatusFailed {
		t.Error("Expected operation status to be updated to failed")
	}
}

// TestInitializeUIIntegration tests UI integration initialization
func TestInitializeUIIntegration(t *testing.T) {
	logger := NewTestLogger(t)

	// Initialize UI integration
	err := InitializeUIIntegration(logger.Logger)
	if err != nil {
		t.Errorf("Expected UI integration initialization to succeed: %v", err)
	}

	// Check that initialization was logged
	if !logger.ContainsMessage("Error recovery UI integration initialized") {
		t.Error("Expected initialization to be logged")
	}
}

// TestInitializeGlobalErrorRecoveryUI tests global error recovery UI initialization
func TestInitializeGlobalErrorRecoveryUI(t *testing.T) {
	logger := NewTestLogger(t)

	// Initialize global error recovery UI
	InitializeGlobalErrorRecoveryUI(logger.Logger)

	// Check that global instance was created
	globalUI := GetGlobalErrorRecoveryUI()
	if globalUI == nil {
		t.Error("Expected global error recovery UI to be created")
	}

	// Check that initialization was logged
	if !logger.ContainsMessage("Global error recovery UI initialized") {
		t.Error("Expected global initialization to be logged")
	}
}

// TestGlobalErrorRecoveryUIConvenienceFunctions tests global convenience functions
func TestGlobalErrorRecoveryUIConvenienceFunctions(t *testing.T) {
	logger := NewTestLogger(t)

	// Initialize global error recovery UI
	InitializeGlobalErrorRecoveryUI(logger.Logger)

	// Test show global recovery panel
	ShowGlobalRecoveryPanel()
	globalUI := GetGlobalErrorRecoveryUI()
	if !globalUI.showRecoveryPanel {
		t.Error("Expected global recovery panel to be shown")
	}

	// Test hide global recovery panel
	HideGlobalRecoveryPanel()
	if globalUI.showRecoveryPanel {
		t.Error("Expected global recovery panel to be hidden")
	}

	// Test toggle global recovery panel
	ToggleGlobalRecoveryPanel()
	if !globalUI.showRecoveryPanel {
		t.Error("Expected global recovery panel to be shown after toggle")
	}

	ToggleGlobalRecoveryPanel()
	if globalUI.showRecoveryPanel {
		t.Error("Expected global recovery panel to be hidden after second toggle")
	}
}

// TestErrorRecoveryUIConcurrentOperations tests concurrent operations
func TestErrorRecoveryUIConcurrentOperations(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Toggle panel concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			errorRecoveryUI.ToggleRecoveryPanel()
			done <- true
		}()
	}

	// Wait for all operations to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Final state should be consistent (panel should be shown since we started with hidden and toggled 10 times)
	if !errorRecoveryUI.showRecoveryPanel {
		t.Error("Expected recovery panel to be shown after even number of toggles")
	}

	// Get operations concurrently
	for i := 0; i < 10; i++ {
		go func() {
			operations := errorRecoveryUI.GetRecoveryOperations()
			_ = operations
			done <- true
		}()
	}

	// Wait for all operations to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Get system health concurrently
	for i := 0; i < 10; i++ {
		go func() {
			health := errorRecoveryUI.GetSystemHealth()
			_ = health
			done <- true
		}()
	}

	// Wait for all operations to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestErrorRecoveryUIPerformance tests performance of UI operations
func TestErrorRecoveryUIPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Measure toggle performance
	start := time.Now()
	for i := 0; i < 1000; i++ {
		errorRecoveryUI.ToggleRecoveryPanel()
	}
	duration := time.Since(start)

	// Toggle should be fast
	if duration > 100*time.Millisecond {
		t.Errorf("Toggle operations took too long: %v", duration)
	}

	// Measure get operations performance
	start = time.Now()
	for i := 0; i < 1000; i++ {
		operations := errorRecoveryUI.GetRecoveryOperations()
		_ = operations
	}
	duration = time.Since(start)

	// Get operations should be fast
	if duration > 100*time.Millisecond {
		t.Errorf("Get operations took too long: %v", duration)
	}

	// Measure get system health performance
	start = time.Now()
	for i := 0; i < 1000; i++ {
		health := errorRecoveryUI.GetSystemHealth()
		_ = health
	}
	duration = time.Since(start)

	// Get system health should be fast
	if duration > 100*time.Millisecond {
		t.Errorf("Get system health took too long: %v", duration)
	}

	// Measure view rendering performance
	start = time.Now()
	for i := 0; i < 1000; i++ {
		view := errorRecoveryUI.View()
		_ = view
	}
	duration = time.Since(start)

	// View rendering should be reasonably fast
	if duration > 500*time.Millisecond {
		t.Errorf("View rendering took too long: %v", duration)
	}
}

// TestRecoveryOperationSerialization tests recovery operation serialization
func TestRecoveryOperationSerialization(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Create a recovery operation with all fields set
	now := time.Now()
	endTime := now.Add(5 * time.Minute)
	testErr := NewValidationError("test error", nil)

	operation := RecoveryOperation{
		ID:          "test_id",
		Type:        "test_type",
		Description: "Test description",
		Status:      StatusCompleted,
		Progress:    75.5,
		StartTime:   now,
		EndTime:     &endTime,
		Error:       testErr,
		Metadata: map[string]interface{}{
			"string_key": "string_value",
			"int_key":    42,
			"bool_key":   true,
		},
	}

	// Add operation to UI
	errorRecoveryUI.recoveryOperations = append(errorRecoveryUI.recoveryOperations, operation)

	// Get operations and verify they match
	operations := errorRecoveryUI.GetRecoveryOperations()
	if len(operations) != 1 {
		t.Fatalf("Expected 1 operation, got %d", len(operations))
	}

	retrievedOp := operations[0]
	if retrievedOp.ID != operation.ID {
		t.Error("Expected ID to match")
	}

	if retrievedOp.Type != operation.Type {
		t.Error("Expected Type to match")
	}

	if retrievedOp.Description != operation.Description {
		t.Error("Expected Description to match")
	}

	if retrievedOp.Status != operation.Status {
		t.Error("Expected Status to match")
	}

	if retrievedOp.Progress != operation.Progress {
		t.Error("Expected Progress to match")
	}

	if retrievedOp.Error.Error() != operation.Error.Error() {
		t.Error("Expected Error to match")
	}

	if retrievedOp.Metadata["string_key"] != operation.Metadata["string_key"] {
		t.Error("Expected string metadata to match")
	}

	if retrievedOp.Metadata["int_key"] != operation.Metadata["int_key"] {
		t.Error("Expected int metadata to match")
	}

	if retrievedOp.Metadata["bool_key"] != operation.Metadata["bool_key"] {
		t.Error("Expected bool metadata to match")
	}
}

// TestSystemHealthStatusSerialization tests system health status serialization
func TestSystemHealthStatusSerialization(t *testing.T) {
	logger := NewTestLogger(t)
	errorRecoveryUI := NewErrorRecoveryUI(logger.Logger)

	// Get system health
	health := errorRecoveryUI.GetSystemHealth()

	// Modify health status
	health.OverallScore = 75
	health.Recommendations = []string{
		"Test recommendation 1",
		"Test recommendation 2",
	}

	// Add feature health
	health.FeatureHealth["test_feature"] = FeatureHealth{
		Score:     80,
		Status:    "degraded",
		LastCheck: time.Now(),
		Issues:    []string{"Test issue 1", "Test issue 2"},
	}

	// Update UI health status
	errorRecoveryUI.healthStatus = health

	// Get health status and verify it matches
	retrievedHealth := errorRecoveryUI.GetSystemHealth()

	if retrievedHealth.OverallScore != health.OverallScore {
		t.Error("Expected OverallScore to match")
	}

	if len(retrievedHealth.Recommendations) != len(health.Recommendations) {
		t.Error("Expected Recommendations length to match")
	}

	if len(retrievedHealth.FeatureHealth) != len(health.FeatureHealth) {
		t.Error("Expected FeatureHealth length to match")
	}

	featureHealth, exists := retrievedHealth.FeatureHealth["test_feature"]
	if !exists {
		t.Error("Expected test_feature to exist in feature health")
	}

	if featureHealth.Score != health.FeatureHealth["test_feature"].Score {
		t.Error("Expected feature Score to match")
	}

	if featureHealth.Status != health.FeatureHealth["test_feature"].Status {
		t.Error("Expected feature Status to match")
	}
}
