package theme

// All 10 Kyanite themes with exact hex codes from KYANITE-STANDARDS.md

var (
	Monochrome = Theme{
		Name:       "Monochrome",
		Primary:    "#E0E0E0",
		Secondary:  "#B0B0B0",
		Accent:     "#FFFFFF",
		Background: "#1A1A1A",
		Text:       "#E8E8E8",
		Success:    "#90C695",
	}

	AmberNight = Theme{
		Name:       "Amber Night",
		Primary:    "#D4A574",
		Secondary:  "#8B7355",
		Accent:     "#C0956E",
		Background: "#2A2421",
		Text:       "#F5E6D3",
		Success:    "#A8C686",
	}

	TwilightMist = Theme{
		Name:       "Twilight Mist",
		Primary:    "#9D84B7",
		Secondary:  "#7B68A6",
		Accent:     "#C9AEE8",
		Background: "#2D263E",
		Text:       "#E8DFF5",
		Success:    "#A8C686",
	}

	IndigoDepths = Theme{
		Name:       "Indigo Depths",
		Primary:    "#4169E1",
		Secondary:  "#5F9EA0",
		Accent:     "#DEB887",
		Background: "#0C0C1E",
		Text:       "#F0E68C",
		Success:    "#5F9EA0",
	}

	ForestPath = Theme{
		Name:       "Forest Path",
		Primary:    "#52B788",
		Secondary:  "#52A068",
		Accent:     "#95D5B2",
		Background: "#1B263B",
		Text:       "#D8F3DC",
		Success:    "#B7E4C7",
	}

	ClayEarth = Theme{
		Name:       "Clay Earth",
		Primary:    "#A0644E",
		Secondary:  "#8B4513",
		Accent:     "#D2691E",
		Background: "#2F1F1F",
		Text:       "#F5E6D3",
		Success:    "#9ACD32",
	}

	IronForge = Theme{
		Name:       "Iron Forge",
		Primary:    "#6B4423",
		Secondary:  "#8B4513",
		Accent:     "#CD853F",
		Background: "#1A1A1A",
		Text:       "#F5DEB3",
		Success:    "#98FB98",
	}

	Sunlight = Theme{
		Name:       "Sunlight",
		Primary:    "#FFB92C",
		Secondary:  "#FF8C00",
		Accent:     "#FFA500",
		Background: "#1F1F1F",
		Text:       "#FFFACD",
		Success:    "#32CD32",
	}

	CyanWave = Theme{
		Name:       "Cyan Wave",
		Primary:    "#00D4FF",
		Secondary:  "#00B8D4",
		Accent:     "#00F5FF",
		Background: "#0D1F26",
		Text:       "#E0FFFF",
		Success:    "#00FF7F",
	}

	ElectricRose = Theme{
		Name:       "Electric Rose",
		Primary:    "#FF0080",
		Secondary:  "#D4005C",
		Accent:     "#FF1493",
		Background: "#1A0E1A",
		Text:       "#FFE0F0",
		Success:    "#39FF14",
	}
)

// AllThemes returns all available themes
func AllThemes() []Theme {
	return []Theme{
		Monochrome,
		AmberNight,
		TwilightMist,
		IndigoDepths,
		ForestPath,
		ClayEarth,
		IronForge,
		Sunlight,
		CyanWave,
		ElectricRose,
	}
}

// GetTheme returns a theme by name
func GetTheme(name string) *Theme {
	themes := AllThemes()
	for _, theme := range themes {
		if theme.Name == name {
			return &theme
		}
	}
	return &AmberNight // Default
}
