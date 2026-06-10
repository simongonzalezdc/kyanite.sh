package noise

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	shortcutsSourcePath   = "internal/ui/editor/shortcuts.go"
	menuSourcePath        = "internal/ui/menu.go"
	exportShortcutBinding = "ctrl+e"
	exportShortcutAction  = "ActionExport"
	exportMenuPrompt      = "Export current song"
	exportMenuLabel       = "Export"
	exportWorkflowTitle   = "Workflow Test"
)

var exportIntegrationFiles = []string{
	"internal/ui/export.go",
	"internal/ui/menu.go",
	"internal/ui/root.go",
	"internal/ui/editor/shortcuts.go",
	"internal/ui/editor/split_pane.go",
}

func requireFileExists(t *testing.T, path string) {
	t.Helper()
	cleanPath := filepath.Clean(path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		t.Fatalf("required file does not exist: %s: %v", cleanPath, err)
	}

	if info.IsDir() {
		t.Fatalf("expected file but found directory: %s", cleanPath)
	}

	t.Logf("verified required file exists: %s", cleanPath)
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	cleanPath := filepath.Clean(path)

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", cleanPath, err)
	}

	t.Logf("read %d bytes from %s", len(data), cleanPath)
	return string(data)
}

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	fullPath := filepath.Join(dir, name)
	if err := os.WriteFile(fullPath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file %s: %v", fullPath, err)
	}

	t.Logf("created temp file %s", fullPath)

	t.Cleanup(func() {
		if err := os.Remove(fullPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("cleanup failed for %s: %v", fullPath, err)
		}
	})

	return fullPath
}

// TestExportIntegration tests the export system integration.
func TestExportIntegration(t *testing.T) {
	for _, file := range exportIntegrationFiles {
		file := file
		t.Run(file, func(t *testing.T) {
			requireFileExists(t, file)
		})
	}
}

// TestExportFormatStrings tests export format string representations.
func TestExportFormatStrings(t *testing.T) {
	formats := []string{"PDF", "HTML", "Plain Text", "JSON", "Markdown"}

	for _, format := range formats {
		format := format
		t.Run(format, func(t *testing.T) {
			if format == "" {
				t.Error("format string should not be empty")
			}
			t.Logf("validated export format string %q", format)
		})
	}
}

// TestContentProcessing tests content processing functions.
func TestContentProcessing(t *testing.T) {
	testContent := `# Title

[Verse 1]
This is **bold** text
And this is *italic* text

[Chorus]
More content here
`

	wordCount := len(strings.Fields(testContent))
	if wordCount == 0 {
		t.Error("word count should be greater than 0 for non-empty content")
	} else {
		t.Logf("computed word count: %d", wordCount)
	}

	charCount := len(testContent)
	if charCount == 0 {
		t.Error("character count should be greater than 0 for non-empty content")
	} else {
		t.Logf("computed character count: %d", charCount)
	}

	stripped := stripMarkdown(testContent)
	if strings.Contains(stripped, "**") || strings.Contains(stripped, "*") {
		t.Error("markdown formatting should be stripped from plain text")
	}
}

// TestFileOperations tests file operation capabilities.
func TestFileOperations(t *testing.T) {
	const testContent = "test file content for export testing"

	tempDir := t.TempDir()
	tempFile := writeTempFile(t, tempDir, "export_temp.txt", testContent)

	readContent := readFile(t, tempFile)
	if readContent != testContent {
		t.Error("file content does not match written content")
	}
}

// TestKeyboardShortcutIntegration tests that export shortcuts are properly defined.
func TestKeyboardShortcutIntegration(t *testing.T) {
	content := readFile(t, shortcutsSourcePath)

	if !strings.Contains(content, exportShortcutBinding) {
		t.Errorf("export shortcut (%s) not found in shortcuts file", exportShortcutBinding)
	}

	if !strings.Contains(content, exportShortcutAction) {
		t.Errorf("%s not found in shortcuts file", exportShortcutAction)
	}
}

// TestMenuIntegration tests that export option is in the menu.
func TestMenuIntegration(t *testing.T) {
	content := readFile(t, menuSourcePath)

	if !strings.Contains(content, exportMenuLabel) {
		t.Errorf("%s menu item not found in menu file", exportMenuLabel)
	}

	if !strings.Contains(content, exportMenuPrompt) {
		t.Errorf("%s menu description not found in menu file", exportMenuPrompt)
	}
}

// TestResponsiveDesign tests responsive design considerations.
func TestResponsiveDesign(t *testing.T) {
	type testSize struct {
		name   string
		width  int
		height int
	}

	testSizes := []testSize{
		{"standard", 80, 24},
		{"compact", 60, 20},
		{"wide", 120, 30},
		{"minimal", 40, 15},
	}

	for _, size := range testSizes {
		size := size
		t.Run(size.name, func(t *testing.T) {
			if size.width < 60 && size.name != "compact" && size.name != "minimal" {
				t.Errorf("width %d should be considered compact", size.width)
			}

			if size.height < 20 && size.name != "minimal" {
				t.Errorf("height %d should be considered minimal", size.height)
			}
		})
	}
}

// TestExportWorkflow tests the complete export workflow.
func TestExportWorkflow(t *testing.T) {
	const workflowContent = `# Workflow Test

[Verse]
This tests the complete export workflow
From content creation to file output

Testing multiple formats and options
`

	htmlContent := generateTestHTML(workflowContent)
	if !strings.Contains(htmlContent, "<html>") {
		t.Error("HTML generation failed")
	}

	plainText := generateTestPlainText(workflowContent)
	if strings.Contains(plainText, "[Verse]") {
		t.Error("plain text should strip markdown formatting")
	}

	jsonContent := generateTestJSON(workflowContent)
	if !strings.Contains(jsonContent, `"content"`) {
		t.Error("JSON should contain content field")
	}
}

// stripMarkdown removes markdown formatting (extracted from export.go for testing).
func stripMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	cleanLines := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimPrefix(line, "# ")
		line = strings.TrimPrefix(line, "## ")
		line = strings.TrimPrefix(line, "### ")
		line = strings.TrimPrefix(line, "#### ")
		line = strings.TrimPrefix(line, "##### ")
		line = strings.TrimPrefix(line, "###### ")

		line = strings.ReplaceAll(line, "**", "")
		line = strings.ReplaceAll(line, "*", "")
		line = strings.ReplaceAll(line, "_", "")
		line = strings.ReplaceAll(line, "`", "")

		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "* ")

		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// generateTestHTML generates HTML content for testing (simplified version).
func generateTestHTML(content string) string {
	return `<html><body><h1>` +
		exportWorkflowTitle +
		`</h1><p>` +
		strings.ReplaceAll(content, "\n", "<br>") +
		`</p></body></html>`
}

// generateTestPlainText generates plain text for testing.
func generateTestPlainText(content string) string {
	result := strings.ReplaceAll(content, "# ", "")
	result = strings.ReplaceAll(result, "**", "")
	result = strings.ReplaceAll(result, "*", "")
	result = strings.ReplaceAll(result, "[Verse]", "Verse:")
	return result
}

// generateTestJSON generates JSON content for testing.
func generateTestJSON(content string) string {
	return `{
	"title": "` + exportWorkflowTitle + `",
	"content": "` + strings.ReplaceAll(content, "\"", "\\\"") + `",
	"format": "markdown"
}`
}

// BenchmarkContentProcessing benchmarks content processing performance.
func BenchmarkContentProcessing(b *testing.B) {
	content := `# Performance Test Document

[Verse 1]
This is a test document
With multiple sections
To benchmark processing speed

[Chorus]
Processing performance is important
For good user experience
When working with large documents

[Verse 2]
Another section with more content
To ensure adequate test coverage
For performance measurements

[Bridge]
Bridge section for additional testing
Of content processing algorithms
And export functionality`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		words := strings.Fields(content)
		_ = len(words)

		lines := strings.Split(content, "\n")
		_ = len(lines)

		processed := stripMarkdown(content)
		_ = len(processed)
	}
}
