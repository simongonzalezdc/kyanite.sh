package dashboard

import (
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ThemePreviewModel showcases the current theme
type ThemePreviewModel struct {
	width       int
	height      int
	currentIdx  int
	themes      []theme.Theme
	showDetails bool
	// Add field for theme change notifications
	themeChangeChan chan theme.Theme
}

// NewThemePreviewModel creates a new theme preview
func NewThemePreviewModel() *ThemePreviewModel {
	themeIDs := theme.ListThemes()
	themes := make([]theme.Theme, len(themeIDs))
	for i, id := range themeIDs {
		themes[i] = theme.GetTheme(id)
	}
	
	return &ThemePreviewModel{
		currentIdx: 0,
		themes:     themes,
		showDetails: false,
		themeChangeChan: make(chan theme.Theme, 10), // Buffered channel for theme changes
	}
}

// Update handles messages for the theme preview
func (m *ThemePreviewModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// CRITICAL: Return a command to trigger re-render
		return nil

	// Handle theme change notifications
	case theme.ThemeChangeMsg:
		// Find the index of the new theme
		for i, t := range m.themes {
			if t.Name == msg.Theme.Name {
				m.currentIdx = i
				break
			}
		}
		return nil
	}

	// CRITICAL: Always return commands
	return tea.Batch(cmds...)
}

// View renders the theme preview
func (m *ThemePreviewModel) View() string {
	if m.width == 0 {
		return "Theme Preview"
	}
	
	currentTheme := m.themes[m.currentIdx]
	
	// Create theme showcase
	title := lipgloss.NewStyle().
		Foreground(currentTheme.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Render(currentTheme.Name)
	
	// Color palette display
	palette := m.renderColorPalette(currentTheme)
	
	// Sample UI elements
	samples := m.renderSampleElements(currentTheme)
	
	// Navigation
	nav := m.renderNavigation()
	
	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		palette,
		"",
		samples,
		"",
		nav,
	)
}

func (m *ThemePreviewModel) renderColorPalette(t theme.Theme) string {
	colors := []struct {
		name  string
		color lipgloss.Color
	}{
		{"Primary", t.Primary},
		{"Secondary", t.Secondary},
		{"Accent", t.Accent},
		{"Background", t.Background},
		{"Text", t.Text},
		{"Success", t.Success},
		{"Warning", t.Warning},
		{"Error", t.Error},
	}
	
	var rows []string
	for i := 0; i < len(colors); i += 4 {
		var row []string
		for j := i; j < i+4 && j < len(colors); j++ {
			colorBox := lipgloss.NewStyle().
				Background(colors[j].color).
				Foreground(t.Text).
				Width(8).
				Height(3).
				Align(lipgloss.Center).
				Bold(true).
				Render(colors[j].name[:3])
			row = append(row, colorBox)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, row...))
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m *ThemePreviewModel) renderSampleElements(t theme.Theme) string {
	button := lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Primary).
		Padding(0, 2).
		Bold(true).
		Render("Primary Button")
	
	secondary := lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Secondary).
		Padding(0, 2).
		Render("Secondary")
	
	accent := lipgloss.NewStyle().
		Foreground(t.Background).
		Background(t.Accent).
		Padding(0, 2).
		Bold(true).
		Render("Accent")
	
	return lipgloss.JoinHorizontal(lipgloss.Left, button, " ", secondary, " ", accent)
}

func (m *ThemePreviewModel) renderNavigation() string {
	nav := lipgloss.NewStyle().
		Foreground(theme.GetManager().Current().Text).
		Render("◀ ▶ Switch Themes")
	
	return nav
}

// NextTheme switches to the next theme
func (m *ThemePreviewModel) NextTheme() {
	m.currentIdx = (m.currentIdx + 1) % len(m.themes)
	// Update the global theme manager
	theme.GetManager().SetTheme(theme.ListThemes()[m.currentIdx])
}

// PreviousTheme switches to the previous theme
func (m *ThemePreviewModel) PreviousTheme() {
	m.currentIdx = (m.currentIdx - 1 + len(m.themes)) % len(m.themes)
	// Update the global theme manager
	theme.GetManager().SetTheme(theme.ListThemes()[m.currentIdx])
}

// SetThemeChangeChannel sets up a channel to listen for theme changes
func (m *ThemePreviewModel) SetThemeChangeChannel(ch <-chan theme.Theme) {
	go func() {
		for newTheme := range ch {
			// Find the index of the new theme
			for i, t := range m.themes {
				if t.Name == newTheme.Name {
					m.currentIdx = i
					break
				}
			}
		}
	}()
}