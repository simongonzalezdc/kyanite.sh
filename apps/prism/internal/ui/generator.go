package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
	appai "github.com/kyanite/prism/internal/ai"
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
	aiClient    *appai.Client
	aiInputMode bool
	aiInputText string
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
		aiClient:       appai.NewClient(),
	}
}

// Init initializes the generator
func (m GeneratorModel) Init() tea.Cmd {
	return nil
}

// Update handles generator messages
func (m GeneratorModel) Update(msg tea.Msg) (GeneratorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case PaletteGeneratedMsg:
		if msg.Err != "" {
			m.err = msg.Err
			m.generatedPalette = nil
		} else {
			m.generatedPalette = msg.Palette
			m.err = ""
		}
		return m, nil
	case AIPaletteGeneratedMsg:
		if msg.Err != "" {
			m.err = msg.Err
			m.status = ""
		} else {
			m.generatedPalette = msg.Palette
			m.err = ""
			m.status = "✓ AI palette generated"
		}
		m.aiInputMode = false
		return m, nil
	case tea.KeyMsg:
		// AI input mode captures text until Enter/Escape
		if m.aiInputMode {
			return m.handleAIInput(msg.String())
		}

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
		case "p":
			// Enter AI palette input mode
			m.aiInputMode = true
			m.aiInputText = ""
			m.status = ""
			m.err = ""
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
			swatch := lipgloss.Style{}.
				Background(lipgloss.Color(c.Hex)).
				Padding(design.SpacingNone, design.SpacingS).
				Render("██")

			line := fmt.Sprintf("%d. %s %s", i+1, swatch, c.Hex)
			b.WriteString(line)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// AI input prompt
	if m.aiInputMode {
		b.WriteString(styles.Secondary.Render("AI Palette — describe your palette:"))
		b.WriteString("\n")
		prompt := fmt.Sprintf("> %s█", m.aiInputText)
		b.WriteString(styles.Primary.Render(prompt))
		b.WriteString("\n\n")
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
	if m.aiInputMode {
		help = styles.Muted.Render("Type description, Enter: Generate • Esc: Cancel")
	} else if m.exportMode {
		help = styles.Muted.Render("↑/↓: Select Format • Enter: Export • Esc: Cancel")
	} else if m.generatedPalette != nil {
		help = styles.Muted.Render("↑/↓: Select • Enter: Generate • P: AI Palette • C: Copy • S: Save • E: Export • Esc: Menu")
	} else {
		help = styles.Muted.Render("↑/↓: Select Rule • Enter: Generate • P: AI Palette • Esc: Menu")
	}
	b.WriteString(help)

	// Wrap in border
	content := styles.Border.Width(ContentWidth).Render(b.String())

	return lipgloss.Place(ScreenWidth, ScreenHeight, lipgloss.Center, lipgloss.Center, content)
}

// PaletteGeneratedMsg carries the result of a palette generation back to Update.
type PaletteGeneratedMsg struct {
	Palette *palette.Palette
	Err     string
}

// generate generates a palette
func (m GeneratorModel) generate() tea.Cmd {
	return func() tea.Msg {
		baseColor, err := color.ParseHex(m.baseColorInput)
		if err != nil {
			return PaletteGeneratedMsg{Err: "Invalid hex color"}
		}

		pal, err := palette.Generate(baseColor, m.rules[m.selectedRule])
		if err != nil {
			return PaletteGeneratedMsg{Err: err.Error()}
		}

		return PaletteGeneratedMsg{Palette: &pal}
	}
}

// AIPaletteGeneratedMsg carries the result of an AI palette generation.
type AIPaletteGeneratedMsg struct {
	Palette *palette.Palette
	Err     string
}

// handleAIInput handles key events while in AI description input mode.
func (m GeneratorModel) handleAIInput(key string) (GeneratorModel, tea.Cmd) {
	switch key {
	case "enter":
		if strings.TrimSpace(m.aiInputText) == "" {
			m.aiInputMode = false
			return m, nil
		}
		m.status = "Generating AI palette…"
		return m, m.generateWithAI(m.aiInputText)
	case "esc":
		m.aiInputMode = false
		m.aiInputText = ""
		m.status = ""
	case "backspace":
		if len(m.aiInputText) > 0 {
			m.aiInputText = m.aiInputText[:len(m.aiInputText)-1]
		}
	default:
		// Only accept printable characters (single rune keys from bubbletea)
		if len(key) == 1 {
			m.aiInputText += key
		}
	}
	return m, nil
}

// generateWithAI sends a description to the AI client and returns a palette.
func (m GeneratorModel) generateWithAI(description string) tea.Cmd {
	return func() tea.Msg {
		if m.aiClient == nil {
			return AIPaletteGeneratedMsg{Err: "AI not available"}
		}

		colors, err := m.aiClient.GeneratePalette(context.Background(), description)
		if err != nil {
			return AIPaletteGeneratedMsg{Err: err.Error()}
		}

		parsedColors := make([]color.Color, 0, len(colors))
		for _, hex := range colors {
			c, parseErr := color.ParseHex(hex)
			if parseErr != nil {
				continue
			}
			parsedColors = append(parsedColors, c)
		}

		if len(parsedColors) == 0 {
			return AIPaletteGeneratedMsg{Err: "AI returned no valid colors"}
		}

		baseColor := parsedColors[0]
		pal := palette.Palette{
			Name:        "AI: " + description,
			Colors:      parsedColors,
			HarmonyRule: "ai-generated",
			BaseColor:   baseColor,
		}

		return AIPaletteGeneratedMsg{Palette: &pal}
	}
}
