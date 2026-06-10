package app

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/prism/internal/ai"
	"github.com/kyanite/prism/internal/theme"
	"github.com/kyanite/prism/internal/ui"
)

// Screen represents different screens in the app
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenWheel
	ScreenGenerator
	ScreenTheory
	ScreenChecker
	ScreenManager
	ScreenSearch
	ScreenHelp
)

// Model is the root Bubble Tea model
type Model struct {
	CurrentScreen Screen
	Width         int
	Height        int
	ThemeManager  *theme.Manager
	SessionID     string
	lastPalette   string
	aiClient      *ai.Client

	// Screen models
	menuModel      ui.MenuModel
	wheelModel     ui.WheelModel
	generatorModel ui.GeneratorModel
	theoryModel    ui.TheoryModel
	checkerModel   ui.CheckerModel
	managerModel   ui.ManagerModel
	searchModel    ui.SearchModel

	showHelp bool
}

// NewModel creates a new root model
func NewModel() Model {
	themeManager := theme.NewManager()
	aiClient := ai.NewClient()
	sessionID := fmt.Sprintf("prism-%d", time.Now().Unix())

	return Model{
		CurrentScreen:  ScreenMenu,
		ThemeManager:   themeManager,
		SessionID:      sessionID,
		aiClient:       aiClient,
		menuModel:      ui.NewMenuModel(themeManager),
		wheelModel:     ui.NewWheelModel(themeManager),
		generatorModel: ui.NewGeneratorModel(themeManager),
		theoryModel:    ui.NewTheoryModel(themeManager),
		checkerModel:   ui.NewCheckerModel(themeManager),
		managerModel:   ui.NewManagerModel(themeManager),
		searchModel:    ui.NewSearchModel(themeManager),
	}
}

// LoadRecentSession attempts to load and print the most recent session.
// Best-effort: failures are silently ignored.
func (m *Model) LoadRecentSession() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sessions, err := m.aiClient.GetRecentSessions(ctx, 1)
	if err != nil || len(sessions) == 0 {
		return
	}

	s := sessions[0]
	if s.Title != "" {
		fmt.Printf("Recent palettes: %s\n", s.Title)
		m.lastPalette = s.Title
	}
}

// SaveCurrentSession persists the current session state before exit.
// Best-effort: failures are silently ignored.
func (m *Model) SaveCurrentSession() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	title := m.lastPalette
	if title == "" {
		title = fmt.Sprintf("session %s", time.Now().Format("2006-01-02"))
	}

	// Best-effort session save
	_ = m.aiClient.SaveSession(ctx, m.SessionID, title, nil)

	// Best-effort cross-app context for other kyanite apps
	summary := fmt.Sprintf("Used palette '%s' in prism", title)
	_ = m.aiClient.SaveCrossAppContext(ctx, "noise", "palette_usage", summary, 0.5)
	_ = m.aiClient.SaveCrossAppContext(ctx, "focus", "palette_usage", summary, 0.5)
	_ = m.aiClient.SaveCrossAppContext(ctx, "syntax", "palette_usage", summary, 0.5)
}

// SetLastPalette records the most recently generated/saved palette name.
func (m *Model) SetLastPalette(name string) {
	m.lastPalette = name
}

// Close releases resources held by the model.
func (m *Model) Close() {
	m.aiClient.Close()
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeys(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case ui.NavigateMsg:
		m.CurrentScreen = Screen(msg.Screen)
		return m, nil
	}

	// Route to current screen
	return m.routeToScreen(msg)
}

// View renders the UI
func (m Model) View() string {
	if m.showHelp {
		return ui.RenderHelp(m.ThemeManager, m.Width, m.Height)
	}

	switch m.CurrentScreen {
	case ScreenMenu:
		return m.menuModel.View()
	case ScreenWheel:
		return m.wheelModel.View()
	case ScreenGenerator:
		return m.generatorModel.View()
	case ScreenTheory:
		return m.theoryModel.View()
	case ScreenChecker:
		return m.checkerModel.View()
	case ScreenManager:
		return m.managerModel.View()
	case ScreenSearch:
		return m.searchModel.View()
	default:
		return "Unknown screen"
	}
}

// handleKeys handles global keyboard shortcuts
func (m Model) handleKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global shortcuts
	switch msg.String() {
	case "ctrl+q":
		return m, tea.Quit

	case "ctrl+h":
		m.showHelp = !m.showHelp
		return m, nil

	case "ctrl+shift+t":
		m.ThemeManager.NextTheme()
		return m, nil

	case "esc":
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		if m.CurrentScreen != ScreenMenu {
			m.CurrentScreen = ScreenMenu
			return m, nil
		}
		return m, tea.Quit
	}

	return m, nil
}

// routeToScreen routes messages to the appropriate screen model
func (m Model) routeToScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.CurrentScreen {
	case ScreenMenu:
		m.menuModel, cmd = m.menuModel.Update(msg)
	case ScreenWheel:
		m.wheelModel, cmd = m.wheelModel.Update(msg)
	case ScreenGenerator:
		m.generatorModel, cmd = m.generatorModel.Update(msg)
	case ScreenTheory:
		m.theoryModel, cmd = m.theoryModel.Update(msg)
	case ScreenChecker:
		m.checkerModel, cmd = m.checkerModel.Update(msg)
	case ScreenManager:
		m.managerModel, cmd = m.managerModel.Update(msg)
	case ScreenSearch:
		m.searchModel, cmd = m.searchModel.Update(msg)
	}

	return m, cmd
}
