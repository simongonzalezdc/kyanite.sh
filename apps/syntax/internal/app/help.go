package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewHelp() string {
	var b strings.Builder

	title := m.Styles.Title.Render("syntax.sh - Help")
	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, title))
	b.WriteString("\n\n")

	sections := []struct {
		title string
		items []string
	}{
		{
			title: "Global Shortcuts",
			items: []string{
				"Ctrl+Q         - Quit application",
				"Ctrl+Shift+T   - Cycle through themes",
				"?              - Show this help",
				"Esc            - Go back/cancel",
			},
		},
		{
			title: "Navigation",
			items: []string{
				"↑/↓ or j/k     - Navigate lists",
				"Enter          - Select/Open",
				"n              - Create new item",
			},
		},
		{
			title: "Text Editor",
			items: []string{
				"i              - Enter INSERT mode",
				"Esc            - Return to NORMAL mode",
				"Ctrl+S         - Save scene",
				"Ctrl+Z         - Undo",
				"Ctrl+Y         - Redo",
				"↑/↓/←/→        - Move cursor (INSERT mode)",
			},
		},
		{
			title: "Project Dashboard",
			items: []string{
				"c              - View characters",
				"s              - View scenes",
				"l              - View locations",
				"t              - View statistics",
				"e              - Export project",
				"h              - Help",
			},
		},
	}

	for _, section := range sections {
		b.WriteString(m.Styles.Heading.Render(section.title))
		b.WriteString("\n")
		for _, item := range section.items {
			b.WriteString(m.Styles.Text.Render("  " + item))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center,
		m.Styles.Text.Render("Press Esc to close help")))

	return b.String()
}

func (m Model) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "?":
		m.CurrentScreen = m.PreviousScreen
		return m, nil
	}
	return m, nil
}
