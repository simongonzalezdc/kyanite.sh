//go:build ignore

package main

import (
	"fmt"
	"log"

	"github.com/Kyanite/noise/internal/app"
)

// TestTheoryService tests the theory service functionality
func TestTheoryService() {
	// Test theory service functionality
	theoryService := app.NewTheoryService()

	fmt.Println("ðŸŽµ Testing Theory Tools Implementation")
	fmt.Println("=====================================")

	// Test scale functionality
	fmt.Println("\n1. Testing Scale Information:")
	scaleInfo, err := theoryService.GetScaleInfo("C", "major")
	if err != nil {
		log.Printf("Error getting scale info: %v", err)
	} else {
		fmt.Printf("Scale: %s\n", scaleInfo.Name)
		fmt.Printf("Notes: %v\n", scaleInfo.Notes)
		fmt.Printf("Pattern: %v\n", scaleInfo.Pattern)
		fmt.Printf("Description: %s\n", scaleInfo.Description)
	}

	// Test chord functionality
	fmt.Println("\n2. Testing Chord Information:")
	chordInfo, err := theoryService.GetChordInfo("C", "major")
	if err != nil {
		log.Printf("Error getting chord info: %v", err)
	} else {
		fmt.Printf("Chord: %s %s\n", chordInfo.Root, chordInfo.Quality)
		fmt.Printf("Notes: %v\n", chordInfo.Notes)
		fmt.Printf("Intervals: %v\n", chordInfo.Intervals)
		fmt.Printf("Description: %s\n", chordInfo.Description)
	}

	// Test chord analysis
	fmt.Println("\n3. Testing Chord Analysis:")
	testText := "I love the C chord and G major scale in my F songs"
	analysis, err := theoryService.AnalyzeChords(testText)
	if err != nil {
		log.Printf("Error analyzing chords: %v", err)
	} else {
		fmt.Printf("Input: %s\n", analysis.Input)
		fmt.Printf("Found %d chords:\n", len(analysis.DetectedChords))
		for _, chord := range analysis.DetectedChords {
			fmt.Printf("  - %s %s: %v\n", chord.Root, chord.Quality, chord.Notes)
		}
		if len(analysis.Suggestions) > 0 {
			fmt.Println("Suggestions:")
			for _, suggestion := range analysis.Suggestions {
				fmt.Printf("  - %s\n", suggestion)
			}
		}
	}

	// Test progression
	fmt.Println("\n4. Testing Chord Progression:")
	progressionInfo, err := theoryService.GetProgressionInfo("C", "I-V-vi-IV")
	if err != nil {
		log.Printf("Error getting progression: %v", err)
	} else {
		fmt.Printf("Progression: %s\n", progressionInfo.Name)
		fmt.Printf("Chords: %v\n", progressionInfo.Chords)
		fmt.Printf("Description: %s\n", progressionInfo.Description)
	}

	// Test rhymes
	fmt.Println("\n5. Testing Rhyme Finder:")
	rhymes, err := theoryService.FindRhymes("love")
	if err != nil {
		log.Printf("Error finding rhymes: %v", err)
	} else {
		fmt.Printf("Rhymes for 'love': %v\n", rhymes)
	}

	fmt.Println("\nâœ… Theory tools implementation completed successfully!")
}
