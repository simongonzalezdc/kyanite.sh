package noise

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kyanite/noise/internal/export"
)

const (
	sampleExportTitle        = "Test Song"
	sampleChordLine          = "C        G        Am        F"
	samplePlainTextLyricLine = "This is the first line of the verse"
	sampleBPMMarkdown        = "**BPM:** 120"
	chordProTitleDirective   = "{title:Test Song}"
	chordProTempoDirective   = "{tempo:120}"
	chordProKeyDirective     = "{key:C}"
	markdownSectionVerse     = "## Verse"
	markdownChordLine        = "`" + sampleChordLine + "`"
	testTitle                = "Test"
)

var sampleExportContent = strings.Join([]string{
	"# My Song",
	"",
	"This is a test song with chords and lyrics.",
	"",
	"[Verse]",
	sampleChordLine,
	samplePlainTextLyricLine,
	"G        C        G         C",
	"This is the second line of the verse",
	"",
	"[Chorus]",
	"F        C        G        Am",
	"This is the chorus line",
	"F        C        G        C",
	"This is the second chorus line",
	"",
	"BPM: 120",
	"Key: C",
}, "\n")

func newExportService(t *testing.T) *export.ExportService {
	t.Helper()
	tempDir := t.TempDir()
	service := export.NewExportService(tempDir)
	t.Logf("initialized export service with tempDir=%s", tempDir)
	return service
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	cleanPath := filepath.Clean(path)
	info, err := os.Stat(cleanPath)
	if err != nil {
		t.Fatalf("expected file to exist: %s: %v", cleanPath, err)
	}
	if info.IsDir() {
		t.Fatalf("expected file, found directory: %s", cleanPath)
	}
	t.Logf("verified export file exists: %s", cleanPath)
}

func readExportFile(t *testing.T, path string) []byte {
	t.Helper()
	cleanPath := filepath.Clean(path)
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		t.Fatalf("failed to read export file %s: %v", cleanPath, err)
	}
	t.Logf("read %d bytes from export file %s", len(data), cleanPath)
	return data
}

// TestExportMarkdown tests Markdown export functionality.
func TestExportMarkdown(t *testing.T) {
	service := newExportService(t)

	outputPath, err := service.ExportToMarkdown(sampleExportContent, sampleExportTitle)
	if err != nil {
		t.Fatalf("failed to export to Markdown: %v", err)
	}

	t.Logf("markdown export produced output at %s", outputPath)

	if filepath.Ext(outputPath) != ".md" {
		t.Errorf("expected .md extension, got %s", filepath.Ext(outputPath))
	}

	assertFileExists(t, outputPath)

	contentStr := string(readExportFile(t, outputPath))

	if !strings.Contains(contentStr, "# Test Song") {
		t.Error("markdown title not found")
	}

	if !strings.Contains(contentStr, sampleBPMMarkdown) {
		t.Error("BPM metadata not found")
	}

	if !strings.Contains(contentStr, markdownSectionVerse) {
		t.Error("section header not formatted correctly")
	}

	if !strings.Contains(contentStr, markdownChordLine) {
		t.Error("chord line not formatted with code formatting")
	}
}

// TestExportPlainText tests plain text export functionality.
func TestExportPlainText(t *testing.T) {
	service := newExportService(t)

	outputPath, err := service.ExportToPlainText(sampleExportContent, sampleExportTitle)
	if err != nil {
		t.Fatalf("failed to export to Plain Text: %v", err)
	}

	t.Logf("plain text export produced output at %s", outputPath)

	if filepath.Ext(outputPath) != ".txt" {
		t.Errorf("expected .txt extension, got %s", filepath.Ext(outputPath))
	}

	assertFileExists(t, outputPath)

	contentStr := string(readExportFile(t, outputPath))

	if strings.Contains(contentStr, "# ") {
		t.Error("markdown headers should be removed in plain text")
	}

	if strings.Contains(contentStr, "**") {
		t.Error("markdown bold formatting should be removed in plain text")
	}

	if strings.Contains(contentStr, sampleChordLine) {
		t.Error("chord lines should be skipped in plain text export")
	}

	if !strings.Contains(contentStr, samplePlainTextLyricLine) {
		t.Error("lyrics should be preserved in plain text export")
	}
}

// TestExportChordPro tests ChordPro export functionality.
func TestExportChordPro(t *testing.T) {
	service := newExportService(t)

	outputPath, err := service.ExportToChordPro(sampleExportContent, sampleExportTitle)
	if err != nil {
		t.Fatalf("failed to export to ChordPro: %v", err)
	}

	t.Logf("ChordPro export produced output at %s", outputPath)

	if filepath.Ext(outputPath) != ".cho" {
		t.Errorf("expected .cho extension, got %s", filepath.Ext(outputPath))
	}

	assertFileExists(t, outputPath)

	contentStr := string(readExportFile(t, outputPath))

	if !strings.Contains(contentStr, chordProTitleDirective) {
		t.Error("ChordPro title directive not found")
	}

	if !strings.Contains(contentStr, chordProTempoDirective) {
		t.Error("ChordPro tempo directive not found")
	}

	if !strings.Contains(contentStr, chordProKeyDirective) {
		t.Error("ChordPro key directive not found")
	}

	if !strings.Contains(contentStr, "{start_of_verse}") {
		t.Error("ChordPro section directive not found")
	}

	if !strings.Contains(contentStr, sampleChordLine) {
		t.Error("chord lines should be preserved in ChordPro export")
	}

	if !strings.Contains(contentStr, samplePlainTextLyricLine) {
		t.Error("lyrics should be preserved in ChordPro export")
	}
}

// TestExportFormatDetection tests format detection logic.
func TestExportFormatDetection(t *testing.T) {
	formatter := export.NewExportFormatter()

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
		{sampleChordLine, false},
	}

	for _, tc := range testCases {
		content := tc.line + "\nThis is a test line"

		options := export.DefaultExportOptions()
		options.Type = export.ExportTypeMarkdown
		options.Title = testTitle

		result, err := formatter.FormatExport(content, options)
		if err != nil {
			t.Fatalf("failed to format export: %v", err)
		}

		hasSectionHeader := strings.Contains(result.Lyrics, "## ")
		t.Logf("format detection for %q produced section header=%t", tc.line, hasSectionHeader)

		if tc.expected && !hasSectionHeader {
			t.Errorf("expected section header for line: %s", tc.line)
		}
		if !tc.expected && hasSectionHeader {
			t.Errorf("did not expect section header for line: %s", tc.line)
		}
	}
}

// TestExportKeyDetection tests key detection from chords.
func TestExportKeyDetection(t *testing.T) {
	content := `C        G        Am        F
This is a test line
G        D        Em        C
Another test line`

	formatter := export.NewExportFormatter()

	options := export.DefaultExportOptions()
	options.Type = export.ExportTypeChordPro
	options.Title = "Test"

	result, err := formatter.FormatExport(content, options)
	if err != nil {
		t.Fatalf("failed to format export: %v", err)
	}

	t.Logf("key detection produced lyrics: %s", result.Lyrics)

	if !strings.Contains(result.Lyrics, "{key:") {
		t.Error("key detection failed - no key directive found")
	}
}

// TestExportBPMDetection tests BPM detection from content.
func TestExportBPMDetection(t *testing.T) {
	content := `BPM: 120
This is a test line
tempo: 140
Another test line`

	formatter := export.NewExportFormatter()

	options := export.DefaultExportOptions()
	options.Type = export.ExportTypeChordPro
	options.Title = "Test"

	result, err := formatter.FormatExport(content, options)
	if err != nil {
		t.Fatalf("failed to format export: %v", err)
	}

	t.Logf("BPM detection produced lyrics: %s", result.Lyrics)

	if !strings.Contains(result.Lyrics, "{tempo:") {
		t.Error("BPM detection failed - no tempo directive found")
	}

	if !strings.Contains(result.Lyrics, "{tempo:120}") {
		t.Error("expected BPM 120 not detected correctly")
	}
}
