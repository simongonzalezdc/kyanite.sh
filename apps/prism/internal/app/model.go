package app

import (
	tea "github.com/charmbracelet/bubbletea"
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

	return Model{
		CurrentScreen:  ScreenMenu,
		ThemeManager:   themeManager,
		menuModel:      ui.NewMenuModel(themeManager),
		wheelModel:     ui.NewWheelModel(themeManager),
		generatorModel: ui.NewGeneratorModel(themeManager),
		theoryModel:    ui.NewTheoryModel(themeManager),
		checkerModel:   ui.NewCheckerModel(themeManager),
		managerModel:   ui.NewManagerModel(themeManager),
		searchModel:    ui.NewSearchModel(themeManager),
	}
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
