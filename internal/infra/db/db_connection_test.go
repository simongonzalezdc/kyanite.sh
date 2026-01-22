package db

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestNewInitializesSchema(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)

	// Core tables required for single-user mode
	// Collaboration tables are only created when EnableCollaboration is true
	requiredTables := []string{
		"songs",
		"versions",
		"writing_stats",
		"kb_entries",
		"projects",
	}

	for _, table := range requiredTables {
		var name string
		err := database.conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Fatalf("expected table %q to exist, query failed: %v", table, err)
		}
		if name != table {
			t.Fatalf("expected table name %q, got %q", table, name)
		}
	}

	// Validate foreign key linkage for versions -> songs
	rows, err := database.conn.Query("PRAGMA foreign_key_list(versions)")
	if err != nil {
		t.Fatalf("failed to inspect foreign key list: %v", err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var (
			id, seq            int
			table, from, to    string
			onUpdate, onDelete string
			match              string
		)
		if scanErr := rows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); scanErr != nil {
			t.Fatalf("failed to scan foreign key row: %v", scanErr)
		}
		if table == "songs" && from == "song_id" && to == "id" && onDelete == "CASCADE" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected versions.song_id to reference songs.id with ON DELETE CASCADE")
	}
}

func TestConnectionPoolConfiguration(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)
	stats := database.conn.Stats()

	if stats.MaxOpenConnections != 10 {
		t.Fatalf("expected MaxOpenConnections=10, got %d", stats.MaxOpenConnections)
	}
	if stats.Idle > 5 {
		t.Fatalf("expected idle connections <= 5, got %d", stats.Idle)
	}
}

func TestConnectionRecoveryAfterFailure(t *testing.T) {
	t.Parallel()

	dataDir := filepath.Join(t.TempDir(), "persistent")
	cfg := Config{DataDir: dataDir}

	database, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if cerr := database.Close(); cerr != nil && cerr != sql.ErrConnDone {
			t.Fatalf("failed to close database: %v", cerr)
		}
	}()

	// Simulate connection failure.
	if err := database.conn.Close(); err != nil {
		t.Fatalf("failed to close underlying connection: %v", err)
	}

	if err := database.Ping(); err == nil {
		t.Fatal("expected ping to fail after closing connection")
	}

	// Recovery: create new connection with same configuration.
	recovered, err := New(cfg)
	if err != nil {
		t.Fatalf("expected recovery to succeed, got error: %v", err)
	}
	t.Cleanup(func() {
		_ = recovered.Close()
	})

	if err := recovered.Ping(); err != nil {
		t.Fatalf("expected recovered database ping to succeed, got: %v", err)
	}
}

func TestValidateConnectionHealth(t *testing.T) {
	t.Parallel()

	database := newTestDB(t)

	if err := database.validateConnection(); err != nil {
		t.Fatalf("expected healthy connection, got error: %v", err)
	}

	// Close the connection to trigger validation failure.
	if err := database.conn.Close(); err != nil {
		t.Fatalf("failed to close connection: %v", err)
	}

	if err := database.validateConnection(); err == nil {
		t.Fatal("expected validateConnection to fail after closing connection")
	}
}
