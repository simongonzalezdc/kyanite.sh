// Package app provides integrated editor services that combine file I/O and database operations
package app

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/errors"
	errutil "github.com/Kyanite/noise/internal/errutil"
	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/infra/files"
	"github.com/Kyanite/noise/internal/logging"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// IntegratedEditorService combines file I/O and database operations for comprehensive song management
type IntegratedEditorService struct {
	// Core services
	songRepo    db.SongRepository
	versionRepo db.VersionRepository
	fileRepo    files.FileRepository

	// Configuration
	config Config
}

// Config holds configuration for the integrated editor service
type Config struct {
	AutoSave         bool
	AutoSaveInterval time.Duration
	BackupOnSave     bool
	SyncToDatabase   bool
}

// NewIntegratedEditorService creates a new integrated editor service
func NewIntegratedEditorService(
	songRepo db.SongRepository,
	versionRepo db.VersionRepository,
	fileRepo files.FileRepository,
	config Config,
) *IntegratedEditorService {
	// Set defaults
	if config.AutoSaveInterval == 0 {
		config.AutoSaveInterval = 5 * time.Minute
	}

	return &IntegratedEditorService{
		songRepo:    songRepo,
		versionRepo: versionRepo,
		fileRepo:    fileRepo,
		config:      config,
	}
}

// LoadOrCreateSong loads a song from file or creates it if it doesn't exist
func (s *IntegratedEditorService) LoadOrCreateSong(filePath string) (*domain.Song, error) {
	if filePath == "" {
		return nil, errors.NewValidationError("file path cannot be empty", nil)
	}

	// First try to load from file
	song, err := s.fileRepo.LoadSong(filePath)
	if err != nil {
		// If file doesn't exist, create a new song
		if fileErr, ok := err.(*errors.AppError); ok && fileErr.HasCategory(errors.CategoryFile) {
			song = &domain.Song{
				Filepath: filePath,
				Metadata: domain.SongMetadata{
					Title:     s.extractTitleFromPath(filePath),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Sections: []domain.Section{},
			}

			// Save the new file
			if err := s.fileRepo.SaveSong(song, filePath); err != nil {
				return nil, errutil.Wrap(err, "create new song file")
			}

			logging.GetDefaultLogger().Info("Created new song file", "path", filePath, "title", song.Metadata.Title)
			return song, nil
		}

		return nil, errutil.Wrap(err, "load song from file")
	}

	// If configured to sync with database, ensure song exists there too
	if s.config.SyncToDatabase {
		if err := s.syncSongToDatabase(song); err != nil {
			logging.GetDefaultLogger().Warn("Failed to sync song to database", "path", filePath, "error", err)
			// Don't fail the operation, just log the warning
		}
	}

	return song, nil
}

// SaveSong saves a song to both file and database (if configured)
func (s *IntegratedEditorService) SaveSong(song *domain.Song) error {
	if song == nil {
		return errors.NewValidationError("song cannot be nil", nil)
	}

	// Update timestamp
	song.Metadata.UpdatedAt = time.Now()

	// Save to file
	if err := s.fileRepo.SaveSong(song, song.Filepath); err != nil {
		return errutil.Wrap(err, "save song to file")
	}

	// Create backup if configured
	if s.config.BackupOnSave {
		if err := s.fileRepo.CreateBackup(song.Filepath); err != nil {
			logging.GetDefaultLogger().Warn("Failed to create backup", "path", song.Filepath, "error", err)
		}
	}

	// Sync to database if configured
	if s.config.SyncToDatabase {
		if err := s.syncSongToDatabase(song); err != nil {
			logging.GetDefaultLogger().Warn("Failed to sync song to database", "path", song.Filepath, "error", err)
		}
	}

	// Create version in database if configured
	if s.config.SyncToDatabase {
		versionName := fmt.Sprintf("Edit %s", time.Now().Format("2006-01-02 15:04:05"))
		content := fmt.Sprintf("Song: %s by %s", song.Metadata.Title, song.Metadata.Artist)

		if _, err := s.versionRepo.SaveVersion(song.ID, content, false, versionName); err != nil {
			logging.GetDefaultLogger().Warn("Failed to save version", "song_id", song.ID, "error", err)
		}
	}

	logging.GetDefaultLogger().Info("Song saved successfully", "path", song.Filepath, "title", song.Metadata.Title)
	return nil
}

// DeleteSong deletes a song from both file and database
func (s *IntegratedEditorService) DeleteSong(filePath string) error {
	if filePath == "" {
		return errors.NewValidationError("file path cannot be empty", nil)
	}

	// Delete from file system
	if err := s.fileRepo.DeleteSong(filePath); err != nil {
		return errutil.Wrap(err, "delete song file")
	}

	// If syncing with database, try to find and delete from database too
	if s.config.SyncToDatabase {
		// Try to find the song in database by filepath
		if song, err := s.songRepo.GetSongByFilepath(filePath); err == nil && song != nil {
			if err := s.songRepo.DeleteSong(song.ID); err != nil {
				logging.GetDefaultLogger().Warn("Failed to delete song from database", "id", song.ID, "error", err)
			}
		}
	}

	logging.GetDefaultLogger().Info("Song deleted successfully", "path", filePath)
	return nil
}

// ListSongs returns a combined list of songs from files and database
func (s *IntegratedEditorService) ListSongs() ([]*domain.Song, error) {
	var songs []*domain.Song

	// Get songs from files
	filePaths, err := s.fileRepo.ListSongs()
	if err != nil {
		logging.GetDefaultLogger().Warn("Failed to list songs from files", "error", err)
	} else {
		for _, path := range filePaths {
			if song, err := s.fileRepo.LoadSong(path); err == nil {
				songs = append(songs, song)
			}
		}
	}

	// If syncing with database, also get songs from database that aren't in files
	if s.config.SyncToDatabase {
		limit := 1000 // Reasonable limit
		dbSongs, err := s.songRepo.ListSongs(limit, 0)
		if err != nil {
			logging.GetDefaultLogger().Warn("Failed to list songs from database", "error", err)
		} else {
			for _, song := range dbSongs {
				// Check if we already have this song from files
				found := false
				for _, existingSong := range songs {
					if existingSong.ID == song.ID {
						found = true
						break
					}
				}

				if !found {
					songs = append(songs, song)
				}
			}
		}
	}

	return songs, nil
}

// SyncFromFile loads a song from file and syncs it to database
func (s *IntegratedEditorService) SyncFromFile(filePath string) error {
	if filePath == "" {
		return errors.NewValidationError("file path cannot be empty", nil)
	}

	song, err := s.fileRepo.LoadSong(filePath)
	if err != nil {
		return errutil.Wrap(err, "load song from file")
	}

	return s.syncSongToDatabase(song)
}

// SyncToFile saves a song from database to file
func (s *IntegratedEditorService) SyncToFile(songID int, filePath string) error {
	if songID <= 0 {
		return errors.NewValidationError("song ID must be positive", nil)
	}
	if filePath == "" {
		return errors.NewValidationError("file path cannot be empty", nil)
	}

	song, err := s.songRepo.GetSong(songID)
	if err != nil {
		return errutil.Wrap(err, "load song from database")
	}

	song.Filepath = filePath
	return s.fileRepo.SaveSong(song, filePath)
}

// GetSongInfo returns comprehensive information about a song
func (s *IntegratedEditorService) GetSongInfo(filePath string) (*SongInfo, error) {
	if filePath == "" {
		return nil, errors.NewValidationError("file path cannot be empty", nil)
	}

	info := &SongInfo{
		FilePath: filePath,
	}

	// Get file info
	fileInfo, err := s.fileRepo.GetFileInfo(filePath)
	if err != nil {
		return nil, errutil.Wrap(err, "get file info")
	}
	info.FileInfo = fileInfo

	// Try to load song for metadata
	if song, err := s.fileRepo.LoadSong(filePath); err == nil {
		info.Song = song
		info.HasContent = true
	}

	// If syncing with database, get database info too
	if s.config.SyncToDatabase {
		if dbSong, err := s.songRepo.GetSongByFilepath(filePath); err == nil && dbSong != nil {
			info.DatabaseSong = dbSong
			info.InDatabase = true
		}
	}

	return info, nil
}

// SongInfo represents comprehensive information about a song
type SongInfo struct {
	FilePath     string                `json:"file_path"`
	FileInfo     *files.FileInfo       `json:"file_info"`
	Song         *domain.Song          `json:"song,omitempty"`
	DatabaseSong *domain.Song          `json:"database_song,omitempty"`
	HasContent   bool                  `json:"has_content"`
	InDatabase   bool                  `json:"in_database"`
}

// syncSongToDatabase syncs a song to the database
func (s *IntegratedEditorService) syncSongToDatabase(song *domain.Song) error {
	// Check if song already exists in database
	existingSong, err := s.songRepo.GetSongByFilepath(song.Filepath)
	if err == nil && existingSong != nil {
		// Update existing song
		song.ID = existingSong.ID
		return s.songRepo.UpdateSong(song)
	}

	// Create new song in database
	newSong, err := s.songRepo.InsertSong(song)
	if err != nil {
		return errutil.Wrap(err, "insert song in database")
	}

	// Update the song with the database ID
	song.ID = newSong.ID
	return nil
}

// extractTitleFromPath extracts a title from a file path
func (s *IntegratedEditorService) extractTitleFromPath(filePath string) string {
	fileName := filepath.Base(filePath)
	// Remove extension and replace underscores/dashes with spaces
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	title = strings.ReplaceAll(title, "_", " ")
	title = strings.ReplaceAll(title, "-", " ")
	c := cases.Title(language.Und)
	return c.String(title)
}
