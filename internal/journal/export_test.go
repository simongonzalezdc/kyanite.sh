package journal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kyanite/focus/pkg/models"
)

func TestExporter_NewExporter(t *testing.T) {
	exporter := NewExporter()

	// Test that exporter creates directory with actual user home
	homeDir, _ := os.UserHomeDir()
	expectedDir := filepath.Join(homeDir, "syntax", "imports")

	if exporter.syntaxDir != expectedDir {
		t.Errorf("Expected syntaxDir to be %s, got %s", expectedDir, exporter.syntaxDir)
	}
}

func TestExporter_ExportToSyntax(t *testing.T) {
	exporter := NewExporter()

	// Create test entry
	entry := &models.JournalEntry{
		ID:           "test-123",
		Date:         "2025-01-01",
		Title:        "Test Entry",
		Content:      "This is a test journal entry for character development.",
		Mood:         "productive",
		Tags:         []string{"test", "character"},
		WordCount:    10,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsPrivate:    false,
		TemplateUsed: "daily_log",
	}

	// Test character export
	err := exporter.ExportToSyntax(entry, models.ExportCharacter)
	if err != nil {
		t.Fatalf("Failed to export character: %v", err)
	}

	// Verify file exists
	filename := filepath.Join(exporter.syntaxDir, "character-2025-01-01.md")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist", filename)
	}

	// Clean up
	os.Remove(filename)
}

func TestExporter_GetExportFilename(t *testing.T) {
	entry := &models.JournalEntry{
		Date: "2025-01-01",
	}

	filename := entry.GetExportFilename(models.ExportCharacter)
	expected := "character-2025-01-01.md"

	if filename != expected {
		t.Errorf("Expected filename %s, got %s", expected, filename)
	}
}

func TestExporter_ListExports(t *testing.T) {
	exporter := NewExporter()

	// Clean up any existing files first
	os.RemoveAll(exporter.syntaxDir)

	// Test with empty directory
	exports, err := exporter.ListExports()
	if err != nil {
		t.Fatalf("Failed to list exports: %v", err)
	}

	if len(exports) != 0 {
		t.Errorf("Expected 0 exports, got %d", len(exports))
	}

	// Create a test file
	testFile := filepath.Join(exporter.syntaxDir, "test.md")
	testDir := filepath.Dir(testFile)
	os.MkdirAll(testDir, 0o755)
	os.WriteFile(testFile, []byte("test content"), 0o644)

	// Test listing with file
	exports, err = exporter.ListExports()
	if err != nil {
		t.Fatalf("Failed to list exports: %v", err)
	}

	if len(exports) != 1 {
		t.Errorf("Expected 1 export, got %d", len(exports))
	}

	// Clean up
	os.RemoveAll(exporter.syntaxDir)
}

func TestExporter_RemoveExport(t *testing.T) {
	exporter := NewExporter()

	// Create a test file
	testFile := filepath.Join(exporter.syntaxDir, "test.md")
	testDir := filepath.Dir(testFile)
	os.MkdirAll(testDir, 0o755)
	os.WriteFile(testFile, []byte("test content"), 0o644)

	// Remove the file
	err := exporter.RemoveExport("test.md")
	if err != nil {
		t.Fatalf("Failed to remove export: %v", err)
	}

	// Verify file doesn't exist
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("Expected file to be removed")
	}

	// Clean up
	os.RemoveAll(exporter.syntaxDir)
}

func TestExporter_GetExportPath(t *testing.T) {
	exporter := NewExporter()
	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, "syntax", "imports")

	if exporter.GetExportPath() != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, exporter.GetExportPath())
	}
}

func TestExporter_formatForExport(t *testing.T) {
	exporter := NewExporter()

	entry := &models.JournalEntry{
		ID:      "test",
		Date:    "2025-01-01",
		Title:   "Test Entry",
		Content: "Test content",
		Mood:    "happy",
		Tags:    []string{"test"},
	}

	// Test different export types
	tests := []struct {
		name          models.ExportType
		shouldContain []string
	}{
		{
			name:          models.ExportCharacter,
			shouldContain: []string{"Character observation", "Character Notes", "Dialogue Seed", "Scene Potential"},
		},
		{
			name:          models.ExportDialogue,
			shouldContain: []string{"Dialogue material", "Dialogue Analysis", "Dialogue Development"},
		},
		{
			name:          models.ExportScene,
			shouldContain: []string{"Scene description", "Scene Elements", "Scene Construction"},
		},
		{
			name:          models.ExportResearch,
			shouldContain: []string{"Research notes", "Research Analysis", "Research Applications"},
		},
	}

	for _, test := range tests {
		t.Run(string(test.name), func(t *testing.T) {
			content := exporter.formatForExport(entry, test.name)

			for _, shouldContain := range test.shouldContain {
				if !strings.Contains(content, shouldContain) {
					t.Errorf("Expected content to contain '%s', but it didn't", shouldContain)
				}
			}
		})
	}
}
