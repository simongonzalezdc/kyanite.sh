package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a color scheme used by the TUI.
type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Background lipgloss.Color
	Text       lipgloss.Color
	Success    lipgloss.Color
}

// GetStyle returns a lipgloss style using given foreground and background.
func (t Theme) GetStyle(fg, bg lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(fg).
		Background(bg)
}

// PrimaryStyle returns a style using the theme primary color as background and theme text as foreground.
func (t Theme) PrimaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Primary)
}

// SecondaryStyle returns a style using the theme secondary color as background and theme text as foreground.
func (t Theme) SecondaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Secondary)
}

// AccentStyle returns a style using the theme accent color as background and theme background as foreground.
func (t Theme) AccentStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Background).
		Background(t.Accent)
}

// SuccessStyle returns a style intended for success messages.
func (t Theme) SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Success).
		Background(t.Background).
		Bold(true)
}

// TextStyle returns a basic text style for the theme.
func (t Theme) TextStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Background)
}