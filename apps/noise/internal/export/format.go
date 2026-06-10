// Package export provides song export formats and services.
package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	errutil "github.com/Kyanite/noise/internal/errutil"
)

// ExportFormatter handles formatting of exports
type ExportFormatter struct {
	version string
}

// NewExportFormatter creates a new export formatter
func NewExportFormatter() *ExportFormatter {
	return &ExportFormatter{
		version: "1.0.0",
	}
}

// FormatExport formats content for export based on options
func (ef *ExportFormatter) FormatExport(content string, options *ExportOptions) (*NoiseExport, error) {
	now := time.Now()

	// Create metadata
	metadata := ExportMetadata{
		Title:      options.Title,
		BPM:        options.BPM,
		Created:    now,
		ExportedAt: now,
	}

	// Create export
	export := &NoiseExport{
		Type:     "noise_sh_idea",
		Version:  ef.version,
		Metadata: metadata,
	}

	// Parse content based on export type
	switch options.Type {
	case ExportTypePattern:
		patterns := ef.extractPatterns(content)
		export.Patterns = patterns

		// Try to extract BPM from content if not provided
		if options.BPM == 0 {
			bpm := ef.extractBPM(content)
			if bpm > 0 {
				export.Metadata.BPM = bpm
			}
		}

	case ExportTypeLyrics:
		lyrics := ef.extractLyrics(content)
		export.Lyrics = lyrics

	case ExportTypeChords:
		chords := ef.extractChords(content)
		export.Chords = chords

	case ExportTypeFull:
		patterns := ef.extractPatterns(content)
		export.Patterns = patterns

		lyrics := ef.extractLyrics(content)
		if lyrics != "" {
			export.Lyrics = lyrics
		}

		chords := ef.extractChords(content)
		if len(chords) > 0 {
			export.Chords = chords
		}

		// Try to extract BPM from content if not provided
		if options.BPM == 0 {
			bpm := ef.extractBPM(content)
			if bpm > 0 {
				export.Metadata.BPM = bpm
			}
		}

	case ExportTypeMarkdown:
		// For Markdown, we'll format the content directly
		formattedContent := ef.formatAsMarkdown(content, options)
		export.Lyrics = formattedContent

	case ExportTypePlainText:
		// For Plain Text, we'll strip formatting
		plainText := ef.formatAsPlainText(content)
		export.Lyrics = plainText

	case ExportTypeChordPro:
		// Determine BPM if not provided
		actualBPM := options.BPM
		if actualBPM == 0 {
			actualBPM = ef.extractBPM(content)
		}
		if actualBPM > 0 {
			export.Metadata.BPM = actualBPM
		}

		// For ChordPro, we'll format with chord directives, using detected BPM where available
		chordPro := ef.formatAsChordPro(content, options, export.Metadata.BPM)
		export.Lyrics = chordPro
	}

	// Add notes if requested
	if options.IncludeNotes {
		notes := ef.extractNotes(content)
		if notes != "" {
			export.Notes = notes
		}
	}

	return export, nil
}

// SaveToFile saves the export to a file
func (ef *ExportFormatter) SaveToFile(export *NoiseExport, outputPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return errutil.Wrap(err, "create directory")
	}

	// Determine file format based on extension
	ext := strings.ToLower(filepath.Ext(outputPath))

	var data []byte
	var err error

	switch ext {
	case ".md":
		// For Markdown, save the lyrics content directly
		data = []byte(export.Lyrics)
	case ".txt":
		// For Plain Text, save the lyrics content directly
		data = []byte(export.Lyrics)
	case ".cho":
		// For ChordPro, save the lyrics content directly
		data = []byte(export.Lyrics)
	default:
		// For existing formats, save as JSON
		data, err = json.MarshalIndent(export, "", "  ")
		if err != nil {
			return errutil.Wrap(err, "marshal export")
		}
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return errutil.Wrap(err, "write file")
	}

	return nil
}

// extractPatterns extracts pattern data from content
func (ef *ExportFormatter) extractPatterns(content string) []string {
	var patterns []string

	// Split content into lines
	lines := strings.Split(content, "\n")

	// Find pattern sections (lines starting with "pattern:" or similar)
	patternRegex := regexp.MustCompile(`^\s*(pattern|code|data):\s*(.*)`)
	inPattern := false
	var currentPattern strings.Builder

	for _, line := range lines {
		matches := patternRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			// Save previous pattern if exists
			if inPattern && currentPattern.Len() > 0 {
				patterns = append(patterns, strings.TrimSpace(currentPattern.String()))
				currentPattern.Reset()
			}

			// Start new pattern
			inPattern = true
			if matches[2] != "" {
				currentPattern.WriteString(matches[2])
			}
		} else if inPattern {
			// Continue pattern if indented or not empty
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) {
				currentPattern.WriteString("\n")
				currentPattern.WriteString(line)
			} else if trimmed != "" {
				// End of pattern
				if currentPattern.Len() > 0 {
					patterns = append(patterns, strings.TrimSpace(currentPattern.String()))
					currentPattern.Reset()
				}
				inPattern = false
			}
		}
	}

	// Save last pattern if exists
	if inPattern && currentPattern.Len() > 0 {
		patterns = append(patterns, strings.TrimSpace(currentPattern.String()))
	}

	// If no patterns found, treat entire content as a pattern
	if len(patterns) == 0 {
		patterns = append(patterns, content)
	}

	return patterns
}

// extractLyrics extracts lyrics from content
func (ef *ExportFormatter) extractLyrics(content string) string {
	// Remove code blocks and patterns
	content = ef.removeCodeBlocks(content)
	content = ef.removePatterns(content)

	// Remove markdown formatting
	content = ef.removeMarkdown(content)

	// Clean up extra whitespace
	content = strings.TrimSpace(content)

	return content
}

// extractChords extracts chords from content
func (ef *ExportFormatter) extractChords(content string) []string {
	var chords []string

	// Look for chord lines (lines with chord symbols)
	lines := strings.Split(content, "\n")
	chordRegex := regexp.MustCompile(`\b[A-G]([#b])?(m|maj|min|dim|aug|sus|add)?\d*\b`)

	for _, line := range lines {
		// Skip lines that are clearly not chord lines
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		// Find all chord symbols in the line
		matches := chordRegex.FindAllString(line, -1)
		if len(matches) > 0 {
			// If this line has multiple chords, it's likely a chord line
			if len(matches) >= 2 || (len(matches) == 1 && len(trimmed) < 20) {
				chords = append(chords, trimmed)
			}
		}
	}

	return chords
}

// extractBPM extracts BPM from content
func (ef *ExportFormatter) extractBPM(content string) int {
	// Look for BPM markers
	bpmRegex := regexp.MustCompile(`(?i)(bpm|tempo)[:\s]+(\d+)`)
	matches := bpmRegex.FindStringSubmatch(content)

	if len(matches) >= 3 {
		bpm, err := strconv.Atoi(matches[2])
		if err == nil && bpm > 0 && bpm < 300 {
			return bpm
		}
	}

	return 0
}

// extractNotes extracts notes from content
func (ef *ExportFormatter) extractNotes(content string) string {
	var notes strings.Builder

	// Look for note sections
	lines := strings.Split(content, "\n")
	inNotes := false
	noteRegex := regexp.MustCompile(`(?i)^\s*(note|notes|comment|idea)[:\s]*(.*)`)

	for _, line := range lines {
		matches := noteRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			inNotes = true
			if matches[2] != "" {
				notes.WriteString(matches[2])
				notes.WriteString("\n")
			}
		} else if inNotes {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "//") {
				notes.WriteString(line)
				notes.WriteString("\n")
			} else if trimmed == "" {
				// Empty line might end the notes section
				// Continue for now in case there are more notes
				notes.WriteString("\n")
			}
		}
	}

	return strings.TrimSpace(notes.String())
}

// removeCodeBlocks removes code blocks from content
func (ef *ExportFormatter) removeCodeBlocks(content string) string {
	codeBlockRegex := regexp.MustCompile("(?s)```[\\s\\S]*?```")
	return codeBlockRegex.ReplaceAllString(content, "")
}

// removePatterns removes pattern sections from content
func (ef *ExportFormatter) removePatterns(content string) string {
	patternRegex := regexp.MustCompile(`(?m)^[ \t]*pattern:.*`)
	lines := strings.Split(content, "\n")
	result := make([]string, 0, 32)
	skip := false

	for _, line := range lines {
		if patternRegex.MatchString(line) {
			skip = true
			continue
		}

		if skip {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				skip = false
				continue
			}

			if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
				continue // Skip indented lines of the pattern
			} else {
				skip = false
			}
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// removeMarkdown removes markdown formatting from content
func (ef *ExportFormatter) removeMarkdown(content string) string {
	// Remove headers
	headerRegex := regexp.MustCompile(`^#+\s+`)
	content = headerRegex.ReplaceAllString(content, "")

	// Remove bold/italic
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	content = boldRegex.ReplaceAllString(content, "$1")

	italicRegex := regexp.MustCompile(`\*([^*]+)\*`)
	content = italicRegex.ReplaceAllString(content, "$1")

	// Remove links
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	content = linkRegex.ReplaceAllString(content, "$1")

	// Remove inline code
	codeRegex := regexp.MustCompile("`([^`]+)`")
	content = codeRegex.ReplaceAllString(content, "$1")

	// Remove blockquotes
	blockquoteRegex := regexp.MustCompile(`^>\s+`)
	content = blockquoteRegex.ReplaceAllString(content, "")

	return content
}

// formatAsMarkdown formats content as Markdown
func (ef *ExportFormatter) formatAsMarkdown(content string, options *ExportOptions) string {
	var builder strings.Builder

	// Add title as header if provided
	if options.Title != "" {
		builder.WriteString("# ")
		builder.WriteString(options.Title)
		builder.WriteString("\n\n")
	}

	// Add metadata if available
	if options.BPM > 0 {
		builder.WriteString("**BPM:** ")
		builder.WriteString(strconv.Itoa(options.BPM))
		builder.WriteString("\n\n")
	}

	// Process content line by line
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			builder.WriteString("\n")
			continue
		}

		// Check for section headers
		if ef.isSectionHeader(trimmed) {
			// Convert to markdown header
			sectionType := ef.extractSectionType(trimmed)
			if sectionType != "" {
				builder.WriteString("## ")
				builder.WriteString(strings.ToUpper(sectionType[:1]) + sectionType[1:])
				builder.WriteString("\n\n")
				continue
			}
		}

		// Check for chord lines
		if ef.isChordLine(trimmed) {
			// Format chords with code formatting for visibility
			builder.WriteString("`")
			builder.WriteString(trimmed)
			builder.WriteString("`\n")
			continue
		}

		// Regular lyric line
		builder.WriteString(trimmed)
		builder.WriteString("\n")
	}

	return builder.String()
}

// formatAsPlainText formats content as plain text
func (ef *ExportFormatter) formatAsPlainText(content string) string {
	var builder strings.Builder

	// Process content line by line
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			builder.WriteString("\n")
			continue
		}

		// Remove markdown formatting
		cleanLine := ef.removeMarkdown(trimmed)

		// Skip chord lines in plain text (optional - could include them)
		if ef.isChordLine(trimmed) {
			continue // Skip chord lines for cleaner lyrics
		}

		builder.WriteString(cleanLine)
		builder.WriteString("\n")
	}

	return builder.String()
}

// formatAsChordPro formats content as ChordPro
func (ef *ExportFormatter) formatAsChordPro(content string, options *ExportOptions, bpm int) string {
	var builder strings.Builder

	// Add ChordPro metadata directives
	if options.Title != "" {
		builder.WriteString("{title:")
		builder.WriteString(options.Title)
		builder.WriteString("}\n")
	}

	actualBPM := options.BPM
	if actualBPM == 0 {
		actualBPM = bpm
	}

	if actualBPM > 0 {
		builder.WriteString("{tempo:")
		builder.WriteString(strconv.Itoa(actualBPM))
		builder.WriteString("}\n")
	}

	// Add key if detected
	key := ef.detectKey(content)
	if key != "" {
		builder.WriteString("{key:")
		builder.WriteString(key)
		builder.WriteString("}\n")
	}

	builder.WriteString("\n")

	// Process content line by line
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			builder.WriteString("\n")
			continue
		}

		// Check for section headers
		if ef.isSectionHeader(trimmed) {
			// Convert to ChordPro directive
			sectionType := ef.extractSectionType(trimmed)
			if sectionType != "" {
				builder.WriteString("{start_of_")
				builder.WriteString(sectionType)
				builder.WriteString("}\n")
				continue
			}
		}

		// Check for chord lines and format them inline with lyrics
		if ef.isChordLine(trimmed) {
			// Look ahead for the corresponding lyric line
			chordLine := trimmed
			builder.WriteString(chordLine)
			builder.WriteString("\n")
			continue
		}

		// Regular lyric line
		builder.WriteString(trimmed)
		builder.WriteString("\n")
	}

	return builder.String()
}

// Helper methods for the new formatters

// isSectionHeader checks if a line is a section header
func (ef *ExportFormatter) isSectionHeader(line string) bool {
	sectionHeaders := []string{
		"verse", "chorus", "bridge", "intro", "outro", "pre-chorus",
		"[verse]", "[chorus]", "[bridge]", "[intro]", "[outro]", "[pre-chorus]",
	}

	lowerLine := strings.ToLower(line)
	for _, header := range sectionHeaders {
		if strings.Contains(lowerLine, header) {
			return true
		}
	}
	return false
}

// extractSectionType extracts the section type from a header line
func (ef *ExportFormatter) extractSectionType(line string) string {
	lowerLine := strings.ToLower(line)

	sections := map[string]string{
		"verse":      "verse",
		"chorus":     "chorus",
		"bridge":     "bridge",
		"intro":      "intro",
		"outro":      "outro",
		"pre-chorus": "pre-chorus",
		"pre chorus": "pre-chorus",
	}

	for key, value := range sections {
		if strings.Contains(lowerLine, key) {
			return value
		}
	}
	return ""
}

// isChordLine checks if a line contains primarily chords
func (ef *ExportFormatter) isChordLine(line string) bool {
	// Count chord symbols vs regular text
	chordRegex := regexp.MustCompile(`\b[A-G]([#b])?(m|maj|min|dim|aug|sus|add)?\d*\b`)
	chords := chordRegex.FindAllString(line, -1)

	// If we have multiple chords and the line is short, it's likely a chord line
	if len(chords) >= 2 && len(strings.TrimSpace(line)) < 30 {
		return true
	}

	// If we have chords and very few non-chord characters
	nonChordChars := chordRegex.ReplaceAllString(line, "")
	nonChordChars = strings.TrimSpace(nonChordChars)

	if len(chords) > 0 && len(nonChordChars) < 5 {
		return true
	}

	return false
}

// detectKey attempts to detect the key of the song from chords
func (ef *ExportFormatter) detectKey(content string) string {
	chordRegex := regexp.MustCompile(`\b([A-G])([#b])?\b`)
	matches := chordRegex.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		return ""
	}

	// Simple key detection - find the most common root note
	rootCount := make(map[string]int)
	for _, match := range matches {
		if len(match) > 1 {
			root := match[1]
			rootCount[root]++
		}
	}

	// Find the most common root
	var mostCommonRoot string
	maxCount := 0
	for root, count := range rootCount {
		if count > maxCount {
			maxCount = count
			mostCommonRoot = root
		}
	}

	return mostCommonRoot
}
