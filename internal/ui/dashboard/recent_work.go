package dashboard

import (
	"fmt"
	"time"

	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RecentItem represents a recent work item (song or project)
type RecentItem struct {
	ID        int
	Title     string
	Type      string // "song" or "project"
	UpdatedAt time.Time
}

// recentWorkLoadedMsg is sent when recent work items are loaded
type recentWorkLoadedMsg struct {
	items []RecentItem
}

// recentWorkLoadErrorMsg is sent when loading fails
type recentWorkLoadErrorMsg struct {
	err error
}

// RecentWorkModel manages the recent work panel
type RecentWorkModel struct {
	width    int
	height   int
	database *db.DB
	items    []RecentItem
	selected int
	hovered  int
	loading  bool
	err      error
}

// NewRecentWorkModel creates a new recent work model
func NewRecentWorkModel() *RecentWorkModel {
	return &RecentWorkModel{
		selected: 0,
		hovered:  -1,
		loading:  true,
	}
}

// SetDatabase sets the database reference for loading recent work
func (m *RecentWorkModel) SetDatabase(database *db.DB) {
	m.database = database
}

// Init initializes the model and loads recent work
func (m *RecentWorkModel) Init() tea.Cmd {
	return m.loadRecentWork()
}

// loadRecentWork loads recent songs from the database
func (m *RecentWorkModel) loadRecentWork() tea.Cmd {
	return func() tea.Msg {
		if m.database == nil {
			return recentWorkLoadErrorMsg{err: fmt.Errorf("database not available")}
		}

		// Load recent songs (limit to 5 for the dashboard panel)
		songs, err := m.database.ListSongs(5, 0)
		if err != nil {
			return recentWorkLoadErrorMsg{err: err}
		}

		var items []RecentItem
		for _, song := range songs {
			items = append(items, RecentItem{
				ID:        song.ID,
				Title:     song.Metadata.Title,
				Type:      "song",
				UpdatedAt: song.Metadata.UpdatedAt,
			})
		}

		return recentWorkLoadedMsg{items: items}
	}
}

// Update handles messages for the recent work panel
func (m *RecentWorkModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil

	case recentWorkLoadedMsg:
		m.items = msg.items
		m.loading = false
		m.err = nil
		return nil

	case recentWorkLoadErrorMsg:
		m.loading = false
		m.err = msg.err
		return nil

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.items)-1 {
				m.selected++
			}
		case "enter":
			return m.openSelectedItem()
		}
	}

	return nil
}

// handleMouse processes mouse events
func (m *RecentWorkModel) handleMouse(msg tea.MouseMsg) tea.Cmd {
	if m.width == 0 || len(m.items) == 0 {
		return nil
	}

	// Calculate item positions
	// Title (1) + empty (1) + items start at row 2
	// Each item is 1 row, button is after items
	itemStartY := 2
	itemEndY := itemStartY + len(m.items) - 1
	buttonY := itemEndY + 2

	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionRelease {
			// Check if click is on an item
			if msg.Y >= itemStartY && msg.Y <= itemEndY {
				m.selected = msg.Y - itemStartY
				return m.openSelectedItem()
			}
			// Check if click is on the button
			if msg.Y >= buttonY && msg.Y <= buttonY+1 {
				return m.openProjectManager()
			}
		} else if msg.Action == tea.MouseActionPress {
			// Visual feedback on press
			if msg.Y >= itemStartY && msg.Y <= itemEndY {
				m.selected = msg.Y - itemStartY
			}
		}

	case tea.MouseButtonWheelUp:
		if m.selected > 0 {
			m.selected--
		}

	case tea.MouseButtonWheelDown:
		if m.selected < len(m.items)-1 {
			m.selected++
		}

	case tea.MouseButtonNone:
		if msg.Action == tea.MouseActionMotion {
			if msg.Y >= itemStartY && msg.Y <= itemEndY {
				m.hovered = msg.Y - itemStartY
			} else if msg.Y >= buttonY && msg.Y <= buttonY+1 {
				m.hovered = len(m.items) // Special index for button
			} else {
				m.hovered = -1
			}
		}
	}

	return nil
}

// openSelectedItem opens the selected song in the editor
func (m *RecentWorkModel) openSelectedItem() tea.Cmd {
	if m.selected < 0 || m.selected >= len(m.items) {
		return nil
	}

	// Navigate to manager to open the song
	// The manager screen handles opening songs
	return func() tea.Msg {
		return ScreenChangeMsg{Screen: 5} // screenManager
	}
}

// openProjectManager navigates to the project manager
func (m *RecentWorkModel) openProjectManager() tea.Cmd {
	return func() tea.Msg {
		return ScreenChangeMsg{Screen: 5} // screenManager
	}
}

// View renders the recent work panel
func (m *RecentWorkModel) View() string {
	if m.width == 0 {
		return "Recent Work"
	}

	t := theme.GetManager().Current()

	title := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Render("Recent Work")

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
			Width(m.width - 2).
			MaxWidth(m.width - 2).
			MaxHeight(m.height - 2).
			Padding(0, 1).
			Render(content)
	}

	// Show error state
	if m.err != nil {
		errorText := lipgloss.NewStyle().
			Foreground(t.Error).
			Render("[X] Could not load recent work")

		content := lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			errorText,
		)

		return lipgloss.NewStyle().
			Width(m.width - 2).
			MaxWidth(m.width - 2).
			MaxHeight(m.height - 2).
			Padding(0, 1).
			Render(content)
	}

	// Show empty state
	if len(m.items) == 0 {
		emptyText := lipgloss.NewStyle().
			Foreground(t.Secondary).
			Faint(true).
			Render("No recent work yet")

		content := lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			emptyText,
		)

		return lipgloss.NewStyle().
			Width(m.width - 2).
			MaxWidth(m.width - 2).
			MaxHeight(m.height - 2).
			Padding(0, 1).
			Render(content)
	}

	// Render recent work items
	var itemViews []string
	for i, item := range m.items {
		// Format time ago
		timeAgo := formatTimeAgo(item.UpdatedAt)

		// Determine style based on selection/hover state
		itemStyle := lipgloss.NewStyle().Foreground(t.Text)
		if i == m.selected {
			itemStyle = itemStyle.
				Foreground(t.Background).
				Background(t.Primary).
				Bold(true)
		} else if i == m.hovered {
			itemStyle = itemStyle.
				Foreground(t.Accent)
		}

		// Icon based on type
		icon := "[~]"
		if item.Type == "project" {
			icon = "[P]"
		}

		itemText := fmt.Sprintf("%s %s - %s", icon, item.Title, timeAgo)
		itemViews = append(itemViews, itemStyle.Render(itemText))
	}

	// Open projects button
	buttonStyle := lipgloss.NewStyle().
		Padding(0, 1).
		MarginTop(1)

	if m.hovered == len(m.items) {
		buttonStyle = buttonStyle.
			Foreground(t.Background).
			Background(t.Primary).
			Bold(true)
	} else {
		buttonStyle = buttonStyle.
			Foreground(t.Text).
			Background(t.Secondary)
	}

	openButton := buttonStyle.Render("[Open Projects]")

	// Mouse hint
	mouseHint := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Faint(true).
		Render("Click to open")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		lipgloss.JoinVertical(lipgloss.Left, itemViews...),
		"",
		openButton,
		mouseHint,
	)

	return lipgloss.NewStyle().
		Width(m.width - 2).
		MaxWidth(m.width - 2).
		MaxHeight(m.height - 2).
		Padding(0, 1).
		Render(content)
}

// formatTimeAgo formats a time as a human-readable "time ago" string
func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1m ago"
		}
		return fmt.Sprintf("%dm ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1h ago"
		}
		return fmt.Sprintf("%dh ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	default:
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1w ago"
		}
		return fmt.Sprintf("%dw ago", weeks)
	}
}
