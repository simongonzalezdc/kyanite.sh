package styles

import (
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
)

// Styles is the centralized, theme-derived style set for focus.
//
// Every style is built once per theme change from the live design.Theme
// (via lipgloss.Style{} compositions seeded with theme colors, plus the
// shared design.TokenSet exposed for consumers). Callers never construct
// styles inline — they read a field from Current() instead, so a theme
// switch repaints the whole app from one place.
type Styles struct {
	Theme  design.Theme
	Tokens design.TokenSet

	// Style getters (returned verbatim by the Get*Style helpers).
	Title   lipgloss.Style
	Box     lipgloss.Style
	Success lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style

	// Render-fn base styles (.Render(text) applied by the render helper).
	FocusStyle     lipgloss.Style
	FocusPink      lipgloss.Style
	FocusBlue      lipgloss.Style
	FocusGreen     lipgloss.Style
	FocusPurple    lipgloss.Style
	ID             lipgloss.Style
	HeaderStyle    lipgloss.Style
	FooterStyle    lipgloss.Style
	Category       lipgloss.Style
	Suggestion     lipgloss.Style
	AIResponse     lipgloss.Style
	Report         lipgloss.Style
	SynthwaveTitle lipgloss.Style
	FocusBoxBase   lipgloss.Style // FocusBox adds BorderForeground(borderColor)
	CyberGrid      lipgloss.Style
	DigitalRain    lipgloss.Style
	TaskDone       lipgloss.Style
	TaskActive     lipgloss.Style
	CyberTag       lipgloss.Style
	Banner         lipgloss.Style
	Loading        lipgloss.Style
	EmptyState     lipgloss.Style
	SynthAI        lipgloss.Style
	SynthReport    lipgloss.Style
	StatsBase      lipgloss.Style // CyberStats outer style; inner spans compose at call time
}

var (
	stylesMu sync.RWMutex
	current  = buildStyles(design.DefaultTheme())
)

// buildStyles derives the full style set from a theme.
// Style modifiers are transcribed verbatim from the previous inline
// constructor chains so output bytes are unchanged.
func buildStyles(t design.Theme) Styles {
	ts := design.NewTokenSet(t)
	return Styles{
		Theme:  t,
		Tokens: ts,

		Title: lipgloss.Style{}.
			Foreground(t.Primary).
			Background(t.Background).
			Bold(true).
			Padding(design.SpacingS).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent),
		Box: lipgloss.Style{}.
			Foreground(t.Text).
			Background(t.Panel).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Border).
			Padding(design.SpacingS),
		Success: lipgloss.Style{}.Foreground(t.Success).Bold(true),
		Error:   lipgloss.Style{}.Foreground(t.Error).Bold(true),
		Warning: lipgloss.Style{}.Foreground(t.Warning).Bold(true),

		FocusStyle:  lipgloss.Style{}.Foreground(t.Primary).Bold(true),
		FocusPink:   lipgloss.Style{}.Foreground(t.Primary),
		FocusBlue:   lipgloss.Style{}.Foreground(t.Secondary),
		FocusGreen:  lipgloss.Style{}.Foreground(t.Success),
		FocusPurple: lipgloss.Style{}.Foreground(t.Accent),
		ID:          lipgloss.Style{}.Foreground(t.Secondary),
		HeaderStyle: lipgloss.Style{}.
			Foreground(t.Primary).
			Bold(true).
			AlignHorizontal(lipgloss.Center).
			Underline(true),
		FooterStyle: lipgloss.Style{}.Foreground(t.Accent).Italic(true),
		Category: lipgloss.Style{}.
			Foreground(t.Secondary).
			Background(t.Background).
			Padding(0, design.SpacingXS).
			Bold(true),
		Suggestion: lipgloss.Style{}.Foreground(t.Success).Bold(true),
		AIResponse: lipgloss.Style{}.Foreground(t.Accent).Italic(true),
		Report: lipgloss.Style{}.
			Foreground(t.Accent).
			Padding(design.SpacingS).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Secondary),
		SynthwaveTitle: lipgloss.Style{}.
			Foreground(t.Primary).
			Background(t.Background).
			Bold(true).
			Italic(true).
			Underline(true).
			Padding(design.SpacingS, design.SpacingM).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent).
			AlignHorizontal(lipgloss.Center),
		FocusBoxBase: lipgloss.Style{}.
			Foreground(t.Accent).
			Background(t.Panel).
			Padding(design.SpacingS, design.SpacingM).
			Border(lipgloss.RoundedBorder()).
			Bold(true),
		CyberGrid: lipgloss.Style{}.
			Foreground(t.Success).
			Background(t.Panel).
			Padding(0, design.SpacingXS).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent),
		DigitalRain: lipgloss.Style{}.
			Foreground(t.Success).
			Background(t.Background).
			Bold(true).
			Italic(true),
		TaskDone:   lipgloss.Style{}.Foreground(t.Success).Background(t.Panel).Bold(true),
		TaskActive: lipgloss.Style{}.Foreground(t.Accent).Background(t.Panel).Bold(true),
		CyberTag: lipgloss.Style{}.
			Foreground(t.Primary).
			Background(t.Panel).
			Padding(0, design.SpacingXS).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Secondary),
		Banner: lipgloss.Style{}.
			Foreground(t.Success).
			Background(t.Background).
			Bold(true).
			AlignHorizontal(lipgloss.Center).
			Padding(design.SpacingS).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent),
		Loading: lipgloss.Style{}.Foreground(t.Accent).Bold(true),
		EmptyState: lipgloss.Style{}.
			Foreground(t.Accent).
			Background(t.Background).
			Bold(true).
			AlignHorizontal(lipgloss.Center).
			Padding(design.SpacingM).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary),
		SynthAI: lipgloss.Style{}.
			Foreground(t.Accent).
			Background(t.Panel).
			Bold(true).
			Italic(true).
			Padding(design.SpacingS).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Secondary),
		SynthReport: lipgloss.Style{}.
			Foreground(t.Primary).
			Background(t.Panel).
			Padding(design.SpacingM).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent).
			Bold(true),
		StatsBase: lipgloss.Style{}.
			Foreground(t.Accent).
			Background(t.Panel).
			Bold(true).
			Padding(design.SpacingS).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary),
	}
}

// Current returns the cached, theme-derived style set.
func Current() Styles {
	stylesMu.RLock()
	defer stylesMu.RUnlock()
	return current
}

// rebuildStyles rebuilds the cached style set from the named theme.
// The caller must hold themeMu for writing (it serializes theme changes);
// rebuildStyles separately takes stylesMu to publish the new set.
func rebuildStyles(name string) {
	t := design.Get(name)
	if t.Name == "" {
		t = design.DefaultTheme()
	}
	s := buildStyles(t)
	stylesMu.Lock()
	current = s
	stylesMu.Unlock()
}
