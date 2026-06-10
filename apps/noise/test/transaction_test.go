package noise

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/infra/db"
)

// TestDatabaseTransactions tests the transaction functionality
func TestDatabaseTransactions(t *testing.T) {
	// Create temporary directory for test database
	tempDir, err := os.MkdirTemp("", "noise.sh_test_transactions")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize database
	dbConfig := db.Config{DataDir: tempDir}
	database, err := db.New(dbConfig)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	t.Run("TestProjectOperations", func(t *testing.T) {
		testProjectTransactionOperations(t, database)
	})

	t.Run("TestSongWithVersionOperations", func(t *testing.T) {
		testSongWithVersionTransactionOperations(t, database)
	})

	t.Run("TestBatchStatsOperations", func(t *testing.T) {
		testBatchStatsTransactionOperations(t, database)
	})

	t.Run("TestTransactionRollback", func(t *testing.T) {
		testTransactionRollback(t, database)
	})

	t.Run("TestExecuteInTransaction", func(t *testing.T) {
		testExecuteInTransaction(t, database)
	})
}

func testProjectTransactionOperations(t *testing.T, database *db.DB) {
	// Create a test project
	project := &domain.Project{
		Name:        "Test Project",
		Description: "A test project for transaction testing",
		SongIDs:     []int{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdProject, err := database.CreateProject(project)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Create a test song first
	song := &domain.Song{
		Filepath: "/test/song.txt",
		Metadata: domain.SongMetadata{
			Title:     "Test Song",
			Artist:    "Test Artist",
			Key:       "C",
			Tempo:     120,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	createdSong, err := database.InsertSong(song)
	if err != nil {
		t.Fatalf("Failed to create song: %v", err)
	}

	// Test transactional AddSongToProject
	err = database.AddSongToProject(createdProject.ID, createdSong.ID)
	if err != nil {
		t.Fatalf("Failed to add song to project: %v", err)
	}

	// Verify the song was added
	updatedProject, err := database.GetProject(createdProject.ID)
	if err != nil {
		t.Fatalf("Failed to get project: %v", err)
	}

	if len(updatedProject.SongIDs) != 1 || updatedProject.SongIDs[0] != createdSong.ID {
		t.Fatalf("Song was not properly added to project")
	}

	// Test transactional RemoveSongFromProject
	err = database.RemoveSongFromProject(createdProject.ID, createdSong.ID)
	if err != nil {
		t.Fatalf("Failed to remove song from project: %v", err)
	}

	// Verify the song was removed
	updatedProject, err = database.GetProject(createdProject.ID)
	if err != nil {
		t.Fatalf("Failed to get project: %v", err)
	}

	if len(updatedProject.SongIDs) != 0 {
		t.Fatalf("Song was not properly removed from project")
	}
}

func testSongWithVersionTransactionOperations(t *testing.T, database *db.DB) {
	// Create a test song
	song := &domain.Song{
		Filepath: "/test/song_with_version.txt",
		Metadata: domain.SongMetadata{
			Title:     "Test Song with Version",
			Artist:    "Test Artist",
			Key:       "G",
			Tempo:     100,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	initialContent := "This is the initial content of the song."

	// Test transactional InsertSongWithVersion
	createdSong, version, err := database.InsertSongWithVersion(song, initialContent)
	if err != nil {
		t.Fatalf("Failed to insert song with version: %v", err)
	}

	// Verify song was created
	if createdSong.ID == 0 {
		t.Fatalf("Song ID should not be zero")
	}

	// Verify version was created
	if version == nil {
		t.Fatalf("Version should not be nil")
	}

	if version.SongID != createdSong.ID {
		t.Fatalf("Version should reference the created song")
	}

	if version.Content != initialContent {
		t.Fatalf("Version content should match initial content")
	}

	// Test transactional UpdateSongWithVersion
	newContent := "This is the updated content of the song."
	err = database.UpdateSongWithVersion(createdSong, newContent, true, "Updated version")
	if err != nil {
		t.Fatalf("Failed to update song with version: %v", err)
	}

	// Verify version was created for the update
	versions, err := database.GetVersions(createdSong.ID, 10)
	if err != nil {
		t.Fatalf("Failed to get versions: %v", err)
	}

	if len(versions) != 2 {
		t.Fatalf("Expected 2 versions, got %d", len(versions))
	}

	// Find the milestone version
	var milestoneVersion *domain.Version
	for _, v := range versions {
		if v.IsMilestone && v.MilestoneName == "Updated version" {
			milestoneVersion = v
			break
		}
	}

	if milestoneVersion == nil {
		t.Fatalf("Milestone version not found")
	}

	if milestoneVersion.Content != newContent {
		t.Fatalf("Milestone version content should match updated content")
	}
}

func testBatchStatsTransactionOperations(t *testing.T, database *db.DB) {
	// Create test stats
	stats1 := &domain.WritingStats{
		Date:             time.Now().AddDate(0, 0, -2),
		WordsWritten:     100,
		SongsCreated:     1,
		SongsEdited:      2,
		AIRequests:       5,
		TimeSpentMinutes: 60,
	}

	stats2 := &domain.WritingStats{
		Date:             time.Now().AddDate(0, 0, -1),
		WordsWritten:     150,
		SongsCreated:     2,
		SongsEdited:      1,
		AIRequests:       3,
		TimeSpentMinutes: 45,
	}

	statsList := []*domain.WritingStats{stats1, stats2}

	// Test batch update
	err := database.BatchUpdateStats(statsList)
	if err != nil {
		t.Fatalf("Failed to batch update stats: %v", err)
	}

	// Verify stats were recorded
	for _, stats := range statsList {
		retrievedStats, err := database.GetStats(stats.Date)
		if err != nil {
			t.Fatalf("Failed to get stats for date %v: %v", stats.Date, err)
		}

		if retrievedStats.WordsWritten != stats.WordsWritten {
			t.Fatalf("Words written mismatch for date %v", stats.Date)
		}

		if retrievedStats.SongsCreated != stats.SongsCreated {
			t.Fatalf("Songs created mismatch for date %v", stats.Date)
		}
	}
}

func testTransactionRollback(t *testing.T, database *db.DB) {
	// Create a test project
	project := &domain.Project{
		Name:        "Rollback Test Project",
		Description: "A test project for rollback testing",
		SongIDs:     []int{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdProject, err := database.CreateProject(project)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Create a test song
	song := &domain.Song{
		Filepath: "/test/rollback_song.txt",
		Metadata: domain.SongMetadata{
			Title:     "Rollback Test Song",
			Artist:    "Test Artist",
			Key:       "D",
			Tempo:     110,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	createdSong, err := database.InsertSong(song)
	if err != nil {
		t.Fatalf("Failed to create song: %v", err)
	}

	// Get initial project state
	initialProject, err := database.GetProject(createdProject.ID)
	if err != nil {
		t.Fatalf("Failed to get initial project: %v", err)
	}

	initialSongCount := len(initialProject.SongIDs)

	// Use ExecuteInTransaction to simulate a failing operation
	ctx := context.Background()
	err = database.ExecuteInTransaction(ctx, func(tx *sql.Tx) error {
		// This should succeed
		_, err := tx.Exec("UPDATE projects SET name = ? WHERE id = ?", "Updated Name", createdProject.ID)
		if err != nil {
			return err
		}

		// This should also succeed
		_, err = tx.Exec("INSERT INTO projects (name, description, song_ids, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			"Temp Project", "Temp Description", "[]", time.Now(), time.Now())
		if err != nil {
			return err
		}

		// Force a failure to test rollback
		return fmt.Errorf("intentional failure for rollback test")
	})

	// The transaction should have failed
	if err == nil {
		t.Fatalf("Expected transaction to fail")
	}

	// Verify rollback worked - project name should not have changed
	rolledBackProject, err := database.GetProject(createdProject.ID)
	if err != nil {
		t.Fatalf("Failed to get project after rollback: %v", err)
	}

	if rolledBackProject.Name == "Updated Name" {
		t.Fatalf("Project name should not have changed due to rollback")
	}

	// Verify song count is unchanged
	if len(rolledBackProject.SongIDs) != initialSongCount {
		t.Fatalf("Song count should not have changed due to rollback")
	}

	// Test successful transaction with AddSongToProject
	err = database.AddSongToProject(createdProject.ID, createdSong.ID)
	if err != nil {
		t.Fatalf("Failed to add song to project: %v", err)
	}

	// Verify the addition worked
	finalProject, err := database.GetProject(createdProject.ID)
	if err != nil {
		t.Fatalf("Failed to get final project: %v", err)
	}

	if len(finalProject.SongIDs) != initialSongCount+1 {
		t.Fatalf("Song should have been added to project")
	}
}

func testExecuteInTransaction(t *testing.T, database *db.DB) {
	// Create test data
	project := &domain.Project{
		Name:        "ExecuteInTransaction Test",
		Description: "Testing ExecuteInTransaction functionality",
		SongIDs:     []int{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdProject, err := database.CreateProject(project)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	song := &domain.Song{
		Filepath: "/test/execute_tx_song.txt",
		Metadata: domain.SongMetadata{
			Title:     "ExecuteInTransaction Song",
			Artist:    "Test Artist",
			Key:       "A",
			Tempo:     130,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	createdSong, err := database.InsertSong(song)
	if err != nil {
		t.Fatalf("Failed to create song: %v", err)
	}

	// Use the createdSong variable to avoid unused variable error
	_ = createdSong

	// Test successful ExecuteInTransaction
	ctx := context.Background()
	err = database.ExecuteInTransaction(ctx, func(tx *sql.Tx) error {
		// Update project name
		_, err := tx.Exec("UPDATE projects SET name = ?, updated_at = ? WHERE id = ?",
			"Updated in Transaction", time.Now(), createdProject.ID)
		if err != nil {
			return err
		}

		// The song addition should be handled by the existing transactional method
		return nil
	})
	if err != nil {
		t.Fatalf("ExecuteInTransaction failed: %v", err)
	}

	// Verify the update worked
	updatedProject, err := database.GetProject(createdProject.ID)
	if err != nil {
		t.Fatalf("Failed to get updated project: %v", err)
	}

	if updatedProject.Name != "Updated in Transaction" {
		t.Fatalf("Project name should have been updated in transaction")
	}
}

// BenchmarkTransactionPerformance benchmarks transaction performance
func BenchmarkTransactionPerformance(b *testing.B) {
	// Create temporary directory for benchmark database
	tempDir, err := os.MkdirTemp("", "noise.sh_benchmark_transactions")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize database
	dbConfig := db.Config{DataDir: tempDir}
	database, err := db.New(dbConfig)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	// Create test data
	project := &domain.Project{
		Name:        "Benchmark Project",
		Description: "Project for benchmarking transactions",
		SongIDs:     []int{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdProject, err := database.CreateProject(project)
	if err != nil {
		b.Fatalf("Failed to create project: %v", err)
	}

	b.ResetTimer()

	b.Run("TransactionalAddSong", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			song := &domain.Song{
				Filepath: fmt.Sprintf("/test/benchmark_song_%d.txt", i),
				Metadata: domain.SongMetadata{
					Title:     fmt.Sprintf("Benchmark Song %d", i),
					Artist:    "Test Artist",
					Key:       "C",
					Tempo:     120,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			createdSong, err := database.InsertSong(song)
			if err != nil {
				b.Fatalf("Failed to create song: %v", err)
			}

			err = database.AddSongToProject(createdProject.ID, createdSong.ID)
			if err != nil {
				b.Fatalf("Failed to add song to project: %v", err)
			}

			// Use the createdSong variable to avoid unused variable error
			_ = createdSong
		}
	})

	b.Run("NonTransactionalAddSong", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			song := &domain.Song{
				Filepath: fmt.Sprintf("/test/benchmark_song_nt_%d.txt", i),
				Metadata: domain.SongMetadata{
					Title:     fmt.Sprintf("Benchmark Song NT %d", i),
					Artist:    "Test Artist",
					Key:       "C",
					Tempo:     120,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			createdSong, err := database.InsertSong(song)
			if err != nil {
				b.Fatalf("Failed to create song: %v", err)
			}

			err = database.AddSongToProjectNonTx(createdProject.ID, createdSong.ID)
			if err != nil {
				b.Fatalf("Failed to add song to project: %v", err)
			}

			// Use the createdSong variable to avoid unused variable error
			_ = createdSong
		}
	})
}
