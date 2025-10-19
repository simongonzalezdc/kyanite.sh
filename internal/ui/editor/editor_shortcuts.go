package editor

import (
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// EditorShortcuts implementation methods

// HandleKey processes a key message and returns the corresponding action
func (s *EditorShortcuts) HandleKey(msg tea.KeyMsg) (ShortcutAction, bool) {
	return s.shortcutManager.HandleKey(msg)
}

// HandleShortcutAction executes a shortcut action
func (s *EditorShortcuts) HandleShortcutAction(action ShortcutAction, state StateManagerInterface) tea.Cmd {
	switch action.Type {
	case ActionToggleLineNumbers:
		state.ToggleLineNumbers()
	case ActionToggleWordWrap:
		state.ToggleWordWrap()
	case ActionToggleAutoIndent:
		state.ToggleAutoIndent()
	case ActionToggleBracketMatching:
		state.ToggleBracketMatching()
	case ActionFind:
		state.SetSearchMode(true)
		state.SetSearchQuery("")
	case ActionReplace:
		if state.IsSearchMode() {
			state.SetReplaceQuery("")
		}
	case ActionFindNext:
		if state.IsSearchMode() {
			state.NextSearchMatch()
		}
	case ActionFindPrev:
		if state.IsSearchMode() {
			state.PreviousSearchMatch()
		}
	case ActionSave:
		err := state.ForceSave()
		if err != nil {
			// Use proper error handling
			// In a full implementation, this would show a user notification
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
		s.moveCursorToStartOfLine(state)
	case ActionEndOfLine:
		s.moveCursorToEndOfLine(state)
	case ActionStartOfFile:
		s.moveCursorToStartOfFile(state)
	case ActionEndOfFile:
		s.moveCursorToEndOfFile(state)
	case ActionPrevWord:
		s.moveCursorToPrevWord(state)
	case ActionNextWord:
		s.moveCursorToNextWord(state)
	case ActionSelectToStartOfLine:
		s.selectToStartOfLine(state)
	case ActionSelectToEndOfLine:
		s.selectToEndOfLine(state)
	case ActionSelectToStartOfFile:
		s.selectToStartOfFile(state)
	case ActionSelectToEndOfFile:
		s.selectToEndOfFile(state)
	case ActionSelectLeft:
		s.selectLeft(state)
	case ActionSelectRight:
		s.selectRight(state)
	case ActionSelectUp:
		s.selectUp(state)
	case ActionSelectDown:
		s.selectDown(state)
	case ActionPageUp:
		s.pageUp(state)
	case ActionPageDown:
		s.pageDown(state)
	case ActionGoToLine:
		s.goToLine(state)
	case ActionNewFile:
		state.NewFile()
	case ActionOpenFile:
		s.openFile(state)
	case ActionSaveAs:
		s.saveAs(state)
	case ActionCloseFile:
		state.CloseFile()
	case ActionQuit:
		// This should be handled by parent model
	case ActionSettings:
		// This should be handled by parent model
	case ActionExport:
		// Export current content
		s.exportContent(state)
	case ActionTheoryTools:
		// This should be handled by parent model
	case ActionAudioTools:
		// This should be handled by parent model
	case ActionToggleHelp:
		// Help mode is handled by shortcut manager
	case ActionChordPicker:
		// This should be handled by parent model
	case ActionBPMTapper:
		// This should be handled by parent model
	case ActionBackToMenu:
		// This should be handled by parent model
	// AI Quick Actions - these should be handled by the AI component
	case ActionAIUnstick:
		// This should be handled by AI component
	case ActionAISpark:
		// This should be handled by AI component
	case ActionAITweak:
		// This should be handled by AI component
	case ActionAICheck:
		// This should be handled by AI component
	}

	return nil
}

// GetShortcutManager returns the shortcut manager
func (s *EditorShortcuts) GetShortcutManager() *ShortcutManager {
	return s.shortcutManager
}

// SetShortcutContext sets the keyboard shortcut context
func (s *EditorShortcuts) SetShortcutContext(context KeyContext) {
	if s.shortcutManager != nil {
		s.shortcutManager.SetContext(context)
	}
}

// GetShortcutHints returns shortcut hints for the status bar
func (s *EditorShortcuts) GetShortcutHints() string {
	if s.shortcutManager != nil {
		return s.shortcutManager.GetStatusBarHints()
	}
	return ""
}

// Cursor movement helper methods

func (s *EditorShortcuts) moveCursorToStartOfLine(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) moveCursorToEndOfLine(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) moveCursorToStartOfFile(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) moveCursorToEndOfFile(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) moveCursorToPrevWord(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) moveCursorToNextWord(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

// Text selection helper methods

func (s *EditorShortcuts) selectToStartOfLine(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) selectToEndOfLine(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) selectToStartOfFile(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) selectToEndOfFile(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) selectLeft(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) selectRight(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) selectUp(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) selectDown(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

// Navigation helper methods

func (s *EditorShortcuts) pageUp(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

func (s *EditorShortcuts) pageDown(state StateManagerInterface) {
	// Implementation would depend on textarea capabilities
	// For now, this is a placeholder
}

// Feature helper methods

func (s *EditorShortcuts) goToLine(state StateManagerInterface) {
	// Implementation for go to line functionality
	// For now, this is a placeholder
}

func (s *EditorShortcuts) openFile(state StateManagerInterface) {
	// For now, we'll use a simple approach - in a full implementation,
	// this would show a file dialog. For this demo, we'll show available files
	// and let the user choose via a simple text input

	// This is a simplified implementation
	// In a full implementation, this would show a file picker
	// For now, we'll just open the first available file as a demo
	state.OpenFile("demo.md")
}

func (s *EditorShortcuts) saveAs(state StateManagerInterface) {
	// For this implementation, we'll use a simple filename
	// In a full implementation, this would show a save dialog
	filename := "untitled.md"
	if state.GetCurrentFilePath() != "" {
		filename = filepath.Base(state.GetCurrentFilePath())
	}

	err := state.SaveAs(filename)
	if err != nil {
		// Use proper error handling instead of printf
		return
	}
}

func (s *EditorShortcuts) exportContent(state StateManagerInterface) {
	// This would need access to the export service
	// For now, this is a placeholder implementation
	// In a full implementation, this would:
	// 1. Get the export service from the main model
	// 2. Get the current content from the state
	// 3. Create a title from the filename or use "Untitled"
	// 4. Call the export service's QuickExport method

	// Placeholder implementation
	content := state.GetText()
	title := "Untitled"
	if state.GetCurrentFilePath() != "" {
		title = filepath.Base(state.GetCurrentFilePath())
		title = strings.TrimSuffix(title, filepath.Ext(title))
	}

	// In a full implementation, this would use the export service
	// For now, we'll just log what would happen
	_ = content
	_ = title
	// _, err := exportService.QuickExport(content, title)
	// if err != nil {
	//     // Use proper error handling
	// }
}
