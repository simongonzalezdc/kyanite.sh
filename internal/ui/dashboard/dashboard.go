package dashboard

import (
	"fmt"
	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ForceRefreshMsg triggers a complete re-render
type ForceRefreshMsg struct{}

// AnimationTickMsg is sent when animations need to be updated
type AnimationTickMsg struct{}

// ScreenChangeMsg represents a screen change message
type ScreenChangeMsg struct {
	Screen int
}

// DashboardModel represents the main dashboard
type DashboardModel struct {
	// Layout and sizing
	width  int
	height int

	// Theme management
	currentTheme theme.Theme
	themeManager *ThemePreviewModel

	// Component panels
	quickActions *QuickActionsModel
	recentWork   *RecentWorkModel
	musicTools   *MusicToolsModel
	aiAssistant  *AIAssistantModel
	systemInfo   *SystemInfoModel

	// UI state
	focusedPanel string
	showAllPanels bool // New field for progressive disclosure

	// Header and footer
	header *HeaderModel
	footer *FooterModel
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel() *DashboardModel {
	return &DashboardModel{
		currentTheme: theme.GetManager().Current(),
		themeManager: NewThemePreviewModel(),
		quickActions: NewQuickActionsModel(),
		recentWork:   NewRecentWorkModel(),
		musicTools:   NewMusicToolsModel(),
		aiAssistant:  NewAIAssistantModel(),
		systemInfo:   NewSystemInfoModel(),
		focusedPanel: "quickactions",
		showAllPanels: false, // Start with simplified view
		header:       NewHeaderModel(),
		footer:       NewFooterModel(),
	}
}

// Init initializes the dashboard model
func (dm *DashboardModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the dashboard
func (dm *DashboardModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		dm.width = msg.Width
		dm.height = msg.Height

		// Update all child components
		if dm.themeManager != nil {
			cmds = append(cmds, dm.themeManager.Update(msg))
		}
		if dm.quickActions != nil {
			cmds = append(cmds, dm.quickActions.Update(msg))
		}
		if dm.recentWork != nil {
			cmds = append(cmds, dm.recentWork.Update(msg))
		}
		if dm.musicTools != nil {
			cmds = append(cmds, dm.musicTools.Update(msg))
		}
		if dm.aiAssistant != nil {
			cmds = append(cmds, dm.aiAssistant.Update(msg))
		}
		if dm.systemInfo != nil {
			cmds = append(cmds, dm.systemInfo.Update(msg))
		}

		// CRITICAL: Force refresh after resize
		cmds = append(cmds, func() tea.Msg { return ForceRefreshMsg{} })

	case AnimationTickMsg:
		// Update all components
		if dm.themeManager != nil {
			cmds = append(cmds, dm.themeManager.Update(msg))
		}
		if dm.quickActions != nil {
			cmds = append(cmds, dm.quickActions.Update(msg))
		}
		if dm.recentWork != nil {
			cmds = append(cmds, dm.recentWork.Update(msg))
		}
		if dm.musicTools != nil {
			cmds = append(cmds, dm.musicTools.Update(msg))
		}
		if dm.aiAssistant != nil {
			cmds = append(cmds, dm.aiAssistant.Update(msg))
		}
		if dm.systemInfo != nil {
			cmds = append(cmds, dm.systemInfo.Update(msg))
		}

	case ForceRefreshMsg:
		// Handle force refresh
		return nil

	case tea.KeyMsg:
		switch msg.String() {
		case "1", "2", "3", "4", "5", "6":
			// Handle quick action shortcuts
			if dm.quickActions != nil {
				cmd := dm.quickActions.HandleKey(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		case "t", "T":
			// Switch theme
			if dm.themeManager != nil {
				dm.themeManager.NextTheme()
				// Update current theme reference
				dm.currentTheme = theme.GetManager().Current()
				// Trigger refresh
				cmds = append(cmds, func() tea.Msg { return ForceRefreshMsg{} })
			}
		case "tab":
			// Navigate between panels
			dm.focusNextPanel()
		case "d", "D": // New key binding for toggling panel visibility
			// Toggle between simplified and full dashboard views
			dm.showAllPanels = !dm.showAllPanels
		case "esc":
			// Return to menu
			return func() tea.Msg {
				return ScreenChangeMsg{Screen: 1} // screenMenu
			}
		}
	}

	return tea.Batch(cmds...)
}

// View renders the dashboard
func (dm *DashboardModel) View() string {
	if dm.width == 0 || dm.height == 0 {
		return fmt.Sprintf("Dashboard loading... (width: %d, height: %d)", dm.width, dm.height)
	}

	// Get layout configuration based on terminal size
	config := GetLayoutConfig(dm.width, dm.height)

	// Header
	header := dm.header.View()

	// Main content based on layout and progressive disclosure
	var content string
	if dm.showAllPanels {
		// Full dashboard view
		switch config.GridCols {
		case 3:
			// Full 3x3 grid
			topRow := lipgloss.JoinHorizontal(lipgloss.Top,
				dm.themeManager.View(),
				dm.quickActions.View(),
				dm.recentWork.View(),
			)
			middleRow := lipgloss.JoinHorizontal(lipgloss.Top,
				dm.musicTools.View(),
				dm.aiAssistant.View(),
				dm.systemInfo.View(),
			)
			content = lipgloss.JoinVertical(lipgloss.Left, topRow, middleRow)

		case 2:
			// 2x3 grid
			topRow := lipgloss.JoinHorizontal(lipgloss.Top,
				dm.themeManager.View(),
				dm.quickActions.View(),
			)
			middleRow := lipgloss.JoinHorizontal(lipgloss.Top,
				dm.recentWork.View(),
				dm.musicTools.View(),
			)
			bottomRow := lipgloss.JoinHorizontal(lipgloss.Top,
				dm.aiAssistant.View(),
				dm.systemInfo.View(),
			)
			content = lipgloss.JoinVertical(lipgloss.Left, topRow, middleRow, bottomRow)

		case 1:
			// Single column with active panel
			content = dm.renderActivePanel()

		default:
			// Minimal menu
			content = dm.renderMinimalMenu()
		}
	} else {
		// Simplified dashboard view with only essential panels
		content = dm.renderSimplifiedView()
	}

	// Footer
	footer := dm.footer.View()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// focusNextPanel moves focus to the next panel
func (dm *DashboardModel) focusNextPanel() {
	panels := []string{"theme", "actions", "recent", "tools", "ai", "info"}
	for i, panel := range panels {
		if panel == dm.focusedPanel && i < len(panels)-1 {
			dm.focusedPanel = panels[i+1]
			return
		}
	}
	// If at the end or not found, go to the beginning
	dm.focusedPanel = panels[0]
}

// renderActivePanel renders the currently active panel for small terminals
func (dm *DashboardModel) renderActivePanel() string {
	switch dm.focusedPanel {
	case "theme":
		return dm.themeManager.View()
	case "actions":
		return dm.quickActions.View()
	case "recent":
		return dm.recentWork.View()
	case "tools":
		return dm.musicTools.View()
	case "ai":
		return dm.aiAssistant.View()
	case "info":
		return dm.systemInfo.View()
	default:
		return dm.quickActions.View()
	}
}

// renderSimplifiedView renders a simplified dashboard view
func (dm *DashboardModel) renderSimplifiedView() string {
	t := theme.GetManager().Current()

	// Show only quick actions and theme preview in simplified view
	topRow := lipgloss.JoinHorizontal(lipgloss.Top,
		dm.themeManager.View(),
		dm.quickActions.View(),
	)

	// Add hint for expanding the view
	expandHint := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Align(lipgloss.Center).
		Render("Press [D] to expand dashboard panels")

	return lipgloss.JoinVertical(lipgloss.Left, topRow, "", expandHint)
}

// GetOnboardingHints returns contextual hints for new users
func (dm *DashboardModel) GetOnboardingHints() []string {
	hints := []string{
		"🎯 Quick tip: Use number keys (1-6) to quickly access actions",
		"🎨 Press 'T' to cycle through beautiful themes",
		"📚 Press 'H' or '?' for help and keyboard shortcuts",
		"🔍 Press 'D' to toggle between simple and detailed views",
	}

	// Show different hints based on current state
	if !dm.showAllPanels {
		hints = append(hints, "💡 Press 'D' to see all dashboard panels")
	} else {
		hints = append(hints, "💡 Press 'D' to simplify the dashboard")
	}

	return hints
}


// renderMinimalMenu renders a minimal menu for very small terminals
func (dm *DashboardModel) renderMinimalMenu() string {
	return dm.quickActions.View()
}