package design

import "sync"

// Manager provides thread-safe theme selection and cycling.
// All four apps (focus, noise, syntax, prism) share this implementation;
// app-specific features (persistence, change notifications) wrap it.
type Manager struct {
	mu          sync.RWMutex
	currentName string
}

// NewManager creates a Manager initialized to the given theme name.
// Falls back to the default theme if name is empty or unregistered.
func NewManager(initial string) *Manager {
	if initial == "" || Get(initial).Name == "" {
		initial = "amber-night"
	}
	return &Manager{currentName: initial}
}

// Current returns the current Theme.
func (m *Manager) Current() Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return resolveTheme(m.currentName)
}

// CurrentName returns the registered name of the current theme.
func (m *Manager) CurrentName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentName
}

// Set changes the current theme by registered name.
// Falls back to the default theme if name is empty or unregistered.
func (m *Manager) Set(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if name == "" || Get(name).Name == "" {
		name = "amber-night"
	}
	m.currentName = name
}

// Next cycles to the next registered theme and returns it.
func (m *Manager) Next() Theme {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cycle(1)
}

// Previous cycles to the previous registered theme and returns it.
func (m *Manager) Previous() Theme {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cycle(-1)
}

// cycle moves the current theme by delta positions (wrapping) and returns it.
func (m *Manager) cycle(delta int) Theme {
	ids := List()
	if len(ids) == 0 {
		return DefaultTheme()
	}
	for i, name := range ids {
		if name == m.currentName {
			m.currentName = ids[(i+delta+len(ids))%len(ids)]
			return resolveTheme(m.currentName)
		}
	}
	m.currentName = ids[0]
	return resolveTheme(m.currentName)
}

// resolveTheme returns the named theme or the default if not found.
func resolveTheme(name string) Theme {
	if t := Get(name); t.Name != "" {
		return t
	}
	return DefaultTheme()
}
