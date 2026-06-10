package ui

import (
	"strings"

	"github.com/kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// =============================================================================
// Types
// =============================================================================

// ShortcutCategory groups related shortcuts
type ShortcutCategory struct {
	Name      string
	Shortcuts []ShortcutItem
}

// ShortcutItem represents a single shortcut
type ShortcutItem struct {
	Key         string
	Description string
}

// ToggleShortcutsOverlayMsg toggles the shortcuts overlay visibility
type ToggleShortcutsOverlayMsg struct{}

// =============================================================================
// Shortcuts Overlay Model
// =============================================================================

// ShortcutsOverlayModel displays a keyboard shortcuts quick reference
type ShortcutsOverlayModel struct {
	visible    bool
	width      int
	height     int
	scrollY    int
	categories []ShortcutCategory

	// Styles
	overlayStyle     lipgloss.Style
	titleStyle       lipgloss.Style
	categoryStyle    lipgloss.Style
	keyStyle         lipgloss.Style
	descriptionStyle lipgloss.Style
	hintStyle        lipgloss.Style
}

// NewShortcutsOverlayModel creates a new shortcuts overlay
func NewShortcutsOverlayModel() *ShortcutsOverlayModel {
	t := theme.GetManager().Current()

	m := &ShortcutsOverlayModel{
		visible:    false,
		width:      80,
		height:     24,
		scrollY:    0,
		categories: defaultShortcuts(),

		overlayStyle: lipgloss.NewStyle().
			Background(t.Background).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary).
			Padding(1, 2),

		titleStyle: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			MarginBottom(1),

		categoryStyle: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true).
			MarginTop(1),

		keyStyle: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.Secondary).
			Padding(0, 1),

		descriptionStyle: lipgloss.NewStyle().
			Foreground(t.Text),

		hintStyle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Italic(true).
			MarginTop(1),
	}

	return m
}

// defaultShortcuts returns the default shortcut categories
func defaultShortcuts() []ShortcutCategory {
	return []ShortcutCategory{
		{
			Name: "Navigation",
			Shortcuts: []ShortcutItem{
				{Key: "Tab", Description: "Next pane"},
				{Key: "Shift+Tab", Description: "Previous pane"},
				{Key: "Ctrl+J/K", Description: "Navigate panes"},
				{Key: "Home/End", Description: "Start/end of line"},
				{Key: "Ctrl+Home/End", Description: "Start/end of file"},
				{Key: "PgUp/PgDn", Description: "Page up/down"},
			},
		},
		{
			Name: "Editing",
			Shortcuts: []ShortcutItem{
				{Key: "Ctrl+C", Description: "Copy"},
				{Key: "Ctrl+V", Description: "Paste"},
				{Key: "Ctrl+X", Description: "Cut"},
				{Key: "Ctrl+Z", Description: "Undo"},
				{Key: "Ctrl+Y", Description: "Redo"},
				{Key: "Ctrl+A", Description: "Select all"},
			},
		},
		{
			Name: "File",
			Shortcuts: []ShortcutItem{
				{Key: "Ctrl+S", Description: "Save"},
				{Key: "Ctrl+Shift+S", Description: "Save as"},
				{Key: "Ctrl+O", Description: "Open"},
				{Key: "Ctrl+N", Description: "New"},
			},
		},
		{
			Name: "Tools",
			Shortcuts: []ShortcutItem{
				{Key: "Ctrl+F", Description: "Chord picker"},
				{Key: "Ctrl+Shift+B", Description: "BPM tapper"},
				{Key: "Ctrl+Shift+F", Description: "Find"},
				{Key: "Ctrl+H", Description: "Replace"},
			},
		},
		{
			Name: "View",
			Shortcuts: []ShortcutItem{
				{Key: "Ctrl++/-", Description: "Zoom in/out"},
				{Key: "Ctrl+0", Description: "Reset zoom"},
				{Key: "F1", Description: "Help"},
				{Key: "?", Description: "Toggle shortcuts"},
				{Key: "Esc", Description: "Close overlay"},
			},
		},
		{
			Name: "Theme",
			Shortcuts: []ShortcutItem{
				{Key: "Ctrl+T", Description: "Cycle themes"},
				{Key: "Ctrl+Shift+T", Description: "Theme picker"},
			},
		},
	}
}

// Init initializes the model
func (m *ShortcutsOverlayModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *ShortcutsOverlayModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ToggleShortcutsOverlayMsg:
		m.visible = !m.visible
		m.scrollY = 0

	case tea.KeyMsg:
		if !m.visible {
			return nil
		}

		switch msg.String() {
		case "esc", "?":
			m.visible = false
		case "up", "k":
			if m.scrollY > 0 {
				m.scrollY--
			}
		case "down", "j":
			m.scrollY++
		case "home":
			m.scrollY = 0
		}
	}

	return nil
}

// View renders the overlay
func (m *ShortcutsOverlayModel) View() string {
	if !m.visible {
		return ""
	}

	// Build content
	var content strings.Builder

	// Title
	content.WriteString(m.titleStyle.Render("⌨ Keyboard Shortcuts"))
	content.WriteString("\n")

	// Categories
	for _, cat := range m.categories {
		content.WriteString(m.categoryStyle.Render(cat.Name))
		content.WriteString("\n")

		for _, shortcut := range cat.Shortcuts {
			key := m.keyStyle.Render(shortcut.Key)
			desc := m.descriptionStyle.Render(shortcut.Description)
			content.WriteString("  " + key + "  " + desc + "\n")
		}
	}

	// Hint
	content.WriteString(m.hintStyle.Render("Press Esc or ? to close"))

	// Calculate overlay dimensions
	overlayWidth := min(60, m.width-4)
	overlayHeight := min(m.height-4, 30)

	// Apply scroll
	lines := strings.Split(content.String(), "\n")
	maxScroll := max(0, len(lines)-overlayHeight+4)
	if m.scrollY > maxScroll {
		m.scrollY = maxScroll
	}

	// Get visible lines
	start := m.scrollY
	end := min(start+overlayHeight-4, len(lines))
	visibleLines := lines[start:end]

	// Render overlay
	overlayContent := strings.Join(visibleLines, "\n")
	overlay := m.overlayStyle.
		Width(overlayWidth).
		MaxHeight(overlayHeight).
		Render(overlayContent)

	// Center the overlay
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		overlay,
	)
}

// =============================================================================
// Public Methods
// =============================================================================

// IsVisible returns whether the overlay is visible
func (m *ShortcutsOverlayModel) IsVisible() bool {
	return m.visible
}

// Show displays the overlay
func (m *ShortcutsOverlayModel) Show() {
	m.visible = true
	m.scrollY = 0
}

// Hide hides the overlay
func (m *ShortcutsOverlayModel) Hide() {
	m.visible = false
}

// Toggle toggles the overlay visibility
func (m *ShortcutsOverlayModel) Toggle() {
	m.visible = !m.visible
	if m.visible {
		m.scrollY = 0
	}
}

// SetDimensions sets the overlay dimensions
func (m *ShortcutsOverlayModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// UpdateTheme refreshes styles when theme changes
func (m *ShortcutsOverlayModel) UpdateTheme() {
	t := theme.GetManager().Current()

	m.overlayStyle = lipgloss.NewStyle().
		Background(t.Background).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		Padding(1, 2)

	m.titleStyle = lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		MarginBottom(1)

	m.categoryStyle = lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true).
		MarginTop(1)

	m.keyStyle = lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Secondary).
		Padding(0, 1)

	m.descriptionStyle = lipgloss.NewStyle().
		Foreground(t.Text)

	m.hintStyle = lipgloss.NewStyle().
		Foreground(t.Secondary).
		Italic(true).
		MarginTop(1)
}
