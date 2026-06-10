package color

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Color represents a color in multiple formats
type Color struct {
	Hex  string
	RGB  RGB
	HSL  HSL
	HSV  HSV
	Name string
}

// RGB represents RGB color space (0-255)
type RGB struct {
	R, G, B int
}

// HSL represents HSL color space
type HSL struct {
	H float64 // 0-360 degrees
	S float64 // 0-100 percent
	L float64 // 0-100 percent
}

// HSV represents HSV color space
type HSV struct {
	H float64 // 0-360
	S float64 // 0-100
	V float64 // 0-100
}

// ParseHex creates a Color from hex string
func ParseHex(hex string) (Color, error) {
	// Remove # if present
	original := hex
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) != 6 {
		return Color{}, fmt.Errorf("invalid hex color '%s': must be 6 characters (e.g., #FF5733 or FF5733), got %d characters", original, len(hex))
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 0)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color '%s': red component '%s' is not valid hexadecimal", original, hex[0:2])
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 0)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color '%s': green component '%s' is not valid hexadecimal", original, hex[2:4])
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 0)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color '%s': blue component '%s' is not valid hexadecimal", original, hex[4:6])
	}

	rgb := RGB{R: int(r), G: int(g), B: int(b)}
	hsl := RGBToHSL(rgb)
	hsv := RGBToHSV(rgb)

	return Color{
		Hex: "#" + strings.ToUpper(hex),
		RGB: rgb,
		HSL: hsl,
		HSV: hsv,
	}, nil
}

// NewFromRGB creates a Color from RGB values
func NewFromRGB(r, g, b int) Color {
	hex := RGBToHex(r, g, b)
	hsl := RGBToHSL(RGB{r, g, b})
	hsv := RGBToHSV(RGB{r, g, b})

	return Color{
		Hex: hex,
		RGB: RGB{r, g, b},
		HSL: hsl,
		HSV: hsv,
	}
}

// NewFromHSL creates a Color from HSL values
func NewFromHSL(h, s, l float64) Color {
	rgb := HSLToRGB(h, s, l)
	hex := RGBToHex(rgb.R, rgb.G, rgb.B)
	hsv := RGBToHSV(rgb)

	return Color{
		Hex: hex,
		RGB: rgb,
		HSL: HSL{h, s, l},
		HSV: hsv,
	}
}

// Complement returns the complementary color (180° hue shift)
func (c Color) Complement() Color {
	h := math.Mod(c.HSL.H+180, 360)
	return NewFromHSL(h, c.HSL.S, c.HSL.L)
}

// Lighten increases lightness by percentage
func (c Color) Lighten(pct float64) Color {
	l := math.Min(100, c.HSL.L+pct)
	return NewFromHSL(c.HSL.H, c.HSL.S, l)
}

// Darken decreases lightness by percentage
func (c Color) Darken(pct float64) Color {
	l := math.Max(0, c.HSL.L-pct)
	return NewFromHSL(c.HSL.H, c.HSL.S, l)
}

// Saturate increases saturation by percentage
func (c Color) Saturate(pct float64) Color {
	s := math.Min(100, c.HSL.S+pct)
	return NewFromHSL(c.HSL.H, s, c.HSL.L)
}

// Desaturate decreases saturation by percentage
func (c Color) Desaturate(pct float64) Color {
	s := math.Max(0, c.HSL.S-pct)
	return NewFromHSL(c.HSL.H, s, c.HSL.L)
}

// IsWarm returns true if the color is considered warm
func (c Color) IsWarm() bool {
	h := c.HSL.H
	return (h >= 0 && h < 90) || (h >= 270 && h <= 360)
}

// IsCool returns true if the color is considered cool
func (c Color) IsCool() bool {
	return !c.IsWarm()
}
