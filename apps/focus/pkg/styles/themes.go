package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/theme"
)

// Theme modes for the app (backward compatibility)
type ThemeMode string

const (
	ThemeSynthwave ThemeMode = "synthwave"
	ThemeLight     ThemeMode = "light"
	ThemePlain     ThemeMode = "plain"
)

// Current theme registry (Kyanite themes)
var currentTheme = theme.Default()

// Light theme colors - research-based best practices for light themes
var (
	LightBg        = lipgloss.Color("#FAFAFA") // Warm white, not harsh
	LightFg        = lipgloss.Color("#2D3748") // Dark blue-gray for readability
	LightAccent    = lipgloss.Color("#4299E1") // Blue accent
	LightSecondary = lipgloss.Color("#805AD5") // Purple accent
	LightBorder    = lipgloss.Color("#CBD5E0") // Light gray border
	LightSuccess   = lipgloss.Color("#38A169") // Muted green
	LightWarning   = lipgloss.Color("#ED8936") // Warm orange
	LightError     = lipgloss.Color("#E53E3E") // Soft red
	LightPanel     = lipgloss.Color("#F7FAFC") // Very light gray for panels
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

// SetThemeByName changes the current theme by display name or id
func SetThemeByName(name string) {
	// Try by display name first
	t := theme.GetThemeByName(name)
	if t.Name == name {
		currentTheme = t
		return
	}
	// Fallback to id (with migration)
	currentTheme = theme.GetTheme(name)
}

// SetTheme changes the current theme (backward compatibility)
func SetTheme(t ThemeMode) {
	switch t {
	case ThemeLight:
		currentTheme = theme.Theme{
			Name:       "Light",
			Background: LightBg,
			Text:       LightFg,
			Accent:     LightAccent,
			Border:     LightBorder,
			Success:    LightSuccess,
			Warning:    LightWarning,
			Error:      LightError,
			Panel:      LightPanel,
		}
	case ThemePlain:
		currentTheme = theme.Theme{
			Name:       "Plain",
			Background: PlainBg,
			Text:       PlainFg,
			Accent:     PlainAccent,
			Border:     PlainBorder,
			Success:    PlainSuccess,
			Warning:    PlainWarning,
			Error:      PlainError,
			Panel:      PlainBg,
		}
	default:
		currentTheme = theme.Default()
	}
}

// GetTheme returns the current theme
func GetTheme() theme.Theme {
	return currentTheme
}

// CycleTheme cycles through all Kyanite themes
func CycleTheme() {
	m := theme.GetManager()
	m.Next()
	currentTheme = m.Current()
}

// Theme-aware color getters
func GetBackground() lipgloss.Color {
	return currentTheme.Background
}

func GetForeground() lipgloss.Color {
	return currentTheme.Text
}

func GetAccent() lipgloss.Color {
	return currentTheme.Accent
}

func GetBorder() lipgloss.Color {
	return currentTheme.Border
}

func GetSuccess() lipgloss.Color {
	return currentTheme.Success
}

func GetWarning() lipgloss.Color {
	return currentTheme.Warning
}

func GetError() lipgloss.Color {
	return currentTheme.Error
}

func GetPanel() lipgloss.Color {
	return currentTheme.Panel
}

// Theme-aware style getters
func GetTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.Primary).
		Background(currentTheme.Background).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(currentTheme.Accent)
}

func GetBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.Text).
		Background(currentTheme.Panel).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(currentTheme.Border).
		Padding(1)
}

func GetPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.Text).
		Background(currentTheme.Panel).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(currentTheme.Border).
		Padding(1)
}

func GetSuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.Success).
		Bold(true)
}

func GetErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.Error).
		Bold(true)
}

func GetWarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.Warning).
		Bold(true)
}

// Update existing styles to be theme-aware
func GetSynthwaveTitle() lipgloss.Style {
	return GetTitleStyle()
}

func GetFocusBox() lipgloss.Style {
	return GetBoxStyle()
}

// GetThemeList returns all available theme names for CLI commands
func GetThemeList() []string {
	return theme.GetThemeNames()
}

// GetCurrentThemeName returns the name of the current theme
func GetCurrentThemeName() string {
	return currentTheme.Name
}
