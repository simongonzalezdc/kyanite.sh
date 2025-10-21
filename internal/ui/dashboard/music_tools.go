package dashboard

import (
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MusicToolsModel manages the music tools panel
type MusicToolsModel struct {
	width     int
	height    int
}

// NewMusicToolsModel creates a new music tools model
func NewMusicToolsModel() *MusicToolsModel {
	return &MusicToolsModel{
		
	}
}

// Update handles messages for the music tools panel
func (m *MusicToolsModel) Update(msg tea.Msg) tea.Cmd {
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

// View renders the music tools panel
func (m *MusicToolsModel) View() string {
	if m.width == 0 {
		return "Music Tools"
	}
	
	t := theme.GetManager().Current()
	
	title := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Render("Music Tools")
	
	// Music tools items
	items := []string{
		"🎼 Chord Progressions",
		"🎸 Scale Reference",
		"📖 Rhyme Dictionary",
		"🎵 Song Templates",
	}
	
	var itemViews []string
	for _, item := range items {
		itemView := lipgloss.NewStyle().
			Foreground(t.Text).
			Render(item)
		itemViews = append(itemViews, itemView)
	}
	
	// Open theory tools button
	openButton := lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Secondary).
		Padding(0, 1).
		MarginTop(1).
		Render("[Open Theory Tools]")
	
	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		lipgloss.JoinVertical(lipgloss.Left, itemViews...),
		"",
		openButton,
	)
	
	return content
}