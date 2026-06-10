package export

import (
	"encoding/json"

	"github.com/kyanite/prism/internal/palette"
)

// KyaniteTheme represents the Kyanite theme format
type KyaniteTheme struct {
	Name           string            `json:"name"`
	KyaniteVersion string            `json:"kyanite_version"`
	CreatedAt      string            `json:"created_at"`
	Theme          map[string]string `json:"theme"`
}

// ExportTheme exports palette in Kyanite theme format
func ExportTheme(p palette.Palette) ([]byte, error) {

	theme := KyaniteTheme{
		Name:           p.Name,
		KyaniteVersion: "1.0",
		CreatedAt:      p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Theme:          make(map[string]string),
	}

	// Map palette colors to theme roles
	if len(p.Colors) > 0 {
		theme.Theme["primary"] = p.Colors[0].Hex
	}
	if len(p.Colors) > 1 {
		theme.Theme["secondary"] = p.Colors[1].Hex
	}
	if len(p.Colors) > 2 {
		theme.Theme["accent"] = p.Colors[2].Hex
	}
	if len(p.Colors) > 3 {
		theme.Theme["background"] = p.Colors[3].Hex
	}
	if len(p.Colors) > 4 {
		theme.Theme["text"] = p.Colors[4].Hex
	}
	if len(p.Colors) > 5 {
		theme.Theme["success"] = p.Colors[5].Hex
	}

	return json.MarshalIndent(theme, "", "  ")
}
