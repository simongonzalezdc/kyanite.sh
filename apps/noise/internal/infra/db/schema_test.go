package db

import (
	"testing"
	"time"
)

func TestSchemaIndexesCreated(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)

	rows, err := database.conn.Query(`SELECT name FROM sqlite_master WHERE type='index'`)
	if err != nil {
		t.Fatalf("failed to query sqlite_master: %v", err)
	}
	defer rows.Close()

	indexes := make(map[string]struct{})
	for rows.Next() {
		var name string
		if scanErr := rows.Scan(&name); scanErr != nil {
			t.Fatalf("failed to scan index row: %v", scanErr)
		}
		indexes[name] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterating index rows failed: %v", err)
	}

	expected := []string{
		"idx_songs_updated",
		"idx_songs_title",
		"idx_versions_song",
		"idx_stats_date",
		"idx_projects_name",
	}

	for _, idx := range expected {
		if _, ok := indexes[idx]; !ok {
			t.Fatalf("expected index %q to exist", idx)
		}
	}
}

func TestSchemaForeignKeyCascade(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)

	song := mustInsertSong(t, database, newTestSong("schema-cascade"))
	mustSaveVersion(t, database, song.ID, "foreign key cascade content")

	if err := database.DeleteSong(song.ID); err != nil {
		t.Fatalf("DeleteSong failed: %v", err)
	}

	var count int
	if err := database.conn.QueryRow(`SELECT COUNT(1) FROM versions WHERE song_id = ?`, song.ID).Scan(&count); err != nil {
		t.Fatalf("failed to count versions: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected cascade delete of versions, found %d remaining", count)
	}
}

func TestSchemaUniqueConstraints(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)
	date := time.Now().Truncate(24 * time.Hour)

	if _, err := database.conn.Exec(`INSERT INTO writing_stats (date) VALUES (?)`, date); err != nil {
		t.Fatalf("expected first insert to succeed, got %v", err)
	}

	if _, err := database.conn.Exec(`INSERT INTO writing_stats (date) VALUES (?)`, date); err == nil {
		t.Fatal("expected unique constraint violation on duplicate date insert")
	}
}
