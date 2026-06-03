// Package files provides repository interfaces for file I/O operations
package files

import (
	"os"
	"time"

	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/errors"
)

// FileRepository defines the interface for file-based song persistence operations
type FileRepository interface {
	// Core CRUD operations
	SaveSong(song *domain.Song, filePath string) error
	LoadSong(filePath string) (*domain.Song, error)
	DeleteSong(filePath string) error

	// Query operations
	ListSongs() ([]string, error)
	SongExists(filePath string) (bool, error)

	// File management
	CreateBackup(filePath string) error
	RestoreBackup(filePath, backupPath string) error
	GetFileInfo(filePath string) (*FileInfo, error)
}

// FileInfo represents metadata about a song file
type FileInfo struct {
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	Hash    string    `json:"hash"`
	Exists  bool      `json:"exists"`
	IsValid bool      `json:"is_valid"`
}

// LoadSong loads a song from the given file path.
func (s *Service) LoadSong(filePath string) (*domain.Song, error) {
	return s.LoadSongFromFile(filePath)
}

// SaveSongToFile saves a song to a file (repository interface)
func (s *Service) SaveSongToFile(song *domain.Song, filePath string) error {
	return s.WriteSong(song, filePath)
}

// LoadSongFromFile loads a song from a file (repository interface)
func (s *Service) LoadSongFromFile(filePath string) (*domain.Song, error) {
	return s.ReadSong(filePath)
}

// DeleteSongFile deletes a song file (repository interface)
func (s *Service) DeleteSongFile(filePath string) error {
	return s.DeleteSong(filePath)
}

// ListSongFiles lists all song files (repository interface)
func (s *Service) ListSongFiles() ([]string, error) {
	return s.ListSongs()
}

// SongExists checks if a song file exists
func (s *Service) SongExists(filePath string) (bool, error) {
	fullPath := s.resolvePath(filePath)

	// Check if file exists in filesystem
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, errors.NewFileError("stat", fullPath, err)
	}

	return true, nil
}

// CreateBackup creates a backup of a song file
func (s *Service) CreateBackup(filePath string) error {
	fullPath := s.resolvePath(filePath)

	// Check if source file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return errors.NewFileError("backup_source_missing", fullPath, err)
	}

	// Create backup filename with timestamp
	backupPath := fullPath + ".backup." + time.Now().Format("20060102_150405")

	// Copy file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return errors.NewFileError("read_for_backup", fullPath, err)
	}

	if err := os.WriteFile(backupPath, content, 0o644); err != nil {
		return errors.NewFileError("write_backup", backupPath, err)
	}

	return nil
}

// RestoreBackup restores a song from backup
func (s *Service) RestoreBackup(filePath, backupPath string) error {
	fullPath := s.resolvePath(filePath)
	fullBackupPath := s.resolvePath(backupPath)

	// Check if backup exists
	if _, err := os.Stat(fullBackupPath); os.IsNotExist(err) {
		return errors.NewFileError("backup_missing", fullBackupPath, err)
	}

	// Copy backup content to original location
	content, err := os.ReadFile(fullBackupPath)
	if err != nil {
		return errors.NewFileError("read_backup", fullBackupPath, err)
	}

	if err := os.WriteFile(fullPath, content, 0o644); err != nil {
		return errors.NewFileError("restore_file", fullPath, err)
	}

	return nil
}

// GetFileInfo returns metadata about a song file
func (s *Service) GetFileInfo(filePath string) (*FileInfo, error) {
	fullPath := s.resolvePath(filePath)

	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return &FileInfo{
			Path:    filePath,
			Exists:  false,
			IsValid: false,
		}, nil
	}
	if err != nil {
		return nil, errors.NewFileError("stat", fullPath, err)
	}

	// Check if it's a valid song file by trying to parse it
	isValid := true
	if content, err := os.ReadFile(fullPath); err == nil {
		// Try to parse as song to validate
		if _, err := s.parseSongContent(string(content), fullPath); err != nil {
			isValid = false
		}
	} else {
		isValid = false
	}

	return &FileInfo{
		Path:    filePath,
		Size:    info.Size(),
		ModTime: info.ModTime(),
		Exists:  true,
		IsValid: isValid,
	}, nil
}
