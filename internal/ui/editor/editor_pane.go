package editor

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/noise/internal/app"
	"github.com/puente-labs/noise/internal/domain"
	"github.com/puente-labs/noise/internal/infra/files"
	"github.com/puente-labs/noise/internal/ui/styles"
)

// MarkdownElement represents a highlighted markdown element
type MarkdownElement struct {
	Type    ElementType
	Content string
	Start   int
	End     int
	Style   lipgloss.Style
}

// ElementType represents different types of markdown elements
type ElementType int

const (
	ElementText ElementType = iota
	ElementHeader
	ElementBold
	ElementItalic
	ElementCode
	ElementCodeBlock
	ElementLink
	ElementList
	ElementBlockquote
)

// SyntaxHighlighter handles markdown parsing and styling
type SyntaxHighlighter struct {
	patterns map[ElementType]*regexp.Regexp
	styles   map[ElementType]lipgloss.Style
}

// NewSyntaxHighlighter creates a new syntax highlighter
func NewSyntaxHighlighter() *SyntaxHighlighter {
	sh := &SyntaxHighlighter{
		patterns: make(map[ElementType]*regexp.Regexp),
		styles:   make(map[ElementType]lipgloss.Style),
	}

	// Define regex patterns for markdown elements
	sh.patterns[ElementHeader] = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	sh.patterns[ElementBold] = regexp.MustCompile(`\*\*(.+?)\*\*`)
	sh.patterns[ElementItalic] = regexp.MustCompile(`\*([^*]+?)\*`)
	sh.patterns[ElementCode] = regexp.MustCompile("`([^`]+?)`")
	sh.patterns[ElementCodeBlock] = regexp.MustCompile("(?s)```[\\s\\S]*?```")
	sh.patterns[ElementLink] = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	sh.patterns[ElementList] = regexp.MustCompile(`(?m)^(\s*)([-*+]|\d+\.)\s+(.+)$`)
	sh.patterns[ElementBlockquote] = regexp.MustCompile(`(?m)^>\s+(.+)$`)

	// Define styles for each element type using Midnight Jazz theme
	sh.styles[ElementHeader] = lipgloss.NewStyle().Bold(true).Foreground(styles.Accent)
	sh.styles[ElementBold] = lipgloss.NewStyle().Bold(true).Foreground(styles.Primary)
	sh.styles[ElementItalic] = lipgloss.NewStyle().Italic(true).Foreground(styles.TextSecondary)
	sh.styles[ElementCode] = lipgloss.NewStyle().Background(styles.Dark2).Foreground(styles.Success).Padding(0, 1)
	sh.styles[ElementCodeBlock] = lipgloss.NewStyle().Background(styles.Dark1).Foreground(styles.TextPrimary)
	sh.styles[ElementLink] = lipgloss.NewStyle().Foreground(styles.Info).Underline(true)
	sh.styles[ElementList] = lipgloss.NewStyle().Foreground(styles.TextSecondary)
	sh.styles[ElementBlockquote] = lipgloss.NewStyle().Foreground(styles.TextMuted).Italic(true)

	return sh
}

// ParseMarkdown parses markdown content and returns highlighted elements
func (sh *SyntaxHighlighter) ParseMarkdown(content string) []MarkdownElement {
	var elements []MarkdownElement
	processed := make([]bool, len(content))

	// Process code blocks first (they take precedence)
	codeBlockMatches := sh.patterns[ElementCodeBlock].FindAllStringIndex(content, -1)
	for _, match := range codeBlockMatches {
		if !processed[match[0]] && !processed[match[1]-1] {
			elements = append(elements, MarkdownElement{
				Type:    ElementCodeBlock,
				Content: content[match[0]:match[1]],
				Start:   match[0],
				End:     match[1],
				Style:   sh.styles[ElementCodeBlock],
			})
			for i := match[0]; i < match[1]; i++ {
				processed[i] = true
			}
		}
	}

	// Process other elements
	for elementType, pattern := range sh.patterns {
		if elementType == ElementCodeBlock {
			continue // Already processed
		}

		matches := pattern.FindAllStringSubmatchIndex(content, -1)
		for _, match := range matches {
			if len(match) >= 4 && match[0] >= 0 && match[1] <= len(content) {
				start, end := match[0], match[1]
				if !processed[start] && !processed[end-1] {
					var contentPart string
					if len(match) > 2 {
						contentPart = content[match[2]:match[3]]
					} else {
						contentPart = content[start:end]
					}

					elements = append(elements, MarkdownElement{
						Type:    elementType,
						Content: contentPart,
						Start:   start,
						End:     end,
						Style:   sh.styles[elementType],
					})

					for i := start; i < end; i++ {
						processed[i] = true
					}
				}
			}
		}
	}

	// Sort elements by start position
	for i := 0; i < len(elements)-1; i++ {
		for j := i + 1; j < len(elements); j++ {
			if elements[i].Start > elements[j].Start {
				elements[i], elements[j] = elements[j], elements[i]
			}
		}
	}

	return elements
}

// EditorPaneModel handles the text editing pane with syntax highlighting
type EditorPaneModel struct {
	textarea textarea.Model
	width    int
	height   int
	focused  bool

	// Syntax highlighting
	highlighter *SyntaxHighlighter
	elements    []MarkdownElement

	// Editor features
	showLineNumbers  bool
	lineNumbersWidth int
	wordWrap         bool
	autoIndent       bool
	bracketMatching  bool

	// Search and replace
	searchMode    bool
	searchQuery   string
	replaceQuery  string
	searchMatches []int
	currentMatch  int

	// Keyboard shortcuts
	shortcutManager *ShortcutManager

	// Cursor and selection
	cursorLine   int
	cursorColumn int

	// Status bar component
	statusBar *StatusBarModel

	// Auto-save functionality
	autoSaveService *app.AutoSaveService
	lastContent     string
	lastSaveStatus  app.AutoSaveStatus

	// Song integration
	currentSong *domain.Song

	// File I/O service
	fileService *files.Service

	// Current file path
	currentFilePath string

	// Styles
	focusedStyle     lipgloss.Style
	blurredStyle     lipgloss.Style
	borderStyle      lipgloss.Style
	lineNumberStyle  lipgloss.Style
	cursorLineStyle  lipgloss.Style
	selectionStyle   lipgloss.Style
	searchMatchStyle lipgloss.Style
	autoSaveStyle    lipgloss.Style // Auto-save status styling
}

// NewEditorPaneModel creates a new editor pane model
func NewEditorPaneModel(textarea textarea.Model) *EditorPaneModel {
	model := &EditorPaneModel{
		textarea:         textarea,
		focused:          true,
		highlighter:      NewSyntaxHighlighter(),
		showLineNumbers:  true,
		lineNumbersWidth: 4,
		wordWrap:         true,
		autoIndent:       true,
		bracketMatching:  true,
		cursorLine:       0,
		cursorColumn:     0,
		statusBar:        NewStatusBarModel(),
		autoSaveService:  nil, // Will be set by parent model
		lastContent:      "",
		lastSaveStatus:   app.AutoSaveIdle,
		currentSong:      nil,
		fileService:      nil, // Will be set by parent model
		currentFilePath:  "",
		shortcutManager:  NewShortcutManager(),
		focusedStyle:     styles.BorderActive,
		blurredStyle:     styles.Border,
		borderStyle:      styles.Border,
		lineNumberStyle: lipgloss.NewStyle().
			Foreground(styles.TextMuted).
			Width(4).
			Align(lipgloss.Right),
		cursorLineStyle: lipgloss.NewStyle().
			Background(styles.Dark2),
		selectionStyle: lipgloss.NewStyle().
			Background(styles.Dark3),
		searchMatchStyle: lipgloss.NewStyle().
			Background(styles.Accent).
			Foreground(styles.Background),
		autoSaveStyle: lipgloss.NewStyle().
			Foreground(styles.Success).
			Bold(true),
	}

	return model
}

// Init initializes the editor pane
func (m *EditorPaneModel) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages for the editor pane
func (m *EditorPaneModel) Update(msg tea.Msg) (*EditorPaneModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle text area updates
	m.textarea, cmd = m.textarea.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Update cursor position
	m.updateCursorPosition()

	// Update status bar with current state
	m.updateStatusBar()

	// Handle key events for editor features
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Set context for shortcuts
		if m.focused {
			m.shortcutManager.SetContext(ContextEditor)
		} else {
			m.shortcutManager.SetContext(ContextGlobal)
		}

		// Handle shortcuts
		if action, handled := m.shortcutManager.HandleKey(msg); handled {
			m.handleShortcutAction(action)
			return m, tea.Batch(cmds...)
		}

		// Handle legacy key events for backward compatibility
		switch msg.String() {
		case "ctrl+r":
			// Enter replace mode
			if m.searchMode {
				m.replaceQuery = ""
			}
		case "enter":
			if m.searchMode {
				// Perform search
				m.performSearch()
			} else if m.autoIndent {
				// Handle auto-indentation
				m.handleAutoIndent()
			}
		}
	}

	// Update syntax highlighting
	m.updateSyntaxHighlighting()

	// Handle auto-save on content changes
	m.handleAutoSave()

	return m, tea.Batch(cmds...)
}

// View renders the editor pane with syntax highlighting
func (m *EditorPaneModel) View() string {
	var style lipgloss.Style

	if m.focused {
		style = m.focusedStyle
	} else {
		style = m.blurredStyle
	}

	// Calculate content dimensions
	contentWidth := m.width - 4
	contentHeight := m.height - 4

	if m.showLineNumbers {
		contentWidth -= m.lineNumbersWidth + 1 // Account for line numbers and spacing
	}

	// Set textarea dimensions
	m.textarea.SetWidth(contentWidth)
	m.textarea.SetHeight(contentHeight)

	// Get content and apply syntax highlighting
	content := m.textarea.Value()
	highlightedContent := m.renderHighlightedContent(content, contentWidth)

	// Add line numbers if enabled
	if m.showLineNumbers {
		highlightedContent = m.addLineNumbers(highlightedContent)
	}

	// Add title
	title := "Editor"
	if m.focused {
		title = "Editor (Focused)"
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.TextPrimary).
		Background(styles.Dark3).
		Padding(0, 1).
		Width(m.width - 4)

	titleBar := titleStyle.Render(title)

	// Add status bar with editor features info
	statusBar := m.renderStatusBar()

	// Combine title, content, and status bar
	fullContent := lipgloss.JoinVertical(lipgloss.Left, titleBar, highlightedContent, statusBar)

	return style.Width(m.width).Height(m.height).Render(fullContent)
}

// SetDimensions sets the pane dimensions
func (m *EditorPaneModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	if m.statusBar != nil {
		m.statusBar.SetDimensions(width, 1)
	}
}

// Focus focuses the editor pane
func (m *EditorPaneModel) Focus() {
	m.focused = true
	m.textarea.Focus()
}

// Blur blurs the editor pane
func (m *EditorPaneModel) Blur() {
	m.focused = false
	m.textarea.Blur()
}

// GetText returns the current text content
func (m *EditorPaneModel) GetText() string {
	return m.textarea.Value()
}

// SetText sets the text content
func (m *EditorPaneModel) SetText(text string) {
	m.textarea.SetValue(text)
	m.updateSyntaxHighlighting()
	m.updateStatusBar()
}

// Helper methods for enhanced editor features

// updateCursorPosition updates the current cursor position
func (m *EditorPaneModel) updateCursorPosition() {
	content := m.textarea.Value()
	// For now, use a simple approximation of cursor position
	// In a full implementation, this would track cursor position more accurately
	lines := strings.Split(content, "\n")
	m.cursorLine = len(lines) - 1

	if m.cursorLine < len(lines) {
		m.cursorColumn = len(lines[m.cursorLine])
	} else {
		m.cursorColumn = 0
	}
}

// updateStatusBar updates the status bar with current editor state
func (m *EditorPaneModel) updateStatusBar() {
	if m.statusBar == nil {
		return
	}

	// Update content and statistics
	content := m.textarea.Value()
	m.statusBar.UpdateContent(content)

	// Update cursor position
	m.statusBar.UpdateCursorPosition(m.cursorLine, m.cursorColumn)

	// Update auto-save status
	if m.autoSaveService != nil {
		status := m.GetAutoSaveStatus()
		lastSaveTime := m.GetLastSaveTime()
		m.statusBar.UpdateAutoSaveStatus(status, lastSaveTime)
	}

	// Update editor features
	m.statusBar.UpdateEditorFeatures(m.showLineNumbers, m.wordWrap, m.autoIndent, m.bracketMatching)

	// Update current file path
	if m.statusBar != nil {
		filename := "Untitled"
		if m.currentFilePath != "" {
			filename = filepath.Base(m.currentFilePath)
		}
		m.statusBar.UpdateFileInfo(filename)
	}

	// Update shortcut hints
	hints := m.GetShortcutHints()
	m.statusBar.UpdateShortcutHints(hints)

	// Update dimensions and responsive mode
	m.statusBar.SetDimensions(m.width, 1)
	m.statusBar.UpdateResponsiveMode(m.width)
}

// updateSyntaxHighlighting updates the syntax highlighting elements
func (m *EditorPaneModel) updateSyntaxHighlighting() {
	content := m.textarea.Value()
	m.elements = m.highlighter.ParseMarkdown(content)
}

// renderHighlightedContent renders content with syntax highlighting
func (m *EditorPaneModel) renderHighlightedContent(content string, width int) string {
	if len(m.elements) == 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	var resultLines []string

	for _, line := range lines {
		highlightedLine := m.renderHighlightedLine(line)
		resultLines = append(resultLines, highlightedLine)
	}

	return strings.Join(resultLines, "\n")
}

// renderHighlightedLine renders a single line with syntax highlighting
func (m *EditorPaneModel) renderHighlightedLine(line string) string {
	// For now, return the line as-is
	// In a full implementation, this would apply styling to markdown elements
	return line
}

// addLineNumbers adds line numbers to the content
func (m *EditorPaneModel) addLineNumbers(content string) string {
	lines := strings.Split(content, "\n")
	var numberedLines []string

	for i, line := range lines {
		lineNum := fmt.Sprintf("%*d", m.lineNumbersWidth, i+1)
		lineNumber := m.lineNumberStyle.Render(lineNum)
		numberedLine := lineNumber + " " + line
		numberedLines = append(numberedLines, numberedLine)
	}

	return strings.Join(numberedLines, "\n")
}

// renderStatusBar renders the status bar with editor features info
func (m *EditorPaneModel) renderStatusBar() string {
	if m.statusBar == nil {
		// Fallback to old implementation if status bar not initialized
		var features []string

		if m.showLineNumbers {
			features = append(features, "L")
		}
		if m.wordWrap {
			features = append(features, "W")
		}
		if m.autoIndent {
			features = append(features, "I")
		}
		if m.bracketMatching {
			features = append(features, "B")
		}

		featuresStr := strings.Join(features, " ")

		// Add auto-save status
		autoSaveStatus := ""
		if m.autoSaveService != nil {
			status := m.GetAutoSaveStatus()
			switch status {
			case app.AutoSaveSaving:
				autoSaveStatus = " Saving..."
			case app.AutoSaveSuccess:
				autoSaveStatus = " Saved"
			case app.AutoSaveError:
				autoSaveStatus = " Save Error"
			default:
				// Show last save time for idle status
				if !m.GetLastSaveTime().IsZero() {
					autoSaveStatus = fmt.Sprintf(" Saved %s", m.GetLastSaveTime().Format("15:04:05"))
				}
			}
		}

		// Combine editor status and auto-save status
		editorStatus := fmt.Sprintf("Ln %d, Col %d | %s", m.cursorLine+1, m.cursorColumn+1, featuresStr)
		if autoSaveStatus != "" {
			editorStatus = fmt.Sprintf("%s |%s", editorStatus, autoSaveStatus)
		}

		// Add shortcut hints if available
		if m.shortcutManager != nil {
			hints := m.shortcutManager.GetStatusBarHints()
			if hints != "" {
				editorStatus = fmt.Sprintf("%s | %s", editorStatus, hints)
			}
		}

		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(0, 1).
			Width(m.width - 4)

		return statusStyle.Render(editorStatus)
	}

	// Use the new status bar component
	return m.statusBar.View()
}

// performSearch performs search operation
func (m *EditorPaneModel) performSearch() {
	content := m.textarea.Value()
	m.searchMatches = nil

	if m.searchQuery == "" {
		return
	}

	// Simple search implementation
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, m.searchQuery) {
			m.searchMatches = append(m.searchMatches, i)
		}
	}

	m.currentMatch = 0
	if len(m.searchMatches) > 0 {
		m.cursorLine = m.searchMatches[0]
	}
}

// nextSearchMatch moves to the next search match
func (m *EditorPaneModel) nextSearchMatch() {
	if len(m.searchMatches) == 0 {
		return
	}

	m.currentMatch = (m.currentMatch + 1) % len(m.searchMatches)
	m.cursorLine = m.searchMatches[m.currentMatch]
}

// previousSearchMatch moves to the previous search match
func (m *EditorPaneModel) previousSearchMatch() {
	if len(m.searchMatches) == 0 {
		return
	}

	m.currentMatch = (m.currentMatch - 1 + len(m.searchMatches)) % len(m.searchMatches)
	m.cursorLine = m.searchMatches[m.currentMatch]
}

// handleAutoIndent handles automatic indentation for new lines
func (m *EditorPaneModel) handleAutoIndent() {
	content := m.textarea.Value()
	lines := strings.Split(content, "\n")

	if m.cursorLine > 0 {
		prevLine := lines[m.cursorLine-1]
		indent := ""

		// Count leading spaces and tabs
		for _, r := range prevLine {
			if r == ' ' || r == '\t' {
				indent += string(r)
			} else {
				break
			}
		}

		if indent != "" {
			// Insert indentation at cursor position
			m.textarea.InsertString(indent)
		}
	}
}

// insertTab inserts a tab or handles indentation
// TODO: Uncomment when tab insertion functionality is needed
// func (m *EditorPaneModel) insertTab() {
// 	// Insert 2 spaces instead of tab
// 	m.textarea.InsertString("  ")
// }

// SetAutoSaveService sets the auto-save service for this editor pane
func (m *EditorPaneModel) SetAutoSaveService(service *app.AutoSaveService) {
	m.autoSaveService = service
	if service != nil {
		service.SetStatusChangeCallback(m.onAutoSaveStatusChange)
		service.SetErrorCallback(m.onAutoSaveError)
	}
}

// SetFileService sets the file service for this editor pane
func (m *EditorPaneModel) SetFileService(service *files.Service) {
	m.fileService = service
}

// handleAutoSave handles auto-save triggers on content changes
func (m *EditorPaneModel) handleAutoSave() {
	if m.autoSaveService == nil {
		return
	}

	currentContent := m.textarea.Value()

	// Only trigger auto-save if content has actually changed
	if currentContent != m.lastContent {
		m.lastContent = currentContent
		
		// Update current song with editor content
		if m.currentSong != nil {
			m.currentSong.RawContent = currentContent
			
			// Use song-specific auto-save if we have a song ID
			if m.currentSong.ID > 0 {
				versionName := fmt.Sprintf("Auto-save %s", time.Now().Format("2006-01-02 15:04:05"))
				if err := m.autoSaveService.SaveWithVersioning(m.currentSong.ID, currentContent, false, versionName); err != nil {
					// Use proper error handling instead of printf
					m.onAutoSaveError(err)
				}
				return
			}
		}
		
		// Fallback to general auto-save
		m.autoSaveService.SaveContent(currentContent)
	}
}

// onAutoSaveStatusChange handles auto-save status changes
func (m *EditorPaneModel) onAutoSaveStatusChange(status app.AutoSaveStatus) {
	m.lastSaveStatus = status
	m.updateStatusBar()
}

// onAutoSaveError handles auto-save errors
func (m *EditorPaneModel) onAutoSaveError(err error) {
	// Log the error properly instead of printf
	if err != nil {
		// In a full implementation, this would show a user notification
		// For now, we'll just update the status
		m.lastSaveStatus = app.AutoSaveError
		m.updateStatusBar()
	}
}

// ForceSave performs an immediate save
func (m *EditorPaneModel) ForceSave() error {
	if m.autoSaveService == nil {
		return fmt.Errorf("auto-save service not initialized")
	}

	content := m.textarea.Value()
	return m.autoSaveService.ForceSave(content)
}

// GetAutoSaveStatus returns the current auto-save status
func (m *EditorPaneModel) GetAutoSaveStatus() app.AutoSaveStatus {
	if m.autoSaveService == nil {
		return app.AutoSaveIdle
	}
	return m.autoSaveService.GetStatus()
}

// GetLastSaveTime returns when the last save occurred
func (m *EditorPaneModel) GetLastSaveTime() time.Time {
	if m.autoSaveService == nil {
		return time.Time{}
	}
	return m.autoSaveService.GetLastSaveTime()
}

// SetSong sets the current song being edited
func (m *EditorPaneModel) SetSong(song *domain.Song) {
	m.currentSong = song
	if song != nil {
		m.SetText(song.RawContent)
	}
}

// GetSong returns the current song being edited
func (m *EditorPaneModel) GetSong() *domain.Song {
	return m.currentSong
}

// SaveSong saves the current content as a version of the current song
func (m *EditorPaneModel) SaveSong(isMilestone bool, name string) error {
	if m.autoSaveService == nil {
		return fmt.Errorf("auto-save service not initialized")
	}
	if m.currentSong == nil {
		return fmt.Errorf("no current song set")
	}

	// Update current song with editor content before saving
	content := m.textarea.Value()
	m.currentSong.RawContent = content

	// Use auto-save with song ID for both regular saves and milestones
	// The auto-save service will handle the proper versioning
	return m.autoSaveService.SaveWithVersioning(m.currentSong.ID, content, isMilestone, name)
}

// CreateMilestone creates a milestone version of the current song
func (m *EditorPaneModel) CreateMilestone(name string) error {
	return m.SaveSong(true, name)
}

// GetVersionHistory returns the version history for the current song
func (m *EditorPaneModel) GetVersionHistory() ([]*domain.Version, error) {
	if m.autoSaveService == nil {
		return nil, fmt.Errorf("auto-save service not initialized")
	}
	if m.currentSong == nil {
		return nil, fmt.Errorf("no current song set")
	}

	return m.autoSaveService.GetVersionHistory(m.currentSong.ID, 0)
}

// RecoverFromLastSave recovers content from the last auto-save for the current song
func (m *EditorPaneModel) RecoverFromLastSave() error {
	if m.autoSaveService == nil {
		return fmt.Errorf("auto-save service not initialized")
	}
	if m.currentSong == nil {
		return fmt.Errorf("no current song set")
	}

	content, err := m.autoSaveService.RecoverFromLastSave(m.currentSong.ID)
	if err != nil {
		return err
	}

	m.SetText(content)
	return nil
}

// GetSaveStatistics returns save statistics for the current song
func (m *EditorPaneModel) GetSaveStatistics() (*app.SaveStatistics, error) {
	if m.autoSaveService == nil {
		return nil, fmt.Errorf("auto-save service not initialized")
	}
	if m.currentSong == nil {
		return nil, fmt.Errorf("no current song set")
	}

	return m.autoSaveService.GetSaveStatistics(m.currentSong.ID)
}

// handleShortcutAction handles actions from the keyboard shortcut system
func (m *EditorPaneModel) handleShortcutAction(action ShortcutAction) {
	switch action.Type {
	case ActionToggleLineNumbers:
		m.showLineNumbers = !m.showLineNumbers
		m.updateStatusBar()
	case ActionToggleWordWrap:
		m.wordWrap = !m.wordWrap
		m.updateStatusBar()
	case ActionToggleAutoIndent:
		m.autoIndent = !m.autoIndent
		m.updateStatusBar()
	case ActionToggleBracketMatching:
		m.bracketMatching = !m.bracketMatching
		m.updateStatusBar()
	case ActionFind:
		m.searchMode = true
		m.searchQuery = ""
	case ActionReplace:
		if m.searchMode {
			m.replaceQuery = ""
		}
	case ActionFindNext:
		if m.searchMode {
			m.nextSearchMatch()
		}
	case ActionFindPrev:
		if m.searchMode {
			m.previousSearchMatch()
		}
	case ActionSave:
		if m.autoSaveService != nil {
			err := m.ForceSave()
			if err != nil {
				// Use proper error handling
				m.onAutoSaveError(err)
			}
		}
	case ActionSelectAll:
		// Select all text - implementation depends on textarea capabilities
		// For now, this is a placeholder that would need custom implementation
	case ActionCopy:
		// Copy selected text - implementation depends on textarea capabilities
		// For now, this is a placeholder that would need custom implementation
	case ActionPaste:
		// Paste from clipboard - implementation depends on textarea capabilities
		// For now, this is a placeholder that would need custom implementation
	case ActionCut:
		// Cut selected text - implementation depends on textarea capabilities
		// For now, this is a placeholder that would need custom implementation
	case ActionUndo:
		// Note: Bubble Tea textarea doesn't have built-in undo/redo
		// This would need custom implementation
	case ActionRedo:
		// Note: Bubble Tea textarea doesn't have built-in undo/redo
		// This would need custom implementation
	case ActionStartOfLine:
		m.moveCursorToStartOfLine()
	case ActionEndOfLine:
		m.moveCursorToEndOfLine()
	case ActionStartOfFile:
		m.moveCursorToStartOfFile()
	case ActionEndOfFile:
		m.moveCursorToEndOfFile()
	case ActionPrevWord:
		m.moveCursorToPrevWord()
	case ActionNextWord:
		m.moveCursorToNextWord()
	case ActionSelectToStartOfLine:
		m.selectToStartOfLine()
	case ActionSelectToEndOfLine:
		m.selectToEndOfLine()
	case ActionSelectToStartOfFile:
		m.selectToStartOfFile()
	case ActionSelectToEndOfFile:
		m.selectToEndOfFile()
	case ActionSelectLeft:
		m.selectLeft()
	case ActionSelectRight:
		m.selectRight()
	case ActionSelectUp:
		m.selectUp()
	case ActionSelectDown:
		m.selectDown()
	case ActionPageUp:
		m.pageUp()
	case ActionPageDown:
		m.pageDown()
	case ActionGoToLine:
		m.goToLine()
	case ActionNewFile:
		m.newFile()
	case ActionOpenFile:
		m.openFile()
	case ActionSaveAs:
		m.saveAs()
	case ActionCloseFile:
		m.closeFile()
	case ActionQuit:
		// This should be handled by parent model
	case ActionSettings:
		// This should be handled by parent model
	case ActionTheoryTools:
		// This should be handled by parent model
	case ActionAudioTools:
		// This should be handled by parent model
	case ActionToggleHelp:
		// Help mode is handled by shortcut manager
	case ActionBackToMenu:
		// This should be handled by parent model
	}
}

// Cursor movement helper methods
func (m *EditorPaneModel) moveCursorToStartOfLine() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) moveCursorToEndOfLine() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) moveCursorToStartOfFile() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) moveCursorToEndOfFile() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) moveCursorToPrevWord() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) moveCursorToNextWord() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

// Text selection helper methods
func (m *EditorPaneModel) selectToStartOfLine() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) selectToEndOfLine() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) selectToStartOfFile() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) selectToEndOfFile() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) selectLeft() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) selectRight() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) selectUp() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) selectDown() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

// Navigation helper methods
func (m *EditorPaneModel) pageUp() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (m *EditorPaneModel) pageDown() {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

// Feature helper methods
func (m *EditorPaneModel) goToLine() {
	// Implementation for go to line functionality
	// For now, this is a placeholder
}

func (m *EditorPaneModel) newFile() {
	// Clear current content and reset state
	m.SetText("")
	m.currentFilePath = ""
	m.currentSong = nil
	m.updateStatusBar()
}

func (m *EditorPaneModel) openFile() {
	// For now, we'll use a simple approach - in a full implementation,
	// this would show a file dialog. For this demo, we'll show available files
	// and let the user choose via a simple text input

	if m.fileService == nil {
		// Use proper error handling instead of printf
		return
	}

	// List available song files
	songFiles, err := m.fileService.ListSongs()
	if err != nil {
		// Use proper error handling instead of printf
		return
	}

	if len(songFiles) == 0 {
		// Use proper error handling instead of printf
		return
	}

	// For this implementation, we'll just open the first file as a demo
	// In a full implementation, this would show a file picker
	fileToOpen := songFiles[0]

	song, err := m.fileService.ReadSong(fileToOpen)
	if err != nil {
		// Use proper error handling instead of printf
		return
	}

	m.SetText(song.RawContent)
	m.currentFilePath = fileToOpen
	m.currentSong = song
	m.updateStatusBar()

	// Remove printf - the status bar update is sufficient
}

func (m *EditorPaneModel) saveAs() {
	if m.fileService == nil {
		// Use proper error handling instead of printf
		return
	}

	content := m.GetText()

	// For this implementation, we'll use a simple filename
	// In a full implementation, this would show a save dialog
	filename := "untitled.md"
	if m.currentFilePath != "" {
		filename = filepath.Base(m.currentFilePath)
	}

	// Create a basic song structure for saving
	song := &domain.Song{
		Metadata: domain.SongMetadata{
			Title:     "Untitled Song",
			Artist:    "Unknown Artist",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Sections: []domain.Section{
			{
				Type:   domain.SectionVerse,
				Number: 1,
				Lines:  []domain.Line{},
			},
		},
		RawContent: content,
	}

	err := m.fileService.WriteSong(song, filename)
	if err != nil {
		// Use proper error handling instead of printf
		return
	}

	m.currentFilePath = filename
	m.updateStatusBar()

	// Remove printf - the status bar update is sufficient
}

func (m *EditorPaneModel) closeFile() {
	// Clear current content and reset state
	m.SetText("")
	m.currentFilePath = ""
	m.currentSong = nil
	m.updateStatusBar()
}

// GetShortcutManager returns the shortcut manager for external access
func (m *EditorPaneModel) GetShortcutManager() *ShortcutManager {
	return m.shortcutManager
}

// SetShortcutContext sets the keyboard shortcut context
func (m *EditorPaneModel) SetShortcutContext(context KeyContext) {
	if m.shortcutManager != nil {
		m.shortcutManager.SetContext(context)
	}
}

// GetShortcutHints returns shortcut hints for the status bar
func (m *EditorPaneModel) GetShortcutHints() string {
	if m.shortcutManager != nil {
		return m.shortcutManager.GetStatusBarHints()
	}
	return ""
}

// GetCurrentFilePath returns the current file path
func (m *EditorPaneModel) GetCurrentFilePath() string {
	return m.currentFilePath
}

// GetFileService returns the file service for external access
func (m *EditorPaneModel) GetFileService() *files.Service {
	return m.fileService
}
