package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

 // Focus Color Palette - Synthwave Cyberpunk Theme
var (
	FocusPink   = lipgloss.Color("#FF71CE")
	FocusBlue   = lipgloss.Color("#00FFFF")
	FocusGreen  = lipgloss.Color("#00FF66")
	FocusPurple = lipgloss.Color("#B967C7")
	DarkBg      = lipgloss.Color("#0F0A19")
	AccentColor = lipgloss.Color("#FFC0CB")
	focusRed    = lipgloss.Color("#FF0055")

	// Base Styles
	focusStyle = lipgloss.NewStyle().
		Foreground(FocusPink).
		Bold(true)

	greenStyle = lipgloss.NewStyle().
		Foreground(FocusGreen).
		Bold(true)

	blueStyle = lipgloss.NewStyle().
		Foreground(FocusBlue)

	purpleStyle = lipgloss.NewStyle().
		Foreground(FocusPurple)
)

func FocusStyle(text string) string {
	return focusStyle.Render(text)
}

func FocusPinkColor(text string) string {
	return lipgloss.NewStyle().Foreground(FocusPink).Render(text)
}

func FocusBlueColor(text string) string {
	return lipgloss.NewStyle().Foreground(FocusBlue).Render(text)
}

func FocusGreenColor(text string) string {
	return lipgloss.NewStyle().Foreground(FocusGreen).Render(text)
}

func FocusPurpleColor(text string) string {
	return lipgloss.NewStyle().Foreground(FocusPurple).Render(text)
}

func PriorityStyle(priority string) string {
	var color lipgloss.Color
	var symbol string

	switch priority {
	case "high":
		color = focusRed
		symbol = "🔥"
	case "medium":
		color = FocusBlue
		symbol = "⚡"
	case "low":
		color = FocusGreen
		symbol = "💤"
	default:
		color = FocusPurple
		symbol = "⚪"
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(symbol + " " + priority)
}

func IDStyle(id string) string {
	return lipgloss.NewStyle().
		Foreground(FocusBlue).
		Render(id)
}

func HeaderStyle(text string) string {
	return lipgloss.NewStyle().
		Foreground(FocusPink).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Underline(true).
		Render(text)
}

func FooterStyle(text string) string {
	return lipgloss.NewStyle().
		Foreground(FocusPurple).
		Italic(true).
		Render(text)
}

func CategoryStyle(category string) string {
	return lipgloss.NewStyle().
		Foreground(FocusBlue).
		Background(DarkBg).
		Padding(0, 1).
		Bold(true).
		Render("🏷️ " + category)
}

func SuggestionStyle(suggestion string) string {
	return lipgloss.NewStyle().
		Foreground(FocusGreen).
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
		Foreground(FocusBlue).
		Italic(true).
		Render("🤖 " + response)
}

func ReportStyle(text string) string {
	return lipgloss.NewStyle().
		Foreground(FocusPurple).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(FocusBlue).
		Render(text)
}
