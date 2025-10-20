package theme

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// ContrastRatio calculates the WCAG contrast ratio between two colors
func ContrastRatio(c1, c2 color.Color) float64 {
	l1 := relativeLuminance(c1)
	l2 := relativeLuminance(c2)
	
	if l1 > l2 {
		return (l1 + 0.05) / (l2 + 0.05)
	}
	return (l2 + 0.05) / (l1 + 0.05)
}

// relativeLuminance calculates the relative luminance of a color
func relativeLuminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	
	// Convert to 0-1 range
	rs := float64(r) / 65535.0
	gs := float64(g) / 65535.0
	bs := float64(b) / 65535.0
	
	// Apply gamma correction
	r8 := gammaCorrect(rs)
	g8 := gammaCorrect(gs)
	b8 := gammaCorrect(bs)
	
	// Calculate luminance
	return 0.2126*r8 + 0.7152*g8 + 0.0722*b8
}

func gammaCorrect(v float64) float64 {
	if v <= 0.03928 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// ValidateTheme checks if a theme meets WCAG AA standards
func ValidateTheme(t Theme) []string {
	warnings := []string{}
	
	// Parse colors
	bgColor := parseHexColor(string(t.Background))
	textColor := parseHexColor(string(t.Text))
	primaryColor := parseHexColor(string(t.Primary))
	
	// Check text contrast (should be 4.5:1 for AA)
	textContrast := ContrastRatio(textColor, bgColor)
	if textContrast < 4.5 {
		warnings = append(warnings, 
			fmt.Sprintf("Text contrast too low: %.2f:1 (need 4.5:1)", textContrast))
	}
	
	// Check UI element contrast (should be 3:1 for AA)
	uiContrast := ContrastRatio(primaryColor, bgColor)
	if uiContrast < 3.0 {
		warnings = append(warnings, 
			fmt.Sprintf("UI contrast too low: %.2f:1 (need 3:1)", uiContrast))
	}
	
	return warnings
}

// parseHexColor converts hex string to color.Color
func parseHexColor(hex string) color.Color {
	hex = strings.TrimPrefix(hex, "#")
	
	if len(hex) != 6 {
		return color.White
	}
	
	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)
	
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}