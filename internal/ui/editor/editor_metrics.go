package editor

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/charmbracelet/lipgloss"
)

// EditorMetrics implementation methods

// UpdateStatusBar updates the status bar with current editor state
func (m *EditorMetrics) UpdateStatusBar(state StateManagerInterface, ai *EditorAI) {
	if m.statusBar == nil {
		return
	}

	// Update content and statistics
	content := state.GetText()
	m.statusBar.UpdateContent(content)

	// Update cursor position
	m.statusBar.UpdateCursorPosition(state.GetCursorLine(), state.GetCursorColumn())

	// Update content type detection
	if ai != nil {
		contentType := ai.AnalyzeContentType(content)
		m.statusBar.UpdateContentType(contentType)
	}

	// Update auto-save status
	autoSaveStatus := state.GetAutoSaveStatus()
	lastSaveTime := state.GetLastSaveTime()
	m.statusBar.UpdateAutoSaveStatus(autoSaveStatus, lastSaveTime)

	// Update editor features
	m.statusBar.UpdateEditorFeatures(
		state.ShowLineNumbers(),
		state.WordWrapEnabled(),
		state.AutoIndentEnabled(),
		state.BracketMatchingEnabled(),
	)

	// Update current file path
	filename := "Untitled"
	if state.GetCurrentFilePath() != "" {
		filename = filepath.Base(state.GetCurrentFilePath())
	}
	m.statusBar.UpdateFileInfo(filename)

	// Update dimensions and responsive mode
	m.statusBar.SetDimensions(m.width, 1)
	m.statusBar.UpdateResponsiveMode(m.width)
}

// RenderStatusBar renders the status bar with editor features info
func (m *EditorMetrics) RenderStatusBar() string {
	if m.statusBar == nil {
		// Fallback to old implementation if status bar not initialized
		return "Status bar not initialized"
	}

	// Use the new status bar component
	return m.statusBar.View()
}

// SetDimensions sets the metrics component dimensions
func (m *EditorMetrics) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	if m.statusBar != nil {
		m.statusBar.SetDimensions(width, 1)
	}
}

// GetDimensions returns the metrics component dimensions
func (m *EditorMetrics) GetDimensions() (int, int) {
	return m.width, m.height
}

// UpdateKnowledgeBaseStatus updates the knowledge base status
func (m *EditorMetrics) UpdateKnowledgeBaseStatus(available bool, statusText string) {
	if m.statusBar != nil {
		m.statusBar.UpdateKnowledgeBaseStatus(available, statusText)
	}
}

// UpdateZoomLevel updates the zoom level indicator in the status bar
func (m *EditorMetrics) UpdateZoomLevel(zoomLevel int) {
	if m.statusBar != nil {
		m.statusBar.UpdateZoomLevel(zoomLevel)
	}
}

// GetStatusBar returns the status bar model
func (m *EditorMetrics) GetStatusBar() *StatusBarModel {
	return m.statusBar
}

// SetStatusBar sets the status bar model
func (m *EditorMetrics) SetStatusBar(statusBar *StatusBarModel) {
	m.statusBar = statusBar
}

// RenderFallbackStatusBar renders a fallback status bar when the main status bar is not available
func (m *EditorMetrics) RenderFallbackStatusBar(state StateManagerInterface, shortcuts *EditorShortcuts) string {
	var features []string

	if state.ShowLineNumbers() {
		features = append(features, "L")
	}
	if state.WordWrapEnabled() {
		features = append(features, "W")
	}
	if state.AutoIndentEnabled() {
		features = append(features, "I")
	}
	if state.BracketMatchingEnabled() {
		features = append(features, "B")
	}

	featuresStr := strings.Join(features, " ")

	// Add auto-save status
	autoSaveStatus := ""
	autoSaveStatusEnum := state.GetAutoSaveStatus()
	switch autoSaveStatusEnum {
	case app.AutoSaveSaving:
		autoSaveStatus = " Saving..."
	case app.AutoSaveSuccess:
		autoSaveStatus = " Saved"
	case app.AutoSaveError:
		autoSaveStatus = " Save Error"
	default:
		// Show last save time for idle status
		if !state.GetLastSaveTime().IsZero() {
			autoSaveStatus = fmt.Sprintf(" Saved %s", state.GetLastSaveTime().Format("15:04:05"))
		}
	}

	// Combine editor status and auto-save status
	editorStatus := fmt.Sprintf("Ln %d, Col %d | %s", state.GetCursorLine()+1, state.GetCursorColumn()+1, featuresStr)
	if autoSaveStatus != "" {
		editorStatus = fmt.Sprintf("%s |%s", editorStatus, autoSaveStatus)
	}

	// Add shortcut hints if available
	if shortcuts != nil {
		hints := shortcuts.GetShortcutHints()
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

// UpdateContentStats updates content statistics
func (m *EditorMetrics) UpdateContentStats(content string) {
	if m.statusBar != nil {
		m.statusBar.UpdateContent(content)
	}
}

// UpdateCursorPosition updates the cursor position in the status bar
func (m *EditorMetrics) UpdateCursorPosition(line, column int) {
	if m.statusBar != nil {
		m.statusBar.UpdateCursorPosition(line, column)
	}
}

// UpdateContentType updates the content type in the status bar
func (m *EditorMetrics) UpdateContentType(contentType string) {
	if m.statusBar != nil {
		m.statusBar.UpdateContentType(contentType)
	}
}

// UpdateAutoSaveStatus updates the auto-save status in the status bar
func (m *EditorMetrics) UpdateAutoSaveStatus(status app.AutoSaveStatus, lastSaveTime time.Time) {
	if m.statusBar != nil {
		m.statusBar.UpdateAutoSaveStatus(status, lastSaveTime)
	}
}

// UpdateEditorFeatures updates the editor features in the status bar
func (m *EditorMetrics) UpdateEditorFeatures(showLineNumbers, wordWrap, autoIndent, bracketMatching bool) {
	if m.statusBar != nil {
		m.statusBar.UpdateEditorFeatures(showLineNumbers, wordWrap, autoIndent, bracketMatching)
	}
}

// UpdateFileInfo updates the file information in the status bar
func (m *EditorMetrics) UpdateFileInfo(filename string) {
	if m.statusBar != nil {
		m.statusBar.UpdateFileInfo(filename)
	}
}

// UpdateShortcutHints updates the shortcut hints in the status bar
func (m *EditorMetrics) UpdateShortcutHints(hints string) {
	if m.statusBar != nil {
		m.statusBar.UpdateShortcutHints(hints)
	}
}

// UpdateResponsiveMode updates the responsive mode in the status bar
func (m *EditorMetrics) UpdateResponsiveMode(width int) {
	if m.statusBar != nil {
		m.statusBar.UpdateResponsiveMode(width)
	}
}

// GetContentStats returns content statistics
func (m *EditorMetrics) GetContentStats() (lines, words, chars int) {
	// Placeholder implementation
	// In a full implementation, this would calculate actual statistics
	return 0, 0, 0
}

// GetCursorPosition returns the current cursor position
func (m *EditorMetrics) GetCursorPosition() (int, int) {
	if m.statusBar != nil {
		return m.statusBar.GetCursorPosition()
	}
	return 0, 0
}

// GetContentType returns the current content type
func (m *EditorMetrics) GetContentType() string {
	// Placeholder implementation
	// In a full implementation, this would return the actual content type
	return ""
}

// GetAutoSaveStatus returns the current auto-save status
func (m *EditorMetrics) GetAutoSaveStatus() (app.AutoSaveStatus, time.Time) {
	// Placeholder implementation
	// In a full implementation, this would return the actual auto-save status
	return app.AutoSaveIdle, time.Time{}
}

// GetEditorFeatures returns the current editor features
func (m *EditorMetrics) GetEditorFeatures() (showLineNumbers, wordWrap, autoIndent, bracketMatching bool) {
	if m.statusBar != nil {
		return m.statusBar.GetEditorFeatures()
	}
	return false, false, false, false
}

// GetFileInfo returns the current file information
func (m *EditorMetrics) GetFileInfo() string {
	// Placeholder implementation
	// In a full implementation, this would return the actual file info
	return ""
}

// GetShortcutHints returns the current shortcut hints
func (m *EditorMetrics) GetShortcutHints() string {
	// Placeholder implementation
	// In a full implementation, this would return the actual shortcut hints
	return ""
}

// GetKnowledgeBaseStatus returns the knowledge base status
func (m *EditorMetrics) GetKnowledgeBaseStatus() (bool, string) {
	// Placeholder implementation
	// In a full implementation, this would return the actual knowledge base status
	return false, ""
}

// IsResponsiveMode returns whether responsive mode is active
func (m *EditorMetrics) IsResponsiveMode() bool {
	// Placeholder implementation
	// In a full implementation, this would return the actual responsive mode status
	return false
}

// Clear clears the metrics component
func (m *EditorMetrics) Clear() {
	// Placeholder implementation
	// In a full implementation, this would clear the status bar
}

// Reset resets the metrics component to its initial state
func (m *EditorMetrics) Reset() {
	// Placeholder implementation
	// In a full implementation, this would reset the status bar
	m.width = 0
	m.height = 0
}
