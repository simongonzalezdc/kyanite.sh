package theme

import "github.com/charmbracelet/lipgloss"

// Registry holds all available built-in themes.
// Keys are kebab-case IDs used for config and switching.
var Registry = make(map[string]Theme)

// DefaultTheme is the ID used when no valid theme is selected.
var DefaultTheme = "amber-night"

func init() {
	// Built-in themes from specification (12)
	Registry["slate"] = Theme{
		Name:       "Slate Mist",
		Primary:    lipgloss.Color("#E0E0E0"),
		Secondary:  lipgloss.Color("#B0B0B0"),
		Accent:     lipgloss.Color("#FFFFFF"),
		Background: lipgloss.Color("#1A1A1A"),
		Text:       lipgloss.Color("#E8E8E8"),
		Success:    lipgloss.Color("#90C695"),
	}

	Registry["amber-night"] = Theme{
		Name:       "Amber Night",
		Primary:    lipgloss.Color("#B8936E"),
		Secondary:  lipgloss.Color("#6D5B8B"),
		Accent:     lipgloss.Color("#E8C547"),
		Background: lipgloss.Color("#12101A"),
		Text:       lipgloss.Color("#F0E6D8"),
		Success:    lipgloss.Color("#7FB89C"),
	}

	Registry["molten-gold"] = Theme{
		Name:       "Molten Gold",
		Primary:    lipgloss.Color("#FFD700"),
		Secondary:  lipgloss.Color("#FF6600"),
		Accent:     lipgloss.Color("#8B00FF"),
		Background: lipgloss.Color("#0D0D0D"),
		Text:       lipgloss.Color("#FFFFFF"),
		Success:    lipgloss.Color("#00FF7F"),
	}

	Registry["clay-roads"] = Theme{
		Name:       "Clay Roads",
		Primary:    lipgloss.Color("#CC5500"),
		Secondary:  lipgloss.Color("#D2691E"),
		Accent:     lipgloss.Color("#FFD700"),
		Background: lipgloss.Color("#1C1410"),
		Text:       lipgloss.Color("#FFF8DC"),
		Success:    lipgloss.Color("#9ACD32"),
	}

	Registry["iron-storm"] = Theme{
		Name:       "Iron Storm",
		Primary:    lipgloss.Color("#DC143C"),
		Secondary:  lipgloss.Color("#708090"),
		Accent:     lipgloss.Color("#FF4500"),
		Background: lipgloss.Color("#0D0D0D"),
		Text:       lipgloss.Color("#E8E8E8"),
		Success:    lipgloss.Color("#4A6741"),
	}

	Registry["jade-tide"] = Theme{
		Name:       "Jade Tide",
		Primary:    lipgloss.Color("#20B2AA"),
		Secondary:  lipgloss.Color("#2E8B87"),
		Accent:     lipgloss.Color("#E0F2F1"),
		Background: lipgloss.Color("#0F1419"),
		Text:       lipgloss.Color("#F5F8FA"),
		Success:    lipgloss.Color("#7FB89C"),
	}

	Registry["sunset-ember"] = Theme{
		Name:       "Sunset Ember",
		Primary:    lipgloss.Color("#FF6B9D"),
		Secondary:  lipgloss.Color("#FF8C42"),
		Accent:     lipgloss.Color("#FFC312"),
		Background: lipgloss.Color("#1e1e2e"),
		Text:       lipgloss.Color("#FFEAA7"),
		Success:    lipgloss.Color("#55E6C1"),
	}

	Registry["forest-whisper"] = Theme{
		Name:       "Forest Whisper",
		Primary:    lipgloss.Color("#52B788"),
		Secondary:  lipgloss.Color("#52A068"),
		Accent:     lipgloss.Color("#95D5B2"),
		Background: lipgloss.Color("#1B263B"),
		Text:       lipgloss.Color("#D8F3DC"),
		Success:    lipgloss.Color("#B7E4C7"),
	}

	Registry["electric-bloom"] = Theme{
		Name:       "Electric Bloom",
		Primary:    lipgloss.Color("#FF0080"),
		Secondary:  lipgloss.Color("#00D4FF"),
		Accent:     lipgloss.Color("#FFE600"),
		Background: lipgloss.Color("#0D0221"),
		Text:       lipgloss.Color("#F0F3FF"),
		Success:    lipgloss.Color("#39FF14"),
	}

	Registry["plasma-pulse"] = Theme{
		Name:       "Plasma Pulse",
		Primary:    lipgloss.Color("#39FF14"),
		Secondary:  lipgloss.Color("#00F5FF"),
		Accent:     lipgloss.Color("#FF1493"),
		Background: lipgloss.Color("#0A0118"),
		Text:       lipgloss.Color("#E0FFFF"),
		Success:    lipgloss.Color("#00FF7F"),
	}

	Registry["indigo-depths"] = Theme{
		Name:       "Indigo Depths",
		Primary:    lipgloss.Color("#4169E1"),
		Secondary:  lipgloss.Color("#5F9EA0"),
		Accent:     lipgloss.Color("#DEB887"),
		Background: lipgloss.Color("#0C0C1E"),
		Text:       lipgloss.Color("#F0E68C"),
		Success:    lipgloss.Color("#5F9EA0"),
	}

	Registry["sage-meadow"] = Theme{
		Name:       "Sage Meadow",
		Primary:    lipgloss.Color("#8FBC8F"),
		Secondary:  lipgloss.Color("#D2B48C"),
		Accent:     lipgloss.Color("#F4A460"),
		Background: lipgloss.Color("#1A1612"),
		Text:       lipgloss.Color("#FAF0E6"),
		Success:    lipgloss.Color("#9ACD32"),
	}

	// Preserve original "Violet Dusk" (existing app theme historically named "Midnight Jazz").
	// Keep name, assets and settings intact; provide ID "violet-dusk".
	Registry["violet-dusk"] = Theme{
		Name:       "Violet Dusk",
		Primary:    lipgloss.Color("#9D84B7"), // Soft purple - main brand
		Secondary:  lipgloss.Color("#5E4B8B"), // Deep purple - secondary actions
		Accent:     lipgloss.Color("#F4D03F"), // Gold - highlights
		Background: lipgloss.Color("#0A0E27"), // Deep navy - main background
		Text:       lipgloss.Color("#E8DFF5"), // Light lavender - main text
		Success:    lipgloss.Color("#52D3AA"), // Mint green - success states
	}

	// Backwards-compatible alias: some configs use "midnight_jazz" (snake) or "midnight-jazz" (kebab).
	Registry["midnight-jazz"] = Registry["violet-dusk"]
	Registry["midnight_jazz"] = Registry["violet-dusk"]
	Registry["purple-jazz"] = Registry["violet-dusk"] // Additional alias for old name
}

// GetTheme returns a theme by id; falls back to DefaultTheme if not found.
func GetTheme(id string) Theme {
	if t, ok := Registry[id]; ok {
		return t
	}
	return Registry[DefaultTheme]
}

// ListThemes returns the list of registered theme IDs (unordered).
func ListThemes() []string {
	out := make([]string, 0, len(Registry))
	for k := range Registry {
		out = append(out, k)
	}
	return out
}