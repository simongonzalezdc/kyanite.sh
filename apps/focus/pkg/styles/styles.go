package styles

import (
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
	"github.com/pterm/pterm"
)

// themeMu protects currentThemeName and derived legacy/synthwave color vars.
var themeMu sync.RWMutex

// currentThemeName holds the active theme name for style derivation.
var currentThemeName = "amber-night"

func currentTheme() design.Theme {
	themeMu.RLock()
	name := currentThemeName
	themeMu.RUnlock()
	t := design.Get(name)
	if t.Name == "" {
		return design.DefaultTheme()
	}
	return t
}

// SetThemeByName changes the current theme by display name or id.
func SetThemeByName(name string) {
	themeMu.Lock()
	defer themeMu.Unlock()
	for _, id := range design.List() {
		t := design.Get(id)
		if t.Name == name || id == name {
			currentThemeName = id
			updateLegacyColors()
			updateSynthwaveColors()
			return
		}
	}
	currentThemeName = "amber-night"
	updateLegacyColors()
	updateSynthwaveColors()
}

// SetTheme is a no-op kept for backward compatibility with unified_dashboard.
// The old ThemeMode values are no longer used; everything goes through design themes.
func SetTheme(_ interface{}) {
	// All themes are now handled by the design module.
}

// GetTheme returns the current design.Theme.
func GetTheme() design.Theme {
	return currentTheme()
}

// CycleTheme cycles through all design module themes.
func CycleTheme() {
	themeMu.Lock()
	defer themeMu.Unlock()
	ids := design.List()
	for i, id := range ids {
		if id == currentThemeName {
			currentThemeName = ids[(i+1)%len(ids)]
			updateLegacyColors()
			updateSynthwaveColors()
			return
		}
	}
	currentThemeName = "amber-night"
	updateLegacyColors()
	updateSynthwaveColors()
}

// RefreshColors is kept for backward compatibility; now a no-op since styles
// are derived on each call from the current design.Theme.
func RefreshColors() {}

// Theme-aware color getters

func GetBackground() lipgloss.Color { return currentTheme().Background }
func GetForeground() lipgloss.Color { return currentTheme().Text }
func GetAccent() lipgloss.Color     { return currentTheme().Accent }
func GetBorder() lipgloss.Color     { return currentTheme().Border }
func GetSuccess() lipgloss.Color    { return currentTheme().Success }
func GetWarning() lipgloss.Color    { return currentTheme().Warning }
func GetError() lipgloss.Color      { return currentTheme().Error }
func GetPanel() lipgloss.Color      { return currentTheme().Panel }

// Legacy color aliases for backward compatibility.
// These map to the current theme's semantic fields.
var (
	FocusPink   lipgloss.Color // primary
	FocusBlue   lipgloss.Color // secondary
	FocusGreen  lipgloss.Color // success
	FocusPurple lipgloss.Color // accent
	DarkBg      lipgloss.Color // background
	AccentColor lipgloss.Color // border
	focusRed    lipgloss.Color // error
)

func init() {
	updateLegacyColors()
}

func updateLegacyColors() {
	t := design.Get(currentThemeName)
	if t.Name == "" {
		t = design.DefaultTheme()
	}
	FocusPink = t.Primary
	FocusBlue = t.Secondary
	FocusGreen = t.Success
	FocusPurple = t.Accent
	DarkBg = t.Background
	AccentColor = t.Border
	focusRed = t.Error
}

// Synthwave color aliases — map to current theme fields for backward compat.
var (
	SynthwavePink   lipgloss.Color // primary
	SynthwaveCyan   lipgloss.Color // accent
	SynthwavePurple lipgloss.Color // secondary
	SynthwaveYellow lipgloss.Color // warning
	SynthwaveGreen  lipgloss.Color // success
	SynthwaveOrange lipgloss.Color // warning
	SynthwaveRed    lipgloss.Color // error
	DeepSpace       lipgloss.Color // background
	DarkVoid        lipgloss.Color // panel
	CyberGrid       lipgloss.Color // panel
	FocusGrid       lipgloss.Color // border
)

func init() {
	updateSynthwaveColors()
}

func updateSynthwaveColors() {
	t := design.Get(currentThemeName)
	if t.Name == "" {
		t = design.DefaultTheme()
	}
	SynthwavePink = t.Primary
	SynthwaveCyan = t.Accent
	SynthwavePurple = t.Secondary
	SynthwaveYellow = t.Warning
	SynthwaveGreen = t.Success
	SynthwaveOrange = t.Warning
	SynthwaveRed = t.Error
	DeepSpace = t.Background
	DarkVoid = t.Panel
	CyberGrid = t.Panel
	FocusGrid = t.Border
}

// Theme-aware style getters

func GetTitleStyle() lipgloss.Style {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Primary).
		Background(t.Background).
		Bold(true).
		Padding(design.SpacingS).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent)
}

func GetBoxStyle() lipgloss.Style {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Panel).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Padding(design.SpacingS)
}

func GetPanelStyle() lipgloss.Style {
	return GetBoxStyle()
}

func GetSuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme().Success).
		Bold(true)
}

func GetErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme().Error).
		Bold(true)
}

func GetWarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme().Warning).
		Bold(true)
}

func GetSynthwaveTitle() lipgloss.Style { return GetTitleStyle() }
func GetFocusBox() lipgloss.Style       { return GetBoxStyle() }

// GetThemeList returns all available theme names for CLI commands.
func GetThemeList() []string { return design.List() }

// GetCurrentThemeName returns the display name of the current theme.
func GetCurrentThemeName() string { return currentTheme().Name }

// Base render functions

func FocusStyle(text string) string {
	return lipgloss.NewStyle().
		Foreground(currentTheme().Primary).
		Bold(true).
		Render(text)
}

func FocusPinkColor(text string) string {
	return lipgloss.NewStyle().Foreground(currentTheme().Primary).Render(text)
}

func FocusBlueColor(text string) string {
	return lipgloss.NewStyle().Foreground(currentTheme().Secondary).Render(text)
}

func FocusGreenColor(text string) string {
	return lipgloss.NewStyle().Foreground(currentTheme().Success).Render(text)
}

func FocusPurpleColor(text string) string {
	return lipgloss.NewStyle().Foreground(currentTheme().Accent).Render(text)
}

func PriorityStyle(priority string) string {
	t := currentTheme()
	var color lipgloss.Color
	var symbol string

	switch priority {
	case "high":
		color = t.Error
		symbol = "🔥"
	case "medium":
		color = t.Warning
		symbol = "⚡"
	case "low":
		color = t.Success
		symbol = "💤"
	default:
		color = t.Accent
		symbol = "⚪"
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(symbol + " " + priority)
}

func IDStyle(id string) string {
	return lipgloss.NewStyle().
		Foreground(currentTheme().Secondary).
		Render(id)
}

func HeaderStyle(text string) string {
	return lipgloss.NewStyle().
		Foreground(currentTheme().Primary).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Underline(true).
		Render(text)
}

func FooterStyle(text string) string {
	return lipgloss.NewStyle().
		Foreground(currentTheme().Accent).
		Italic(true).
		Render(text)
}

func CategoryStyle(category string) string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Secondary).
		Background(t.Background).
		Padding(0, design.SpacingXS).
		Bold(true).
		Render("🏷️ " + category)
}

func SuggestionStyle(suggestion string) string {
	return lipgloss.NewStyle().
		Foreground(currentTheme().Success).
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
		Foreground(currentTheme().Accent).
		Italic(true).
		Render("🤖 " + response)
}

func ReportStyle(text string) string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Accent).
		Padding(design.SpacingS).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Secondary).
		Render(text)
}

// Synthwave-style render functions (theme-aware replacements)

func SynthwaveTitle(text string) string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Primary).
		Background(t.Background).
		Bold(true).
		Italic(true).
		Underline(true).
		Padding(design.SpacingS, design.SpacingM).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent).
		AlignHorizontal(lipgloss.Center).
		Render(text)
}

func Title(text string) string { return SynthwaveTitle(text) }

func FocusBox(text string, borderColor lipgloss.Color) string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Accent).
		Background(t.Panel).
		Padding(design.SpacingS, design.SpacingM).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Bold(true).
		Render(text)
}

func CyberGridBox(text string) string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Success).
		Background(t.Panel).
		Padding(0, design.SpacingXS).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent).
		Render("▶ " + text)
}

func HolographicText(text string) string {
	t := currentTheme()
	colors := []lipgloss.Color{
		t.Primary, t.Accent, t.Secondary,
		t.Warning, t.Success, t.Warning,
	}

	result := ""
	for i, char := range text {
		color := colors[i%len(colors)]
		result += lipgloss.NewStyle().
			Foreground(color).
			Bold(true).
			Background(t.Background).
			Render(string(char))
	}
	return result
}

func DigitalRain(text string) string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Success).
		Background(t.Background).
		Bold(true).
		Italic(true).
		Render("▓ " + text + " ▓")
}

func PriorityExplosion(priority string) string {
	t := currentTheme()
	switch priority {
	case "high":
		return lipgloss.NewStyle().
			Foreground(t.Error).
			Background(t.Panel).
			Bold(true).
			Render("🔥 HIGH PRIORITY 🔥")
	case "medium":
		return lipgloss.NewStyle().
			Foreground(t.Warning).
			Background(t.Panel).
			Bold(true).
			Render("⚡ MEDIUM ⚡")
	case "low":
		return lipgloss.NewStyle().
			Foreground(t.Success).
			Background(t.Background).
			Bold(true).
			Render("💤 LOW PRIORITY 💤")
	default:
		return lipgloss.NewStyle().
			Foreground(t.Secondary).
			Render("◉ NORMAL")
	}
}

func TaskStatus(status string) string {
	t := currentTheme()
	if status == "completed" {
		return lipgloss.NewStyle().
			Foreground(t.Success).
			Background(t.Panel).
			Bold(true).
			Render("✅ MISSION COMPLETE ✅")
	}
	return lipgloss.NewStyle().
		Foreground(t.Accent).
		Background(t.Panel).
		Bold(true).
		Render("◯ ACTIVE MISSION ◯")
}

func CyberTag(tag string) string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Primary).
		Background(t.Panel).
		Padding(0, design.SpacingXS).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Secondary).
		Render("🏷️ " + tag)
}

func Header() string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Success).
		Background(t.Background).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Padding(design.SpacingS).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent).
		Render("focus.sh Task Management")
}

func CyberStats(active, completed, total int) string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Accent).
		Background(t.Panel).
		Bold(true).
		Padding(design.SpacingS).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		Render(
			"📊 GRID STATUS: "+
				lipgloss.NewStyle().Foreground(t.Warning).Render("⚡ "+string(rune(active)))+
				" ACTIVE | "+
				lipgloss.NewStyle().Foreground(t.Success).Render("✅ "+string(rune(completed)))+
				" COMPLETED | "+
				lipgloss.NewStyle().Foreground(t.Secondary).Render("🌟 "+string(rune(total)))+
				" TOTAL",
		)
}

func LoadingMessage() string {
	return lipgloss.NewStyle().
		Foreground(currentTheme().Accent).
		Bold(true).
		Render("⠋ Loading...")
}

func EmptyStateMessage() string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Accent).
		Background(t.Background).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Padding(design.SpacingM).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		Render("No tasks found.\n\nCreate your first task:\nfocus add \"task description\"")
}

// Pterm integration for CLI output — these remain as-is.

func SynthwaveGetPriorityColor(priority string) *pterm.Style {
	switch priority {
	case "high":
		return &pterm.Style{
			pterm.FgRed,
			pterm.Bold,
			pterm.BgLightBlue,
		}
	case "medium":
		return &pterm.Style{
			pterm.FgYellow,
			pterm.Bold,
			pterm.BgMagenta,
		}
	case "low":
		return &pterm.Style{
			pterm.FgGreen,
			pterm.Bold,
			pterm.BgCyan,
		}
	default:
		return &pterm.Style{
			pterm.FgCyan,
			pterm.Bold,
			pterm.BgBlue,
		}
	}
}

func SynthwaveAIResponseStyle(response string) string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Accent).
		Background(t.Panel).
		Bold(true).
		Italic(true).
		Padding(design.SpacingS).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Secondary).
		Render("🤖 AI SYNTHESIS: " + response)
}

func SynthwaveReportStyle(text string) string {
	t := currentTheme()
	return lipgloss.NewStyle().
		Foreground(t.Primary).
		Background(t.Panel).
		Padding(design.SpacingM).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent).
		Bold(true).
		Render("📋 ANALYSIS REPORT:\n\n" + text)
}

// ThemeMode and constants kept for backward compat with unified_dashboard.
type ThemeMode string

const (
	ThemeSynthwave ThemeMode = "synthwave"
	ThemeLight     ThemeMode = "light"
	ThemePlain     ThemeMode = "plain"
)
