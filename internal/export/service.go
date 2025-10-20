package export

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	errutil "github.com/Kyanite/noise/internal/errutil"
)

// ExportService handles the export functionality
type ExportService struct {
	formatter *ExportFormatter
	outputDir string
}

// NewExportService creates a new export service
func NewExportService(outputDir string) *ExportService {
	return &ExportService{
		formatter: NewExportFormatter(),
		outputDir: outputDir,
	}
}

// Export exports content based on options
func (es *ExportService) Export(content string, options *ExportOptions) (string, error) {
	// Format the export
	export, err := es.formatter.FormatExport(content, options)
	if err != nil {
		return "", errutil.Wrap(err, "format export")
	}
	
	// Generate output path if not provided
	outputPath := options.OutputPath
	if outputPath == "" {
		outputPath = es.generateOutputPath(options)
	}
	
	// Ensure the path is absolute
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(es.outputDir, outputPath)
	}
	
	// Save to file
	if err := es.formatter.SaveToFile(export, outputPath); err != nil {
		return "", errutil.Wrap(err, "save export")
	}
	
	return outputPath, nil
}

// QuickExport performs a quick export with default options
func (es *ExportService) QuickExport(content string, title string) (string, error) {
	options := DefaultExportOptions()
	options.Title = title
	
	return es.Export(content, options)
}

// ExportToPattern exports content as a pattern
func (es *ExportService) ExportToPattern(content string, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypePattern
	options.Title = title
	
	return es.Export(content, options)
}

// ExportToLyrics exports content as lyrics
func (es *ExportService) ExportToLyrics(content string, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypeLyrics
	options.Title = title
	
	return es.Export(content, options)
}

// ExportToChords exports content as chords
func (es *ExportService) ExportToChords(content string, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypeChords
	options.Title = title
	
	return es.Export(content, options)
}

// ExportFull exports all content types
func (es *ExportService) ExportFull(content string, title string, bpm int, includeNotes bool) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypeFull
	options.Title = title
	options.BPM = bpm
	options.IncludeNotes = includeNotes
	
	return es.Export(content, options)
}

// ExportToMarkdown exports content as Markdown
func (es *ExportService) ExportToMarkdown(content string, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypeMarkdown
	options.Title = title
	
	return es.Export(content, options)
}

// ExportToPlainText exports content as plain text
func (es *ExportService) ExportToPlainText(content string, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypePlainText
	options.Title = title
	
	return es.Export(content, options)
}

// ExportToChordPro exports content as ChordPro
func (es *ExportService) ExportToChordPro(content string, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypeChordPro
	options.Title = title
	
	return es.Export(content, options)
}

// ListExports returns a list of all exports in the output directory
func (es *ExportService) ListExports() ([]string, error) {
	var exports []string
	
	// Ensure output directory exists
	if err := os.MkdirAll(es.outputDir, 0755); err != nil {
		return nil, errutil.Wrap(err, "create output directory")
	}
	
	// Read directory
	entries, err := os.ReadDir(es.outputDir)
	if err != nil {
		return nil, errutil.Wrap(err, "read output directory")
	}
	
	// Filter for JSON files
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			exports = append(exports, entry.Name())
		}
	}
	
	return exports, nil
}

// DeleteExport deletes an export file
func (es *ExportService) DeleteExport(filename string) error {
	// Ensure the path is absolute
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(es.outputDir, filename)
	}
	
	// Delete file
	if err := os.Remove(filename); err != nil {
		return errutil.Wrap(err, "delete export")
	}
	
	return nil
}

// GetExportPath returns the full path to an export file
func (es *ExportService) GetExportPath(filename string) string {
	if !filepath.IsAbs(filename) {
		return filepath.Join(es.outputDir, filename)
	}
	return filename
}

// SetOutputDir sets the output directory for exports
func (es *ExportService) SetOutputDir(outputDir string) {
	es.outputDir = outputDir
}

// GetOutputDir returns the current output directory
func (es *ExportService) GetOutputDir() string {
	return es.outputDir
}

// generateOutputPath generates a unique output path for an export
func (es *ExportService) generateOutputPath(options *ExportOptions) string {
	// Create filename from title and timestamp
	title := options.Title
	if title == "" {
		title = "untitled"
	}
	
	// Sanitize title
	sanitized := sanitizeFilename(title)
	
	// Add timestamp
	timestamp := time.Now().Format("20060102_150405")
	
	// Determine file extension based on export type
	var extension string
	switch options.Type {
	case ExportTypeMarkdown:
		extension = "md"
	case ExportTypePlainText:
		extension = "txt"
	case ExportTypeChordPro:
		extension = "cho"
	default:
		extension = "json" // Keep JSON for existing formats
	}
	
	filename := fmt.Sprintf("%s_%s.%s", sanitized, timestamp, extension)
	
	return filename
}

// sanitizeFilename sanitizes a string for use as a filename
func sanitizeFilename(name string) string {
	// Replace spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")
	
	// Remove invalid characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalid {
		name = strings.ReplaceAll(name, char, "")
	}
	
	// Ensure it's not empty
	if name == "" {
		name = "untitled"
	}
	
	return name
}