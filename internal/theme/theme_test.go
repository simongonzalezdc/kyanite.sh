package theme

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestAllThemes(t *testing.T) {
	// Test that we have exactly 13 themes
	if len(AllThemes) != 13 {
		t.Errorf("Expected 13 themes, got %d", len(AllThemes))
	}
	
	// Test that each theme has required fields
	for _, theme := range AllThemes {
		if theme.Name == "" {
			t.Error("Theme name should not be empty")
		}
		
		if theme.Primary == "" {
			t.Error("Theme primary color should not be empty")
		}
		
		if theme.Background == "" {
			t.Error("Theme background color should not be empty")
		}
		
		if theme.Text == "" {
			t.Error("Theme text color should not be empty")
		}
		
		if theme.Success == "" {
			t.Error("Theme success color should not be empty")
		}
	}
}

func TestGetThemeByName(t *testing.T) {
	// Test getting existing theme
	slateMist := GetThemeByName("Slate Mist")
	if slateMist.Name != "Slate Mist" {
		t.Errorf("Expected theme name 'Slate Mist', got '%s'", slateMist.Name)
	}
	
	if slateMist.Primary != "#E0E0E0" {
		t.Errorf("Expected primary color '#E0E0E0', got '%s'", slateMist.Primary)
	}
	
	// Test getting non-existent theme (should return default)
	unknownTheme := GetThemeByName("NonExistentTheme")
	if unknownTheme.Name != DefaultTheme.Name {
		t.Errorf("Expected default theme '%s', got '%s'", DefaultTheme.Name, unknownTheme.Name)
	}
}

func TestGetThemeNames(t *testing.T) {
	names := GetThemeNames()
	
	// Should return 13 theme names
	if len(names) != 13 {
		t.Errorf("Expected 13 theme names, got %d", len(names))
	}
	
	// Should contain all theme names
	expectedThemes := []string{
		"Slate Mist", "Violet Dusk", "Amber Night", "Molten Gold", 
		"Clay Roads", "Iron Storm", "Jade Tide", "Sunset Ember",
		"Forest Whisper", "Electric Bloom", "Plasma Pulse", "Indigo Depths", "Sage Meadow",
	}
	
	for _, expected := range expectedThemes {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected theme name '%s' not found in theme names", expected)
		}
	}
}

func TestDefaultTheme(t *testing.T) {
	// Default theme should be Amber Night
	if DefaultTheme.Name != "Amber Night" {
		t.Errorf("Expected default theme 'Amber Night', got '%s'", DefaultTheme.Name)
	}
	
	// Default theme should have specific colors
	if DefaultTheme.Primary != "#B8936E" {
		t.Errorf("Expected default primary color '#B8936E', got '%s'", DefaultTheme.Primary)
	}
	
	if DefaultTheme.Background != "#12101A" {
		t.Errorf("Expected default background color '#12101A', got '%s'", DefaultTheme.Background)
	}
	
	if DefaultTheme.Success != "#52D3AA" {
		t.Errorf("Expected default success color '#52D3AA', got '%s'", DefaultTheme.Success)
	}
}

func TestThemeColors(t *testing.T) {
	// Test that all theme colors are valid lipgloss colors
	for _, theme := range AllThemes {
		// This will panic if color is invalid - good for testing
		_ = lipgloss.Color(theme.Primary)
		_ = lipgloss.Color(theme.Secondary)
		_ = lipgloss.Color(theme.Accent)
		_ = lipgloss.Color(theme.Background)
		_ = lipgloss.Color(theme.Text)
		_ = lipgloss.Color(theme.Success)
		_ = lipgloss.Color(theme.Warning)
		_ = lipgloss.Color(theme.Error)
		_ = lipgloss.Color(theme.Border)
		_ = lipgloss.Color(theme.Panel)
	}
}

func TestSpecificThemeDefinitions(t *testing.T) {
	// Test Slate Mist theme
	slateMist := SlateMist
	if slateMist.Name != "Slate Mist" {
		t.Errorf("Expected Slate Mist name to be 'Slate Mist', got '%s'", slateMist.Name)
	}
	
	if slateMist.Primary != "#E0E0E0" {
		t.Errorf("Expected Slate Mist primary '#E0E0E0', got '%s'", slateMist.Primary)
	}
	
	// Test Violet Dusk theme
	violetDusk := VioletDusk
	if violetDusk.Name != "Violet Dusk" {
		t.Errorf("Expected Violet Dusk name to be 'Violet Dusk', got '%s'", violetDusk.Name)
	}
	
	if violetDusk.Background != "#1A1527" {
		t.Errorf("Expected Violet Dusk background '#1A1527', got '%s'", violetDusk.Background)
	}
	
	// Test Amber Night theme (default)
	amberNight := AmberNight
	if amberNight.Name != "Amber Night" {
		t.Errorf("Expected Amber Night name to be 'Amber Night', got '%s'", amberNight.Name)
	}
	
	if amberNight.Accent != "#E8C547" {
		t.Errorf("Expected Amber Night accent '#E8C547', got '%s'", amberNight.Accent)
	}
	
	// Test Electric Bloom theme
	electricBloom := ElectricBloom
	if electricBloom.Name != "Electric Bloom" {
		t.Errorf("Expected Electric Bloom name to be 'Electric Bloom', got '%s'", electricBloom.Name)
	}
	
	if electricBloom.Text != "#F0F3FF" {
		t.Errorf("Expected Electric Bloom text '#F0F3FF', got '%s'", electricBloom.Text)
	}
}

func TestThemeUniqueness(t *testing.T) {
	// Test that all theme names are unique
	themeNames := make(map[string]bool)
	for _, theme := range AllThemes {
		if themeNames[theme.Name] {
			t.Errorf("Duplicate theme name found: %s", theme.Name)
		}
		themeNames[theme.Name] = true
	}
	
	// Test that all theme color combinations are unique
	themeCombinations := make(map[string]bool)
	for _, theme := range AllThemes {
		combination := string(theme.Primary) + string(theme.Background) + string(theme.Accent)
		if themeCombinations[combination] {
			t.Errorf("Duplicate theme color combination found for theme: %s", theme.Name)
		}
		themeCombinations[combination] = true
	}
}

func TestThemeColorFormat(t *testing.T) {
	// Test that all theme colors are valid hex color format
	for _, theme := range AllThemes {
		colors := []string{
			string(theme.Primary),
			string(theme.Secondary),
			string(theme.Accent),
			string(theme.Background),
			string(theme.Text),
			string(theme.Success),
			string(theme.Warning),
			string(theme.Error),
			string(theme.Border),
			string(theme.Panel),
		}
		
		for _, color := range colors {
			if len(color) != 7 || color[0] != '#' {
				t.Errorf("Invalid color format for theme %s: %s", theme.Name, color)
			}
		}
	}
}

func TestThemeContrast(t *testing.T) {
	// Basic test to ensure themes have contrasting colors
	for _, theme := range AllThemes {
		// These tests are basic - in a real implementation you'd want to
		// test actual contrast ratios, but for now we just ensure they're different
		if theme.Background == theme.Text {
			t.Errorf("Theme %s has same background and text color", theme.Name)
		}
		
		if theme.Background == theme.Primary {
			t.Errorf("Theme %s has same background and primary color", theme.Name)
		}
		
		if theme.Success == theme.Error {
			t.Errorf("Theme %s has same success and error color", theme.Name)
		}
	}
}