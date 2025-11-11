package journal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyanite/focus/pkg/models"
)

// Exporter handles journal exports to syntax.sh
type Exporter struct {
	syntaxDir string
}

// NewExporter creates a new journal exporter
func NewExporter() *Exporter {
	// Default syntax.sh imports directory
	homeDir, err := os.UserHomeDir()

	syntaxDir := filepath.Join(".", "syntax", "imports") // Default fallback
	if err == nil {
		syntaxDir = filepath.Join(homeDir, "syntax", "imports")
	}

	return &Exporter{
		syntaxDir: syntaxDir,
	}
}

// ExportToSyntax exports a journal entry to syntax.sh format
func (e *Exporter) ExportToSyntax(entry *models.JournalEntry, exportType models.ExportType) error {
	// Ensure syntax directory exists
	if err := os.MkdirAll(e.syntaxDir, 0o755); err != nil {
		return fmt.Errorf("failed to create syntax imports directory: %w", err)
	}

	// Generate content based on export type
	content := e.formatForExport(entry, exportType)

	// Generate filename
	filename := entry.GetExportFilename(exportType)
	filePath := filepath.Join(e.syntaxDir, filename)

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	return nil
}

// formatForExport formats a journal entry for syntax.sh import
func (e *Exporter) formatForExport(entry *models.JournalEntry, exportType models.ExportType) string {
	var header strings.Builder
	var content strings.Builder

	// Header section
	header.WriteString("# Imported from focus.sh\n\n")
	header.WriteString(fmt.Sprintf("**Source:** Journal entry from %s\n", entry.Date))
	header.WriteString(fmt.Sprintf("**Export Type:** %s\n", e.getExportTypeName(exportType)))
	header.WriteString(fmt.Sprintf("**Tags:** %s\n", strings.Join(entry.Tags, ", ")))
	if entry.Mood != "" {
		header.WriteString(fmt.Sprintf("**Mood:** %s\n", entry.Mood))
	}
	if entry.Title != "" {
		header.WriteString(fmt.Sprintf("**Title:** %s\n", entry.Title))
	}

	// Content section
	content.WriteString("\n---\n\n")
	content.WriteString("## Original Journal Entry\n\n")
	content.WriteString(entry.Content)
	content.WriteString("\n\n---\n\n")
	content.WriteString("## Story Application\n\n")
	content.WriteString(e.getStoryApplicationPrompts(exportType))
	content.WriteString("\n\n")
	content.WriteString("**Next Steps:** Develop this into ")
	content.WriteString(e.getExportTypeName(exportType))
	content.WriteString(" notes, scene, or dialogue as needed.\n")

	return header.String() + content.String()
}

// getExportTypeName returns a user-friendly name for the export type
func (e *Exporter) getExportTypeName(exportType models.ExportType) string {
	switch exportType {
	case models.ExportCharacter:
		return "Character observation"
	case models.ExportDialogue:
		return "Dialogue material"
	case models.ExportScene:
		return "Scene description"
	case models.ExportResearch:
		return "Research notes"
	default:
		return "Journal content"
	}
}

// getStoryApplicationPrompts returns specific prompts based on export type
func (e *Exporter) getStoryApplicationPrompts(exportType models.ExportType) string {
	switch exportType {
	case models.ExportCharacter:
		return `### Character Notes
How does this character approach problems? Do they:
- Break complex situations into manageable parts?
- Show emotional intelligence or introspection?
- Have moments of clarity or breakthrough?
- Demonstrate specific personality traits?

### Dialogue Seed
Consider how this character might express these thoughts in conversation:

### Scene Potential
Imagine a scene where this character is experiencing or reflecting on these events.`

	case models.ExportDialogue:
		return `### Dialogue Analysis
What conversational patterns emerge?
- What topics or themes are important?
- How does emotion affect communication?
- What questions or conflicts arise?

### Dialogue Development
Consider how this could become:
- A conversation between characters
- Internal monologue
- A speech or presentation
- An interview or interrogation`

	case models.ExportScene:
		return `### Scene Elements
What physical and emotional elements are present?
- Setting and atmosphere
- Character movements and actions
- Sensory details and observations
- Emotional undercurrents

### Scene Construction
Think about how this could become:
- A fully realized scene with multiple characters
- A montage or sequence of events
- A flashback or memory
- A dramatic moment of change`

	case models.ExportResearch:
		return `### Research Analysis
What insights or discoveries are documented?
- Patterns or connections
- Questions raised
- Hypotheses formed
- Data or observations collected

### Research Applications
Consider how this could inform:
- World-building elements
- Technical or scientific accuracy
- Historical context
- Character knowledge or expertise`

	default:
		return `### Creative Applications
How can this journal entry inspire creative work?
- What themes or motifs emerge?
- What emotions or ideas could be explored?
- What characters or situations are suggested?

### Development Ideas
Think about how this material could become:
- Character inspiration
- Scene setting
- Plot points
- Thematic elements`
	}
}

// ListExports returns a list of existing export files
func (e *Exporter) ListExports() ([]string, error) {
	files, err := os.ReadDir(e.syntaxDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var exports []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			exports = append(exports, file.Name())
		}
	}

	return exports, nil
}

// RemoveExport removes an existing export file
func (e *Exporter) RemoveExport(filename string) error {
	filePath := filepath.Join(e.syntaxDir, filename)
	return os.Remove(filePath)
}

// GetExportPath returns the path where exports are stored
func (e *Exporter) GetExportPath() string {
	return e.syntaxDir
}
