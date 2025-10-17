package app

import (
	"strings"
)

// TheoryService handles music theory operations
type TheoryService struct{}

// NewTheoryService creates a new theory service
func NewTheoryService() *TheoryService {
	return &TheoryService{}
}

// GetScale returns notes for a given scale
func (s *TheoryService) GetScale(key string, scaleType string) ([]string, error) {
	// Placeholder implementation
	// In a full implementation, this would use the go-music-theory library
	scales := map[string][]string{
		"C:major": {"C", "D", "E", "F", "G", "A", "B"},
		"G:major": {"G", "A", "B", "C", "D", "E", "F#"},
		"D:major": {"D", "E", "F#", "G", "A", "B", "C#"},
	}

	keyScale := key + ":" + scaleType
	if scale, exists := scales[keyScale]; exists {
		return scale, nil
	}

	// Default to C major
	return scales["C:major"], nil
}

// GetChord returns notes for a given chord
func (s *TheoryService) GetChord(root string, chordType string) ([]string, error) {
	// Placeholder implementation
	chords := map[string][]string{
		"C:major": {"C", "E", "G"},
		"G:major": {"G", "B", "D"},
		"D:major": {"D", "F#", "A"},
		"A:minor": {"A", "C", "E"},
		"E:minor": {"E", "G", "B"},
		"B:minor": {"B", "D", "F#"},
	}

	chord := root + ":" + chordType
	if notes, exists := chords[chord]; exists {
		return notes, nil
	}

	// Default to C major
	return chords["C:major"], nil
}

// GetProgression returns a chord progression for a key
func (s *TheoryService) GetProgression(key string, pattern string) ([]string, error) {
	// Placeholder implementation
	progressions := map[string][]string{
		"C:I-V-vi-IV": {"C", "G", "Am", "F"},
		"G:I-V-vi-IV": {"G", "D", "Em", "C"},
		"A:I-V-vi-IV": {"A", "E", "F#m", "D"},
	}

	progKey := key + ":" + pattern
	if progression, exists := progressions[progKey]; exists {
		return progression, nil
	}

	// Default to I-V-vi-IV in C
	return progressions["C:I-V-vi-IV"], nil
}

// FindRhymes finds rhyming words for a given word
func (s *TheoryService) FindRhymes(word string) ([]string, error) {
	// Placeholder implementation
	// In a full implementation, this would use the pronouncing library
	rhymes := map[string][]string{
		"love":  {"dove", "glove", "above", "shove"},
		"time":  {"rhyme", "climb", "dime", "lime"},
		"heart": {"start", "part", "art", "chart"},
		"night": {"light", "right", "fight", "bright"},
	}

	if rhymeList, exists := rhymes[word]; exists {
		return rhymeList, nil
	}

	// Return empty slice if no rhymes found
	return []string{}, nil
}

// CountSyllables counts syllables in a word
func (s *TheoryService) CountSyllables(word string) (int, error) {
	// Placeholder implementation
	// In a full implementation, this would use the pronouncing library

	// Simple heuristic-based syllable counting
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
	// Placeholder implementation for prosody analysis
	// In a full implementation, this would analyze stress patterns

	words := strings.Fields(line)
	syllableCount := 0

	for _, word := range words {
		syllables, _ := s.CountSyllables(word)
		syllableCount += syllables
	}

	return syllableCount, nil
}
