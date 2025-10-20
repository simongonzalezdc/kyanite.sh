package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Clean TUI Model
type Model struct {
	store         *store.Store
	engine        *engine.Engine
	tasks         []models.Task
	currentView   string
	selectedIndex int
	width         int
	height        int
}

func NewModel() *Model {
	store := store.New(utils.GetStoragePath())
	engine := engine.New(store)

	tasks, _ := engine.ListTasks("all")

	// Convert engine tasks to models
	modelTasks := make([]models.Task, len(tasks))
	for i, task := range tasks {
		modelTasks[i] = models.Task{
			ID:          task.ID,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	return &Model{
		store:         store,
		engine:        engine,
		tasks:         modelTasks,
		currentView:   "dashboard",
		selectedIndex: 0,
	}
}

// Tea initialization
func (m *Model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type TickMsg time.Time

// Main update loop
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}

		case "down", "j":
			if m.selectedIndex < len(m.tasks)-1 {
				m.selectedIndex++
			}

		case "r":
			m.refreshTasks()

		case "1":
			m.currentView = "dashboard"
		case "2":
			m.currentView = "matrix"
		case "3":
			m.currentView = "about"
		}

	case TickMsg:
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// Main view rendering
func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return styles.LoadingMessage()
	}

	switch m.currentView {
	case "dashboard":
		return m.renderDashboard()
	case "matrix":
		return m.renderMatrix()
	case "about":
		return m.renderGlitchRoom()
	default:
		return m.renderDashboard()
	}
}

func (m *Model) refreshTasks() {
	tasks, _ := m.engine.ListTasks("all")

	// Convert engine tasks to models
	modelTasks := make([]models.Task, len(tasks))
	for i, task := range tasks {
		modelTasks[i] = models.Task{
			ID:          task.ID,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}
	m.tasks = modelTasks
}

// Dashboard View
func (m *Model) renderDashboard() string {
	var content strings.Builder

	// Clean header
	header := styles.Header()
	content.WriteString(header)
	content.WriteString("\n\n")

	// Task Grid
	if len(m.tasks) == 0 {
		content.WriteString(styles.EmptyStateMessage())
	} else {
		// Render tasks with clean styling
		for i, task := range m.tasks {
			taskStyle := m.getTaskStyle(i, task)
			content.WriteString(taskStyle)
			content.WriteString("\n")

			// Add metadata
			if len(task.Categories) > 0 {
				tags := make([]string, len(task.Categories))
				for j, cat := range task.Categories {
					tags[j] = styles.CyberTag(cat)
				}
				content.WriteString("   " + strings.Join(tags, " "))
				content.WriteString("\n")
			}

			// Task ID
			idStyle := lipgloss.NewStyle().
				Foreground(styles.SynthwaveCyan).
				Background(styles.DeepSpace).
				Render("ID: " + task.ID)
			content.WriteString("   " + idStyle)
			content.WriteString("\n\n")
		}
	}

	// Stats Footer
	active := 0
	completed := 0
	for _, task := range m.tasks {
		if task.Status == "completed" {
			completed++
		} else {
			active++
		}
	}
	stats := styles.CyberStats(active, completed, len(m.tasks))
	content.WriteString("\n")
	content.WriteString(stats)
	content.WriteString("\n\n")

	// Controls
	controls := m.renderControls()
	content.WriteString(controls)

	return content.String()
}

// Clean Matrix View
func (m *Model) renderMatrix() string {
	var content strings.Builder

	title := lipgloss.NewStyle().
		Foreground(styles.SynthwaveGreen).
		Background(styles.DeepSpace).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Render("Task Matrix")

	content.WriteString(title)
	content.WriteString("\n\n")

	// Display tasks in grid format
	if len(m.tasks) == 0 {
		content.WriteString(styles.EmptyStateMessage())
	} else {
		for i, task := range m.tasks {
			taskLine := fmt.Sprintf("%02d: %s", i+1, task.Description)
			taskStyle := lipgloss.NewStyle().
				Foreground(styles.SynthwaveCyan).
				Background(styles.DeepSpace).
				Render(taskLine)
			content.WriteString(taskStyle)
			content.WriteString("\n")
		}
	}

	// Overlay task info
	if len(m.tasks) > 0 && m.selectedIndex < len(m.tasks) {
		task := m.tasks[m.selectedIndex]
		taskInfo := styles.FocusBox(
			fmt.Sprintf("SELECTED: %s\nSTATUS: %s\nPRIORITY: %s",
				task.Description,
				task.Status,
				task.Priority,
			),
			styles.SynthwavePink,
		)
		content.WriteString("\n")
		content.WriteString(taskInfo)
	}

	return content.String()
}

// Stats View with Holographic Effects
func (m *Model) renderStats() string {
	var content strings.Builder

	title := styles.HolographicText("📊 MISSION ANALYTICS")
	content.WriteString(title)
	content.WriteString("\n\n")

	// Calculate statistics
	active := 0
	completed := 0
	highPriority := 0
	mediumPriority := 0
	lowPriority := 0

	for _, task := range m.tasks {
		if task.Status == "completed" {
			completed++
		} else {
			active++
		}

		switch task.Priority {
		case "high":
			highPriority++
		case "medium":
			mediumPriority++
		case "low":
			lowPriority++
		}
	}

	total := len(m.tasks)
	completionRate := float64(completed) / float64(total) * 100

	// Render stats with style
	stats := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"TOTAL MISSIONS", fmt.Sprintf("%d", total), styles.SynthwaveCyan},
		{"ACTIVE MISSIONS", fmt.Sprintf("%d", active), styles.SynthwaveYellow},
		{"COMPLETED MISSIONS", fmt.Sprintf("%d", completed), styles.SynthwaveGreen},
		{"COMPLETION RATE", fmt.Sprintf("%.1f%%", completionRate), styles.SynthwavePink},
		{"HIGH PRIORITY", fmt.Sprintf("%d", highPriority), styles.SynthwaveRed},
		{"MEDIUM PRIORITY", fmt.Sprintf("%d", mediumPriority), styles.SynthwaveOrange},
		{"LOW PRIORITY", fmt.Sprintf("%d", lowPriority), styles.SynthwaveGreen},
	}

	for _, stat := range stats {
		statLine := lipgloss.NewStyle().
			Foreground(stat.color).
			Background(styles.DarkVoid).
			Bold(true).
			Padding(0, 2).
			Render(fmt.Sprintf("▸ %s: %s", stat.label, stat.value))
		content.WriteString(statLine)
		content.WriteString("\n")
	}

	// Progress bar with synthwave style
	progressWidth := 40
	filled := int(float64(progressWidth) * completionRate / 100)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", progressWidth-filled)

	progressBar := lipgloss.NewStyle().
		Foreground(styles.SynthwaveGreen).
		Background(styles.DeepSpace).
		Render("PROGRESS: [" + bar + "]")
	content.WriteString("\n")
	content.WriteString(progressBar)

	return content.String()
}

// Clean About View
func (m *Model) renderGlitchRoom() string {
	var content strings.Builder

	title := styles.Title("About focus.sh")
	content.WriteString(title)
	content.WriteString("\n\n")

	about := `A clean, professional task management tool
that follows the Kyanite Suite design principles.

Features:
• Clean, distraction-free interface
• Consistent theming
• Purposeful animations
• Professional appearance`

	aboutStyle := lipgloss.NewStyle().
		Foreground(styles.SynthwaveCyan).
		Background(styles.DeepSpace).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SynthwavePink).
		Render(about)

	content.WriteString(aboutStyle)
	return content.String()
}

// Clean Task Styling
func (m *Model) getTaskStyle(index int, task models.Task) string {
	var prefix string
	var style lipgloss.Style

	if index == m.selectedIndex {
		prefix = "▶ "
		style = lipgloss.NewStyle().
			Foreground(styles.SynthwavePink).
			Background(styles.DarkVoid).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(styles.SynthwaveCyan)
	} else {
		prefix = "  "
		style = lipgloss.NewStyle().
			Foreground(styles.SynthwaveCyan).
			Background(styles.DeepSpace)
	}

	taskText := fmt.Sprintf("%s%s", prefix, task.Description)
	if task.Status == "completed" {
		taskText += " ✓"
	}

	priorityText := styles.PriorityExplosion(task.Priority)

	return style.Render(taskText + " " + priorityText)
}

func (m *Model) renderControls() string {
	controls := []string{
		"↑↓: Navigate",
		"SPACE: Toggle Glitch",
		"R: Refresh",
		"1-4: Switch Views",
		"Q: Quit",
	}

	controlStyle := lipgloss.NewStyle().
		Foreground(styles.SynthwavePurple).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SynthwaveCyan).
		Render("CONTROLS: " + strings.Join(controls, " | "))

	return controlStyle
}
