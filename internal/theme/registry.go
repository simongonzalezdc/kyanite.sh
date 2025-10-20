package theme

import "github.com/charmbracelet/lipgloss"

var (
	Monochrome = Theme{
		Name:       "Monochrome",
		Primary:    lipgloss.Color("#FFFFFF"),
		Secondary:  lipgloss.Color("#999999"),
		Accent:     lipgloss.Color("#FFFFFF"),
		Background: lipgloss.Color("#000000"),
		Text:       lipgloss.Color("#FFFFFF"),
		Success:    lipgloss.Color("#CCCCCC"),
		Warning:    lipgloss.Color("#888888"),
		Error:      lipgloss.Color("#666666"),
		Border:     lipgloss.Color("#999999"),
		Panel:      lipgloss.Color("#000000"),
	}

	AmberNight = Theme{
		Name:       "Amber Night",
		Primary:    lipgloss.Color("#D4A574"),
		Secondary:  lipgloss.Color("#9D84B7"),
		Accent:     lipgloss.Color("#F4D03F"),
		Background: lipgloss.Color("#0A0E27"),
		Text:       lipgloss.Color("#E8DFF5"),
		Success:    lipgloss.Color("#52D3AA"),
		Warning:    lipgloss.Color("#FFA502"),
		Error:      lipgloss.Color("#EA2027"),
		Border:     lipgloss.Color("#9D84B7"),
		Panel:      lipgloss.Color("#15182F"),
	}

	TwilightMist = Theme{
		Name:       "Twilight Mist",
		Primary:    lipgloss.Color("#B8A3C9"),
		Secondary:  lipgloss.Color("#8E7B9D"),
		Accent:     lipgloss.Color("#D4C5E0"),
		Background: lipgloss.Color("#151520"),
		Text:       lipgloss.Color("#E8E4F0"),
		Success:    lipgloss.Color("#90C695"),
		Warning:    lipgloss.Color("#C9A87C"),
		Error:      lipgloss.Color("#C77777"),
		Border:     lipgloss.Color("#8E7B9D"),
		Panel:      lipgloss.Color("#1F1F2A"),
	}

	IndigoDepths = Theme{
		Name:       "Indigo Depths",
		Primary:    lipgloss.Color("#4169E1"),
		Secondary:  lipgloss.Color("#5F9EA0"),
		Accent:     lipgloss.Color("#87CEEB"),
		Background: lipgloss.Color("#0A0A1A"),
		Text:       lipgloss.Color("#E6F2FF"),
		Success:    lipgloss.Color("#52D3AA"),
		Warning:    lipgloss.Color("#FFB84D"),
		Error:      lipgloss.Color("#FF6B6B"),
		Border:     lipgloss.Color("#5F9EA0"),
		Panel:      lipgloss.Color("#14142A"),
	}

	ForestPath = Theme{
		Name:       "Forest Path",
		Primary:    lipgloss.Color("#8FBC8F"),
		Secondary:  lipgloss.Color("#6B8E6B"),
		Accent:     lipgloss.Color("#B4D7B4"),
		Background: lipgloss.Color("#1A1F1A"),
		Text:       lipgloss.Color("#E8F5E8"),
		Success:    lipgloss.Color("#90EE90"),
		Warning:    lipgloss.Color("#DAA520"),
		Error:      lipgloss.Color("#CD5C5C"),
		Border:     lipgloss.Color("#6B8E6B"),
		Panel:      lipgloss.Color("#232823"),
	}

	ClayEarth = Theme{
		Name:       "Clay Earth",
		Primary:    lipgloss.Color("#A0522D"),
		Secondary:  lipgloss.Color("#8B4513"),
		Accent:     lipgloss.Color("#DEB887"),
		Background: lipgloss.Color("#1A1410"),
		Text:       lipgloss.Color("#F5E6D3"),
		Success:    lipgloss.Color("#8FBC8F"),
		Warning:    lipgloss.Color("#CD853F"),
		Error:      lipgloss.Color("#CD5C5C"),
		Border:     lipgloss.Color("#8B4513"),
		Panel:      lipgloss.Color("#231C18"),
	}

	IronForge = Theme{
		Name:       "Iron Forge",
		Primary:    lipgloss.Color("#DC143C"),
		Secondary:  lipgloss.Color("#4A4A4A"),
		Accent:     lipgloss.Color("#FF6347"),
		Background: lipgloss.Color("#1A0A0A"),
		Text:       lipgloss.Color("#FFE6E6"),
		Success:    lipgloss.Color("#90C695"),
		Warning:    lipgloss.Color("#FFB84D"),
		Error:      lipgloss.Color("#FF4444"),
		Border:     lipgloss.Color("#4A4A4A"),
		Panel:      lipgloss.Color("#220F0F"),
	}

	Sunlight = Theme{
		Name:       "Sunlight",
		Primary:    lipgloss.Color("#FFD700"),
		Secondary:  lipgloss.Color("#DAA520"),
		Accent:     lipgloss.Color("#FFF8DC"),
		Background: lipgloss.Color("#0F0F0A"),
		Text:       lipgloss.Color("#FFFACD"),
		Success:    lipgloss.Color("#98D982"),
		Warning:    lipgloss.Color("#FF9800"),
		Error:      lipgloss.Color("#D32F2F"),
		Border:     lipgloss.Color("#DAA520"),
		Panel:      lipgloss.Color("#18180F"),
	}

	CyanWave = Theme{
		Name:       "Cyan Wave",
		Primary:    lipgloss.Color("#00CED1"),
		Secondary:  lipgloss.Color("#4682B4"),
		Accent:     lipgloss.Color("#7FFFD4"),
		Background: lipgloss.Color("#0A1418"),
		Text:       lipgloss.Color("#E0F7FA"),
		Success:    lipgloss.Color("#52D3AA"),
		Warning:    lipgloss.Color("#FFB84D"),
		Error:      lipgloss.Color("#FF6B6B"),
		Border:     lipgloss.Color("#4682B4"),
		Panel:      lipgloss.Color("#132126"),
	}

	ElectricRose = Theme{
		Name:       "Electric Rose",
		Primary:    lipgloss.Color("#FF1493"),
		Secondary:  lipgloss.Color("#C71585"),
		Accent:     lipgloss.Color("#00CED1"),
		Background: lipgloss.Color("#1A0A1A"),
		Text:       lipgloss.Color("#FFF0F8"),
		Success:    lipgloss.Color("#52D3AA"),
		Warning:    lipgloss.Color("#FFB84D"),
		Error:      lipgloss.Color("#FF4444"),
		Border:     lipgloss.Color("#C71585"),
		Panel:      lipgloss.Color("#230F23"),
	}
)

var Registry = map[string]Theme{
	"monochrome":    Monochrome,
	"amber-night":   AmberNight,
	"twilight-mist": TwilightMist,
	"indigo-depths": IndigoDepths,
	"forest-path":   ForestPath,
	"clay-earth":    ClayEarth,
	"iron-forge":    IronForge,
	"sunlight":      Sunlight,
	"cyan-wave":     CyanWave,
	"electric-rose": ElectricRose,
}

func Default() Theme { return AmberNight }

func GetTheme(id string) Theme {
	id = migrateThemeID(id)
	if t, ok := Registry[id]; ok {
		return t
	}
	return Default()
}

func ListThemes() []string {
	return []string{"monochrome", "amber-night", "twilight-mist", "indigo-depths", "forest-path", "clay-earth", "iron-forge", "sunlight", "cyan-wave", "electric-rose"}
}

func GetThemeByName(name string) Theme {
	for _, t := range Registry {
		if t.Name == name {
			return t
		}
	}
	return Default()
}

func GetThemeNames() []string {
	names := make([]string, 0, len(Registry))
	for _, id := range ListThemes() {
		names = append(names, Registry[id].Name)
	}
	return names
}

func migrateThemeID(oldID string) string {
	m := map[string]string{
		"slate-mist":     "twilight-mist",
		"violet-dusk":    "twilight-mist",
		"molten-gold":    "sunlight",
		"clay-roads":     "clay-earth",
		"iron-storm":     "iron-forge",
		"jade-tide":      "cyan-wave",
		"sunset-ember":   "electric-rose",
		"forest-whisper": "forest-path",
		"electric-bloom": "electric-rose",
		"plasma-pulse":   "amber-night",
		"sage-meadow":    "forest-path",
		"synthwave":      "amber-night",
		"light":          "monochrome",
		"plain":          "amber-night",
	}
	if nid, ok := m[oldID]; ok {
		return nid
	}
	return oldID
}
