package dashboard

import (
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SystemInfoModel manages the system info panel
type SystemInfoModel struct {
	width     int
	height    int
}

// NewSystemInfoModel creates a new system info model
func NewSystemInfoModel() *SystemInfoModel {
	return &SystemInfoModel{
		
	}
}

// Update handles messages for the system info panel
func (m *SystemInfoModel) Update(msg tea.Msg) tea.Cmd {
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

// View renders the system info panel
func (m *SystemInfoModel) View() string {
	if m.width == 0 {
		return "System Info"
	}
	
	t := theme.GetManager().Current()
	
	title := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Render("System Info")
	
	// Statistics
	statsTitle := lipgloss.NewStyle().
		Foreground(t.Text).
		Bold(true).
		Render("[STATS] Statistics:")
	
	stats := []string{
		"- Songs: 24",
		"- Projects: 5",
		"- Hours: 47",
	}
	
	var statViews []string
	for _, stat := range stats {
		statView := lipgloss.NewStyle().
			Foreground(t.Text).
			Render(stat)
		statViews = append(statViews, statView)
	}
	
	// Storage info
	storage := lipgloss.NewStyle().
		Foreground(t.Text).
		Render("[DISK] Storage: 2.3 MB")

	// Performance info
	performance := lipgloss.NewStyle().
		Foreground(t.Success).
		Render("[PERF] Performance: Good")
	
	// Open settings button
	settingsButton := lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Secondary).
		Padding(0, 1).
		MarginTop(1).
		Render("[Open Settings]")
	
	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		statsTitle,
		lipgloss.JoinVertical(lipgloss.Left, statViews...),
		"",
		storage,
		performance,
		"",
		settingsButton,
	)

	// Use allocated panel dimensions directly
	return lipgloss.NewStyle().
		Width(m.width - 2).
		MaxWidth(m.width - 2).
		MaxHeight(m.height - 2).
		Padding(0, 1).
		Render(content)
}