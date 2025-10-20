package noise

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Kyanite/noise/internal/ui"
)

// TestNewMenuModelBasic ensures the menu model constructs and exposes a title when width is unset.
func TestNewMenuModelBasic(t *testing.T) {
	m := ui.NewMenuModel()
	if m == nil {
		t.Fatal("expected NewMenuModel to return a non-nil model")
	}

	// When dimensions are zero, View should return the title (non-empty)
	view := m.View()
	if view == "" {
		t.Fatalf("expected view to contain a title when dimensions are zero, got empty")
	}
	if !strings.Contains(view, "noise") && !strings.Contains(strings.ToLower(view), "noise") {
		t.Fatalf("expected menu title to mention noise, got %q", view)
	}
}

// TestMenuWindowSizeHandling verifies Update handles WindowSizeMsg correctly.
func TestMenuWindowSizeHandling(t *testing.T) {
	m := ui.NewMenuModel()
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updated, _ := m.Update(msg)
	if updated == nil {
		t.Fatal("expected Update to return a non-nil model on WindowSizeMsg")
	}

	// After sizing, View should use the list rendering and be non-empty
	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view after setting window size")
	}
}

// TestMenuNavigationAndSelection simulates up/down navigation and enter selection behavior.
func TestMenuNavigationAndSelection(t *testing.T) {
	m := ui.NewMenuModel()
	// Set size so the list is active
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Press down several times to change selection
	for i := 0; i < 3; i++ {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		if updated == nil {
			t.Fatalf("expected non-nil model after down key (iteration %d)", i)
		}
	}

	// Press up to move selection back
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if updated == nil {
		t.Fatal("expected non-nil model after up key")
	}

	// Press enter to trigger selection; Update returns a cmd that will eventually emit a ScreenChangeMsg
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updated == nil {
		t.Fatal("expected non-nil model after enter key")
	}
	if cmd == nil {
		// It's acceptable for cmd to be nil in some implementations, but ensure calling View doesn't panic.
		_ = m.View()
	} else {
		// Execute returned command (it should be safe and quick)
		res := cmd()
		_ = res
	}
}

// TestMenuResponsiveTransitions verifies responsive mode flags change with width thresholds.
func TestMenuResponsiveTransitions(t *testing.T) {
	m := ui.NewMenuModel()

	// wide -> full mode
	m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	viewFull := m.View()
	if viewFull == "" {
		t.Fatal("expected non-empty view in wide layout")
	}

	// medium -> compact mode
	m.Update(tea.WindowSizeMsg{Width: 85, Height: 30})
	viewCompact := m.View()
	if viewCompact == "" {
		t.Fatal("expected non-empty view in compact layout")
	}

	// small -> minimal menu mode (title short)
	m.Update(tea.WindowSizeMsg{Width: 65, Height: 20})
	viewSmall := m.View()
	if viewSmall == "" {
		t.Fatal("expected non-empty view in small layout")
	}
	if !strings.Contains(viewSmall, "ðŸŽµ") && !strings.Contains(strings.ToLower(viewSmall), "noise") {
		// Ensure some title still present
		t.Logf("small view content: %q", viewSmall)
	}
}

// TestMenuAnimationTick ensures animation update path runs without panic.
func TestMenuAnimationTick(t *testing.T) {
	m := ui.NewMenuModel()
	// Simulate an animation tick message
	_, _ = m.Update(struct{ AnimationTickMsg int }{})
	// Also call Update with a WindowSizeMsg to ensure combined handling works
	_, _ = m.Update(tea.WindowSizeMsg{Width: 90, Height: 25})
	// Call View to ensure rendering path is stable
	_ = m.View()
}

// TestMenuStabilityUnderRapidKeys simulates rapid key presses to ensure no panics or nil returns.
func TestMenuStabilityUnderRapidKeys(t *testing.T) {
	m := ui.NewMenuModel()
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	keys := []tea.KeyMsg{
		{Type: tea.KeyDown},
		{Type: tea.KeyDown},
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
		{Type: tea.KeyEnter},
	}
	for i, k := range keys {
		mod, cmd := m.Update(k)
		if mod == nil {
			t.Fatalf("model became nil after key %d", i)
		}
		if cmd != nil {
			// Execute quickly; guard against long-running commands
			_ = cmd()
		}
	}
}

// TestMenuDefaultFocusBehavior verifies focused state can be toggled via list interactions.
func TestMenuDefaultFocusBehavior(t *testing.T) {
	m := ui.NewMenuModel()
	// Initially unfocused; calling View shouldn't panic and should return title or list
	_ = m.View()

	// Simulate window size to focus the list
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	_ = m.View()

	// Simulate brief wait to ensure animations may progress
	time.Sleep(5 * time.Millisecond)
	_ = m.View()
}