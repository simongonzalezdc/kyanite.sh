package errors

import (
	"fmt"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/logging"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	// NotificationError represents an error notification
	NotificationError NotificationType = "error"
	// NotificationWarning represents a warning notification
	NotificationWarning NotificationType = "warning"
	// NotificationInfo represents an informational notification
	NotificationInfo NotificationType = "info"
	// NotificationSuccess represents a success notification
	NotificationSuccess NotificationType = "success"
)

// NotificationLevel represents the severity level of a notification
type NotificationLevel int

const (
	LevelLow NotificationLevel = iota
	LevelMedium
	LevelHigh
	LevelCritical
)

// Notification represents a user notification
type Notification struct {
	ID        string               `json:"id"`
	Type      NotificationType     `json:"type"`
	Level     NotificationLevel    `json:"level"`
	Title     string               `json:"title"`
	Message   string               `json:"message"`
	Error     *AppError            `json:"error,omitempty"`
	Actions   []NotificationAction `json:"actions,omitempty"`
	Duration  time.Duration        `json:"duration"`
	CreatedAt time.Time            `json:"created_at"`
	Dismissed bool                 `json:"dismissed"`
}

// NotificationAction represents an action that can be taken on a notification
type NotificationAction struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Type    string `json:"type"` // "primary", "secondary", "danger"
	Handler func() error
}

// NotificationManager manages user notifications
type NotificationManager struct {
	notifications    map[string]*Notification
	logger           logging.Logger
	maxNotifications int
	defaultDuration  time.Duration
	mu               sync.RWMutex
	onUpdate         func([]*Notification)
}

// NotificationUIModel handles the UI for notifications
type NotificationUIModel struct {
	manager       *NotificationManager
	showDetails   bool
	selectedIndex int
	spinner       spinner.Model
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(logger logging.Logger) *NotificationManager {
	return &NotificationManager{
		notifications:    make(map[string]*Notification),
		logger:           logger,
		maxNotifications: 10,
		defaultDuration:  5 * time.Second,
	}
}

// SetOnUpdate sets a callback function to be called when notifications are updated
func (nm *NotificationManager) SetOnUpdate(callback func([]*Notification)) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	nm.onUpdate = callback
}

// ShowError displays an error notification
func (nm *NotificationManager) ShowError(title, message string, err *AppError) string {
	return nm.ShowNotification(&Notification{
		ID:        generateNotificationID(),
		Type:      NotificationError,
		Level:     nm.getLevelFromError(err),
		Title:     title,
		Message:   message,
		Error:     err,
		Duration:  nm.getDurationFromError(err),
		CreatedAt: time.Now(),
	})
}

// ShowWarning displays a warning notification
func (nm *NotificationManager) ShowWarning(title, message string) string {
	return nm.ShowNotification(&Notification{
		ID:        generateNotificationID(),
		Type:      NotificationWarning,
		Level:     LevelMedium,
		Title:     title,
		Message:   message,
		Duration:  4 * time.Second,
		CreatedAt: time.Now(),
	})
}

// ShowInfo displays an informational notification
func (nm *NotificationManager) ShowInfo(title, message string) string {
	return nm.ShowNotification(&Notification{
		ID:        generateNotificationID(),
		Type:      NotificationInfo,
		Level:     LevelLow,
		Title:     title,
		Message:   message,
		Duration:  3 * time.Second,
		CreatedAt: time.Now(),
	})
}

// ShowSuccess displays a success notification
func (nm *NotificationManager) ShowSuccess(title, message string) string {
	return nm.ShowNotification(&Notification{
		ID:        generateNotificationID(),
		Type:      NotificationSuccess,
		Level:     LevelLow,
		Title:     title,
		Message:   message,
		Duration:  3 * time.Second,
		CreatedAt: time.Now(),
	})
}

// ShowNotification displays a notification with custom settings
func (nm *NotificationManager) ShowNotification(notification *Notification) string {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	// Set default duration if not specified
	// Note: Duration of 0 is valid for critical errors (never auto-dismiss),
	// so we only apply the default if no error is attached or it's not critical
	if notification.Duration == 0 && (notification.Error == nil || notification.Error.Severity != SeverityCritical) {
		notification.Duration = nm.defaultDuration
	}

	// Add to notifications map
	nm.notifications[notification.ID] = notification

	// Limit number of notifications
	if len(nm.notifications) > nm.maxNotifications {
		nm.removeOldestNotification()
	}

	// Log the notification
	nm.logger.Infof("Notification shown: [%s] %s - %s", notification.Type, notification.Title, notification.Message)

	// Notify UI of update
	nm.notifyUpdate()

	return notification.ID
}

// DismissNotification dismisses a notification by ID
func (nm *NotificationManager) DismissNotification(id string) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if notification, exists := nm.notifications[id]; exists {
		notification.Dismissed = true
		nm.logger.Debugf("Notification dismissed: %s", id)
		nm.notifyUpdate()
	}
}

// DismissAllNotifications dismisses all notifications
func (nm *NotificationManager) DismissAllNotifications() {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	for id, notification := range nm.notifications {
		notification.Dismissed = true
		nm.logger.Debugf("Notification dismissed: %s", id)
	}

	nm.notifyUpdate()
}

// GetActiveNotifications returns all active (non-dismissed) notifications
func (nm *NotificationManager) GetActiveNotifications() []*Notification {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	var active []*Notification
	for _, notification := range nm.notifications {
		if !notification.Dismissed {
			active = append(active, notification)
		}
	}

	return active
}

// GetNotificationsByType returns notifications of a specific type
func (nm *NotificationManager) GetNotificationsByType(notificationType NotificationType) []*Notification {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	var filtered []*Notification
	for _, notification := range nm.notifications {
		if notification.Type == notificationType && !notification.Dismissed {
			filtered = append(filtered, notification)
		}
	}

	return filtered
}

// CleanupExpiredNotifications removes notifications that have exceeded their duration
func (nm *NotificationManager) CleanupExpiredNotifications() {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	now := time.Now()
	expired := make([]string, 0)

	for id, notification := range nm.notifications {
		if !notification.Dismissed && now.Sub(notification.CreatedAt) > notification.Duration {
			expired = append(expired, id)
		}
	}

	for _, id := range expired {
		delete(nm.notifications, id)
		nm.logger.Debugf("Expired notification removed: %s", id)
	}

	if len(expired) > 0 {
		nm.notifyUpdate()
	}
}

// Helper functions

func (nm *NotificationManager) getLevelFromError(err *AppError) NotificationLevel {
	if err == nil {
		return LevelMedium
	}

	switch err.Severity {
	case SeverityCritical:
		return LevelCritical
	case SeverityHigh:
		return LevelHigh
	case SeverityMedium:
		return LevelMedium
	default:
		return LevelLow
	}
}

func (nm *NotificationManager) getDurationFromError(err *AppError) time.Duration {
	if err == nil {
		return nm.defaultDuration
	}

	switch err.Severity {
	case SeverityCritical:
		return 0 // Never auto-dismiss critical errors
	case SeverityHigh:
		return 10 * time.Second
	case SeverityMedium:
		return 7 * time.Second
	default:
		return 5 * time.Second
	}
}

func (nm *NotificationManager) removeOldestNotification() {
	var oldestID string
	var oldestTime time.Time

	for id, notification := range nm.notifications {
		if oldestID == "" || notification.CreatedAt.Before(oldestTime) {
			oldestID = id
			oldestTime = notification.CreatedAt
		}
	}

	if oldestID != "" {
		delete(nm.notifications, oldestID)
		nm.logger.Debugf("Old notification removed: %s", oldestID)
	}
}

func (nm *NotificationManager) notifyUpdate() {
	if nm.onUpdate != nil {
		active := nm.GetActiveNotifications()
		nm.onUpdate(active)
	}
}

func generateNotificationID() string {
	return fmt.Sprintf("notif_%d", time.Now().UnixNano())
}

// UI Model for notifications

// NewNotificationUIModel creates a new notification UI model
func NewNotificationUIModel(manager *NotificationManager) *NotificationUIModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &NotificationUIModel{
		manager:       manager,
		showDetails:   false,
		selectedIndex: 0,
		spinner:       s,
	}
}

// Init initializes the notification UI model
func (m *NotificationUIModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the notification UI
func (m *NotificationUIModel) Update(msg tea.Msg) (*NotificationUIModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			m.showDetails = !m.showDetails
		case "j", "down":
			m.selectNext()
		case "k", "up":
			m.selectPrevious()
		case "enter":
			m.handleAction()
		case "esc":
			m.DismissSelected()
		}
	case NotificationUpdateMsg:
		// Update received from manager
		cmd = m.spinner.Tick
	}

	return m, cmd
}

// View renders the notification UI
func (m *NotificationUIModel) View() string {
	notifications := m.manager.GetActiveNotifications()

	if len(notifications) == 0 {
		return ""
	}

	var sections []string

	for i, notification := range notifications {
		selected := i == m.selectedIndex && m.showDetails
		section := m.renderNotification(notification, selected)
		sections = append(sections, section)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// Helper functions for UI

func (m *NotificationUIModel) renderNotification(notification *Notification, selected bool) string {
	var style lipgloss.Style

	// Choose style based on notification type
	switch notification.Type {
	case NotificationError:
		style = errorStyle
	case NotificationWarning:
		style = warningStyle
	case NotificationInfo:
		style = infoStyle
	case NotificationSuccess:
		style = successStyle
	default:
		style = infoStyle
	}

	// Add selection indicator if selected
	indicator := " "
	if selected {
		indicator = "â–¶ "
	}

	// Render title and message
	title := style.Render(indicator + notification.Title)
	message := notification.Message

	// Add error details if showing details and error exists
	if selected && notification.Error != nil {
		message += "\n\nError Details:"
		message += fmt.Sprintf("\n  Code: %s", notification.Error.Code)
		message += fmt.Sprintf("\n  Category: %s", notification.Error.Category)
		message += fmt.Sprintf("\n  Severity: %s", notification.Error.Severity)

		if notification.Error.Operation != "" {
			message += fmt.Sprintf("\n  Operation: %s", notification.Error.Operation)
		}

		if notification.Error.Component != "" {
			message += fmt.Sprintf("\n  Component: %s", notification.Error.Component)
		}

		// Add recovery suggestion if available
		if notification.Error.Recovery != RecoveryNone {
			message += fmt.Sprintf("\n  Recovery: %s", notification.Error.Recovery)
		}
	}

	content := title + "\n" + message

	// Add border and padding
	border := lipgloss.NormalBorder()
	if selected {
		border = lipgloss.ThickBorder()
	}

	return lipgloss.NewStyle().
		Border(border).
		Padding(0, 1).
		Margin(0, 0, 1, 0).
		Render(content)
}

func (m *NotificationUIModel) selectNext() {
	notifications := m.manager.GetActiveNotifications()
	if len(notifications) > 0 {
		m.selectedIndex = (m.selectedIndex + 1) % len(notifications)
	}
}

func (m *NotificationUIModel) selectPrevious() {
	notifications := m.manager.GetActiveNotifications()
	if len(notifications) > 0 {
		m.selectedIndex = (m.selectedIndex - 1 + len(notifications)) % len(notifications)
	}
}

func (m *NotificationUIModel) handleAction() {
	notifications := m.manager.GetActiveNotifications()
	if m.selectedIndex < len(notifications) {
		notification := notifications[m.selectedIndex]

		// Handle default actions based on notification type
		switch notification.Type {
		case NotificationError:
			if notification.Error != nil && notification.Error.CanRecover(RecoveryRetry) {
				// Retry the operation
				m.manager.logger.Info("Retrying operation...")
			}
		}
	}
}

func (m *NotificationUIModel) DismissSelected() {
	notifications := m.manager.GetActiveNotifications()
	if m.selectedIndex < len(notifications) {
		notification := notifications[m.selectedIndex]
		m.manager.DismissNotification(notification.ID)
	}
}

// Styles for different notification types

var (
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF")).
			Bold(false)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)
)

// Message types for Bubble Tea

// NotificationUpdateMsg is sent when notifications are updated
type NotificationUpdateMsg struct {
	Notifications []*Notification
	Timestamp     time.Time
}

// Global notification manager instance
var globalNotificationManager *NotificationManager

// SetGlobalNotificationManager sets the global notification manager
func SetGlobalNotificationManager(manager *NotificationManager) {
	globalNotificationManager = manager
}

// GetGlobalNotificationManager returns the global notification manager
func GetGlobalNotificationManager() *NotificationManager {
	return globalNotificationManager
}

// Convenience functions for global notifications

// ShowGlobalError displays an error notification using the global manager
func ShowGlobalError(title, message string, err *AppError) string {
	if globalNotificationManager != nil {
		return globalNotificationManager.ShowError(title, message, err)
	}
	return ""
}

// ShowGlobalWarning displays a warning notification using the global manager
func ShowGlobalWarning(title, message string) string {
	if globalNotificationManager != nil {
		return globalNotificationManager.ShowWarning(title, message)
	}
	return ""
}

// ShowGlobalInfo displays an informational notification using the global manager
func ShowGlobalInfo(title, message string) string {
	if globalNotificationManager != nil {
		return globalNotificationManager.ShowInfo(title, message)
	}
	return ""
}

// ShowGlobalSuccess displays a success notification using the global manager
func ShowGlobalSuccess(title, message string) string {
	if globalNotificationManager != nil {
		return globalNotificationManager.ShowSuccess(title, message)
	}
	return ""
}
