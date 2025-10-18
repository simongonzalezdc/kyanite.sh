package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/noise/internal/ui/styles"
)

// TerminalSize represents terminal dimensions
type TerminalSize struct {
	Width  int
	Height int
}

// ResponsiveBreakpoint defines responsive breakpoints for different terminal sizes
type ResponsiveBreakpoint struct {
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
	Name      string
}

// Predefined breakpoints for common terminal sizes
var (
	BreakpointUltraCompact = ResponsiveBreakpoint{0, 80, 0, 24, "ultra-compact"}
	BreakpointCompact      = ResponsiveBreakpoint{80, 100, 24, 30, "compact"}
	BreakpointMedium       = ResponsiveBreakpoint{100, 120, 30, 40, "medium"}
	BreakpointLarge        = ResponsiveBreakpoint{120, 160, 40, 60, "large"}
	BreakpointUltraWide    = ResponsiveBreakpoint{160, 9999, 60, 9999, "ultra-wide"}
)

// ResponsiveLayoutManager manages responsive layout behavior
type ResponsiveLayoutManager struct {
	currentSize      TerminalSize
	minSize          TerminalSize
	breakpoints      []ResponsiveBreakpoint
	activeBreakpoint ResponsiveBreakpoint
	sizeWarnings     []string
	showWarnings     bool
}

// NewResponsiveLayoutManager creates a new responsive layout manager
func NewResponsiveLayoutManager() *ResponsiveLayoutManager {
	return &ResponsiveLayoutManager{
		currentSize: TerminalSize{Width: 80, Height: 24},
		minSize:     TerminalSize{Width: 80, Height: 24},
		breakpoints: []ResponsiveBreakpoint{
			BreakpointUltraCompact,
			BreakpointCompact,
			BreakpointMedium,
			BreakpointLarge,
			BreakpointUltraWide,
		},
		showWarnings: true,
	}
}

// UpdateSize updates the current terminal size and validates it
func (r *ResponsiveLayoutManager) UpdateSize(width, height int) {
	r.currentSize.Width = width
	r.currentSize.Height = height
	r.activeBreakpoint = r.getActiveBreakpoint()
	r.updateSizeWarnings()
}

// GetCurrentSize returns the current terminal size
func (r *ResponsiveLayoutManager) GetCurrentSize() TerminalSize {
	return r.currentSize
}

// GetActiveBreakpoint returns the currently active breakpoint
func (r *ResponsiveLayoutManager) GetActiveBreakpoint() ResponsiveBreakpoint {
	return r.activeBreakpoint
}

// IsMinimumSize checks if current size meets minimum requirements
func (r *ResponsiveLayoutManager) IsMinimumSize() bool {
	return r.currentSize.Width >= r.minSize.Width && r.currentSize.Height >= r.minSize.Height
}

// IsOptimalSize checks if current size is optimal for the best experience
func (r *ResponsiveLayoutManager) IsOptimalSize() bool {
	return r.currentSize.Width >= 100 && r.currentSize.Height >= 30
}

// GetSizeWarnings returns current size warnings
func (r *ResponsiveLayoutManager) GetSizeWarnings() []string {
	if !r.showWarnings {
		return nil
	}
	return r.sizeWarnings
}

// ShouldShowCompactLayout determines if compact layout should be used
func (r *ResponsiveLayoutManager) ShouldShowCompactLayout() bool {
	return r.activeBreakpoint == BreakpointUltraCompact || r.activeBreakpoint == BreakpointCompact
}

// ShouldCollapsePanes determines if panes should be collapsed for very small terminals
func (r *ResponsiveLayoutManager) ShouldCollapsePanes() bool {
	return r.currentSize.Width < 100
}

// GetRecommendedSplitRatio returns the recommended split ratio for current size
func (r *ResponsiveLayoutManager) GetRecommendedSplitRatio() float64 {
	switch r.activeBreakpoint {
	case BreakpointUltraCompact:
		return 0.7 // Favor editor pane more in very small terminals
	case BreakpointCompact:
		return 0.6
	case BreakpointMedium:
		return 0.5
	case BreakpointLarge:
		return 0.45
	case BreakpointUltraWide:
		return 0.4
	default:
		return 0.5
	}
}

// GetStatusBarMode returns the appropriate status bar mode for current size
func (r *ResponsiveLayoutManager) GetStatusBarMode() StatusBarMode {
	if r.currentSize.Width < 100 {
		return StatusBarCompact
	}
	return StatusBarFull
}

// GetMenuMode returns the appropriate menu mode for current size
func (r *ResponsiveLayoutManager) GetMenuMode() MenuMode {
	if r.currentSize.Width < 90 {
		return MenuCompact
	}
	return MenuFull
}

// getActiveBreakpoint determines which breakpoint is currently active
func (r *ResponsiveLayoutManager) getActiveBreakpoint() ResponsiveBreakpoint {
	for _, bp := range r.breakpoints {
		if r.currentSize.Width >= bp.MinWidth && r.currentSize.Width <= bp.MaxWidth &&
			r.currentSize.Height >= bp.MinHeight && r.currentSize.Height <= bp.MaxHeight {
			return bp
		}
	}
	// Default to compact if no breakpoint matches
	return BreakpointCompact
}

// updateSizeWarnings updates size-related warnings
func (r *ResponsiveLayoutManager) updateSizeWarnings() {
	r.sizeWarnings = nil

	if !r.IsMinimumSize() {
		r.sizeWarnings = append(r.sizeWarnings,
			fmt.Sprintf("Terminal too small! Minimum size is %dx%d, current is %dx%d",
				r.minSize.Width, r.minSize.Height, r.currentSize.Width, r.currentSize.Height))
	}

	if r.currentSize.Width < 100 || r.currentSize.Height < 30 {
		r.sizeWarnings = append(r.sizeWarnings,
			"For optimal experience, use at least 100x30 terminal")
	}

	if r.currentSize.Width > 200 {
		r.sizeWarnings = append(r.sizeWarnings,
			"Ultra-wide terminal detected. Consider using multiple columns for better layout")
	}
}

// RenderSizeWarning renders size warnings for display
func (r *ResponsiveLayoutManager) RenderSizeWarning() string {
	if len(r.sizeWarnings) == 0 {
		return ""
	}

	warningStyle := lipgloss.NewStyle().
		Foreground(styles.Warning).
		Background(styles.Dark1).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Warning)

	warningText := strings.Join(r.sizeWarnings, "\n")
	return warningStyle.Render("⚠️  " + warningText)
}

// HandleWindowSizeMsg handles window size messages and returns appropriate commands
func (r *ResponsiveLayoutManager) HandleWindowSizeMsg(msg tea.WindowSizeMsg) tea.Cmd {
	r.UpdateSize(msg.Width, msg.Height)

	// Return nil if size is acceptable, or a command to show warning
	if !r.IsMinimumSize() {
		return func() tea.Msg {
			return SizeValidationMsg{
				IsValid:  false,
				Width:    msg.Width,
				Height:   msg.Height,
				Warnings: r.GetSizeWarnings(),
			}
		}
	}

	return func() tea.Msg {
		return SizeValidationMsg{
			IsValid: true,
			Width:   msg.Width,
			Height:  msg.Height,
		}
	}
}

// SizeValidationMsg represents a size validation message
type SizeValidationMsg struct {
	IsValid  bool
	Width    int
	Height   int
	Warnings []string
}

// StatusBarMode represents different status bar display modes
type StatusBarMode int

const (
	StatusBarFull StatusBarMode = iota
	StatusBarCompact
	StatusBarMinimal
)

// MenuMode represents different menu display modes
type MenuMode int

const (
	MenuFull MenuMode = iota
	MenuCompact
	MenuScrollable
)

// Global responsive layout manager instance
var GlobalResponsiveManager = NewResponsiveLayoutManager()

// GetMinimumResolutionOptimizations returns specific optimizations for 80x24 terminals
func (r *ResponsiveLayoutManager) GetMinimumResolutionOptimizations() map[string]interface{} {
	if !r.IsMinimumSize() {
		return nil
	}

	optimizations := make(map[string]interface{})

	// At minimum resolution, prioritize essential features
	optimizations["hide_line_numbers"] = r.currentSize.Width < 90
	optimizations["minimal_status_bar"] = r.currentSize.Width < 85
	optimizations["compact_menus"] = r.currentSize.Width < 90
	optimizations["essential_help_only"] = r.currentSize.Width < 85
	optimizations["reduced_padding"] = true
	optimizations["minimal_animations"] = r.currentSize.Width < 85

	return optimizations
}

// IsUltraCompact returns true if terminal is in ultra-compact mode
func (r *ResponsiveLayoutManager) IsUltraCompact() bool {
	return r.activeBreakpoint == BreakpointUltraCompact
}

// GetOptimalDimensions returns recommended dimensions for the best experience
func (r *ResponsiveLayoutManager) GetOptimalDimensions() TerminalSize {
	return TerminalSize{Width: 100, Height: 30}
}

// GetLayoutEfficiency returns a score (0-100) for layout efficiency at current size
func (r *ResponsiveLayoutManager) GetLayoutEfficiency() int {
	if !r.IsMinimumSize() {
		return 0
	}

	baseScore := 50 // Base score for meeting minimum requirements

	// Bonus points for optimal size
	if r.IsOptimalSize() {
		baseScore += 30
	}

	// Bonus points for larger terminals
	if r.currentSize.Width >= 120 {
		baseScore += 10
	}
	if r.currentSize.Height >= 40 {
		baseScore += 10
	}

	// Penalty for very small terminals
	if r.currentSize.Width < 90 || r.currentSize.Height < 28 {
		baseScore -= 10
	}

	// Cap at 100
	if baseScore > 100 {
		baseScore = 100
	}

	return baseScore
}
