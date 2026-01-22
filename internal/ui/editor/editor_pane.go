package editor

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/app/ai"
	"github.com/Kyanite/noise/internal/data"
	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/export"
	"github.com/Kyanite/noise/internal/infra/files"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/Kyanite/noise/internal/ui/dimension"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// EditorPaneModel handles the text editing pane with syntax highlighting
// This is the refactored version that coordinates between components
type EditorPaneModel struct {
	// Component instances
	state     *EditorState
	shortcuts *EditorShortcuts
	ai        *EditorAI
	metrics   *EditorMetrics

	// UI components that remain in main model
	chordPicker *chordPickerModel
	bpmTapper   *bpmTapperModel

	// Syntax highlighting
	highlighter *SyntaxHighlighter
	elements    []MarkdownElement

	// Theme management
	themeManager interface{} // Using interface{} to avoid import cycle

	// Styles
	focusedStyle     lipgloss.Style
	blurredStyle     lipgloss.Style
	borderStyle      lipgloss.Style
	lineNumberStyle  lipgloss.Style
	cursorLineStyle  lipgloss.Style
	selectionStyle   lipgloss.Style
	searchMatchStyle lipgloss.Style
	autoSaveStyle    lipgloss.Style

	// Dimensions
	width  int
	height int
}

// NewEditorPaneModel creates a new editor pane model with refactored components
func NewEditorPaneModel(textarea textarea.Model) *EditorPaneModel {
	teaModel := &textarea
	t := theme.GetManager().Current()
	model := &EditorPaneModel{
		state:     NewEditorState(teaModel),
		shortcuts: NewEditorShortcuts(),
		ai:        NewEditorAI(),
		metrics:   NewEditorMetrics(),

		highlighter: NewSyntaxHighlighter(),
		chordPicker: NewChordPickerModel(),
		bpmTapper:   NewBPMTapperModel(),

		focusedStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary),
		blurredStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Secondary),
		borderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Secondary),
		lineNumberStyle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Width(4).
			Align(lipgloss.Right),
		cursorLineStyle: lipgloss.NewStyle().
			Background(t.Background),
		selectionStyle: lipgloss.NewStyle().
			Background(t.Background),
		searchMatchStyle: lipgloss.NewStyle().
			Background(t.Accent).
			Foreground(t.Background),
		autoSaveStyle: lipgloss.NewStyle().
			Foreground(t.Success).
			Bold(true),
	}

	return model
}

// Init initializes the editor pane
func (m *EditorPaneModel) Init() tea.Cmd {
	return m.state.Init()
}

// Update handles messages for the editor pane
func (m *EditorPaneModel) Update(msg tea.Msg) (*EditorPaneModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle state component updates
	stateCmd := m.state.Update(msg)
	if stateCmd != nil {
		cmds = append(cmds, stateCmd)
	}

	// Handle chord picker updates
	if m.chordPicker.IsVisible() {
		m.chordPicker, cmd = m.chordPicker.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Handle BPM tapper updates
	if m.bpmTapper.IsVisible() {
		m.bpmTapper, cmd = m.bpmTapper.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Update cursor position
	m.state.UpdateCursorPosition()

	// Update status bar with current state
	m.metrics.UpdateStatusBar(m.state, m.ai)

	// Handle chord picker messages
	switch msg := msg.(type) {
	case ShowChordPickerMsg:
		if m.chordPicker != nil {
			m.chordPicker.visible = true
			m.chordPicker.insertCallback = msg.InsertCallback
			m.chordPicker.selectedIdx = 0
			m.chordPicker.activeMood = "all"
			m.chordPicker.showAll = true
			m.chordPicker.animationTime = time.Now()

			// Load progressions if not already loaded
			if !m.chordPicker.loaded {
				if progressions, err := data.GetAllChordProgressions(); err == nil {
					m.chordPicker.progressions = progressions
					m.chordPicker.filteredProg = progressions
					m.chordPicker.loaded = true
				}
			}
		}
	case HideChordPickerMsg:
		if m.chordPicker != nil {
			m.chordPicker.visible = false
		}
	case ShowBPMTapperMsg:
		if m.bpmTapper != nil {
			m.bpmTapper.visible = true
			m.bpmTapper.setBMPCallback = m.setBPM
			m.bpmTapper.reset()
		}
	case HideBPMTapperMsg:
		if m.bpmTapper != nil {
			m.bpmTapper.visible = false
		}
	}

	// Handle key events for editor features
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Set context for shortcuts
		if m.state.IsFocused() {
			m.shortcuts.SetShortcutContext(ContextEditor)
		} else {
			m.shortcuts.SetShortcutContext(ContextGlobal)
		}

		// Cancel active AI overlays with Escape even if the shortcut is handled elsewhere
		if msg.Type == tea.KeyEsc {
			if m.cancelActiveAIModes() {
				m.metrics.UpdateStatusBar(m.state, m.ai)
				return m, tea.Batch(cmds...)
			}
		}

		// Handle shortcuts (but do not let shortcuts intercept numeric selections when AI modes are active)
		if !m.ai.IsRapidBrainstorm() && !m.ai.IsContinueMode() && !m.ai.IsVariationMode() {
			if action, handled := m.shortcuts.HandleKey(msg); handled {
				if m.handleAIShortcut(action) {
					m.metrics.UpdateStatusBar(m.state, m.ai)
					return m, tea.Batch(cmds...)
				}

				cmd = m.shortcuts.HandleShortcutAction(action, m.state)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
				return m, tea.Batch(cmds...)
			}
		}

		// Handle rapid prototyping key events
		if m.ai.IsRapidBrainstorm() {
			// Robust selection parsing: accept Rune-based KeyMsgs, simple single-char strings,
			// or the Esc key. This covers various ways tests construct KeyMsg values.
			if len(msg.Runes) > 0 && msg.Runes[0] >= '1' && msg.Runes[0] <= '9' {
				index := int(msg.Runes[0] - '0')
				if index >= 1 && index <= 3 {
					m.ai.SelectBrainstormAngle(index-1, m.state)
					return m, tea.Batch(cmds...)
				}
			}
			if s := msg.String(); len(s) == 1 && s[0] >= '1' && s[0] <= '9' {
				index := int(s[0] - '0')
				if index >= 1 && index <= 3 {
					m.ai.SelectBrainstormAngle(index-1, m.state)
					return m, tea.Batch(cmds...)
				}
			}
			if msg.Type == tea.KeyEsc || strings.ToLower(msg.String()) == "esc" {
				m.ai.CancelBrainstormMode()
				return m, tea.Batch(cmds...)
			}
		}

		if m.ai.IsContinueMode() {
			// Robust selection parsing for continue suggestions
			if len(msg.Runes) > 0 && msg.Runes[0] >= '1' && msg.Runes[0] <= '9' {
				index := int(msg.Runes[0] - '0')
				if index >= 1 && index <= 3 {
					m.ai.SelectContinueSuggestion(index-1, m.state)
					return m, tea.Batch(cmds...)
				}
			}
			if s := msg.String(); len(s) == 1 && s[0] >= '1' && s[0] <= '9' {
				index := int(s[0] - '0')
				if index >= 1 && index <= 3 {
					m.ai.SelectContinueSuggestion(index-1, m.state)
					return m, tea.Batch(cmds...)
				}
			}
			if msg.Type == tea.KeyEsc || strings.ToLower(msg.String()) == "esc" {
				m.ai.CancelContinueMode()
				return m, tea.Batch(cmds...)
			}
		}

		if m.ai.IsVariationMode() {
			// Robust selection parsing for variation options
			if len(msg.Runes) > 0 && msg.Runes[0] >= '1' && msg.Runes[0] <= '9' {
				index := int(msg.Runes[0] - '0')
				if index >= 1 && index <= 3 {
					m.ai.SelectVariation(index-1, m.state)
					return m, tea.Batch(cmds...)
				}
			}
			if s := msg.String(); len(s) == 1 && s[0] >= '1' && s[0] <= '9' {
				index := int(s[0] - '0')
				if index >= 1 && index <= 3 {
					m.ai.SelectVariation(index-1, m.state)
					return m, tea.Batch(cmds...)
				}
			}
			if msg.Type == tea.KeyEsc || strings.ToLower(msg.String()) == "esc" {
				m.ai.CancelVariationMode()
				return m, tea.Batch(cmds...)
			}
		}

		// Handle rapid prototyping hotkeys
		switch msg.String() {
		case "ctrl+1":
			// Switch to Sketch mode
			m.state.SetEditorMode(ModeSketch)
			return m, tea.Batch(cmds...)
		case "ctrl+2":
			// Switch to Draft mode
			m.state.SetEditorMode(ModeDraft)
			return m, tea.Batch(cmds...)
		case "ctrl+3":
			// Switch to Polish mode
			m.state.SetEditorMode(ModePolish)
			return m, tea.Batch(cmds...)
		case "ctrl+k":
			// Keep draft (save from scratch mode)
			if m.state.IsScratchMode() {
				// For now, just show a notification
				// In a full implementation, this would prompt for title and save
				m.state.SetScratchMode(false)
			}
			return m, tea.Batch(cmds...)
		}

		// Handle legacy key events for backward compatibility
		switch msg.String() {
		case "enter":
			if m.state.IsSearchMode() {
				// Perform search
				m.state.PerformSearch()
			} else if m.state.AutoIndentEnabled() {
				// Handle auto-indentation
				m.state.HandleAutoIndent()
			}
		}

		// Forward any unhandled key events to the underlying textarea so that
		// rune input and editing keystrokes update the buffer as expected by tests.
		if t := m.state.GetTextarea(); t != nil {
			updatedTA, taCmd := (*t).Update(msg)
			// Apply updated textarea state back to the editor state
			*t = updatedTA
			if taCmd != nil {
				cmds = append(cmds, taCmd)
			}
		}
	}

	// Update syntax highlighting
	m.updateSyntaxHighlighting()

	// Handle auto-save on content changes
	m.state.HandleAutoSave()

	return m, tea.Batch(cmds...)
}

// View renders the editor pane with syntax highlighting
func (m *EditorPaneModel) View() string {
	var style lipgloss.Style

	if m.state.IsFocused() {
		style = m.focusedStyle
	} else {
		style = m.blurredStyle
	}

	// Calculate content dimensions
	contentWidth := m.width - 4
	contentHeight := m.height - 4

	if m.state.ShowLineNumbers() {
		contentWidth -= 4 // Account for line numbers and spacing
	}

	// Set textarea dimensions
	textarea := m.state.GetTextarea()
	textarea.SetWidth(contentWidth)
	textarea.SetHeight(contentHeight)

	// Get content and apply syntax highlighting
	content := m.state.GetText()
	highlightedContent := m.renderHighlightedContent(content, contentWidth)

	// Add line numbers if enabled
	if m.state.ShowLineNumbers() {
		highlightedContent = m.addLineNumbers(highlightedContent)
	}

	// Add title with mode and scratch indicators
	title := "Editor"
	if m.state.IsFocused() {
		title = "Editor (Focused)"
	}

	// Add mode indicator
	modeStr := ""
	switch m.state.GetEditorMode() {
	case ModeSketch:
		modeStr = " [SKETCH]"
	case ModeDraft:
		modeStr = " [DRAFT]"
	case ModePolish:
		modeStr = " [POLISH]"
	}

	// Add scratch mode indicator
	if m.state.IsScratchMode() {
		modeStr += " [SCRATCH]"
	}

	title += modeStr

	t := theme.GetManager().Current()
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Text).
		Background(t.Background).
		Padding(0, 1).
		Width(m.width - 4)

	titleBar := titleStyle.Render(title)

	// Add rapid prototyping UI overlays
	var overlayContent string
	if m.ai.IsRapidBrainstorm() {
		overlayContent = m.ai.RenderOverlays(m.width)
	} else if m.ai.IsContinueMode() {
		overlayContent = m.ai.RenderOverlays(m.width)
	} else if m.ai.IsVariationMode() {
		overlayContent = m.ai.RenderOverlays(m.width)
	} else if m.chordPicker.IsVisible() {
		overlayContent = m.chordPicker.View()
	} else if m.bpmTapper.IsVisible() {
		overlayContent = m.bpmTapper.View()
	}

	// Combine content and overlay
	if overlayContent != "" {
		highlightedContent = lipgloss.JoinVertical(lipgloss.Left, highlightedContent, overlayContent)
	}

	// Add status bar with editor features info
	statusBar := m.metrics.RenderStatusBar()

	// Combine title, content, and status bar
	fullContent := lipgloss.JoinVertical(lipgloss.Left, titleBar, highlightedContent, statusBar)

	return style.Width(m.width).Height(m.height).Render(fullContent)
}

// SetDimensions sets the pane dimensions
func (m *EditorPaneModel) SetDimensions(width, height int) {
	dimension.Set(&m.width, &m.height, width, height)
	m.state.SetDimensions(width, height)
	m.metrics.SetDimensions(width, height)
}

// GetDimensions returns the pane dimensions
func (m *EditorPaneModel) GetDimensions() (int, int) {
	return m.width, m.height
}

// Focus focuses the editor pane
func (m *EditorPaneModel) Focus() {
	m.state.Focus()
}

// Blur blurs the editor pane
func (m *EditorPaneModel) Blur() {
	m.state.Blur()
}

// Public API methods - delegate to components

// GetText returns the current text content
func (m *EditorPaneModel) GetText() string {
	return m.state.GetText()
}

// SetText sets the text content
func (m *EditorPaneModel) SetText(text string) {
	m.state.SetText(text)
	m.updateSyntaxHighlighting()
}

// GetSong returns the current song being edited
func (m *EditorPaneModel) GetSong() *domain.Song {
	return m.state.GetSong()
}

// SetSong sets the current song being edited
func (m *EditorPaneModel) SetSong(song *domain.Song) {
	m.state.SetSong(song)
}

// Service injection methods

// SetAutoSaveService sets the auto-save service
func (m *EditorPaneModel) SetAutoSaveService(service *app.AutoSaveService) {
	m.state.SetAutoSaveService(service)
}

// SetFileService sets the file service
func (m *EditorPaneModel) SetFileService(service *files.Service) {
	m.state.SetFileService(service)
}

// SetExportService sets the export service
func (m *EditorPaneModel) SetExportService(service *export.ExportService) {
	m.state.SetExportService(service)
}

// SetThemeManager sets the theme manager
func (m *EditorPaneModel) SetThemeManager(manager interface{}) {
	m.themeManager = manager
	m.state.SetThemeManager(manager)
}

// SetAIAgent sets the AI agent
func (m *EditorPaneModel) SetAIAgent(agent *ai.QuickIdeaAgent) {
	m.ai.SetAIAgent(agent)
	if m.metrics != nil {
		m.metrics.UpdateKnowledgeBaseStatus(false, "")
		m.metrics.UpdateStatusBar(m.state, m.ai)
	}
}

// SetAIService sets the AI service
func (m *EditorPaneModel) SetAIService(service *app.AIService) {
	m.ai.SetAIService(service)
}

// GetAIAgent returns the current AI agent
func (m *EditorPaneModel) GetAIAgent() *ai.QuickIdeaAgent {
	return m.ai.GetAIAgent()
}

// Shortcut management methods

// GetShortcutManager returns the shortcut manager
func (m *EditorPaneModel) GetShortcutManager() *ShortcutManager {
	return m.shortcuts.GetShortcutManager()
}

// SetShortcutContext sets the keyboard shortcut context
func (m *EditorPaneModel) SetShortcutContext(context KeyContext) {
	m.shortcuts.SetShortcutContext(context)
}

// GetShortcutHints returns shortcut hints for the status bar
func (m *EditorPaneModel) GetShortcutHints() string {
	return m.shortcuts.GetShortcutHints()
}

// HasStatusBar indicates whether the editor has an initialized status bar component
func (m *EditorPaneModel) HasStatusBar() bool {
	return m.metrics != nil && m.metrics.GetStatusBar() != nil
}

// IsContinueMode reports whether the AI continue mode is active.
func (m *EditorPaneModel) IsContinueMode() bool {
	return m.ai.IsContinueMode()
}

// IsVariationMode reports whether the AI variation mode is active.
func (m *EditorPaneModel) IsVariationMode() bool {
	return m.ai.IsVariationMode()
}

// IsRapidBrainstorm reports whether the AI rapid brainstorm mode is active.
func (m *EditorPaneModel) IsRapidBrainstorm() bool {
	return m.ai.IsRapidBrainstorm()
}

// GetBrainstormTheme returns the current brainstorm theme.
func (m *EditorPaneModel) GetBrainstormTheme() string {
	return m.ai.GetBrainstormTheme()
}

// GetVariationOriginal returns the original text used for variation mode.
func (m *EditorPaneModel) GetVariationOriginal() string {
	return m.ai.GetVariationOriginal()
}

// UpdateStatusBar refreshes status metrics manually
func (m *EditorPaneModel) UpdateStatusBar() {
	if m.metrics != nil {
		m.metrics.UpdateStatusBar(m.state, m.ai)
	}
}

// UpdateZoomLevel propagates preview zoom level to the status bar
func (m *EditorPaneModel) UpdateZoomLevel(zoomLevel int) {
	if m.metrics != nil {
		m.metrics.UpdateZoomLevel(zoomLevel)
	}
}

// UpdateShortcutHints updates shortcut hints in the status bar
func (m *EditorPaneModel) UpdateShortcutHints(hints string) {
	if m.metrics != nil {
		m.metrics.UpdateShortcutHints(hints)
	}
}

// GetStatusBar exposes the underlying status bar for package-level consumers
func (m *EditorPaneModel) GetStatusBar() *StatusBarModel {
	if m.metrics == nil {
		return nil
	}
	return m.metrics.GetStatusBar()
}

// File operations

// GetCurrentFilePath returns the current file path
func (m *EditorPaneModel) GetCurrentFilePath() string {
	return m.state.GetCurrentFilePath()
}

// GetFileService returns the file service
func (m *EditorPaneModel) GetFileService() *files.Service {
	return m.state.GetFileService()
}

// Editor mode methods

// SetScratchMode sets the scratch mode
func (m *EditorPaneModel) SetScratchMode(scratchMode bool) {
	m.state.SetScratchMode(scratchMode)
}

// IsScratchMode returns whether the editor is in scratch mode
func (m *EditorPaneModel) IsScratchMode() bool {
	return m.state.IsScratchMode()
}

// SetEditorMode sets the editor mode
func (m *EditorPaneModel) SetEditorMode(mode EditorMode) {
	m.state.SetEditorMode(mode)
}

// GetEditorMode returns the current editor mode
func (m *EditorPaneModel) GetEditorMode() EditorMode {
	return m.state.GetEditorMode()
}

// AI methods

// StartRapidBrainstorm starts rapid brainstorm
func (m *EditorPaneModel) StartRapidBrainstorm(theme string) {
	m.ai.StartRapidBrainstorm(theme)
}

// StartContinueMode starts continue mode
func (m *EditorPaneModel) StartContinueMode() {
	m.ai.StartContinueMode()
	// Generate continue suggestions based on current content
	m.ai.GenerateContinueSuggestions(m.state)
}

// StartVariationMode starts variation mode
func (m *EditorPaneModel) StartVariationMode(selectedText string) {
	m.ai.StartVariationMode(selectedText)
}

// ForceSave performs an immediate save through the state component
func (m *EditorPaneModel) ForceSave() error {
	return m.state.ForceSave()
}

// Private helper methods

// updateSyntaxHighlighting updates the syntax highlighting elements
func (m *EditorPaneModel) updateSyntaxHighlighting() {
	content := m.state.GetText()
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
		lineNum := fmt.Sprintf("%*d", 4, i+1)
		lineNumber := m.lineNumberStyle.Render(lineNum)
		numberedLine := lineNumber + " " + line
		numberedLines = append(numberedLines, numberedLine)
	}

	return strings.Join(numberedLines, "\n")
}

// insertChords inserts the selected chords into the editor
func (m *EditorPaneModel) insertChords(chords []string) {
	// Format chords as a string
	chordStr := strings.Join(chords, " - ")

	// Get current content
	currentContent := m.state.GetText()

	// Add chords to content with proper formatting
	if currentContent != "" {
		chordStr = "\n\n" + chordStr
	}

	// Insert chords at cursor position
	m.state.SetText(currentContent + chordStr)

	// Update syntax highlighting
	m.updateSyntaxHighlighting()
}

// setBPM sets the BPM in the current pattern
func (m *EditorPaneModel) setBPM(bpm int) {
	// For now, just add a comment with the BPM
	// In a full implementation, this would update the pattern data
	bpmComment := fmt.Sprintf("\n\n<!-- BPM: %d -->", bpm)

	// Get current content
	currentContent := m.state.GetText()

	// Add BPM comment to content
	m.state.SetText(currentContent + bpmComment)

	// Update syntax highlighting
	m.updateSyntaxHighlighting()
}

// AI helper methods for backward compatibility

// CancelBrainstormMode cancels brainstorm mode
func (m *EditorPaneModel) CancelBrainstormMode() {
	m.ai.CancelBrainstormMode()
}

// CancelContinueMode cancels continue mode
func (m *EditorPaneModel) CancelContinueMode() {
	m.ai.CancelContinueMode()
}

// CancelVariationMode cancels variation mode
func (m *EditorPaneModel) CancelVariationMode() {
	m.ai.CancelVariationMode()
}

// GetBrainstormAngles returns brainstorm angles
func (m *EditorPaneModel) GetBrainstormAngles() []string {
	return m.ai.GetBrainstormAngles()
}

// GetContinueSuggestions returns continue suggestions
func (m *EditorPaneModel) GetContinueSuggestions() []string {
	return m.ai.GetContinueSuggestions()
}

// GetVariationOptions returns variation options
func (m *EditorPaneModel) GetVariationOptions() []string {
	return m.ai.GetVariationOptions()
}

// handleAIShortcut processes AI-related shortcut actions.
func (m *EditorPaneModel) handleAIShortcut(action ShortcutAction) bool {
	switch action.Type {
	case ActionAIUnstick:
		m.ai.StartContinueMode()
		return true

	case ActionAISpark:
		theme := m.deriveBrainstormTheme()
		m.ai.StartRapidBrainstorm(theme)
		return true

	case ActionAITweak:
		selection := m.deriveVariationSelection()
		if selection == "" {
			return true
		}
		m.ai.StartVariationMode(selection)
		return true

	case ActionAICheck:
		m.ai.PerformQualityCheck(m.state)
		return true

	case ActionBackToMenu:
		return m.cancelActiveAIModes()
	}

	return false
}

func (m *EditorPaneModel) cancelActiveAIModes() bool {
	var cancelled bool

	if m.ai.IsContinueMode() {
		m.ai.CancelContinueMode()
		cancelled = true
	}

	if m.ai.IsVariationMode() {
		m.ai.CancelVariationMode()
		cancelled = true
	}

	if m.ai.IsRapidBrainstorm() {
		m.ai.CancelBrainstormMode()
		cancelled = true
	}

	return cancelled
}

// deriveBrainstormTheme returns a non-empty theme for brainstorm mode.
func (m *EditorPaneModel) deriveBrainstormTheme() string {
	text := strings.TrimSpace(m.state.GetText())
	if text == "" {
		return "Song Inspiration"
	}

	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return trimmed
		}
	}

	return "Song Inspiration"
}

// deriveVariationSelection returns representative text for variation mode.
func (m *EditorPaneModel) deriveVariationSelection() string {
	text := m.state.GetText()
	if strings.TrimSpace(text) == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed != "" {
			return trimmed
		}
	}

	return ""
}

// SelectBrainstormAngle selects a brainstorm angle
func (m *EditorPaneModel) SelectBrainstormAngle(index int) {
	m.ai.SelectBrainstormAngle(index, m.state)
}

// SelectContinueSuggestion selects a continue suggestion
func (m *EditorPaneModel) SelectContinueSuggestion(index int) {
	m.ai.SelectContinueSuggestion(index, m.state)
}

// SelectVariation selects a variation
func (m *EditorPaneModel) SelectVariation(index int) {
	m.ai.SelectVariation(index, m.state)
}

// PerformQualityCheck performs a quality check
func (m *EditorPaneModel) PerformQualityCheck() {
	m.ai.PerformQualityCheck(m.state)
}
