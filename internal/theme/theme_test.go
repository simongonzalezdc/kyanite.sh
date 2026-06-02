package theme

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestRegistryCount(t *testing.T) {
	if len(Registry) != 10 {
		t.Errorf("Expected 10 themes, got %d", len(Registry))
	}
}

func TestDefaultTheme(t *testing.T) {
	def := Default()
	if def.Name != "Amber Night" {
		t.Errorf("Expected default theme 'Amber Night', got '%s'", def.Name)
	}
}

func TestGetThemeByIDAndName(t *testing.T) {
	// By id
	t1 := GetTheme("amber-night")
	if t1.Name != "Amber Night" {
		t.Errorf("ID amber-night should map to Amber Night, got %s", t1.Name)
	}
	// By name
	t2 := GetThemeByName("Amber Night")
	if t2.Name != "Amber Night" {
		t.Errorf("Name lookup failed, got %s", t2.Name)
	}
	// Unknown -> default
	unk := GetTheme("does-not-exist")
	if unk.Name != Default().Name {
		t.Errorf("Unknown id should return default, got %s", unk.Name)
	}
}

func TestGetThemeNames(t *testing.T) {
	names := GetThemeNames()
	if len(names) != 10 {
		t.Errorf("Expected 10 theme names, got %d", len(names))
	}
	found := false
	for _, n := range names {
		if n == "Amber Night" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Amber Night should be in theme names")
	}
}

func TestThemeColorsFormatAndNonEmpty(t *testing.T) {
	for _, th := range Registry {
		colors := []string{
			string(th.Primary), string(th.Secondary), string(th.Accent),
			string(th.Background), string(th.Text), string(th.Success),
			string(th.Warning), string(th.Error), string(th.Border), string(th.Panel),
		}
		for _, c := range colors {
			if len(c) == 0 {
				t.Errorf("Empty color in theme %s", th.Name)
			}
			if len(c) != 7 || c[0] != '#' {
				t.Errorf("Invalid color format in theme %s: %s", th.Name, c)
			}
		}
		// Ensure lipgloss accepts them
		_ = lipgloss.Color(th.Primary)
		_ = lipgloss.Color(th.Background)
	}
}
