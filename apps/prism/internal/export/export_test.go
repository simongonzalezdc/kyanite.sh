package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	color "github.com/kyanite/prism/internal/color"
	palette "github.com/kyanite/prism/internal/palette"
)

// createTestPalette creates a test palette for export tests
func createTestPalette() *palette.Palette {
	baseColor, _ := color.ParseHex("#FF5733")
	pal, _ := palette.Generate(baseColor, palette.Triadic)
	pal.Name = "Sunset Triadic"
	pal.Description = "A warm triadic palette"
	pal.Tags = []string{"warm", "vibrant"}
	return &pal
}

// TestExportIntegration tests all export formats
func TestExportIntegration(t *testing.T) {
	pal := createTestPalette()

	t.Run("ExportJSON", func(t *testing.T) {
		data, err := ExportJSON(*pal)
		if err != nil {
			t.Fatalf("ExportJSON failed: %v", err)
		}

		// Verify it's valid JSON
		var parsed map[string]interface{}
		err = json.Unmarshal(data, &parsed)
		if err != nil {
			t.Fatalf("Exported JSON is invalid: %v", err)
		}

		// Verify key fields are present
		if _, ok := parsed["id"]; !ok {
			t.Error("Exported JSON missing 'id' field")
		}
		if _, ok := parsed["name"]; !ok {
			t.Error("Exported JSON missing 'name' field")
		}
		if _, ok := parsed["colors"]; !ok {
			t.Error("Exported JSON missing 'colors' field")
		}
		if _, ok := parsed["harmony_rule"]; !ok {
			t.Error("Exported JSON missing 'harmony_rule' field")
		}

		// Verify name matches
		if parsed["name"] != pal.Name {
			t.Errorf("Exported name = %v, want %s", parsed["name"], pal.Name)
		}
	})

	t.Run("ExportJSONCompact", func(t *testing.T) {
		data, err := ExportJSONCompact(*pal)
		if err != nil {
			t.Fatalf("ExportJSONCompact failed: %v", err)
		}

		// Verify it's valid JSON
		var parsed map[string]interface{}
		err = json.Unmarshal(data, &parsed)
		if err != nil {
			t.Fatalf("Exported compact JSON is invalid: %v", err)
		}

		// Compact should have no extra whitespace
		if strings.Contains(string(data), "\n") {
			t.Error("Compact JSON should not contain newlines")
		}
	})

	t.Run("ExportCSS", func(t *testing.T) {
		css := ExportCSS(*pal)

		// Verify CSS structure
		if !strings.Contains(css, ":root {") {
			t.Error("CSS should contain ':root {' declaration")
		}

		if !strings.Contains(css, "--color-primary:") {
			t.Error("CSS should contain primary color variable")
		}

		if !strings.Contains(css, "--color-primary-rgb:") {
			t.Error("CSS should contain RGB color variable")
		}

		// Verify palette name is in comments
		if !strings.Contains(css, pal.Name) {
			t.Error("CSS should contain palette name in comments")
		}

		// Verify harmony rule is in comments
		if !strings.Contains(css, string(pal.HarmonyRule)) {
			t.Error("CSS should contain harmony rule in comments")
		}

		// Count color variables - should have one for each color
		for i := range pal.Colors {
			switch i {
			case 0:
				if !strings.Contains(css, "--color-primary:") {
					t.Error("CSS should contain primary color")
				}
			case 1:
				if !strings.Contains(css, "--color-secondary:") {
					t.Error("CSS should contain secondary color")
				}
			}
		}
	})

	t.Run("ExportTOML", func(t *testing.T) {
		toml := ExportTOML(*pal)

		// Verify TOML structure
		if !strings.Contains(toml, "name =") {
			t.Error("TOML should contain name field")
		}

		if !strings.Contains(toml, "harmony_rule =") {
			t.Error("TOML should contain harmony_rule field")
		}

		if !strings.Contains(toml, "[[colors]]") {
			t.Error("TOML should contain colors array")
		}

		// Verify palette name
		if !strings.Contains(toml, pal.Name) {
			t.Error("TOML should contain palette name")
		}

		// Verify at least one color hex value
		if !strings.Contains(toml, "hex =") {
			t.Error("TOML should contain hex color values")
		}

		// Count color entries
		colorCount := strings.Count(toml, "[[colors]]")
		if colorCount != len(pal.Colors) {
			t.Errorf("TOML should contain %d colors, got %d", len(pal.Colors), colorCount)
		}
	})

	t.Run("ExportTheme", func(t *testing.T) {
		data, err := ExportTheme(*pal)
		if err != nil {
			t.Fatalf("ExportTheme failed: %v", err)
		}

		// Verify it's valid JSON
		var theme map[string]interface{}
		err = json.Unmarshal(data, &theme)
		if err != nil {
			t.Fatalf("Exported theme is invalid JSON: %v", err)
		}

		// Verify theme structure
		if _, ok := theme["name"]; !ok {
			t.Error("Theme missing 'name' field")
		}
		if _, ok := theme["kyanite_version"]; !ok {
			t.Error("Theme missing 'kyanite_version' field")
		}
		if _, ok := theme["theme"]; !ok {
			t.Error("Theme missing 'theme' field")
		}

		// Verify theme name
		if theme["name"] != pal.Name {
			t.Errorf("Theme name = %v, want %s", theme["name"], pal.Name)
		}

		// Verify theme colors
		themeColors, ok := theme["theme"].(map[string]interface{})
		if !ok {
			t.Fatal("Theme 'theme' field should be an object")
		}

		// Should have color mappings
		if len(themeColors) == 0 {
			t.Error("Theme should have at least one color mapping")
		}

		// Verify primary color if palette has colors
		if len(pal.Colors) > 0 {
			if _, ok := themeColors["primary"]; !ok {
				t.Error("Theme should have 'primary' color")
			}
		}
	})
}

// TestExportFormatsComparison tests consistency across formats
func TestExportFormatsComparison(t *testing.T) {
	pal := createTestPalette()

	jsonData, _ := ExportJSON(*pal)
	jsonCompact, _ := ExportJSONCompact(*pal)
	css := ExportCSS(*pal)
	toml := ExportTOML(*pal)
	theme, _ := ExportTheme(*pal)

	// All formats should contain the palette name
	formats := map[string]string{
		"JSON":        string(jsonData),
		"JSONCompact": string(jsonCompact),
		"CSS":         css,
		"TOML":        toml,
		"Theme":       string(theme),
	}

	for name, content := range formats {
		if !strings.Contains(content, pal.Name) || !strings.Contains(content, "Sunset") {
			t.Errorf("%s format should contain palette name '%s'", name, pal.Name)
		}
	}

	// All formats should reference the colors
	firstColorHex := pal.Colors[0].Hex
	colorFormats := map[string]string{
		"JSON":        string(jsonData),
		"JSONCompact": string(jsonCompact),
		"CSS":         css,
		"TOML":        toml,
		"Theme":       string(theme),
	}

	for name, content := range colorFormats {
		if !strings.Contains(strings.ToUpper(content), strings.ToUpper(firstColorHex)) {
			t.Errorf("%s format should contain first color %s", name, firstColorHex)
		}
	}
}

// TestExportEdgeCases tests export with edge cases
func TestExportEdgeCases(t *testing.T) {
	t.Run("EmptyPalette", func(t *testing.T) {
		emptyPal := &palette.Palette{
			ID:          "empty-test",
			Name:        "Empty",
			Colors:      []color.Color{},
			HarmonyRule: "monochromatic",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// JSON export should work with empty colors
		jsonData, err := ExportJSON(*emptyPal)
		if err != nil {
			t.Errorf("ExportJSON should handle empty palette: %v", err)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(jsonData, &parsed); err != nil {
			t.Fatalf("Exported JSON should be valid: %v", err)
		}
		colors := parsed["colors"].([]interface{})
		if len(colors) != 0 {
			t.Errorf("Empty palette should have 0 colors, got %d", len(colors))
		}

		// CSS export should work
		css := ExportCSS(*emptyPal)
		if !strings.Contains(css, ":root {") {
			t.Error("CSS should still have :root even with empty palette")
		}

		// TOML export should work
		toml := ExportTOML(*emptyPal)
		if !strings.Contains(toml, "name =") {
			t.Error("TOML should contain name even with empty palette")
		}

		// Theme export should work
		_, err = ExportTheme(*emptyPal)
		if err != nil {
			t.Errorf("ExportTheme should handle empty palette: %v", err)
		}
	})

	t.Run("SingleColorPalette", func(t *testing.T) {
		singleColor, _ := color.ParseHex("#FF0000")
		singlePal := &palette.Palette{
			ID:          "single-test",
			Name:        "Single Red",
			Colors:      []color.Color{singleColor},
			HarmonyRule: "monochromatic",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		css := ExportCSS(*singlePal)
		if !strings.Contains(css, "--color-primary: #FF0000") {
			t.Error("CSS should contain the single color")
		}

		toml := ExportTOML(*singlePal)
		colorCount := strings.Count(toml, "[[colors]]")
		if colorCount != 1 {
			t.Errorf("TOML should have exactly 1 color, got %d", colorCount)
		}
	})

	t.Run("SpecialCharactersInName", func(t *testing.T) {
		baseColor, _ := color.ParseHex("#FF5733")
		pal, _ := palette.Generate(baseColor, palette.Triadic)
		pal.Name = "Test \"Quotes\" & Symbols <>"
		pal.Description = "Test 'description' with symbols"

		// JSON should handle special characters
		jsonData, err := ExportJSON(pal)
		if err != nil {
			t.Fatalf("ExportJSON should handle special characters: %v", err)
		}

		var parsed map[string]interface{}
		err = json.Unmarshal(jsonData, &parsed)
		if err != nil {
			t.Errorf("Exported JSON with special chars should be valid: %v", err)
		}

		// CSS should handle special characters in comments
		css := ExportCSS(pal)
		if !strings.Contains(css, "Test") {
			t.Error("CSS should contain the palette name")
		}
	})

	t.Run("LargePalette", func(t *testing.T) {
		// Create palette with many colors
		largeColors := []color.Color{}
		for i := 0; i < 20; i++ {
			c, _ := color.ParseHex("#FF5733")
			largeColors = append(largeColors, c)
		}

		largePal := &palette.Palette{
			ID:          "large-test",
			Name:        "Large Palette",
			Colors:      largeColors,
			HarmonyRule: "custom",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// All exports should handle large palette
		_, err := ExportJSON(*largePal)
		if err != nil {
			t.Errorf("ExportJSON should handle large palette: %v", err)
		}

		css := ExportCSS(*largePal)
		if len(css) == 0 {
			t.Error("CSS export should produce output for large palette")
		}

		toml := ExportTOML(*largePal)
		colorCount := strings.Count(toml, "[[colors]]")
		if colorCount != 20 {
			t.Errorf("TOML should have 20 colors, got %d", colorCount)
		}
	})
}

// TestExportRoundTrip tests exporting and re-importing
func TestExportRoundTrip(t *testing.T) {
	originalPal := createTestPalette()

	// Export to JSON
	jsonData, err := ExportJSON(*originalPal)
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}

	// Parse back to palette
	var importedPal palette.Palette
	err = json.Unmarshal(jsonData, &importedPal)
	if err != nil {
		t.Fatalf("Failed to parse exported JSON: %v", err)
	}

	// Verify data matches
	if importedPal.Name != originalPal.Name {
		t.Errorf("Round-trip name = %s, want %s", importedPal.Name, originalPal.Name)
	}

	if importedPal.HarmonyRule != originalPal.HarmonyRule {
		t.Errorf("Round-trip harmony_rule = %s, want %s", importedPal.HarmonyRule, originalPal.HarmonyRule)
	}

	if len(importedPal.Colors) != len(originalPal.Colors) {
		t.Errorf("Round-trip colors count = %d, want %d", len(importedPal.Colors), len(originalPal.Colors))
	}

	// Verify first color matches
	if len(importedPal.Colors) > 0 && len(originalPal.Colors) > 0 {
		if importedPal.Colors[0].Hex != originalPal.Colors[0].Hex {
			t.Errorf("Round-trip first color = %s, want %s",
				importedPal.Colors[0].Hex, originalPal.Colors[0].Hex)
		}
	}
}
