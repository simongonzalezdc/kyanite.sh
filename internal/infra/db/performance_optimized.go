package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/constants"
	"github.com/Kyanite/noise/internal/domain"
	appErrors "github.com/Kyanite/noise/internal/errors"
	"github.com/Kyanite/noise/internal/logging"
)

// PerformanceOptimizedDB extends the base DB with performance optimizations
type PerformanceOptimizedDB struct {
	*DB

	// Connection pool management
	poolConfig PoolConfig

	// Query optimization
	preparedStatements map[string]*sql.Stmt
	stmtMutex          sync.RWMutex

	// Performance monitoring
	metrics *PerformanceMetrics

	// Connection health monitoring
	healthChecker *ConnectionHealthChecker
}

// PoolConfig defines database connection pool configuration
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// PerformanceMetrics tracks database performance
type PerformanceMetrics struct {
	QueryCount       int64
	TotalQueryTime   time.Duration
	SlowQueries      int64
	ConnectionErrors int64
	mutex            sync.RWMutex
}

// ConnectionHealthChecker monitors connection health
type ConnectionHealthChecker struct {
	db            *PerformanceOptimizedDB
	checkInterval time.Duration
	stopChan      chan struct{}
	healthStatus  bool
	mutex         sync.RWMutex
}

// NewPerformanceOptimizedDB creates a new performance-optimized database instance
func NewPerformanceOptimizedDB(cfg Config) (*PerformanceOptimizedDB, error) {
	// Create base database
	baseDB, err := New(cfg)
	if err != nil {
		return nil, err
	}

	// Default performance configuration
	poolConfig := PoolConfig{
		MaxOpenConns:    25,               // Increased from 10
		MaxIdleConns:    10,               // Increased from 5
		ConnMaxLifetime: 30 * time.Minute, // Reduced from 1 hour
		ConnMaxIdleTime: 5 * time.Minute,  // New: connection idle timeout
	}

	optimizedDB := &PerformanceOptimizedDB{
		DB:                 baseDB,
		poolConfig:         poolConfig,
		preparedStatements: make(map[string]*sql.Stmt),
		metrics:            &PerformanceMetrics{},
	}

	// Apply connection pool optimizations
	if err := optimizedDB.optimizeConnectionPool(); err != nil {
		return nil, fmt.Errorf("failed to optimize connection pool: %w", err)
	}

	// Prepare frequently used statements
	if err := optimizedDB.prepareCommonStatements(); err != nil {
		logging.GetDefaultLogger().Warn("Failed to prepare some statements", "error", err)
	}

	// Start health monitoring
	optimizedDB.startHealthMonitoring()

	// Create performance indexes
	if err := optimizedDB.createPerformanceIndexes(); err != nil {
		logging.GetDefaultLogger().Warn("Failed to create some performance indexes", "error", err)
	}

	logging.GetDefaultLogger().Info("Performance-optimized database initialized",
		"max_open_conns", poolConfig.MaxOpenConns,
		"max_idle_conns", poolConfig.MaxIdleConns)

	return optimizedDB, nil
}

// optimizeConnectionPool applies performance optimizations to the connection pool
func (db *PerformanceOptimizedDB) optimizeConnectionPool() error {
	if db.conn == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Apply optimized pool settings
	db.conn.SetMaxOpenConns(db.poolConfig.MaxOpenConns)
	db.conn.SetMaxIdleConns(db.poolConfig.MaxIdleConns)
	db.conn.SetConnMaxLifetime(db.poolConfig.ConnMaxLifetime)
	db.conn.SetConnMaxIdleTime(db.poolConfig.ConnMaxIdleTime)

	return nil
}

// prepareCommonStatements prepares frequently used SQL statements
func (db *PerformanceOptimizedDB) prepareCommonStatements() error {
	if db.conn == nil {
		return fmt.Errorf("database connection is nil")
	}

	statements := map[string]string{
		"get_song":      `SELECT id, filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at FROM songs WHERE id = ?`,
		"list_songs":    `SELECT id, filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at FROM songs ORDER BY updated_at DESC LIMIT ? OFFSET ?`,
		"search_songs":  `SELECT id, filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at FROM songs WHERE title LIKE ? OR artist LIKE ? OR tags LIKE ? ORDER BY updated_at DESC LIMIT ?`,
		"get_versions":  `SELECT id, song_id, content, is_milestone, milestone_name, created_at FROM versions WHERE song_id = ? ORDER BY created_at DESC LIMIT ?`,
		"get_project":   `SELECT id, name, description, song_ids, created_at, updated_at FROM projects WHERE id = ?`,
		"list_projects": `SELECT id, name, description, song_ids, created_at, updated_at FROM projects ORDER BY updated_at DESC`,
		"record_stats":  `INSERT OR REPLACE INTO writing_stats (date, words_written, songs_created, songs_edited, ai_requests, time_spent_minutes) VALUES (?, ?, ?, ?, ?, ?)`,
		"get_stats":     `SELECT id, date, words_written, songs_created, songs_edited, ai_requests, time_spent_minutes FROM writing_stats WHERE date = ?`,
	}

	db.stmtMutex.Lock()
	defer db.stmtMutex.Unlock()

	for name, query := range statements {
		stmt, err := db.conn.Prepare(query)
		if err != nil {
			logging.GetDefaultLogger().Warn("Failed to prepare statement", "name", name, "error", err)
			continue
		}
		db.preparedStatements[name] = stmt
	}

	logging.GetDefaultLogger().Info("Prepared database statements", "count", len(db.preparedStatements))
	return nil
}

// createPerformanceIndexes creates indexes to improve query performance
func (db *PerformanceOptimizedDB) createPerformanceIndexes() error {
	if db.conn == nil {
		return fmt.Errorf("database connection is nil")
	}

	indexes := []struct {
		name     string
		query    string
		critical bool
	}{
		{
			name:     "idx_songs_title",
			query:    "CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title)",
			critical: true,
		},
		{
			name:     "idx_songs_artist",
			query:    "CREATE INDEX IF NOT EXISTS idx_songs_artist ON songs(artist)",
			critical: true,
		},
		{
			name:     "idx_songs_filepath",
			query:    "CREATE INDEX IF NOT EXISTS idx_songs_filepath ON songs(filepath)",
			critical: true,
		},
		{
			name:     "idx_songs_tags",
			query:    "CREATE INDEX IF NOT EXISTS idx_songs_tags ON songs(tags)",
			critical: false,
		},
		{
			name:     "idx_songs_updated_at",
			query:    "CREATE INDEX IF NOT EXISTS idx_songs_updated_at ON songs(updated_at DESC)",
			critical: true,
		},
		{
			name:     "idx_versions_song_id_created",
			query:    "CREATE INDEX IF NOT EXISTS idx_versions_song_id_created ON versions(song_id, created_at DESC)",
			critical: true,
		},
		{
			name:     "idx_versions_song_milestone",
			query:    "CREATE INDEX IF NOT EXISTS idx_versions_song_milestone ON versions(song_id, is_milestone, created_at DESC)",
			critical: true,
		},
		{
			name:     "idx_projects_updated_at",
			query:    "CREATE INDEX IF NOT EXISTS idx_projects_updated_at ON projects(updated_at DESC)",
			critical: false,
		},
		{
			name:     "idx_writing_stats_date",
			query:    "CREATE INDEX IF NOT EXISTS idx_writing_stats_date ON writing_stats(date)",
			critical: false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimeout)
	defer cancel()

	for _, index := range indexes {
		start := time.Now()
		_, err := db.conn.ExecContext(ctx, index.query)
		duration := time.Since(start)

		if err != nil {
			if index.critical {
				return fmt.Errorf("failed to create critical index %s: %w", index.name, err)
			}
			logging.GetDefaultLogger().Warn("Failed to create index", "name", index.name, "error", err)
		} else {
			logging.GetDefaultLogger().Debug("Created index", "name", index.name, "duration", duration)
		}
	}

	return nil
}

// GetSongOptimized retrieves a song using prepared statements and performance monitoring
func (db *PerformanceOptimizedDB) GetSongOptimized(id int) (*domain.Song, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		db.recordQueryMetrics("get_song", duration)
	}()

	// Input validation
	if id <= 0 {
		return nil, appErrors.NewValidationError("song ID must be positive", nil)
	}

	// Validate database connection
	if err := db.validateConnection(); err != nil {
		db.metrics.mutex.Lock()
		db.metrics.ConnectionErrors++
		db.metrics.mutex.Unlock()
		return nil, err
	}

	db.stmtMutex.RLock()
	stmt, exists := db.preparedStatements["get_song"]
	db.stmtMutex.RUnlock()

	if !exists {
		// Fallback to regular query if statement not prepared
		return db.GetSong(id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortContextTimeout) // Reduced timeout
	defer cancel()

	row := stmt.QueryRowContext(ctx, id)

	var song domain.Song
	var tagsJSON string

	err := row.Scan(
		&song.ID,
		&song.Filepath,
		&song.Metadata.Title,
		&song.Metadata.Artist,
		&song.Metadata.Key,
		&song.Metadata.Tempo,
		&song.Metadata.TimeSignature,
		&song.Metadata.Structure,
		&tagsJSON,
		&song.Metadata.CreatedAt,
		&song.Metadata.UpdatedAt,
	)
	if err != nil {
		dbErr := db.handleDatabaseError("GetSongOptimized", err)
		if err == sql.ErrNoRows {
			logging.GetDefaultLogger().Debug("Song not found", "id", id)
			return nil, appErrors.NewDatabaseError("song_not_found", fmt.Errorf("song with ID %d not found", id)).WithOperation("GetSongOptimized").WithComponent("repository")
		}

		logging.GetDefaultLogger().Error("Failed to get song", "id", id, "error", dbErr)
		return nil, dbErr
	}

	// Unmarshal tags with error handling
	song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
	if err != nil {
		dbErr := appErrors.NewDatabaseError("unmarshal_tags", err).WithOperation("GetSongOptimized").WithComponent("repository")
		logging.GetDefaultLogger().Error("Failed to unmarshal song tags", "id", id, "error", dbErr)
		return nil, dbErr
	}

	// Load sections (placeholder)
	song.Sections = []domain.Section{}

	logging.GetDefaultLogger().Debug("Song retrieved successfully (optimized)", "id", song.ID, "title", song.Metadata.Title)

	return &song, nil
}

// ListSongsOptimized retrieves songs with pagination using optimized queries
func (db *PerformanceOptimizedDB) ListSongsOptimized(limit, offset int) ([]*domain.Song, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		db.recordQueryMetrics("list_songs", duration)
	}()

	if limit <= 0 {
		limit = 50 // default limit
	}

	// Validate database connection
	if err := db.validateConnection(); err != nil {
		db.metrics.mutex.Lock()
		db.metrics.ConnectionErrors++
		db.metrics.mutex.Unlock()
		return nil, err
	}

	db.stmtMutex.RLock()
	stmt, exists := db.preparedStatements["list_songs"]
	db.stmtMutex.RUnlock()

	if !exists {
		// Fallback to regular query
		return db.ListSongs(limit, offset)
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	rows, err := stmt.QueryContext(ctx, limit, offset)
	if err != nil {
		return nil, appErrors.NewDatabaseError("list_songs_optimized", err).WithOperation("ListSongsOptimized").WithComponent("repository")
	}
	defer rows.Close()

	var songs []*domain.Song
	for rows.Next() {
		var song domain.Song
		var tagsJSON string

		err := rows.Scan(
			&song.ID,
			&song.Filepath,
			&song.Metadata.Title,
			&song.Metadata.Artist,
			&song.Metadata.Key,
			&song.Metadata.Tempo,
			&song.Metadata.TimeSignature,
			&song.Metadata.Structure,
			&tagsJSON,
			&song.Metadata.CreatedAt,
			&song.Metadata.UpdatedAt,
		)
		if err != nil {
			return nil, appErrors.NewDatabaseError("scan_song_row", err).WithOperation("ListSongsOptimized").WithComponent("repository")
		}

		// Unmarshal tags
		song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
		if err != nil {
			return nil, appErrors.NewDatabaseError("unmarshal_tags", err).WithOperation("ListSongsOptimized").WithComponent("repository")
		}

		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		return nil, appErrors.NewDatabaseError("iterate_song_rows", err).WithOperation("ListSongsOptimized").WithComponent("repository")
	}

	return songs, nil
}

// SearchSongsOptimized searches songs using optimized full-text search
func (db *PerformanceOptimizedDB) SearchSongsOptimized(query string, limit int) ([]*domain.Song, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		db.recordQueryMetrics("search_songs", duration)
	}()

	if limit <= 0 {
		limit = 20 // default limit
	}

	// Validate database connection
	if err := db.validateConnection(); err != nil {
		db.metrics.mutex.Lock()
		db.metrics.ConnectionErrors++
		db.metrics.mutex.Unlock()
		return nil, err
	}

	db.stmtMutex.RLock()
	stmt, exists := db.preparedStatements["search_songs"]
	db.stmtMutex.RUnlock()

	if !exists {
		// Fallback to regular search
		return db.SearchSongs(query, limit)
	}

	// Optimized search query
	searchQuery := "%" + query + "%"
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	rows, err := stmt.QueryContext(ctx, searchQuery, searchQuery, searchQuery, limit)
	if err != nil {
		return nil, appErrors.NewDatabaseError("search_songs_optimized", err).WithOperation("SearchSongsOptimized").WithComponent("repository")
	}
	defer rows.Close()

	var songs []*domain.Song
	for rows.Next() {
		var song domain.Song
		var tagsJSON string

		err := rows.Scan(
			&song.ID,
			&song.Filepath,
			&song.Metadata.Title,
			&song.Metadata.Artist,
			&song.Metadata.Key,
			&song.Metadata.Tempo,
			&song.Metadata.TimeSignature,
			&song.Metadata.Structure,
			&tagsJSON,
			&song.Metadata.CreatedAt,
			&song.Metadata.UpdatedAt,
		)
		if err != nil {
			return nil, appErrors.NewDatabaseError("scan_song_row", err).WithOperation("SearchSongsOptimized").WithComponent("repository")
		}

		// Unmarshal tags
		song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
		if err != nil {
			return nil, appErrors.NewDatabaseError("unmarshal_tags", err).WithOperation("SearchSongsOptimized").WithComponent("repository")
		}

		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		return nil, appErrors.NewDatabaseError("iterate_search_results", err).WithOperation("SearchSongsOptimized").WithComponent("repository")
	}

	return songs, nil
}

// recordQueryMetrics records performance metrics for queries
func (db *PerformanceOptimizedDB) recordQueryMetrics(queryType string, duration time.Duration) {
	db.metrics.mutex.Lock()
	defer db.metrics.mutex.Unlock()

	db.metrics.QueryCount++
	db.metrics.TotalQueryTime += duration

	// Track slow queries (>100ms)
	if duration > 100*time.Millisecond {
		db.metrics.SlowQueries++
		logging.GetDefaultLogger().Warn("Slow query detected",
			"type", queryType,
			"duration", duration,
			"threshold", "100ms")
	}

	// Log performance warnings for very slow queries
	if duration > 500*time.Millisecond {
		logging.GetDefaultLogger().Error("Very slow query detected",
			"type", queryType,
			"duration", duration,
			"threshold", "500ms")
	}
}

// GetPerformanceMetrics returns current performance metrics
func (db *PerformanceOptimizedDB) GetPerformanceMetrics() *PerformanceMetrics {
	db.metrics.mutex.RLock()
	defer db.metrics.mutex.RUnlock()

	return db.metrics
}

// startHealthMonitoring starts the connection health checker
func (db *PerformanceOptimizedDB) startHealthMonitoring() {
	db.healthChecker = &ConnectionHealthChecker{
		db:            db,
		checkInterval: 30 * time.Second,
		stopChan:      make(chan struct{}),
		healthStatus:  true,
	}

	go db.healthChecker.start()
}

// IsHealthy returns the current health status of the database
func (db *PerformanceOptimizedDB) IsHealthy() bool {
	if db.healthChecker == nil {
		return true
	}

	db.healthChecker.mutex.RLock()
	defer db.healthChecker.mutex.RUnlock()

	return db.healthChecker.healthStatus
}

// Close closes the database and cleans up resources
func (db *PerformanceOptimizedDB) Close() error {
	// Stop health monitoring
	if db.healthChecker != nil {
		close(db.healthChecker.stopChan)
	}

	// Close prepared statements
	db.stmtMutex.Lock()
	for name, stmt := range db.preparedStatements {
		if err := stmt.Close(); err != nil {
			logging.GetDefaultLogger().Warn("Failed to close prepared statement", "name", name, "error", err)
		}
	}
	db.stmtMutex.Unlock()

	// Close base database
	return db.DB.Close()
}

// ConnectionHealthChecker implementation
func (hc *ConnectionHealthChecker) start() {
	ticker := time.NewTicker(hc.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hc.checkHealth()
		case <-hc.stopChan:
			return
		}
	}
}

func (hc *ConnectionHealthChecker) checkHealth() {
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortContextTimeout)
	defer cancel()

	err := hc.db.conn.PingContext(ctx)

	hc.mutex.Lock()
	hc.healthStatus = err == nil
	hc.mutex.Unlock()

	if err != nil {
		logging.GetDefaultLogger().Warn("Database health check failed", "error", err)

		// Attempt to re-establish connection
		if recoveryErr := hc.db.optimizeConnectionPool(); recoveryErr != nil {
			logging.GetDefaultLogger().Error("Failed to recover database connection", "error", recoveryErr)
		}
	}
}
