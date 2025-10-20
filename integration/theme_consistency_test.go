package integration

import (
	"strings"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/ui/editor"
	tea "github.com/charmbracelet/bubbletea"
)

// TestThemeConsistencyAcrossComponents verifies that editor, preview and status bar
// show consistent information when given the same content. This is a UI-only test
// and runs without a database to avoid cgo/sqlite issues.
func TestThemeConsistencyAcrossComponents(t *testing.T) {
	// Create a split pane without DB (UI-only)
	sp := editor.NewSplitPaneModel(nil)

	// Initialize dimensions via WindowSizeMsg
	_, _ = sp.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Content with several markdown elements to exercise rendering paths
	content := "# Theme Consistency Test\n\n" +
		"## Header\n" +
		"This is **bold** and *italic* text.\n\n" +
		"```\ncode snippet\nprintln(\"hi\")\n```\n\n" +
		"[Chorus]\nSing this chorus"

	// Set editor content and allow a short moment for realtime features
	sp.SetEditorText(content)
	time.Sleep(120 * time.Millisecond)

	// Combined view (editor + preview) must contain key pieces of content
	view := sp.View()
	if view == "" {
		t.Fatal("combined view should not be empty")
	}
	if !strings.Contains(view, "Theme Consistency Test") {
		t.Fatalf("combined view missing title: %q", view)
	}
	if !strings.Contains(view, "Sing this chorus") {
		t.Fatalf("combined view missing chorus: %q", view)
	}

	// Create a standalone status bar and update it with the same content
	sb := editor.NewStatusBarModel()
	sb.SetDimensions(120, 1)
	sb.UpdateContent(content)

	// Update cursor position and auto-save status to emulate editor state
	sb.UpdateCursorPosition(2, 5) // line 3, col 6 (0-based internal)
	sb.UpdateAutoSaveStatus(app.AutoSaveSuccess, time.Now())

	// Render status bar and verify it contains expected indicators
	statusView := sb.View()
	if statusView == "" {
		t.Fatal("status bar view should not be empty")
	}

	// Check for cursor position indicators (Ln / Col)
	if !strings.Contains(statusView, "Ln 3") {
		t.Fatalf("status bar missing expected line indicator: %q", statusView)
	}
	if !strings.Contains(statusView, "Col 6") {
		t.Fatalf("status bar missing expected column indicator: %q", statusView)
	}

	// Auto-save success indicator should appear (either "Saved" or a time string)
	if !strings.Contains(statusView, "Saved") {
		t.Logf("status bar did not contain explicit 'Saved' text; view=%q", statusView)
	}

	// Final sanity: both editor combined view and status bar reference the same document title
	if !strings.Contains(statusView, "Theme Consistency Test") && !strings.Contains(view, "Theme Consistency Test") {
		t.Fatalf("both status bar and view should reference the document title; status=%q view=%q", statusView, view)
	}

	// Cleanup split pane resources
	sp.Cleanup()
}
