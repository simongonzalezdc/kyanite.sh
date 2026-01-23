package app

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Kyanite/noise/internal/constants"
	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/errors"
	errutil "github.com/Kyanite/noise/internal/errutil"
	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/logging"
	"gopkg.in/yaml.v3"
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
			return nil, errutil.Wrap(err, "save song to database")
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

	ctx, cancel := context.WithTimeout(context.Background(), constants.SaveOperationTimeout)
	defer cancel()

	// Create backup before saving (safety measure)
	originalSong := *song

	// Update timestamp
	song.Metadata.UpdatedAt = time.Now()

	// Attempt to save with retry logic
	var err error
	retryTimer := time.NewTimer(constants.SaveRetryDelayBase)
	defer retryTimer.Stop()

	for attempt := 1; attempt <= constants.MaxSaveRetries; attempt++ {
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
			if attempt < constants.MaxSaveRetries && dbErr.CanRecover(errors.RecoveryRetry) {
				wait := time.Duration(attempt) * constants.SaveRetryDelayBase
				logging.GetDefaultLogger().Info("Retrying save operation", "id", song.ID, "attempt", attempt+1, "wait", wait)

				retryTimer.Reset(wait)
				select {
				case <-ctx.Done():
					err = errors.NewDatabaseError("save_timeout", ctx.Err()).WithOperation("SaveSong").WithComponent("editor_service")
					goto handleError
				case <-retryTimer.C:
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

	autoSaveTimeout := 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), autoSaveTimeout)
	defer cancel()

	// Create auto-save version with error handling
	versionName := fmt.Sprintf("Auto-save %s", time.Now().Format("2006-01-02 15:04:05"))

	// Attempt to save version with retry logic
	var err error
	retryTimer := time.NewTimer(constants.AutoSaveRetryDelay)
	defer retryTimer.Stop()

	for attempt := 1; attempt <= constants.MaxAutoSaveRetries; attempt++ {
		// Serialize full song content (YAML frontmatter + body)
		content, err := serializeSongToMarkdown(song)
		if err != nil {
			logging.GetDefaultLogger().Error("Failed to serialize song for auto-save", "song_id", song.ID, "error", err)
			return err
		}

		_, err = s.versionRepo.SaveVersion(song.ID, content, false, versionName)
		if err == nil {
			logging.GetDefaultLogger().Debug("Auto-save completed successfully", "song_id", song.ID, "title", song.Metadata.Title, "attempt", attempt)
			return nil
		}

		// Handle database errors
		if dbErr, ok := err.(*errors.AppError); ok {
			logging.GetDefaultLogger().Warn("Auto-save attempt failed", "song_id", song.ID, "attempt", attempt, "error", dbErr)

			// Check if retryable
			if attempt < constants.MaxAutoSaveRetries && dbErr.CanRecover(errors.RecoveryRetry) {
				wait := time.Duration(attempt) * constants.AutoSaveRetryDelay
				logging.GetDefaultLogger().Debug("Retrying auto-save", "song_id", song.ID, "attempt", attempt+1, "wait", wait)

				retryTimer.Reset(wait)
				select {
				case <-ctx.Done():
					goto handleError
				case <-retryTimer.C:
					// Continue to next retry attempt
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
	// Input validation
	if song == nil {
		return errors.NewValidationError("song cannot be nil for milestone", nil)
	}
	if song.ID <= 0 {
		return errors.NewValidationError("song must have valid ID for milestone", nil)
	}

	// Default milestone name
	if strings.TrimSpace(name) == "" {
		autoSaveTimestampFormat := "2006-01-02 15:04:05"
		name = fmt.Sprintf("Milestone %s", time.Now().Format(autoSaveTimestampFormat))
	}

	// Serialize full song content
	content, serr := serializeSongToMarkdown(song)
	if serr != nil {
		logging.GetDefaultLogger().Error("Failed to serialize song for milestone", "song_id", song.ID, "error", serr)
		return serr
	}

	_, err := s.versionRepo.SaveVersion(song.ID, content, true, name)
	if err != nil {
		logging.GetDefaultLogger().Error("Failed to save milestone version", "song_id", song.ID, "error", err)
		return err
	}

	logging.GetDefaultLogger().Info("Milestone created successfully", "song_id", song.ID, "name", name)
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

	// Parse the version content (expected to be full markdown with YAML frontmatter)
	restored, perr := parseMarkdownToSong(version.Content)
	if perr != nil {
		logging.GetDefaultLogger().Error("Failed to parse version content during restore", "version_id", versionID, "error", perr)
		return nil, perr
	}

	// Assign canonical fields
	restored.ID = songID
	restored.Filepath = ""

	// Ensure timestamps: prefer metadata but fall back to version created time
	if restored.Metadata.CreatedAt.IsZero() {
		restored.Metadata.CreatedAt = version.CreatedAt
	}
	if restored.Metadata.UpdatedAt.IsZero() {
		restored.Metadata.UpdatedAt = version.CreatedAt
	}

	restored.RawContent = version.Content

	// Validate restored song
	if restored.ID != songID {
		err := errors.NewValidationError("restored song ID mismatch", nil).WithOperation("RestoreVersion").WithComponent("editor_service")
		logging.GetDefaultLogger().Error("Restored song ID mismatch", "expected", songID, "actual", restored.ID, "error", err)
		return nil, err
	}

	// Log successful restoration
	logging.GetDefaultLogger().Info("Version restored successfully", "song_id", songID, "version_id", versionID, "title", restored.Metadata.Title)

	// Show success notification
	errors.ShowGlobalSuccess("Version Restored", fmt.Sprintf("Song has been restored to version from %s.", version.CreatedAt.Format("2006-01-02 15:04:05")))

	return restored, nil
}

// Helper: serialize a domain.Song into markdown with YAML frontmatter
func serializeSongToMarkdown(song *domain.Song) (string, error) {
	if song == nil {
		return "", errors.NewValidationError("song cannot be nil", nil)
	}

	// YAML frontmatter for metadata
	yamlBytes, err := yaml.Marshal(song.Metadata)
	if err != nil {
		return "", errors.NewParsingError("yaml marshal", err)
	}

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(string(yamlBytes))
	sb.WriteString("---\n\n")

	// If RawContent exists, prefer returning that (it's already full markdown)
	if strings.TrimSpace(song.RawContent) != "" {
		// Ensure frontmatter is present; RawContent may already include frontmatter.
		// If RawContent already includes the same frontmatter, we still return RawContent to preserve user edits.
		return song.RawContent, nil
	}

	// Otherwise serialize sections to a simple markdown body
	for _, section := range song.Sections {
		sb.WriteString("## ")
		sb.WriteString(string(section.Type))
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf("%d", section.Number))
		sb.WriteString("\n\n")
		for _, line := range section.Lines {
			sb.WriteString(line.Text)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// Helper: parse markdown+YAML into domain.Song (simplified)
func parseMarkdownToSong(content string) (*domain.Song, error) {
	if strings.TrimSpace(content) == "" {
		return &domain.Song{
			Sections:   []domain.Section{},
			RawContent: content,
		}, nil
	}

	// Extract frontmatter
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return &domain.Song{RawContent: content}, nil
	}

	var front string
	var body string

	if len(lines) >= 2 && strings.HasPrefix(lines[0], "---") {
		// find closing ---
		endIndex := 0
		for i := 1; i < len(lines); i++ {
			if strings.HasPrefix(lines[i], "---") {
				endIndex = i
				break
			}
		}
		if endIndex > 0 {
			front = strings.Join(lines[1:endIndex], "\n")
			body = strings.Join(lines[endIndex+1:], "\n")
		} else {
			// no closing delimiter - treat all as body
			body = content
		}
	} else {
		body = content
	}

	var metadata domain.SongMetadata
	if strings.TrimSpace(front) != "" {
		if err := yaml.Unmarshal([]byte(front), &metadata); err != nil {
			return nil, errors.NewParsingError("yaml unmarshal", err)
		}
	}

	// Parse body into a single section (simplified)
	var sectionLines []domain.Line
	for _, ln := range strings.Split(body, "\n") {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		sectionLines = append(sectionLines, domain.Line{
			Text: ln,
		})
	}

	var sections []domain.Section
	if len(sectionLines) > 0 {
		sections = []domain.Section{
			{
				Type:   domain.SectionVerse,
				Number: 1,
				Lines:  sectionLines,
			},
		}
	} else {
		sections = []domain.Section{}
	}

	return &domain.Song{
		Metadata:   metadata,
		Sections:   sections,
		RawContent: content,
	}, nil
}
