// Package theme provides theme management for the noise app,
// delegating to the shared design module for theme definitions and tokens.
package theme

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/kyanite/design"
)

// Theme is an alias for design.Theme so existing call sites compile unchanged.
type Theme = design.Theme

// ThemeChangeMsg represents a message for theme changes.
type ThemeChangeMsg struct {
	Theme Theme
}

// DefaultThemeID is the ID used when no valid theme is selected.
const DefaultThemeID = "amber-night"

// Default returns the default theme (Amber Night).
func Default() Theme { return design.DefaultTheme() }

// GetTheme returns a theme by ID, falling back to default.
func GetTheme(id string) Theme {
	t := design.Get(id)
	if t.Name == "" {
		return design.DefaultTheme()
	}
	return t
}

// ListThemes returns all theme IDs from the design module registry.
func ListThemes() []string { return design.List() }

// Manager handles theme selection, persistence, and switching.
// It wraps the shared design.Manager with noise-specific persistence
// and change notification.
type Manager struct {
	inner           *design.Manager
	themeChangeChan chan Theme
}

var (
	globalManager *Manager
	once          sync.Once
)

// GetManager returns the global theme manager (singleton).
func GetManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{
			inner:           design.NewManager(DefaultThemeID),
			themeChangeChan: make(chan Theme, 10),
		}
		_ = globalManager.LoadThemePreference()
	})
	return globalManager
}

// Current returns the current theme.
func (m *Manager) Current() Theme { return m.inner.Current() }

// SetTheme sets the current theme by ID and persists the preference.
func (m *Manager) SetTheme(id string) {
	m.inner.Set(id)
	t := m.inner.Current()

	select {
	case m.themeChangeChan <- t:
	default:
	}

	go func(savedID string) {
		defer func() {
			if r := recover(); r != nil {
				_ = r
			}
		}()
		_ = saveThemePreferenceByID(savedID)
	}(id)
}

// Next cycles to the next theme.
func (m *Manager) Next() Theme {
	m.inner.Next()
	return m.inner.Current()
}

// Previous cycles to the previous theme.
func (m *Manager) Previous() Theme {
	m.inner.Previous()
	return m.inner.Current()
}

// GetThemeChangeChannel returns a channel for listening to theme changes.
func (m *Manager) GetThemeChangeChannel() <-chan Theme {
	return m.themeChangeChan
}

// SaveThemePreference saves the current theme preference to file.
func (m *Manager) SaveThemePreference() error {
	return saveThemePreferenceByID(m.inner.CurrentName())
}

// LoadThemePreference loads theme preference from file.
func (m *Manager) LoadThemePreference() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configFile := filepath.Join(homeDir, ".config", "noise", "theme.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	var pref themePreference
	if err := json.Unmarshal(data, &pref); err != nil {
		return err
	}
	if pref.ThemeID != "" {
		m.SetTheme(pref.ThemeID)
	}
	return nil
}

type themePreference struct {
	ThemeID string `json:"theme_id"`
}

func saveThemePreferenceByID(themeID string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(homeDir, ".config", "noise")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(themePreference{ThemeID: themeID}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(configDir, "theme.json"), data, 0o644)
}

// ValidateTheme checks if a theme meets accessibility standards.
func ValidateTheme(t Theme) []string {
	err := design.ValidateThemeAA(t)
	if err == nil {
		return nil
	}
	return []string{err.Error()}
}

// Registry is a compatibility shim that returns all registered themes
// as a map keyed by theme name.
var Registry map[string]Theme

func init() {
	Registry = make(map[string]Theme)
	for _, name := range design.List() {
		Registry[name] = design.Get(name)
	}
}
