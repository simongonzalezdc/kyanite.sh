package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/export"
)

type ExportFormat int

const (
	ExportMarkdown ExportFormat = iota
	ExportHTML
	ExportPDF
	ExportDOCX
)

func (m Model) viewExport() string {
	if m.CurrentProject == nil {
		return "No project loaded"
	}

	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(fmt.Sprintf(" %s - Export ", m.CurrentProject.Title))
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	b.WriteString(m.Styles.Heading.Render("📤 Export Project"))
	b.WriteString("\n\n")

	b.WriteString(m.Styles.Text.Render("Select export format:\n\n"))

	formats := []struct {
		name        string
		description string
		available   bool
	}{
		{"Markdown", "Single .md file with all scenes", true},
		{"HTML", "Formatted web page", true},
		{"PDF", "Professional manuscript format", true},
		{"DOCX", "Microsoft Word document", true},
	}

	for i, format := range formats {
		prefix := "  "
		if i == m.SelectedIndex {
			prefix = "> "
		}

		style := m.Styles.Text
		if i == m.SelectedIndex {
			style = m.Styles.Accent
		}

		if !format.available {
			style = m.Styles.Text.Faint(true)
		}

		b.WriteString(style.Render(fmt.Sprintf("%s%d. %s - %s\n", prefix, i+1, format.name, format.description)))
	}

	b.WriteString("\n")
	b.WriteString(m.Styles.Text.Render("Press Enter to export, Esc to cancel\n"))

	// Footer
	if m.Message != "" {
		b.WriteString("\n")
		b.WriteString(m.Styles.Success.Render(m.Message))
	}

	if m.Error != nil {
		b.WriteString("\n")
		b.WriteString(m.Styles.Error.Render(fmt.Sprintf("Error: %v", m.Error)))
	}

	return b.String()
}

func (m Model) handleExportKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.SelectedIndex > 0 {
			m.SelectedIndex--
		}
		return m, nil

	case "down", "j":
		if m.SelectedIndex < 3 {
			m.SelectedIndex++
		}
		return m, nil

	case "enter":
		return m.performExport()

	case "esc":
		m.SelectedIndex = 0
		m.Message = ""
		m.Error = nil
		m.CurrentScreen = ScreenEditor
		return m, nil
	}

	return m, nil
}

func (m Model) performExport() (tea.Model, tea.Cmd) {
	if m.CurrentProject == nil {
		m.Error = fmt.Errorf("no project loaded")
		return m, nil
	}

	var data []byte
	var err error
	var filename string
	var format string

	switch m.SelectedIndex {
	case 0: // Markdown
		data, err = export.ExportMarkdown(m.CurrentProject, m.CurrentProject.Scenes)
		format = "markdown"
		filename = fmt.Sprintf("%s.md", sanitizeFilename(m.CurrentProject.Title))

	case 1: // HTML
		data, err = export.ExportHTML(m.CurrentProject, m.CurrentProject.Scenes)
		format = "html"
		filename = fmt.Sprintf("%s.html", sanitizeFilename(m.CurrentProject.Title))

	case 2: // PDF
		data, err = export.ExportPDF(m.CurrentProject, m.CurrentProject.Scenes)
		format = "pdf"
		filename = fmt.Sprintf("%s.pdf", sanitizeFilename(m.CurrentProject.Title))

	case 3: // DOCX
		data, err = export.ExportDOCX(m.CurrentProject, m.CurrentProject.Scenes)
		format = "docx"
		filename = fmt.Sprintf("%s.docx", sanitizeFilename(m.CurrentProject.Title))

	default:
		m.Error = fmt.Errorf("invalid export format")
		return m, nil
	}

	if err != nil {
		m.Error = fmt.Errorf("export failed: %w", err)
		return m, nil
	}

	// Determine output path
	outputDir := filepath.Join(m.CurrentProject.Directory, "exports")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		m.Error = fmt.Errorf("failed to create export directory: %w", err)
		return m, nil
	}

	outputPath := filepath.Join(outputDir, filename)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		m.Error = fmt.Errorf("failed to write export file: %w", err)
		return m, nil
	}

	m.Message = fmt.Sprintf("Exported to %s as %s", format, outputPath)
	m.Error = nil

	return m, nil
}

func sanitizeFilename(input string) string {
	// Replace invalid filename characters
	replacements := map[rune]rune{
		'/':  '-',
		'\\': '-',
		':':  '-',
		'*':  '-',
		'?':  '-',
		'"':  '-',
		'<':  '-',
		'>':  '-',
		'|':  '-',
	}

	var result strings.Builder
	for _, r := range input {
		if replacement, found := replacements[r]; found {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(r)
		}
	}

	sanitized := strings.TrimSpace(result.String())
	if len(sanitized) > 200 {
		sanitized = sanitized[:200]
	}

	return sanitized
}
