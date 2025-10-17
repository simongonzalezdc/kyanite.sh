package app

import (
	"fmt"
	"time"

	"github.com/puente-labs/lyricforge/internal/domain"
	"github.com/puente-labs/lyricforge/internal/infra/db"
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

// LoadSong loads a song by ID
func (s *EditorService) LoadSong(id int) (*domain.Song, error) {
	return s.songRepo.GetSong(id)
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

// SaveSong saves a song to file and database
func (s *EditorService) SaveSong(song *domain.Song) error {
	// Update timestamp
	song.Metadata.UpdatedAt = time.Now()

	// Update in database
	if err := s.songRepo.UpdateSong(song); err != nil {
		return fmt.Errorf("failed to update song in database: %w", err)
	}

	return nil
}

// CreateSong creates a new song
func (s *EditorService) CreateSong(title, artist string) (*domain.Song, error) {
	song := &domain.Song{
		Metadata: domain.SongMetadata{
			Title:     title,
			Artist:    artist,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Sections: []domain.Section{},
	}

	// Save to database
	savedSong, err := s.songRepo.InsertSong(song)
	if err != nil {
		return nil, fmt.Errorf("failed to save new song: %w", err)
	}

	return savedSong, nil
}

// AutoSave creates an auto-save version of the song
func (s *EditorService) AutoSave(song *domain.Song) error {
	// For foundation phase, just return success
	// In a full implementation, this would serialize and save the song content
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

// RestoreVersion restores a song to a specific version
func (s *EditorService) RestoreVersion(songID int, versionID int) (*domain.Song, error) {
	// Get the version content
	version, err := s.versionRepo.GetVersion(versionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
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

	return song, nil
}

