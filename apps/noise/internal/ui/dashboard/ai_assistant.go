// Package dashboard provides the dashboard UI components.
package dashboard

import (
	"github.com/kyanite/noise/internal/app"
	"github.com/kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TriggerBrainstormMsg tells the editor to start AI brainstorm
type TriggerBrainstormMsg struct {
	Theme string // Optional theme for brainstorm
}

// AIAssistantModel manages the AI assistant panel
type AIAssistantModel struct {
	width     int
	height    int
	aiService *app.AIService // Reference to AI service
	focused   bool           // Track focus state
	selected  int            // Track selected item (0=brainstorm button)
	hovered   int            // Track mouse hover state
}

// NewAIAssistantModel creates a new AI assistant model
func NewAIAssistantModel() *AIAssistantModel {
	return &AIAssistantModel{
		selected: 0,
		hovered:  -1,
	}
}

// SetAIService sets the AI service reference
func (m *AIAssistantModel) SetAIService(aiService *app.AIService) {
	m.aiService = aiService
}

// Update handles messages for the AI assistant panel
func (m *AIAssistantModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil

	case tea.MouseMsg:
		// Handle mouse events for the AI assistant panel
		return m.handleMouse(msg)

	case tea.KeyMsg:
		// Handle keyboard events when focused
		if m.focused {
			switch msg.String() {
			case "enter":
				return m.triggerBrainstorm()
			case "up", "k":
				// No navigation needed with single button
				return nil
			case "down", "j":
				// No navigation needed with single button
				return nil
			}
		}
	}

	return nil
}

// handleMouse processes mouse events and returns any resulting command
func (m *AIAssistantModel) handleMouse(msg tea.MouseMsg) tea.Cmd {
	if m.width == 0 {
		return nil
	}

	// Calculate button position (button is near the bottom of the panel)
	// Title (1) + empty (1) + status (1) + empty (1) + suggestions title (1) +
	// 3 suggestions (3) + empty (1) + button area starts around line 9
	buttonStartY := 9
	buttonEndY := 11

	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionRelease {
			// Check if click is in the button area
			if msg.Y >= buttonStartY && msg.Y <= buttonEndY {
				return m.triggerBrainstorm()
			}
		} else if msg.Action == tea.MouseActionPress {
			// Visual feedback on press
			if msg.Y >= buttonStartY && msg.Y <= buttonEndY {
				m.selected = 0
			}
		}

	case tea.MouseButtonWheelUp, tea.MouseButtonWheelDown:
		// Could scroll through suggestions in the future
		return nil

	case tea.MouseButtonNone:
		// Mouse motion - update hover state
		if msg.Action == tea.MouseActionMotion {
			if msg.Y >= buttonStartY && msg.Y <= buttonEndY {
				m.hovered = 0
			} else {
				m.hovered = -1
			}
		}
	}

	return nil
}

// triggerBrainstorm sends a message to start the AI brainstorm
func (m *AIAssistantModel) triggerBrainstorm() tea.Cmd {
	return func() tea.Msg {
		return TriggerBrainstormMsg{}
	}
}

// SetFocused sets the focus state of the panel
func (m *AIAssistantModel) SetFocused(focused bool) {
	m.focused = focused
}

// GetAIStatus returns the current AI service status
func (m *AIAssistantModel) GetAIStatus() string {
	if m.aiService == nil {
		return "Not configured"
	}
	if m.aiService.IsAvailable() {
		return "Ready"
	}
	return "Not available"
}

// View renders the AI assistant panel
func (m *AIAssistantModel) View() string {
	if m.width == 0 {
		return "AI Assistant"
	}

	t := theme.GetManager().Current()

	title := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Render("AI Assistant")

	// AI status - show actual status from service
	statusText := m.GetAIStatus()
	statusColor := t.Success
	if statusText != "Ready" {
		statusColor = t.Warning
	}
	status := lipgloss.NewStyle().
		Foreground(statusColor).
		Render("[AI] Status: " + statusText)

	// Quick suggestions
	suggestionsTitle := lipgloss.NewStyle().
		Foreground(t.Text).
		Bold(true).
		Render("[TIP] Quick Suggestions:")

	suggestions := []string{
		"- \"Summer memories...\"",
		"- \"Dancing through...\"",
		"- \"Where the wind...\"",
	}

	var suggestionViews []string
	for _, suggestion := range suggestions {
		suggestionView := lipgloss.NewStyle().
			Foreground(t.Text).
			Faint(true).
			Render(suggestion)
		suggestionViews = append(suggestionViews, suggestionView)
	}

	// Start brainstorm button with hover/selected states
	buttonStyle := lipgloss.NewStyle().
		Padding(0, 1).
		MarginTop(1)

	if m.hovered == 0 || m.selected == 0 {
		// Hover/selected state - more prominent
		buttonStyle = buttonStyle.
			Foreground(t.Background).
			Background(t.Primary).
			Bold(true)
	} else {
		// Normal state
		buttonStyle = buttonStyle.
			Foreground(t.Text).
			Background(t.Secondary)
	}

	startButton := buttonStyle.Render("[Start Brainstorm]")

	// Mouse hint
	mouseHint := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Faint(true).
		Render("Click button or press Enter")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		status,
		"",
		suggestionsTitle,
		lipgloss.JoinVertical(lipgloss.Left, suggestionViews...),
		"",
		startButton,
		mouseHint,
	)

	// Use allocated panel dimensions directly
	return lipgloss.NewStyle().
		Width(m.width-2).
		MaxWidth(m.width-2).
		MaxHeight(m.height-2).
		Padding(0, 1).
		Render(content)
}
