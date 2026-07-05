package theme

import "github.com/kyanite/design"

// Manager wraps the shared design.Manager with prism-specific convenience methods.
type Manager struct {
	inner *design.Manager
}

// NewManager creates a new theme manager with the default theme.
func NewManager() *Manager {
	return &Manager{inner: design.NewManager("amber-night")}
}

// SetTheme sets the current theme.
func (m *Manager) SetTheme(t *design.Theme) { m.inner.Set(t.Name) }

// CurrentTheme returns a pointer to the current theme.
func (m *Manager) CurrentTheme() *design.Theme {
	t := m.inner.Current()
	return &t
}

// NextTheme cycles to the next theme in the registry.
func (m *Manager) NextTheme() *design.Theme {
	t := m.inner.Next()
	return &t
}

// AllThemes returns all registered themes from the design module.
func AllThemes() []design.Theme {
	names := design.List()
	themes := make([]design.Theme, len(names))
	for i, name := range names {
		themes[i] = design.Get(name)
	}
	return themes
}

// GetTheme returns a theme by name, or the default theme if not found.
func GetTheme(name string) *design.Theme {
	t := design.Get(name)
	if t.Name == "" {
		t = design.DefaultTheme()
	}
	return &t
}
