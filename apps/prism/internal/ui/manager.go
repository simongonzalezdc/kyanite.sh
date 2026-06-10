package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/prism/internal/palette"
	"github.com/kyanite/prism/internal/storage"
	"github.com/kyanite/prism/internal/theme"
)

// ManagerModel represents the palette manager screen
type ManagerModel struct {
	themeManager *theme.Manager
	styles       Styles
	palettes     []palette.Palette
	selected     int
	err          string
}

// NewManagerModel creates a new manager model
func NewManagerModel(tm *theme.Manager) ManagerModel {
	return ManagerModel{
		themeManager: tm,
		styles:       NewStyles(tm.CurrentTheme()),
		selected:     0,
	}
}

// Init initializes the manager
func (m ManagerModel) Init() tea.Cmd {
	return m.loadPalettes()
}

// Update handles manager messages
func (m ManagerModel) Update(msg tea.Msg) (ManagerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.palettes)-1 {
				m.selected++
			}
		case "d":
			if len(m.palettes) > 0 {
				return m, m.deletePalette()
			}
		case "r":
			return m, m.loadPalettes()
		}
	}

	return m, nil
}

// View renders the manager
func (m ManagerModel) View() string {
	styles := NewStyles(m.themeManager.CurrentTheme())

	var b strings.Builder

	// Title
	title := styles.Title.Render("Palette Manager")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Palettes list
	if len(m.palettes) == 0 {
		b.WriteString(styles.Muted.Render("No saved palettes"))
		b.WriteString("\n\n")
	} else {
		b.WriteString(styles.Secondary.Render(fmt.Sprintf("Saved Palettes (%d):", len(m.palettes))))
		b.WriteString("\n")

		for i, pal := range m.palettes {
			style := styles.Unselected
			cursor := "  "
			if i == m.selected {
				style = styles.Selected
				cursor = "▸ "
			}

			line := fmt.Sprintf("%s%s (%s - %d colors)", cursor, pal.Name, pal.HarmonyRule, len(pal.Colors))
			b.WriteString(style.Render(line))
			b.WriteString("\n")

			// Show color swatches for selected
			if i == m.selected {
				var swatches strings.Builder
				swatches.WriteString("  ")
				for _, c := range pal.Colors {
					swatch := lipgloss.NewStyle().
						Background(lipgloss.Color(c.Hex)).
						Padding(0, 1).
						Render("██")
					swatches.WriteString(swatch)
					swatches.WriteString(" ")
				}
				b.WriteString(swatches.String())
				b.WriteString("\n")
			}
		}
		b.WriteString("\n")
	}

	// Error
	if m.err != "" {
		b.WriteString(styles.Error.Render("Error: " + m.err))
		b.WriteString("\n\n")
	}

	// Help
	help := styles.Muted.Render("↑/↓: Navigate • D: Delete • R: Refresh • Esc: Menu")
	b.WriteString(help)

	// Wrap in border
	content := styles.Border.Width(ContentWidth).Render(b.String())

	return lipgloss.Place(ScreenWidth, ScreenHeight, lipgloss.Center, lipgloss.Center, content)
}

// loadPalettes loads saved palettes
func (m ManagerModel) loadPalettes() tea.Cmd {
	return func() tea.Msg {
		palettes, err := storage.ListPalettes()
		if err != nil {
			m.err = err.Error()
			return nil
		}

		m.palettes = palettes
		m.err = ""
		return nil
	}
}

// deletePalette deletes the selected palette
func (m ManagerModel) deletePalette() tea.Cmd {
	return func() tea.Msg {
		if len(m.palettes) == 0 {
			return nil
		}

		pal := m.palettes[m.selected]
		err := storage.DeletePalette(pal.ID)
		if err != nil {
			m.err = err.Error()
			return nil
		}

		// Reload palettes
		return m.loadPalettes()
	}
}
