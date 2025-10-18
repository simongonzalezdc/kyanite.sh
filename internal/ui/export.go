package ui

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/noise/internal/ui/styles"
)

// ExportFormat represents different export formats
type ExportFormat int

const (
	FormatPDF ExportFormat = iota
	FormatHTML
	FormatPlainText
	FormatJSON
	FormatMarkdown
	FormatChordPro
)

// String returns the string representation of the export format
func (f ExportFormat) String() string {
	switch f {
	case FormatPDF:
		return "PDF"
	case FormatHTML:
		return "HTML"
	case FormatPlainText:
		return "Plain Text"
	case FormatJSON:
		return "JSON"
	case FormatMarkdown:
		return "Markdown"
	case FormatChordPro:
		return "ChordPro"
	default:
		return "Unknown"
	}
}

// ExportOptions contains options for exporting content
type ExportOptions struct {
	Format          ExportFormat
	IncludeMetadata bool
	IncludeStyles   bool
	TemplateName    string
	OutputPath      string
	CustomTemplate  string
}

// ExportMetadata contains metadata about the exported content
type ExportMetadata struct {
	Title          string    `json:"title"`
	Author         string    `json:"author,omitempty"`
	ExportTime     time.Time `json:"export_time"`
	WordCount      int       `json:"word_count"`
	CharacterCount int       `json:"character_count"`
	Format         string    `json:"format"`
	Application    string    `json:"application"`
}

// ExportResult contains the result of an export operation
type ExportResult struct {
	Success      bool
	OutputPath   string
	ErrorMessage string
	Format       ExportFormat
}

// ExportModel handles the export functionality
type ExportModel struct {
	width        int
	height       int
	focused      bool
	selected     int
	options      []ExportFormat
	showProgress bool
	progressMsg  string
	result       *ExportResult

	// Content to export
	content  string
	metadata *ExportMetadata

	// Styles
	focusedStyle  lipgloss.Style
	blurredStyle  lipgloss.Style
	selectedStyle lipgloss.Style
	errorStyle    lipgloss.Style
	successStyle  lipgloss.Style
}

// NewExportModel creates a new export model
func NewExportModel(content string) *ExportModel {
	metadata := &ExportMetadata{
		Title:          "Untitled Document",
		ExportTime:     time.Now(),
		WordCount:      countWords(content),
		CharacterCount: len(content),
		Application:    "noise.sh",
	}

	model := &ExportModel{
		content:      content,
		metadata:     metadata,
		options:      []ExportFormat{FormatPDF, FormatHTML, FormatPlainText, FormatJSON, FormatMarkdown, FormatChordPro},
		selected:     0,
		showProgress: false,
		focusedStyle: styles.BorderActive,
		blurredStyle: styles.Border,
		selectedStyle: lipgloss.NewStyle().
			Foreground(styles.Background).
			Background(styles.Primary).
			Bold(true).
			Padding(0, 1),
		errorStyle: lipgloss.NewStyle().
			Foreground(styles.Background).
			Background(styles.Error).
			Bold(true).
			Padding(0, 1),
		successStyle: lipgloss.NewStyle().
			Foreground(styles.Background).
			Background(styles.Success).
			Bold(true).
			Padding(0, 1),
	}

	return model
}

// Init initializes the export model
func (m *ExportModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the export model
func (m *ExportModel) Update(msg tea.Msg) (*ExportModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.focused && !m.showProgress {
				m.selected = max(0, m.selected-1)
			}
		case "down", "j":
			if m.focused && !m.showProgress {
				m.selected = min(len(m.options)-1, m.selected+1)
			}
		case "enter":
			if m.focused && !m.showProgress {
				return m, m.performExport()
			}
		case "esc":
			if !m.showProgress {
				return m, func() tea.Msg {
					return BackMsg{}
				}
			}
		}
	}

	return m, nil
}

// View renders the export interface
func (m *ExportModel) View() string {
	var style lipgloss.Style
	if m.focused {
		style = m.focusedStyle
	} else {
		style = m.blurredStyle
	}

	title := styles.TitleGradient("📤 Export Options")
	title = lipgloss.NewStyle().Bold(true).Padding(0, 2).Render(title)

	// Format options
	var options []string
	for i, format := range m.options {
		var option string
		if i == m.selected && m.focused {
			option = m.selectedStyle.Render("▶ " + format.String())
		} else {
			option = "  " + format.String()
		}

		// Add format descriptions
		switch format {
		case FormatPDF:
			option += " - Portable document with formatting"
		case FormatHTML:
			option += " - Web page with styling"
		case FormatPlainText:
			option += " - Clean text without formatting"
		case FormatJSON:
			option += " - Structured data format"
		case FormatMarkdown:
			option += " - Original markdown format"
		case FormatChordPro:
			option += " - ChordPro format for lyrics with chords"
		}

		options = append(options, option)
	}

	optionsView := strings.Join(options, "\n")

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Render("\n↑↓ Navigate • Enter: Export • Esc: Back")

	// Progress indicator
	progressView := ""
	if m.showProgress {
		progressView = lipgloss.NewStyle().
			Foreground(styles.Accent).
			Render("\n⏳ " + m.progressMsg + "...")
	}

	// Result message
	resultView := ""
	if m.result != nil {
		if m.result.Success {
			resultView = m.successStyle.Render(fmt.Sprintf("\n✅ Export successful: %s", m.result.OutputPath))
		} else {
			resultView = m.errorStyle.Render(fmt.Sprintf("\n❌ Export failed: %s", m.result.ErrorMessage))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		optionsView,
		instructions,
		progressView,
		resultView,
	)

	return style.Width(m.width).Height(m.height).Render(content)
}

// performExport performs the actual export operation
func (m *ExportModel) performExport() tea.Cmd {
	return func() tea.Msg {
		m.showProgress = true
		m.progressMsg = "Preparing export"

		format := m.options[m.selected]
		m.metadata.Format = format.String()

		// Generate filename with timestamp
		timestamp := time.Now().Format("20060102_150405")
		ext := m.getFileExtension(format)
		filename := fmt.Sprintf("noise.sh_export_%s.%s", timestamp, ext)

		// Default output path
		if m.result == nil {
			m.result = &ExportResult{Format: format}
		}
		m.result.OutputPath = filepath.Join(".", filename)

		// Perform export based on format
		success := false
		var errorMsg string

		switch format {
		case FormatPDF:
			success, errorMsg = m.exportToPDF()
		case FormatHTML:
			success, errorMsg = m.exportToHTML()
		case FormatPlainText:
			success, errorMsg = m.exportToPlainText()
		case FormatJSON:
			success, errorMsg = m.exportToJSON()
		case FormatMarkdown:
			success, errorMsg = m.exportToMarkdown()
		case FormatChordPro:
			success, errorMsg = m.exportToChordPro()
		}

		m.showProgress = false
		m.result.Success = success
		if !success {
			m.result.ErrorMessage = errorMsg
		}

		return ExportCompleteMsg{Result: *m.result}
	}
}

// getFileExtension returns the appropriate file extension for the format
func (m *ExportModel) getFileExtension(format ExportFormat) string {
	switch format {
	case FormatPDF:
		return "pdf"
	case FormatHTML:
		return "html"
	case FormatPlainText:
		return "txt"
	case FormatJSON:
		return "json"
	case FormatMarkdown:
		return "md"
	case FormatChordPro:
		return "cho"
	default:
		return "txt"
	}
}

// exportToPDF exports content to PDF format
func (m *ExportModel) exportToPDF() (bool, string) {
	// For now, we'll create an HTML file that can be converted to PDF
	// In a full implementation, you would use a PDF library like gofpdf or chromedp
	htmlContent, err := m.generateHTMLContent(true)
	if err != nil {
		return false, fmt.Sprintf("Failed to generate PDF content: %v", err)
	}

	// Write HTML content (can be converted to PDF using external tools)
	err = os.WriteFile(m.result.OutputPath, []byte(htmlContent), 0644)
	if err != nil {
		return false, fmt.Sprintf("Failed to write PDF file: %v", err)
	}

	return true, ""
}

// exportToHTML exports content to HTML format
func (m *ExportModel) exportToHTML() (bool, string) {
	htmlContent, err := m.generateHTMLContent(true)
	if err != nil {
		return false, fmt.Sprintf("Failed to generate HTML: %v", err)
	}

	err = os.WriteFile(m.result.OutputPath, []byte(htmlContent), 0644)
	if err != nil {
		return false, fmt.Sprintf("Failed to write HTML file: %v", err)
	}

	return true, ""
}

// exportToPlainText exports content to plain text format
func (m *ExportModel) exportToPlainText() (bool, string) {
	// Strip markdown formatting for plain text
	textContent := m.stripMarkdown(m.content)

	err := os.WriteFile(m.result.OutputPath, []byte(textContent), 0644)
	if err != nil {
		return false, fmt.Sprintf("Failed to write text file: %v", err)
	}

	return true, ""
}

// exportToJSON exports content to JSON format
func (m *ExportModel) exportToJSON() (bool, string) {
	jsonData := struct {
		Metadata ExportMetadata `json:"metadata"`
		Content  string         `json:"content"`
	}{
		Metadata: *m.metadata,
		Content:  m.content,
	}

	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return false, fmt.Sprintf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(m.result.OutputPath, jsonBytes, 0644)
	if err != nil {
		return false, fmt.Sprintf("Failed to write JSON file: %v", err)
	}

	return true, ""
}

// exportToMarkdown exports content to markdown format (as-is)
func (m *ExportModel) exportToMarkdown() (bool, string) {
	err := os.WriteFile(m.result.OutputPath, []byte(m.content), 0644)
	if err != nil {
		return false, fmt.Sprintf("Failed to write markdown file: %v", err)
	}

	return true, ""
}

// exportToChordPro exports content to ChordPro format
func (m *ExportModel) exportToChordPro() (bool, string) {
	// Convert content to ChordPro format
	chordProContent := m.convertToChordPro(m.content)
	
	err := os.WriteFile(m.result.OutputPath, []byte(chordProContent), 0644)
	if err != nil {
		return false, fmt.Sprintf("Failed to write ChordPro file: %v", err)
	}

	return true, ""
}

// convertToChordPro converts content to ChordPro format
func (m *ExportModel) convertToChordPro(content string) string {
	var builder strings.Builder
	
	// Add ChordPro metadata directives
	if m.metadata.Title != "" {
		builder.WriteString("{title:")
		builder.WriteString(m.metadata.Title)
		builder.WriteString("}\n")
	}
	
	// Add tempo if available (we don't have BPM in the UI metadata, so we'll extract it)
	bpm := m.extractBPM(content)
	if bpm > 0 {
		builder.WriteString("{tempo:")
		builder.WriteString(strconv.Itoa(bpm))
		builder.WriteString("}\n")
	}
	
	// Add key if detected
	key := m.detectKey(content)
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
		if m.isSectionHeader(trimmed) {
			// Convert to ChordPro directive
			sectionType := m.extractSectionType(trimmed)
			if sectionType != "" {
				// Close previous section if open
				builder.WriteString("{end_of_")
				builder.WriteString(sectionType)
				builder.WriteString("}\n")
				// Start new section
				builder.WriteString("{start_of_")
				builder.WriteString(sectionType)
				builder.WriteString("}\n")
				continue
			}
		}
		
		// Check for chord lines and format them inline with lyrics
		if m.isChordLine(trimmed) {
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

// Helper methods for ChordPro conversion

// extractBPM extracts BPM from content
func (m *ExportModel) extractBPM(content string) int {
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

// detectKey attempts to detect the key of the song from chords
func (m *ExportModel) detectKey(content string) string {
	chordRegex := regexp.MustCompile(`\b([A-G])(#|b)?\b`)
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

// isSectionHeader checks if a line is a section header
func (m *ExportModel) isSectionHeader(line string) bool {
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
func (m *ExportModel) extractSectionType(line string) string {
	lowerLine := strings.ToLower(line)
	
	// Remove brackets if present
	line = strings.ReplaceAll(line, "[", "")
	line = strings.ReplaceAll(line, "]", "")
	
	sections := map[string]string{
		"verse": "verse",
		"chorus": "chorus",
		"bridge": "bridge",
		"intro": "intro",
		"outro": "outro",
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
func (m *ExportModel) isChordLine(line string) bool {
	// Count chord symbols vs regular text
	chordRegex := regexp.MustCompile(`\b[A-G](#|b)?(m|maj|min|dim|aug|sus|add)?[0-9]*\b`)
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

// generateHTMLContent generates HTML content with optional styling
func (m *ExportModel) generateHTMLContent(includeStyles bool) (string, error) {
	var htmlBuilder strings.Builder

	// HTML header
	htmlBuilder.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	htmlBuilder.WriteString("<meta charset=\"UTF-8\">\n")
	htmlBuilder.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	htmlBuilder.WriteString(fmt.Sprintf("<title>%s</title>\n", html.EscapeString(m.metadata.Title)))

	if includeStyles {
		// Add embedded styles
		htmlBuilder.WriteString("<style>\n")
		htmlBuilder.WriteString(m.getEmbeddedStyles())
		htmlBuilder.WriteString("</style>\n")
	}

	htmlBuilder.WriteString("</head>\n<body>\n")

	// Add metadata if enabled
	if m.result != nil && m.result.Success { // This indicates includeMetadata option
		htmlBuilder.WriteString("<div class=\"metadata\">\n")
		htmlBuilder.WriteString(fmt.Sprintf("<h1>%s</h1>\n", html.EscapeString(m.metadata.Title)))
		if m.metadata.Author != "" {
			htmlBuilder.WriteString(fmt.Sprintf("<p><strong>Author:</strong> %s</p>\n", html.EscapeString(m.metadata.Author)))
		}
		htmlBuilder.WriteString(fmt.Sprintf("<p><strong>Exported:</strong> %s</p>\n", m.metadata.ExportTime.Format("January 2, 2006 at 3:04 PM")))
		htmlBuilder.WriteString(fmt.Sprintf("<p><strong>Words:</strong> %d</p>\n", m.metadata.WordCount))
		htmlBuilder.WriteString("</div>\n\n")
	}

	// Convert markdown to HTML (simplified conversion)
	htmlContent := m.markdownToHTML(m.content)
	htmlBuilder.WriteString(htmlContent)

	htmlBuilder.WriteString("\n</body>\n</html>")

	return htmlBuilder.String(), nil
}

// getEmbeddedStyles returns CSS styles for HTML export
func (m *ExportModel) getEmbeddedStyles() string {
	return `
body {
	font-family: 'Georgia', 'Times New Roman', serif;
	line-height: 1.6;
	color: #333;
	max-width: 800px;
	margin: 0 auto;
	padding: 20px;
	background: #fff;
}

.metadata {
	border-bottom: 2px solid #333;
	padding-bottom: 20px;
	margin-bottom: 30px;
}

.metadata h1 {
	margin: 0 0 10px 0;
	color: #222;
	font-size: 2em;
}

.metadata p {
	margin: 5px 0;
	color: #666;
}

h1, h2, h3, h4, h5, h6 {
	color: #222;
	margin-top: 30px;
	margin-bottom: 15px;
}

h1 { font-size: 1.8em; }
h2 { font-size: 1.6em; }
h3 { font-size: 1.4em; }
h4 { font-size: 1.2em; }

p {
	margin-bottom: 15px;
}

pre, code {
	background: #f5f5f5;
	border: 1px solid #ddd;
	border-radius: 3px;
	font-family: 'Courier New', monospace;
}

pre {
	padding: 15px;
	overflow-x: auto;
	margin: 20px 0;
}

code {
	padding: 2px 5px;
}

blockquote {
	border-left: 4px solid #ddd;
	padding-left: 20px;
	margin: 20px 0;
	color: #666;
	font-style: italic;
}

ul, ol {
	padding-left: 30px;
	margin-bottom: 15px;
}

li {
	margin-bottom: 5px;
}

a {
	color: #0066cc;
	text-decoration: none;
}

a:hover {
	text-decoration: underline;
}

hr {
	border: none;
	border-top: 2px solid #ddd;
	margin: 30px 0;
}

table {
	border-collapse: collapse;
	width: 100%;
	margin: 20px 0;
}

th, td {
	border: 1px solid #ddd;
	padding: 8px 12px;
	text-align: left;
}

th {
	background: #f5f5f5;
	font-weight: bold;
}

/* Lyric-specific styles */
.verse, .chorus, .bridge {
	margin: 20px 0;
	padding: 15px;
	border-radius: 5px;
}

.verse {
	background: #f9f9f9;
	border-left: 4px solid #28a745;
}

.chorus {
	background: #fff3cd;
	border-left: 4px solid #ffc107;
	font-weight: bold;
}

.bridge {
	background: #d1ecf1;
	border-left: 4px solid #17a2b8;
}
`
}

// markdownToHTML converts markdown to HTML (simplified)
func (m *ExportModel) markdownToHTML(content string) string {
	lines := strings.Split(content, "\n")
	var htmlLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Headers
		if strings.HasPrefix(line, "# ") {
			content := strings.TrimPrefix(line, "# ")
			htmlLines = append(htmlLines, fmt.Sprintf("<h1>%s</h1>", html.EscapeString(content)))
		} else if strings.HasPrefix(line, "## ") {
			content := strings.TrimPrefix(line, "## ")
			htmlLines = append(htmlLines, fmt.Sprintf("<h2>%s</h2>", html.EscapeString(content)))
		} else if strings.HasPrefix(line, "### ") {
			content := strings.TrimPrefix(line, "### ")
			htmlLines = append(htmlLines, fmt.Sprintf("<h3>%s</h3>", html.EscapeString(content)))
		} else if strings.HasPrefix(line, "#### ") {
			content := strings.TrimPrefix(line, "#### ")
			htmlLines = append(htmlLines, fmt.Sprintf("<h4>%s</h4>", html.EscapeString(content)))
		} else if strings.HasPrefix(line, "##### ") {
			content := strings.TrimPrefix(line, "##### ")
			htmlLines = append(htmlLines, fmt.Sprintf("<h5>%s</h5>", html.EscapeString(content)))
		} else if strings.HasPrefix(line, "###### ") {
			content := strings.TrimPrefix(line, "###### ")
			htmlLines = append(htmlLines, fmt.Sprintf("<h6>%s</h6>", html.EscapeString(content)))

			// Code blocks
		} else if strings.HasPrefix(line, "```") {
			htmlLines = append(htmlLines, "<pre><code>")
		} else if line == "```" && len(htmlLines) > 0 && htmlLines[len(htmlLines)-1] == "<pre><code>" {
			htmlLines = append(htmlLines, "</code></pre>")

			// Bold and italic
		} else if strings.Contains(line, "**") {
			// Simple bold conversion
			line = strings.Replace(line, "**", "<strong>", 1)
			line = strings.Replace(line, "**", "</strong>", 1)
			htmlLines = append(htmlLines, fmt.Sprintf("<p>%s</p>", line))
		} else if strings.Contains(line, "*") && !strings.HasPrefix(line, "* ") {
			// Simple italic conversion (avoid list items)
			line = strings.Replace(line, "*", "<em>", 1)
			line = strings.Replace(line, "*", "</em>", 1)
			htmlLines = append(htmlLines, fmt.Sprintf("<p>%s</p>", line))

			// Lists
		} else if strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "- ") {
			if len(htmlLines) == 0 || !strings.HasPrefix(htmlLines[len(htmlLines)-1], "<ul>") {
				htmlLines = append(htmlLines, "<ul>")
			}
			content := strings.TrimPrefix(strings.TrimPrefix(line, "* "), "- ")
			htmlLines = append(htmlLines, fmt.Sprintf("<li>%s</li>", html.EscapeString(content)))
		} else if strings.HasPrefix(line, "1. ") {
			if len(htmlLines) == 0 || !strings.HasPrefix(htmlLines[len(htmlLines)-1], "<ol>") {
				htmlLines = append(htmlLines, "<ol>")
			}
			content := strings.TrimPrefix(line, "1. ")
			htmlLines = append(htmlLines, fmt.Sprintf("<li>%s</li>", html.EscapeString(content)))

			// Empty line - close lists
		} else if line == "" {
			if len(htmlLines) > 0 {
				if strings.HasPrefix(htmlLines[len(htmlLines)-1], "<li>") {
					htmlLines = append(htmlLines, "</ul>")
				} else if strings.HasPrefix(htmlLines[len(htmlLines)-1], "<li>") {
					htmlLines = append(htmlLines, "</ol>")
				}
			}

			// Regular paragraphs
		} else if line != "" {
			htmlLines = append(htmlLines, fmt.Sprintf("<p>%s</p>", html.EscapeString(line)))
		}
	}

	// Close any open lists
	for i := len(htmlLines) - 1; i >= 0; i-- {
		if strings.HasPrefix(htmlLines[i], "<li>") {
			htmlLines = append(htmlLines, "</ul>")
			break
		}
	}

	return strings.Join(htmlLines, "\n")
}

// stripMarkdown removes markdown formatting for plain text export
func (m *ExportModel) stripMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	var cleanLines []string

	for _, line := range lines {
		// Remove markdown headers
		line = strings.TrimPrefix(line, "# ")
		line = strings.TrimPrefix(line, "## ")
		line = strings.TrimPrefix(line, "### ")
		line = strings.TrimPrefix(line, "#### ")
		line = strings.TrimPrefix(line, "##### ")
		line = strings.TrimPrefix(line, "###### ")

		// Remove markdown formatting
		line = strings.ReplaceAll(line, "**", "")
		line = strings.ReplaceAll(line, "*", "")
		line = strings.ReplaceAll(line, "_", "")
		line = strings.ReplaceAll(line, "`", "")

		// Remove list markers
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "* ")

		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// Helper functions
func countWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Focus focuses the export model
func (m *ExportModel) Focus() {
	m.focused = true
}

// Blur blurs the export model
func (m *ExportModel) Blur() {
	m.focused = false
}

// SetDimensions sets the dimensions of the export model and adapts for responsive layout
func (m *ExportModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height

	// Adapt layout based on terminal size for responsive behavior
	if width < 80 {
		// Very narrow terminal - use compact layout
		m.width = max(60, width)
	} else if width > 160 {
		// Very wide terminal - use expanded layout
		m.width = min(120, width)
	}

	if height < 20 {
		// Very short terminal - reduce content
		m.height = max(15, height)
	}
}

// GetSelectedFormat returns the currently selected export format
func (m *ExportModel) GetSelectedFormat() ExportFormat {
	return m.options[m.selected]
}

// SetContent sets the content to be exported
func (m *ExportModel) SetContent(content string) {
	m.content = content
	m.metadata.WordCount = countWords(content)
	m.metadata.CharacterCount = len(content)
}

// Messages for communication with parent
type ExportCompleteMsg struct {
	Result ExportResult
}

type BackMsg struct{}
