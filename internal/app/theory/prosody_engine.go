package theory

import (
	"regexp"
	"strings"
)

// Stress represents the rhythmic weight of a syllable
type Stress int

const (
	Unstressed Stress = iota
	Stressed
	Secondary
)

// Syllable represents a single unit of pronunciation
type Syllable struct {
	Text   string
	Stress Stress
}

// LineAnalysis contains the rhythmic breakdown of a single line of lyrics
type LineAnalysis struct {
	Text      string
	Syllables []Syllable
	Meter     string // e.g., "Iambic Pentameter"
	Count     int
}

// ProsodyEngine handles rhythmic and melodic analysis of lyrics
type ProsodyEngine struct {
	// Cache for syllables to improve performance
	cache map[string]int
}

// NewProsodyEngine creates a new prosody engine
func NewProsodyEngine() *ProsodyEngine {
	return &ProsodyEngine{
		cache: make(map[string]int),
	}
}

// CountSyllables estimates the number of syllables in a line of text
func (e *ProsodyEngine) CountSyllables(text string) int {
	words := strings.Fields(text)
	total := 0
	for _, word := range words {
		total += e.countWordSyllables(word)
	}
	return total
}

// AnalyzeLine performs a detailed rhythmic analysis of a line
func (e *ProsodyEngine) AnalyzeLine(text string) LineAnalysis {
	words := strings.Fields(text)
	var syllables []Syllable
	totalCount := 0

	for _, word := range words {
		count := e.countWordSyllables(word)
		totalCount += count

		// TODO: Implement more sophisticated stress detection
		// For now, use a simplified alternating pattern heuristic
		for i := 0; i < count; i++ {
			stress := Unstressed
			if (totalCount-count+i)%2 == 1 {
				stress = Stressed
			}
			syllables = append(syllables, Syllable{
				Text:   word, // In a full implementation, we'd split the word
				Stress: stress,
			})
		}
	}

	return LineAnalysis{
		Text:      text,
		Syllables: syllables,
		Count:     totalCount,
		Meter:     e.detectMeter(syllables),
	}
}

// countWordSyllables uses a heuristic to count syllables in an English word
func (e *ProsodyEngine) countWordSyllables(word string) int {
	word = strings.ToLower(regexp.MustCompile(`[^a-z]`).ReplaceAllString(word, ""))
	if len(word) == 0 {
		return 0
	}
	if len(word) <= 3 {
		return 1
	}

	if count, ok := e.cache[word]; ok {
		return count
	}

	// Basic heuristic: count vowel groups
	vowelGroups := regexp.MustCompile(`[aeiouy]+`)
	matches := vowelGroups.FindAllString(word, -1)
	count := len(matches)

	// Adjustments
	if strings.HasSuffix(word, "e") {
		// Silent 'e' at the end (usually)
		if !strings.HasSuffix(word, "le") && count > 1 {
			count--
		}
	}

	// Handle common prefixes/suffixes
	if strings.HasSuffix(word, "ing") {
		// Already counted in vowel groups, but ensure it's distinct
	}

	if count < 1 {
		count = 1
	}

	e.cache[word] = count
	return count
}

// detectMeter attempts to identify the poetic meter of the line
func (e *ProsodyEngine) detectMeter(syllables []Syllable) string {
	if len(syllables) == 0 {
		return "Unknown"
	}

	// Simplified meter detection
	// Iambic: unstressed-stressed
	// Trochaic: stressed-unstressed

	iambicScore := 0
	trochaicScore := 0

	for i, s := range syllables {
		if (i%2 == 0 && s.Stress == Unstressed) || (i%2 == 1 && s.Stress == Stressed) {
			iambicScore++
		}
		if (i%2 == 0 && s.Stress == Stressed) || (i%2 == 1 && s.Stress == Unstressed) {
			trochaicScore++
		}
	}

	if iambicScore > trochaicScore && iambicScore > len(syllables)/2 {
		return "Iambic"
	}
	if trochaicScore > iambicScore && trochaicScore > len(syllables)/2 {
		return "Trochaic"
	}

	return "Irregular"
}
