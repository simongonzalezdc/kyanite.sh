// Package theory provides music theory helpers for harmony and prosody analysis.
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

// FormatChord ensures a chord string is in a standard format.
// Normalizes chord notation to a consistent, widely-recognized format.
func (h *HarmonyLib) FormatChord(chordStr string) string {
	if chordStr == "" {
		return ""
	}

	// Trim whitespace
	chordStr = strings.TrimSpace(chordStr)

	// Handle slash chords (e.g., "Am/E")
	var bassNote string
	if slashIdx := strings.Index(chordStr, "/"); slashIdx > 0 {
		bassNote = chordStr[slashIdx+1:]
		chordStr = chordStr[:slashIdx]
	}

	// Extract root note
	root := h.extractRoot(chordStr)
	if root == "" {
		return chordStr // Return original if we can't parse it
	}

	// Normalize the root note
	root = h.normalizeRoot(root)

	// Get the suffix (chord type)
	suffix := strings.TrimPrefix(chordStr, h.extractRoot(chordStr))
	suffix = h.normalizeSuffix(suffix)

	// Reconstruct chord
	result := root + suffix

	// Add bass note if present
	if bassNote != "" {
		bassNote = h.normalizeRoot(strings.TrimSpace(bassNote))
		result += "/" + bassNote
	}

	return result
}

// extractRoot extracts the root note from a chord string
func (h *HarmonyLib) extractRoot(chordStr string) string {
	if len(chordStr) == 0 {
		return ""
	}

	// Check for notes with accidentals first (e.g., "C#", "Bb")
	if len(chordStr) >= 2 {
		twoChar := strings.ToUpper(string(chordStr[0])) + string(chordStr[1])
		for _, r := range h.validRoots {
			if twoChar == r {
				return twoChar
			}
		}
	}

	// Check single character roots
	oneChar := strings.ToUpper(string(chordStr[0]))
	for _, r := range h.validRoots {
		if oneChar == r {
			return oneChar
		}
	}

	return ""
}

// normalizeRoot normalizes a root note to a consistent format
func (h *HarmonyLib) normalizeRoot(root string) string {
	if len(root) == 0 {
		return ""
	}

	// Capitalize the letter
	root = strings.ToUpper(string(root[0])) + root[1:]

	// Normalize accidentals: prefer sharps in some keys, flats in others
	// Standard convention: use sharps for C, G, D, A, E keys; flats for F, Bb, Eb, Ab, Db keys
	// For simplicity, we'll normalize to the more common enharmonic spelling

	// Common enharmonic normalizations
	enharmonicMap := map[string]string{
		"Cb": "B",
		"E#": "F",
		"Fb": "E",
		"B#": "C",
		// Keep Db, Eb, Ab, Bb, Gb as flats (common in jazz/pop)
		// Keep C#, F#, G# as sharps
	}

	if normalized, ok := enharmonicMap[root]; ok {
		return normalized
	}

	return root
}

// normalizeSuffix normalizes chord type suffixes to standard notation
func (h *HarmonyLib) normalizeSuffix(suffix string) string {
	// Handle empty suffix (major chord)
	if suffix == "" {
		return ""
	}

	// Normalize common variations to standard notation
	normalizations := []struct {
		patterns []string
		standard string
	}{
		// Minor variations
		{[]string{"minor", "MIN", "Min", "-"}, "m"},
		// Major 7 variations
		{[]string{"major7", "MAJOR7", "Major7", "maj7", "Maj7", "MAJ7", "Δ7", "M7"}, "maj7"},
		// Major variations (for explicit major notation)
		{[]string{"major", "MAJOR", "Major", "maj", "MAJ", "Maj", "M", "Δ"}, ""},
		// Minor 7 variations
		{[]string{"minor7", "MINOR7", "Minor7", "min7", "Min7", "MIN7", "-7"}, "m7"},
		// Diminished variations
		{[]string{"diminished", "DIMINISHED", "Diminished", "dim", "DIM", "Dim", "°", "o"}, "dim"},
		// Augmented variations
		{[]string{"augmented", "AUGMENTED", "Augmented", "aug", "AUG", "Aug", "+"}, "aug"},
		// Suspended variations
		{[]string{"sus2", "SUS2", "Sus2", "suspended2"}, "sus2"},
		{[]string{"sus4", "SUS4", "Sus4", "suspended4", "sus"}, "sus4"},
		// 7th variations
		{[]string{"dom7", "DOM7", "Dom7"}, "7"},
		// 9th variations
		{[]string{"major9", "MAJOR9", "Major9"}, "maj9"},
		{[]string{"minor9", "MINOR9", "Minor9", "min9", "Min9"}, "m9"},
		// 6th
		{[]string{"major6", "MAJOR6", "Major6", "maj6"}, "6"},
		{[]string{"minor6", "MINOR6", "Minor6", "min6"}, "m6"},
		// Add9
		{[]string{"add9", "ADD9", "Add9", "add2"}, "add9"},
		// 11th and 13th
		{[]string{"11th", "11TH"}, "11"},
		{[]string{"13th", "13TH"}, "13"},
		// Half-diminished (minor 7 flat 5)
		{[]string{"m7b5", "m7♭5", "ø", "ø7", "half-dim", "halfdim"}, "m7b5"},
		// Diminished 7
		{[]string{"dim7", "°7", "o7"}, "dim7"},
	}

	result := suffix
	for _, n := range normalizations {
		for _, pattern := range n.patterns {
			if strings.EqualFold(result, pattern) || strings.HasPrefix(strings.ToLower(result), strings.ToLower(pattern)) {
				// Replace the pattern with standard notation
				remainder := result[len(pattern):]
				result = n.standard + remainder
				break
			}
		}
	}

	return result
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
