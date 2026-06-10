package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// EditorMode represents the editor state
type EditorMode int

const (
	EditorModeNormal EditorMode = iota
	EditorModeInsert
	EditorModeSearch
	EditorModeReplace
)

func (m Model) viewTextEditor() string {
	if m.CurrentProject == nil || m.CurrentScene == nil {
		return "No scene loaded"
	}

	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(fmt.Sprintf(" %s - Editing: %s | %s ",
		m.CurrentProject.Title, m.CurrentScene.Name, m.CurrentTheme.Name))
	b.WriteString(titleBar)
	b.WriteString("\n")

	// Calculate dimensions for split pane
	contentHeight := m.Height - 4 // Minus title and status bars
	paneWidth := m.Width / 2

	// Editor pane (left)
	editorContent := m.renderEditorPane(paneWidth, contentHeight)

	// Preview pane (right)
	previewContent := m.renderPreviewPane(paneWidth, contentHeight)

	// Combine panes side by side
	editorLines := strings.Split(editorContent, "\n")
	previewLines := strings.Split(previewContent, "\n")

	maxLines := max(len(editorLines), len(previewLines))
	for i := 0; i < maxLines && i < contentHeight; i++ {
		var editorLine, previewLine string
		if i < len(editorLines) {
			editorLine = editorLines[i]
		}
		if i < len(previewLines) {
			previewLine = previewLines[i]
		}

		// Pad to pane width
		editorLine = padRight(editorLine, paneWidth)
		previewLine = padRight(previewLine, paneWidth)

		b.WriteString(editorLine)
		b.WriteString("│")
		b.WriteString(previewLine)
		b.WriteString("\n")
	}

	// Search input box (overlay when in search mode)
	if m.EditorMode == EditorModeSearch && m.InputMode {
		searchBox := m.Styles.Accent.Render(fmt.Sprintf(" Find: %s█ ", m.InputValue))
		b.WriteString(searchBox)
		b.WriteString("\n")
	}

	// Replace input box (overlay when in replace mode)
	if m.EditorMode == EditorModeReplace {
		if m.InputMode {
			// Entering search term or replace term
			if m.ReplaceValue == "" {
				// First input: search term
				replaceBox := m.Styles.Accent.Render(fmt.Sprintf(" Find: %s█ | Replace: ", m.InputValue))
				b.WriteString(replaceBox)
			} else {
				// Second input: replace term
				replaceBox := m.Styles.Accent.Render(fmt.Sprintf(" Find: %s | Replace: %s█ ", m.ReplaceValue, m.InputValue))
				b.WriteString(replaceBox)
			}
			b.WriteString("\n")
		} else {
			// Show replace controls
			replaceBox := m.Styles.Accent.Render(fmt.Sprintf(" Find: %s | Replace: %s | r: Replace | a: Replace All | n: Next | p: Previous | Esc: Cancel ", m.ReplaceValue, m.InputValue))
			b.WriteString(replaceBox)
			b.WriteString("\n")
		}
	}

	// Status bar
	mode := "NORMAL"
	switch m.EditorMode {
	case EditorModeInsert:
		mode = "INSERT"
	case EditorModeSearch:
		mode = "SEARCH"
	case EditorModeReplace:
		mode = "REPLACE"
	}

	line, col := m.Buffer.CursorPosition()
	wordCount := len(strings.Fields(m.Buffer.GetContent()))

	aiStatus := ""
	if m.AIClient != nil && m.AIClient.IsEnabled() {
		aiStatus = " | AI: ON"
	}

	spellStatus := ""
	if m.SpellChecker != nil && m.SpellChecker.IsEnabled() {
		spellStatus = " | Spell: ON"
	}

	// Save status indicator
	saveStatus := ""
	switch m.SaveStatus {
	case SaveStatusSaved:
		elapsed := time.Since(m.LastSaveTime)
		if elapsed < time.Minute {
			saveStatus = fmt.Sprintf(" | Saved %ds ago", int(elapsed.Seconds()))
		} else {
			saveStatus = fmt.Sprintf(" | Saved %dm ago", int(elapsed.Minutes()))
		}
	case SaveStatusSaving:
		saveStatus = " | Saving..."
	case SaveStatusUnsaved:
		if m.Buffer.IsModified() {
			saveStatus = " | Unsaved changes"
		}
	}

	// Search info
	searchInfo := ""
	if m.EditorMode == EditorModeSearch || m.EditorMode == EditorModeReplace {
		searchTerm, current, total := m.Buffer.GetSearchInfo()
		if total > 0 {
			searchInfo = fmt.Sprintf(" | Search: '%s' (%d/%d)", searchTerm, current, total)
		} else if searchTerm != "" {
			searchInfo = fmt.Sprintf(" | Search: '%s' (no matches)", searchTerm)
		} else {
			searchInfo = " | Enter search term"
		}
	}

	statusBar := m.Styles.StatusBar.Render(fmt.Sprintf(
		" %s | Line %d:%d | Words: %d%s%s%s%s | Ctrl+F: Find | Ctrl+H: Replace | Ctrl+L: Spell | Ctrl+S: Save | Esc: Exit ",
		mode, line+1, col+1, wordCount, aiStatus, spellStatus, saveStatus, searchInfo))
	b.WriteString(statusBar)

	return b.String()
}

func (m Model) renderEditorPane(width, height int) string {
	if m.Buffer == nil {
		return "Loading..."
	}

	var b strings.Builder
	b.WriteString(m.Styles.Heading.Render(" EDITOR"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", width-2))
	b.WriteString("\n")

	line, col := m.Buffer.CursorPosition()
	lines := m.Buffer.GetLines()

	// Simple viewport - show lines around cursor
	startLine := max(0, line-height/2)
	endLine := min(len(lines), startLine+height-3)

	for i := startLine; i < endLine; i++ {
		lineText := lines[i]

		// Show line numbers
		lineNum := fmt.Sprintf("%4d ", i+1)

		// Highlight current line
		if i == line {
			lineNum = m.Styles.Accent.Render(lineNum)

			// Show cursor
			if m.EditorMode == EditorModeInsert && col <= len(lineText) {
				before := lineText[:col]
				cursor := "█"
				after := ""
				if col < len(lineText) {
					after = lineText[col:]
				}
				lineText = before + m.Styles.Accent.Render(cursor) + after
			}
		} else {
			lineNum = m.Styles.Text.Faint(true).Render(lineNum)
		}

		b.WriteString(lineNum)
		b.WriteString(lineText)
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderPreviewPane(width, height int) string {
	if m.Buffer == nil {
		return "No preview"
	}

	var b strings.Builder
	b.WriteString(m.Styles.Heading.Render(" PREVIEW"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", width-2))
	b.WriteString("\n")

	// Simple markdown rendering
	content := m.Buffer.GetContent()
	lines := strings.Split(content, "\n")

	for i := 0; i < min(len(lines), height-3); i++ {
		line := lines[i]
		rendered := m.renderMarkdownLine(line)
		b.WriteString(rendered)
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderMarkdownLine(line string) string {
	// Handle headers
	if strings.HasPrefix(line, "# ") {
		return m.Styles.Heading.Bold(true).Render(strings.TrimPrefix(line, "# "))
	} else if strings.HasPrefix(line, "## ") {
		return m.Styles.Heading.Render(strings.TrimPrefix(line, "## "))
	} else if strings.HasPrefix(line, "### ") {
		return m.Styles.Accent.Render(strings.TrimPrefix(line, "### "))
	}

	// Handle lists
	if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
		content := m.renderInlineMarkdown(line[2:])
		return m.Styles.Text.Render("  • " + content)
	}

	// Handle code blocks (simple detection)
	if strings.HasPrefix(line, "```") {
		return m.Styles.Text.Faint(true).Render(line)
	}

	// Handle indented code
	if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
		return m.Styles.Text.Faint(true).Render(line)
	}

	// Handle blockquotes
	if strings.HasPrefix(line, "> ") {
		content := m.renderInlineMarkdown(line[2:])
		return m.Styles.Text.Italic(true).Render("  " + content)
	}

	// Render inline markdown (bold, italic, code, links)
	rendered := m.renderInlineMarkdown(line)
	return m.Styles.Text.Render(rendered)
}

// renderInlineMarkdown handles inline markdown formatting
func (m Model) renderInlineMarkdown(text string) string {
	var result strings.Builder
	i := 0

	for i < len(text) {
		// Check for inline code `code`
		if text[i] == '`' {
			end := strings.Index(text[i+1:], "`")
			if end != -1 {
				end += i + 1
				// Render code without further processing
				result.WriteString(text[i : end+1])
				i = end + 1
				continue
			}
		}

		// Check for links [text](url)
		if text[i] == '[' {
			closeBracket := strings.Index(text[i+1:], "]")
			if closeBracket != -1 {
				closeBracket += i + 1
				if closeBracket+1 < len(text) && text[closeBracket+1] == '(' {
					closeParen := strings.Index(text[closeBracket+2:], ")")
					if closeParen != -1 {
						// Extract link text
						linkText := text[i+1 : closeBracket]
						result.WriteString(linkText)
						i = closeBracket + closeParen + 3
						continue
					}
				}
			}
		}

		// Check for bold **text**
		if i+1 < len(text) && text[i:i+2] == "**" {
			end := strings.Index(text[i+2:], "**")
			if end != -1 {
				end += i + 2
				// Extract and render bold text
				result.WriteString(text[i : end+2])
				i = end + 2
				continue
			}
		}

		// Check for italic *text*
		if text[i] == '*' {
			end := strings.Index(text[i+1:], "*")
			if end != -1 {
				end += i + 1
				// Extract and render italic text
				result.WriteString(text[i : end+1])
				i = end + 1
				continue
			}
		}

		// Regular character
		result.WriteByte(text[i])
		i++
	}

	return result.String()
}

func (m Model) handleTextEditorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle search mode
	if m.EditorMode == EditorModeSearch {
		switch msg.String() {
		case "esc":
			m.EditorMode = EditorModeNormal
			m.Buffer.ClearSearch()
			m.ResetInputState()
			return m, nil
		case "enter":
			// Perform search
			if m.InputValue != "" {
				count := m.Buffer.Find(m.InputValue, false)
				if count == 0 {
					m.Message = "No matches found"
				} else {
					m.Message = fmt.Sprintf("Found %d matches", count)
				}
			}
			m.InputMode = false
			return m, nil
		case "ctrl+n", "n":
			// Find next
			m.Buffer.FindNext()
			return m, nil
		case "ctrl+p", "p":
			// Find previous
			m.Buffer.FindPrevious()
			return m, nil
		case "backspace":
			if m.InputMode && len(m.InputValue) > 0 {
				m.InputValue = m.InputValue[:len(m.InputValue)-1]
			}
			return m, nil
		default:
			// Add character to search input
			if m.InputMode && len(msg.String()) == 1 {
				m.InputValue += msg.String()
			}
			return m, nil
		}
	}

	// Handle replace mode
	if m.EditorMode == EditorModeReplace {
		if m.InputMode {
			// Entering search or replace term
			switch msg.String() {
			case "esc":
				m.EditorMode = EditorModeNormal
				m.Buffer.ClearSearch()
				m.ResetInputState()
				return m, nil
			case "enter":
				if m.ReplaceValue == "" {
					// First enter: save search term, prompt for replace term
					m.ReplaceValue = m.InputValue
					m.InputValue = ""
				} else {
					// Second enter: perform search
					if m.ReplaceValue != "" {
						count := m.Buffer.Find(m.ReplaceValue, false)
						if count == 0 {
							m.Message = "No matches found"
						} else {
							m.Message = fmt.Sprintf("Found %d matches", count)
						}
					}
					m.InputMode = false
				}
				return m, nil
			case "backspace":
				if len(m.InputValue) > 0 {
					m.InputValue = m.InputValue[:len(m.InputValue)-1]
				}
				return m, nil
			default:
				// Add character to input
				if len(msg.String()) == 1 {
					m.InputValue += msg.String()
				}
				return m, nil
			}
		} else {
			// Replace controls active
			switch msg.String() {
			case "esc":
				m.EditorMode = EditorModeNormal
				m.Buffer.ClearSearch()
				m.ResetInputState()
				return m, nil
			case "r":
				// Replace current match
				if m.Buffer.ReplaceCurrent(m.InputValue) {
					searchTerm, current, total := m.Buffer.GetSearchInfo()
					if total > 0 {
						m.Message = fmt.Sprintf("Replaced. '%s' (%d/%d remaining)", searchTerm, current, total)
					} else {
						m.Message = "All matches replaced"
						m.EditorMode = EditorModeNormal
						m.ResetInputState()
					}
				}
				m.LastEditTime = time.Now()
				m.SaveStatus = SaveStatusUnsaved
				return m, nil
			case "a":
				// Replace all matches
				count := m.Buffer.ReplaceAll(m.ReplaceValue, m.InputValue, false)
				m.Message = fmt.Sprintf("Replaced %d occurrences", count)
				m.EditorMode = EditorModeNormal
				m.ResetInputState()
				m.LastEditTime = time.Now()
				m.SaveStatus = SaveStatusUnsaved
				return m, nil
			case "n":
				// Find next
				m.Buffer.FindNext()
				return m, nil
			case "p":
				// Find previous
				m.Buffer.FindPrevious()
				return m, nil
			}
		}
		return m, nil
	}

	// Mode switching
	if m.EditorMode == EditorModeNormal {
		switch msg.String() {
		case "i":
			m.EditorMode = EditorModeInsert
			return m, nil
		case "esc":
			// Save and exit
			m.ExitEditor(true)
			return m, nil
		case "ctrl+f":
			// Enter search mode
			m.EditorMode = EditorModeSearch
			m.InputMode = true
			m.InputValue = ""
			return m, nil
		case "ctrl+h":
			// Enter replace mode
			m.EditorMode = EditorModeReplace
			m.InputMode = true
			m.InputValue = ""
			m.ReplaceValue = ""
			return m, nil
		case "n":
			// Find next (if search active)
			m.Buffer.FindNext()
			return m, nil
		case "N":
			// Find previous (if search active)
			m.Buffer.FindPrevious()
			return m, nil
		case "ctrl+s":
			// Manual save
			m.SaveCurrentScene()
			return m, nil
		case "ctrl+z":
			m.Buffer.Undo()
			return m, nil
		case "ctrl+y":
			m.Buffer.Redo()
			return m, nil
		case "ctrl+a":
			// Show AI suggestion menu
			if m.AIClient != nil && m.AIClient.IsEnabled() {
				m.PreviousScreen = ScreenTextEditor
				m.CurrentScreen = ScreenAISuggestion
				m.SelectedIndex = 0
			} else {
				m.Message = "AI Assistant is disabled. Enable in config."
			}
			return m, nil
		case "ctrl+l":
			// Toggle spell check
			if m.SpellChecker != nil {
				m.SpellChecker.Toggle()
				if m.SpellChecker.IsEnabled() {
					m.Message = "Spell check enabled"
				} else {
					m.Message = "Spell check disabled"
				}
			}
			return m, nil
		}
	} else if m.EditorMode == EditorModeInsert {
		switch msg.String() {
		case "esc":
			m.EditorMode = EditorModeNormal
			return m, nil
		case "enter":
			m.Buffer.InsertNewline()
			m.LastEditTime = time.Now()
			m.SaveStatus = SaveStatusUnsaved
			return m, nil
		case "backspace":
			m.Buffer.DeleteChar()
			m.LastEditTime = time.Now()
			m.SaveStatus = SaveStatusUnsaved
			return m, nil
		case "up":
			m.Buffer.MoveCursorUp()
			return m, nil
		case "down":
			m.Buffer.MoveCursorDown()
			return m, nil
		case "left":
			m.Buffer.MoveCursorLeft()
			return m, nil
		case "right":
			m.Buffer.MoveCursorRight()
			return m, nil
		default:
			// Regular character input
			if len(msg.String()) == 1 {
				m.Buffer.InsertRune(rune(msg.String()[0]))
				m.LastEditTime = time.Now()
				m.SaveStatus = SaveStatusUnsaved
				return m, nil
			}
		}
	}

	return m, nil
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

func padRight(s string, width int) string {
	// Remove ANSI codes for length calculation
	visibleLen := len(stripANSI(s))
	if visibleLen >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-visibleLen)
}

func stripANSI(s string) string {
	// Simple ANSI stripping (for length calculation)
	result := strings.Builder{}
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
		} else if !inEscape {
			result.WriteRune(r)
		}
	}
	return result.String()
}
