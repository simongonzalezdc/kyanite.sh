package theme

import (
	"testing"

	"github.com/kyanite/design"
)

func TestManagerCurrent(t *testing.T) {
	m := GetManager()
	cur := m.Current()
	if cur.Name == "" {
		t.Error("Current theme should have a name")
	}
}

func TestManagerSetTheme(t *testing.T) {
	m := GetManager()
	m.SetTheme("electric-rose")
	cur := m.Current()
	if cur.Name != "electric-rose" {
		t.Errorf("Expected electric-rose, got %s", cur.Name)
	}
	// Reset
	m.SetTheme("amber-night")
}

func TestManagerNext(t *testing.T) {
	m := GetManager()
	m.SetTheme("amber-night")
	next := m.Next()
	if next.Name == "" {
		t.Error("Next theme should have a name")
	}
	// Reset
	m.SetTheme("amber-night")
}

func TestManagerPrevious(t *testing.T) {
	m := GetManager()
	m.SetTheme("amber-night")
	prev := m.Previous()
	if prev.Name == "" {
		t.Error("Previous theme should have a name")
	}
	// Reset
	m.SetTheme("amber-night")
}

func TestMigrateThemeID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"synthwave", "electric-rose"},
		{"light", "monochrome"},
		{"plain", "amber-night"},
		{"twilight-mist", "twilight-mist"},
		{"amber-night", "amber-night"},
		{"electric-rose", "electric-rose"},
	}
	for _, tt := range tests {
		result := migrateThemeID(tt.input)
		if result != tt.expected {
			t.Errorf("migrateThemeID(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestListThemes(t *testing.T) {
	themes := ListThemes()
	if len(themes) < 6 {
		t.Errorf("Expected at least 6 themes, got %d", len(themes))
	}
}

func TestDefault(t *testing.T) {
	def := Default()
	if def.Name != "amber-night" {
		t.Errorf("Expected amber-night, got %s", def.Name)
	}
}

func TestGetThemeByName(t *testing.T) {
	t2 := GetThemeByName("amber-night")
	if t2.Name != "amber-night" {
		t.Errorf("Name lookup failed, got %s", t2.Name)
	}
	unk := GetThemeByName("does-not-exist")
	if unk.Name != design.DefaultTheme().Name {
		t.Errorf("Unknown name should return default, got %s", unk.Name)
	}
}
