package editor

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/puente-labs/noise/internal/app"
	"github.com/puente-labs/noise/internal/domain"
	"github.com/puente-labs/noise/internal/export"
	"github.com/puente-labs/noise/internal/infra/files"
)

// EditorState implementation methods

// Init initializes the editor state component
func (s *EditorState) Init() tea.Cmd {
	// textarea.Model doesn't implement tea.Model interface, so no Init needed
	return nil
}

// Update handles messages for the editor state component
func (s *EditorState) Update(msg tea.Msg) tea.Cmd {
	// textarea.Model doesn't implement tea.Model interface, so no Update needed
	return nil
}

// GetText returns the current text content
func (s *EditorState) GetText() string {
	return (*s.textarea).Value()
}

// SetText sets the text content
func (s *EditorState) SetText(text string) {
	(*s.textarea).SetValue(text)
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

// UpdateCursorPosition updates the current cursor position
func (s *EditorState) UpdateCursorPosition() {
	content := s.GetText()
	// For now, use a simple approximation of cursor position
	// In a full implementation, this would track cursor position more accurately
	lines := strings.Split(content, "\n")
	s.cursorLine = len(lines) - 1

	if s.cursorLine < len(lines) {
		s.cursorColumn = len(lines[s.cursorLine])
	} else {
		s.cursorColumn = 0
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
				versionName := fmt.Sprintf("Auto-save %s", time.Now().Format("2006-01-02 15:04:05"))
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
