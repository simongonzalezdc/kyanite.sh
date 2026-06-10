package ui

import (
	"time"

	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// =============================================================================
// Constants
// =============================================================================

const (
	// ToastDurationDefault is the default duration a toast is visible
	ToastDurationDefault = 3 * time.Second
	// ToastDurationLong is for important messages
	ToastDurationLong = 5 * time.Second
	// ToastDurationShort is for quick feedback
	ToastDurationShort = 2 * time.Second
)

// =============================================================================
// Types
// =============================================================================

// ToastType represents the type of toast notification
type ToastType int

const (
	// ToastInfo is for informational messages
	ToastInfo ToastType = iota
	// ToastSuccess is for success messages
	ToastSuccess
	// ToastWarning is for warning messages
	ToastWarning
	// ToastError is for error messages
	ToastError
)

// ToastMsg is the message type for showing a toast
type ToastMsg struct {
	Type     ToastType
	Message  string
	Duration time.Duration
}

// toastDismissMsg is sent when a toast should be dismissed
type toastDismissMsg struct {
	id int
}

// Toast represents a single toast notification
type Toast struct {
	ID       int
	Type     ToastType
	Message  string
	Duration time.Duration
	ShowTime time.Time
}

// =============================================================================
// Toast Model
// =============================================================================

// ToastModel manages toast notifications in the TUI
type ToastModel struct {
	toasts     []Toast
	nextID     int
	width      int
	maxVisible int

	// Styles
	infoStyle      lipgloss.Style
	successStyle   lipgloss.Style
	warningStyle   lipgloss.Style
	errorStyle     lipgloss.Style
	containerStyle lipgloss.Style
}

// NewToastModel creates a new toast notification manager
func NewToastModel() *ToastModel {
	t := theme.GetManager().Current()

	return &ToastModel{
		toasts:     []Toast{},
		nextID:     0,
		width:      80,
		maxVisible: 3,

		infoStyle: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.Background).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Secondary).
			Padding(0, 1),

		successStyle: lipgloss.NewStyle().
			Foreground(t.Background).
			Background(t.Success).
			Bold(true).
			Padding(0, 1),

		warningStyle: lipgloss.NewStyle().
			Foreground(t.Background).
			Background(t.Accent).
			Bold(true).
			Padding(0, 1),

		errorStyle: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.Error).
			Bold(true).
			Padding(0, 1),

		containerStyle: lipgloss.NewStyle().
			Padding(0, 1),
	}
}

// Init initializes the toast model
func (m *ToastModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the toast model
func (m *ToastModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width

	case ToastMsg:
		return m.addToast(msg)

	case toastDismissMsg:
		m.removeToast(msg.id)
	}

	return nil
}

// View renders the toast notifications
func (m *ToastModel) View() string {
	if len(m.toasts) == 0 {
		return ""
	}

	// Render visible toasts (most recent at bottom)
	visible := m.toasts
	if len(visible) > m.maxVisible {
		visible = visible[len(visible)-m.maxVisible:]
	}

	var rendered []string
	for _, toast := range visible {
		rendered = append(rendered, m.renderToast(toast))
	}

	// Join toasts vertically
	result := lipgloss.JoinVertical(lipgloss.Left, rendered...)
	return m.containerStyle.Render(result)
}

// =============================================================================
// Public Methods
// =============================================================================

// SetWidth sets the width for toast rendering
func (m *ToastModel) SetWidth(width int) {
	m.width = width
}

// HasToasts returns true if there are visible toasts
func (m *ToastModel) HasToasts() bool {
	return len(m.toasts) > 0
}

// ClearAll removes all toasts
func (m *ToastModel) ClearAll() {
	m.toasts = []Toast{}
}

// UpdateTheme refreshes styles when theme changes
func (m *ToastModel) UpdateTheme() {
	t := theme.GetManager().Current()

	m.infoStyle = lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Background).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Secondary).
		Padding(0, 1)

	m.successStyle = lipgloss.NewStyle().
		Foreground(t.Background).
		Background(t.Success).
		Bold(true).
		Padding(0, 1)

	m.warningStyle = lipgloss.NewStyle().
		Foreground(t.Background).
		Background(t.Accent).
		Bold(true).
		Padding(0, 1)

	m.errorStyle = lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Error).
		Bold(true).
		Padding(0, 1)
}

// =============================================================================
// Private Methods
// =============================================================================

// addToast adds a new toast and returns a command to dismiss it later
func (m *ToastModel) addToast(msg ToastMsg) tea.Cmd {
	duration := msg.Duration
	if duration == 0 {
		duration = ToastDurationDefault
	}

	toast := Toast{
		ID:       m.nextID,
		Type:     msg.Type,
		Message:  msg.Message,
		Duration: duration,
		ShowTime: time.Now(),
	}
	m.nextID++
	m.toasts = append(m.toasts, toast)

	// Return command to dismiss this toast after duration
	id := toast.ID
	return tea.Tick(duration, func(time.Time) tea.Msg {
		return toastDismissMsg{id: id}
	})
}

// removeToast removes a toast by ID
func (m *ToastModel) removeToast(id int) {
	for i, toast := range m.toasts {
		if toast.ID == id {
			m.toasts = append(m.toasts[:i], m.toasts[i+1:]...)
			return
		}
	}
}

// renderToast renders a single toast
func (m *ToastModel) renderToast(toast Toast) string {
	icon := m.getIcon(toast.Type)
	content := icon + " " + toast.Message

	// Truncate if too long
	maxWidth := m.width - 4
	if len(content) > maxWidth && maxWidth > 10 {
		content = content[:maxWidth-3] + "..."
	}

	style := m.getStyle(toast.Type)
	return style.Width(m.width - 4).Render(content)
}

// getIcon returns the icon for a toast type
func (m *ToastModel) getIcon(t ToastType) string {
	switch t {
	case ToastSuccess:
		return "[OK]"
	case ToastWarning:
		return "[!]"
	case ToastError:
		return "[X]"
	default:
		return "[i]"
	}
}

// getStyle returns the style for a toast type
func (m *ToastModel) getStyle(t ToastType) lipgloss.Style {
	switch t {
	case ToastSuccess:
		return m.successStyle
	case ToastWarning:
		return m.warningStyle
	case ToastError:
		return m.errorStyle
	default:
		return m.infoStyle
	}
}

// =============================================================================
// Helper Functions for Creating Toast Messages
// =============================================================================

// InfoToast creates an info toast message
func InfoToast(message string) ToastMsg {
	return ToastMsg{
		Type:     ToastInfo,
		Message:  message,
		Duration: ToastDurationDefault,
	}
}

// SuccessToast creates a success toast message
func SuccessToast(message string) ToastMsg {
	return ToastMsg{
		Type:     ToastSuccess,
		Message:  message,
		Duration: ToastDurationDefault,
	}
}

// WarningToast creates a warning toast message
func WarningToast(message string) ToastMsg {
	return ToastMsg{
		Type:     ToastWarning,
		Message:  message,
		Duration: ToastDurationLong,
	}
}

// ErrorToast creates an error toast message
func ErrorToast(message string) ToastMsg {
	return ToastMsg{
		Type:     ToastError,
		Message:  message,
		Duration: ToastDurationLong,
	}
}
