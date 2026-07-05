package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kyanite/focus/pkg/models"
)

// JournalStorage handles journal entry persistence with in-memory caching
type JournalStorage struct {
	filePath   string
	cache      []*models.JournalEntry
	cacheIdx   map[string]int // ID -> index mapping for O(1) lookups
	cacheDirty bool
	cacheMutex sync.RWMutex
	cacheValid bool // Track if cache is loaded
}

// NewJournalStorage creates a new journal storage instance
func NewJournalStorage() *JournalStorage {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	storageDir := filepath.Join(homeDir, ".focus")

	js := &JournalStorage{
		filePath: filepath.Join(storageDir, "journal.json"),
		cacheIdx: make(map[string]int),
	}

	// Pre-load cache
	_ = js.ensureCached()

	return js
}

// NewJournalStorageWithPath creates a journal storage instance using a custom directory.
// This is intended for tests that should not write to the real home directory.
func NewJournalStorageWithPath(dir string) *JournalStorage {
	js := &JournalStorage{
		filePath: filepath.Join(dir, "journal.json"),
		cacheIdx: make(map[string]int),
	}
	_ = js.ensureCached()
	return js
}

// ensureCached loads entries into cache if not already loaded
func (js *JournalStorage) ensureCached() error {
	js.cacheMutex.Lock()
	defer js.cacheMutex.Unlock()

	if js.cacheValid {
		return nil // Already cached
	}

	// Check if file exists
	if _, err := os.Stat(js.filePath); os.IsNotExist(err) {
		// Initialize empty cache
		js.cache = []*models.JournalEntry{}
		js.cacheIdx = make(map[string]int)
		js.cacheValid = true
		return nil
	}

	data, err := os.ReadFile(js.filePath)
	if err != nil {
		return fmt.Errorf("failed to read journal file: %w", err)
	}

	var entries []*models.JournalEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to parse journal entries: %w", err)
	}

	// Build cache and index
	js.cache = entries
	js.cacheIdx = make(map[string]int)
	for i, entry := range entries {
		js.cacheIdx[entry.ID] = i
	}
	js.cacheValid = true

	return nil
}

// LoadEntries loads all journal entries from storage (now uses cache)
func (js *JournalStorage) LoadEntries() ([]*models.JournalEntry, error) {
	if err := js.ensureCached(); err != nil {
		return nil, err
	}

	js.cacheMutex.RLock()
	defer js.cacheMutex.RUnlock()

	// Return copy of cache to prevent external modification
	result := make([]*models.JournalEntry, len(js.cache))
	copy(result, js.cache)
	return result, nil
}

// SaveEntries saves all journal entries to storage
func (js *JournalStorage) SaveEntries(entries []*models.JournalEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal journal entries: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(js.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create journal directory: %w", err)
	}

	if err := os.WriteFile(js.filePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write journal file: %w", err)
	}

	// Update cache
	js.cacheMutex.Lock()
	defer js.cacheMutex.Unlock()
	js.cache = entries
	js.cacheIdx = make(map[string]int)
	for i, entry := range entries {
		js.cacheIdx[entry.ID] = i
	}
	js.cacheValid = true

	return nil
}

// flushCacheLocked saves cache to disk (called when cache is dirty).
// Caller must hold js.cacheMutex.
func (js *JournalStorage) flushCacheLocked() error {
	if !js.cacheDirty {
		return nil
	}

	entries := make([]*models.JournalEntry, len(js.cache))
	copy(entries, js.cache)

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal journal entries: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(js.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create journal directory: %w", err)
	}

	if err := os.WriteFile(js.filePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write journal file: %w", err)
	}

	js.cacheDirty = false
	return nil
}

// AddEntry adds a new journal entry (O(1) with cache)
func (js *JournalStorage) AddEntry(entry *models.JournalEntry) error {
	if err := js.ensureCached(); err != nil {
		return err
	}

	js.cacheMutex.Lock()
	defer js.cacheMutex.Unlock()

	// Check if entry with same ID already exists (O(1) lookup)
	if idx, exists := js.cacheIdx[entry.ID]; exists {
		// Update existing entry
		js.cache[idx] = entry
	} else {
		// Add new entry
		js.cacheIdx[entry.ID] = len(js.cache)
		js.cache = append(js.cache, entry)
	}

	js.cacheDirty = true
	return js.flushCacheLocked()
}

// GetEntryByID retrieves a journal entry by ID (O(1) with cache)
func (js *JournalStorage) GetEntryByID(id string) (*models.JournalEntry, error) {
	if err := js.ensureCached(); err != nil {
		return nil, err
	}

	js.cacheMutex.RLock()
	defer js.cacheMutex.RUnlock()

	// O(1) lookup via index
	if idx, exists := js.cacheIdx[id]; exists {
		return js.cache[idx], nil
	}

	return nil, fmt.Errorf("journal entry with ID %s not found", id)
}

// GetEntryByDate retrieves a journal entry by date
func (js *JournalStorage) GetEntryByDate(date string) (*models.JournalEntry, error) {
	entries, err := js.LoadEntries()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.Date == date {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("journal entry for date %s not found", date)
}

// UpdateEntry updates an existing journal entry (O(1) with cache)
func (js *JournalStorage) UpdateEntry(entry *models.JournalEntry) error {
	if err := js.ensureCached(); err != nil {
		return err
	}

	js.cacheMutex.Lock()
	defer js.cacheMutex.Unlock()

	// O(1) lookup via index
	if idx, exists := js.cacheIdx[entry.ID]; exists {
		js.cache[idx] = entry
		js.cacheDirty = true
		return js.flushCacheLocked()
	}

	return fmt.Errorf("journal entry with ID %s not found", entry.ID)
}

// DeleteEntry deletes a journal entry by ID (O(1) with cache)
func (js *JournalStorage) DeleteEntry(id string) error {
	if err := js.ensureCached(); err != nil {
		return err
	}

	js.cacheMutex.Lock()
	defer js.cacheMutex.Unlock()

	// O(1) lookup via index
	idx, exists := js.cacheIdx[id]
	if !exists {
		return fmt.Errorf("journal entry with ID %s not found", id)
	}

	// Remove entry from slice
	js.cache = append(js.cache[:idx], js.cache[idx+1:]...)

	// Rebuild index (since indices shifted)
	delete(js.cacheIdx, id)
	for i := idx; i < len(js.cache); i++ {
		js.cacheIdx[js.cache[i].ID] = i
	}

	js.cacheDirty = true
	return js.flushCacheLocked()
}

// SearchEntries searches journal entries by keyword in content, title, or tags
func (js *JournalStorage) SearchEntries(keyword string) ([]*models.JournalEntry, error) {
	entries, err := js.LoadEntries()
	if err != nil {
		return nil, err
	}

	keyword = strings.ToLower(keyword)
	var results []*models.JournalEntry

	for _, entry := range entries {
		// Search in title
		if strings.Contains(strings.ToLower(entry.Title), keyword) {
			results = append(results, entry)
			continue
		}

		// Search in content
		if strings.Contains(strings.ToLower(entry.Content), keyword) {
			results = append(results, entry)
			continue
		}

		// Search in tags
		for _, tag := range entry.Tags {
			if strings.Contains(strings.ToLower(tag), keyword) {
				results = append(results, entry)
				break
			}
		}
	}

	return results, nil
}

// GetEntriesByTag retrieves all entries with a specific tag
func (js *JournalStorage) GetEntriesByTag(tag string) ([]*models.JournalEntry, error) {
	entries, err := js.LoadEntries()
	if err != nil {
		return nil, err
	}

	var results []*models.JournalEntry
	for _, entry := range entries {
		if entry.HasTag(tag) {
			results = append(results, entry)
		}
	}

	return results, nil
}

// GetEntriesByDateRange retrieves entries within a date range (inclusive)
func (js *JournalStorage) GetEntriesByDateRange(startDate, endDate string) ([]*models.JournalEntry, error) {
	entries, err := js.LoadEntries()
	if err != nil {
		return nil, err
	}

	var results []*models.JournalEntry
	for _, entry := range entries {
		if entry.Date >= startDate && entry.Date <= endDate {
			results = append(results, entry)
		}
	}

	return results, nil
}

// GetAllTags returns all unique tags used in journal entries
func (js *JournalStorage) GetAllTags() ([]string, error) {
	entries, err := js.LoadEntries()
	if err != nil {
		return nil, err
	}

	tagSet := make(map[string]bool)
	for _, entry := range entries {
		for _, tag := range entry.Tags {
			tagSet[tag] = true
		}
	}

	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetStats returns statistics about journal entries
func (js *JournalStorage) GetStats() (map[string]int, error) {
	entries, err := js.LoadEntries()
	if err != nil {
		return nil, err
	}

	stats := map[string]int{
		"total_entries": len(entries),
		"total_words":   0,
	}

	for _, entry := range entries {
		stats["total_words"] += entry.WordCount
	}

	return stats, nil
}
