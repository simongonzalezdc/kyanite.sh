package noise

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/ui/editor"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(input string) string {
	return ansiEscape.ReplaceAllString(input, "")
}

func TestStatusBarUpdateContentStatistics(t *testing.T) {
	m := editor.NewStatusBarModel()
	m.SetDimensions(120, 1)
	content := "Hello world\nSecond line"
	m.UpdateContent(content)
	wc, cc, lc := m.GetDocumentStats()

	if wc != 4 {
		t.Fatalf("expected 4 words, got %d", wc)
	}
	if lc != 2 {
		t.Fatalf("expected 2 lines, got %d", lc)
	}
	if cc != len(strings.TrimSpace(content)) {
		t.Fatalf("expected %d chars, got %d", len(strings.TrimSpace(content)), cc)
	}
}

func TestStatusBarCursorAndFeatures(t *testing.T) {
	m := editor.NewStatusBarModel()
	m.UpdateCursorPosition(4, 10)
	line, col := m.GetCursorPosition()
	if line != 4 || col != 10 {
		t.Fatalf("expected cursor 4,10 got %d,%d", line, col)
	}

	m.UpdateEditorFeatures(false, false, false, false)
	a, b, c, d := m.GetEditorFeatures()
	if a || b || c || d {
		t.Fatalf("expected all features false got %v,%v,%v,%v", a, b, c, d)
	}
}

func TestStatusBarAutoSaveView(t *testing.T) {
	m := editor.NewStatusBarModel()
	m.SetDimensions(200, 1)

	// Saving state should render "Saving"
	m.UpdateAutoSaveStatus(app.AutoSaveSaving, time.Time{})
	v := m.View()
	plain := stripANSI(v)
	if !strings.Contains(plain, "Saving") {
		t.Fatalf("expected view to contain Saving, got %q", plain)
	}

	// Success state with timestamp should render "Saved"
	ttime := time.Date(2020, 1, 2, 15, 4, 5, 0, time.UTC)
	m.UpdateAutoSaveStatus(app.AutoSaveSuccess, ttime)
	v = m.View()
	plain = stripANSI(v)
	if !strings.Contains(plain, "Saved") {
		t.Fatalf("expected view to contain Saved, got %q", plain)
	}
}

func TestStatusBarContentTypeAndKB(t *testing.T) {
	m := editor.NewStatusBarModel()
	m.SetDimensions(200, 1)

	m.UpdateContentType("lyrics")
	m.UpdateKnowledgeBaseStatus(true, "KB: Ready")
	v := m.View()
	plain := stripANSI(v)

	if !strings.Contains(strings.ToLower(plain), "lyrics") {
		t.Fatalf("expected view to contain content type lyrics, got %q", plain)
	}
	if !strings.Contains(plain, "KB: Ready") {
		t.Fatalf("expected view to contain KB: Ready, got %q", plain)
	}
}

func TestStatusBarResponsiveViews(t *testing.T) {
	m := editor.NewStatusBarModel()

	// Minimal mode
	m.SetDimensions(70, 1)
	m.UpdateResponsiveMode(70)
	v := m.View()
	if !strings.Contains(v, "Ln") {
		t.Fatalf("expected minimal view to contain Ln, got %q", v)
	}

	// Compact mode
	m.SetDimensions(90, 1)
	m.UpdateResponsiveMode(90)
	v = m.View()
	if !strings.Contains(v, "Ln") {
		t.Fatalf("expected compact view to contain Ln, got %q", v)
	}

	// Full mode
	m.SetDimensions(140, 1)
	m.UpdateResponsiveMode(140)
	v = m.View()
	if v == "" {
		t.Fatalf("expected full view to be non-empty")
	}
}

func TestStatusBarShortcutHintsAndZoom(t *testing.T) {
	m := editor.NewStatusBarModel()
	m.SetDimensions(200, 1)

	m.UpdateShortcutHints("Ctrl+S")
	m.UpdateZoomLevel(125)
	v := m.View()
	plain := stripANSI(v)

	if !strings.Contains(plain, "Ctrl+S") {
		t.Fatalf("expected view to contain shortcut hints, got %q", plain)
	}

	if !strings.Contains(plain, "125%") {
		t.Fatalf("expected view to contain zoom level, got %q", plain)
	}
}

func TestStatusBarGettersAndMode(t *testing.T) {
	m := editor.NewStatusBarModel()
	m.SetDimensions(120, 1)

	m.UpdateFileInfo("song.txt")
	// Ensure getters run without panic
	_, _, _ = m.GetDocumentStats()

	m.UpdateContentType("patterns")
	v := m.View()
	if !strings.Contains(strings.ToLower(v), "patterns") {
		t.Fatalf("expected view to contain patterns, got %q", v)
	}
}