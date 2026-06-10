package ui

import (
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/prism/internal/clipboard"
	"github.com/kyanite/prism/internal/color"
	"github.com/kyanite/prism/internal/palette"
	"github.com/kyanite/prism/internal/theme"
)

// WheelModel represents the color wheel screen
type WheelModel struct {
	themeManager *theme.Manager
	styles       Styles
	currentHue   float64
	saturation   float64
	lightness    float64
	status       string
}

// NewWheelModel creates a new color wheel model
func NewWheelModel(tm *theme.Manager) WheelModel {
	return WheelModel{
		themeManager: tm,
		styles:       NewStyles(tm.CurrentTheme()),
		currentHue:   0,
		saturation:   100,
		lightness:    50,
	}
}

// Init initializes the color wheel
func (m WheelModel) Init() tea.Cmd {
	return nil
}

// Update handles color wheel messages
func (m WheelModel) Update(msg tea.Msg) (WheelModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.currentHue = math.Mod(m.currentHue-5+360, 360)
			m.status = ""
		case "right", "l":
			m.currentHue = math.Mod(m.currentHue+5, 360)
			m.status = ""
		case "up", "k":
			m.lightness = math.Min(100, m.lightness+5)
			m.status = ""
		case "down", "j":
			m.lightness = math.Max(0, m.lightness-5)
			m.status = ""
		case "+":
			m.saturation = math.Min(100, m.saturation+10)
			m.status = ""
		case "-":
			m.saturation = math.Max(0, m.saturation-10)
			m.status = ""
		case "c":
			currentColor := color.NewFromHSL(m.currentHue, m.saturation, m.lightness)
			if err := clipboard.Copy(currentColor.Hex); err == nil {
				m.status = fmt.Sprintf("✓ Copied %s to clipboard", currentColor.Hex)
			} else {
				m.status = "✗ Clipboard unavailable - use export instead"
			}
		}
	}

	return m, nil
}

// View renders the color wheel
func (m WheelModel) View() string {
	styles := NewStyles(m.themeManager.CurrentTheme())

	var b strings.Builder

	// Title
	title := styles.Title.Render("Color Wheel")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Current color
	currentColor := color.NewFromHSL(m.currentHue, m.saturation, m.lightness)

	// Color visualization
	colorBar := m.renderColorBar()
	b.WriteString(colorBar)
	b.WriteString("\n\n")

	// Current hue indicator
	indicator := fmt.Sprintf("Current: %.0f°", m.currentHue)
	b.WriteString(styles.Accent.Render(indicator))
	b.WriteString("\n\n")

	// Current color info
	b.WriteString(styles.Primary.Render("Current Color:"))
	b.WriteString("\n")

	swatch := lipgloss.NewStyle().
		Background(lipgloss.Color(currentColor.Hex)).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(1, 4).
		Render("████")
	b.WriteString(swatch)
	fmt.Fprintf(&b, "  %s", currentColor.Hex)
	b.WriteString("\n")

	fmt.Fprintf(&b, "RGB: (%d, %d, %d)\n", currentColor.RGB.R, currentColor.RGB.G, currentColor.RGB.B)
	fmt.Fprintf(&b, "HSL: (%.0f°, %.0f%%, %.0f%%)\n", currentColor.HSL.H, currentColor.HSL.S, currentColor.HSL.L)
	b.WriteString("\n")

	// Related colors
	b.WriteString(styles.Secondary.Render("Related Colors:"))
	b.WriteString("\n")

	complement := currentColor.Complement()
	b.WriteString(m.renderColorLine("Complementary:", complement))

	analogous1 := color.NewFromHSL(math.Mod(m.currentHue-palette.AnalogousAngle+360, 360), m.saturation, m.lightness)
	b.WriteString(m.renderColorLine("Analogous -30°:", analogous1))

	analogous2 := color.NewFromHSL(math.Mod(m.currentHue+palette.AnalogousAngle, 360), m.saturation, m.lightness)
	b.WriteString(m.renderColorLine("Analogous +30°:", analogous2))

	b.WriteString("\n")

	// Status message
	if m.status != "" {
		b.WriteString(styles.Success.Render(m.status))
		b.WriteString("\n\n")
	}

	// Help
	help := styles.Muted.Render("←/→: Hue • ↑/↓: Lightness • +/-: Saturation • C: Copy • Esc: Menu")
	b.WriteString(help)

	// Wrap in border
	content := styles.Border.Width(ContentWidth).Render(b.String())

	return lipgloss.Place(ScreenWidth, ScreenHeight, lipgloss.Center, lipgloss.Center, content)
}

// renderColorBar renders the hue spectrum
func (m WheelModel) renderColorBar() string {
	var bar strings.Builder

	width := 60
	for i := 0; i < width; i++ {
		hue := float64(i) * 360.0 / float64(width)
		c := color.NewFromHSL(hue, 100, 50)

		style := lipgloss.NewStyle().
			Background(lipgloss.Color(c.Hex))

		// Highlight current hue
		char := "█"
		if math.Abs(hue-m.currentHue) < 10 {
			char = "▓"
		}

		bar.WriteString(style.Render(char))
	}

	return bar.String()
}

// renderColorLine renders a color swatch with info
func (m WheelModel) renderColorLine(label string, c color.Color) string {
	swatch := lipgloss.NewStyle().
		Background(lipgloss.Color(c.Hex)).
		Padding(0, 2).
		Render("██")

	return fmt.Sprintf("%s %s %s\n", swatch, c.Hex, label)
}
