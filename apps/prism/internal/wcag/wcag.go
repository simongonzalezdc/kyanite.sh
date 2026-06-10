// Package wcag provides WCAG 2.1 contrast utilities for Prism.
// It wraps the shared design module's WCAG functions, preserving the
// color.Color-based API used by Prism's internal packages.
package wcag

import (
	"fmt"
	"math"

	"github.com/kyanite/design"
	"github.com/kyanite/prism/internal/color"
)

// WCAG 2.1 contrast ratio thresholds, re-exported from the design module.
const (
	ContrastRatioAA       = design.ContrastRatioAA
	ContrastRatioAAA      = design.ContrastRatioAAA
	ContrastRatioAALarge  = design.ContrastRatioAALarge
	ContrastRatioAAALarge = design.ContrastRatioAAALarge

	// MinimumContrast is the absolute minimum for palette validation (3:1).
	MinimumContrast = design.ContrastRatioAALarge
)

// ContrastResult represents the result of a WCAG contrast calculation.
type ContrastResult struct {
	Ratio       float64
	Level       string // "AAA", "AA", "FAIL"
	IsLargeText bool
	PassedAA    bool
	PassedAAA   bool
}

// Summary returns a human-readable summary.
func (r ContrastResult) Summary() string {
	return fmt.Sprintf("%.2f:1 - WCAG %s", r.Ratio, r.Level)
}

// CalculateContrast calculates the WCAG 2.1 contrast ratio between two colors.
func CalculateContrast(fg, bg color.Color) float64 {
	ratio, _ := design.ContrastRatio(fg.Hex, bg.Hex)
	return ratio
}

// Validate checks if contrast meets WCAG standards.
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

// RelativeLuminance calculates relative luminance per WCAG 2.1.
func RelativeLuminance(r, g, b int) float64 {
	rSRGB := float64(r) / 255.0
	gSRGB := float64(g) / 255.0
	bSRGB := float64(b) / 255.0
	return 0.2126*linearize(rSRGB) + 0.7152*linearize(gSRGB) + 0.0722*linearize(bSRGB)
}

func linearize(channel float64) float64 {
	if channel <= 0.03928 {
		return channel / 12.92
	}
	return math.Pow((channel+0.055)/1.055, 2.4)
}

// IsPassingAASmall returns true if passing AA for small text.
func IsPassingAASmall(contrast float64) bool {
	return contrast >= ContrastRatioAA
}

// IsPassingAALarge returns true if passing AA for large text.
func IsPassingAALarge(contrast float64) bool {
	return contrast >= ContrastRatioAALarge
}

// IsPassingAAASmall returns true if passing AAA for small text.
func IsPassingAAASmall(contrast float64) bool {
	return contrast >= ContrastRatioAAA
}

// IsPassingAAALarge returns true if passing AAA for large text.
func IsPassingAAALarge(contrast float64) bool {
	return contrast >= ContrastRatioAAALarge
}
