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
	Chords   []string       `json:"chords,omitempty"`
	Lyrics   string         `json:"lyrics,omitempty"`
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

// ExportType represents different types of exports.
const (
	ExportTypePattern ExportType = iota
	ExportTypeLyrics
	ExportTypeChords
	ExportTypeFull
	ExportTypeMarkdown
	ExportTypePlainText
	ExportTypeChordPro
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
	case ExportTypeMarkdown:
		return "markdown"
	case ExportTypePlainText:
		return "plaintext"
	case ExportTypeChordPro:
		return "chordpro"
	default:
		return "unknown"
	}
}

// ExportOptions contains options for the export
type ExportOptions struct {
	Title        string
	Type         ExportType
	OutputPath   string
	BPM          int
	IncludeNotes bool
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
