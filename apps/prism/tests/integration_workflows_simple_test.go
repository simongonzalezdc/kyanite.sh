package tests

import (
	"testing"

	"github.com/kyanite/prism/internal/color"
	"github.com/kyanite/prism/internal/export"
	"github.com/kyanite/prism/internal/palette"
	"github.com/kyanite/prism/internal/wcag"
)

// TestSimpleColorWorkflow tests basic color operations
func TestSimpleColorWorkflow(t *testing.T) {
	// Parse a color
	c, err := color.ParseHex("#FF5733")
	if err != nil {
		t.Fatalf("ParseHex failed: %v", err)
	}

	// Generate palette
	pal, err := palette.Generate(c, palette.Triadic)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Validate
	minContrast, ok := palette.ValidatePaletteContrast(pal)
	t.Logf("Min contrast: %.2f, ok: %v", minContrast, ok)

	// Export
	jsonData, _ := export.ExportJSON(pal)
	if len(jsonData) == 0 {
		t.Error("JSON export should produce data")
	}

	css := export.ExportCSS(pal)
	if len(css) == 0 {
		t.Error("CSS export should produce data")
	}
}

// TestWCAGValidation tests WCAG contrast validation
func TestWCAGValidation(t *testing.T) {
	white, _ := color.ParseHex("#FFFFFF")
	black, _ := color.ParseHex("#000000")

	ratio := wcag.CalculateContrast(white, black)
	if ratio < 20.0 {
		t.Errorf("White/black contrast = %.2f, expected ~21", ratio)
	}

	result := wcag.Validate(white, black)
	if result.Level != "AAA" {
		t.Errorf("Level = %s, want AAA", result.Level)
	}
}
