package tests

import (
	"testing"

	"github.com/kyanite/prism/internal/color"
	"github.com/kyanite/prism/internal/palette"
	"github.com/kyanite/prism/internal/wcag"
)

// Benchmark color conversions
func BenchmarkRGBToHSL(b *testing.B) {
	rgb := color.RGB{R: 255, G: 87, B: 51}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = color.RGBToHSL(rgb)
	}
}

func BenchmarkHSLToRGB(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = color.HSLToRGB(12, 100, 60)
	}
}

func BenchmarkParseHex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = color.ParseHex("#FF5733")
	}
}

func BenchmarkRGBToHex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = color.RGBToHex(255, 87, 51)
	}
}

// Benchmark color operations
func BenchmarkColorLighten(b *testing.B) {
	c, _ := color.ParseHex("#FF5733")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Lighten(20)
	}
}

func BenchmarkColorDarken(b *testing.B) {
	c, _ := color.ParseHex("#FF5733")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Darken(20)
	}
}

func BenchmarkColorSaturate(b *testing.B) {
	c, _ := color.ParseHex("#FF5733")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Saturate(20)
	}
}

// Benchmark palette generation
func BenchmarkGenerateMonochromatic(b *testing.B) {
	baseColor, _ := color.ParseHex("#FF5733")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = palette.GenerateMonochromatic(baseColor)
	}
}

func BenchmarkGenerateComplementary(b *testing.B) {
	baseColor, _ := color.ParseHex("#FF5733")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = palette.GenerateComplementary(baseColor)
	}
}

func BenchmarkGenerateTriadic(b *testing.B) {
	baseColor, _ := color.ParseHex("#FF5733")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = palette.GenerateTriadic(baseColor)
	}
}

func BenchmarkGenerateTetradic(b *testing.B) {
	baseColor, _ := color.ParseHex("#FF5733")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = palette.GenerateTetradic(baseColor)
	}
}

func BenchmarkGenerateAllRules(b *testing.B) {
	baseColor, _ := color.ParseHex("#FF5733")
	rules := []palette.HarmonyRule{
		palette.Monochromatic,
		palette.Complementary,
		palette.Analogous,
		palette.Triadic,
		palette.Tetradic,
		palette.SplitComplementary,
		palette.Square,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, rule := range rules {
			_, _ = palette.Generate(baseColor, rule)
		}
	}
}

// Benchmark WCAG calculations
func BenchmarkRelativeLuminance(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wcag.RelativeLuminance(255, 87, 51)
	}
}

func BenchmarkCalculateContrast(b *testing.B) {
	fg, _ := color.ParseHex("#FFFFFF")
	bg, _ := color.ParseHex("#000000")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wcag.CalculateContrast(fg, bg)
	}
}

func BenchmarkValidate(b *testing.B) {
	fg, _ := color.ParseHex("#FFFFFF")
	bg, _ := color.ParseHex("#000000")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wcag.Validate(fg, bg)
	}
}

// Benchmark color search
func BenchmarkSearchColors(b *testing.B) {
	// Ensure colors are loaded
	_ = color.LoadNamedColors()

	b.Run("ExactMatch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = color.SearchColors("red")
		}
	})

	b.Run("PrefixMatch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = color.SearchColors("blu")
		}
	})

	b.Run("ContainsMatch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = color.SearchColors("sky")
		}
	})

	b.Run("NoMatch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = color.SearchColors("xyzabc")
		}
	})
}

// Benchmark palette contrast validation
func BenchmarkValidatePaletteContrast(b *testing.B) {
	baseColor, _ := color.ParseHex("#FF5733")
	pal, _ := palette.Generate(baseColor, palette.Monochromatic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = palette.ValidatePaletteContrast(pal)
	}
}

// Benchmark complete workflow
func BenchmarkCompleteWorkflow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Parse color
		baseColor, _ := color.ParseHex("#FF5733")

		// Generate palette
		pal, _ := palette.Generate(baseColor, palette.Triadic)

		// Validate contrast
		_, _ = palette.ValidatePaletteContrast(pal)

		// Check WCAG for first two colors
		if len(pal.Colors) >= 2 {
			_ = wcag.Validate(pal.Colors[0], pal.Colors[1])
		}
	}
}
