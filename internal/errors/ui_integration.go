package errors

import (
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/lyricforge/internal/logging"
)

// ErrorRecoveryUI handles the UI integration for error recovery
type ErrorRecoveryUI struct {
	gracefulDegradation *EnhancedGracefulDegradation
	corruptionDetector  *FileCorruptionDetector
	logger              *logging.Logger

	// UI state
	showRecoveryPanel  bool
	selectedRecovery   int
	recoveryOperations []RecoveryOperation
	healthStatus       SystemHealthStatus

	// UI components
	spinner     spinner.Model
	progressBar *ProgressBar

	mu sync.RWMutex
}

// RecoveryOperation represents an ongoing recovery operation
type RecoveryOperation struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Status      RecoveryStatus         `json:"status"`
	Progress    float64                `json:"progress"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     *time.Time             `json:"end_time,omitempty"`
	Error       *AppError              `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// RecoveryStatus represents the status of a recovery operation
type RecoveryStatus string

const (
	StatusPending   RecoveryStatus = "pending"
	StatusRunning   RecoveryStatus = "running"
	StatusCompleted RecoveryStatus = "completed"
	StatusFailed    RecoveryStatus = "failed"
	StatusCancelled RecoveryStatus = "cancelled"
)

// SystemHealthStatus represents overall system health
type SystemHealthStatus struct {
	OverallScore    int                      `json:"overall_score"`
	FeatureHealth   map[string]FeatureHealth `json:"feature_health"`
	LastUpdated     time.Time                `json:"last_updated"`
	Recommendations []string                 `json:"recommendations"`
}

// FeatureHealth represents the health of a specific feature
type FeatureHealth struct {
	Score         int        `json:"score"`
	Status        string     `json:"status"`
	LastCheck     time.Time  `json:"last_check"`
	Issues        []string   `json:"issues"`
	DegradedSince *time.Time `json:"degraded_since,omitempty"`
}

// ProgressBar represents a simple progress bar component
type ProgressBar struct {
	Width    int
	Progress float64
	Style    lipgloss.Style
}

// NewErrorRecoveryUI creates a new error recovery UI handler
func NewErrorRecoveryUI(logger *logging.Logger) *ErrorRecoveryUI {
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Create progress bar style
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1).
		Margin(0, 1)

	ui := &ErrorRecoveryUI{
		logger:             logger,
		showRecoveryPanel:  false,
		selectedRecovery:   0,
		recoveryOperations: make([]RecoveryOperation, 0),
		spinner:            s,
		progressBar: &ProgressBar{
			Width:    50,
			Progress: 0.0,
			Style:    progressStyle,
		},
	}

	// Initialize health status
	ui.healthStatus = SystemHealthStatus{
		OverallScore:    100,
		FeatureHealth:   make(map[string]FeatureHealth),
		LastUpdated:     time.Now(),
		Recommendations: make([]string, 0),
	}

	return ui
}

// InitializeUIIntegration initializes the UI integration with all error recovery components
func InitializeUIIntegration(logger *logging.Logger) error {
	// Initialize enhanced notification manager
	InitializeEnhancedNotifications(*logger)

	// Initialize enhanced graceful degradation
	InitializeEnhancedGracefulDegradation(logger)

	// Initialize corruption detector
	InitializeCorruptionDetector(logger)

	// Initialize enhanced reporting
	InitializeEnhancedReporting(logger)

	logger.Info("Error recovery UI integration initialized")
	return nil
}

// ShowRecoveryPanel shows the error recovery panel
func (eru *ErrorRecoveryUI) ShowRecoveryPanel() {
	eru.mu.Lock()
	defer eru.mu.Unlock()

	eru.showRecoveryPanel = true
	eru.refreshRecoveryOperations()
	eru.updateSystemHealth()
}

// HideRecoveryPanel hides the error recovery panel
func (eru *ErrorRecoveryUI) HideRecoveryPanel() {
	eru.mu.Lock()
	defer eru.mu.Unlock()

	eru.showRecoveryPanel = false
}

// ToggleRecoveryPanel toggles the visibility of the recovery panel
func (eru *ErrorRecoveryUI) ToggleRecoveryPanel() {
	eru.mu.Lock()
	defer eru.mu.Unlock()

	eru.showRecoveryPanel = !eru.showRecoveryPanel

	if eru.showRecoveryPanel {
		eru.refreshRecoveryOperations()
		eru.updateSystemHealth()
	}
}

// refreshRecoveryOperations refreshes the list of recovery operations
func (eru *ErrorRecoveryUI) refreshRecoveryOperations() {
	eru.recoveryOperations = []RecoveryOperation{}

	// Add file corruption recovery operations
	corruptedFiles, err := eru.corruptionDetector.ScanDirectory(".")
	if err == nil && len(corruptedFiles) > 0 {
		eru.recoveryOperations = append(eru.recoveryOperations, RecoveryOperation{
			ID:          "file_recovery",
			Type:        "file_recovery",
			Description: fmt.Sprintf("Recover %d corrupted files", len(corruptedFiles)),
			Status:      StatusPending,
			Progress:    0.0,
			StartTime:   time.Now(),
			Metadata:    map[string]interface{}{"files": corruptedFiles},
		})
	}

	// Add degraded feature recovery operations
	degradedFeatures := eru.gracefulDegradation.GetDegradedFeatures()
	for _, feature := range degradedFeatures {
		if feature.State == FeatureDegraded {
			eru.recoveryOperations = append(eru.recoveryOperations, RecoveryOperation{
				ID:          fmt.Sprintf("feature_recovery_%s", feature.Name),
				Type:        "feature_recovery",
				Description: fmt.Sprintf("Recover %s feature", feature.Description),
				Status:      StatusPending,
				Progress:    0.0,
				StartTime:   time.Now(),
				Metadata:    map[string]interface{}{"feature": feature.Name},
			})
		}
	}
}

// updateSystemHealth updates the overall system health status
func (eru *ErrorRecoveryUI) updateSystemHealth() {
	health := eru.gracefulDegradation.GetSystemHealth()

	eru.healthStatus.OverallScore = health["health_score"].(int)
	eru.healthStatus.LastUpdated = time.Now()

	// Update feature health
	eru.healthStatus.FeatureHealth = make(map[string]FeatureHealth)
	for featureName, featureInfo := range eru.gracefulDegradation.features {
		var status string
		switch featureInfo.State {
		case FeatureEnabled:
			status = "healthy"
		case FeatureDegraded:
			status = "degraded"
		case FeatureDisabled:
			status = "disabled"
		case FeatureFailed:
			status = "failed"
		}

		featureHealth := FeatureHealth{
			Score:     eru.calculateFeatureScore(featureInfo),
			Status:    status,
			LastCheck: time.Now(),
			Issues:    []string{},
		}

		if featureInfo.LastFailure != nil {
			featureHealth.Issues = append(featureHealth.Issues, featureInfo.FailureReason)
		}

		if featureInfo.State == FeatureDegraded {
			featureHealth.DegradedSince = featureInfo.LastFailure
		}

		eru.healthStatus.FeatureHealth[featureName] = featureHealth
	}

	// Generate recommendations
	eru.generateRecommendations()
}

// calculateFeatureScore calculates a health score for a feature
func (eru *ErrorRecoveryUI) calculateFeatureScore(feature *FeatureInfo) int {
	baseScore := 100

	// Deduct points for failures
	baseScore -= feature.FailureCount * 10

	// Deduct points for degraded state
	if feature.State == FeatureDegraded {
		baseScore -= 30
	} else if feature.State == FeatureDisabled {
		baseScore -= 50
	} else if feature.State == FeatureFailed {
		baseScore -= 80
	}

	// Bonus points for recent recovery
	if feature.LastFailure == nil {
		baseScore += 10
	} else if time.Since(*feature.LastFailure) > time.Hour {
		baseScore += 5
	}

	if baseScore < 0 {
		baseScore = 0
	}
	if baseScore > 100 {
		baseScore = 100
	}

	return baseScore
}

// generateRecommendations generates system health recommendations
func (eru *ErrorRecoveryUI) generateRecommendations() {
	eru.healthStatus.Recommendations = []string{}

	if eru.healthStatus.OverallScore < 50 {
		eru.healthStatus.Recommendations = append(eru.healthStatus.Recommendations,
			"System health is critical. Consider restarting the application.")
	}

	degradedCount := 0
	for _, feature := range eru.healthStatus.FeatureHealth {
		if feature.Status == "degraded" || feature.Status == "failed" {
			degradedCount++
		}
	}

	if degradedCount > 0 {
		eru.healthStatus.Recommendations = append(eru.healthStatus.Recommendations,
			fmt.Sprintf("%d features are degraded. Check the recovery panel for details.", degradedCount))
	}

	if len(eru.recoveryOperations) > 0 {
		eru.healthStatus.Recommendations = append(eru.healthStatus.Recommendations,
			"Recovery operations are available. Open the recovery panel to resolve issues.")
	}
}

// ExecuteRecoveryOperation executes a recovery operation
func (eru *ErrorRecoveryUI) ExecuteRecoveryOperation(operationID string) error {
	eru.mu.Lock()
	defer eru.mu.Unlock()

	// Find the operation
	var operation *RecoveryOperation
	for i := range eru.recoveryOperations {
		if eru.recoveryOperations[i].ID == operationID {
			operation = &eru.recoveryOperations[i]
			break
		}
	}

	if operation == nil {
		return fmt.Errorf("recovery operation not found: %s", operationID)
	}

	// Update status
	operation.Status = StatusRunning
	eru.logger.Info("Starting recovery operation", "operation_id", operationID, "type", operation.Type)

	// Execute based on type
	var err error
	switch operation.Type {
	case "file_recovery":
		err = eru.executeFileRecovery(operation)
	case "feature_recovery":
		err = eru.executeFeatureRecovery(operation)
	default:
		err = fmt.Errorf("unknown recovery operation type: %s", operation.Type)
	}

	// Update operation status
	if err != nil {
		operation.Status = StatusFailed
		operation.Error = NewAppError("RECOVERY_FAILED",
			fmt.Sprintf("Recovery operation failed: %v", err), err,
			CategoryFile, SeverityMedium, RecoveryManual)
		eru.logger.Error("Recovery operation failed", "operation_id", operationID, "error", err)
	} else {
		operation.Status = StatusCompleted
		now := time.Now()
		operation.EndTime = &now
		operation.Progress = 100.0
		eru.logger.Info("Recovery operation completed", "operation_id", operationID)
	}

	return err
}

// executeFileRecovery executes file recovery operations
func (eru *ErrorRecoveryUI) executeFileRecovery(operation *RecoveryOperation) error {
	files, ok := operation.Metadata["files"].([]string)
	if !ok {
		return fmt.Errorf("no files specified for recovery")
	}

	recovered := 0
	for _, file := range files {
		if err := eru.corruptionDetector.RecoverFile(file); err != nil {
			eru.logger.Error("Failed to recover file", "file", file, "error", err)
		} else {
			recovered++
			operation.Progress = float64(recovered) / float64(len(files)) * 100.0
		}
	}

	if recovered == 0 {
		return fmt.Errorf("failed to recover any files")
	}

	return nil
}

// executeFeatureRecovery executes feature recovery operations
func (eru *ErrorRecoveryUI) executeFeatureRecovery(operation *RecoveryOperation) error {
	featureName, ok := operation.Metadata["feature"].(string)
	if !ok {
		return fmt.Errorf("no feature specified for recovery")
	}

	return eru.gracefulDegradation.AttemptFeatureRecovery(featureName)
}

// GetRecoveryOperations returns the current recovery operations
func (eru *ErrorRecoveryUI) GetRecoveryOperations() []RecoveryOperation {
	eru.mu.RLock()
	defer eru.mu.RUnlock()

	// Return a copy
	operations := make([]RecoveryOperation, len(eru.recoveryOperations))
	copy(operations, eru.recoveryOperations)
	return operations
}

// GetSystemHealth returns the current system health status
func (eru *ErrorRecoveryUI) GetSystemHealth() SystemHealthStatus {
	eru.mu.RLock()
	defer eru.mu.RUnlock()

	return eru.healthStatus
}

// Update handles UI updates for the error recovery system
func (eru *ErrorRecoveryUI) Update(msg tea.Msg) (*ErrorRecoveryUI, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if eru.showRecoveryPanel {
				eru.ToggleRecoveryPanel()
			}
		case "ctrl+r":
			eru.refreshRecoveryOperations()
		case "enter":
			if eru.showRecoveryPanel && eru.selectedRecovery < len(eru.recoveryOperations) {
				operation := eru.recoveryOperations[eru.selectedRecovery]
				if operation.Status == StatusPending {
					go func() {
						if err := eru.ExecuteRecoveryOperation(operation.ID); err != nil {
							eru.logger.Error("Recovery operation failed", "operation_id", operation.ID, "error", err)
						}
					}()
				}
			}
		case "j", "down":
			if eru.showRecoveryPanel && len(eru.recoveryOperations) > 0 {
				eru.selectedRecovery = (eru.selectedRecovery + 1) % len(eru.recoveryOperations)
			}
		case "k", "up":
			if eru.showRecoveryPanel && len(eru.recoveryOperations) > 0 {
				eru.selectedRecovery = (eru.selectedRecovery - 1 + len(eru.recoveryOperations)) % len(eru.recoveryOperations)
			}
		}
	case RecoveryStatusUpdateMsg:
		// Update recovery operation status
		eru.refreshRecoveryOperations()
		cmd = eru.spinner.Tick
	}

	return eru, cmd
}

// View renders the error recovery UI
func (eru *ErrorRecoveryUI) View() string {
	if !eru.showRecoveryPanel {
		return eru.renderStatusBar()
	}

	return eru.renderRecoveryPanel()
}

// renderStatusBar renders the status bar showing system health
func (eru *ErrorRecoveryUI) renderStatusBar() string {
	health := eru.GetSystemHealth()

	// Health indicator
	var healthColor string
	var healthIcon string

	switch {
	case health.OverallScore >= 80:
		healthColor = "#00FF00"
		healthIcon = "●"
	case health.OverallScore >= 60:
		healthColor = "#FFA500"
		healthIcon = "●"
	case health.OverallScore >= 40:
		healthColor = "#FF6600"
		healthIcon = "●"
	default:
		healthColor = "#FF0000"
		healthIcon = "●"
	}

	healthIndicator := lipgloss.NewStyle().
		Foreground(lipgloss.Color(healthColor)).
		Bold(true).
		Render(fmt.Sprintf("%s %d%%", healthIcon, health.OverallScore))

	// Recovery operations indicator
	operations := eru.GetRecoveryOperations()
	pendingOps := 0
	for _, op := range operations {
		if op.Status == StatusPending {
			pendingOps++
		}
	}

	var recoveryIndicator string
	if pendingOps > 0 {
		recoveryIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Render(fmt.Sprintf("⚠ %d recovery operations available (press 'r' to view)", pendingOps))
	}

	// Combine indicators
	indicators := []string{healthIndicator}
	if recoveryIndicator != "" {
		indicators = append(indicators, recoveryIndicator)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, indicators...)
}

// renderRecoveryPanel renders the full recovery panel
func (eru *ErrorRecoveryUI) renderRecoveryPanel() string {
	var sections []string

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00BFFF")).
		Render("🔧 Error Recovery Panel")
	sections = append(sections, header)

	// System health
	health := eru.GetSystemHealth()
	healthSection := eru.renderHealthSection(health)
	sections = append(sections, healthSection)

	// Recovery operations
	operations := eru.GetRecoveryOperations()
	if len(operations) > 0 {
		operationsSection := eru.renderOperationsSection(operations)
		sections = append(sections, operationsSection)
	} else {
		noOps := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Render("✅ No recovery operations needed")
		sections = append(sections, noOps)
	}

	// Footer
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Render("Press 'r' to close • ↑/↓ to navigate • Enter to execute • Ctrl+R to refresh")
	sections = append(sections, footer)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderHealthSection renders the system health section
func (eru *ErrorRecoveryUI) renderHealthSection(health SystemHealthStatus) string {
	var lines []string

	lines = append(lines, fmt.Sprintf("Overall Health: %d%%", health.OverallScore))

	// Feature health summary
	healthy := 0
	degraded := 0
	disabled := 0

	for _, feature := range health.FeatureHealth {
		switch feature.Status {
		case "healthy":
			healthy++
		case "degraded":
			degraded++
		case "disabled", "failed":
			disabled++
		}
	}

	lines = append(lines, fmt.Sprintf("Features: %d healthy, %d degraded, %d disabled",
		healthy, degraded, disabled))

	// Recommendations
	if len(health.Recommendations) > 0 {
		lines = append(lines, "Recommendations:")
		for _, rec := range health.Recommendations {
			lines = append(lines, fmt.Sprintf("  • %s", rec))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// renderOperationsSection renders the recovery operations section
func (eru *ErrorRecoveryUI) renderOperationsSection(operations []RecoveryOperation) string {
	var lines []string

	lines = append(lines, "Recovery Operations:")

	for i, op := range operations {
		selected := i == eru.selectedRecovery

		// Status indicator
		var statusIcon string
		var statusColor string

		switch op.Status {
		case StatusPending:
			statusIcon = "⏳"
			statusColor = "#FFA500"
		case StatusRunning:
			statusIcon = "🔄"
			statusColor = "#00BFFF"
		case StatusCompleted:
			statusIcon = "✅"
			statusColor = "#00FF00"
		case StatusFailed:
			statusIcon = "❌"
			statusColor = "#FF0000"
		case StatusCancelled:
			statusIcon = "⏹️"
			statusColor = "#888888"
		}

		// Progress bar for running operations
		var progressBar string
		if op.Status == StatusRunning {
			progressBar = eru.renderProgressBar(op.Progress)
		}

		// Format operation line
		indicator := "  "
		if selected {
			indicator = "▶ "
		}

		line := fmt.Sprintf("%s%s %s", indicator, statusIcon, op.Description)
		if progressBar != "" {
			line += "\n" + progressBar
		}

		// Apply selection styling
		if selected {
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFF00")).
				Background(lipgloss.Color("#333333")).
				Padding(0, 1).
				Render(line)
		} else {
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color(statusColor)).
				Render(line)
		}

		lines = append(lines, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// renderProgressBar renders a simple progress bar
func (eru *ErrorRecoveryUI) renderProgressBar(progress float64) string {
	filled := int(float64(eru.progressBar.Width) * progress / 100.0)
	empty := eru.progressBar.Width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}

	return eru.progressBar.Style.Render(fmt.Sprintf("[%s] %.1f%%", bar, progress))
}

// Message types for Bubble Tea

// RecoveryStatusUpdateMsg is sent when recovery operations are updated
type RecoveryStatusUpdateMsg struct {
	Operations []RecoveryOperation
	Timestamp  time.Time
}

// Global error recovery UI instance
var globalErrorRecoveryUI *ErrorRecoveryUI

// InitializeGlobalErrorRecoveryUI initializes the global error recovery UI
func InitializeGlobalErrorRecoveryUI(logger *logging.Logger) {
	globalErrorRecoveryUI = NewErrorRecoveryUI(logger)
	logger.Info("Global error recovery UI initialized")
}

// GetGlobalErrorRecoveryUI returns the global error recovery UI
func GetGlobalErrorRecoveryUI() *ErrorRecoveryUI {
	return globalErrorRecoveryUI
}

// Convenience functions for global access

// ShowGlobalRecoveryPanel shows the global recovery panel
func ShowGlobalRecoveryPanel() {
	if globalErrorRecoveryUI != nil {
		globalErrorRecoveryUI.ShowRecoveryPanel()
	}
}

// HideGlobalRecoveryPanel hides the global recovery panel
func HideGlobalRecoveryPanel() {
	if globalErrorRecoveryUI != nil {
		globalErrorRecoveryUI.HideRecoveryPanel()
	}
}

// ToggleGlobalRecoveryPanel toggles the global recovery panel
func ToggleGlobalRecoveryPanel() {
	if globalErrorRecoveryUI != nil {
		globalErrorRecoveryUI.ToggleRecoveryPanel()
	}
}
