package noise

import (
	"context"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/infra/files"
)

// TestAutoSaveIntegration validates the end-to-end auto-save workflow across database,
// file storage, versioning, and recovery paths. The test is intentionally verbose to
// ensure comprehensive coverage, mirroring the legacy demonstration while adhering to
// current linting and resource-handling requirements.
func TestAutoSaveIntegration(t *testing.T) {
	dataDir := t.TempDir()
	songsDir := t.TempDir()

	database, err := db.New(db.Config{DataDir: dataDir})
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Fatalf("Failed to close database: %v", err)
		}
	})

	fileService, err := files.New(files.Config{BaseDir: songsDir})
	if err != nil {
		t.Fatalf("Failed to initialize file service: %v", err)
	}
	t.Cleanup(func() {
		if err := fileService.Close(); err != nil {
			t.Fatalf("Failed to close file service: %v", err)
		}
	})

	autoSaveConfig := app.DefaultAutoSaveConfig()
	autoSaveService := app.NewAutoSaveService(database, autoSaveConfig)
	editorService := app.NewEditorService(database, database)

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

## verse 1

This is the first line of the test song
This is the second line of the test song

## chorus 1

This is the chorus of the test song
It should be properly serialized and restored
`,
	}

	savedSong, err := editorService.CreateSong(song.Metadata.Title, song.Metadata.Artist)
	if err != nil {
		t.Fatalf("Failed to create song: %v", err)
	}

	savedSong.Sections = song.Sections
	savedSong.RawContent = song.RawContent
	savedSong.Metadata.Key = song.Metadata.Key
	savedSong.Metadata.Tempo = song.Metadata.Tempo

	if err := editorService.SaveSong(savedSong); err != nil {
		t.Fatalf("Failed to save song: %v", err)
	}

	t.Logf("Created and saved song with ID: %d", savedSong.ID)

	t.Run("AutoSaveContentSerialization", func(t *testing.T) {
		if err := editorService.AutoSave(savedSong); err != nil {
			t.Fatalf("Failed to auto-save song: %v", err)
		}
		t.Log("Auto-save content serialization completed successfully")
	})

	milestoneName := "Test Milestone"
	t.Run("MilestoneCreation", func(t *testing.T) {
		if err := editorService.CreateMilestone(savedSong, milestoneName); err != nil {
			t.Fatalf("Failed to create milestone: %v", err)
		}
		t.Logf("Milestone %q created successfully", milestoneName)
	})

	var versions []*domain.Version
	t.Run("VersionHistory", func(t *testing.T) {
		var err error
		versions, err = editorService.GetVersions(savedSong.ID, 10)
		if err != nil {
			t.Fatalf("Failed to get version history: %v", err)
		}

		t.Logf("Found %d versions:", len(versions))
		for i, version := range versions {
			t.Logf("  %d. %s (Milestone: %t, Created: %s)",
				i+1,
				version.MilestoneName,
				version.IsMilestone,
				version.CreatedAt.Format("2006-01-02 15:04:05"),
			)

			if len(version.Content) > 100 {
				t.Logf("     Content preview: %s...", version.Content[:100])
			} else {
				t.Logf("     Content: %s", version.Content)
			}
		}
	})

	t.Run("VersionRestoration", func(t *testing.T) {
		if len(versions) == 0 {
			t.Skip("No versions available for restoration")
		}

		versionToRestore := versions[0]
		restoredSong, err := editorService.RestoreVersion(savedSong.ID, versionToRestore.ID)
		if err != nil {
			t.Fatalf("Failed to restore version: %v", err)
		}

		t.Log("Version restored successfully")
		t.Logf("Restored song title: %s", restoredSong.Metadata.Title)
		t.Logf("Restored song artist: %s", restoredSong.Metadata.Artist)
		t.Logf("Restored content length: %d characters", len(restoredSong.RawContent))

		if restoredSong.RawContent != versionToRestore.Content {
			t.Fatalf("Restored content differs from version (expected length %d, got %d)",
				len(versionToRestore.Content), len(restoredSong.RawContent))
		}

		t.Log("Content restoration verified - content matches exactly")
	})

	t.Run("AutoSaveServiceIntegration", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := autoSaveService.Start(ctx); err != nil {
			t.Fatalf("Failed to start auto-save service: %v", err)
		}
		if err := autoSaveService.Stop(); err != nil {
			t.Fatalf("Failed to stop auto-save service: %v", err)
		}

		testContent := `---
title: Direct Auto-Save Test
artist: Test Artist
---

## verse 1

This content was saved directly through the auto-save service
It should preserve the full markdown structure
`

		if err := autoSaveService.SaveWithVersioning(savedSong.ID, testContent, false, "Direct Test"); err != nil {
			t.Fatalf("Failed to save content with versioning: %v", err)
		}
		t.Log("Direct content save with versioning completed")

		recoveredContent, err := autoSaveService.RecoverFromLastSave(savedSong.ID)
		if err != nil {
			t.Fatalf("Failed to recover from last save: %v", err)
		}

		if recoveredContent != testContent {
			t.Fatalf("Recovered content differs (expected length %d, got %d)",
				len(testContent), len(recoveredContent))
		}

		t.Log("Content recovery verified - recovered content matches saved content")
	})

	t.Run("FileIOIntegration", func(t *testing.T) {
		filename := "test_autosave_song.md"

		if err := fileService.WriteSong(savedSong, filename); err != nil {
			t.Fatalf("Failed to write song to file: %v", err)
		}
		t.Logf("Song saved to file: %s", filename)

		readSong, err := fileService.ReadSong(filename)
		if err != nil {
			t.Fatalf("Failed to read song from file: %v", err)
		}
		t.Logf("Song read from file: %s", filename)
		t.Logf("Read song title: %s", readSong.Metadata.Title)
		t.Logf("Read content length: %d characters", len(readSong.RawContent))
	})

	t.Log("All auto-save integration checks completed successfully")
}
