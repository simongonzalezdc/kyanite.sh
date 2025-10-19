// Package db provides database persistence and repository implementations for the noise.sh application.
// It supports both SQLite with CGO and in-memory fallback storage for environments where CGO is unavailable.
package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	errutil "github.com/puente-labs/noise/internal/errutil"
	"github.com/puente-labs/noise/internal/domain"
)

// DB wraps the database connection with helper methods
type DB struct {
	conn *sql.DB

	// In-memory fallback storage when SQLite/CGo is unavailable or when SQL persistence
	// cannot be used (e.g., CGO disabled or foreign key constraints unsupported).
	versionMutex  sync.Mutex
	versions      []*domain.Version
	nextVersionID int

	fallbackMu      sync.RWMutex
	versionFallback bool
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

	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return nil, errutil.Wrap(err, "prepare data directory")
	}

	dbPath := filepath.Join(cfg.DataDir, "noise.sh.db")

	conn, err := sql.Open(sqliteDriverName, dbPath)
	if err != nil {
		if isDriverUnavailable(err) {
			return newFallbackDB(), nil
		}
		return nil, errutil.Wrap(err, "open database")
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		if isDriverUnavailable(err) {
			return newFallbackDB(), nil
		}
		return nil, errutil.Wrap(err, "ping database")
	}

	// Configure connection pool
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(time.Hour)

	// Enable foreign keys (may fail when the driver has limited PRAGMA support)
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		conn.Close()
		if isDriverUnavailable(err) {
			return newFallbackDB(), nil
		}
		return nil, errutil.Wrap(err, "enable foreign keys")
	}

	// Create schema
	if err := initializeSchema(conn); err != nil {
		conn.Close()
		return nil, errutil.Wrap(err, "initialize schema")
	}

	return &DB{
		conn:            conn,
		versions:        make([]*domain.Version, 0),
		nextVersionID:   1,
		versionFallback: false,
	}, nil
}

func newFallbackDB() *DB {
	return &DB{
		conn:            nil,
		versions:        make([]*domain.Version, 0),
		nextVersionID:   1,
		versionFallback: true,
	}
}

func isDriverUnavailable(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	return strings.Contains(errMsg, "requires cgo") ||
		strings.Contains(errMsg, "driver not found") ||
		strings.Contains(errMsg, "no such driver") ||
		strings.Contains(errMsg, "unknown driver") ||
		strings.Contains(errMsg, "binary was built with 'CGO_ENABLED=0'")
}

func isForeignKeyViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "foreign key")
}

// Close closes the database connection
func (db *DB) shouldUseVersionFallback() bool {
	if db == nil {
		return true
	}

	db.fallbackMu.RLock()
	defer db.fallbackMu.RUnlock()
	return db.conn == nil || db.versionFallback
}

func (db *DB) enableVersionFallback() {
	if db == nil {
		return
	}

	db.fallbackMu.Lock()
	db.versionFallback = true
	db.fallbackMu.Unlock()
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
	if db == nil || db.conn == nil {
		return nil
	}
	return db.conn.Ping()
}

// initializeSchema creates all tables if they don't exist
func initializeSchema(conn *sql.DB) error {
	if _, err := conn.Exec(Schema); err != nil {
		return errutil.Wrap(err, "execute schema")
	}
	return nil
}

// getDefaultDataDir returns the default data directory
func getDefaultDataDir() string {
	homeDir, err := getHomeDir()
	if err != nil {
		return "./data" // fallback
	}
	return filepath.Join(homeDir, ".noise")
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
	InsertSongWithVersion(song *domain.Song, initialContent string) (*domain.Song, *domain.Version, error)
	GetSong(id int) (*domain.Song, error)
	GetSongByFilepath(filepath string) (*domain.Song, error)
	UpdateSong(song *domain.Song) error
	UpdateSongWithVersion(song *domain.Song, newContent string, isMilestone bool, milestoneName string) error
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
	BatchUpdateStats(statsList []*domain.WritingStats) error
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

// TransactionRepository defines the interface for transaction operations
type TransactionRepository interface {
	ExecuteInTransaction(ctx context.Context, fn func(*sql.Tx) error) error
}

// Ensure DB implements the repository interfaces
var _ SongRepository = (*DB)(nil)
var _ VersionRepository = (*DB)(nil)
var _ StatsRepository = (*DB)(nil)
var _ ProjectRepository = (*DB)(nil)
var _ TransactionRepository = (*DB)(nil)

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

// In-memory fallback methods for version management

// saveVersionInMemory saves a version in memory when database is unavailable
func (db *DB) saveVersionInMemory(songID int, content string, isMilestone bool, name string, createdAt time.Time) *domain.Version {
	if db == nil {
		return nil
	}

	db.versionMutex.Lock()
	defer db.versionMutex.Unlock()

	version := &domain.Version{
		ID:            db.nextVersionID,
		SongID:        songID,
		Content:       content,
		IsMilestone:   isMilestone,
		MilestoneName: name,
		CreatedAt:     createdAt,
	}

	db.versions = append(db.versions, version)
	db.nextVersionID++

	return version
}

// getVersionsInMemory retrieves version history for a song from memory
func (db *DB) getVersionsInMemory(songID int, limit int) []*domain.Version {
	if db == nil {
		return nil
	}

	db.versionMutex.Lock()
	defer db.versionMutex.Unlock()

	var versions []*domain.Version
	count := 0

	// Get versions in reverse order (newest first)
	for i := len(db.versions) - 1; i >= 0 && count < limit; i-- {
		if db.versions[i].SongID == songID {
			versions = append(versions, db.versions[i])
			count++
		}
	}

	return versions
}

// getVersionInMemory retrieves a specific version by ID from memory
func (db *DB) getVersionInMemory(id int) (*domain.Version, error) {
	if db == nil {
		return nil, fmt.Errorf("database is nil")
	}

	db.versionMutex.Lock()
	defer db.versionMutex.Unlock()

	for _, version := range db.versions {
		if version.ID == id {
			return version, nil
		}
	}

	return nil, fmt.Errorf("version with ID %d not found", id)
}

// deleteVersionInMemory deletes a version by ID from memory
func (db *DB) deleteVersionInMemory(id int) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}

	db.versionMutex.Lock()
	defer db.versionMutex.Unlock()

	for i, version := range db.versions {
		if version.ID == id {
			db.versions = append(db.versions[:i], db.versions[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("version with ID %d not found", id)
}
