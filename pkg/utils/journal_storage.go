package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyanite/focus/pkg/models"
)

// JournalStorage handles journal entry persistence
type JournalStorage struct {
	filePath string
}

// NewJournalStorage creates a new journal storage instance
func NewJournalStorage() *JournalStorage {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	storageDir := filepath.Join(homeDir, ".focus")

	return &JournalStorage{
		filePath: filepath.Join(storageDir, "journal.json"),
	}
}

// LoadEntries loads all journal entries from storage
func (js *JournalStorage) LoadEntries() ([]*models.JournalEntry, error) {
	// Check if file exists
	if _, err := os.Stat(js.filePath); os.IsNotExist(err) {
		// Create empty journal file
		return []*models.JournalEntry{}, js.SaveEntries([]*models.JournalEntry{})
	}

	data, err := os.ReadFile(js.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read journal file: %w", err)
	}

	var entries []*models.JournalEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse journal entries: %w", err)
	}

	return entries, nil
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

	if err := os.WriteFile(js.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write journal file: %w", err)
	}

	return nil
}

// AddEntry adds a new journal entry
func (js *JournalStorage) AddEntry(entry *models.JournalEntry) error {
	entries, err := js.LoadEntries()
	if err != nil {
		return err
	}

	// Check if entry with same ID already exists
	for i, existing := range entries {
		if existing.ID == entry.ID {
			// Update existing entry
			entries[i] = entry
			return js.SaveEntries(entries)
		}
	}

	// Add new entry
	entries = append(entries, entry)
	return js.SaveEntries(entries)
}

// GetEntryByID retrieves a journal entry by ID
func (js *JournalStorage) GetEntryByID(id string) (*models.JournalEntry, error) {
	entries, err := js.LoadEntries()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.ID == id {
			return entry, nil
		}
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

// UpdateEntry updates an existing journal entry
func (js *JournalStorage) UpdateEntry(entry *models.JournalEntry) error {
	entries, err := js.LoadEntries()
	if err != nil {
		return err
	}

	for i, existing := range entries {
		if existing.ID == entry.ID {
			entries[i] = entry
			return js.SaveEntries(entries)
		}
	}

	return fmt.Errorf("journal entry with ID %s not found", entry.ID)
}

// DeleteEntry deletes a journal entry by ID
func (js *JournalStorage) DeleteEntry(id string) error {
	entries, err := js.LoadEntries()
	if err != nil {
		return err
	}

	for i, entry := range entries {
		if entry.ID == id {
			// Remove entry from slice
			entries = append(entries[:i], entries[i+1:]...)
			return js.SaveEntries(entries)
		}
	}

	return fmt.Errorf("journal entry with ID %s not found", id)
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
