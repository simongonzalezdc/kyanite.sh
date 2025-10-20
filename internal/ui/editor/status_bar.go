package editor

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/Kyanite/noise/internal/ui/dimension"
	"github.com/charmbracelet/lipgloss"
)

// StatusBarSection represents a section of the status bar
type StatusBarSection struct {
	Content string
	Style   lipgloss.Style
	Width   int
}

// StatusBarModel represents the comprehensive status bar
type StatusBarModel struct {
	width  int
	height int

	// Content information
	content        string
	cursorLine     int
	cursorColumn   int
	wordCount      int
	characterCount int
	lineCount      int

	// Status information
	autoSaveStatus app.AutoSaveStatus
	lastSaveTime   time.Time
	editorMode     string
	fileName       string
	zoomLevel      int
	contentType    string // New field for content type detection
	kbAvailable    bool   // Knowledge base availability status
	kbStatus       string // Knowledge base status message

	// Editor features
	showLineNumbers bool
	wordWrap        bool
	autoIndent      bool
	bracketMatching bool

	// Active shortcuts hints
	shortcutHints string

	// Performance optimization
	lastContentHash        uint64
	lastUpdateTime         time.Time
	updateThrottleDuration time.Duration // Renamed from updateThrottleMs to avoid Ms suffix

	// Responsive layout
	compactMode     bool
	showMinimalInfo bool

	// Styles for different sections
	leftSectionStyle   lipgloss.Style
	centerSectionStyle lipgloss.Style
	rightSectionStyle  lipgloss.Style

	// Status-specific styles
	autoSaveSavingStyle  lipgloss.Style
	autoSaveSuccessStyle lipgloss.Style
	autoSaveErrorStyle   lipgloss.Style
	autoSaveIdleStyle    lipgloss.Style
	modeIndicatorStyle   lipgloss.Style
	shortcutHintStyle    lipgloss.Style
}

// NewStatusBarModel creates a new status bar model
func NewStatusBarModel() *StatusBarModel {
	model := &StatusBarModel{
		width:                  0,
		height:                 1,
		content:                "",
		cursorLine:             0,
		cursorColumn:           0,
		wordCount:              0,
		characterCount:         0,
		lineCount:              0,
		autoSaveStatus:         app.AutoSaveIdle,
		lastSaveTime:           time.Time{},
		editorMode:             "Normal",
		fileName:               "Untitled",
		zoomLevel:              100,
		contentType:            "Unknown",
		kbAvailable:            false,
		kbStatus:               "KB: Unavailable",
		showLineNumbers:        true,
		wordWrap:               true,
		autoIndent:             true,
		bracketMatching:        true,
		shortcutHints:          "",
		lastContentHash:        0,
		lastUpdateTime:         time.Time{},
		updateThrottleDuration: 100 * time.Millisecond, // Throttle updates to 100ms
		compactMode:            false,
		showMinimalInfo:        false,

		// Styles will be initialized with the theme
	}

	// Initialize styles with theme
	t := theme.GetManager().Current()

	// Initialize styles
	model.leftSectionStyle = lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Background).
		Padding(0, 1)

	model.centerSectionStyle = lipgloss.NewStyle().
		Foreground(t.Secondary).
		Background(t.Background).
		Padding(0, 1)

	model.rightSectionStyle = lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Background).
		Padding(0, 1)

	model.autoSaveSavingStyle = lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true)

	model.autoSaveSuccessStyle = lipgloss.NewStyle().
		Foreground(t.Success)

	model.autoSaveErrorStyle = lipgloss.NewStyle().
		Foreground(t.Error).
		Bold(true)

	model.autoSaveIdleStyle = lipgloss.NewStyle().
		Foreground(t.Text)

	model.modeIndicatorStyle = lipgloss.NewStyle().
		Foreground(t.Error).
		Bold(true)

	model.shortcutHintStyle = lipgloss.NewStyle().
		Foreground(t.Secondary).
		Italic(true)

	return model
}

// SetDimensions sets the status bar dimensions
func (m *StatusBarModel) SetDimensions(width, height int) {
	dimension.Set(&m.width, &m.height, width, height)
}

func (m *StatusBarModel) GetDimensions() (int, int) {
	return m.width, m.height
}

// UpdateContent updates the content and recalculates statistics
func (m *StatusBarModel) UpdateContent(content string) {
	// Compute hash and detect actual change first
	currentHash := m.hashContent(content)
	changed := !(currentHash == m.lastContentHash && content == m.content)

	// If nothing changed and we updated recently, skip work
	now := time.Now()
	if !changed && now.Sub(m.lastUpdateTime) < m.updateThrottleDuration {
		return
	}

	// Always update when content changed; otherwise only when throttle window passed
	if changed || now.Sub(m.lastUpdateTime) >= m.updateThrottleDuration {
		m.content = content
		m.lastContentHash = currentHash
		m.lastUpdateTime = now
		m.calculateStatistics()
	}
}

// UpdateCursorPosition updates the cursor position
func (m *StatusBarModel) UpdateCursorPosition(line, column int) {
	// Only update if position has actually changed
	if m.cursorLine != line || m.cursorColumn != column {
		m.cursorLine = line
		m.cursorColumn = column
	}
}

// UpdateAutoSaveStatus updates the auto-save status and timestamp
func (m *StatusBarModel) UpdateAutoSaveStatus(status app.AutoSaveStatus, lastSaveTime time.Time) {
	// Only update if status or time has actually changed
	if m.autoSaveStatus != status || !m.lastSaveTime.Equal(lastSaveTime) {
		m.autoSaveStatus = status
		m.lastSaveTime = lastSaveTime
	}
}

// UpdateEditorFeatures updates the editor feature flags
func (m *StatusBarModel) UpdateEditorFeatures(showLineNumbers, wordWrap, autoIndent, bracketMatching bool) {
	// Only update if features have actually changed
	if m.showLineNumbers != showLineNumbers || m.wordWrap != wordWrap ||
		m.autoIndent != autoIndent || m.bracketMatching != bracketMatching {
		m.showLineNumbers = showLineNumbers
		m.wordWrap = wordWrap
		m.autoIndent = autoIndent
		m.bracketMatching = bracketMatching
	}
}

// UpdateFileInfo updates file information
func (m *StatusBarModel) UpdateFileInfo(fileName string) {
	m.fileName = fileName
}

// UpdateZoomLevel updates the preview zoom level
func (m *StatusBarModel) UpdateZoomLevel(zoomLevel int) {
	m.zoomLevel = zoomLevel
}

// UpdateShortcutHints updates the active shortcut hints
func (m *StatusBarModel) UpdateShortcutHints(hints string) {
	m.shortcutHints = hints
}

// UpdateContentType updates the detected content type
func (m *StatusBarModel) UpdateContentType(contentType string) {
	// Only update if content type has actually changed
	if m.contentType != contentType {
		m.contentType = contentType
	}
}

// UpdateKnowledgeBaseStatus updates the knowledge base availability status
func (m *StatusBarModel) UpdateKnowledgeBaseStatus(available bool, status string) {
	// Only update if status has actually changed
	if m.kbAvailable != available || m.kbStatus != status {
		m.kbAvailable = available
		m.kbStatus = status
	}
}

// hashContent calculates a simple hash of the content for change detection
func (m *StatusBarModel) hashContent(content string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(content))
	return h.Sum64()
}

// calculateStatistics calculates word count, character count, and line count
func (m *StatusBarModel) calculateStatistics() {
	content := strings.TrimSpace(m.content)

	// Calculate line count
	if content == "" {
		m.lineCount = 0
	} else {
		m.lineCount = strings.Count(content, "\n") + 1
	}

	// Calculate character count (including spaces)
	m.characterCount = len(content)

	// Calculate word count
	words := regexp.MustCompile(`\S+`).FindAllString(content, -1)
	m.wordCount = len(words)
}

// renderLeftSection renders the left section (file info and position)
func (m *StatusBarModel) renderLeftSection() StatusBarSection {
	// File name and cursor position
	position := fmt.Sprintf("%s | Ln %d, Col %d", m.fileName, m.cursorLine+1, m.cursorColumn+1)

	// Add mode indicator if not normal
	if m.editorMode != "Normal" {
		position = fmt.Sprintf("%s | %s", position, m.modeIndicatorStyle.Render(m.editorMode))
	}

	// Add content type indicator
	if m.contentType != "" && m.contentType != "Unknown" {
		contentTypeStyle := m.getContentTypeStyle(m.contentType)
		position = fmt.Sprintf("%s | %s", position, contentTypeStyle.Render(m.contentType))
	}

	content := m.leftSectionStyle.Render(position)
	return StatusBarSection{
		Content: content,
		Style:   m.leftSectionStyle,
		Width:   len(position) + 4, // Account for padding
	}
}

// renderCenterSection renders the center section (document statistics)
func (m *StatusBarModel) renderCenterSection() StatusBarSection {
	// Document statistics
	stats := fmt.Sprintf("%d words | %d chars | %d lines",
		m.wordCount, m.characterCount, m.lineCount)

	content := m.centerSectionStyle.Render(stats)
	return StatusBarSection{
		Content: content,
		Style:   m.centerSectionStyle,
		Width:   len(stats) + 4, // Account for padding
	}
}

// renderRightSection renders the right section (status indicators and hints)
func (m *StatusBarModel) renderRightSection() StatusBarSection {
	var indicators []string

	// Editor features
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

	if len(features) > 0 {
		indicators = append(indicators, strings.Join(features, ""))
	}

	// Auto-save status
	autoSaveText := m.renderAutoSaveStatus()
	if autoSaveText != "" {
		indicators = append(indicators, autoSaveText)
	}

	// Zoom level
	if m.zoomLevel != 100 {
		indicators = append(indicators, fmt.Sprintf("%d%%", m.zoomLevel))
	}

	// Shortcut hints
	if m.shortcutHints != "" {
		hints := m.shortcutHintStyle.Render(m.shortcutHints)
		indicators = append(indicators, hints)
	}

	// Theme indicator (shows active theme id)
	if mgr := theme.GetManager(); mgr != nil {
		themeID := mgr.Current().Name
		if themeID != "" {
			themeLabel := m.modeIndicatorStyle.Render("Theme: " + themeID)
			indicators = append(indicators, themeLabel)
		}

		// Add content type indicator if not already shown in left section
		if m.contentType != "" && m.contentType != "Unknown" && m.editorMode == "Normal" {
			contentTypeStyle := m.getCompactContentTypeStyle(m.contentType)
			indicators = append(indicators, contentTypeStyle.Render(m.contentType))
		}

		// Add knowledge base status indicator
		if m.kbStatus != "" {
			kbStyle := m.getKnowledgeBaseStatusStyle(m.kbAvailable)
			kbIndicator := kbStyle.Render(m.kbStatus)
			indicators = append(indicators, kbIndicator)
		}
	}

	content := strings.Join(indicators, " | ")
	content = m.rightSectionStyle.Render(content)

	return StatusBarSection{
		Content: content,
		Style:   m.rightSectionStyle,
		Width:   len(content) + 4, // Account for padding
	}
}

// renderAutoSaveStatus renders the auto-save status with appropriate styling
func (m *StatusBarModel) renderAutoSaveStatus() string {
	var text string

	switch m.autoSaveStatus {
	case app.AutoSaveSaving:
		text = "Saving..."
		return m.autoSaveSavingStyle.Render(text)
	case app.AutoSaveSuccess:
		if !m.lastSaveTime.IsZero() {
			text = fmt.Sprintf("Saved %s", m.lastSaveTime.Format("15:04:05"))
		} else {
			text = "Saved"
		}
		return m.autoSaveSuccessStyle.Render(text)
	case app.AutoSaveError:
		text = "Save Error"
		return m.autoSaveErrorStyle.Render(text)
	case app.AutoSaveIdle:
		if !m.lastSaveTime.IsZero() {
			text = fmt.Sprintf("Saved %s", m.lastSaveTime.Format("15:04:05"))
			return m.autoSaveIdleStyle.Render(text)
		}
		return ""
	default:
		return ""
	}
}

// View renders the status bar
func (m *StatusBarModel) View() string {
	if m.width == 0 {
		return ""
	}

	// Use responsive rendering based on terminal width
	if m.showMinimalInfo {
		// Minimal mode: only essential information
		essential := m.renderMinimalView()
		return essential
	}

	if m.compactMode {
		// Compact mode: abbreviated information
		compact := m.renderCompactView()
		return compact
	}

	// Full mode: all information
	return m.renderFullView()
}

// renderMinimalView renders the most essential information only
func (m *StatusBarModel) renderMinimalView() string {
	// Show only cursor position and critical indicators
	parts := []string{}

	// Cursor position (most essential)
	parts = append(parts, fmt.Sprintf("Ln%d", m.cursorLine+1))

	// Auto-save status if critical
	autoSaveText := m.renderAutoSaveStatus()
	if autoSaveText != "" && (m.autoSaveStatus == app.AutoSaveSaving || m.autoSaveStatus == app.AutoSaveError) {
		parts = append(parts, autoSaveText)
	}

	content := strings.Join(parts, " â”‚ ")
	content = m.centerSectionStyle.Render(content)

	// Ensure proper width
	if len(content) > m.width {
		content = content[:m.width]
	} else if len(content) < m.width {
		padding := strings.Repeat(" ", m.width-len(content))
		content += padding
	}

	return content
}

// renderCompactView renders abbreviated but complete information
func (m *StatusBarModel) renderCompactView() string {
	// Render sections with responsive content
	leftSection := m.renderCompactLeftSection()
	centerSection := m.renderCompactCenterSection()
	rightSection := m.renderCompactRightSection()

	// Combine sections with compact spacing
	fullContent := leftSection + "â”‚" + centerSection + "â”‚" + rightSection

	// Ensure the content fits exactly within the width
	if len(fullContent) > m.width {
		fullContent = fullContent[:m.width]
	} else if len(fullContent) < m.width {
		padding := strings.Repeat(" ", m.width-len(fullContent))
		fullContent += padding
	}

	return fullContent
}

// renderFullView renders the complete status bar
func (m *StatusBarModel) renderFullView() string {
	// Render sections
	leftSection := m.renderLeftSection()
	centerSection := m.renderCenterSection()
	rightSection := m.renderRightSection()

	// Calculate available space for each section
	usedWidth := leftSection.Width + centerSection.Width + rightSection.Width
	availableWidth := m.width - usedWidth

	// Distribute extra space proportionally
	if availableWidth > 0 {
		// Add extra space to sections that need it
		extraSpace := availableWidth / 3
		if extraSpace > 0 {
			leftSection.Width += extraSpace
			centerSection.Width += extraSpace
			rightSection.Width += extraSpace
		}
	}

	// Create spacer sections to fill remaining space
	leftSpacer := strings.Repeat(" ", max(0, m.width/3-leftSection.Width/2))
	rightSpacer := strings.Repeat(" ", max(0, m.width/3-centerSection.Width/2))

	// Combine sections with proper spacing
	fullContent := leftSection.Content + leftSpacer + centerSection.Content + rightSpacer + rightSection.Content

	// Ensure the content fits exactly within the width
	if len(fullContent) > m.width {
		fullContent = fullContent[:m.width]
	} else if len(fullContent) < m.width {
		padding := strings.Repeat(" ", m.width-len(fullContent))
		fullContent += padding
	}

	return fullContent
}

// renderCompactLeftSection renders a compact version of the left section
func (m *StatusBarModel) renderCompactLeftSection() string {
	// Show only filename and cursor position
	position := fmt.Sprintf("%s Ln%d", m.fileName, m.cursorLine+1)
	return m.leftSectionStyle.Render(position)
}

// renderCompactCenterSection renders a compact version of the center section
func (m *StatusBarModel) renderCompactCenterSection() string {
	// Show only word count and line count
	stats := fmt.Sprintf("%dw %dln", m.wordCount, m.lineCount)
	return m.centerSectionStyle.Render(stats)
}

// renderCompactRightSection renders a compact version of the right section
func (m *StatusBarModel) renderCompactRightSection() string {
	var indicators []string

	// Editor features (abbreviated)
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
	if len(features) > 0 {
		indicators = append(indicators, strings.Join(features, ""))
	}

	// Auto-save status
	autoSaveText := m.renderAutoSaveStatus()
	if autoSaveText != "" {
		indicators = append(indicators, autoSaveText)
	}

	// Zoom level (only if not 100%)
	if m.zoomLevel != 100 {
		indicators = append(indicators, fmt.Sprintf("%d%%", m.zoomLevel))
	}

	return m.rightSectionStyle.Render(strings.Join(indicators, "â”‚"))
}

// Helper function to get maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GetDocumentStats returns the current document statistics
func (m *StatusBarModel) GetDocumentStats() (wordCount, characterCount, lineCount int) {
	return m.wordCount, m.characterCount, m.lineCount
}

// GetCursorPosition returns the current cursor position
func (m *StatusBarModel) GetCursorPosition() (line, column int) {
	return m.cursorLine, m.cursorColumn
}

// GetAutoSaveInfo returns the current auto-save information
func (m *StatusBarModel) GetAutoSaveInfo() (status app.AutoSaveStatus, lastSaveTime time.Time) {
	return m.autoSaveStatus, m.lastSaveTime
}

// GetEditorFeatures returns the current editor feature flags
func (m *StatusBarModel) GetEditorFeatures() (showLineNumbers, wordWrap, autoIndent, bracketMatching bool) {
	return m.showLineNumbers, m.wordWrap, m.autoIndent, m.bracketMatching
}

// NeedsRender checks if the status bar needs to be re-rendered
func (m *StatusBarModel) NeedsRender() bool {
	// Check if enough time has passed since last update
	now := time.Now()
	return now.Sub(m.lastUpdateTime) >= m.updateThrottleDuration
}

// SetUpdateThrottle sets the update throttle duration
func (m *StatusBarModel) SetUpdateThrottle(duration time.Duration) {
	m.updateThrottleDuration = duration
}

// UpdateResponsiveMode updates the responsive display mode based on terminal width
func (m *StatusBarModel) UpdateResponsiveMode(width int) {
	// Enable compact mode for smaller terminals
	compactMode := width < 100
	minimalMode := width < 80

	// Only update if mode has actually changed
	if m.compactMode != compactMode || m.showMinimalInfo != minimalMode {
		m.compactMode = compactMode
		m.showMinimalInfo = minimalMode
	}
}

// renderCompactSection renders a compact version of a section
// TODO: Uncomment when needed
// func (m *StatusBarModel) renderCompactSection(content string) string {
// 	// For very small terminals, show only essential info
// 	if m.showMinimalInfo {
// 		// Show only cursor position and basic indicators
// 		indicators := []string{}
//
// 		// Auto-save status (most important)
// 		autoSaveText := m.renderAutoSaveStatus()
// 		if autoSaveText != "" {
// 			indicators = append(indicators, autoSaveText)
// 		}
//
// 		// Editor features (abbreviated)
// 		var features []string
// 		if m.showLineNumbers {
// 			features = append(features, "L")
// 		}
// 		if m.wordWrap {
// 			features = append(features, "W")
// 		}
// 		if len(features) > 0 {
// 			indicators = append(indicators, strings.Join(features, ""))
// 		}
//
// 		if len(indicators) > 0 {
// 			return m.rightSectionStyle.Render(strings.Join(indicators, "â”‚"))
// 		}
// 		return ""
// 	}
//
// 	// Compact mode: show abbreviated information
// 	if m.compactMode {
// 		// Show essential info only
// 		parts := []string{}
//
// 		// Cursor position (essential)
// 		parts = append(parts, fmt.Sprintf("Ln%d", m.cursorLine+1))
//
// 		// Word count (if space allows)
// 		if len(content) < 20 {
// 			parts = append(parts, fmt.Sprintf("%dw", m.wordCount))
// 		}
//
// 		// Auto-save status
// 		autoSaveText := m.renderAutoSaveStatus()
// 		if autoSaveText != "" {
// 			parts = append(parts, autoSaveText)
// 		}
//
// 		return m.centerSectionStyle.Render(strings.Join(parts, " â”‚ "))
// 	}
//
// 	// Full mode: show all information
// 	return content
// }

// getContentTypeStyle returns the appropriate style for the content type
func (m *StatusBarModel) getContentTypeStyle(contentType string) lipgloss.Style {
	switch strings.ToLower(contentType) {
	case "lyrics":
		t := theme.GetManager().Current()
		return lipgloss.NewStyle().
			Foreground(t.Success).
			Bold(true)
	case "patterns":
		t := theme.GetManager().Current()
		return lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true)
	case "mixed":
		t := theme.GetManager().Current()
		return lipgloss.NewStyle().
			Foreground(t.Error).
			Bold(true)
	default:
		t := theme.GetManager().Current()
		return lipgloss.NewStyle().
			Foreground(t.Text).
			Bold(true)
	}
}

// getCompactContentTypeStyle returns a compact style for content type indicators
func (m *StatusBarModel) getCompactContentTypeStyle(contentType string) lipgloss.Style {
	switch strings.ToLower(contentType) {
	case "lyrics":
		t := theme.GetManager().Current()
		return lipgloss.NewStyle().
			Foreground(t.Success).
			Bold(true)
	case "patterns":
		t := theme.GetManager().Current()
		return lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true)
	case "mixed":
		t := theme.GetManager().Current()
		return lipgloss.NewStyle().
			Foreground(t.Error).
			Bold(true)
	default:
		t := theme.GetManager().Current()
		return lipgloss.NewStyle().
			Foreground(t.Text)
	}
}

// getKnowledgeBaseStatusStyle returns the appropriate style for knowledge base status
func (m *StatusBarModel) getKnowledgeBaseStatusStyle(available bool) lipgloss.Style {
	if available {
		t := theme.GetManager().Current()
		return lipgloss.NewStyle().
			Foreground(t.Success).
			Bold(true)
	} else {
		t := theme.GetManager().Current()
		return lipgloss.NewStyle().
			Foreground(t.Text).
			Italic(true)
	}
}
