package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/audio"
	"github.com/kyanite/focus/pkg/calendar"
	"github.com/kyanite/focus/pkg/glow"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/internal/ai"
)

// Tick message for real-time updates
type tickMsg time.Time

// Spinner message for loading states
type spinnerTickMsg time.Time

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
	tasks          []DashboardTask
	activeTimer    *TimerSession
	timer          timer.Model
	stopwatch      stopwatch.Model
	help           help.Model
	keys           keyMap
	quitting       bool
	workTime       time.Duration
	breakTime      time.Duration
	sessionType    sessionType
	sessions       int
	currentView    mainView
	selectedTask   int
	width          int
	height         int
	
	// Calendar management
	cal            *calendar.Calendar
	calRenderer    *calendar.Renderer
	calSelectedDate time.Time
	calViewMode    string // "month", "week", "day"
	
	// Chat assistant components
	chatActive     bool
	chatInput      string
	chatHistory    []string
	chatViewport   viewport.Model
	
	// Task entry components
	taskEntryMode  bool
	taskInput      string
	
	// Notes management
	notesMode      bool
	notesInput     string
	editingTask    *DashboardTask
	
	// Enhanced filtering
	filterMode     bool
	filterStatus   string // "all", "pending", "completed"
	filterPriority string // "all", "high", "medium", "low"
	
	// Loading and spinner states
	loadingState   loadingState
	spinnerFrame   int
	
	// Glow styler for enhanced markdown rendering
	glowStyler     *glow.GlowStyler
	
	// AI manager for real chat integration
	aiManager      *ai.Manager
	aiStatus       string // "online", "offline", "checking"
	lastAICheck     time.Time
	aiThinking      bool   // Whether AI is currently responding
	aiSpinnerFrame  int    // Current spinner frame for AI response
	
	// Digital artifacts/glitches
	glitchCount    int
	showGlitch     bool
	glitchMessage  string
	
	// Theme management
	currentThemeIndex int
	themes           []styles.ThemeMode
	themeNames       []string
	
	// Settings management
	settingsMode     bool
	settingsInput    string
	audioEnabled     bool
	workDuration     time.Duration
	breakDuration    time.Duration
	
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
	Notes       string    // Add notes support
}

type TimerSession struct {
	Task        DashboardTask
	Mode        string // "work" or "break"
	Duration    time.Duration
	TimeLeft    time.Duration
	StartTime   time.Time
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
	settingsView
)

type keyMap struct {
	start        key.Binding
	stop         key.Binding
	reset        key.Binding
	quit         key.Binding
	help         key.Binding
	switchSession key.Binding
	up           key.Binding
	down         key.Binding
	enter        key.Binding
	tab          key.Binding
	focusMode    key.Binding
	chat         key.Binding
	chatSend     key.Binding
	chatBack     key.Binding
	addTask      key.Binding
	confirmAdd   key.Binding
	cancelAdd    key.Binding
	glitchTest   key.Binding
	completeTask key.Binding
	priorityTask key.Binding
	themeCycle   key.Binding
	notes        key.Binding
	calendarKey  key.Binding
	navCalPrev   key.Binding
	navCalNext   key.Binding
	settingsKey  key.Binding
	filterKey    key.Binding
	audioToggleKey   key.Binding
	workDurationKey  key.Binding
	breakDurationKey key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.start, k.stop, k.reset, k.focusMode, k.chat, k.calendarKey, k.settingsKey, k.filterKey, k.addTask, k.completeTask, k.priorityTask, k.notes, k.themeCycle, k.help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.start, k.stop, k.reset},
		{k.up, k.down, k.enter},
		{k.focusMode, k.tab, k.chat, k.calendarKey, k.settingsKey, k.filterKey, k.navCalPrev, k.navCalNext, k.addTask, k.completeTask, k.priorityTask, k.notes, k.themeCycle, k.help, k.quit},
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
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q", "quit"),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
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
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
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
		glitchTest: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "test glitch"),
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
			key.WithHelp("ctrl+t", "cycle theme"),
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
	}
}

var (
	// DYNAMIC COLORS - Will be updated based on theme
	synthPink      lipgloss.Color
	synthBlue      lipgloss.Color
	synthPurple    lipgloss.Color
	synthYellow    lipgloss.Color
	synthCyan      lipgloss.Color
	synthGreen     lipgloss.Color
	synthRed       lipgloss.Color
	darkBg         lipgloss.Color
	darkPanel      lipgloss.Color
	gridLine       lipgloss.Color

	// DYNAMIC STYLES - Will be recreated on theme change
	titleStyle     lipgloss.Style
	sectionTitleStyle lipgloss.Style
	taskStyle      lipgloss.Style
	selectedTaskStyle lipgloss.Style
	completedTaskStyle lipgloss.Style
	statsStyle     lipgloss.Style
	helpStyle      lipgloss.Style
	glitchStyle    lipgloss.Style
	statusBarStyle lipgloss.Style
	chatInputStyle lipgloss.Style
	chatMessageStyle lipgloss.Style
	chatUserStyle  lipgloss.Style
	taskInputStyle lipgloss.Style
	priorityInputStyle lipgloss.Style
)

func NewMainModel(tasks []DashboardTask) MainModel {
	t := timer.NewWithInterval(time.Minute*25, time.Second)
	s := stopwatch.New()
	vp := viewport.New(40, 10)

	m := MainModel{
		tasks:       tasks,
		timer:       t,
		stopwatch:   s,
		help:        help.New(),
		keys:        newKeyMap(),
		workTime:    time.Minute * 25,
		breakTime:   time.Minute * 5,
		sessionType: workSession,
		sessions:    0,
		currentView: dashboardView,
		selectedTask: 0,
		chatHistory: []string{"🤖 SynthWave AI Assistant Ready!"},
		chatViewport: vp,
		glitchMessage: " 💾 GRID ERROR DETECTED 💾 ",
		
		// Theme management
		currentThemeIndex: 0,
		themes: []styles.ThemeMode{
			styles.ThemeSynthwave,
			styles.ThemeLight,
			styles.ThemePlain,
		},
		themeNames: []string{"Synthwave", "Light", "Plain"},
		
		// Settings management
		settingsMode: false,
		audioEnabled: true,
		workDuration: time.Minute * 25,
		breakDuration: time.Minute * 5,
		
		// Enhanced filtering
		filterMode: false,
		filterStatus: "all",
		filterPriority: "all",
		
		// Loading and spinner states
		loadingState: startingUp,
		spinnerFrame: 0,
		
		// AI manager initialization
		aiManager: ai.New(),
		aiStatus: "checking",
		lastAICheck: time.Now(),
		aiThinking: false,
		aiSpinnerFrame: 0,
		
		// Calendar management
		calViewMode: "month",
		cal: calendar.New("synthwave"),
		calSelectedDate: time.Now(),
	}
	
		// Initialize theme colors and styles
	m.updateTheme()
	
	// Auto-launch Ollama at startup
	if m.aiManager != nil {
		go func() {
			if !m.aiManager.IsOllamaAvailable() {
				if err := m.aiManager.LaunchOllama(); err == nil {
					m.aiStatus = "online"
				} else {
					m.aiStatus = "offline"
				}
			} else {
				m.aiStatus = "online"
			}
		}()
	}
	
	// Initialize calendar renderer after theme is set
	m.calRenderer = calendar.NewRenderer("synthwave", 80, 20)
	
	// Initialize glow styler for markdown notes
	m.glowStyler = glow.NewGlowStyler("synthwave")
	
	return m
}

func (m MainModel) Init() tea.Cmd {
	// Start real-time clock ticker and spinner
	return tea.Batch(
		tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }),
		tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg { return spinnerTickMsg(t) }),
	)
}

// Spinner frames
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func (m MainModel) getSpinner() string {
	if m.loadingState == notLoading {
		return ""
	}
	return spinnerFrames[m.spinnerFrame%len(spinnerFrames)]
}

// updateTheme updates all colors and styles based on current theme
func (m *MainModel) updateTheme() {
	// Set the theme globally
	styles.SetTheme(m.themes[m.currentThemeIndex])
	
	// Update colors based on theme
	switch m.themes[m.currentThemeIndex] {
	case styles.ThemeSynthwave:
		synthPink = lipgloss.Color("#FF10F0")
		synthBlue = lipgloss.Color("#00FFF0")
		synthPurple = lipgloss.Color("#BD10E0")
		synthYellow = lipgloss.Color("#FFF01F")
		synthCyan = lipgloss.Color("#00FFFF")
		synthGreen = lipgloss.Color("#39FF14")
		synthRed = lipgloss.Color("#FF0040")
		darkBg = lipgloss.Color("#0A0014")
		darkPanel = lipgloss.Color("#1A0033")  // More purple, less grey
		gridLine = lipgloss.Color("#4A0080")   // Brighter purple border
		
	case styles.ThemeLight:
		// MATRIX BLACK AND GREEN THEME
		synthPink = lipgloss.Color("#00FF00")  // Matrix green
		synthBlue = lipgloss.Color("#00CC00")   // Darker green
		synthPurple = lipgloss.Color("#009900") // Even darker green
		synthYellow = lipgloss.Color("#00FF00") // Matrix green highlight
		synthCyan = lipgloss.Color("#00FFAA")  // Green cyan
		synthGreen = lipgloss.Color("#00FF00") // Bright green
		synthRed = lipgloss.Color("#FF0000")   // Matrix red (alerts)
		darkBg = lipgloss.Color("#000000")      // Pure black background
		darkPanel = lipgloss.Color("#0A1A0A")  // Very dark green tint
		gridLine = lipgloss.Color("#003300")    // Dark green grid lines
		
	case styles.ThemePlain:
		synthPink = lipgloss.Color("")
		synthBlue = lipgloss.Color("")
		synthPurple = lipgloss.Color("")
		synthYellow = lipgloss.Color("")
		synthCyan = lipgloss.Color("")
		synthGreen = lipgloss.Color("")
		synthRed = lipgloss.Color("")
		darkBg = lipgloss.Color("")
		darkPanel = lipgloss.Color("")
		gridLine = lipgloss.Color("")
	}
	
	// Recreate all styles with new colors
	m.recreateStyles()
}

// recreateStyles recreates all styles with current theme colors
func (m *MainModel) recreateStyles() {
	titleStyle = lipgloss.NewStyle().
		Foreground(synthPink).
		Background(styles.GetBackground()).
		Bold(true).
		Italic(true).
		Padding(0, 3).
		MarginTop(1).
		MarginBottom(1).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(synthBlue).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		Align(lipgloss.Center).
		Underline(true).
		Faint(false)

	sectionTitleStyle = lipgloss.NewStyle().
		Foreground(synthCyan).
		Bold(true).
		Italic(true).
		MarginBottom(1).
		MarginTop(1).
		Padding(0, 1).
		Background(styles.GetBoxStyle().GetBackground()).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthBlue)

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
		Foreground(synthBlue).
		Background(styles.GetBoxStyle().GetBackground()).
		Bold(true).
		Italic(true).
		Padding(0, 2).
		MarginBottom(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(synthPink).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(true).
		BorderBottom(true).
		Underline(true)

	completedTaskStyle = lipgloss.NewStyle().
		Foreground(synthPurple).
		Strikethrough(true).
		Italic(true).
		Faint(true)

	statsStyle = lipgloss.NewStyle().
		Foreground(synthYellow).
		Background(styles.GetBoxStyle().GetBackground()).
		Bold(true).
		Padding(1, 2).
		Margin(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthPink).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(true).
		BorderBottom(true)

	helpStyle = lipgloss.NewStyle().
		Foreground(synthRed).
		Italic(true).
		Bold(true).
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthRed)

	glitchStyle = lipgloss.NewStyle().
		Foreground(synthRed).
		Background(styles.GetBoxStyle().GetBackground()).
		Bold(true).
		Blink(true).
		Italic(true).
		Padding(0, 2).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(synthPink)

	statusBarStyle = lipgloss.NewStyle().
		Foreground(synthCyan).
		Background(styles.GetBoxStyle().GetBackground()).
		Bold(true).
		Padding(0, 2).
		MarginTop(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthBlue).
		BorderBottom(true)

	chatInputStyle = lipgloss.NewStyle().
		Foreground(synthPink).
		Background(styles.GetBoxStyle().GetBackground()).
		Bold(true).
		Padding(0, 2).
		MarginTop(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(synthBlue).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		BorderTop(true)

	chatMessageStyle = lipgloss.NewStyle().
		Foreground(synthGreen).
		Italic(true).
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(0, 1).
		MarginBottom(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthGreen).
		BorderLeft(true).
		BorderRight(true)

	chatUserStyle = lipgloss.NewStyle().
		Foreground(synthBlue).
		Bold(true).
		Italic(true).
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(0, 1).
		MarginBottom(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthBlue).
		BorderLeft(true).
		BorderRight(true)

	taskInputStyle = lipgloss.NewStyle().
		Foreground(synthCyan).
		Background(styles.GetBoxStyle().GetBackground()).
		Bold(true).
		Padding(0, 2).
		MarginTop(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(synthGreen).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		BorderTop(true)

	priorityInputStyle = lipgloss.NewStyle().
		Foreground(synthYellow).
		Background(styles.GetBoxStyle().GetBackground()).
		Bold(true).
		Italic(true).
		Padding(0, 2).
		MarginTop(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(synthYellow).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		BorderTop(true)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle digital artifacts/glitches
	if m.showGlitch {
		m.glitchCount++
		if m.glitchCount > 10 {
			m.showGlitch = false
			m.glitchCount = 0
		}
	}
	
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
			case key.Matches(msg, m.keys.confirmAdd): // Use enter
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
			case key.Matches(msg, m.keys.confirmAdd): // Use enter key
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
				if msg.Type == tea.KeyRunes {
					m.notesInput += msg.String()
				} else if msg.Type == tea.KeyBackspace {
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
				case key.Matches(msg, m.keys.confirmAdd):
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
				if msg.Type == tea.KeyRunes {
					m.taskInput += msg.String()
				} else if msg.Type == tea.KeyBackspace {
					if len(m.taskInput) > 0 {
						m.taskInput = m.taskInput[:len(m.taskInput)-1]
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
					m.chatHistory = append(m.chatHistory, 
						chatUserStyle.Render("You: ")+m.chatInput,
						"🤖 AI is thinking...")
					
					// Process chat message in background
					go func() {
						response := m.processChatMessage(m.chatInput)
						// Replace thinking message with actual response
						if len(m.chatHistory) >= 2 {
							m.chatHistory[len(m.chatHistory)-1] = chatMessageStyle.Render("AI: "+response)
						}
						m.aiThinking = false
						m.chatInput = ""
					}()
				}
				return m, nil
				
			case key.Matches(msg, m.keys.tab):
				// Allow tab navigation from chat view
				m.currentView = settingsView
				audio.PlaySound(audio.SoundNavigate)
				return m, nil
				
			default:
				// Handle text input
				if msg.Type == tea.KeyRunes {
					m.chatInput += msg.String()
				} else if msg.Type == tea.KeyBackspace {
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
	if m.aiStatus == "offline" && time.Since(m.lastAICheck) < 11*time.Second {
		if m.aiManager != nil {
			go func() {
				if err := m.aiManager.LaunchOllama(); err == nil {
					m.aiStatus = "online"
				}
			}()
		}
	}
	
	// Handle normal dashboard mode
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global keys that should work in ANY mode (except task entry)
		switch {
		case key.Matches(msg, m.keys.quit):
			m.quitting = true
			return m, tea.Quit
			
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
			
		case key.Matches(msg, m.keys.chat):
			if !m.taskEntryMode {
				m.currentView = chatView
				audio.PlaySound(audio.SoundNavigate)
			}
			return m, nil
			
		case key.Matches(msg, m.keys.tab):
			// ALWAYS ALLOW TAB NAVIGATION - NO BLOCKING CONDITIONS
			// COMPLETE USER-CENTRIC TAB NAVIGATION
			// Order: Dashboard → Timer → Calendar → Notes → Chat → Settings → Dashboard
			switch m.currentView {
			case dashboardView:
				// Tasks are part of dashboard view, so next is Timer
				if len(m.tasks) > 0 && m.selectedTask >= 0 {
					// Start timer for selected task
					task := m.tasks[m.selectedTask]
					m.currentView = focusView
					m.activeTimer = &TimerSession{
						Task: task,
						Mode: "work",
					}
				} else {
					m.currentView = focusView
				}
			case focusView:
				m.currentView = calendarView
			case calendarView:
				m.currentView = notesView
			case notesView:
				m.currentView = chatView
			case chatView:
				m.currentView = settingsView
			case settingsView:
				m.currentView = dashboardView
			default:
				m.currentView = dashboardView
			}
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
			// ALWAYS ALLOW TAB NAVIGATION - NO BLOCKING CONDITIONS
			// COMPLETE USER-CENTRIC TAB NAVIGATION
			// Order: Dashboard → Timer → Calendar → Notes → Chat → Settings → Dashboard
			switch m.currentView {
			case dashboardView:
				// Tasks are part of dashboard view, so next is Timer
				if len(m.tasks) > 0 && m.selectedTask >= 0 {
					// Start timer for selected task
					task := m.tasks[m.selectedTask]
					m.currentView = focusView
					m.activeTimer = &TimerSession{
						Task: task,
						Mode: "work",
					}
				} else {
					m.currentView = focusView
				}
			case focusView:
				m.currentView = calendarView
			case calendarView:
				m.currentView = notesView
			case notesView:
				m.currentView = chatView
			case chatView:
				m.currentView = settingsView
			case settingsView:
				m.currentView = dashboardView
			default:
				m.currentView = dashboardView
			}
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

		case key.Matches(msg, m.keys.calendarKey):
			if !m.taskEntryMode {
				m.currentView = calendarView
				// Load tasks into calendar when switching to calendar view
				if m.cal != nil {
					m.loadTasksIntoCalendar()
				}
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

		case key.Matches(msg, m.keys.glitchTest):
			audio.PlaySound(audio.SoundGlitch)
			m.showGlitch = true
			m.glitchCount = 0
			return m, nil

		case key.Matches(msg, m.keys.themeCycle):
			// Cycle through Kyanite themes
			styles.CycleTheme()
			currentTheme := styles.GetTheme()
			
			// Update all colors and styles
			m.updateTheme()
			
			// Add theme change message to chat history
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

		case key.Matches(msg, m.keys.notes):
			if !m.taskEntryMode && len(m.tasks) > 0 {
				// Launch notes for selected task
				task := &m.tasks[m.selectedTask]
				m.editingTask = task
				m.notesMode = true
				m.notesInput = task.Notes
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
				if m.sessionType == workSession {
					// Stop the timer (Bubble Tea handles stopping)
					// Timer is already stopped when TimeoutMsg is not sent
				}
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
	}

	m.help.ShowAll = false
	return m, nil
}

func (m MainModel) processChatMessage(message string) string {
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
	
	// Call real AI
	response, err := m.aiManager.ChatAssistant(context.Background(), message, taskDescriptions)
	if err != nil {
		// Fallback response if AI fails
		return fmt.Sprintf("🤖 AI unavailable (%s). However, you have %d tasks. Try starting Ollama with: ollama serve", err.Error(), len(m.tasks))
	}
	
	return response
}

func (m MainModel) formatChatHistory() string {
	content := ""
	for _, msg := range m.chatHistory {
		content += msg + "\n"
	}
	return content
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

func (m MainModel) getFilteredTasks() []DashboardTask {
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

func (m MainModel) renderProgressBar() string {
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
	bar.WriteString(fmt.Sprintf("%d%%", percentage))

	return bar.String()
}

func (m MainModel) renderStats() string {
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
	
	// Render enhanced stats with borders and flair
	statsContent := ""
	for i, stat := range stats {
		statStyle := lipgloss.NewStyle().
			Foreground(synthYellow).
			Bold(true).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(0, 1).
			MarginLeft(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthPink)
		
		if i%2 == 0 {
			statStyle = statStyle.Foreground(synthCyan).BorderForeground(synthBlue)
		}
		
		statsContent += statStyle.Render("  "+stat+"  ")
		if i < len(stats)-1 {
			statsContent += " "
		}
		if i%2 == 1 || i == len(stats)-1 {
			statsContent += "\n"
		}
	}
	
	// Enhanced stats container with more flair
	statsContainer := lipgloss.NewStyle().
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(1, 2).
		Margin(1).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(synthPink).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)
	
	return statsContainer.Render(statsContent)
}

func (m MainModel) renderTaskList() string {
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
				Render("😴 No missions in the grid - ready for new quests!")
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
			priorityColor := synthGreen
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

func (m MainModel) renderFocusView() string {
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

func (m MainModel) renderChatView() string {
	var b strings.Builder

	// Enhanced chat header with more flair
	header := lipgloss.NewStyle().
		Foreground(synthPink).
		Background(styles.GetBoxStyle().GetBackground()).
		Bold(true).
		Italic(true).
		Align(lipgloss.Center).
		Padding(1, 3).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(synthCyan).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		Underline(true).
		Render("💬 SYNTHWAVE CHAT ASSISTANT 🤖")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Enhanced chat history with visual separators
	for i, msg := range m.chatHistory {
		if m.aiThinking && i == len(m.chatHistory)-1 && strings.Contains(msg, "AI is thinking") {
			// Replace with animated spinner
			spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			spinner := spinners[m.aiSpinnerFrame%len(spinners)]
			thinkingMsg := fmt.Sprintf("%s 🤖 AI is thinking%s", spinner, strings.Repeat(".", (m.aiSpinnerFrame/3)%4))
			thinkingStyled := lipgloss.NewStyle().
				Foreground(synthCyan).
				Italic(true).
				Bold(true).
				Background(styles.GetBoxStyle().GetBackground()).
				Padding(0, 2).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(synthBlue).
				Render(thinkingMsg)
			b.WriteString(thinkingStyled)
		} else {
			// Add visual styling to each message
			if strings.Contains(msg, "You:") {
				msgStyle := lipgloss.NewStyle().
					Foreground(synthBlue).
					Bold(true).
					Background(styles.GetBoxStyle().GetBackground()).
					Padding(0, 1).
					MarginLeft(2).
					MarginRight(2).
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(synthBlue).
					Render(msg)
				b.WriteString(msgStyle)
			} else if strings.Contains(msg, "AI:") {
				msgStyle := lipgloss.NewStyle().
					Foreground(synthGreen).
					Italic(true).
					Background(styles.GetBoxStyle().GetBackground()).
					Padding(0, 1).
					MarginLeft(2).
					MarginRight(2).
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(synthGreen).
					Render(msg)
				b.WriteString(msgStyle)
			} else {
				msgStyle := lipgloss.NewStyle().
					Foreground(synthPurple).
					Italic(true).
					Background(styles.GetBoxStyle().GetBackground()).
					Padding(0, 1).
					MarginLeft(2).
					MarginRight(2).
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(synthPurple).
					Render(msg)
				b.WriteString(msgStyle)
			}
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Enhanced input prompt with more flair
	inputBox := lipgloss.NewStyle().
		Foreground(synthPink).
		Background(styles.GetBoxStyle().GetBackground()).
		Bold(true).
		Padding(0, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(synthCyan).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		BorderTop(true).
		Render("You: " + m.chatInput + "█")
	b.WriteString(inputBox)
	b.WriteString("\n\n")
	
	// Enhanced controls with more visual flair
	controls := lipgloss.NewStyle().
		Foreground(synthBlue).
		Italic(true).
		Background(styles.GetBoxStyle().GetBackground()).
		Bold(true).
		Padding(0, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthYellow).
		Render("⌨️ [Enter] Send 📤  [Esc] Back 🔙  [Tab] Next View 📋")
	b.WriteString(controls)

	return b.String()
}

func (m MainModel) renderTaskEntry() string {
	var b strings.Builder
	
	if m.taskInput == "PRIORITY_MODE" {
		// Priority change mode
		header := lipgloss.NewStyle().
			Foreground(synthYellow).
			Bold(true).
			Align(lipgloss.Center).
			Render("🎚️ CHANGE TASK PRIORITY")
		b.WriteString(header)
		b.WriteString("\n\n")
		b.WriteString(priorityInputStyle.Render("Select Priority:"))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("  [1] Low 💤"))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(synthGreen).Render("  [2] Medium ⚡"))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(synthRed).Render("  [3] High 🔥"))
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(synthYellow).Render("[Esc] Cancel"))
	} else {
		// Task entry mode
		header := lipgloss.NewStyle().
			Foreground(synthGreen).
			Bold(true).
			Align(lipgloss.Center).
			Render("➕ ADD NEW MISSION")
		b.WriteString(header)
		b.WriteString("\n\n")
		b.WriteString(taskInputStyle.Render("Mission: " + m.taskInput + "█"))
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("[Enter] Confirm  [Esc] Cancel"))
	}
	
	return b.String()
}

func (m MainModel) renderCalendarView() string {
	var b strings.Builder
	
	// Calendar header
	header := lipgloss.NewStyle().
		Foreground(synthPink).
		Bold(true).
		Align(lipgloss.Center).
		Render("📅 SYNTHWAVE CALENDAR")
	b.WriteString(header)
	b.WriteString("\n\n")
	
	// Navigation info
	navInfo := lipgloss.NewStyle().
		Foreground(synthCyan).
		Render("←/h: Previous Month  →/l: Next Month  Tab: Switch View")
	b.WriteString(navInfo)
	b.WriteString("\n\n")
	
	// Calendar content
	if m.cal != nil && m.calRenderer != nil {
		// Update selected date
		m.cal.SelectedDate = m.calSelectedDate
		
		// Render based on view mode
		var calendarContent string
		switch m.calViewMode {
		case "month":
			calendarContent = m.calRenderer.RenderMonth(m.cal)
		case "week":
			calendarContent = m.calRenderer.RenderWeek(m.cal)
		case "day":
			calendarContent = m.calRenderer.RenderDay(m.cal)
		default:
			calendarContent = m.calRenderer.RenderMonth(m.cal)
		}
		
		b.WriteString(lipgloss.NewStyle().
			Foreground(synthGreen).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthBlue).
			Render(calendarContent))
	} else {
		// Fallback if calendar not initialized
		b.WriteString(lipgloss.NewStyle().
			Foreground(synthYellow).
			Render("📅 Calendar initializing..."))
	}
	
	b.WriteString("\n\n")
	
	// Current month info
	monthYear := m.calSelectedDate.Format("January 2006")
	monthInfo := lipgloss.NewStyle().
		Foreground(synthYellow).
		Bold(true).
		Render("📅 Currently viewing: " + monthYear)
	b.WriteString(monthInfo)
	b.WriteString("\n\n")
	
	// Help text
	helpText := lipgloss.NewStyle().Foreground(synthBlue).Render("Press 'Tab' to return to dashboard")
	b.WriteString(helpText)
	
	return b.String()
}

func (m MainModel) renderNotesEditor() string {
	var b strings.Builder
	
	if m.editingTask == nil {
		return "No task selected for notes"
	}
	
	// Notes editor header
	header := lipgloss.NewStyle().
		Foreground(synthCyan).
		Bold(true).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("📝 EDITING NOTES FOR: %s", m.editingTask.Description))
	b.WriteString(header)
	b.WriteString("\n\n")
	
	// Current notes display with Glow markdown rendering
	b.WriteString(lipgloss.NewStyle().Foreground(synthGreen).Bold(true).Render("📝 Current notes (Markdown):"))
	b.WriteString("\n")
	
	// Show notes using Glow markdown renderer
	currentNotes := m.editingTask.Notes
	if currentNotes == "" {
		currentNotes = "*No notes yet - start writing in markdown!*\n\n## Examples:\n- **Bold text**\n- *Italic text*\n- `Code snippets`\n- # Headers\n- [Links](url)"
	}
	
	// Use Glow for enhanced markdown rendering
	if m.glowStyler != nil {
		glowContent := m.glowStyler.RenderSectionWithGlow(
			"📝 Notes Content",
			currentNotes,
			"#00FFF0", // Cyan accent
		)
		b.WriteString(glowContent)
	} else {
		// Fallback to basic rendering
		notesStyle := lipgloss.NewStyle().
			Foreground(synthYellow).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(synthBlue).
			Width(70).
			Render(currentNotes)
		b.WriteString(notesStyle)
	}
	b.WriteString("\n\n")
	
	// Input area
	b.WriteString(lipgloss.NewStyle().Foreground(synthGreen).Render("Enter new notes:"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(synthCyan).
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(synthGreen).
		Width(70).
		Height(6).
		Render(m.notesInput + "█"))
	b.WriteString("\n\n")
	
	// Instructions
	instructions := []string{
		"[Enter] Save notes",
		"[Esc] Cancel",
		"Type to add notes",
	}
	
	for _, instruction := range instructions {
		b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  " + instruction))
		b.WriteString("\n")
	}
	
	return b.String()
}

func (m MainModel) renderFilterEditor() string {
	var b strings.Builder
	
	// Filter editor header
	header := lipgloss.NewStyle().
		Foreground(synthYellow).
		Bold(true).
		Align(lipgloss.Center).
		Render("🔍 TASK FILTER")
	b.WriteString(header)
	b.WriteString("\n\n")
	
	// Current filter status
	b.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("Current filter:"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthGreen).
		Render(fmt.Sprintf("Status: %s | Priority: %s", m.filterStatus, m.filterPriority)))
	b.WriteString("\n\n")
	
	// Filter options
	b.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("Status options:"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [1] All tasks"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [2] Pending only"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [3] Completed only"))
	b.WriteString("\n\n")
	
	b.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("Priority options:"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [H] High priority"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [M] Medium priority"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [L] Low priority"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [A] All priorities"))
	b.WriteString("\n\n")
	
	// Instructions
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("[Enter] Apply filter  [Esc] Cancel"))
	
	return b.String()
}

func (m MainModel) renderSettingsView() string {
	var b strings.Builder
	
	// Settings header
	header := lipgloss.NewStyle().
		Foreground(synthPink).
		Bold(true).
		Align(lipgloss.Center).
		Render("⚙️ SYNTHWAVE SETTINGS")
	b.WriteString(header)
	b.WriteString("\n\n")
	
	// Current settings display
	settings := []string{
		fmt.Sprintf("🎨 Theme: %s", m.themeNames[m.currentThemeIndex]),
		fmt.Sprintf("🔊 Audio: %t", m.audioEnabled),
		fmt.Sprintf("⏰ Work Duration: %v", m.workDuration),
		fmt.Sprintf("☕ Break Duration: %v", m.breakDuration),
		fmt.Sprintf("📅 Calendar View: %s", m.calViewMode),
	}
	
	for _, setting := range settings {
		b.WriteString(lipgloss.NewStyle().
			Foreground(synthCyan).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(synthBlue).
			Render(setting))
		b.WriteString("\n")
	}
	
	b.WriteString("\n")
	
	// Controls info
	controls := []string{
		"[T] Cycle theme",
		"[M] Toggle audio",
		"[W] Change work duration",
		"[B] Change break duration",
		"[Tab] Switch view",
	}
	
	b.WriteString(sectionTitleStyle.Render("⚡ SETTINGS CONTROLS:"))
	b.WriteString("\n")
	for _, control := range controls {
		b.WriteString(lipgloss.NewStyle().Foreground(synthGreen).Render("  " + control))
		b.WriteString("\n")
	}
	
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("Press 'Tab' to return to dashboard"))
	
	return b.String()
}

func (m MainModel) renderDashboard() string {
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

	// Show glitch if active
	if m.showGlitch {
		b.WriteString(glitchStyle.Render(m.glitchMessage))
		b.WriteString("\n\n")
	}

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
			Width(m.width - 4).
			Height(m.height - 4).
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
			Width(m.width - 4).
			Height(m.height - 4).
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
			Width(m.width - 4).
			Height(m.height - 4).
			Align(lipgloss.Center, lipgloss.Center).
			Render(overlay)
			
		b.WriteString(container)
		return b.String()
	}

	// Main dashboard layout
	var leftContent, centerContent, rightContent string

	// Left column - Task list
	leftContent = sectionTitleStyle.Render("📋 MISSION BOARD:") + "\n\n" + m.renderTaskList()

	// Center column - Selected task
	centerContent = sectionTitleStyle.Render("🎯 SELECTED MISSION:") + "\n"
	if len(m.tasks) > 0 && m.selectedTask < len(m.tasks) {
		task := m.tasks[m.selectedTask]
		status := "PENDING"
		if task.Status == "completed" {
			status = "COMPLETED ✅"
		}
		
		centerContent += lipgloss.NewStyle().
			Foreground(synthCyan).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthBlue).
			Render(task.Description) + "\n\n"
		centerContent += lipgloss.NewStyle().Foreground(synthGreen).Render(
			fmt.Sprintf("Status: %s\nPriority: %s", status, strings.ToUpper(task.Priority)))
	} else {
		centerContent += "😴 No missions selected"
	}
	centerContent += "\n\n" + lipgloss.NewStyle().Foreground(synthBlue).Render(
		"💡 [ENTER] Focus mode\n💡 [A] Add mission\n💡 [C] Chat assistant\n💡 [D] Complete\n💡 [P] Priority\n💡 [G] Glitch")

	// Right column - Stats and quick actions
	rightContent = m.renderStats()
	rightContent += "\n\n" + sectionTitleStyle.Render("⚡ REAL-TIME DATA:") + "\n"
	rightContent += lipgloss.NewStyle().Foreground(synthCyan).Render("🕐 " + time.Now().Format("15:04:05")) + "\n"
	rightContent += lipgloss.NewStyle().Foreground(synthYellow).Render("📅 " + time.Now().Format("2006-01-02")) + "\n\n"
	rightContent += sectionTitleStyle.Render("⚡ QUICK ACTIONS:") + "\n"
	actions := []string{
		"[A] Add mission",
		"[D] Complete task", 
		"[P] Change priority",
		"[C] Chat assistant", 
		"[F] Focus mode",
		"[G] Test glitch",
		"[↑/↓] Navigate",
	}
	for _, action := range actions {
		rightContent += lipgloss.NewStyle().Foreground(synthCyan).Render("  " + action) + "\n"
	}

	// Apply width constraints safely
	if m.width > 80 { // Only apply if we have enough width
		colWidth := (m.width - 6) / 3 // Account for padding
		if colWidth < 30 { // Minimum column width
			colWidth = 30
		}
		
		leftStyle := lipgloss.NewStyle().Width(colWidth)
		centerStyle := lipgloss.NewStyle().Width(colWidth)
		rightStyle := lipgloss.NewStyle().Width(colWidth)
		
		leftContent = leftStyle.Render(leftContent)
		centerContent = centerStyle.Render(centerContent)
		rightContent = rightStyle.Render(rightContent)
	}
	
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftContent, centerContent, rightContent)
	b.WriteString(layout)

	return b.String()
}

func (m MainModel) View() string {
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
		case settingsView:
			content = m.renderSettingsView()
		default:
			content = m.renderDashboard()
		}
	}

	// Main container with theme-aware background
	mainContainer := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(synthPink).
		Padding(1).
		Background(styles.GetBackground()).
		Width(m.width - 4).
		Height(m.height - 4)

	// Add help at the bottom
	helpText := helpStyle.Render("\n🎮 " + m.help.View(m.keys))

	return mainContainer.Render(content + helpText)
}

// Overlay helper functions
func (m MainModel) renderTaskEntryOverlay() string {
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
		Width(m.width - 10).
		Height(m.height - 15).
		Align(lipgloss.Center, lipgloss.Center).
		Render(container)
	
	return outerContainer
}

func (m MainModel) renderNotesOverlay() string {
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

func (m MainModel) renderFilterOverlay() string {
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

func StartMainDashboard(tasks []DashboardTask) error {
	model := NewMainModel(tasks)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
