package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/prism/internal/theme"
)

// Styles holds all UI styles
type Styles struct {
	Theme *theme.Theme

	// Layout
	Border    lipgloss.Style
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	StatusBar lipgloss.Style

	// Text
	Primary   lipgloss.Style
	Secondary lipgloss.Style
	Accent    lipgloss.Style
	Success   lipgloss.Style
	Error     lipgloss.Style
	Muted     lipgloss.Style

	// Interactive
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	Focused    lipgloss.Style

	// Components
	ColorSwatch lipgloss.Style
	HelpKey     lipgloss.Style
	HelpDesc    lipgloss.Style
}

// NewStyles creates styles based on theme
func NewStyles(t *theme.Theme) Styles {
	s := Styles{Theme: t}

	// Layout styles
	s.Border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.Primary)).
		Padding(1, 2)

	s.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Primary)).
		Bold(true).
		Padding(0, 1)

	s.Subtitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Secondary)).
		Italic(true)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Text)).
		Background(lipgloss.Color(t.Background)).
		Padding(0, 1)

	// Text styles
	s.Primary = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Primary))

	s.Secondary = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Secondary))

	s.Accent = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Accent))

	s.Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Success))

	s.Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	s.Muted = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Text)).
		Faint(true)

	// Interactive styles
	s.Selected = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Background)).
		Background(lipgloss.Color(t.Primary)).
		Bold(true).
		Padding(0, 1)

	s.Unselected = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Text)).
		Padding(0, 1)

	s.Focused = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(t.Accent))

	// Component styles
	s.ColorSwatch = lipgloss.NewStyle().
		Padding(0, 2).
		Margin(0, 1)

	s.HelpKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Accent)).
		Bold(true)

	s.HelpDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Text))

	return s
}
