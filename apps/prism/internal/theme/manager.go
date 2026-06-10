package theme

import "github.com/kyanite/design"

// Manager handles theme switching using the shared design module.
type Manager struct {
	current design.Theme
}

// NewManager creates a new theme manager with the default theme.
func NewManager() *Manager {
	return &Manager{
		current: design.DefaultTheme(),
	}
}

// SetTheme sets the current theme.
func (m *Manager) SetTheme(t *design.Theme) {
	m.current = *t
}

// CurrentTheme returns the current theme.
func (m *Manager) CurrentTheme() *design.Theme {
	return &m.current
}

// NextTheme cycles to the next theme in the registry.
func (m *Manager) NextTheme() *design.Theme {
	names := design.List()
	for i, name := range names {
		if name == m.current.Name {
			next := (i + 1) % len(names)
			m.current = design.Get(names[next])
			return &m.current
		}
	}
	return &m.current
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
