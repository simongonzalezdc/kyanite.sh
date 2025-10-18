package styles

import "github.com/charmbracelet/lipgloss"

// Midnight Jazz Theme Colors
var (
	// Primary Colors
	Primary   = lipgloss.Color("#9D84B7") // Soft purple - main brand
	Secondary = lipgloss.Color("#5E4B8B") // Deep purple - secondary actions
	Accent    = lipgloss.Color("#F4D03F") // Gold - highlights, important elements

	// Functional Colors
	Success = lipgloss.Color("#52D3AA") // Mint green - success states
	Warning = lipgloss.Color("#FFA500") // Orange - warnings
	Error   = lipgloss.Color("#FF6347") // Tomato - errors
	Info    = lipgloss.Color("#87CEEB") // Sky blue - info messages

	// Background & Text Colors
	Background    = lipgloss.Color("#0A0E27") // Deep navy - main background
	TextPrimary   = lipgloss.Color("#E8DFF5") // Light lavender - main text
	TextSecondary = lipgloss.Color("#9D84B7") // Soft purple - secondary text
	TextMuted     = lipgloss.Color("#5E4B8B") // Deep purple - muted text
	TextAccent    = lipgloss.Color("#F4D03F") // Gold - emphasized text

	// Extended Palette
	Purple1 = lipgloss.Color("#9D84B7") // Lightest purple
	Purple2 = lipgloss.Color("#5E4B8B") // Medium purple
	Purple3 = lipgloss.Color("#3D2C8D") // Darkest purple

	Gold1 = lipgloss.Color("#F4D03F") // Bright gold
	Gold2 = lipgloss.Color("#D4AF37") // Muted gold

	Dark1 = lipgloss.Color("#0A0E27") // Main background
	Dark2 = lipgloss.Color("#1A1E37") // Lighter background
	Dark3 = lipgloss.Color("#2A2E47") // Border/divider

	BorderColor = lipgloss.Color("#2A2E47")
)

// Base Styles
var (
	// Title style for app header
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		MarginBottom(1).
		Padding(0, 2)

	// Subtitle for sections
	Subtitle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Italic(true).
			MarginBottom(1)

	// Main content text
	Text = lipgloss.NewStyle().
		Foreground(TextPrimary)

	// Muted text for hints/secondary info
	Muted = lipgloss.NewStyle().
		Foreground(TextMuted)

	// Emphasized text
	Emphasis = lipgloss.NewStyle().
			Foreground(TextAccent).
			Bold(true)
)

// Border Styles
var (
	// Standard border for panels
	Border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(1, 2)

	// Active/focused border
	BorderActive = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	// Thick border for emphasis
	BorderThick = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(Accent).
			Padding(1, 2)
)

// Button Styles
var (
	// Primary action button
	ButtonPrimary = lipgloss.NewStyle().
			Background(Primary).
			Foreground(Background).
			Bold(true).
			Padding(0, 3).
			MarginRight(1)

	// Secondary button
	ButtonSecondary = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Foreground(Secondary).
			Padding(0, 2).
			MarginRight(1)

	// Accent/highlight button
	ButtonAccent = lipgloss.NewStyle().
			Background(Accent).
			Foreground(Background).
			Bold(true).
			Padding(0, 3).
			MarginRight(1)

	// Disabled button
	ButtonDisabled = lipgloss.NewStyle().
			Background(Dark3).
			Foreground(TextMuted).
			Padding(0, 3).
			MarginRight(1)
)

// Status Styles
var (
	StatusSuccess = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	StatusWarning = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)

	StatusError = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	StatusInfo = lipgloss.NewStyle().
			Foreground(Info)
)

// Editor Component Styles
var (
	// Split pane editor - left side (editing)
	EditorPane = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	// Split pane preview - right side (rendered)
	PreviewPane = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(1, 2)

	// Status bar at bottom
	StatusBar = lipgloss.NewStyle().
			Background(Dark2).
			Foreground(TextPrimary).
			Padding(0, 2)

	// Cursor indicator
	Cursor = lipgloss.NewStyle().
		Background(Primary).
		Foreground(Background)

	// Divider between panes
	Divider = lipgloss.NewStyle().
		Foreground(BorderColor).
		Width(1).
		Height(1)
)

// Typography Styles
var (
	H1 = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		MarginBottom(1).
		Padding(0, 2)

	H2 = lipgloss.NewStyle().
		Bold(true).
		Foreground(Secondary).
		MarginBottom(1).
		Padding(0, 1).
		Underline(true)

	H3 = lipgloss.NewStyle().
		Bold(true).
		Foreground(TextAccent).
		MarginBottom(1)

	Bold = lipgloss.NewStyle().Bold(true)

	Italic = lipgloss.NewStyle().Italic(true)

	Underline = lipgloss.NewStyle().Underline(true)

	Code = lipgloss.NewStyle().
		Foreground(Accent).
		Background(Dark2).
		Padding(0, 1)

	Quote = lipgloss.NewStyle().
		Foreground(TextSecondary).
		Italic(true).
		BorderLeft(true).
		BorderForeground(Primary).
		PaddingLeft(2)
)

// List & Menu Styles
var (
	// Normal list item
	ListItem = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Padding(0, 2)

	// Selected list item
	ListItemSelected = lipgloss.NewStyle().
				Background(Primary).
				Foreground(Background).
				Bold(true).
				Padding(0, 2)
)

// Card & Panel Styles
var (
	// Standard card for grouped content
	Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(2, 4).
		MarginBottom(1)

	// Highlighted card (for important info)
	CardHighlight = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Accent).
			Padding(2, 4).
			MarginBottom(1)

	// Success card
	CardSuccess = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Success).
			Padding(2, 4).
			MarginBottom(1)

	// Error card
	CardError = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Error).
			Padding(2, 4).
			MarginBottom(1)
)

// Badge & Tag Styles
var (
	// Tag for metadata
	Tag = lipgloss.NewStyle().
		Background(Secondary).
		Foreground(TextPrimary).
		Padding(0, 1).
		MarginRight(1)

	// Badge for counts/numbers
	Badge = lipgloss.NewStyle().
		Background(Accent).
		Foreground(Background).
		Bold(true).
		Padding(0, 1)
)

// Helper Functions
func Gradient(text string, colors []lipgloss.Color) string {
	var result string
	textRunes := []rune(text)
	colorCount := len(colors)

	for i, char := range textRunes {
		if char == ' ' {
			result += " "
			continue
		}
		colorIdx := (i * colorCount) / len(textRunes)
		if colorIdx >= colorCount {
			colorIdx = colorCount - 1
		}
		style := lipgloss.NewStyle().Foreground(colors[colorIdx])
		result += style.Render(string(char))
	}

	return result
}

// Purple to gold gradient for titles
func TitleGradient(text string) string {
	colors := []lipgloss.Color{
		Primary,
		lipgloss.Color("#7D6497"),
		lipgloss.Color("#8D7447"),
		Accent,
	}
	return Gradient(text, colors)
}

// Status badge with color coding
func StatusBadge(status string, color lipgloss.Color) string {
	return lipgloss.NewStyle().
		Background(color).
		Foreground(Background).
		Bold(true).
		Padding(0, 1).
		Render(status)
}

// List item with icon
func ListItemWithIcon(icon string, text string, selected bool) string {
	iconStyle := lipgloss.NewStyle().Foreground(Accent)
	textStyle := ListItem
	if selected {
		textStyle = ListItemSelected
		iconStyle = iconStyle.Background(Primary)
	}
	return iconStyle.Render(icon) + " " + textStyle.Render(text)
}

// ToHex converts a lipgloss.Color to hex string for use in external configurations
func ToHex(color lipgloss.Color) string {
	return string(color)
}
