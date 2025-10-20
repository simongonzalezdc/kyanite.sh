package errors

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/logging"
)

// ActionableError represents an error with suggested actions
type ActionableError struct {
	*AppError
	SuggestedActions []UserAction
	ContextInfo      map[string]string
	RelatedErrors    []*AppError
}

// UserAction represents an action the user can take to resolve an error
type UserAction struct {
	ID          string            `json:"id"`
	Label       string            `json:"label"`
	Description string            `json:"description"`
	Type        ActionType        `json:"type"`
	Priority    int               `json:"priority"` // Higher number = higher priority
	Handler     func() error      `json:"-"`
	Parameters  map[string]string `json:"parameters,omitempty"`
}

// ActionType represents the type of user action
type ActionType string

const (
	ActionRetry     ActionType = "retry"
	ActionRestore   ActionType = "restore"
	ActionConfigure ActionType = "configure"
	ActionContact   ActionType = "contact"
	ActionIgnore    ActionType = "ignore"
	ActionRestart   ActionType = "restart"
	ActionCheck     ActionType = "check"
	ActionUpdate    ActionType = "update"
	ActionBackup    ActionType = "backup"
	ActionRecover   ActionType = "recover"
)

// EnhancedNotificationManager extends the basic notification manager with actionable errors
type EnhancedNotificationManager struct {
	*NotificationManager
	actionHandlers map[string]func(action *UserAction) error
	logger         *logging.Logger
}

// NewEnhancedNotificationManager creates a new enhanced notification manager
func NewEnhancedNotificationManager(logger logging.Logger) *EnhancedNotificationManager {
	baseManager := NewNotificationManager(logger)

	return &EnhancedNotificationManager{
		NotificationManager: baseManager,
		actionHandlers:      make(map[string]func(action *UserAction) error),
		logger:              &logger,
	}
}

// ShowActionableError displays an error with suggested actions
func (enm *EnhancedNotificationManager) ShowActionableError(title string, err *AppError) string {
	actionableErr := enm.createActionableError(err)

	// Create notification with actions
	notification := &Notification{
		ID:        generateNotificationID(),
		Type:      enm.getNotificationType(err),
		Level:     enm.getLevelFromError(err),
		Title:     title,
		Message:   enm.generateActionableMessage(actionableErr),
		Error:     err,
		Actions:   enm.convertActions(actionableErr.SuggestedActions),
		Duration:  enm.getDurationFromError(err),
		CreatedAt: time.Now(),
	}

	return enm.ShowNotification(notification)
}

// createActionableError creates an actionable error from an AppError
func (enm *EnhancedNotificationManager) createActionableError(err *AppError) *ActionableError {
	actionable := &ActionableError{
		AppError:         err,
		SuggestedActions: make([]UserAction, 0),
		ContextInfo:      make(map[string]string),
		RelatedErrors:    make([]*AppError, 0),
	}

	// Generate suggested actions based on error category and type
	enm.generateSuggestedActions(actionable)

	// Add context information
	enm.addContextInfo(actionable)

	return actionable
}

// generateSuggestedActions generates suggested actions based on error characteristics
func (enm *EnhancedNotificationManager) generateSuggestedActions(actionable *ActionableError) {
	err := actionable.AppError

	switch err.Category {
	case CategoryFile:
		enm.generateFileErrorActions(actionable)
	case CategoryDatabase:
		enm.generateDatabaseErrorActions(actionable)
	case CategoryNetwork:
		enm.generateNetworkErrorActions(actionable)
	case CategoryPermission:
		enm.generatePermissionErrorActions(actionable)
	case CategoryResource:
		enm.generateResourceErrorActions(actionable)
	case CategoryParsing:
		enm.generateParsingErrorActions(actionable)
	case CategoryValidation:
		enm.generateValidationErrorActions(actionable)
	case CategoryConfiguration:
		enm.generateConfigurationErrorActions(actionable)
	default:
		enm.generateGenericErrorActions(actionable)
	}
}

// generateFileErrorActions generates actions for file-related errors
func (enm *EnhancedNotificationManager) generateFileErrorActions(actionable *ActionableError) {
	err := actionable.AppError

	// Add retry action for recoverable file errors
	if err.CanRecover(RecoveryRetry) {
		actionable.SuggestedActions = append(actionable.SuggestedActions, UserAction{
			ID:          "retry_operation",
			Label:       "Retry",
			Description: "Retry the file operation",
			Type:        ActionRetry,
			Priority:    10,
			Handler: func() error {
				enm.logger.Info("User initiated retry for file operation")
				return nil // Handler would be implemented by caller
			},
		})
	}

	// Add recovery action for corrupted files
	if strings.Contains(err.Message, "corrupt") || strings.Contains(err.Message, "invalid") {
		actionable.SuggestedActions = append(actionable.SuggestedActions, UserAction{
			ID:          "recover_file",
			Label:       "Recover File",
			Description: "Attempt to recover the file from backup",
			Type:        ActionRecover,
			Priority:    9,
			Handler: func() error {
				enm.logger.Info("User initiated file recovery")
				// This would integrate with the corruption detector
				return nil
			},
		})
	}

	// Add check file action
	actionable.SuggestedActions = append(actionable.SuggestedActions, UserAction{
		ID:          "check_file",
		Label:       "Check File",
		Description: "Open file location to check manually",
		Type:        ActionCheck,
		Priority:    5,
		Handler: func() error {
			enm.logger.Info("User requested to check file location")
			return nil
		},
	})
}

// generateDatabaseErrorActions generates actions for database-related errors
func (enm *EnhancedNotificationManager) generateDatabaseErrorActions(actionable *ActionableError) {
	err := actionable.AppError

	// Add retry action for retryable database errors
	if err.CanRecover(RecoveryRetry) {
		actionable.SuggestedActions = append(actionable.SuggestedActions, UserAction{
			ID:          "retry_db_operation",
			Label:       "Retry",
			Description: "Retry the database operation",
			Type:        ActionRetry,
			Priority:    10,
			Handler: func() error {
				enm.logger.Info("User initiated database operation retry")
				return nil
			},
		})
	}

	// Add restart action for connection issues
	if strings.Contains(err.Message, "connection") || strings.Contains(err.Message, "timeout") {
		actionable.SuggestedActions = append(actionable.SuggestedActions, UserAction{
			ID:          "restart_application",
			Label:       "Restart Application",
			Description: "Restart the application to reestablish database connection",
			Type:        ActionRestart,
			Priority:    8,
			Handler: func() error {
				enm.logger.Info("User requested application restart")
				return nil
			},
		})
	}
}

// generateNetworkErrorActions generates actions for network-related errors
func (enm *EnhancedNotificationManager) generateNetworkErrorActions(actionable *ActionableError) {
	actionable.SuggestedActions = append(actionable.SuggestedActions, []UserAction{
		{
			ID:          "check_connection",
			Label:       "Check Connection",
			Description: "Check your internet connection",
			Type:        ActionCheck,
			Priority:    10,
			Handler: func() error {
				enm.logger.Info("User requested connection check")
				return nil
			},
		},
		{
			ID:          "retry_network",
			Label:       "Retry",
			Description: "Retry the network operation",
			Type:        ActionRetry,
			Priority:    9,
			Handler: func() error {
				enm.logger.Info("User initiated network retry")
				return nil
			},
		},
	}...)
}

// generatePermissionErrorActions generates actions for permission-related errors
func (enm *EnhancedNotificationManager) generatePermissionErrorActions(actionable *ActionableError) {
	actionable.SuggestedActions = append(actionable.SuggestedActions, UserAction{
		ID:          "check_permissions",
		Label:       "Check Permissions",
		Description: "Check file or folder permissions",
		Type:        ActionCheck,
		Priority:    10,
		Handler: func() error {
			enm.logger.Info("User requested permission check")
			return nil
		},
	})
}

// generateResourceErrorActions generates actions for resource-related errors
func (enm *EnhancedNotificationManager) generateResourceErrorActions(actionable *ActionableError) {
	actionable.SuggestedActions = append(actionable.SuggestedActions, []UserAction{
		{
			ID:          "free_memory",
			Label:       "Free Memory",
			Description: "Close other applications to free up memory",
			Type:        ActionConfigure,
			Priority:    8,
			Handler: func() error {
				enm.logger.Info("User requested memory optimization")
				return nil
			},
		},
		{
			ID:          "restart_application",
			Label:       "Restart Application",
			Description: "Restart the application to free up resources",
			Type:        ActionRestart,
			Priority:    7,
			Handler: func() error {
				enm.logger.Info("User requested application restart for resource issues")
				return nil
			},
		},
	}...)
}

// generateParsingErrorActions generates actions for parsing-related errors
func (enm *EnhancedNotificationManager) generateParsingErrorActions(actionable *ActionableError) {
	actionable.SuggestedActions = append(actionable.SuggestedActions, []UserAction{
		{
			ID:          "check_file_format",
			Label:       "Check File Format",
			Description: "Ensure the file is in the correct format",
			Type:        ActionCheck,
			Priority:    10,
			Handler: func() error {
				enm.logger.Info("User requested file format check")
				return nil
			},
		},
		{
			ID:          "recover_file",
			Label:       "Recover File",
			Description: "Attempt to recover the corrupted file",
			Type:        ActionRecover,
			Priority:    9,
			Handler: func() error {
				enm.logger.Info("User initiated file recovery for parsing error")
				return nil
			},
		},
	}...)
}

// generateValidationErrorActions generates actions for validation-related errors
func (enm *EnhancedNotificationManager) generateValidationErrorActions(actionable *ActionableError) {
	actionable.SuggestedActions = append(actionable.SuggestedActions, UserAction{
		ID:          "fix_validation",
		Label:       "Fix Issues",
		Description: "Review and fix the validation errors",
		Type:        ActionConfigure,
		Priority:    10,
		Handler: func() error {
			enm.logger.Info("User requested validation fix")
			return nil
		},
	})
}

// generateConfigurationErrorActions generates actions for configuration-related errors
func (enm *EnhancedNotificationManager) generateConfigurationErrorActions(actionable *ActionableError) {
	actionable.SuggestedActions = append(actionable.SuggestedActions, []UserAction{
		{
			ID:          "check_config",
			Label:       "Check Configuration",
			Description: "Review your configuration settings",
			Type:        ActionConfigure,
			Priority:    10,
			Handler: func() error {
				enm.logger.Info("User requested configuration check")
				return nil
			},
		},
		{
			ID:          "reset_config",
			Label:       "Reset to Defaults",
			Description: "Reset configuration to default values",
			Type:        ActionConfigure,
			Priority:    5,
			Handler: func() error {
				enm.logger.Info("User requested configuration reset")
				return nil
			},
		},
	}...)
}

// generateGenericErrorActions generates actions for generic/unknown errors
func (enm *EnhancedNotificationManager) generateGenericErrorActions(actionable *ActionableError) {
	actionable.SuggestedActions = append(actionable.SuggestedActions, []UserAction{
		{
			ID:          "retry_generic",
			Label:       "Retry",
			Description: "Retry the operation",
			Type:        ActionRetry,
			Priority:    7,
			Handler: func() error {
				enm.logger.Info("User initiated generic retry")
				return nil
			},
		},
		{
			ID:          "contact_support",
			Label:       "Contact Support",
			Description: "Contact support if the issue persists",
			Type:        ActionContact,
			Priority:    3,
			Handler: func() error {
				enm.logger.Info("User requested support contact")
				return nil
			},
		},
	}...)
}

// addContextInfo adds helpful context information to the actionable error
func (enm *EnhancedNotificationManager) addContextInfo(actionable *ActionableError) {
	err := actionable.AppError

	// Add operation context
	if err.Operation != "" {
		actionable.ContextInfo["operation"] = err.Operation
	}

	// Add component context
	if err.Component != "" {
		actionable.ContextInfo["component"] = err.Component
	}

	// Add timing context
	if !err.Timestamp.IsZero() {
		actionable.ContextInfo["timestamp"] = err.Timestamp.Format(time.RFC3339)
	}

	// Add recovery context
	if err.RecoveryAttempts > 0 {
		actionable.ContextInfo["recovery_attempts"] = fmt.Sprintf("%d", err.RecoveryAttempts)
	}

	// Add specific context based on error category
	switch err.Category {
	case CategoryFile:
		if filepath, exists := err.Metadata["filepath"]; exists {
			actionable.ContextInfo["file_path"] = filepath
		}
	case CategoryDatabase:
		if err.Operation != "" {
			actionable.ContextInfo["db_operation"] = err.Operation
		}
	case CategoryNetwork:
		actionable.ContextInfo["network_issue"] = "true"
	case CategoryResource:
		if resource, exists := err.Metadata["resource"]; exists {
			actionable.ContextInfo["resource_type"] = resource
		}
	}
}

// generateActionableMessage generates a user-friendly message with actionable guidance
func (enm *EnhancedNotificationManager) generateActionableMessage(actionable *ActionableError) string {
	err := actionable.AppError

	// Start with the base message
	message := err.Message

	// Add context information
	if len(actionable.ContextInfo) > 0 {
		contextParts := make([]string, 0, len(actionable.ContextInfo))
		for key, value := range actionable.ContextInfo {
			contextParts = append(contextParts, fmt.Sprintf("%s: %s", key, value))
		}
		message += "\n\nContext: " + strings.Join(contextParts, ", ")
	}

	// Add action guidance
	if len(actionable.SuggestedActions) > 0 {
		message += "\n\nSuggested actions:"
		for _, action := range actionable.SuggestedActions {
			if action.Priority >= 7 { // Only show high-priority actions in message
				message += fmt.Sprintf("\nâ€¢ %s: %s", action.Label, action.Description)
			}
		}
	}

	// Add recovery information
	if err.Recovery != RecoveryNone {
		message += fmt.Sprintf("\n\nRecovery strategy: %s", err.Recovery)
		if err.RecoveryAttempts > 0 {
			message += fmt.Sprintf(" (attempted %d times)", err.RecoveryAttempts)
		}
	}

	return message
}

// convertActions converts UserActions to NotificationActions
func (enm *EnhancedNotificationManager) convertActions(userActions []UserAction) []NotificationAction {
	actions := make([]NotificationAction, len(userActions))

	for i, userAction := range userActions {
		actions[i] = NotificationAction{
			ID:    userAction.ID,
			Label: userAction.Label,
			Type:  string(userAction.Type),
			Handler: func() error {
				return enm.executeAction(&userAction)
			},
		}
	}

	return actions
}

// executeAction executes a user action
func (enm *EnhancedNotificationManager) executeAction(action *UserAction) error {
	enm.logger.Info("Executing user action", "action_id", action.ID, "action_type", action.Type)

	if action.Handler != nil {
		return action.Handler()
	}

	return nil
}

// getNotificationType determines the notification type from an error
func (enm *EnhancedNotificationManager) getNotificationType(err *AppError) NotificationType {
	switch err.Severity {
	case SeverityCritical:
		return NotificationError
	case SeverityHigh:
		return NotificationError
	case SeverityMedium:
		return NotificationWarning
	default:
		return NotificationInfo
	}
}

// Global enhanced notification manager
var globalEnhancedNotificationManager *EnhancedNotificationManager

// InitializeEnhancedNotifications initializes the global enhanced notification manager
func InitializeEnhancedNotifications(logger logging.Logger) {
	globalEnhancedNotificationManager = NewEnhancedNotificationManager(logger)
	logger.Info("Enhanced notifications initialized")
}

// GetGlobalEnhancedNotificationManager returns the global enhanced notification manager
func GetGlobalEnhancedNotificationManager() *EnhancedNotificationManager {
	return globalEnhancedNotificationManager
}

// ShowGlobalActionableError displays an actionable error using the global manager
func ShowGlobalActionableError(title string, err *AppError) string {
	if globalEnhancedNotificationManager != nil {
		return globalEnhancedNotificationManager.ShowActionableError(title, err)
	}
	return ""
}
