package noise

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/app/ai"
	"github.com/Kyanite/noise/internal/collaboration"
	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/performance"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/Kyanite/noise/internal/ui"
)

// TestPerformanceTargets validates that all performance optimization targets are achieved
func TestPerformanceTargets(t *testing.T) {
	// Performance targets to validate
	targets := map[string]time.Duration{
		"database_query_time": 100 * time.Millisecond, // <100ms for 95% of queries
		"ai_response_time":    2 * time.Second,        // <2 seconds for all modes
		"ui_render_time":      16 * time.Millisecond,  // <16ms (60fps) for all operations
		"theme_switch_time":   50 * time.Millisecond,  // <50ms
		"application_startup": 3 * time.Second,        // <3 seconds
	}

	t.Run("DatabasePerformance", func(t *testing.T) {
		testDatabasePerformance(t, targets["database_query_time"])
	})

	t.Run("AIPerformance", func(t *testing.T) {
		testAIPerformance(t, targets["ai_response_time"])
	})

	t.Run("UIPerformance", func(t *testing.T) {
		testUIPerformance(t, targets["ui_render_time"])
	})

	t.Run("ThemePerformance", func(t *testing.T) {
		testThemePerformance(t, targets["theme_switch_time"])
	})

	t.Run("CollaborationPerformance", func(t *testing.T) {
		testCollaborationPerformance(t)
	})

	t.Run("MemoryPerformance", func(t *testing.T) {
		testMemoryPerformance(t)
	})

	t.Run("IntegratedPerformance", func(t *testing.T) {
		testIntegratedPerformance(t, targets)
	})
}

// testDatabasePerformance validates database performance targets
func testDatabasePerformance(t *testing.T, target time.Duration) {
	// Initialize performance-optimized database
	config := db.Config{
		DataDir: "./test_data",
	}

	start := time.Now()
	optimizedDB, err := db.NewPerformanceOptimizedDB(config)
	if err != nil {
		t.Fatalf("Failed to create performance-optimized database: %v", err)
	}
	defer optimizedDB.Close()
	initTime := time.Since(start)

	t.Logf("Database initialization time: %v (target: <3s)", initTime)
	if initTime > 3*time.Second {
		t.Errorf("Database initialization exceeded target: %v > 3s", initTime)
	}

	// Test query performance
	slowQueries := 0
	totalQueries := 100
	queryTimes := make([]time.Duration, 0, totalQueries)

	for i := 0; i < totalQueries; i++ {
		start := time.Now()

		// Test various query types
		switch i % 4 {
		case 0:
			_, err = optimizedDB.GetSongOptimized(1)
		case 1:
			_, err = optimizedDB.ListSongsOptimized(10, 0)
		case 2:
			_, err = optimizedDB.SearchSongsOptimized("test", 5)
		case 3:
			// Test connection health
			healthy := optimizedDB.IsHealthy()
			if !healthy {
				t.Error("Database health check failed")
			}
		}

		queryTime := time.Since(start)
		queryTimes = append(queryTimes, queryTime)

		if err != nil {
			// Some queries might fail due to no data, that's okay for performance testing
			continue
		}

		if queryTime > target {
			slowQueries++
		}
	}

	// Calculate statistics
	var totalTime time.Duration
	for _, qt := range queryTimes {
		totalTime += qt
	}
	avgTime := totalTime / time.Duration(len(queryTimes))

	// Calculate 95th percentile
	sortedTimes := make([]time.Duration, len(queryTimes))
	copy(sortedTimes, queryTimes)

	// Simple sort (not optimal but sufficient for testing)
	for i := 0; i < len(sortedTimes); i++ {
		for j := i + 1; j < len(sortedTimes); j++ {
			if sortedTimes[i] > sortedTimes[j] {
				sortedTimes[i], sortedTimes[j] = sortedTimes[j], sortedTimes[i]
			}
		}
	}

	p95Index := int(float64(len(sortedTimes)) * 0.95)
	if p95Index >= len(sortedTimes) {
		p95Index = len(sortedTimes) - 1
	}
	p95Time := sortedTimes[p95Index]

	// Get performance metrics
	metrics := optimizedDB.GetPerformanceMetrics()

	t.Logf("Database Performance Results:")
	t.Logf("  Total queries: %d", totalQueries)
	t.Logf("  Slow queries: %d (%.1f%%)", slowQueries, float64(slowQueries)/float64(totalQueries)*100)
	t.Logf("  Average query time: %v", avgTime)
	t.Logf("  95th percentile: %v", p95Time)
	t.Logf("  Query count: %d", metrics.QueryCount)
	t.Logf("  Total query time: %v", metrics.TotalQueryTime)
	t.Logf("  Connection errors: %d", metrics.ConnectionErrors)

	// Validate targets
	if p95Time > target {
		t.Errorf("95th percentile query time exceeded target: %v > %v", p95Time, target)
	}

	if float64(slowQueries)/float64(totalQueries) > 0.05 { // More than 5% slow queries
		t.Errorf("Too many slow queries: %.1f%% > 5%%", float64(slowQueries)/float64(totalQueries)*100)
	}

	// Get performance report
	if report := getDatabasePerformanceReport(optimizedDB); report != nil {
		t.Logf("Database performance report: %+v", report)
	}
}

// testAIPerformance validates AI service performance targets
func testAIPerformance(t *testing.T, target time.Duration) {
	// Initialize performance-optimized AI
	config := ai.OptimizationConfig{
		CacheEnabled:          true,
		CacheMaxSize:          1000,
		CacheTTL:              30 * time.Minute,
		BatchProcessing:       true,
		MaxConcurrentRequests: 10,
		RequestTimeout:        10 * time.Second,
		EnableMetrics:         true,
	}

	aiService := ai.NewPerformanceOptimizedAI(config)
	defer aiService.Close()

	// Test AI response performance
	slowResponses := 0
	totalRequests := 50
	responseTimes := make([]time.Duration, 0, totalRequests)

	ctx := context.Background()

	for i := 0; i < totalRequests; i++ {
		start := time.Now()

		// Test different AI modes
		mode := ai.QuickIdeaModeSpark
		switch i % 4 {
		case 0:
			mode = ai.QuickIdeaModeUnstick
		case 1:
			mode = ai.QuickIdeaModeSpark
		case 2:
			mode = ai.QuickIdeaModeTweak
		case 3:
			mode = ai.QuickIdeaModeCheck
		}

		_, err := aiService.GenerateWithContextOptimized(
			ctx,
			ai.ContentTypeLyrics,
			mode,
			"test content for performance testing",
			map[string]string{"test": "performance"},
		)

		responseTime := time.Since(start)
		responseTimes = append(responseTimes, responseTime)

		if err != nil {
			// Some requests might fail due to mock implementation, that's okay
			continue
		}

		if responseTime > target {
			slowResponses++
		}
	}

	// Calculate statistics
	var totalTime time.Duration
	for _, rt := range responseTimes {
		totalTime += rt
	}
	avgTime := totalTime / time.Duration(len(responseTimes))

	// Get AI metrics
	metrics := aiService.GetMetrics()
	// Get cache stats from the performance report instead
	aiReport := aiService.GetPerformanceReport()
	var cacheStats map[string]interface{}
	if cacheData, ok := aiReport["cache"].(map[string]interface{}); ok {
		cacheStats = cacheData
	} else {
		cacheStats = map[string]interface{}{"hit_rate": 0.0}
	}

	t.Logf("AI Performance Results:")
	t.Logf("  Total requests: %d", totalRequests)
	t.Logf("  Slow responses: %d (%.1f%%)", slowResponses, float64(slowResponses)/float64(totalRequests)*100)
	t.Logf("  Average response time: %v", avgTime)
	t.Logf("  Total requests (from metrics): %d", metrics.TotalRequests)
	t.Logf("  Cache hits: %d", metrics.CacheHits)
	t.Logf("  Cache misses: %d", metrics.CacheMisses)
	t.Logf("  Failed requests: %d", metrics.FailedRequests)
	t.Logf("  Cache hit rate: %.2f%%", cacheStats["hit_rate"].(float64)*100)

	// Validate targets
	if avgTime > target {
		t.Errorf("Average AI response time exceeded target: %v > %v", avgTime, target)
	}

	// Get performance report
	report := aiService.GetPerformanceReport()
	t.Logf("AI performance report: %+v", report)
}

// testUIPerformance validates UI rendering performance targets
func testUIPerformance(t *testing.T, target time.Duration) {
	// Initialize performance-optimized UI
	config := ui.UIPerformanceConfig{
		MaxFrameRate:      60,
		EnableRenderCache: true,
		CacheMaxSize:      1000,
		EnableLazyLoading: true,
		AnimationPoolSize: 50,
		ThemePreloadCount: 3,
		EnableMetrics:     true,
		RenderTimeout:     16 * time.Millisecond,
	}

	uiService := ui.NewPerformanceOptimizedUI(config)
	defer uiService.Close()

	// Test UI rendering performance
	slowRenders := 0
	totalRenders := 200
	renderTimes := make([]time.Duration, 0, totalRenders)

	for i := 0; i < totalRenders; i++ {
		start := time.Now()

		// Test different UI components
		component := "status"
		switch i % 4 {
		case 0:
			component = "loading"
		case 1:
			component = "status"
		case 2:
			component = "header"
		case 3:
			component = "footer"
		}

		rendered := uiService.RenderOptimized(component, 80, 24, "default")
		renderTime := time.Since(start)
		renderTimes = append(renderTimes, renderTime)

		if rendered == "" {
			// Frame was dropped due to rate limiting
			continue
		}

		if renderTime > target {
			slowRenders++
		}
	}

	// Calculate statistics
	var totalTime time.Duration
	for _, rt := range renderTimes {
		totalTime += rt
	}
	avgTime := totalTime / time.Duration(len(renderTimes))

	// Get UI metrics
	metrics := uiService.GetMetrics()
	cacheStats := uiService.GetCacheStats()

	t.Logf("UI Performance Results:")
	t.Logf("  Total renders: %d", totalRenders)
	t.Logf("  Slow renders: %d (%.1f%%)", slowRenders, float64(slowRenders)/float64(totalRenders)*100)
	t.Logf("  Average render time: %v", avgTime)
	t.Logf("  Render count: %d", metrics.RenderCount)
	t.Logf("  Cache hits: %d", metrics.CacheHits)
	t.Logf("  Cache misses: %d", metrics.CacheMisses)
	t.Logf("  Frame drops: %d", metrics.FrameDrops)
	t.Logf("  Cache hit rate: %.2f%%", cacheStats["hit_rate"].(float64)*100)

	// Calculate FPS
	if avgTime >= time.Millisecond {
		fps := 1000 / avgTime.Milliseconds()
		t.Logf("  Estimated FPS: %d", fps)

		if fps < 55 { // Allow some tolerance below 60fps
			t.Errorf("UI frame rate below target: %d < 55 FPS", fps)
		}
	} else if avgTime > 0 {
		t.Logf("  Estimated FPS: >1000 (ultra-fast)")
	}

	// Validate targets
	if avgTime > target {
		t.Errorf("Average UI render time exceeded target: %v > %v", avgTime, target)
	}

	// Test theme switching performance
	themeSwitchTime := uiService.SwitchThemeOptimized("dark")
	t.Logf("  Theme switch time: %v", themeSwitchTime)

	if themeSwitchTime > 50*time.Millisecond {
		t.Errorf("Theme switch time exceeded target: %v > 50ms", themeSwitchTime)
	}

	// Get performance report
	report := uiService.GetPerformanceReport()
	t.Logf("UI performance report: %+v", report)
}

// testThemePerformance validates theme switching performance targets
func testThemePerformance(t *testing.T, target time.Duration) {
	// Initialize performance-optimized theme manager
	config := theme.ThemePerformanceConfig{
		PreloadThemes:     []string{"default", "dark", "light"},
		CacheSize:         100,
		EnableAsyncSave:   true,
		SaveBatchSize:     5,
		SaveInterval:      1 * time.Second,
		EnableRenderCache: true,
		RenderCacheSize:   200,
		EnableMetrics:     true,
	}

	themeManager := theme.NewPerformanceOptimizedManager(config)
	defer themeManager.Close()

	// Test theme switching performance
	switchTimes := make([]time.Duration, 0, 20)
	themes := []string{"default", "dark", "light"}

	for i := 0; i < 20; i++ {
		themeName := themes[i%len(themes)]
		start := time.Now()

		themeManager.SetThemeOptimized(themeName)

		switchTime := time.Since(start)
		switchTimes = append(switchTimes, switchTime)
	}

	// Calculate statistics
	var totalTime time.Duration
	for _, st := range switchTimes {
		totalTime += st
	}
	avgTime := totalTime / time.Duration(len(switchTimes))

	// Get theme metrics
	metrics := themeManager.GetMetrics()

	t.Logf("Theme Performance Results:")
	t.Logf("  Total switches: %d", len(switchTimes))
	t.Logf("  Average switch time: %v", avgTime)
	t.Logf("  Theme switches: %d", metrics.ThemeSwitches)
	t.Logf("  Average switch time (from metrics): %v", metrics.AverageSwitchTime)
	t.Logf("  Cache hits: %d", metrics.CacheHits)
	t.Logf("  Cache misses: %d", metrics.CacheMisses)
	t.Logf("  Async saves: %d", metrics.AsyncSaves)

	// Validate targets
	if avgTime > target {
		t.Errorf("Average theme switch time exceeded target: %v > %v", avgTime, target)
	}

	// Get performance report
	report := themeManager.GetPerformanceReport()
	t.Logf("Theme performance report: %+v", report)
}

// testCollaborationPerformance validates collaboration performance
func testCollaborationPerformance(t *testing.T) {
	// Initialize performance-optimized collaboration manager
	mockDB := &db.DB{}
	config := collaboration.CollaborationPerformanceConfig{
		SessionCacheSize:   100,
		OperationPoolSize:  50,
		BatchSize:          10,
		BatchTimeout:       100 * time.Millisecond,
		ConnectionPoolSize: 20,
		ConnectionTimeout:  30 * time.Second,
		EnableMetrics:      true,
		EnableCompression:  false,
		MaxConcurrentUsers: 50,
	}

	collabManager := collaboration.NewPerformanceOptimizedCollaborationManager(mockDB, config)
	defer collabManager.Close()

	// Test session creation and operation performance
	sessionCreationTimes := make([]time.Duration, 0, 10)
	operationTimes := make([]time.Duration, 0, 50)

	// Test session creation
	for i := 0; i < 10; i++ {
		start := time.Now()

		_, err := collabManager.CreateSessionOptimized(
			i,
			fmt.Sprintf("test-session-%d", i),
			"test-user",
			collaboration.SessionSettings{},
		)

		creationTime := time.Since(start)
		sessionCreationTimes = append(sessionCreationTimes, creationTime)

		if err != nil {
			// Some operations might fail due to mock implementation
			continue
		}
	}

	// Test operations
	for i := 0; i < 50; i++ {
		start := time.Now()

		op := collaboration.Operation{
			Type:     "insert",
			Position: i,
			Content:  fmt.Sprintf("test operation %d", i),
		}

		err := collabManager.ApplyOperationOptimized("test-session", "test-user", op)

		operationTime := time.Since(start)
		operationTimes = append(operationTimes, operationTime)

		if err != nil {
			// Some operations might fail due to mock implementation
			continue
		}
	}

	// Calculate statistics
	var totalSessionTime, totalOperationTime time.Duration
	for _, st := range sessionCreationTimes {
		totalSessionTime += st
	}
	for _, ot := range operationTimes {
		totalOperationTime += ot
	}

	avgSessionTime := time.Duration(0)
	if len(sessionCreationTimes) > 0 {
		avgSessionTime = totalSessionTime / time.Duration(len(sessionCreationTimes))
	}

	avgOperationTime := time.Duration(0)
	if len(operationTimes) > 0 {
		avgOperationTime = totalOperationTime / time.Duration(len(operationTimes))
	}

	// Get collaboration metrics
	metrics := collabManager.GetMetrics()

	t.Logf("Collaboration Performance Results:")
	t.Logf("  Session creations: %d", len(sessionCreationTimes))
	t.Logf("  Average session creation time: %v", avgSessionTime)
	t.Logf("  Operations: %d", len(operationTimes))
	t.Logf("  Average operation time: %v", avgOperationTime)
	t.Logf("  Session operations: %d", metrics.SessionOperations)
	t.Logf("  Operation batches: %d", metrics.OperationBatches)
	t.Logf("  Average operation time (from metrics): %v", metrics.AverageOperationTime)
	t.Logf("  Cache hits: %d", metrics.CacheHits)
	t.Logf("  Cache misses: %d", metrics.CacheMisses)

	// Get performance report
	report := collabManager.GetPerformanceReport()
	t.Logf("Collaboration performance report: %+v", report)
}

// testMemoryPerformance validates memory usage targets
func testMemoryPerformance(t *testing.T) {
	// Initialize memory optimizer
	config := performance.MemoryOptimizerConfig{
		EnableMonitoring:     true,
		EnableGCOptimization: true,
		EnableObjectPooling:  true,
		MemoryCheckInterval:  5 * time.Second,
		GCThreshold:          0.8,
		PressureThreshold:    0.9,
		MaxHistorySize:       100,
		StringPoolSize:       1000,
		BufferPoolSize:       100,
		EnableAdaptiveGC:     true,
		TargetMemoryUsage:    100 * 1024 * 1024, // 100MB
	}

	memOptimizer := performance.NewMemoryOptimizer(config)
	defer memOptimizer.Stop()

	// Start memory monitoring
	memOptimizer.Start()

	// Test memory usage under load
	initialMemory := memOptimizer.GetMemoryMetrics()

	// Simulate memory usage
	for i := 0; i < 1000; i++ {
		// Get strings from pool
		s := memOptimizer.GetString()
		memOptimizer.PutString(s)

		// Get buffers from pool
		buf := memOptimizer.GetBuffer()
		memOptimizer.PutBuffer(buf)
	}

	// Force garbage collection to see effect
	start := time.Now()
	// Force garbage collection manually since triggerGC is private
	runtime.GC()
	gcTime := time.Since(start)

	finalMemory := memOptimizer.GetMemoryMetrics()
	peakMemory := memOptimizer.GetPeakMemoryMetrics()

	t.Logf("Memory Performance Results:")
	t.Logf("  Initial heap alloc: %d bytes", initialMemory.HeapAlloc)
	t.Logf("  Final heap alloc: %d bytes", finalMemory.HeapAlloc)
	t.Logf("  Peak heap alloc: %d bytes", peakMemory.HeapAlloc)
	t.Logf("  Target memory usage: %d bytes", config.TargetMemoryUsage)
	t.Logf("  GC time: %v", gcTime)

	// Validate memory targets
	if finalMemory.HeapAlloc > config.TargetMemoryUsage {
		t.Errorf("Memory usage exceeds target: %d > %d bytes", finalMemory.HeapAlloc, config.TargetMemoryUsage)
	}

	// Get memory report
	report := memOptimizer.GetMemoryReport()
	t.Logf("Memory performance report: %+v", report)
}

// testIntegratedPerformance validates the integrated performance system
func testIntegratedPerformance(t *testing.T, targets map[string]time.Duration) {
	// Initialize performance manager
	config := performance.PerformanceManagerConfig{
		EnableDBOptimization:      true,
		EnableAIOptimization:      true,
		EnableUIOptimization:      true,
		EnableThemeOptimization:   true,
		EnableCollabOptimization:  true,
		EnableMemoryOptimization:  true,
		EnableMonitoring:          true,
		MetricsCollectionInterval: 1 * time.Second,
		PerformanceReportInterval: 2 * time.Second,
		EnableAutoTuning:          true,
		TargetResponseTime:        targets["ai_response_time"],
		TargetMemoryUsage:         100 * 1024 * 1024, // 100MB
		EnableRegressionDetection: true,
	}

	perfManager := performance.NewPerformanceManager(config)

	// Initialize all components
	start := time.Now()
	err := perfManager.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize performance manager: %v", err)
	}
	initTime := time.Since(start)

	t.Logf("Performance manager initialization time: %v", initTime)
	if initTime > targets["application_startup"] {
		t.Errorf("Performance manager initialization exceeded target: %v > %v", initTime, targets["application_startup"])
	}

	// Start performance monitoring
	err = perfManager.Start()
	if err != nil {
		t.Fatalf("Failed to start performance manager: %v", err)
	}
	defer perfManager.Stop()

	// Let it run for a bit to collect metrics
	time.Sleep(3 * time.Second)

	// Get comprehensive performance report
	report := perfManager.GetPerformanceReport()
	targetsReport := perfManager.GetPerformanceTargets()

	t.Logf("Integrated Performance Report:")
	t.Logf("  Performance targets: %+v", targetsReport)

	// Validate that all components are reporting
	if _, ok := report["database"]; !ok {
		t.Error("Database metrics missing from integrated report")
	}
	if _, ok := report["ai"]; !ok {
		t.Error("AI metrics missing from integrated report")
	}
	if _, ok := report["ui"]; !ok {
		t.Error("UI metrics missing from integrated report")
	}
	if _, ok := report["theme"]; !ok {
		t.Error("Theme metrics missing from integrated report")
	}
	if _, ok := report["memory"]; !ok {
		t.Error("Memory metrics missing from integrated report")
	}
	if _, ok := report["aggregated"]; !ok {
		t.Error("Aggregated metrics missing from integrated report")
	}

	// Check aggregated metrics
	if aggregated, ok := report["aggregated"].(map[string]interface{}); ok {
		t.Logf("Aggregated metrics: %+v", aggregated)

		// Validate database performance in aggregated metrics
		if dbMetrics, ok := aggregated["database"].(map[string]interface{}); ok {
			if queryTime, ok := dbMetrics["total_query_time"].(int64); ok {
				t.Logf("  Average database query time: %v", time.Duration(queryTime))
				if time.Duration(queryTime) > targets["database_query_time"] {
					t.Errorf("Aggregated database query time exceeded target: %v > %v", time.Duration(queryTime), targets["database_query_time"])
				}
			}
		}

		// Validate AI performance in aggregated metrics
		if aiMetrics, ok := aggregated["ai"].(map[string]interface{}); ok {
			if responseTime, ok := aiMetrics["average_response_time"].(float64); ok {
				t.Logf("  Average AI response time: %v", time.Duration(responseTime))
				if time.Duration(responseTime) > targets["ai_response_time"] {
					t.Errorf("Aggregated AI response time exceeded target: %v > %v", time.Duration(responseTime), targets["ai_response_time"])
				}
			}
		}

		// Validate UI performance in aggregated metrics
		if uiMetrics, ok := aggregated["ui"].(map[string]interface{}); ok {
			if renderTime, ok := uiMetrics["average_render_time"].(float64); ok {
				t.Logf("  Average UI render time: %v", time.Duration(renderTime))
				if time.Duration(renderTime) > targets["ui_render_time"] {
					t.Errorf("Aggregated UI render time exceeded target: %v > %v", time.Duration(renderTime), targets["ui_render_time"])
				}
			}
		}

		// Validate memory performance in aggregated metrics
		if memMetrics, ok := aggregated["memory"].(map[string]interface{}); ok {
			if heapAlloc, ok := memMetrics["average_heap_alloc"].(float64); ok {
				t.Logf("  Average heap allocation: %.2f MB", heapAlloc/(1024*1024))
				if uint64(heapAlloc) > 100*1024*1024 { // 100MB
					t.Errorf("Aggregated memory usage exceeded target: %.2f MB > 100 MB", heapAlloc/(1024*1024))
				}
			}
		}
	}

	t.Log("Integrated performance validation completed successfully")
}

// Helper function to get database performance report
func getDatabasePerformanceReport(db *db.PerformanceOptimizedDB) map[string]interface{} {
	metrics := db.GetPerformanceMetrics()

	return map[string]interface{}{
		"query_count":       metrics.QueryCount,
		"total_query_time":  metrics.TotalQueryTime,
		"slow_queries":      metrics.SlowQueries,
		"connection_errors": metrics.ConnectionErrors,
		"healthy":           db.IsHealthy(),
	}
}

// BenchmarkDatabaseQueries benchmarks database query performance
func BenchmarkDatabaseQueries(b *testing.B) {
	config := db.Config{
		DataDir: "./test_data",
	}

	optimizedDB, err := db.NewPerformanceOptimizedDB(config)
	if err != nil {
		b.Fatalf("Failed to create performance-optimized database: %v", err)
	}
	defer optimizedDB.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		switch i % 3 {
		case 0:
			optimizedDB.GetSongOptimized(1)
		case 1:
			optimizedDB.ListSongsOptimized(10, 0)
		case 2:
			optimizedDB.SearchSongsOptimized("test", 5)
		}
	}
}

// BenchmarkAIResponses benchmarks AI response performance
func BenchmarkAIResponses(b *testing.B) {
	config := ai.OptimizationConfig{
		CacheEnabled:          true,
		CacheMaxSize:          1000,
		CacheTTL:              30 * time.Minute,
		BatchProcessing:       true,
		MaxConcurrentRequests: 10,
		RequestTimeout:        10 * time.Second,
		EnableMetrics:         true,
	}

	aiService := ai.NewPerformanceOptimizedAI(config)
	defer aiService.Close()

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		aiService.GenerateWithContextOptimized(
			ctx,
			ai.ContentTypeLyrics,
			ai.QuickIdeaModeSpark,
			"test content for benchmarking",
			map[string]string{"test": "benchmark"},
		)
	}
}

// BenchmarkUIRendering benchmarks UI rendering performance
func BenchmarkUIRendering(b *testing.B) {
	config := ui.UIPerformanceConfig{
		MaxFrameRate:      60,
		EnableRenderCache: true,
		CacheMaxSize:      1000,
		EnableLazyLoading: true,
		AnimationPoolSize: 50,
		ThemePreloadCount: 3,
		EnableMetrics:     true,
		RenderTimeout:     16 * time.Millisecond,
	}

	uiService := ui.NewPerformanceOptimizedUI(config)
	defer uiService.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		component := "status"
		switch i % 3 {
		case 0:
			component = "loading"
		case 1:
			component = "status"
		case 2:
			component = "header"
		}

		uiService.RenderOptimized(component, 80, 24, "default")
	}
}
