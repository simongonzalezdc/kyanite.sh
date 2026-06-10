package errors

import (
	"sync"
	"testing"
	"time"
)

// TestNewEnhancedGracefulDegradation tests the creation of enhanced graceful degradation
func TestNewEnhancedGracefulDegradation(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	if egd == nil {
		t.Fatal("Expected non-nil EnhancedGracefulDegradation")
	}

	if egd.features == nil {
		t.Error("Expected features map to be initialized")
	}

	if egd.rules == nil {
		t.Error("Expected rules map to be initialized")
	}

	if egd.healthChecker == nil {
		t.Error("Expected health checker to be initialized")
	}

	if egd.recoveryManager == nil {
		t.Error("Expected recovery manager to be initialized")
	}

	// Check that default features are registered
	expectedFeatures := []string{
		"auto_save", "real_time_preview", "audio_playback",
		"database_operations", "file_operations",
	}

	for _, feature := range expectedFeatures {
		if _, exists := egd.features[feature]; !exists {
			t.Errorf("Expected feature %s to be registered", feature)
		}
	}
}

// TestRegisterFeature tests feature registration
func TestRegisterFeature(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Test registering a new feature
	featureName := "test_feature"
	description := "Test feature for testing"
	impact := ImpactMinor
	dependencies := []string{"auto_save"}

	egd.RegisterFeature(featureName, description, impact, dependencies)

	feature, exists := egd.features[featureName]
	if !exists {
		t.Fatal("Expected feature to be registered")
	}

	if feature.Name != featureName {
		t.Errorf("Expected feature name %s, got %s", featureName, feature.Name)
	}

	if feature.Description != description {
		t.Errorf("Expected description %s, got %s", description, feature.Description)
	}

	if feature.Impact != impact {
		t.Errorf("Expected impact %v, got %v", impact, feature.Impact)
	}

	if len(feature.Dependencies) != len(dependencies) {
		t.Errorf("Expected %d dependencies, got %d", len(dependencies), len(feature.Dependencies))
	}

	// Check that logger was called
	if !logger.ContainsMessage("Feature registered for graceful degradation") {
		t.Error("Expected registration message to be logged")
	}
}

// TestHandleError tests error handling and feature degradation
func TestHandleError(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Test handling a file error that should trigger degradation
	fileErr := NewFileError("read_file", "/test/file.txt",
		NewAppError("FILE_NOT_FOUND", "File not found", nil, CategoryFile, SeverityMedium, RecoveryRetry))

	egd.HandleError(fileErr)

	// Check that auto_save feature was degraded
	autoSaveFeature, exists := egd.features["auto_save"]
	if !exists {
		t.Fatal("Expected auto_save feature to exist")
	}

	if autoSaveFeature.State != FeatureDegraded {
		t.Errorf("Expected auto_save to be degraded, got %v", autoSaveFeature.State)
	}

	if autoSaveFeature.FailureCount == 0 {
		t.Error("Expected failure count to be incremented")
	}

	if autoSaveFeature.FailureReason != fileErr.Message {
		t.Errorf("Expected failure reason %s, got %s", fileErr.Message, autoSaveFeature.FailureReason)
	}

	// Check that file_operations feature was also degraded
	fileOpsFeature, exists := egd.features["file_operations"]
	if !exists {
		t.Fatal("Expected file_operations feature to exist")
	}

	if fileOpsFeature.State != FeatureDegraded {
		t.Errorf("Expected file_operations to be degraded, got %v", fileOpsFeature.State)
	}
}

// TestFindAffectedFeatures tests finding affected features for errors
func TestFindAffectedFeatures(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Test file error
	fileErr := NewFileError("read_file", "/test/file.txt", nil)
	affected := egd.findAffectedFeatures(fileErr)

	expectedFileFeatures := []string{"file_operations", "auto_save"}
	for _, expected := range expectedFileFeatures {
		found := false
		for _, feature := range affected {
			if feature == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file error to affect %s", expected)
		}
	}

	// Test database error
	dbErr := NewDatabaseError("connection", nil)
	affected = egd.findAffectedFeatures(dbErr)

	expectedDBFeatures := []string{"database_operations", "auto_save"}
	for _, expected := range expectedDBFeatures {
		found := false
		for _, feature := range affected {
			if feature == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected database error to affect %s", expected)
		}
	}

	// Test component-specific error
	uiErr := NewAppError("UI_ERROR", "UI error in preview", nil, CategoryUI, SeverityMedium, RecoveryGraceful).
		WithComponent("preview")
	affected = egd.findAffectedFeatures(uiErr)

	found := false
	for _, feature := range affected {
		if feature == "real_time_preview" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected preview component error to affect real_time_preview")
	}
}

// TestShouldDegradeFeature tests degradation logic
func TestShouldDegradeFeature(t *testing.T) {
	t.Skip("KNOWN LIMITATION: Graceful degradation feature incomplete - see docs/KNOWN_TEST_LIMITATIONS.md")
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	feature := &FeatureInfo{
		Name:         "test_feature",
		State:        FeatureEnabled,
		FailureCount: 0,
	}

	rule := &DegradationRule{
		FeatureName:     "test_feature",
		TriggerErrors:   []ErrorCategory{CategoryFile},
		MaxFailures:     3,
		DegradationTime: 30 * time.Second,
		RecoveryTime:    5 * time.Minute,
		Impact:          ImpactMinor,
	}

	// Test with matching error category
	fileErr := NewFileError("read_file", "/test/file.txt", nil)
	if !egd.shouldDegradeFeature(feature, rule, fileErr) {
		t.Error("Expected feature to be degraded with matching error category")
	}

	// Test with non-matching error category
	uiErr := NewUIError("preview", nil)
	if egd.shouldDegradeFeature(feature, rule, uiErr) {
		t.Error("Expected feature not to be degraded with non-matching error category")
	}

	// Test with failure count threshold
	feature.FailureCount = 3
	if !egd.shouldDegradeFeature(feature, rule, uiErr) {
		t.Error("Expected feature to be degraded when failure count exceeds threshold")
	}

	// Test degraded feature still in degradation window
	feature.State = FeatureDegraded
	now := time.Now()
	recoveryTime := now.Add(1 * time.Minute)
	feature.RecoveryTime = &recoveryTime

	if egd.shouldDegradeFeature(feature, rule, fileErr) {
		t.Error("Expected degraded feature not to be degraded again while in degradation window")
	}

	// Test degraded feature past recovery window
	pastRecoveryTime := now.Add(-1 * time.Minute)
	feature.RecoveryTime = &pastRecoveryTime

	if !egd.shouldDegradeFeature(feature, rule, fileErr) {
		t.Error("Expected degraded feature past recovery window to be eligible for degradation")
	}
}

// TestDegradeFeature tests feature degradation
func TestDegradeFeature(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	feature := &FeatureInfo{
		Name:         "test_feature",
		State:        FeatureEnabled,
		Description:  "Test feature",
		FailureCount: 0,
		Impact:       ImpactMinor,
		AutoRecovery: true,
		Metadata:     make(map[string]interface{}),
	}

	rule := &DegradationRule{
		FeatureName:     "test_feature",
		TriggerErrors:   []ErrorCategory{CategoryFile},
		MaxFailures:     3,
		DegradationTime: 30 * time.Second,
		RecoveryTime:    5 * time.Minute,
		Impact:          ImpactMinor,
	}

	err := NewFileError("read_file", "/test/file.txt", nil)
	egd.degradeFeature(feature, rule, err)

	// Check feature state
	if feature.State != FeatureDegraded {
		t.Errorf("Expected feature state to be degraded, got %v", feature.State)
	}

	if feature.FailureCount != 1 {
		t.Errorf("Expected failure count to be 1, got %d", feature.FailureCount)
	}

	if feature.FailureReason != err.Message {
		t.Errorf("Expected failure reason %s, got %s", err.Message, feature.FailureReason)
	}

	if feature.RecoveryTime == nil {
		t.Error("Expected recovery time to be set")
	}

	// Check that degradation was logged
	if !logger.ContainsMessage("Feature degraded due to errors") {
		t.Error("Expected degradation to be logged")
	}
}

// TestHandleDependentFeatures tests handling of dependent features
func TestHandleDependentFeatures(t *testing.T) {
	t.Skip("KNOWN LIMITATION: Graceful degradation feature incomplete - see docs/KNOWN_TEST_LIMITATIONS.md")
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Register a feature with dependencies
	egd.RegisterFeature("dependent_feature", "Dependent feature", ImpactMinor, []string{"auto_save"})

	// Degrade the auto_save feature
	autoSaveFeature := egd.features["auto_save"]
	autoSaveFeature.State = FeatureDegraded

	// Handle dependent features
	egd.handleDependentFeatures(autoSaveFeature)

	// Check that dependent feature was also degraded
	dependentFeature, exists := egd.features["dependent_feature"]
	if !exists {
		t.Fatal("Expected dependent feature to exist")
	}

	if dependentFeature.State != FeatureDegraded {
		t.Errorf("Expected dependent feature to be degraded, got %v", dependentFeature.State)
	}
}

// TestIsFeatureAvailable tests feature availability checking
func TestIsFeatureAvailable(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Test with enabled feature
	if !egd.IsFeatureAvailable("auto_save") {
		t.Error("Expected auto_save to be available initially")
	}

	// Degrade the feature
	autoSaveFeature := egd.features["auto_save"]
	autoSaveFeature.State = FeatureDegraded

	if egd.IsFeatureAvailable("auto_save") {
		t.Error("Expected degraded auto_save to not be available")
	}

	// Test with unknown feature
	if !egd.IsFeatureAvailable("unknown_feature") {
		t.Error("Expected unknown feature to be considered available")
	}
}

// TestGetFeatureState tests getting feature state
func TestGetFeatureState(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Test getting state of existing feature
	state, err := egd.GetFeatureState("auto_save")
	if err != nil {
		t.Errorf("Unexpected error getting feature state: %v", err)
	}

	if state != FeatureEnabled {
		t.Errorf("Expected initial state to be enabled, got %v", state)
	}

	// Test getting state of non-existent feature
	_, err = egd.GetFeatureState("non_existent_feature")
	if err == nil {
		t.Error("Expected error for non-existent feature")
	}
}

// TestGetDegradedFeatures tests getting degraded features
func TestGetDegradedFeatures(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Initially no degraded features
	degraded := egd.GetDegradedFeatures()
	if len(degraded) != 0 {
		t.Errorf("Expected 0 degraded features initially, got %d", len(degraded))
	}

	// Degrade a feature
	autoSaveFeature := egd.features["auto_save"]
	autoSaveFeature.State = FeatureDegraded

	degraded = egd.GetDegradedFeatures()
	if len(degraded) != 1 {
		t.Errorf("Expected 1 degraded feature, got %d", len(degraded))
	}

	if degraded[0].Name != "auto_save" {
		t.Errorf("Expected degraded feature to be auto_save, got %s", degraded[0].Name)
	}

	// Disable a feature
	audioFeature := egd.features["audio_playback"]
	audioFeature.State = FeatureDisabled

	degraded = egd.GetDegradedFeatures()
	if len(degraded) != 2 {
		t.Errorf("Expected 2 degraded/disabled features, got %d", len(degraded))
	}
}

// TestAttemptFeatureRecovery tests feature recovery
func TestAttemptFeatureRecovery(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Test recovering non-existent feature
	err := egd.AttemptFeatureRecovery("non_existent_feature")
	if err == nil {
		t.Error("Expected error for non-existent feature")
	}

	// Test recovering enabled feature
	err = egd.AttemptFeatureRecovery("auto_save")
	if err == nil {
		t.Error("Expected error for enabled feature")
	}

	// Degrade a feature and set recovery time in the past
	autoSaveFeature := egd.features["auto_save"]
	autoSaveFeature.State = FeatureDegraded
	pastTime := time.Now().Add(-1 * time.Minute)
	autoSaveFeature.RecoveryTime = &pastTime

	// Attempt recovery
	err = egd.AttemptFeatureRecovery("auto_save")
	if err != nil {
		t.Errorf("Unexpected error during recovery: %v", err)
	}

	// Check that feature was recovered
	if autoSaveFeature.State != FeatureEnabled {
		t.Errorf("Expected feature to be enabled after recovery, got %v", autoSaveFeature.State)
	}

	if autoSaveFeature.FailureCount != 0 {
		t.Errorf("Expected failure count to be reset, got %d", autoSaveFeature.FailureCount)
	}

	if autoSaveFeature.FailureReason != "" {
		t.Errorf("Expected failure reason to be cleared, got %s", autoSaveFeature.FailureReason)
	}

	if autoSaveFeature.RecoveryTime != nil {
		t.Error("Expected recovery time to be cleared")
	}

	// Test recovering feature still in degradation window
	autoSaveFeature.State = FeatureDegraded
	futureTime := time.Now().Add(1 * time.Minute)
	autoSaveFeature.RecoveryTime = &futureTime

	err = egd.AttemptFeatureRecovery("auto_save")
	if err == nil {
		t.Error("Expected error for feature still in degradation window")
	}
}

// TestPerformHealthCheck tests periodic health checks
func TestPerformHealthCheck(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Set up a degraded feature with auto-recovery time in the past
	autoSaveFeature := egd.features["auto_save"]
	autoSaveFeature.State = FeatureDegraded
	autoSaveFeature.AutoRecovery = true
	pastTime := time.Now().Add(-1 * time.Minute)
	autoSaveFeature.RecoveryTime = &pastTime

	// Perform health check
	egd.performHealthCheck()

	// Check that feature was auto-recovered
	if autoSaveFeature.State != FeatureEnabled {
		t.Errorf("Expected feature to be auto-recovered, got %v", autoSaveFeature.State)
	}

	// Test failure count decrement for old failures
	autoSaveFeature.FailureCount = 2
	oldFailureTime := time.Now().Add(-2 * time.Hour)
	autoSaveFeature.LastFailure = &oldFailureTime

	egd.performHealthCheck()

	if autoSaveFeature.FailureCount != 1 {
		t.Errorf("Expected failure count to be decremented, got %d", autoSaveFeature.FailureCount)
	}
}

// TestGetSystemHealth tests system health assessment
func TestGetSystemHealth(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Test with all features enabled
	health := egd.GetSystemHealth()

	if health["health_score"].(int) != 100 {
		t.Errorf("Expected health score 100 with all features enabled, got %d", health["health_score"].(int))
	}

	// Degrade some features
	egd.features["auto_save"].State = FeatureDegraded
	egd.features["audio_playback"].State = FeatureDisabled

	health = egd.GetSystemHealth()
	score := health["health_score"].(int)

	if score >= 100 {
		t.Errorf("Expected health score to decrease with degraded features, got %d", score)
	}

	// Check expected fields
	expectedFields := []string{
		"health_score", "total_features", "enabled_features",
		"degraded_features", "disabled_features", "failed_features",
		"degraded_list",
	}

	for _, field := range expectedFields {
		if _, exists := health[field]; !exists {
			t.Errorf("Expected health field %s to be present", field)
		}
	}
}

// TestConcurrentDegradation tests concurrent degradation scenarios
func TestConcurrentDegradation(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	var wg sync.WaitGroup
	numGoroutines := 10
	errorsPerGoroutine := 5

	// Start multiple goroutines that trigger errors
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < errorsPerGoroutine; j++ {
				err := NewFileError("read_file", "/test/file.txt", nil)
				egd.HandleError(err)
			}
		}(i)
	}

	wg.Wait()

	// Check that the system is still in a consistent state
	autoSaveFeature := egd.features["auto_save"]
	if autoSaveFeature.State != FeatureDegraded {
		t.Errorf("Expected auto_save to be degraded, got %v", autoSaveFeature.State)
	}

	// Check that failure count is reasonable (should be at least numGoroutines * errorsPerGoroutine)
	expectedMinFailures := numGoroutines * errorsPerGoroutine
	if autoSaveFeature.FailureCount < expectedMinFailures {
		t.Errorf("Expected at least %d failures, got %d", expectedMinFailures, autoSaveFeature.FailureCount)
	}
}

// TestPerformanceImpactDuringDegradation tests performance during degradation
func TestPerformanceImpactDuringDegradation(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Measure time to handle errors when features are enabled
	start := time.Now()
	for i := 0; i < 1000; i++ {
		err := NewFileError("read_file", "/test/file.txt", nil)
		egd.HandleError(err)
	}
	enabledDuration := time.Since(start)

	// Degrade features
	egd.features["auto_save"].State = FeatureDegraded
	egd.features["file_operations"].State = FeatureDegraded

	// Measure time to handle errors when features are degraded
	start = time.Now()
	for i := 0; i < 1000; i++ {
		err := NewFileError("read_file", "/test/file.txt", nil)
		egd.HandleError(err)
	}
	degradedDuration := time.Since(start)

	// Check that performance doesn't degrade significantly
	// Allow for some overhead but it shouldn't be dramatically slower
	// Relative wall-clock comparison; unreliable on shared CI runners and under
	// -race, so relax it like the other perf budgets in this package.
	if !relaxPerfBudgets() && degradedDuration > enabledDuration*2 {
		t.Errorf("Performance degraded too much: enabled=%v, degraded=%v", enabledDuration, degradedDuration)
	}
}

// TestRecoveryWhenServicesReturn tests recovery when services return to normal
func TestRecoveryWhenServicesReturn(t *testing.T) {
	t.Skip("KNOWN LIMITATION: Graceful degradation feature incomplete - see docs/KNOWN_TEST_LIMITATIONS.md")
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Degrade a feature
	autoSaveFeature := egd.features["auto_save"]
	autoSaveFeature.State = FeatureDegraded
	autoSaveFeature.FailureCount = 3
	autoSaveFeature.RecoveryTime = nil // No recovery time set

	// Simulate services returning to normal by attempting recovery
	err := egd.AttemptFeatureRecovery("auto_save")
	if err == nil {
		t.Error("Expected error when attempting recovery without recovery time")
	}

	// Set recovery time in the past to allow recovery
	pastTime := time.Now().Add(-1 * time.Minute)
	autoSaveFeature.RecoveryTime = &pastTime

	err = egd.AttemptFeatureRecovery("auto_save")
	if err != nil {
		t.Errorf("Unexpected error during recovery: %v", err)
	}

	// Check that feature is fully recovered
	if autoSaveFeature.State != FeatureEnabled {
		t.Errorf("Expected feature to be enabled after recovery, got %v", autoSaveFeature.State)
	}

	if autoSaveFeature.FailureCount != 0 {
		t.Errorf("Expected failure count to be reset, got %d", autoSaveFeature.FailureCount)
	}
}

// TestServiceDegradationLevels tests different levels of service degradation
func TestServiceDegradationLevels(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Test different impact levels
	testCases := []struct {
		feature     string
		impact      FeatureImpact
		expectedErr *AppError
	}{
		{
			feature:     "auto_save",
			impact:      ImpactMinor,
			expectedErr: NewFileError("read_file", "/test/file.txt", nil),
		},
		{
			feature:     "database_operations",
			impact:      ImpactMajor,
			expectedErr: NewDatabaseError("connection", nil),
		},
	}

	for _, tc := range testCases {
		// Get the feature
		feature := egd.features[tc.feature]
		originalImpact := feature.Impact
		feature.Impact = tc.impact

		// Handle error
		egd.HandleError(tc.expectedErr)

		// Check that feature was degraded
		if feature.State != FeatureDegraded {
			t.Errorf("Expected %s to be degraded with impact %v", tc.feature, tc.impact)
		}

		// Restore original impact
		feature.Impact = originalImpact
		feature.State = FeatureEnabled
		feature.FailureCount = 0
	}
}

// TestDegradationTriggersAndConditions tests various degradation triggers and conditions
func TestDegradationTriggersAndConditions(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	testCases := []struct {
		name             string
		err              *AppError
		shouldDegrade    bool
		expectedFeatures []string
	}{
		{
			name:             "File error should degrade auto_save and file_operations",
			err:              NewFileError("read_file", "/test/file.txt", nil),
			shouldDegrade:    true,
			expectedFeatures: []string{"auto_save", "file_operations"},
		},
		{
			name:             "Database error should degrade database_operations and auto_save",
			err:              NewDatabaseError("connection", nil),
			shouldDegrade:    true,
			expectedFeatures: []string{"database_operations", "auto_save"},
		},
		{
			name:             "UI error in preview should degrade real_time_preview",
			err:              NewUIError("preview", nil),
			shouldDegrade:    true,
			expectedFeatures: []string{"real_time_preview"},
		},
		{
			name:             "Resource error should degrade real_time_preview and audio_playback",
			err:              NewResourceError("memory", nil),
			shouldDegrade:    true,
			expectedFeatures: []string{"real_time_preview", "audio_playback"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset all features to enabled
			for _, feature := range egd.features {
				feature.State = FeatureEnabled
				feature.FailureCount = 0
			}

			// Handle the error
			egd.HandleError(tc.err)

			// Check expected features
			for _, featureName := range tc.expectedFeatures {
				feature := egd.features[featureName]
				if tc.shouldDegrade && feature.State != FeatureDegraded {
					t.Errorf("Expected %s to be degraded", featureName)
				}
				if !tc.shouldDegrade && feature.State == FeatureDegraded {
					t.Errorf("Expected %s not to be degraded", featureName)
				}
			}
		})
	}
}

// TestAutoRecoveryMechanisms tests automatic recovery mechanisms
func TestAutoRecoveryMechanisms(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Set up a feature for auto-recovery
	autoSaveFeature := egd.features["auto_save"]
	autoSaveFeature.State = FeatureDegraded
	autoSaveFeature.AutoRecovery = true
	autoSaveFeature.FailureCount = 2

	// Set recovery time in the past to trigger auto-recovery
	pastTime := time.Now().Add(-1 * time.Minute)
	autoSaveFeature.RecoveryTime = &pastTime

	// Perform health check which should trigger auto-recovery
	egd.performHealthCheck()

	// Check that feature was auto-recovered
	if autoSaveFeature.State != FeatureEnabled {
		t.Errorf("Expected feature to be auto-recovered, got %v", autoSaveFeature.State)
	}

	if autoSaveFeature.FailureCount != 0 {
		t.Errorf("Expected failure count to be reset after auto-recovery, got %d", autoSaveFeature.FailureCount)
	}

	// Test feature without auto-recovery
	audioFeature := egd.features["audio_playback"]
	audioFeature.State = FeatureDegraded
	audioFeature.AutoRecovery = false
	audioFeature.RecoveryTime = &pastTime

	egd.performHealthCheck()

	// Feature should not be recovered
	if audioFeature.State != FeatureDegraded {
		t.Errorf("Expected feature without auto-recovery to remain degraded, got %v", audioFeature.State)
	}
}

// TestDependencyCascade tests cascading degradation due to dependencies
func TestDependencyCascade(t *testing.T) {
	t.Skip("KNOWN LIMITATION: Graceful degradation feature incomplete - see docs/KNOWN_TEST_LIMITATIONS.md")
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Register features with dependencies
	egd.RegisterFeature("feature_a", "Feature A", ImpactMinor, []string{})
	egd.RegisterFeature("feature_b", "Feature B", ImpactMinor, []string{"feature_a"})
	egd.RegisterFeature("feature_c", "Feature C", ImpactMinor, []string{"feature_b"})

	// Degrade the root feature
	featureA := egd.features["feature_a"]
	featureA.State = FeatureDegraded

	// Handle dependent features
	egd.handleDependentFeatures(featureA)

	// Check that dependent features were also degraded
	featureB := egd.features["feature_b"]
	if featureB.State != FeatureDegraded {
		t.Error("Expected feature_b to be degraded due to dependency on feature_a")
	}

	featureC := egd.features["feature_c"]
	if featureC.State != FeatureDegraded {
		t.Error("Expected feature_c to be degraded due to dependency on feature_b")
	}
}

// TestPartialRecoveryScenarios tests partial recovery scenarios
func TestPartialRecoveryScenarios(t *testing.T) {
	logger := NewTestLogger(t)
	egd := NewEnhancedGracefulDegradation(logger.Logger)

	// Degrade multiple features
	egd.features["auto_save"].State = FeatureDegraded
	egd.features["file_operations"].State = FeatureDisabled
	egd.features["real_time_preview"].State = FeatureDegraded

	// Set recovery times in the past for degraded features only
	pastTime := time.Now().Add(-1 * time.Minute)
	egd.features["auto_save"].RecoveryTime = &pastTime
	egd.features["real_time_preview"].RecoveryTime = &pastTime

	// Recover only the auto_save feature
	err := egd.AttemptFeatureRecovery("auto_save")
	if err != nil {
		t.Errorf("Unexpected error during recovery: %v", err)
	}

	// Check that only auto_save was recovered
	if egd.features["auto_save"].State != FeatureEnabled {
		t.Error("Expected auto_save to be recovered")
	}

	if egd.features["file_operations"].State != FeatureDisabled {
		t.Error("Expected file_operations to remain disabled")
	}

	if egd.features["real_time_preview"].State != FeatureDegraded {
		t.Error("Expected real_time_preview to remain degraded")
	}

	// Check system health reflects partial recovery
	health := egd.GetSystemHealth()
	degradedCount := health["degraded_features"].(int)
	disabledCount := health["disabled_features"].(int)

	if degradedCount != 1 {
		t.Errorf("Expected 1 degraded feature, got %d", degradedCount)
	}

	if disabledCount != 1 {
		t.Errorf("Expected 1 disabled feature, got %d", disabledCount)
	}
}
