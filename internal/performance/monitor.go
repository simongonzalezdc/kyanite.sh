package performance

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/app/ai"
	"github.com/Kyanite/noise/internal/collaboration"
	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/logging"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/Kyanite/noise/internal/ui"
)

// Monitor provides comprehensive performance monitoring for the noise.sh application
type Monitor struct {
	// Component monitors
	dbMonitor     *DatabaseMonitor
	aiMonitor     *AIMonitor
	uiMonitor     *UIMonitor
	themeMonitor  *ThemeMonitor
	collabMonitor *CollaborationMonitor

	// System monitoring
	systemMonitor *SystemMonitor
	memoryMonitor *MemoryMonitor

	// Aggregated metrics
	aggregatedMetrics *AggregatedMetrics

	// Configuration
	config MonitorConfig

	// Synchronization
	mutex    sync.RWMutex
	stopChan chan struct{}

	// Performance regression detection
	regressionDetector *RegressionDetector
}

// MonitorConfig defines performance monitoring configuration
type MonitorConfig struct {
	Enabled                bool            `json:"enabled"`
	CollectInterval        time.Duration   `json:"collect_interval"`
	ReportInterval         time.Duration   `json:"report_interval"`
	EnableRegression       bool            `json:"enable_regression"`
	RegressionThreshold    float64         `json:"regression_threshold"`
	EnableAlerts           bool            `json:"enable_alerts"`
	AlertThresholds        AlertThresholds `json:"alert_thresholds"`
	HistorySize            int             `json:"history_size"`
	EnableAutoOptimization bool            `json:"enable_auto_optimization"`
}

// AlertThresholds defines thresholds for performance alerts
type AlertThresholds struct {
	CPUUsage             float64       `json:"cpu_usage"`
	MemoryUsage          float64       `json:"memory_usage"`
	DatabaseQueryTime    time.Duration `json:"database_query_time"`
	AIResponseTime       time.Duration `json:"ai_response_time"`
	UIRenderTime         time.Duration `json:"ui_render_time"`
	ThemeSwitchTime      time.Duration `json:"theme_switch_time"`
	CollaborationLatency time.Duration `json:"collaboration_latency"`
}

// DatabaseMonitor monitors database performance
type DatabaseMonitor struct {
	connectionPool     *db.PerformanceOptimizedDB
	metrics            *db.PerformanceMetrics
	slowQueryThreshold time.Duration
}

// AIMonitor monitors AI service performance
type AIMonitor struct {
	service           *ai.PerformanceOptimizedAI
	metrics           ai.AIMetrics
	responseThreshold time.Duration
}

// UIMonitor monitors UI performance
type UIMonitor struct {
	service         *ui.PerformanceOptimizedUI
	metrics         ui.UIMetrics
	renderThreshold time.Duration
}

// ThemeMonitor monitors theme performance
type ThemeMonitor struct {
	service         *theme.PerformanceOptimizedManager
	metrics         theme.ThemeMetrics
	switchThreshold time.Duration
}

// CollaborationMonitor monitors collaboration performance
type CollaborationMonitor struct {
	service          *collaboration.PerformanceOptimizedCollaborationManager
	metrics          collaboration.CollaborationMetrics
	latencyThreshold time.Duration
}

// SystemMonitor monitors system-level performance
type SystemMonitor struct {
	cpuUsage    float64
	lastCPUTime time.Time
	mutex       sync.RWMutex
}

// MemoryMonitor monitors memory usage
type MemoryMonitor struct {
	allocations uint64
	gcCount     uint32
	lastGCTime  time.Time
	mutex       sync.RWMutex
}

// AggregatedMetrics combines metrics from all components
type AggregatedMetrics struct {
	Timestamp        time.Time                          `json:"timestamp"`
	Database         db.PerformanceMetrics              `json:"database"`
	AI               ai.AIMetrics                       `json:"ai"`
	UI               ui.UIMetrics                       `json:"ui"`
	Theme            theme.ThemeMetrics                 `json:"theme"`
	Collaboration    collaboration.CollaborationMetrics `json:"collaboration"`
	System           SystemMetrics                      `json:"system"`
	Memory           MemoryMetrics                      `json:"memory"`
	OverallScore     float64                            `json:"overall_score"`
	PerformanceLevel PerformanceLevel                   `json:"performance_level"`
}

// SystemMetrics represents system-level metrics
type SystemMetrics struct {
	CPUUsage   float64 `json:"cpu_usage"`
	GoRoutines int     `json:"go_routines"`
	GCCount    uint32  `json:"gc_count"`
}

// MemoryMetrics represents memory usage metrics
type MemoryMetrics struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	HeapAlloc  uint64 `json:"heap_alloc"`
	HeapSys    uint64 `json:"heap_sys"`
	HeapInuse  uint64 `json:"heap_inuse"`
}

// PerformanceLevel represents overall performance level
type PerformanceLevel string

const (
	PerformanceExcellent PerformanceLevel = "excellent"
	PerformanceGood      PerformanceLevel = "good"
	PerformanceFair      PerformanceLevel = "fair"
	PerformancePoor      PerformanceLevel = "poor"
)

// RegressionDetector detects performance regressions
type RegressionDetector struct {
	baseline   *AggregatedMetrics
	history    []*AggregatedMetrics
	maxHistory int
	threshold  float64
	mutex      sync.RWMutex
}

// PerformanceAlert represents a performance alert
type PerformanceAlert struct {
	Timestamp   time.Time `json:"timestamp"`
	Component   string    `json:"component"`
	Metric      string    `json:"metric"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
}

// NewMonitor creates a new performance monitor
func NewMonitor(config MonitorConfig) *Monitor {
	// Set defaults if not provided
	if config.CollectInterval == 0 {
		config.CollectInterval = 5 * time.Second
	}
	if config.ReportInterval == 0 {
		config.ReportInterval = 30 * time.Second
	}
	if config.RegressionThreshold == 0 {
		config.RegressionThreshold = 0.15 // 15% degradation threshold
	}
	if config.HistorySize == 0 {
		config.HistorySize = 100
	}

	monitor := &Monitor{
		config:             config,
		aggregatedMetrics:  &AggregatedMetrics{},
		stopChan:           make(chan struct{}),
		regressionDetector: NewRegressionDetector(config.RegressionThreshold, config.HistorySize),
	}

	// Initialize default alert thresholds
	if config.AlertThresholds.CPUUsage == 0 {
		config.AlertThresholds.CPUUsage = 80.0
	}
	if config.AlertThresholds.MemoryUsage == 0 {
		config.AlertThresholds.MemoryUsage = 80.0
	}
	if config.AlertThresholds.DatabaseQueryTime == 0 {
		config.AlertThresholds.DatabaseQueryTime = 100 * time.Millisecond
	}
	if config.AlertThresholds.AIResponseTime == 0 {
		config.AlertThresholds.AIResponseTime = 2 * time.Second
	}
	if config.AlertThresholds.UIRenderTime == 0 {
		config.AlertThresholds.UIRenderTime = 16 * time.Millisecond
	}
	if config.AlertThresholds.ThemeSwitchTime == 0 {
		config.AlertThresholds.ThemeSwitchTime = 50 * time.Millisecond
	}
	if config.AlertThresholds.CollaborationLatency == 0 {
		config.AlertThresholds.CollaborationLatency = 100 * time.Millisecond
	}

	return monitor
}

// Start begins performance monitoring
func (m *Monitor) Start() {
	if !m.config.Enabled {
		logging.GetDefaultLogger().Info("Performance monitoring disabled")
		return
	}

	logging.GetDefaultLogger().Info("Starting performance monitoring",
		"collect_interval", m.config.CollectInterval,
		"report_interval", m.config.ReportInterval)

	// Start collection goroutine
	go m.collectionLoop()

	// Start reporting goroutine
	go m.reportingLoop()

	// Start regression detection
	if m.config.EnableRegression {
		go m.regressionLoop()
	}
}

// Stop stops performance monitoring
func (m *Monitor) Stop() {
	logging.GetDefaultLogger().Info("Stopping performance monitoring")
	close(m.stopChan)
}

// SetDatabaseMonitor sets the database monitor
func (m *Monitor) SetDatabaseMonitor(db *db.PerformanceOptimizedDB) {
	m.dbMonitor = &DatabaseMonitor{
		connectionPool:     db,
		slowQueryThreshold: 100 * time.Millisecond,
	}
}

// SetAIMonitor sets the AI monitor
func (m *Monitor) SetAIMonitor(ai *ai.PerformanceOptimizedAI) {
	m.aiMonitor = &AIMonitor{
		service:           ai,
		responseThreshold: 2 * time.Second,
	}
}

// SetUIMonitor sets the UI monitor
func (m *Monitor) SetUIMonitor(ui *ui.PerformanceOptimizedUI) {
	m.uiMonitor = &UIMonitor{
		service:         ui,
		renderThreshold: 16 * time.Millisecond,
	}
}

// SetThemeMonitor sets the theme monitor
func (m *Monitor) SetThemeMonitor(theme *theme.PerformanceOptimizedManager) {
	m.themeMonitor = &ThemeMonitor{
		service:         theme,
		switchThreshold: 50 * time.Millisecond,
	}
}

// SetCollaborationMonitor sets the collaboration monitor
func (m *Monitor) SetCollaborationMonitor(collab *collaboration.PerformanceOptimizedCollaborationManager) {
	m.collabMonitor = &CollaborationMonitor{
		service:          collab,
		latencyThreshold: 100 * time.Millisecond,
	}
}

// collectionLoop collects metrics at regular intervals
func (m *Monitor) collectionLoop() {
	ticker := time.NewTicker(m.config.CollectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.collectMetrics()
		case <-m.stopChan:
			return
		}
	}
}

// reportingLoop generates performance reports at regular intervals
func (m *Monitor) reportingLoop() {
	ticker := time.NewTicker(m.config.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.generateReport()
		case <-m.stopChan:
			return
		}
	}
}

// regressionLoop detects performance regressions
func (m *Monitor) regressionLoop() {
	ticker := time.NewTicker(m.config.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.detectRegressions()
		case <-m.stopChan:
			return
		}
	}
}

// collectMetrics collects metrics from all components
func (m *Monitor) collectMetrics() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	metrics := &AggregatedMetrics{
		Timestamp: time.Now(),
	}

	// Collect database metrics
	if m.dbMonitor != nil && m.dbMonitor.connectionPool != nil {
		metrics.Database = m.dbMonitor.connectionPool.GetPerformanceMetrics()
	}

	// Collect AI metrics
	if m.aiMonitor != nil && m.aiMonitor.service != nil {
		metrics.AI = m.aiMonitor.service.GetMetrics()
	}

	// Collect UI metrics
	if m.uiMonitor != nil && m.uiMonitor.service != nil {
		metrics.UI = m.uiMonitor.service.GetMetrics()
	}

	// Collect theme metrics
	if m.themeMonitor != nil && m.themeMonitor.service != nil {
		metrics.Theme = m.themeMonitor.service.GetMetrics()
	}

	// Collect collaboration metrics
	if m.collabMonitor != nil && m.collabMonitor.service != nil {
		metrics.Collaboration = m.collabMonitor.service.GetMetrics()
	}

	// Collect system metrics
	metrics.System = m.collectSystemMetrics()

	// Collect memory metrics
	metrics.Memory = m.collectMemoryMetrics()

	// Calculate overall performance score
	metrics.OverallScore = m.calculatePerformanceScore(metrics)
	metrics.PerformanceLevel = m.determinePerformanceLevel(metrics.OverallScore)

	m.aggregatedMetrics = metrics

	// Check for alerts
	if m.config.EnableAlerts {
		m.checkAlerts(metrics)
	}
}

// collectSystemMetrics collects system-level metrics
func (m *Monitor) collectSystemMetrics() SystemMetrics {
	return SystemMetrics{
		CPUUsage:   m.getCPUUsage(),
		GoRoutines: runtime.NumGoroutine(),
		GCCount:    0, // Simplified - would need proper GC stats collection
	}
}

// collectMemoryMetrics collects memory usage metrics
func (m *Monitor) collectMemoryMetrics() MemoryMetrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return MemoryMetrics{
		Alloc:      memStats.Alloc,
		TotalAlloc: memStats.TotalAlloc,
		Sys:        memStats.Sys,
		HeapAlloc:  memStats.HeapAlloc,
		HeapSys:    memStats.HeapSys,
		HeapInuse:  memStats.HeapInuse,
	}
}

// getCPUUsage estimates CPU usage (simplified implementation)
func (m *Monitor) getCPUUsage() float64 {
	// This is a simplified implementation
	// In a real application, you'd use system-specific APIs
	return float64(runtime.NumGoroutine()) / 1000.0 * 100.0
}

// calculatePerformanceScore calculates an overall performance score
func (m *Monitor) calculatePerformanceScore(metrics *AggregatedMetrics) float64 {
	score := 100.0

	// Database performance impact (30% weight)
	if metrics.Database.QueryCount > 0 {
		avgQueryTime := float64(metrics.Database.TotalQueryTime.Nanoseconds()) / float64(metrics.Database.QueryCount) / 1e6
		if avgQueryTime > 100 { // 100ms threshold
			score -= (avgQueryTime - 100) * 0.3
		}
	}

	// AI performance impact (20% weight)
	if metrics.AI.TotalRequests > 0 {
		aiScore := 100.0
		if metrics.AI.AverageResponseTime > 2000*time.Millisecond { // 2s threshold
			aiScore -= float64(metrics.AI.AverageResponseTime-2000*time.Millisecond) / 1e6 * 0.2
		}
		score = score*0.8 + aiScore*0.2
	}

	// UI performance impact (20% weight)
	if metrics.UI.RenderCount > 0 {
		uiScore := 100.0
		if metrics.UI.AverageRenderTime > 16*time.Millisecond { // 16ms threshold
			uiScore -= float64(metrics.UI.AverageRenderTime-16*time.Millisecond) / 1e6 * 0.2
		}
		score = score*0.8 + uiScore*0.2
	}

	// Memory usage impact (15% weight)
	if metrics.Memory.HeapAlloc > 0 {
		memScore := 100.0
		// Simple memory pressure calculation
		memPressure := float64(metrics.Memory.HeapAlloc) / (100 * 1024 * 1024) // 100MB baseline
		if memPressure > 1.0 {
			memScore -= (memPressure - 1.0) * 10
		}
		score = score*0.85 + memScore*0.15
	}

	// System performance impact (15% weight)
	sysScore := 100.0
	if metrics.System.CPUUsage > 80 {
		sysScore -= (metrics.System.CPUUsage - 80) * 0.5
	}
	score = score*0.85 + sysScore*0.15

	// Ensure score is within bounds
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// determinePerformanceLevel determines performance level from score
func (m *Monitor) determinePerformanceLevel(score float64) PerformanceLevel {
	switch {
	case score >= 90:
		return PerformanceExcellent
	case score >= 75:
		return PerformanceGood
	case score >= 60:
		return PerformanceFair
	default:
		return PerformancePoor
	}
}

// checkAlerts checks for performance alerts
func (m *Monitor) checkAlerts(metrics *AggregatedMetrics) {
	alerts := []PerformanceAlert{}

	// Check CPU usage
	if metrics.System.CPUUsage > m.config.AlertThresholds.CPUUsage {
		alerts = append(alerts, PerformanceAlert{
			Timestamp:   time.Now(),
			Component:   "system",
			Metric:      "cpu_usage",
			Value:       metrics.System.CPUUsage,
			Threshold:   m.config.AlertThresholds.CPUUsage,
			Severity:    "warning",
			Description: fmt.Sprintf("CPU usage %.1f%% exceeds threshold %.1f%%", metrics.System.CPUUsage, m.config.AlertThresholds.CPUUsage),
		})
	}

	// Check database query time
	if metrics.Database.QueryCount > 0 {
		avgQueryTime := float64(metrics.Database.TotalQueryTime.Nanoseconds()) / float64(metrics.Database.QueryCount) / 1e6
		if avgQueryTime > float64(m.config.AlertThresholds.DatabaseQueryTime.Milliseconds()) {
			alerts = append(alerts, PerformanceAlert{
				Timestamp:   time.Now(),
				Component:   "database",
				Metric:      "query_time",
				Value:       avgQueryTime,
				Threshold:   float64(m.config.AlertThresholds.DatabaseQueryTime.Milliseconds()),
				Severity:    "warning",
				Description: fmt.Sprintf("Database query time %.1fms exceeds threshold %.1fms", avgQueryTime, float64(m.config.AlertThresholds.DatabaseQueryTime.Milliseconds())),
			})
		}
	}

	// Log alerts
	for _, alert := range alerts {
		logging.GetDefaultLogger().Warn("Performance alert",
			"component", alert.Component,
			"metric", alert.Metric,
			"value", alert.Value,
			"threshold", alert.Threshold,
			"description", alert.Description)
	}
}

// generateReport generates a performance report
func (m *Monitor) generateReport() {
	m.mutex.RLock()
	metrics := m.aggregatedMetrics
	m.mutex.RUnlock()

	report := map[string]interface{}{
		"timestamp":         metrics.Timestamp,
		"overall_score":     metrics.OverallScore,
		"performance_level": metrics.PerformanceLevel,
		"database":          metrics.Database,
		"ai":                metrics.AI,
		"ui":                metrics.UI,
		"theme":             metrics.Theme,
		"collaboration":     metrics.Collaboration,
		"system":            metrics.System,
		"memory":            metrics.Memory,
	}

	reportJSON, _ := json.MarshalIndent(report, "", "  ")
	logging.GetDefaultLogger().Info("Performance report", "report", string(reportJSON))
}

// detectRegressions detects performance regressions
func (m *Monitor) detectRegressions() {
	m.mutex.RLock()
	current := m.aggregatedMetrics
	m.mutex.RUnlock()

	regressions := m.regressionDetector.Detect(current)

	for _, regression := range regressions {
		logging.GetDefaultLogger().Error("Performance regression detected",
			"component", regression.Component,
			"metric", regression.Metric,
			"baseline", regression.BaselineValue,
			"current", regression.CurrentValue,
			"degradation", regression.Degradation)
	}
}

// GetMetrics returns the current aggregated metrics
func (m *Monitor) GetMetrics() AggregatedMetrics {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return *m.aggregatedMetrics
}

// GetPerformanceReport returns a comprehensive performance report
func (m *Monitor) GetPerformanceReport() map[string]interface{} {
	metrics := m.GetMetrics()

	report := map[string]interface{}{
		"metrics": metrics,
		"config":  m.config,
	}

	// Add regression information if enabled
	if m.config.EnableRegression {
		report["regressions"] = m.regressionDetector.GetHistory()
	}

	return report
}

// NewRegressionDetector creates a new regression detector
func NewRegressionDetector(threshold float64, maxHistory int) *RegressionDetector {
	return &RegressionDetector{
		history:    make([]*AggregatedMetrics, 0, maxHistory),
		maxHistory: maxHistory,
		threshold:  threshold,
	}
}

// Detect detects performance regressions
func (rd *RegressionDetector) Detect(current *AggregatedMetrics) []RegressionAlert {
	rd.mutex.Lock()
	defer rd.mutex.Unlock()

	alerts := []RegressionAlert{}

	// Set baseline if this is the first measurement
	if rd.baseline == nil {
		rd.baseline = current
		return alerts
	}

	// Add to history
	rd.history = append(rd.history, current)
	if len(rd.history) > rd.maxHistory {
		rd.history = rd.history[1:]
	}

	// Compare with baseline
	alerts = append(alerts, rd.compareMetrics(rd.baseline, current)...)

	return alerts
}

// RegressionAlert represents a performance regression alert
type RegressionAlert struct {
	Timestamp     time.Time `json:"timestamp"`
	Component     string    `json:"component"`
	Metric        string    `json:"metric"`
	BaselineValue float64   `json:"baseline_value"`
	CurrentValue  float64   `json:"current_value"`
	Degradation   float64   `json:"degradation"`
}

// compareMetrics compares metrics and detects regressions
func (rd *RegressionDetector) compareMetrics(baseline, current *AggregatedMetrics) []RegressionAlert {
	alerts := []RegressionAlert{}

	// Compare database performance
	if baseline.Database.QueryCount > 0 && current.Database.QueryCount > 0 {
		baselineAvg := float64(baseline.Database.TotalQueryTime.Nanoseconds()) / float64(baseline.Database.QueryCount)
		currentAvg := float64(current.Database.TotalQueryTime.Nanoseconds()) / float64(current.Database.QueryCount)

		if degradation := (currentAvg - baselineAvg) / baselineAvg; degradation > rd.threshold {
			alerts = append(alerts, RegressionAlert{
				Timestamp:     time.Now(),
				Component:     "database",
				Metric:        "average_query_time",
				BaselineValue: baselineAvg,
				CurrentValue:  currentAvg,
				Degradation:   degradation,
			})
		}
	}

	// Compare AI performance
	if baseline.AI.TotalRequests > 0 && current.AI.TotalRequests > 0 {
		baselineAvg := float64(baseline.AI.AverageResponseTime.Nanoseconds()) / 1e6
		currentAvg := float64(current.AI.AverageResponseTime.Nanoseconds()) / 1e6

		if degradation := (currentAvg - baselineAvg) / baselineAvg; degradation > rd.threshold {
			alerts = append(alerts, RegressionAlert{
				Timestamp:     time.Now(),
				Component:     "ai",
				Metric:        "average_response_time",
				BaselineValue: baselineAvg,
				CurrentValue:  currentAvg,
				Degradation:   degradation,
			})
		}
	}

	return alerts
}

// GetHistory returns the regression detection history
func (rd *RegressionDetector) GetHistory() []*AggregatedMetrics {
	rd.mutex.RLock()
	defer rd.mutex.RUnlock()

	history := make([]*AggregatedMetrics, len(rd.history))
	copy(history, rd.history)

	return history
}
