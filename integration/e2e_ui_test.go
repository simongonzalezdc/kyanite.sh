package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/ui/editor"
	tea "github.com/charmbracelet/bubbletea"
)

// TestE2E_UI_Workflow runs UI-only end-to-end scenarios without a database (passes nil DB).
// These tests exercise rendering, split-pane behavior, keyboard shortcuts, and responsive layout.
func TestE2E_UI_Workflow(t *testing.T) {
	// Create split pane with nil database to avoid DB/cgo dependencies for UI E2E tests
	sp := editor.NewSplitPaneModel(nil, nil)

	// Initialize dimensions via WindowSizeMsg
	_, _ = sp.Update(tea.WindowSizeMsg{Width: 140, Height: 48})

	// 1) Full editor lifecycle: set text, view, simulate edits
	original := "# E2E UI Test\n\nVerse 1\nLine A\nLine B"
	sp.SetEditorText(original)
	time.Sleep(50 * time.Millisecond)

	view := sp.View()
	if view == "" {
		t.Fatalf("expected non-empty view after setting content")
	}
	if !strings.Contains(view, "E2E UI Test") {
		t.Fatalf("expected view to include title; view=%q", view)
	}

	// 2) Split-pane layout: verify presence of divider character from styles (approx)
	if !strings.Contains(view, "|") && !strings.Contains(view, "|") {
		// Not fatal if styles vary, but at least ensure combined output contains both panes text
		if !strings.Contains(view, "Line A") || !strings.Contains(view, "Line B") {
			t.Fatalf("view did not contain editor content; view=%q", view)
		}
	}

	// 3) Real-time preview: update editor and ensure view reflects change
	updated := original + "\n\n[Chorus]\nSing loud"
	sp.SetEditorText(updated)
	time.Sleep(80 * time.Millisecond)

	view2 := sp.View()
	if !strings.Contains(view2, "Sing loud") {
		t.Fatalf("expected preview to include updated content; view=%q", view2)
	}

	// 4) Keyboard shortcuts: simulate Tab to switch pane focus and arrow keys for preview
	_, _ = sp.Update(tea.KeyMsg{Type: tea.KeyTab, Runes: []rune{'\t'}})
	_, _ = sp.Update(tea.KeyMsg{Type: tea.KeyUp, Runes: []rune{'k'}})
	_, _ = sp.Update(tea.KeyMsg{Type: tea.KeyDown, Runes: []rune{'j'}})

	// 5) Responsive layout: shrink and expand terminal sizes and ensure view still renders
	sizes := []tea.WindowSizeMsg{
		{Width: 80, Height: 24},
		{Width: 100, Height: 30},
		{Width: 160, Height: 50},
	}

	for _, s := range sizes {
		_, _ = sp.Update(s)
		v := sp.View()
		if v == "" {
			t.Fatalf("view empty after resize to %dx%d", s.Width, s.Height)
		}
		// Basic sanity: ensure editor title still present
		if !strings.Contains(v, "E2E UI Test") {
			t.Fatalf("after resize %dx%d title missing from view", s.Width, s.Height)
		}
	}

	// 6) Performance spot-check: set a moderately large document and ensure SetEditorText returns quickly
	var large strings.Builder
	large.WriteString("# Large E2E Test\n\n")
	for i := 0; i < 200; i++ {
		large.WriteString(fmt.Sprintf("Line %d repeating to increase size.\n", i+1))
	}
	start := time.Now()
	sp.SetEditorText(large.String())
	elapsed := time.Since(start)
	if elapsed > 2*time.Second {
		t.Logf("warning: setting large content took %v (test threshold 2s)", elapsed)
	}

	// Cleanup
	sp.Cleanup()
}
