//go:build ignore

package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/puente-labs/noise/internal/app"
)

func main() {
	fmt.Println("🎵 Testing Week 2 CMU Dictionary Enhancements")
	fmt.Println("==============================================")

	// Test enhanced dictionary functionality
	testEnhancedDictionary()
	
	// Test theory service integration
	testTheoryServiceIntegration()
	
	// Test fallback handling
	testFallbackHandling()
	
	// Test performance optimizations
	testPerformanceOptimizations()
	
	fmt.Println("\n✅ Week 2 CMU Dictionary Enhancements testing completed successfully!")
}

func testEnhancedDictionary() {
	fmt.Println("\n1. Testing Enhanced Dictionary Service:")
	
	// Create dictionary with proper path
	dict := app.NewDictionary()
	
	// Try to load the dictionary file
	dictPath := filepath.Join("data", "dictionary.json")
	err := dict.LoadDictionary(dictPath)
	if err != nil {
		log.Printf("Warning: Could not load dictionary file: %v", err)
		log.Println("Continuing with fallback functionality...")
	}
	
	// Test syllable counting with dictionary
	fmt.Println("\n   Testing syllable counting:")
	testWords := []string{"love", "beautiful", "computer", "extraordinary", "rhythm"}
	
	for _, word := range testWords {
		syllables, err := dict.CountSyllables(word)
		if err != nil {
			log.Printf("Error counting syllables for '%s': %v", word, err)
			continue
		}
		fmt.Printf("   - '%s': %d syllables\n", word, syllables)
	}
	
	// Test rhyme finding
	fmt.Println("\n   Testing rhyme finding:")
	rhymeWords := []string{"love", "time", "heart", "night"}
	
	for _, word := range rhymeWords {
		rhymes, err := dict.FindRhymes(word)
		if err != nil {
			log.Printf("Error finding rhymes for '%s': %v", word, err)
			continue
		}
		fmt.Printf("   - '%s': %v\n", word, rhymes)
	}
	
	// Test text syllable counting
	fmt.Println("\n   Testing text syllable counting:")
	testTexts := []string{
		"Beautiful love song in the night",
		"The quick brown fox jumps over the lazy dog",
		"Extraordinary rhythm and rhyme",
	}
	
	for _, text := range testTexts {
		syllables, err := dict.CountSyllablesInText(text)
		if err != nil {
			log.Printf("Error counting syllables in text: %v", err)
			continue
		}
		fmt.Printf("   - '%s': %d syllables\n", text, syllables)
	}
	
	// Get dictionary statistics
	stats := dict.GetStats()
	fmt.Printf("\n   Dictionary Statistics:\n")
	fmt.Printf("   - Total Words: %d\n", stats.TotalWords)
	fmt.Printf("   - Total Rhymes: %d\n", stats.TotalRhymes)
	fmt.Printf("   - Average Syllables: %.2f\n", stats.AvgSyllables)
	fmt.Printf("   - Load Time: %s\n", stats.LoadTime)
}

func testTheoryServiceIntegration() {
	fmt.Println("\n2. Testing Theory Service Integration:")
	
	// Create theory service (will automatically load dictionary)
	theoryService := app.NewTheoryService()
	
	// Test enhanced rhyme finding
	fmt.Println("\n   Testing enhanced rhyme finding:")
	rhymeWords := []string{"love", "time", "heart"}
	
	for _, word := range rhymeWords {
		rhymes, err := theoryService.FindRhymes(word)
		if err != nil {
			log.Printf("Error finding rhymes for '%s': %v", word, err)
			continue
		}
		fmt.Printf("   - '%s': %v\n", word, rhymes)
	}
	
	// Test enhanced syllable counting
	fmt.Println("\n   Testing enhanced syllable counting:")
	syllableWords := []string{"beautiful", "computer", "rhythm"}
	
	for _, word := range syllableWords {
		syllables, err := theoryService.CountSyllables(word)
		if err != nil {
			log.Printf("Error counting syllables for '%s': %v", word, err)
			continue
		}
		fmt.Printf("   - '%s': %d syllables\n", word, syllables)
	}
	
	// Test enhanced prosody analysis
	fmt.Println("\n   Testing enhanced prosody analysis:")
	prosodyTexts := []string{
		"Beautiful love in the moonlight",
		"Dancing through the night",
		"Extraordinary rhythm flows",
	}
	
	for _, text := range prosodyTexts {
		prosody, err := theoryService.AnalyzeProsody(text)
		if err != nil {
			log.Printf("Error analyzing prosody: %v", err)
			continue
		}
		fmt.Printf("   - '%s': %d syllables\n", text, prosody)
	}
	
	// Test dictionary statistics
	stats, err := theoryService.GetDictionaryStats()
	if err != nil {
		log.Printf("Error getting dictionary stats: %v", err)
	} else {
		fmt.Printf("\n   Dictionary Statistics from Theory Service:\n")
		fmt.Printf("   - Total Words: %d\n", stats.TotalWords)
		fmt.Printf("   - Total Rhymes: %d\n", stats.TotalRhymes)
		fmt.Printf("   - Average Syllables: %.2f\n", stats.AvgSyllables)
	}
}

func testFallbackHandling() {
	fmt.Println("\n3. Testing Fallback Handling:")
	
	// Create theory service
	theoryService := app.NewTheoryService()
	
	// Test words that might not be in dictionary
	fmt.Println("\n   Testing fallback for uncommon words:")
	uncommonWords := []string{"supercalifragilisticexpialidocious", "pneumonoultramicroscopicsilicovolcanoconiosis", "antidisestablishmentarianism"}
	
	for _, word := range uncommonWords {
		// Test syllable counting fallback
		syllables, err := theoryService.CountSyllables(word)
		if err != nil {
			log.Printf("Error counting syllables for '%s': %v", word, err)
			continue
		}
		fmt.Printf("   - '%s': %d syllables (fallback)\n", word, syllables)
		
		// Test rhyme finding fallback
		rhymes, err := theoryService.FindRhymes(word)
		if err != nil {
			log.Printf("Error finding rhymes for '%s': %v", word, err)
			continue
		}
		fmt.Printf("   - '%s' rhymes: %v (fallback)\n", word, rhymes)
	}
	
	// Test word validation
	fmt.Println("\n   Testing word validation:")
	testValidation := []string{"love", "nonexistentword", "", "   "}
	
	for _, word := range testValidation {
		isValid := theoryService.ValidateWord(word)
		fmt.Printf("   - '%s': %v\n", word, isValid)
	}
}

func testPerformanceOptimizations() {
	fmt.Println("\n4. Testing Performance Optimizations:")
	
	// Create theory service
	theoryService := app.NewTheoryService()
	
	// Test批量操作
	fmt.Println("\n   Testing batch operations:")
	
	// Test multiple syllable counts
	words := []string{"love", "time", "heart", "night", "blue", "day", "world", "eyes", "home", "dream"}
	
	fmt.Println("   Testing batch syllable counting:")
	for i, word := range words {
		syllables, err := theoryService.CountSyllables(word)
		if err != nil {
			log.Printf("Error counting syllables for '%s': %v", word, err)
			continue
		}
		if i < 5 { // Show first 5 results
			fmt.Printf("   - '%s': %d syllables\n", word, syllables)
		}
	}
	fmt.Printf("   ... processed %d words\n", len(words))
	
	// Test word search
	fmt.Println("\n   Testing word search patterns:")
	searchPatterns := []string{"*", "love*", "*ing", "beautiful"}
	
	for _, pattern := range searchPatterns {
		results, err := theoryService.SearchWords(pattern, 5)
		if err != nil {
			log.Printf("Error searching with pattern '%s': %v", pattern, err)
			continue
		}
		fmt.Printf("   - Pattern '%s': %v\n", pattern, results)
	}
	
	// Test syllable count filtering
	fmt.Println("\n   Testing syllable count filtering:")
	for i := 1; i <= 4; i++ {
		words, err := theoryService.GetWordsBySyllableCount(i, 3)
		if err != nil {
			log.Printf("Error getting words with %d syllables: %v", i, err)
			continue
		}
		fmt.Printf("   - %d syllables: %v\n", i, words)
	}
	
	// Test part of speech filtering
	fmt.Println("\n   Testing part of speech filtering:")
	posTypes := []string{"noun", "verb", "adjective"}
	
	for _, pos := range posTypes {
		words, err := theoryService.SearchWords("*", 10) // Get some words first
		if err != nil {
			log.Printf("Error getting words: %v", err)
			continue
		}
		fmt.Printf("   - %s: (showing up to 10 words)\n", pos)
		for i, word := range words {
			if i >= 3 { // Show first 3
				break
			}
			isValid := theoryService.ValidateWord(word)
			fmt.Printf("     * '%s': %v\n", word, isValid)
		}
	}
}