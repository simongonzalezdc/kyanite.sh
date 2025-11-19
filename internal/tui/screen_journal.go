package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/styles"
	"strings"
)

func (m *MainModel) renderJournalView() string {
	var b strings.Builder

	// CONSISTENT: Standard header styling
	header := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Background(styles.GetPanel()).
		Bold(true).
		Align(lipgloss.Center).
		Padding(1, 3).                         // CONSISTENT: 1 vertical, 3 horizontal
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT: RoundedBorder
		BorderForeground(styles.GetBorder()).
		Render("📝 Journal")
	b.WriteString(header)
	b.WriteString("\n\n")

	// CONSISTENT: Journal content styling
	journalBox := lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Background(styles.GetPanel()).
		Padding(1, 2).                    // CONSISTENT: 1 vertical, 2 horizontal
		Margin(1, 0, 1, 0).               // CONSISTENT: 1 top/bottom
		Border(lipgloss.RoundedBorder()). // CONSISTENT: RoundedBorder
		BorderForeground(styles.GetBorder()).
		Width(m.width - 6).
		Height(m.height - 10)

	journalContent := lipgloss.NewStyle().
		Render(`Journal entries will appear here.

Available commands:
• focus journal new     - Create new entry
• focus journal list     - List all entries  
• focus journal view     - View specific entry
• focus journal search   - Search entries

Press Tab to navigate back to dashboard.`)

	b.WriteString(journalBox.Render(journalContent))
	b.WriteString("\n\n")

	// CONSISTENT: Instructions styling
	instructions := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Italic(true).
		Align(lipgloss.Center).
		Render("Press Tab to return to dashboard, or use focus journal commands in terminal.")

	b.WriteString(instructions)

	return b.String()
}
