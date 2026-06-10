package design

import "github.com/charmbracelet/lipgloss"

// TypographySet provides 4 text hierarchy levels derived from a Theme.
// These cover the standard text roles across all kyanite.sh apps.
type TypographySet struct {
	Title   lipgloss.Style // Bold + accent color (page titles)
	Heading lipgloss.Style // Bold + primary color (section headings)
	Body    lipgloss.Style // Normal text color (body text)
	Muted   lipgloss.Style // Faint + dim (deemphasized text)
}

// NewTypographySet creates a TypographySet from a Theme.
func NewTypographySet(t Theme) TypographySet {
	return TypographySet{
		Title: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true),
		Heading: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),
		Body: lipgloss.NewStyle().
			Foreground(t.Text),
		Muted: lipgloss.NewStyle().
			Foreground(t.Muted).
			Faint(true),
	}
}
