package models

import (
	"fmt"
	"strings"
	"time"
)

// JournalEntry represents a single journal entry
type JournalEntry struct {
	ID           string    `json:"id"`
	Date         string    `json:"date"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Mood         string    `json:"mood"`
	Tags         []string  `json:"tags"`
	WordCount    int       `json:"word_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsPrivate    bool      `json:"is_private"`
	TemplateUsed string    `json:"template_used"`
}

// ExportType represents the type of export for syntax.sh
type ExportType string

const (
	ExportCharacter ExportType = "character"
	ExportDialogue  ExportType = "dialogue"
	ExportScene     ExportType = "scene"
	ExportResearch  ExportType = "research"
)

// String returns the string representation of ExportType
func (e ExportType) String() string {
	return string(e)
}

// JournalTemplate represents a journal template
type JournalTemplate struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Prompts     []string `json:"prompts"`
	Tags        []string `json:"suggested_tags"`
}

// Predefined journal templates
var JournalTemplates = []JournalTemplate{
	{
		Name:        "morning_pages",
		Description: "Stream of consciousness morning writing",
		Prompts:     []string{"What's on your mind?", "What are your intentions for today?", "What are you grateful for?"},
		Tags:        []string{"morning", "reflection", "gratitude"},
	},
	{
		Name:        "evening_reflection",
		Description: "End-of-day reflection on accomplishments and learnings",
		Prompts:     []string{"What went well today?", "What did you learn?", "What could be better tomorrow?"},
		Tags:        []string{"evening", "reflection", "growth"},
	},
	{
		Name:        "daily_log",
		Description: "Simple daily log of events and feelings",
		Prompts:     []string{"What happened today?", "How did you feel?", "Key highlights?"},
		Tags:        []string{"daily", "log", "events"},
	},
}

// NewJournalEntry creates a new journal entry with current timestamp
func NewJournalEntry(date, title, content, mood, template string, tags []string) *JournalEntry {
	now := time.Now()
	entry := &JournalEntry{
		ID:           generateID(),
		Date:         date,
		Title:        title,
		Content:      content,
		Mood:         mood,
		Tags:         tags,
		WordCount:    countWords(content),
		CreatedAt:    now,
		UpdatedAt:    now,
		IsPrivate:    false,
		TemplateUsed: template,
	}
	
	if entry.Date == "" {
		entry.Date = now.Format("2006-01-02")
	}
	
	return entry
}

// UpdateContent updates the entry content and metadata
func (j *JournalEntry) UpdateContent(content string) {
	j.Content = content
	j.WordCount = countWords(content)
	j.UpdatedAt = time.Now()
}

// UpdateTags updates the entry tags
func (j *JournalEntry) UpdateTags(tags []string) {
	j.Tags = tags
	j.UpdatedAt = time.Now()
}

// HasTag checks if the entry has a specific tag
func (j *JournalEntry) HasTag(tag string) bool {
	for _, t := range j.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// GetExportFilename generates a filename for syntax.sh export
func (j *JournalEntry) GetExportFilename(exportType ExportType) string {
	return fmt.Sprintf("%s-%s.md", exportType, j.Date)
}

// Helper functions
func generateID() string {
	return fmt.Sprintf("journal_%d", time.Now().UnixNano())
}

func countWords(content string) int {
	if content == "" {
		return 0
	}
	words := strings.Fields(content)
	return len(words)
}