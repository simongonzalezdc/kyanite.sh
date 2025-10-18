package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDictionary(t *testing.T) {
	dict := NewDictionary()
	if dict == nil {
		t.Fatal("NewDictionary() returned nil")
	}
	
	if dict.IsLoaded() {
		t.Error("NewDictionary() should not be loaded initially")
	}
}

func TestDictionaryLoad(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			},
			"beautiful": {
				"syllables": 3,
				"pronunciation": "B Y UW1 T AH0 F AH0 L",
				"rhymes": ["dutiful", "grateful"],
				"pos": ["adjective"]
			}
		}
	}`
	
	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	dict := NewDictionary()
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	if !dict.IsLoaded() {
		t.Error("Dictionary should be loaded after successful LoadDictionary call")
	}
}

func TestDictionaryLoadInvalidFile(t *testing.T) {
	dict := NewDictionary()
	err := dict.LoadDictionary("nonexistent_file.json")
	if err == nil {
		t.Error("Expected error when loading nonexistent file")
	}
}

func TestGetWordEntry(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			}
		}
	}`
	
	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	dict := NewDictionary()
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	// Test existing word
	entry, exists := dict.GetWordEntry("test")
	if !exists {
		t.Error("Expected word 'test' to exist")
	}
	if entry.Syllables != 1 {
		t.Errorf("Expected 1 syllable, got %d", entry.Syllables)
	}
	
	// Test case insensitive lookup
	entry, exists = dict.GetWordEntry("TEST")
	if !exists {
		t.Error("Expected word 'TEST' to exist (case insensitive)")
	}
	
	// Test non-existent word
	_, exists = dict.GetWordEntry("nonexistent")
	if exists {
		t.Error("Expected word 'nonexistent' to not exist")
	}
}

func TestCountSyllables(t *testing.T) {
	dict := NewDictionary()
	
	// Test with unloaded dictionary (should use fallback)
	syllables, err := dict.CountSyllables("beautiful")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if syllables <= 0 {
		t.Errorf("Expected positive syllable count for 'beautiful', got %d", syllables)
	}
	
	// Test with loaded dictionary
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			},
			"beautiful": {
				"syllables": 3,
				"pronunciation": "B Y UW1 T AH0 F AH0 L",
				"rhymes": ["dutiful", "grateful"],
				"pos": ["adjective"]
			}
		}
	}`
	
	err = os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	// Test word from dictionary
	syllables, err = dict.CountSyllables("beautiful")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if syllables != 3 {
		t.Errorf("Expected 3 syllables for 'beautiful', got %d", syllables)
	}
	
	// Test word not in dictionary (should use fallback)
	syllables, err = dict.CountSyllables("computer")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if syllables <= 0 {
		t.Errorf("Expected positive syllable count for 'computer', got %d", syllables)
	}
}

func TestCountSyllablesHeuristic(t *testing.T) {
	dict := NewDictionary()
	
	testCases := []struct {
		word      string
		expected  int
		tolerance int
	}{
		{"hello", 2, 0},
		{"world", 1, 0},
		{"computer", 3, 1},
		{"beautiful", 3, 1},
		{"queue", 1, 0},
		{"one", 1, 0},
		{"two", 1, 0},
		{"three", 1, 0},
		{"", 0, 0},
		{"a", 1, 0},
		{"the", 1, 0},
	}
	
	for _, tc := range testCases {
		result := dict.countSyllablesHeuristic(tc.word)
		if tc.tolerance == 0 && result != tc.expected {
			t.Errorf("Expected %d syllables for '%s', got %d", tc.expected, tc.word, result)
		} else if tc.tolerance > 0 {
			diff := result - tc.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tc.tolerance {
				t.Errorf("Expected %d±%d syllables for '%s', got %d", tc.expected, tc.tolerance, tc.word, result)
			}
		}
	}
}

func TestFindRhymes(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			},
			"love": {
				"syllables": 1,
				"pronunciation": "L AH1 V",
				"rhymes": ["dove", "glove", "above"],
				"pos": ["verb", "noun"]
			}
		}
	}`
	
	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	dict := NewDictionary()
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	// Test word with rhymes in dictionary
	rhymes, err := dict.FindRhymes("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedRhymes := []string{"best", "rest", "west"}
	if len(rhymes) != len(expectedRhymes) {
		t.Errorf("Expected %d rhymes, got %d", len(expectedRhymes), len(rhymes))
	}
	
	// Test word not in dictionary (should use fallback)
	rhymes, err = dict.FindRhymes("computer")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Should return some results or empty array, but not error
	if rhymes == nil {
		t.Error("Expected rhymes array, got nil")
	}
}

func TestFindRhymesPhonetic(t *testing.T) {
	dict := NewDictionary()
	
	// Test with empty string
	rhymes := dict.findRhymesPhonetic("")
	if len(rhymes) != 0 {
		t.Errorf("Expected empty rhymes for empty string, got %v", rhymes)
	}
	
	// Test with short word
	rhymes = dict.findRhymesPhonetic("a")
	if len(rhymes) != 0 {
		t.Errorf("Expected empty rhymes for short word, got %v", rhymes)
	}
}

func TestCountSyllablesInText(t *testing.T) {
	dict := NewDictionary()
	
	testCases := []struct {
		text     string
		expected int
	}{
		{"hello world", 3},
		{"the quick brown fox", 6},
		{"", 0},
		{"   ", 0},
		{"a", 1},
		{"hello, world!", 3}, // Should handle punctuation
	}
	
	for _, tc := range testCases {
		result, err := dict.CountSyllablesInText(tc.text)
		if err != nil {
			t.Errorf("Unexpected error for text '%s': %v", tc.text, err)
		}
		if result != tc.expected {
			t.Errorf("Expected %d syllables for '%s', got %d", tc.expected, tc.text, result)
		}
	}
}

func TestGetStats(t *testing.T) {
	dict := NewDictionary()
	
	// Test with unloaded dictionary
	stats := dict.GetStats()
	if stats.TotalWords != 0 {
		t.Errorf("Expected 0 total words for unloaded dictionary, got %d", stats.TotalWords)
	}
	
	// Test with loaded dictionary
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			},
			"beautiful": {
				"syllables": 3,
				"pronunciation": "B Y UW1 T AH0 F AH0 L",
				"rhymes": ["dutiful", "grateful"],
				"pos": ["adjective"]
			}
		}
	}`
	
	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	stats = dict.GetStats()
	if stats.TotalWords != 2 {
		t.Errorf("Expected 2 total words, got %d", stats.TotalWords)
	}
	if stats.TotalRhymes != 5 { // 3 for test + 2 for beautiful
		t.Errorf("Expected 5 total rhymes, got %d", stats.TotalRhymes)
	}
	if stats.AvgSyllables != 2.0 { // (1 + 3) / 2
		t.Errorf("Expected avg syllables 2.0, got %f", stats.AvgSyllables)
	}
}

func TestAddWord(t *testing.T) {
	dict := NewDictionary()
	
	entry := WordEntry{
		Syllables:    2,
		Pronunciation: "T EH1 S T",
		Rhymes:       []string{"best", "rest"},
		POS:          []string{"noun"},
	}
	
	err := dict.AddWord("testword", entry)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	// Verify word was added
	retrieved, exists := dict.GetWordEntry("testword")
	if !exists {
		t.Error("Word was not added")
	}
	if retrieved.Syllables != 2 {
		t.Errorf("Expected 2 syllables, got %d", retrieved.Syllables)
	}
}

func TestSearchWords(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			},
			"testing": {
				"syllables": 2,
				"pronunciation": "T EH1 S T IH0 NG",
				"rhymes": ["resting", "besting"],
				"pos": ["verb"]
			},
			"best": {
				"syllables": 1,
				"pronunciation": "B EH1 S T",
				"rhymes": ["test", "rest"],
				"pos": ["adjective", "noun"]
			}
		}
	}`
	
	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	dict := NewDictionary()
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	// Test exact match
	words, err := dict.SearchWords("test", 10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(words) != 1 || words[0] != "test" {
		t.Errorf("Expected ['test'], got %v", words)
	}
	
	// Test wildcard
	words, err = dict.SearchWords("test*", 10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(words) != 2 {
		t.Errorf("Expected 2 words matching 'test*', got %d", len(words))
	}
	
	// Test limit
	words, err = dict.SearchWords("*", 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(words) != 1 {
		t.Errorf("Expected 1 word with limit, got %d", len(words))
	}
}

func TestGetWordsBySyllableCount(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			},
			"beautiful": {
				"syllables": 3,
				"pronunciation": "B Y UW1 T AH0 F AH0 L",
				"rhymes": ["dutiful", "grateful"],
				"pos": ["adjective"]
			},
			"best": {
				"syllables": 1,
				"pronunciation": "B EH1 S T",
				"rhymes": ["test", "rest"],
				"pos": ["adjective", "noun"]
			}
		}
	}`
	
	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	dict := NewDictionary()
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	// Test 1 syllable
	words, err := dict.GetWordsBySyllableCount(1, 10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(words) != 2 {
		t.Errorf("Expected 2 words with 1 syllable, got %d", len(words))
	}
	
	// Test 3 syllables
	words, err = dict.GetWordsBySyllableCount(3, 10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(words) != 1 || words[0] != "beautiful" {
		t.Errorf("Expected ['beautiful'], got %v", words)
	}
	
	// Test no matches
	words, err = dict.GetWordsBySyllableCount(5, 10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(words) != 0 {
		t.Errorf("Expected no words with 5 syllables, got %v", words)
	}
}

func TestGetWordsByPartOfSpeech(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			},
			"beautiful": {
				"syllables": 3,
				"pronunciation": "B Y UW1 T AH0 F AH0 L",
				"rhymes": ["dutiful", "grateful"],
				"pos": ["adjective"]
			},
			"best": {
				"syllables": 1,
				"pronunciation": "B EH1 S T",
				"rhymes": ["test", "rest"],
				"pos": ["adjective", "noun"]
			}
		}
	}`
	
	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	dict := NewDictionary()
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	// Test nouns
	words, err := dict.GetWordsByPartOfSpeech("noun", 10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(words) != 2 {
		t.Errorf("Expected 2 nouns, got %d", len(words))
	}
	
	// Test adjectives
	words, err = dict.GetWordsByPartOfSpeech("adjective", 10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(words) != 2 {
		t.Errorf("Expected 2 adjectives, got %d", len(words))
	}
	
	// Test no matches
	words, err = dict.GetWordsByPartOfSpeech("adverb", 10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(words) != 0 {
		t.Errorf("Expected no adverbs, got %v", words)
	}
}

func TestValidateWord(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			}
		}
	}`
	
	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	dict := NewDictionary()
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	// Test existing word
	if !dict.ValidateWord("test") {
		t.Error("Expected 'test' to be valid")
	}
	
	// Test non-existing word
	if dict.ValidateWord("nonexistent") {
		t.Error("Expected 'nonexistent' to be invalid")
	}
	
	// Test case insensitive
	if !dict.ValidateWord("TEST") {
		t.Error("Expected 'TEST' to be valid (case insensitive)")
	}
	
	// Test empty string
	if dict.ValidateWord("") {
		t.Error("Expected empty string to be invalid")
	}
}

func TestGetRandomWord(t *testing.T) {
	dict := NewDictionary()
	
	// Test with empty dictionary
	_, _, err := dict.GetRandomWord()
	if err == nil {
		t.Error("Expected error for empty dictionary")
	}
	
	// Test with loaded dictionary
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			}
		}
	}`
	
	err = os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	word, entry, err := dict.GetRandomWord()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if word != "test" {
		t.Errorf("Expected 'test', got '%s'", word)
	}
	if entry.Syllables != 1 {
		t.Errorf("Expected 1 syllable, got %d", entry.Syllables)
	}
}

func TestGetPhoneticEnding(t *testing.T) {
	dict := NewDictionary()
	
	testCases := []struct {
		word     string
		expected string
	}{
		{"test", "est"},
		{"hello", "llo"},
		{"a", "a"},
		{"ab", "ab"},
		{"", ""},
	}
	
	for _, tc := range testCases {
		result := dict.getPhoneticEnding(tc.word)
		if result != tc.expected {
			t.Errorf("Expected '%s' for '%s', got '%s'", tc.expected, tc.word, result)
		}
	}
}

func TestDictionaryConcurrency(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")
	
	testDict := `{
		"words": {
			"test": {
				"syllables": 1,
				"pronunciation": "T EH1 S T",
				"rhymes": ["best", "rest", "west"],
				"pos": ["noun", "verb"]
			}
		}
	}`
	
	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}
	
	dict := NewDictionary()
	err = dict.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}
	
	// Test concurrent reads
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, exists := dict.GetWordEntry("test")
			if !exists {
				t.Error("Expected word to exist")
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}