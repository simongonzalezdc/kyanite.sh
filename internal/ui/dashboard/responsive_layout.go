package dashboard

import (
	"github.com/Kyanite/noise/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// LayoutConfig defines the layout configuration for different screen sizes
type LayoutConfig struct {
	PanelWidths  map[string]int
	PanelHeights map[string]int
	GridCols     int
	GridRows     int
	ShowPanels   map[string]bool
}

// GetLayoutConfig returns the appropriate layout configuration for the given screen size
func GetLayoutConfig(width, height int) LayoutConfig {
	switch {
	case width >= 240 && height >= 30:
		// Ultra-wide terminals (240+) - full 6-column layout
		return LayoutConfig{
			PanelWidths: map[string]int{
				"theme":   width / 6,
				"actions": width / 6,
				"recent":  width / 6,
				"tools":   width / 6,
				"ai":      width / 6,
				"info":    width / 6,
			},
			PanelHeights: map[string]int{
				"theme":   height - 4,
				"actions": height - 4,
				"recent":  height - 4,
				"tools":   height - 4,
				"ai":      height - 4,
				"info":    height - 4,
			},
			GridCols: 6,
			GridRows: 1,
			ShowPanels: map[string]bool{
				"theme": true, "actions": true, "recent": true,
				"tools": true, "ai": true, "info": true,
			},
		}

	case width >= 180 && height >= 30:
		// Wide terminals (180-239) - 4-column layout
		return LayoutConfig{
			PanelWidths: map[string]int{
				"theme":   width / 4,
				"actions": width / 4,
				"recent":  width / 4,
				"tools":   width / 4,
				"ai":      width / 4,
				"info":    width / 4,
			},
			PanelHeights: map[string]int{
				"theme":   height / 2,
				"actions": height / 2,
				"recent":  height / 2,
				"tools":   height / 2,
				"ai":      height / 2,
				"info":    height / 2,
			},
			GridCols: 4,
			GridRows: 2,
			ShowPanels: map[string]bool{
				"theme": true, "actions": true, "recent": true,
				"tools": true, "ai": true, "info": true,
			},
		}

	case width >= 120 && height >= 30:
		// Large terminals - full 3-column grid
		return LayoutConfig{
			PanelWidths: map[string]int{
				"theme":   width / 3,
				"actions": width / 3,
				"recent":  width / 3,
				"tools":   width / 3,
				"ai":      width / 3,
				"info":    width / 3,
			},
			PanelHeights: map[string]int{
				"theme":   height / 2,
				"actions": height / 2,
				"recent":  height / 2,
				"tools":   height / 2,
				"ai":      height / 2,
				"info":    height / 2,
			},
			GridCols: 3,
			GridRows: 2,
			ShowPanels: map[string]bool{
				"theme": true, "actions": true, "recent": true,
				"tools": true, "ai": true, "info": true,
			},
		}

	case width >= 80 && height >= 24:
		// Medium terminals - 2x3 grid
		return LayoutConfig{
			PanelWidths: map[string]int{
				"theme":   width / 2,
				"actions": width / 2,
				"recent":  width / 2,
				"tools":   width / 2,
				"ai":      width / 2,
				"info":    width / 2,
			},
			PanelHeights: map[string]int{
				"theme":   height / 3,
				"actions": height / 3,
				"recent":  height / 3,
				"tools":   height / 3,
				"ai":      height / 3,
				"info":    height / 3,
			},
			GridCols: 2,
			GridRows: 3,
			ShowPanels: map[string]bool{
				"theme": true, "actions": true, "recent": true,
				"tools": true, "ai": true, "info": true,
			},
		}

	case width >= 60 && height >= 20:
		// Small terminals - single column with tabs
		return LayoutConfig{
			PanelWidths: map[string]int{
				"active": width - 4,
			},
			PanelHeights: map[string]int{
				"active": height - 8,
			},
			GridCols: 1,
			GridRows: 1,
			ShowPanels: map[string]bool{
				"active": true,
			},
		}

	default:
		// Minimal terminals - essential features only
		return LayoutConfig{
			PanelWidths: map[string]int{
				"menu": width - 2,
			},
			PanelHeights: map[string]int{
				"menu": height - 4,
			},
			GridCols: 1,
			GridRows: 1,
			ShowPanels: map[string]bool{
				"menu": true,
			},
		}
	}
}

// HeaderModel represents the dashboard header
type HeaderModel struct {
	width  int
	height int
}

// NewHeaderModel creates a new header model
func NewHeaderModel() *HeaderModel {
	return &HeaderModel{}
}

// View renders the header
func (hm *HeaderModel) View() string {
	title := "[~] noise.sh"
	themeName := theme.GetManager().Current().Name
	status := "[*] 3 Songs, 2 Projects"
	help := "[F1] Help"
	quit := "[Q]uit"
	
	return lipgloss.JoinHorizontal(lipgloss.Left,
		title,
		"    ",
		themeName,
		"    ",
		status,
		"    ",
		help,
		"    ",
		quit,
	)
}

// FooterModel represents the dashboard footer
type FooterModel struct {
	width  int
	height int
}

// NewFooterModel creates a new footer model
func NewFooterModel() *FooterModel {
	return &FooterModel{}
}

// View renders the footer
func (fm *FooterModel) View() string {
	actions := "[1-6] Quick Actions"
	navigate := "[Tab] Navigate Panels"
	theme := "[T] Switch Theme"
	help := "[H] Help"
	menu := "[ESC] Menu"
	time := "23:45"
	version := "v1.0.0"
	
	return lipgloss.JoinHorizontal(lipgloss.Left,
		actions,
		"    ",
		navigate,
		"    ",
		theme,
		"    ",
		help,
		"    ",
		menu,
		"    ",
		time,
		"    ",
		version,
	)
}