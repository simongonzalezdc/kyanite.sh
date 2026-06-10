package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/ai"
	_ "github.com/kyanite/syntax/internal/character" // Used in story.Project
	"github.com/kyanite/syntax/internal/editor"
	"github.com/kyanite/syntax/internal/location"
	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/spellcheck"
	"github.com/kyanite/syntax/internal/storage"
	"github.com/kyanite/syntax/internal/story"
	"github.com/kyanite/syntax/internal/theme"
)

// Screen represents different app screens
type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenProjectList
	ScreenEditor
	ScreenCharacters
	ScreenScenes
	ScreenLocations
	ScreenLocationEditor
	ScreenTextEditor
	ScreenHelp
	ScreenStats
	ScreenExport
	ScreenAISuggestion
	ScreenRelationshipMap
	ScreenSceneValidation
	ScreenBackups
)

// Model is the root Bubble Tea model
// It follows a centralized state management pattern where all application state
// is held in one place. While this creates a large struct, it simplifies state
// management and message passing in the Bubble Tea architecture.
type Model struct {
	// ============================================================
	// UI State - Screen management and window dimensions
	// ============================================================
	CurrentScreen  Screen // Current active screen
	PreviousScreen Screen // Previous screen for back navigation
	Width          int    // Terminal width
	Height         int    // Terminal height

	// ============================================================
	// Theme State - Visual styling and appearance
	// ============================================================
	ThemeManager *theme.Manager // Manages theme switching
	CurrentTheme theme.Theme    // Active theme
	Styles       theme.Styles   // Computed styles from theme

	// ============================================================
	// Project State - Story projects and data
	// ============================================================
	CurrentProject *story.Project     // Active project being edited
	Projects       []*story.Project   // List of available projects

	// ============================================================
	// Input State - User input handling
	// ============================================================
	SelectedIndex int    // Index of selected item in lists
	InputMode     bool   // Whether input field is active
	InputValue    string // Current input field value
	ReplaceValue  string // Replacement value for find & replace

	// ============================================================
	// Editor State - Text editing and buffer management
	// ============================================================
	CurrentScene    *scene.Scene       // Scene being edited
	CurrentLocation *location.Location // Location being edited
	Buffer          *editor.Buffer     // Text editor buffer
	EditorMode      EditorMode         // Editor mode (normal/insert/search/replace)

	// ============================================================
	// Feature Integrations - AI and spell checking
	// ============================================================
	AIClient     *ai.Client           // AI assistant client
	AISuggestion *ai.Suggestion       // Current AI suggestion
	AIGenerating bool                 // Whether AI is generating
	SpellChecker *spellcheck.Checker  // Spell check integration

	// ============================================================
	// Auto-save State - Automatic saving and status
	// ============================================================
	LastSaveTime time.Time   // Time of last save
	SaveStatus   SaveStatus  // Current save status
	LastEditTime time.Time   // Time of last edit (for auto-save)

	// ============================================================
	// User Feedback - Messages and errors
	// ============================================================
	Message string // Status message to display
	Error   error  // Current error if any
}

// SaveStatus represents the current save state
type SaveStatus int

const (
	SaveStatusSaved SaveStatus = iota
	SaveStatusSaving
	SaveStatusUnsaved
)

// AutoSaveTickMsg is sent periodically to trigger auto-save check
type AutoSaveTickMsg struct{}

// AutoSaveCompleteMsg is sent when auto-save completes
type AutoSaveCompleteMsg struct {
	Err error
}

// NewModel creates a new root model
func NewModel() Model {
	themeManager := theme.NewManager("monochrome")
	currentTheme := themeManager.GetCurrent()

	return Model{
		CurrentScreen: ScreenWelcome,
		ThemeManager:  themeManager,
		CurrentTheme:  currentTheme,
		Styles:        currentTheme.ApplyTheme(),
		SelectedIndex: 0,
		SaveStatus:    SaveStatusSaved,
		LastSaveTime:  time.Now(),
		SpellChecker:  spellcheck.NewChecker(),
	}
}

// autoSaveTick creates a periodic tick command for auto-save
func autoSaveTick() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return AutoSaveTickMsg{}
	})
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Start auto-save ticker
	return autoSaveTick()
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case AISuggestionMsg:
		m = m.HandleAISuggestionMsg(msg)
		return m, nil

	case AutoSaveTickMsg:
		// Check if auto-save is needed
		cmd := m.checkAutoSave()
		// Schedule next tick
		return m, tea.Batch(cmd, autoSaveTick())

	case AutoSaveCompleteMsg:
		if msg.Err != nil {
			m.SaveStatus = SaveStatusUnsaved
			m.Error = msg.Err
		} else {
			m.SaveStatus = SaveStatusSaved
			m.LastSaveTime = time.Now()
			if m.Buffer != nil {
				m.Buffer.SetModified(false)
			}
		}
		return m, nil
	}

	return m, nil
}

// checkAutoSave performs auto-save if needed
func (m *Model) checkAutoSave() tea.Cmd {
	if !m.CanAutoSave() {
		return nil
	}

	// Perform auto-save
	m.SaveStatus = SaveStatusSaving
	projectDir := m.CurrentProject.Directory
	scene := m.CurrentScene
	content := m.Buffer.GetContent()

	return func() tea.Msg {
		scene.Content = content
		err := storage.SaveScene(projectDir, scene)
		return AutoSaveCompleteMsg{Err: err}
	}
}

// ensureDataLoaded ensures necessary data is loaded for the current screen
func (m *Model) ensureDataLoaded() {
	if m.CurrentProject == nil {
		return
	}

	switch m.CurrentScreen {
	case ScreenScenes, ScreenSceneValidation:
		if m.CurrentProject.Scenes == nil || len(m.CurrentProject.Scenes) == 0 {
			scenes, err := storage.LoadAllScenes(m.CurrentProject.Directory)
			if err == nil {
				m.CurrentProject.Scenes = scenes
			}
		}
	case ScreenCharacters, ScreenRelationshipMap:
		if m.CurrentProject.Characters == nil || len(m.CurrentProject.Characters) == 0 {
			chars, err := storage.LoadAllCharacters(m.CurrentProject.Directory)
			if err == nil {
				m.CurrentProject.Characters = chars
			}
		}
	case ScreenLocations:
		if m.CurrentProject.Locations == nil || len(m.CurrentProject.Locations) == 0 {
			locs, err := storage.LoadAllLocations(m.CurrentProject.Directory)
			if err == nil {
				m.CurrentProject.Locations = locs
			}
		}
	}
}

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Store previous screen to detect changes
	prevScreen := m.CurrentScreen

	// Global shortcuts
	switch msg.String() {
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit

	case "ctrl+shift+t":
		// Cycle theme
		m.CurrentTheme = m.ThemeManager.NextTheme()
		m.Styles = m.CurrentTheme.ApplyTheme()
		m.Message = fmt.Sprintf("Theme: %s", m.CurrentTheme.Name)
		return m, nil

	case "?", "h":
		// Show help
		m.PreviousScreen = m.CurrentScreen
		m.CurrentScreen = ScreenHelp
		m.ensureDataLoaded()
		return m, nil
	}

	// Screen-specific shortcuts
	switch m.CurrentScreen {
	case ScreenWelcome:
		return m.handleWelcomeKeys(msg)
	case ScreenProjectList:
		return m.handleProjectListKeys(msg)
	case ScreenEditor:
		return m.handleEditorKeys(msg)
	case ScreenCharacters:
		updatedModel, cmd := m.handleCharactersKeys(msg)
		m = updatedModel.(Model)
		if m.CurrentScreen != prevScreen {
			m.ensureDataLoaded()
		}
		return m, cmd
	case ScreenScenes:
		updatedModel, cmd := m.handleScenesKeys(msg)
		m = updatedModel.(Model)
		if m.CurrentScreen != prevScreen {
			m.ensureDataLoaded()
		}
		return m, cmd
	case ScreenLocations:
		return m.handleLocationsKeys(msg)
	case ScreenLocationEditor:
		return m.handleLocationEditorKeys(msg)
	case ScreenTextEditor:
		return m.handleTextEditorKeys(msg)
	case ScreenHelp:
		return m.handleHelpKeys(msg)
	case ScreenStats:
		return m.handleStatsKeys(msg)
	case ScreenExport:
		return m.handleExportKeys(msg)
	case ScreenAISuggestion:
		return m.handleAISuggestionKeys(msg)
	case ScreenRelationshipMap:
		updatedModel, cmd := m.handleRelationshipMapKeys(msg)
		m = updatedModel.(Model)
		if m.CurrentScreen != prevScreen {
			m.ensureDataLoaded()
		}
		return m, cmd
	case ScreenSceneValidation:
		updatedModel, cmd := m.handleSceneValidationKeys(msg)
		m = updatedModel.(Model)
		if m.CurrentScreen != prevScreen {
			m.ensureDataLoaded()
		}
		return m, cmd
	case ScreenBackups:
		return m.handleBackupsKeys(msg)
	}

	return m, nil
}

// View renders the current screen
func (m Model) View() string {
	switch m.CurrentScreen {
	case ScreenWelcome:
		return m.viewWelcome()
	case ScreenProjectList:
		return m.viewProjectList()
	case ScreenEditor:
		return m.viewEditor()
	case ScreenCharacters:
		return m.viewCharacters()
	case ScreenScenes:
		return m.viewScenes()
	case ScreenLocations:
		return m.viewLocations()
	case ScreenLocationEditor:
		return m.viewLocationEditor()
	case ScreenTextEditor:
		return m.viewTextEditor()
	case ScreenHelp:
		return m.viewHelp()
	case ScreenStats:
		return m.viewStats()
	case ScreenExport:
		return m.viewExport()
	case ScreenAISuggestion:
		return m.viewAISuggestion()
	case ScreenRelationshipMap:
		return m.viewRelationshipMap()
	case ScreenSceneValidation:
		return m.viewSceneValidation()
	case ScreenBackups:
		return m.viewBackups()
	default:
		return "Unknown screen"
	}
}

// ============================================================
// Helper Methods - Encapsulate common state checks and operations
// ============================================================

// HasProject returns whether a project is currently loaded
func (m *Model) HasProject() bool {
	return m.CurrentProject != nil
}

// HasUnsavedChanges returns whether there are unsaved changes in the editor
func (m *Model) HasUnsavedChanges() bool {
	return m.Buffer != nil && m.Buffer.IsModified()
}

// IsEditorActive returns whether the text editor is currently active
func (m *Model) IsEditorActive() bool {
	return m.CurrentScreen == ScreenTextEditor && m.Buffer != nil && m.CurrentScene != nil
}

// CanAutoSave returns whether auto-save can be performed
func (m *Model) CanAutoSave() bool {
	return m.IsEditorActive() &&
	       m.HasProject() &&
	       m.HasUnsavedChanges() &&
	       time.Since(m.LastEditTime) >= 3*time.Second
}

// ResetInputState clears all input-related state
func (m *Model) ResetInputState() {
	m.InputMode = false
	m.InputValue = ""
	m.ReplaceValue = ""
}

// GetContentDimensions returns the available dimensions for content rendering
// after accounting for title and status bars
func (m *Model) GetContentDimensions() (width, height int) {
	return m.Width, m.Height - 4 // Minus title and status bars
}

// SetMessage sets a status message for the user
func (m *Model) SetMessage(msg string) {
	m.Message = msg
}

// SetError sets an error and clears any status message
func (m *Model) SetError(err error) {
	m.Error = err
	if err != nil {
		m.Message = ""
	}
}

// ClearFeedback clears both messages and errors
func (m *Model) ClearFeedback() {
	m.Message = ""
	m.Error = nil
}

// SaveCurrentScene saves the current scene with error handling
// Returns true if save was successful, false otherwise
func (m *Model) SaveCurrentScene() bool {
	if !m.HasProject() || m.CurrentScene == nil || m.Buffer == nil {
		return false
	}

	m.SaveStatus = SaveStatusSaving
	m.CurrentScene.Content = m.Buffer.GetContent()
	err := storage.SaveScene(m.CurrentProject.Directory, m.CurrentScene)

	if err != nil {
		m.SetError(err)
		m.SaveStatus = SaveStatusUnsaved
		return false
	}

	m.Buffer.SetModified(false)
	m.SaveStatus = SaveStatusSaved
	m.LastSaveTime = time.Now()
	m.SetMessage("Saved")
	return true
}

// ExitEditor cleanly exits the text editor, saving if needed
func (m *Model) ExitEditor(saveIfModified bool) {
	if saveIfModified && m.HasUnsavedChanges() {
		m.SaveCurrentScene()
	}

	if m.Buffer != nil {
		m.Buffer.ClearSearch()
	}

	m.CurrentScene = nil
	m.Buffer = nil
	m.CurrentScreen = ScreenScenes
}
