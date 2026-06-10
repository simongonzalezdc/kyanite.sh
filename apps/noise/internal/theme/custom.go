// Package theme provides theme management for the noise app,
// delegating to the shared design module for theme definitions and tokens.
package theme

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
	"github.com/pelletier/go-toml/v2"
)

// CustomTheme represents a user-defined theme from TOML.
type CustomTheme struct {
	Name   string      `toml:"name"`
	Colors ThemeColors `toml:"colors"`
}

// ThemeColors holds hex color values.
type ThemeColors struct {
	Primary    string `toml:"primary"`
	Secondary  string `toml:"secondary"`
	Accent     string `toml:"accent"`
	Background string `toml:"background"`
	Text       string `toml:"text"`
	Success    string `toml:"success"`
	Warning    string `toml:"warning"`
	Error      string `toml:"error"`
	Border     string `toml:"border"`
	Panel      string `toml:"panel"`
	Muted      string `toml:"muted"`
}

// LoadCustomThemes loads custom themes from ~/.config/[tool]/themes/
// and registers them via design.RegisterCustom (which validates WCAG AA).
func LoadCustomThemes(toolName string) (map[string]Theme, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	themesDir := filepath.Join(homeDir, ".config", toolName, "themes")
	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
		return map[string]Theme{}, nil
	}

	customThemes := make(map[string]Theme)

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".toml" {
			continue
		}

		var ct CustomTheme
		filePath := filepath.Join(themesDir, entry.Name())

		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		if err := toml.Unmarshal(data, &ct); err != nil {
			continue
		}
		if ct.Name == "" {
			continue
		}

		t := design.Theme{
			Name:       ct.Name,
			Primary:    lipgloss.Color(ct.Colors.Primary),
			Secondary:  lipgloss.Color(ct.Colors.Secondary),
			Accent:     lipgloss.Color(ct.Colors.Accent),
			Background: lipgloss.Color(ct.Colors.Background),
			Text:       lipgloss.Color(ct.Colors.Text),
			Success:    lipgloss.Color(ct.Colors.Success),
			Warning:    lipgloss.Color(ct.Colors.Warning),
			Error:      lipgloss.Color(ct.Colors.Error),
			Border:     lipgloss.Color(ct.Colors.Border),
			Panel:      lipgloss.Color(ct.Colors.Panel),
			Muted:      lipgloss.Color(ct.Colors.Muted),
		}

		// Use RegisterCustom (returns error, not panic) for TOML-loaded themes.
		if err := design.RegisterCustom(t); err != nil {
			continue // Skip themes that fail WCAG validation
		}

		themeID := entry.Name()[:len(entry.Name())-5]
		customThemes[themeID] = t
	}

	return customThemes, nil
}
