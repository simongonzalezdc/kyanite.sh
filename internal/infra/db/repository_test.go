package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/domain"
	appErrors "github.com/Kyanite/noise/internal/errors"
)

func TestInsertSongCreatesPersistedRecord(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)
	song := newTestSong("insert-basic")

	inserted, err := database.InsertSong(song)
	if err != nil {
		t.Fatalf("insert song failed: %v", err)
	}

	if inserted.ID == 0 {
		t.Fatal("expected inserted song to have non-zero ID")
	}

	fetched, err := database.GetSong(inserted.ID)
	if err != nil {
		t.Fatalf("get song failed: %v", err)
	}

	if fetched.Metadata.Title != song.Metadata.Title {
		t.Fatalf("expected title %q, got %q", song.Metadata.Title, fetched.Metadata.Title)
	}

	if fetched.Filepath != song.Filepath {
		t.Fatalf("expected filepath %q, got %q", song.Filepath, fetched.Filepath)
	}
}

func TestInsertSongValidationFailures(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)

	_, err := database.InsertSong(nil)
	var validationErr *appErrors.AppError
	if !errors.As(err, &validationErr) || !validationErr.HasCategory(appErrors.CategoryValidation) {
		t.Fatalf("expected validation error when inserting nil song, got %v", err)
	}

	badSong := newTestSong("validation")
	badSong.Metadata.Title = ""

	_, err = database.InsertSong(badSong)
	validationErr = nil
	if !errors.As(err, &validationErr) || !validationErr.HasCategory(appErrors.CategoryValidation) {
		t.Fatalf("expected validation error for missing title, got %v", err)
	}
}

func TestInsertSongValidationScenarios(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)

	testCases := []struct {
		name        string
		modify      func(*domain.Song)
		expectError bool
		category    appErrors.ErrorCategory
	}{
		{
			name:        "valid-song",
			expectError: false,
		},
		{
			name: "missing-title",
			modify: func(song *domain.Song) {
				song.Metadata.Title = ""
			},
			expectError: true,
			category:    appErrors.CategoryValidation,
		},
		{
			name: "invalid-filepath",
			modify: func(song *domain.Song) {
				song.Filepath = "bad:path?"
			},
			expectError: true,
			category:    appErrors.CategoryValidation,
		},
		{
			name: "invalid-tempo",
			modify: func(song *domain.Song) {
				song.Metadata.Tempo = 400
			},
			expectError: true,
			category:    appErrors.CategoryValidation,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			song := newTestSong("scenario-" + tc.name)
			if tc.modify != nil {
				tc.modify(song)
			}

			inserted, err := database.InsertSong(song)
			if tc.expectError {
				var appErr *appErrors.AppError
				if !errors.As(err, &appErr) {
					t.Fatalf("expected app error, got %v", err)
				}
				if tc.category != "" && !appErr.HasCategory(tc.category) {
					t.Fatalf("expected category %q, got %q", tc.category, appErr.Category)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected insert to succeed, got %v", err)
			}
			if inserted == nil || inserted.ID == 0 {
				t.Fatal("expected inserted song with assigned ID")
			}
		})
	}
}

func TestInsertSongDuplicateFilepath(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)
	song := newTestSong("duplicate")

	if _, err := database.InsertSong(song); err != nil {
		t.Fatalf("initial insert failed: %v", err)
	}

	_, err := database.InsertSong(song)
	if err == nil {
		t.Fatal("expected duplicate insert to return error")
	}

	var dbErr *appErrors.AppError
	if !errors.As(err, &dbErr) {
		t.Fatalf("expected database error type, got %T", err)
	}
	if !strings.Contains(strings.ToLower(dbErr.Error()), "duplicate") {
		t.Fatalf("expected duplicate constraint error, got %v", dbErr)
	}
}

func TestInsertSongWithVersionTransaction(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)
	song := newTestSong("with-version")

	insertedSong, version, err := database.InsertSongWithVersion(song, "initial content")
	if err != nil {
		t.Fatalf("InsertSongWithVersion failed: %v", err)
	}

	if insertedSong.ID == 0 || version.ID == 0 {
		t.Fatal("expected song and version IDs to be assigned")
	}

	if version.SongID != insertedSong.ID {
		t.Fatalf("expected version song ID %d, got %d", insertedSong.ID, version.SongID)
	}

	fetchedVersions, err := database.GetVersions(insertedSong.ID, 10)
	if err != nil {
		t.Fatalf("GetVersions failed: %v", err)
	}
	if len(fetchedVersions) != 1 {
		t.Fatalf("expected exactly one version, got %d", len(fetchedVersions))
	}
	if fetchedVersions[0].Content != "initial content" {
		t.Fatalf("expected version content %q, got %q", "initial content", fetchedVersions[0].Content)
	}
}

func TestUpdateSongWithVersionIsAtomic(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)
	song := mustInsertSong(t, database, newTestSong("update-atomic"))

	song.Metadata.Title = "Updated Title"
	err := database.UpdateSongWithVersion(song, "updated content", true, "Milestone 1")
	if err != nil {
		t.Fatalf("UpdateSongWithVersion failed: %v", err)
	}

	fetched, err := database.GetSong(song.ID)
	if err != nil {
		t.Fatalf("GetSong failed: %v", err)
	}
	if fetched.Metadata.Title != "Updated Title" {
		t.Fatalf("expected updated title, got %q", fetched.Metadata.Title)
	}

	versions, err := database.GetVersions(song.ID, 10)
	if err != nil {
		t.Fatalf("GetVersions failed: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("expected single version snapshot, got %d", len(versions))
	}
	if !versions[0].IsMilestone || versions[0].MilestoneName != "Milestone 1" {
		t.Fatalf("expected milestone snapshot with name 'Milestone 1', got %+v", versions[0])
	}
}

func TestExecuteInTransactionRollbackOnError(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)

	err := database.ExecuteInTransaction(context.Background(), func(tx *sql.Tx) error {
		if _, execErr := tx.Exec(`INSERT INTO songs (filepath, title, created_at, updated_at) VALUES (?, ?, ?, ?)`,
			"tx/song.noise", "Transactional Song", time.Now(), time.Now()); execErr != nil {
			return execErr
		}
		return errors.New("force rollback")
	})
	if err == nil {
		t.Fatal("expected ExecuteInTransaction to propagate error")
	}

	_, err = database.GetSongByFilepath("tx/song.noise")
	if err == nil {
		t.Fatal("expected song insert to be rolled back on error")
	}
}

func TestSongCRUDOperations(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)

	s1 := mustInsertSong(t, database, newTestSong("crud-1"))
	s2 := mustInsertSong(t, database, newTestSong("crud-2"))

	list, err := database.ListSongs(10, 0)
	if err != nil {
		t.Fatalf("ListSongs failed: %v", err)
	}
	if len(list) < 2 {
		t.Fatalf("expected at least two songs, got %d", len(list))
	}

	s1.Metadata.Title = "Crud Update Title"
	if err := database.UpdateSong(s1); err != nil {
		t.Fatalf("UpdateSong failed: %v", err)
	}

	updated, err := database.GetSong(s1.ID)
	if err != nil {
		t.Fatalf("GetSong failed: %v", err)
	}
	if updated.Metadata.Title != "Crud Update Title" {
		t.Fatalf("expected updated title, got %q", updated.Metadata.Title)
	}

	if err := database.DeleteSong(s2.ID); err != nil {
		t.Fatalf("DeleteSong failed: %v", err)
	}

	if err := database.DeleteSong(999999); err == nil {
		t.Fatal("expected error when deleting non-existent song")
	}
}

func TestSearchSongsMatchesTitleAndTags(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)

	song := newTestSong("searchable")
	song.Metadata.Tags = []string{"ambient", "focus"}
	song.Metadata.Title = "Ambient Focus Song"

	mustInsertSong(t, database, song)

	results, err := database.SearchSongs("Ambient", 5)
	if err != nil {
		t.Fatalf("SearchSongs failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected search to return at least one result")
	}

	found := false
	for _, res := range results {
		if res.ID == song.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected inserted song to be returned by search")
	}
}

func TestSaveVersionAndGetVersion(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)
	song := mustInsertSong(t, database, newTestSong("versioning"))

	version, err := database.SaveVersion(song.ID, "content v1", false, "snapshot-1")
	if err != nil {
		t.Fatalf("SaveVersion failed: %v", err)
	}

	fetched, err := database.GetVersion(version.ID)
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}

	if fetched.Content != "content v1" {
		t.Fatalf("expected version content %q, got %q", "content v1", fetched.Content)
	}

	if err := database.DeleteVersion(version.ID); err != nil {
		t.Fatalf("DeleteVersion failed: %v", err)
	}

	if _, err := database.GetVersion(version.ID); err == nil {
		t.Fatal("expected error retrieving deleted version")
	}
}

func TestStatsRepositoryOperations(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)
	date := time.Now().Truncate(24 * time.Hour)

	stats := newTestStats(date)
	if err := database.RecordStats(stats); err != nil {
		t.Fatalf("RecordStats failed: %v", err)
	}

	fetched, err := database.GetStats(date)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if fetched.WordsWritten != stats.WordsWritten {
		t.Fatalf("expected words written %d, got %d", stats.WordsWritten, fetched.WordsWritten)
	}

	fetched.WordsWritten = 250
	if err := database.UpdateStats(fetched); err != nil {
		t.Fatalf("UpdateStats failed: %v", err)
	}

	rangeStats, err := database.GetStatsRange(date.Add(-24*time.Hour), date.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("GetStatsRange failed: %v", err)
	}
	if len(rangeStats) != 1 {
		t.Fatalf("expected stats range to return single record, got %d", len(rangeStats))
	}

	if err := database.BatchUpdateStats([]*domain.WritingStats{
		{
			Date:             date.Add(24 * time.Hour),
			WordsWritten:     500,
			SongsCreated:     1,
			SongsEdited:      1,
			AIRequests:       2,
			TimeSpentMinutes: 90,
		},
	}); err != nil {
		t.Fatalf("BatchUpdateStats failed: %v", err)
	}
}

func TestProjectRepositoryOperations(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)
	project := &domain.Project{
		Name:        "Collab Project",
		Description: "Test project for repository operations",
		SongIDs:     []int{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	created, err := database.CreateProject(project)
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	fetched, err := database.GetProject(created.ID)
	if err != nil {
		t.Fatalf("GetProject failed: %v", err)
	}
	if fetched.Name != project.Name {
		t.Fatalf("expected project name %q, got %q", project.Name, fetched.Name)
	}

	song := mustInsertSong(t, database, newTestSong("project-song"))
	if err := database.AddSongToProject(created.ID, song.ID); err != nil {
		t.Fatalf("AddSongToProject failed: %v", err)
	}

	fetched, err = database.GetProject(created.ID)
	if err != nil {
		t.Fatalf("GetProject failed after add: %v", err)
	}
	if len(fetched.SongIDs) != 1 || fetched.SongIDs[0] != song.ID {
		t.Fatalf("expected song ID %d in project, got %v", song.ID, fetched.SongIDs)
	}

	if err := database.RemoveSongFromProject(created.ID, song.ID); err != nil {
		t.Fatalf("RemoveSongFromProject failed: %v", err)
	}

	if err := database.DeleteProject(created.ID); err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}
}

func TestValidateConnectionDetectsClosedDB(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)
	if err := database.validateConnection(); err != nil {
		t.Fatalf("expected healthy connection, got %v", err)
	}

	if err := database.conn.Close(); err != nil {
		t.Fatalf("failed to close connection: %v", err)
	}

	if err := database.validateConnection(); err == nil {
		t.Fatal("expected validation failure after closing connection")
	}
}
