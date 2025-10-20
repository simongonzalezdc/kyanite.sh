package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a color theme for the focus.sh application
type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Background lipgloss.Color
	Text       lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Border     lipgloss.Color
	Panel      lipgloss.Color
}
