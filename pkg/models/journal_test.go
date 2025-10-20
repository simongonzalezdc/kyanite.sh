package models

import (
	"strings"
	"testing"
	"time"
)

func TestNewJournalEntry(t *testing.T) {
	// Test basic entry creation
	entry := NewJournalEntry("2025-01-01", "Test Title", "Test content", "happy", "daily_log", []string{"test"})
	
	if entry.Date != "2025-01-01" {
		t.Errorf("Expected date '2025-01-01', got '%s'", entry.Date)
	}
	
	if entry.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", entry.Title)
	}
	
	if entry.Content != "Test content" {
		t.Errorf("Expected content 'Test content', got '%s'", entry.Content)
	}
	
	if entry.Mood != "happy" {
		t.Errorf("Expected mood 'happy', got '%s'", entry.Mood)
	}
	
	if entry.TemplateUsed != "daily_log" {
		t.Errorf("Expected template 'daily_log', got '%s'", entry.TemplateUsed)
	}
	
	if len(entry.Tags) != 1 || entry.Tags[0] != "test" {
		t.Errorf("Expected tags ['test'], got %v", entry.Tags)
	}
	
	if entry.WordCount != 2 {
		t.Errorf("Expected word count 2, got %d", entry.WordCount)
	}
	
	if entry.IsPrivate {
		t.Error("Expected IsPrivate to be false")
	}
}

func TestNewJournalEntry_DefaultDate(t *testing.T) {
	entry := NewJournalEntry("", "Title", "Content", "mood", "template", []string{})
	
	expectedDate := time.Now().Format("2006-01-02")
	if entry.Date != expectedDate {
		t.Errorf("Expected date '%s', got '%s'", expectedDate, entry.Date)
	}
}

func TestNewJournalEntry_EmptyContent(t *testing.T) {
	entry := NewJournalEntry("2025-01-01", "Title", "", "mood", "template", []string{})
	
	if entry.WordCount != 0 {
		t.Errorf("Expected word count 0 for empty content, got %d", entry.WordCount)
	}
}

func TestJournalEntry_UpdateContent(t *testing.T) {
	entry := NewJournalEntry("2025-01-01", "Title", "Original content", "mood", "template", []string{})
	
	originalUpdatedAt := entry.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	
	entry.UpdateContent("Updated content")
	
	if entry.Content != "Updated content" {
		t.Errorf("Expected content 'Updated content', got '%s'", entry.Content)
	}
	
	if entry.WordCount != 2 {
		t.Errorf("Expected word count 2 for updated content, got %d", entry.WordCount)
	}
	
	if !entry.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestJournalEntry_UpdateTags(t *testing.T) {
	entry := NewJournalEntry("2025-01-01", "Title", "Content", "mood", "template", []string{"original"})
	
	originalUpdatedAt := entry.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	
	newTags := []string{"updated1", "updated2"}
	entry.UpdateTags(newTags)
	
	if len(entry.Tags) != 2 || entry.Tags[0] != "updated1" || entry.Tags[1] != "updated2" {
		t.Errorf("Expected tags %v, got %v", newTags, entry.Tags)
	}
	
	if !entry.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestJournalEntry_HasTag(t *testing.T) {
	entry := NewJournalEntry("2025-01-01", "Title", "Content", "mood", "template", []string{"work", "productive"})
	
	// Test existing tag
	if !entry.HasTag("work") {
		t.Error("Expected HasTag to return true for existing tag")
	}
	
	if !entry.HasTag("productive") {
		t.Error("Expected HasTag to return true for existing tag")
	}
	
	// Test non-existing tag
	if entry.HasTag("nonexistent") {
		t.Error("Expected HasTag to return false for non-existing tag")
	}
	
	// Test empty tags
	emptyTagsEntry := NewJournalEntry("2025-01-01", "Title", "Content", "mood", "template", []string{})
	if emptyTagsEntry.HasTag("anything") {
		t.Error("Expected HasTag to return false for entry with no tags")
	}
}

func TestJournalEntry_GetExportFilename(t *testing.T) {
	entry := NewJournalEntry("2025-01-01", "Title", "Content", "mood", "template", []string{})
	
	filename := entry.GetExportFilename(ExportCharacter)
	expected := "character-2025-01-01.md"
	
	if filename != expected {
		t.Errorf("Expected filename '%s', got '%s'", expected, filename)
	}
}

func TestJournalTemplates(t *testing.T) {
	// Test that templates exist
	if len(JournalTemplates) == 0 {
		t.Error("Expected JournalTemplates to contain templates")
	}
	
	// Test specific templates
	var foundMorningPages bool
	var foundEveningReflection bool
	var foundDailyLog bool
	
	for _, template := range JournalTemplates {
		switch template.Name {
		case "morning_pages":
			foundMorningPages = true
		case "evening_reflection":
			foundEveningReflection = true
		case "daily_log":
			foundDailyLog = true
		}
	}
	
	if !foundMorningPages {
		t.Error("Expected morning_pages template")
	}
	
	if !foundEveningReflection {
		t.Error("Expected evening_reflection template")
	}
	
	if !foundDailyLog {
		t.Error("Expected daily_log template")
	}
}

func TestJournalTemplateStructure(t *testing.T) {
	// Test that templates have required fields
	for _, template := range JournalTemplates {
		if template.Name == "" {
			t.Error("Template name should not be empty")
		}
		
		if template.Description == "" {
			t.Error("Template description should not be empty")
		}
		
		if len(template.Prompts) == 0 {
			t.Error("Template should have at least one prompt")
		}
		
		if len(template.Tags) == 0 {
			t.Error("Template should have suggested tags")
		}
	}
}

func TestExportType_String(t *testing.T) {
	tests := []struct {
		exportType ExportType
		expected    string
	}{
		{ExportCharacter, "character"},
		{ExportDialogue, "dialogue"},
		{ExportScene, "scene"},
		{ExportResearch, "research"},
	}
	
	for _, test := range tests {
		if test.exportType.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.exportType.String())
		}
	}
}

func TestJournalEntry_IDGeneration(t *testing.T) {
	// Add a small delay to ensure unique IDs
	time.Sleep(1 * time.Millisecond)
	entry1 := NewJournalEntry("2025-01-01", "Title", "Content", "mood", "template", []string{})
	time.Sleep(1 * time.Millisecond)
	entry2 := NewJournalEntry("2025-01-01", "Title", "Content", "mood", "template", []string{})
	
	// IDs should be unique
	if entry1.ID == entry2.ID {
		t.Errorf("Expected unique IDs for different entries, both got: %s", entry1.ID)
	}
	
	// ID should start with "journal_"
	if !strings.HasPrefix(entry1.ID, "journal_") {
		t.Errorf("Expected ID to start with 'journal_', got '%s'", entry1.ID)
	}
	
	if !strings.HasPrefix(entry2.ID, "journal_") {
		t.Errorf("Expected ID to start with 'journal_', got '%s'", entry2.ID)
	}
}

func TestJournalEntry_WordCountCalculation(t *testing.T) {
	tests := []struct {
		content   string
		wordCount int
	}{
		{"", 0},
		{"word", 1},
		{"word word", 2},
		{"word   word", 2}, // multiple spaces
		{"word\nword", 2}, // newline
		{"   word   word   ", 2}, // leading/trailing spaces
		{"Hello world, how are you?", 5},
		{"This is a test of the word count function.", 9},
	}
	
	for _, test := range tests {
		entry := NewJournalEntry("2025-01-01", "Title", test.content, "mood", "template", []string{})
		if entry.WordCount != test.wordCount {
			t.Errorf("Expected word count %d for content '%s', got %d", test.wordCount, test.content, entry.WordCount)
		}
	}
}