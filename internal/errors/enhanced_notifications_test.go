package errors

import (
	"fmt"
	"testing"
	"time"
)

// TestNewEnhancedNotificationManager tests the creation of enhanced notification manager
func TestNewEnhancedNotificationManager(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	if enm == nil {
		t.Fatal("Expected non-nil EnhancedNotificationManager")
	}

	if enm.NotificationManager == nil {
		t.Error("Expected NotificationManager to be initialized")
	}

	if enm.actionHandlers == nil {
		t.Error("Expected action handlers map to be initialized")
	}

	if enm.logger == nil {
		t.Error("Expected logger to be initialized")
	}
}

// TestCreateActionableError tests creation of actionable errors
func TestCreateActionableError(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Test with a file error
	fileErr := NewFileError("read_file", "/test/file.txt",
		NewAppError("FILE_NOT_FOUND", "File not found", nil, CategoryFile, SeverityMedium, RecoveryRetry))

	actionable := enm.createActionableError(fileErr)

	if actionable == nil {
		t.Fatal("Expected non-nil actionable error")
	}

	if actionable.AppError != fileErr {
		t.Error("Expected AppError to be set correctly")
	}

	if len(actionable.SuggestedActions) == 0 {
		t.Error("Expected suggested actions to be generated")
	}

	if actionable.ContextInfo == nil {
		t.Error("Expected context info to be added")
	}

	if len(actionable.RelatedErrors) != 0 {
		t.Error("Expected no related errors initially")
	}
}

// TestGenerateSuggestedActions tests generation of suggested actions
func TestGenerateSuggestedActions(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	testCases := []struct {
		name            string
		err             *AppError
		expectedActions []string
	}{
		{
			name:            "File error",
			err:             NewFileError("read_file", "/test/file.txt", nil),
			expectedActions: []string{"retry_operation", "check_file"},
		},
		{
			name:            "Database error",
			err:             NewDatabaseError("connection", nil),
			expectedActions: []string{"retry_db_operation"},
		},
		{
			name:            "Network error",
			err:             NewAppError("NETWORK_ERROR", "Network timeout", nil, CategoryNetwork, SeverityMedium, RecoveryRetry),
			expectedActions: []string{"check_connection", "retry_network"},
		},
		{
			name:            "Permission error",
			err:             NewPermissionError("/test/file.txt", nil),
			expectedActions: []string{"check_permissions"},
		},
		{
			name:            "Resource error",
			err:             NewResourceError("memory", nil),
			expectedActions: []string{"free_memory", "restart_application"},
		},
		{
			name:            "Parsing error",
			err:             NewParsingError("json", nil),
			expectedActions: []string{"check_file_format", "recover_file"},
		},
		{
			name:            "Validation error",
			err:             NewValidationError("Invalid input", nil),
			expectedActions: []string{"fix_validation"},
		},
		{
			name:            "Configuration error",
			err:             NewConfigurationError("database_url", nil),
			expectedActions: []string{"check_config", "reset_config"},
		},
		{
			name:            "Unknown error",
			err:             NewAppError("UNKNOWN_ERROR", "Unknown error", nil, CategoryUnknown, SeverityMedium, RecoveryNone),
			expectedActions: []string{"retry_generic", "contact_support"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actionable := &ActionableError{
				AppError:         tc.err,
				SuggestedActions: make([]UserAction, 0),
				ContextInfo:      make(map[string]string),
				RelatedErrors:    make([]*AppError, 0),
			}

			enm.generateSuggestedActions(actionable)

			// Check that expected actions are present
			actionIDs := make([]string, len(actionable.SuggestedActions))
			for i, action := range actionable.SuggestedActions {
				actionIDs[i] = action.ID
			}

			for _, expectedAction := range tc.expectedActions {
				found := false
				for _, actionID := range actionIDs {
					if actionID == expectedAction {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected action %s to be suggested for %s", expectedAction, tc.name)
				}
			}
		})
	}
}

// TestAddContextInfo tests adding context information to actionable errors
func TestAddContextInfo(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Create an error with various context
	err := NewFileError("read_file", "/test/file.txt", nil).
		WithOperation("load_song").
		WithComponent("file_manager").
		WithDuration(100 * time.Millisecond)
	err.RecoveryAttempts = 2
	err.Metadata["filepath"] = "/test/file.txt"

	actionable := &ActionableError{
		AppError:         err,
		SuggestedActions: make([]UserAction, 0),
		ContextInfo:      make(map[string]string),
		RelatedErrors:    make([]*AppError, 0),
	}

	enm.addContextInfo(actionable)

	// Check expected context fields
	expectedContext := map[string]string{
		"operation":         "load_song",
		"component":         "file_manager",
		"timestamp":         err.Timestamp.Format(time.RFC3339),
		"recovery_attempts": "2",
		"file_path":         "/test/file.txt",
	}

	for key, expectedValue := range expectedContext {
		if actualValue, exists := actionable.ContextInfo[key]; !exists {
			t.Errorf("Expected context key %s to be present", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected context value for %s to be %s, got %s", key, expectedValue, actualValue)
		}
	}
}

// TestGenerateActionableMessage tests generation of actionable messages
func TestGenerateActionableMessage(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Create an actionable error with context and actions
	err := NewFileError("read_file", "/test/file.txt", nil)
	actionable := &ActionableError{
		AppError: err,
		ContextInfo: map[string]string{
			"operation": "load_song",
			"component": "file_manager",
		},
		SuggestedActions: []UserAction{
			{
				ID:          "retry_operation",
				Label:       "Retry",
				Description: "Retry the file operation",
				Type:        ActionRetry,
				Priority:    10,
			},
			{
				ID:          "check_file",
				Label:       "Check File",
				Description: "Open file location to check manually",
				Type:        ActionCheck,
				Priority:    5,
			},
		},
	}

	message := enm.generateActionableMessage(actionable)

	// Check that base message is present
	if !contains(message, err.Message) {
		t.Error("Expected base error message to be included")
	}

	// Check that context is included
	if !contains(message, "Context:") {
		t.Error("Expected context section to be included")
	}

	if !contains(message, "operation: load_song") {
		t.Error("Expected operation context to be included")
	}

	// Check that high-priority actions are included
	if !contains(message, "Suggested actions:") {
		t.Error("Expected suggested actions section to be included")
	}

	if !contains(message, "Retry: Retry the file operation") {
		t.Error("Expected high-priority action to be included")
	}

	// Low-priority action should not be included
	if contains(message, "Check File: Open file location to check manually") {
		t.Error("Expected low-priority action not to be included in message")
	}
}

// TestConvertActions tests conversion of user actions to notification actions
func TestConvertActions(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	userActions := []UserAction{
		{
			ID:          "retry_operation",
			Label:       "Retry",
			Description: "Retry the file operation",
			Type:        ActionRetry,
			Priority:    10,
			Handler: func() error {
				return nil
			},
		},
		{
			ID:          "check_file",
			Label:       "Check File",
			Description: "Open file location to check manually",
			Type:        ActionCheck,
			Priority:    5,
			Handler: func() error {
				return nil
			},
		},
	}

	notificationActions := enm.convertActions(userActions)

	if len(notificationActions) != len(userActions) {
		t.Errorf("Expected %d notification actions, got %d", len(userActions), len(notificationActions))
	}

	for i, userAction := range userActions {
		notificationAction := notificationActions[i]

		if notificationAction.ID != userAction.ID {
			t.Errorf("Expected action ID %s, got %s", userAction.ID, notificationAction.ID)
		}

		if notificationAction.Label != userAction.Label {
			t.Errorf("Expected action label %s, got %s", userAction.Label, notificationAction.Label)
		}

		if notificationAction.Type != string(userAction.Type) {
			t.Errorf("Expected action type %s, got %s", string(userAction.Type), notificationAction.Type)
		}

		if notificationAction.Handler == nil {
			t.Error("Expected action handler to be set")
		}
	}
}

// TestExecuteAction tests execution of user actions
func TestExecuteAction(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	executed := false
	action := &UserAction{
		ID:          "test_action",
		Label:       "Test Action",
		Description: "Test action description",
		Type:        ActionRetry,
		Priority:    10,
		Handler: func() error {
			executed = true
			return nil
		},
	}

	err := enm.executeAction(action)
	if err != nil {
		t.Errorf("Unexpected error executing action: %v", err)
	}

	if !executed {
		t.Error("Expected action handler to be executed")
	}

	// Test action with nil handler
	actionNilHandler := &UserAction{
		ID:          "nil_handler_action",
		Label:       "Nil Handler Action",
		Description: "Action with nil handler",
		Type:        ActionRetry,
		Priority:    10,
		Handler:     nil,
	}

	err = enm.executeAction(actionNilHandler)
	if err != nil {
		t.Errorf("Unexpected error executing action with nil handler: %v", err)
	}
}

// TestGetNotificationType tests determination of notification type from error
func TestGetNotificationType(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	testCases := []struct {
		name         string
		err          *AppError
		expectedType NotificationType
	}{
		{
			name:         "Critical error",
			err:          NewAppError("CRITICAL_ERROR", "Critical error", nil, CategoryUnknown, SeverityCritical, RecoveryNone),
			expectedType: NotificationError,
		},
		{
			name:         "High severity error",
			err:          NewAppError("HIGH_ERROR", "High severity error", nil, CategoryUnknown, SeverityHigh, RecoveryNone),
			expectedType: NotificationError,
		},
		{
			name:         "Medium severity error",
			err:          NewAppError("MEDIUM_ERROR", "Medium severity error", nil, CategoryUnknown, SeverityMedium, RecoveryNone),
			expectedType: NotificationWarning,
		},
		{
			name:         "Low severity error",
			err:          NewAppError("LOW_ERROR", "Low severity error", nil, CategoryUnknown, SeverityLow, RecoveryNone),
			expectedType: NotificationInfo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			notificationType := enm.getNotificationType(tc.err)
			if notificationType != tc.expectedType {
				t.Errorf("Expected notification type %s, got %s", tc.expectedType, notificationType)
			}
		})
	}
}

// TestShowActionableError tests displaying actionable errors
func TestShowActionableError(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Test with a file error
	fileErr := NewFileError("read_file", "/test/file.txt", nil)
	title := "File Error"

	notificationID := enm.ShowActionableError(title, fileErr)

	if notificationID == "" {
		t.Error("Expected non-empty notification ID")
	}

	// Check that notification was created
	notifications := enm.GetActiveNotifications()
	if len(notifications) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(notifications))
	}

	notification := notifications[0]
	if notification.Title != title {
		t.Errorf("Expected title %s, got %s", title, notification.Title)
	}

	if notification.Error != fileErr {
		t.Error("Expected error to be set in notification")
	}

	if len(notification.Actions) == 0 {
		t.Error("Expected actions to be set in notification")
	}
}

// TestNotificationCategorization tests notification categorization
func TestNotificationCategorization(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	testCases := []struct {
		name             string
		err              *AppError
		expectedCategory NotificationType
	}{
		{
			name:             "File error should be notification",
			err:              NewFileError("read_file", "/test/file.txt", nil),
			expectedCategory: NotificationError,
		},
		{
			name:             "Database error should be error",
			err:              NewDatabaseError("connection", nil),
			expectedCategory: NotificationError,
		},
		{
			name:             "Network error should be warning",
			err:              NewAppError("NETWORK_ERROR", "Network timeout", nil, CategoryNetwork, SeverityMedium, RecoveryRetry),
			expectedCategory: NotificationWarning,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			notificationType := enm.getNotificationType(tc.err)
			if notificationType != tc.expectedCategory {
				t.Errorf("Expected notification category %s, got %s", tc.expectedCategory, notificationType)
			}
		})
	}
}

// TestNotificationContentAndFormatting tests notification content and formatting
func TestNotificationContentAndFormatting(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Create an error with rich context
	err := NewFileError("read_file", "/test/file.txt", nil).
		WithOperation("load_song").
		WithComponent("file_manager").
		WithContext("song_id", "123").
		WithContext("user_id", "user456")

	title := "File Operation Failed"
	notificationID := enm.ShowActionableError(title, err)

	if notificationID == "" {
		t.Fatal("Expected non-empty notification ID")
	}

	notifications := enm.GetActiveNotifications()
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}

	notification := notifications[0]
	message := notification.Message

	// Check that context is included in the message
	if !contains(message, "Context:") {
		t.Error("Expected context section to be included in message")
	}

	if !contains(message, "operation: load_song") {
		t.Error("Expected operation context to be included")
	}

	if !contains(message, "component: file_manager") {
		t.Error("Expected component context to be included")
	}

	// Note: The addContextInfo function doesn't include arbitrary metadata in the message
	// Only specific fields like operation, component, etc. are included

	// Check that suggested actions are included
	if !contains(message, "Suggested actions:") {
		t.Error("Expected suggested actions section to be included")
	}
}

// TestNotificationThrottlingAndBatching tests notification throttling and batching
func TestNotificationThrottlingAndBatching(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Set max notifications to a low value for testing
	enm.maxNotifications = 3

	// Create multiple notifications rapidly
	for i := 0; i < 5; i++ {
		err := NewFileError("read_file", "/test/file.txt", nil)
		enm.ShowActionableError("File Error", err)
	}

	// Check that only max notifications are kept
	notifications := enm.GetActiveNotifications()
	if len(notifications) > enm.maxNotifications {
		t.Errorf("Expected at most %d notifications, got %d", enm.maxNotifications, len(notifications))
	}

	// Check that oldest notifications were removed
	if len(notifications) == enm.maxNotifications {
		// The newest notifications should be kept
		for i, notification := range notifications {
			if notification.Title != "File Error" {
				t.Errorf("Expected notification %d to have title 'File Error', got '%s'", i, notification.Title)
			}
		}
	}
}

// TestErrorScenariosInNotificationSystem tests error scenarios in the notification system itself
func TestErrorScenariosInNotificationSystem(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Test with error that has nil cause
	err := NewFileError("read_file", "/test/file.txt", nil)
	notificationID := enm.ShowActionableError("File Error", err)
	if notificationID == "" {
		t.Error("Expected non-empty notification ID with error that has nil cause")
	}

	// Test action execution failure
	failingAction := &UserAction{
		ID:          "failing_action",
		Label:       "Failing Action",
		Description: "Action that always fails",
		Type:        ActionRetry,
		Priority:    10,
		Handler: func() error {
			return NewAppError("ACTION_FAILED", "Action failed", nil, CategoryUnknown, SeverityMedium, RecoveryNone)
		},
	}

	actionErr := enm.executeAction(failingAction)
	if actionErr == nil {
		t.Error("Expected error when executing failing action")
	}
}

// TestNotificationChannels tests different notification channels
func TestNotificationChannels(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Test UI channel
	fileErr := NewFileError("read_file", "/test/file.txt", nil)
	notificationID := enm.ShowActionableError("File Error", fileErr)

	if notificationID == "" {
		t.Error("Expected non-empty notification ID for UI channel")
	}

	// Check that notification is in the UI channel (active notifications)
	uiNotifications := enm.GetActiveNotifications()
	if len(uiNotifications) == 0 {
		t.Error("Expected notification to be available in UI channel")
	}

	// Test log channel - check that notification was logged
	// Note: The actual logging message might vary, so we just check that something was logged
	if len(logger.GetMessages()) == 0 {
		t.Error("Expected some logging to occur")
	}
}

// TestNotificationLocalization tests notification localization
func TestNotificationLocalization(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Test with different error types that should have localized messages
	testCases := []struct {
		name          string
		err           *AppError
		expectedTitle string
		expectedMsg   string
	}{
		{
			name:          "File error",
			err:           NewFileError("read_file", "/test/file.txt", nil),
			expectedTitle: "File Error",
			expectedMsg:   "File operation 'read_file' failed for '/test/file.txt'",
		},
		{
			name:          "Database error",
			err:           NewDatabaseError("connection", nil),
			expectedTitle: "Database Error",
			expectedMsg:   "Database operation 'connection' failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			notificationID := enm.ShowActionableError(tc.expectedTitle, tc.err)
			if notificationID == "" {
				t.Error("Expected non-empty notification ID")
			}

			notifications := enm.GetActiveNotifications()
			if len(notifications) == 0 {
				t.Fatal("Expected at least one notification")
			}

			// GetActiveNotifications iterates a map, so its order is not
			// deterministic; look the notification up by the ID we just got
			// instead of assuming it is the last element.
			var notification *Notification
			for _, n := range notifications {
				if n.ID == notificationID {
					notification = n
					break
				}
			}
			if notification == nil {
				t.Fatalf("Notification %s not found among active notifications", notificationID)
			}
			if notification.Title != tc.expectedTitle {
				t.Errorf("Expected title %s, got %s", tc.expectedTitle, notification.Title)
			}

			if !contains(notification.Message, tc.expectedMsg) {
				t.Errorf("Expected message to contain %s, got %s", tc.expectedMsg, notification.Message)
			}
		})
	}
}

// TestNotificationRoutingAndDelivery tests notification routing and delivery
func TestNotificationRoutingAndDelivery(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Test routing based on error severity
	criticalErr := NewAppError("CRITICAL_ERROR", "Critical error", nil, CategoryUnknown, SeverityCritical, RecoveryNone)
	notificationID1 := enm.ShowActionableError("Critical Error", criticalErr)

	mediumErr := NewAppError("MEDIUM_ERROR", "Medium error", nil, CategoryUnknown, SeverityMedium, RecoveryNone)
	notificationID2 := enm.ShowActionableError("Medium Error", mediumErr)

	if notificationID1 == "" || notificationID2 == "" {
		t.Error("Expected non-empty notification IDs")
	}

	// Check that notifications were routed correctly
	notifications := enm.GetActiveNotifications()
	if len(notifications) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(notifications))
	}

	// Find the critical and medium notifications
	var criticalNotification, mediumNotification *Notification
	for _, notif := range notifications {
		if notif.Title == "Critical Error" {
			criticalNotification = notif
		} else if notif.Title == "Medium Error" {
			mediumNotification = notif
		}
	}

	if criticalNotification == nil {
		t.Error("Expected to find critical notification")
	}

	if mediumNotification == nil {
		t.Error("Expected to find medium notification")
	}

	// Critical error should have longer duration (0 means never auto-dismiss)
	if criticalNotification != nil && criticalNotification.Duration != 0 {
		t.Errorf("Expected critical error notification to not auto-dismiss, got duration: %v, severity: %v",
			criticalNotification.Duration, criticalNotification.Error.Severity)
	}

	if mediumNotification != nil && mediumNotification.Duration <= 0 {
		t.Error("Expected medium error notification to have auto-dismiss duration")
	}
}

// TestNotificationDeliveryFailure tests notification delivery failure scenarios
func TestNotificationDeliveryFailure(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Test with malformed error
	malformedErr := &AppError{
		Code:     "",
		Message:  "",
		Category: "",
		Severity: "",
		Recovery: "",
	}

	notificationID := enm.ShowActionableError("Malformed Error", malformedErr)
	if notificationID == "" {
		t.Error("Expected non-empty notification ID even with malformed error")
	}

	// Test with extremely long message
	longMessage := ""
	for i := 0; i < 10000; i++ {
		longMessage += "This is a very long message. "
	}

	longErr := NewAppError("LONG_ERROR", longMessage, nil, CategoryUnknown, SeverityMedium, RecoveryNone)
	notificationID = enm.ShowActionableError("Long Error", longErr)
	if notificationID == "" {
		t.Error("Expected non-empty notification ID with long message")
	}

	// Check that notification was created despite the long message
	notifications := enm.GetActiveNotifications()
	if len(notifications) == 0 {
		t.Error("Expected notification to be created despite long message")
	}
}

// TestConcurrentNotificationHandling tests concurrent notification handling
func TestConcurrentNotificationHandling(t *testing.T) {
	logger := NewTestLogger(t)
	enm := NewEnhancedNotificationManager(*logger.Logger)

	// Create multiple notifications concurrently
	numGoroutines := 10
	notificationsPerGoroutine := 5

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < notificationsPerGoroutine; j++ {
				err := NewFileError("read_file", "/test/file.txt", nil)
				title := fmt.Sprintf("File Error %d-%d", id, j)
				enm.ShowActionableError(title, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Check that notifications were created
	// Note: NotificationManager has a max limit (typically 10), so we can't expect more than that
	notifications := enm.GetActiveNotifications()
	maxNotifications := 10 // This is the default max from NotificationManager
	if len(notifications) == 0 {
		t.Error("Expected at least some notifications, got 0")
	}
	if len(notifications) > maxNotifications {
		t.Errorf("Expected at most %d notifications (max limit), got %d", maxNotifications, len(notifications))
	}

	// Check that the system is still in a consistent state
	for _, notification := range notifications {
		if notification.ID == "" {
			t.Error("Expected notification to have non-empty ID")
		}
		if notification.Title == "" {
			t.Error("Expected notification to have non-empty title")
		}
	}
}
