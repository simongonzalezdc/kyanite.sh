package dashboard

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SystemStats holds the system statistics
type SystemStats struct {
	SongCount    int
	ProjectCount int
	StorageBytes int64
}

// systemStatsLoadedMsg is sent when stats are loaded
type systemStatsLoadedMsg struct {
	stats SystemStats
}

// systemStatsLoadErrorMsg is sent when loading fails
type systemStatsLoadErrorMsg struct {
	err error
}

// SystemInfoModel manages the system info panel
type SystemInfoModel struct {
	width    int
	height   int
	database *db.DB
	stats    SystemStats
	loading  bool
	err      error
	hovered  int // -1 = none, 0 = settings button
}

// NewSystemInfoModel creates a new system info model
func NewSystemInfoModel() *SystemInfoModel {
	return &SystemInfoModel{
		loading: true,
		hovered: -1,
	}
}

// SetDatabase sets the database reference
func (m *SystemInfoModel) SetDatabase(database *db.DB) {
	m.database = database
}

// Init initializes the model and loads stats
func (m *SystemInfoModel) Init() tea.Cmd {
	return m.loadStats()
}

// loadStats loads statistics from the database
func (m *SystemInfoModel) loadStats() tea.Cmd {
	return func() tea.Msg {
		if m.database == nil {
			return systemStatsLoadErrorMsg{err: fmt.Errorf("database not available")}
		}

		var stats SystemStats

		// Count songs
		songs, err := m.database.ListSongs(1000, 0) // Get all songs to count
		if err == nil {
			stats.SongCount = len(songs)
		}

		// Count projects
		projects, err := m.database.ListProjects()
		if err == nil {
			stats.ProjectCount = len(projects)
		}

		// Calculate storage size (data directory)
		stats.StorageBytes = m.calculateStorageSize()

		return systemStatsLoadedMsg{stats: stats}
	}
}

// calculateStorageSize calculates the size of the data directory
func (m *SystemInfoModel) calculateStorageSize() int64 {
	// Check common data locations
	dataDirs := []string{
		"data",
		filepath.Join(os.Getenv("HOME"), ".noise"),
	}

	var totalSize int64
	for _, dir := range dataDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() {
					totalSize += info.Size()
				}
				return nil
			})
		}
	}

	return totalSize
}

// Update handles messages for the system info panel
func (m *SystemInfoModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil

	case systemStatsLoadedMsg:
		m.stats = msg.stats
		m.loading = false
		m.err = nil
		return nil

	case systemStatsLoadErrorMsg:
		m.loading = false
		m.err = msg.err
		return nil

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.openSettings()
		}
	}

	return nil
}

// handleMouse processes mouse events
func (m *SystemInfoModel) handleMouse(msg tea.MouseMsg) tea.Cmd {
	if m.width == 0 {
		return nil
	}

	// Button is near the bottom of the panel
	buttonY := 10 // Approximate position

	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionRelease {
			if msg.Y >= buttonY && msg.Y <= buttonY+1 {
				return m.openSettings()
			}
		}

	case tea.MouseButtonNone:
		if msg.Action == tea.MouseActionMotion {
			if msg.Y >= buttonY && msg.Y <= buttonY+1 {
				m.hovered = 0
			} else {
				m.hovered = -1
			}
		}
	}

	return nil
}

// openSettings navigates to the settings screen
func (m *SystemInfoModel) openSettings() tea.Cmd {
	return func() tea.Msg {
		return ScreenChangeMsg{Screen: 7} // screenSettings
	}
}

// View renders the system info panel
func (m *SystemInfoModel) View() string {
	if m.width == 0 {
		return "System Info"
	}

	t := theme.GetManager().Current()

	title := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Render("System Info")

	// Show loading state
	if m.loading {
		loadingText := lipgloss.NewStyle().
			Foreground(t.Secondary).
			Faint(true).
			Render("Loading...")

		content := lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			loadingText,
		)

		return lipgloss.NewStyle().
			Width(m.width-2).
			MaxWidth(m.width-2).
			MaxHeight(m.height-2).
			Padding(0, 1).
			Render(content)
	}

	// Statistics section
	statsTitle := lipgloss.NewStyle().
		Foreground(t.Text).
		Bold(true).
		Render("[STATS] Statistics:")

	stats := []string{
		fmt.Sprintf("- Songs: %d", m.stats.SongCount),
		fmt.Sprintf("- Projects: %d", m.stats.ProjectCount),
	}

	var statViews []string
	for _, stat := range stats {
		statView := lipgloss.NewStyle().
			Foreground(t.Text).
			Render(stat)
		statViews = append(statViews, statView)
	}

	// Storage info
	storageStr := formatBytes(m.stats.StorageBytes)
	storage := lipgloss.NewStyle().
		Foreground(t.Text).
		Render(fmt.Sprintf("[DISK] Storage: %s", storageStr))

	// Performance indicator (simple heuristic)
	perfStatus := "Good"
	perfColor := t.Success
	if m.stats.SongCount > 500 {
		perfStatus = "Heavy"
		perfColor = t.Warning
	}
	performance := lipgloss.NewStyle().
		Foreground(perfColor).
		Render(fmt.Sprintf("[PERF] Performance: %s", perfStatus))

	// Open settings button
	buttonStyle := lipgloss.NewStyle().
		Padding(0, 1).
		MarginTop(1)

	if m.hovered == 0 {
		buttonStyle = buttonStyle.
			Foreground(t.Background).
			Background(t.Primary).
			Bold(true)
	} else {
		buttonStyle = buttonStyle.
			Foreground(t.Text).
			Background(t.Secondary)
	}

	settingsButton := buttonStyle.Render("[Open Settings]")

	// Mouse hint
	mouseHint := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Faint(true).
		Render("Click to configure")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		statsTitle,
		lipgloss.JoinVertical(lipgloss.Left, statViews...),
		"",
		storage,
		performance,
		"",
		settingsButton,
		mouseHint,
	)

	return lipgloss.NewStyle().
		Width(m.width-2).
		MaxWidth(m.width-2).
		MaxHeight(m.height-2).
		Padding(0, 1).
		Render(content)
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	if unit > 0 {
		for n := bytes / unit; n >= unit; n /= unit {
			div *= unit
			exp++
		}
	}

	units := []string{"KB", "MB", "GB", "TB"}
	if div > 0 && exp < len(units) {
		return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
	}
	return fmt.Sprintf("%d B", bytes)
}
