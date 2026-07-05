package styles

import (
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
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
			rebuildStyles(currentThemeName)
			return
		}
	}
	currentThemeName = "amber-night"
	updateLegacyColors()
	updateSynthwaveColors()
	rebuildStyles(currentThemeName)
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
			rebuildStyles(currentThemeName)
			return
		}
	}
	currentThemeName = "amber-night"
	updateLegacyColors()
	updateSynthwaveColors()
	rebuildStyles(currentThemeName)
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

// GetThemeList returns all available theme names for CLI commands.
func GetThemeList() []string { return design.List() }

// GetCurrentThemeName returns the display name of the current theme.
func GetCurrentThemeName() string { return currentTheme().Name }

// Theme-aware style getters. All derive from the cached, theme-driven Styles
// in tokens.go; app code never constructs styles inline.

func GetTitleStyle() lipgloss.Style   { return Current().Title }
func GetBoxStyle() lipgloss.Style     { return Current().Box }
func GetPanelStyle() lipgloss.Style   { return Current().Box }
func GetSuccessStyle() lipgloss.Style { return Current().Success }
func GetErrorStyle() lipgloss.Style   { return Current().Error }
func GetWarningStyle() lipgloss.Style { return Current().Warning }

func GetSynthwaveTitle() lipgloss.Style { return Current().SynthwaveTitle }
func GetFocusBox() lipgloss.Style       { return Current().Box }

// Base render functions

func FocusStyle(text string) string       { return Current().FocusStyle.Render(text) }
func FocusPinkColor(text string) string   { return Current().FocusPink.Render(text) }
func FocusBlueColor(text string) string   { return Current().FocusBlue.Render(text) }
func FocusGreenColor(text string) string  { return Current().FocusGreen.Render(text) }
func FocusPurpleColor(text string) string { return Current().FocusPurple.Render(text) }

func PriorityStyle(priority string) string {
	t := Current().Theme
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

	return lipgloss.Style{}.
		Foreground(color).
		Bold(true).
		Render(symbol + " " + priority)
}

func IDStyle(id string) string { return Current().ID.Render(id) }

func HeaderStyle(text string) string { return Current().HeaderStyle.Render(text) }
func FooterStyle(text string) string { return Current().FooterStyle.Render(text) }

func CategoryStyle(category string) string {
	return Current().Category.Render("🏷️ " + category)
}

func SuggestionStyle(suggestion string) string {
	return Current().Suggestion.Render("🎯 " + suggestion)
}

func AIResponseStyle(response string) string {
	return Current().AIResponse.Render("🤖 " + response)
}

func ReportStyle(text string) string { return Current().Report.Render(text) }

// Synthwave-style render functions (theme-aware, drawn from the cached Styles).

func SynthwaveTitle(text string) string { return Current().SynthwaveTitle.Render(text) }

func Title(text string) string { return SynthwaveTitle(text) }

func FocusBox(text string, borderColor lipgloss.Color) string {
	return Current().FocusBoxBase.BorderForeground(borderColor).Render(text)
}

func CyberGridBox(text string) string {
	return Current().CyberGrid.Render("▶ " + text)
}

func HolographicText(text string) string {
	t := Current().Theme
	colors := []lipgloss.Color{
		t.Primary, t.Accent, t.Secondary,
		t.Warning, t.Success, t.Warning,
	}

	result := ""
	for i, char := range text {
		color := colors[i%len(colors)]
		result += lipgloss.Style{}.
			Foreground(color).
			Bold(true).
			Background(t.Background).
			Render(string(char))
	}
	return result
}

func DigitalRain(text string) string {
	return Current().DigitalRain.Render("▓ " + text + " ▓")
}

func PriorityExplosion(priority string) string {
	t := Current().Theme
	switch priority {
	case "high":
		return lipgloss.Style{}.
			Foreground(t.Error).
			Background(t.Panel).
			Bold(true).
			Render("🔥 HIGH PRIORITY 🔥")
	case "medium":
		return lipgloss.Style{}.
			Foreground(t.Warning).
			Background(t.Panel).
			Bold(true).
			Render("⚡ MEDIUM ⚡")
	case "low":
		return lipgloss.Style{}.
			Foreground(t.Success).
			Background(t.Background).
			Bold(true).
			Render("💤 LOW PRIORITY 💤")
	default:
		return lipgloss.Style{}.
			Foreground(t.Secondary).
			Render("◉ NORMAL")
	}
}

func TaskStatus(status string) string {
	s := Current()
	if status == "completed" {
		return s.TaskDone.Render("✅ MISSION COMPLETE ✅")
	}
	return s.TaskActive.Render("◯ ACTIVE MISSION ◯")
}

func CyberTag(tag string) string {
	return Current().CyberTag.Render("🏷️ " + tag)
}

func Header() string {
	return Current().Banner.Render("focus.sh Task Management")
}

func CyberStats(active, completed, total int) string {
	s := Current()
	t := s.Theme
	return s.StatsBase.Render(
		"📊 GRID STATUS: " +
			lipgloss.Style{}.Foreground(t.Warning).Render("⚡ "+string(rune(active))) +
			" ACTIVE | " +
			lipgloss.Style{}.Foreground(t.Success).Render("✅ "+string(rune(completed))) +
			" COMPLETED | " +
			lipgloss.Style{}.Foreground(t.Secondary).Render("🌟 "+string(rune(total))) +
			" TOTAL",
	)
}

func LoadingMessage() string {
	return Current().Loading.Render("⠋ Loading...")
}

func EmptyStateMessage() string {
	return Current().EmptyState.Render("No tasks found.\n\nCreate your first task:\nfocus add \"task description\"")
}

func SynthwaveAIResponseStyle(response string) string {
	return Current().SynthAI.Render("🤖 AI SYNTHESIS: " + response)
}

func SynthwaveReportStyle(text string) string {
	return Current().SynthReport.Render("📋 ANALYSIS REPORT:\n\n" + text)
}

// ThemeMode and constants kept for backward compat with unified_dashboard.
type ThemeMode string

const (
	ThemeSynthwave ThemeMode = "synthwave"
	ThemeLight     ThemeMode = "light"
	ThemePlain     ThemeMode = "plain"
)
