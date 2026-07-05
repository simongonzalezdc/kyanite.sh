package errors

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/noise/internal/theme"
	"github.com/muesli/termenv"
)

// init forces TrueColor so Render() emits deterministic 24-bit ANSI escapes
// in this headless test. The GetForeground property assertions hold without
// it; the profile only makes the rendered bytes inspectable.
func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}

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

// fgRe matches every TrueColor foreground run. lipgloss coalesces SGR
// attributes into one escape (e.g. "\x1b[1;38;2;R;G;Bm" for bold+fg), so we
// scan every match and ask whether the expected color is present.
var fgRe = regexp.MustCompile(`38;2;(\d{1,3});(\d{1,3});(\d{1,3})`)

func near(a, b int) bool {
	if a -= b; a < 0 {
		a = -a
	}
	return a <= 3
}

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

// restoreTheme captures the user's current noise theme and restores it after
// the test, persisting synchronously so the final on-disk state is the
// original. SetTheme persists asynchronously via a goroutine; SaveThemePreference
// writes synchronously and therefore wins the final ~/.config/noise/theme.json
// content. t.Cleanup runs in LIFO order: SetTheme (registered last) runs
// first, then SaveThemePreference.
func restoreTheme(t *testing.T) {
	mgr := theme.GetManager()
	origName := mgr.Current().Name
	t.Cleanup(func() {
		if err := mgr.SaveThemePreference(); err != nil {
			t.Logf("restore: SaveThemePreference: %v", err)
		}
	})
	t.Cleanup(func() {
		mgr.SetTheme(origName)
	})
}

// TestNotificationStyleTracksLiveTheme proves the ed98c2e fix: notificationStyle
// reads the live manager theme on every call, so each type's foreground follows
// the active theme instead of being pinned to design.DefaultTheme() at package
// load.
func TestNotificationStyleTracksLiveTheme(t *testing.T) {
	restoreTheme(t)
	mgr := theme.GetManager()

	themes := []string{"amber-night", "indigo-depths", "forest-path"}
	for _, name := range themes {
		t.Run(name, func(t *testing.T) {
			mgr.SetTheme(name)
			cur := mgr.Current()

			cases := []struct {
				label string
				nt    NotificationType
				want  lipgloss.Color
			}{
				{"Error", NotificationError, cur.Error},
				{"Warning", NotificationWarning, cur.Warning},
				{"Success", NotificationSuccess, cur.Success},
				{"Info", NotificationInfo, cur.Accent},
			}
			for _, c := range cases {
				style := notificationStyle(c.nt)
				if got := style.GetForeground(); got != c.want {
					t.Errorf("%s foreground: theme=%s got %v want %s", c.label, name, got, c.want)
					continue
				}
				if !containsFG(style.Render(c.label), c.want) {
					t.Errorf("%s render: theme=%s output missing TrueColor color ~%s",
						c.label, name, c.want)
				}
			}
		})
	}
}

// TestNotificationStyleDiffersAcrossTheme proves the foreground actually
// changes with the theme — not a pinned constant.
func TestNotificationStyleDiffersAcrossTheme(t *testing.T) {
	restoreTheme(t)
	mgr := theme.GetManager()

	mgr.SetTheme("amber-night")
	amberErr := notificationStyle(NotificationError).GetForeground()
	amberWarn := notificationStyle(NotificationWarning).GetForeground()

	mgr.SetTheme("indigo-depths")
	indigoErr := notificationStyle(NotificationError).GetForeground()
	indigoWarn := notificationStyle(NotificationWarning).GetForeground()

	if amberErr == indigoErr {
		t.Errorf("Error foreground did not change across themes: %s", amberErr)
	}
	if amberWarn == indigoWarn {
		t.Errorf("Warning foreground did not change across themes: %s", amberWarn)
	}
}

// TestRenderSample emits TrueColor-rendered notification samples to the test
// log. Capture with: go test -v -run TestRenderSample ./internal/errors/
func TestRenderSample(t *testing.T) {
	restoreTheme(t)
	mgr := theme.GetManager()
	for _, name := range []string{"amber-night", "indigo-depths", "forest-path"} {
		mgr.SetTheme(name)
		t.Logf("=== theme: %s ===", name)
		t.Logf("error  : %s", notificationStyle(NotificationError).Render("  upload failed  "))
		t.Logf("warn   : %s", notificationStyle(NotificationWarning).Render("  low disk  "))
		t.Logf("success: %s", notificationStyle(NotificationSuccess).Render("  saved  "))
		t.Logf("info   : %s", notificationStyle(NotificationInfo).Render("  synced  "))
	}
}
