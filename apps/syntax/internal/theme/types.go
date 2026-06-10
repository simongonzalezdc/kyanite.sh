package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a color theme for the application
type Theme struct {
	Name       string
	Primary    string
	Secondary  string
	Accent     string
	Background string
	Text       string
	Success    string
}

// Styles contains all Lipgloss styles for a theme
type Styles struct {
	Base          lipgloss.Style
	Title         lipgloss.Style
	Heading       lipgloss.Style
	Text          lipgloss.Style
	Accent        lipgloss.Style
	Success       lipgloss.Style
	Error         lipgloss.Style
	Border        lipgloss.Style
	StatusBar     lipgloss.Style
	EditorPane    lipgloss.Style
	PreviewPane   lipgloss.Style
	MenuSelected  lipgloss.Style
	MenuUnselected lipgloss.Style
}

// ApplyTheme creates Lipgloss styles from a theme
func (t *Theme) ApplyTheme() Styles {
	return Styles{
		Base: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Text)).
			Background(lipgloss.Color(t.Background)),

		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Accent)).
			Background(lipgloss.Color(t.Background)).
			Bold(true).
			Padding(0, 1),

		Heading: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Primary)).
			Bold(true),

		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Text)),

		Accent: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Accent)),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Success)),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")),

		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(t.Secondary)),

		StatusBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Background)).
			Background(lipgloss.Color(t.Primary)).
			Padding(0, 1),

		EditorPane: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Text)).
			Background(lipgloss.Color(t.Background)).
			Padding(1),

		PreviewPane: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Text)).
			Background(lipgloss.Color(t.Background)).
			Padding(1),

		MenuSelected: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Background)).
			Background(lipgloss.Color(t.Accent)).
			Padding(0, 1),

		MenuUnselected: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Text)).
			Padding(0, 1),
	}
}
