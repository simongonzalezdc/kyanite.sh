package lyricforge_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/puente-labs/lyricforge/internal/ui"
)

// TestResponsiveLayoutManagerCreation tests responsive layout manager creation
func TestResponsiveLayoutManagerCreation(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	if manager == nil {
		t.Fatal("Expected responsive layout manager to be created, got nil")
	}

	// Test initial size
	size := manager.GetCurrentSize()
	if size.Width != 80 || size.Height != 24 {
		t.Errorf("Expected initial size 80x24, got %dx%d", size.Width, size.Height)
	}

	// Test initial breakpoint
	breakpoint := manager.GetActiveBreakpoint()
	if breakpoint.Name != "ultra-compact" {
		t.Errorf("Expected initial breakpoint to be ultra-compact, got %s", breakpoint.Name)
	}
}

// TestResponsiveBreakpoints tests breakpoint detection
func TestResponsiveBreakpoints(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test ultra-compact breakpoint
	manager.UpdateSize(60, 20)
	breakpoint := manager.GetActiveBreakpoint()
	if breakpoint.Name != "ultra-compact" {
		t.Errorf("Expected ultra-compact breakpoint for 60x20, got %s", breakpoint.Name)
	}

	// Test compact breakpoint
	manager.UpdateSize(90, 25)
	breakpoint = manager.GetActiveBreakpoint()
	if breakpoint.Name != "compact" {
		t.Errorf("Expected compact breakpoint for 90x25, got %s", breakpoint.Name)
	}

	// Test medium breakpoint
	manager.UpdateSize(110, 35)
	breakpoint = manager.GetActiveBreakpoint()
	if breakpoint.Name != "medium" {
		t.Errorf("Expected medium breakpoint for 110x35, got %s", breakpoint.Name)
	}

	// Test large breakpoint
	manager.UpdateSize(140, 50)
	breakpoint = manager.GetActiveBreakpoint()
	if breakpoint.Name != "large" {
		t.Errorf("Expected large breakpoint for 140x50, got %s", breakpoint.Name)
	}

	// Test ultra-wide breakpoint
	manager.UpdateSize(180, 80)
	breakpoint = manager.GetActiveBreakpoint()
	if breakpoint.Name != "ultra-wide" {
		t.Errorf("Expected ultra-wide breakpoint for 180x80, got %s", breakpoint.Name)
	}
}

// TestResponsiveSizeValidation tests size validation
func TestResponsiveSizeValidation(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test minimum size
	manager.UpdateSize(80, 24)
	if !manager.IsMinimumSize() {
		t.Error("Expected 80x24 to meet minimum size requirements")
	}

	// Test below minimum size
	manager.UpdateSize(70, 20)
	if manager.IsMinimumSize() {
		t.Error("Expected 70x20 to be below minimum size requirements")
	}

	// Test optimal size
	manager.UpdateSize(100, 30)
	if !manager.IsOptimalSize() {
		t.Error("Expected 100x30 to be optimal size")
	}

	// Test non-optimal size
	manager.UpdateSize(90, 25)
	if manager.IsOptimalSize() {
		t.Error("Expected 90x25 to not be optimal size")
	}
}

// TestResponsiveSplitRatio tests split ratio recommendations
func TestResponsiveSplitRatio(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test ultra-compact split ratio
	manager.UpdateSize(60, 20)
	ratio := manager.GetRecommendedSplitRatio()
	if ratio != 0.7 {
		t.Errorf("Expected split ratio 0.7 for ultra-compact, got %f", ratio)
	}

	// Test compact split ratio
	manager.UpdateSize(90, 25)
	ratio = manager.GetRecommendedSplitRatio()
	if ratio != 0.6 {
		t.Errorf("Expected split ratio 0.6 for compact, got %f", ratio)
	}

	// Test medium split ratio
	manager.UpdateSize(110, 35)
	ratio = manager.GetRecommendedSplitRatio()
	if ratio != 0.5 {
		t.Errorf("Expected split ratio 0.5 for medium, got %f", ratio)
	}

	// Test large split ratio
	manager.UpdateSize(140, 50)
	ratio = manager.GetRecommendedSplitRatio()
	if ratio != 0.45 {
		t.Errorf("Expected split ratio 0.45 for large, got %f", ratio)
	}

	// Test ultra-wide split ratio
	manager.UpdateSize(180, 80)
	ratio = manager.GetRecommendedSplitRatio()
	if ratio != 0.4 {
		t.Errorf("Expected split ratio 0.4 for ultra-wide, got %f", ratio)
	}
}

// TestResponsiveLayoutModes tests layout mode recommendations
func TestResponsiveLayoutModes(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test compact layout
	manager.UpdateSize(60, 20)
	if !manager.ShouldShowCompactLayout() {
		t.Error("Expected compact layout for 60x20")
	}

	if !manager.ShouldCollapsePanes() {
		t.Error("Expected panes to collapse for 60x20")
	}

	// Test non-compact layout
	manager.UpdateSize(120, 40)
	if manager.ShouldShowCompactLayout() {
		t.Error("Expected non-compact layout for 120x40")
	}

	if manager.ShouldCollapsePanes() {
		t.Error("Expected panes to not collapse for 120x40")
	}

	// Test status bar mode
	statusBarMode := manager.GetStatusBarMode()
	if statusBarMode != ui.StatusBarCompact {
		t.Errorf("Expected compact status bar for 120x40, got %v", statusBarMode)
	}

	// Test menu mode
	menuMode := manager.GetMenuMode()
	if menuMode != ui.MenuFull {
		t.Errorf("Expected full menu for 120x40, got %v", menuMode)
	}
}

// TestResponsiveSizeWarnings tests size warning generation
func TestResponsiveSizeWarnings(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test no warnings for optimal size
	manager.UpdateSize(120, 40)
	warnings := manager.GetSizeWarnings()
	if len(warnings) != 0 {
		t.Errorf("Expected no warnings for optimal size, got %d warnings", len(warnings))
	}

	// Test warnings for small size
	manager.UpdateSize(70, 20)
	warnings = manager.GetSizeWarnings()
	if len(warnings) == 0 {
		t.Error("Expected warnings for small size")
	}

	// Test warnings for ultra-wide
	manager.UpdateSize(220, 60)
	warnings = manager.GetSizeWarnings()
	if len(warnings) == 0 {
		t.Error("Expected warnings for ultra-wide terminal")
	}
}

// TestResponsiveSizeWarningRendering tests size warning rendering
func TestResponsiveSizeWarningRendering(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test no warning rendering for optimal size
	manager.UpdateSize(120, 40)
	warning := manager.RenderSizeWarning()
	if warning != "" {
		t.Errorf("Expected no warning for optimal size, got '%s'", warning)
	}

	// Test warning rendering for small size
	manager.UpdateSize(70, 20)
	warning = manager.RenderSizeWarning()
	if warning == "" {
		t.Error("Expected warning for small size")
	}

	// Verify warning contains expected content
	if !contains(warning, "⚠️") {
		t.Error("Expected warning to contain warning emoji")
	}
}

// TestResponsiveWindowSizeMsgHandling tests window size message handling
func TestResponsiveWindowSizeMsgHandling(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test valid size message
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	cmd := manager.HandleWindowSizeMsg(msg)

	if cmd == nil {
		t.Error("Expected command to be returned")
	}

	// Execute command
	sizeMsg := cmd()
	if sizeMsg == nil {
		t.Error("Expected size validation message to be returned")
	}

	// Test invalid size message
	msg = tea.WindowSizeMsg{Width: 70, Height: 20}
	cmd = manager.HandleWindowSizeMsg(msg)

	// Execute command
	sizeMsg = cmd()
	if sizeMsg == nil {
		t.Error("Expected size validation message to be returned")
	}
}

// TestResponsiveMinimumResolutionOptimizations tests minimum resolution optimizations
func TestResponsiveMinimumResolutionOptimizations(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test optimizations for minimum size
	manager.UpdateSize(80, 24)
	optimizations := manager.GetMinimumResolutionOptimizations()

	if optimizations == nil {
		t.Error("Expected optimizations for minimum size")
	}

	// Test specific optimizations
	if hideLineNumbers, ok := optimizations["hide_line_numbers"].(bool); ok && hideLineNumbers {
		t.Error("Expected line numbers to be shown at 80x24")
	}

	if minimalStatusBar, ok := optimizations["minimal_status_bar"].(bool); ok && !minimalStatusBar {
		t.Error("Expected minimal status bar at 80x24")
	}

	// Test no optimizations for below minimum size
	manager.UpdateSize(70, 20)
	optimizations = manager.GetMinimumResolutionOptimizations()

	if optimizations != nil {
		t.Error("Expected no optimizations for below minimum size")
	}
}

// TestResponsiveUltraCompactDetection tests ultra-compact mode detection
func TestResponsiveUltraCompactDetection(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test ultra-compact detection
	manager.UpdateSize(60, 20)
	if !manager.IsUltraCompact() {
		t.Error("Expected ultra-compact mode for 60x20")
	}

	// Test non-ultra-compact detection
	manager.UpdateSize(90, 25)
	if manager.IsUltraCompact() {
		t.Error("Expected non-ultra-compact mode for 90x25")
	}
}

// TestResponsiveOptimalDimensions tests optimal dimensions
func TestResponsiveOptimalDimensions(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test optimal dimensions
	dimensions := manager.GetOptimalDimensions()
	if dimensions.Width != 100 || dimensions.Height != 30 {
		t.Errorf("Expected optimal dimensions 100x30, got %dx%d", dimensions.Width, dimensions.Height)
	}
}

// TestResponsiveLayoutEfficiency tests layout efficiency calculation
func TestResponsiveLayoutEfficiency(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test efficiency for below minimum size
	manager.UpdateSize(70, 20)
	efficiency := manager.GetLayoutEfficiency()
	if efficiency != 0 {
		t.Errorf("Expected efficiency 0 for below minimum size, got %d", efficiency)
	}

	// Test efficiency for minimum size
	manager.UpdateSize(80, 24)
	efficiency = manager.GetLayoutEfficiency()
	if efficiency <= 0 || efficiency > 100 {
		t.Errorf("Expected efficiency between 0 and 100 for minimum size, got %d", efficiency)
	}

	// Test efficiency for optimal size
	manager.UpdateSize(100, 30)
	efficiency = manager.GetLayoutEfficiency()
	if efficiency < 80 {
		t.Errorf("Expected high efficiency for optimal size, got %d", efficiency)
	}

	// Test efficiency for large size
	manager.UpdateSize(140, 50)
	efficiency = manager.GetLayoutEfficiency()
	if efficiency < 80 {
		t.Errorf("Expected high efficiency for large size, got %d", efficiency)
	}
}

// TestResponsiveEdgeCases tests edge cases and boundary conditions
func TestResponsiveEdgeCases(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test zero size
	manager.UpdateSize(0, 0)
	if manager.IsMinimumSize() {
		t.Error("Expected zero size to not meet minimum requirements")
	}

	// Test extremely large size
	manager.UpdateSize(1000, 1000)
	if !manager.IsOptimalSize() {
		t.Error("Expected extremely large size to be considered optimal")
	}

	// Test boundary conditions
	manager.UpdateSize(80, 24) // Exactly minimum size
	if !manager.IsMinimumSize() {
		t.Error("Expected exactly minimum size to meet requirements")
	}

	manager.UpdateSize(100, 30) // Exactly optimal size
	if !manager.IsOptimalSize() {
		t.Error("Expected exactly optimal size to be optimal")
	}
}

// TestResponsiveWarningsToggle tests warnings toggle functionality
func TestResponsiveWarningsToggle(t *testing.T) {
	manager := ui.NewResponsiveLayoutManager()

	// Test warnings enabled by default
	manager.UpdateSize(70, 20)
	warnings := manager.GetSizeWarnings()
	if len(warnings) == 0 {
		t.Error("Expected warnings to be enabled by default")
	}

	// Test warnings can be disabled (if method is available)
	// Note: This assumes there's a SetShowWarnings method or similar
	// If not available, this test can be adjusted
}

// BenchmarkResponsiveBreakpointDetection benchmarks breakpoint detection performance
func BenchmarkResponsiveBreakpointDetection(b *testing.B) {
	manager := ui.NewResponsiveLayoutManager()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Test various sizes
		width := 60 + (i % 140)
		height := 20 + (i % 60)
		manager.UpdateSize(width, height)
		_ = manager.GetActiveBreakpoint()
	}
}

// BenchmarkResponsiveSplitRatioCalculation benchmarks split ratio calculation performance
func BenchmarkResponsiveSplitRatioCalculation(b *testing.B) {
	manager := ui.NewResponsiveLayoutManager()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Test various sizes
		width := 60 + (i % 140)
		height := 20 + (i % 60)
		manager.UpdateSize(width, height)
		_ = manager.GetRecommendedSplitRatio()
	}
}

// BenchmarkResponsiveLayoutEfficiency benchmarks layout efficiency calculation performance
func BenchmarkResponsiveLayoutEfficiency(b *testing.B) {
	manager := ui.NewResponsiveLayoutManager()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Test various sizes
		width := 60 + (i % 140)
		height := 20 + (i % 60)
		manager.UpdateSize(width, height)
		_ = manager.GetLayoutEfficiency()
	}
}
