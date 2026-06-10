package ui

import (
	"fmt"
	"time"

	"github.com/kyanite/noise/internal/app"
	"github.com/kyanite/noise/internal/collaboration"
	"github.com/kyanite/noise/internal/config"
	errutil "github.com/kyanite/noise/internal/errutil"
	"github.com/kyanite/noise/internal/infra/db"
	"github.com/kyanite/noise/internal/infra/sync"
	"github.com/kyanite/noise/internal/logging"
	"github.com/kyanite/noise/internal/plugins"
	"github.com/kyanite/noise/internal/theme"
	"github.com/kyanite/noise/internal/ui/dashboard"
	"github.com/kyanite/noise/internal/ui/editor"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// QuickStartConfig contains configuration for quick start mode
type QuickStartConfig struct {
	Theme          string
	ScratchMode    bool
	AutoBrainstorm bool
}

// Screen represents different screens in the application
type screen int

const (
	screenSplash screen = iota
	screenMenu
	screenDashboard
	screenEditor
	screenExport
	screenTheory
	screenAudio
	screenManager
	screenSettings
	screenLoading
)

// ForceRefreshMsg triggers a complete re-render
type ForceRefreshMsg struct{}

// RefreshTimerMsg triggers periodic refresh
type RefreshTimerMsg struct{}

// VoiceTranscriptionMsg carries the result of voice-to-text transcription
type VoiceTranscriptionMsg struct {
	Text      string
	Cancelled bool
	Error     error
}

// RootModel is the main application model that handles routing between screens
type RootModel struct {
	// Current state
	currentScreen screen
	width         int
	height        int

	// Quit confirmation dialog state
	confirmingQuit bool

	// Database connection
	database *db.DB

	// Configuration
	config *config.Config

	// Child models
	splash    *SplashModel
	menu      *MenuModel
	dashboard interface {
		Init() tea.Cmd
		Update(tea.Msg) tea.Cmd
		View() string
	}
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

	// Quick start configuration
	quickStartConfig *QuickStartConfig

	// AI Service
	aiService *app.AIService

	// Voice-to-text service
	voiceService   *app.VoiceService
	voiceIndicator *VoiceIndicatorModel

	// PWA sync service
	syncServer *sync.SyncServer
	syncStatus *SyncStatusModel

	// Toast notification system
	toast *ToastModel
}

// NewRootModel creates a new root model with initialized state
func NewRootModel(pluginManager *plugins.DefaultManager) *RootModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.GetManager().Current().Primary)

	return &RootModel{
		currentScreen:     screenSplash,
		loading:           true,
		spinner:           s,
		animation:         NewAnimationManager(),
		responsiveManager: NewResponsiveLayoutManager(),
		// Collaboration managers are initialized in initializeCollaborationSystem
		// only when the feature flag is enabled (config.Features.EnableCollaboration)
		pluginManager: pluginManager,
	}
}

// Init initializes the root model
func (m *RootModel) Init() tea.Cmd {
	// Enable TUI mode to suppress log output that would corrupt the terminal
	logging.EnableTUIMode()

	return tea.Batch(
		tea.EnterAltScreen,
		m.initializeApp(),
		m.spinner.Tick,
		// Use tea.Tick for periodic refresh instead of blocking sleep
		m.scheduleRefresh(),
	)
}

// scheduleRefresh returns a command that schedules the next refresh tick
// Uses tea.Tick to avoid blocking the event loop
func (m *RootModel) scheduleRefresh() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return RefreshTimerMsg{}
	})
}

// initializeApp initializes the application (database, secrets, AI service)
func (m *RootModel) initializeApp() tea.Cmd {
	return func() tea.Msg {
		// Initialize database
		database, err := db.New(db.Config{})
		if err != nil {
			return initErrorMsg{err: errutil.Wrap(err, "initialize database")}
		}

		// Initialize AI Service
		cfg, loadErr := config.Load()
		if loadErr != nil {
			logging.Warnf("Failed to load config for AI service initialization: %v", loadErr)
			cfg = config.DefaultConfig()
		}
		aiService := app.NewAIService(cfg)

		return initSuccessMsg{database: database, aiService: aiService}
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

		// CRITICAL: Forward WindowSizeMsg to ALL child models
		// This ensures all screens have proper dimensions when navigated to
		if m.dashboard != nil {
			if cmd := m.dashboard.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		if m.menu != nil {
			m.menu.Update(msg)
		}
		if m.theory != nil {
			m.theory.Update(msg)
		}
		if m.export != nil {
			m.export.Update(msg)
		}
		if m.audio != nil {
			m.audio.Update(msg)
		}
		if m.settings != nil {
			m.settings.Update(msg)
		}
		if m.editor != nil {
			m.editor.Update(msg)
		}
		if m.helpPane != nil {
			m.helpPane.Update(msg)
		}
		if m.manager != nil {
			m.manager.Update(msg)
		}

	case ForceRefreshMsg:
		// Force a complete re-render
		return m, nil

	case dashboard.ForceRefreshMsg:
		// Handle force refresh from dashboard (different type due to package boundary)
		return m, nil

	case RefreshTimerMsg:
		// Schedule next refresh using tea.Tick (non-blocking)
		cmds = append(cmds, m.scheduleRefresh())

	case tea.KeyMsg:
		// Handle quit confirmation dialog
		if m.confirmingQuit {
			switch msg.String() {
			case "y", "Y", "enter":
				// User confirmed quit
				return m, m.performQuit()
			case "n", "N", "esc":
				// User cancelled quit
				m.confirmingQuit = false
				return m, nil
			}
			// Ignore other keys while confirming
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			// Check for unsaved changes before quitting
			if m.hasUnsavedChanges() {
				m.confirmingQuit = true
				return m, nil
			}
			return m, m.performQuit()
		case "ctrl+d":
			// Global voice dictation toggle
			if m.voiceService != nil && m.voiceService.IsAvailable() {
				return m, m.handleVoiceDictation()
			}
		case "esc":
			// Cancel voice dictation if active
			if m.voiceService != nil && m.voiceService.GetState() == app.VoiceStateRecording {
				_ = m.voiceService.CancelDictation()
				return m, nil
			}
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
		default:
			// FOCUS TRAPPING: When help mode is active, route keys to help pane only
			if m.helpMode && m.helpPane != nil {
				_, helpCmd := m.helpPane.Update(msg)
				if helpCmd != nil {
					cmds = append(cmds, helpCmd)
				}
				// Don't propagate to underlying screen
				return m, tea.Batch(cmds...)
			}
			// After any key press, ensure we refresh
			cmds = append(cmds, func() tea.Msg { return ForceRefreshMsg{} })
		}

	case VoiceTranscriptionMsg:
		// Handle voice transcription result
		if msg.Text != "" && m.editor != nil && m.currentScreen == screenEditor {
			// Insert transcribed text at cursor in editor
			m.editor.InsertText(msg.Text)
		}
		if m.voiceIndicator != nil {
			m.voiceIndicator.SetState(app.VoiceStateIdle)
		}

	case initSuccessMsg:
		m.database = msg.database
		m.aiService = msg.aiService
		m.loading = false

		// If quick start is configured, go directly to editor
		// Otherwise, start with the dashboard as the primary interface
		if m.quickStartConfig != nil {
			m.currentScreen = screenEditor
		} else {
			m.currentScreen = screenDashboard
		}

		// Initialize collaboration system
		m.initializeCollaborationSystem()

		// Initialize child models
		m.initializeChildModels()

		// CRITICAL: Call Init() on child models and collect their commands
		// Without this, animations and background tasks won't start
		childCmds := m.initializeChildModelCommands()
		cmds = append(cmds, childCmds...)

		// If quick start is configured, configure the editor
		if m.quickStartConfig != nil {
			m.configureQuickStart()
		}

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
			logging.Warn("Size validation issues detected for terminal size")
		}

	case ScreenChangeMsg:
		// Handle screen changes from menu
		m.currentScreen = msg.Screen
		// Focus the export screen when entering it
		if msg.Screen == screenExport && m.export != nil {
			m.export.Focus()
		}

	case dashboard.ScreenChangeMsg:
		// Handle screen changes from dashboard (different type due to package boundary)
		m.currentScreen = screen(msg.Screen)
		// Focus the export screen when entering it
		if screen(msg.Screen) == screenExport && m.export != nil {
			m.export.Focus()
		}

	case editor.ScreenChangeMsg:
		// Handle screen changes from editor/split_pane (different type due to package boundary)
		m.currentScreen = screen(msg.Screen)
		// Focus the export screen when entering it
		if screen(msg.Screen) == screenExport && m.export != nil {
			m.export.Focus()
		}

	case dashboard.AnimationTickMsg:
		// Handle animation ticks from dashboard (different type due to package boundary)
		// Forward to dashboard for animation updates
		if m.dashboard != nil {
			if cmd := m.dashboard.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case dashboard.TriggerBrainstormMsg:
		// Handle AI brainstorm trigger from dashboard
		// Switch to editor and start brainstorm mode
		m.currentScreen = screenEditor
		if m.editor != nil {
			// Use provided theme or derive from context
			theme := msg.Theme
			if theme == "" {
				theme = "creative inspiration"
			}
			m.editor.StartRapidBrainstorm(theme)
		}
		return m, nil

	case ToastMsg:
		// Handle global toast notifications
		if m.toast != nil {
			cmd := m.toast.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case toastDismissMsg:
		// Handle toast dismissal
		if m.toast != nil {
			cmd := m.toast.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
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

	// Initialize theme system from config
	if cfg.UI.Theme != "" {
		theme.GetManager().SetTheme(cfg.UI.Theme)
	}

	m.splash = NewSplashModel()
	m.menu = NewMenuModel()
	m.editor = NewEditorModel(m.database, m.aiService)
	m.export = NewExportModel("") // Content will be set when entering export screen
	m.theory = NewTheoryModel()
	m.audio = NewAudioModel()
	m.manager = NewManagerModel(m.database)
	m.settings = NewSettingsModel(m.config)

	// Initialize dashboard
	m.dashboard = dashboard.NewDashboardModel()

	// Pass AI service and database to dashboard for its panels
	if dashboardModel, ok := m.dashboard.(*dashboard.DashboardModel); ok {
		dashboardModel.SetAIService(m.aiService)
		dashboardModel.SetDatabase(m.database)
	}

	// CRITICAL: Pass current dimensions to ALL child models immediately after creation
	// This ensures all screens have valid dimensions when navigated to
	if m.width > 0 && m.height > 0 {
		sizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
		if m.dashboard != nil {
			m.dashboard.Update(sizeMsg)
		}
		if m.menu != nil {
			m.menu.Update(sizeMsg)
		}
		if m.theory != nil {
			m.theory.Update(sizeMsg)
		}
		if m.export != nil {
			m.export.Update(sizeMsg)
		}
		if m.audio != nil {
			m.audio.Update(sizeMsg)
		}
		if m.settings != nil {
			m.settings.Update(sizeMsg)
		}
		if m.editor != nil {
			m.editor.Update(sizeMsg)
		}
		if m.manager != nil {
			m.manager.Update(sizeMsg)
		}
	}

	// Initialize help system
	m.helpPane = editor.NewHelpPaneModel(nil)

	// Initialize toast notification system
	m.toast = NewToastModel()

	// Initialize voice service if enabled
	m.initializeVoiceService(cfg)

	// Initialize sync service if enabled (auto-creates directories)
	m.initializeSyncService(cfg)
}

// initializeChildModelCommands calls Init() on all child models and returns their commands
// CRITICAL: This must be called after initializeChildModels() to start animations and background tasks
func (m *RootModel) initializeChildModelCommands() []tea.Cmd {
	var cmds []tea.Cmd

	// Initialize editor (includes split pane, editor pane, preview pane)
	if m.editor != nil {
		if cmd := m.editor.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Initialize manager (starts project loading and spinner)
	if m.manager != nil {
		if cmd := m.manager.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Initialize dashboard
	if m.dashboard != nil {
		if cmd := m.dashboard.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Initialize other models that have Init() implementations
	if m.settings != nil {
		if cmd := m.settings.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if m.theory != nil {
		if cmd := m.theory.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if m.export != nil {
		if cmd := m.export.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if m.helpPane != nil {
		if cmd := m.helpPane.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return cmds
}

// initializeVoiceService initializes the voice-to-text service
// Automatically downloads the whisper model on first run if needed
func (m *RootModel) initializeVoiceService(cfg *config.Config) {
	if !cfg.Voice.Enabled {
		logging.Debug("Voice service disabled in config")
		return
	}

	// Create voice indicator
	m.voiceIndicator = NewVoiceIndicatorModel()

	// Create voice service with auto-setup (downloads model if needed)
	// This runs in background to avoid blocking app startup
	go func() {
		voiceService, err := app.NewVoiceServiceWithAutoSetup(cfg, logging.GetDefaultLogger(),
			func(status string, progress float64) {
				// Update indicator with download status
				if m.voiceIndicator != nil {
					if progress < 1.0 && progress > 0 {
						m.voiceIndicator.SetError(fmt.Sprintf("%.0f%% %s", progress*100, status))
					} else if progress >= 1.0 {
						m.voiceIndicator.ClearError()
					}
				}
				logging.Debugf("Voice setup: %s (%.0f%%)", status, progress*100)
			})
		if err != nil {
			logging.Warnf("Failed to create voice service: %v", err)
			if m.voiceIndicator != nil {
				m.voiceIndicator.SetError("Voice unavailable")
			}
			return
		}

		// Set up callbacks
		voiceService.OnStateChange(func(state app.VoiceState) {
			if m.voiceIndicator != nil {
				m.voiceIndicator.SetState(state)
			}
		})

		voiceService.OnLevelChange(func(level float32) {
			if m.voiceIndicator != nil {
				m.voiceIndicator.SetLevel(level)
			}
		})

		m.voiceService = voiceService

		// Initialize whisper engine (loads model)
		if err := voiceService.Initialize(); err != nil {
			logging.Warnf("Failed to initialize voice service: %v", err)
			if m.voiceIndicator != nil {
				m.voiceIndicator.SetError("Voice unavailable")
			}
		} else {
			logging.Info("Voice service ready")
			if m.voiceIndicator != nil {
				m.voiceIndicator.ClearError()
			}
		}
	}()
}

// initializeSyncService initializes the PWA sync service
// Automatically creates necessary directories
func (m *RootModel) initializeSyncService(cfg *config.Config) {
	if !cfg.Sync.Enabled {
		logging.Debug("Sync service disabled in config")
		return
	}

	// Create sync server with auto-setup (creates directories automatically)
	syncServer, err := sync.NewSyncServerWithAutoSetup(
		cfg.GetDataDir(),
		cfg.Sync.Port,
		logging.GetDefaultLogger(),
	)
	if err != nil {
		logging.Warnf("Failed to create sync server: %v", err)
		return
	}

	m.syncServer = syncServer
	m.syncStatus = NewSyncStatusModel(syncServer)

	// Auto-start if configured
	if cfg.Sync.AutoStart {
		if err := syncServer.Start(); err != nil {
			logging.Warnf("Failed to auto-start sync server: %v", err)
		} else {
			logging.Infof("Sync server started at %s", syncServer.GetLocalURL())
		}
	}
}

// initializeCollaborationSystem initializes the collaboration system
// Only initializes if collaboration feature is enabled in config
func (m *RootModel) initializeCollaborationSystem() {
	// Check if collaboration feature is enabled
	if m.config == nil || !m.config.IsCollaborationEnabled() {
		// Collaboration is disabled - skip initialization
		// This is the default for single-user mode
		return
	}

	// Initialize collaboration managers
	m.collaborationManager = collaboration.NewCollaborationManager(m.database)
	m.presenceManager = collaboration.NewPresenceManager()
	m.sessionManager = collaboration.NewSessionManager()
	m.invitationManager = collaboration.NewInvitationManager()
	m.conflictResolver = collaboration.NewConflictResolver()

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
	case screenDashboard:
		return m.updateDashboard(msg)
	case screenSettings:
		return m.updateSettings(msg)
	}
	return nil
}

// View renders the current screen
func (m *RootModel) View() string {
	// Show quit confirmation dialog if active
	if m.confirmingQuit {
		return m.renderQuitConfirmation()
	}

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
	case screenDashboard:
		content = m.renderDashboard()
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

	// CRITICAL: Constrain all output to terminal dimensions
	// This prevents rendering artifacts from content overflow
	if m.width > 0 && m.height > 0 {
		constrainStyle := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			MaxWidth(m.width).
			MaxHeight(m.height)
		content = constrainStyle.Render(content)
	}

	// Render toast notifications as overlay (if any)
	if m.toast != nil {
		toastView := m.toast.View()
		if toastView != "" {
			// Position toasts at the top of the screen
			toastStyle := lipgloss.NewStyle().
				MarginTop(1).
				Align(lipgloss.Right).
				Width(m.width)
			content = lipgloss.JoinVertical(lipgloss.Left,
				toastStyle.Render(toastView),
				content,
			)
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
		Foreground(theme.GetManager().Current().Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width).
		Height(m.height)

	loadingText := "[~] noise.sh [~]\n\n"
	if m.errorMsg != "" {
		loadingText += "Error: " + m.errorMsg + "\n\nPress any key to exit..."
	} else {
		loadingText += "Initializing...\n\n" + m.spinner.View()
	}

	return loadingStyle.Render(loadingText)
}

func (m *RootModel) renderDashboard() string {
	if m.dashboard != nil {
		return m.dashboard.View()
	}
	return "Dashboard loading..."
}

// renderQuitConfirmation renders the quit confirmation dialog
func (m *RootModel) renderQuitConfirmation() string {
	t := theme.GetManager().Current()

	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Warning).
		Background(t.Background).
		Padding(1, 3).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Foreground(t.Warning).
		Bold(true)

	textStyle := lipgloss.NewStyle().
		Foreground(t.Text)

	hintStyle := lipgloss.NewStyle().
		Foreground(t.Secondary).
		MarginTop(1)

	title := titleStyle.Render("[!] Unsaved Changes")
	message := textStyle.Render("You have unsaved changes that will be lost.")
	hint := hintStyle.Render("[Y] Quit anyway  [N/Esc] Cancel")

	dialog := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		message,
		"",
		hint,
	)

	dialogBox := dialogStyle.Render(dialog)

	// Center on screen
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialogBox,
	)
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
	database  *db.DB
	aiService *app.AIService
}

type initErrorMsg struct {
	err error
}

// ScreenChangeMsg requests a change to the given screen.
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

func (m *RootModel) updateDashboard(msg tea.Msg) tea.Cmd {
	if m.dashboard != nil {
		cmd := m.dashboard.Update(msg)
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

// handleVoiceDictation toggles voice dictation recording
func (m *RootModel) handleVoiceDictation() tea.Cmd {
	if m.voiceService == nil {
		return nil
	}

	state := m.voiceService.GetState()

	switch state {
	case app.VoiceStateIdle:
		// Start recording
		if err := m.voiceService.StartDictation(); err != nil {
			logging.Errorf("Failed to start dictation: %v", err)
			if m.voiceIndicator != nil {
				m.voiceIndicator.SetError("Mic error")
			}
			return nil
		}
		if m.voiceIndicator != nil {
			m.voiceIndicator.SetState(app.VoiceStateRecording)
		}
		// Return a command that updates the duration periodically
		return m.voiceDurationTick()

	case app.VoiceStateRecording:
		// Stop recording and transcribe
		if m.voiceIndicator != nil {
			m.voiceIndicator.SetState(app.VoiceStateProcessing)
		}
		return func() tea.Msg {
			text, err := m.voiceService.StopDictation()
			if err != nil {
				logging.Errorf("Failed to stop dictation: %v", err)
				return VoiceTranscriptionMsg{Error: err}
			}
			return VoiceTranscriptionMsg{Text: text}
		}

	default:
		return nil
	}
}

// voiceDurationTick returns a command that updates the voice indicator duration
func (m *RootModel) voiceDurationTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		if m.voiceService != nil && m.voiceService.GetState() == app.VoiceStateRecording {
			if m.voiceIndicator != nil {
				m.voiceIndicator.SetDuration(m.voiceService.GetDuration())
			}
			// Check max duration
			if m.voiceService.CheckMaxDuration() {
				// Auto-stop at max duration
				return tea.KeyMsg{Type: tea.KeyCtrlD}
			}
		}
		return nil
	})
}

// GetVoiceService returns the voice service for external access
func (m *RootModel) GetVoiceService() *app.VoiceService {
	return m.voiceService
}

// IsVoiceAvailable returns whether voice-to-text is available
func (m *RootModel) IsVoiceAvailable() bool {
	return m.voiceService != nil && m.voiceService.IsAvailable()
}

// SetQuickStart configures the application for quick start mode
func (m *RootModel) SetQuickStart(config *QuickStartConfig) {
	m.quickStartConfig = config
}

// configureQuickStart applies quick start configuration to the editor
func (m *RootModel) configureQuickStart() {
	if m.quickStartConfig == nil || m.editor == nil {
		return
	}

	// Configure editor for scratch mode
	if m.editor.GetSplitPane() != nil {
		// Set scratch mode - this will be implemented in the editor pane
		// For now, we'll use a placeholder method that will be added later
		splitPane := m.editor.GetSplitPane()
		if splitPane != nil {
			// Set scratch mode and trigger brainstorm if needed
			// These methods will be implemented in the editor package
			splitPane.SetQuickStartConfig(m.quickStartConfig.Theme, m.quickStartConfig.ScratchMode, m.quickStartConfig.AutoBrainstorm)
		}
	}
}

// hasUnsavedChanges checks if the editor has unsaved changes
func (m *RootModel) hasUnsavedChanges() bool {
	if m.editor == nil {
		return false
	}
	splitPane := m.editor.GetSplitPane()
	if splitPane == nil {
		return false
	}
	return splitPane.HasUnsavedChanges()
}

// performQuit performs the actual quit operation with cleanup
func (m *RootModel) performQuit() tea.Cmd {
	if m.database != nil {
		m.database.Close()
	}
	// Clean up animation manager
	if m.animation != nil {
		m.animation.Close()
	}
	// Clean up voice service
	if m.voiceService != nil {
		m.voiceService.Close()
	}
	// Clean up sync server
	if m.syncServer != nil {
		_ = m.syncServer.Stop()
	}
	// Restore normal logging before exit
	logging.DisableTUIMode()
	return tea.Quit
}
