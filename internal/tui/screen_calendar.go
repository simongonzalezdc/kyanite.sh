package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/styles"
	"strings"
)

func (m *MainModel) renderCalendarView() string {
	var b strings.Builder

	// CONSISTENT: Standard header styling
	header := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Bold(true).
		Align(lipgloss.Center).
		Padding(1, 3). // CONSISTENT: 1 vertical, 3 horizontal
		Render("📅 Calendar")
	b.WriteString(header)
	b.WriteString("\n\n")

	// CONSISTENT: Navigation styling
	navStyle := lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Background(styles.GetPanel()).
		Padding(1, 2).                    // CONSISTENT: 1 vertical, 2 horizontal
		Margin(1, 0, 1, 0).               // CONSISTENT: 1 top/bottom
		Border(lipgloss.RoundedBorder()). // CONSISTENT: RoundedBorder
		BorderForeground(styles.GetBorder())

	navControls := "←/h: Previous Month  •  →/l: Next Month  • Tab: Switch View"
	b.WriteString(navStyle.Render(navControls))
	b.WriteString("\n\n")

	// CONSISTENT: Calendar container styling
	if m.cal != nil && m.calRenderer != nil {
		m.cal.SelectedDate = m.calSelectedDate
		var calendarContent string
		if m.width > 60 { // Only render calendar if we have enough space
			calendarContent = m.calSelectedDate.Format("January 2006")
		} else {
			// Fallback for small screens
			calendarContent = m.calSelectedDate.Format("2006-01-02")
		}

		calendarBox := lipgloss.NewStyle().
			Foreground(styles.GetForeground()).
			Background(styles.GetPanel()).
			Padding(1, 2).                    // CONSISTENT: 1 vertical, 2 horizontal
			Border(lipgloss.RoundedBorder()). // CONSISTENT: RoundedBorder
			BorderForeground(styles.GetBorder()).
			Width(m.width - 6).
			Height(m.height - 15)

		b.WriteString(calendarBox.Render(calendarContent))
	} else {
		// CONSISTENT: Placeholder styling
		placeholder := lipgloss.NewStyle().
			Foreground(styles.GetForeground()).
			Background(styles.GetPanel()).
			Padding(1, 2).                    // CONSISTENT: 1 vertical, 2 horizontal
			Border(lipgloss.RoundedBorder()). // CONSISTENT: RoundedBorder
			BorderForeground(styles.GetBorder()).
			Align(lipgloss.Center).
			Render("Calendar loading...")
		b.WriteString(placeholder)
	}

	b.WriteString("\n\n")

	// CONSISTENT: Instructions styling
	instructions := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Italic(true).
		Align(lipgloss.Center).
		Render("Navigate months with arrow keys. Press Tab to switch views. Press Escape to return.")

	b.WriteString(instructions)

	return b.String()
}
