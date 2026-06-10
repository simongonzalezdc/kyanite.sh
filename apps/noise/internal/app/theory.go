package app

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Kyanite/noise/internal/logging"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TheoryService handles music theory operations
type TheoryService struct {
	dictionary *Dictionary
}

// ChordInfo represents information about a chord
type ChordInfo struct {
	Root        string   `json:"root"`
	Quality     string   `json:"quality"`
	Notes       []string `json:"notes"`
	Intervals   []string `json:"intervals"`
	Description string   `json:"description"`
}

// ScaleInfo represents information about a scale
type ScaleInfo struct {
	Name        string   `json:"name"`
	Notes       []string `json:"notes"`
	Pattern     []string `json:"pattern"`
	Description string   `json:"description"`
}

// ModeInfo represents information about a musical mode
type ModeInfo struct {
	Name        string   `json:"name"`
	ParentScale string   `json:"parent_scale"`
	Notes       []string `json:"notes"`
	Pattern     []string `json:"pattern"`
	Description string   `json:"description"`
}

// ChordAnalysis represents the result of chord analysis
type ChordAnalysis struct {
	Input          string      `json:"input"`
	DetectedChords []ChordInfo `json:"detected_chords"`
	Suggestions    []string    `json:"suggestions"`
	Key            string      `json:"key,omitempty"`
}

// NewTheoryService creates a new theory service
func NewTheoryService() *TheoryService {
	dict := NewDictionary()
	// Try to load the dictionary file
	if err := dict.LoadDictionary("data/dictionary.json"); err != nil {
		// If loading fails, continue with empty dictionary (will use fallbacks)
		// In production, we log this error for visibility
		logging.Warnf("TheoryService: failed to load dictionary: %v", err)
	}

	return &TheoryService{
		dictionary: dict,
	}
}

// GetScale returns notes for a given scale
func (s *TheoryService) GetScale(key string, scaleType string) ([]string, error) {
	scaleInfo, err := s.GetScaleInfo(key, scaleType)
	if err != nil {
		return nil, err
	}
	return scaleInfo.Notes, nil
}

// GetScaleInfo returns detailed information about a scale
func (s *TheoryService) GetScaleInfo(key string, scaleType string) (*ScaleInfo, error) {
	// Normalize key (capitalize first letter)
	c := cases.Title(language.Und)
	key = c.String(strings.ToLower(key))

	// Scale patterns (semitones from root)
	scalePatterns := map[string][]int{
		"major":      {0, 2, 4, 5, 7, 9, 11},
		"minor":      {0, 2, 3, 5, 7, 8, 10},
		"harmonic":   {0, 2, 3, 5, 7, 8, 11},
		"melodic":    {0, 2, 3, 5, 7, 9, 11},
		"pentatonic": {0, 2, 4, 7, 9},
		"blues":      {0, 3, 5, 6, 7, 10},
		"dorian":     {0, 2, 3, 5, 7, 9, 10},
		"phrygian":   {0, 1, 3, 5, 7, 8, 10},
		"lydian":     {0, 2, 4, 6, 7, 9, 11},
		"mixolydian": {0, 2, 4, 5, 7, 9, 10},
		"locrian":    {0, 1, 3, 5, 6, 8, 10},
	}

	pattern, exists := scalePatterns[strings.ToLower(scaleType)]
	if !exists {
		return nil, fmt.Errorf("unknown scale type: %s", scaleType)
	}

	// Note names in chromatic scale
	notes := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

	// Find root note index
	rootIndex := -1
	for i, note := range notes {
		if note == key || (len(key) > 1 && note == key+"#") {
			rootIndex = i
			break
		}
	}

	if rootIndex == -1 {
		return nil, fmt.Errorf("unknown key: %s", key)
	}

	// Generate scale notes
	scaleNotes := make([]string, len(pattern))
	intervalNames := make([]string, len(pattern))

	intervalMap := map[int]string{
		0: "Root", 1: "Minor 2nd", 2: "Major 2nd", 3: "Minor 3rd",
		4: "Major 3rd", 5: "Perfect 4th", 6: "Tritone", 7: "Perfect 5th",
		8: "Minor 6th", 9: "Major 6th", 10: "Minor 7th", 11: "Major 7th",
	}

	for i, semitone := range pattern {
		noteIndex := (rootIndex + semitone) % 12
		scaleNotes[i] = notes[noteIndex]
		intervalNames[i] = intervalMap[semitone]
	}

	description := fmt.Sprintf("%s %s scale", key, c.String(scaleType))

	return &ScaleInfo{
		Name:        fmt.Sprintf("%s %s", key, scaleType),
		Notes:       scaleNotes,
		Pattern:     intervalNames,
		Description: description,
	}, nil
}

// GetChord returns notes for a given chord
func (s *TheoryService) GetChord(root string, chordType string) ([]string, error) {
	chordInfo, err := s.GetChordInfo(root, chordType)
	if err != nil {
		return nil, err
	}
	return chordInfo.Notes, nil
}

// GetChordInfo returns detailed information about a chord
func (s *TheoryService) GetChordInfo(root string, chordType string) (*ChordInfo, error) {
	// Normalize root note
	c := cases.Title(language.Und)
	root = c.String(strings.ToLower(root))

	// Chord patterns (semitones from root)
	chordPatterns := map[string]struct {
		Intervals   []int
		Description string
	}{
		"major":    {[]int{0, 4, 7}, "Major triad"},
		"maj":      {[]int{0, 4, 7}, "Major triad"},
		"M":        {[]int{0, 4, 7}, "Major triad"},
		"minor":    {[]int{0, 3, 7}, "Minor triad"},
		"min":      {[]int{0, 3, 7}, "Minor triad"},
		"m":        {[]int{0, 3, 7}, "Minor triad"},
		"dim":      {[]int{0, 3, 6}, "Diminished triad"},
		"aug":      {[]int{0, 4, 8}, "Augmented triad"},
		"sus2":     {[]int{0, 2, 7}, "Suspended 2nd"},
		"sus4":     {[]int{0, 5, 7}, "Suspended 4th"},
		"maj7":     {[]int{0, 4, 7, 11}, "Major 7th"},
		"7":        {[]int{0, 4, 7, 10}, "Dominant 7th"},
		"min7":     {[]int{0, 3, 7, 10}, "Minor 7th"},
		"m7":       {[]int{0, 3, 7, 10}, "Minor 7th"},
		"dim7":     {[]int{0, 3, 6, 9}, "Diminished 7th"},
		"halfdim7": {[]int{0, 3, 6, 10}, "Half-diminished 7th"},
		"m7b5":     {[]int{0, 3, 6, 10}, "Half-diminished 7th"},
		"maj9":     {[]int{0, 4, 7, 11, 14}, "Major 9th"},
		"9":        {[]int{0, 4, 7, 10, 14}, "Dominant 9th"},
		"min9":     {[]int{0, 3, 7, 10, 14}, "Minor 9th"},
		"add9":     {[]int{0, 4, 7, 14}, "Add 9"},
		"6":        {[]int{0, 4, 7, 9}, "Major 6th"},
		"min6":     {[]int{0, 3, 7, 9}, "Minor 6th"},
	}

	pattern, exists := chordPatterns[strings.ToLower(chordType)]
	if !exists {
		return nil, fmt.Errorf("unknown chord type: %s", chordType)
	}

	// Note names in chromatic scale
	notes := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

	// Find root note index
	rootIndex := -1
	for i, note := range notes {
		if note == root {
			rootIndex = i
			break
		}
	}

	if rootIndex == -1 {
		return nil, fmt.Errorf("unknown root note: %s", root)
	}

	// Generate chord notes
	chordNotes := make([]string, len(pattern.Intervals))
	intervalNames := make([]string, len(pattern.Intervals))

	intervalMap := map[int]string{
		0: "Root", 1: "Minor 2nd", 2: "Major 2nd", 3: "Minor 3rd",
		4: "Major 3rd", 5: "Perfect 4th", 6: "Tritone", 7: "Perfect 5th",
		8: "Minor 6th", 9: "Major 6th", 10: "Minor 7th", 11: "Major 7th",
		12: "Octave", 14: "Major 9th",
	}

	for i, semitone := range pattern.Intervals {
		noteIndex := (rootIndex + semitone) % 12
		chordNotes[i] = notes[noteIndex]
		intervalNames[i] = intervalMap[semitone]
	}

	return &ChordInfo{
		Root:        root,
		Quality:     chordType,
		Notes:       chordNotes,
		Intervals:   intervalNames,
		Description: fmt.Sprintf("%s %s: %s", root, chordType, pattern.Description),
	}, nil
}

// GetProgression returns a chord progression for a key
func (s *TheoryService) GetProgression(key string, pattern string) ([]string, error) {
	progressionInfo, err := s.GetProgressionInfo(key, pattern)
	if err != nil {
		return nil, err
	}
	return progressionInfo.Chords, nil
}

// ProgressionInfo represents a chord progression
type ProgressionInfo struct {
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Chords      []string `json:"chords"`
	Description string   `json:"description"`
}

// GetProgressionInfo returns detailed information about a chord progression
func (s *TheoryService) GetProgressionInfo(key string, pattern string) (*ProgressionInfo, error) {
	// Normalize key
	c := cases.Title(language.Und)
	key = c.String(strings.ToLower(key))

	// Common chord progressions
	progressions := map[string]struct {
		Chords      []string
		Description string
	}{
		"I-V-vi-IV": {
			[]string{"I", "V", "vi", "IV"},
			"The most popular chord progression in modern music",
		},
		"ii-V-I": {
			[]string{"ii", "V", "I"},
			"Classic jazz progression",
		},
		"I-IV-V-IV": {
			[]string{"I", "IV", "V", "IV"},
			"Blues progression",
		},
		"I-vi-IV-V": {
			[]string{"I", "vi", "IV", "V"},
			"50s doo-wop progression",
		},
		"vi-IV-I-V": {
			[]string{"vi", "IV", "I", "V"},
			"Pachelbel's Canon progression",
		},
	}

	patternKey := strings.ToUpper(pattern)
	prog, exists := progressions[patternKey]
	if !exists {
		return nil, fmt.Errorf("unknown progression pattern: %s", pattern)
	}

	// Convert roman numerals to actual chords based on key
	chordNames := s.romanToChords(key, prog.Chords)

	return &ProgressionInfo{
		Name:        fmt.Sprintf("%s %s", key, pattern),
		Key:         key,
		Chords:      chordNames,
		Description: prog.Description,
	}, nil
}

// romanToChords converts roman numeral chords to actual chord names
func (s *TheoryService) romanToChords(key string, romanNumerals []string) []string {
	// Key signatures (major keys)
	keySignatures := map[string][]string{
		"C":  {"C", "Dm", "Em", "F", "G", "Am", "Bdim"},
		"C#": {"C#", "D#m", "Fm", "F#", "G#", "A#m", "Cdim"},
		"D":  {"D", "Em", "F#m", "G", "A", "Bm", "C#dim"},
		"Eb": {"Eb", "Fm", "Gm", "Ab", "Bb", "Cm", "Ddim"},
		"E":  {"E", "F#m", "G#m", "A", "B", "C#m", "D#dim"},
		"F":  {"F", "Gm", "Am", "Bb", "C", "Dm", "Edim"},
		"F#": {"F#", "G#m", "A#m", "B", "C#", "D#m", "Fdim"},
		"G":  {"G", "Am", "Bm", "C", "D", "Em", "F#dim"},
		"Ab": {"Ab", "Bbm", "Cm", "Db", "Eb", "Fm", "Gdim"},
		"A":  {"A", "Bm", "C#m", "D", "E", "F#m", "G#dim"},
		"Bb": {"Bb", "Cm", "Dm", "Eb", "F", "Gm", "Adim"},
		"B":  {"B", "C#m", "D#m", "E", "F#", "G#m", "A#dim"},
	}

	chords, exists := keySignatures[key]
	if !exists {
		// Default to C major if key not found
		chords = keySignatures["C"]
	}

	result := make([]string, len(romanNumerals))
	romanToIndex := map[string]int{
		"I": 0, "i": 0,
		"II": 1, "ii": 1,
		"III": 2, "iii": 2,
		"IV": 3, "iv": 3,
		"V": 4, "v": 4,
		"VI": 5, "vi": 5,
		"VII": 6, "vii": 6,
	}

	for i, roman := range romanNumerals {
		if index, exists := romanToIndex[roman]; exists && index < len(chords) {
			result[i] = chords[index]
		} else {
			result[i] = roman // Keep original if not found
		}
	}

	return result
}

// AnalyzeChords analyzes text input for chord patterns
func (s *TheoryService) AnalyzeChords(input string) (*ChordAnalysis, error) {
	// Regular expression to find chord patterns
	// Matches common chord notation: C, Cm, Cmaj7, C#m7, etc.
	chordRegex := regexp.MustCompile(`\b([A-G])(#|b)?((?:maj|min|dim|aug|m|M|sus|add|dom)?[0-9]*|°|ø)?\b`)

	matches := chordRegex.FindAllStringSubmatch(input, -1)

	var detectedChords []ChordInfo
	var suggestions []string
	keys := make(map[string]bool)

	for _, match := range matches {
		if len(match) >= 2 {
			root := match[1]
			accidental := ""
			quality := "major" // default

			if len(match) > 2 && match[2] != "" {
				accidental = match[2]
			}

			if len(match) > 3 && match[3] != "" {
				quality = match[3]
			}

			// Normalize accidental notation
			if accidental == "#" {
				root = root + "#"
			} else if accidental == "b" {
				root = root + "b"
			}

			// Normalize quality
			switch quality {
			case "m", "min":
				quality = "minor"
			case "M", "maj":
				quality = "major"
			case "°":
				quality = "dim"
			case "ø":
				quality = "halfdim7"
			default:
				if strings.Contains(quality, "maj") {
					quality = "major"
				} else if strings.Contains(quality, "min") || strings.Contains(quality, "m") {
					quality = "minor"
				}
			}

			// Get chord info
			chordInfo, err := s.GetChordInfo(root, quality)
			if err == nil {
				detectedChords = append(detectedChords, *chordInfo)
				keys[root] = true
			}
		}
	}

	// Determine possible keys
	var possibleKeys []string
	for key := range keys {
		possibleKeys = append(possibleKeys, key)
	}
	sort.Strings(possibleKeys)

	// Generate suggestions
	if len(detectedChords) > 0 {
		suggestions = append(suggestions, fmt.Sprintf("Found %d chords in the text", len(detectedChords)))
		if len(possibleKeys) > 0 {
			suggestions = append(suggestions, fmt.Sprintf("Possible keys: %s", strings.Join(possibleKeys, ", ")))
		}
	} else {
		suggestions = append(suggestions, "No chords detected in the text")
	}

	return &ChordAnalysis{
		Input:          input,
		DetectedChords: detectedChords,
		Suggestions:    suggestions,
	}, nil
}

// GetModes returns all modes for a given scale
func (s *TheoryService) GetModes(parentScale string) ([]*ModeInfo, error) {
	modes := []*ModeInfo{
		{
			Name:        "Ionian",
			ParentScale: parentScale,
			Pattern:     []string{"Major", "Major", "Minor", "Major", "Major", "Minor", "Minor"},
			Description: "The major scale - happy, bright sound",
		},
		{
			Name:        "Dorian",
			ParentScale: parentScale,
			Pattern:     []string{"Major", "Minor", "Minor", "Major", "Major", "Minor", "Major"},
			Description: "Minor scale with a major 6th - funky, jazzy sound",
		},
		{
			Name:        "Phrygian",
			ParentScale: parentScale,
			Pattern:     []string{"Minor", "Minor", "Major", "Major", "Minor", "Major", "Major"},
			Description: "Minor scale with a flat 2nd - Spanish, exotic sound",
		},
		{
			Name:        "Lydian",
			ParentScale: parentScale,
			Pattern:     []string{"Major", "Major", "Major", "Minor", "Major", "Minor", "Minor"},
			Description: "Major scale with a sharp 4th - dreamy, floating sound",
		},
		{
			Name:        "Mixolydian",
			ParentScale: parentScale,
			Pattern:     []string{"Major", "Major", "Minor", "Major", "Minor", "Minor", "Major"},
			Description: "Major scale with a flat 7th - bluesy, rock sound",
		},
		{
			Name:        "Aeolian",
			ParentScale: parentScale,
			Pattern:     []string{"Major", "Minor", "Minor", "Major", "Minor", "Minor", "Major"},
			Description: "The natural minor scale - sad, melancholic sound",
		},
		{
			Name:        "Locrian",
			ParentScale: parentScale,
			Pattern:     []string{"Minor", "Minor", "Major", "Minor", "Minor", "Major", "Major"},
			Description: "Diminished scale with a flat 5th - tense, unstable sound",
		},
	}

	// Generate actual notes for each mode based on the parent scale
	_, err := s.GetScaleInfo("C", parentScale) // Use C as reference for validation
	if err != nil {
		return nil, err
	}

	// Mode patterns (semitones from root for each mode)
	modePatterns := [][]int{
		{0, 2, 4, 5, 7, 9, 11}, // Ionian (Major)
		{0, 2, 3, 5, 7, 9, 10}, // Dorian
		{0, 1, 3, 5, 7, 8, 10}, // Phrygian
		{0, 2, 4, 6, 7, 9, 11}, // Lydian
		{0, 2, 4, 5, 7, 9, 10}, // Mixolydian
		{0, 2, 3, 5, 7, 8, 10}, // Aeolian (Natural Minor)
		{0, 1, 3, 5, 6, 8, 10}, // Locrian
	}

	notes := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

	for i, mode := range modes {
		if i < len(modePatterns) {
			pattern := modePatterns[i]
			modeNotes := make([]string, len(pattern))
			for j, semitone := range pattern {
				noteIndex := semitone % 12
				modeNotes[j] = notes[noteIndex]
			}
			mode.Notes = modeNotes
		}
	}

	return modes, nil
}

// GetCommonScales returns a list of commonly used scales
func (s *TheoryService) GetCommonScales() []*ScaleInfo {
	commonScales := []struct {
		Key         string
		ScaleType   string
		Description string
	}{
		{"C", "major", "The most common key - no sharps or flats"},
		{"G", "major", "Popular for folk and country music"},
		{"A", "minor", "Natural minor scale - melancholic sound"},
		{"D", "major", "Bright and happy sounding key"},
		{"E", "minor", "Popular in rock and pop music"},
		{"F", "major", "Warm and mellow sounding key"},
		{"B", "minor", "Rich and expressive minor key"},
		{"C", "pentatonic", "Five-note scale used in blues and rock"},
		{"A", "pentatonic", "Common scale for improvisation"},
		{"C", "blues", "Six-note scale with blue notes"},
	}

	var scales []*ScaleInfo
	for _, cs := range commonScales {
		scaleInfo, err := s.GetScaleInfo(cs.Key, cs.ScaleType)
		if err == nil {
			scaleInfo.Description = cs.Description
			scales = append(scales, scaleInfo)
		}
	}

	return scales
}

// FindRhymes finds rhyming words for a given word
func (s *TheoryService) FindRhymes(word string) ([]string, error) {
	// Use the enhanced dictionary if available
	if s.dictionary != nil && s.dictionary.IsLoaded() {
		return s.dictionary.FindRhymes(word)
	}

	// Fallback to static rhyme dictionary
	rhymes := map[string][]string{
		"love":  {"dove", "glove", "above", "shove", "of", "rough", "tough"},
		"time":  {"rhyme", "climb", "dime", "lime", "crime", "prime", "sublime"},
		"heart": {"start", "part", "art", "chart", "smart", "cart", "depart"},
		"night": {"light", "right", "fight", "bright", "sight", "might", "tight"},
		"blue":  {"true", "new", "few", "shoe", "through", "you", "view"},
		"day":   {"way", "say", "play", "stay", "away", "okay", "display"},
		"world": {"word", "heard", "bird", "third", "curved", "served", "learned"},
		"eyes":  {"skies", "lies", "rise", "wise", "size", "prize", "surprise"},
		"home":  {"alone", "stone", "phone", "tone", "zone", "throne", "unknown"},
		"dream": {"seem", "team", "stream", "cream", "scheme", "theme", "supreme"},
		"rain":  {"pain", "main", "chain", "train", "brain", "explain", "maintain"},
		"song":  {"long", "strong", "wrong", "along", "belong", "throng", "prolong"},
		"fire":  {"desire", "wire", "hire", "inspire", "entire", "require", "acquire"},
		"road":  {"load", "code", "mode", "node", "episode", "overload", "download"},
		"star":  {"far", "car", "bar", "scar", "guitar", "cigar", "bizarre"},
	}

	// Normalize input
	word = strings.ToLower(strings.TrimSpace(word))

	if rhymeList, exists := rhymes[word]; exists {
		return rhymeList, nil
	}

	// Return empty slice if no rhymes found
	return []string{}, nil
}

// CountSyllables counts syllables in a word
func (s *TheoryService) CountSyllables(word string) (int, error) {
	// Use the enhanced dictionary if available
	if s.dictionary != nil && s.dictionary.IsLoaded() {
		return s.dictionary.CountSyllables(word)
	}

	// Fallback to heuristic-based syllable counting
	vowels := "aeiouy"
	syllables := 0
	prevWasVowel := false

	for _, char := range word {
		isVowel := false
		for _, vowel := range vowels {
			if char == vowel {
				isVowel = true
				break
			}
		}

		if isVowel && !prevWasVowel {
			syllables++
		}
		prevWasVowel = isVowel
	}

	// Handle silent 'e'
	if len(word) > 2 && word[len(word)-1] == 'e' && syllables > 1 {
		syllables--
	}

	if syllables == 0 {
		syllables = 1 // Every word has at least one syllable
	}

	return syllables, nil
}

// AnalyzeProsody analyzes the prosody of a line
func (s *TheoryService) AnalyzeProsody(line string) (int, error) {
	// Use the enhanced dictionary if available
	if s.dictionary != nil && s.dictionary.IsLoaded() {
		return s.dictionary.CountSyllablesInText(line)
	}

	// Fallback to word-by-word analysis
	words := strings.Fields(line)
	syllableCount := 0

	for _, word := range words {
		syllables, _ := s.CountSyllables(word)
		syllableCount += syllables
	}

	return syllableCount, nil
}

// GetDictionaryStats returns statistics about the dictionary
func (s *TheoryService) GetDictionaryStats() (DictionaryStats, error) {
	if s.dictionary != nil && s.dictionary.IsLoaded() {
		return s.dictionary.GetStats(), nil
	}

	return DictionaryStats{}, fmt.Errorf("dictionary not loaded")
}

// ValidateWord checks if a word exists in the dictionary
func (s *TheoryService) ValidateWord(word string) bool {
	if s.dictionary != nil && s.dictionary.IsLoaded() {
		return s.dictionary.ValidateWord(word)
	}

	// Fallback: check if it's a non-empty string
	return strings.TrimSpace(word) != ""
}

// SearchWords searches for words matching a pattern
func (s *TheoryService) SearchWords(pattern string, limit int) ([]string, error) {
	if s.dictionary != nil && s.dictionary.IsLoaded() {
		return s.dictionary.SearchWords(pattern, limit)
	}

	return []string{}, fmt.Errorf("dictionary not loaded")
}

// GetWordsBySyllableCount returns words with a specific syllable count
func (s *TheoryService) GetWordsBySyllableCount(syllables int, limit int) ([]string, error) {
	if s.dictionary != nil && s.dictionary.IsLoaded() {
		return s.dictionary.GetWordsBySyllableCount(syllables, limit)
	}

	return []string{}, fmt.Errorf("dictionary not loaded")
}
