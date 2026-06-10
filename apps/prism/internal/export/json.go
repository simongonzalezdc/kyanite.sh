package export

import (
	"encoding/json"

	"github.com/kyanite/prism/internal/palette"
)

// ExportJSON exports palette to JSON format
func ExportJSON(p palette.Palette) ([]byte, error) {
	// Format matches PRD specification (RFC 8259 compliant)
	return json.MarshalIndent(p, "", "  ")
}

// ExportJSONCompact exports palette to compact JSON
func ExportJSONCompact(p palette.Palette) ([]byte, error) {
	return json.Marshal(p)
}
