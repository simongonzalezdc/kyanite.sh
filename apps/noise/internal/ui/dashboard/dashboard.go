package dashboard

import (
	"fmt"

	"github.com/kyanite/noise/internal/app"
	"github.com/kyanite/noise/internal/infra/db"
	"github.com/kyanite/noise/internal/theme"
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
	focusedPanel  string
	showAllPanels bool // New field for progressive disclosure

	// Header and footer
	header *HeaderModel
	footer *FooterModel
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel() *DashboardModel {
	return &DashboardModel{
		currentTheme:  theme.GetManager().Current(),
		themeManager:  NewThemePreviewModel(),
		quickActions:  NewQuickActionsModel(),
		recentWork:    NewRecentWorkModel(),
		musicTools:    NewMusicToolsModel(),
		aiAssistant:   NewAIAssistantModel(),
		systemInfo:    NewSystemInfoModel(),
		focusedPanel:  "quickactions",
		showAllPanels: false, // Start with simplified view
		header:        NewHeaderModel(),
		footer:        NewFooterModel(),
	}
}

// Init initializes the dashboard model
func (dm *DashboardModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Initialize recent work panel (loads data from database)
	if dm.recentWork != nil {
		if cmd := dm.recentWork.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Initialize system info panel
	if dm.systemInfo != nil {
		if cmd := dm.systemInfo.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// SetAIService passes the AI service to the dashboard's AI assistant panel
func (dm *DashboardModel) SetAIService(aiService *app.AIService) {
	if dm.aiAssistant != nil {
		dm.aiAssistant.SetAIService(aiService)
	}
}

// SetDatabase passes the database to panels that need it
func (dm *DashboardModel) SetDatabase(database *db.DB) {
	if dm.recentWork != nil {
		dm.recentWork.SetDatabase(database)
	}
	if dm.systemInfo != nil {
		dm.systemInfo.SetDatabase(database)
	}
}

// Update handles messages for the dashboard
func (dm *DashboardModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		dm.width = msg.Width
		dm.height = msg.Height

		// Get layout config to determine panel dimensions
		config := GetLayoutConfig(msg.Width, msg.Height)

		// Calculate panel dimensions based on grid configuration
		panelWidth := msg.Width / config.GridCols
		panelHeight := (msg.Height - 4) / config.GridRows // Subtract header/footer space

		// Create panel-specific size message
		panelMsg := tea.WindowSizeMsg{Width: panelWidth, Height: panelHeight}

		// Update all child components with panel dimensions
		if dm.themeManager != nil {
			cmds = append(cmds, dm.themeManager.Update(panelMsg))
		}
		if dm.quickActions != nil {
			cmds = append(cmds, dm.quickActions.Update(panelMsg))
		}
		if dm.recentWork != nil {
			cmds = append(cmds, dm.recentWork.Update(panelMsg))
		}
		if dm.musicTools != nil {
			cmds = append(cmds, dm.musicTools.Update(panelMsg))
		}
		if dm.aiAssistant != nil {
			cmds = append(cmds, dm.aiAssistant.Update(panelMsg))
		}
		if dm.systemInfo != nil {
			cmds = append(cmds, dm.systemInfo.Update(panelMsg))
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

	case tea.MouseMsg:
		// Forward mouse events to interactive child components
		// This enables click, scroll, and hover interactions
		if dm.quickActions != nil {
			if cmd := dm.quickActions.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		if dm.aiAssistant != nil {
			if cmd := dm.aiAssistant.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		if dm.themeManager != nil {
			if cmd := dm.themeManager.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		if dm.recentWork != nil {
			if cmd := dm.recentWork.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		if dm.musicTools != nil {
			if cmd := dm.musicTools.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

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

	// Header (fixed height: 1 line)
	header := dm.header.View()

	// Footer (fixed height: 1 line)
	footer := dm.footer.View()

	// Calculate available height for content (subtract header + footer + margins)
	contentHeight := dm.height - 4 // 1 header + 1 footer + 2 margin lines

	// Main content based on layout and progressive disclosure
	var content string
	if dm.showAllPanels {
		// Full dashboard view
		switch config.GridCols {
		case 6:
			// Ultra-wide: all 6 panels in a single row
			content = lipgloss.JoinHorizontal(lipgloss.Top,
				dm.themeManager.View(),
				dm.quickActions.View(),
				dm.recentWork.View(),
				dm.musicTools.View(),
				dm.aiAssistant.View(),
				dm.systemInfo.View(),
			)

		case 4:
			// Wide: 4 columns, 2 rows
			topRow := lipgloss.JoinHorizontal(lipgloss.Top,
				dm.themeManager.View(),
				dm.quickActions.View(),
				dm.recentWork.View(),
				dm.musicTools.View(),
			)
			bottomRow := lipgloss.JoinHorizontal(lipgloss.Top,
				dm.aiAssistant.View(),
				dm.systemInfo.View(),
			)
			content = lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)

		case 3:
			// Large: 3 columns, 2 rows
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
			// Medium: 2 columns, 3 rows
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

	// Constrain content to available height
	contentStyle := lipgloss.NewStyle().
		MaxHeight(contentHeight).
		MaxWidth(dm.width)
	content = contentStyle.Render(content)

	// Join all sections
	output := lipgloss.JoinVertical(lipgloss.Left, header, content, footer)

	// CRITICAL: Constrain final output to exact terminal dimensions
	// This prevents overflow that causes rendering artifacts
	finalStyle := lipgloss.NewStyle().
		Width(dm.width).
		Height(dm.height).
		MaxWidth(dm.width).
		MaxHeight(dm.height)

	return finalStyle.Render(output)
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

	// Calculate panel dimensions
	panelWidth := dm.width / 2
	panelHeight := dm.height - 6 // Leave room for header, footer, hint

	// Apply constraints to panels
	themeView := lipgloss.NewStyle().
		MaxWidth(panelWidth).
		MaxHeight(panelHeight).
		Render(dm.themeManager.View())

	actionsView := lipgloss.NewStyle().
		MaxWidth(panelWidth).
		MaxHeight(panelHeight).
		Render(dm.quickActions.View())

	// Show only quick actions and theme preview in simplified view
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, themeView, actionsView)

	// Add hint for expanding the view
	expandHint := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Align(lipgloss.Center).
		MaxWidth(dm.width).
		Render("Press [D] to expand dashboard panels")

	return lipgloss.JoinVertical(lipgloss.Left, topRow, "", expandHint)
}

// GetOnboardingHints returns contextual hints for new users
func (dm *DashboardModel) GetOnboardingHints() []string {
	hints := []string{
		"[TIP] Use number keys (1-6) to quickly access actions",
		"[TIP] Press 'T' to cycle through beautiful themes",
		"[TIP] Press 'H' or '?' for help and keyboard shortcuts",
		"[TIP] Press 'D' to toggle between simple and detailed views",
	}

	// Show different hints based on current state
	if !dm.showAllPanels {
		hints = append(hints, "[TIP] Press 'D' to see all dashboard panels")
	} else {
		hints = append(hints, "[TIP] Press 'D' to simplify the dashboard")
	}

	return hints
}

// renderMinimalMenu renders a minimal menu for very small terminals
func (dm *DashboardModel) renderMinimalMenu() string {
	return dm.quickActions.View()
}
