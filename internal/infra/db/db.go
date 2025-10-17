package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/puente-labs/lyricforge/internal/domain"
)

// DB wraps the database connection with helper methods
type DB struct {
	conn *sql.DB
}

// Config holds database configuration
type Config struct {
	DataDir string
}

// New creates a new database connection and initializes the schema
func New(cfg Config) (*DB, error) {
	if cfg.DataDir == "" {
		cfg.DataDir = getDefaultDataDir()
	}

	dbPath := filepath.Join(cfg.DataDir, "lyricforge.db")

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(time.Hour)

	// Enable foreign keys
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Create schema
	if err := initializeSchema(conn); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// Ping tests the database connection
func (db *DB) Ping() error {
	return db.conn.Ping()
}

// initializeSchema creates all tables if they don't exist
func initializeSchema(conn *sql.DB) error {
	if _, err := conn.Exec(Schema); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}

// getDefaultDataDir returns the default data directory
func getDefaultDataDir() string {
	homeDir, err := getHomeDir()
	if err != nil {
		return "./data" // fallback
	}
	return filepath.Join(homeDir, ".lyricforge")
}

// getHomeDir returns the user's home directory (cross-platform)
func getHomeDir() (string, error) {
	// This is a simplified implementation
	// In a real application, you'd use os.UserHomeDir() (Go 1.12+)
	// or a cross-platform library like github.com/mitchellh/go-homedir
	return "./data", nil
}

// SongRepository defines the interface for song persistence operations
type SongRepository interface {
	InsertSong(song *domain.Song) (*domain.Song, error)
	GetSong(id int) (*domain.Song, error)
	GetSongByFilepath(filepath string) (*domain.Song, error)
	UpdateSong(song *domain.Song) error
	DeleteSong(id int) error
	ListSongs(limit, offset int) ([]*domain.Song, error)
	SearchSongs(query string, limit int) ([]*domain.Song, error)
}

// VersionRepository defines the interface for version history operations
type VersionRepository interface {
	SaveVersion(songID int, content string, isMilestone bool, name string) (*domain.Version, error)
	GetVersions(songID int, limit int) ([]*domain.Version, error)
	GetVersion(id int) (*domain.Version, error)
	DeleteVersion(id int) error
}

// StatsRepository defines the interface for writing statistics operations
type StatsRepository interface {
	RecordStats(stats *domain.WritingStats) error
	GetStats(date time.Time) (*domain.WritingStats, error)
	GetStatsRange(start, end time.Time) ([]*domain.WritingStats, error)
	UpdateStats(stats *domain.WritingStats) error
}

// ProjectRepository defines the interface for project operations
type ProjectRepository interface {
	CreateProject(project *domain.Project) (*domain.Project, error)
	GetProject(id int) (*domain.Project, error)
	UpdateProject(project *domain.Project) error
	DeleteProject(id int) error
	ListProjects() ([]*domain.Project, error)
	AddSongToProject(projectID, songID int) error
	RemoveSongFromProject(projectID, songID int) error
}

// Ensure DB implements the repository interfaces
var _ SongRepository = (*DB)(nil)
var _ VersionRepository = (*DB)(nil)
var _ StatsRepository = (*DB)(nil)
var _ ProjectRepository = (*DB)(nil)

// Helper functions for JSON marshaling/unmarshaling

// unmarshalStringArray converts JSON string array to []string
func unmarshalStringArray(jsonStr string) ([]string, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var arr []string
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		return nil, err
	}
	return arr, nil
}

// marshalStringArray converts []string to JSON string
func marshalStringArray(arr []string) (string, error) {
	if arr == nil {
		return "", nil
	}
	jsonBytes, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// unmarshalIntArray converts JSON int array to []int
func unmarshalIntArray(jsonStr string) ([]int, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var arr []int
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		return nil, err
	}
	return arr, nil
}

// marshalIntArray converts []int to JSON string
func marshalIntArray(arr []int) (string, error) {
	if arr == nil {
		return "", nil
	}
	jsonBytes, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
