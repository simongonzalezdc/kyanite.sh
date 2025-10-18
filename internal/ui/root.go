package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/noise/internal/collaboration"
	"github.com/puente-labs/noise/internal/config"
	"github.com/puente-labs/noise/internal/infra/db"
	"github.com/puente-labs/noise/internal/plugins"
	"github.com/puente-labs/noise/internal/ui/editor"
	"github.com/puente-labs/noise/internal/ui/styles"
	"github.com/puente-labs/noise/internal/theme"
)

// Screen represents different screens in the application
type screen int

const (
	screenSplash screen = iota
	screenMenu
	screenEditor
	screenExport
	screenTheory
	screenAudio
	screenManager
	screenSettings
	screenLoading
)

// RootModel is the main application model that handles routing between screens
type RootModel struct {
	// Current state
	currentScreen screen
	width         int
	height        int

	// Database connection
	database *db.DB

	// Configuration
	config *config.Config

	// Child models
	splash   *SplashModel
	menu     *MenuModel
	editor   *EditorModel
	export   *ExportModel
	theory   *TheoryModel
	audio    *AudioModel
	manager  *ManagerModel
	settings *SettingsModel

	// Help system
	helpMode bool
	helpPane *editor.HelpPaneModel

	// Loading state
	loading  bool
	spinner  spinner.Model
	errorMsg string

	// Animation system
	animation *AnimationManager

	// Responsive layout system
	responsiveManager *ResponsiveLayoutManager

	// Collaboration system
	collaborationManager *collaboration.CollaborationManager
	presenceManager      *collaboration.PresenceManager
	sessionManager       *collaboration.SessionManager
	invitationManager    *collaboration.InvitationManager
	conflictResolver     *collaboration.ConflictResolver

	// Plugin system
	pluginManager *plugins.DefaultManager
}

// NewRootModel creates a new root model with initialized state
func NewRootModel(pluginManager *plugins.DefaultManager) *RootModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.Primary)

	return &RootModel{
		currentScreen:        screenSplash,
		loading:              true,
		spinner:              s,
		animation:            NewAnimationManager(),
		responsiveManager:    NewResponsiveLayoutManager(),
		collaborationManager: collaboration.NewCollaborationManager(nil), // Database will be set after initialization
		presenceManager:      collaboration.NewPresenceManager(),
		sessionManager:       collaboration.NewSessionManager(),
		invitationManager:    collaboration.NewInvitationManager(),
		conflictResolver:     collaboration.NewConflictResolver(),
		pluginManager:        pluginManager,
	}
}

// Init initializes the root model
func (m *RootModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.initializeApp(),
		m.spinner.Tick,
	)
}

// initializeApp initializes the application (database, etc.)
func (m *RootModel) initializeApp() tea.Cmd {
	return func() tea.Msg {
		// Initialize database
		database, err := db.New(db.Config{})
		if err != nil {
			return initErrorMsg{err: fmt.Errorf("failed to initialize database: %w", err)}
		}

		return initSuccessMsg{database: database}
	}
}

// Update handles messages and updates the model
func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Handle responsive layout updates
		sizeCmd := m.responsiveManager.HandleWindowSizeMsg(msg)
		if sizeCmd != nil {
			cmds = append(cmds, sizeCmd)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.database != nil {
				m.database.Close()
			}
			// Clean up animation manager
			if m.animation != nil {
				m.animation.Close()
			}
			return m, tea.Quit
		case "esc":
			if m.helpMode {
				m.helpMode = false
				// Start fade out animation for help
				m.animation.FadeTransition("help_fade_out", 1.0)
				return m, nil
			} else if m.currentScreen != screenMenu {
				// Start screen transition animation
				m.animation.SlideTransition("screen_transition", 1.0)
				m.currentScreen = screenMenu
				return m, nil
			}
		case "f1", "?":
			// Toggle help mode
			m.helpMode = !m.helpMode
			return m, nil
		}

	case initSuccessMsg:
		m.database = msg.database
		m.loading = false
		m.currentScreen = screenMenu

		// Initialize collaboration system
		m.initializeCollaborationSystem()

		// Initialize child models
		m.initializeChildModels()

	case initErrorMsg:
		m.errorMsg = msg.err.Error()
		m.loading = false

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case AnimationTickMsg:
		// Update global animations
		cmd := m.animation.Update()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case SizeValidationMsg:
		// Handle size validation results
		if !msg.IsValid {
			// Could show a temporary warning overlay or log the warnings
			// For now, we'll just ensure the app continues to work
			// TODO: Implement size validation warnings display
			// Log validation issues for debugging
			fmt.Printf("Size validation issues detected for terminal size\n")
		}

	case ScreenChangeMsg:
		// Handle screen changes from menu
		m.currentScreen = msg.Screen
	}

	// Update current screen
	cmd := m.updateCurrentScreen(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// initializeChildModels initializes all child models
func (m *RootModel) initializeChildModels() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// Use default config if loading fails
		cfg = config.DefaultConfig()
	}
	m.config = cfg

	// Initialize theme system from config (applies styles globally)
	theme.InitFromConfig(cfg)

	m.splash = NewSplashModel()
	m.menu = NewMenuModel()
	m.editor = NewEditorModel(m.database)
	m.export = NewExportModel("") // Content will be set when entering export screen
	m.theory = NewTheoryModel()
	m.audio = NewAudioModel()
	m.manager = NewManagerModel(m.database)
	m.settings = NewSettingsModel(m.config)

	// Initialize help system
	m.helpPane = editor.NewHelpPaneModel(nil) // Shortcut manager will be set when needed
}

// initializeCollaborationSystem initializes the collaboration system
func (m *RootModel) initializeCollaborationSystem() {
	// Set database for collaboration manager
	m.collaborationManager = collaboration.NewCollaborationManager(m.database)

	// Set up UI callbacks for collaboration events
	m.collaborationManager.SetUICallbacks(
		m.onSessionUpdate,
		m.onPresenceUpdate,
		m.onConflictDetected,
	)

	// Set up presence manager callbacks
	m.presenceManager.SetPresenceCallbacks(
		m.onPresenceUpdate,
		func(sessionID, userID string) {
			// Handle user joined for presence manager
		},
		func(sessionID, userID string) {
			// Handle user left for presence manager
		},
	)

	// Set up session manager callbacks
	m.sessionManager.SetSessionCallbacks(
		m.onSessionCreated,
		m.onSessionEnded,
		m.onUserJoined,
		m.onUserLeft,
	)
}

// updateCurrentScreen routes messages to the current screen's model
func (m *RootModel) updateCurrentScreen(msg tea.Msg) tea.Cmd {
	switch m.currentScreen {
	case screenSplash:
		return m.updateSplash(msg)
	case screenMenu:
		return m.updateMenu(msg)
	case screenEditor:
		return m.updateEditor(msg)
	case screenExport:
		return m.updateExport(msg)
	case screenTheory:
		return m.updateTheory(msg)
	case screenAudio:
		return m.updateAudio(msg)
	case screenManager:
		// Propagate OpenSongMsg from manager to editor by setting editor state
		switch v := msg.(type) {
		case OpenSongMsg:
			if m.editor != nil {
				// Set the editor content and current song
				if v.Song != nil {
					// Set editor text and current song via exported helper
					m.editor.SetEditorText(v.Song.RawContent)
					if m.editor.GetSplitPane() != nil {
						m.editor.GetSplitPane().SetCurrentSong(v.Song)
					}
					// Switch to editor screen
					m.currentScreen = screenEditor
					return nil
				}
			}
		}
		return m.updateManager(msg)
	case screenSettings:
		return m.updateSettings(msg)
	}
	return nil
}

// View renders the current screen
func (m *RootModel) View() string {
	// Show help overlay if help mode is active
	if m.helpMode {
		return m.renderHelp()
	}

	var content string

	switch m.currentScreen {
	case screenSplash:
		content = m.renderSplash()
	case screenMenu:
		content = m.renderMenu()
	case screenEditor:
		content = m.renderEditor()
	case screenTheory:
		content = m.renderTheory()
	case screenAudio:
		content = m.renderAudio()
	case screenManager:
		content = m.renderManager()
	case screenSettings:
		content = m.renderSettings()
	case screenExport:
		content = m.renderExport()
	case screenLoading:
		content = m.renderLoading()
	default:
		content = "Unknown screen"
	}

	// Add size warnings if terminal size is not optimal
	if warnings := m.responsiveManager.GetSizeWarnings(); len(warnings) > 0 {
		warningBox := m.responsiveManager.RenderSizeWarning()
		if warningBox != "" {
			// Add warning at the top of the content
			content = warningBox + "\n" + content
		}
	}

	return content
}

// renderSplash renders the splash screen
func (m *RootModel) renderSplash() string {
	if m.loading {
		return m.renderLoading()
	}
	return m.splash.View()
}

// renderMenu renders the main menu
func (m *RootModel) renderMenu() string {
	return m.menu.View()
}

// renderEditor renders the editor screen
func (m *RootModel) renderEditor() string {
	return m.editor.View()
}

// renderTheory renders the theory tools screen
func (m *RootModel) renderTheory() string {
	return m.theory.View()
}

// renderAudio renders the audio tools screen
func (m *RootModel) renderAudio() string {
	return m.audio.View()
}

// renderManager renders the project manager screen
func (m *RootModel) renderManager() string {
	return m.manager.View()
}

// renderSettings renders the settings screen
func (m *RootModel) renderSettings() string {
	return m.settings.View()
}

// renderLoading renders the loading screen
func (m *RootModel) renderLoading() string {
	loadingStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width).
		Height(m.height)

	loadingText := "🎵 noise.sh 🎵\n\n"
	if m.errorMsg != "" {
		loadingText += "Error: " + m.errorMsg + "\n\nPress any key to exit..."
	} else {
		loadingText += "Initializing...\n\n" + m.spinner.View()
	}

	return loadingStyle.Render(loadingText)
}

// renderHelp renders the help screen
func (m *RootModel) renderHelp() string {
	// Set up help pane with current editor's shortcut manager if available
	if m.editor != nil && m.editor.GetSplitPane() != nil {
		shortcutManager := m.editor.GetSplitPane().GetShortcutManager()
		if shortcutManager != nil {
			m.helpPane.SetShortcutManager(shortcutManager)
		}
	}

	// Set dimensions for help pane
	m.helpPane.SetDimensions(m.width, m.height)

	// Update help pane (for any animations or state changes)
	_, _ = m.helpPane.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})

	return m.helpPane.View()
}

// Message types for initialization
type initSuccessMsg struct {
	database *db.DB
}

type initErrorMsg struct {
	err error
}

// Message type for screen changes
type ScreenChangeMsg struct {
	Screen screen
}

// Screen update methods (placeholders for now)
func (m *RootModel) updateSplash(msg tea.Msg) tea.Cmd {
	if m.splash != nil {
		return m.splash.Update(msg)
	}
	return nil
}

func (m *RootModel) updateMenu(msg tea.Msg) tea.Cmd {
	if m.menu != nil {
		_, cmd := m.menu.Update(msg)
		return cmd
	}
	return nil
}

func (m *RootModel) updateEditor(msg tea.Msg) tea.Cmd {
	if m.editor != nil {
		_, cmd := m.editor.Update(msg)
		return cmd
	}
	return nil
}

func (m *RootModel) updateTheory(msg tea.Msg) tea.Cmd {
	if m.theory != nil {
		_, cmd := m.theory.Update(msg)
		return cmd
	}
	return nil
}

func (m *RootModel) updateAudio(msg tea.Msg) tea.Cmd {
	if m.audio != nil {
		_, cmd := m.audio.Update(msg)
		return cmd
	}
	return nil
}

func (m *RootModel) updateManager(msg tea.Msg) tea.Cmd {
	if m.manager != nil {
		_, cmd := m.manager.Update(msg)
		return cmd
	}
	return nil
}

func (m *RootModel) updateSettings(msg tea.Msg) tea.Cmd {
	if m.settings != nil {
		_, cmd := m.settings.Update(msg)
		return cmd
	}
	return nil
}

func (m *RootModel) updateExport(msg tea.Msg) tea.Cmd {
	if m.export != nil {
		// Handle special messages for export screen
		switch msg.(type) {
		case ExportCompleteMsg:
			// Export completed, return to editor
			m.currentScreen = screenEditor
			return nil
		case BackMsg:
			// User pressed back, return to menu
			m.currentScreen = screenMenu
			return nil
		}

		_, cmd := m.export.Update(msg)
		return cmd
	}
	return nil
}

func (m *RootModel) renderExport() string {
	// Set export content from current editor content if available
	if m.editor != nil {
		content := m.editor.GetEditorText()
		m.export.SetContent(content)
	}

	// Set dimensions for responsive layout
	m.export.SetDimensions(m.width, m.height)

	return m.export.View()
}

// Collaboration event handlers

func (m *RootModel) onSessionUpdate(event collaboration.SessionUpdateEvent) {
	// Handle session updates - could trigger UI refresh or notifications
}

func (m *RootModel) onPresenceUpdate(event collaboration.PresenceUpdateEvent) {
	// Handle presence updates - could update presence indicators
}

func (m *RootModel) onConflictDetected(event collaboration.ConflictEvent) {
	// Handle conflict detection - could show conflict resolution UI
}

func (m *RootModel) onUserJoined(session *collaboration.Session, participant *collaboration.Participant) {
	// Handle user joining session
}

func (m *RootModel) onUserLeft(session *collaboration.Session, participant *collaboration.Participant) {
	// Handle user leaving session
}

func (m *RootModel) onSessionCreated(session *collaboration.Session) {
	// Handle session creation
}

func (m *RootModel) onSessionEnded(session *collaboration.Session) {
	// Handle session ending
}
