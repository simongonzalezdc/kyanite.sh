package dashboard

import (
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RecentWorkModel manages the recent work panel
type RecentWorkModel struct {
	width     int
	height    int
}

// NewRecentWorkModel creates a new recent work model
func NewRecentWorkModel() *RecentWorkModel {
	return &RecentWorkModel{
		
	}
}

// Update handles messages for the recent work panel
func (m *RecentWorkModel) Update(msg tea.Msg) tea.Cmd {
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

// View renders the recent work panel
func (m *RecentWorkModel) View() string {
	if m.width == 0 {
		return "Recent Work"
	}
	
	t := theme.GetManager().Current()
	
	title := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Render("Recent Work")
	
	// Sample recent work items
	items := []string{
		"[P] Rock Opera - Last edited: 2h ago",
		"[~] Midnight Dreams - Last edited: 1d ago",
		"[~] Summer Breeze - Last edited: 3d ago",
	}
	
	var itemViews []string
	for _, item := range items {
		itemView := lipgloss.NewStyle().
			Foreground(t.Text).
			Render(item)
		itemViews = append(itemViews, itemView)
	}
	
	// Resume all projects button
	resumeButton := lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Secondary).
		Padding(0, 1).
		MarginTop(1).
		Render("[Resume All Projects]")
	
	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		lipgloss.JoinVertical(lipgloss.Left, itemViews...),
		"",
		resumeButton,
	)

	// Use allocated panel dimensions directly
	return lipgloss.NewStyle().
		Width(m.width - 2).
		MaxWidth(m.width - 2).
		MaxHeight(m.height - 2).
		Padding(0, 1).
		Render(content)
}