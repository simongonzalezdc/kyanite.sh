package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
)

// Styles contains all Lipgloss styles for a theme.
// It wraps the design module's TokenSet and maps syntax-specific styles
// to the appropriate standard tokens or extensions.
type Styles struct {
	// Standard tokens mapped directly from TokenSet
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

// applyTheme creates Styles from a design.Theme by deriving all styles
// from the design module's TokenSet standard tokens.
func applyTheme(t design.Theme) Styles {
	ts := design.NewTokenSet(t)

	return Styles{
		Base:     ts.Base,
		Title:    ts.Title,
		Heading:  ts.Heading,
		Text:     ts.Body,
		Accent:   ts.Accent,
		Success:  ts.Active.Foreground(t.Success).Background(t.Background).Bold(false),
		Error:    ts.Error,
		Border:   ts.Border,
		StatusBar:     ts.Selected.Bold(false).Padding(0, 1),
		EditorPane:    ts.Body.Padding(1),
		PreviewPane:   ts.Body.Padding(1),
		MenuSelected:  ts.Active.Padding(0, 1),
		MenuUnselected: ts.Body.Padding(0, 1),
	}
}

// Manager wraps the shared design.Manager with syntax-specific theme methods.
type Manager struct {
	inner *design.Manager
}

// NewManager creates a new theme manager backed by the shared design.Manager.
func NewManager(initialTheme string) *Manager {
	return &Manager{inner: design.NewManager(initialTheme)}
}

// GetCurrent returns the current design.Theme.
func (m *Manager) GetCurrent() design.Theme { return m.inner.Current() }

// GetCurrentName returns the current theme name.
func (m *Manager) GetCurrentName() string { return m.inner.CurrentName() }

// NextTheme cycles to the next theme.
func (m *Manager) NextTheme() design.Theme { return m.inner.Next() }

// PrevTheme cycles to the previous theme.
func (m *Manager) PrevTheme() design.Theme { return m.inner.Previous() }

// SetTheme sets a specific theme by name.
func (m *Manager) SetTheme(name string) design.Theme {
	m.inner.Set(name)
	return m.inner.Current()
}
