package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

// NEON Color Palette - Synthwave Cyberpunk Theme
var (
	NeonPink    = lipgloss.Color("#FF71CE")
	NeonBlue    = lipgloss.Color("#00FFFF")
	NeonGreen   = lipgloss.Color("#00FF66")
	NeonPurple  = lipgloss.Color("#B967C7")
	DarkBg      = lipgloss.Color("#0F0A19")
	AccentColor = lipgloss.Color("#FFC0CB")
	neonRed     = lipgloss.Color("#FF0055")

	// Base Styles
	neonStyle = lipgloss.NewStyle().
		Foreground(NeonPink).
		Bold(true)

	greenStyle = lipgloss.NewStyle().
		Foreground(NeonGreen).
		Bold(true)

	blueStyle = lipgloss.NewStyle().
		Foreground(NeonBlue)

	purpleStyle = lipgloss.NewStyle().
		Foreground(NeonPurple)
)

func NeonStyle(text string) string {
	return neonStyle.Render(text)
}

func NeonPinkColor(text string) string {
	return lipgloss.NewStyle().Foreground(NeonPink).Render(text)
}

func NeonBlueColor(text string) string {
	return lipgloss.NewStyle().Foreground(NeonBlue).Render(text)
}

func NeonGreenColor(text string) string {
	return lipgloss.NewStyle().Foreground(NeonGreen).Render(text)
}

func NeonPurpleColor(text string) string {
	return lipgloss.NewStyle().Foreground(NeonPurple).Render(text)
}

func PriorityStyle(priority string) string {
	var color lipgloss.Color
	var symbol string

	switch priority {
	case "high":
		color = neonRed
		symbol = "🔥"
	case "medium":
		color = NeonBlue
		symbol = "⚡"
	case "low":
		color = NeonGreen
		symbol = "💤"
	default:
		color = NeonPurple
		symbol = "⚪"
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(symbol + " " + priority)
}

func IDStyle(id string) string {
	return lipgloss.NewStyle().
		Foreground(NeonBlue).
		Render(id)
}

func HeaderStyle(text string) string {
	return lipgloss.NewStyle().
		Foreground(NeonPink).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Underline(true).
		Render(text)
}

func FooterStyle(text string) string {
	return lipgloss.NewStyle().
		Foreground(NeonPurple).
		Italic(true).
		Render(text)
}

func CategoryStyle(category string) string {
	return lipgloss.NewStyle().
		Foreground(NeonBlue).
		Background(DarkBg).
		Padding(0, 1).
		Bold(true).
		Render("🏷️ " + category)
}

func SuggestionStyle(suggestion string) string {
	return lipgloss.NewStyle().
		Foreground(NeonGreen).
		Bold(true).
		Render("🎯 " + suggestion)
}

func GetPriorityColor(priority string) *pterm.Style {
	switch priority {
	case "high":
		return &pterm.Style{pterm.FgRed}
	case "medium":
		return &pterm.Style{pterm.FgCyan}
	case "low":
		return &pterm.Style{pterm.FgGreen}
	default:
		return &pterm.Style{pterm.FgWhite}
	}
}

func AIResponseStyle(response string) string {
	return lipgloss.NewStyle().
		Foreground(NeonBlue).
		Italic(true).
		Render("🤖 " + response)
}

func ReportStyle(text string) string {
	return lipgloss.NewStyle().
		Foreground(NeonPurple).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(NeonBlue).
		Render(text)
}
