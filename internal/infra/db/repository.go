package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/puente-labs/lyricforge/internal/domain"
)

// InsertSong inserts a new song into the database
func (db *DB) InsertSong(song *domain.Song) (*domain.Song, error) {
	tagsJSON, err := marshalStringArray(song.Metadata.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		INSERT INTO songs (filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.conn.Exec(query,
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
	if err != nil {
		return nil, fmt.Errorf("failed to insert song: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inserted song ID: %w", err)
	}

	song.ID = int(id)
	return song, nil
}

// GetSong retrieves a song by ID
func (db *DB) GetSong(id int) (*domain.Song, error) {
	query := `
		SELECT id, filepath, title, artist, key, tempo, time_signature, structure, tags, created_at, updated_at
		FROM songs WHERE id = ?`

	row := db.conn.QueryRow(query, id)

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
			return nil, fmt.Errorf("song with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get song: %w", err)
	}

	// Unmarshal tags
	song.Metadata.Tags, err = unmarshalStringArray(tagsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	// Load sections (this would be implemented with a separate sections table in a full implementation)
	song.Sections = []domain.Section{} // Placeholder

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
	query := `
		INSERT INTO versions (song_id, content, is_milestone, milestone_name, created_at)
		VALUES (?, ?, ?, ?, ?)`

	result, err := db.conn.Exec(query, songID, content, isMilestone, name, time.Now())
	if err != nil {
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
		CreatedAt:     time.Now(),
	}

	return version, nil
}

// GetVersions retrieves version history for a song
func (db *DB) GetVersions(songID int, limit int) ([]*domain.Version, error) {
	if limit <= 0 {
		limit = 50 // default limit
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

// AddSongToProject adds a song to a project
func (db *DB) AddSongToProject(projectID, songID int) error {
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

// RemoveSongFromProject removes a song from a project
func (db *DB) RemoveSongFromProject(projectID, songID int) error {
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
