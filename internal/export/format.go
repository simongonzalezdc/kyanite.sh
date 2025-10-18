package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
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
		patterns, err := ef.extractPatterns(content)
		if err != nil {
			return nil, fmt.Errorf("failed to extract patterns: %w", err)
		}
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
		patterns, err := ef.extractPatterns(content)
		if err != nil {
			return nil, fmt.Errorf("failed to extract patterns: %w", err)
		}
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
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// extractPatterns extracts pattern data from content
func (ef *ExportFormatter) extractPatterns(content string) ([]string, error) {
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
	
	return patterns, nil
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
	chordRegex := regexp.MustCompile(`\b[A-G](#|b)?(m|maj|min|dim|aug|sus|add)?[0-9]*\b`)
	
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
	patternRegex := regexp.MustCompile("(?m)^[ \\t]*pattern:.*")
	lines := strings.Split(content, "\n")
	var result []string
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