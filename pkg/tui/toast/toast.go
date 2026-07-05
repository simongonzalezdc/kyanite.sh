// Package toast provides a Bubble Tea toast notification component
// shared across all kyanite.sh apps.
//
// Toasts are short-lived overlay messages (info, success, warning, error)
// that auto-dismiss after a configurable duration. They stack vertically
// and cap at a maximum number of visible items.
//
// Usage:
//
//	toasts := toast.New(designManager)
//	// In your Update():
//	case tea.WindowSizeMsg:
//	    toasts.SetWidth(msg.Width)
//	case someEvent:
//	    return toasts.Add(toast.Success("Saved!"))
//
//	// In your View():
//	    return toasts.View() + mainContent
package toast

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/design"
	"github.com/kyanite/design/icons"
)

// Duration constants for common toast display times.
const (
	DurationDefault = 3 * time.Second
	DurationLong    = 5 * time.Second
	DurationShort   = 2 * time.Second
)

// Type represents the toast severity level.
type Type int

const (
	Info Type = iota
	Success
	Warning
	Error
)

// Msg is a Bubble Tea message that triggers a toast.
type Msg struct {
	Type     Type
	Message  string
	Duration time.Duration
}

// dismissMsg is an internal message to auto-remove a toast.
type dismissMsg struct{ id int }

// Item represents a single visible toast.
type Item struct {
	ID       int
	Type     Type
	Message  string
	Duration time.Duration
	ShowTime time.Time
}

// Model manages a stack of toast notifications.
type Model struct {
	items     []Item
	nextID    int
	width     int
	maxVisible int
	dm        *design.Manager

	infoStyle    lipgloss.Style
	successStyle lipgloss.Style
	warningStyle lipgloss.Style
	errorStyle   lipgloss.Style
}

// New creates a toast model using the given design manager for theming.
func New(dm *design.Manager) *Model {
	m := &Model{
		items:      []Item{},
		width:      80,
		maxVisible: 3,
		dm:         dm,
	}
	m.applyTheme()
	return m
}

// Init satisfies tea.Model (no-op).
func (m *Model) Init() tea.Cmd { return nil }

// Update handles Bubble Tea messages.
func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case Msg:
		return m.add(msg)
	case dismissMsg:
		m.remove(msg.id)
	}
	return nil
}

// View renders the stacked toast notifications.
func (m *Model) View() string {
	if len(m.items) == 0 {
		return ""
	}

	visible := m.items
	if len(visible) > m.maxVisible {
		visible = visible[len(visible)-m.maxVisible:]
	}

	var rendered []string
	for _, item := range visible {
		rendered = append(rendered, m.renderItem(item))
	}

	container := lipgloss.NewStyle().Padding(0, 1)
	return container.Render(lipgloss.JoinVertical(lipgloss.Left, rendered...))
}

// SetWidth updates the render width.
func (m *Model) SetWidth(w int) { m.width = w }

// HasToasts reports whether any toasts are visible.
func (m *Model) HasToasts() bool { return len(m.items) > 0 }

// ClearAll removes all toasts.
func (m *Model) ClearAll() { m.items = []Item{} }

// RefreshTheme re-applies styles from the current theme (call after theme change).
func (m *Model) RefreshTheme() { m.applyTheme() }

// Add returns a tea.Cmd that adds a toast and schedules its dismissal.
// This is the primary way to show a toast from an Update method.
func (m *Model) Add(msg Msg) tea.Cmd {
	return func() tea.Msg { return msg }
}

// --- Constructors for common toast types ---

// InfoMsg creates an informational toast.
func InfoMsg(message string) Msg {
	return Msg{Type: Info, Message: message, Duration: DurationDefault}
}

// SuccessMsg creates a success toast.
func SuccessMsg(message string) Msg {
	return Msg{Type: Success, Message: message, Duration: DurationDefault}
}

// WarningMsg creates a warning toast.
func WarningMsg(message string) Msg {
	return Msg{Type: Warning, Message: message, Duration: DurationLong}
}

// ErrorMsg creates an error toast.
func ErrorMsg(message string) Msg {
	return Msg{Type: Error, Message: message, Duration: DurationLong}
}

// --- Private ---

func (m *Model) add(msg Msg) tea.Cmd {
	dur := msg.Duration
	if dur == 0 {
		dur = DurationDefault
	}

	item := Item{
		ID:       m.nextID,
		Type:     msg.Type,
		Message:  msg.Message,
		Duration: dur,
		ShowTime: time.Now(),
	}
	m.nextID++
	m.items = append(m.items, item)

	id := item.ID
	return tea.Tick(dur, func(time.Time) tea.Msg { return dismissMsg{id: id} })
}

func (m *Model) remove(id int) {
	for i, t := range m.items {
		if t.ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return
		}
	}
}

// renderPadding is the horizontal padding subtracted from width for toast content.
const renderPadding = 4

func (m *Model) renderItem(item Item) string {
	icon := iconFor(item.Type)
	content := icon + " " + item.Message

	maxW := m.width - renderPadding
	if len(content) > maxW && maxW > 10 {
		content = content[:maxW-3] + "..."
	}

	return styleFor(m, item.Type).Width(m.width - renderPadding).Render(content)
}

func iconFor(t Type) string {
	switch t {
	case Success:
		return icons.GetIcon("success")
	case Warning:
		return icons.GetIcon("warning")
	case Error:
		return icons.GetIcon("error")
	default:
		return icons.GetIcon("info")
	}
}

func styleFor(m *Model, t Type) lipgloss.Style {
	switch t {
	case Success:
		return m.successStyle
	case Warning:
		return m.warningStyle
	case Error:
		return m.errorStyle
	default:
		return m.infoStyle
	}
}

func (m *Model) applyTheme() {
	if m.dm == nil {
		return
	}
	t := m.dm.Current()

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
