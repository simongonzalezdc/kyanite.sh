// Package files provides file I/O operations for the noise.sh application.
// It handles reading and writing markdown files with YAML frontmatter,
// integrating with the existing error handling and domain systems.
package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/errors"
	"gopkg.in/yaml.v3"
)

// Service provides file I/O operations for markdown files with YAML frontmatter
type Service struct {
	// Configuration
	config Config

	// Caching for file operations
	cacheMutex sync.RWMutex
	fileCache  map[string]*cachedFile

	// File watchers for real-time updates
	watchersMutex sync.RWMutex
	watchers      map[string][]FileWatcher
}

// Config holds file service configuration
type Config struct {
	BaseDir     string
	AutoSave    bool
	AutoSaveInterval time.Duration
	BackupCount int
}

// FileWatcher defines the interface for file change notifications
type FileWatcher interface {
	OnFileChanged(filePath string, event FileEvent)
	OnFileDeleted(filePath string)
}

// FileEvent represents different types of file events
type FileEvent string

const (
	EventCreated FileEvent = "created"
	EventModified FileEvent = "modified"
	EventDeleted FileEvent = "deleted"
	EventRenamed FileEvent = "renamed"
)

// cachedFile represents a cached file with metadata
type cachedFile struct {
	content   string
	metadata  domain.SongMetadata
	modTime   time.Time
	size      int64
	lastRead  time.Time
}

// New creates a new file service instance
func New(config Config) (*Service, error) {
	if config.BaseDir == "" {
		config.BaseDir = getDefaultSongsDir()
	}

	// Ensure base directory exists
	if err := os.MkdirAll(config.BaseDir, 0o755); err != nil {
		return nil, errors.NewFileError("mkdir", config.BaseDir, err)
	}

	// Set defaults
	if config.AutoSaveInterval == 0 {
		config.AutoSaveInterval = 5 * time.Minute
	}
	if config.BackupCount == 0 {
		config.BackupCount = 5
	}

	return &Service{
		config:    config,
		fileCache: make(map[string]*cachedFile),
		watchers:  make(map[string][]FileWatcher),
	}, nil
}

// getDefaultSongsDir returns the default directory for song files
func getDefaultSongsDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./songs" // fallback
	}
	return filepath.Join(homeDir, ".noise", "songs")
}

// ReadSong reads a song file and parses it into a domain.Song
func (s *Service) ReadSong(filePath string) (*domain.Song, error) {
	// Validate file path
	if filePath == "" {
		return nil, errors.NewValidationError("file path cannot be empty", nil)
	}

	// Ensure path is absolute or resolve relative to base dir
	fullPath := s.resolvePath(filePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, errors.NewFileError("stat", fullPath, err)
	}

	// Read file content with size limit to prevent memory issues
	const maxFileSize = 10 * 1024 * 1024 // 10MB limit
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return nil, errors.NewFileError("stat", fullPath, err)
	}
	
	if fileInfo.Size() > maxFileSize {
		return nil, errors.NewFileError("size_limit", fullPath,
			fmt.Errorf("file size %d exceeds maximum allowed size %d", fileInfo.Size(), maxFileSize))
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, errors.NewFileError("read", fullPath, err)
	}

	// Parse the markdown content with YAML frontmatter
	song, err := s.parseSongContent(string(content), fullPath)
	if err != nil {
		return nil, err
	}

	// Update file metadata
	if fileInfo, err := os.Stat(fullPath); err == nil {
		s.updateFileCache(fullPath, string(content), song.Metadata, fileInfo)
	}

	return song, nil
}

// WriteSong writes a song to a file with proper YAML frontmatter
func (s *Service) WriteSong(song *domain.Song, filePath string) error {
	if song == nil {
		return errors.NewValidationError("song cannot be nil", nil)
	}

	if filePath == "" {
		return errors.NewValidationError("file path cannot be empty", nil)
	}

	// Ensure path is absolute or resolve relative to base dir
	fullPath := s.resolvePath(filePath)

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return errors.NewFileError("mkdir", dir, err)
	}

	// Determine content to write. Prefer RawContent to preserve user-authored formatting.
	content := song.RawContent
	if strings.TrimSpace(content) == "" {
		var err error
		content, err = s.serializeSong(song)
		if err != nil {
			return err
		}
	} else if !strings.HasPrefix(strings.TrimSpace(content), "---") {
		// Ensure frontmatter is present when writing raw content.
		yamlBytes, err := yaml.Marshal(song.Metadata)
		if err != nil {
			return errors.NewParsingError("yaml marshal", err)
		}

		var builder strings.Builder
		builder.WriteString("---\n")
		builder.WriteString(string(yamlBytes))
		builder.WriteString("---\n\n")
		builder.WriteString(strings.TrimLeft(content, "\n"))
		content = builder.String()
	}

	// Write to file with size limit
	if len(content) > 10*1024*1024 { // 10MB limit
		return errors.NewFileError("size_limit", fullPath,
			fmt.Errorf("content size %d exceeds maximum allowed size %d", len(content), 10*1024*1024))
	}

	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return errors.NewFileError("write", fullPath, err)
	}

	// Update cache
	if fileInfo, err := os.Stat(fullPath); err == nil {
		s.updateFileCache(fullPath, content, song.Metadata, fileInfo)
	}

	// Notify watchers
	s.notifyWatchers(fullPath, EventModified)

	return nil
}

// DeleteSong removes a song file
func (s *Service) DeleteSong(filePath string) error {
	if filePath == "" {
		return errors.NewValidationError("file path cannot be empty", nil)
	}

	fullPath := s.resolvePath(filePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return errors.NewFileError("stat", fullPath, err)
	}

	// Remove from cache
	s.removeFromCache(fullPath)

	// Delete file
	if err := os.Remove(fullPath); err != nil {
		return errors.NewFileError("delete", fullPath, err)
	}

	// Notify watchers
	s.notifyWatchers(fullPath, EventDeleted)

	return nil
}

// ListSongs returns a list of all song files in the base directory
func (s *Service) ListSongs() ([]string, error) {
	var files []string

	err := filepath.Walk(s.config.BaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only include .md files
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".md") {
			// Return relative path from base directory
			relPath, err := filepath.Rel(s.config.BaseDir, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, errors.NewFileError("list", s.config.BaseDir, err)
	}

	return files, nil
}

// WatchFile adds a file watcher for the specified file
func (s *Service) WatchFile(filePath string, watcher FileWatcher) error {
	if filePath == "" {
		return errors.NewValidationError("file path cannot be empty", nil)
	}
	if watcher == nil {
		return errors.NewValidationError("watcher cannot be nil", nil)
	}

	fullPath := s.resolvePath(filePath)

	s.watchersMutex.Lock()
	defer s.watchersMutex.Unlock()

	s.watchers[fullPath] = append(s.watchers[fullPath], watcher)
	return nil
}

// UnwatchFile removes a file watcher
func (s *Service) UnwatchFile(filePath string, watcher FileWatcher) error {
	if filePath == "" {
		return errors.NewValidationError("file path cannot be empty", nil)
	}

	fullPath := s.resolvePath(filePath)

	s.watchersMutex.Lock()
	defer s.watchersMutex.Unlock()

	watchers := s.watchers[fullPath]
	for i, w := range watchers {
		if w == watcher {
			s.watchers[fullPath] = append(watchers[:i], watchers[i+1:]...)
			break
		}
	}

	return nil
}

// resolvePath resolves a file path relative to the base directory
func (s *Service) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(s.config.BaseDir, path)
}

// updateFileCache updates the file cache with new content and metadata
func (s *Service) updateFileCache(fullPath, content string, metadata domain.SongMetadata, info os.FileInfo) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.fileCache[fullPath] = &cachedFile{
		content:  content,
		metadata: metadata,
		modTime:  info.ModTime(),
		size:     info.Size(),
		lastRead: time.Now(),
	}
}

// removeFromCache removes a file from the cache
func (s *Service) removeFromCache(fullPath string) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	delete(s.fileCache, fullPath)
}

// notifyWatchers notifies all watchers of a file event
func (s *Service) notifyWatchers(fullPath string, event FileEvent) {
	s.watchersMutex.RLock()
	watchers := s.watchers[fullPath]
	s.watchersMutex.RUnlock()

	for _, watcher := range watchers {
		// Use a goroutine to avoid blocking the file operation
		go func(w FileWatcher) {
			if w != nil {
				w.OnFileChanged(fullPath, event)
			}
		}(watcher)
	}
}

// GetCachedFile returns cached file information if available
func (s *Service) GetCachedFile(filePath string) (*cachedFile, bool) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	cached, exists := s.fileCache[s.resolvePath(filePath)]
	return cached, exists
}

// ClearCache clears the entire file cache
func (s *Service) ClearCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.fileCache = make(map[string]*cachedFile)
}

// Close gracefully shuts down the file service
func (s *Service) Close() error {
	s.ClearCache()

	s.watchersMutex.Lock()
	// Clear all watchers to prevent further notifications
	s.watchers = make(map[string][]FileWatcher)
	s.watchersMutex.Unlock()

	return nil
}

// parseSongContent parses markdown content with YAML frontmatter into a Song
func (s *Service) parseSongContent(content, fullPath string) (*domain.Song, error) {
	// Split frontmatter and body
	frontmatter, body, err := s.extractFrontmatter(content)
	if err != nil {
		return nil, errors.NewParsingError("frontmatter", err)
	}

	// Parse YAML frontmatter
	var metadata domain.SongMetadata
	if frontmatter != "" {
		if err := yaml.Unmarshal([]byte(frontmatter), &metadata); err != nil {
			return nil, errors.NewParsingError("yaml", err)
		}
	}

	// Set default timestamps if not provided
	if metadata.CreatedAt.IsZero() {
		metadata.CreatedAt = time.Now()
	}
	if metadata.UpdatedAt.IsZero() {
		metadata.UpdatedAt = time.Now()
	}

	// Parse sections from body (simplified - just split by double newlines for now)
	sections := s.parseSections(body)

	// Get relative path for the song
	relPath, err := filepath.Rel(s.config.BaseDir, fullPath)
	if err != nil {
		relPath = fullPath
	}

	return &domain.Song{
		Filepath:   relPath,
		Metadata:   metadata,
		Sections:   sections,
		RawContent: content,
	}, nil
}

// extractFrontmatter extracts YAML frontmatter from markdown content
func (s *Service) extractFrontmatter(content string) (frontmatter, body string, err error) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return "", "", nil
	}

	// Check if content starts with frontmatter delimiter (---)
	if len(lines) < 2 || !strings.HasPrefix(lines[0], "---") {
		return "", content, nil
	}

	// Find closing delimiter
	var endIndex int
	for i := 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "---") {
			endIndex = i
			break
		}
	}

	if endIndex == 0 {
		// No closing delimiter found, treat as regular content
		return "", content, nil
	}

	// Extract frontmatter and body
	frontmatter = strings.Join(lines[1:endIndex], "\n")
	body = strings.Join(lines[endIndex+1:], "\n")

	return frontmatter, body, nil
}

// parseSections parses markdown content into song sections (simplified implementation)
func (s *Service) parseSections(content string) []domain.Section {
	if content == "" {
		return nil
	}

	// For now, create a single section with all content
	// This is a simplified implementation - a full implementation would parse
	// different section types (verse, chorus, etc.) based on markdown headers

	lines := strings.Split(content, "\n")
	var sectionLines []domain.Line

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Simple syllable count estimation (rough approximation)
		words := strings.Fields(line)
		syllableCount := 0
		for _, word := range words {
			syllableCount += s.estimateSyllables(word)
		}

		sectionLines = append(sectionLines, domain.Line{
			Text:      line,
			Syllables: syllableCount,
		})
	}

	if len(sectionLines) == 0 {
		return nil
	}

	return []domain.Section{
		{
			Type:   domain.SectionVerse, // Default to verse
			Number: 1,
			Lines:  sectionLines,
		},
	}
}

// estimateSyllables provides a rough estimation of syllables in a word
func (s *Service) estimateSyllables(word string) int {
	word = strings.ToLower(word)
	vowels := "aeiouy"

	count := 0
	prevWasVowel := false

	for _, char := range word {
		isVowel := strings.ContainsRune(vowels, char)
		if isVowel && !prevWasVowel {
			count++
		}
		prevWasVowel = isVowel
	}

	// Handle silent e
	if strings.HasSuffix(word, "e") && len(word) > 1 {
		count--
	}

	// Every word has at least one syllable
	if count <= 0 {
		count = 1
	}

	return count
}

// serializeSong converts a Song back to markdown with YAML frontmatter
func (s *Service) serializeSong(song *domain.Song) (string, error) {
	var result strings.Builder

	// Write YAML frontmatter
	result.WriteString("---\n")

	yamlBytes, err := yaml.Marshal(song.Metadata)
	if err != nil {
		return "", errors.NewParsingError("yaml marshal", err)
	}

	result.WriteString(string(yamlBytes))
	result.WriteString("---\n\n")

	// Write sections (simplified - just write line text)
	for _, section := range song.Sections {
		// Add section header
		result.WriteString("## ")
		result.WriteString(string(section.Type))
		result.WriteString(" ")
		result.WriteString(fmt.Sprintf("%d", section.Number))
		result.WriteString("\n\n")

		// Write lines
		for _, line := range section.Lines {
			result.WriteString(line.Text)
			result.WriteString("\n")
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}
