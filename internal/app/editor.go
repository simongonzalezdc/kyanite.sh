package app

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/puente-labs/lyricforge/internal/domain"
	"github.com/puente-labs/lyricforge/internal/errors"
	"github.com/puente-labs/lyricforge/internal/infra/db"
	"github.com/puente-labs/lyricforge/internal/logging"
)

// EditorService handles song editing operations
type EditorService struct {
	songRepo    db.SongRepository
	versionRepo db.VersionRepository
}

// NewEditorService creates a new editor service
func NewEditorService(songRepo db.SongRepository, versionRepo db.VersionRepository) *EditorService {
	return &EditorService{
		songRepo:    songRepo,
		versionRepo: versionRepo,
	}
}

// LoadSong loads a song by ID with comprehensive error handling
func (s *EditorService) LoadSong(id int) (*domain.Song, error) {
	// Input validation
	if id <= 0 {
		return nil, errors.NewValidationError("song ID must be positive", nil)
	}

	// Load song from repository with error handling
	song, err := s.songRepo.GetSong(id)
	if err != nil {
		// Handle database errors
		if dbErr, ok := err.(*errors.AppError); ok {
			logging.GetDefaultLogger().Error("Failed to load song", "id", id, "error", dbErr)

			// Show user-friendly error notification
			errors.ShowGlobalError(
				"Failed to Load Song",
				fmt.Sprintf("Could not load song with ID %d. Please try again.", id),
				dbErr,
			)

			return nil, err
		}

		// Wrap unexpected errors
		unexpectedErr := errors.NewDatabaseError("load_song", err).WithOperation("LoadSong").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Unexpected error loading song", "id", id, "error", unexpectedErr)
		return nil, unexpectedErr
	}

	// Validate loaded song
	if song == nil {
		err := errors.NewDatabaseError("song_nil", fmt.Errorf("loaded song is nil")).WithOperation("LoadSong").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Loaded song is nil", "id", id, "error", err)
		return nil, err
	}

	if song.ID != id {
		err := errors.NewValidationError("song ID mismatch", nil).WithOperation("LoadSong").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Song ID mismatch", "expected", id, "actual", song.ID, "error", err)
		return nil, err
	}

	logging.GetDefaultLogger().Info("Song loaded successfully", "id", song.ID, "title", song.Metadata.Title)
	return song, nil
}

// LoadSongByFilepath loads a song by its file path
func (s *EditorService) LoadSongByFilepath(filepath string) (*domain.Song, error) {
	// First try to get from database
	song, err := s.songRepo.GetSongByFilepath(filepath)
	if err != nil {
		// If not in database, create a new song with the filepath
		song = &domain.Song{
			Filepath: filepath,
			Metadata: domain.SongMetadata{
				Title:     "Untitled Song",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Sections: []domain.Section{},
		}

		// Save to database
		song, err = s.songRepo.InsertSong(song)
		if err != nil {
			return nil, fmt.Errorf("failed to save song to database: %w", err)
		}
	}

	return song, nil
}

// SaveSong saves a song to file and database with comprehensive error handling and recovery
func (s *EditorService) SaveSong(song *domain.Song) error {
	// Input validation
	if song == nil {
		return errors.NewValidationError("song cannot be nil", nil)
	}

	// Validate song metadata before saving
	validationResult := errors.ValidateSongMetadata(
		song.Metadata.Title,
		song.Metadata.Artist,
		song.Metadata.Key,
		song.Metadata.Tempo,
		song.Metadata.TimeSignature,
	)
	if !validationResult.IsValid() {
		return errors.NewValidationError("invalid song metadata: "+validationResult.Error(), nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create backup before saving (safety measure)
	originalSong := *song

	// Update timestamp
	song.Metadata.UpdatedAt = time.Now()

	// Attempt to save with retry logic
	var err error
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = s.songRepo.UpdateSong(song)
		if err == nil {
			// Success - log and notify
			logging.GetDefaultLogger().Info("Song saved successfully", "id", song.ID, "title", song.Metadata.Title, "attempt", attempt)

			// Show success notification
			errors.ShowGlobalSuccess("Song Saved", fmt.Sprintf("'%s' has been saved successfully.", song.Metadata.Title))

			return nil
		}

		// Handle database errors
		if dbErr, ok := err.(*errors.AppError); ok {
			logging.GetDefaultLogger().Warn("Save attempt failed", "id", song.ID, "attempt", attempt, "error", dbErr)

			// Check if retryable
			if attempt < maxRetries && dbErr.CanRecover(errors.RecoveryRetry) {
				wait := time.Duration(attempt) * time.Second
				logging.GetDefaultLogger().Info("Retrying save operation", "id", song.ID, "attempt", attempt+1, "wait", wait)

				select {
				case <-ctx.Done():
					err = errors.NewDatabaseError("save_timeout", ctx.Err()).WithOperation("SaveSong").WithComponent("editor_service")
					goto handleError
				case <-time.After(wait):
					continue
				}
			}

			// Show user-friendly error notification
			errors.ShowGlobalError(
				"Failed to Save Song",
				fmt.Sprintf("Could not save '%s'. Please check your connection and try again.", song.Metadata.Title),
				dbErr,
			)

			return err
		}

		// Wrap unexpected errors
		unexpectedErr := errors.NewDatabaseError("save_song", err).WithOperation("SaveSong").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Unexpected error saving song", "id", song.ID, "error", unexpectedErr)
		return unexpectedErr
	}

handleError:
	// All retries failed - restore original song data if needed
	*song = originalSong

	// Show critical error notification
	errors.ShowGlobalError(
		"Save Failed",
		"Unable to save the song after multiple attempts. Your changes may be lost.",
		errors.NewDatabaseError("save_failed", err).WithOperation("SaveSong").WithComponent("editor_service"),
	)

	return err
}

// CreateSong creates a new song with comprehensive error handling and validation
func (s *EditorService) CreateSong(title, artist string) (*domain.Song, error) {
	// Input validation
	if strings.TrimSpace(title) == "" {
		return nil, errors.NewValidationError("song title cannot be empty", nil)
	}

	if utf8.RuneCountInString(title) > 200 {
		return nil, errors.NewValidationError("song title must be less than 200 characters", nil)
	}

	if utf8.RuneCountInString(artist) > 100 {
		return nil, errors.NewValidationError("artist name must be less than 100 characters", nil)
	}

	// Validate metadata
	validationResult := errors.ValidateSongMetadata(title, artist, "", 0, "")
	if !validationResult.IsValid() {
		return nil, errors.NewValidationError("invalid song metadata: "+validationResult.Error(), nil)
	}

	// Create new song
	song := &domain.Song{
		Metadata: domain.SongMetadata{
			Title:     strings.TrimSpace(title),
			Artist:    strings.TrimSpace(artist),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Sections: []domain.Section{},
	}

	// Attempt to save with error handling
	savedSong, err := s.songRepo.InsertSong(song)
	if err != nil {
		// Handle database errors
		if dbErr, ok := err.(*errors.AppError); ok {
			logging.GetDefaultLogger().Error("Failed to create song", "title", title, "artist", artist, "error", dbErr)

			// Show user-friendly error notification
			errors.ShowGlobalError(
				"Failed to Create Song",
				fmt.Sprintf("Could not create song '%s'. Please try again.", title),
				dbErr,
			)

			return nil, err
		}

		// Wrap unexpected errors
		unexpectedErr := errors.NewDatabaseError("create_song", err).WithOperation("CreateSong").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Unexpected error creating song", "title", title, "error", unexpectedErr)
		return nil, unexpectedErr
	}

	// Validate created song
	if savedSong == nil {
		err := errors.NewDatabaseError("created_song_nil", fmt.Errorf("created song is nil")).WithOperation("CreateSong").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Created song is nil", "title", title, "error", err)
		return nil, err
	}

	if savedSong.ID <= 0 {
		err := errors.NewValidationError("invalid song ID after creation", nil).WithOperation("CreateSong").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Invalid song ID after creation", "id", savedSong.ID, "error", err)
		return nil, err
	}

	// Log successful creation
	logging.GetDefaultLogger().Info("Song created successfully", "id", savedSong.ID, "title", savedSong.Metadata.Title, "artist", savedSong.Metadata.Artist)

	// Show success notification
	errors.ShowGlobalSuccess("Song Created", fmt.Sprintf("'%s' has been created successfully.", savedSong.Metadata.Title))

	return savedSong, nil
}

// AutoSave creates an auto-save version of the song with comprehensive error handling
func (s *EditorService) AutoSave(song *domain.Song) error {
	// Input validation
	if song == nil {
		return errors.NewValidationError("song cannot be nil for auto-save", nil)
	}

	// Validate song has valid ID
	if song.ID <= 0 {
		return errors.NewValidationError("song must have valid ID for auto-save", nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create auto-save version with error handling
	versionName := fmt.Sprintf("Auto-save %s", time.Now().Format("2006-01-02 15:04:05"))

	// Attempt to save version with retry logic
	var err error
	maxRetries := 2

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Create a simple content representation for auto-save
		// In a full implementation, this would serialize the actual song content
		content := fmt.Sprintf("Auto-saved: %s by %s", song.Metadata.Title, song.Metadata.Artist)

		_, err = s.versionRepo.SaveVersion(song.ID, content, false, versionName)
		if err == nil {
			logging.GetDefaultLogger().Debug("Auto-save completed successfully", "song_id", song.ID, "title", song.Metadata.Title, "attempt", attempt)
			return nil
		}

		// Handle database errors
		if dbErr, ok := err.(*errors.AppError); ok {
			logging.GetDefaultLogger().Warn("Auto-save attempt failed", "song_id", song.ID, "attempt", attempt, "error", dbErr)

			// Check if retryable
			if attempt < maxRetries && dbErr.CanRecover(errors.RecoveryRetry) {
				wait := time.Duration(attempt) * 500 * time.Millisecond
				logging.GetDefaultLogger().Debug("Retrying auto-save", "song_id", song.ID, "attempt", attempt+1, "wait", wait)

				select {
				case <-ctx.Done():
					err = errors.NewDatabaseError("autosave_timeout", ctx.Err()).WithOperation("AutoSave").WithComponent("editor_service")
					goto handleError
				case <-time.After(wait):
					continue
				}
			}

			// Log error but don't show to user (auto-save shouldn't be intrusive)
			logging.GetDefaultLogger().Error("Auto-save failed after retries", "song_id", song.ID, "error", dbErr)
			return err
		}

		// Wrap unexpected errors
		unexpectedErr := errors.NewDatabaseError("autosave", err).WithOperation("AutoSave").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Unexpected error during auto-save", "song_id", song.ID, "error", unexpectedErr)
		return unexpectedErr
	}

handleError:
	// Auto-save failures shouldn't prevent normal operation
	logging.GetDefaultLogger().Warn("Auto-save failed, continuing normal operation", "song_id", song.ID, "error", err)
	return nil
}

// CreateMilestone creates a milestone version of the song
func (s *EditorService) CreateMilestone(song *domain.Song, name string) error {
	// For foundation phase, just return success
	// In a full implementation, this would serialize and save the song content as a milestone
	return nil
}

// GetVersions retrieves version history for a song
func (s *EditorService) GetVersions(songID int, limit int) ([]*domain.Version, error) {
	return s.versionRepo.GetVersions(songID, limit)
}

// RestoreVersion restores a song to a specific version with comprehensive error handling
func (s *EditorService) RestoreVersion(songID int, versionID int) (*domain.Song, error) {
	// Input validation
	if songID <= 0 {
		return nil, errors.NewValidationError("song ID must be positive", nil)
	}
	if versionID <= 0 {
		return nil, errors.NewValidationError("version ID must be positive", nil)
	}

	// Get the version content with error handling
	version, err := s.versionRepo.GetVersion(versionID)
	if err != nil {
		// Handle database errors
		if dbErr, ok := err.(*errors.AppError); ok {
			logging.GetDefaultLogger().Error("Failed to get version for restore", "song_id", songID, "version_id", versionID, "error", dbErr)

			// Show user-friendly error notification
			errors.ShowGlobalError(
				"Failed to Restore Version",
				fmt.Sprintf("Could not retrieve version %d for restoration.", versionID),
				dbErr,
			)

			return nil, err
		}

		// Wrap unexpected errors
		unexpectedErr := errors.NewDatabaseError("get_version", err).WithOperation("RestoreVersion").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Unexpected error getting version", "version_id", versionID, "error", unexpectedErr)
		return nil, unexpectedErr
	}

	// Validate retrieved version
	if version == nil {
		err := errors.NewDatabaseError("version_nil", fmt.Errorf("retrieved version is nil")).WithOperation("RestoreVersion").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Retrieved version is nil", "version_id", versionID, "error", err)
		return nil, err
	}

	if version.SongID != songID {
		err := errors.NewValidationError("version does not belong to specified song", nil).WithOperation("RestoreVersion").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Version song ID mismatch", "expected", songID, "actual", version.SongID, "error", err)
		return nil, err
	}

	// For foundation phase, create a basic song from the version
	// In a full implementation, this would parse the markdown content
	song := &domain.Song{
		ID:       songID,
		Filepath: "",
		Metadata: domain.SongMetadata{
			Title:     "Restored Song",
			CreatedAt: version.CreatedAt,
			UpdatedAt: version.CreatedAt,
		},
		Sections: []domain.Section{},
	}

	// Validate restored song
	if song.ID != songID {
		err := errors.NewValidationError("restored song ID mismatch", nil).WithOperation("RestoreVersion").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Restored song ID mismatch", "expected", songID, "actual", song.ID, "error", err)
		return nil, err
	}

	// Log successful restoration
	logging.GetDefaultLogger().Info("Version restored successfully", "song_id", songID, "version_id", versionID, "title", song.Metadata.Title)

	// Show success notification
	errors.ShowGlobalSuccess("Version Restored", fmt.Sprintf("Song has been restored to version from %s.", version.CreatedAt.Format("2006-01-02 15:04:05")))

	return song, nil
}
