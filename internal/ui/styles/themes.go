package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a complete color theme
type Theme struct {
	Name        string
	Description string
	Colors      ThemeColors
}

// ThemeColors contains all color definitions for a theme
type ThemeColors struct {
	// Primary Colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color

	// Functional Colors
	Success lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color
	Info    lipgloss.Color

	// Background & Text Colors
	Background    lipgloss.Color
	TextPrimary   lipgloss.Color
	TextSecondary lipgloss.Color
	TextMuted     lipgloss.Color
	TextAccent    lipgloss.Color

	// Extended Palette
	Dark1 lipgloss.Color
	Dark2 lipgloss.Color
	Dark3 lipgloss.Color
}

// Predefined themes
var (
	// Theme 1: Midnight Jazz (current default)
	MidnightJazzTheme = Theme{
		Name:        "Midnight Jazz",
		Description: "Deep navy background with purple and gold accents",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#9D84B7"), // Soft purple
			Secondary: lipgloss.Color("#5E4B8B"), // Deep purple
			Accent:    lipgloss.Color("#F4D03F"), // Gold
			Success:   lipgloss.Color("#52D3AA"), // Mint green
			Warning:   lipgloss.Color("#FFA500"), // Orange
			Error:     lipgloss.Color("#FF6347"), // Tomato
			Info:      lipgloss.Color("#87CEEB"), // Sky blue
			Background:    lipgloss.Color("#0A0E27"), // Deep navy
			TextPrimary:   lipgloss.Color("#E8DFF5"), // Light lavender
			TextSecondary: lipgloss.Color("#9D84B7"), // Soft purple
			TextMuted:     lipgloss.Color("#5E4B8B"), // Deep purple
			TextAccent:    lipgloss.Color("#F4D03F"), // Gold
			Dark1:         lipgloss.Color("#0A0E27"), // Main background
			Dark2:         lipgloss.Color("#1A1E37"), // Lighter background
			Dark3:         lipgloss.Color("#2A2E47"), // Border/divider
		},
	}

	// Theme 2: Neon Dreams
	NeonDreamsTheme = Theme{
		Name:        "Neon Dreams",
		Description: "Dark background with vibrant neon colors",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#FF00FF"), // Magenta
			Secondary: lipgloss.Color("#00FFFF"), // Cyan
			Accent:    lipgloss.Color("#FFFF00"), // Yellow
			Success:   lipgloss.Color("#00FF00"), // Green
			Warning:   lipgloss.Color("#FFA500"), // Orange
			Error:     lipgloss.Color("#FF0000"), // Red
			Info:      lipgloss.Color("#00FFFF"), // Cyan
			Background:    lipgloss.Color("#000000"), // Black
			TextPrimary:   lipgloss.Color("#FFFFFF"), // White
			TextSecondary: lipgloss.Color("#CC00CC"), // Light magenta
			TextMuted:     lipgloss.Color("#666666"), // Gray
			TextAccent:    lipgloss.Color("#FFFF00"), // Yellow
			Dark1:         lipgloss.Color("#000000"), // Black
			Dark2:         lipgloss.Color("#1A1A1A"), // Dark gray
			Dark3:         lipgloss.Color("#333333"), // Medium gray
		},
	}

	// Theme 3: Forest Retreat
	ForestRetreatTheme = Theme{
		Name:        "Forest Retreat",
		Description: "Natural green tones with earthy accents",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#228B22"), // Forest green
			Secondary: lipgloss.Color("#8B4513"), // Saddle brown
			Accent:    lipgloss.Color("#FFD700"), // Gold
			Success:   lipgloss.Color("#32CD32"), // Lime green
			Warning:   lipgloss.Color("#FF8C00"), // Dark orange
			Error:     lipgloss.Color("#8B0000"), // Dark red
			Info:      lipgloss.Color("#4682B4"), // Steel blue
			Background:    lipgloss.Color("#1A2F1A"), // Dark green
			TextPrimary:   lipgloss.Color("#F0FFF0"), // Honeydew
			TextSecondary: lipgloss.Color("#90EE90"), // Light green
			TextMuted:     lipgloss.Color("#556B2F"), // Dark olive green
			TextAccent:    lipgloss.Color("#FFD700"), // Gold
			Dark1:         lipgloss.Color("#1A2F1A"), // Dark green
			Dark2:         lipgloss.Color("#2A3F2A"), // Medium green
			Dark3:         lipgloss.Color("#3A4F3A"), // Light green
		},
	}

	// Theme 4: Ocean Blue
	OceanBlueTheme = Theme{
		Name:        "Ocean Blue",
		Description: "Calming blue tones inspired by the ocean",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#4682B4"), // Steel blue
			Secondary: lipgloss.Color("#1E90FF"), // Dodger blue
			Accent:    lipgloss.Color("#00CED1"), // Dark turquoise
			Success:   lipgloss.Color("#40E0D0"), // Turquoise
			Warning:   lipgloss.Color("#FFD700"), // Gold
			Error:     lipgloss.Color("#DC143C"), // Crimson
			Info:      lipgloss.Color("#87CEEB"), // Sky blue
			Background:    lipgloss.Color("#0F1F2F"), // Dark blue
			TextPrimary:   lipgloss.Color("#E0F7FA"), // Light cyan
			TextSecondary: lipgloss.Color("#87CEEB"), // Sky blue
			TextMuted:     lipgloss.Color("#4682B4"), // Steel blue
			TextAccent:    lipgloss.Color("#00CED1"), // Dark turquoise
			Dark1:         lipgloss.Color("#0F1F2F"), // Dark blue
			Dark2:         lipgloss.Color("#1F2F3F"), // Medium blue
			Dark3:         lipgloss.Color("#2F3F4F"), // Light blue
		},
	}

	// Theme 5: Sunset Glow
	SunsetGlowTheme = Theme{
		Name:        "Sunset Glow",
		Description: "Warm oranges and pinks inspired by sunset",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#FF6B6B"), // Light coral
			Secondary: lipgloss.Color("#FF8E53"), // Sunset orange
			Accent:    lipgloss.Color("#FEC260"), // Saffron
			Success:   lipgloss.Color("#06FFA5"), // Spring green
			Warning:   lipgloss.Color("#FFB700"), // Amber
			Error:     lipgloss.Color("#C9184A"), // Raspberry
			Info:      lipgloss.Color("#748FFC"), // Periwinkle
			Background:    lipgloss.Color("#2D1B69"), // Deep indigo
			TextPrimary:   lipgloss.Color("#FFEDBC"), // Peach puff
			TextSecondary: lipgloss.Color("#FFB6C1"), // Light pink
			TextMuted:     lipgloss.Color("#FF6B6B"), // Light coral
			TextAccent:    lipgloss.Color("#FEC260"), // Saffron
			Dark1:         lipgloss.Color("#2D1B69"), // Deep indigo
			Dark2:         lipgloss.Color("#3D2B79"), // Medium indigo
			Dark3:         lipgloss.Color("#4D3B89"), // Light indigo
		},
	}

	// Theme 6: Arctic Frost
	ArcticFrostTheme = Theme{
		Name:        "Arctic Frost",
		Description: "Cool blues and whites inspired by arctic landscapes",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#5F9EA0"), // Cadet blue
			Secondary: lipgloss.Color("#B0E0E6"), // Powder blue
			Accent:    lipgloss.Color("#E0FFFF"), // Light cyan
			Success:   lipgloss.Color("#98FB98"), // Pale green
			Warning:   lipgloss.Color("#F0E68C"), // Khaki
			Error:     lipgloss.Color("#CD5C5C"), // Indian red
			Info:      lipgloss.Color("#ADD8E6"), // Light blue
			Background:    lipgloss.Color("#1E2A3A"), // Dark slate blue
			TextPrimary:   lipgloss.Color("#F8F8FF"), // Ghost white
			TextSecondary: lipgloss.Color("#B0C4DE"), // Light steel blue
			TextMuted:     lipgloss.Color("#5F9EA0"), // Cadet blue
			TextAccent:    lipgloss.Color("#E0FFFF"), // Light cyan
			Dark1:         lipgloss.Color("#1E2A3A"), // Dark slate blue
			Dark2:         lipgloss.Color("#2E3A4A"), // Medium slate blue
			Dark3:         lipgloss.Color("#3E4A5A"), // Light slate blue
		},
	}

	// Theme 7: Monochrome
	MonochromeTheme = Theme{
		Name:        "Monochrome",
		Description: "Classic black and white theme",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#FFFFFF"), // White
			Secondary: lipgloss.Color("#CCCCCC"), // Light gray
			Accent:    lipgloss.Color("#FFFFFF"), // White
			Success:   lipgloss.Color("#00FF00"), // Green
			Warning:   lipgloss.Color("#FFFF00"), // Yellow
			Error:     lipgloss.Color("#FF0000"), // Red
			Info:      lipgloss.Color("#00FFFF"), // Cyan
			Background:    lipgloss.Color("#000000"), // Black
			TextPrimary:   lipgloss.Color("#FFFFFF"), // White
			TextSecondary: lipgloss.Color("#CCCCCC"), // Light gray
			TextMuted:     lipgloss.Color("#666666"), // Gray
			TextAccent:    lipgloss.Color("#FFFFFF"), // White
			Dark1:         lipgloss.Color("#000000"), // Black
			Dark2:         lipgloss.Color("#1A1A1A"), // Dark gray
			Dark3:         lipgloss.Color("#333333"), // Medium gray
		},
	}

	// Theme 8: Royal Purple
	RoyalPurpleTheme = Theme{
		Name:        "Royal Purple",
		Description: "Rich purples and golds inspired by royalty",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#7851A9"), // Medium purple
			Secondary: lipgloss.Color("#6A0DAD"), // Purple
			Accent:    lipgloss.Color("#FFD700"), // Gold
			Success:   lipgloss.Color("#9370DB"), // Medium purple
			Warning:   lipgloss.Color("#DAA520"), // Goldenrod
			Error:     lipgloss.Color("#8B008B"), // Dark magenta
			Info:      lipgloss.Color("#87CEEB"), // Sky blue
			Background:    lipgloss.Color("#2E0854"), // Dark slate magenta
			TextPrimary:   lipgloss.Color("#F5F5DC"), // Beige
			TextSecondary: lipgloss.Color("#DDA0DD"), // Plum
			TextMuted:     lipgloss.Color("#7851A9"), // Medium purple
			TextAccent:    lipgloss.Color("#FFD700"), // Gold
			Dark1:         lipgloss.Color("#2E0854"), // Dark slate magenta
			Dark2:         lipgloss.Color("#3E1864"), // Medium magenta
			Dark3:         lipgloss.Color("#4E2874"), // Light magenta
		},
	}

	// Theme 9: Cyberpunk
	CyberpunkTheme = Theme{
		Name:        "Cyberpunk",
		Description: "Futuristic theme with electric colors",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#00FFFF"), // Cyan
			Secondary: lipgloss.Color("#FF00FF"), // Magenta
			Accent:    lipgloss.Color("#FFFF00"), // Yellow
			Success:   lipgloss.Color("#00FF00"), // Green
			Warning:   lipgloss.Color("#FF4500"), // Orange red
			Error:     lipgloss.Color("#FF0000"), // Red
			Info:      lipgloss.Color("#1E90FF"), // Dodger blue
			Background:    lipgloss.Color("#0D0221"), // Navy
			TextPrimary:   lipgloss.Color("#F0F8FF"), // Alice blue
			TextSecondary: lipgloss.Color("#00FFFF"), // Cyan
			TextMuted:     lipgloss.Color("#8B008B"), // Dark magenta
			TextAccent:    lipgloss.Color("#FFFF00"), // Yellow
			Dark1:         lipgloss.Color("#0D0221"), // Navy
			Dark2:         lipgloss.Color("#1D1231"), // Dark purple
			Dark3:         lipgloss.Color("#2D2241"), // Medium purple
		},
	}

	// Theme 10: Coffee Shop
	CoffeeShopTheme = Theme{
		Name:        "Coffee Shop",
		Description: "Warm browns and creams inspired by coffee",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#6F4E37"), // Coffee
			Secondary: lipgloss.Color("#8B4513"), // Saddle brown
			Accent:    lipgloss.Color("#D2691E"), // Chocolate
			Success:   lipgloss.Color("#A0522D"), // Sienna
			Warning:   lipgloss.Color("#CD853F"), // Peru
			Error:     lipgloss.Color("#8B4513"), // Saddle brown
			Info:      lipgloss.Color("#DEB887"), // Burlywood
			Background:    lipgloss.Color("#3C2415"), // Dark coffee
			TextPrimary:   lipgloss.Color("#FFF8DC"), // Cornsilk
			TextSecondary: lipgloss.Color("#F5DEB3"), // Wheat
			TextMuted:     lipgloss.Color("#A0522D"), // Sienna
			TextAccent:    lipgloss.Color("#D2691E"), // Chocolate
			Dark1:         lipgloss.Color("#3C2415"), // Dark coffee
			Dark2:         lipgloss.Color("#4C3425"), // Medium coffee
			Dark3:         lipgloss.Color("#5C4435"), // Light coffee
		},
	}

	// Theme 11: Desert Sunset
	DesertSunsetTheme = Theme{
		Name:        "Desert Sunset",
		Description: "Warm oranges and reds inspired by desert landscapes",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#FF7F50"), // Coral
			Secondary: lipgloss.Color("#FF6347"), // Tomato
			Accent:    lipgloss.Color("#FF4500"), // Orange red
			Success:   lipgloss.Color("#FFD700"), // Gold
			Warning:   lipgloss.Color("#FF8C00"), // Dark orange
			Error:     lipgloss.Color("#B22222"), // Fire brick
			Info:      lipgloss.Color("#F0E68C"), // Khaki
			Background:    lipgloss.Color("#4A2C2A"), // Dark red
			TextPrimary:   lipgloss.Color("#FFF5EE"), // Seashell
			TextSecondary: lipgloss.Color("#FFE4B5"), // Moccasin
			TextMuted:     lipgloss.Color("#CD5C5C"), // Indian red
			TextAccent:    lipgloss.Color("#FF4500"), // Orange red
			Dark1:         lipgloss.Color("#4A2C2A"), // Dark red
			Dark2:         lipgloss.Color("#5A3C3A"), // Medium red
			Dark3:         lipgloss.Color("#6A4C4A"), // Light red
		},
	}

	// Theme 12: Zen Garden
	ZenGardenTheme = Theme{
		Name:        "Zen Garden",
		Description: "Calming greens and grays inspired by zen gardens",
		Colors: ThemeColors{
			Primary:   lipgloss.Color("#808080"), // Gray
			Secondary: lipgloss.Color("#A9A9A9"), // Dark gray
			Accent:    lipgloss.Color("#90EE90"), // Light green
			Success:   lipgloss.Color("#228B22"), // Forest green
			Warning:   lipgloss.Color("#DAA520"), // Goldenrod
			Error:     lipgloss.Color("#8B4513"), // Saddle brown
			Info:      lipgloss.Color("#B0C4DE"), // Light steel blue
			Background:    lipgloss.Color("#2F4F4F"), // Dark slate gray
			TextPrimary:   lipgloss.Color("#F5F5F5"), // White smoke
			TextSecondary: lipgloss.Color("#D3D3D3"), // Light gray
			TextMuted:     lipgloss.Color("#808080"), // Gray
			TextAccent:    lipgloss.Color("#90EE90"), // Light green
			Dark1:         lipgloss.Color("#2F4F4F"), // Dark slate gray
			Dark2:         lipgloss.Color("#3F5F5F"), // Medium gray
			Dark3:         lipgloss.Color("#4F6F6F"), // Light gray
		},
	}
)

// AllThemes contains all predefined themes
var AllThemes = []Theme{
	MidnightJazzTheme,
	NeonDreamsTheme,
	ForestRetreatTheme,
	OceanBlueTheme,
	SunsetGlowTheme,
	ArcticFrostTheme,
	MonochromeTheme,
	RoyalPurpleTheme,
	CyberpunkTheme,
	CoffeeShopTheme,
	DesertSunsetTheme,
	ZenGardenTheme,
}