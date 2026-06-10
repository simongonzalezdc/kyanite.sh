package design

import "github.com/charmbracelet/lipgloss"

func registerThemes() {
	RegisterBuiltIn(Theme{
		Name:       "monochrome",
		Primary:    lipgloss.Color("#E0E0E0"),
		Secondary:  lipgloss.Color("#909090"),
		Accent:     lipgloss.Color("#FFFFFF"),
		Background: lipgloss.Color("#0D0D0D"),
		Text:       lipgloss.Color("#E8E8E8"),
		Success:    lipgloss.Color("#A8D8A8"),
		Warning:    lipgloss.Color("#D4C090"),
		Error:      lipgloss.Color("#E08080"),
		Border:     lipgloss.Color("#606060"),
		Panel:      lipgloss.Color("#1A1A1A"),
		Muted:      lipgloss.Color("#707070"),
	})

	RegisterBuiltIn(Theme{
		Name:       "amber-night",
		Primary:    lipgloss.Color("#D4A574"),
		Secondary:  lipgloss.Color("#9D84B7"),
		Accent:     lipgloss.Color("#F4D03F"),
		Background: lipgloss.Color("#0A0E27"),
		Text:       lipgloss.Color("#E8DFF5"),
		Success:    lipgloss.Color("#52D3AA"),
		Warning:    lipgloss.Color("#FFA502"),
		Error:      lipgloss.Color("#FF3333"),
		Border:     lipgloss.Color("#9D84B7"),
		Panel:      lipgloss.Color("#15182F"),
		Muted:      lipgloss.Color("#7B6F99"),
	})

	RegisterBuiltIn(Theme{
		Name:       "indigo-depths",
		Primary:    lipgloss.Color("#7B9EFF"),
		Secondary:  lipgloss.Color("#5F9EA0"),
		Accent:     lipgloss.Color("#87CEEB"),
		Background: lipgloss.Color("#0A0A1A"),
		Text:       lipgloss.Color("#E6F2FF"),
		Success:    lipgloss.Color("#52D3AA"),
		Warning:    lipgloss.Color("#FFB84D"),
		Error:      lipgloss.Color("#FF6B6B"),
		Border:     lipgloss.Color("#5F9EA0"),
		Panel:      lipgloss.Color("#14142A"),
		Muted:      lipgloss.Color("#6B7FA8"),
	})

	RegisterBuiltIn(Theme{
		Name:       "forest-path",
		Primary:    lipgloss.Color("#8FBC8F"),
		Secondary:  lipgloss.Color("#6B8E6B"),
		Accent:     lipgloss.Color("#B4D7B4"),
		Background: lipgloss.Color("#0E1A0E"),
		Text:       lipgloss.Color("#E8F5E8"),
		Success:    lipgloss.Color("#90EE90"),
		Warning:    lipgloss.Color("#DAA520"),
		Error:      lipgloss.Color("#CD5C5C"),
		Border:     lipgloss.Color("#6B8E6B"),
		Panel:      lipgloss.Color("#162216"),
		Muted:      lipgloss.Color("#5A7A5A"),
	})

	RegisterBuiltIn(Theme{
		Name:       "cyan-wave",
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
		Muted:      lipgloss.Color("#4A7A8C"),
	})

	RegisterBuiltIn(Theme{
		Name:       "electric-rose",
		Primary:    lipgloss.Color("#FF1493"),
		Secondary:  lipgloss.Color("#C71585"),
		Accent:     lipgloss.Color("#00CED1"),
		Background: lipgloss.Color("#140A14"),
		Text:       lipgloss.Color("#FFF0F8"),
		Success:    lipgloss.Color("#52D3AA"),
		Warning:    lipgloss.Color("#FFB84D"),
		Error:      lipgloss.Color("#FF4444"),
		Border:     lipgloss.Color("#C71585"),
		Panel:      lipgloss.Color("#230F23"),
		Muted:      lipgloss.Color("#8C4A7A"),
	})

	RegisterBuiltIn(Theme{
		Name:       "twilight-mist",
		Primary:    lipgloss.Color("#B8A3C9"),
		Secondary:  lipgloss.Color("#8E7B9D"),
		Accent:     lipgloss.Color("#D4C5E0"),
		Background: lipgloss.Color("#151520"),
		Text:       lipgloss.Color("#E8E4F0"),
		Success:    lipgloss.Color("#90C695"),
		Warning:    lipgloss.Color("#C9A87C"),
		Error:      lipgloss.Color("#C77777"),
		Border:     lipgloss.Color("#8E7B9D"),
		Panel:      lipgloss.Color("#1A1A2A"),
		Muted:      lipgloss.Color("#7A6E8D"),
	})

	RegisterBuiltIn(Theme{
		Name:       "clay-earth",
		Primary:    lipgloss.Color("#C87840"),
		Secondary:  lipgloss.Color("#8B4513"),
		Accent:     lipgloss.Color("#DEB887"),
		Background: lipgloss.Color("#1A1410"),
		Text:       lipgloss.Color("#F5E6D3"),
		Success:    lipgloss.Color("#8FBC8F"),
		Warning:    lipgloss.Color("#CD853F"),
		Error:      lipgloss.Color("#CD5C5C"),
		Border:     lipgloss.Color("#8B4513"),
		Panel:      lipgloss.Color("#221A12"),
		Muted:      lipgloss.Color("#8B7355"),
	})

	RegisterBuiltIn(Theme{
		Name:       "iron-forge",
		Primary:    lipgloss.Color("#F04060"),
		Secondary:  lipgloss.Color("#4A4A4A"),
		Accent:     lipgloss.Color("#FF6347"),
		Background: lipgloss.Color("#1A0A0A"),
		Text:       lipgloss.Color("#FFE6E6"),
		Success:    lipgloss.Color("#90C695"),
		Warning:    lipgloss.Color("#FFB84D"),
		Error:      lipgloss.Color("#FF4444"),
		Border:     lipgloss.Color("#4A4A4A"),
		Panel:      lipgloss.Color("#241010"),
		Muted:      lipgloss.Color("#7A5050"),
	})

	RegisterBuiltIn(Theme{
		Name:       "sunlight",
		Primary:    lipgloss.Color("#FFD700"),
		Secondary:  lipgloss.Color("#DAA520"),
		Accent:     lipgloss.Color("#FFF8DC"),
		Background: lipgloss.Color("#0F0F0A"),
		Text:       lipgloss.Color("#FFFACD"),
		Success:    lipgloss.Color("#98D982"),
		Warning:    lipgloss.Color("#FF9800"),
		Error:      lipgloss.Color("#E84040"),
		Border:     lipgloss.Color("#DAA520"),
		Panel:      lipgloss.Color("#1A1A10"),
		Muted:      lipgloss.Color("#8B8B5A"),
	})
}
