package export

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	errutil "github.com/Kyanite/noise/internal/errutil"
)

// ExportService handles the export functionality.
type ExportService struct {
	formatter *ExportFormatter
	outputDir string
}

// NewExportService creates a new export service.
func NewExportService(outputDir string) *ExportService {
	sanitized := sanitizeOutputDirectory(outputDir)
	return &ExportService{
		formatter: NewExportFormatter(),
		outputDir: sanitized,
	}
}

// Export exports content based on options.
func (es *ExportService) Export(content string, options *ExportOptions) (string, error) {
	exportData, err := es.formatter.FormatExport(content, options)
	if err != nil {
		return "", errutil.Wrap(err, "format export")
	}

	outputPath := options.OutputPath
	if outputPath == "" {
		outputPath = es.generateOutputPath(options)
	}

	baseDir := es.outputDir
	if baseDir == "" {
		baseDir = "."
	}
	if !filepath.IsAbs(baseDir) {
		baseDir, err = filepath.Abs(baseDir)
		if err != nil {
			return "", errutil.Wrap(err, "resolve export directory")
		}
		es.outputDir = baseDir
	}

	cleanPath := filepath.Clean(outputPath)
	if !filepath.IsAbs(cleanPath) {
		cleanPath = filepath.Join(baseDir, cleanPath)
	}

	relPath, err := filepath.Rel(baseDir, cleanPath)
	if err != nil {
		return "", errutil.Wrap(err, "validate export path")
	}
	if relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) || strings.HasPrefix(relPath, "../") {
		return "", errutil.Wrap(fmt.Errorf("output path escapes export directory"), "validate export path")
	}

	if err := es.formatter.SaveToFile(exportData, cleanPath); err != nil {
		return "", errutil.Wrap(err, "save export")
	}

	return cleanPath, nil
}

// QuickExport performs a quick export with default options.
func (es *ExportService) QuickExport(content, title string) (string, error) {
	options := DefaultExportOptions()
	options.Title = title

	return es.Export(content, options)
}

// ExportToPattern exports content as a pattern.
func (es *ExportService) ExportToPattern(content, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypePattern
	options.Title = title

	return es.Export(content, options)
}

// ExportToLyrics exports content as lyrics.
func (es *ExportService) ExportToLyrics(content, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypeLyrics
	options.Title = title

	return es.Export(content, options)
}

// ExportToChords exports content as chords.
func (es *ExportService) ExportToChords(content, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypeChords
	options.Title = title

	return es.Export(content, options)
}

// ExportFull exports all content types.
func (es *ExportService) ExportFull(content, title string, bpm int, includeNotes bool) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypeFull
	options.Title = title
	options.BPM = bpm
	options.IncludeNotes = includeNotes

	return es.Export(content, options)
}

// ExportToMarkdown exports content as Markdown.
func (es *ExportService) ExportToMarkdown(content, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypeMarkdown
	options.Title = title

	return es.Export(content, options)
}

// ExportToPlainText exports content as plain text.
func (es *ExportService) ExportToPlainText(content, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypePlainText
	options.Title = title

	return es.Export(content, options)
}

// ExportToChordPro exports content as ChordPro.
func (es *ExportService) ExportToChordPro(content, title string) (string, error) {
	options := DefaultExportOptions()
	options.Type = ExportTypeChordPro
	options.Title = title

	return es.Export(content, options)
}

// ListExports returns a list of all exports in the output directory.
func (es *ExportService) ListExports() ([]string, error) {
	var exports []string

	if err := os.MkdirAll(es.outputDir, 0o700); err != nil {
		return nil, errutil.Wrap(err, "create output directory")
	}

	entries, err := os.ReadDir(es.outputDir)
	if err != nil {
		return nil, errutil.Wrap(err, "read output directory")
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			exports = append(exports, entry.Name())
		}
	}

	return exports, nil
}

// DeleteExport deletes an export file.
func (es *ExportService) DeleteExport(filename string) error {
	baseDir, err := filepath.Abs(es.outputDir)
	if err != nil {
		return errutil.Wrap(err, "resolve export directory")
	}

	cleanFilename := filepath.Clean(filename)
	if !filepath.IsAbs(cleanFilename) {
		cleanFilename = filepath.Join(baseDir, cleanFilename)
	}

	relPath, err := filepath.Rel(baseDir, cleanFilename)
	if err != nil {
		return errutil.Wrap(err, "validate export path")
	}
	if relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) || strings.HasPrefix(relPath, "../") {
		return errutil.Wrap(fmt.Errorf("export filename escapes output directory"), "validate export path")
	}

	if err := os.Remove(cleanFilename); err != nil {
		return errutil.Wrap(err, "delete export")
	}

	return nil
}

// GetExportPath returns the full path to an export file.
func (es *ExportService) GetExportPath(filename string) string {
	cleanFilename := filepath.Clean(filename)
	if filepath.IsAbs(cleanFilename) {
		return cleanFilename
	}
	return filepath.Join(es.outputDir, cleanFilename)
}

// SetOutputDir sets the output directory for exports.
func (es *ExportService) SetOutputDir(outputDir string) {
	es.outputDir = sanitizeOutputDirectory(outputDir)
}

// GetOutputDir returns the current output directory.
func (es *ExportService) GetOutputDir() string {
	return es.outputDir
}

// generateOutputPath generates a unique output path for an export.
func (es *ExportService) generateOutputPath(options *ExportOptions) string {
	title := options.Title
	if title == "" {
		title = "untitled"
	}

	sanitized := sanitizeFilename(title)
	timestamp := time.Now().Format("20060102_150405")

	var extension string
	switch options.Type {
	case ExportTypeMarkdown:
		extension = "md"
	case ExportTypePlainText:
		extension = "txt"
	case ExportTypeChordPro:
		extension = "cho"
	default:
		extension = "json"
	}

	return fmt.Sprintf("%s_%s.%s", sanitized, timestamp, extension)
}

func sanitizeOutputDirectory(outputDir string) string {
	if outputDir == "" {
		outputDir = "."
	}
	absDir, err := filepath.Abs(outputDir)
	if err != nil {
		return filepath.Clean(outputDir)
	}
	return absDir
}

// sanitizeFilename sanitizes a string for use as a filename.
func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, " ", "_")

	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalid {
		name = strings.ReplaceAll(name, char, "")
	}

	if name == "" {
		name = "untitled"
	}

	return name
}
