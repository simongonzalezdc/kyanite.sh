package theme

import (
	"sync"

	"github.com/kyanite/design"
)

// Theme is an alias for design.Theme so consumers can migrate gradually.
type Theme = design.Theme

var (
	globalManager *Manager
	once          sync.Once
)

// Manager stores the current theme name and delegates to the design module.
type Manager struct {
	mu          sync.RWMutex
	currentName string
}

// GetManager returns the singleton theme manager.
func GetManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{currentName: "amber-night"}
	})
	return globalManager
}

// SetTheme sets the current theme by name (with migration of old IDs).
func (m *Manager) SetTheme(id string) {
	id = migrateThemeID(id)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentName = id
}

// Current returns the current design.Theme.
func (m *Manager) Current() Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t := design.Get(m.currentName)
	if t.Name == "" {
		return design.DefaultTheme()
	}
	return t
}

// Next cycles to the next theme and returns it.
func (m *Manager) Next() Theme {
	m.mu.Lock()
	defer m.mu.Unlock()
	ids := design.List()
	for i, name := range ids {
		if name == m.currentName {
			m.currentName = ids[(i+1)%len(ids)]
			t := design.Get(m.currentName)
			if t.Name == "" {
				return design.DefaultTheme()
			}
			return t
		}
	}
	m.currentName = "amber-night"
	return design.DefaultTheme()
}

// Previous cycles to the previous theme and returns it.
func (m *Manager) Previous() Theme {
	m.mu.Lock()
	defer m.mu.Unlock()
	ids := design.List()
	for i, name := range ids {
		if name == m.currentName {
			m.currentName = ids[(i-1+len(ids))%len(ids)]
			t := design.Get(m.currentName)
			if t.Name == "" {
				return design.DefaultTheme()
			}
			return t
		}
	}
	m.currentName = "amber-night"
	return design.DefaultTheme()
}

// migrateThemeID maps old Focus theme IDs to design module theme names.
func migrateThemeID(oldID string) string {
	m := map[string]string{
		"slate-mist":     "amber-night",
		"violet-dusk":    "amber-night",
		"molten-gold":    "amber-night",
		"clay-roads":     "amber-night",
		"iron-storm":     "amber-night",
		"jade-tide":      "cyan-wave",
		"sunset-ember":   "electric-rose",
		"forest-whisper": "forest-path",
		"electric-bloom": "electric-rose",
		"plasma-pulse":   "amber-night",
		"sage-meadow":    "forest-path",
		"synthwave":      "electric-rose",
		"light":          "monochrome",
		"plain":          "amber-night",
	}
	if nid, ok := m[oldID]; ok {
		return nid
	}
	return oldID
}

// Default returns the default theme.
func Default() Theme { return design.DefaultTheme() }

// ListThemes returns all design module theme names.
func ListThemes() []string { return design.List() }

// GetTheme returns a theme by name.
func GetTheme(name string) Theme {
	name = migrateThemeID(name)
	t := design.Get(name)
	if t.Name == "" {
		return design.DefaultTheme()
	}
	return t
}

// GetThemeByName returns a theme by its display Name field.
func GetThemeByName(name string) Theme {
	for _, themeName := range design.List() {
		t := design.Get(themeName)
		if t.Name == name {
			return t
		}
	}
	return design.DefaultTheme()
}

// GetThemeNames returns display names of all themes.
func GetThemeNames() []string {
	names := make([]string, 0, len(design.List()))
	for _, id := range design.List() {
		t := design.Get(id)
		if t.Name != "" {
			names = append(names, t.Name)
		}
	}
	return names
}
