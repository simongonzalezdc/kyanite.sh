// Package theme provides theme definitions, registration, and management.
package theme

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
)

// CustomTheme represents a user-defined theme from TOML
type CustomTheme struct {
	Name   string      `toml:"name"`
	Colors ThemeColors `toml:"colors"`
}

// ThemeColors holds hex color values
type ThemeColors struct {
	Primary    string `toml:"primary"`
	Secondary  string `toml:"secondary"`
	Accent     string `toml:"accent"`
	Background string `toml:"background"`
	Text       string `toml:"text"`
	Success    string `toml:"success"`
	Warning    string `toml:"warning"`
	Error      string `toml:"error"`
}

// LoadCustomThemes loads custom themes from ~/.config/[tool]/themes/
func LoadCustomThemes(toolName string) (map[string]Theme, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	themesDir := filepath.Join(homeDir, ".config", toolName, "themes")

	// Check if themes directory exists
	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
		return map[string]Theme{}, nil
	}

	customThemes := make(map[string]Theme)

	// Read all .toml files in themes directory
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

		// Read and decode TOML file
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files we can't read
		}

		if err := toml.Unmarshal(data, &ct); err != nil {
			continue // Skip invalid TOML files
		}

		// Validate required fields
		if ct.Name == "" {
			continue
		}

		// Convert to Theme
		theme := Theme{
			Name:       ct.Name,
			Primary:    lipgloss.Color(ct.Colors.Primary),
			Secondary:  lipgloss.Color(ct.Colors.Secondary),
			Accent:     lipgloss.Color(ct.Colors.Accent),
			Background: lipgloss.Color(ct.Colors.Background),
			Text:       lipgloss.Color(ct.Colors.Text),
			Success:    lipgloss.Color(ct.Colors.Success),
			Warning:    lipgloss.Color(ct.Colors.Warning),
			Error:      lipgloss.Color(ct.Colors.Error),
		}

		// Use filename without extension as ID
		themeID := entry.Name()[:len(entry.Name())-5]
		customThemes[themeID] = theme
	}

	return customThemes, nil
}
