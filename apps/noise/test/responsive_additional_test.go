package noise

import (
	"strings"
	"testing"

	"github.com/Kyanite/noise/internal/ui"
)

func TestResponsiveManagerBasics(t *testing.T) {
	r := ui.NewResponsiveLayoutManager()
	if r == nil {
		t.Fatal("expected NewResponsiveLayoutManager to return non-nil")
	}

	size := r.GetCurrentSize()
	if size.Width <= 0 || size.Height <= 0 {
		t.Fatalf("expected sensible default size, got %v", size)
	}

	// Default active breakpoint should be one of the predefined names
	name := r.GetActiveBreakpoint().Name
	if name == "" {
		t.Fatal("expected active breakpoint to have a name")
	}
}

func TestResponsiveUpdateAndBreakpoints(t *testing.T) {
	r := ui.NewResponsiveLayoutManager()

	// Update to medium size and assert breakpoint
	r.UpdateSize(110, 35)
	if bp := r.GetActiveBreakpoint(); !strings.EqualFold(bp.Name, "medium") {
		t.Fatalf("expected medium breakpoint for 110x35, got %q", bp.Name)
	}

	// Update to ultra-compact
	r.UpdateSize(70, 20)
	if !r.IsUltraCompact() {
		t.Fatalf("expected IsUltraCompact true for 70x20, got false")
	}

	// Update to ultra-wide
	r.UpdateSize(200, 80)
	if !strings.Contains(strings.ToLower(r.GetActiveBreakpoint().Name), "ultra") {
		t.Fatalf("expected ultra breakpoint for 200x80, got %q", r.GetActiveBreakpoint().Name)
	}
}

func TestResponsiveModesAndRecommendations(t *testing.T) {
	r := ui.NewResponsiveLayoutManager()

	// Small terminals
	r.UpdateSize(80, 24)
	if !r.ShouldShowCompactLayout() {
		t.Fatalf("expected compact layout for 80x24")
	}
	if !r.ShouldCollapsePanes() {
		t.Fatalf("expected panes to collapse for width < 100")
	}

	// Larger terminals
	r.UpdateSize(120, 40)
	if r.ShouldShowCompactLayout() {
		t.Fatalf("did not expect compact layout for 120x40")
	}
	if r.ShouldCollapsePanes() {
		t.Fatalf("did not expect panes to collapse for 120x40")
	}

	// Recommended split ratio varies with breakpoint
	r.UpdateSize(70, 20)
	if ratio := r.GetRecommendedSplitRatio(); ratio <= 0 || ratio >= 1 {
		t.Fatalf("expected recommended split ratio in (0,1), got %v", ratio)
	}
}

func TestResponsiveSizeWarningsAndRendering(t *testing.T) {
	r := ui.NewResponsiveLayoutManager()

	// Too small
	r.UpdateSize(40, 10)
	warns := r.GetSizeWarnings()
	if len(warns) == 0 {
		t.Fatalf("expected warnings for 40x10")
	}
	rendered := r.RenderSizeWarning()
	if rendered == "" {
		t.Fatalf("expected non-empty rendered size warning")
	}

	// Toggle to acceptable size
	r.UpdateSize(120, 40)
	warns = r.GetSizeWarnings()
	if len(warns) == 0 {
		// acceptable - may be empty when size is fine
	} else {
		t.Logf("warnings present at 120x40: %v", warns)
	}
}

func TestHandleWindowSizeMsgProducesValidationMsg(t *testing.T) {
	r := ui.NewResponsiveLayoutManager()

	// Very small -> invalid
	cmd := r.HandleWindowSizeMsg(struct{ Width, Height int }{Width: 40, Height: 10})
	if cmd == nil {
		t.Fatalf("expected non-nil command for invalid size handling")
	}
	// Execute the command to get the message and inspect it
	msg := cmd()
	if sv, ok := msg.(ui.SizeValidationMsg); ok {
		if sv.IsValid {
			t.Fatalf("expected IsValid=false for 40x10")
		}
		if sv.Width != 40 || sv.Height != 10 {
			t.Fatalf("expected SizeValidationMsg to reflect requested size, got %dx%d", sv.Width, sv.Height)
		}
	} else {
		t.Fatalf("expected SizeValidationMsg from command, got %T", msg)
	}

	// Acceptable size -> valid
	cmd = r.HandleWindowSizeMsg(struct{ Width, Height int }{Width: 120, Height: 40})
	msg = cmd()
	if sv, ok := msg.(ui.SizeValidationMsg); ok {
		if !sv.IsValid {
			t.Fatalf("expected IsValid=true for 120x40")
		}
	} else {
		t.Fatalf("expected SizeValidationMsg from command, got %T", msg)
	}
}

func TestMinimumResolutionOptimizationsAndEfficiency(t *testing.T) {
	r := ui.NewResponsiveLayoutManager()

	// At minimum 80x24
	r.UpdateSize(80, 24)
	opts := r.GetMinimumResolutionOptimizations()
	if opts == nil {
		t.Fatalf("expected optimizations map at minimum size")
	}
	if _, ok := opts["reduced_padding"]; !ok {
		t.Fatalf("expected reduced_padding key in optimizations")
	}

	// Layout efficiency score should be within 0-100
	score := r.GetLayoutEfficiency()
	if score < 0 || score > 100 {
		t.Fatalf("expected layout efficiency within 0-100, got %d", score)
	}
}

func TestOptimalDimensionsAndHelpers(t *testing.T) {
	r := ui.NewResponsiveLayoutManager()
	opt := r.GetOptimalDimensions()
	if opt.Width <= 0 || opt.Height <= 0 {
		t.Fatalf("expected positive optimal dimensions, got %v", opt)
	}
	if eff := r.GetLayoutEfficiency(); eff >= 0 {
		// just ensure method runs
		_ = eff
	}
}
