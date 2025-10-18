package export

import (
	"time"
)

// NoiseExport represents the main export structure for noise.sh
type NoiseExport struct {
	Type     string         `json:"type"`
	Version  string         `json:"version"`
	Metadata ExportMetadata `json:"metadata"`
	Patterns []string       `json:"patterns"`
	Lyrics   string         `json:"lyrics,omitempty"`
	Chords   []string       `json:"chords,omitempty"`
	Notes    string         `json:"notes,omitempty"`
}

// ExportMetadata contains metadata about the export
type ExportMetadata struct {
	Title      string    `json:"title"`
	BPM        int       `json:"bpm"`
	Created    time.Time `json:"created"`
	ExportedAt time.Time `json:"exported_at"`
}

// ExportType represents different types of exports
type ExportType int

const (
	ExportTypePattern ExportType = iota
	ExportTypeLyrics
	ExportTypeChords
	ExportTypeFull
)

// String returns a string representation of the export type
func (et ExportType) String() string {
	switch et {
	case ExportTypePattern:
		return "pattern"
	case ExportTypeLyrics:
		return "lyrics"
	case ExportTypeChords:
		return "chords"
	case ExportTypeFull:
		return "full"
	default:
		return "unknown"
	}
}

// ExportOptions contains options for the export
type ExportOptions struct {
	Type        ExportType
	Title       string
	BPM         int
	IncludeNotes bool
	OutputPath  string
}

// DefaultExportOptions returns default export options
func DefaultExportOptions() *ExportOptions {
	return &ExportOptions{
		Type:         ExportTypePattern,
		Title:        "Untitled",
		BPM:          120,
		IncludeNotes: false,
		OutputPath:   "export.json",
	}
}