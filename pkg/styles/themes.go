package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme modes for the app
type ThemeMode string

const (
	ThemeSynthwave ThemeMode = "synthwave"
	ThemeLight     ThemeMode = "light"
	ThemePlain     ThemeMode = "plain"
)

// Current theme (can be changed at runtime)
var currentTheme ThemeMode = ThemeSynthwave

// Light theme colors - research-based best practices for light themes
var (
	LightBg        = lipgloss.Color("#FAFAFA")      // Warm white, not harsh
	LightFg        = lipgloss.Color("#2D3748")      // Dark blue-gray for readability
	LightAccent    = lipgloss.Color("#4299E1")      // Blue accent
	LightSecondary = lipgloss.Color("#805AD5")      // Purple accent
	LightBorder    = lipgloss.Color("#CBD5E0")      // Light gray border
	LightSuccess   = lipgloss.Color("#38A169")      // Muted green
	LightWarning   = lipgloss.Color("#ED8936")      // Warm orange
	LightError     = lipgloss.Color("#E53E3E")      // Soft red
	LightPanel     = lipgloss.Color("#F7FAFC")      // Very light gray for panels
)

// Plain theme colors (terminal default)
var (
	PlainBg        = lipgloss.Color("")
	PlainFg        = lipgloss.Color("")
	PlainAccent    = lipgloss.Color("")
	PlainSecondary = lipgloss.Color("")
	PlainBorder    = lipgloss.Color("")
	PlainSuccess   = lipgloss.Color("")
	PlainWarning   = lipgloss.Color("")
	PlainError     = lipgloss.Color("")
)

// SetTheme changes the current theme
func SetTheme(theme ThemeMode) {
	currentTheme = theme
}

// GetTheme returns the current theme
func GetTheme() ThemeMode {
	return currentTheme
}

// Theme-aware color getters
func GetBackground() lipgloss.Color {
	switch currentTheme {
	case ThemeLight:
		return LightBg
	case ThemePlain:
		return PlainBg
	default:
		return DeepSpace
	}
}

func GetForeground() lipgloss.Color {
	switch currentTheme {
	case ThemeLight:
		return LightFg
	case ThemePlain:
		return PlainFg
	default:
		return SynthwavePink
	}
}

func GetAccent() lipgloss.Color {
	switch currentTheme {
	case ThemeLight:
		return LightAccent
	case ThemePlain:
		return PlainAccent
	default:
		return SynthwaveCyan
	}
}

func GetBorder() lipgloss.Color {
	switch currentTheme {
	case ThemeLight:
		return LightBorder
	case ThemePlain:
		return PlainBorder
	default:
		return SynthwaveCyan
	}
}

func GetSuccess() lipgloss.Color {
	switch currentTheme {
	case ThemeLight:
		return LightSuccess
	case ThemePlain:
		return PlainSuccess
	default:
		return SynthwaveGreen
	}
}

func GetWarning() lipgloss.Color {
	switch currentTheme {
	case ThemeLight:
		return LightWarning
	case ThemePlain:
		return PlainWarning
	default:
		return SynthwaveYellow
	}
}

func GetError() lipgloss.Color {
	switch currentTheme {
	case ThemeLight:
		return LightError
	case ThemePlain:
		return PlainError
	default:
		return SynthwaveRed
	}
}

// Theme-aware style getters
func GetTitleStyle() lipgloss.Style {
	switch currentTheme {
	case ThemeLight:
		return lipgloss.NewStyle().
			Foreground(LightFg).
			Background(LightBg).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(LightAccent)
	case ThemePlain:
		return lipgloss.NewStyle().
			Bold(true)
	default:
		return lipgloss.NewStyle().
			Foreground(SynthwavePink).
			Background(DeepSpace).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(SynthwaveCyan)
	}
}

func GetBoxStyle() lipgloss.Style {
	switch currentTheme {
	case ThemeLight:
		return lipgloss.NewStyle().
			Foreground(LightFg).
			Background(LightPanel).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(LightBorder).
			Padding(1)
	case ThemePlain:
		return lipgloss.NewStyle()
	default:
		return lipgloss.NewStyle().
			Foreground(SynthwaveCyan).
			Background(DarkVoid).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(SynthwaveCyan).
			Padding(1)
	}
}

func GetPanelStyle() lipgloss.Style {
	switch currentTheme {
	case ThemeLight:
		return lipgloss.NewStyle().
			Foreground(LightFg).
			Background(LightPanel).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(LightBorder).
			Padding(1)
	case ThemePlain:
		return lipgloss.NewStyle()
	default:
		return lipgloss.NewStyle().
			Foreground(SynthwaveCyan).
			Background(DarkVoid).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SynthwaveCyan).
			Padding(1)
	}
}

func GetSuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(GetSuccess()).
		Bold(true)
}

func GetErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(GetError()).
		Bold(true)
}

func GetWarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(GetWarning()).
		Bold(true)
}

// Update existing styles to be theme-aware
func GetSynthwaveTitle() lipgloss.Style {
	return GetTitleStyle()
}

func GetFocusBox() lipgloss.Style {
	return GetBoxStyle()
}
