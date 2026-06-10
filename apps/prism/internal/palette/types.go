// Package palette provides color palette generation using color theory harmony rules.
// It implements classical color harmony algorithms including monochromatic, complementary,
// analogous, triadic, tetradic, split-complementary, and square harmonies.
package palette

import (
	"time"

	"github.com/kyanite/prism/internal/color"
)

// Palette represents a color palette with metadata and harmony information.
// It contains a collection of colors generated from a base color using a specific harmony rule.
type Palette struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Colors      []color.Color `json:"colors"`
	HarmonyRule string        `json:"harmony_rule"`
	BaseColor   color.Color   `json:"base_color"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Tags        []string      `json:"tags,omitempty"`
}

// HarmonyRule defines color harmony rules
type HarmonyRule string

const (
	Monochromatic      HarmonyRule = "monochromatic"
	Complementary      HarmonyRule = "complementary"
	Analogous          HarmonyRule = "analogous"
	Triadic            HarmonyRule = "triadic"
	Tetradic           HarmonyRule = "tetradic"
	SplitComplementary HarmonyRule = "split_complementary"
	Square             HarmonyRule = "square"
)

// AllRules returns all available harmony rules
func AllRules() []HarmonyRule {
	return []HarmonyRule{
		Monochromatic,
		Complementary,
		Analogous,
		Triadic,
		Tetradic,
		SplitComplementary,
		Square,
	}
}
