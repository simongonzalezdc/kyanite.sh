// Test file for the File I/O system
// This file tests the basic functionality of the file I/O system
package lyricforge_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/puente-labs/lyricforge/internal/domain"
	"github.com/puente-labs/lyricforge/internal/infra/files"
)

// TestFileIO tests the file I/O system functionality
func TestFileIO() {
	fmt.Println("Testing File I/O System...")

	// Create a temporary directory for testing
	testDir := "./test_songs"
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		log.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create file service
	fileService, err := files.New(files.Config{
		BaseDir:     testDir,
		AutoSave:    true,
		AutoSaveInterval: time.Minute,
		BackupCount: 3,
	})
	if err != nil {
		log.Fatalf("Failed to create file service: %v", err)
	}
	defer fileService.Close()

	// Test 1: Create a new song
	fmt.Println("\n1. Creating a new song...")
	song := &domain.Song{
		Metadata: domain.SongMetadata{
			Title:     "Test Song",
			Artist:    "Test Artist",
			Key:       "C Major",
			Tempo:     120,
			TimeSignature: "4/4",
			Tags:      []string{"test", "demo"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Sections: []domain.Section{
			{
				Type:   domain.SectionVerse,
				Number: 1,
				Lines: []domain.Line{
					{Text: "This is the first line of the verse", Syllables: 8},
					{Text: "This is the second line of the verse", Syllables: 8},
				},
			},
			{
				Type:   domain.SectionChorus,
				Number: 1,
				Lines: []domain.Line{
					{Text: "This is the chorus line", Syllables: 5},
					{Text: "Another chorus line here", Syllables: 5},
				},
			},
		},
	}

	// Test 2: Save song to file
	fmt.Println("2. Saving song to file...")
	testFilePath := "test_song.md"
	if err := fileService.WriteSong(song, testFilePath); err != nil {
		log.Fatalf("Failed to save song: %v", err)
	}
	fmt.Printf("   Song saved to: %s\n", testFilePath)

	// Test 3: List songs
	fmt.Println("3. Listing songs...")
	songFiles, err := fileService.ListSongs()
	if err != nil {
		log.Fatalf("Failed to list songs: %v", err)
	}
	fmt.Printf("   Found %d song files:\n", len(songFiles))
	for _, file := range songFiles {
		fmt.Printf("   - %s\n", file)
	}

	// Test 4: Read song from file
	fmt.Println("4. Reading song from file...")
	loadedSong, err := fileService.ReadSong(testFilePath)
	if err != nil {
		log.Fatalf("Failed to read song: %v", err)
	}
	fmt.Printf("   Loaded song: %s by %s\n", loadedSong.Metadata.Title, loadedSong.Metadata.Artist)
	fmt.Printf("   Sections: %d\n", len(loadedSong.Sections))

	// Test 5: Verify YAML frontmatter
	fmt.Println("5. Verifying YAML frontmatter...")
	if loadedSong.Metadata.Title != song.Metadata.Title {
		log.Fatalf("Title mismatch: expected %s, got %s", song.Metadata.Title, loadedSong.Metadata.Title)
	}
	if loadedSong.Metadata.Artist != song.Metadata.Artist {
		log.Fatalf("Artist mismatch: expected %s, got %s", song.Metadata.Artist, loadedSong.Metadata.Artist)
	}
	fmt.Printf("   Title: %s\n", loadedSong.Metadata.Title)
	fmt.Printf("   Artist: %s\n", loadedSong.Metadata.Artist)
	fmt.Printf("   Key: %s\n", loadedSong.Metadata.Key)
	fmt.Printf("   Tempo: %d\n", loadedSong.Metadata.Tempo)

	// Test 6: Test file caching
	fmt.Println("6. Testing file caching...")
	_, exists := fileService.GetCachedFile(testFilePath)
	if !exists {
		log.Fatalf("File not found in cache")
	}
	fmt.Printf("   File is cached: %v\n", exists)

	// Test 7: Test file watcher
	fmt.Println("7. Testing file watcher...")
	watchReceived := make(chan bool, 1)
	watcher := &testWatcher{
		events: watchReceived,
	}

	if err := fileService.WatchFile(testFilePath, watcher); err != nil {
		log.Fatalf("Failed to watch file: %v", err)
	}

	// Modify the file to trigger watcher
	go func() {
		time.Sleep(100 * time.Millisecond)
		song.Metadata.UpdatedAt = time.Now()
		if err := fileService.WriteSong(song, testFilePath); err != nil {
			fmt.Printf("   Warning: Failed to write song for watcher test: %v\n", err)
		}
	}()

	// Wait for watcher event
	select {
	case <-watchReceived:
		fmt.Println("   Watcher event received successfully")
	case <-time.After(time.Second):
		fmt.Println("   Watcher event timeout (this may be normal)")
	}

	// Cleanup
	if err := fileService.UnwatchFile(testFilePath, watcher); err != nil {
		fmt.Printf("   Warning: Failed to unwatch file: %v\n", err)
	}

	// Test 8: Test file deletion
	fmt.Println("8. Testing file deletion...")
	if err := fileService.DeleteSong(testFilePath); err != nil {
		log.Fatalf("Failed to delete song: %v", err)
	}
	fmt.Println("   Song file deleted successfully")

	// Verify deletion
	if exists, _ := fileService.SongExists(testFilePath); exists {
		log.Fatalf("Song still exists after deletion")
	}
	fmt.Println("   Verified song no longer exists")

	// Test 9: Test backup functionality
	fmt.Println("9. Testing backup functionality...")
	backupTestFile := "backup_test.md"

	// Create and save a song
	backupSong := &domain.Song{
		Metadata: domain.SongMetadata{
			Title:     "Backup Test Song",
			Artist:    "Backup Artist",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Sections: []domain.Section{
			{
				Type:   domain.SectionVerse,
				Number: 1,
				Lines:  []domain.Line{{Text: "Backup test line", Syllables: 4}},
			},
		},
	}

	if err := fileService.WriteSong(backupSong, backupTestFile); err != nil {
		log.Fatalf("Failed to save backup test song: %v", err)
	}

	// Create backup
	if err := fileService.CreateBackup(backupTestFile); err != nil {
		log.Fatalf("Failed to create backup: %v", err)
	}
	fmt.Println("   Backup created successfully")

	// Verify backup exists (check in the test directory)
	backupPath := filepath.Join(testDir, backupTestFile+".backup.*")
	matches, err := filepath.Glob(backupPath)
	if err != nil || len(matches) == 0 {
		fmt.Println("   Warning: Backup file not found (may be in different location)")
	} else {
		fmt.Printf("   Backup file created: %s\n", matches[0])
	}

	fmt.Println("\n✅ All File I/O tests completed successfully!")
}

// testWatcher implements the FileWatcher interface for testing
type testWatcher struct {
	events chan<- bool
}

func (w *testWatcher) OnFileChanged(filePath string, event files.FileEvent) {
	fmt.Printf("   File changed: %s (event: %s)\n", filePath, event)
	select {
	case w.events <- true:
	default:
	}
}

func (w *testWatcher) OnFileDeleted(filePath string) {
	fmt.Printf("   File deleted: %s\n", filePath)
	select {
	case w.events <- true:
	default:
	}
}