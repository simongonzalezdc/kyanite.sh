package editor

import (
	"context"
	"fmt"
	"sort"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/config"
	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/infra/files"
	"github.com/Kyanite/noise/internal/logging"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Screen represents different screens in the application (local copy to avoid import cycle)
type screen int

const (
	screenEditor screen = iota
	screenExport
	screenTheory
	screenAudio
)

// ScreenChangeMsg represents a message to change screens (local copy to avoid import cycle)
type ScreenChangeMsg struct {
	Screen screen
}

// SplitPaneModel represents the main split-pane editor layout
type SplitPaneModel struct {
	// Layout dimensions
	width      int
	height     int
	splitRatio float64 // Ratio of left pane to total width (0.0-1.0)

	// Child components
	editorPane  *EditorPaneModel
	previewPane *PreviewPaneModel
	fileDialog  *FileDialogModel

	// State
	focusedPane     FocusedPane
	database        *db.DB
	autoSaveService *app.AutoSaveService
	fileService     *files.Service
	ctx             context.Context
	cancel          context.CancelFunc

	// Keyboard shortcuts
	shortcutManager *ShortcutManager

	// Performance optimizations
	lastUpdateLength int

	// Mode-specific layouts
	sketchLayout *SketchLayout
	draftLayout  *DraftLayout
	polishLayout *PolishLayout

	// Styles
	dividerStyle lipgloss.Style
}

// FocusedPane represents which pane currently has focus
type FocusedPane int

const (
	EditorPane FocusedPane = iota
	PreviewPane
)

// NewSplitPaneModel creates a new split-pane editor model
func NewSplitPaneModel(database *db.DB) *SplitPaneModel {
	// Initialize text area for editor pane
	editorTA := textarea.New()
	editorTA.Placeholder = "Start writing your lyrics..."
	editorTA.Focus()

	// Default split ratio (50/50)
	splitRatio := 0.5

	// Initialize context and services
	ctx, cancel := context.WithCancel(context.Background())
	autoSaveService := app.NewAutoSaveService(database, nil) // Use default config

	// Initialize file service
	fileService, err := files.New(files.Config{
		BaseDir:  "./songs", // Default songs directory
		AutoSave: false,     // We'll handle saving manually for now
	})
	if err != nil {
		logging.Warnf("Failed to initialize file service: %v", err)
		// Continue without file service - operations will be no-ops
		fileService = nil
	}

	t := theme.GetManager().Current()
	model := &SplitPaneModel{
		splitRatio:      splitRatio,
		editorPane:      NewEditorPaneModel(editorTA),
		previewPane:     NewPreviewPaneModel(),
		fileDialog:      NewFileDialogModel(DialogOpen, "Open File", "./songs", []string{".md", ".txt"}),
		focusedPane:     EditorPane,
		database:        database,
		autoSaveService: autoSaveService,
		fileService:     fileService,
		ctx:             ctx,
		cancel:          cancel,
		shortcutManager: NewShortcutManager(),
		dividerStyle: lipgloss.NewStyle().
			Foreground(t.Secondary),
		sketchLayout: NewSketchLayout(),
		draftLayout:  NewDraftLayout(),
		polishLayout: NewPolishLayout(),
	}

	// Set up services for editor pane
	model.editorPane.SetAutoSaveService(autoSaveService)
	model.editorPane.SetFileService(fileService)

	// Start the auto-save service
	if err := autoSaveService.Start(ctx); err != nil {
		logging.Warnf("Failed to start auto-save service: %v", err)
	}

	return model
}

// Init initializes the split-pane model
func (m *SplitPaneModel) Init() tea.Cmd {
	// Set up file dialog callbacks
	m.setupFileDialogCallbacks()

	return tea.Batch(
		m.editorPane.Init(),
		m.previewPane.Init(),
		m.fileDialog.Init(),
	)
}

// Update handles messages for the split-pane layout
func (m *SplitPaneModel) Update(msg tea.Msg) (*SplitPaneModel, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle file dialog first (highest priority)
	if m.fileDialog.IsVisible() {
		var cmd tea.Cmd
		m.fileDialog, cmd = m.fileDialog.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updatePaneDimensions()
		// Update file dialog dimensions
		m.fileDialog.SetDimensions(m.width, m.height)

	case tea.KeyMsg:
		// Set context based on focused pane
		switch m.focusedPane {
		case EditorPane:
			m.shortcutManager.SetContext(ContextEditor)
			m.editorPane.SetShortcutContext(ContextEditor)
		case PreviewPane:
			m.shortcutManager.SetContext(ContextPreview)
			m.previewPane.SetShortcutContext(ContextPreview)
		}

		// Handle shortcuts
		if action, handled := m.shortcutManager.HandleKey(msg); handled {
			return m.handleShortcutAction(action)
		}

		// Handle legacy key events for backward compatibility
		switch msg.String() {
		case "tab":
			// Switch focus between panes
			m.switchFocus()
			return m, nil
		}
	}

	// Update focused pane
	var cmd tea.Cmd
	switch m.focusedPane {
	case EditorPane:
		m.editorPane, cmd = m.editorPane.Update(msg)
		cmds = append(cmds, cmd)

		// Update preview with current editor content for real-time preview
		// Only update if content has actually changed to avoid unnecessary re-renders
		editorContent := m.editorPane.GetText()
		currentContent := m.previewPane.GetContent()

		// Skip update if content hasn't changed (for large documents, this saves significant processing)
		if editorContent == currentContent && len(editorContent) > 10000 { // 10KB threshold
			// For very large documents, only update if content length has changed
			if m.lastUpdateLength == len(editorContent) {
				// Content length hasn't changed, skip update
				return m, nil
			}
		}
		m.lastUpdateLength = len(editorContent)

		if m.previewPane.GetRealtimeManager() != nil {
			m.previewPane.GetRealtimeManager().UpdateContent(editorContent, ChangeSourceEditor)
		} else {
			m.previewPane.SetContent(editorContent)
		}

		// Update status bar with current zoom level
		if m.editorPane.HasStatusBar() {
			zoomLevel := m.previewPane.GetZoomLevel()
			m.editorPane.UpdateZoomLevel(zoomLevel)
		}

	case PreviewPane:
		m.previewPane, cmd = m.previewPane.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the split-pane layout
func (m *SplitPaneModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// If file dialog is visible, render it as overlay
	if m.fileDialog.IsVisible() {
		// Render the background first
		background := m.renderDefaultLayout()

		// Render file dialog overlay
		dialogView := m.fileDialog.View()

		// Combine background with dialog overlay
		return lipgloss.JoinVertical(lipgloss.Left, background, dialogView)
	}

	// Get the current editor mode
	editorMode := m.editorPane.GetEditorMode()

	// Set dimensions for layouts
	m.sketchLayout.SetDimensions(m.width, m.height)
	m.draftLayout.SetDimensions(m.width, m.height)
	m.polishLayout.SetDimensions(m.width, m.height)

	// Get content from panes
	editorContent := m.editorPane.View()
	previewContent := m.previewPane.View()

	// Render based on mode
	switch editorMode {
	case ModeSketch:
		// Sketch mode: Editor + AI panel
		brainstormContent := "AI Assistant\n\nCtrl+G: Continue\nCtrl+V: Variations\n\nTheme brainstorming will appear here..."
		return m.sketchLayout.Render(editorContent, brainstormContent)

	case ModeDraft:
		// Draft mode: Editor + Preview + Theory
		theoryContent := "Theory Tools\n\nRhyme Dictionary\nSyllable Counter\n\nTools will appear here..."
		return m.draftLayout.Render(editorContent, previewContent, theoryContent)

	case ModePolish:
		// Polish mode: Full suite
		theoryContent := "Theory Tools\n\nCircle of Fifths\nChord Progressions\n\nTools will appear here..."
		critiqueContent := "AI Critique\n\nQuality analysis will appear here..."
		return m.polishLayout.Render(editorContent, previewContent, theoryContent, critiqueContent)

	default:
		// Fallback to original layout
		return m.renderDefaultLayout()
	}
}

// renderDefaultLayout renders the default split-pane layout
func (m *SplitPaneModel) renderDefaultLayout() string {
	// Calculate responsive split ratio based on terminal width
	splitRatio := m.calculateResponsiveSplitRatio()

	// Calculate pane dimensions
	leftWidth := int(float64(m.width) * splitRatio)
	rightWidth := m.width - leftWidth - 1 // -1 for divider

	// Ensure minimum widths based on responsive breakpoints
	minPaneWidth := m.getMinimumPaneWidth()

	if leftWidth < minPaneWidth {
		leftWidth = minPaneWidth
		rightWidth = m.width - leftWidth - 1
	}
	if rightWidth < minPaneWidth {
		rightWidth = minPaneWidth
		leftWidth = m.width - rightWidth - 1
	}

	// Render panes
	editorView := m.editorPane.View()
	editorView = lipgloss.NewStyle().Width(leftWidth).Height(m.height).Render(editorView)

	previewView := m.previewPane.View()
	previewView = lipgloss.NewStyle().Width(rightWidth).Height(m.height).Render(previewView)

	// Create divider
	divider := m.dividerStyle.Render("â”‚")

	// Combine views horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, editorView, divider, previewView)
}

// updatePaneDimensions updates the dimensions of child panes
func (m *SplitPaneModel) updatePaneDimensions() {
	// Use responsive split ratio for dimension calculations
	splitRatio := m.calculateResponsiveSplitRatio()
	leftWidth := int(float64(m.width) * splitRatio)
	rightWidth := m.width - leftWidth - 1

	// Apply minimum width constraints
	minPaneWidth := m.getMinimumPaneWidth()
	if leftWidth < minPaneWidth {
		leftWidth = minPaneWidth
		rightWidth = m.width - leftWidth - 1
	}
	if rightWidth < minPaneWidth {
		rightWidth = minPaneWidth
		leftWidth = m.width - rightWidth - 1
	}

	// Update editor pane dimensions
	m.editorPane.SetDimensions(leftWidth, m.height)

	// Update preview pane dimensions
	m.previewPane.SetDimensions(rightWidth, m.height)
}

// switchFocus switches focus between panes
func (m *SplitPaneModel) switchFocus() {
	if m.focusedPane == EditorPane {
		m.focusedPane = PreviewPane
		m.editorPane.Blur()
		m.previewPane.Focus()
	} else {
		m.focusedPane = EditorPane
		m.previewPane.Blur()
		m.editorPane.Focus()
	}
}

// GetEditorText returns the current editor text
func (m *SplitPaneModel) GetEditorText() string {
	return m.editorPane.GetText()
}

// SetEditorText sets the editor text
func (m *SplitPaneModel) SetEditorText(text string) {
	m.editorPane.SetText(text)
}

// SetCurrentSong sets the given song into the editor pane and updates editor text.
// This is an exported helper so callers in other packages (e.g., root UI) can open a song.
func (m *SplitPaneModel) SetCurrentSong(song *domain.Song) {
	if m.editorPane == nil {
		return
	}
	if song == nil {
		m.editorPane.SetSong(nil)
		return
	}
	// Set song on editor pane (will populate RawContent into the textarea)
	m.editorPane.SetSong(song)
	// Update editor text explicitly if RawContent exists
	if song.RawContent != "" {
		m.editorPane.SetText(song.RawContent)
	}
}

// Cleanup cleans up resources when the model is destroyed
func (m *SplitPaneModel) Cleanup() {
	if m.cancel != nil {
		m.cancel()
	}
	if m.autoSaveService != nil {
		if err := m.autoSaveService.Stop(); err != nil {
			logging.Warnf("Error stopping auto-save service: %v", err)
		}
	}
}

// handleShortcutAction handles actions from the keyboard shortcut system
func getThemeCycleOrder() []string {
	ids := theme.ListThemes()
	if len(ids) == 0 {
		return nil
	}

	sort.Strings(ids)

	seenNames := make(map[string]struct{}, len(ids))
	result := make([]string, 0, len(ids))

	for _, id := range ids {
		th := theme.GetTheme(id)
		if _, exists := seenNames[th.Name]; exists {
			continue
		}
		seenNames[th.Name] = struct{}{}
		result = append(result, id)
	}

	return result
}

func (m *SplitPaneModel) rotateTheme(delta int) {
	ids := getThemeCycleOrder()
	if len(ids) == 0 {
		return
	}

	currentTheme := theme.GetManager().Current()
	currentIdx := -1
	for i, id := range ids {
		th := theme.GetTheme(id)
		if th.Name == currentTheme.Name {
			currentIdx = i
			break
		}
	}
	if currentIdx == -1 {
		currentIdx = 0
	}

	nextIdx := (currentIdx + delta) % len(ids)
	if nextIdx < 0 {
		nextIdx += len(ids)
	}
	nextID := ids[nextIdx]

	theme.GetManager().SetTheme(nextID)

	if cfg, err := config.Load(); err == nil {
		cfg.UI.Theme = nextID
		_ = cfg.Save()
	}

	if m.editorPane != nil {
		if m.editorPane.HasStatusBar() {
			m.editorPane.UpdateShortcutHints("Theme: " + nextID)
		}
		m.editorPane.UpdateStatusBar()
	}
}

func (m *SplitPaneModel) handleShortcutAction(action ShortcutAction) (*SplitPaneModel, tea.Cmd) {
	switch action.Type {
	case ActionNextPane:
		m.switchFocus()
		return m, nil
	case ActionPrevPane:
		// Switch to the other pane (reverse of next)
		if m.focusedPane == EditorPane {
			m.switchFocus() // This will switch to PreviewPane
		} else {
			m.switchFocus() // This will switch to EditorPane
		}
		return m, nil
	case ActionBackToMenu:
		// This should be handled by the root model
		// For now, we'll just return without action
		return m, nil
	case ActionQuit:
		// This should be handled by the root model
		return m, tea.Quit
	case ActionSettings:
		// This should be handled by the root model
		return m, nil
	case ActionTheoryTools:
		// Navigate to theory tools screen
		return m, func() tea.Msg {
			return ScreenChangeMsg{Screen: screenTheory}
		}
	case ActionAudioTools:
		// Navigate to audio tools screen
		return m, func() tea.Msg {
			return ScreenChangeMsg{Screen: screenAudio}
		}
	case ActionThemeCycle:
		m.rotateTheme(1)
		return m, nil
	case ActionNextTheme:
		m.rotateTheme(1)
		return m, nil
	case ActionPreviousTheme:
		m.rotateTheme(-1)
		return m, nil
	case ActionNewFile:
		// Clear current content
		m.editorPane.SetText("")
		return m, nil
	case ActionOpenFile:
		// Show open file dialog
		m.showOpenFileDialog()
		return m, nil
	case ActionSave:
		// Trigger save in editor pane
		if m.autoSaveService != nil {
			err := m.editorPane.ForceSave()
			if err != nil {
				// Handle save error - could log or show user notification
				logging.Errorf("Save failed: %v", err)
			}
		}
		return m, nil
	case ActionSaveAs:
		// Show save as file dialog
		m.showSaveAsDialog()
		return m, nil
	case ActionExport:
		// Navigate to export screen
		return m, func() tea.Msg {
			return ScreenChangeMsg{Screen: screenExport}
		}
	case ActionCloseFile:
		// Clear current content
		m.editorPane.SetText("")
		return m, nil
	case ActionToggleHelp:
		// Toggle help mode in shortcut manager
		m.shortcutManager.SetHelpMode(!m.shortcutManager.IsHelpMode())
		return m, nil
	default:
		// For other actions, pass them to the focused pane
		switch m.focusedPane {
		case EditorPane:
			// Actions are already handled in the editor pane's Update method
			return m, nil
		case PreviewPane:
			// Preview pane would handle its own actions
			return m, nil
		}
		return m, nil
	}
}

// GetShortcutManager returns the shortcut manager for external access
func (m *SplitPaneModel) GetShortcutManager() *ShortcutManager {
	return m.shortcutManager
}

// SetShortcutContext sets the keyboard shortcut context
func (m *SplitPaneModel) SetShortcutContext(context KeyContext) {
	if m.shortcutManager != nil {
		m.shortcutManager.SetContext(context)
	}
}

// calculateResponsiveSplitRatio calculates the optimal split ratio based on terminal width
func (m *SplitPaneModel) calculateResponsiveSplitRatio() float64 {
	// Use responsive breakpoints for different terminal sizes
	if m.width < 100 {
		return 0.7 // Favor editor pane more in very small terminals
	} else if m.width < 120 {
		return 0.6
	} else if m.width < 160 {
		return 0.5
	} else {
		return 0.45 // Favor preview pane more in ultra-wide terminals
	}
}

// getMinimumPaneWidth returns the minimum pane width based on terminal size
func (m *SplitPaneModel) getMinimumPaneWidth() int {
	if m.width < 100 {
		return 10 // Smaller minimum for compact terminals
	}
	return 15 // Standard minimum for larger terminals
}

// SetQuickStartConfig configures the split pane for quick start mode
func (m *SplitPaneModel) SetQuickStartConfig(theme string, scratchMode bool, autoBrainstorm bool) {
	if m.editorPane == nil {
		return
	}

	// Set scratch mode
	m.editorPane.SetScratchMode(scratchMode)

	// If auto-brainstorm is enabled and theme is provided, start brainstorm
	if autoBrainstorm && theme != "" {
		m.editorPane.StartRapidBrainstorm(theme)
	}
}

// File dialog helper methods

// setupFileDialogCallbacks sets up the callbacks for the file dialog
func (m *SplitPaneModel) setupFileDialogCallbacks() {
	// Open file callback
	m.fileDialog.SetConfirmCallback(func(path string) error {
		if err := m.editorPane.state.OpenFile(path); err != nil {
			logging.Errorf("Failed to open file: %v", err)
			return fmt.Errorf("failed to open file: %w", err)
		}
		logging.Infof("Opened file: %s", path)
		return nil
	})

	// Cancel callback
	m.fileDialog.SetCancelCallback(func() {
		logging.Debugf("File dialog cancelled")
	})
}

// showOpenFileDialog shows the open file dialog
func (m *SplitPaneModel) showOpenFileDialog() {
	// Create a new open file dialog
	currentPath := m.editorPane.GetCurrentFilePath()
	if currentPath == "" {
		currentPath = "./songs"
	}

	m.fileDialog = NewFileDialogModel(DialogOpen, "Open File", currentPath, []string{".md", ".txt"})
	m.fileDialog.SetDimensions(m.width, m.height)
	m.setupFileDialogCallbacks()
	m.fileDialog.Show()
}

// showSaveAsDialog shows the save as file dialog
func (m *SplitPaneModel) showSaveAsDialog() {
	// Create a new save as dialog
	currentPath := m.editorPane.GetCurrentFilePath()
	if currentPath == "" {
		currentPath = "./songs/untitled.md"
	}

	m.fileDialog = NewFileDialogModel(DialogSaveAs, "Save As", currentPath, []string{".md", ".txt"})
	m.fileDialog.SetDimensions(m.width, m.height)

	// Set up save as callback
	m.fileDialog.SetConfirmCallback(func(path string) error {
		if err := m.editorPane.state.SaveAs(path); err != nil {
			logging.Errorf("Failed to save file: %v", err)
			return fmt.Errorf("failed to save file: %w", err)
		}
		logging.Infof("Saved file: %s", path)
		return nil
	})

	m.fileDialog.SetCancelCallback(func() {
		logging.Debugf("Save as dialog cancelled")
	})

	m.fileDialog.Show()
}
