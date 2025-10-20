package errors

import (
	"fmt"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/logging"
)

// FeatureState represents the current state of a feature
type FeatureState int

const (
	FeatureEnabled FeatureState = iota
	FeatureDegraded
	FeatureDisabled
	FeatureFailed
)

// FeatureInfo contains information about a feature's state
type FeatureInfo struct {
	Name          string                 `json:"name"`
	State         FeatureState           `json:"state"`
	Description   string                 `json:"description"`
	FailureReason string                 `json:"failure_reason,omitempty"`
	LastFailure   *time.Time             `json:"last_failure,omitempty"`
	FailureCount  int                    `json:"failure_count"`
	RecoveryTime  *time.Time             `json:"recovery_time,omitempty"`
	Dependencies  []string               `json:"dependencies"`
	Impact        FeatureImpact          `json:"impact"`
	AutoRecovery  bool                   `json:"auto_recovery"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// FeatureImpact represents the impact of a feature failure
type FeatureImpact int

const (
	ImpactNone FeatureImpact = iota
	ImpactMinor
	ImpactModerate
	ImpactMajor
	ImpactCritical
)

// DegradationRule defines when and how to degrade a feature
type DegradationRule struct {
	FeatureName     string
	TriggerErrors   []ErrorCategory
	MaxFailures     int
	DegradationTime time.Duration
	RecoveryTime    time.Duration
	Impact          FeatureImpact
}

// EnhancedGracefulDegradation provides advanced graceful degradation capabilities
type EnhancedGracefulDegradation struct {
	*GracefulDegradation
	features        map[string]*FeatureInfo
	rules           map[string]*DegradationRule
	logger          *logging.Logger
	mu              sync.RWMutex
	healthChecker   *FeatureHealthChecker
	recoveryManager *FeatureRecoveryManager
}

// FeatureHealthChecker monitors feature health
type FeatureHealthChecker struct {
	features map[string]*FeatureInfo
	logger   *logging.Logger
}

// FeatureRecoveryManager handles automatic feature recovery
type FeatureRecoveryManager struct {
	features map[string]*FeatureInfo
	logger   *logging.Logger
}

// NewEnhancedGracefulDegradation creates a new enhanced graceful degradation handler
func NewEnhancedGracefulDegradation(logger *logging.Logger) *EnhancedGracefulDegradation {
	baseDegradation := NewGracefulDegradation(logger)

	egd := &EnhancedGracefulDegradation{
		GracefulDegradation: baseDegradation,
		features:            make(map[string]*FeatureInfo),
		rules:               make(map[string]*DegradationRule),
		logger:              logger,
	}

	// Initialize health checker and recovery manager
	egd.healthChecker = &FeatureHealthChecker{
		features: egd.features,
		logger:   logger,
	}

	egd.recoveryManager = &FeatureRecoveryManager{
		features: egd.features,
		logger:   logger,
	}

	// Register default degradation rules
	egd.registerDefaultRules()

	// Start background monitoring
	go egd.startMonitoring()

	return egd
}

// registerDefaultRules registers default degradation rules for common features
func (egd *EnhancedGracefulDegradation) registerDefaultRules() {
	defaultRules := []*DegradationRule{
		{
			FeatureName:     "auto_save",
			TriggerErrors:   []ErrorCategory{CategoryFile, CategoryDatabase, CategoryResource},
			MaxFailures:     3,
			DegradationTime: 30 * time.Second,
			RecoveryTime:    5 * time.Minute,
			Impact:          ImpactMinor,
		},
		{
			FeatureName:     "real_time_preview",
			TriggerErrors:   []ErrorCategory{CategoryResource, CategoryUI},
			MaxFailures:     2,
			DegradationTime: 15 * time.Second,
			RecoveryTime:    2 * time.Minute,
			Impact:          ImpactModerate,
		},
		{
			FeatureName:     "audio_playback",
			TriggerErrors:   []ErrorCategory{CategoryResource, CategoryFile},
			MaxFailures:     2,
			DegradationTime: 20 * time.Second,
			RecoveryTime:    3 * time.Minute,
			Impact:          ImpactModerate,
		},
		{
			FeatureName:     "database_operations",
			TriggerErrors:   []ErrorCategory{CategoryDatabase, CategoryNetwork},
			MaxFailures:     3,
			DegradationTime: 45 * time.Second,
			RecoveryTime:    10 * time.Minute,
			Impact:          ImpactMajor,
		},
		{
			FeatureName:     "file_operations",
			TriggerErrors:   []ErrorCategory{CategoryFile, CategoryPermission},
			MaxFailures:     2,
			DegradationTime: 20 * time.Second,
			RecoveryTime:    5 * time.Minute,
			Impact:          ImpactModerate,
		},
	}

	for _, rule := range defaultRules {
		egd.rules[rule.FeatureName] = rule

		// Initialize feature info
		egd.features[rule.FeatureName] = &FeatureInfo{
			Name:         rule.FeatureName,
			State:        FeatureEnabled,
			Description:  fmt.Sprintf("%s feature", rule.FeatureName),
			FailureCount: 0,
			Impact:       rule.Impact,
			AutoRecovery: true,
			Metadata:     make(map[string]interface{}),
		}
	}
}

// RegisterFeature registers a new feature for graceful degradation
func (egd *EnhancedGracefulDegradation) RegisterFeature(featureName, description string, impact FeatureImpact, dependencies []string) {
	egd.mu.Lock()
	defer egd.mu.Unlock()

	feature := &FeatureInfo{
		Name:         featureName,
		State:        FeatureEnabled,
		Description:  description,
		FailureCount: 0,
		Dependencies: dependencies,
		Impact:       impact,
		AutoRecovery: true,
		Metadata:     make(map[string]interface{}),
	}

	egd.features[featureName] = feature
	egd.logger.Info("Feature registered for graceful degradation", "feature", featureName, "impact", impact)
}

// HandleError handles an error and applies graceful degradation if needed
func (egd *EnhancedGracefulDegradation) HandleError(err *AppError) {
	egd.mu.Lock()
	defer egd.mu.Unlock()

	// Find affected features based on error category and component
	affectedFeatures := egd.findAffectedFeatures(err)

	for _, featureName := range affectedFeatures {
		feature, exists := egd.features[featureName]
		if !exists {
			continue
		}

		// Check if this error type should trigger degradation for this feature
		rule, hasRule := egd.rules[featureName]
		if hasRule && egd.shouldDegradeFeature(feature, rule, err) {
			egd.degradeFeature(feature, rule, err)
		} else {
			// Still record the failure for monitoring
			feature.FailureCount++
			now := time.Now()
			feature.LastFailure = &now
		}
	}
}

// findAffectedFeatures finds features affected by an error
func (egd *EnhancedGracefulDegradation) findAffectedFeatures(err *AppError) []string {
	var affected []string

	// Map error categories to features
	categoryFeatureMap := map[ErrorCategory][]string{
		CategoryFile:       {"file_operations", "auto_save"},
		CategoryDatabase:   {"database_operations", "auto_save"},
		CategoryNetwork:    {"database_operations"},
		CategoryResource:   {"real_time_preview", "audio_playback", "auto_save"},
		CategoryUI:         {"real_time_preview"},
		CategoryPermission: {"file_operations"},
	}

	// Add features based on error category
	if features, exists := categoryFeatureMap[err.Category]; exists {
		affected = append(affected, features...)
	}

	// Add features based on component
	if err.Component != "" {
		componentFeatureMap := map[string][]string{
			"editor":       {"auto_save", "real_time_preview"},
			"preview":      {"real_time_preview"},
			"audio":        {"audio_playback"},
			"database":     {"database_operations"},
			"file_manager": {"file_operations"},
		}

		if features, exists := componentFeatureMap[err.Component]; exists {
			affected = append(affected, features...)
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var unique []string
	for _, feature := range affected {
		if !seen[feature] {
			seen[feature] = true
			unique = append(unique, feature)
		}
	}

	return unique
}

// shouldDegradeFeature determines if a feature should be degraded
func (egd *EnhancedGracefulDegradation) shouldDegradeFeature(feature *FeatureInfo, rule *DegradationRule, err *AppError) bool {
	// Check if error category matches trigger errors
	for _, triggerCategory := range rule.TriggerErrors {
		if err.Category == triggerCategory {
			return true
		}
	}

	// Check failure count threshold
	if feature.FailureCount >= rule.MaxFailures {
		return true
	}

	// Check if feature is already degraded and still within degradation window
	if feature.State == FeatureDegraded && feature.RecoveryTime != nil {
		if time.Now().Before(*feature.RecoveryTime) {
			return false // Still in degradation period
		}
	}

	return false
}

// degradeFeature degrades a feature according to its rule
func (egd *EnhancedGracefulDegradation) degradeFeature(feature *FeatureInfo, rule *DegradationRule, err *AppError) {
	now := time.Now()

	// Update feature state
	feature.State = FeatureDegraded
	feature.FailureCount++
	feature.LastFailure = &now
	feature.FailureReason = err.Message

	// Set recovery time
	recoveryTime := now.Add(rule.RecoveryTime)
	feature.RecoveryTime = &recoveryTime

	// Log degradation
	egd.logger.Warn("Feature degraded due to errors",
		"feature", feature.Name,
		"state", feature.State,
		"failures", feature.FailureCount,
		"recovery_time", recoveryTime,
		"impact", feature.Impact,
	)

	// Show user notification based on impact
	switch feature.Impact {
	case ImpactCritical:
		ShowGlobalError("Critical Feature Degraded",
			fmt.Sprintf("The %s feature has been temporarily disabled due to repeated errors. This may significantly impact functionality.", feature.Description),
			NewAppError("FEATURE_DEGRADED", fmt.Sprintf("Feature %s degraded", feature.Name), err, CategoryUI, SeverityHigh, RecoveryGraceful))
	case ImpactMajor:
		ShowGlobalWarning("Major Feature Degraded",
			fmt.Sprintf("The %s feature has been temporarily degraded due to errors. Some functionality may be limited.", feature.Description))
	case ImpactModerate:
		ShowGlobalInfo("Feature Degraded",
			fmt.Sprintf("The %s feature has been temporarily degraded. Performance may be affected.", feature.Description))
	default:
		// Minor impact - just log, don't notify user
		egd.logger.Info("Feature degraded with minor impact", "feature", feature.Name)
	}

	// Check dependencies and degrade them if necessary
	egd.handleDependentFeatures(feature)
}

// handleDependentFeatures handles features that depend on a degraded feature
func (egd *EnhancedGracefulDegradation) handleDependentFeatures(feature *FeatureInfo) {
	for _, depName := range feature.Dependencies {
		if depFeature, exists := egd.features[depName]; exists {
			if depFeature.State == FeatureEnabled {
				// Create a dependency error
				depErr := NewAppError("DEPENDENCY_FAILED",
					fmt.Sprintf("Feature %s failed due to dependency on %s", depName, feature.Name),
					nil, CategoryUI, SeverityMedium, RecoveryGraceful)

				// Find the dependency's rule or use default
				depRule := egd.rules[depName]
				if depRule == nil {
					depRule = &DegradationRule{
						MaxFailures:     1,
						DegradationTime: 10 * time.Second,
						RecoveryTime:    1 * time.Minute,
						Impact:          ImpactModerate,
					}
				}

				egd.degradeFeature(depFeature, depRule, depErr)
			}
		}
	}
}

// IsFeatureAvailable checks if a feature is available (enabled and not degraded)
func (egd *EnhancedGracefulDegradation) IsFeatureAvailable(featureName string) bool {
	egd.mu.RLock()
	defer egd.mu.RUnlock()

	feature, exists := egd.features[featureName]
	if !exists {
		return true // Unknown features are considered available
	}

	return feature.State == FeatureEnabled
}

// GetFeatureState returns the current state of a feature
func (egd *EnhancedGracefulDegradation) GetFeatureState(featureName string) (FeatureState, error) {
	egd.mu.RLock()
	defer egd.mu.RUnlock()

	feature, exists := egd.features[featureName]
	if !exists {
		return FeatureEnabled, fmt.Errorf("feature not found: %s", featureName)
	}

	return feature.State, nil
}

// GetDegradedFeatures returns a list of currently degraded features
func (egd *EnhancedGracefulDegradation) GetDegradedFeatures() []*FeatureInfo {
	egd.mu.RLock()
	defer egd.mu.RUnlock()

	var degraded []*FeatureInfo
	for _, feature := range egd.features {
		if feature.State == FeatureDegraded || feature.State == FeatureDisabled {
			degraded = append(degraded, feature)
		}
	}

	return degraded
}

// AttemptFeatureRecovery attempts to recover a degraded feature
func (egd *EnhancedGracefulDegradation) AttemptFeatureRecovery(featureName string) error {
	egd.mu.Lock()
	defer egd.mu.Unlock()

	feature, exists := egd.features[featureName]
	if !exists {
		return fmt.Errorf("feature not found: %s", featureName)
	}

	if feature.State == FeatureEnabled {
		return fmt.Errorf("feature is not degraded: %s", featureName)
	}

	// Check if enough time has passed for recovery attempt
	if feature.RecoveryTime != nil && time.Now().Before(*feature.RecoveryTime) {
		return fmt.Errorf("feature is still in degradation period: %s", featureName)
	}

	// Attempt recovery
	feature.State = FeatureEnabled
	feature.FailureCount = 0
	feature.FailureReason = ""
	feature.LastFailure = nil
	feature.RecoveryTime = nil

	egd.logger.Info("Feature recovery attempted", "feature", featureName)

	// Show success notification
	ShowGlobalSuccess("Feature Recovered",
		fmt.Sprintf("The %s feature has been restored.", feature.Description))

	return nil
}

// startMonitoring starts background monitoring of feature health
func (egd *EnhancedGracefulDegradation) startMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		egd.performHealthCheck()
		egd.attemptAutoRecoveries()
	}
}

// performHealthCheck performs periodic health checks on all features
func (egd *EnhancedGracefulDegradation) performHealthCheck() {
	egd.mu.Lock()
	defer egd.mu.Unlock()

	now := time.Now()

	for _, feature := range egd.features {
		// Check if degraded features should be auto-recovered
		if feature.State == FeatureDegraded && feature.AutoRecovery {
			if feature.RecoveryTime != nil && now.After(*feature.RecoveryTime) {
				// Auto-recovery time reached
				feature.State = FeatureEnabled
				feature.FailureCount = 0
				feature.FailureReason = ""
				feature.RecoveryTime = nil

				egd.logger.Info("Feature auto-recovered", "feature", feature.Name)

				// Show info notification for auto-recovery
				ShowGlobalInfo("Feature Auto-Recovered",
					fmt.Sprintf("The %s feature has been automatically restored.", feature.Description))
			}
		}

		// Reset failure count for features that haven't failed recently
		if feature.LastFailure != nil && now.Sub(*feature.LastFailure) > time.Hour {
			if feature.FailureCount > 0 {
				feature.FailureCount--
				egd.logger.Debug("Feature failure count decremented", "feature", feature.Name, "count", feature.FailureCount)
			}
		}
	}
}

// attemptAutoRecoveries attempts automatic recovery of failed features
func (egd *EnhancedGracefulDegradation) attemptAutoRecoveries() {
	// This would implement automatic recovery strategies
	// For now, it's handled in performHealthCheck
}

// GetSystemHealth returns an overall system health assessment
func (egd *EnhancedGracefulDegradation) GetSystemHealth() map[string]interface{} {
	egd.mu.RLock()
	defer egd.mu.RUnlock()

	totalFeatures := len(egd.features)
	disabledFeatures := 0
	degradedFeatures := 0
	criticalFailures := 0

	for _, feature := range egd.features {
		switch feature.State {
		case FeatureDisabled:
			disabledFeatures++
		case FeatureDegraded:
			degradedFeatures++
		case FeatureFailed:
			criticalFailures++
		}
	}

	healthScore := 100
	if totalFeatures > 0 {
		penalty := (disabledFeatures * 50) + (degradedFeatures * 20) + (criticalFailures * 100)
		healthScore = max(0, 100-penalty/totalFeatures)
	}

	return map[string]interface{}{
		"health_score":      healthScore,
		"total_features":    totalFeatures,
		"enabled_features":  totalFeatures - disabledFeatures,
		"degraded_features": degradedFeatures,
		"disabled_features": disabledFeatures,
		"failed_features":   criticalFailures,
		"degraded_list":     egd.GetDegradedFeatures(),
	}
}

// Global enhanced graceful degradation instance
var globalEnhancedGracefulDegradation *EnhancedGracefulDegradation

// InitializeEnhancedGracefulDegradation initializes the global enhanced graceful degradation handler
func InitializeEnhancedGracefulDegradation(logger *logging.Logger) {
	globalEnhancedGracefulDegradation = NewEnhancedGracefulDegradation(logger)
	logger.Info("Enhanced graceful degradation initialized")
}

// GetGlobalEnhancedGracefulDegradation returns the global enhanced graceful degradation handler
func GetGlobalEnhancedGracefulDegradation() *EnhancedGracefulDegradation {
	return globalEnhancedGracefulDegradation
}

// Helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
