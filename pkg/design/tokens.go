package design

import "github.com/charmbracelet/lipgloss"

// Spacing constants define the linear 0-4 character cell scale.
// This is the terminal equivalent of the 4px grid used in web design systems.
const (
	SpacingNone = 0 // No spacing
	SpacingXS   = 1 // Tight: inline badges, icon gaps
	SpacingS    = 2 // Standard: content padding, list items
	SpacingM    = 3 // Comfortable: panel sections
	SpacingL    = 4 // Generous: between major sections
)

// TokenSet provides 10 standard semantic style tokens plus app-specific extensions.
// Standard tokens cover the common visual roles across all kyanite.sh apps.
// Apps that need additional styles (e.g., EditorPane, PreviewPane) use the
// Extensions map with their own keys.
type TokenSet struct {
	// Standard tokens (always present)
	Base     lipgloss.Style // Default content style (Text on Background)
	Title    lipgloss.Style // Page title (bold + accent fg)
	Heading  lipgloss.Style // Section heading (bold + primary fg)
	Body     lipgloss.Style // Body text (text fg on background)
	Muted    lipgloss.Style // Deemphasized text (muted fg, faint)
	Accent   lipgloss.Style // Accent highlight (accent fg)
	Border   lipgloss.Style // Border elements (border color)
	Selected lipgloss.Style // Selected item (inverted primary/bg)
	Active   lipgloss.Style // Active/focused element (accent bg + text fg)
	Error    lipgloss.Style // Error state (error fg, bold)

	// Extensions holds app-specific styles that don't map to the 10 standard tokens.
	// Keys are convention-based strings like "StatusBar", "EditorPane", etc.
	Extensions map[string]lipgloss.Style
}

// NewTokenSet creates a TokenSet from a Theme, deriving all 10 standard tokens
// and initializing an empty Extensions map.
func NewTokenSet(t Theme) TokenSet {
	return TokenSet{
		Base: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.Background),
		Title: lipgloss.NewStyle().
			Foreground(t.Accent).
			Background(t.Background).
			Bold(true),
		Heading: lipgloss.NewStyle().
			Foreground(t.Primary).
			Background(t.Background).
			Bold(true),
		Body: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.Background),
		Muted: lipgloss.NewStyle().
			Foreground(t.Muted).
			Background(t.Background).
			Faint(true),
		Accent: lipgloss.NewStyle().
			Foreground(t.Accent).
			Background(t.Background),
		Border: lipgloss.NewStyle().
			Foreground(t.Border).
			BorderForeground(t.Border),
		Selected: lipgloss.NewStyle().
			Foreground(t.Background).
			Background(t.Primary).
			Bold(true),
		Active: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.Accent).
			Bold(true),
		Error: lipgloss.NewStyle().
			Foreground(t.Error).
			Background(t.Background).
			Bold(true),
		Extensions: make(map[string]lipgloss.Style),
	}
}

// Ext returns a style from the Extensions map, or a zero style if not found.
func (ts TokenSet) Ext(key string) lipgloss.Style {
	if ts.Extensions == nil {
		return lipgloss.Style{}
	}
	return ts.Extensions[key]
}

// SetExt sets an extension style and returns the updated TokenSet.
func (ts TokenSet) SetExt(key string, style lipgloss.Style) TokenSet {
	if ts.Extensions == nil {
		ts.Extensions = make(map[string]lipgloss.Style)
	}
	ts.Extensions[key] = style
	return ts
}
