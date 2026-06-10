package dashboard

import (
	"github.com/kyanite/design/icons"
	"github.com/kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MusicTool represents a music tool item
type MusicTool struct {
	ID          string
	Title       string
	Icon        string
	Description string
}

// MusicToolsModel manages the music tools panel
type MusicToolsModel struct {
	width    int
	height   int
	tools    []MusicTool
	selected int
	hovered  int
}

// NewMusicToolsModel creates a new music tools model
func NewMusicToolsModel() *MusicToolsModel {
	tools := []MusicTool{
		{ID: "chords", Title: "Chord Progressions", Icon: "[C]", Description: "Common chord progressions"},
		{ID: "scales", Title: "Scale Reference", Icon: icons.GetIcon("settings"), Description: "Musical scales and modes"},
		{ID: "rhymes", Title: "Rhyme Dictionary", Icon: "[R]", Description: "Find rhyming words"},
		{ID: "templates", Title: "Song Templates", Icon: "[T]", Description: "Song structure templates"},
	}

	return &MusicToolsModel{
		tools:    tools,
		selected: 0,
		hovered:  -1,
	}
}

// Update handles messages for the music tools panel
func (m *MusicToolsModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
			if m.selected < len(m.tools) {
				m.selected++
			}
		case "enter":
			return m.activateSelected()
		case "ctrl+shift+c":
			m.selected = 0
			return m.openTheoryTools()
		case "ctrl+shift+s":
			m.selected = 1
			return m.openTheoryTools()
		case "ctrl+shift+r":
			m.selected = 2
			return m.openTheoryTools()
		case "ctrl+shift+t":
			m.selected = 3
			return m.openTheoryTools()
		}
	}

	return nil
}

// handleMouse processes mouse events
func (m *MusicToolsModel) handleMouse(msg tea.MouseMsg) tea.Cmd {
	if m.width == 0 {
		return nil
	}

	// Calculate item positions
	// Title (1) + empty (1) + items start at row 2
	itemStartY := 2
	itemEndY := itemStartY + len(m.tools) - 1
	buttonY := itemEndY + 2

	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionRelease {
			// Check if click is on a tool item
			if msg.Y >= itemStartY && msg.Y <= itemEndY {
				m.selected = msg.Y - itemStartY
				return m.openTheoryTools()
			}
			// Check if click is on the button
			if msg.Y >= buttonY && msg.Y <= buttonY+1 {
				return m.openTheoryTools()
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
		if m.selected < len(m.tools) {
			m.selected++
		}

	case tea.MouseButtonNone:
		if msg.Action == tea.MouseActionMotion {
			if msg.Y >= itemStartY && msg.Y <= itemEndY {
				m.hovered = msg.Y - itemStartY
			} else if msg.Y >= buttonY && msg.Y <= buttonY+1 {
				m.hovered = len(m.tools) // Special index for button
			} else {
				m.hovered = -1
			}
		}
	}

	return nil
}

// activateSelected activates the selected item
func (m *MusicToolsModel) activateSelected() tea.Cmd {
	if m.selected == len(m.tools) {
		// Button is selected
		return m.openTheoryTools()
	}
	if m.selected >= 0 && m.selected < len(m.tools) {
		return m.openTheoryTools()
	}
	return nil
}

// openTheoryTools navigates to the theory tools screen
func (m *MusicToolsModel) openTheoryTools() tea.Cmd {
	return func() tea.Msg {
		return ScreenChangeMsg{Screen: 4} // screenTheory
	}
}

// View renders the music tools panel
func (m *MusicToolsModel) View() string {
	if m.width == 0 {
		return "Music Tools"
	}

	t := theme.GetManager().Current()

	title := lipgloss.Style{}.
		Foreground(t.Primary).
		Bold(true).
		Render("Music Tools")

	// Render tool items
	var itemViews []string
	for i, tool := range m.tools {
		itemStyle := lipgloss.Style{}.Foreground(t.Text)

		if i == m.selected {
			itemStyle = itemStyle.
				Foreground(t.Background).
				Background(t.Primary).
				Bold(true)
		} else if i == m.hovered {
			itemStyle = itemStyle.
				Foreground(t.Accent)
		}

		itemText := tool.Icon + " " + tool.Title
		itemViews = append(itemViews, itemStyle.Render(itemText))
	}

	// Open theory tools button
	buttonStyle := lipgloss.Style{}.
		Padding(0, 1).
		MarginTop(1)

	if m.hovered == len(m.tools) || m.selected == len(m.tools) {
		buttonStyle = buttonStyle.
			Foreground(t.Background).
			Background(t.Primary).
			Bold(true)
	} else {
		buttonStyle = buttonStyle.
			Foreground(t.Text).
			Background(t.Secondary)
	}

	openButton := buttonStyle.Render("[Open Theory Tools]")

	// Mouse hint
	mouseHint := lipgloss.Style{}.
		Foreground(t.Secondary).
		Faint(true).
		Render("Click or use Ctrl+Shift+C/S/R/T")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		lipgloss.JoinVertical(lipgloss.Left, itemViews...),
		"",
		openButton,
		mouseHint,
	)

	return lipgloss.Style{}.
		Width(m.width-2).
		MaxWidth(m.width-2).
		MaxHeight(m.height-2).
		Padding(0, 1).
		Render(content)
}
