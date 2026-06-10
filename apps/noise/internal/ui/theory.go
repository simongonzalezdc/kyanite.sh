package ui

import (
	"fmt"
	"strings"

	"github.com/kyanite/noise/internal/app"
	"github.com/kyanite/noise/internal/theme"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TheoryModel handles the music theory tools screen
type TheoryModel struct {
	// Layout
	width  int
	height int

	// Tabs for different theory tools
	activeTab int
	tabNames  []string

	// Input fields
	scaleInput       textinput.Model
	chordInput       textinput.Model
	progressionInput textinput.Model
	analysisInput    textinput.Model
	rhymeInput       textinput.Model

	// Viewport for scrollable content
	viewport viewport.Model

	// Help system
	help help.Model
	keys theoryKeyMap

	// Theory service
	theoryService *app.TheoryService

	// Focus management
	focusedElement FocusedElement
}

// FocusedElement represents which UI element has focus
type FocusedElement int

const (
	ElementScaleInput FocusedElement = iota
	ElementChordInput
	ElementProgressionInput
	ElementAnalysisInput
	ElementRhymeInput
	ElementViewport
)

// Tab indices
const (
	TabScales int = iota
	TabChords
	TabProgressions
	TabAnalysis
	TabRhymes
)

// NewTheoryModel creates a new theory model
func NewTheoryModel() *TheoryModel {
	// Initialize theory service
	theoryService := app.NewTheoryService()

	// Initialize inputs
	scaleInput := textinput.New()
	scaleInput.Placeholder = "Enter key and scale type (e.g., C major)"
	scaleInput.Focus()

	chordInput := textinput.New()
	chordInput.Placeholder = "Enter chord (e.g., Cmaj7)"

	progressionInput := textinput.New()
	progressionInput.Placeholder = "Enter key and pattern (e.g., C I-V-vi-IV)"

	analysisInput := textinput.New()
	analysisInput.Placeholder = "Enter text to analyze for chords"

	rhymeInput := textinput.New()
	rhymeInput.Placeholder = "Enter word to find rhymes"

	// Initialize viewport
	vp := viewport.New(0, 0)

	// Tab names
	tabNames := []string{"Scales", "Chords", "Progressions", "Analysis", "Rhymes"}

	m := &TheoryModel{
		activeTab:        0,
		tabNames:         tabNames,
		scaleInput:       scaleInput,
		chordInput:       chordInput,
		progressionInput: progressionInput,
		analysisInput:    analysisInput,
		rhymeInput:       rhymeInput,
		viewport:         vp,
		help:             help.New(),
		keys:             theoryKeys,
		theoryService:    theoryService,
		focusedElement:   ElementScaleInput,
	}

	return m
}

// Init initializes the theory model
func (m *TheoryModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the theory tools
func (m *TheoryModel) Update(msg tea.Msg) (*TheoryModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()

	case tea.MouseMsg:
		// Handle mouse events for theory tabs
		switch msg.Button {
		case tea.MouseButtonLeft:
			if msg.Action == tea.MouseActionRelease {
				// Check if click is in tab bar area (typically Y <= 3)
				if msg.Y <= 3 {
					// Calculate which tab was clicked based on X position
					// Tabs are roughly evenly spaced
					if len(m.tabNames) > 0 {
						tabWidth := m.width / len(m.tabNames)
						if tabWidth > 0 {
							clickedTab := msg.X / tabWidth
							if clickedTab >= 0 && clickedTab < len(m.tabNames) {
								m.activeTab = clickedTab
							}
						}
					}
				}
			}

		case tea.MouseButtonWheelUp:
			// Scroll viewport up
			m.viewport.LineUp(3)

		case tea.MouseButtonWheelDown:
			// Scroll viewport down
			m.viewport.LineDown(3)
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Tab):
			m.nextTab()
			return m, nil
		case key.Matches(msg, m.keys.ShiftTab):
			m.prevTab()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			return m, m.handleEnter()
		case key.Matches(msg, m.keys.Escape):
			return m, nil // Return to menu
		case key.Matches(msg, m.keys.Help):
			// Toggle help (for now just return)
			return m, nil
		}
	}

	// Update inputs based on focus
	switch m.focusedElement {
	case ElementScaleInput:
		var cmd tea.Cmd
		m.scaleInput, cmd = m.scaleInput.Update(msg)
		cmds = append(cmds, cmd)
	case ElementChordInput:
		var cmd tea.Cmd
		m.chordInput, cmd = m.chordInput.Update(msg)
		cmds = append(cmds, cmd)
	case ElementProgressionInput:
		var cmd tea.Cmd
		m.progressionInput, cmd = m.progressionInput.Update(msg)
		cmds = append(cmds, cmd)
	case ElementAnalysisInput:
		var cmd tea.Cmd
		m.analysisInput, cmd = m.analysisInput.Update(msg)
		cmds = append(cmds, cmd)
	case ElementRhymeInput:
		var cmd tea.Cmd
		m.rhymeInput, cmd = m.rhymeInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update viewport
	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}

// View renders the theory tools screen
func (m *TheoryModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading theory tools..."
	}

	// Header
	header := m.renderHeader()

	// Tabs
	tabs := m.renderTabs()

	// Content area
	content := m.renderContent()

	// Footer with help
	footer := m.renderFooter()

	// Combine everything
	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		tabs,
		content,
		footer,
	)
}

// renderHeader renders the header section
func (m *TheoryModel) renderHeader() string {
	t := theme.GetManager().Current()
	title := titleGradient("[~] Music Theory Tools", t)
	title = lipgloss.Style{}.
		Width(m.width).
		Align(lipgloss.Center).
		Padding(1, 0).
		Render(title)

	return title
}

// renderTabs renders the tab navigation
func (m *TheoryModel) renderTabs() string {
	var tabItems []string

	t := theme.GetManager().Current()
	for i, name := range m.tabNames {
		var style lipgloss.Style
		if i == m.activeTab {
			style = lipgloss.Style{}.
				Background(t.Primary).
				Foreground(t.Background).
				Bold(true).
				Padding(0, 2)
		} else {
			style = lipgloss.Style{}.
				Foreground(t.Secondary).
				Background(t.Background).
				Padding(0, 2)
		}

		tab := style.Render(fmt.Sprintf(" %s ", name))
		tabItems = append(tabItems, tab)
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top, tabItems...)
	tabs = lipgloss.Style{}.
		Width(m.width).
		Align(lipgloss.Center).
		Padding(0, 2).
		Render(tabs)

	return tabs
}

// renderContent renders the main content area
func (m *TheoryModel) renderContent() string {
	content := m.getCurrentTabContent()

	// Set viewport content
	m.viewport.SetContent(content)

	// Calculate viewport dimensions
	contentHeight := m.height - 8 // Account for header, tabs, and footer
	if contentHeight < 5 {
		contentHeight = 5
	}

	m.viewport.Width = m.width - 4
	m.viewport.Height = contentHeight

	return lipgloss.Style{}.
		Padding(1, 2).
		Render(m.viewport.View())
}

// getCurrentTabContent returns content for the active tab
func (m *TheoryModel) getCurrentTabContent() string {
	switch m.activeTab {
	case TabScales:
		return m.renderScalesTab()
	case TabChords:
		return m.renderChordsTab()
	case TabProgressions:
		return m.renderProgressionsTab()
	case TabAnalysis:
		return m.renderAnalysisTab()
	case TabRhymes:
		return m.renderRhymesTab()
	default:
		return "Unknown tab"
	}
}

// renderScalesTab renders the scales tab content
func (m *TheoryModel) renderScalesTab() string {
	var content strings.Builder

	t := theme.GetManager().Current()
	// Input section
	inputStyle := lipgloss.Style{}.Foreground(t.Primary).Bold(true)
	content.WriteString(inputStyle.Render("Scale Lookup:\n"))
	content.WriteString(m.scaleInput.View())
	content.WriteString("\n\n")

	// Results section
	if m.scaleInput.Value() != "" {
		parts := strings.Fields(m.scaleInput.Value())
		if len(parts) >= 2 {
			key := parts[0]
			scaleType := parts[1]

			scaleInfo, err := m.theoryService.GetScaleInfo(key, scaleType)
			if err != nil {
				content.WriteString(fmt.Sprintf("Error: %v\n", err))
			} else {
				content.WriteString(fmt.Sprintf("Scale: %s\n", scaleInfo.Name))
				content.WriteString(fmt.Sprintf("Notes: %s\n", strings.Join(scaleInfo.Notes, ", ")))
				content.WriteString(fmt.Sprintf("Intervals: %s\n", strings.Join(scaleInfo.Pattern, ", ")))
				content.WriteString(fmt.Sprintf("Description: %s\n\n", scaleInfo.Description))
			}
		}
	}

	// Common scales section
	content.WriteString(inputStyle.Render("Common Scales:\n"))
	commonScales := m.theoryService.GetCommonScales()
	for _, scale := range commonScales {
		content.WriteString(fmt.Sprintf("- %s: %s\n", scale.Name, strings.Join(scale.Notes, ", ")))
	}

	return content.String()
}

// renderChordsTab renders the chords tab content
func (m *TheoryModel) renderChordsTab() string {
	var content strings.Builder

	t := theme.GetManager().Current()
	// Input section
	inputStyle := lipgloss.Style{}.Foreground(t.Primary).Bold(true)
	content.WriteString(inputStyle.Render("Chord Lookup:\n"))
	content.WriteString(m.chordInput.View())
	content.WriteString("\n\n")

	// Results section
	if m.chordInput.Value() != "" {
		parts := strings.Fields(m.chordInput.Value())
		if len(parts) >= 2 {
			root := parts[0]
			chordType := parts[1]

			chordInfo, err := m.theoryService.GetChordInfo(root, chordType)
			if err != nil {
				content.WriteString(fmt.Sprintf("Error: %v\n", err))
			} else {
				content.WriteString(fmt.Sprintf("Chord: %s %s\n", chordInfo.Root, chordInfo.Quality))
				content.WriteString(fmt.Sprintf("Notes: %s\n", strings.Join(chordInfo.Notes, ", ")))
				content.WriteString(fmt.Sprintf("Intervals: %s\n", strings.Join(chordInfo.Intervals, ", ")))
				content.WriteString(fmt.Sprintf("Description: %s\n", chordInfo.Description))
			}
		}
	}

	return content.String()
}

// renderProgressionsTab renders the progressions tab content
func (m *TheoryModel) renderProgressionsTab() string {
	var content strings.Builder

	t := theme.GetManager().Current()
	// Input section
	inputStyle := lipgloss.Style{}.Foreground(t.Primary).Bold(true)
	content.WriteString(inputStyle.Render("Chord Progression:\n"))
	content.WriteString(m.progressionInput.View())
	content.WriteString("\n\n")

	// Results section
	if m.progressionInput.Value() != "" {
		parts := strings.Fields(m.progressionInput.Value())
		if len(parts) >= 2 {
			key := parts[0]
			pattern := parts[1]

			progressionInfo, err := m.theoryService.GetProgressionInfo(key, pattern)
			if err != nil {
				content.WriteString(fmt.Sprintf("Error: %v\n", err))
			} else {
				content.WriteString(fmt.Sprintf("Progression: %s\n", progressionInfo.Name))
				content.WriteString(fmt.Sprintf("Chords: %s\n", strings.Join(progressionInfo.Chords, " - ")))
				content.WriteString(fmt.Sprintf("Description: %s\n", progressionInfo.Description))
			}
		}
	}

	return content.String()
}

// renderAnalysisTab renders the chord analysis tab content
func (m *TheoryModel) renderAnalysisTab() string {
	var content strings.Builder

	t := theme.GetManager().Current()
	// Input section
	inputStyle := lipgloss.Style{}.Foreground(t.Primary).Bold(true)
	content.WriteString(inputStyle.Render("Chord Analysis:\n"))
	content.WriteString(m.analysisInput.View())
	content.WriteString("\n\n")

	// Results section
	if m.analysisInput.Value() != "" {
		analysis, err := m.theoryService.AnalyzeChords(m.analysisInput.Value())
		if err != nil {
			content.WriteString(fmt.Sprintf("Error: %v\n", err))
		} else {
			content.WriteString(fmt.Sprintf("Found %d chords:\n", len(analysis.DetectedChords)))
			for _, chord := range analysis.DetectedChords {
				content.WriteString(fmt.Sprintf("- %s %s: %s\n", chord.Root, chord.Quality, strings.Join(chord.Notes, ", ")))
			}

			if len(analysis.Suggestions) > 0 {
				content.WriteString("\nSuggestions:\n")
				for _, suggestion := range analysis.Suggestions {
					content.WriteString(fmt.Sprintf("- %s\n", suggestion))
				}
			}
		}
	}

	return content.String()
}

// renderRhymesTab renders the rhymes tab content
func (m *TheoryModel) renderRhymesTab() string {
	var content strings.Builder

	t := theme.GetManager().Current()
	// Input section
	inputStyle := lipgloss.Style{}.Foreground(t.Primary).Bold(true)
	content.WriteString(inputStyle.Render("Rhyme Finder:\n"))
	content.WriteString(m.rhymeInput.View())
	content.WriteString("\n\n")

	// Results section
	if m.rhymeInput.Value() != "" {
		rhymes, err := m.theoryService.FindRhymes(m.rhymeInput.Value())
		if err != nil {
			content.WriteString(fmt.Sprintf("Error: %v\n", err))
		} else if len(rhymes) > 0 {
			content.WriteString(fmt.Sprintf("Rhymes for \"%s\":\n", m.rhymeInput.Value()))
			content.WriteString(strings.Join(rhymes, ", "))
		} else {
			content.WriteString(fmt.Sprintf("No rhymes found for \"%s\"", m.rhymeInput.Value()))
		}
	}

	return content.String()
}

// renderFooter renders the footer with help information
func (m *TheoryModel) renderFooter() string {
	helpText := m.help.ShortHelpView(m.keys.ShortHelp())
	t := theme.GetManager().Current()
	helpText = lipgloss.Style{}.
		Foreground(t.Secondary).
		Padding(0, 2).
		Render(helpText)

	return helpText
}

// updateLayout updates the layout when dimensions change
func (m *TheoryModel) updateLayout() {
	// Update viewport dimensions with responsive considerations
	contentHeight := m.height - 8
	if contentHeight < 5 {
		contentHeight = 5
	}

	// Ensure minimum width for readability
	minWidth := 60
	if m.width < minWidth {
		// For very narrow terminals, reduce padding
		m.viewport.Width = m.width - 2
	} else {
		m.viewport.Width = m.width - 4
	}

	m.viewport.Height = contentHeight
}

// nextTab switches to the next tab
func (m *TheoryModel) nextTab() {
	m.activeTab = (m.activeTab + 1) % len(m.tabNames)
	m.updateFocusForTab()
}

// prevTab switches to the previous tab
func (m *TheoryModel) prevTab() {
	m.activeTab = (m.activeTab - 1 + len(m.tabNames)) % len(m.tabNames)
	m.updateFocusForTab()
}

// updateFocusForTab updates focus based on the active tab
func (m *TheoryModel) updateFocusForTab() {
	switch m.activeTab {
	case TabScales:
		m.focusedElement = ElementScaleInput
		m.scaleInput.Focus()
		m.chordInput.Blur()
		m.progressionInput.Blur()
		m.analysisInput.Blur()
		m.rhymeInput.Blur()
	case TabChords:
		m.focusedElement = ElementChordInput
		m.chordInput.Focus()
		m.scaleInput.Blur()
		m.progressionInput.Blur()
		m.analysisInput.Blur()
		m.rhymeInput.Blur()
	case TabProgressions:
		m.focusedElement = ElementProgressionInput
		m.progressionInput.Focus()
		m.scaleInput.Blur()
		m.chordInput.Blur()
		m.analysisInput.Blur()
		m.rhymeInput.Blur()
	case TabAnalysis:
		m.focusedElement = ElementAnalysisInput
		m.analysisInput.Focus()
		m.scaleInput.Blur()
		m.chordInput.Blur()
		m.progressionInput.Blur()
		m.rhymeInput.Blur()
	case TabRhymes:
		m.focusedElement = ElementRhymeInput
		m.rhymeInput.Focus()
		m.scaleInput.Blur()
		m.chordInput.Blur()
		m.progressionInput.Blur()
		m.analysisInput.Blur()
	}
}

// handleEnter processes enter key presses
func (m *TheoryModel) handleEnter() tea.Cmd {
	// Refresh the current tab content
	return nil
}

// Key map for theory tools
type theoryKeyMap struct {
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Help     key.Binding
}

var theoryKeys = theoryKeyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next tab"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev tab"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Help: key.NewBinding(
		key.WithKeys("f1", "?"),
		key.WithHelp("f1", "help"),
	),
}

func (k theoryKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.ShiftTab, k.Enter, k.Escape, k.Help}
}

func (k theoryKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.ShiftTab},
		{k.Enter, k.Escape},
		{k.Help},
	}
}
