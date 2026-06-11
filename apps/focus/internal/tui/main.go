package tui

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kyanite/focus/internal/ai"
	"github.com/kyanite/focus/internal/theme"
	"github.com/kyanite/tui/aipanel"
	"github.com/kyanite/focus/pkg/audio"
	"github.com/kyanite/focus/pkg/calendar"
	"github.com/kyanite/focus/pkg/glow"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
)

// Tick message for real-time updates
type tickMsg time.Time

// Spinner message for loading states
type spinnerTickMsg time.Time

// AI response message for chat
type aiResponseMsg struct {
	response  string
	userInput string
}

// AI status message for status updates from goroutines
type aiStatusMsg struct {
	status string
}

// Loading state
type loadingState int

const (
	notLoading loadingState = iota
	startingUp
	savingTask
	loadingCalendar
	processingAI
)

type MainModel struct {
	tasks        []DashboardTask
	activeTimer  *TimerSession
	timer        timer.Model
	stopwatch    stopwatch.Model
	help         help.Model
	keys         keyMap
	showHelp     bool
	quitting     bool
	workTime     time.Duration
	breakTime    time.Duration
	sessionType  sessionType
	sessions     int
	currentView  mainView
	selectedTask int
	width        int
	height       int

	// Calendar management
	cal             *calendar.Calendar
	calRenderer     *calendar.Renderer
	calSelectedDate time.Time
	calViewMode     string // "month", "week", "day"

	// Chat assistant components
	chatInput    string
	chatHistory  []string
	chatViewport viewport.Model

	// Task entry components
	taskEntryMode bool
	taskInput     string

	// Notes management
	notesMode   bool
	notesInput  string
	editingTask *DashboardTask

	// Enhanced filtering
	filterMode     bool
	filterStatus   string // "all", "pending", "completed"
	filterPriority string // "all", "high", "medium", "low"

	// Loading and spinner states
	loadingState loadingState
	spinnerFrame int

	// Glow styler for enhanced markdown rendering
	glowStyler *glow.GlowStyler

	// AI manager for real chat integration
	aiManager      *ai.Manager
	aiStatus       string // "online", "offline", "checking"
	lastAICheck    time.Time
	aiThinking     bool // Whether AI is currently responding
	aiSpinnerFrame int  // Current spinner frame for AI response

	// AI Dashboard panel
	aiPanel       aipanel.Model
	aiPanelInput  string // text input for AI panel queries
	aiPanelActive bool   // whether AI panel input is focused

	// Visual effects
	glitchCount int

	glitchMessage string

	// Theme management

	// Settings management
	settingsMode  bool
	audioEnabled  bool
	workDuration  time.Duration
	breakDuration time.Duration

	// Mutex for protecting concurrent access to shared fields
	mu sync.RWMutex

	// Calendar management (temporarily disabled)
	// calendarView     calendar.ViewMode
	// calendar         *calendar.Calendar
	// calendarRenderer *calendar.Renderer
}

type DashboardTask struct {
	ID          string
	Description string
	Status      string
	Priority    string
	CreatedAt   time.Time
	Deadline    *time.Time // Add deadline support
	Categories  []string
	Notes       string // Add notes support
}

type TimerSession struct {
	Task      DashboardTask
	Mode      string // "work" or "break"
	Duration  time.Duration
	TimeLeft  time.Duration
	StartTime time.Time
}

type sessionType int

const (
	workSession sessionType = iota
	breakSession
)

type mainView int

const (
	dashboardView mainView = iota
	focusView
	chatView
	calendarView
	notesView
	journalView
	settingsView
)

type keyMap struct {
	start         key.Binding
	stop          key.Binding
	reset         key.Binding
	quit          key.Binding
	help          key.Binding
	switchSession key.Binding
	up            key.Binding
	down          key.Binding
	enter         key.Binding
	tab           key.Binding
	focusMode     key.Binding
	chat          key.Binding
	chatSend      key.Binding
	chatBack      key.Binding
	addTask       key.Binding
	confirmAdd    key.Binding
	cancelAdd     key.Binding

	completeTask     key.Binding
	priorityTask     key.Binding
	themeCycle       key.Binding
	journalKey       key.Binding
	notes            key.Binding
	calendarKey      key.Binding
	navCalPrev       key.Binding
	navCalNext       key.Binding
	settingsKey      key.Binding
	filterKey        key.Binding
	palette          key.Binding
	save             key.Binding
	audioToggleKey   key.Binding
	workDurationKey  key.Binding
	breakDurationKey key.Binding
	aiPanelKey       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.start, k.stop, k.reset, k.focusMode, k.chat, k.calendarKey, k.settingsKey, k.filterKey, k.addTask, k.completeTask, k.priorityTask, k.journalKey, k.notes, k.themeCycle, k.palette, k.help, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.start, k.stop, k.reset},
		{k.up, k.down, k.enter},
		{k.focusMode, k.tab, k.chat, k.calendarKey, k.settingsKey, k.filterKey, k.navCalPrev, k.navCalNext, k.addTask, k.completeTask, k.priorityTask, k.notes, k.themeCycle, k.palette, k.save, k.help, k.quit},
	}
}

func newKeyMap() keyMap {
	return keyMap{
		start: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "start timer"),
		),
		stop: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "stop timer"),
		),
		reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset timer"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q", "ctrl+q"),
			key.WithHelp("q/Ctrl+Q", "quit"),
		),
		help: key.NewBinding(
			key.WithKeys("?", "ctrl+h"),
			key.WithHelp("?/Ctrl+H", "help"),
		),
		switchSession: key.NewBinding(
			key.WithKeys("b", "space"),
			key.WithHelp("b/space", "switch mode"),
		),
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/focus"),
		),
		tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch view"),
		),
		focusMode: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "focus mode"),
		),
		chat: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "chat assistant"),
		),
		chatSend: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "send message"),
		),
		chatBack: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back to dashboard"),
		),
		addTask: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add task"),
		),
		confirmAdd: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm add"),
		),
		cancelAdd: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel add"),
		),

		completeTask: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "complete task"),
		),
		priorityTask: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "change priority"),
		),
		themeCycle: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("Ctrl+T", "cycle theme"),
		),
		journalKey: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "journal"),
		),
		notes: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "notes"),
		),
		calendarKey: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "calendar"),
		),
		navCalPrev: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "prev month"),
		),
		navCalNext: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next month"),
		),
		settingsKey: key.NewBinding(
			key.WithKeys(","),
			key.WithHelp(",", "settings"),
		),
		audioToggleKey: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "toggle audio"),
		),
		workDurationKey: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "work duration"),
		),
		breakDurationKey: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "break duration"),
		),
		filterKey: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		palette: key.NewBinding(
			key.WithKeys("ctrl+/"),
			key.WithHelp("Ctrl+/", "command palette"),
		),
		save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("Ctrl+S", "save"),
		),
		aiPanelKey: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("Ctrl+A", "ai dashboard"),
		),
	}
}

var (
	// DYNAMIC COLORS - Will be updated based on theme
	synthPink   lipgloss.Color
	synthBlue   lipgloss.Color
	synthPurple lipgloss.Color
	synthYellow lipgloss.Color
	synthCyan   lipgloss.Color
	synthGreen  lipgloss.Color
	synthRed    lipgloss.Color
	gridLine    lipgloss.Color

	// DYNAMIC STYLES - Will be recreated on theme change
	titleStyle         lipgloss.Style
	sectionTitleStyle  lipgloss.Style
	taskStyle          lipgloss.Style
	selectedTaskStyle  lipgloss.Style
	completedTaskStyle lipgloss.Style
	statsStyle         lipgloss.Style
	helpStyle          lipgloss.Style
	glitchStyle        lipgloss.Style
	statusBarStyle     lipgloss.Style
	chatInputStyle     lipgloss.Style
	chatMessageStyle   lipgloss.Style
	chatUserStyle      lipgloss.Style
	taskInputStyle     lipgloss.Style
	priorityInputStyle lipgloss.Style
)

func NewMainModel(tasks []DashboardTask) *MainModel {
	t := timer.NewWithInterval(time.Minute*25, time.Second)
	s := stopwatch.New()
	vp := viewport.New(40, 10)

	m := &MainModel{
		tasks:         tasks,
		timer:         t,
		stopwatch:     s,
		help:          help.New(),
		keys:          newKeyMap(),
		workTime:      time.Minute * 25,
		breakTime:     time.Minute * 5,
		sessionType:   workSession,
		sessions:      0,
		currentView:   dashboardView,
		selectedTask:  0,
		chatHistory:   []string{"🤖 SynthWave AI Assistant Ready!"},
		chatViewport:  vp,
		glitchMessage: " 💾 GRID ERROR DETECTED 💾 ",

		// Theme management

		// Settings management
		settingsMode:  false,
		audioEnabled:  true,
		workDuration:  time.Minute * 25,
		breakDuration: time.Minute * 5,

		// Enhanced filtering
		filterMode:     false,
		filterStatus:   "all",
		filterPriority: "all",

		// Loading and spinner states
		loadingState: startingUp,
		spinnerFrame: 0,

		// AI manager initialization
		aiManager:      ai.New(),
		aiStatus:       "checking",
		lastAICheck:    time.Now(),
		aiThinking:     false,
		aiSpinnerFrame: 0,

		// AI Dashboard panel
		aiPanel: aipanel.New(nil, 50, 20),

		// Calendar management
		calViewMode:     "month",
		cal:             calendar.New(theme.GetManager().Current().Name),
		calSelectedDate: time.Now(),
	}

	// Initialize theme colors and styles
	m.updateTheme()

	// AI status will be checked via command in Init()
	m.aiStatus = "checking"

	// Initialize calendar renderer after theme is set
	m.calRenderer = calendar.NewRenderer(theme.GetManager().Current().Name, 80, 20)

	// Initialize glow styler for markdown notes
	m.glowStyler = glow.NewGlowStyler("synthwave")

	// Wire Brain to AI panel for streaming generation
	bp := ai.NewBrainProvider()
	if bp != nil && bp.Brain() != nil {
		m.aiPanel = aipanel.New(bp.Brain(), 50, 20)
	}

	return m
}

func (m *MainModel) Init() tea.Cmd {
	// Start real-time clock ticker, spinner, and AI status check
	return tea.Batch(
		tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }),
		tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg { return spinnerTickMsg(t) }),
		m.checkAIStatusCmd(),
	)
}

// checkAIStatusCmd returns a command that checks AI status and returns a message
func (m *MainModel) checkAIStatusCmd() tea.Cmd {
	if m.aiManager == nil {
		return nil
	}
	return func() tea.Msg {
		status := "offline"
		if !m.aiManager.IsOllamaAvailable() {
			if err := m.aiManager.LaunchOllama(); err == nil {
				status = "online"
			}
		} else {
			status = "online"
		}
		return aiStatusMsg{status: status}
	}
}

// Spinner frames
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func (m *MainModel) getSpinner() string {
	if m.loadingState == notLoading {
		return ""
	}
	return spinnerFrames[m.spinnerFrame%len(spinnerFrames)]
}

// updateTheme updates all colors and styles based on current theme
func (m *MainModel) updateTheme() {
	// Sync dynamic colors with global theme manager
	t := styles.GetTheme()
	synthPink = t.Accent
	synthBlue = t.Primary
	synthPurple = t.Secondary
	synthYellow = t.Warning
	synthCyan = t.Accent
	synthGreen = t.Success
	synthRed = t.Error
	gridLine = t.Border
	// Recreate all styles with new colors
	m.recreateStyles()
}

// recreateStyles recreates all styles with current theme colors
func (m *MainModel) recreateStyles() {
	titleStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Background(styles.GetBackground()).
		Bold(true).
		Italic(true).
		Padding(1, 3).                         // CONSISTENT: 1 vertical, 3 horizontal
		Margin(1, 1, 1, 1).                    // CONSISTENT: 1 all around
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT: RoundedBorder
		BorderForeground(styles.GetBorder()).
		Align(lipgloss.Center).
		Underline(true).
		Faint(false)

	sectionTitleStyle = lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Bold(true).
		Italic(true).
		Padding(1, 2).      // CONSISTENT: 1 vertical, 2 horizontal
		Margin(1, 0, 1, 0). // CONSISTENT: 1 top/bottom, 0 sides
		Background(styles.GetPanel()).
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT: RoundedBorder
		BorderForeground(styles.GetBorder())

	taskStyle = lipgloss.NewStyle().
		Foreground(synthGreen).
		Padding(0, 2).
		MarginBottom(1).
		Background(styles.GetBoxStyle().GetBackground()).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(gridLine).
		BorderLeft(true).
		BorderRight(true)

	selectedTaskStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Background(styles.GetPanel()).
		Bold(true).
		Italic(true).
		Padding(1, 2).                         // CONSISTENT: 1 vertical, 2 horizontal
		Margin(0, 0, 1, 0).                    // CONSISTENT: 1 bottom only
		BorderStyle(lipgloss.RoundedBorder()). // CHANGED: RoundedBorder
		BorderForeground(styles.GetBorder()).
		Underline(true)

	completedTaskStyle = lipgloss.NewStyle().
		Foreground(synthPurple).
		Strikethrough(true).
		Italic(true).
		Faint(true)

	statsStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Background(styles.GetPanel()).
		Bold(true).
		Padding(1, 2).                         // CONSISTENT
		Margin(1, 1, 1, 1).                    // CONSISTENT
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT
		BorderForeground(styles.GetBorder())

	helpStyle = lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Italic(true).
		Bold(true).
		Background(styles.GetPanel()).
		Padding(1, 2).                         // CONSISTENT
		Margin(1, 0, 1, 0).                    // CONSISTENT
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT
		BorderForeground(styles.GetBorder())

	glitchStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Background(styles.GetPanel()).
		Bold(true).
		Blink(true).
		Italic(true).
		Padding(0, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.GetBorder())

	statusBarStyle = lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Background(styles.GetPanel()).
		Bold(true).
		Padding(1, 2).                         // CONSISTENT
		Margin(1, 0, 1, 0).                    // CONSISTENT
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT
		BorderForeground(styles.GetBorder()).
		BorderBottom(true)

	chatInputStyle = lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Background(styles.GetPanel()).
		Bold(true).
		Padding(1, 2).                         // CONSISTENT
		Margin(1, 0, 1, 0).                    // CONSISTENT
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT
		BorderForeground(styles.GetBorder())

	chatMessageStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Italic(true).
		Background(styles.GetPanel()).
		Padding(1, 2).                         // CONSISTENT
		Margin(0, 0, 1, 0).                    // CONSISTENT
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT
		BorderForeground(styles.GetAccent()).
		BorderLeft(true).
		BorderRight(true)

	chatUserStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Bold(true).
		Italic(true).
		Background(styles.GetPanel()).
		Padding(1, 2).                         // CONSISTENT
		Margin(0, 0, 1, 0).                    // CONSISTENT
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT
		BorderForeground(styles.GetAccent())

	taskInputStyle = lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Background(styles.GetPanel()).
		Bold(true).
		Padding(1, 2).                         // CONSISTENT
		Margin(1, 0, 1, 0).                    // CONSISTENT
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT
		BorderForeground(styles.GetBorder())

	priorityInputStyle = lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Background(styles.GetPanel()).
		Bold(true).
		Italic(true).
		Padding(1, 2).                         // CONSISTENT
		Margin(1, 0, 1, 0).                    // CONSISTENT
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT
		BorderForeground(styles.GetAccent())
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle priority change mode
	if m.taskEntryMode && m.taskInput == "PRIORITY_MODE" {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case msg.String() == "1":
				if len(m.tasks) > 0 && m.selectedTask < len(m.tasks) {
					m.tasks[m.selectedTask].Priority = "low"
					m.taskInput = ""
					m.taskEntryMode = false
				}
				return m, nil
			case msg.String() == "2":
				if len(m.tasks) > 0 && m.selectedTask < len(m.tasks) {
					m.tasks[m.selectedTask].Priority = "medium"
					m.taskInput = ""
					m.taskEntryMode = false
				}
				return m, nil
			case msg.String() == "3":
				if len(m.tasks) > 0 && m.selectedTask < len(m.tasks) {
					m.tasks[m.selectedTask].Priority = "high"
					m.taskInput = ""
					m.taskEntryMode = false
				}
				return m, nil
			case key.Matches(msg, m.keys.cancelAdd):
				m.taskInput = ""
				m.taskEntryMode = false
				return m, nil
			}
		}
	}

	// Handle filter mode
	if m.filterMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.cancelAdd): // Use escape
				m.filterMode = false
				return m, nil
			case key.Matches(msg, m.keys.confirmAdd), key.Matches(msg, m.keys.save): // Apply filter
				m.filterMode = false
				audio.PlaySound(audio.SoundSuccess)
				return m, nil
			case msg.String() == "1":
				m.filterStatus = "all"
				return m, nil
			case msg.String() == "2":
				m.filterStatus = "pending"
				return m, nil
			case msg.String() == "3":
				m.filterStatus = "completed"
				return m, nil
			case msg.String() == "h":
				m.filterPriority = "high"
				return m, nil
			case msg.String() == "m":
				m.filterPriority = "medium"
				return m, nil
			case msg.String() == "l":
				m.filterPriority = "low"
				return m, nil
			case msg.String() == "a":
				m.filterPriority = "all"
				return m, nil
			}
		}
	}

	// Handle notes mode
	if m.notesMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.cancelAdd): // Use escape key
				m.notesMode = false
				m.notesInput = ""
				m.editingTask = nil
				return m, nil
			case key.Matches(msg, m.keys.confirmAdd), key.Matches(msg, m.keys.save): // Save notes
				if m.editingTask != nil {
					m.editingTask.Notes = m.notesInput
					audio.PlaySound(audio.SoundSuccess)
				}
				m.notesMode = false
				m.notesInput = ""
				m.editingTask = nil
				return m, nil
			default:
				// Handle text input for notes
				switch msg.Type {
				case tea.KeyRunes:
					m.notesInput += msg.String()
				case tea.KeyBackspace:
					if len(m.notesInput) > 0 {
						m.notesInput = m.notesInput[:len(m.notesInput)-1]
					}
				}
				return m, nil
			}
		}
	}

	// Handle task entry mode
	if m.taskEntryMode && m.taskInput != "PRIORITY_MODE" {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.confirmAdd), key.Matches(msg, m.keys.save):
				if m.taskInput != "" {
					// Add task (simplified - in real implementation this would call the engine)
					audio.PlaySound(audio.SoundSuccess)

					newTask := DashboardTask{
						ID:          fmt.Sprintf("%d", time.Now().Unix()),
						Description: m.taskInput,
						Status:      "pending",
						Priority:    "medium",
						CreatedAt:   time.Now(),
						Deadline:    nil, // Would be set by AI in real implementation
						Categories:  []string{},
						Notes:       "", // Initialize empty notes
					}
					m.tasks = append(m.tasks, newTask)
					m.taskInput = ""
					m.taskEntryMode = false
				}
				return m, nil

			case key.Matches(msg, m.keys.cancelAdd):
				m.taskInput = ""
				m.taskEntryMode = false
				return m, nil

			default:
				// Handle text input
				switch msg.Type {
				case tea.KeyRunes:
					m.taskInput += msg.String()
				case tea.KeyBackspace:
					if len(m.taskInput) > 0 {
						m.taskInput = m.taskInput[:len(m.taskInput)-1]
					}
				}
				return m, nil
			}
		}
	}


	// Handle AI panel input mode
	if m.aiPanelActive && m.aiPanel.Visible() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.aiPanelKey):
				// Ctrl+A also closes the panel from input mode
				m.aiPanel = m.aiPanel.Toggle()
				m.aiPanelActive = false
				m.aiPanelInput = ""
				return m, nil
			case key.Matches(msg, m.keys.cancelAdd): // Escape
				m.aiPanelActive = false
				m.aiPanelInput = ""
				return m, nil
			case key.Matches(msg, m.keys.confirmAdd), key.Matches(msg, m.keys.enter):
				if m.aiPanelInput != "" {
					// Custom query
					prompt := m.buildAIQuery(m.aiPanelInput)
					m.aiPanel = m.aiPanel.StartStream("AI Dashboard")
					m.aiPanelInput = ""
					return m, m.aiPanel.Generate(prompt)
				} else {
					// Daily briefing (Enter with no input)
					prompt := m.buildDailyBriefing()
					m.aiPanel = m.aiPanel.StartStream("Daily Briefing")
					return m, m.aiPanel.Generate(prompt)
				}
			default:
				// Handle text input for AI panel
				switch msg.Type {
				case tea.KeyRunes:
					m.aiPanelInput += msg.String()
				case tea.KeyBackspace:
					if len(m.aiPanelInput) > 0 {
						m.aiPanelInput = m.aiPanelInput[:len(m.aiPanelInput)-1]
					}
				}
				return m, nil
			}
		}
	}
	// Handle chat input mode
	if m.currentView == chatView {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.quit):
				m.quitting = true
				return m, tea.Quit

			case key.Matches(msg, m.keys.chatBack):
				m.currentView = dashboardView
				return m, nil

			case key.Matches(msg, m.keys.chatSend):
				if m.chatInput != "" {
					// Start AI thinking indicator
					m.aiThinking = true
					userInput := m.chatInput
					m.chatInput = "" // Clear input immediately
					m.chatHistory = append(m.chatHistory,
						chatUserStyle.Render("You: ")+userInput,
						"🤖 AI is thinking...")

					// Process chat message in background using tea.Cmd (thread-safe)
					return m, m.processChatCmd(userInput)
				}
				return m, nil

			default:
				// Handle text input
				switch msg.Type {
				case tea.KeyRunes:
					m.chatInput += msg.String()
				case tea.KeyBackspace:
					if len(m.chatInput) > 0 {
						m.chatInput = m.chatInput[:len(m.chatInput)-1]
					}
				}
				return m, nil
			}
		}
	}

	// Handle AI status checking (every 10 seconds)
	if time.Since(m.lastAICheck) > 10*time.Second {
		if m.aiManager != nil {
			if m.aiManager.IsOllamaAvailable() {
				m.aiStatus = "online"
			} else {
				m.aiStatus = "offline"
			}
		} else {
			m.aiStatus = "offline"
		}
		m.lastAICheck = time.Now()
	}

	// Animate AI spinner when thinking
	if m.aiThinking {
		m.aiSpinnerFrame = (m.aiSpinnerFrame + 1) % 8
	}

	// Auto-launch Ollama if offline and haven't tried recently
	if m.aiStatus == "offline" && time.Since(m.lastAICheck) > 11*time.Second {
		m.lastAICheck = time.Now()
		return m, m.checkAIStatusCmd()
	}

	// Handle normal dashboard mode
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global keys that should work in ANY mode (except task entry)
		switch {
		case key.Matches(msg, m.keys.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.help):
			m.showHelp = !m.showHelp
			return m, nil
		case key.Matches(msg, m.keys.palette):
			m.chatHistory = append(m.chatHistory, "⌨ Command palette coming soon")
			audio.PlaySound(audio.SoundNavigate)
			return m, nil

		case key.Matches(msg, m.keys.aiPanelKey):
			m.aiPanel = m.aiPanel.Toggle()
			m.aiPanelActive = m.aiPanel.Visible()
			audio.PlaySound(audio.SoundNavigate)
			return m, nil

		case key.Matches(msg, m.keys.calendarKey):
			if !m.taskEntryMode {
				m.currentView = calendarView
				// Load tasks into calendar when switching to calendar view
				if m.cal != nil {
					m.loadTasksIntoCalendar()
				}
				audio.PlaySound(audio.SoundNavigate)
			}
			return m, nil

		case key.Matches(msg, m.keys.notes):
			// ENHANCED NOTES MODE - Always create new note task
			if !m.taskEntryMode {
				// Always create a new note task regardless of current view
				newNote := DashboardTask{
					ID:          fmt.Sprintf("note_%d", len(m.tasks)+1),
					Description: "📝 New Note - Edit to add content",
					Priority:    "medium",
					Status:      "pending",
					CreatedAt:   time.Now(),
					Deadline:    nil,
					Categories:  []string{"notes"},
					Notes:       "",
				}
				m.tasks = append(m.tasks, newNote)
				m.editingTask = &m.tasks[len(m.tasks)-1]
				m.notesMode = true
				m.notesInput = m.editingTask.Notes
				m.currentView = notesView // Switch to notes view for editing
				audio.PlaySound(audio.SoundNavigate)
			}
			return m, nil

		case key.Matches(msg, m.keys.tab):
			next := map[mainView]mainView{dashboardView: focusView, focusView: calendarView, calendarView: notesView, notesView: chatView, chatView: settingsView, settingsView: dashboardView}
			m.currentView = next[m.currentView]
			audio.PlaySound(audio.SoundNavigate)
			return m, nil

		// case key.Matches(msg, m.keys.calendarKey):
		// 	if !m.taskEntryMode {
		// 		// Toggle calendar view
		// 		if m.currentView == calendarView {
		// 			m.currentView = dashboardView
		// 		} else {
		// 			m.currentView = calendarView
		// 			// Load tasks into calendar
		// 			m.loadTasksIntoCalendar()
		// 		}
		// 	}
		// 	return m, nil

		// case key.Matches(msg, m.keys.navCalendarUp):
		// 	if m.currentView == calendarView && !m.taskEntryMode {
		// 		m.calendar.NavigateDate("prev")
		// 	}
		// 	return m, nil

		// case key.Matches(msg, m.keys.navCalendarDown):
		// 	if m.currentView == calendarView && !m.taskEntryMode {
		// 		m.calendar.NavigateDate("next")
		// 	}
		// 	return m, nil

		case key.Matches(msg, m.keys.up):
			if len(m.tasks) > 0 && !m.taskEntryMode {
				audio.PlaySound(audio.SoundNavigate)
				m.selectedTask--
				if m.selectedTask < 0 {
					m.selectedTask = len(m.tasks) - 1
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.down):
			if len(m.tasks) > 0 && !m.taskEntryMode {
				audio.PlaySound(audio.SoundNavigate)
				m.selectedTask++
				if m.selectedTask >= len(m.tasks) {
					m.selectedTask = 0
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.enter):
			if len(m.tasks) > 0 && !m.taskEntryMode {
				m.currentView = focusView
				task := m.tasks[m.selectedTask]
				m.activeTimer = &TimerSession{
					Task: task,
					Mode: "work",
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.tab):
			next := map[mainView]mainView{dashboardView: focusView, focusView: calendarView, calendarView: notesView, notesView: journalView, journalView: chatView, chatView: settingsView, settingsView: dashboardView}
			m.currentView = next[m.currentView]
			audio.PlaySound(audio.SoundNavigate)
			return m, nil

		case key.Matches(msg, m.keys.journalKey):
			// Switch to journal view
			m.currentView = journalView
			audio.PlaySound(audio.SoundNavigate)
			return m, nil

		case key.Matches(msg, m.keys.focusMode):
			if len(m.tasks) > 0 && !m.taskEntryMode {
				m.currentView = focusView
				task := m.tasks[m.selectedTask]
				m.activeTimer = &TimerSession{
					Task: task,
					Mode: "work",
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.chat):
			if !m.taskEntryMode {
				m.currentView = chatView
			}
			return m, nil

		case key.Matches(msg, m.keys.navCalPrev):
			if m.currentView == calendarView && !m.taskEntryMode {
				m.calSelectedDate = m.calSelectedDate.AddDate(0, -1, 0)
				audio.PlaySound(audio.SoundNavigate)
			}
			return m, nil

		case key.Matches(msg, m.keys.navCalNext):
			if m.currentView == calendarView && !m.taskEntryMode {
				m.calSelectedDate = m.calSelectedDate.AddDate(0, 1, 0)
				audio.PlaySound(audio.SoundNavigate)
			}
			return m, nil

		case key.Matches(msg, m.keys.settingsKey):
			if !m.taskEntryMode {
				m.currentView = settingsView
				audio.PlaySound(audio.SoundNavigate)
			}
			return m, nil

		case key.Matches(msg, m.keys.filterKey):
			if !m.taskEntryMode {
				m.filterMode = !m.filterMode
				audio.PlaySound(audio.SoundNavigate)
			}
			return m, nil

		case key.Matches(msg, m.keys.addTask):
			if !m.taskEntryMode {
				audio.PlaySound(audio.SoundTaskAdd)
				m.taskEntryMode = true
				m.taskInput = ""
			}
			return m, nil

		case key.Matches(msg, m.keys.completeTask):
			if len(m.tasks) > 0 && !m.taskEntryMode {
				audio.PlaySound(audio.SoundTaskComplete)
				if m.tasks[m.selectedTask].Status == "completed" {
					m.tasks[m.selectedTask].Status = "pending"
				} else {
					m.tasks[m.selectedTask].Status = "completed"
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.priorityTask):
			if len(m.tasks) > 0 && !m.taskEntryMode {
				m.taskEntryMode = true
				m.taskInput = "PRIORITY_MODE"
			}
			return m, nil

		case key.Matches(msg, m.keys.themeCycle):
			styles.CycleTheme()
			currentTheme := styles.GetTheme()
			m.updateTheme()
			themeMsg := fmt.Sprintf("🎨 Theme changed to %s", currentTheme.Name)
			m.chatHistory = append(m.chatHistory, themeMsg)
			return m, nil

		// Settings view specific keys
		case key.Matches(msg, m.keys.audioToggleKey):
			if m.currentView == settingsView {
				if audio.IsAudioEnabled() {
					audio.DisableAudio()
					status := "muted"
					msg := fmt.Sprintf("🔊 Audio %s", status)
					m.chatHistory = append(m.chatHistory, msg)
				} else {
					audio.EnableAudio()
					status := "enabled"
					msg := fmt.Sprintf("🔊 Audio %s", status)
					m.chatHistory = append(m.chatHistory, msg)
				}
				audio.PlaySound(audio.SoundNavigate)
			}
			return m, nil

		case key.Matches(msg, m.keys.workDurationKey):
			if m.currentView == settingsView {
				// Cycle work duration: 25 -> 30 -> 45 -> 60 -> 25 minutes
				switch m.workDuration {
				case 25:
					m.workDuration = 30
				case 30:
					m.workDuration = 45
				case 45:
					m.workDuration = 60
				default:
					m.workDuration = 25
				}
				msg := fmt.Sprintf("⏱️ Work duration: %d minutes", m.workDuration)
				m.chatHistory = append(m.chatHistory, msg)
				audio.PlaySound(audio.SoundNavigate)
			}
			return m, nil

		case key.Matches(msg, m.keys.breakDurationKey):
			if m.currentView == settingsView {
				// Cycle break duration: 5 -> 10 -> 15 -> 20 -> 5 minutes
				switch m.breakDuration {
				case 5:
					m.breakDuration = 10
				case 10:
					m.breakDuration = 15
				case 15:
					m.breakDuration = 20
				default:
					m.breakDuration = 5
				}
				msg := fmt.Sprintf("☕ Break duration: %d minutes", m.breakDuration)
				m.chatHistory = append(m.chatHistory, msg)
				audio.PlaySound(audio.SoundNavigate)
			}
			return m, nil

		case key.Matches(msg, m.keys.start):
			if m.currentView == focusView && !m.taskEntryMode {
				audio.PlaySound(audio.SoundTimerStart)
				if m.sessionType == workSession {
					return m, m.timer.Init()
				}
				return m, m.stopwatch.Init()
			}
			return m, nil

		case key.Matches(msg, m.keys.stop):
			if m.currentView == focusView && !m.taskEntryMode {
				audio.PlaySound(audio.SoundTimerStop)
				// Work-session timer stops automatically when TimeoutMsg is not sent
			}
			return m, nil

		case key.Matches(msg, m.keys.reset):
			if m.currentView == focusView && !m.taskEntryMode {
				if m.sessionType == workSession {
					m.timer = timer.NewWithInterval(m.workTime, time.Second)
				} else {
					m.stopwatch = stopwatch.New()
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.switchSession):
			if m.currentView == focusView && !m.taskEntryMode {
				if m.sessionType == workSession {
					m.sessionType = breakSession
					m.stopwatch = stopwatch.New()
				} else {
					m.sessionType = workSession
					m.timer = timer.NewWithInterval(m.workTime, time.Second)
				}
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.chatViewport.Width > 0 {
			m.chatViewport.Width = msg.Width - 10
			m.chatViewport.Height = msg.Height - 15
		}
		return m, nil

	case timer.TickMsg:
		if m.currentView == focusView {
			var cmd tea.Cmd
			m.timer, cmd = m.timer.Update(msg)
			return m, cmd
		}

	case timer.StartStopMsg:
		if m.currentView == focusView {
			var cmd tea.Cmd
			m.timer, cmd = m.timer.Update(msg)
			return m, cmd
		}

	case timer.TimeoutMsg:
		if m.currentView == focusView {
			m.sessions++
			m.sessionType = breakSession
			m.stopwatch = stopwatch.New()
			return m, m.stopwatch.Init()
		}

	case stopwatch.TickMsg:
		if m.currentView == focusView {
			var cmd tea.Cmd
			m.stopwatch, cmd = m.stopwatch.Update(msg)
			return m, cmd
		}

	case stopwatch.StartStopMsg:
		if m.currentView == focusView {
			var cmd tea.Cmd
			m.stopwatch, cmd = m.stopwatch.Update(msg)
			return m, cmd
		}

	case tickMsg:
		// Handle real-time updates
		// Update current time for display
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })

	case spinnerTickMsg:
		// Handle spinner animation
		m.spinnerFrame++
		if m.loadingState == startingUp {
			// Stop initial loading after a few frames
			if m.spinnerFrame > 10 {
				m.loadingState = notLoading
			}
		}
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg { return spinnerTickMsg(t) })

	case aiResponseMsg:
		// Handle AI chat response (thread-safe)
		if len(m.chatHistory) >= 2 {
			// Replace the "thinking..." message with actual response
			m.chatHistory[len(m.chatHistory)-1] = chatMessageStyle.Render("AI: " + msg.response)
		}
		m.aiThinking = false
		return m, nil

	case aiStatusMsg:
		// Handle AI status update (thread-safe via message passing)
		m.aiStatus = msg.status
		return m, nil

	case aipanel.StreamChunk:
		// Route stream chunks to the AI panel sub-model
		var cmd tea.Cmd
		m.aiPanel, cmd = m.aiPanel.Update(msg)
		return m, cmd

	case aipanel.ErrorMsg:
		// Route errors to the AI panel sub-model
		var cmd tea.Cmd
		m.aiPanel, cmd = m.aiPanel.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *MainModel) processChatMessage(message string) string {
	// Use real AI manager for chat responses
	if m.aiManager == nil {
		return "🤖 AI manager not initialized"
	}

	// Extract task descriptions for context
	taskDescriptions := make([]string, len(m.tasks))
	for i, task := range m.tasks {
		status := "pending"
		if task.Status == "completed" {
			status = "completed"
		}
		taskDescriptions[i] = fmt.Sprintf("%s (%s, %s)", task.Description, status, task.Priority)
	}

	// Load cross-app context (best-effort, skip if unavailable)
	if crossApp := m.aiManager.GetCrossAppContext(context.Background(), 3); len(crossApp) > 0 {
		var summaries []string
		for _, c := range crossApp {
			summaries = append(summaries, c.Summary)
		}
		message = "Context from other apps: " + strings.Join(summaries, "; ") + "\n\n" + message
	}

	// Call real AI
	response, err := m.aiManager.ChatAssistant(context.Background(), message, taskDescriptions)
	if err != nil {
		// Fallback response if AI fails
		return fmt.Sprintf("🤖 AI unavailable (%s). However, you have %d tasks. Try starting Ollama with: ollama serve", err.Error(), len(m.tasks))
	}

	return response
}

// processChatCmd runs the AI chat processing in a goroutine and returns the result as a message
func (m *MainModel) processChatCmd(userInput string) tea.Cmd {
	return func() tea.Msg {
		response := m.processChatMessage(userInput)
		return aiResponseMsg{
			response:  response,
			userInput: userInput,
		}
	}
}

func (m *MainModel) formatChatHistory() string {
	var b strings.Builder
	for _, msg := range m.chatHistory {
		b.WriteString(msg)
		b.WriteString("\n")
	}
	return b.String()
}

func (m *MainModel) loadTasksIntoCalendar() {
	if m.cal == nil {
		return
	}

	// Convert DashboardTask to models.Task for calendar
	modelTasks := make([]models.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		modelTask := models.Task{
			ID:          task.ID,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			CreatedAt:   task.CreatedAt,
			Deadline:    time.Time{}, // DashboardTask doesn't have deadline yet
			Categories:  task.Categories,
		}
		if task.Deadline != nil {
			modelTask.Deadline = *task.Deadline
		}
		modelTasks = append(modelTasks, modelTask)
	}

	m.cal.LoadTasks(modelTasks)
	m.cal.SelectedDate = m.calSelectedDate
}

func (m *MainModel) getFilteredTasks() []DashboardTask {
	if !m.filterMode {
		return m.tasks
	}

	var filtered []DashboardTask
	for _, task := range m.tasks {
		// Filter by status
		if m.filterStatus != "all" {
			if m.filterStatus == "pending" && task.Status != "pending" {
				continue
			}
			if m.filterStatus == "completed" && task.Status != "completed" {
				continue
			}
		}

		// Filter by priority
		if m.filterPriority != "all" {
			if task.Priority != m.filterPriority {
				continue
			}
		}

		filtered = append(filtered, task)
	}

	return filtered
}

func (m *MainModel) renderProgressBar() string {
	var totalTime time.Duration
	var elapsed time.Duration

	if m.sessionType == workSession {
		totalTime = m.workTime
		elapsed = m.workTime - m.timer.Timeout
	} else {
		totalTime = m.breakTime
		if m.stopwatch.Running() {
			elapsed = m.stopwatch.Elapsed()
		} else {
			elapsed = 0
		}
	}

	progress := 0.0
	if totalTime > 0 {
		progress = float64(elapsed) / float64(totalTime)
		if progress > 1.0 {
			progress = 1.0
		}
	}

	barWidth := 30
	filledWidth := int(float64(barWidth) * progress)
	emptyWidth := barWidth - filledWidth

	var bar strings.Builder
	bar.WriteString("[")

	// Filled portion (synth green)
	for i := 0; i < filledWidth; i++ {
		bar.WriteString(lipgloss.NewStyle().Foreground(synthGreen).Render("="))
	}

	// Empty portion (dark)
	for i := 0; i < emptyWidth; i++ {
		bar.WriteString(lipgloss.NewStyle().Foreground(gridLine).Render("-"))
	}

	bar.WriteString("] ")
	percentage := int(progress * 100)
	fmt.Fprintf(&bar, "%d%%", percentage)

	return bar.String()
}

func (m *MainModel) renderStats() string {
	total := len(m.tasks)
	completed := 0
	pending := 0

	for _, task := range m.tasks {
		if task.Status == "completed" {
			completed++
		} else {
			pending++
		}
	}

	stats := []string{
		fmt.Sprintf("🔢 Total: %d", total),
		fmt.Sprintf("✅ Done: %d", completed),
		fmt.Sprintf("⏳ Pending: %d", pending),
	}

	// Enhanced stats with AI status indicator and more flair
	var aiStatusIcon string
	switch m.aiStatus {
	case "online":
		aiStatusIcon = "🤖 AI: 🟢 Online"
	case "offline":
		aiStatusIcon = "🤖 AI: 🔴 Offline"
	default:
		aiStatusIcon = "🤖 AI: 🟡 Checking..."
	}
	stats = append(stats, aiStatusIcon)

	if total > 0 {
		completion := float64(completed) / float64(total) * 100
		progressEmoji := "🎯"
		if completion >= 80 {
			progressEmoji = "🔥"
		} else if completion >= 50 {
			progressEmoji = "⚡"
		}
		stats = append(stats, fmt.Sprintf("%s Rate: %.0f%%", progressEmoji, completion))
	}

	// Add animated time indicator
	currentTime := fmt.Sprintf("🕐 %s", time.Now().Format("15:04:05"))
	stats = append(stats, currentTime)

	// Render enhanced stats with proper alignment and theme colors
	var statsContent strings.Builder
	for i, stat := range stats {
		statStyle := lipgloss.NewStyle().
			Bold(true).
			Background(styles.GetPanel()).
			Padding(0, 1).
			MarginLeft(1).
			BorderStyle(lipgloss.RoundedBorder())

		// Alternate colors for visual distinction
		if i%2 == 0 {
			// Even indexes (0, 2, 4): tan color
			statStyle = statStyle.
				Foreground(styles.GetAccent()).
				BorderForeground(styles.GetBorder())
		} else {
			// Odd indexes (1, 3, 5): yellow color
			statStyle = statStyle.
				Foreground(styles.GetAccent()).
				BorderForeground(styles.GetAccent())
		}

		// Ensure consistent width for alignment
		maxWidth := 25 // Reasonable width for stat boxes
		statText := lipgloss.NewStyle().Width(maxWidth).Render("  " + stat + "  ")
		statsContent.WriteString(statStyle.Render(statText))

		if i < len(stats)-1 {
			statsContent.WriteString(" ")
		}

		// Break lines properly: 3 per line
		if (i+1)%3 == 0 {
			statsContent.WriteString("\n")
		}
	}

	// Enhanced stats container with more flair
	statsContainer := lipgloss.NewStyle().
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(1, 2).
		Margin(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthPink).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)

	return statsContainer.Render(statsContent.String())
}

func (m *MainModel) renderTaskList() string {
	var b strings.Builder

	// Use filtered tasks
	tasksToShow := m.getFilteredTasks()
	taskCount := len(tasksToShow)

	if taskCount == 0 {
		if m.filterMode {
			emptyStyle := lipgloss.NewStyle().
				Foreground(synthYellow).
				Italic(true).
				Bold(true).
				Background(styles.GetBoxStyle().GetBackground()).
				Padding(1, 2).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(synthYellow).
				Render("🔍 No tasks match current filter")
			return emptyStyle
		} else {
			emptyStyle := lipgloss.NewStyle().
				Foreground(synthBlue).
				Italic(true).
				Bold(true).
				Background(styles.GetBoxStyle().GetBackground()).
				Padding(1, 2).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(synthBlue).
				Render("😴 No tasks in the grid - ready for new items!")
			return emptyStyle
		}
	}

	// Show enhanced filter info if in filter mode
	if m.filterMode {
		filterInfo := fmt.Sprintf("🔍 Filter: %s status, %s priority", m.filterStatus, m.filterPriority)
		filterStyle := lipgloss.NewStyle().
			Foreground(synthYellow).
			Bold(true).
			Italic(true).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthYellow).
			Render(filterInfo)
		b.WriteString(filterStyle)
		b.WriteString("\n")
	}

	for i, task := range tasksToShow {
		var taskLine string
		if task.Status == "completed" {
			// Enhanced completed task style
			completedPrefix := lipgloss.NewStyle().
				Foreground(synthGreen).
				Bold(true).
				Render("✅")
			taskDesc := lipgloss.NewStyle().
				Foreground(synthPurple).
				Strikethrough(true).
				Italic(true).
				Faint(true).
				Render(task.Description)
			taskLine = fmt.Sprintf("%s %s", completedPrefix, taskDesc)
		} else {
			// Enhanced priority symbols with colors
			prioritySymbol := ""
			var priorityColor lipgloss.Color
			switch task.Priority {
			case "high":
				prioritySymbol = "🔥"
				priorityColor = synthRed
			case "medium":
				prioritySymbol = "⚡"
				priorityColor = synthYellow
			case "low":
				prioritySymbol = "💤"
				priorityColor = synthPurple
			default:
				prioritySymbol = "⚪"
				priorityColor = synthCyan
			}

			priorityStyled := lipgloss.NewStyle().
				Foreground(priorityColor).
				Bold(true).
				Render(prioritySymbol)

			taskDesc := lipgloss.NewStyle().
				Foreground(synthGreen).
				Render(task.Description)

			taskLine = fmt.Sprintf("%s %s", priorityStyled, taskDesc)
		}

		// Apply selection highlighting
		if i == m.selectedTask {
			taskLine = selectedTaskStyle.Render(taskLine)
		} else {
			taskLine = taskStyle.Render(taskLine)
		}

		b.WriteString(taskLine)
		b.WriteString("\n")
	}

	return b.String()
}

func (m *MainModel) renderFocusView() string {
	var b strings.Builder

	if m.activeTimer == nil {
		return "No active timer"
	}

	// LARGE HEADER
	header := lipgloss.NewStyle().
		Foreground(synthPink).
		Bold(true).
		Align(lipgloss.Center).
		Render("CYBER FOCUS MODE")
	b.WriteString(header)
	b.WriteString("\n\n")

	// CURRENT MISSION
	b.WriteString(sectionTitleStyle.Render("🎯 CURRENT MISSION:"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(synthGreen).
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthBlue).
		Render(m.activeTimer.Task.Description))
	b.WriteString("\n\n")

	// SESSION INFO
	b.WriteString(sectionTitleStyle.Render("📊 SESSION DATA:"))
	b.WriteString("\n")
	sessionInfo := fmt.Sprintf("💾 Sessions Completed: %d", m.sessions)
	b.WriteString(lipgloss.NewStyle().Foreground(synthYellow).Bold(true).Render(sessionInfo))
	b.WriteString("\n\n")

	// MODE STATUS
	b.WriteString(sectionTitleStyle.Render("💾 CURRENT MODE:"))
	b.WriteString("\n")
	if m.sessionType == workSession {
		b.WriteString(lipgloss.NewStyle().
			Foreground(synthGreen).
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthGreen).
			Padding(0, 1).
			Blink(true).
			Render("🔥 FOCUS MODE ACTIVATED"))
	} else {
		b.WriteString(lipgloss.NewStyle().
			Foreground(synthBlue).
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthBlue).
			Padding(0, 1).
			Blink(true).
			Render("🌌 RELAXATION MODE ENGAGED"))
	}
	b.WriteString("\n\n")

	// TIMER DISPLAY
	b.WriteString(sectionTitleStyle.Render("⏰ TIME:"))
	b.WriteString("\n")
	timerDisplay := ""
	if m.sessionType == workSession {
		timerDisplay = m.timer.View()
	} else {
		timerDisplay = m.stopwatch.View()
	}

	b.WriteString(lipgloss.NewStyle().
		Foreground(synthGreen).
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(2, 4).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthGreen).
		Align(lipgloss.Center).
		Bold(true).
		Render(timerDisplay))
	b.WriteString("\n\n")

	// PROGRESS BAR
	b.WriteString(sectionTitleStyle.Render("📊 PROGRESS:"))
	b.WriteString("\n")
	b.WriteString(m.renderProgressBar())
	b.WriteString("\n\n")

	return b.String()
}

func (m *MainModel) renderDashboard() string {
	var b strings.Builder

	// Simple decorative lines - make responsive to width
	headerLineLength := m.width - 4
	if headerLineLength < 20 {
		headerLineLength = 20
	}
	headerLine := strings.Repeat("═", headerLineLength)
	headerLine = lipgloss.NewStyle().
		Foreground(synthPink).
		Render(headerLine)
	b.WriteString(headerLine)
	b.WriteString("\n")

	titleLine := lipgloss.NewStyle().
		Foreground(synthBlue).
		Bold(true).
		Align(lipgloss.Center).
		Render("focus.sh")
	b.WriteString(titleLine)

	// Add spinner if loading
	if m.loadingState != notLoading {
		spinnerLine := lipgloss.NewStyle().
			Foreground(synthYellow).
			Align(lipgloss.Center).
			Render(m.getSpinner() + " Loading...")
		b.WriteString(spinnerLine)
		b.WriteString("\n")
	}

	b.WriteString("\n")

	bottomLine := lipgloss.NewStyle().
		Foreground(synthPink).
		Render("════════════════════════════════════════════════════════════════════════════════════════════════════════")
	b.WriteString(bottomLine)
	b.WriteString("\n\n")

	// Task entry overlay
	if m.taskEntryMode {
		// Responsive overlay sizing
		overlayWidth := m.width - 20
		if overlayWidth < 60 {
			overlayWidth = 60
		}
		if overlayWidth > 80 {
			overlayWidth = 80
		}

		overlayHeight := m.height - 10
		if overlayHeight < 10 {
			overlayHeight = 10
		}
		if overlayHeight > 15 {
			overlayHeight = 15
		}

		overlay := lipgloss.NewStyle().
			Foreground(synthGreen).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthGreen).
			Align(lipgloss.Center).
			Width(overlayWidth).
			Height(overlayHeight).
			Render(m.renderTaskEntry())

		// Center the overlay
		container := lipgloss.NewStyle().
			Width(m.width-4).
			Height(m.height-4).
			Align(lipgloss.Center, lipgloss.Center).
			Render(overlay)

		b.WriteString(container)
		return b.String()
	}

	// Notes overlay
	if m.notesMode && m.editingTask != nil {
		// Responsive overlay sizing
		overlayWidth := m.width - 20
		if overlayWidth < 60 {
			overlayWidth = 60
		}
		if overlayWidth > 100 {
			overlayWidth = 100
		}

		overlayHeight := m.height - 10
		if overlayHeight < 15 {
			overlayHeight = 15
		}
		if overlayHeight > 25 {
			overlayHeight = 25
		}

		overlay := lipgloss.NewStyle().
			Foreground(synthCyan).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthCyan).
			Align(lipgloss.Center).
			Width(overlayWidth).
			Height(overlayHeight).
			Render(m.renderNotesEditor())

		// Center the overlay
		container := lipgloss.NewStyle().
			Width(m.width-4).
			Height(m.height-4).
			Align(lipgloss.Center, lipgloss.Center).
			Render(overlay)

		b.WriteString(container)
		return b.String()
	}

	// Filter overlay
	if m.filterMode {
		// Responsive overlay sizing
		overlayWidth := m.width - 20
		if overlayWidth < 60 {
			overlayWidth = 60
		}
		if overlayWidth > 80 {
			overlayWidth = 80
		}

		overlayHeight := m.height - 10
		if overlayHeight < 15 {
			overlayHeight = 15
		}
		if overlayHeight > 20 {
			overlayHeight = 20
		}

		overlay := lipgloss.NewStyle().
			Foreground(synthYellow).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthYellow).
			Align(lipgloss.Center).
			Width(overlayWidth).
			Height(overlayHeight).
			Render(m.renderFilterEditor())

		// Center the overlay
		container := lipgloss.NewStyle().
			Width(m.width-4).
			Height(m.height-4).
			Align(lipgloss.Center, lipgloss.Center).
			Render(overlay)

		b.WriteString(container)
		return b.String()
	}

	// Main dashboard layout - using strings.Builder for performance
	var leftContent strings.Builder
	var centerContent strings.Builder
	var rightContent strings.Builder

	// Left column - Task list
	leftContent.WriteString(sectionTitleStyle.Render("📋 MISSION BOARD:"))
	leftContent.WriteString("\n\n")
	leftContent.WriteString(m.renderTaskList())

	// Center column - Selected task
	centerContent.WriteString(sectionTitleStyle.Render("🎯 SELECTED MISSION:"))
	centerContent.WriteString("\n")
	if len(m.tasks) > 0 && m.selectedTask < len(m.tasks) {
		task := m.tasks[m.selectedTask]
		status := "PENDING"
		if task.Status == "completed" {
			status = "COMPLETED ✅"
		}

		centerContent.WriteString(lipgloss.NewStyle().
			Foreground(synthCyan).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthBlue).
			Render(task.Description))
		centerContent.WriteString("\n\n")
		centerContent.WriteString(lipgloss.NewStyle().Foreground(synthGreen).Render(
			fmt.Sprintf("Status: %s\nPriority: %s", status, strings.ToUpper(task.Priority))))
	} else {
		centerContent.WriteString("😴 No missions selected")
	}
	centerContent.WriteString("\n\n")
	centerContent.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render(
		"💡 [ENTER] Focus mode\n💡 [A] Add mission\n💡 [C] Chat assistant\n💡 [Ctrl+A] AI Dashboard\n💡 [D] Complete\n💡 [P] Priority\n💡 [G] Glitch"))

	// Right column - Stats and quick actions
	rightContent.WriteString(m.renderStats())
	rightContent.WriteString("\n\n")
	rightContent.WriteString(sectionTitleStyle.Render("⚡ REAL-TIME DATA:"))
	rightContent.WriteString("\n")
	rightContent.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("🕐 " + time.Now().Format("15:04:05")))
	rightContent.WriteString("\n")
	rightContent.WriteString(lipgloss.NewStyle().Foreground(synthYellow).Render("📅 " + time.Now().Format("2006-01-02")))
	rightContent.WriteString("\n\n")
	rightContent.WriteString(sectionTitleStyle.Render("⚡ QUICK ACTIONS:"))
	rightContent.WriteString("\n")
	actions := []string{
		"[A] Add task",
		"[D] Complete task",
		"[P] Change priority",
		"[C] Chat assistant",
		"[Ctrl+A] AI Dashboard",
		"[F] Focus mode",
		"[↑/↓] Navigate",
	}
	for _, action := range actions {
		rightContent.WriteString(lipgloss.NewStyle().Foreground(styles.GetAccent()).Render("  " + action))
		rightContent.WriteString("\n")
	}

	// Convert builders to strings
	leftStr := leftContent.String()
	centerStr := centerContent.String()
	rightStr := rightContent.String()

	// Apply width constraints safely
	if m.width > 80 { // Only apply if we have enough width
		colWidth := (m.width - 6) / 3 // Account for padding
		if colWidth < 25 {            // Minimum column width for small windows
			colWidth = 25
		}

		leftStyle := lipgloss.NewStyle().Width(colWidth)
		centerStyle := lipgloss.NewStyle().Width(colWidth)
		rightStyle := lipgloss.NewStyle().Width(colWidth)

		leftStr = leftStyle.Render(leftStr)
		centerStr = centerStyle.Render(centerStr)
		rightStr = rightStyle.Render(rightStr)
	} else {
		// For small windows, stack vertically
		leftStyle := lipgloss.NewStyle().Width(m.width - 4)
		centerStyle := lipgloss.NewStyle().Width(m.width - 4)
		rightStyle := lipgloss.NewStyle().Width(m.width - 4)

		leftStr = leftStyle.Render(leftStr)
		centerStr = centerStyle.Render(centerStr)
		rightStr = rightStyle.Render(rightStr)
	}

	// Use appropriate layout based on window size
	var layout string
	if m.width > 80 {
		layout = lipgloss.JoinHorizontal(lipgloss.Top, leftStr, centerStr, rightStr)
	} else {
		// Small windows: stack vertically
		layout = lipgloss.JoinVertical(lipgloss.Left, leftStr, centerStr, rightStr)
	}
	b.WriteString(layout)

	return b.String()
}

func (m *MainModel) View() string {
	if m.quitting {
		return "\n  Powering down the grid... 🌌\n\n"
	}

	// Handle initial size - make more room for complex layouts
	if m.width == 0 {
		m.width = 160
		m.height = 50
	}

	// Ensure minimum size to prevent layout issues
	if m.width < 120 {
		m.width = 120
	}
	if m.height < 40 {
		m.height = 40
	}

	var content string

	// Determine what to display based on current view and modes
	switch {
	case m.taskEntryMode:
		// Task entry overlay
		content = m.renderTaskEntryOverlay()
	case m.notesMode && m.editingTask != nil:
		// Notes overlay
		content = m.renderNotesOverlay()
	case m.filterMode:
		// Filter overlay
		content = m.renderFilterOverlay()
	default:
		// Normal view switching with complete user-centric flow
		switch m.currentView {
		case dashboardView:
			content = m.renderDashboard()
		case focusView:
			if m.activeTimer != nil {
				content = m.renderFocusView()
			} else {
				content = "No active timer"
			}
		case chatView:
			content = m.renderChatView()
		case calendarView:
			content = m.renderCalendarView()
		case notesView:
			content = m.renderNotesEditor() // Notes as standalone view
		case journalView:
			content = m.renderJournalView()
		case settingsView:
			content = m.renderSettingsView()
		default:
			content = m.renderDashboard()
		}
	}

	// Add help bar at the bottom
	helpBar := m.renderHelpBar()

	// Render main content, optionally with AI panel as side panel
	mainContent := lipgloss.JoinVertical(lipgloss.Left,
		content,
		helpBar,
	)

	if m.aiPanel.Visible() {
		// Update panel dimensions
		m.aiPanel = m.aiPanel.SetSize(m.width/3, m.height-4)

		// Render AI panel with input indicator
		panelView := m.aiPanel.View()
		if m.aiPanelActive {
			inputLine := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f4a261")).
				Render("> " + m.aiPanelInput + "█")
			panelView = panelView + "\n" + inputLine
		}

		// Place panel on the right side
		mainStyled := lipgloss.NewStyle().
			Width(m.width * 2 / 3).
			Render(mainContent)
		panelStyled := lipgloss.NewStyle().
			Width(m.width / 3).
			Render(panelView)

		fullContent := lipgloss.NewStyle().
			Height(m.height).
			Width(m.width).
			Render(lipgloss.JoinHorizontal(lipgloss.Top, mainStyled, panelStyled))

		return fullContent
	}

	fullContent := lipgloss.NewStyle().
		Height(m.height).
		Width(m.width).
		Render(mainContent)

	return fullContent
}

// Overlay helper functions
func (m *MainModel) renderTaskEntryOverlay() string {
	var overlay strings.Builder

	// Task entry content
	overlay.WriteString(m.renderTaskEntry())

	// Center container
	container := lipgloss.NewStyle().
		Foreground(synthGreen).
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthGreen).
		Align(lipgloss.Center).
		Width(60).
		Height(10).
		Render(overlay.String())

	// Wrap in outer container
	outerContainer := lipgloss.NewStyle().
		Width(m.width-10).
		Height(m.height-15).
		Align(lipgloss.Center, lipgloss.Center).
		Render(container)

	return outerContainer
}

func (m *MainModel) renderNotesOverlay() string {
	var overlay strings.Builder

	// Notes content
	overlay.WriteString(m.renderNotesEditor())

	// Center container with improved sizing
	containerWidth := 80
	containerHeight := 25

	// Adjust for window size
	if m.width > 0 && m.width-20 > containerWidth {
		containerWidth = m.width - 20
	}
	if m.height > 0 && m.height-10 > containerHeight {
		containerHeight = m.height - 10
	}

	container := lipgloss.NewStyle().
		Foreground(synthCyan).
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthCyan).
		Align(lipgloss.Center).
		Width(containerWidth).
		Height(containerHeight).
		Render(overlay.String())

	// Wrap in outer container
	outerContainer := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(container)

	return outerContainer
}

func (m *MainModel) renderFilterOverlay() string {
	var overlay strings.Builder

	// Filter content
	overlay.WriteString(m.renderFilterEditor())

	// Center container with improved sizing
	containerWidth := 60
	containerHeight := 20

	// Adjust for window size
	if m.width > 0 && m.width-30 > containerWidth {
		containerWidth = m.width - 30
	}
	if m.height > 0 && m.height-15 > containerHeight {
		containerHeight = m.height - 15
	}

	container := lipgloss.NewStyle().
		Foreground(synthYellow).
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthYellow).
		Align(lipgloss.Center).
		Width(containerWidth).
		Height(containerHeight).
		Render(overlay.String())

	// Wrap in outer container
	outerContainer := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(container)

	return outerContainer
}


// buildDailyBriefing constructs a prompt for the daily productivity briefing.
func (m *MainModel) buildDailyBriefing() string {
	now := time.Now()
	today := now.Format("2006-01-02")

	var dueToday, overdue, completedToday int
	var high, medium, low int
	var taskList []string

	for _, t := range m.tasks {
		taskList = append(taskList, t.Description)
		switch t.Status {
		case "completed":
			completedToday++
		case "pending":
			if t.Deadline != nil {
				if t.Deadline.Format("2006-01-02") == today {
					dueToday++
				} else if t.Deadline.Before(now) {
					overdue++
				}
			}
		}
		switch t.Priority {
		case "high":
			high++
		case "medium":
			medium++
		case "low":
			low++
		}
	}

	return fmt.Sprintf(
		"Generate a daily productivity briefing. Tasks: %v. Completed today: %d. Due today: %d. Overdue: %d. Priority distribution: high=%d, medium=%d, low=%d. Be concise, 3-5 lines. Suggest one actionable improvement.",
		taskList, completedToday, dueToday, overdue, high, medium, low,
	)
}

// buildAIQuery constructs a prompt for a natural language query to the AI.
func (m *MainModel) buildAIQuery(query string) string {
	var recentTasks []string
	for i, t := range m.tasks {
		if i >= 5 {
			break
		}
		recentTasks = append(recentTasks, t.Description)
	}
	return fmt.Sprintf(
		"You are a helpful productivity assistant. Answer concisely. Context: %v. Question: %s",
		recentTasks, query,
	)
}
func StartMainDashboard(tasks []DashboardTask) error {
	// Kyanite Suite Dashboard Initialization
	fmt.Println("🌌 Kyanite Suite - focus.sh Dashboard")
	fmt.Println("✨ Loading task management interface...")

	model := NewMainModel(tasks)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
