package theme

// Manager handles theme selection and switching
type Manager struct {
	currentTheme string
	currentIndex int
}

// NewManager creates a new theme manager
func NewManager(initialTheme string) *Manager {
	// Find index of initial theme
	idx := 0
	for i, name := range ThemeNames {
		if name == initialTheme {
			idx = i
			break
		}
	}

	return &Manager{
		currentTheme: initialTheme,
		currentIndex: idx,
	}
}

// GetCurrent returns the current theme
func (m *Manager) GetCurrent() Theme {
	return GetTheme(m.currentTheme)
}

// GetCurrentName returns the current theme name
func (m *Manager) GetCurrentName() string {
	return m.currentTheme
}

// NextTheme cycles to the next theme
func (m *Manager) NextTheme() Theme {
	m.currentIndex = (m.currentIndex + 1) % len(ThemeNames)
	m.currentTheme = ThemeNames[m.currentIndex]
	return m.GetCurrent()
}

// PrevTheme cycles to the previous theme
func (m *Manager) PrevTheme() Theme {
	m.currentIndex--
	if m.currentIndex < 0 {
		m.currentIndex = len(ThemeNames) - 1
	}
	m.currentTheme = ThemeNames[m.currentIndex]
	return m.GetCurrent()
}

// SetTheme sets a specific theme by name
func (m *Manager) SetTheme(name string) Theme {
	// Find the index
	for i, themeName := range ThemeNames {
		if themeName == name {
			m.currentIndex = i
			m.currentTheme = name
			break
		}
	}
	return m.GetCurrent()
}
