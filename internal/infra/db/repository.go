package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/puente-labs/lyricforge/internal/domain"
	appErrors "github.com/puente-labs/lyricforge/internal/errors"
	"github.com/puente-labs/lyricforge/internal/logging"
)

// InsertSong inserts a new song into the database with comprehensive error handling
func (db *DB) InsertSong(song *domain.Song) (*domain.Song, error) {
	// Input validation
	if song == nil {
		return nil, appErrors.NewValidationError("song cannot be nil", nil)
	}

	// Validate song metadata
	validationResult := appErrors.ValidateSongMetadata(
		song.Metadata.Title,
		song.Metadata.Artist,
		song.Metadata.Key,
		song.Metadata.Tempo,
		song.Metadata.TimeSignature,
	)
	if !validationResult.IsValid() {
		return nil, appErrors.NewValidationError("invalid song metadata: "+validationResult.Error(), nil)
	}

	// Validate filepath if provided
	if song.Filepath != "" {
		fileValidation := appErrors.ValidateFilepath(song.Filepath)
		if !fileValidation.IsValid() {
			return nil, appErrors.NewValidationError("invalid filepath: "+fileValidation.Error(), nil)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Marshal tags with error handling
	tagsJSON, err := marshalStringArray(song.Metadata.Tags)
	if err != nil {
		dbErr := appErrors.NewDatabaseError("marshal_tags", err).WithOperation("InsertSong").WithComponent("repository")
		logging.GetDefaultLogger().Error("Failed to marshal song tags", "error", dbErr)
		return nil, dbErr
	}

	query := `
		INSERT INTO songs (filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Execute with retry logic for transient errors
	var result sql.Result
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err = db.conn.ExecContext(ctx, query,
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

		if err == nil {
			break
		}

		// Check if this is a retryable error
		if attempt < maxRetries && db.isRetryableError(err) {
			logging.GetDefaultLogger().Warnf("Database insert attempt %d failed, retrying: %v", attempt, err)
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		// Non-retryable error or max retries reached
		dbErr := appErrors.NewDatabaseError("insert_song", err).WithOperation("InsertSong").WithComponent("repository")
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			dbErr = appErrors.NewDatabaseError("duplicate_song", err).WithOperation("InsertSong").WithComponent("repository")
		} else if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			dbErr = appErrors.NewDatabaseError("foreign_key_violation", err).WithOperation("InsertSong").WithComponent("repository")
		}

		logging.GetDefaultLogger().Error("Failed to insert song after retries", "error", dbErr, "attempts", attempt)
		return nil, dbErr
	}

	// Get the inserted ID with error handling
	id, err := result.LastInsertId()
	if err != nil {
		dbErr := appErrors.NewDatabaseError("get_last_insert_id", err).WithOperation("InsertSong").WithComponent("repository")
		logging.GetDefaultLogger().Error("Failed to get inserted song ID", "error", dbErr)
		return nil, dbErr
	}

	song.ID = int(id)

	// Log successful insertion
	logging.GetDefaultLogger().Info("Song inserted successfully", "id", song.ID, "title", song.Metadata.Title)

	return song, nil
}

// InsertSongWithVersion inserts a new song and creates an initial version atomically
func (db *DB) InsertSongWithVersion(song *domain.Song, initialContent string) (*domain.Song, *domain.Version, error) {
	// Input validation
	if song == nil {
		return nil, nil, appErrors.NewValidationError("song cannot be nil", nil)
	}

	// Validate song metadata
	validationResult := appErrors.ValidateSongMetadata(
		song.Metadata.Title,
		song.Metadata.Artist,
		song.Metadata.Key,
		song.Metadata.Tempo,
		song.Metadata.TimeSignature,
	)
	if !validationResult.IsValid() {
		return nil, nil, appErrors.NewValidationError("invalid song metadata: "+validationResult.Error(), nil)
	}

	// Validate filepath if provided
	if song.Filepath != "" {
		fileValidation := appErrors.ValidateFilepath(song.Filepath)
		if !fileValidation.IsValid() {
			return nil, nil, appErrors.NewValidationError("invalid filepath: "+fileValidation.Error(), nil)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if tx != nil {
				db.rollbackTransaction(tx, "InsertSongWithVersion")
			}
		}
	}()

	// Marshal tags with error handling
	tagsJSON, err := marshalStringArray(song.Metadata.Tags)
	if err != nil {
		dbErr := appErrors.NewDatabaseError("marshal_tags", err).WithOperation("InsertSongWithVersion").WithComponent("repository")
		logging.GetDefaultLogger().Error("Failed to marshal song tags", "error", dbErr)
		return nil, nil, dbErr
	}

	query := `
		INSERT INTO songs (filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Execute with retry logic for transient errors
	var result sql.Result
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err = tx.ExecContext(ctx, query,
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

		if err == nil {
			break
		}

		// Check if this is a retryable error
		if attempt < maxRetries && db.isRetryableError(err) {
			logging.GetDefaultLogger().Warnf("Database insert attempt %d failed, retrying: %v", attempt, err)
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		// Non-retryable error or max retries reached
		dbErr := appErrors.NewDatabaseError("insert_song", err).WithOperation("InsertSongWithVersion").WithComponent("repository")
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			dbErr = appErrors.NewDatabaseError("duplicate_song", err).WithOperation("InsertSongWithVersion").WithComponent("repository")
		} else if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			dbErr = appErrors.NewDatabaseError("foreign_key_violation", err).WithOperation("InsertSongWithVersion").WithComponent("repository")
		}

		logging.GetDefaultLogger().Error("Failed to insert song after retries", "error", dbErr, "attempts", attempt)
		return nil, nil, dbErr
	}

	// Get the inserted ID with error handling
	id, err := result.LastInsertId()
	if err != nil {
		dbErr := appErrors.NewDatabaseError("get_last_insert_id", err).WithOperation("InsertSongWithVersion").WithComponent("repository")
		logging.GetDefaultLogger().Error("Failed to get inserted song ID", "error", dbErr)
		return nil, nil, dbErr
	}

	song.ID = int(id)

	// Create initial version within the same transaction
	version := &domain.Version{
		SongID:        song.ID,
		Content:       initialContent,
		IsMilestone:   false,
		MilestoneName: "",
		CreatedAt:     time.Now(),
	}

	versionQuery := `
		INSERT INTO versions (song_id, content, is_milestone, milestone_name, created_at)
		VALUES (?, ?, ?, ?, ?)`

	versionResult, err := tx.ExecContext(ctx, versionQuery,
		version.SongID,
		version.Content,
		version.IsMilestone,
		version.MilestoneName,
		version.CreatedAt,
	)
	if err != nil {
		dbErr := appErrors.NewDatabaseError("insert_version", err).WithOperation("InsertSongWithVersion").WithComponent("repository")
		logging.GetDefaultLogger().Error("Failed to insert initial version", "error", dbErr)
		return nil, nil, dbErr
	}

	versionID, err := versionResult.LastInsertId()
	if err != nil {
		dbErr := appErrors.NewDatabaseError("get_version_insert_id", err).WithOperation("InsertSongWithVersion").WithComponent("repository")
		logging.GetDefaultLogger().Error("Failed to get inserted version ID", "error", dbErr)
		return nil, nil, dbErr
	}

	version.ID = int(versionID)

	// Commit transaction
	if err = db.commitTransaction(tx, "InsertSongWithVersion"); err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log successful insertion
	logging.GetDefaultLogger().Info("Song with initial version inserted successfully",
		"song_id", song.ID,
		"version_id", version.ID,
		"title", song.Metadata.Title)

	return song, version, nil
}

// GetSong retrieves a song by ID with comprehensive error handling
func (db *DB) GetSong(id int) (*domain.Song, error) {
	// Input validation
	if id <= 0 {
		return nil, appErrors.NewValidationError("song ID must be positive", nil)
	}

	// Validate database connection
	if err := db.validateConnection(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		SELECT id, filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at
		FROM songs WHERE id = ?`

	row := db.conn.QueryRowContext(ctx, query, id)

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
		dbErr := db.handleDatabaseError("GetSong", err)
		if errors.Is(err, sql.ErrNoRows) {
			logging.GetDefaultLogger().Debug("Song not found", "id", id)
			return nil, appErrors.NewDatabaseError("song_not_found", fmt.Errorf("song with ID %d not found", id)).WithOperation("GetSong").WithComponent("repository")
		}

		logging.GetDefaultLogger().Error("Failed to get song", "id", id, "error", dbErr)
		return nil, dbErr
	}

	// Unmarshal tags with error handling
	song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
	if err != nil {
		dbErr := appErrors.NewDatabaseError("unmarshal_tags", err).WithOperation("GetSong").WithComponent("repository")
		logging.GetDefaultLogger().Error("Failed to unmarshal song tags", "id", id, "error", dbErr)
		return nil, dbErr
	}

	// Load sections (this would be implemented with a separate sections table in a full implementation)
	song.Sections = []domain.Section{} // Placeholder

	// Log successful retrieval
	logging.GetDefaultLogger().Debug("Song retrieved successfully", "id", song.ID, "title", song.Metadata.Title)

	return &song, nil
}

// GetSongByFilepath retrieves a song by its file path
func (db *DB) GetSongByFilepath(filepath string) (*domain.Song, error) {
	query := `
		SELECT id, filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at
		FROM songs WHERE filepath = ?`

	row := db.conn.QueryRow(query, filepath)

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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("song with filepath %s not found", filepath)
		}
		return nil, fmt.Errorf("failed to get song by filepath: %w", err)
	}

	// Unmarshal tags
	song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return &song, nil
}

// UpdateSong updates an existing song
func (db *DB) UpdateSong(song *domain.Song) error {
	tagsJSON, err := marshalStringArray(song.Metadata.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		UPDATE songs
		SET filepath = ?, title = ?, artist = ?, key = ?, tempo = ?, time_signature = ?, structure = ?, tags = ?, updated_at = ?
		WHERE id = ?`

	_, err = db.conn.Exec(query,
		song.Filepath,
		song.Metadata.Title,
		song.Metadata.Artist,
		song.Metadata.Key,
		song.Metadata.Tempo,
		song.Metadata.TimeSignature,
		song.Metadata.Structure,
		tagsJSON,
		time.Now(),
		song.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update song: %w", err)
	}

	return nil
}

// DeleteSong deletes a song by ID
func (db *DB) DeleteSong(id int) error {
	query := `DELETE FROM songs WHERE id = ?`

	result, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("song with ID %d not found", id)
	}

	return nil
}

// ListSongs retrieves songs with pagination
func (db *DB) ListSongs(limit, offset int) ([]*domain.Song, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}

	query := `
		SELECT id, filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at
		FROM songs ORDER BY updated_at DESC LIMIT ? OFFSET ?`

	rows, err := db.conn.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list songs: %w", err)
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
			return nil, fmt.Errorf("failed to scan song row: %w", err)
		}

		// Unmarshal tags
		song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating song rows: %w", err)
	}

	return songs, nil
}

// SearchSongs searches songs by title, artist, or content
func (db *DB) SearchSongs(query string, limit int) ([]*domain.Song, error) {
	if limit <= 0 {
		limit = 20 // default limit
	}

	// Simple search implementation - in a full implementation, you'd use FTS
	searchQuery := "%" + query + "%"
	sqlQuery := `
		SELECT id, filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at
		FROM songs
		WHERE title LIKE ? OR artist LIKE ? OR tags LIKE ?
		ORDER BY updated_at DESC LIMIT ?`

	rows, err := db.conn.Query(sqlQuery, searchQuery, searchQuery, searchQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search songs: %w", err)
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
			return nil, fmt.Errorf("failed to scan song row: %w", err)
		}

		// Unmarshal tags
		song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return songs, nil
}

// SaveVersion saves a version snapshot of a song
func (db *DB) SaveVersion(songID int, content string, isMilestone bool, name string) (*domain.Version, error) {
	createdAt := time.Now()

	if db.shouldUseVersionFallback() {
		return db.saveVersionInMemory(songID, content, isMilestone, name, createdAt), nil
	}

	query := `
		INSERT INTO versions (song_id, content, is_milestone, milestone_name, created_at)
		VALUES (?, ?, ?, ?, ?)`

	result, err := db.conn.Exec(query, songID, content, isMilestone, name, createdAt)
	if err != nil {
		if isForeignKeyViolation(err) || isDriverUnavailable(err) {
			db.enableVersionFallback()
			return db.saveVersionInMemory(songID, content, isMilestone, name, createdAt), nil
		}
		return nil, fmt.Errorf("failed to save version: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inserted version ID: %w", err)
	}

	version := &domain.Version{
		ID:            int(id),
		SongID:        songID,
		Content:       content,
		IsMilestone:   isMilestone,
		MilestoneName: name,
		CreatedAt:     createdAt,
	}

	return version, nil
}

// GetVersions retrieves version history for a song
func (db *DB) GetVersions(songID int, limit int) ([]*domain.Version, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}

	// In-memory fallback
	if db.shouldUseVersionFallback() {
		return db.getVersionsInMemory(songID, limit), nil
	}

	query := `
		SELECT id, song_id, content, is_milestone, milestone_name, created_at
		FROM versions WHERE song_id = ? ORDER BY created_at DESC LIMIT ?`

	rows, err := db.conn.Query(query, songID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}
	defer rows.Close()

	var versions []*domain.Version
	for rows.Next() {
		var version domain.Version

		err := rows.Scan(
			&version.ID,
			&version.SongID,
			&version.Content,
			&version.IsMilestone,
			&version.MilestoneName,
			&version.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan version row: %w", err)
		}

		versions = append(versions, &version)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating version rows: %w", err)
	}

	return versions, nil
}

// GetVersion retrieves a specific version by ID
func (db *DB) GetVersion(id int) (*domain.Version, error) {
	// In-memory fallback
	if db.shouldUseVersionFallback() {
		return db.getVersionInMemory(id)
	}

	query := `
		SELECT id, song_id, content, is_milestone, milestone_name, created_at
		FROM versions WHERE id = ?`

	row := db.conn.QueryRow(query, id)

	var version domain.Version
	err := row.Scan(
		&version.ID,
		&version.SongID,
		&version.Content,
		&version.IsMilestone,
		&version.MilestoneName,
		&version.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("version with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	return &version, nil
}

// DeleteVersion deletes a version by ID
func (db *DB) DeleteVersion(id int) error {
	// In-memory fallback
	if db.shouldUseVersionFallback() {
		return db.deleteVersionInMemory(id)
	}

	query := `DELETE FROM versions WHERE id = ?`

	result, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete version: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("version with ID %d not found", id)
	}

	return nil
}

// UpdateSongWithVersion updates a song and creates a version snapshot atomically
func (db *DB) UpdateSongWithVersion(song *domain.Song, newContent string, isMilestone bool, milestoneName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if tx != nil {
				db.rollbackTransaction(tx, "UpdateSongWithVersion")
			}
		}
	}()

	// Update song within transaction
	tagsJSON, err := marshalStringArray(song.Metadata.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		UPDATE songs
		SET filepath = ?, title = ?, artist = ?, key = ?, tempo = ?, time_signature = ?, structure = ?, tags = ?, updated_at = ?
		WHERE id = ?`

	_, err = tx.ExecContext(ctx, query,
		song.Filepath,
		song.Metadata.Title,
		song.Metadata.Artist,
		song.Metadata.Key,
		song.Metadata.Tempo,
		song.Metadata.TimeSignature,
		song.Metadata.Structure,
		tagsJSON,
		time.Now(),
		song.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update song: %w", err)
	}

	// Create version snapshot within the same transaction
	version := &domain.Version{
		SongID:        song.ID,
		Content:       newContent,
		IsMilestone:   isMilestone,
		MilestoneName: milestoneName,
		CreatedAt:     time.Now(),
	}

	versionQuery := `
		INSERT INTO versions (song_id, content, is_milestone, milestone_name, created_at)
		VALUES (?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, versionQuery,
		version.SongID,
		version.Content,
		version.IsMilestone,
		version.MilestoneName,
		version.CreatedAt,
	)
	if err != nil {
		dbErr := appErrors.NewDatabaseError("insert_version", err).WithOperation("UpdateSongWithVersion").WithComponent("repository")
		logging.GetDefaultLogger().Error("Failed to insert version", "error", dbErr)
		return fmt.Errorf("failed to save version: %w", err)
	}

	// Commit transaction
	if err = db.commitTransaction(tx, "UpdateSongWithVersion"); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logging.GetDefaultLogger().Info("Song updated with version successfully",
		"song_id", song.ID,
		"title", song.Metadata.Title,
		"is_milestone", isMilestone,
		"milestone_name", milestoneName)

	return nil
}

// BatchUpdateStats atomically updates multiple writing statistics records
func (db *DB) BatchUpdateStats(statsList []*domain.WritingStats) error {
	if len(statsList) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if tx != nil {
				db.rollbackTransaction(tx, "BatchUpdateStats")
			}
		}
	}()

	// Use UPSERT pattern for each stats record within the transaction
	query := `
		INSERT OR REPLACE INTO writing_stats (date, words_written, songs_created, songs_edited, ai_requests, time_spent_minutes)
		VALUES (?, ?, ?, ?, ?, ?)`

	for _, stats := range statsList {
		_, err = tx.ExecContext(ctx, query,
			stats.Date,
			stats.WordsWritten,
			stats.SongsCreated,
			stats.SongsEdited,
			stats.AIRequests,
			stats.TimeSpentMinutes,
		)
		if err != nil {
			return fmt.Errorf("failed to update stats for date %s: %w", stats.Date.Format("2006-01-02"), err)
		}
	}

	// Commit transaction
	if err = db.commitTransaction(tx, "BatchUpdateStats"); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logging.GetDefaultLogger().Info("Batch stats update completed successfully", "count", len(statsList))
	return nil
}

// ExecuteInTransaction executes a function within a database transaction
// This provides a flexible way to perform multiple operations atomically
func (db *DB) ExecuteInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if tx != nil {
				db.rollbackTransaction(tx, "ExecuteInTransaction")
			}
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = db.commitTransaction(tx, "ExecuteInTransaction"); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RecordStats records or updates writing statistics for a date
func (db *DB) RecordStats(stats *domain.WritingStats) error {
	// Use UPSERT pattern (INSERT OR REPLACE for SQLite)
	query := `
		INSERT OR REPLACE INTO writing_stats (date, words_written, songs_created, songs_edited, ai_requests, time_spent_minutes)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.conn.Exec(query,
		stats.Date,
		stats.WordsWritten,
		stats.SongsCreated,
		stats.SongsEdited,
		stats.AIRequests,
		stats.TimeSpentMinutes,
	)
	if err != nil {
		return fmt.Errorf("failed to record stats: %w", err)
	}

	return nil
}

// GetStats retrieves writing statistics for a specific date
func (db *DB) GetStats(date time.Time) (*domain.WritingStats, error) {
	query := `SELECT id, date, words_written, songs_created, songs_edited, ai_requests, time_spent_minutes FROM writing_stats WHERE date = ?`

	row := db.conn.QueryRow(query, date)

	var stats domain.WritingStats
	err := row.Scan(
		&stats.ID,
		&stats.Date,
		&stats.WordsWritten,
		&stats.SongsCreated,
		&stats.SongsEdited,
		&stats.AIRequests,
		&stats.TimeSpentMinutes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stats for date %s not found", date.Format("2006-01-02"))
		}
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return &stats, nil
}

// GetStatsRange retrieves writing statistics for a date range
func (db *DB) GetStatsRange(start, end time.Time) ([]*domain.WritingStats, error) {
	query := `
		SELECT id, date, words_written, songs_created, songs_edited, ai_requests, time_spent_minutes
		FROM writing_stats WHERE date BETWEEN ? AND ? ORDER BY date`

	rows, err := db.conn.Query(query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats range: %w", err)
	}
	defer rows.Close()

	var stats []*domain.WritingStats
	for rows.Next() {
		var stat domain.WritingStats

		err := rows.Scan(
			&stat.ID,
			&stat.Date,
			&stat.WordsWritten,
			&stat.SongsCreated,
			&stat.SongsEdited,
			&stat.AIRequests,
			&stat.TimeSpentMinutes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stats row: %w", err)
		}

		stats = append(stats, &stat)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stats rows: %w", err)
	}

	return stats, nil
}

// UpdateStats updates existing writing statistics
func (db *DB) UpdateStats(stats *domain.WritingStats) error {
	query := `
		UPDATE writing_stats
		SET words_written = ?, songs_created = ?, songs_edited = ?, ai_requests = ?, time_spent_minutes = ?
		WHERE id = ?`

	_, err := db.conn.Exec(query,
		stats.WordsWritten,
		stats.SongsCreated,
		stats.SongsEdited,
		stats.AIRequests,
		stats.TimeSpentMinutes,
		stats.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update stats: %w", err)
	}

	return nil
}

// CreateProject creates a new project
func (db *DB) CreateProject(project *domain.Project) (*domain.Project, error) {
	songIDsJSON, err := marshalIntArray(project.SongIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal song IDs: %w", err)
	}

	query := `
		INSERT INTO projects (name, description, song_ids, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`

	result, err := db.conn.Exec(query,
		project.Name,
		project.Description,
		songIDsJSON,
		project.CreatedAt,
		project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inserted project ID: %w", err)
	}

	project.ID = int(id)
	return project, nil
}

// GetProject retrieves a project by ID
func (db *DB) GetProject(id int) (*domain.Project, error) {
	query := `
		SELECT id, name, description, song_ids, created_at, updated_at
		FROM projects WHERE id = ?`

	row := db.conn.QueryRow(query, id)

	var project domain.Project
	var songIDsJSON string

	err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&songIDsJSON,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Unmarshal song IDs
	project.SongIDs, err = unmarshalIntArray(songIDsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal song IDs: %w", err)
	}

	return &project, nil
}

// UpdateProject updates an existing project
func (db *DB) UpdateProject(project *domain.Project) error {
	songIDsJSON, err := marshalIntArray(project.SongIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal song IDs: %w", err)
	}

	query := `
		UPDATE projects
		SET name = ?, description = ?, song_ids = ?, updated_at = ?
		WHERE id = ?`

	_, err = db.conn.Exec(query,
		project.Name,
		project.Description,
		songIDsJSON,
		time.Now(),
		project.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

// DeleteProject deletes a project by ID
func (db *DB) DeleteProject(id int) error {
	query := `DELETE FROM projects WHERE id = ?`

	result, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project with ID %d not found", id)
	}

	return nil
}

// ListProjects retrieves all projects
func (db *DB) ListProjects() ([]*domain.Project, error) {
	query := `
		SELECT id, name, description, song_ids, created_at, updated_at
		FROM projects ORDER BY updated_at DESC`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var project domain.Project
		var songIDsJSON string

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&songIDsJSON,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}

		// Unmarshal song IDs
		project.SongIDs, err = unmarshalIntArray(songIDsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal song IDs: %w", err)
		}

		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project rows: %w", err)
	}

	return projects, nil
}

// AddSongToProject adds a song to a project (transactional)
func (db *DB) AddSongToProject(projectID, songID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if tx != nil {
				db.rollbackTransaction(tx, "AddSongToProject")
			}
		}
	}()

	// Get current project within transaction
	project, err := db.getProjectInTx(tx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Check if song is already in project
	for _, id := range project.SongIDs {
		if id == songID {
			// Commit the transaction since no changes are needed
			if commitErr := db.commitTransaction(tx, "AddSongToProject"); commitErr != nil {
				return fmt.Errorf("failed to commit transaction: %w", commitErr)
			}
			return nil // Already in project
		}
	}

	// Add song ID to project
	project.SongIDs = append(project.SongIDs, songID)

	// Update project within transaction
	if err = db.updateProjectInTx(tx, project); err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	// Commit transaction
	if err = db.commitTransaction(tx, "AddSongToProject"); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logging.GetDefaultLogger().Info("Song added to project successfully", "project_id", projectID, "song_id", songID)
	return nil
}

// AddSongToProjectNonTx adds a song to a project (non-transactional for compatibility)
func (db *DB) AddSongToProjectNonTx(projectID, songID int) error {
	// Get current project
	project, err := db.GetProject(projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Check if song is already in project
	for _, id := range project.SongIDs {
		if id == songID {
			return nil // Already in project
		}
	}

	// Add song ID to project
	project.SongIDs = append(project.SongIDs, songID)

	// Update project
	return db.UpdateProject(project)
}

// isRetryableError determines if a database error is retryable
func (db *DB) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// Retryable errors
	retryablePatterns := []string{
		"timeout",
		"connection reset",
		"connection refused",
		"temporary failure",
		"server closed",
		"database is locked",
		"database locked",
		"deadlock",
		"lock wait timeout",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// validateConnection checks if the database connection is healthy
func (db *DB) validateConnection() error {
	if db.conn == nil {
		return appErrors.NewDatabaseError("connection_nil", fmt.Errorf("database connection is nil")).WithComponent("repository")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.conn.PingContext(ctx); err != nil {
		return appErrors.NewDatabaseError("ping_failed", err).WithComponent("repository")
	}

	return nil
}

// beginTransaction starts a database transaction with error handling
func (db *DB) beginTransaction(ctx context.Context) (*sql.Tx, error) {
	if err := db.validateConnection(); err != nil {
		return nil, err
	}

	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, appErrors.NewDatabaseError("begin_transaction", err).WithOperation("BeginTransaction").WithComponent("repository")
	}

	return tx, nil
}

// rollbackTransaction rolls back a transaction with error handling
func (db *DB) rollbackTransaction(tx *sql.Tx, operation string) {
	if tx == nil {
		return
	}

	if err := tx.Rollback(); err != nil {
		logging.GetDefaultLogger().Error("Failed to rollback transaction", "operation", operation, "error", err)
	}
}

// commitTransaction commits a transaction with error handling
func (db *DB) commitTransaction(tx *sql.Tx, operation string) error {
	if tx == nil {
		return appErrors.NewDatabaseError("commit_nil_tx", fmt.Errorf("transaction is nil")).WithOperation(operation).WithComponent("repository")
	}

	if err := tx.Commit(); err != nil {
		return appErrors.NewDatabaseError("commit_transaction", err).WithOperation(operation).WithComponent("repository")
	}

	return nil
}

// handleDatabaseError creates appropriate error types based on database error conditions
func (db *DB) handleDatabaseError(operation string, err error) error {
	if err == nil {
		return nil
	}

	dbErr := appErrors.NewDatabaseError(operation, err).WithComponent("repository")

	// Categorize specific database errors
	errStr := strings.ToLower(err.Error())

	switch {
	case strings.Contains(errStr, "no such table") || strings.Contains(errStr, "doesn't exist"):
		return appErrors.NewDatabaseError("table_not_found", err).WithOperation(operation).WithComponent("repository")
	case strings.Contains(errStr, "unique constraint") || strings.Contains(errStr, "duplicate entry"):
		return appErrors.NewDatabaseError("duplicate_entry", err).WithOperation(operation).WithComponent("repository")
	case strings.Contains(errStr, "foreign key constraint"):
		return appErrors.NewDatabaseError("foreign_key_violation", err).WithOperation(operation).WithComponent("repository")
	case strings.Contains(errStr, "check constraint"):
		return appErrors.NewDatabaseError("check_constraint", err).WithOperation(operation).WithComponent("repository")
	case strings.Contains(errStr, "syntax error") || strings.Contains(errStr, "near"):
		return appErrors.NewDatabaseError("sql_syntax_error", err).WithOperation(operation).WithComponent("repository")
	case err == sql.ErrNoRows:
		return appErrors.NewDatabaseError("no_rows", err).WithOperation(operation).WithComponent("repository")
	case strings.Contains(errStr, "locked"):
		return appErrors.NewDatabaseError("database_locked", err).WithOperation(operation).WithComponent("repository")
	default:
		return dbErr
	}
}

// getProjectInTx retrieves a project within a transaction
func (db *DB) getProjectInTx(tx *sql.Tx, id int) (*domain.Project, error) {
	query := `
		SELECT id, name, description, song_ids, created_at, updated_at
		FROM projects WHERE id = ?`

	row := tx.QueryRow(query, id)

	var project domain.Project
	var songIDsJSON string

	err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&songIDsJSON,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Unmarshal song IDs
	project.SongIDs, err = unmarshalIntArray(songIDsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal song IDs: %w", err)
	}

	return &project, nil
}

// updateProjectInTx updates a project within a transaction
func (db *DB) updateProjectInTx(tx *sql.Tx, project *domain.Project) error {
	songIDsJSON, err := marshalIntArray(project.SongIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal song IDs: %w", err)
	}

	query := `
		UPDATE projects
		SET name = ?, description = ?, song_ids = ?, updated_at = ?
		WHERE id = ?`

	_, err = tx.Exec(query,
		project.Name,
		project.Description,
		songIDsJSON,
		time.Now(),
		project.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

// RemoveSongFromProject removes a song from a project (transactional)
func (db *DB) RemoveSongFromProject(projectID, songID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if tx != nil {
				db.rollbackTransaction(tx, "RemoveSongFromProject")
			}
		}
	}()

	// Get current project within transaction
	project, err := db.getProjectInTx(tx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Remove song ID from project
	removed := false
	for i, id := range project.SongIDs {
		if id == songID {
			project.SongIDs = append(project.SongIDs[:i], project.SongIDs[i+1:]...)
			removed = true
			break
		}
	}

	if !removed {
		// Song wasn't in project, commit and return success
		if commitErr := db.commitTransaction(tx, "RemoveSongFromProject"); commitErr != nil {
			return fmt.Errorf("failed to commit transaction: %w", commitErr)
		}
		return nil
	}

	// Update project within transaction
	if err = db.updateProjectInTx(tx, project); err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	// Commit transaction
	if err = db.commitTransaction(tx, "RemoveSongFromProject"); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logging.GetDefaultLogger().Info("Song removed from project successfully", "project_id", projectID, "song_id", songID)
	return nil
}

// RemoveSongFromProjectNonTx removes a song from a project (non-transactional for compatibility)
func (db *DB) RemoveSongFromProjectNonTx(projectID, songID int) error {
	// Get current project
	project, err := db.GetProject(projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Remove song ID from project
	for i, id := range project.SongIDs {
		if id == songID {
			project.SongIDs = append(project.SongIDs[:i], project.SongIDs[i+1:]...)
			break
		}
	}

	// Update project
	return db.UpdateProject(project)
}
