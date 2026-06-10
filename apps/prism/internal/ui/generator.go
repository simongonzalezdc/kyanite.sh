package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/prism/internal/clipboard"
	"github.com/kyanite/prism/internal/color"
	"github.com/kyanite/prism/internal/export"
	"github.com/kyanite/prism/internal/palette"
	"github.com/kyanite/prism/internal/storage"
	"github.com/kyanite/prism/internal/theme"
)

// GeneratorModel represents the palette generator screen
type GeneratorModel struct {
	themeManager     *theme.Manager
	styles           Styles
	baseColorInput   string
	selectedRule     int
	rules            []palette.HarmonyRule
	generatedPalette *palette.Palette
	err              string
	status           string
	exportMode       bool
	selectedExport   int
	exportFormats    []string
}

// NewGeneratorModel creates a new generator model
func NewGeneratorModel(tm *theme.Manager) GeneratorModel {
	return GeneratorModel{
		themeManager:   tm,
		styles:         NewStyles(tm.CurrentTheme()),
		baseColorInput: "#FF0080",
		selectedRule:   0,
		rules:          palette.AllRules(),
		exportFormats:  []string{"JSON", "CSS Variables", "TOML", "Kyanite Theme"},
	}
}

// Init initializes the generator
func (m GeneratorModel) Init() tea.Cmd {
	return nil
}

// Update handles generator messages
func (m GeneratorModel) Update(msg tea.Msg) (GeneratorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.exportMode {
			return m.handleExportMode(msg.String())
		}

		switch msg.String() {
		case "up", "k":
			if m.selectedRule > 0 {
				m.selectedRule--
			}
			m.status = ""
		case "down", "j":
			if m.selectedRule < len(m.rules)-1 {
				m.selectedRule++
			}
			m.status = ""
		case "enter", " ":
			return m, m.generate()
		case "c":
			if m.generatedPalette != nil {
				m.copyPalette()
			}
		case "s":
			if m.generatedPalette != nil {
				m.savePalette()
			}
		case "e":
			if m.generatedPalette != nil {
				m.exportMode = true
				m.selectedExport = 0
				m.status = ""
			}
		}
	}

	return m, nil
}

func (m *GeneratorModel) handleExportMode(key string) (GeneratorModel, tea.Cmd) {
	switch key {
	case "up", "k":
		if m.selectedExport > 0 {
			m.selectedExport--
		}
	case "down", "j":
		if m.selectedExport < len(m.exportFormats)-1 {
			m.selectedExport++
		}
	case "enter", " ":
		m.exportPalette()
		m.exportMode = false
	case "esc":
		m.exportMode = false
		m.status = ""
	}
	return *m, nil
}

func (m *GeneratorModel) copyPalette() {
	hexes := make([]string, len(m.generatedPalette.Colors))
	for i, c := range m.generatedPalette.Colors {
		hexes[i] = c.Hex
	}
	text := strings.Join(hexes, ", ")

	if err := clipboard.Copy(text); err == nil {
		m.status = fmt.Sprintf("✓ Copied %d colors to clipboard", len(hexes))
	} else {
		m.status = "✗ Clipboard unavailable"
	}
	m.err = ""
}

func (m *GeneratorModel) savePalette() {
	if err := storage.SavePalette(*m.generatedPalette); err == nil {
		m.status = fmt.Sprintf("✓ Saved palette: %s", m.generatedPalette.Name)
	} else {
		m.status = fmt.Sprintf("✗ Failed to save: %v", err)
	}
	m.err = ""
}

func (m *GeneratorModel) exportPalette() {
	var output string
	var err error

	switch m.selectedExport {
	case 0: // JSON
		data, exportErr := export.ExportJSON(*m.generatedPalette)
		if exportErr != nil {
			err = exportErr
		} else {
			output = string(data)
		}
	case 1: // CSS
		output = export.ExportCSS(*m.generatedPalette)
	case 2: // TOML
		output = export.ExportTOML(*m.generatedPalette)
	case 3: // Kyanite
		data, exportErr := export.ExportTheme(*m.generatedPalette)
		if exportErr != nil {
			err = exportErr
		} else {
			output = string(data)
		}
	}

	if err == nil && clipboard.Copy(output) == nil {
		m.status = fmt.Sprintf("✓ Exported as %s (copied to clipboard)", m.exportFormats[m.selectedExport])
	} else {
		m.status = "✗ Export failed or clipboard unavailable"
	}
	m.err = ""
}

// View renders the generator
func (m GeneratorModel) View() string {
	styles := NewStyles(m.themeManager.CurrentTheme())

	var b strings.Builder

	// Title
	title := styles.Title.Render("Palette Generator")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Base color
	b.WriteString(styles.Primary.Render("Base Color: "))
	b.WriteString(m.baseColorInput)
	b.WriteString("\n\n")

	// Harmony rules
	b.WriteString(styles.Secondary.Render("Select Harmony Rule:"))
	b.WriteString("\n")

	for i, rule := range m.rules {
		style := styles.Unselected
		cursor := "  "
		if i == m.selectedRule {
			style = styles.Selected
			cursor = "▸ "
		}

		line := fmt.Sprintf("%s%s", cursor, rule)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Generated palette
	if m.generatedPalette != nil {
		b.WriteString(styles.Success.Render("Generated Palette:"))
		b.WriteString("\n")

		for i, c := range m.generatedPalette.Colors {
			swatch := lipgloss.NewStyle().
				Background(lipgloss.Color(c.Hex)).
				Padding(0, 2).
				Render("██")

			line := fmt.Sprintf("%d. %s %s", i+1, swatch, c.Hex)
			b.WriteString(line)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Export mode menu
	if m.exportMode {
		b.WriteString(styles.Secondary.Render("Select Export Format:"))
		b.WriteString("\n")
		for i, format := range m.exportFormats {
			style := styles.Unselected
			cursor := "  "
			if i == m.selectedExport {
				style = styles.Selected
				cursor = "▸ "
			}
			line := fmt.Sprintf("%s%s", cursor, format)
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}
		b.WriteString("\n")
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
	if m.exportMode {
		help = styles.Muted.Render("↑/↓: Select Format • Enter: Export • Esc: Cancel")
	} else if m.generatedPalette != nil {
		help = styles.Muted.Render("↑/↓: Select • Enter: Generate • C: Copy • S: Save • E: Export • Esc: Menu")
	} else {
		help = styles.Muted.Render("↑/↓: Select Rule • Enter: Generate • Esc: Menu")
	}
	b.WriteString(help)

	// Wrap in border
	content := styles.Border.Width(ContentWidth).Render(b.String())

	return lipgloss.Place(ScreenWidth, ScreenHeight, lipgloss.Center, lipgloss.Center, content)
}

// generate generates a palette
func (m GeneratorModel) generate() tea.Cmd {
	return func() tea.Msg {
		baseColor, err := color.ParseHex(m.baseColorInput)
		if err != nil {
			m.err = "Invalid hex color"
			return nil
		}

		pal, err := palette.Generate(baseColor, m.rules[m.selectedRule])
		if err != nil {
			m.err = err.Error()
			return nil
		}

		m.generatedPalette = &pal
		m.err = ""
		return nil
	}
}
