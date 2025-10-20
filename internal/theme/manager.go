package theme

import (
	"sync"

	"github.com/Kyanite/noise/internal/config"
)

// Manager handles runtime theme switching in a thread-safe way.
// It stores both the active Theme and its ID so callers can query the current ID.
type Manager struct {
	mu        sync.RWMutex
	current   Theme
	currentID string
}

// NewManager creates a manager initialized to the given theme id (falls back to default).
func NewManager(initID string) *Manager {
	id := initID
	if id == "" {
		id = DefaultTheme
	}
	return &Manager{
		current:   GetTheme(id),
		currentID: id,
	}
}

// Current returns the currently active theme.
func (m *Manager) Current() Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

// CurrentID returns the currently active theme id.
func (m *Manager) CurrentID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentID
}

// SetTheme sets the current theme by id. Returns true if the id existed.
func (m *Manager) SetTheme(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t, ok := Registry[id]; ok {
		m.current = t
		m.currentID = id
		return true
	}
	// fallback: try default
	m.current = Registry[DefaultTheme]
	m.currentID = DefaultTheme
	return false
}

// SetThemeIfExists sets theme only if the id exists and returns whether it changed.
func (m *Manager) SetThemeIfExists(id string) (Theme, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t, ok := Registry[id]; ok {
		m.current = t
		m.currentID = id
		return m.current, true
	}
	return m.current, false
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

// Global singleton
var (
	globalManager *Manager
	once          sync.Once
)

// GetManager returns the global theme manager singleton.
// If not initialized, it uses DefaultTheme.
func GetManager() *Manager {
	once.Do(func() {
		globalManager = NewManager(DefaultTheme)
	})
	return globalManager
}