package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ai "github.com/kyanite/ai"
	"github.com/kyanite/design"
	aipanel "github.com/kyanite/tui/aipanel"
	"github.com/kyanite/prism/internal/clipboard"
	"github.com/kyanite/prism/internal/color"
	"github.com/kyanite/prism/internal/export"
	"github.com/kyanite/prism/internal/palette"
	"github.com/kyanite/prism/internal/storage"
	"github.com/kyanite/prism/internal/theme"
)

// MoodPalette holds a single named palette from a mood generation.
type MoodPalette struct {
	Name   string   `json:"name"`
	Colors []string `json:"colors"`
}

// A11yAnalysis holds the result of an accessibility analysis.
type A11yAnalysis struct {
	Raw string
}


// AIMode represents the current AI panel sub-mode.
type AIMode int

const (
	AIModeNone AIMode = iota
	AIModePalette
	AIModeMood
	AIModeA11y
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
	brain            *ai.Brain
	aiInputMode      bool
	aiInputText      string
	aiMode           AIMode
	aiPanel          aipanel.Model
	aiPanelVisible   bool
	exportFormats    []string
	crossAppContext string // Cross-app context from the Brain, loaded when AI panel opens
}

// NewGeneratorModel creates a new generator model bound to the given
// AI brain. The brain may be nil if the shared kyanite/ai Brain could
// not be initialized (e.g. NUCBox unreachable at startup); the AI
// panel and AI generation features simply no-op or return errors.
func NewGeneratorModel(tm *theme.Manager, brain *ai.Brain) GeneratorModel {
	var panel aipanel.Model
	if brain != nil {
		panel = aipanel.New(brain, ContentWidth, ScreenHeight)
	}
	return GeneratorModel{
		themeManager:   tm,
		styles:         NewStyles(tm.CurrentTheme()),
		baseColorInput: "#FF0080",
		selectedRule:   0,
		rules:          palette.AllRules(),
		exportFormats:  []string{"JSON", "CSS Variables", "TOML", "Kyanite Theme"},
		brain:          brain,
		aiPanel:        panel,
	}
}

// Init initializes the generator
func (m GeneratorModel) Init() tea.Cmd {
	return nil
}

// Update handles generator messages
func (m GeneratorModel) Update(msg tea.Msg) (GeneratorModel, tea.Cmd) {
	// Route aipanel messages first when panel is visible
	switch msg := msg.(type) {
	case aipanel.StreamChunk:
		var cmd tea.Cmd
		m.aiPanel, cmd = m.aiPanel.Update(msg)
		return m, cmd
	case aipanel.ErrorMsg:
		var cmd tea.Cmd
		m.aiPanel, cmd = m.aiPanel.Update(msg)
		m.aiPanelVisible = m.aiPanel.Visible()
		return m, cmd
	}

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
	case AIMoodGeneratedMsg:
		if msg.Err != "" {
			m.err = msg.Err
			m.status = ""
		} else {
			m.generatedPalette = msg.Palette
			m.err = ""
			m.status = "✓ AI mood palette generated"
		}
		m.aiInputMode = false
		return m, nil
	case AIA11yGeneratedMsg:
		if msg.Err != "" {
			m.err = msg.Err
			m.status = ""
		} else {
			m.err = ""
			m.status = "✓ A11y analysis complete"
			m.aiPanel = m.aiPanel.StartStream("A11y Analysis")
			m.aiPanel = m.aiPanel.Toggle() // make visible
			// Write the analysis text directly into the panel content
			m.aiPanel.Content() // ensure initialized
		}
		m.aiInputMode = false
		return m, nil
	case crossAppContextLoadedMsg:
		m.crossAppContext = msg.context
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
			// Enter AI palette input mode (backward compatible)
			m.aiInputMode = true
			m.aiInputText = ""
			m.aiMode = AIModePalette
			m.aiPanelVisible = false
			m.status = ""
			m.err = ""
		case "ctrl+a":
			// Toggle AI panel
			if m.aiPanelVisible {
				m.aiPanelVisible = false
				m.aiPanel = m.aiPanel.Toggle()
				m.status = ""
			} else {
				m.aiInputMode = true
				m.aiInputText = ""
				m.aiMode = AIModePalette
				m.aiPanelVisible = true
				m.aiPanel = m.aiPanel.Toggle()
				m.status = ""
				m.err = ""
				// Load cross-app context when panel opens
				if m.brain != nil {
					return m, m.loadCrossAppContext()
				}
			}
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

// resolveAIMode determines the AI sub-mode from the input text prefix.
func resolveAIMode(input string) (AIMode, string) {
	if strings.HasPrefix(input, "mood:") {
		return AIModeMood, strings.TrimSpace(strings.TrimPrefix(input, "mood:"))
	}
	if strings.HasPrefix(input, "a11y:") {
		return AIModeA11y, strings.TrimSpace(strings.TrimPrefix(input, "a11y:"))
	}
	return AIModePalette, input
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

	// AI panel (when visible via Ctrl+A)
	if m.aiPanelVisible {
		b.WriteString(m.aiPanel.View())
		b.WriteString("\n")
	}

	// AI input prompt
	if m.aiInputMode {
		modeLabel := "AI Palette"
		switch m.aiMode {
		case AIModeMood:
			modeLabel = "AI Mood"
		case AIModeA11y:
			modeLabel = "AI A11y"
		}

		hint := "describe your palette"
		switch m.aiMode {
		case AIModeMood:
			hint = "describe the mood"
		case AIModeA11y:
			hint = "paste colors to check"
		default:
			hint = "describe your palette (prefix mood: or a11y: to switch mode)"
		}

		b.WriteString(styles.Secondary.Render(fmt.Sprintf("◆ %s — %s:", modeLabel, hint)))
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
		help = styles.Muted.Render("Type description, Enter: Generate • Esc: Cancel • Prefix mood: or a11y: to switch mode")
	} else if m.exportMode {
		help = styles.Muted.Render("↑/↓: Select Format • Enter: Export • Esc: Cancel")
	} else if m.generatedPalette != nil {
		help = styles.Muted.Render("↑/↓: Select • Enter: Generate • P: AI Palette • Ctrl+A: AI Panel • C: Copy • S: Save • E: Export • Esc: Menu")
	} else {
		help = styles.Muted.Render("↑/↓: Select Rule • Enter: Generate • P: AI Palette • Ctrl+A: AI Panel • Esc: Menu")
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

// AIMoodGeneratedMsg carries the result of an AI mood palette generation.
type AIMoodGeneratedMsg struct {
	Palette *palette.Palette
	Err     string
}

// AIA11yGeneratedMsg carries the result of an AI accessibility analysis.
type AIA11yGeneratedMsg struct {
	Analysis *A11yAnalysis
	Err      string
}

// handleAIInput handles key events while in AI description input mode.
func (m GeneratorModel) handleAIInput(key string) (GeneratorModel, tea.Cmd) {
	switch key {
	case "enter":
		if strings.TrimSpace(m.aiInputText) == "" {
			m.aiInputMode = false
			return m, nil
		}
		// Resolve mode from prefix (allows overriding the default mode)
		resolvedMode, input := resolveAIMode(m.aiInputText)
		m.aiMode = resolvedMode

		switch m.aiMode {
		case AIModeMood:
			m.status = "Generating mood palettes…"
			return m, m.generateMoodWithAI(input)
		case AIModeA11y:
			m.status = "Analyzing accessibility…"
			return m, m.generateA11yWithAI(input)
		default:
			m.status = "Generating AI palette…"
			return m, m.generateWithAI(input)
		}
	case "esc":
		m.aiInputMode = false
		m.aiPanelVisible = false
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

// hexRegex matches hex color codes like #FF5733
var hexRegex = regexp.MustCompile(`#[0-9A-Fa-f]{6}`)

// paletteResult is the JSON shape returned by PrismPalettePrompt.
type paletteResult struct {
	Colors []string `json:"colors"`
	Name   string   `json:"name"`
}

// generateWithAI sends a description to the AI brain and returns a palette.
func (m GeneratorModel) generateWithAI(description string) tea.Cmd {
	return func() tea.Msg {
		if m.brain == nil {
			return AIPaletteGeneratedMsg{Err: "AI not available"}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if !m.brain.IsLLMAvailable(ctx) {
			return AIPaletteGeneratedMsg{Err: "ai: NUCBox Ollama server unreachable"}
		}

		prompt := ai.PrismPalettePrompt(description + formatCrossAppContext(m.crossAppContext))
		resp, err := m.brain.Generate(ctx, prompt, ai.WithJSONMode())
		if err != nil {
			return AIPaletteGeneratedMsg{Err: fmt.Sprintf("ai: palette generation failed: %v", err)}
		}

		cleaned := stripMarkdownFence(resp)
		var result paletteResult
		if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
			return AIPaletteGeneratedMsg{Err: fmt.Sprintf("ai: failed to parse palette response: %v", err)}
		}
		if len(result.Colors) == 0 {
			return AIPaletteGeneratedMsg{Err: "ai: no colors returned"}
		}

		parsedColors := make([]color.Color, 0, len(result.Colors))
		for _, hex := range result.Colors {
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

// moodPaletteResult is the JSON shape returned by PrismMoodPalettePrompt.
type moodPaletteResult struct {
	Palettes []MoodPalette `json:"palettes"`
}

// generateMoodWithAI sends a mood description to the AI brain and
// returns the first palette.
func (m GeneratorModel) generateMoodWithAI(mood string) tea.Cmd {
	return func() tea.Msg {
		if m.brain == nil {
			return AIMoodGeneratedMsg{Err: "AI not available"}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()

		if !m.brain.IsLLMAvailable(ctx) {
			return AIMoodGeneratedMsg{Err: "ai: NUCBox Ollama server unreachable"}
		}

		prompt := ai.PrismMoodPalettePrompt(mood + formatCrossAppContext(m.crossAppContext))
		resp, err := m.brain.Generate(ctx, prompt, ai.WithJSONMode())
		if err != nil {
			return AIMoodGeneratedMsg{Err: fmt.Sprintf("ai: mood palette generation failed: %v", err)}
		}

		cleaned := stripMarkdownFence(resp)
		var result moodPaletteResult
		if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
			return AIMoodGeneratedMsg{Err: fmt.Sprintf("ai: failed to parse mood palette response: %v", err)}
		}
		if len(result.Palettes) == 0 {
			return AIMoodGeneratedMsg{Err: "ai: no palettes returned"}
		}

		// Use the first palette from the mood results
		first := result.Palettes[0]
		parsedColors := make([]color.Color, 0, len(first.Colors))
		for _, hex := range first.Colors {
			c, parseErr := color.ParseHex(hex)
			if parseErr != nil {
				continue
			}
			parsedColors = append(parsedColors, c)
		}

		if len(parsedColors) == 0 {
			return AIMoodGeneratedMsg{Err: "AI returned no valid colors for mood"}
		}

		baseColor := parsedColors[0]
		pal := palette.Palette{
			Name:        "Mood: " + first.Name,
			Colors:      parsedColors,
			HarmonyRule: "ai-mood",
			BaseColor:   baseColor,
		}

		return AIMoodGeneratedMsg{Palette: &pal}
	}
}

// generateA11yWithAI sends colors for accessibility analysis via the AI brain.
func (m GeneratorModel) generateA11yWithAI(colorsInput string) tea.Cmd {
	return func() tea.Msg {
		if m.brain == nil {
			return AIA11yGeneratedMsg{Err: "AI not available"}
		}

		// If no explicit colors provided, use the current palette
		colors := colorsInput
		if m.generatedPalette != nil && strings.TrimSpace(colors) == "" {
			hexes := make([]string, len(m.generatedPalette.Colors))
			for i, c := range m.generatedPalette.Colors {
				hexes[i] = c.Hex
			}
			colors = strings.Join(hexes, ", ")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if !m.brain.IsLLMAvailable(ctx) {
			return AIA11yGeneratedMsg{Err: "ai: NUCBox Ollama server unreachable"}
		}

		prompt := ai.PrismA11yCheckPrompt(colors)
		resp, err := m.brain.Generate(ctx, prompt)
		if err != nil {
			return AIA11yGeneratedMsg{Err: fmt.Sprintf("ai: a11y analysis failed: %v", err)}
		}

		return AIA11yGeneratedMsg{Analysis: &A11yAnalysis{Raw: strings.TrimSpace(resp)}}
	}
}
// crossAppContextLoadedMsg carries cross-app context loaded from the Brain.
type crossAppContextLoadedMsg struct {
	context string
}

// loadCrossAppContext loads context from other kyanite apps via the Brain.
// Best-effort: returns an empty-string message if unavailable.
func (m GeneratorModel) loadCrossAppContext() tea.Cmd {
	return func() tea.Msg {
		if m.brain == nil {
			return crossAppContextLoadedMsg{context: ""}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		entries, err := m.brain.GetCrossAppContext(ctx, 3)
		if err != nil || len(entries) == 0 {
			return crossAppContextLoadedMsg{context: ""}
		}

		var parts []string
		for _, e := range entries {
			switch e.SourceApp {
			case "syntax":
				parts = append(parts, fmt.Sprintf("Story themes: %s", e.Summary))
			case "noise":
				parts = append(parts, fmt.Sprintf("Musical mood: %s", e.Summary))
			}
		}
		return crossAppContextLoadedMsg{context: strings.Join(parts, "\n")}
	}
}

// formatCrossAppContext formats cross-app context for inclusion in AI prompts.
// Returns an empty string if no context is available.
func formatCrossAppContext(ctx string) string {
	if ctx == "" {
		return ""
	}
	return "\n\nInspiration from your other apps:\n" + ctx
}

// stripMarkdownFence removes a leading or trailing ```json / ``` fence that
// the LLM may wrap around JSON output, then trims surrounding whitespace.
func stripMarkdownFence(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
