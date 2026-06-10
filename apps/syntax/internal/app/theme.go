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

// Manager handles theme selection and switching using the shared design module.
type Manager struct {
	currentTheme string
	currentIndex int
	themeNames   []string
}

// NewManager creates a new theme manager backed by the design module registry.
func NewManager(initialTheme string) *Manager {
	themeNames := design.List()

	idx := 0
	for i, name := range themeNames {
		if name == initialTheme {
			idx = i
			break
		}
	}

	return &Manager{
		currentTheme: initialTheme,
		currentIndex: idx,
		themeNames:   themeNames,
	}
}

// GetCurrent returns the current design.Theme.
func (m *Manager) GetCurrent() design.Theme {
	return design.Get(m.currentTheme)
}

// GetCurrentName returns the current theme name.
func (m *Manager) GetCurrentName() string {
	return m.currentTheme
}

// NextTheme cycles to the next theme.
func (m *Manager) NextTheme() design.Theme {
	m.currentIndex = (m.currentIndex + 1) % len(m.themeNames)
	m.currentTheme = m.themeNames[m.currentIndex]
	return m.GetCurrent()
}

// PrevTheme cycles to the previous theme.
func (m *Manager) PrevTheme() design.Theme {
	m.currentIndex--
	if m.currentIndex < 0 {
		m.currentIndex = len(m.themeNames) - 1
	}
	m.currentTheme = m.themeNames[m.currentIndex]
	return m.GetCurrent()
}

// SetTheme sets a specific theme by name.
func (m *Manager) SetTheme(name string) design.Theme {
	for i, themeName := range m.themeNames {
		if themeName == name {
			m.currentIndex = i
			m.currentTheme = name
			break
		}
	}
	return m.GetCurrent()
}
