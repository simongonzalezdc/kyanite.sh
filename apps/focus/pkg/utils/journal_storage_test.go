package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/kyanite/focus/pkg/models"
)

func TestJournalStorage_NewJournalStorage(t *testing.T) {
	storage := NewJournalStorage()

	// Test that storage creates correct file path
	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, ".focus", "journal.json")

	if storage.filePath != expectedPath {
		t.Errorf("Expected filePath to be %s, got %s", expectedPath, storage.filePath)
	}
}

func TestJournalStorage_SaveAndLoadEntries(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	// Create test entries
	entries := []*models.JournalEntry{
		{
			ID:      "test-1",
			Date:    "2025-01-01",
			Title:   "First Entry",
			Content: "First content",
			Mood:    "happy",
			Tags:    []string{"test"},
		},
		{
			ID:      "test-2",
			Date:    "2025-01-02",
			Title:   "Second Entry",
			Content: "Second content",
			Mood:    "productive",
			Tags:    []string{"work"},
		},
	}

	// Save entries
	err := storage.SaveEntries(entries)
	if err != nil {
		t.Fatalf("Failed to save entries: %v", err)
	}

	// Load entries
	loaded, err := storage.LoadEntries()
	if err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	// Verify entries
	if len(loaded) != len(entries) {
		t.Errorf("Expected %d entries, got %d", len(entries), len(loaded))
	}

	for i, entry := range entries {
		if loaded[i].ID != entry.ID {
			t.Errorf("Expected ID %s, got %s", entry.ID, loaded[i].ID)
		}
		if loaded[i].Date != entry.Date {
			t.Errorf("Expected Date %s, got %s", entry.Date, loaded[i].Date)
		}
		if loaded[i].Content != entry.Content {
			t.Errorf("Expected Content %s, got %s", entry.Content, loaded[i].Content)
		}
	}
}

func TestJournalStorage_AddEntry(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	entry := &models.JournalEntry{
		ID:      "test-add",
		Date:    "2025-01-01",
		Title:   "Test Add Entry",
		Content: "Test content",
		Mood:    "testy",
		Tags:    []string{"test"},
	}

	// Add entry
	err := storage.AddEntry(entry)
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Load and verify
	loaded, err := storage.LoadEntries()
	if err != nil {
		t.Fatalf("Failed to load entries: %v", err)
	}

	if len(loaded) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(loaded))
	}

	if loaded[0].ID != entry.ID {
		t.Errorf("Expected ID %s, got %s", entry.ID, loaded[0].ID)
	}
}

func TestJournalStorage_ConcurrentAddsPersistAllEntries(t *testing.T) {
	storage := &JournalStorage{
		filePath: filepath.Join(t.TempDir(), "journal.json"),
		cacheIdx: make(map[string]int),
	}

	const count = 20
	var wg sync.WaitGroup
	errs := make(chan error, count)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs <- storage.AddEntry(&models.JournalEntry{
				ID:      fmt.Sprintf("entry-%02d", i),
				Date:    fmt.Sprintf("2025-01-%02d", i+1),
				Title:   fmt.Sprintf("Entry %d", i),
				Content: "Concurrent entry",
			})
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("AddEntry failed: %v", err)
		}
	}

	data, err := os.ReadFile(storage.filePath)
	if err != nil {
		t.Fatalf("failed to read persisted journal: %v", err)
	}

	var entries []*models.JournalEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("failed to parse persisted journal: %v", err)
	}

	if len(entries) != count {
		t.Fatalf("expected %d persisted entries, got %d", count, len(entries))
	}
}

func TestJournalStorage_GetEntryByID(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	entry := &models.JournalEntry{
		ID:      "test-get",
		Date:    "2025-01-01",
		Title:   "Test Get Entry",
		Content: "Test content",
		Mood:    "testy",
		Tags:    []string{"test"},
	}

	// Add entry
	err := storage.AddEntry(entry)
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Get by ID
	retrieved, err := storage.GetEntryByID("test-get")
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}

	if retrieved.ID != entry.ID {
		t.Errorf("Expected ID %s, got %s", entry.ID, retrieved.ID)
	}

	// Test non-existent entry
	_, err = storage.GetEntryByID("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent entry")
	}
}

func TestJournalStorage_GetEntryByDate(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	entry := &models.JournalEntry{
		ID:      "test-date",
		Date:    "2025-01-01",
		Title:   "Test Date Entry",
		Content: "Test content",
		Mood:    "testy",
		Tags:    []string{"test"},
	}

	// Add entry
	err := storage.AddEntry(entry)
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Get by date
	retrieved, err := storage.GetEntryByDate("2025-01-01")
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}

	if retrieved.ID != entry.ID {
		t.Errorf("Expected ID %s, got %s", entry.ID, retrieved.ID)
	}

	// Test non-existent date
	_, err = storage.GetEntryByDate("2025-99-99")
	if err == nil {
		t.Error("Expected error for non-existent date")
	}
}

func TestJournalStorage_UpdateEntry(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	entry := &models.JournalEntry{
		ID:      "test-update",
		Date:    "2025-01-01",
		Title:   "Original Title",
		Content: "Original content",
		Mood:    "original",
		Tags:    []string{"original"},
	}

	// Add entry
	err := storage.AddEntry(entry)
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Update entry
	entry.Title = "Updated Title"
	entry.Content = "Updated content"
	entry.Mood = "updated"
	entry.Tags = []string{"updated"}
	entry.UpdatedAt = time.Now()

	err = storage.UpdateEntry(entry)
	if err != nil {
		t.Fatalf("Failed to update entry: %v", err)
	}

	// Verify update
	retrieved, err := storage.GetEntryByID("test-update")
	if err != nil {
		t.Fatalf("Failed to get updated entry: %v", err)
	}

	if retrieved.Title != "Updated Title" {
		t.Errorf("Expected updated title, got %s", retrieved.Title)
	}

	if retrieved.Content != "Updated content" {
		t.Errorf("Expected updated content, got %s", retrieved.Content)
	}
}

func TestJournalStorage_SearchEntries(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	entries := []*models.JournalEntry{
		{
			ID:      "test-search-1",
			Date:    "2025-01-01",
			Title:   "Programming Journey",
			Content: "Today I learned Go programming",
			Mood:    "excited",
			Tags:    []string{"programming", "golang"},
		},
		{
			ID:      "test-search-2",
			Date:    "2025-01-02",
			Title:   "Reflections",
			Content: "Thinking about life and goals",
			Mood:    "thoughtful",
			Tags:    []string{"life", "reflection"},
		},
		{
			ID:      "test-search-3",
			Date:    "2025-01-03",
			Title:   "Random thoughts",
			Content: "Nothing special here",
			Mood:    "neutral",
			Tags:    []string{"random"},
		},
	}

	// Add entries
	for _, entry := range entries {
		err := storage.AddEntry(entry)
		if err != nil {
			t.Fatalf("Failed to add entry: %v", err)
		}
	}

	// Test search in content
	results, err := storage.SearchEntries("programming")
	if err != nil {
		t.Fatalf("Failed to search entries: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'programming', got %d", len(results))
	}

	// Test search in title
	results, err = storage.SearchEntries("Journey")
	if err != nil {
		t.Fatalf("Failed to search entries: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'Journey', got %d", len(results))
	}

	// Test search in tags
	results, err = storage.SearchEntries("golang")
	if err != nil {
		t.Fatalf("Failed to search entries: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'golang', got %d", len(results))
	}

	// Test search with no results
	results, err = storage.SearchEntries("nonexistent")
	if err != nil {
		t.Fatalf("Failed to search entries: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results for 'nonexistent', got %d", len(results))
	}
}

func TestJournalStorage_GetEntriesByTag(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	entries := []*models.JournalEntry{
		{
			ID:      "test-tag-1",
			Date:    "2025-01-01",
			Title:   "Work Entry",
			Content: "Today was productive",
			Mood:    "happy",
			Tags:    []string{"work", "productive"},
		},
		{
			ID:      "test-tag-2",
			Date:    "2025-01-02",
			Title:   "Personal Entry",
			Content: "Personal thoughts",
			Mood:    "reflective",
			Tags:    []string{"personal", "reflection"},
		},
		{
			ID:      "test-tag-3",
			Date:    "2025-01-03",
			Title:   "Another Work Entry",
			Content: "More work stuff",
			Mood:    "motivated",
			Tags:    []string{"work"},
		},
	}

	// Add entries
	for _, entry := range entries {
		err := storage.AddEntry(entry)
		if err != nil {
			t.Fatalf("Failed to add entry: %v", err)
		}
	}

	// Test get entries by tag
	results, err := storage.GetEntriesByTag("work")
	if err != nil {
		t.Fatalf("Failed to get entries by tag: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 entries for tag 'work', got %d", len(results))
	}

	// Test tag that doesn't exist
	results, err = storage.GetEntriesByTag("nonexistent")
	if err != nil {
		t.Fatalf("Failed to get entries by tag: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 entries for tag 'nonexistent', got %d", len(results))
	}
}

func TestJournalStorage_GetEntriesByDateRange(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	entries := []*models.JournalEntry{
		{
			ID:      "test-range-1",
			Date:    "2025-01-01",
			Title:   "First Day",
			Content: "Content 1",
			Mood:    "good",
			Tags:    []string{},
		},
		{
			ID:      "test-range-2",
			Date:    "2025-01-02",
			Title:   "Second Day",
			Content: "Content 2",
			Mood:    "better",
			Tags:    []string{},
		},
		{
			ID:      "test-range-3",
			Date:    "2025-01-03",
			Title:   "Third Day",
			Content: "Content 3",
			Mood:    "best",
			Tags:    []string{},
		},
		{
			ID:      "test-range-4",
			Date:    "2025-01-04",
			Title:   "Fourth Day",
			Content: "Content 4",
			Mood:    "amazing",
			Tags:    []string{},
		},
	}

	// Add entries
	for _, entry := range entries {
		err := storage.AddEntry(entry)
		if err != nil {
			t.Fatalf("Failed to add entry: %v", err)
		}
	}

	// Test date range
	results, err := storage.GetEntriesByDateRange("2025-01-02", "2025-01-03")
	if err != nil {
		t.Fatalf("Failed to get entries by date range: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 entries for range 2025-01-02 to 2025-01-03, got %d", len(results))
	}

	// Test single date
	results, err = storage.GetEntriesByDateRange("2025-01-01", "2025-01-01")
	if err != nil {
		t.Fatalf("Failed to get entries by date range: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 entry for single date, got %d", len(results))
	}

	// Test empty range
	results, err = storage.GetEntriesByDateRange("2025-99-99", "2025-99-99")
	if err != nil {
		t.Fatalf("Failed to get entries by date range: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 entries for empty range, got %d", len(results))
	}
}

func TestJournalStorage_GetAllTags(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	entries := []*models.JournalEntry{
		{
			ID:      "test-tags-1",
			Date:    "2025-01-01",
			Title:   "Entry 1",
			Content: "Content 1",
			Mood:    "happy",
			Tags:    []string{"work", "productive", "learning"},
		},
		{
			ID:      "test-tags-2",
			Date:    "2025-01-02",
			Title:   "Entry 2",
			Content: "Content 2",
			Mood:    "thoughtful",
			Tags:    []string{"personal", "reflection", "learning"},
		},
	}

	// Add entries
	for _, entry := range entries {
		err := storage.AddEntry(entry)
		if err != nil {
			t.Fatalf("Failed to add entry: %v", err)
		}
	}

	// Get all tags
	tags, err := storage.GetAllTags()
	if err != nil {
		t.Fatalf("Failed to get all tags: %v", err)
	}

	// Should contain unique tags
	expectedTags := []string{"work", "productive", "learning", "personal", "reflection"}
	if len(tags) != len(expectedTags) {
		t.Errorf("Expected %d unique tags, got %d", len(expectedTags), len(tags))
	}

	for _, expectedTag := range expectedTags {
		found := false
		for _, tag := range tags {
			if tag == expectedTag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected tag '%s' not found in tags", expectedTag)
		}
	}
}

func TestJournalStorage_GetStats(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	entries := []*models.JournalEntry{
		{
			ID:        "test-stats-1",
			Date:      "2025-01-01",
			Title:     "Entry 1",
			Content:   "Content with five words",
			Mood:      "good",
			Tags:      []string{},
			WordCount: 5, // Set explicitly to test stats
		},
		{
			ID:        "test-stats-2",
			Date:      "2025-01-02",
			Title:     "Entry 2",
			Content:   "Content with three words",
			Mood:      "better",
			Tags:      []string{},
			WordCount: 3, // Set explicitly to test stats
		},
	}

	// Add entries
	for _, entry := range entries {
		err := storage.AddEntry(entry)
		if err != nil {
			t.Fatalf("Failed to add entry: %v", err)
		}
	}

	// Get stats
	stats, err := storage.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	// Verify stats
	if stats["total_entries"] != 2 {
		t.Errorf("Expected 2 total entries, got %d", stats["total_entries"])
	}

	// Note: stats["total_words"] calculates from Content, not from WordCount field
	if stats["total_words"] != 8 { // 5 + 3 from content
		t.Errorf("Expected 8 total words from content, got %d", stats["total_words"])
	}
}

func TestJournalStorage_DeleteEntry(t *testing.T) {
	storage := NewJournalStorageWithPath(t.TempDir())

	entry := &models.JournalEntry{
		ID:      "test-delete",
		Date:    "2025-01-01",
		Title:   "Test Delete Entry",
		Content: "Test content",
		Mood:    "testy",
		Tags:    []string{"test"},
	}

	// Add entry
	err := storage.AddEntry(entry)
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Verify entry exists
	_, err = storage.GetEntryByID("test-delete")
	if err != nil {
		t.Fatalf("Failed to get entry before delete: %v", err)
	}

	// Delete entry
	err = storage.DeleteEntry("test-delete")
	if err != nil {
		t.Fatalf("Failed to delete entry: %v", err)
	}

	// Verify entry is gone
	_, err = storage.GetEntryByID("test-delete")
	if err == nil {
		t.Error("Expected error when getting deleted entry")
	}

	// Test delete non-existent entry
	err = storage.DeleteEntry("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent entry")
	}
}
