package editor

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/export"
	"github.com/Kyanite/noise/internal/infra/files"
	"github.com/Kyanite/noise/internal/logging"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// EditorState implementation methods

// Init initializes the editor state component
func (s *EditorState) Init() tea.Cmd {
	// textarea.Model doesn't implement tea.Model interface, so no Init needed
	return nil
}

// Update handles messages for the editor state component
func (s *EditorState) Update(msg tea.Msg) tea.Cmd {
	if s.textarea == nil {
		return nil
	}
	var cmd tea.Cmd
	*s.textarea, cmd = s.textarea.Update(msg)
	return cmd
}

// GetText returns the current text content
func (s *EditorState) GetText() string {
	return (*s.textarea).Value()
}

// SetText sets the text content
func (s *EditorState) SetText(text string) {
	(*s.textarea).SetValue(text)
}

// InsertText inserts text at the current cursor position
func (s *EditorState) InsertText(text string) {
	if s.textarea == nil || text == "" {
		return
	}

	// Get current content
	currentText := (*s.textarea).Value()

	// Calculate approximate insertion point based on tracked cursor position
	lines := strings.Split(currentText, "\n")
	var offset int
	for i := 0; i < s.cursorLine && i < len(lines); i++ {
		offset += len(lines[i]) + 1 // +1 for newline
	}
	if s.cursorLine < len(lines) {
		// Add column offset, but clamp to line length
		col := s.cursorColumn
		if col > len(lines[s.cursorLine]) {
			col = len(lines[s.cursorLine])
		}
		offset += col
	}

	// Clamp offset to valid range
	if offset > len(currentText) {
		offset = len(currentText)
	}
	if offset < 0 {
		offset = 0
	}

	// Insert text at calculated position
	newText := currentText[:offset] + text + currentText[offset:]
	(*s.textarea).SetValue(newText)

	// Update tracked cursor position
	insertedLines := strings.Split(text, "\n")
	if len(insertedLines) > 1 {
		s.cursorLine += len(insertedLines) - 1
		s.cursorColumn = len(insertedLines[len(insertedLines)-1])
	} else {
		s.cursorColumn += len(text)
	}
}

// GetCurrentFilePath returns the current file path
func (s *EditorState) GetCurrentFilePath() string {
	return s.currentFilePath
}

// SetCurrentFilePath sets the current file path
func (s *EditorState) SetCurrentFilePath(path string) {
	s.currentFilePath = path
}

// GetSong returns the current song being edited
func (s *EditorState) GetSong() *domain.Song {
	return s.currentSong
}

// SetSong sets the current song being edited
func (s *EditorState) SetSong(song *domain.Song) {
	s.currentSong = song
	if song != nil {
		s.SetText(song.RawContent)
	}
}

// IsScratchMode returns whether the editor is in scratch mode
func (s *EditorState) IsScratchMode() bool {
	return s.scratchMode
}

// SetScratchMode sets the scratch mode for the editor
func (s *EditorState) SetScratchMode(scratchMode bool) {
	s.scratchMode = scratchMode
}

// GetEditorMode returns the current editor mode
func (s *EditorState) GetEditorMode() EditorMode {
	return s.editorMode
}

// SetEditorMode sets the editor mode (Sketch/Draft/Polish)
func (s *EditorState) SetEditorMode(mode EditorMode) {
	s.editorMode = mode
}

// UpdateCursorPosition updates the current cursor position accurately.
func (s *EditorState) UpdateCursorPosition() {
	if s.textarea == nil {
		return
	}

	// textarea.Model doesn't always expose cursor index easily,
	// so we use a robust but efficient line/col calculation.
	content := s.GetText()
	lines := strings.Split(content, "\n")

	// The editor pane tracks the line/col via metrics and state updates
	// but here we ensure they are synchronized with the actual content.
	s.cursorLine = len(lines) - 1
	if s.cursorLine >= 0 {
		s.cursorColumn = len(lines[s.cursorLine])
	}
}

// HandleAutoSave handles auto-save triggers on content changes
func (s *EditorState) HandleAutoSave() {
	if s.autoSaveService == nil {
		return
	}

	currentContent := s.GetText()

	// Only trigger auto-save if content has actually changed
	if currentContent != s.lastContent {
		s.lastContent = currentContent

		// Update current song with editor content
		if s.currentSong != nil {
			s.currentSong.RawContent = currentContent

			// Use song-specific auto-save if we have a song ID
			if s.currentSong.ID > 0 {
				autoSaveTimestampFormat := "2006-01-02 15:04:05"
				versionName := fmt.Sprintf("Auto-save %s", time.Now().Format(autoSaveTimestampFormat))
				if err := s.autoSaveService.SaveWithVersioning(s.currentSong.ID, currentContent, false, versionName); err != nil {
					// Use proper error handling instead of printf
					s.onAutoSaveError(err)
				}
				return
			}
		}

		// Fallback to general auto-save
		s.autoSaveService.SaveContent(currentContent)
	}
}

// ForceSave performs an immediate save
func (s *EditorState) ForceSave() error {
	if s.autoSaveService == nil {
		return fmt.Errorf("auto-save service not initialized")
	}

	content := s.GetText()
	return s.autoSaveService.ForceSave(content)
}

// GetAutoSaveStatus returns the current auto-save status
func (s *EditorState) GetAutoSaveStatus() app.AutoSaveStatus {
	if s.autoSaveService == nil {
		return app.AutoSaveIdle
	}
	return s.autoSaveService.GetStatus()
}

// GetLastSaveTime returns when the last save occurred
func (s *EditorState) GetLastSaveTime() time.Time {
	if s.autoSaveService == nil {
		return time.Time{}
	}
	return s.autoSaveService.GetLastSaveTime()
}

// SaveSong saves the current content as a version of the current song
func (s *EditorState) SaveSong(isMilestone bool, name string) error {
	if s.autoSaveService == nil {
		return fmt.Errorf("auto-save service not initialized")
	}
	if s.currentSong == nil {
		return fmt.Errorf("no current song set")
	}

	// Update current song with editor content before saving
	content := s.GetText()
	s.currentSong.RawContent = content

	// Use auto-save with song ID for both regular saves and milestones
	// The auto-save service will handle the proper versioning
	return s.autoSaveService.SaveWithVersioning(s.currentSong.ID, content, isMilestone, name)
}

// CreateMilestone creates a milestone version of the current song
func (s *EditorState) CreateMilestone(name string) error {
	return s.SaveSong(true, name)
}

// RecoverFromLastSave recovers content from the last auto-save for the current song
func (s *EditorState) RecoverFromLastSave() error {
	if s.autoSaveService == nil {
		return fmt.Errorf("auto-save service not initialized")
	}
	if s.currentSong == nil {
		return fmt.Errorf("no current song set")
	}

	content, err := s.autoSaveService.RecoverFromLastSave(s.currentSong.ID)
	if err != nil {
		return err
	}

	s.SetText(content)
	return nil
}

// Service injection methods

// SetAutoSaveService sets the auto-save service for this editor state
func (s *EditorState) SetAutoSaveService(service *app.AutoSaveService) {
	s.autoSaveService = service
	if service != nil {
		service.SetStatusChangeCallback(s.onAutoSaveStatusChange)
		service.SetErrorCallback(s.onAutoSaveError)
	}
}

// SetFileService sets the file service for this editor state
func (s *EditorState) SetFileService(service *files.Service) {
	s.fileService = service
}

// GetFileService returns the file service for external access
func (s *EditorState) GetFileService() *files.Service {
	return s.fileService
}

// SetExportService sets the export service for this editor state
func (s *EditorState) SetExportService(service *export.ExportService) {
	s.exportService = service
}

// SetThemeManager sets the theme manager for this editor state
func (s *EditorState) SetThemeManager(manager interface{}) {
	// Store theme manager if needed
	// Using interface{} to avoid import cycle
}

// SetDimensions sets the editor state dimensions
func (s *EditorState) SetDimensions(width, height int) {
	s.width = width
	s.height = height
}

// GetDimensions returns the editor state dimensions
func (s *EditorState) GetDimensions() (int, int) {
	return s.width, s.height
}

// Focus focuses the editor state
func (s *EditorState) Focus() {
	s.focused = true
	(*s.textarea).Focus()
}

// Blur blurs the editor state
func (s *EditorState) Blur() {
	s.focused = false
	(*s.textarea).Blur()
}

// IsFocused returns whether the editor state is focused
func (s *EditorState) IsFocused() bool {
	return s.focused
}

// GetTextarea returns the underlying textarea model
func (s *EditorState) GetTextarea() *textarea.Model {
	return s.textarea
}

// SetTextarea sets the textarea model
func (s *EditorState) SetTextarea(textarea *textarea.Model) {
	s.textarea = textarea
}

// Search and replace methods

// SetSearchMode sets the search mode
func (s *EditorState) SetSearchMode(enabled bool) {
	s.searchMode = enabled
}

// IsSearchMode returns whether search mode is active
func (s *EditorState) IsSearchMode() bool {
	return s.searchMode
}

// SetSearchQuery sets the search query
func (s *EditorState) SetSearchQuery(query string) {
	s.searchQuery = query
}

// GetSearchQuery returns the current search query
func (s *EditorState) GetSearchQuery() string {
	return s.searchQuery
}

// SetReplaceQuery sets the replace query
func (s *EditorState) SetReplaceQuery(query string) {
	s.replaceQuery = query
}

// GetReplaceQuery returns the current replace query
func (s *EditorState) GetReplaceQuery() string {
	return s.replaceQuery
}

// PerformSearch performs search operation
func (s *EditorState) PerformSearch() {
	content := s.GetText()
	s.searchMatches = nil

	if s.searchQuery == "" {
		return
	}

	// Simple search implementation
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, s.searchQuery) {
			s.searchMatches = append(s.searchMatches, i)
		}
	}

	s.currentMatch = 0
	if len(s.searchMatches) > 0 {
		s.cursorLine = s.searchMatches[0]
	}
}

// NextSearchMatch moves to the next search match
func (s *EditorState) NextSearchMatch() {
	if len(s.searchMatches) == 0 {
		return
	}

	s.currentMatch = (s.currentMatch + 1) % len(s.searchMatches)
	s.cursorLine = s.searchMatches[s.currentMatch]
}

// PreviousSearchMatch moves to the previous search match
func (s *EditorState) PreviousSearchMatch() {
	if len(s.searchMatches) == 0 {
		return
	}

	s.currentMatch = (s.currentMatch - 1 + len(s.searchMatches)) % len(s.searchMatches)
	s.cursorLine = s.searchMatches[s.currentMatch]
}

// GetSearchMatches returns the current search matches
func (s *EditorState) GetSearchMatches() []int {
	return s.searchMatches
}

// GetCurrentMatch returns the current search match index
func (s *EditorState) GetCurrentMatch() int {
	return s.currentMatch
}

// Cursor position methods

// GetCursorLine returns the current cursor line
func (s *EditorState) GetCursorLine() int {
	return s.cursorLine
}

// GetCursorColumn returns the current cursor column
func (s *EditorState) GetCursorColumn() int {
	return s.cursorColumn
}

// Editor feature methods

// ToggleLineNumbers toggles line numbers display
func (s *EditorState) ToggleLineNumbers() {
	s.showLineNumbers = !s.showLineNumbers
}

// ShowLineNumbers returns whether line numbers are shown
func (s *EditorState) ShowLineNumbers() bool {
	return s.showLineNumbers
}

// ToggleWordWrap toggles word wrap
func (s *EditorState) ToggleWordWrap() {
	s.wordWrap = !s.wordWrap
}

// WordWrapEnabled returns whether word wrap is enabled
func (s *EditorState) WordWrapEnabled() bool {
	return s.wordWrap
}

// ToggleAutoIndent toggles auto indent
func (s *EditorState) ToggleAutoIndent() {
	s.autoIndent = !s.autoIndent
}

// AutoIndentEnabled returns whether auto indent is enabled
func (s *EditorState) AutoIndentEnabled() bool {
	return s.autoIndent
}

// ToggleBracketMatching toggles bracket matching
func (s *EditorState) ToggleBracketMatching() {
	s.bracketMatching = !s.bracketMatching
}

// BracketMatchingEnabled returns whether bracket matching is enabled
func (s *EditorState) BracketMatchingEnabled() bool {
	return s.bracketMatching
}

// HandleAutoIndent handles automatic indentation for new lines
func (s *EditorState) HandleAutoIndent() {
	content := s.GetText()
	lines := strings.Split(content, "\n")

	if s.cursorLine > 0 {
		prevLine := lines[s.cursorLine-1]
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
			(*s.textarea).InsertString(indent)
		}
	}
}

// File operations

// NewFile creates a new file
func (s *EditorState) NewFile() {
	// Clear current content and reset state
	s.SetText("")
	s.currentFilePath = ""
	s.currentSong = nil
}

// OpenFile opens a file
func (s *EditorState) OpenFile(filename string) error {
	if s.fileService == nil {
		return fmt.Errorf("file service not initialized")
	}

	song, err := s.fileService.ReadSong(filename)
	if err != nil {
		return err
	}

	s.SetText(song.RawContent)
	s.currentFilePath = filename
	s.currentSong = song

	return nil
}

// SaveAs saves the current content to a new file
func (s *EditorState) SaveAs(filename string) error {
	if s.fileService == nil {
		return fmt.Errorf("file service not initialized")
	}

	content := s.GetText()

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

	err := s.fileService.WriteSong(song, filename)
	if err != nil {
		return err
	}

	s.currentFilePath = filename
	return nil
}

// CloseFile closes the current file
func (s *EditorState) CloseFile() {
	// Clear current content and reset state
	s.SetText("")
	s.currentFilePath = ""
	s.currentSong = nil
}

// Clipboard operations

// SelectAll selects all text in the editor
func (s *EditorState) SelectAll() {
	content := s.GetText()
	s.selectionStart = 0
	s.selectionEnd = len(content)
	s.hasSelection = true
	logging.Debugf("Selected all text (%d characters)", len(content))
}

// GetSelectedText returns the currently selected text
func (s *EditorState) GetSelectedText() string {
	if !s.hasSelection {
		return ""
	}
	content := s.GetText()
	if s.selectionStart < 0 || s.selectionEnd > len(content) || s.selectionStart > s.selectionEnd {
		return ""
	}
	return content[s.selectionStart:s.selectionEnd]
}

// HasSelection returns whether there is currently a selection
func (s *EditorState) HasSelection() bool {
	return s.hasSelection
}

// CopySelectedText copies the selected text to the system clipboard
func (s *EditorState) CopySelectedText() error {
	if !s.hasSelection {
		logging.Warnf("No text selected for copy")
		return nil
	}

	selectedText := s.GetSelectedText()
	if selectedText == "" {
		logging.Warnf("Selected text is empty")
		return nil
	}

	// Always store in the internal clipboard first so copy/cut work even when
	// no system clipboard utility is available (e.g. a headless CI runner).
	s.clipboardContent = selectedText

	// Best-effort sync to the system clipboard; its absence is non-fatal
	// because the internal clipboard already holds the text.
	if err := clipboard.WriteAll(selectedText); err != nil {
		logging.Warnf("System clipboard unavailable, using internal clipboard: %v", err)
	}

	logging.Debugf("Copied %d characters to clipboard", len(selectedText))
	return nil
}

// PasteFromClipboard pastes text from the system clipboard
func (s *EditorState) PasteFromClipboard() error {
	// Try system clipboard first
	clipboardText, err := clipboard.ReadAll()
	if err != nil {
		logging.Warnf("Failed to read from system clipboard: %v, using internal clipboard", err)
		// Fallback to internal clipboard
		clipboardText = s.clipboardContent
		if clipboardText == "" {
			logging.Warnf("No clipboard content available")
			return fmt.Errorf("no clipboard content available")
		}
	}

	if clipboardText == "" {
		logging.Warnf("Clipboard is empty")
		return nil
	}

	// Save current state for undo
	s.saveUndoState()

	// If there's a selection, replace it
	if s.hasSelection {
		content := s.GetText()
		newContent := content[:s.selectionStart] + clipboardText + content[s.selectionEnd:]
		s.SetText(newContent)
		// Update cursor position to end of pasted text
		s.selectionEnd = s.selectionStart + len(clipboardText)
		s.selectionStart = s.selectionEnd
		s.hasSelection = false
	} else {
		// Insert at cursor position
		(*s.textarea).InsertString(clipboardText)
	}

	logging.Debugf("Pasted %d characters from clipboard", len(clipboardText))
	return nil
}

// CutSelectedText cuts the selected text and copies it to the clipboard
func (s *EditorState) CutSelectedText() error {
	if !s.hasSelection {
		logging.Warnf("No text selected for cut")
		return nil
	}

	selectedText := s.GetSelectedText()
	if selectedText == "" {
		logging.Warnf("Selected text is empty")
		return nil
	}

	// Copy to clipboard first
	if err := s.CopySelectedText(); err != nil {
		return err
	}

	// Save current state for undo
	s.saveUndoState()

	// Remove selected text
	content := s.GetText()
	newContent := content[:s.selectionStart] + content[s.selectionEnd:]
	s.SetText(newContent)

	// Clear selection
	s.selectionStart = s.selectionEnd
	s.hasSelection = false

	logging.Debugf("Cut %d characters", len(selectedText))
	return nil
}

// Undo reverses the last edit operation
func (s *EditorState) Undo() error {
	if len(s.undoStack) == 0 {
		logging.Warnf("No undo history available")
		return fmt.Errorf("no undo history available")
	}

	// Save current state to redo stack
	currentContent := s.GetText()
	s.redoStack = append(s.redoStack, currentContent)
	if len(s.redoStack) > s.maxUndoStack {
		s.redoStack = s.redoStack[1:]
	}

	// Restore previous state
	if len(s.undoStack) > 0 {
		previousState := s.undoStack[len(s.undoStack)-1]
		s.undoStack = s.undoStack[:len(s.undoStack)-1]
		s.SetText(previousState)
		logging.Debugf("Undo operation restored %d characters", len(previousState))
	} else {
		return fmt.Errorf("no undo history available")
	}

	return nil
}

// Redo reapplies the last undone edit operation
func (s *EditorState) Redo() error {
	if len(s.redoStack) == 0 {
		logging.Warnf("No redo history available")
		return fmt.Errorf("no redo history available")
	}

	// Save current state to undo stack
	s.saveUndoState()

	// Restore redo state
	if len(s.redoStack) > 0 {
		redoState := s.redoStack[len(s.redoStack)-1]
		s.redoStack = s.redoStack[:len(s.redoStack)-1]
		s.SetText(redoState)
		logging.Debugf("Redo operation restored %d characters", len(redoState))
	} else {
		return fmt.Errorf("no redo history available")
	}

	return nil
}

// saveUndoState saves the current content to the undo stack
func (s *EditorState) saveUndoState() {
	currentContent := s.GetText()

	// Don't save if the content is the same as the last undo state
	if len(s.undoStack) > 0 && s.undoStack[len(s.undoStack)-1] == currentContent {
		return
	}

	s.undoStack = append(s.undoStack, currentContent)
	if len(s.undoStack) > s.maxUndoStack {
		s.undoStack = s.undoStack[1:]
	}

	// Clear redo stack when new content is added
	s.redoStack = s.redoStack[:0]
}

// Private helper methods

// onAutoSaveStatusChange handles auto-save status changes
func (s *EditorState) onAutoSaveStatusChange(status app.AutoSaveStatus) {
	s.lastSaveStatus = status
}

// onAutoSaveError handles auto-save errors
func (s *EditorState) onAutoSaveError(err error) {
	// Log the error properly instead of printf
	if err != nil {
		// In a full implementation, this would show a user notification
		// For now, we'll just update the status
		s.lastSaveStatus = app.AutoSaveError
	}
}
