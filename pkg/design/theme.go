package design

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// Theme represents a color theme with 11 semantic color fields.
// All fields use lipgloss.Color for direct compatibility with lipgloss styles.
type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Background lipgloss.Color
	Text       lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Border     lipgloss.Color
	Panel      lipgloss.Color
	Muted      lipgloss.Color
}

// registry holds all registered themes, keyed by name.
var (
	registry   = make(map[string]Theme)
	registryMu sync.RWMutex
)

// RegisterBuiltIn registers a built-in theme. It panics if the theme fails
// WCAG AA contrast validation for any text-on-background pair.
// Use this for the curated themes at init time.
func RegisterBuiltIn(t Theme) {
	if t.Name == "" {
		panic("design: theme must have a name")
	}
	if err := ValidateThemeAA(t); err != nil {
		panic(fmt.Sprintf("design: built-in theme %q fails WCAG AA: %s", t.Name, err))
	}
	registryMu.Lock()
	registry[t.Name] = t
	registryMu.Unlock()
}

// RegisterCustom registers a user-defined theme. It returns an error if the
// theme fails WCAG AA contrast validation, instead of panicking.
// Use this for TOML/JSON-loaded custom themes at runtime.
func RegisterCustom(t Theme) error {
	if t.Name == "" {
		return fmt.Errorf("design: theme must have a name")
	}
	if err := ValidateThemeAA(t); err != nil {
		return fmt.Errorf("theme %q fails WCAG AA: %w", t.Name, err)
	}
	registryMu.Lock()
	registry[t.Name] = t
	registryMu.Unlock()
	return nil
}

// Get returns the theme with the given name, or the zero Theme if not found.
func Get(name string) Theme {
	registryMu.RLock()
	defer registryMu.RUnlock()
	return registry[name]
}

// List returns all registered theme names.
func List() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// DefaultTheme returns the default theme (amber-night).
func DefaultTheme() Theme {
	registryMu.RLock()
	t, ok := registry["amber-night"]
	registryMu.RUnlock()
	if !ok {
		panic("design: default theme \"amber-night\" not registered")
	}
	return t
}

func init() {
	registerThemes()
}
