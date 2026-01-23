package dashboard

import (
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AIAssistantModel manages the AI assistant panel
type AIAssistantModel struct {
	width     int
	height    int
	
}
// NewAIAssistantModel creates a new AI assistant model
func NewAIAssistantModel() *AIAssistantModel {
	return &AIAssistantModel{
		
	}
}

// Update handles messages for the AI assistant panel
func (m *AIAssistantModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// CRITICAL: Return a command to trigger re-render
		return nil
		
	}
	
	// CRITICAL: Always return commands
	return tea.Batch(cmds...)
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
	
	// AI status
	status := lipgloss.NewStyle().
		Foreground(t.Success).
		Render("🤖 AI Status: Ready")
	
	// Quick suggestions
	suggestionsTitle := lipgloss.NewStyle().
		Foreground(t.Text).
		Bold(true).
		Render("💡 Quick Suggestions:")
	
	suggestions := []string{
		"• \"Summer memories...\"",
		"• \"Dancing through...\"",
		"• \"Where the wind...\"",
	}
	
	var suggestionViews []string
	for _, suggestion := range suggestions {
		suggestionView := lipgloss.NewStyle().
			Foreground(t.Text).
			Faint(true).
			Render(suggestion)
		suggestionViews = append(suggestionViews, suggestionView)
	}
	
	// Start brainstorm button
	startButton := lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Secondary).
		Padding(0, 1).
		MarginTop(1).
		Render("[Start Brainstorm]")
	
	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		status,
		"",
		suggestionsTitle,
		lipgloss.JoinVertical(lipgloss.Left, suggestionViews...),
		"",
		startButton,
	)

	// Use allocated panel dimensions directly
	return lipgloss.NewStyle().
		Width(m.width - 2).
		MaxWidth(m.width - 2).
		MaxHeight(m.height - 2).
		Padding(0, 1).
		Render(content)
}