package editor

import (
	"strings"
	"unicode"
)

// PreviewStats calculates statistics for preview content
type PreviewStats struct {
	Lines      int
	Words      int
	Characters int
	Syllables  int
	Sentences  int
}

// CalculateStats computes statistics for the given content
func CalculateStats(content string) PreviewStats {
	stats := PreviewStats{}

	if content == "" {
		return stats
	}

	// Count lines
	stats.Lines = strings.Count(content, "\n") + 1

	// Count words, characters, sentences
	inWord := false
	for _, r := range content {
		stats.Characters++

		if unicode.IsSpace(r) {
			inWord = false
		} else {
			if !inWord {
				stats.Words++
				inWord = true
			}
		}

		// Simple sentence detection
		if r == '.' || r == '!' || r == '?' {
			stats.Sentences++
		}
	}

	// Estimate syllables (very basic heuristic)
	stats.Syllables = estimateSyllables(content)

	return stats
}

// estimateSyllables provides a rough syllable count
func estimateSyllables(content string) int {
	// Very basic: count vowel groups
	vowels := "aeiouAEIOU"
	syllableCount := 0
	inVowelGroup := false

	for _, r := range content {
		if strings.ContainsRune(vowels, r) {
			if !inVowelGroup {
				syllableCount++
				inVowelGroup = true
			}
		} else {
			inVowelGroup = false
		}
	}

	// Minimum 1 syllable per word
	words := len(strings.Fields(content))
	if syllableCount < words {
		syllableCount = words
	}

	return syllableCount
}

// ReadingTime estimates reading time in minutes
func (s PreviewStats) ReadingTime() float64 {
	// Average reading speed: 200 words per minute
	if s.Words == 0 {
		return 0
	}
	return float64(s.Words) / 200.0
}

// SpeakingTime estimates speaking time in minutes
func (s PreviewStats) SpeakingTime() float64 {
	// Average speaking speed: 150 words per minute
	if s.Words == 0 {
		return 0
	}
	return float64(s.Words) / 150.0
}

// AverageWordLength calculates average word length
func (s PreviewStats) AverageWordLength() float64 {
	if s.Words == 0 {
		return 0
	}
	return float64(s.Characters) / float64(s.Words)
}

// AverageSentenceLength calculates average sentence length in words
func (s PreviewStats) AverageSentenceLength() float64 {
	if s.Sentences == 0 {
		return 0
	}
	return float64(s.Words) / float64(s.Sentences)
}
