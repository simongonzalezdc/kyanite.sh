package theme

// Manager handles theme switching
type Manager struct {
	currentTheme *Theme
}

// NewManager creates a new theme manager
func NewManager() *Manager {
	return &Manager{
		currentTheme: &AmberNight,
	}
}

// SetTheme sets the current theme
func (m *Manager) SetTheme(theme *Theme) {
	m.currentTheme = theme
}

// CurrentTheme returns the current theme
func (m *Manager) CurrentTheme() *Theme {
	return m.currentTheme
}

// NextTheme cycles to the next theme
func (m *Manager) NextTheme() *Theme {
	themes := AllThemes()
	for i, theme := range themes {
		if theme.Name == m.currentTheme.Name {
			next := (i + 1) % len(themes)
			m.currentTheme = &themes[next]
			return m.currentTheme
		}
	}
	return m.currentTheme
}
