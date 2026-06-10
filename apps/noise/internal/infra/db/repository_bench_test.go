package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/kyanite/noise/internal/domain"
)

// BenchmarkInsertSong measures the performance of song insertion operations
func BenchmarkInsertSong(b *testing.B) {
	database := newTestDB(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		song := newTestSong(fmt.Sprintf("bench-%d", i))
		if _, err := database.InsertSong(song); err != nil {
			b.Fatalf("insert song failed: %v", err)
		}
	}
}

// BenchmarkGetSong measures the performance of song retrieval operations
func BenchmarkGetSong(b *testing.B) {
	database := newTestDB(b)

	// Pre-insert songs for benchmarking
	songs := make([]*domain.Song, 100)
	for i := 0; i < 100; i++ {
		song := newTestSong(fmt.Sprintf("bench-get-%d", i))
		inserted, err := database.InsertSong(song)
		if err != nil {
			b.Fatalf("insert song failed: %v", err)
		}
		songs[i] = inserted
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		song := songs[i%100]
		if _, err := database.GetSong(song.ID); err != nil {
			b.Fatalf("get song failed: %v", err)
		}
	}
}

// BenchmarkUpdateSong measures the performance of song update operations
func BenchmarkUpdateSong(b *testing.B) {
	database := newTestDB(b)

	// Pre-insert songs for benchmarking
	songs := make([]*domain.Song, 100)
	for i := 0; i < 100; i++ {
		song := newTestSong(fmt.Sprintf("bench-update-%d", i))
		inserted, err := database.InsertSong(song)
		if err != nil {
			b.Fatalf("insert song failed: %v", err)
		}
		songs[i] = inserted
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		song := songs[i%100]
		song.Metadata.Title = fmt.Sprintf("Updated Title %d", i)
		if err := database.UpdateSong(song); err != nil {
			b.Fatalf("update song failed: %v", err)
		}
	}
}

// BenchmarkListSongs measures the performance of song listing operations
func BenchmarkListSongs(b *testing.B) {
	database := newTestDB(b)

	// Pre-insert songs for benchmarking
	for i := 0; i < 100; i++ {
		song := newTestSong(fmt.Sprintf("bench-list-%d", i))
		if _, err := database.InsertSong(song); err != nil {
			b.Fatalf("insert song failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := database.ListSongs(50, 0); err != nil {
			b.Fatalf("list songs failed: %v", err)
		}
	}
}

// BenchmarkSearchSongs measures the performance of song search operations
func BenchmarkSearchSongs(b *testing.B) {
	database := newTestDB(b)

	// Pre-insert songs for benchmarking
	for i := 0; i < 100; i++ {
		song := newTestSong(fmt.Sprintf("bench-search-%d", i))
		song.Metadata.Title = fmt.Sprintf("Searchable Song %d", i)
		song.Metadata.Tags = []string{"searchable", fmt.Sprintf("tag-%d", i)}
		if _, err := database.InsertSong(song); err != nil {
			b.Fatalf("insert song failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := database.SearchSongs("Searchable", 20); err != nil {
			b.Fatalf("search songs failed: %v", err)
		}
	}
}

// BenchmarkSaveVersion measures the performance of version save operations
func BenchmarkSaveVersion(b *testing.B) {
	database := newTestDB(b)

	// Pre-insert a song for versioning
	song := mustInsertSong(b, database, newTestSong("bench-version"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		content := fmt.Sprintf("Version content %d", i)
		if _, err := database.SaveVersion(song.ID, content, false, fmt.Sprintf("snapshot-%d", i)); err != nil {
			b.Fatalf("save version failed: %v", err)
		}
	}
}

// BenchmarkGetVersions measures the performance of version retrieval operations
func BenchmarkGetVersions(b *testing.B) {
	database := newTestDB(b)

	// Pre-insert a song and versions
	song := mustInsertSong(b, database, newTestSong("bench-get-versions"))
	for i := 0; i < 50; i++ {
		content := fmt.Sprintf("Version content %d", i)
		mustSaveVersion(b, database, song.ID, content)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := database.GetVersions(song.ID, 20); err != nil {
			b.Fatalf("get versions failed: %v", err)
		}
	}
}

// BenchmarkInsertSongWithVersion measures the performance of atomic song+version insertion
func BenchmarkInsertSongWithVersion(b *testing.B) {
	database := newTestDB(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		song := newTestSong(fmt.Sprintf("bench-atomic-%d", i))
		content := fmt.Sprintf("Initial content %d", i)
		if _, _, err := database.InsertSongWithVersion(song, content); err != nil {
			b.Fatalf("insert song with version failed: %v", err)
		}
	}
}

// BenchmarkUpdateSongWithVersion measures the performance of atomic song update+version creation
func BenchmarkUpdateSongWithVersion(b *testing.B) {
	database := newTestDB(b)

	// Pre-insert songs for benchmarking
	songs := make([]*domain.Song, 100)
	for i := 0; i < 100; i++ {
		song := newTestSong(fmt.Sprintf("bench-update-version-%d", i))
		inserted, err := database.InsertSong(song)
		if err != nil {
			b.Fatalf("insert song failed: %v", err)
		}
		songs[i] = inserted
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		song := songs[i%100]
		song.Metadata.Title = fmt.Sprintf("Updated Title %d", i)
		content := fmt.Sprintf("Updated content %d", i)
		if err := database.UpdateSongWithVersion(song, content, true, fmt.Sprintf("Milestone %d", i)); err != nil {
			b.Fatalf("update song with version failed: %v", err)
		}
	}
}

// BenchmarkRecordStats measures the performance of stats recording operations
func BenchmarkRecordStats(b *testing.B) {
	database := newTestDB(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		date := time.Now().AddDate(0, 0, i%365) // Vary dates to avoid conflicts
		stats := newTestStats(date)
		if err := database.RecordStats(stats); err != nil {
			b.Fatalf("record stats failed: %v", err)
		}
	}
}

// BenchmarkGetStats measures the performance of stats retrieval operations
func BenchmarkGetStats(b *testing.B) {
	database := newTestDB(b)

	// Pre-insert stats for benchmarking
	dates := make([]time.Time, 100)
	for i := 0; i < 100; i++ {
		date := time.Now().AddDate(0, 0, i-50) // Spread dates around today
		stats := newTestStats(date)
		if err := database.RecordStats(stats); err != nil {
			b.Fatalf("record stats failed: %v", err)
		}
		dates[i] = date
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		date := dates[i%100]
		if _, err := database.GetStats(date); err != nil {
			b.Fatalf("get stats failed: %v", err)
		}
	}
}

// BenchmarkBatchUpdateStats measures the performance of batch stats update operations
func BenchmarkBatchUpdateStats(b *testing.B) {
	database := newTestDB(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		statsList := make([]*domain.WritingStats, 10)
		for j := 0; j < 10; j++ {
			date := time.Now().AddDate(0, 0, (i*10+j)%365)
			statsList[j] = newTestStats(date)
		}
		if err := database.BatchUpdateStats(statsList); err != nil {
			b.Fatalf("batch update stats failed: %v", err)
		}
	}
}

// BenchmarkCreateProject measures the performance of project creation operations
func BenchmarkCreateProject(b *testing.B) {
	database := newTestDB(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		project := &domain.Project{
			Name:        fmt.Sprintf("Benchmark Project %d", i),
			Description: fmt.Sprintf("Description for project %d", i),
			SongIDs:     []int{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if _, err := database.CreateProject(project); err != nil {
			b.Fatalf("create project failed: %v", err)
		}
	}
}

// BenchmarkAddSongToProject measures the performance of adding songs to projects
func BenchmarkAddSongToProject(b *testing.B) {
	database := newTestDB(b)

	// Pre-insert projects and songs for benchmarking
	projects := make([]*domain.Project, 10)
	for i := 0; i < 10; i++ {
		project := &domain.Project{
			Name:        fmt.Sprintf("Benchmark Project %d", i),
			Description: fmt.Sprintf("Description for project %d", i),
			SongIDs:     []int{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		created, err := database.CreateProject(project)
		if err != nil {
			b.Fatalf("create project failed: %v", err)
		}
		projects[i] = created
	}

	songs := make([]*domain.Song, 100)
	for i := 0; i < 100; i++ {
		song := newTestSong(fmt.Sprintf("bench-project-song-%d", i))
		inserted, err := database.InsertSong(song)
		if err != nil {
			b.Fatalf("insert song failed: %v", err)
		}
		songs[i] = inserted
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		project := projects[i%10]
		song := songs[i%100]
		if err := database.AddSongToProject(project.ID, song.ID); err != nil {
			b.Fatalf("add song to project failed: %v", err)
		}
	}
}

// BenchmarkExecuteInTransaction measures the performance of transaction operations
func BenchmarkExecuteInTransaction(b *testing.B) {
	database := newTestDB(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := database.ExecuteInTransaction(context.Background(), func(tx *sql.Tx) error {
			// Insert a song within the transaction
			song := newTestSong(fmt.Sprintf("bench-tx-%d", i))
			tagsJSON, _ := marshalStringArray(song.Metadata.Tags)

			_, err := tx.Exec(`
				INSERT INTO songs (filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				song.Filepath,
				song.Metadata.Title,
				song.Metadata.Artist,
				song.Metadata.Key,
				song.Metadata.Tempo,
				song.Metadata.TimeSignature,
				song.Metadata.Structure,
				tagsJSON,
				song.Metadata.CreatedAt,
				song.Metadata.UpdatedAt,
			)
			return err
		})
		if err != nil {
			b.Fatalf("execute in transaction failed: %v", err)
		}
	}
}
