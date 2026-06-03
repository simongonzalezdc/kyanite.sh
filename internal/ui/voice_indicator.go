package ui

import (
	"fmt"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// VoiceIndicatorModel displays the current voice recording status
type VoiceIndicatorModel struct {
	state      app.VoiceState
	duration   time.Duration
	level      float32 // Audio level 0.0 to 1.0
	processing bool
	error      string
	width      int

	// Styles
	containerStyle  lipgloss.Style
	recordingStyle  lipgloss.Style
	processingStyle lipgloss.Style
	errorStyle      lipgloss.Style
	levelBarStyle   lipgloss.Style
	levelFillStyle  lipgloss.Style
}

// VoiceStateMsg updates the voice indicator state
type VoiceStateMsg struct {
	State    app.VoiceState
	Duration time.Duration
	Level    float32
	Error    string
}

// NewVoiceIndicatorModel creates a new voice indicator
func NewVoiceIndicatorModel() *VoiceIndicatorModel {
	t := theme.GetManager().Current()

	return &VoiceIndicatorModel{
		state: app.VoiceStateIdle,
		containerStyle: lipgloss.NewStyle().
			Padding(0, 1),
		recordingStyle: lipgloss.NewStyle().
			Foreground(t.Error). // Red for recording
			Bold(true),
		processingStyle: lipgloss.NewStyle().
			Foreground(t.Warning). // Yellow for processing
			Bold(true),
		errorStyle: lipgloss.NewStyle().
			Foreground(t.Error),
		levelBarStyle: lipgloss.NewStyle().
			Foreground(t.Secondary),
		levelFillStyle: lipgloss.NewStyle().
			Foreground(t.Primary),
	}
}

// Init initializes the voice indicator
func (m *VoiceIndicatorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *VoiceIndicatorModel) Update(msg tea.Msg) (*VoiceIndicatorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width

	case VoiceStateMsg:
		m.state = msg.State
		m.duration = msg.Duration
		m.level = msg.Level
		m.error = msg.Error
	}

	return m, nil
}

// View renders the voice indicator
func (m *VoiceIndicatorModel) View() string {
	if m.state == app.VoiceStateIdle && m.error == "" {
		return "" // Don't show when idle
	}

	var content string

	switch m.state {
	case app.VoiceStateRecording:
		// Recording indicator with duration and level meter
		durStr := formatDuration(m.duration)
		levelBar := m.renderLevelBar()
		content = m.recordingStyle.Render(fmt.Sprintf("[REC] %s %s", durStr, levelBar))

	case app.VoiceStateProcessing:
		content = m.processingStyle.Render("[...] Transcribing...")

	case app.VoiceStateError:
		if m.error != "" {
			content = m.errorStyle.Render(fmt.Sprintf("[X] %s", m.error))
		}

	default:
		return ""
	}

	return m.containerStyle.Render(content)
}

// renderLevelBar creates a visual audio level meter
func (m *VoiceIndicatorModel) renderLevelBar() string {
	const barWidth = 10
	filled := int(m.level * float32(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}

	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "#"
		} else {
			bar += "-"
		}
	}

	return fmt.Sprintf("[%s]", bar)
}

// SetState updates the indicator state
func (m *VoiceIndicatorModel) SetState(state app.VoiceState) {
	m.state = state
}

// SetDuration updates the recording duration
func (m *VoiceIndicatorModel) SetDuration(d time.Duration) {
	m.duration = d
}

// SetLevel updates the audio level
func (m *VoiceIndicatorModel) SetLevel(level float32) {
	m.level = level
}

// SetError sets an error message
func (m *VoiceIndicatorModel) SetError(err string) {
	m.error = err
	if err != "" {
		m.state = app.VoiceStateError
	}
}

// ClearError clears the error message
func (m *VoiceIndicatorModel) ClearError() {
	m.error = ""
	if m.state == app.VoiceStateError {
		m.state = app.VoiceStateIdle
	}
}

// IsVisible returns whether the indicator should be displayed
func (m *VoiceIndicatorModel) IsVisible() bool {
	return m.state != app.VoiceStateIdle || m.error != ""
}

// formatDuration formats a duration as M:SS or MM:SS
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	s := (d % time.Minute) / time.Second
	return fmt.Sprintf("%d:%02d", m, s)
}

// VoiceStatusBar provides a compact status bar component for voice status
type VoiceStatusBar struct {
	indicator *VoiceIndicatorModel
	available bool
	shortcut  string
}

// NewVoiceStatusBar creates a status bar voice component
func NewVoiceStatusBar(shortcut string) *VoiceStatusBar {
	return &VoiceStatusBar{
		indicator: NewVoiceIndicatorModel(),
		shortcut:  shortcut,
	}
}

// SetAvailable sets whether voice is available
func (sb *VoiceStatusBar) SetAvailable(available bool) {
	sb.available = available
}

// Update handles messages
func (sb *VoiceStatusBar) Update(msg tea.Msg) tea.Cmd {
	_, cmd := sb.indicator.Update(msg)
	return cmd
}

// View renders the status bar component
func (sb *VoiceStatusBar) View() string {
	t := theme.GetManager().Current()

	// If recording or processing, show indicator
	if sb.indicator.IsVisible() {
		return sb.indicator.View()
	}

	// Show availability hint
	if sb.available {
		style := lipgloss.NewStyle().Foreground(t.Secondary)
		return style.Render(fmt.Sprintf("[MIC] %s: Dictate", sb.shortcut))
	}

	return ""
}

// GetIndicator returns the underlying indicator model
func (sb *VoiceStatusBar) GetIndicator() *VoiceIndicatorModel {
	return sb.indicator
}
