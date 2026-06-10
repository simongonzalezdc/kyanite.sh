package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/prism/internal/clipboard"
	"github.com/kyanite/prism/internal/color"
	"github.com/kyanite/prism/internal/theme"
)

// SearchModel represents the color name search screen
type SearchModel struct {
	themeManager   *theme.Manager
	styles         Styles
	textInput      textinput.Model
	results        []color.NamedColor
	selectedResult int
	status         string
	err            string
}

// NewSearchModel creates a new search model
func NewSearchModel(tm *theme.Manager) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Enter color name (e.g., red, blue, coral)..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 50

	return SearchModel{
		themeManager: tm,
		styles:       NewStyles(tm.CurrentTheme()),
		textInput:    ti,
		results:      []color.NamedColor{},
	}
}

// Init initializes the search screen
func (m SearchModel) Init() tea.Cmd {
	return nil
}

// Update handles search messages
func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedResult > 0 {
				m.selectedResult--
			}
			m.status = ""
		case "down", "j":
			if m.selectedResult < len(m.results)-1 {
				m.selectedResult++
			}
			m.status = ""
		case "enter":
			return m, m.search()
		case "c":
			if len(m.results) > 0 && m.selectedResult < len(m.results) {
				m.copyColor()
			}
		default:
			// Handle text input
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the search screen
func (m SearchModel) View() string {
	styles := NewStyles(m.themeManager.CurrentTheme())

	var b strings.Builder

	// Title
	title := styles.Title.Render("Color Name Search")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Search input
	b.WriteString(styles.Primary.Render("Search:"))
	b.WriteString("\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	// Results
	if len(m.results) > 0 {
		b.WriteString(styles.Secondary.Render(fmt.Sprintf("Results (%d):", len(m.results))))
		b.WriteString("\n")

		// Show up to 10 results
		maxResults := len(m.results)
		if maxResults > 10 {
			maxResults = 10
		}

		for i := 0; i < maxResults; i++ {
			result := m.results[i]
			style := styles.Unselected
			cursor := "  "
			if i == m.selectedResult {
				style = styles.Selected
				cursor = "▸ "
			}

			swatch := lipgloss.NewStyle().
				Background(lipgloss.Color(result.Hex)).
				Padding(0, 2).
				Render("██")

			line := fmt.Sprintf("%s%s %s - %s", cursor, swatch, result.Name, result.Hex)
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}

		if len(m.results) > 10 {
			b.WriteString(styles.Muted.Render(fmt.Sprintf("... and %d more", len(m.results)-10)))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	} else if m.textInput.Value() != "" {
		b.WriteString(styles.Muted.Render("No colors found. Try a different search term."))
		b.WriteString("\n\n")
	}

	// Error
	if m.err != "" {
		b.WriteString(styles.Error.Render("Error: " + m.err))
		b.WriteString("\n\n")
	}

	// Status
	if m.status != "" {
		b.WriteString(styles.Success.Render(m.status))
		b.WriteString("\n\n")
	}

	// Help
	var help string
	if len(m.results) > 0 {
		help = styles.Muted.Render("↑/↓: Navigate • C: Copy • Enter: Search • Esc: Menu")
	} else {
		help = styles.Muted.Render("Type to search • Enter: Search • Esc: Menu")
	}
	b.WriteString(help)

	// Wrap in border
	content := styles.Border.Width(ContentWidth).Render(b.String())

	return lipgloss.Place(ScreenWidth, ScreenHeight, lipgloss.Center, lipgloss.Center, content)
}

// search performs the color search
func (m SearchModel) search() tea.Cmd {
	return func() tea.Msg {
		query := m.textInput.Value()
		if query == "" {
			m.results = []color.NamedColor{}
			m.selectedResult = 0
			return nil
		}

		results := color.SearchColors(query)
		m.results = results
		m.selectedResult = 0
		m.err = ""

		if len(results) == 0 {
			m.status = ""
		} else {
			m.status = ""
		}

		return nil
	}
}

// copyColor copies the selected color to clipboard
func (m *SearchModel) copyColor() {
	if m.selectedResult >= len(m.results) {
		return
	}

	selected := m.results[m.selectedResult]

	if err := clipboard.Copy(selected.Hex); err == nil {
		m.status = fmt.Sprintf("✓ Copied %s (%s) to clipboard", selected.Name, selected.Hex)
	} else {
		m.status = "✗ Clipboard unavailable"
	}
	m.err = ""
}
