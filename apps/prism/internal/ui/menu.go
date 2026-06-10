package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/prism/internal/theme"
)

// MenuModel represents the main menu
type MenuModel struct {
	themeManager *theme.Manager
	styles       Styles
	selected     int
	options      []menuOption
}

type menuOption struct {
	key    string
	label  string
	screen int
}

// NewMenuModel creates a new menu model
func NewMenuModel(tm *theme.Manager) MenuModel {
	options := []menuOption{
		{"c", "Color Wheel", 1},
		{"g", "Generate Palette", 2},
		{"s", "Search Colors", 6},
		{"l", "Learn Color Theory", 3},
		{"a", "Check Accessibility (WCAG)", 4},
		{"m", "Manage Palettes", 5},
		{"q", "Quit", -1},
	}

	return MenuModel{
		themeManager: tm,
		styles:       NewStyles(tm.CurrentTheme()),
		options:      options,
		selected:     0,
	}
}

// Init initializes the menu
func (m MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles menu messages
func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.options)-1 {
				m.selected++
			}
		case "enter", " ":
			return m, m.navigate()
		case "c":
			return m, func() tea.Msg { return Navigate(1) }
		case "g":
			return m, func() tea.Msg { return Navigate(2) }
		case "s":
			return m, func() tea.Msg { return Navigate(6) }
		case "l":
			return m, func() tea.Msg { return Navigate(3) }
		case "a":
			return m, func() tea.Msg { return Navigate(4) }
		case "m":
			return m, func() tea.Msg { return Navigate(5) }
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the menu
func (m MenuModel) View() string {
	styles := NewStyles(m.themeManager.CurrentTheme())

	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.themeManager.CurrentTheme().Primary)).
		Bold(true).
		Padding(1).
		Align(lipgloss.Center).
		Width(60).
		Render("╔═══════════════════════════════════════╗\n║         PRISM.SH v1.0.0              ║\n║    Color Palette Design Tool          ║\n╚═══════════════════════════════════════╝")

	b.WriteString(title)
	b.WriteString("\n\n")

	// Menu options
	for i, opt := range m.options {
		style := styles.Unselected
		cursor := "  "
		if i == m.selected {
			style = styles.Selected
			cursor = "▸ "
		}

		line := fmt.Sprintf("%s[%s] %s", cursor, opt.key, opt.label)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Help text
	help := styles.Muted.Render("↑/↓: Navigate • Enter: Select • Ctrl+H: Help • Ctrl+Q: Quit")
	b.WriteString(help)

	// Wrap in border
	content := styles.Border.Width(ContentWidthNarrow).Render(b.String())

	return lipgloss.Place(ScreenWidth, ScreenHeight, lipgloss.Center, lipgloss.Center, content)
}

// navigate returns a command to navigate to the selected screen
func (m MenuModel) navigate() tea.Cmd {
	opt := m.options[m.selected]
	if opt.screen == -1 {
		return tea.Quit
	}
	return func() tea.Msg {
		return Navigate(opt.screen)
	}
}
