package theme

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// ContrastRatio calculates the contrast ratio between two colors
func ContrastRatio(c1, c2 color.Color) float64 {
	l1 := relativeLuminance(c1)
	l2 := relativeLuminance(c2)
	if l1 > l2 {
		return (l1 + 0.05) / (l2 + 0.05)
	}
	return (l2 + 0.05) / (l1 + 0.05)
}

func relativeLuminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	rs := float64(r) / 65535.0
	gs := float64(g) / 65535.0
	bs := float64(b) / 65535.0
	return 0.2126*gamma(rs) + 0.7152*gamma(gs) + 0.0722*gamma(bs)
}

func gamma(v float64) float64 {
	if v <= 0.03928 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// ValidateTheme validates a theme and returns warnings
func ValidateTheme(t Theme) []string {
	warnings := []string{}
	bg := parseHex(string(t.Background))
	tx := parseHex(string(t.Text))
	pr := parseHex(string(t.Primary))
	if cr := ContrastRatio(tx, bg); cr < 4.5 {
		warnings = append(warnings, fmt.Sprintf("Text contrast too low: %.2f:1", cr))
	}
	if cr := ContrastRatio(pr, bg); cr < 3.0 {
		warnings = append(warnings, fmt.Sprintf("UI contrast too low: %.2f:1", cr))
	}
	return warnings
}

func parseHex(hex string) color.Color {
	h := strings.TrimPrefix(hex, "#")
	if len(h) != 6 {
		return color.White
	}
	r, errR := strconv.ParseUint(h[0:2], 16, 8)
	g, errG := strconv.ParseUint(h[2:4], 16, 8)
	b, errB := strconv.ParseUint(h[4:6], 16, 8)
	if errR != nil || errG != nil || errB != nil {
		return color.White
	}
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}
