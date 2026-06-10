package design

import (
	"fmt"
	"math"
	"strconv"
)

// WCAG 2.1 contrast ratio thresholds.
const (
	ContrastRatioAA       = 4.5  // AA normal text
	ContrastRatioAALarge  = 3.0  // AA large text
	ContrastRatioAAA      = 7.0  // AAA normal text
	ContrastRatioAAALarge = 4.5  // AAA large text
)

// ContrastRatio calculates the WCAG 2.1 contrast ratio between two hex colors.
// Colors must be in "#RRGGBB" format.
func ContrastRatio(fg, bg string) (float64, error) {
	fgR, fgG, fgB, err := parseHexColor(fg)
	if err != nil {
		return 0, fmt.Errorf("fg color: %w", err)
	}
	bgR, bgG, bgB, err := parseHexColor(bg)
	if err != nil {
		return 0, fmt.Errorf("bg color: %w", err)
	}

	l1 := relativeLuminance(fgR, fgG, fgB)
	l2 := relativeLuminance(bgR, bgG, bgB)

	lighter := math.Max(l1, l2)
	darker := math.Min(l1, l2)

	return (lighter + 0.05) / (darker + 0.05), nil
}

// relativeLuminance calculates relative luminance per WCAG 2.1.
func relativeLuminance(r, g, b uint8) float64 {
	rSRGB := float64(r) / 255.0
	gSRGB := float64(g) / 255.0
	bSRGB := float64(b) / 255.0

	return 0.2126*linearize(rSRGB) + 0.7152*linearize(gSRGB) + 0.0722*linearize(bSRGB)
}

// linearize applies gamma correction per WCAG 2.1.
func linearize(channel float64) float64 {
	if channel <= 0.03928 {
		return channel / 12.92
	}
	return math.Pow((channel+0.055)/1.055, 2.4)
}

// parseHexColor parses a "#RRGGBB" hex color string into R, G, B components.
func parseHexColor(hex string) (uint8, uint8, uint8, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0, fmt.Errorf("invalid hex color %q: expected #RRGGBB", hex)
	}
	v, err := strconv.ParseUint(hex[1:], 16, 24)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid hex color %q: %w", hex, err)
	}
	return uint8((v >> 16) & 0xFF), uint8((v >> 8) & 0xFF), uint8(v & 0xFF), nil
}

// ValidateThemeAA validates that all text-on-background color pairs in a Theme
// meet WCAG AA contrast requirements (4.5:1 for normal text).
func ValidateThemeAA(t Theme) error {
	bg := string(t.Background)
	fgColors := []struct {
		name string
		val  string
	}{
		{"Text", string(t.Text)},
		{"Primary", string(t.Primary)},
		{"Success", string(t.Success)},
		{"Error", string(t.Error)},
		{"Warning", string(t.Warning)},
	}

	var errs []string
	for _, fg := range fgColors {
		ratio, err := ContrastRatio(fg.val, bg)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", fg.name, err))
			continue
		}
		if ratio < ContrastRatioAA {
			errs = append(errs, fmt.Sprintf("%s on Background: %.2f:1 (need %.1f:1)", fg.name, ratio, ContrastRatioAA))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("WCAG AA failures: %v", errs)
	}
	return nil
}

// MustContrastRatio calculates the contrast ratio and panics on error.
func MustContrastRatio(fg, bg string) float64 {
	r, err := ContrastRatio(fg, bg)
	if err != nil {
		panic(err)
	}
	return r
}
