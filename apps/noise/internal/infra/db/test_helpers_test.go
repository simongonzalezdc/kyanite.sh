package db

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyanite/noise/internal/domain"
)

func newTestConfig(t testing.TB) Config {
	t.Helper()
	dataDir := filepath.Join(t.TempDir(), "db")
	return Config{DataDir: dataDir}
}

func newTestDB(t testing.TB) *DB {
	t.Helper()
	database, err := New(newTestConfig(t))
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Logf("Warning: failed to close test database: %v", err)
		}
	})
	return database
}

func newTestSong(suffix string) *domain.Song {
	now := time.Now().UTC().Truncate(time.Second)
	return &domain.Song{
		Filepath: fmt.Sprintf("song-%s.noise", suffix),
		Metadata: domain.SongMetadata{
			Title:         fmt.Sprintf("Song %s", suffix),
			Artist:        "Test Artist",
			Key:           "C",
			Tempo:         120,
			TimeSignature: "4/4",
			Structure:     "Intro-Verse-Chorus",
			Tags:          []string{"test", suffix},
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		Sections: []domain.Section{},
	}
}

func mustInsertSong(t testing.TB, database *DB, song *domain.Song) *domain.Song {
	t.Helper()
	inserted, err := database.InsertSong(song)
	if err != nil {
		t.Fatalf("insert song failed: %v", err)
	}
	return inserted
}

func mustSaveVersion(t testing.TB, database *DB, songID int, content string) *domain.Version {
	t.Helper()
	version, err := database.SaveVersion(songID, content, false, fmt.Sprintf("snapshot-%d", time.Now().UnixNano()))
	if err != nil {
		t.Fatalf("save version failed: %v", err)
	}
	return version
}

func mustGetVersions(t testing.TB, database *DB, songID, limit int) []*domain.Version {
	t.Helper()
	versions, err := database.GetVersions(songID, limit)
	if err != nil {
		t.Fatalf("get versions failed: %v", err)
	}
	return versions
}

func newTestStats(date time.Time) *domain.WritingStats {
	return &domain.WritingStats{
		Date:             date,
		WordsWritten:     100,
		SongsCreated:     2,
		SongsEdited:      3,
		AIRequests:       5,
		TimeSpentMinutes: 60,
	}
}
