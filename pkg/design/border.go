package design

import "github.com/charmbracelet/lipgloss"

// RoundedBox returns a lipgloss.Style with rounded borders and configurable
// width and padding. This is the standard box style for all kyanite.sh apps.
func RoundedBox(width int, padding int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Get("") /* uses default */ .Border).
		Padding(padding).
		Width(width)
}

// RoundedBoxThemed returns a RoundedBox using the given theme's border color.
func RoundedBoxThemed(t Theme, width int, padding int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Padding(padding).
		Width(width)
}

// Panel returns a styled panel with rounded border, background fill, and padding.
func Panel(t Theme, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Background(t.Panel).
		Foreground(t.Text).
		Padding(SpacingS).
		Width(width)
}
