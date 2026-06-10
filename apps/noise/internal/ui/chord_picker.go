package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/kyanite/noise/internal/data"
	"github.com/kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChordPickerModel handles the chord picker UI component
type ChordPickerModel struct {
	width        int
	height       int
	visible      bool
	progressions []data.ChordProgression
	filteredProg []data.ChordProgression
	selectedIdx  int
	activeMood   string
	showAll      bool
	loaded       bool

	// Callback to insert chords into editor
	insertCallback func(chords []string)

	// Animation for visual feedback
	animationTime time.Time
	pulseProgress float64

	// Styles
	containerStyle   lipgloss.Style
	headerStyle      lipgloss.Style
	selectedStyle    lipgloss.Style
	normalStyle      lipgloss.Style
	moodStyle        lipgloss.Style
	descriptionStyle lipgloss.Style
	instructionStyle lipgloss.Style
}

// ChordPickerMsg represents messages for the chord picker
type ChordPickerMsg struct {
	Type string
	Data interface{}
}

// ShowChordPickerMsg shows the chord picker
type ShowChordPickerMsg struct {
	InsertCallback func(chords []string)
}

// HideChordPickerMsg hides the chord picker
type HideChordPickerMsg struct{}

// NewChordPickerModel creates a new chord picker model
func NewChordPickerModel() *ChordPickerModel {
	t := theme.GetManager().Current()

	return &ChordPickerModel{
		visible:       false,
		selectedIdx:   0,
		activeMood:    "all",
		showAll:       true,
		loaded:        false,
		animationTime: time.Now(),

		containerStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary).
			Background(t.Background).
			Padding(1, 2),

		headerStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Primary).
			Align(lipgloss.Center).
			MarginBottom(1),

		selectedStyle: lipgloss.NewStyle().
			Background(t.Primary).
			Foreground(t.Background).
			Bold(true).
			Padding(0, 1),

		normalStyle: lipgloss.NewStyle().
			Foreground(t.Text).
			Padding(0, 1),

		moodStyle: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true).
			Padding(0, 1),

		descriptionStyle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Italic(true).
			MarginLeft(2),

		instructionStyle: lipgloss.NewStyle().
			Foreground(t.Text).
			Align(lipgloss.Center).
			MarginTop(1),
	}
}

// Init initializes the chord picker model
func (m *ChordPickerModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the chord picker
func (m *ChordPickerModel) Update(msg tea.Msg) (*ChordPickerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ShowChordPickerMsg:
		m.visible = true
		m.insertCallback = msg.InsertCallback
		m.selectedIdx = 0
		m.activeMood = "all"
		m.showAll = true
		m.animationTime = time.Now()

		// Load progressions if not already loaded
		if !m.loaded {
			if progressions, err := data.GetAllChordProgressions(); err == nil {
				m.progressions = progressions
				m.filteredProg = progressions
				m.loaded = true
			}
		}

		return m, nil

	case HideChordPickerMsg:
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

		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}

		case "down", "j":
			if m.selectedIdx < len(m.filteredProg)-1 {
				m.selectedIdx++
			}

		case "1", "2", "3", "4":
			// Mood selection
			moods := []string{"all", "happy", "sad", "tense", "chill"}
			idx := 0
			switch msg.String() {
			case "1":
				idx = 0
			case "2":
				idx = 1
			case "3":
				idx = 2
			case "4":
				idx = 3
			}

			if idx < len(moods) {
				m.setMoodFilter(moods[idx])
			}

		case "alt+r":
			// Random selection
			if len(m.filteredProg) > 0 {
				m.selectedIdx = rand.Intn(len(m.filteredProg))
				m.animationTime = time.Now()
			}

		case "enter", " ":
			// Select and insert chords
			if m.selectedIdx >= 0 && m.selectedIdx < len(m.filteredProg) {
				selected := m.filteredProg[m.selectedIdx]
				if m.insertCallback != nil {
					m.insertCallback(selected.Chords)
				}
				m.visible = false
			}
		}

	case AnimationTickMsg:
		// Update pulse animation
		if time.Since(m.animationTime) < time.Second {
			m.pulseProgress = float64(time.Since(m.animationTime)) / float64(time.Second)
		} else {
			m.pulseProgress = 0
		}
	}

	return m, nil
}

// View renders the chord picker
func (m *ChordPickerModel) View() string {
	if !m.visible || !m.loaded {
		return ""
	}

	// Calculate dimensions
	maxWidth := m.width - 10
	if maxWidth < 60 {
		maxWidth = 60
	}

	maxHeight := m.height - 10
	if maxHeight < 20 {
		maxHeight = 20
	}

	// Build content
	var content strings.Builder

	// Header
	header := m.headerStyle.Render("[~] Quick Chord Picker")
	content.WriteString(header)
	content.WriteString("\n\n")

	// Mood filter
	content.WriteString(m.renderMoodFilter())
	content.WriteString("\n\n")

	// Progressions list
	content.WriteString(m.renderProgressionsList(maxWidth, maxHeight-8))

	// Instructions
	content.WriteString("\n\n")
	content.WriteString(m.renderInstructions())

	// Apply container style
	return m.containerStyle.Width(maxWidth).Render(content.String())
}

// renderMoodFilter renders the mood filter buttons
func (m *ChordPickerModel) renderMoodFilter() string {
	moods := []string{"all", "happy", "sad", "tense", "chill"}
	var buttons []string

	for i, mood := range moods {
		var style lipgloss.Style
		if mood == m.activeMood {
			t := theme.GetManager().Current()
			style = m.moodStyle.Background(t.Primary).Foreground(t.Background)
		} else {
			style = m.moodStyle
		}

		button := style.Render(fmt.Sprintf("[%d] %s", i+1, strings.Title(mood)))
		buttons = append(buttons, button)
	}

	return strings.Join(buttons, " ")
}

// renderProgressionsList renders the list of chord progressions
func (m *ChordPickerModel) renderProgressionsList(maxWidth, maxHeight int) string {
	if len(m.filteredProg) == 0 {
		return "No progressions found for this mood."
	}

	var lines []string

	// Calculate visible items
	visibleItems := maxHeight
	if visibleItems > len(m.filteredProg) {
		visibleItems = len(m.filteredProg)
	}

	// Calculate scroll position
	startIdx := m.selectedIdx
	if startIdx+visibleItems > len(m.filteredProg) {
		startIdx = len(m.filteredProg) - visibleItems
	}
	if startIdx < 0 {
		startIdx = 0
	}

	// Render visible items
	for i := startIdx; i < startIdx+visibleItems && i < len(m.filteredProg); i++ {
		progression := m.filteredProg[i]

		var style lipgloss.Style
		if i == m.selectedIdx {
			style = m.selectedStyle
		} else {
			style = m.normalStyle
		}

		// Format: [index] Name - Chords
		chordStr := strings.Join(progression.Chords, " - ")
		lineText := fmt.Sprintf("%d. %s - %s", i+1, progression.Name, chordStr)

		// Truncate if too long
		if len(lineText) > maxWidth-4 {
			lineText = lineText[:maxWidth-7] + "..."
		}

		line := style.Render(lineText)
		lines = append(lines, line)

		// Add description for selected item
		if i == m.selectedIdx {
			desc := m.descriptionStyle.Render(progression.Description)
			lines = append(lines, desc)
		}
	}

	return strings.Join(lines, "\n")
}

// renderInstructions renders the instruction text
func (m *ChordPickerModel) renderInstructions() string {
	instructions := []string{
		"[1-5] Filter by mood",
		"[Alt+R] Random selection",
		"[UpDown] Navigate",
		"[Enter/Space] Insert chords",
		"[Esc] Cancel",
	}

	return m.instructionStyle.Render(strings.Join(instructions, " | "))
}

// setMoodFilter sets the mood filter and updates the filtered progressions
func (m *ChordPickerModel) setMoodFilter(mood string) {
	m.activeMood = mood
	m.showAll = (mood == "all")
	m.selectedIdx = 0

	if m.showAll {
		m.filteredProg = m.progressions
	} else {
		var filtered []data.ChordProgression
		for _, prog := range m.progressions {
			if prog.Mood == mood {
				filtered = append(filtered, prog)
			}
		}
		m.filteredProg = filtered
	}
}

// Show shows the chord picker with the given callback
func (m *ChordPickerModel) Show(callback func(chords []string)) tea.Cmd {
	return func() tea.Msg {
		return ShowChordPickerMsg{InsertCallback: callback}
	}
}

// Hide hides the chord picker
func (m *ChordPickerModel) Hide() tea.Cmd {
	return func() tea.Msg {
		return HideChordPickerMsg{}
	}
}

// IsVisible returns whether the chord picker is currently visible
func (m *ChordPickerModel) IsVisible() bool {
	return m.visible
}

// GetSelectedProgression returns the currently selected progression
func (m *ChordPickerModel) GetSelectedProgression() *data.ChordProgression {
	if m.selectedIdx >= 0 && m.selectedIdx < len(m.filteredProg) {
		return &m.filteredProg[m.selectedIdx]
	}
	return nil
}
