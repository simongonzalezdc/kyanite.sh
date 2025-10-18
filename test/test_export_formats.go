package noise

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/puente-labs/noise/internal/export"
)

func TestExportMarkdown(t *testing.T) {
	// Create a temporary directory for test exports
	tempDir, err := os.MkdirTemp("", "noise_export_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create export service
	service := export.NewExportService(tempDir)

	// Test content with various elements
	content := `# My Song

This is a test song with chords and lyrics.

[Verse]
C        G        Am        F
This is the first line of the verse
G        C        G         C
This is the second line of the verse

[Chorus]
F        C        G        Am
This is the chorus line
F        C        G        C
This is the second chorus line

BPM: 120
Key: C`

	// Export as Markdown
	outputPath, err := service.ExportToMarkdown(content, "Test Song")
	if err != nil {
		t.Fatalf("Failed to export to Markdown: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Export file was not created: %s", outputPath)
	}

	// Verify file extension
	if filepath.Ext(outputPath) != ".md" {
		t.Errorf("Expected .md extension, got %s", filepath.Ext(outputPath))
	}

	// Read and verify content
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	contentStr := string(data)
	
	// Check for markdown formatting
	if !strings.Contains(contentStr, "# Test Song") {
		t.Error("Markdown title not found")
	}
	
	if !strings.Contains(contentStr, "**BPM:** 120") {
		t.Error("BPM metadata not found")
	}
	
	if !strings.Contains(contentStr, "## Verse") {
		t.Error("Section header not formatted correctly")
	}
	
	if !strings.Contains(contentStr, "`C        G        Am        F`") {
		t.Error("Chord line not formatted with code formatting")
	}
}

func TestExportPlainText(t *testing.T) {
	// Create a temporary directory for test exports
	tempDir, err := os.MkdirTemp("", "noise_export_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create export service
	service := export.NewExportService(tempDir)

	// Test content with various elements
	content := `# My Song

This is a test song with chords and lyrics.

[Verse]
C        G        Am        F
This is the first line of the verse
G        C        G         C
This is the second line of the verse

[Chorus]
F        C        G        Am
This is the chorus line
F        C        G        C
This is the second chorus line

BPM: 120
Key: C`

	// Export as Plain Text
	outputPath, err := service.ExportToPlainText(content, "Test Song")
	if err != nil {
		t.Fatalf("Failed to export to Plain Text: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Export file was not created: %s", outputPath)
	}

	// Verify file extension
	if filepath.Ext(outputPath) != ".txt" {
		t.Errorf("Expected .txt extension, got %s", filepath.Ext(outputPath))
	}

	// Read and verify content
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	contentStr := string(data)
	
	// Check that markdown formatting is removed
	if strings.Contains(contentStr, "# ") {
		t.Error("Markdown headers should be removed in plain text")
	}
	
	if strings.Contains(contentStr, "**") {
		t.Error("Markdown bold formatting should be removed in plain text")
	}
	
	// Check that chord lines are skipped
	if strings.Contains(contentStr, "C        G        Am        F") {
		t.Error("Chord lines should be skipped in plain text export")
	}
	
	// Check that lyrics are preserved
	if !strings.Contains(contentStr, "This is the first line of the verse") {
		t.Error("Lyrics should be preserved in plain text export")
	}
}

func TestExportChordPro(t *testing.T) {
	// Create a temporary directory for test exports
	tempDir, err := os.MkdirTemp("", "noise_export_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create export service
	service := export.NewExportService(tempDir)

	// Test content with various elements
	content := `# My Song

This is a test song with chords and lyrics.

[Verse]
C        G        Am        F
This is the first line of the verse
G        C        G         C
This is the second line of the verse

[Chorus]
F        C        G        Am
This is the chorus line
F        C        G        C
This is the second chorus line

BPM: 120
Key: C`

	// Export as ChordPro
	outputPath, err := service.ExportToChordPro(content, "Test Song")
	if err != nil {
		t.Fatalf("Failed to export to ChordPro: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Export file was not created: %s", outputPath)
	}

	// Verify file extension
	if filepath.Ext(outputPath) != ".cho" {
		t.Errorf("Expected .cho extension, got %s", filepath.Ext(outputPath))
	}

	// Read and verify content
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	contentStr := string(data)
	
	// Check for ChordPro metadata directives
	if !strings.Contains(contentStr, "{title:Test Song}") {
		t.Error("ChordPro title directive not found")
	}
	
	if !strings.Contains(contentStr, "{tempo:120}") {
		t.Error("ChordPro tempo directive not found")
	}
	
	if !strings.Contains(contentStr, "{key:C}") {
		t.Error("ChordPro key directive not found")
	}
	
	// Check for ChordPro section directives
	if !strings.Contains(contentStr, "{start_of_verse}") {
		t.Error("ChordPro section directive not found")
	}
	
	// Check that chords are preserved
	if !strings.Contains(contentStr, "C        G        Am        F") {
		t.Error("Chord lines should be preserved in ChordPro export")
	}
	
	// Check that lyrics are preserved
	if !strings.Contains(contentStr, "This is the first line of the verse") {
		t.Error("Lyrics should be preserved in ChordPro export")
	}
}

func TestExportFormatDetection(t *testing.T) {
	// Test the helper methods for format detection
	
	// Create a formatter to test the helper methods
	formatter := export.NewExportFormatter()
	
	// Test section header detection
	testCases := []struct {
		line     string
		expected bool
	}{
		{"[Verse]", true},
		{"[Chorus]", true},
		{"[Bridge]", true},
		{"Verse", true},
		{"Chorus:", true},
		{"This is a lyric line", false},
		{"C        G        Am        F", false},
	}
	
	for _, tc := range testCases {
		// We can't directly test private methods, but we can test through the export
		// functionality which uses these methods internally
		content := tc.line + "\nThis is a test line"
		
		// Test with Markdown export
		options := export.DefaultExportOptions()
		options.Type = export.ExportTypeMarkdown
		options.Title = "Test"
		
		result, err := formatter.FormatExport(content, options)
		if err != nil {
			t.Fatalf("Failed to format export: %v", err)
		}
		
		// Check if section was detected and formatted
		hasSectionHeader := strings.Contains(result.Lyrics, "## ")
		if tc.expected && !hasSectionHeader {
			t.Errorf("Expected section header for line: %s", tc.line)
		}
		if !tc.expected && hasSectionHeader {
			t.Errorf("Did not expect section header for line: %s", tc.line)
		}
	}
}

func TestExportKeyDetection(t *testing.T) {
	// Test key detection from chord content
	content := `C        G        Am        F
This is a test line
G        D        Em        C
Another test line`
	
	// Create a formatter to test key detection
	formatter := export.NewExportFormatter()
	
	// Test with ChordPro export
	options := export.DefaultExportOptions()
	options.Type = export.ExportTypeChordPro
	options.Title = "Test"
	
	result, err := formatter.FormatExport(content, options)
	if err != nil {
		t.Fatalf("Failed to format export: %v", err)
	}
	
	// Check if key was detected (should be C or G based on frequency)
	if !strings.Contains(result.Lyrics, "{key:") {
		t.Error("Key detection failed - no key directive found")
	}
}

func TestExportBPMDetection(t *testing.T) {
	// Test BPM detection from content
	content := `BPM: 120
This is a test line
tempo: 140
Another test line`
	
	// Create a formatter to test BPM detection
	formatter := export.NewExportFormatter()
	
	// Test with ChordPro export
	options := export.DefaultExportOptions()
	options.Type = export.ExportTypeChordPro
	options.Title = "Test"
	
	result, err := formatter.FormatExport(content, options)
	if err != nil {
		t.Fatalf("Failed to format export: %v", err)
	}
	
	// Check if BPM was detected
	if !strings.Contains(result.Lyrics, "{tempo:") {
		t.Error("BPM detection failed - no tempo directive found")
	}
	
	// Should detect the first BPM (120)
	if !strings.Contains(result.Lyrics, "{tempo:120}") {
		t.Error("Expected BPM 120 not detected correctly")
	}
}