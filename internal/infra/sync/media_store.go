package sync

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MediaStore handles storage of synced media files
type MediaStore struct {
	basePath   string
	voicePath  string
	photosPath string
}

// NewMediaStore creates a new media store
func NewMediaStore(basePath string) (*MediaStore, error) {
	voicePath := filepath.Join(basePath, "voice")
	photosPath := filepath.Join(basePath, "photos")

	// Create directories
	if err := os.MkdirAll(voicePath, DefaultDirPermission); err != nil {
		return nil, fmt.Errorf("failed to create voice directory: %w", err)
	}
	if err := os.MkdirAll(photosPath, DefaultDirPermission); err != nil {
		return nil, fmt.Errorf("failed to create photos directory: %w", err)
	}

	return &MediaStore{
		basePath:   basePath,
		voicePath:  voicePath,
		photosPath: photosPath,
	}, nil
}

// SaveVoiceMemo saves a voice memo and returns its path
func (ms *MediaStore) SaveVoiceMemo(deviceID string, data []byte) (string, error) {
	filename := generateMediaFilename(deviceID, "webm")
	path := filepath.Join(ms.voicePath, filename)

	if err := os.WriteFile(path, data, DefaultFilePermission); err != nil {
		return "", fmt.Errorf("failed to write voice memo: %w", err)
	}

	return path, nil
}

// SavePhoto saves a photo and returns its path
func (ms *MediaStore) SavePhoto(deviceID string, data []byte) (string, error) {
	filename := generateMediaFilename(deviceID, "jpg")
	path := filepath.Join(ms.photosPath, filename)

	if err := os.WriteFile(path, data, DefaultFilePermission); err != nil {
		return "", fmt.Errorf("failed to write photo: %w", err)
	}

	return path, nil
}

// GetPath returns the full path for a media ID
func (ms *MediaStore) GetPath(mediaID string) string {
	// Try voice first
	voicePath := filepath.Join(ms.voicePath, mediaID)
	if _, err := os.Stat(voicePath); err == nil {
		return voicePath
	}

	// Try photos
	photosPath := filepath.Join(ms.photosPath, mediaID)
	if _, err := os.Stat(photosPath); err == nil {
		return photosPath
	}

	return ""
}

// Delete removes a media file
func (ms *MediaStore) Delete(mediaID string) error {
	path := ms.GetPath(mediaID)
	if path == "" {
		return fmt.Errorf("media not found: %s", mediaID)
	}

	return os.Remove(path)
}

// ListVoiceMemos returns all voice memo files
func (ms *MediaStore) ListVoiceMemos() ([]string, error) {
	return listFiles(ms.voicePath)
}

// ListPhotos returns all photo files
func (ms *MediaStore) ListPhotos() ([]string, error) {
	return listFiles(ms.photosPath)
}

// GetBasePath returns the base path for media storage
func (ms *MediaStore) GetBasePath() string {
	return ms.basePath
}

// GetVoicePath returns the voice memos directory
func (ms *MediaStore) GetVoicePath() string {
	return ms.voicePath
}

// GetPhotosPath returns the photos directory
func (ms *MediaStore) GetPhotosPath() string {
	return ms.photosPath
}

// Cleanup removes old media files older than the specified duration
// Returns count of deleted files and any errors encountered (aggregated)
func (ms *MediaStore) Cleanup(olderThan time.Duration) (int, error) {
	cutoff := time.Now().Add(-olderThan)
	count := 0
	var errs []error

	// Cleanup voice memos
	voiceFiles, err := listFiles(ms.voicePath)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to list voice files: %w", err))
	}
	for _, file := range voiceFiles {
		filePath := filepath.Join(ms.voicePath, file)
		info, err := os.Stat(filePath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to stat %s: %w", file, err))
			continue
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(filePath); err != nil {
				errs = append(errs, fmt.Errorf("failed to remove %s: %w", file, err))
			} else {
				count++
			}
		}
	}

	// Cleanup photos
	photoFiles, err := listFiles(ms.photosPath)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to list photo files: %w", err))
	}
	for _, file := range photoFiles {
		filePath := filepath.Join(ms.photosPath, file)
		info, err := os.Stat(filePath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to stat %s: %w", file, err))
			continue
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(filePath); err != nil {
				errs = append(errs, fmt.Errorf("failed to remove %s: %w", file, err))
			} else {
				count++
			}
		}
	}

	// Return aggregated error if any occurred
	if len(errs) > 0 {
		return count, fmt.Errorf("cleanup encountered %d errors, first: %w", len(errs), errs[0])
	}
	return count, nil
}

// generateMediaFilename creates a unique filename for media
func generateMediaFilename(deviceID, ext string) string {
	timestamp := time.Now().Format("20060102_150405")
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback: use timestamp nanoseconds for uniqueness if crypto/rand fails
		ns := time.Now().UnixNano()
		randomBytes[0] = byte(ns)
		randomBytes[1] = byte(ns >> 8)
		randomBytes[2] = byte(ns >> 16)
		randomBytes[3] = byte(ns >> 24)
	}

	// Truncate device ID to 8 chars or use full if shorter
	shortDeviceID := deviceID
	if len(shortDeviceID) > 8 {
		shortDeviceID = shortDeviceID[:8]
	}

	return fmt.Sprintf("%s_%s_%x.%s", timestamp, shortDeviceID, randomBytes, ext)
}

// listFiles returns all files in a directory
func listFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}
