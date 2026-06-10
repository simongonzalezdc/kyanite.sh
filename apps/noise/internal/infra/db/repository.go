package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/constants"
	"github.com/Kyanite/noise/internal/domain"
	appErrors "github.com/Kyanite/noise/internal/errors"
	errutil "github.com/Kyanite/noise/internal/errutil"
	"github.com/Kyanite/noise/internal/logging"
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

	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimeout)
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
	for attempt := 1; attempt <= constants.MaxDBRetries; attempt++ {
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
		if attempt < constants.MaxDBRetries && db.isRetryableError(err) {
			logging.GetDefaultLogger().Warnf("Database insert attempt %d failed, retrying: %v", attempt, err)
			time.Sleep(time.Duration(attempt) * constants.DefaultRetryDelay)
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

	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimeout)
	defer cancel()

	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return nil, nil, errutil.Wrap(err, "begin transaction")
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
	for attempt := 1; attempt <= constants.MaxDBRetries; attempt++ {
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
		if attempt < constants.MaxDBRetries && db.isRetryableError(err) {
			logging.GetDefaultLogger().Warnf("Database insert attempt %d failed, retrying: %v", attempt, err)
			time.Sleep(time.Duration(attempt) * constants.DefaultRetryDelay)
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
		return nil, nil, errutil.Wrap(err, "commit transaction")
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

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
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
	if err := db.validateConnection(); err != nil {
		return nil, err
	}

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
		return nil, errutil.Wrap(err, "get song by filepath")
	}

	// Unmarshal tags
	song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
	if err != nil {
		return nil, errutil.Wrap(err, "unmarshal tags")
	}

	return &song, nil
}

// UpdateSong updates an existing song
func (db *DB) UpdateSong(song *domain.Song) error {
	tagsJSON, err := marshalStringArray(song.Metadata.Tags)
	if err != nil {
		return errutil.Wrap(err, "marshal tags")
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
		return errutil.Wrap(err, "update song")
	}

	return nil
}

// DeleteSong deletes a song by ID
func (db *DB) DeleteSong(id int) error {
	query := `DELETE FROM songs WHERE id = ?`

	result, err := db.conn.Exec(query, id)
	if err != nil {
		return errutil.Wrap(err, "delete song")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errutil.Wrap(err, "get rows affected")
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
		return nil, errutil.Wrap(err, "list songs")
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
			return nil, errutil.Wrap(err, "scan song row")
		}

		// Unmarshal tags
		song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
		if err != nil {
			return nil, errutil.Wrap(err, "unmarshal tags")
		}

		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		return nil, errutil.Wrap(err, "iterate song rows")
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
		return nil, errutil.Wrap(err, "search songs")
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
			return nil, errutil.Wrap(err, "scan song row")
		}

		// Unmarshal tags
		song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
		if err != nil {
			return nil, errutil.Wrap(err, "unmarshal tags")
		}

		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		return nil, errutil.Wrap(err, "iterate search results")
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
		return nil, errutil.Wrap(err, "save version")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errutil.Wrap(err, "get inserted version ID")
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
		return nil, errutil.Wrap(err, "get versions")
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
			return nil, errutil.Wrap(err, "scan version row")
		}

		versions = append(versions, &version)
	}

	if err = rows.Err(); err != nil {
		return nil, errutil.Wrap(err, "iterate version rows")
	}

	return versions, nil
}

// GetVersion retrieves a specific version by ID
func (db *DB) GetVersion(id int) (*domain.Version, error) {
	// In-memory fallback
	if db.shouldUseVersionFallback() {
		return db.getVersionInMemory(id)
	}

	if err := db.validateConnection(); err != nil {
		return nil, err
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
		return nil, errutil.Wrap(err, "get version")
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
		return errutil.Wrap(err, "delete version")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errutil.Wrap(err, "get rows affected")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("version with ID %d not found", id)
	}

	return nil
}

// UpdateSongWithVersion updates a song and creates a version snapshot atomically
func (db *DB) UpdateSongWithVersion(song *domain.Song, newContent string, isMilestone bool, milestoneName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimeout)
	defer cancel()

	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return errutil.Wrap(err, "begin transaction")
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
		return errutil.Wrap(err, "marshal tags")
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
		return errutil.Wrap(err, "update song")
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
		return errutil.Wrap(err, "save version")
	}

	// Commit transaction
	if err = db.commitTransaction(tx, "UpdateSongWithVersion"); err != nil {
		return errutil.Wrap(err, "commit transaction")
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

	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimeout)
	defer cancel()

	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return errutil.Wrap(err, "begin transaction")
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
			return errutil.Wrapf(err, "update stats for date %s", stats.Date.Format("2006-01-02"))
		}
	}

	// Commit transaction
	if err = db.commitTransaction(tx, "BatchUpdateStats"); err != nil {
		return errutil.Wrap(err, "commit transaction")
	}

	logging.GetDefaultLogger().Info("Batch stats update completed successfully", "count", len(statsList))
	return nil
}

// ExecuteInTransaction executes a function within a database transaction
// This provides a flexible way to perform multiple operations atomically
func (db *DB) ExecuteInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return errutil.Wrap(err, "begin transaction")
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
		return errutil.Wrap(err, "commit transaction")
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
		return errutil.Wrap(err, "record stats")
	}

	return nil
}

// GetStats retrieves writing statistics for a specific date
func (db *DB) GetStats(date time.Time) (*domain.WritingStats, error) {
	if err := db.validateConnection(); err != nil {
		return nil, err
	}

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
		return nil, errutil.Wrap(err, "get stats")
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
		return nil, errutil.Wrap(err, "get stats range")
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
			return nil, errutil.Wrap(err, "scan stats row")
		}

		stats = append(stats, &stat)
	}

	if err = rows.Err(); err != nil {
		return nil, errutil.Wrap(err, "iterate stats rows")
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
		return errutil.Wrap(err, "update stats")
	}

	return nil
}

// CreateProject creates a new project
func (db *DB) CreateProject(project *domain.Project) (*domain.Project, error) {
	songIDsJSON, err := marshalIntArray(project.SongIDs)
	if err != nil {
		return nil, errutil.Wrap(err, "marshal song IDs")
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
		return nil, errutil.Wrap(err, "create project")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errutil.Wrap(err, "get inserted project ID")
	}

	project.ID = int(id)
	return project, nil
}

// GetProject retrieves a project by ID
func (db *DB) GetProject(id int) (*domain.Project, error) {
	if err := db.validateConnection(); err != nil {
		return nil, err
	}

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
		return nil, errutil.Wrap(err, "get project")
	}

	// Unmarshal song IDs
	project.SongIDs, err = unmarshalIntArray(songIDsJSON)
	if err != nil {
		return nil, errutil.Wrap(err, "unmarshal song IDs")
	}

	return &project, nil
}

// UpdateProject updates an existing project
func (db *DB) UpdateProject(project *domain.Project) error {
	songIDsJSON, err := marshalIntArray(project.SongIDs)
	if err != nil {
		return errutil.Wrap(err, "marshal song IDs")
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
		return errutil.Wrap(err, "update project")
	}

	return nil
}

// DeleteProject deletes a project by ID
func (db *DB) DeleteProject(id int) error {
	query := `DELETE FROM projects WHERE id = ?`

	result, err := db.conn.Exec(query, id)
	if err != nil {
		return errutil.Wrap(err, "delete project")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errutil.Wrap(err, "get rows affected")
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
		return nil, errutil.Wrap(err, "list projects")
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
			return nil, errutil.Wrap(err, "scan project row")
		}

		// Unmarshal song IDs
		project.SongIDs, err = unmarshalIntArray(songIDsJSON)
		if err != nil {
			return nil, errutil.Wrap(err, "unmarshal song IDs")
		}

		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		return nil, errutil.Wrap(err, "iterate project rows")
	}

	return projects, nil
}

// AddSongToProject adds a song to a project (transactional)
func (db *DB) AddSongToProject(projectID, songID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimeout)
	defer cancel()

	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return errutil.Wrap(err, "begin transaction")
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
		return errutil.Wrap(err, "get project")
	}

	// Check if song is already in project
	for _, id := range project.SongIDs {
		if id == songID {
			// Commit the transaction since no changes are needed
			if commitErr := db.commitTransaction(tx, "AddSongToProject"); commitErr != nil {
				return errutil.Wrap(commitErr, "commit transaction")
			}
			return nil // Already in project
		}
	}

	// Add song ID to project
	project.SongIDs = append(project.SongIDs, songID)

	// Update project within transaction
	if err = db.updateProjectInTx(tx, project); err != nil {
		return errutil.Wrap(err, "update project")
	}

	// Commit transaction
	if err = db.commitTransaction(tx, "AddSongToProject"); err != nil {
		return errutil.Wrap(err, "commit transaction")
	}

	logging.GetDefaultLogger().Info("Song added to project successfully", "project_id", projectID, "song_id", songID)
	return nil
}

// AddSongToProjectNonTx adds a song to a project (non-transactional for compatibility)
func (db *DB) AddSongToProjectNonTx(projectID, songID int) error {
	// Get current project
	project, err := db.GetProject(projectID)
	if err != nil {
		return errutil.Wrap(err, "get project")
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

// RemoveSongFromProject removes a song from a project (transactional)
func (db *DB) RemoveSongFromProject(projectID, songID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimeout)
	defer cancel()

	tx, err := db.beginTransaction(ctx)
	if err != nil {
		return errutil.Wrap(err, "begin transaction")
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
		return errutil.Wrap(err, "get project")
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
			return errutil.Wrap(commitErr, "commit transaction")
		}
		return nil
	}

	// Update project within transaction
	if err = db.updateProjectInTx(tx, project); err != nil {
		return errutil.Wrap(err, "update project")
	}

	// Commit transaction
	if err = db.commitTransaction(tx, "RemoveSongFromProject"); err != nil {
		return errutil.Wrap(err, "commit transaction")
	}

	logging.GetDefaultLogger().Info("Song removed from project successfully", "project_id", projectID, "song_id", songID)
	return nil
}

// RemoveSongFromProjectNonTx removes a song from a project (non-transactional for compatibility)
func (db *DB) RemoveSongFromProjectNonTx(projectID, songID int) error {
	// Get current project
	project, err := db.GetProject(projectID)
	if err != nil {
		return errutil.Wrap(err, "get project")
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
