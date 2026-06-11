package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/internal/wizards"
	"github.com/kyanite/config"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"
)

type UnifiedAction int

const (
	ActionAddTask UnifiedAction = iota
	ActionListTasks
	ActionWizardTask
	ActionInteractiveTask
	ActionFilterTasks
	ActionNotes
	ActionChat
	ActionInspire
	ActionConfig
	ActionEnhancedConfig
	ActionTheme
	ActionStats
	ActionHelp
	ActionQuit
)

type UnifiedDashboardModel struct {
	width         int
	height        int
	taskEngine    *engine.Engine
	config        *config.FocusConfig
	currentView   string
	menuItems     []list.Item
	taskList      []models.Task
	menu          list.Model
	keyMap        UnifiedKeyMap
	statusMessage string
}

type UnifiedKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Enter   key.Binding
	Back    key.Binding
	Help    key.Binding
	Quit    key.Binding
	Config  key.Binding
	Theme   key.Binding
	Refresh key.Binding
}

var (
	menuTitleStyle    lipgloss.Style
	selectedItemStyle lipgloss.Style
	normalItemStyle   lipgloss.Style
	statusStyle       lipgloss.Style
	errorStyle        lipgloss.Style
)

func DefaultUnifiedKeyMap() UnifiedKeyMap {
	return UnifiedKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "back"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l", "enter"),
			key.WithHelp("→/l/↵", "select"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("↵/space", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("escape", "q"),
			key.WithHelp("esc/q", "back"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+q"),
			key.WithHelp("Ctrl+Q", "quit"),
		),
		Config: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "config"),
		),
		Theme: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "theme"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
	}
}

func NewUnifiedDashboardModel() UnifiedDashboardModel {
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	taskEngine := engine.New(repo)

	root, _ := config.Load()
	var cfg *config.FocusConfig
	if root != nil {
		focus := root.Focus
		cfg = &focus
	}

	initial := styles.ThemeSynthwave
	if cfg != nil {
		switch cfg.Theme {
		case string(styles.ThemeLight):
			initial = styles.ThemeLight
		case string(styles.ThemePlain):
			initial = styles.ThemePlain
		}
	}
	applyThemeStyles(initial)

	menuItems := []list.Item{
		UnifiedMenuItem{title: "🎯 Add Task", description: "Create a new task", action: ActionAddTask},
		UnifiedMenuItem{title: "📋 List Tasks", description: "View all tasks", action: ActionListTasks},
		UnifiedMenuItem{title: "🧙‍♂️ Task Wizard", description: "Advanced task creation", action: ActionWizardTask},
		UnifiedMenuItem{title: "🎯 Interactive Task", description: "Gum-powered task creation", action: ActionInteractiveTask},
		UnifiedMenuItem{title: "🔍 Filter Tasks", description: "Search and filter tasks", action: ActionFilterTasks},
		UnifiedMenuItem{title: "📝 Notes", description: "Manage task notes", action: ActionNotes},
		UnifiedMenuItem{title: "🤖 AI Chat", description: "Get AI assistance", action: ActionChat},
		UnifiedMenuItem{title: "🔮 Get Inspired", description: "AI task suggestions", action: ActionInspire},
		UnifiedMenuItem{title: "⚙️ Configuration", description: "Basic settings", action: ActionConfig},
		UnifiedMenuItem{title: "🔧 Enhanced Config", description: "Advanced configuration", action: ActionEnhancedConfig},
		UnifiedMenuItem{title: "🎨 Theme", description: "Change visual theme", action: ActionTheme},
		UnifiedMenuItem{title: "📊 Statistics", description: "View productivity stats", action: ActionStats},
		UnifiedMenuItem{title: "❓ Help", description: "Show help information", action: ActionHelp},
		UnifiedMenuItem{title: "🚪 Quit", description: "Exit focus.sh", action: ActionQuit},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedItemStyle
	delegate.Styles.SelectedDesc = selectedItemStyle
	delegate.Styles.NormalTitle = normalItemStyle
	delegate.Styles.NormalDesc = normalItemStyle

	menu := list.New(menuItems, delegate, 0, 0)
	menu.Title = "🌌 focus.sh Unified Dashboard"
	menu.SetShowStatusBar(false)
	menu.SetFilteringEnabled(false)
	menu.SetShowPagination(false)
	menu.SetShowHelp(false)

	return UnifiedDashboardModel{
		taskEngine:  taskEngine,
		config:      cfg,
		currentView: "menu",
		menuItems:   menuItems,
		menu:        menu,
		keyMap:      DefaultUnifiedKeyMap(),
	}
}

type UnifiedMenuItem struct {
	title       string
	description string
	action      UnifiedAction
}

func (i UnifiedMenuItem) Title() string       { return i.title }
func (i UnifiedMenuItem) Description() string { return i.description }
func (i UnifiedMenuItem) FilterValue() string { return i.title }

func (m UnifiedDashboardModel) Init() tea.Cmd {
	return tea.Batch(
		tea.Cmd(func() tea.Msg { return LoadConfigMsg{} }),
		tea.Cmd(func() tea.Msg { return RefreshTasksMsg{} }),
	)
}

func (m UnifiedDashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menu.SetWidth(msg.Width - 4)
		m.menu.SetHeight(msg.Height - 10)

	case LoadConfigMsg:
		root, err := config.Load()
		if err == nil {
			focus := root.Focus
			m.config = &focus
			th := styles.ThemeSynthwave
			switch focus.Theme {
			case string(styles.ThemeLight):
				th = styles.ThemeLight
			case string(styles.ThemePlain):
				th = styles.ThemePlain
			}
			m.applyTheme(th)
			m.statusMessage = "Configuration loaded"
		} else {
			m.statusMessage = "Failed to load configuration"
		}

	case RefreshTasksMsg:
		tasks, err := m.taskEngine.ListTasks("all")
		if err == nil {
			m.taskList = tasks
			m.statusMessage = fmt.Sprintf("Loaded %d tasks", len(tasks))
		} else {
			m.statusMessage = "Failed to load tasks"
		}

	case UnifiedActionMsg:
		return m, m.executeAction(msg.Action)

	case ThemeCycleMsg:
		if m.config != nil && m.config.Theme != "" {
			styles.SetThemeByName(m.config.Theme)
		}

		// Cycle to next theme
		styles.CycleTheme()
		next := styles.GetTheme()

		m.applyThemeString(next.Name)
		if m.config != nil {
			m.config.Theme = next.Name
			// Persist: load full root, replace Focus, save.
			if root, err := config.Load(); err == nil {
				root.Focus = *m.config
				_ = config.Save(root)
			}
		}
		m.statusMessage = fmt.Sprintf("Theme changed to %s", next.Name)

	case StatusMsg:
		m.statusMessage = msg.Message
	}

	if m.currentView == "menu" {
		m.menu, cmd = m.menu.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Help):
			m.currentView = "help"
		case key.Matches(msg, m.keyMap.Config):
			return m, tea.Cmd(func() tea.Msg { return UnifiedActionMsg{Action: ActionConfig} })
		case key.Matches(msg, m.keyMap.Theme):
			return m, tea.Cmd(func() tea.Msg { return UnifiedActionMsg{Action: ActionTheme} })
		case key.Matches(msg, m.keyMap.Refresh):
			return m, tea.Cmd(func() tea.Msg { return RefreshTasksMsg{} })
		case key.Matches(msg, m.keyMap.Back):
			if m.currentView != "menu" {
				m.currentView = "menu"
				m.statusMessage = "Back to main menu"
			}
		case key.Matches(msg, m.keyMap.Enter) && m.currentView == "menu":
			if selectedItem, ok := m.menu.SelectedItem().(UnifiedMenuItem); ok {
				return m, tea.Cmd(func() tea.Msg { return UnifiedActionMsg{Action: selectedItem.action} })
			}
		}
	}

	return m, cmd
}

func (m UnifiedDashboardModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var content strings.Builder

	header := fmt.Sprintf("%s\n%s", menuTitleStyle.Render("🌌 focus.sh Unified Dashboard"), strings.Repeat("─", m.width))
	content.WriteString(header)

	switch m.currentView {
	case "menu":
		content.WriteString(m.menu.View())
	case "help":
		content.WriteString(m.renderHelp())
	case "tasks":
		content.WriteString(m.renderTasks())
	}

	footer := fmt.Sprintf("\n%s\n%s", strings.Repeat("─", m.width), m.renderFooter())
	content.WriteString(footer)

	if m.statusMessage != "" {
		status := fmt.Sprintf("\n%s %s", statusStyle.Render("✓"), m.statusMessage)
		content.WriteString(status)
	}

	return content.String()
}

func (m UnifiedDashboardModel) executeAction(action UnifiedAction) tea.Cmd {
	return func() tea.Msg {
		switch action {
		case ActionAddTask:
			return StatusMsg{Message: "Launching add command..."}
		case ActionWizardTask:
			go func() {
				// T4-06: catch any panic from the wizard so a TUI crash doesn't kill the process.
				defer func() { _ = recover() }()
				err := wizards.TaskCreationWizard()
				if err != nil {
					m.statusMessage = fmt.Sprintf("Wizard error: %v", err)
				} else {
					m.statusMessage = "Task created successfully"
				}
			}()
			return StatusMsg{Message: "Launching task wizard..."}
		case ActionInteractiveTask:
			return StatusMsg{Message: "🎯 Interactive task creation - Use: focus interactive"}
		case ActionFilterTasks:
			return StatusMsg{Message: "🔍 Task filtering - Use: focus filter"}
		case ActionNotes:
			return StatusMsg{Message: "📝 Notes manager - Use: focus notes [task-id]"}
		case ActionChat:
			return StatusMsg{Message: "🤖 AI chat - Use: focus chat"}
		case ActionInspire:
			return StatusMsg{Message: "🔮 AI inspiration - Use: focus inspire"}
		case ActionConfig:
			return StatusMsg{Message: "🔧 Configuration wizard - Use: focus enhanced-config"}
		case ActionEnhancedConfig:
			return StatusMsg{Message: "Enhanced configuration coming soon..."}
		case ActionTheme:
			return ThemeCycleMsg{}
		case ActionStats:
			return StatusMsg{Message: "Loading statistics..."}
		case ActionHelp:
			m.currentView = "help"
			return StatusMsg{Message: "Showing help"}
		case ActionQuit:
			return tea.Quit()
		default:
			return StatusMsg{Message: fmt.Sprintf("Action %d not implemented", action)}
		}
	}
}

func applyThemeStyles(theme styles.ThemeMode) {
	styles.SetTheme(theme)
	menuTitleStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Background(styles.GetBackground()).
		Padding(0, 2).
		Bold(true)
	selectedItemStyle = lipgloss.NewStyle().
		Foreground(styles.GetBackground()).
		Background(styles.GetAccent()).
		Bold(true)
	normalItemStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent())
	statusStyle = lipgloss.NewStyle().
		Foreground(styles.GetSuccess()).
		Bold(true)
	errorStyle = lipgloss.NewStyle().
		Foreground(styles.GetError()).
		Bold(true)
}

func (m *UnifiedDashboardModel) applyTheme(theme styles.ThemeMode) {
	applyThemeStyles(theme)
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedItemStyle
	delegate.Styles.SelectedDesc = selectedItemStyle
	delegate.Styles.NormalTitle = normalItemStyle
	delegate.Styles.NormalDesc = normalItemStyle
	selected := m.menu.Index()
	newMenu := list.New(m.menuItems, delegate, 0, 0)
	newMenu.Title = "🌌 focus.sh Unified Dashboard"
	newMenu.SetShowStatusBar(false)
	newMenu.SetFilteringEnabled(false)
	newMenu.SetShowPagination(false)
	newMenu.SetShowHelp(false)
	newMenu.SetWidth(m.width - 4)
	newMenu.SetHeight(m.height - 10)
	if selected >= 0 {
		newMenu.Select(selected)
	}
	m.menu = newMenu
}

func applyThemeStylesByName(themeName string) {
	styles.SetThemeByName(themeName)
	menuTitleStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Background(styles.GetBackground()).
		Padding(0, 2).
		Bold(true)
	selectedItemStyle = lipgloss.NewStyle().
		Foreground(styles.GetBackground()).
		Background(styles.GetAccent()).
		Bold(true)
	normalItemStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent())
	statusStyle = lipgloss.NewStyle().
		Foreground(styles.GetSuccess()).
		Bold(true)
	errorStyle = lipgloss.NewStyle().
		Foreground(styles.GetError()).
		Bold(true)
}

func (m *UnifiedDashboardModel) applyThemeString(themeName string) {
	applyThemeStylesByName(themeName)
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedItemStyle
	delegate.Styles.SelectedDesc = selectedItemStyle
	delegate.Styles.NormalTitle = normalItemStyle
	delegate.Styles.NormalDesc = normalItemStyle
	selected := m.menu.Index()
	newMenu := list.New(m.menuItems, delegate, 0, 0)
	newMenu.Title = "🌌 focus.sh Unified Dashboard"
	newMenu.SetShowStatusBar(false)
	newMenu.SetFilteringEnabled(false)
	newMenu.SetShowPagination(false)
	newMenu.SetShowHelp(false)
	newMenu.SetWidth(m.width - 4)
	newMenu.SetHeight(m.height - 10)
	if selected >= 0 {
		newMenu.Select(selected)
	}
	m.menu = newMenu
}

func (m UnifiedDashboardModel) renderHelp() string {
	helpText := `
🌌 focus.sh Unified Dashboard - Help

🎯 TASK MANAGEMENT
  Add Task        → Create new tasks with AI parsing
  List Tasks      → View and manage existing tasks
  Task Wizard     → Advanced form-based task creation
  Interactive     → Gum-powered interactive task creation
  Filter Tasks    → Search and filter tasks dynamically

🤖 AI FEATURES
  AI Chat         → Get intelligent assistance
  Get Inspired    → AI-powered task suggestions

⚙️ CONFIGURATION
  Configuration   → Basic settings management
  Enhanced Config → Advanced configuration wizard
  Theme           → Switch visual themes

📊 UTILITIES
  Statistics      → View productivity metrics
  Help            → Show this help screen

🎹 KEYBOARD SHORTCUTS
  ↑/k            → Move up
  ↓/j            → Move down
  Enter/Space    → Select item
  c              → Quick config access
  t              → Change theme
  r              → Refresh data
  Esc/q          → Go back
  ?              → Show help
  Ctrl+Q         → Quit application
`
	return helpText
}

func (m UnifiedDashboardModel) renderTasks() string {
	if len(m.taskList) == 0 {
		return "No tasks found. Press Esc to return to menu."
	}
	var content strings.Builder
	content.WriteString("📋 TASK OVERVIEW\n\n")
	for _, task := range m.taskList {
		status := "⏳"
		if task.Status == "completed" {
			status = "✅"
		}
		priority := "🟢"
		switch task.Priority {
		case "high":
			priority = "🔴"
		case "medium":
			priority = "🟡"
		}
		line := fmt.Sprintf("%s %s %s %s\n", status, priority, task.Description, strings.Join(task.Categories, ", "))
		content.WriteString(line)
	}
	return content.String()
}

func (m UnifiedDashboardModel) renderFooter() string {
	var configInfo strings.Builder
	if m.config != nil {
		fmt.Fprintf(&configInfo, "🤖 %s | ", m.config.AI.Provider)
		fmt.Fprintf(&configInfo, "🎨 %s | ", m.config.Theme)
	}
	taskCount := fmt.Sprintf("📋 %d tasks", len(m.taskList))
	return fmt.Sprintf("%s%s | %s | [?]Help [c]Config [Ctrl+Shift+T]Theme [r]Refresh [Ctrl+Q]Quit", configInfo.String(), taskCount, m.currentView)
}

type UnifiedActionMsg struct{ Action UnifiedAction }

type LoadConfigMsg struct{}

type RefreshTasksMsg struct{}

type StatusMsg struct{ Message string }

type ThemeCycleMsg struct{}

func StartUnifiedDashboard() error {
	model := NewUnifiedDashboardModel()
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
