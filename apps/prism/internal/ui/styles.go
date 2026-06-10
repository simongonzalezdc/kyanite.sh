package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
)

// Styles holds all UI styles derived from the design module.
type Styles struct {
	Theme  *design.Theme
	Tokens design.TokenSet
	Typo   design.TypographySet

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

// NewStyles creates styles based on theme using the design module's
// TokenSet, TypographySet, and Extensions for prism-specific styles.
func NewStyles(t *design.Theme) Styles {
	tokens := design.NewTokenSet(*t)
	typo := design.NewTypographySet(*t)

	// Extensions for prism-specific styles
	tokens = tokens.
		SetExt("Subtitle", lipgloss.Style{}.Foreground(t.Secondary).Italic(true)).
		SetExt("StatusBar", lipgloss.Style{}.Foreground(t.Text).Background(t.Background).Padding(design.SpacingNone, design.SpacingXS)).
		SetExt("Secondary", lipgloss.Style{}.Foreground(t.Secondary)).
		SetExt("Success", lipgloss.Style{}.Foreground(t.Success)).
		SetExt("Unselected", lipgloss.Style{}.Foreground(t.Text).Padding(design.SpacingNone, design.SpacingXS)).
		SetExt("Focused", lipgloss.Style{}.Border(lipgloss.RoundedBorder()).BorderForeground(t.Accent)).
		SetExt("ColorSwatch", lipgloss.Style{}.Padding(design.SpacingNone, design.SpacingS).Margin(design.SpacingNone, design.SpacingXS)).
		SetExt("HelpKey", lipgloss.Style{}.Foreground(t.Accent).Bold(true)).
		SetExt("HelpDesc", lipgloss.Style{}.Foreground(t.Text))

	return Styles{
		Theme:  t,
		Tokens: tokens,
		Typo:   typo,

		// Layout — standard tokens + additional config
		Border: tokens.Border.
			Border(lipgloss.RoundedBorder()).
			Padding(design.SpacingXS, design.SpacingS),
		Title:    typo.Title.Padding(design.SpacingNone, design.SpacingXS),
		Subtitle: tokens.Ext("Subtitle"),
		StatusBar: tokens.Ext("StatusBar"),

		// Text — standard tokens and extensions
		Primary:   tokens.Heading,
		Secondary: tokens.Ext("Secondary"),
		Accent:    tokens.Accent,
		Success:   tokens.Ext("Success"),
		Error:     tokens.Error,
		Muted:     tokens.Muted,

		// Interactive
		Selected: tokens.Selected.
			Padding(design.SpacingNone, design.SpacingXS),
		Unselected: tokens.Ext("Unselected"),
		Focused:    tokens.Ext("Focused"),

		// Components
		ColorSwatch: tokens.Ext("ColorSwatch"),
		HelpKey:     tokens.Ext("HelpKey"),
		HelpDesc:    tokens.Ext("HelpDesc"),
	}
}
