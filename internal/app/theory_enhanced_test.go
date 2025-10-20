package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnhancedTheoryService(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")

	testDict := `{
		"words": {
			"love": {
				"syllables": 1,
				"pronunciation": "L AH1 V",
				"rhymes": ["dove", "glove", "above", "shove", "of", "rough", "tough"],
				"pos": ["verb", "noun"]
			},
			"beautiful": {
				"syllables": 3,
				"pronunciation": "B Y UW1 T AH0 F AH0 L",
				"rhymes": ["dutiful", "grateful", "hateful", "painful"],
				"pos": ["adjective"]
			},
			"computer": {
				"syllables": 3,
				"pronunciation": "K AH0 M P Y UW1 T ER0",
				"rhymes": ["commuter", "disputer"],
				"pos": ["noun"]
			}
		}
	}`

	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}

	// Create theory service with custom dictionary path
	theoryService := &TheoryService{
		dictionary: NewDictionary(),
	}

	err = theoryService.dictionary.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}

	// Test enhanced rhyme finding
	rhymes, err := theoryService.FindRhymes("love")
	if err != nil {
		t.Errorf("Unexpected error finding rhymes: %v", err)
	}
	expectedRhymes := []string{"dove", "glove", "above", "shove", "of", "rough", "tough"}
	if len(rhymes) != len(expectedRhymes) {
		t.Errorf("Expected %d rhymes, got %d", len(expectedRhymes), len(rhymes))
	}

	// Test enhanced syllable counting
	syllables, err := theoryService.CountSyllables("beautiful")
	if err != nil {
		t.Errorf("Unexpected error counting syllables: %v", err)
	}
	if syllables != 3 {
		t.Errorf("Expected 3 syllables for 'beautiful', got %d", syllables)
	}

	// Test fallback for word not in dictionary
	syllables, err = theoryService.CountSyllables("extraordinary")
	if err != nil {
		t.Errorf("Unexpected error counting syllables for fallback: %v", err)
	}
	if syllables <= 0 {
		t.Errorf("Expected positive syllable count for fallback word, got %d", syllables)
	}

	// Test enhanced prosody analysis
	prosody, err := theoryService.AnalyzeProsody("beautiful love")
	if err != nil {
		t.Errorf("Unexpected error analyzing prosody: %v", err)
	}
	expectedProsody := 4 // 3 + 1
	if prosody != expectedProsody {
		t.Errorf("Expected %d syllables for prosody, got %d", expectedProsody, prosody)
	}
}

func TestTheoryServiceWithUnloadedDictionary(t *testing.T) {
	// Create theory service without loading dictionary
	theoryService := NewTheoryService()

	// Should still work with fallbacks
	_, err := theoryService.FindRhymes("love")
	if err != nil {
		t.Errorf("Unexpected error finding rhymes with fallback: %v", err)
	}
	// Should return some rhymes from static dictionary

	syllables, err := theoryService.CountSyllables("beautiful")
	if err != nil {
		t.Errorf("Unexpected error counting syllables with fallback: %v", err)
	}
	if syllables <= 0 {
		t.Errorf("Expected positive syllable count with fallback, got %d", syllables)
	}
}

func TestGetDictionaryStats(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")

	testDict := `{
		"words": {
			"love": {
				"syllables": 1,
				"pronunciation": "L AH1 V",
				"rhymes": ["dove", "glove", "above"],
				"pos": ["verb", "noun"]
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

	theoryService := &TheoryService{
		dictionary: NewDictionary(),
	}

	err = theoryService.dictionary.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}

	stats, err := theoryService.GetDictionaryStats()
	if err != nil {
		t.Errorf("Unexpected error getting dictionary stats: %v", err)
	}

	if stats.TotalWords != 2 {
		t.Errorf("Expected 2 total words, got %d", stats.TotalWords)
	}

	if stats.TotalRhymes != 5 { // 3 for love + 2 for beautiful
		t.Errorf("Expected 5 total rhymes, got %d", stats.TotalRhymes)
	}

	if stats.AvgSyllables != 2.0 { // (1 + 3) / 2
		t.Errorf("Expected avg syllables 2.0, got %f", stats.AvgSyllables)
	}
}

func TestTheoryServiceValidateWord(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")

	testDict := `{
		"words": {
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

	theoryService := &TheoryService{
		dictionary: NewDictionary(),
	}

	err = theoryService.dictionary.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}

	// Test existing word
	if !theoryService.ValidateWord("love") {
		t.Error("Expected 'love' to be valid")
	}

	// Test non-existing word
	if theoryService.ValidateWord("nonexistent") {
		t.Error("Expected 'nonexistent' to be invalid")
	}

	// Test case insensitive
	if !theoryService.ValidateWord("LOVE") {
		t.Error("Expected 'LOVE' to be valid (case insensitive)")
	}
}

func TestTheoryServiceSearchWords(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")

	testDict := `{
		"words": {
			"love": {
				"syllables": 1,
				"pronunciation": "L AH1 V",
				"rhymes": ["dove", "glove", "above"],
				"pos": ["verb", "noun"]
			},
			"lovely": {
				"syllables": 2,
				"pronunciation": "L AH1 V L IY0",
				"rhymes": ["above", "dove"],
				"pos": ["adjective"]
			}
		}
	}`

	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}

	theoryService := &TheoryService{
		dictionary: NewDictionary(),
	}

	err = theoryService.dictionary.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}

	// Test wildcard search
	words, err := theoryService.SearchWords("lov*", 10)
	if err != nil {
		t.Errorf("Unexpected error searching words: %v", err)
	}
	if len(words) != 2 {
		t.Errorf("Expected 2 words matching 'lov*', got %d", len(words))
	}
}

func TestTheoryServiceGetWordsBySyllableCount(t *testing.T) {
	// Create a temporary dictionary file for testing
	tempDir := t.TempDir()
	dictFile := filepath.Join(tempDir, "test_dict.json")

	testDict := `{
		"words": {
			"love": {
				"syllables": 1,
				"pronunciation": "L AH1 V",
				"rhymes": ["dove", "glove", "above"],
				"pos": ["verb", "noun"]
			},
			"beautiful": {
				"syllables": 3,
				"pronunciation": "B Y UW1 T AH0 F AH0 L",
				"rhymes": ["dutiful", "grateful"],
				"pos": ["adjective"]
			},
			"lovely": {
				"syllables": 2,
				"pronunciation": "L AH1 V L IY0",
				"rhymes": ["above", "dove"],
				"pos": ["adjective"]
			}
		}
	}`

	err := os.WriteFile(dictFile, []byte(testDict), 0644)
	if err != nil {
		t.Fatalf("Failed to write test dictionary file: %v", err)
	}

	theoryService := &TheoryService{
		dictionary: NewDictionary(),
	}

	err = theoryService.dictionary.LoadDictionary(dictFile)
	if err != nil {
		t.Fatalf("Failed to load dictionary: %v", err)
	}

	// Test 1 syllable
	words, err := theoryService.GetWordsBySyllableCount(1, 10)
	if err != nil {
		t.Errorf("Unexpected error getting words by syllable count: %v", err)
	}
	if len(words) != 1 || words[0] != "love" {
		t.Errorf("Expected ['love'], got %v", words)
	}

	// Test 2 syllables
	words, err = theoryService.GetWordsBySyllableCount(2, 10)
	if err != nil {
		t.Errorf("Unexpected error getting words by syllable count: %v", err)
	}
	if len(words) != 1 || words[0] != "lovely" {
		t.Errorf("Expected ['lovely'], got %v", words)
	}

	// Test 3 syllables
	words, err = theoryService.GetWordsBySyllableCount(3, 10)
	if err != nil {
		t.Errorf("Unexpected error getting words by syllable count: %v", err)
	}
	if len(words) != 1 || words[0] != "beautiful" {
		t.Errorf("Expected ['beautiful'], got %v", words)
	}
}

func TestTheoryServiceIntegration(t *testing.T) {
	// Test the full integration with the actual dictionary file
	theoryService := NewTheoryService()

	// Test that the service works even if dictionary fails to load
	_, err := theoryService.FindRhymes("love")
	if err != nil {
		t.Errorf("Unexpected error finding rhymes: %v", err)
	}
	// Should return some rhymes from static fallback

	syllables, err := theoryService.CountSyllables("beautiful")
	if err != nil {
		t.Errorf("Unexpected error counting syllables: %v", err)
	}
	if syllables <= 0 {
		t.Errorf("Expected positive syllable count, got %d", syllables)
	}

	prosody, err := theoryService.AnalyzeProsody("beautiful love song")
	if err != nil {
		t.Errorf("Unexpected error analyzing prosody: %v", err)
	}
	if prosody <= 0 {
		t.Errorf("Expected positive prosody count, got %d", prosody)
	}
}
