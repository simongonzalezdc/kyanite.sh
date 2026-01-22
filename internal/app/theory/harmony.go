package theory

import (
	"strings"
)

// Chord represents a musical chord
type Chord struct {
	Root       string
	Type       string // e.g., "maj7", "m9", "sus4"
	Extensions []string
}

// Scale represents a musical scale
type Scale struct {
	Name  string
	Root  string
	Notes []string
}

// HarmonyLib provides music theory validation and inspiration
type HarmonyLib struct {
	// Standard chord database
	validRoots []string
	validTypes []string
}

// NewHarmonyLib creates a new harmony library
func NewHarmonyLib() *HarmonyLib {
	return &HarmonyLib{
		validRoots: []string{"C", "C#", "Db", "D", "D#", "Eb", "E", "F", "F#", "Gb", "G", "G#", "Ab", "A", "A#", "Bb", "B"},
		validTypes: []string{"maj", "min", "maj7", "m7", "7", "dim", "aug", "sus2", "sus4", "maj9", "m9", "add9", "11", "13"},
	}
}

// ValidateChord checks if a chord string is theoretically valid
func (h *HarmonyLib) ValidateChord(chordStr string) bool {
	if chordStr == "" {
		return false
	}

	// Basic parsing: Root + Type
	// This is a naive implementation, real world is more complex
	root := ""
	for _, r := range h.validRoots {
		if strings.HasPrefix(chordStr, r) {
			if len(r) > len(root) {
				root = r
			}
		}
	}

	if root == "" {
		return false
	}

	suffix := strings.TrimPrefix(chordStr, root)
	if suffix == "" {
		return true // Basic major chord
	}

	// Handle "m" for minor
	if strings.HasPrefix(suffix, "m") && !strings.HasPrefix(suffix, "maj") {
		suffix = "min" + strings.TrimPrefix(suffix, "m")
	}

	for _, t := range h.validTypes {
		if strings.Contains(suffix, t) {
			return true
		}
	}

	// Allow some common extensions even if not in validTypes
	extensions := []string{"6", "b5", "#5", "b9", "#9", "#11", "b13"}
	for _, ext := range extensions {
		if strings.Contains(suffix, ext) {
			return true
		}
	}

	return false
}

// GetScaleSuggestion provides a scale suggestion based on a mood
func (h *HarmonyLib) GetScaleSuggestion(mood string) Scale {
	switch strings.ToLower(mood) {
	case "happy", "bright", "joyful":
		return Scale{Name: "Major", Root: "C", Notes: []string{"C", "D", "E", "F", "G", "A", "B"}}
	case "sad", "dark", "melancholy":
		return Scale{Name: "Natural Minor", Root: "A", Notes: []string{"A", "B", "C", "D", "E", "F", "G"}}
	case "tense", "mysterious", "industrial":
		return Scale{Name: "Phrygian", Root: "E", Notes: []string{"E", "F", "G", "A", "B", "C", "D"}}
	case "chill", "jazzy", "soulful":
		return Scale{Name: "Dorian", Root: "D", Notes: []string{"D", "E", "F", "G", "A", "B", "C"}}
	default:
		return Scale{Name: "Chromatic", Root: "C", Notes: []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}}
	}
}

// FormatChord ensures a chord string is in a standard format
func (h *HarmonyLib) FormatChord(chordStr string) string {
	// TODO: Implement normalization logic
	return strings.Title(chordStr)
}

// SuggestProgression returns a progression based on mood and logic
func (h *HarmonyLib) SuggestProgression(mood string) []string {
	// This will be augmented by AI, but we provide base logic here
	switch strings.ToLower(mood) {
	case "happy":
		return []string{"I", "IV", "V", "I"} // Placeholder for real chords
	case "sad":
		return []string{"im", "VI", "III", "VII"}
	default:
		return []string{"I", "V", "vim", "IV"}
	}
}
