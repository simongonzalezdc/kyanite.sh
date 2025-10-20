package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
	"time"

	errutil "github.com/Kyanite/noise/internal/errutil"
)

// WordEntry represents a word entry in the dictionary
type WordEntry struct {
	Syllables     int      `json:"syllables"`
	Pronunciation string   `json:"pronunciation"`
	Rhymes        []string `json:"rhymes"`
	POS           []string `json:"pos"` // Parts of speech
}

// Dictionary represents the enhanced dictionary service
type Dictionary struct {
	words    map[string]WordEntry
	rhymeMap map[string][]string // Reverse mapping for faster rhyme lookups
	mutex    sync.RWMutex
	loaded   bool
	loadTime time.Time
}

// DictionaryStats provides statistics about the dictionary
type DictionaryStats struct {
	TotalWords   int     `json:"total_words"`
	TotalRhymes  int     `json:"total_rhymes"`
	AvgSyllables float64 `json:"avg_syllables"`
	LoadTime     string  `json:"load_time"`
	LastUpdated  string  `json:"last_updated"`
	CacheHits    int64   `json:"cache_hits"`
	CacheMisses  int64   `json:"cache_misses"`
}

// NewDictionary creates a new dictionary service
func NewDictionary() *Dictionary {
	return &Dictionary{
		words:    make(map[string]WordEntry),
		rhymeMap: make(map[string][]string),
	}
}

// LoadDictionary loads the dictionary from the JSON file
func (d *Dictionary) LoadDictionary(filePath string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Read the dictionary file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errutil.Wrap(err, "read dictionary file")
	}

	// Parse the JSON
	var dictData struct {
		Words map[string]WordEntry `json:"words"`
	}

	if err := json.Unmarshal(data, &dictData); err != nil {
		return errutil.Wrap(err, "parse dictionary JSON")
	}

	// Load the words
	d.words = dictData.Words

	// Build reverse rhyme mapping for faster lookups
	d.buildRhymeMap()

	d.loaded = true
	d.loadTime = time.Now()

	return nil
}

// buildRhymeMap creates a reverse mapping for faster rhyme lookups
func (d *Dictionary) buildRhymeMap() {
	d.rhymeMap = make(map[string][]string)

	for word, entry := range d.words {
		// Ensure entry.Rhymes is never nil
		rhymes := entry.Rhymes
		if rhymes == nil {
			rhymes = []string{}
		}

		for _, rhyme := range rhymes {
			// Normalize the rhyme word
			normalizedRhyme := strings.ToLower(strings.TrimSpace(rhyme))
			if normalizedRhyme != "" {
				if d.rhymeMap[normalizedRhyme] == nil {
					d.rhymeMap[normalizedRhyme] = []string{}
				}
				d.rhymeMap[normalizedRhyme] = append(d.rhymeMap[normalizedRhyme], word)
			}
		}
	}
}

// IsLoaded returns whether the dictionary has been loaded
func (d *Dictionary) IsLoaded() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.loaded
}

// GetWordEntry returns the dictionary entry for a word
func (d *Dictionary) GetWordEntry(word string) (WordEntry, bool) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Normalize the word
	normalizedWord := strings.ToLower(strings.TrimSpace(word))
	entry, exists := d.words[normalizedWord]
	return entry, exists
}

// CountSyllables counts syllables in a word with fallback algorithms
func (d *Dictionary) CountSyllables(word string) (int, error) {
	// First try to get from dictionary
	if entry, exists := d.GetWordEntry(word); exists {
		return entry.Syllables, nil
	}

	// Fallback to heuristic-based counting
	return d.countSyllablesHeuristic(word), nil
}

// countSyllablesHeuristic provides a fallback syllable counting method
func (d *Dictionary) countSyllablesHeuristic(word string) int {
	if word == "" {
		return 0
	}

	// Normalize the word
	word = strings.ToLower(strings.TrimSpace(word))

	// Handle common exceptions
	exceptions := map[string]int{
		"queue":    1,
		"one":      1,
		"two":      1,
		"three":    1,
		"four":     1,
		"five":     1,
		"seven":    2,
		"eight":    1,
		"nine":     1,
		"ten":      1,
		"hundred":  2,
		"thousand": 2,
		"million":  3,
		"billion":  3,
	}

	if syllables, exists := exceptions[word]; exists {
		return syllables
	}

	// Remove non-alphabetic characters
	re := regexp.MustCompile(`[^a-z]`)
	word = re.ReplaceAllString(word, "")

	// Handle empty string after cleaning
	if word == "" {
		return 1 // Single character words like "a" or "I"
	}

	vowels := "aeiouy"
	syllables := 0
	prevWasVowel := false

	for i, char := range word {
		isVowel := strings.ContainsRune(vowels, char)

		if isVowel && !prevWasVowel {
			syllables++
		}
		prevWasVowel = isVowel

		// Handle special cases for better accuracy
		if i > 0 {
			prevChar := word[i-1]
			currChar := char

			// Handle silent 'e' - only if it's at the end and follows a consonant
			if i == len(word)-1 && currChar == 'e' && !isVowel {
				// Check if previous character is a consonant followed by a single vowel
				if i >= 2 {
					prevPrevChar := word[i-2]
					if strings.ContainsRune(vowels, rune(prevPrevChar)) &&
						!strings.ContainsRune(vowels, rune(prevChar)) &&
						syllables > 1 {
						syllables--
					}
				}
			}

			// Handle 'le' at the end - this typically creates an extra syllable
			if i == len(word)-1 && currChar == 'e' && prevChar == 'l' && syllables > 1 {
				// Only add syllable if the 'l' is preceded by a consonant
				if i >= 2 {
					prevPrevChar := word[i-2]
					if !strings.ContainsRune(vowels, rune(prevPrevChar)) {
						syllables++
					}
				}
			}

			// Handle 'y' as vowel when it functions as a vowel
			if currChar == 'y' && !isVowel && prevWasVowel {
				// 'y' after vowel is usually a consonant, but can create diphthongs
				// Don't add extra syllable here as it's handled by main vowel counting
			}
		}
	}

	// Every word has at least one syllable
	if syllables == 0 {
		syllables = 1
	}

	return syllables
}

// FindRhymes finds rhyming words for a given word
func (d *Dictionary) FindRhymes(word string) ([]string, error) {
	// First try to get from dictionary
	if entry, exists := d.GetWordEntry(word); exists {
		// Ensure we never return nil, return empty slice instead
		if entry.Rhymes == nil {
			return []string{}, nil
		}
		return entry.Rhymes, nil
	}

	// Fallback to phonetic-based rhyme finding
	rhymes := d.findRhymesPhonetic(word)
	// Ensure we never return nil, return empty slice instead
	if rhymes == nil {
		return []string{}, nil
	}
	return rhymes, nil
}

// findRhymesPhonetic provides a fallback rhyme finding method
func (d *Dictionary) findRhymesPhonetic(word string) []string {
	if word == "" {
		return []string{}
	}

	// Normalize the word
	word = strings.ToLower(strings.TrimSpace(word))

	// Check reverse rhyme map
	if rhymes, exists := d.rhymeMap[word]; exists {
		return rhymes
	}

	// Simple phonetic matching based on word endings
	var matches []string
	wordEnd := d.getPhoneticEnding(word)

	if wordEnd == "" {
		return []string{}
	}

	// Find words with similar phonetic endings
	for dictWord := range d.words {
		if dictWord == word {
			continue
		}

		dictEnd := d.getPhoneticEnding(dictWord)
		if dictEnd == wordEnd && len(dictWord) > 2 {
			matches = append(matches, dictWord)
		}
	}

	return matches
}

// getPhoneticEnding extracts the phonetic ending of a word for rhyme matching
func (d *Dictionary) getPhoneticEnding(word string) string {
	if len(word) < 2 {
		return word
	}

	// Get the last 2-3 characters for basic matching
	if len(word) >= 3 {
		return word[len(word)-3:]
	}
	return word
}

// CountSyllablesInText counts syllables in a line of text
func (d *Dictionary) CountSyllablesInText(text string) (int, error) {
	if text == "" {
		return 0, nil
	}

	// Split into words
	words := strings.Fields(text)
	totalSyllables := 0

	for _, word := range words {
		// Clean the word of punctuation
		re := regexp.MustCompile(`[^\w']`)
		cleanWord := re.ReplaceAllString(word, "")

		if cleanWord != "" {
			syllables, err := d.CountSyllables(cleanWord)
			if err != nil {
				continue // Skip words that can't be processed
			}
			totalSyllables += syllables
		}
	}

	return totalSyllables, nil
}

// GetStats returns statistics about the dictionary
func (d *Dictionary) GetStats() DictionaryStats {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	stats := DictionaryStats{
		TotalWords:  len(d.words),
		LoadTime:    d.loadTime.Format(time.RFC3339),
		LastUpdated: d.loadTime.Format("2006-01-02"),
	}

	// Calculate total rhymes and average syllables
	totalRhymes := 0
	totalSyllables := 0

	for _, entry := range d.words {
		totalRhymes += len(entry.Rhymes)
		totalSyllables += entry.Syllables
	}

	stats.TotalRhymes = totalRhymes

	if stats.TotalWords > 0 {
		stats.AvgSyllables = float64(totalSyllables) / float64(stats.TotalWords)
	}

	return stats
}

// AddWord adds a new word to the dictionary
func (d *Dictionary) AddWord(word string, entry WordEntry) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Normalize the word
	normalizedWord := strings.ToLower(strings.TrimSpace(word))

	// Ensure entry.Rhymes is never nil
	if entry.Rhymes == nil {
		entry.Rhymes = []string{}
	}

	d.words[normalizedWord] = entry

	// Update rhyme map
	for _, rhyme := range entry.Rhymes {
		normalizedRhyme := strings.ToLower(strings.TrimSpace(rhyme))
		if normalizedRhyme != "" {
			if d.rhymeMap[normalizedRhyme] == nil {
				d.rhymeMap[normalizedRhyme] = []string{}
			}
			d.rhymeMap[normalizedRhyme] = append(d.rhymeMap[normalizedRhyme], normalizedWord)
		}
	}

	return nil
}

// SearchWords searches for words matching a pattern
func (d *Dictionary) SearchWords(pattern string, limit int) ([]string, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var matches []string

	// Handle exact match vs pattern matching
	var regexPattern string
	if strings.Contains(pattern, "*") {
		// Use wildcard matching if pattern contains asterisks
		regexPattern = strings.ReplaceAll(pattern, "*", ".*")
	} else {
		// Use exact matching for plain strings (anchored to match whole word)
		regexPattern = "^" + regexp.QuoteMeta(pattern) + "$"
	}

	regex, err := regexp.Compile("(?i)" + regexPattern)
	if err != nil {
		return nil, errutil.Wrap(err, "invalid search pattern")
	}

	for word := range d.words {
		if regex.MatchString(word) {
			matches = append(matches, word)
			if limit > 0 && len(matches) >= limit {
				break
			}
		}
	}

	return matches, nil
}

// GetRandomWord returns a random word from the dictionary
func (d *Dictionary) GetRandomWord() (string, WordEntry, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if len(d.words) == 0 {
		return "", WordEntry{}, fmt.Errorf("dictionary is empty")
	}

	// Simple random selection (not cryptographically secure)
	count := 0
	for word, entry := range d.words {
		if count == len(d.words)/2 { // Simple way to get a "random" word
			return word, entry, nil
		}
		count++
	}

	// Fallback to first word
	for word, entry := range d.words {
		return word, entry, nil
	}

	return "", WordEntry{}, fmt.Errorf("no words found in dictionary")
}

// ValidateWord checks if a word exists in the dictionary
func (d *Dictionary) ValidateWord(word string) bool {
	_, exists := d.GetWordEntry(word)
	return exists
}

// GetWordsBySyllableCount returns words with a specific syllable count
func (d *Dictionary) GetWordsBySyllableCount(syllables int, limit int) ([]string, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var matches []string

	for word, entry := range d.words {
		if entry.Syllables == syllables {
			matches = append(matches, word)
			if limit > 0 && len(matches) >= limit {
				break
			}
		}
	}

	return matches, nil
}

// GetWordsByPartOfSpeech returns words by part of speech
func (d *Dictionary) GetWordsByPartOfSpeech(pos string, limit int) ([]string, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var matches []string
	targetPOS := strings.ToLower(pos)

	for word, entry := range d.words {
		for _, wordPOS := range entry.POS {
			if strings.ToLower(wordPOS) == targetPOS {
				matches = append(matches, word)
				if limit > 0 && len(matches) >= limit {
					break
				}
				break
			}
		}
	}

	return matches, nil
}
