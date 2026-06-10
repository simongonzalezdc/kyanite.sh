package design

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestAllThemesRegister(t *testing.T) {
	themes := List()
	if len(themes) != 10 {
		t.Fatalf("expected 10 themes, got %d: %v", len(themes), themes)
	}
	expected := []string{"monochrome", "amber-night", "indigo-depths", "forest-path", "cyan-wave", "electric-rose", "twilight-mist", "clay-earth", "iron-forge", "sunlight"}
	for _, name := range expected {
		th := Get(name)
		if th.Name != name {
			t.Errorf("theme %q not found or has wrong name", name)
		}
	}
}

func TestRegisterBuiltInPanicsOnBadTheme(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic for failing WCAG AA theme")
		}
	}()
	RegisterBuiltIn(Theme{
		Name:       "bad-theme",
		Primary:    lipgloss.Color("#000000"),
		Background: lipgloss.Color("#000000"),
		Text:       lipgloss.Color("#010101"),
		Success:    lipgloss.Color("#020202"),
		Warning:    lipgloss.Color("#030303"),
		Error:      lipgloss.Color("#040404"),
	})
}

func TestRegisterCustomReturnsError(t *testing.T) {
	err := RegisterCustom(Theme{
		Name:       "custom-bad",
		Primary:    lipgloss.Color("#000000"),
		Background: lipgloss.Color("#000000"),
		Text:       lipgloss.Color("#010101"),
		Success:    lipgloss.Color("#020202"),
		Warning:    lipgloss.Color("#030303"),
		Error:      lipgloss.Color("#040404"),
	})
	if err == nil {
		t.Error("expected error for failing WCAG AA custom theme")
	}
}

func TestMustGetPanicsOnMissing(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic for missing theme")
		}
	}()
	MustGet("nonexistent")
}

func TestDefaultTheme(t *testing.T) {
	th := DefaultTheme()
	if th.Name != "amber-night" {
		t.Errorf("expected amber-night, got %s", th.Name)
	}
}

func TestContrastRatio(t *testing.T) {
	ratio, err := ContrastRatio("#FFFFFF", "#000000")
	if err != nil {
		t.Fatal(err)
	}
	if ratio < 20.0 {
		t.Errorf("white on black should be >= 21:1, got %.2f", ratio)
	}
}
