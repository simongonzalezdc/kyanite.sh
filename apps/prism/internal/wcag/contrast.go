// Package wcag implements WCAG 2.1 (Web Content Accessibility Guidelines) contrast algorithms.
// It provides accurate contrast ratio calculations using relative luminance with gamma correction,
// and validates color combinations against AA and AAA accessibility standards.
//
// WCAG 2.1 Requirements:
// - Level AA: Minimum contrast ratio of 4.5:1 for normal text, 3:1 for large text
// - Level AAA: Minimum contrast ratio of 7:1 for normal text, 4.5:1 for large text
package wcag

import (
	"fmt"
	"math"

	"github.com/kyanite/prism/internal/color"
)

// ContrastResult represents the result of a WCAG contrast calculation.
// It includes the contrast ratio, compliance level, and pass/fail status for AA and AAA standards.
type ContrastResult struct {
	Ratio       float64
	Level       string // "AAA", "AA", "FAIL"
	IsLargeText bool
	PassedAA    bool
	PassedAAA   bool
}

// CalculateContrast calculates the WCAG 2.1 contrast ratio between two colors
func CalculateContrast(fg, bg color.Color) float64 {
	l1 := RelativeLuminance(fg.RGB.R, fg.RGB.G, fg.RGB.B)
	l2 := RelativeLuminance(bg.RGB.R, bg.RGB.G, bg.RGB.B)

	lighter := math.Max(l1, l2)
	darker := math.Min(l1, l2)

	return (lighter + 0.05) / (darker + 0.05)
}

// RelativeLuminance calculates relative luminance per WCAG 2.1
func RelativeLuminance(r, g, b int) float64 {
	// Convert to sRGB (0-1)
	rSRGB := float64(r) / 255.0
	gSRGB := float64(g) / 255.0
	bSRGB := float64(b) / 255.0

	// Apply gamma correction (linearize)
	rLinear := linearize(rSRGB)
	gLinear := linearize(gSRGB)
	bLinear := linearize(bSRGB)

	// Calculate relative luminance
	return 0.2126*rLinear + 0.7152*gLinear + 0.0722*bLinear
}

// linearize applies gamma correction
func linearize(channel float64) float64 {
	if channel <= 0.03928 {
		return channel / 12.92
	}
	return math.Pow((channel+0.055)/1.055, 2.4)
}

// Validate checks if contrast meets WCAG standards
func Validate(fg, bg color.Color) ContrastResult {
	ratio := CalculateContrast(fg, bg)

	passedAA := ratio >= ContrastRatioAA
	passedAAA := ratio >= ContrastRatioAAA

	var level string
	if passedAAA {
		level = "AAA"
	} else if passedAA {
		level = "AA"
	} else {
		level = "FAIL"
	}

	return ContrastResult{
		Ratio:     ratio,
		Level:     level,
		PassedAA:  passedAA,
		PassedAAA: passedAAA,
	}
}

// Summary returns a human-readable summary
func (r ContrastResult) Summary() string {
	return fmt.Sprintf("%.2f:1 - WCAG %s", r.Ratio, r.Level)
}

// IsPassingAASmall returns true if passing AA for small text
func IsPassingAASmall(contrast float64) bool {
	return contrast >= ContrastRatioAA
}

// IsPassingAALarge returns true if passing AA for large text
func IsPassingAALarge(contrast float64) bool {
	return contrast >= ContrastRatioAALarge
}

// IsPassingAAASmall returns true if passing AAA for small text
func IsPassingAAASmall(contrast float64) bool {
	return contrast >= ContrastRatioAAA
}

// IsPassingAAALarge returns true if passing AAA for large text
func IsPassingAAALarge(contrast float64) bool {
	return contrast >= ContrastRatioAAALarge
}
