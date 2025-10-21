package dashboard

import (
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Action represents a quick action
type Action struct {
	ID          string
	Title       string
	Description string
	Icon        string
	Shortcut    string
	Animated    bool
}

// QuickActionsModel manages the quick actions grid
type QuickActionsModel struct {
	width     int
	height    int
	actions   []Action
	selected  int
}

// NewQuickActionsModel creates a new quick actions model
func NewQuickActionsModel() *QuickActionsModel {
	actions := []Action{
		{ID: "new", Title: "New Song", Description: "Create a new song", Icon: "🎵", Shortcut: "1", Animated: true},
		{ID: "open", Title: "Open Project", Description: "Browse existing projects", Icon: "📁", Shortcut: "2", Animated: true},
		{ID: "ai", Title: "AI Brainstorm", Description: "Get AI assistance", Icon: "🤖", Shortcut: "3", Animated: true},
		{ID: "export", Title: "Export", Description: "Export your work", Icon: "📤", Shortcut: "4", Animated: true},
		{ID: "theory", Title: "Theory Tools", Description: "Music theory reference", Icon: "🎼", Shortcut: "5", Animated: true},
		{ID: "audio", Title: "Audio Tools", Description: "Metronome & playback", Icon: "🎧", Shortcut: "6", Animated: true},
	}
	
	return &QuickActionsModel{
		actions:   actions,
		selected:  0,
		
	}
}

// Update handles messages for the quick actions
func (m *QuickActionsModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.selected = 0
		case "2":
			m.selected = 1
		case "3":
			m.selected = 2
		case "4":
			m.selected = 3
		case "5":
			m.selected = 4
		case "6":
			m.selected = 5
		}
	}
	
	return tea.Batch(cmds...)
}

// HandleKey handles key presses for quick actions
func (m *QuickActionsModel) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "1":
		return func() tea.Msg {
			return ScreenChangeMsg{Screen: 2} // screenEditor
		}
	case "2":
		return func() tea.Msg {
			return ScreenChangeMsg{Screen: 5} // screenManager
		}
	case "3":
		return func() tea.Msg {
			return ScreenChangeMsg{Screen: 2} // screenEditor
		}
	case "4":
		return func() tea.Msg {
			return ScreenChangeMsg{Screen: 3} // screenExport
		}
	case "5":
		return func() tea.Msg {
			return ScreenChangeMsg{Screen: 4} // screenTheory
		}
	case "6":
		return func() tea.Msg {
			return ScreenChangeMsg{Screen: 5} // screenAudio
		}
	}
	return nil
}

// View renders the quick actions grid
func (m *QuickActionsModel) View() string {
	if m.width == 0 {
		return "Quick Actions"
	}

	t := theme.GetManager().Current()

	// Calculate grid layout
	cols := 2
	if m.width > 100 {
		cols = 3
	}

	var rows []string
	for i := 0; i < len(m.actions); i += cols {
		var row []string
		for j := i; j < i+cols && j < len(m.actions); j++ {
			action := m.actions[j]
			card := m.renderActionCard(action, j == m.selected)
			row = append(row, card)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, row...))
	}

	title := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Render("Quick Actions")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		lipgloss.JoinVertical(lipgloss.Left, rows...),
	)
}

func (m *QuickActionsModel) renderActionCard(action Action, selected bool) string {
	t := theme.GetManager().Current()
	
	// Base card style
	baseStyle := lipgloss.NewStyle().
		Width(m.width/3-4).
		Height(6).
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Margin(1)
	
	if selected {
		baseStyle = baseStyle.
			BorderForeground(t.Primary).
			Background(lipgloss.Color(string(t.Primary))).
			Foreground(t.Background)
	} else {
		baseStyle = baseStyle.
			BorderForeground(t.Secondary).
			Background(t.Background).
			Foreground(t.Text)
	}
	
	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(action.Icon + " " + action.Title),
		"",
		lipgloss.NewStyle().Faint(true).Render(action.Description),
		"",
		lipgloss.NewStyle().Align(lipgloss.Right).Render("[" + action.Shortcut + "]"),
	)
	
	return baseStyle.Render(content)
}