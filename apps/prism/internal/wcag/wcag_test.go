package wcag

import (
	"math"
	"testing"

	color "github.com/kyanite/prism/internal/color"
)

func TestWCAGContrast(t *testing.T) {
	tests := []struct {
		name      string
		fg        string
		bg        string
		wantRatio float64
		wantLevel string
		tolerance float64
	}{
		{
			name:      "White on Black",
			fg:        "#FFFFFF",
			bg:        "#000000",
			wantRatio: 21.0,
			wantLevel: "AAA",
			tolerance: 0.05,
		},
		{
			name:      "Black on White",
			fg:        "#000000",
			bg:        "#FFFFFF",
			wantRatio: 21.0,
			wantLevel: "AAA",
			tolerance: 0.05,
		},
		{
			name:      "Electric Pink on Dark",
			fg:        "#FF0080",
			bg:        "#0D0221",
			wantRatio: 5.2,
			wantLevel: "AA",
			tolerance: 0.3,
		},
		{
			name:      "Gray on White (AA)",
			fg:        "#767676",
			bg:        "#FFFFFF",
			wantRatio: 4.5,
			wantLevel: "AA",
			tolerance: 0.1,
		},
		{
			name:      "Same Color (FAIL)",
			fg:        "#FF0000",
			bg:        "#FF0000",
			wantRatio: 1.0,
			wantLevel: "FAIL",
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fg, _ := color.ParseHex(tt.fg)
			bg, _ := color.ParseHex(tt.bg)

			ratio := CalculateContrast(fg, bg)
			if math.Abs(ratio-tt.wantRatio) > tt.tolerance {
				t.Errorf("CalculateContrast() = %.2f, want %.2f (±%.2f)", ratio, tt.wantRatio, tt.tolerance)
			}

			result := Validate(fg, bg)
			if result.Level != tt.wantLevel {
				t.Errorf("Level = %s, want %s", result.Level, tt.wantLevel)
			}
		})
	}
}

func TestWCAGThresholds(t *testing.T) {
	tests := []struct {
		ratio        float64
		wantAASmall  bool
		wantAALarge  bool
		wantAAASmall bool
		wantAAALarge bool
	}{
		{21.0, true, true, true, true},    // Maximum
		{7.0, true, true, true, true},     // AAA small threshold
		{4.5, true, true, false, true},    // AA small / AAA large threshold
		{3.0, false, true, false, false},  // AA large threshold
		{1.0, false, false, false, false}, // Minimum (fail all)
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := IsPassingAASmall(tt.ratio); got != tt.wantAASmall {
				t.Errorf("IsPassingAASmall(%.1f) = %v, want %v", tt.ratio, got, tt.wantAASmall)
			}
			if got := IsPassingAALarge(tt.ratio); got != tt.wantAALarge {
				t.Errorf("IsPassingAALarge(%.1f) = %v, want %v", tt.ratio, got, tt.wantAALarge)
			}
			if got := IsPassingAAASmall(tt.ratio); got != tt.wantAAASmall {
				t.Errorf("IsPassingAAASmall(%.1f) = %v, want %v", tt.ratio, got, tt.wantAAASmall)
			}
			if got := IsPassingAAALarge(tt.ratio); got != tt.wantAAALarge {
				t.Errorf("IsPassingAAALarge(%.1f) = %v, want %v", tt.ratio, got, tt.wantAAALarge)
			}
		})
	}
}

func TestRelativeLuminance(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b int
		want    float64
	}{
		{"White", 255, 255, 255, 1.0},
		{"Black", 0, 0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RelativeLuminance(tt.r, tt.g, tt.b)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("RelativeLuminance(%d, %d, %d) = %.3f, want %.3f", tt.r, tt.g, tt.b, got, tt.want)
			}
		})
	}
}

func TestContrastResultSummary(t *testing.T) {
	result := ContrastResult{
		Ratio:     7.5,
		Level:     "AAA",
		PassedAA:  true,
		PassedAAA: true,
	}

	summary := result.Summary()
	expected := "7.50:1 - WCAG AAA"
	if summary != expected {
		t.Errorf("Summary() = %q, want %q", summary, expected)
	}
}
