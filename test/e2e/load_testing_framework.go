package e2e

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/app/ai"
	"github.com/Kyanite/noise/internal/collaboration"
)

// LoadTestConfig defines configuration for load tests
type LoadTestConfig struct {
	ConcurrentUsers      int           `json:"concurrent_users"`
	OperationsPerUser    int           `json:"operations_per_user"`
	TestDuration         time.Duration `json:"test_duration"`
	RampUpTime           time.Duration `json:"ramp_up_time"`
	TargetSuccessRate    float64       `json:"target_success_rate"`
	MaxOperationTime     time.Duration `json:"max_operation_time"`
	EnableAIIntegration  bool          `json:"enable_ai_integration"`
	EnableCollaboration  bool          `json:"enable_collaboration"`
}

// DefaultLoadTestConfig returns a default configuration for load testing
func DefaultLoadTestConfig() LoadTestConfig {
	return LoadTestConfig{
		ConcurrentUsers:      5,
		OperationsPerUser:    10,
		TestDuration:         30 * time.Second,
		RampUpTime:           5 * time.Second,
		TargetSuccessRate:    0.95,
		MaxOperationTime:     100 * time.Millisecond,
		EnableAIIntegration:  true,
		EnableCollaboration:  true,
	}
}

// LoadTestResult contains the results of a load test
type LoadTestResult struct {
	TotalOperations       int           `json:"total_operations"`
	SuccessfulOperations  int           `json:"successful_operations"`
	FailedOperations      int           `json:"failed_operations"`
	SuccessRate           float64       `json:"success_rate"`
	AverageOperationTime  time.Duration `json:"average_operation_time"`
	MinOperationTime      time.Duration `json:"min_operation_time"`
	MaxOperationTime      time.Duration `json:"max_operation_time"`
	TotalTestTime         time.Duration `json:"total_test_time"`
	OperationsPerSecond   float64       `json:"operations_per_second"`
	AIRequests            int           `json:"ai_requests"`
	CollaborationSessions int           `json:"collaboration_sessions"`
	DatabaseOperations    int           `json:"database_operations"`
	ErrorCounts           map[string]int `json:"error_counts"`
}

// LoadTester provides a framework for load testing the noise.sh application
type LoadTester struct {
	config      LoadTestConfig
	setup       *E2ETestSetup
	results     []LoadTestResult
	resultMutex sync.Mutex
}

// NewLoadTester creates a new load tester instance
func NewLoadTester(t *testing.T, config LoadTestConfig) *LoadTester {
	setup := NewE2ETestSetup(t)
	
	return &LoadTester{
		config:  config,
		setup:   setup,
		results: make([]LoadTestResult, 0),
	}
}

// RunLoadTest executes a load test with the configured parameters
func (lt *LoadTester) RunLoadTest(t *testing.T) LoadTestResult {
	t.Logf("Starting load test with %d concurrent users, %d operations per user", 
		lt.config.ConcurrentUsers, lt.config.OperationsPerUser)
	
	startTime := time.Now()
	
	// Channels for coordination
	userChan := make(chan int, lt.config.ConcurrentUsers)
	resultChan := make(chan LoadTestResult, lt.config.ConcurrentUsers)
	
	// Start users with ramp-up
	go lt.rampUpUsers(userChan)
	
	// Wait group for all users
	var wg sync.WaitGroup
	wg.Add(lt.config.ConcurrentUsers)
	
	// Launch user goroutines
	for i := 0; i < lt.config.ConcurrentUsers; i++ {
		go func(userID int) {
			defer wg.Done()
			result := lt.runUserOperations(userID)
			resultChan <- result
		}(<-userChan)
	}
	
	// Wait for all users to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Collect results
	var combinedResult LoadTestResult
	var operationTimes []time.Duration
	errorCounts := make(map[string]int)
	
	for result := range resultChan {
		combinedResult.TotalOperations += result.TotalOperations
		combinedResult.SuccessfulOperations += result.SuccessfulOperations
		combinedResult.FailedOperations += result.FailedOperations
		combinedResult.AIRequests += result.AIRequests
		combinedResult.CollaborationSessions += result.CollaborationSessions
		combinedResult.DatabaseOperations += result.DatabaseOperations
		
		// Track operation times for statistics
		if result.AverageOperationTime > 0 {
			operationTimes = append(operationTimes, result.AverageOperationTime)
		}
		
		// Aggregate error counts
		for errorType, count := range result.ErrorCounts {
			errorCounts[errorType] += count
		}
	}
	
	// Calculate final statistics
	combinedResult.TotalTestTime = time.Since(startTime)
	
	if combinedResult.TotalOperations > 0 {
		combinedResult.SuccessRate = float64(combinedResult.SuccessfulOperations) / float64(combinedResult.TotalOperations)
		combinedResult.OperationsPerSecond = float64(combinedResult.TotalOperations) / combinedResult.TotalTestTime.Seconds()
	}
	
	// Calculate operation time statistics
	if len(operationTimes) > 0 {
		minTime := operationTimes[0]
		maxTime := operationTimes[0]
		totalTime := time.Duration(0)
		
		for _, opTime := range operationTimes {
			totalTime += opTime
			if opTime < minTime {
				minTime = opTime
			}
			if opTime > maxTime {
				maxTime = opTime
			}
		}
		
		combinedResult.MinOperationTime = minTime
		combinedResult.MaxOperationTime = maxTime
		combinedResult.AverageOperationTime = totalTime / time.Duration(len(operationTimes))
	}
	
	combinedResult.ErrorCounts = errorCounts
	
	// Store result
	lt.resultMutex.Lock()
	lt.results = append(lt.results, combinedResult)
	lt.resultMutex.Unlock()
	
	// Log results
	t.Logf("Load test completed:")
	t.Logf("  - Total operations: %d", combinedResult.TotalOperations)
	t.Logf("  - Success rate: %.2f%%", combinedResult.SuccessRate*100)
	t.Logf("  - Average operation time: %v", combinedResult.AverageOperationTime)
	t.Logf("  - Operations per second: %.2f", combinedResult.OperationsPerSecond)
	t.Logf("  - AI requests: %d", combinedResult.AIRequests)
	t.Logf("  - Collaboration sessions: %d", combinedResult.CollaborationSessions)
	
	return combinedResult
}

// rampUpUsers gradually starts users to simulate realistic load
func (lt *LoadTester) rampUpUsers(userChan chan<- int) {
	if lt.config.ConcurrentUsers <= 1 {
		for i := 0; i < lt.config.ConcurrentUsers; i++ {
			userChan <- i
		}
		return
	}
	
	delay := lt.config.RampUpTime / time.Duration(lt.config.ConcurrentUsers-1)
	
	for i := 0; i < lt.config.ConcurrentUsers; i++ {
		userChan <- i
		if i < lt.config.ConcurrentUsers-1 {
			time.Sleep(delay)
		}
	}
}

// runUserOperations simulates a single user performing operations
func (lt *LoadTester) runUserOperations(userID int) LoadTestResult {
	result := LoadTestResult{
		ErrorCounts: make(map[string]int),
	}
	
	editorSvc := app.NewEditorService(lt.setup.Database, lt.setup.Database)
	
	for j := 0; j < lt.config.OperationsPerUser; j++ {
		opStart := time.Now()
		
		// Create a song with unique data
		songTitle := fmt.Sprintf("Load Test Song %d-%d", userID, j)
		artistName := fmt.Sprintf("User %d", userID)
		uniqueFilepath := fmt.Sprintf("/tmp/load_test_%d_%d.txt", userID, j)
		
		song, err := editorSvc.CreateSong(songTitle, artistName)
		if err != nil {
			result.FailedOperations++
			result.ErrorCounts["create_song"]++
			continue
		}
		
		// Set unique filepath to avoid constraint violations
		song.Filepath = uniqueFilepath
		err = lt.setup.Database.UpdateSong(song)
		if err != nil {
			// Log but don't fail - the song was created successfully
		}
		
		// Add content and save version
		content := fmt.Sprintf(`# Load Test Song %d-%d

[Verse]
This is song %d-%d
Created by user %d
With test content

[Chorus]
Load test chorus for song %d-%d
`, userID, j, userID, j, userID, userID, j)
		
		_, err = lt.setup.Database.SaveVersion(song.ID, content, j == lt.config.OperationsPerUser-1, fmt.Sprintf("Version %d", j+1))
		if err != nil {
			result.FailedOperations++
			result.ErrorCounts["save_version"]++
			continue
		}
		
		result.DatabaseOperations++
		
		// Test AI integration if enabled
		if lt.config.EnableAIIntegration {
			req := ai.QuickRequest{
				Mode:    ai.QuickIdeaModeSpark,
				Context: "load test theme",
				Options: map[string]string{"theme": "testing"},
			}
			
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, err = lt.setup.AIAgent.Generate(ctx, req)
			cancel()
			
			if err != nil {
				result.ErrorCounts["ai_request"]++
			} else {
				result.AIRequests++
			}
		}
		
		// Test collaboration if enabled
		if lt.config.EnableCollaboration && j%3 == 0 { // Test collaboration less frequently
			sessionName := fmt.Sprintf("Load Test Session %d", userID)
			session, err := lt.setup.SessionManager.CreateSession(
				song.ID, 
				sessionName, 
				fmt.Sprintf("user%d", userID),
				collaboration.SessionSettings{
					MaxParticipants:  3,
					AutoSaveInterval: 30 * time.Second,
					ConflictStrategy: "merge",
					RequireApproval:  false,
				},
			)
			
			if err != nil {
				result.ErrorCounts["create_session"]++
			} else {
				result.CollaborationSessions++
				
				// End the session
				lt.setup.SessionManager.EndSession(session.ID)
			}
		}
		
		result.SuccessfulOperations++
		opDuration := time.Since(opStart)
		result.AverageOperationTime += opDuration
		
		// Check if operation exceeded maximum time
		if opDuration > lt.config.MaxOperationTime {
			result.ErrorCounts["slow_operation"]++
		}
	}
	
	result.TotalOperations = result.SuccessfulOperations + result.FailedOperations
	
	// Calculate average operation time
	if result.SuccessfulOperations > 0 {
		result.AverageOperationTime /= time.Duration(result.SuccessfulOperations)
	}
	
	return result
}

// GetResults returns all load test results
func (lt *LoadTester) GetResults() []LoadTestResult {
	lt.resultMutex.Lock()
	defer lt.resultMutex.Unlock()
	
	// Return a copy to prevent external modification
	results := make([]LoadTestResult, len(lt.results))
	copy(results, lt.results)
	return results
}

// Cleanup cleans up resources used by the load tester
func (lt *LoadTester) Cleanup() {
	if lt.setup != nil {
		lt.setup.Cleanup()
	}
}
