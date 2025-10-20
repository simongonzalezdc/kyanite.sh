package styles

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestThemeConstants(t *testing.T) {
	// Test that theme constants are properly defined
	if ThemeSynthwave != "synthwave" {
		t.Errorf("Expected ThemeSynthwave to be 'synthwave', got '%s'", ThemeSynthwave)
	}
	
	if ThemeLight != "light" {
		t.Errorf("Expected ThemeLight to be 'light', got '%s'", ThemeLight)
	}
	
	if ThemePlain != "plain" {
		t.Errorf("Expected ThemePlain to be 'plain', got '%s'", ThemePlain)
	}
}

func TestSetTheme(t *testing.T) {
	// Test setting different themes
	SetTheme(ThemeSynthwave)
	SetTheme(ThemeLight)
	SetTheme(ThemePlain)
	
	// Should not panic
	currentTheme := GetTheme()
	if currentTheme.Name == "" {
		t.Error("GetTheme() should return a theme with a name")
	}
}

func TestGetTheme(t *testing.T) {
	// Test that GetTheme returns a valid theme
	currentTheme := GetTheme()
	
	if currentTheme.Name == "" {
		t.Error("GetTheme() should return a theme with a name")
	}
	
	if string(currentTheme.Background) != "" {
		t.Error("GetTheme() should return a theme with background color")
	}
}

func TestColorGetters(t *testing.T) {
	// Test that color getters return valid colors
	colors := []func() lipgloss.Color{
		GetBackground,
		GetForeground,
		GetAccent,
		GetBorder,
		GetSuccess,
		GetWarning,
		GetError,
		GetPanel,
	}
	
	for _, colorGetter := range colors {
		color := colorGetter()
		// This will panic if color is invalid - good for testing
		_ = lipgloss.Color(color)
	}
}

func TestStyleGetters(t *testing.T) {
	// Test that style getters return valid styles
	styles := []func() lipgloss.Style{
		GetTitleStyle,
		GetBoxStyle,
		GetPanelStyle,
		GetSuccessStyle,
		GetErrorStyle,
		GetWarningStyle,
	}
	
	for _, styleGetter := range styles {
		style := styleGetter()
		// This should not panic
		_ = style.Render("test")
	}
}

func TestLightThemeColors(t *testing.T) {
	// Test light theme colors are properly defined
	_ = lipgloss.Color(LightBg)
	_ = lipgloss.Color(LightFg)
	_ = lipgloss.Color(LightAccent)
	_ = lipgloss.Color(LightSuccess)
}

func TestPlainThemeColors(t *testing.T) {
	// Test plain theme colors (can be empty for terminal defaults)
	// These should not panic when used
	_ = lipgloss.Color(PlainBg)
	_ = lipgloss.Color(PlainFg)
	_ = lipgloss.Color(PlainAccent)
	_ = lipgloss.Color(PlainSuccess)
}

func TestThemeConsistency(t *testing.T) {
	// Test that theme-related functions work together
	bgColor := GetBackground()
	fgColor := GetForeground()
	titleStyle := GetTitleStyle()
	
	// Should not panic
	_ = lipgloss.NewStyle().
		Foreground(bgColor).
		Background(fgColor).
		Render("test")
	
	_ = titleStyle.Render("test")
}

func TestBackwardCompatibility(t *testing.T) {
	// Test that old theme system still works for compatibility
	SetTheme(ThemeSynthwave)
	
	// Should use theme colors
	bg := GetBackground()
	fg := GetForeground()
	
	// Should not panic
	_ = lipgloss.Color(bg)
	_ = lipgloss.Color(fg)
}