package domain

import (
	"math"
	"strings"
)

// CalculateQualityScore calculates the 7-dimension quality score for a song
func (s *Song) CalculateQualityScore() *QualityScore {
	score := &QualityScore{
		Feedback: make(map[string]string),
	}

	// Calculate each dimension
	score.Specificity = s.calculateSpecificity()
	score.Originality = s.calculateOriginality()
	score.EmotionalResonance = s.calculateEmotionalResonance()
	score.Prosody = s.calculateProsody()
	score.Coherence = s.calculateCoherence()
	score.VoiceConsistency = s.calculateVoiceConsistency()
	score.SurpriseFactor = s.calculateSurpriseFactor()

	// Calculate total score (weighted average)
	score.Total = (score.Specificity*0.15 +
		score.Originality*0.15 +
		score.EmotionalResonance*0.20 +
		score.Prosody*0.15 +
		score.Coherence*0.15 +
		score.VoiceConsistency*0.10 +
		score.SurpriseFactor*0.10)

	// Generate feedback for each dimension
	s.generateFeedback(score)

	return score
}

// calculateSpecificity measures use of concrete vs abstract language
func (s *Song) calculateSpecificity() float64 {
	if len(s.Sections) == 0 {
		return 0.0
	}

	concreteWords := []string{
		"red", "blue", "green", "black", "white", "yellow", "purple", "orange",
		"wood", "metal", "glass", "stone", "water", "fire", "wind", "earth",
		"coffee", "rain", "snow", "sun", "moon", "star", "tree", "flower",
		"walk", "run", "jump", "sing", "dance", "cry", "laugh", "smile",
		"house", "car", "phone", "computer", "guitar", "piano", "drum",
	}

	abstractWords := []string{
		"love", "hate", "happy", "sad", "good", "bad", "beautiful", "ugly",
		"hope", "fear", "dream", "reality", "time", "space", "energy", "power",
		"freedom", "peace", "war", "life", "death", "soul", "spirit", "mind",
	}

	totalWords := 0
	concreteCount := 0
	abstractCount := 0

	for _, section := range s.Sections {
		for _, line := range section.Lines {
			words := strings.Fields(strings.ToLower(line.Text))
			totalWords += len(words)

			for _, word := range words {
				word = strings.Trim(word, ".,!?;\"'()")
				for _, concrete := range concreteWords {
					if word == concrete {
						concreteCount++
						break
					}
				}
				for _, abstract := range abstractWords {
					if word == abstract {
						abstractCount++
						break
					}
				}
			}
		}
	}

	if totalWords == 0 {
		return 0.0
	}

	concreteRatio := float64(concreteCount) / float64(totalWords)
	abstractRatio := float64(abstractCount) / float64(totalWords)

	// Higher specificity score for more concrete language
	return math.Min(100.0, (concreteRatio/(concreteRatio+abstractRatio+0.01))*100.0)
}

// calculateOriginality measures use of unique language patterns
func (s *Song) calculateOriginality() float64 {
	if len(s.Sections) == 0 {
		return 0.0
	}

	// Common clichés to avoid
	cliches := []string{
		"break my heart", "piece of my heart", "follow your heart",
		"from the bottom of my heart", "heart of gold", "heart skips a beat",
		"love at first sight", "match made in heaven", "meant to be",
		"better than ever", "once upon a time", "happily ever after",
		"live life to the fullest", "life is short", "time heals all wounds",
		"what goes around comes around", "when life gives you lemons",
	}

	totalLines := 0
	clicheCount := 0

	for _, section := range s.Sections {
		for _, line := range section.Lines {
			totalLines++
			lineLower := strings.ToLower(line.Text)

			for _, cliche := range cliches {
				if strings.Contains(lineLower, cliche) {
					clicheCount++
					break
				}
			}
		}
	}

	if totalLines == 0 {
		return 100.0
	}

	// Lower score for more clichés
	clicheRatio := float64(clicheCount) / float64(totalLines)
	return math.Max(0.0, (1.0-clicheRatio)*100.0)
}

// calculateEmotionalResonance measures emotional language and imagery
func (s *Song) calculateEmotionalResonance() float64 {
	if len(s.Sections) == 0 {
		return 0.0
	}

	emotionalWords := []string{
		"feel", "heart", "soul", "pain", "joy", "love", "hate", "hope", "fear",
		"dream", "cry", "laugh", "smile", "touch", "warm", "cold", "fire", "ice",
		"bitter", "sweet", "soft", "hard", "bright", "dark", "light", "shadow",
		"whisper", "scream", "silent", "loud", "gentle", "rough", "smooth", "sharp",
	}

	sensoryWords := []string{
		"see", "hear", "touch", "taste", "smell", "feel", "look", "sound",
		"red", "blue", "green", "yellow", "soft", "hard", "warm", "cold",
		"sweet", "sour", "bitter", "salty", "loud", "quiet", "bright", "dark",
		"smooth", "rough", "sharp", "dull", "heavy", "light", "hot", "cold",
	}

	totalWords := 0
	emotionalCount := 0
	sensoryCount := 0

	for _, section := range s.Sections {
		for _, line := range section.Lines {
			words := strings.Fields(strings.ToLower(line.Text))
			totalWords += len(words)

			for _, word := range words {
				word = strings.Trim(word, ".,!?;\"'()")

				for _, emotional := range emotionalWords {
					if word == emotional {
						emotionalCount++
						break
					}
				}

				for _, sensory := range sensoryWords {
					if word == sensory {
						sensoryCount++
						break
					}
				}
			}
		}
	}

	if totalWords == 0 {
		return 0.0
	}

	emotionalRatio := float64(emotionalCount) / float64(totalWords)
	sensoryRatio := float64(sensoryCount) / float64(totalWords)

	// Score based on emotional and sensory language density
	score := (emotionalRatio + sensoryRatio) * 100.0
	return math.Min(100.0, score)
}

// calculateProsody measures rhyme and rhythm consistency
func (s *Song) calculateProsody() float64 {
	if len(s.Sections) == 0 {
		return 0.0
	}

	// Simple syllable-based prosody analysis
	var totalSyllables int
	var linesWithSyllables int

	for _, section := range s.Sections {
		for _, line := range section.Lines {
			if line.Syllables > 0 {
				totalSyllables += line.Syllables
				linesWithSyllables++
			}
		}
	}

	if linesWithSyllables == 0 {
		return 0.0
	}

	avgSyllables := float64(totalSyllables) / float64(linesWithSyllables)

	// Score based on consistency (lower variance = higher score)
	// This is a simplified implementation
	variance := 0.0
	for _, section := range s.Sections {
		for _, line := range section.Lines {
			if line.Syllables > 0 {
				diff := float64(line.Syllables) - avgSyllables
				variance += diff * diff
			}
		}
	}

	variance /= float64(linesWithSyllables)
	consistency := math.Max(0.0, 1.0-(variance/avgSyllables))

	return consistency * 100.0
}

// calculateCoherence measures how well sections flow together
func (s *Song) calculateCoherence() float64 {
	if len(s.Sections) <= 1 {
		return 100.0 // Single section is perfectly coherent
	}

	// Simple coherence based on section ordering and transitions
	// This is a placeholder for more sophisticated analysis

	coherenceScore := 100.0

	// Penalize unusual section orders
	expectedOrder := map[SectionType][]SectionType{
		SectionIntro:     {SectionVerse, SectionChorus},
		SectionVerse:     {SectionChorus, SectionPreChorus, SectionBridge},
		SectionPreChorus: {SectionChorus},
		SectionChorus:    {SectionVerse, SectionBridge, SectionOutro},
		SectionBridge:    {SectionChorus, SectionOutro},
		SectionOutro:     {},
	}

	for i, section := range s.Sections {
		if i == 0 {
			continue // Skip first section
		}

		prevSection := s.Sections[i-1]
		validTransitions, exists := expectedOrder[prevSection.Type]

		if exists {
			valid := false
			for _, validType := range validTransitions {
				if section.Type == validType {
					valid = true
					break
				}
			}
			if !valid {
				coherenceScore -= 10.0 // Penalty for unexpected transitions
			}
		}
	}

	return math.Max(0.0, coherenceScore)
}

// calculateVoiceConsistency measures consistency of language and style
func (s *Song) calculateVoiceConsistency() float64 {
	if len(s.Sections) == 0 {
		return 0.0
	}

	// Simple consistency based on word usage patterns
	// This is a placeholder for more sophisticated analysis

	wordFreq := make(map[string]int)
	totalWords := 0

	for _, section := range s.Sections {
		for _, line := range section.Lines {
			words := strings.Fields(strings.ToLower(line.Text))
			for _, word := range words {
				word = strings.Trim(word, ".,!?;\"'()")
				if len(word) > 3 { // Only consider meaningful words
					wordFreq[word]++
					totalWords++
				}
			}
		}
	}

	if totalWords == 0 {
		return 0.0
	}

	// Calculate consistency based on word distribution
	// Higher score for more consistent word usage
	var consistency float64
	for _, count := range wordFreq {
		freq := float64(count) / float64(totalWords)
		consistency += freq * freq
	}

	return math.Min(100.0, consistency*1000.0)
}

// calculateSurpriseFactor measures unexpected elements and fresh imagery
func (s *Song) calculateSurpriseFactor() float64 {
	if len(s.Sections) == 0 {
		return 0.0
	}

	// Simple surprise factor based on unique word usage
	uniqueWords := make(map[string]bool)
	totalWords := 0

	for _, section := range s.Sections {
		for _, line := range section.Lines {
			words := strings.Fields(strings.ToLower(line.Text))
			for _, word := range words {
				word = strings.Trim(word, ".,!?;\"'()")
				if len(word) > 3 {
					uniqueWords[word] = true
					totalWords++
				}
			}
		}
	}

	if totalWords == 0 {
		return 0.0
	}

	// Higher surprise factor for more unique words
	uniquenessRatio := float64(len(uniqueWords)) / float64(totalWords)
	return math.Min(100.0, uniquenessRatio*200.0)
}

// generateFeedback creates actionable feedback for each dimension
func (s *Song) generateFeedback(score *QualityScore) {
	// Specificity feedback
	if score.Specificity < 60.0 {
		score.Feedback["specificity"] = "Try replacing abstract concepts with concrete, sensory details. Instead of 'I feel sad,' try 'Tears trace salty paths down my cheeks.'"
	} else {
		score.Feedback["specificity"] = "Good use of concrete imagery and specific details."
	}

	// Originality feedback
	if score.Originality < 60.0 {
		score.Feedback["originality"] = "Watch for clichés and overused phrases. Try fresh metaphors and unique perspectives."
	} else {
		score.Feedback["originality"] = "Original language and fresh imagery - well done!"
	}

	// Emotional resonance feedback
	if score.EmotionalResonance < 60.0 {
		score.Feedback["emotional_resonance"] = "Incorporate more emotional and sensory language to connect with listeners."
	} else {
		score.Feedback["emotional_resonance"] = "Strong emotional connection through vivid language."
	}

	// Prosody feedback
	if score.Prosody < 60.0 {
		score.Feedback["prosody"] = "Work on consistent syllable counts and rhyme patterns for better flow."
	} else {
		score.Feedback["prosody"] = "Good rhythmic consistency and natural flow."
	}

	// Coherence feedback
	if score.Coherence < 60.0 {
		score.Feedback["coherence"] = "Review section transitions - ensure each section leads naturally to the next."
	} else {
		score.Feedback["coherence"] = "Sections flow together naturally and logically."
	}

	// Voice consistency feedback
	if score.VoiceConsistency < 60.0 {
		score.Feedback["voice_consistency"] = "Maintain consistent language style and tone throughout the song."
	} else {
		score.Feedback["voice_consistency"] = "Consistent voice and style maintained throughout."
	}

	// Surprise factor feedback
	if score.SurpriseFactor < 60.0 {
		score.Feedback["surprise_factor"] = "Add some unexpected elements or fresh imagery to engage listeners."
	} else {
		score.Feedback["surprise_factor"] = "Good use of surprising and memorable elements."
	}
}
