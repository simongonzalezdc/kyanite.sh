package editor

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// bpmTapperModel handles the BPM tapper UI component
type bpmTapperModel struct {
	width   int
	height  int
	visible bool

	// Tap timing
	tapTimes   []time.Time
	currentBPM int
	lastTap    time.Time

	// Visual feedback
	tapHistory []bool
	maxHistory int

	// Callback to set BPM in pattern
	setBMPCallback func(bpm int)

	// Styles
	containerStyle   lipgloss.Style
	headerStyle      lipgloss.Style
	bpmStyle         lipgloss.Style
	tapStyle         lipgloss.Style
	instructionStyle lipgloss.Style
}

// NewBPMTapperModel creates a new BPM tapper model
func NewBPMTapperModel() *bpmTapperModel {
	t := theme.GetManager().Current()
	return &bpmTapperModel{
		visible:    false,
		currentBPM: 120, // Default 120 BPM
		maxHistory: 8,
		tapHistory: make([]bool, 8),

		containerStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent).
			Background(t.Background).
			Padding(1, 2),

		headerStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent).
			Align(lipgloss.Center).
			MarginBottom(1),

		bpmStyle: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1),

		tapStyle: lipgloss.NewStyle().
			Foreground(t.Success).
			Align(lipgloss.Center).
			MarginBottom(1),

		instructionStyle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Align(lipgloss.Center).
			MarginTop(1),
	}
}

// Init initializes the BPM tapper model
func (m *bpmTapperModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the BPM tapper
func (m *bpmTapperModel) Update(msg tea.Msg) (*bpmTapperModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ShowBPMTapperMsg:
		m.visible = true
		m.setBMPCallback = msg.SetBMPCallback
		m.reset()

		return m, nil

	case HideBPMTapperMsg:
		m.visible = false
		return m, nil

	case tea.KeyMsg:
		if !m.visible {
			return m, nil
		}

		switch msg.String() {
		case "esc":
			m.visible = false
			return m, nil

		case " ", "enter":
			// Record a tap
			m.recordTap()

		case "c":
			// Confirm and set BPM
			if m.setBMPCallback != nil {
				m.setBMPCallback(m.currentBPM)
			}
			m.visible = false
			return m, nil

		case "r":
			// Reset taps
			m.reset()
		}
	}

	return m, nil
}

// View renders the BPM tapper
func (m *bpmTapperModel) View() string {
	if !m.visible {
		return ""
	}

	// Calculate dimensions
	maxWidth := m.width - 10
	if maxWidth < 50 {
		maxWidth = 50
	}

	// Build content
	var content strings.Builder

	// Header
	header := m.headerStyle.Render("[~] BPM Tapper")
	content.WriteString(header)
	content.WriteString("\n\n")

	// BPM display
	bpmText := fmt.Sprintf("%d BPM", m.currentBPM)
	bpmDisplay := m.bpmStyle.Render(bpmText)
	content.WriteString(bpmDisplay)
	content.WriteString("\n\n")

	// Tap history
	tapHistory := m.renderTapHistory()
	content.WriteString(tapHistory)
	content.WriteString("\n\n")

	// Instructions
	instructions := m.renderInstructions()
	content.WriteString(instructions)

	// Apply container style
	return m.containerStyle.Width(maxWidth).Render(content.String())
}

// renderTapHistory renders the visual tap history
func (m *bpmTapperModel) renderTapHistory() string {
	var symbols []string

	for i, tapped := range m.tapHistory {
		if tapped {
			// Show more recent taps as brighter
			if i >= len(m.tapHistory)-3 {
				symbols = append(symbols, "*")
			} else {
				symbols = append(symbols, "o")
			}
		} else {
			symbols = append(symbols, "o")
		}
	}

	tapHistory := strings.Join(symbols, " ")
	return m.tapStyle.Render(tapHistory)
}

// renderInstructions renders the instruction text
func (m *bpmTapperModel) renderInstructions() string {
	instructions := []string{
		"[Space/Enter] Tap tempo",
		"[C] Confirm BPM",
		"[R] Reset",
		"[Esc] Cancel",
	}

	return m.instructionStyle.Render(strings.Join(instructions, " | "))
}

// recordTap records a tap and calculates BPM
func (m *bpmTapperModel) recordTap() {
	now := time.Now()

	// Add to tap times
	m.tapTimes = append(m.tapTimes, now)

	// Keep only the last 8 taps
	if len(m.tapTimes) > 8 {
		m.tapTimes = m.tapTimes[len(m.tapTimes)-8:]
	}

	// Update tap history
	m.tapHistory = append(m.tapHistory[1:], true)

	// Calculate BPM if we have at least 2 taps
	if len(m.tapTimes) >= 2 {
		m.calculateBPM()
	}

	m.lastTap = now
}

// calculateBPM calculates BPM from tap intervals
func (m *bpmTapperModel) calculateBPM() {
	if len(m.tapTimes) < 2 {
		return
	}

	// Use last 4-8 taps for calculation
	recentTaps := m.tapTimes
	if len(recentTaps) > 8 {
		recentTaps = m.tapTimes[len(m.tapTimes)-8:]
	}

	// Calculate average interval
	var totalInterval time.Duration
	for i := 1; i < len(recentTaps); i++ {
		interval := recentTaps[i].Sub(recentTaps[i-1])
		totalInterval += interval
	}

	avgInterval := totalInterval / time.Duration(len(recentTaps)-1)
	if avgInterval <= 0 {
		return
	}

	// Convert to BPM (60 seconds / average interval in seconds)
	m.currentBPM = int(60.0 / avgInterval.Seconds())

	// Clamp to reasonable range
	if m.currentBPM < 60 {
		m.currentBPM = 60
	} else if m.currentBPM > 200 {
		m.currentBPM = 200
	}
}

// reset resets the tap state
func (m *bpmTapperModel) reset() {
	m.tapTimes = nil
	m.currentBPM = 120
	m.lastTap = time.Time{}
	m.tapHistory = make([]bool, m.maxHistory)
}

// Show shows the BPM tapper with the given callback
func (m *bpmTapperModel) Show(callback func(bpm int)) tea.Cmd {
	return func() tea.Msg {
		return ShowBPMTapperMsg{SetBMPCallback: callback}
	}
}

// Hide hides the BPM tapper
func (m *bpmTapperModel) Hide() tea.Cmd {
	return func() tea.Msg {
		return HideBPMTapperMsg{}
	}
}

// IsVisible returns whether the BPM tapper is currently visible
func (m *bpmTapperModel) IsVisible() bool {
	return m.visible
}

// GetCurrentBPM returns the current BPM
func (m *bpmTapperModel) GetCurrentBPM() int {
	return m.currentBPM
}

// Message types for BPM tapper
type ShowBPMTapperMsg struct {
	SetBMPCallback func(bpm int)
}

type HideBPMTapperMsg struct{}
