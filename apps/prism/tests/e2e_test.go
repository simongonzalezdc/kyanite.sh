package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kyanite/prism/internal/color"
	"github.com/kyanite/prism/internal/export"
	"github.com/kyanite/prism/internal/palette"
	"github.com/kyanite/prism/internal/storage"
	"github.com/kyanite/prism/internal/wcag"
)

func cleanupTempDir(t *testing.T, path string) {
	t.Helper()
	t.Cleanup(func() {
		if err := os.RemoveAll(path); err != nil {
			t.Errorf("Failed to remove temp dir %s: %v", path, err)
		}
	})
}

// TestE2E_CreateAndExportPalette tests the complete workflow of creating and exporting a palette
func TestE2E_CreateAndExportPalette(t *testing.T) {
	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "prism-e2e-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	cleanupTempDir(t, tmpDir)

	// Override config directory
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Step 1: User enters a color
	userInput := "#FF5733"
	t.Logf("User input: %s", userInput)

	// Step 2: Parse the color
	baseColor, err := color.ParseHex(userInput)
	if err != nil {
		t.Fatalf("Failed to parse user color: %v", err)
	}
	t.Logf("Parsed color: %s (RGB: %d,%d,%d)", baseColor.Hex, baseColor.RGB.R, baseColor.RGB.G, baseColor.RGB.B)

	// Step 3: Generate palette with chosen harmony
	chosenHarmony := palette.Triadic
	pal, err := palette.Generate(baseColor, chosenHarmony)
	if err != nil {
		t.Fatalf("Failed to generate palette: %v", err)
	}
	pal.Name = "My Beautiful Triadic Palette"
	pal.Description = "Created for my website redesign"
	pal.Tags = []string{"website", "warm", "vibrant"}
	t.Logf("Generated %s palette with %d colors", pal.HarmonyRule, len(pal.Colors))

	// Step 4: Validate WCAG contrast for the palette
	if len(pal.Colors) >= 2 {
		// Check contrast between first two colors
		ratio := wcag.CalculateContrast(pal.Colors[0], pal.Colors[1])
		result := wcag.Validate(pal.Colors[0], pal.Colors[1])
		t.Logf("Color 1 vs Color 2: %s (ratio: %.2f)", result.Summary(), ratio)

		// Full palette validation
		minContrast, ok := palette.ValidatePaletteContrast(pal)
		t.Logf("Palette validation: min contrast %.2f, ok: %v", minContrast, ok)
	}

	// Step 5: Save the palette
	err = storage.SavePalette(pal)
	if err != nil {
		t.Fatalf("Failed to save palette: %v", err)
	}
	t.Logf("Saved palette with ID: %s", pal.ID)

	// Step 6: Export to CSS for use in website
	cssOutput := export.ExportCSS(pal)
	cssFile := filepath.Join(tmpDir, "palette.css")
	err = os.WriteFile(cssFile, []byte(cssOutput), 0644)
	if err != nil {
		t.Fatalf("Failed to write CSS file: %v", err)
	}
	t.Logf("Exported to CSS: %s", cssFile)

	// Step 7: Export to JSON for backup
	jsonOutput, err := export.ExportJSON(pal)
	if err != nil {
		t.Fatalf("Failed to export JSON: %v", err)
	}
	jsonFile := filepath.Join(tmpDir, "palette.json")
	err = os.WriteFile(jsonFile, jsonOutput, 0644)
	if err != nil {
		t.Fatalf("Failed to write JSON file: %v", err)
	}
	t.Logf("Exported to JSON: %s", jsonFile)

	// Step 8: Verify we can load it back later
	loadedPal, err := storage.LoadPalette(pal.ID)
	if err != nil {
		t.Fatalf("Failed to load saved palette: %v", err)
	}

	if loadedPal.Name != pal.Name {
		t.Errorf("Loaded palette name = %s, want %s", loadedPal.Name, pal.Name)
	}

	// Step 9: List all saved palettes
	allPalettes, err := storage.ListPalettes()
	if err != nil {
		t.Fatalf("Failed to list palettes: %v", err)
	}
	t.Logf("Total saved palettes: %d", len(allPalettes))

	// Verify our palette is in the list
	found := false
	for _, p := range allPalettes {
		if p.ID == pal.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("Saved palette should appear in ListPalettes")
	}

	t.Log("✓ E2E workflow completed successfully")
}

// TestE2E_CompareMultiplePalettes tests workflow of comparing different harmonies
func TestE2E_CompareMultiplePalettes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "prism-e2e-compare-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	cleanupTempDir(t, tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// User wants to try different harmonies for the same base color
	baseColor, _ := color.ParseHex("#3498DB") // Nice blue

	harmonies := []palette.HarmonyRule{
		palette.Complementary,
		palette.Triadic,
		palette.Analogous,
	}

	results := make(map[string]*palette.Palette)

	// Generate all harmonies
	for _, harmony := range harmonies {
		pal, err := palette.Generate(baseColor, harmony)
		if err != nil {
			t.Fatalf("Failed to generate %s palette: %v", harmony, err)
		}
		pal.Name = string(harmony) + " Blue"
		results[string(harmony)] = &pal

		// Save each palette
		err = storage.SavePalette(pal)
		if err != nil {
			t.Fatalf("Failed to save %s palette: %v", harmony, err)
		}

		t.Logf("Generated %s: %d colors", harmony, len(pal.Colors))
	}

	// Compare accessibility
	for name, pal := range results {
		if len(pal.Colors) >= 2 {
			ratio := wcag.CalculateContrast(pal.Colors[0], pal.Colors[1])
			result := wcag.Validate(pal.Colors[0], pal.Colors[1])
			t.Logf("%s - First pair contrast: %.2f (%s)", name, ratio, result.Level)
		}
	}

	// User decides on complementary, exports it
	chosenPalette := results["complementary"]
	cssOutput := export.ExportCSS(*chosenPalette)

	if !strings.Contains(cssOutput, "--color-primary:") {
		t.Error("CSS export should contain color variables")
	}

	// Delete the other palettes
	for name, pal := range results {
		if name != "complementary" {
			err = storage.DeletePalette(pal.ID)
			if err != nil {
				t.Fatalf("Failed to delete %s palette: %v", name, err)
			}
		}
	}

	// Verify only complementary remains
	allPalettes, _ := storage.ListPalettes()
	if len(allPalettes) != 1 {
		t.Errorf("Should have 1 palette remaining, got %d", len(allPalettes))
	}

	t.Log("✓ E2E comparison workflow completed successfully")
}

// TestE2E_ColorModificationWorkflow tests modifying colors and regenerating palettes
func TestE2E_ColorModificationWorkflow(t *testing.T) {
	// Start with a color that's too dark
	darkColor, _ := color.ParseHex("#1A1A1A")
	t.Logf("Starting color: %s", darkColor.Hex)

	// User wants to lighten it
	lightened := darkColor.Lighten(40)
	t.Logf("After lightening: %s", lightened.Hex)

	// Still not saturated enough
	vibrant := lightened.Saturate(30)
	t.Logf("After saturating: %s", vibrant.Hex)

	// Generate palette from modified color
	pal, _ := palette.Generate(vibrant, palette.Analogous)
	pal.Name = "Lightened Analogous"

	// Check if the modifications improved contrast for dark-on-light alternatives.
	black, _ := color.ParseHex("#000000")
	originalContrast := wcag.CalculateContrast(darkColor, black)
	modifiedContrast := wcag.CalculateContrast(vibrant, black)

	t.Logf("Original contrast: %.2f", originalContrast)
	t.Logf("Modified contrast: %.2f", modifiedContrast)

	if modifiedContrast <= originalContrast {
		t.Error("Lightening should improve contrast with black")
	}

	// Export for use
	tomlOutput := export.ExportTOML(pal)
	if !strings.Contains(tomlOutput, "harmony_rule") {
		t.Error("TOML export should contain harmony_rule")
	}

	t.Log("✓ E2E color modification workflow completed successfully")
}

// TestE2E_ThemeCreation tests creating a complete theme
func TestE2E_ThemeCreation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "prism-e2e-theme-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	cleanupTempDir(t, tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// User wants to create a complete theme for their app
	// They start with a brand color
	brandColor, _ := color.ParseHex("#6366F1") // Indigo

	// Generate a tetradic palette for variety
	pal, _ := palette.Generate(brandColor, palette.Tetradic)
	pal.Name = "Brand Theme 2024"
	pal.Description = "Main application theme"

	// Save the palette
	if err := storage.SavePalette(pal); err != nil {
		t.Fatalf("Failed to save palette: %v", err)
	}

	// Export as theme
	themeData, err := export.ExportTheme(pal)
	if err != nil {
		t.Fatalf("Failed to export theme: %v", err)
	}

	// Parse theme to verify structure
	var theme map[string]interface{}
	err = json.Unmarshal(themeData, &theme)
	if err != nil {
		t.Fatalf("Theme JSON is invalid: %v", err)
	}

	// Verify theme has required roles
	themeColors := theme["theme"].(map[string]interface{})
	requiredRoles := []string{"primary"}
	for _, role := range requiredRoles {
		if _, ok := themeColors[role]; !ok {
			t.Errorf("Theme missing required role: %s", role)
		}
	}

	// Also export as CSS for web
	cssData := export.ExportCSS(pal)
	cssFile := filepath.Join(tmpDir, "theme.css")
	if err := os.WriteFile(cssFile, []byte(cssData), 0644); err != nil {
		t.Fatalf("Failed to write CSS theme: %v", err)
	}

	// And JSON for config
	jsonData, err := export.ExportJSON(pal)
	if err != nil {
		t.Fatalf("Failed to export JSON theme: %v", err)
	}
	jsonFile := filepath.Join(tmpDir, "theme.json")
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		t.Fatalf("Failed to write JSON theme: %v", err)
	}

	// Verify files exist
	if _, err := os.Stat(cssFile); os.IsNotExist(err) {
		t.Error("CSS theme file should exist")
	}
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		t.Error("JSON theme file should exist")
	}

	t.Log("✓ E2E theme creation workflow completed successfully")
}

// TestE2E_AccessibilityFocusedWorkflow tests creating an accessible palette
func TestE2E_AccessibilityFocusedWorkflow(t *testing.T) {
	// User wants to create an accessible palette
	// Start with a medium blue
	baseColor, _ := color.ParseHex("#4A90E2")

	// Generate monochromatic for consistency
	pal, _ := palette.Generate(baseColor, palette.Monochromatic)
	pal.Name = "Accessible Blues"

	// Validate against white background
	white, _ := color.ParseHex("#FFFFFF")
	accessibleColors := []color.Color{}

	t.Log("Testing colors for AA compliance on white background:")
	for i, c := range pal.Colors {
		ratio := wcag.CalculateContrast(c, white)
		result := wcag.Validate(c, white)

		t.Logf("  Color %d (%s): %.2f - %s", i+1, c.Hex, ratio, result.Level)

		if result.PassedAA {
			accessibleColors = append(accessibleColors, c)
		}
	}

	t.Logf("Found %d AA-compliant colors out of %d", len(accessibleColors), len(pal.Colors))

	// If we don't have enough accessible colors, darken some
	if len(accessibleColors) < 3 {
		t.Log("Need more accessible colors, darkening base color...")
		darkerBase := baseColor.Darken(30)
		pal, _ = palette.Generate(darkerBase, palette.Monochromatic)

		// Re-validate
		accessibleColors = []color.Color{}
		for _, c := range pal.Colors {
			result := wcag.Validate(c, white)
			if result.PassedAA {
				accessibleColors = append(accessibleColors, c)
			}
		}
		t.Logf("After darkening: %d AA-compliant colors", len(accessibleColors))
	}

	// Create a palette with only accessible colors
	accessiblePalette := &palette.Palette{
		ID:          "accessible-" + pal.ID,
		Name:        "Accessible " + pal.Name,
		Description: "AA-compliant colors only",
		Colors:      accessibleColors,
		HarmonyRule: pal.HarmonyRule,
		BaseColor:   pal.BaseColor,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Tags:        []string{"accessible", "wcag-aa"},
	}

	// Export for use
	cssOutput := export.ExportCSS(*accessiblePalette)
	if !strings.Contains(cssOutput, "/* Accessible") {
		t.Error("CSS should note this is an accessible palette")
	}

	t.Log("✓ E2E accessibility workflow completed successfully")
}

// TestE2E_BulkPaletteManagement tests managing multiple palettes
func TestE2E_BulkPaletteManagement(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "prism-e2e-bulk-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	cleanupTempDir(t, tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create multiple palettes for a project
	projectColors := []string{
		"#E74C3C", // Red
		"#3498DB", // Blue
		"#2ECC71", // Green
		"#F39C12", // Orange
		"#9B59B6", // Purple
	}

	createdPalettes := []string{}

	// Generate palette for each color
	for i, hexColor := range projectColors {
		baseColor, _ := color.ParseHex(hexColor)
		pal, _ := palette.Generate(baseColor, palette.Complementary)
		pal.Name = "Project Palette " + string(rune('A'+i))
		pal.Tags = []string{"project", "option-" + string(rune('A'+i))}

		err := storage.SavePalette(pal)
		if err != nil {
			t.Fatalf("Failed to save palette %d: %v", i, err)
		}

		createdPalettes = append(createdPalettes, pal.ID)
		t.Logf("Created palette %s", pal.Name)
	}

	// List all palettes
	allPalettes, err := storage.ListPalettes()
	if err != nil {
		t.Fatalf("Failed to list palettes: %v", err)
	}

	if len(allPalettes) != len(projectColors) {
		t.Errorf("Should have %d palettes, got %d", len(projectColors), len(allPalettes))
	}

	// User decides to keep only the first two
	for i := 2; i < len(createdPalettes); i++ {
		err = storage.DeletePalette(createdPalettes[i])
		if err != nil {
			t.Fatalf("Failed to delete palette %d: %v", i, err)
		}
		t.Logf("Deleted palette %s", createdPalettes[i])
	}

	// Verify deletion
	remainingPalettes, err := storage.ListPalettes()
	if err != nil {
		t.Fatalf("Failed to list remaining palettes: %v", err)
	}
	if len(remainingPalettes) != 2 {
		t.Errorf("Should have 2 palettes remaining, got %d", len(remainingPalettes))
	}

	// Export remaining palettes
	exportDir := filepath.Join(tmpDir, "exports")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		t.Fatalf("Failed to create export dir: %v", err)
	}

	for _, pal := range remainingPalettes {
		// Export to multiple formats
		jsonData, err := export.ExportJSON(pal)
		if err != nil {
			t.Fatalf("Failed to export palette JSON: %v", err)
		}
		jsonFile := filepath.Join(exportDir, pal.ID+".json")
		if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
			t.Fatalf("Failed to write palette JSON: %v", err)
		}

		cssData := export.ExportCSS(pal)
		cssFile := filepath.Join(exportDir, pal.ID+".css")
		if err := os.WriteFile(cssFile, []byte(cssData), 0644); err != nil {
			t.Fatalf("Failed to write palette CSS: %v", err)
		}

		t.Logf("Exported palette %s to multiple formats", pal.Name)
	}

	// Verify exports
	files, err := os.ReadDir(exportDir)
	if err != nil {
		t.Fatalf("Failed to read export dir: %v", err)
	}

	if len(files) != 4 { // 2 palettes × 2 formats
		t.Errorf("Should have 4 export files, got %d", len(files))
	}

	t.Log("✓ E2E bulk management workflow completed successfully")
}

// TestE2E_ConfigurationManagement tests app configuration
func TestE2E_ConfigurationManagement(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "prism-e2e-config-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	cleanupTempDir(t, tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// First time user - should get defaults
	cfg, err := storage.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	defaultCfg := storage.DefaultConfig()
	if cfg.Theme != defaultCfg.Theme {
		t.Errorf("Default theme = %s, want %s", cfg.Theme, defaultCfg.Theme)
	}

	// User changes settings
	cfg.Theme = "custom-dark"
	cfg.AutoSave = false
	cfg.AutoSaveInterval = 120

	err = storage.SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Load again to verify persistence
	loadedCfg, err := storage.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig (second) failed: %v", err)
	}

	if loadedCfg.Theme != "custom-dark" {
		t.Errorf("Loaded theme = %s, want custom-dark", loadedCfg.Theme)
	}

	if loadedCfg.AutoSave != false {
		t.Errorf("Loaded AutoSave = %v, want false", loadedCfg.AutoSave)
	}

	if loadedCfg.AutoSaveInterval != 120 {
		t.Errorf("Loaded AutoSaveInterval = %d, want 120", loadedCfg.AutoSaveInterval)
	}

	t.Log("✓ E2E configuration workflow completed successfully")
}
