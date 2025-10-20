package tui

import (
	"fmt"
	"strings"
	"time"
	"math/rand"

	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/models"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// symbolGlitchText replaces letters with pure symbols (no "tetter" effect)
func symbolGlitchText(text string) string {
	symbolPatterns := []string{
		"⚡⚡⚡ ◈◈◈ ◆◆◆ ◊◊◊ ⚡⚡⚡",
		"◆◆◆ ◊◊◊ ⚡⚡⚡ ◈◈◈ ◆◆◆",
		"◈◈◈ ◆◆◆ ◊◊◊ ⚡⚡⚡ ◈◈◈",
		"◊◊◊ ⚡⚡⚡ ◈◈◈ ◆◆◆ ◊◊◊",
		"⚡◆◈◊ ◈⚡◆◈ ◊◈◆⚡ ◊⚡◆◈",
		"◆◈⚡◊ ◆◈◊⚡ ◈⚡◆◊ ◈◆⚡◊",
		"◈◊⚡◆ ◊◈⚡◆ ◊⚡◆◈ ⚡◈◊◆",
	}
	
	// Return a random symbol pattern instead of corrupting text
	if rand.Float32() < 0.3 { // 30% chance to replace with symbols
		return symbolPatterns[rand.Intn(len(symbolPatterns))]
	}
	
	// Otherwise return original text with minimal character replacement
	result := ""
	symbols := []string{"⚡", "◈", "◆", "◊", "◇", "○", "●", "□", "■", "△", "▽"}
	
	for _, char := range text {
		if rand.Float32() < 0.05 { // 5% chance to replace individual chars
			result += symbols[rand.Intn(len(symbols))]
		} else {
			result += string(char)
		}
	}
	return result
}

// Synthwave TUI Model - Maximum Visual Impact
type Model struct {
	store         *store.Store
	engine        *engine.Engine
	tasks         []models.Task
	currentView   string
	selectedIndex int
	width         int
	height        int
	glitchActive  bool
	animationTick int
	digitalNoise  []string
}

// Initialize random for effects
var glitchRand = rand.New(rand.NewSource(time.Now().UnixNano()))

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
	
	// Generate digital noise artifacts
	noise := make([]string, 20)
	for i := range noise {
		noise[i] = styles.CreateDigitalArtifact()
	}
	
	return &Model{
		store:         store,
		engine:        engine,
		tasks:         modelTasks,
		currentView:   "dashboard",
		selectedIndex: 0,
		glitchActive:  false,
		digitalNoise:  noise,
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
			
		case " ":
			m.glitchActive = !m.glitchActive
			
		case "r":
			m.refreshTasks()
			
		case "1":
			m.currentView = "dashboard"
		case "2":
			m.currentView = "matrix"
		case "3":
			m.currentView = "stats"
		case "4":
			m.currentView = "glitch"
		}
		
	case TickMsg:
		m.animationTick++
		// Random glitch activation
		if glitchRand.Float32() < 0.05 { // 5% chance per tick
			m.glitchActive = true
		}
		if m.glitchActive && glitchRand.Float32() < 0.3 { // 30% chance to deactivate
			m.glitchActive = false
		}
		
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
		return styles.LoadingAnimation()
	}
	
	switch m.currentView {
	case "dashboard":
		return m.renderDashboard()
	case "matrix":
		return m.renderMatrix()
	case "stats":
		return m.renderStats()
	case "glitch":
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

// Dashboard View with Maximum Impact
func (m *Model) renderDashboard() string {
	var content strings.Builder
	
	// Header with glitch effects
	header := styles.MatrixHeader()
	if m.glitchActive {
		header = styles.GlitchTitle("SYNTHWAVE MISSION MATRIX")
	}
	content.WriteString(header)
	content.WriteString("\n\n")
	
	// Task Grid
	if len(m.tasks) == 0 {
		content.WriteString(styles.EmptyStateMessage())
	} else {
		// Render tasks with cyber styling
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
			
			// Task ID with digital effect
			idStyle := lipgloss.NewStyle().
				Foreground(styles.SynthwaveCyan).
				Background(styles.DeepSpace).
				Render("🔑 ID: " + task.ID)
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
	
	// Add digital noise overlay if glitching
	if m.glitchActive {
		overlay := m.renderDigitalNoise()
		content = strings.Builder{}
		content.WriteString(overlay)
		content.WriteString("\n")
		content.WriteString(content.String())
	}
	
	return content.String()
}

// Matrix View - Pure Cyberpunk Aesthetic
func (m *Model) renderMatrix() string {
	var content strings.Builder
	
	title := lipgloss.NewStyle().
		Foreground(styles.SynthwaveGreen).
		Background(styles.DeepSpace).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Render("🌐 DIGITAL MATRIX INTERFACE 🌐")
	
	content.WriteString(title)
	content.WriteString("\n\n")
	
	// Create matrix grid effect
	gridWidth := 60
	gridHeight := 15
	
	for y := 0; y < gridHeight; y++ {
		line := ""
		for x := 0; x < gridWidth; x++ {
			if glitchRand.Float32() < 0.1 {
				line += styles.RandomGlitchChar()
			} else if glitchRand.Float32() < 0.05 {
				line += string(rune('0' + glitchRand.Intn(10)))
			} else {
				line += " "
			}
		}
		
		lineStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveGreen).
			Background(styles.DeepSpace).
			Render(line)
		
		content.WriteString(lineStyle)
		content.WriteString("\n")
	}
	
	// Overlay task info
	if len(m.tasks) > 0 && m.selectedIndex < len(m.tasks) {
		task := m.tasks[m.selectedIndex]
		taskInfo := styles.NeonBox(
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

// Glitch Room - Maximum Visual Chaos
func (m *Model) renderGlitchRoom() string {
	var content strings.Builder
	
	title := styles.GlitchTitle("⚡ GLITCH DIMENSION ⚡")
	content.WriteString(title)
	content.WriteString("\n\n")
	
	// Pure visual chaos
	glitchLines := []string{
		"SYSTEM MALFUNCTION DETECTED",
		"REALITY COMPROMISED",
		"DIGITAL ARTIFACTS ACTIVE",
		"NEURAL INTERFACE UNSTABLE",
		"SYNTHWAVE PROTOCOLS CORRUPTED",
	}
	
	for i, line := range glitchLines {
		var glitchedLine string
		
		switch i % 4 {
		case 0:
			glitchedLine = symbolGlitchText(line)
		case 1:
			glitchedLine = styles.RGBSplitText(line)
		case 2:
			glitchedLine = styles.HolographicText(line)
		default:
			glitchedLine = styles.DigitalRain(line)
		}
		
		lineStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", 
				glitchRand.Intn(255), 
				glitchRand.Intn(255), 
				glitchRand.Intn(255)))).
			Bold(glitchRand.Intn(2) == 0).
			Italic(glitchRand.Intn(2) == 0).
			Render(glitchedLine)
		
		content.WriteString(lineStyle)
		content.WriteString("\n")
	}
	
	// Add random artifacts
	for i := 0; i < 10; i++ {
		artifact := styles.CreateDigitalArtifact()
		artifactStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", 
				glitchRand.Intn(255), 
				glitchRand.Intn(255), 
				glitchRand.Intn(255)))).
			Render(artifact)
		content.WriteString(artifactStyle + " ")
	}
	
	content.WriteString("\n\n")
	
	// Interactive element
	interactive := styles.CyberGridBox("PRESS SPACE TO STABILIZE REALITY")
	content.WriteString(interactive)
	
	return content.String()
}

// Helper Functions
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
	
	if m.glitchActive {
		style = style.Background(lipgloss.Color("#FF00FF"))
	}
	
	taskText := fmt.Sprintf("%s%s", prefix, task.Description)
	if task.Status == "completed" {
		taskText += " ✅"
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

func (m *Model) renderDigitalNoise() string {
	var noise strings.Builder
	
	for i := 0; i < 5; i++ {
		line := ""
		for j := 0; j < m.width; j++ {
			if glitchRand.Float32() < 0.3 {
				line += styles.RandomGlitchChar()
			} else {
				line += " "
			}
		}
		
		lineStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", 
				glitchRand.Intn(255), 
				glitchRand.Intn(255), 
				glitchRand.Intn(255)))).
			Render(line)
		
		noise.WriteString(lineStyle)
		noise.WriteString("\n")
	}
	
	return noise.String()
}
