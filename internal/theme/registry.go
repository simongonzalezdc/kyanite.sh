package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a complete color theme for focus.sh
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
}

// All 13 Kyanite Suite themes
var (
	SlateMist = Theme{
		Name:       "Slate Mist",
		Primary:    "#E0E0E0",
		Secondary:  "#B0B0B0",
		Accent:     "#FFFFFF",
		Background: "#1A1A1A",
		Text:       "#E8E8E8",
		Success:    "#90C695",
		Warning:    "#F7D08A",
		Error:      "#E8838F",
		Border:     "#B0B0B0",
		Panel:      "#252525",
	}

	VioletDusk = Theme{
		Name:       "Violet Dusk",
		Primary:    "#9D84B7",
		Secondary:  "#6D5B7F",
		Accent:     "#D4A5FF",
		Background: "#1A1527",
		Text:       "#E8DFF5",
		Success:    "#9FE2BF",
		Warning:    "#F4C796",
		Error:      "#E88FAF",
		Border:     "#6D5B7F",
		Panel:      "#2A1F3A",
	}

	AmberNight = Theme{
		Name:       "Amber Night",
		Primary:    "#B8936E",
		Secondary:  "#6D5B8B",
		Accent:     "#E8C547",
		Background: "#12101A",
		Text:       "#F0E6D8",
		Success:    "#52D3AA",
		Warning:    "#F4A460",
		Error:      "#E8967A",
		Border:     "#6D5B8B",
		Panel:      "#1E1A2A",
	}

	MoltenGold = Theme{
		Name:       "Molten Gold",
		Primary:    "#D4AF37",
		Secondary:  "#8B5A2B",
		Accent:     "#FFD700",
		Background: "#0F0905",
		Text:       "#FFF8DC",
		Success:    "#98D8C8",
		Warning:    "#FFB347",
		Error:      "#CD5C5C",
		Border:     "#8B5A2B",
		Panel:      "#1F1410",
	}

	ClayRoads = Theme{
		Name:       "Clay Roads",
		Primary:    "#8B4513",
		Secondary:  "#A0522D",
		Accent:     "#CD853F",
		Background: "#1A1410",
		Text:       "#F5DEB3",
		Success:    "#7FB89C",
		Warning:    "#DEB887",
		Error:      "#CD5C5C",
		Border:     "#A0522D",
		Panel:      "#2A1F1A",
	}

	IronStorm = Theme{
		Name:       "Iron Storm",
		Primary:    "#DC143C",
		Secondary:  "#696969",
		Accent:     "#FF6347",
		Background: "#0B0C0E",
		Text:       "#F5F5F5",
		Success:    "#90EE90",
		Warning:    "#FFD700",
		Error:      "#FF6B6B",
		Border:     "#696969",
		Panel:      "#1A1C1E",
	}

	JadeTide = Theme{
		Name:       "Jade Tide",
		Primary:    "#20B2AA",
		Secondary:  "#2E8B87",
		Accent:     "#E0F2F1",
		Background: "#0F1419",
		Text:       "#F5F8FA",
		Success:    "#7FB89C",
		Warning:    "#FFD93D",
		Error:      "#FF6B6B",
		Border:     "#2E8B87",
		Panel:      "#1A2429",
	}

	SunsetEmber = Theme{
		Name:       "Sunset Ember",
		Primary:    "#FF6B9D",
		Secondary:  "#FF8C42",
		Accent:     "#FFC312",
		Background: "#1e1e2e",
		Text:       "#FFEAA7",
		Success:    "#55E6C1",
		Warning:    "#FFA502",
		Error:      "#FF6348",
		Border:     "#FF8C42",
		Panel:      "#2E2E3E",
	}

	ForestWhisper = Theme{
		Name:       "Forest Whisper",
		Primary:    "#52B788",
		Secondary:  "#52A068",
		Accent:     "#95D5B2",
		Background: "#1B263B",
		Text:       "#D8F3DC",
		Success:    "#B7E4C7",
		Warning:    "#E9D8A6",
		Error:      "#E07A5F",
		Border:     "#52A068",
		Panel:      "#2B364B",
	}

	ElectricBloom = Theme{
		Name:       "Electric Bloom",
		Primary:    "#FF0080",
		Secondary:  "#00D4FF",
		Accent:     "#FFE600",
		Background: "#0D0221",
		Text:       "#F0F3FF",
		Success:    "#39FF14",
		Warning:    "#FFB700",
		Error:      "#FF1744",
		Border:     "#00D4FF",
		Panel:      "#1D0A31",
	}

	PlasmaPulse = Theme{
		Name:       "Plasma Pulse",
		Primary:    "#39FF14",
		Secondary:  "#00F5FF",
		Accent:     "#FF1493",
		Background: "#0A0118",
		Text:       "#E0FFFF",
		Success:    "#00FF7F",
		Warning:    "#FFAA00",
		Error:      "#FF073A",
		Border:     "#00F5FF",
		Panel:      "#1A0A28",
	}

	IndigoDepths = Theme{
		Name:       "Indigo Depths",
		Primary:    "#4B0082",
		Secondary:  "#36648B",
		Accent:     "#9370DB",
		Background: "#0A0E27",
		Text:       "#E6E6FA",
		Success:    "#7FB89C",
		Warning:    "#DAA520",
		Error:      "#FF69B4",
		Border:     "#36648B",
		Panel:      "#1A1E37",
	}

	SageMeadow = Theme{
		Name:       "Sage Meadow",
		Primary:    "#9CAF88",
		Secondary:  "#B8B8A3",
		Accent:     "#D4D4AF",
		Background: "#1B1B1B",
		Text:       "#E8E8D0",
		Success:    "#A8D5BA",
		Warning:    "#E6B89C",
		Error:      "#D4A5A5",
		Border:     "#B8B8A3",
		Panel:      "#2B2B2B",
	}
)

// AllThemes slice containing all available themes
var AllThemes = []Theme{
	SlateMist,
	VioletDusk,
	AmberNight,
	MoltenGold,
	ClayRoads,
	IronStorm,
	JadeTide,
	SunsetEmber,
	ForestWhisper,
	ElectricBloom,
	PlasmaPulse,
	IndigoDepths,
	SageMeadow,
}

// DefaultTheme is the default theme for focus.sh
var DefaultTheme = AmberNight

// GetThemeByName returns a theme by its name, or DefaultTheme if not found
func GetThemeByName(name string) Theme {
	for _, theme := range AllThemes {
		if theme.Name == name {
			return theme
		}
	}
	return DefaultTheme
}

// GetThemeNames returns a slice of all theme names
func GetThemeNames() []string {
	names := make([]string, len(AllThemes))
	for i, theme := range AllThemes {
		names[i] = theme.Name
	}
	return names
}