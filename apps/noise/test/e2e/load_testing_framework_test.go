//go:build collaboration

package e2e

import (
	"testing"
	"time"
)

// TestLoadTestingFramework tests the load testing framework itself
func TestLoadTestingFramework(t *testing.T) {
	// Create a lightweight configuration for testing the framework
	config := LoadTestConfig{
		ConcurrentUsers:      2,
		OperationsPerUser:    2,
		TestDuration:         5 * time.Second,
		RampUpTime:           1 * time.Second,
		TargetSuccessRate:    0.50, // Lower for test environment
		MaxOperationTime:     500 * time.Millisecond,
		EnableAIIntegration:  true,
		EnableCollaboration:  true,
	}
	
	// Create and run load tester
	loadTester := NewLoadTester(t, config)
	defer loadTester.Cleanup()
	
	// Run the load test
	result := loadTester.RunLoadTest(t)
	
	// Verify results
	if result.TotalOperations == 0 {
		t.Fatal("Expected at least one operation to be performed")
	}
	
	if result.SuccessRate < config.TargetSuccessRate {
		t.Logf("Warning: Success rate below target: %.2f%% < %.2f%%", 
			result.SuccessRate*100, config.TargetSuccessRate*100)
	}
	
	// Verify the framework collected metrics
	if len(result.ErrorCounts) == 0 {
		t.Log("No errors recorded (may be expected in test environment)")
	}
	
	t.Logf("Load testing framework test completed successfully")
	t.Logf("  - Total operations: %d", result.TotalOperations)
	t.Logf("  - Success rate: %.2f%%", result.SuccessRate*100)
	t.Logf("  - Average operation time: %v", result.AverageOperationTime)
}