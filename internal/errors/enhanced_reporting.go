package errors

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/logging"
)

// ErrorReportEnhancer enhances error reports with additional context
type ErrorReportEnhancer struct {
	systemInfo     *SystemInfo
	userContext    *UserContext
	sessionContext *SessionContext
	logger         *logging.Logger
	mu             sync.RWMutex
}

// SystemInfo contains system-level information
type SystemInfo struct {
	OS          string            `json:"os"`
	Arch        string            `json:"arch"`
	CPUCount    int               `json:"cpu_count"`
	MemoryTotal uint64            `json:"memory_total"`
	MemoryFree  uint64            `json:"memory_free"`
	GoVersion   string            `json:"go_version"`
	AppVersion  string            `json:"app_version"`
	BuildTime   string            `json:"build_time"`
	Environment map[string]string `json:"environment"`
	Metadata    map[string]string `json:"metadata"`
}

// UserContext contains user-related information
type UserContext struct {
	UserID      string            `json:"user_id,omitempty"`
	SessionID   string            `json:"session_id,omitempty"`
	Preferences map[string]string `json:"preferences,omitempty"`
	Permissions []string          `json:"permissions,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// SessionContext contains session-related information
type SessionContext struct {
	StartTime    time.Time         `json:"start_time"`
	LastActivity time.Time         `json:"last_activity"`
	Duration     time.Duration     `json:"duration"`
	Actions      []string          `json:"actions"`
	ErrorCount   int               `json:"error_count"`
	FeatureUsage map[string]int    `json:"feature_usage"`
	Metadata     map[string]string `json:"metadata"`
}

// NewErrorReportEnhancer creates a new error report enhancer
func NewErrorReportEnhancer(logger *logging.Logger) *ErrorReportEnhancer {
	enhancer := &ErrorReportEnhancer{
		logger: logger,
	}

	// Initialize system info
	enhancer.systemInfo = enhancer.collectSystemInfo()

	// Initialize contexts
	enhancer.userContext = &UserContext{
		Preferences: make(map[string]string),
		Permissions: make([]string, 0),
		Metadata:    make(map[string]string),
	}

	enhancer.sessionContext = &SessionContext{
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		Actions:      make([]string, 0),
		FeatureUsage: make(map[string]int),
		Metadata:     make(map[string]string),
	}

	return enhancer
}

// collectSystemInfo collects system information
func (ere *ErrorReportEnhancer) collectSystemInfo() *SystemInfo {
	info := &SystemInfo{
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		GoVersion:   runtime.Version(),
		AppVersion:  "1.0.0", // This would come from build info
		BuildTime:   time.Now().Format(time.RFC3339),
		Environment: make(map[string]string),
		Metadata:    make(map[string]string),
	}

	// Collect CPU information
	info.CPUCount = runtime.NumCPU()

	// Collect memory information (approximate)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	info.MemoryTotal = memStats.TotalAlloc
	info.MemoryFree = memStats.Alloc

	// Collect environment information
	// Note: runtime.GOROOT() is deprecated since Go 1.24
	// For demo purposes, we'll use a placeholder
	info.Environment["GOROOT"] = "runtime.GOROOT() deprecated - use 'go env GOROOT'"

	return info
}

// EnhanceReport enhances an error report with additional context
func (ere *ErrorReportEnhancer) EnhanceReport(report *ErrorReport) *EnhancedErrorReport {
	ere.mu.RLock()
	defer ere.mu.RUnlock()

	// Create enhanced report
	enhanced := &EnhancedErrorReport{
		ErrorReport:    report,
		SystemInfo:     ere.systemInfo,
		UserContext:    ere.userContext,
		SessionContext: ere.sessionContext,
		CollectedAt:    time.Now(),
		ContextualData: make(map[string]interface{}),
	}

	// Add contextual data based on error
	ere.addContextualData(enhanced)

	// Add stack trace if not present
	if enhanced.ErrorReport.Error.StackTrace == "" {
		enhanced.ErrorReport.Error.StackTrace = ere.captureStackTrace()
	}

	return enhanced
}

// addContextualData adds contextual data based on the error
func (ere *ErrorReportEnhancer) addContextualData(enhanced *EnhancedErrorReport) {
	err := enhanced.Error

	// Add error-specific context
	switch err.Category {
	case CategoryFile:
		ere.addFileErrorContext(enhanced)
	case CategoryDatabase:
		ere.addDatabaseErrorContext(enhanced)
	case CategoryNetwork:
		ere.addNetworkErrorContext(enhanced)
	case CategoryResource:
		ere.addResourceErrorContext(enhanced)
	case CategoryUI:
		ere.addUIErrorContext(enhanced)
	}

	// Add operation context
	if err.Operation != "" {
		enhanced.ContextualData["operation_context"] = map[string]interface{}{
			"operation":   err.Operation,
			"component":   err.Component,
			"duration_ms": err.Duration.Milliseconds(),
			"timestamp":   err.Timestamp,
		}
	}

	// Add recovery context
	if err.RecoveryAttempts > 0 {
		enhanced.ContextualData["recovery_context"] = map[string]interface{}{
			"strategy":     err.Recovery,
			"attempts":     err.RecoveryAttempts,
			"max_retries":  err.MaxRetries,
			"is_retryable": err.IsRetryable(),
		}
	}

	// Add session context
	ere.sessionContext.ErrorCount++
	ere.sessionContext.LastActivity = time.Now()
	ere.sessionContext.Duration = time.Since(ere.sessionContext.StartTime)

	// Track feature usage if component is specified
	if err.Component != "" {
		ere.sessionContext.FeatureUsage[err.Component]++
	}
}

// addFileErrorContext adds context for file errors
func (ere *ErrorReportEnhancer) addFileErrorContext(enhanced *EnhancedErrorReport) {
	err := enhanced.Error

	if filepath, exists := err.Metadata["filepath"]; exists {
		enhanced.ContextualData["file_context"] = map[string]interface{}{
			"filepath":    filepath,
			"exists":      ere.fileExists(filepath),
			"permissions": ere.getFilePermissions(filepath),
			"size_bytes":  ere.getFileSize(filepath),
		}
	}
}

// addDatabaseErrorContext adds context for database errors
func (ere *ErrorReportEnhancer) addDatabaseErrorContext(enhanced *EnhancedErrorReport) {
	enhanced.ContextualData["database_context"] = map[string]interface{}{
		"connection_status": "unknown", // Would check actual connection
		"pool_status":       "unknown", // Would check connection pool
		"last_operation":    enhanced.Error.Operation,
	}
}

// addNetworkErrorContext adds context for network errors
func (ere *ErrorReportEnhancer) addNetworkErrorContext(enhanced *EnhancedErrorReport) {
	enhanced.ContextualData["network_context"] = map[string]interface{}{
		"connectivity": ere.checkConnectivity(),
		"dns_status":   "unknown", // Would check DNS resolution
		"proxy_config": ere.getProxyConfig(),
	}
}

// addResourceErrorContext adds context for resource errors
func (ere *ErrorReportEnhancer) addResourceErrorContext(enhanced *EnhancedErrorReport) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	enhanced.ContextualData["resource_context"] = map[string]interface{}{
		"memory_usage_mb": memStats.Alloc / 1024 / 1024,
		"memory_total_mb": memStats.TotalAlloc / 1024 / 1024,
		"goroutines":      runtime.NumGoroutine(),
		"gc_pauses":       memStats.NumGC,
		"resource_type":   ere.getResourceType(enhanced.Error),
	}
}

// addUIErrorContext adds context for UI errors
func (ere *ErrorReportEnhancer) addUIErrorContext(enhanced *EnhancedErrorReport) {
	enhanced.ContextualData["ui_context"] = map[string]interface{}{
		"component":     enhanced.Error.Component,
		"ui_framework":  "bubbletea", // Would detect actual framework
		"terminal_size": "unknown",   // Would get terminal dimensions
		"color_support": "unknown",   // Would check color support
	}
}

// Helper methods for context collection

func (ere *ErrorReportEnhancer) fileExists(filepath string) bool {
	// This would check if file exists (simplified for demo)
	return true
}

func (ere *ErrorReportEnhancer) getFilePermissions(filepath string) string {
	// This would get file permissions (simplified for demo)
	return "readable"
}

func (ere *ErrorReportEnhancer) getFileSize(filepath string) int64 {
	// This would get file size (simplified for demo)
	return 0
}

func (ere *ErrorReportEnhancer) checkConnectivity() string {
	// This would check network connectivity (simplified for demo)
	return "online"
}

func (ere *ErrorReportEnhancer) getProxyConfig() string {
	// This would get proxy configuration (simplified for demo)
	return "none"
}

func (ere *ErrorReportEnhancer) getResourceType(err *AppError) string {
	if resource, exists := err.Metadata["resource"]; exists {
		return resource
	}
	return "unknown"
}

func (ere *ErrorReportEnhancer) captureStackTrace() string {
	// Capture stack trace (simplified for demo)
	return "Stack trace would be captured here"
}

// EnhancedErrorReport extends ErrorReport with additional context
type EnhancedErrorReport struct {
	*ErrorReport
	SystemInfo     *SystemInfo            `json:"system_info"`
	UserContext    *UserContext           `json:"user_context"`
	SessionContext *SessionContext        `json:"session_context"`
	CollectedAt    time.Time              `json:"collected_at"`
	ContextualData map[string]interface{} `json:"contextual_data"`
	Analysis       *ErrorAnalysis         `json:"analysis,omitempty"`
	Suggestions    []ErrorSuggestion      `json:"suggestions,omitempty"`
}

// ErrorAnalysis contains analysis of the error
type ErrorAnalysis struct {
	RootCause      string                 `json:"root_cause"`
	RelatedErrors  []*AppError            `json:"related_errors"`
	PatternMatch   string                 `json:"pattern_match"`
	SeverityScore  int                    `json:"severity_score"`
	RecoveryChance float64                `json:"recovery_chance"`
	ContextualTags []string               `json:"contextual_tags"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// ErrorSuggestion represents a suggestion for resolving the error
type ErrorSuggestion struct {
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Action      string                 `json:"action"`
	Priority    int                    `json:"priority"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// EnhancedErrorReporter extends the basic error reporter with enhancement capabilities
type EnhancedErrorReporter struct {
	*CompositeErrorReporter
	enhancer       *ErrorReportEnhancer
	logger         *logging.Logger
	analysisEngine *ErrorAnalysisEngine
}

// NewEnhancedErrorReporter creates a new enhanced error reporter
func NewEnhancedErrorReporter(logger *logging.Logger) *EnhancedErrorReporter {
	composite := NewCompositeErrorReporter([]ErrorReporter{}, logger)
	enhancer := NewErrorReportEnhancer(logger)

	reporter := &EnhancedErrorReporter{
		CompositeErrorReporter: composite,
		enhancer:               enhancer,
		logger:                 logger,
		analysisEngine:         NewErrorAnalysisEngine(logger),
	}

	return reporter
}

// Report sends an enhanced error report to all configured reporters
func (eer *EnhancedErrorReporter) Report(ctx context.Context, report *ErrorReport) error {
	// Skip if error is nil
	if report.Error == nil {
		return nil
	}

	// Enhance the report
	enhanced := eer.enhancer.EnhanceReport(report)

	// Perform error analysis
	analysis := eer.analysisEngine.AnalyzeError(enhanced)
	enhanced.Analysis = analysis

	// Generate suggestions
	suggestions := eer.analysisEngine.GenerateSuggestions(enhanced)
	enhanced.Suggestions = suggestions

	// Log enhanced report
	eer.logger.Info("Enhanced error report generated",
		"error_code", report.Error.Code,
		"severity_score", analysis.SeverityScore,
		"suggestions_count", len(suggestions),
	)

	// Send to all configured reporters
	return eer.CompositeErrorReporter.Report(ctx, enhanced.ErrorReport)
}

// ErrorAnalysisEngine performs analysis on errors
type ErrorAnalysisEngine struct {
	patterns map[string]*ErrorPattern
	logger   *logging.Logger
}

// ErrorPattern represents a pattern for error analysis
type ErrorPattern struct {
	Name        string
	Description string
	Indicators  []string
	RootCause   string
	Solutions   []string
}

// NewErrorAnalysisEngine creates a new error analysis engine
func NewErrorAnalysisEngine(logger *logging.Logger) *ErrorAnalysisEngine {
	engine := &ErrorAnalysisEngine{
		patterns: make(map[string]*ErrorPattern),
		logger:   logger,
	}

	engine.registerDefaultPatterns()
	return engine
}

// registerDefaultPatterns registers default error patterns
func (eae *ErrorAnalysisEngine) registerDefaultPatterns() {
	patterns := []*ErrorPattern{
		{
			Name:        "file_permission_denied",
			Description: "File permission errors",
			Indicators:  []string{"permission denied", "access denied", "forbidden"},
			RootCause:   "Insufficient file system permissions",
			Solutions:   []string{"check_file_permissions", "run_as_administrator", "change_file_location"},
		},
		{
			Name:        "database_connection_failed",
			Description: "Database connection failures",
			Indicators:  []string{"connection refused", "timeout", "unreachable"},
			RootCause:   "Database server unavailable or misconfigured",
			Solutions:   []string{"check_database_server", "verify_connection_string", "check_network_connectivity"},
		},
		{
			Name:        "memory_exhaustion",
			Description: "Memory exhaustion errors",
			Indicators:  []string{"out of memory", "allocation failed", "heap exhausted"},
			RootCause:   "Insufficient memory or memory leak",
			Solutions:   []string{"restart_application", "increase_memory_limit", "check_for_memory_leaks"},
		},
		{
			Name:        "file_corruption",
			Description: "File corruption issues",
			Indicators:  []string{"corrupt", "invalid format", "parse error"},
			RootCause:   "File data corruption or format mismatch",
			Solutions:   []string{"restore_from_backup", "check_file_integrity", "recreate_file"},
		},
	}

	for _, pattern := range patterns {
		eae.patterns[pattern.Name] = pattern
	}
}

// AnalyzeError analyzes an error and returns analysis results
func (eae *ErrorAnalysisEngine) AnalyzeError(report *EnhancedErrorReport) *ErrorAnalysis {
	analysis := &ErrorAnalysis{
		ContextualTags: make([]string, 0),
		Metadata:       make(map[string]interface{}),
	}

	err := report.Error

	// Calculate severity score
	analysis.SeverityScore = eae.calculateSeverityScore(err)

	// Find matching pattern
	pattern := eae.findMatchingPattern(err)
	if pattern != nil {
		analysis.PatternMatch = pattern.Name
		analysis.RootCause = pattern.RootCause
		analysis.ContextualTags = append(analysis.ContextualTags, pattern.Indicators...)
	}

	// Estimate recovery chance
	analysis.RecoveryChance = eae.estimateRecoveryChance(err, report)

	// Add metadata
	analysis.Metadata["error_category"] = string(err.Category)
	analysis.Metadata["error_severity"] = string(err.Severity)
	analysis.Metadata["has_stack_trace"] = err.StackTrace != ""
	analysis.Metadata["recovery_strategy"] = string(err.Recovery)

	return analysis
}

// calculateSeverityScore calculates a numeric severity score
func (eae *ErrorAnalysisEngine) calculateSeverityScore(err *AppError) int {
	baseScore := map[ErrorSeverity]int{
		SeverityLow:      1,
		SeverityMedium:   3,
		SeverityHigh:     7,
		SeverityCritical: 10,
	}

	score := baseScore[err.Severity]

	// Adjust based on category
	categoryMultiplier := map[ErrorCategory]float64{
		CategoryValidation:    0.5,
		CategoryUI:            0.8,
		CategoryConfiguration: 1.2,
		CategoryPermission:    1.5,
		CategoryFile:          1.8,
		CategoryNetwork:       2.0,
		CategoryResource:      2.5,
		CategoryDatabase:      3.0,
		CategoryParsing:       2.2,
		CategoryUnknown:       1.0,
	}

	if multiplier, exists := categoryMultiplier[err.Category]; exists {
		score = int(float64(score) * multiplier)
	}

	// Adjust based on recovery strategy
	if err.Recovery == RecoveryNone {
		score += 2
	} else if err.Recovery == RecoveryManual {
		score++
	}

	return score
}

// findMatchingPattern finds a matching error pattern
func (eae *ErrorAnalysisEngine) findMatchingPattern(err *AppError) *ErrorPattern {
	message := err.Message

	for _, pattern := range eae.patterns {
		for _, indicator := range pattern.Indicators {
			if len(message) >= len(indicator) {
				// Simple substring matching (in production, use regex)
				for i := 0; i <= len(message)-len(indicator); i++ {
					match := true
					for j := 0; j < len(indicator); j++ {
						if message[i+j] != indicator[j] {
							match = false
							break
						}
					}
					if match {
						return pattern
					}
				}
			}
		}
	}

	return nil
}

// estimateRecoveryChance estimates the chance of successful recovery
func (eae *ErrorAnalysisEngine) estimateRecoveryChance(err *AppError, report *EnhancedErrorReport) float64 {
	baseChance := 0.5

	// Adjust based on recovery strategy
	switch err.Recovery {
	case RecoveryNone:
		baseChance = 0.0
	case RecoveryManual:
		baseChance = 0.3
	case RecoveryRetry:
		baseChance = 0.7
	case RecoveryFallback:
		baseChance = 0.8
	case RecoveryGraceful:
		baseChance = 0.9
	}

	// Adjust based on error category
	categoryAdjustment := map[ErrorCategory]float64{
		CategoryValidation:    0.2,
		CategoryUI:            0.1,
		CategoryConfiguration: -0.1,
		CategoryPermission:    -0.2,
		CategoryFile:          0.0,
		CategoryNetwork:       -0.1,
		CategoryResource:      -0.2,
		CategoryDatabase:      -0.3,
		CategoryParsing:       -0.1,
		CategoryUnknown:       0.0,
	}

	if adjustment, exists := categoryAdjustment[err.Category]; exists {
		baseChance += adjustment
	}

	// Adjust based on retry attempts
	if err.RecoveryAttempts > 0 {
		baseChance -= float64(err.RecoveryAttempts) * 0.1
	}

	// Clamp between 0 and 1
	if baseChance < 0 {
		baseChance = 0
	}
	if baseChance > 1 {
		baseChance = 1
	}

	return baseChance
}

// GenerateSuggestions generates suggestions for resolving the error
func (eae *ErrorAnalysisEngine) GenerateSuggestions(report *EnhancedErrorReport) []ErrorSuggestion {
	var suggestions []ErrorSuggestion

	err := report.Error
	analysis := report.Analysis

	if analysis == nil {
		return suggestions
	}

	// Generate suggestions based on pattern
	if analysis.PatternMatch != "" {
		if pattern, exists := eae.patterns[analysis.PatternMatch]; exists {
			for i, solution := range pattern.Solutions {
				suggestions = append(suggestions, ErrorSuggestion{
					Type:        "solution",
					Title:       fmt.Sprintf("Solution %d", i+1),
					Description: solution,
					Action:      eae.getActionForSolution(solution),
					Priority:    10 - i,
				})
			}
		}
	}

	// Generate suggestions based on recovery strategy
	if err.Recovery == RecoveryRetry && err.IsRetryable() {
		suggestions = append(suggestions, ErrorSuggestion{
			Type:        "recovery",
			Title:       "Retry Operation",
			Description: "The operation can be retried automatically",
			Action:      "retry",
			Priority:    10,
		})
	}

	// Generate suggestions based on error category
	switch err.Category {
	case CategoryFile:
		suggestions = append(suggestions, ErrorSuggestion{
			Type:        "diagnostic",
			Title:       "Check File System",
			Description: "Verify file system integrity and available space",
			Action:      "check_filesystem",
			Priority:    8,
		})
	case CategoryDatabase:
		suggestions = append(suggestions, ErrorSuggestion{
			Type:        "diagnostic",
			Title:       "Check Database Connection",
			Description: "Verify database server is running and accessible",
			Action:      "check_database",
			Priority:    9,
		})
	case CategoryResource:
		suggestions = append(suggestions, ErrorSuggestion{
			Type:        "optimization",
			Title:       "Free System Resources",
			Description: "Close other applications and free up memory",
			Action:      "free_resources",
			Priority:    7,
		})
	}

	return suggestions
}

// getActionForSolution returns an action string for a solution
func (eae *ErrorAnalysisEngine) getActionForSolution(solution string) string {
	actionMap := map[string]string{
		"check_file_permissions":     "check_permissions",
		"run_as_administrator":       "run_as_admin",
		"change_file_location":       "change_location",
		"check_database_server":      "check_database",
		"verify_connection_string":   "verify_config",
		"check_network_connectivity": "check_network",
		"restart_application":        "restart",
		"increase_memory_limit":      "increase_memory",
		"check_for_memory_leaks":     "check_memory",
		"restore_from_backup":        "restore_backup",
		"check_file_integrity":       "check_integrity",
		"recreate_file":              "recreate",
	}

	if action, exists := actionMap[solution]; exists {
		return action
	}

	return "generic_solution"
}

// Global enhanced error reporter instance
var globalEnhancedErrorReporter *EnhancedErrorReporter

// InitializeEnhancedReporting initializes the global enhanced error reporter
func InitializeEnhancedReporting(logger *logging.Logger) {
	globalEnhancedErrorReporter = NewEnhancedErrorReporter(logger)
	logger.Info("Enhanced error reporting initialized")
}

// GetGlobalEnhancedErrorReporter returns the global enhanced error reporter
func GetGlobalEnhancedErrorReporter() *EnhancedErrorReporter {
	return globalEnhancedErrorReporter
}
