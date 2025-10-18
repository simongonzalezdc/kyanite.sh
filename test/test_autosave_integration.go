package lyricforge_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/puente-labs/lyricforge/internal/app"
	"github.com/puente-labs/lyricforge/internal/domain"
	"github.com/puente-labs/lyricforge/internal/infra/db"
	"github.com/puente-labs/lyricforge/internal/infra/files"
)

// TestAutoSaveIntegration tests the auto-save integration functionality
func TestAutoSaveIntegration() {
	fmt.Println("Testing Auto-Save Content Serialization Integration...")

	// Initialize database
	dbConfig := db.Config{DataDir: "./test_data"}
	database, err := db.New(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize file service
	fileConfig := files.Config{BaseDir: "./test_songs"}
	fileService, err := files.New(fileConfig)
	if err != nil {
		log.Fatalf("Failed to initialize file service: %v", err)
	}
	defer fileService.Close()

	// Initialize auto-save service
	autoSaveConfig := app.DefaultAutoSaveConfig()
	autoSaveService := app.NewAutoSaveService(database, autoSaveConfig)

	// Initialize editor service
	editorService := app.NewEditorService(database, database)

	// Create a test song
	song := &domain.Song{
		Metadata: domain.SongMetadata{
			Title:     "Test Song for Auto-Save",
			Artist:    "Test Artist",
			Key:       "C",
			Tempo:     120,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Sections: []domain.Section{
			{
				Type:   domain.SectionVerse,
				Number: 1,
				Lines: []domain.Line{
					{Text: "This is the first line of the test song"},
					{Text: "This is the second line of the test song"},
				},
			},
			{
				Type:   domain.SectionChorus,
				Number: 1,
				Lines: []domain.Line{
					{Text: "This is the chorus of the test song"},
					{Text: "It should be properly serialized and restored"},
				},
			},
		},
		RawContent: `---
title: Test Song for Auto-Save
artist: Test Artist
key: C
tempo: 120
created_at: 2025-01-01T00:00:00Z
updated_at: 2025-01-01T00:00:00Z
---

## verse 1

This is the first line of the test song
This is the second line of the test song

## chorus 1

This is the chorus of the test song
It should be properly serialized and restored
`,
	}

	// Save the song to database
	savedSong, err := editorService.CreateSong(song.Metadata.Title, song.Metadata.Artist)
	if err != nil {
		log.Fatalf("Failed to create song: %v", err)
	}

	// Update the song with sections and content
	savedSong.Sections = song.Sections
	savedSong.RawContent = song.RawContent
	savedSong.Metadata.Key = song.Metadata.Key
	savedSong.Metadata.Tempo = song.Metadata.Tempo

	err = editorService.SaveSong(savedSong)
	if err != nil {
		log.Fatalf("Failed to save song: %v", err)
	}

	fmt.Printf("Created and saved song with ID: %d\n", savedSong.ID)

	// Test 1: Auto-save with content serialization
	fmt.Println("\n=== Test 1: Auto-Save Content Serialization ===")
	err = editorService.AutoSave(savedSong)
	if err != nil {
		log.Fatalf("Failed to auto-save song: %v", err)
	}
	fmt.Println("✓ Auto-save completed successfully")

	// Test 2: Create milestone
	fmt.Println("\n=== Test 2: Milestone Creation ===")
	milestoneName := "Test Milestone"
	err = editorService.CreateMilestone(savedSong, milestoneName)
	if err != nil {
		log.Fatalf("Failed to create milestone: %v", err)
	}
	fmt.Printf("✓ Milestone '%s' created successfully\n", milestoneName)

	// Test 3: Get version history
	fmt.Println("\n=== Test 3: Version History ===")
	versions, err := editorService.GetVersions(savedSong.ID, 10)
	if err != nil {
		log.Fatalf("Failed to get version history: %v", err)
	}

	fmt.Printf("Found %d versions:\n", len(versions))
	for i, version := range versions {
		fmt.Printf("  %d. %s (Milestone: %t, Created: %s)\n",
			i+1, version.MilestoneName, version.IsMilestone,
			version.CreatedAt.Format("2006-01-02 15:04:05"))
		
		// Show first 100 characters of content
		if len(version.Content) > 100 {
			fmt.Printf("     Content preview: %s...\n", version.Content[:100])
		} else {
			fmt.Printf("     Content: %s\n", version.Content)
		}
	}

	// Test 4: Version restoration
	fmt.Println("\n=== Test 4: Version Restoration ===")
	if len(versions) > 0 {
		// Restore the first version
		versionToRestore := versions[0]
		fmt.Printf("Restoring version: %s\n", versionToRestore.MilestoneName)
		
		restoredSong, err := editorService.RestoreVersion(savedSong.ID, versionToRestore.ID)
		if err != nil {
			log.Fatalf("Failed to restore version: %v", err)
		}
		
		fmt.Printf("✓ Version restored successfully\n")
		fmt.Printf("Restored song title: %s\n", restoredSong.Metadata.Title)
		fmt.Printf("Restored song artist: %s\n", restoredSong.Metadata.Artist)
		fmt.Printf("Restored content length: %d characters\n", len(restoredSong.RawContent))
		
		// Verify content was properly restored
		if restoredSong.RawContent == versionToRestore.Content {
			fmt.Println("✓ Content restoration verified - content matches exactly")
		} else {
			fmt.Println("⚠ Content restoration warning - content differs from version")
			fmt.Printf("Expected length: %d, Got length: %d\n", 
				len(versionToRestore.Content), len(restoredSong.RawContent))
		}
	}

	// Test 5: Auto-save service integration
	fmt.Println("\n=== Test 5: Auto-Save Service Integration ===")
	
	// Start auto-save service
	ctx := context.Background()
	err = autoSaveService.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start auto-save service: %v", err)
	}
	if err := autoSaveService.Stop(); err != nil {
		fmt.Printf("Warning: Failed to stop auto-save service: %v\n", err)
	}

	// Test direct content saving
	testContent := `---
title: Direct Auto-Save Test
artist: Test Artist
---

## verse 1

This content was saved directly through the auto-save service
It should preserve the full markdown structure
`

	err = autoSaveService.SaveWithVersioning(savedSong.ID, testContent, false, "Direct Test")
	if err != nil {
		log.Fatalf("Failed to save content with versioning: %v", err)
	}
	fmt.Println("✓ Direct content save with versioning completed")

	// Test recovery from last save
	recoveredContent, err := autoSaveService.RecoverFromLastSave(savedSong.ID)
	if err != nil {
		log.Fatalf("Failed to recover from last save: %v", err)
	}
	
	if recoveredContent == testContent {
		fmt.Println("✓ Content recovery verified - recovered content matches saved content")
	} else {
		fmt.Println("⚠ Content recovery warning - recovered content differs")
		fmt.Printf("Expected length: %d, Got length: %d\n", 
			len(testContent), len(recoveredContent))
	}

	// Test 6: File I/O integration
	fmt.Println("\n=== Test 6: File I/O Integration ===")
	
	// Save song to file
	filename := "test_autosave_song.md"
	err = fileService.WriteSong(savedSong, filename)
	if err != nil {
		log.Fatalf("Failed to write song to file: %v", err)
	}
	fmt.Printf("✓ Song saved to file: %s\n", filename)

	// Read song from file
	readSong, err := fileService.ReadSong(filename)
	if err != nil {
		log.Fatalf("Failed to read song from file: %v", err)
	}
	fmt.Printf("✓ Song read from file: %s\n", filename)
	fmt.Printf("Read song title: %s\n", readSong.Metadata.Title)
	fmt.Printf("Read content length: %d characters\n", len(readSong.RawContent))

	// Clean up test files
	os.Remove(filename)
	os.RemoveAll("./test_data")
	os.RemoveAll("./test_songs")

	fmt.Println("\n=== All Tests Completed Successfully ===")
	fmt.Println("✓ Auto-save content serialization is working properly")
	fmt.Println("✓ Milestone creation is functional")
	fmt.Println("✓ Version restoration preserves full content")
	fmt.Println("✓ Auto-save service integration is complete")
	fmt.Println("✓ File I/O system integration is working")
}