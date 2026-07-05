package styles

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
	"github.com/muesli/termenv"
)

// init forces the TrueColor profile so Render() emits deterministic 24-bit
// ANSI escapes regardless of whether stdout is a TTY. Without this, a
// headless test run renders with no color and the byte-level assertions
// below become vacuous. The property assertions (GetForeground) hold either
// way; the forced profile only makes the *rendered* bytes inspectable.
func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}

// hexRGB parses a "#RRGGBB" lipgloss.Color into its 8-bit channels.
func hexRGB(c lipgloss.Color) (int, int, int) {
	hex := strings.TrimPrefix(string(c), "#")
	if len(hex) != 6 {
		return -1, -1, -1
	}
	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)
	return int(r), int(g), int(b)
}

// fgRe matches every TrueColor foreground run "38;2;R;G;B" inside the
// rendered bytes. lipgloss coalesces SGR attributes into one escape
// (e.g. "\x1b[1;38;2;R;G;Bm" for bold+fg, or "\x1b[3;38;2;R;G;Bm" for
// italic+fg), and bordered styles emit the border color before the text
// color — so we cannot just read the first escape. Scanning every match
// and asking "is the expected color among them?" handles both cases.
var fgRe = regexp.MustCompile(`38;2;(\d{1,3});(\d{1,3});(\d{1,3})`)

// near is a per-channel tolerance. lipgloss/termenv's TrueColor path can
// shift a channel by ±1 vs the raw hex (a color-pipeline rounding artifact,
// not a theming bug), so we compare within a small band rather than exactly.
func near(a, b int) bool {
	if a -= b; a < 0 {
		a = -a
	}
	return a <= 3
}

// containsFG reports whether rendered bytes carry the given theme color as
// a TrueColor foreground (within tolerance), anywhere in the output.
func containsFG(rendered string, want lipgloss.Color) bool {
	wr, wg, wb := hexRGB(want)
	if wr < 0 {
		return false
	}
	for _, m := range fgRe.FindAllStringSubmatch(rendered, -1) {
		r, _ := strconv.Atoi(m[1])
		g, _ := strconv.Atoi(m[2])
		b, _ := strconv.Atoi(m[3])
		if near(r, wr) && near(g, wg) && near(b, wb) {
			return true
		}
	}
	return false
}

// renderThemes spans three visually distinct palettes. If any style is
// accidentally pinned to a single theme (the bug class the centralized
// hub removed), the track + differ assertions fail.
var renderThemes = []string{"amber-night", "indigo-depths", "forest-path"}

// TestStylesTrackLiveTheme proves every focus bridge style reads the LIVE
// theme: after SetThemeByName, each style's foreground matches the theme
// color (property), and the rendered bytes carry that theme's TrueColor
// color (render). Regression guard for "styles pinned to
// design.DefaultTheme() at package load".
func TestStylesTrackLiveTheme(t *testing.T) {
	original := GetCurrentThemeName()
	defer SetThemeByName(original)

	for _, name := range renderThemes {
		t.Run(name, func(t *testing.T) {
			SetThemeByName(name)
			want := design.Get(name)
			s := Current()

			cases := []struct {
				label  string
				got    lipgloss.TerminalColor
				want   lipgloss.Color
				render lipgloss.Style
				text   string
			}{
				{"Title", s.Title.GetForeground(), want.Primary, s.Title, "  FOCUS  "},
				{"Banner", s.Banner.GetForeground(), want.Success, s.Banner, "  ready  "},
				{"AIResponse", s.AIResponse.GetForeground(), want.Accent, s.AIResponse, "  suggestion  "},
				{"SynthwaveTitle", s.SynthwaveTitle.GetForeground(), want.Primary, s.SynthwaveTitle, "  synth  "},
				{"FocusPink", s.FocusPink.GetForeground(), want.Primary, s.FocusPink, "  pink  "},
				{"FooterStyle", s.FooterStyle.GetForeground(), want.Accent, s.FooterStyle, "  footer  "},
			}
			for _, c := range cases {
				// GetForeground returns the TerminalColor interface; the stored
				// value is the concrete lipgloss.Color we seeded the style with.
				gc, ok := c.got.(lipgloss.Color)
				if !ok || gc != c.want {
					t.Errorf("%s foreground: theme=%s got %v want %s", c.label, name, c.got, c.want)
					continue
				}
				if !containsFG(c.render.Render(c.text), c.want) {
					t.Errorf("%s render: theme=%s output missing TrueColor color ~%s",
						c.label, name, c.want)
				}
			}
		})
	}
}

// TestStylesDifferAcrossTheme proves the Primary/Accent-derived styles
// actually CHANGE with the theme — not merely equal to one constant value.
// (Success-derived styles like Banner are excluded here because amber-night
// and indigo-depths intentionally share the same Success color #52D3AA.)
func TestStylesDifferAcrossTheme(t *testing.T) {
	original := GetCurrentThemeName()
	defer SetThemeByName(original)

	SetThemeByName("amber-night")
	amberTitle := Current().Title.GetForeground()
	amberAI := Current().AIResponse.GetForeground()
	amberSynth := Current().SynthwaveTitle.GetForeground()

	SetThemeByName("indigo-depths")
	indigoTitle := Current().Title.GetForeground()
	indigoAI := Current().AIResponse.GetForeground()
	indigoSynth := Current().SynthwaveTitle.GetForeground()

	if amberTitle == indigoTitle {
		t.Errorf("Title foreground did not change across themes: %s", amberTitle)
	}
	if amberAI == indigoAI {
		t.Errorf("AIResponse foreground did not change across themes: %s", amberAI)
	}
	if amberSynth == indigoSynth {
		t.Errorf("SynthwaveTitle foreground did not change across themes: %s", amberSynth)
	}
}

// TestRenderSample emits TrueColor-rendered samples to the test log. Capture
// with: go test -v -run TestRenderSample ./pkg/styles/   and pipe to a file
// to eyeball the theming in any color terminal.
func TestRenderSample(t *testing.T) {
	original := GetCurrentThemeName()
	defer SetThemeByName(original)
	for _, name := range renderThemes {
		SetThemeByName(name)
		s := Current()
		t.Logf("=== theme: %s ===", name)
		t.Logf("title : %s", s.Title.Render("  FOCUS  "))
		t.Logf("banner: %s", s.Banner.Render("  ready  "))
		t.Logf("ai    : %s", s.AIResponse.Render("  suggestion  "))
		t.Logf("synth : %s", s.SynthwaveTitle.Render("  synth  "))
		t.Logf("footer: %s", s.FooterStyle.Render("  footer  "))
	}
}
