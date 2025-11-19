package theme

import (
	"encoding/json"
	"github.com/Kyanite/noise/internal/config"
	"os"
	"path/filepath"
	"sync"
)

var (
	globalManager *Manager
	once          sync.Once
)

// Manager handles theme selection and switching
type Manager struct {
	mu      sync.RWMutex
	current Theme
	// Add a channel for theme change notifications
	themeChangeChan chan Theme
}

// GetManager returns the global theme manager (singleton)
func GetManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{
			current:         Default(),
			themeChangeChan: make(chan Theme, 10), // Buffered channel for theme changes
		}
		// Load saved theme preference
		if err := globalManager.LoadThemePreference(); err != nil {
			// If loading fails, stick with default theme
		}
	})
	return globalManager
}

// SetTheme sets the current theme by ID
func (m *Manager) SetTheme(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.current = GetTheme(id)

	// Notify listeners of theme change
	select {
	case m.themeChangeChan <- m.current:
	default:
		// Channel is full, drop the notification
	}

	// Save preference asynchronously with a copy of the current theme name
	currentThemeName := m.current.Name
	go func() {
		if err := m.saveThemePreferenceByName(currentThemeName); err != nil {
			// Log error but don't fail the theme change
			// In a real app, you'd use proper logging
		}
	}()
}

// Current returns the current theme
func (m *Manager) Current() Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

// Next cycles to the next theme in the registry
func (m *Manager) Next() Theme {
	m.mu.Lock()
	defer m.mu.Unlock()

	themes := ListThemes()
	currentID := ""

	// Find current theme ID
	for id, theme := range Registry {
		if theme.Name == m.current.Name {
			currentID = id
			break
		}
	}

	// Find next theme
	for i, id := range themes {
		if id == currentID {
			nextIndex := (i + 1) % len(themes)
			m.current = Registry[themes[nextIndex]]
			break
		}
	}

	// Notify listeners of theme change
	select {
	case m.themeChangeChan <- m.current:
	default:
		// Channel is full, drop the notification
	}

	// Save preference asynchronously with a copy of the current theme name
	currentThemeName := m.current.Name
	go func() {
		if err := m.saveThemePreferenceByName(currentThemeName); err != nil {
			// Log error but don't fail the theme change
		}
	}()

	return m.current
}

// Previous cycles to the previous theme in the registry
func (m *Manager) Previous() Theme {
	m.mu.Lock()
	defer m.mu.Unlock()

	themes := ListThemes()
	currentID := ""

	// Find current theme ID
	for id, theme := range Registry {
		if theme.Name == m.current.Name {
			currentID = id
			break
		}
	}

	// Find previous theme
	for i, id := range themes {
		if id == currentID {
			prevIndex := (i - 1 + len(themes)) % len(themes)
			m.current = Registry[themes[prevIndex]]
			break
		}
	}

	// Notify listeners of theme change
	select {
	case m.themeChangeChan <- m.current:
	default:
		// Channel is full, drop the notification
	}

	// Save preference asynchronously with a copy of the current theme name
	currentThemeName := m.current.Name
	go func() {
		if err := m.saveThemePreferenceByName(currentThemeName); err != nil {
			// Log error but don't fail the theme change
		}
	}()

	return m.current
}

// ApplyConfig applies theme configuration from the app config (UI.Theme).
// If cfg is nil or theme not present, falls back to DefaultTheme.
func (m *Manager) ApplyConfig(cfg *config.Config) {
	if cfg == nil {
		m.SetTheme(DefaultTheme)
		return
	}
	// prefer UI.Theme if set
	themeID := cfg.UI.Theme
	if themeID == "" {
		m.SetTheme(DefaultTheme)
		return
	}
	m.SetTheme(themeID)
}

// ThemePreference represents a saved theme preference
type ThemePreference struct {
	ThemeID string `json:"theme_id"`
}

// SaveThemePreference saves the current theme preference to file
func (m *Manager) SaveThemePreference() error {
	m.mu.RLock()
	themeName := m.current.Name
	m.mu.RUnlock()
	return m.saveThemePreferenceByName(themeName)
}

// saveThemePreferenceByName saves a theme preference by theme name (thread-safe helper)
func (m *Manager) saveThemePreferenceByName(themeName string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "noise")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "theme.json")

	// Find theme ID by name
	currentID := ""
	for id, theme := range Registry {
		if theme.Name == themeName {
			currentID = id
			break
		}
	}

	pref := ThemePreference{
		ThemeID: currentID,
	}

	data, err := json.MarshalIndent(pref, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

// LoadThemePreference loads theme preference from file
func (m *Manager) LoadThemePreference() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFile := filepath.Join(homeDir, ".config", "noise", "theme.json")

	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil // No preference saved, use default
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	var pref ThemePreference
	if err := json.Unmarshal(data, &pref); err != nil {
		return err
	}

	if pref.ThemeID != "" {
		m.SetTheme(pref.ThemeID)
	}

	return nil
}

// GetThemeChangeChannel returns a channel for listening to theme changes
func (m *Manager) GetThemeChangeChannel() <-chan Theme {
	return m.themeChangeChan
}
