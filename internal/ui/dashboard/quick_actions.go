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
	hovered   int // Track mouse hover for visual feedback
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
		actions:  actions,
		selected: 0,
		hovered:  -1, // No hover initially
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

	case tea.MouseMsg:
		// Handle mouse events for quick actions
		cmd := m.handleMouse(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// handleMouse processes mouse events and returns any resulting command
func (m *QuickActionsModel) handleMouse(msg tea.MouseMsg) tea.Cmd {
	if m.width == 0 {
		return nil
	}

	// Calculate grid layout
	cols := 2
	if m.width < 40 {
		cols = 1
	}

	// Calculate card dimensions (must match renderActionCard)
	cardWidth := (m.width / cols) - 6
	cardHeight := 7 // 5 content + 2 border

	// Account for panel padding and title
	titleHeight := 3 // "Quick Actions" + empty line

	// Calculate which action is at the mouse position
	// X position maps to column
	col := msg.X / (cardWidth + 2) // +2 for margin
	if col >= cols {
		col = cols - 1
	}
	if col < 0 {
		col = 0
	}

	// Y position maps to row (accounting for title)
	row := (msg.Y - titleHeight) / cardHeight
	if row < 0 {
		row = 0
	}

	// Calculate action index
	actionIdx := row*cols + col
	if actionIdx >= len(m.actions) {
		actionIdx = len(m.actions) - 1
	}
	if actionIdx < 0 {
		actionIdx = 0
	}

	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionRelease {
			// Click - select and activate
			m.selected = actionIdx
			return m.activateAction(actionIdx)
		} else if msg.Action == tea.MouseActionPress {
			// Press - just select (visual feedback)
			m.selected = actionIdx
		}

	case tea.MouseButtonWheelUp:
		// Scroll up - select previous action
		if m.selected > 0 {
			m.selected--
		}

	case tea.MouseButtonWheelDown:
		// Scroll down - select next action
		if m.selected < len(m.actions)-1 {
			m.selected++
		}

	case tea.MouseButtonNone:
		// Mouse motion - update hover state
		if msg.Action == tea.MouseActionMotion {
			m.hovered = actionIdx
		}
	}

	return nil
}

// activateAction triggers the action associated with the given index
func (m *QuickActionsModel) activateAction(idx int) tea.Cmd {
	if idx < 0 || idx >= len(m.actions) {
		return nil
	}

	action := m.actions[idx]
	switch action.ID {
	case "new":
		return func() tea.Msg { return ScreenChangeMsg{Screen: 2} } // screenEditor
	case "open":
		return func() tea.Msg { return ScreenChangeMsg{Screen: 5} } // screenManager
	case "ai":
		return func() tea.Msg { return ScreenChangeMsg{Screen: 2} } // screenEditor
	case "export":
		return func() tea.Msg { return ScreenChangeMsg{Screen: 3} } // screenExport
	case "theory":
		return func() tea.Msg { return ScreenChangeMsg{Screen: 4} } // screenTheory
	case "audio":
		return func() tea.Msg { return ScreenChangeMsg{Screen: 5} } // screenAudio (using manager slot)
	}
	return nil
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

	// Calculate grid layout based on panel width
	cols := 2
	if m.width < 40 {
		cols = 1
	}

	var rows []string
	for i := 0; i < len(m.actions); i += cols {
		var row []string
		for j := i; j < i+cols && j < len(m.actions); j++ {
			action := m.actions[j]
			isSelected := j == m.selected
			isHovered := j == m.hovered && j != m.selected
			card := m.renderActionCard(action, isSelected, isHovered)
			row = append(row, card)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, row...))
	}

	title := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Render("Quick Actions")

	// Add mouse hint
	mouseHint := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Faint(true).
		Render("Click to select")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		mouseHint,
		lipgloss.JoinVertical(lipgloss.Left, rows...),
	)

	// Use allocated panel dimensions directly
	return lipgloss.NewStyle().
		Width(m.width - 2).
		MaxWidth(m.width - 2).
		MaxHeight(m.height - 2).
		Padding(0, 1).
		Render(content)
}

func (m *QuickActionsModel) renderActionCard(action Action, selected, hovered bool) string {
	t := theme.GetManager().Current()

	// Calculate card width based on panel width
	// Use 2 columns if panel is wide enough, otherwise 1
	cols := 2
	if m.width < 40 {
		cols = 1
	}
	cardWidth := (m.width / cols) - 6

	// Base card style
	baseStyle := lipgloss.NewStyle().
		Width(cardWidth).
		Height(5).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Margin(0, 1)

	if selected {
		// Selected state - most prominent
		baseStyle = baseStyle.
			BorderForeground(t.Primary).
			Background(lipgloss.Color(string(t.Primary))).
			Foreground(t.Background)
	} else if hovered {
		// Hovered state - subtle highlight
		baseStyle = baseStyle.
			BorderForeground(t.Accent).
			Background(t.Background).
			Foreground(t.Text)
	} else {
		// Normal state
		baseStyle = baseStyle.
			BorderForeground(t.Secondary).
			Background(t.Background).
			Foreground(t.Text)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(action.Icon+" "+action.Title),
		"",
		lipgloss.NewStyle().Faint(true).Render(action.Description),
		"",
		lipgloss.NewStyle().Align(lipgloss.Right).Render("["+action.Shortcut+"]"),
	)

	return baseStyle.Render(content)
}