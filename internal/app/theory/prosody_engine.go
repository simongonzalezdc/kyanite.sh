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
	// Stress lexicon for common words (value is stress pattern: 1=stressed, 0=unstressed)
	stressLexicon map[string][]Stress
}

// NewProsodyEngine creates a new prosody engine
func NewProsodyEngine() *ProsodyEngine {
	e := &ProsodyEngine{
		cache:         make(map[string]int),
		stressLexicon: make(map[string][]Stress),
	}
	e.initializeStressLexicon()
	return e
}

// initializeStressLexicon populates the stress lexicon with common English words
// and their stress patterns. This provides more accurate prosody analysis.
func (e *ProsodyEngine) initializeStressLexicon() {
	// Common monosyllabic words (stressed in isolation)
	monosyllabic := []string{
		"love", "heart", "soul", "dream", "night", "day", "sun", "moon",
		"life", "light", "dark", "time", "sky", "rain", "wind", "song",
		"world", "hand", "eye", "face", "mind", "voice", "way", "home",
	}
	for _, word := range monosyllabic {
		e.stressLexicon[word] = []Stress{Stressed}
	}

	// Unstressed function words (usually unstressed in context)
	functionWords := []string{
		"the", "a", "an", "and", "or", "but", "in", "on", "at", "to",
		"for", "of", "with", "by", "as", "is", "it", "be", "are", "was",
		"were", "been", "being", "have", "has", "had", "do", "does", "did",
		"will", "would", "could", "should", "may", "might", "must", "can",
		"that", "this", "these", "those", "my", "your", "his", "her", "its",
		"our", "their", "which", "who", "what", "when", "where", "why", "how",
	}
	for _, word := range functionWords {
		e.stressLexicon[word] = []Stress{Unstressed}
	}

	// Common two-syllable words with stress patterns
	twoSyllable := map[string][]Stress{
		// Stress on first syllable (trochaic)
		"never":    {Stressed, Unstressed},
		"ever":     {Stressed, Unstressed},
		"over":     {Stressed, Unstressed},
		"under":    {Stressed, Unstressed},
		"after":    {Stressed, Unstressed},
		"better":   {Stressed, Unstressed},
		"city":     {Stressed, Unstressed},
		"morning":  {Stressed, Unstressed},
		"evening":  {Stressed, Unstressed},
		"summer":   {Stressed, Unstressed},
		"winter":   {Stressed, Unstressed},
		"water":    {Stressed, Unstressed},
		"river":    {Stressed, Unstressed},
		"music":    {Stressed, Unstressed},
		"moment":   {Stressed, Unstressed},
		"silent":   {Stressed, Unstressed},
		"broken":   {Stressed, Unstressed},
		"fallen":   {Stressed, Unstressed},
		"golden":   {Stressed, Unstressed},
		"heaven":   {Stressed, Unstressed},
		"dancing":  {Stressed, Unstressed},
		"singing":  {Stressed, Unstressed},
		"running":  {Stressed, Unstressed},
		"waiting":  {Stressed, Unstressed},
		"holding":  {Stressed, Unstressed},
		"letting":  {Stressed, Unstressed},
		"feeling":  {Stressed, Unstressed},
		"dreaming": {Stressed, Unstressed},

		// Stress on second syllable (iambic)
		"again":   {Unstressed, Stressed},
		"away":    {Unstressed, Stressed},
		"along":   {Unstressed, Stressed},
		"alone":   {Unstressed, Stressed},
		"around":  {Unstressed, Stressed},
		"before":  {Unstressed, Stressed},
		"behind":  {Unstressed, Stressed},
		"below":   {Unstressed, Stressed},
		"beyond":  {Unstressed, Stressed},
		"between": {Unstressed, Stressed},
		"tonight": {Unstressed, Stressed},
		"today":   {Unstressed, Stressed},
		"above":   {Unstressed, Stressed},
		"inside":  {Unstressed, Stressed},
		"outside": {Unstressed, Stressed},
		"believe": {Unstressed, Stressed},
		"forget":  {Unstressed, Stressed},
		"begin":   {Unstressed, Stressed},
		"become":  {Unstressed, Stressed},
	}
	for word, stress := range twoSyllable {
		e.stressLexicon[word] = stress
	}

	// Common three-syllable words
	threeSyllable := map[string][]Stress{
		"beautiful":  {Stressed, Unstressed, Unstressed},
		"wonderful":  {Stressed, Unstressed, Unstressed},
		"everything": {Stressed, Unstressed, Unstressed},
		"anything":   {Stressed, Unstressed, Unstressed},
		"everyone":   {Stressed, Unstressed, Unstressed},
		"anyone":     {Stressed, Unstressed, Unstressed},
		"memory":     {Stressed, Unstressed, Unstressed},
		"melody":     {Stressed, Unstressed, Unstressed},
		"yesterday":  {Stressed, Unstressed, Unstressed},
		"tomorrow":   {Unstressed, Stressed, Unstressed},
		"forever":    {Unstressed, Stressed, Unstressed},
		"together":   {Unstressed, Stressed, Unstressed},
		"remember":   {Unstressed, Stressed, Unstressed},
		"surrender":  {Unstressed, Stressed, Unstressed},
		"impossible": {Unstressed, Stressed, Unstressed},
		"emotion":    {Unstressed, Stressed, Unstressed},
	}
	for word, stress := range threeSyllable {
		e.stressLexicon[word] = stress
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
		wordSyllables := e.analyzeWordStress(word)
		syllables = append(syllables, wordSyllables...)
		totalCount += len(wordSyllables)
	}

	return LineAnalysis{
		Text:      text,
		Syllables: syllables,
		Count:     totalCount,
		Meter:     e.detectMeter(syllables),
	}
}

// analyzeWordStress determines the stress pattern for a word.
// Uses lexicon lookup first, then falls back to rules-based detection.
func (e *ProsodyEngine) analyzeWordStress(word string) []Syllable {
	cleanWord := strings.ToLower(regexp.MustCompile(`[^a-z]`).ReplaceAllString(word, ""))
	if len(cleanWord) == 0 {
		return nil
	}

	// Split word into syllables first
	syllableParts := e.splitIntoSyllables(cleanWord)

	// Try lexicon lookup
	if stressPattern, ok := e.stressLexicon[cleanWord]; ok {
		return e.applyStressPattern(syllableParts, stressPattern)
	}

	// Apply rules-based stress detection
	stressPattern := e.inferStressPattern(cleanWord, len(syllableParts))
	return e.applyStressPattern(syllableParts, stressPattern)
}

// splitIntoSyllables attempts to split a word into its syllables
func (e *ProsodyEngine) splitIntoSyllables(word string) []string {
	if len(word) <= 3 {
		return []string{word}
	}

	// Find vowel clusters as syllable nuclei
	vowelPattern := regexp.MustCompile(`[aeiouy]+`)
	vowelMatches := vowelPattern.FindAllStringIndex(word, -1)

	if len(vowelMatches) == 0 {
		return []string{word}
	}

	if len(vowelMatches) == 1 {
		return []string{word}
	}

	// Build syllables by splitting between vowel clusters
	var syllables []string
	lastEnd := 0

	for i, match := range vowelMatches {
		if i == 0 {
			continue
		}

		// Find split point between previous vowel end and current vowel start
		prevEnd := vowelMatches[i-1][1]
		currStart := match[0]

		// Split in the middle of consonant cluster
		splitPoint := (prevEnd + currStart) / 2
		if splitPoint < prevEnd {
			splitPoint = prevEnd
		}

		syllables = append(syllables, word[lastEnd:splitPoint])
		lastEnd = splitPoint
	}

	// Add final syllable
	syllables = append(syllables, word[lastEnd:])

	// Handle silent 'e' at end
	if len(syllables) > 1 && syllables[len(syllables)-1] == "e" {
		// Merge with previous syllable
		syllables[len(syllables)-2] += syllables[len(syllables)-1]
		syllables = syllables[:len(syllables)-1]
	}

	return syllables
}

// inferStressPattern uses rules to determine stress pattern for unknown words
func (e *ProsodyEngine) inferStressPattern(word string, syllableCount int) []Stress {
	if syllableCount <= 0 {
		return nil
	}

	if syllableCount == 1 {
		return []Stress{Stressed}
	}

	pattern := make([]Stress, syllableCount)

	// Rule 1: Two-syllable nouns/adjectives usually stress first syllable
	// Rule 2: Words with common prefixes (un-, re-, be-, a-) stress second syllable
	// Rule 3: Words with certain suffixes have predictable stress

	// Check for common unstressed prefixes
	unstressedPrefixes := []string{"un", "re", "be", "de", "pre", "pro", "a", "dis", "mis", "ex"}
	hasUnstressedPrefix := false
	for _, prefix := range unstressedPrefixes {
		if strings.HasPrefix(word, prefix) && len(word) > len(prefix)+2 {
			hasUnstressedPrefix = true
			break
		}
	}

	// Check for suffixes that attract stress
	stressAttractingSuffixes := []string{"tion", "sion", "ic", "ical", "ity", "ious", "eous"}
	hasStressSuffix := false
	for _, suffix := range stressAttractingSuffixes {
		if strings.HasSuffix(word, suffix) {
			hasStressSuffix = true
			break
		}
	}

	// Check for suffixes that repel stress (stress on preceding syllable)
	stressRepellingSuffixes := []string{"ing", "ed", "er", "est", "ly", "ness", "ment", "ful", "less"}
	hasRepellingSuffix := false
	for _, suffix := range stressRepellingSuffixes {
		if strings.HasSuffix(word, suffix) {
			hasRepellingSuffix = true
			break
		}
	}

	// Apply rules
	if syllableCount == 2 {
		if hasUnstressedPrefix {
			pattern[0] = Unstressed
			pattern[1] = Stressed
		} else {
			// Default: stress first syllable for two-syllable words
			pattern[0] = Stressed
			pattern[1] = Unstressed
		}
	} else {
		// Three or more syllables
		if hasStressSuffix {
			// Stress penultimate syllable
			for i := range pattern {
				if i == syllableCount-2 {
					pattern[i] = Stressed
				} else {
					pattern[i] = Unstressed
				}
			}
		} else if hasRepellingSuffix {
			// Stress antepenultimate (third from last) or first available
			stressPos := syllableCount - 3
			if stressPos < 0 {
				stressPos = 0
			}
			for i := range pattern {
				if i == stressPos {
					pattern[i] = Stressed
				} else {
					pattern[i] = Unstressed
				}
			}
		} else if hasUnstressedPrefix {
			// Stress second syllable
			pattern[0] = Unstressed
			pattern[1] = Stressed
			for i := 2; i < len(pattern); i++ {
				pattern[i] = Unstressed
			}
		} else {
			// Default: stress first syllable
			pattern[0] = Stressed
			for i := 1; i < len(pattern); i++ {
				// Add secondary stress for longer words
				if i == syllableCount-2 && syllableCount > 3 {
					pattern[i] = Secondary
				} else {
					pattern[i] = Unstressed
				}
			}
		}
	}

	return pattern
}

// applyStressPattern creates syllables with the given stress pattern
func (e *ProsodyEngine) applyStressPattern(syllableParts []string, stressPattern []Stress) []Syllable {
	result := make([]Syllable, len(syllableParts))

	for i, part := range syllableParts {
		stress := Unstressed
		if i < len(stressPattern) {
			stress = stressPattern[i]
		}
		result[i] = Syllable{
			Text:   part,
			Stress: stress,
		}
	}

	return result
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
