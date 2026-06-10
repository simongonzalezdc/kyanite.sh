// Package color provides color representation, conversion, and manipulation utilities.
// It supports RGB, HSL, and HSV color spaces with accurate conversion algorithms.
package color

import (
	"fmt"
	"math"
)

// RGBToHex converts RGB color values to a hexadecimal color code.
// Parameters r, g, b should be in range 0-255.
// Returns a 7-character hex string (e.g., "#FF5733").
func RGBToHex(r, g, b int) string {
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// RGBToHSL converts RGB color values (0-255) to HSL color space.
// Returns HSL with H in range 0-360, S and L in range 0-100.
// Uses the standard HSL conversion algorithm with proper handling of edge cases.
func RGBToHSL(rgb RGB) HSL {
	r := float64(rgb.R) / 255.0
	g := float64(rgb.G) / 255.0
	b := float64(rgb.B) / 255.0

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	delta := max - min

	// Lightness
	l := (max + min) / 2.0

	// Saturation
	var s float64
	if delta == 0 {
		s = 0
	} else {
		if l < 0.5 {
			s = delta / (max + min)
		} else {
			s = delta / (2.0 - max - min)
		}
	}

	// Hue
	var h float64
	if delta == 0 {
		h = 0
	} else {
		switch max {
		case r:
			h = math.Mod((g-b)/delta, 6.0)
		case g:
			h = ((b - r) / delta) + 2.0
		case b:
			h = ((r - g) / delta) + 4.0
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}

	return HSL{
		H: h,
		S: s * 100,
		L: l * 100,
	}
}

// HSLToRGB converts HSL to RGB (0-255)
func HSLToRGB(h, s, l float64) RGB {
	s = s / 100.0
	l = l / 100.0

	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return RGB{
		R: int(math.Round((r + m) * 255)),
		G: int(math.Round((g + m) * 255)),
		B: int(math.Round((b + m) * 255)),
	}
}

// RGBToHSV converts RGB to HSV
func RGBToHSV(rgb RGB) HSV {
	r := float64(rgb.R) / 255.0
	g := float64(rgb.G) / 255.0
	b := float64(rgb.B) / 255.0

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	delta := max - min

	// Hue (same as HSL)
	var h float64
	if delta == 0 {
		h = 0
	} else {
		switch max {
		case r:
			h = math.Mod((g-b)/delta, 6.0)
		case g:
			h = ((b - r) / delta) + 2.0
		case b:
			h = ((r - g) / delta) + 4.0
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}

	// Saturation
	var s float64
	if max == 0 {
		s = 0
	} else {
		s = delta / max
	}

	// Value
	v := max

	return HSV{
		H: h,
		S: s * 100,
		V: v * 100,
	}
}

// HSVToRGB converts HSV to RGB
func HSVToRGB(h, s, v float64) RGB {
	s = s / 100.0
	v = v / 100.0

	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return RGB{
		R: int(math.Round((r + m) * 255)),
		G: int(math.Round((g + m) * 255)),
		B: int(math.Round((b + m) * 255)),
	}
}
