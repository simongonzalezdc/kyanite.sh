// Package db provides database persistence and repository implementations for the noise.sh application.
// It supports both SQLite with CGO and in-memory fallback storage for environments where CGO is unavailable.
package db

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/domain"
	errutil "github.com/Kyanite/noise/internal/errutil"
	"github.com/Kyanite/noise/internal/logging"
)

// isTestMode returns true if running in test mode.
// Checks both testing.Testing() and the GO_TEST_MODE environment variable.
func isTestMode() bool {
	return testing.Testing() || os.Getenv("GO_TEST_MODE") != ""
}

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
	DataDir             string
	EnableCollaboration bool // When true, collaboration tables are created
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

	// Performance PRAGMAs for SQLite optimization (2026 best practices)
	// WAL mode: Better concurrent read performance, crash recovery
	if _, err := conn.Exec("PRAGMA journal_mode = WAL"); err != nil {
		// WAL mode failure is non-fatal, log and continue
		logging.GetDefaultLogger().Warn("Failed to enable WAL mode", "error", err)
	}

	// Synchronous NORMAL: Good balance of safety and performance with WAL
	if _, err := conn.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		logging.GetDefaultLogger().Warn("Failed to set synchronous mode", "error", err)
	}

	// Cache size: 64MB (negative value = KB)
	if _, err := conn.Exec("PRAGMA cache_size = -65536"); err != nil {
		logging.GetDefaultLogger().Warn("Failed to set cache size", "error", err)
	}

	// Busy timeout: 5 seconds for lock acquisition
	if _, err := conn.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		logging.GetDefaultLogger().Warn("Failed to set busy timeout", "error", err)
	}

	// Create schema (collaboration tables only if feature is enabled)
	if err := initializeSchema(conn, cfg.EnableCollaboration); err != nil {
		conn.Close()
		return nil, errutil.Wrap(err, "initialize schema")
	}

	logger := logging.GetDefaultLogger()

	schemaHash := sha256.Sum256([]byte(Schema))
	fingerprint := hex.EncodeToString(schemaHash[:])
	if len(fingerprint) > 12 {
		fingerprint = fingerprint[:12]
	}

	var userVersion int
	if err := conn.QueryRow("PRAGMA user_version").Scan(&userVersion); err != nil {
		logger.Warn("Database user_version unavailable", "path", dbPath, "error", err)
		userVersion = -1
	}

	// Only log migration table warnings when not in test mode to reduce noise
	if !isTestMode() {
		var migrationTableCount int
		if err := conn.QueryRow(`SELECT COUNT(1) FROM sqlite_master WHERE type='table' AND name='schema_migrations'`).Scan(&migrationTableCount); err != nil {
			logger.Warn("Failed to inspect schema_migrations table", "path", dbPath, "error", err)
		} else if migrationTableCount == 0 {
			logger.Warn("Schema migrations table missing", "path", dbPath)
		} else {
			logger.Info("Schema migrations table detected", "path", dbPath, "entries", migrationTableCount)
		}

		logger.Info("Database schema initialized", "path", dbPath, "fingerprint", fingerprint, "user_version", userVersion)
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

// Query executes a query that returns rows
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if db == nil || db.conn == nil {
		return nil, sql.ErrConnDone
	}
	return db.conn.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	if db == nil || db.conn == nil {
		return nil
	}
	return db.conn.QueryRow(query, args...)
}

// Exec executes a query that doesn't return rows
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if db == nil || db.conn == nil {
		return nil, sql.ErrConnDone
	}
	return db.conn.Exec(query, args...)
}

// initializeSchema creates all tables if they don't exist
func initializeSchema(conn *sql.DB, enableCollaboration bool) error {
	// Always execute the core schema
	if _, err := conn.Exec(Schema); err != nil {
		return errutil.Wrap(err, "execute core schema")
	}

	// Only execute collaboration schema when the feature is enabled
	if enableCollaboration {
		if _, err := conn.Exec(CollaborationSchema); err != nil {
			return errutil.Wrap(err, "execute collaboration schema")
		}
	}

	// Always execute the Muse agent schema (AI companion feature)
	if _, err := conn.Exec(MuseAgentSchema); err != nil {
		return errutil.Wrap(err, "execute muse agent schema")
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
