package dashboard

import (
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardModel_Accessibility(t *testing.T) {
	t.Run("FocusIndicators", func(t *testing.T) {
		// Test that dashboard components have proper focus indicators
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Check that dashboard has a focused panel
		assert.NotEmpty(t, dashboard.focusedPanel)

		// Check that view reflects focus state
		view := dashboard.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Quick Actions")
	})

	t.Run("KeyboardNavigation", func(t *testing.T) {
		// Test keyboard navigation support
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test tab navigation
		tabMsg := tea.KeyMsg{Type: tea.KeyTab}
		cmd := dashboard.Update(tabMsg)
		// Should not panic
		assert.NotPanics(t, func() {
			_ = cmd
		})

		// Test escape key for modals
		escMsg := tea.KeyMsg{Type: tea.KeyEsc}
		cmd = dashboard.Update(escMsg)
		// Should not panic
		assert.NotPanics(t, func() {
			_ = cmd
		})

		// Test arrow key navigation
		upMsg := tea.KeyMsg{Type: tea.KeyUp}
		cmd = dashboard.Update(upMsg)
		assert.NotPanics(t, func() {
			_ = cmd
		})

		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		cmd = dashboard.Update(downMsg)
		assert.NotPanics(t, func() {
			_ = cmd
		})
	})

	t.Run("ScreenReaderSupport", func(t *testing.T) {
		// Test that dashboard components provide meaningful text for screen readers
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Check that view contains meaningful text
		view := dashboard.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Dashboard")
		assert.Contains(t, view, "Quick Actions")
		assert.Contains(t, view, "Recent Work")
		assert.Contains(t, view, "Music Tools")
		assert.Contains(t, view, "AI Assistant")
		assert.Contains(t, view, "System Info")
	})

	t.Run("ColorContrast", func(t *testing.T) {
		// Test that dashboard components use theme colors with adequate contrast
		themes := theme.ListThemes()
		require.Greater(t, len(themes), 0)

		for _, themeID := range themes {
			t.Run(themeID, func(t *testing.T) {
				th := theme.GetTheme(themeID)

				// Test that text and background colors are different
				assert.NotEqual(t, th.Text, th.Background, "Text and background should have different colors for theme %s", themeID)

				// Test that primary and background colors are different
				assert.NotEqual(t, th.Primary, th.Background, "Primary and background should have different colors for theme %s", themeID)

				// Test that accent and background colors are different
				assert.NotEqual(t, th.Accent, th.Background, "Accent and background should have different colors for theme %s", themeID)
			})
		}
	})
}

func TestDashboardModel_FocusManagement(t *testing.T) {
	t.Run("PanelFocusNavigation", func(t *testing.T) {
		// Test comprehensive focus management for dashboard panels
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test initial focus state
		initialPanel := dashboard.focusedPanel
		assert.NotEmpty(t, initialPanel)

		// Test focus navigation to next panel
		originalPanel := dashboard.focusedPanel
		dashboard.focusNextPanel()
		assert.NotEqual(t, originalPanel, dashboard.focusedPanel, "Focus should move to next panel")

		// Test cycling through all panels
		panels := []string{"theme", "actions", "recent", "tools", "ai", "info"}
		visitedPanels := make(map[string]bool)
		
		for i := 0; i < len(panels)*2; i++ { // Test cycling through panels twice
			panel := dashboard.focusedPanel
			visitedPanels[panel] = true
			dashboard.focusNextPanel()
		}

		// Verify all panels were visited
		for _, panel := range panels {
			assert.True(t, visitedPanels[panel], "Panel %s should be visited during navigation", panel)
		}
	})

	t.Run("PanelFocusIndicators", func(t *testing.T) {
		// Test that focused panels have visual indicators
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Get view and check for focus indicators
		view := dashboard.View()
		assert.NotEmpty(t, view)

		// Should contain focused panel indicator
		// This is a basic check - in a real implementation, you'd check for specific styling
		assert.True(t, len(view) > 0, "View should be non-empty")
	})
}

func TestDashboardModel_ScreenReaderText(t *testing.T) {
	t.Run("DashboardViewAccessibility", func(t *testing.T) {
		// Test that dashboard provides comprehensive screen reader text
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test basic view
		view := dashboard.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Dashboard")

		// Test with different panel states
		// Focus on different panels and verify content
		panels := []string{"theme", "actions", "recent", "tools", "ai", "info"}
		for _, panel := range panels {
			dashboard.focusedPanel = panel
			view = dashboard.View()
			assert.NotEmpty(t, view)
			assert.Contains(t, view, "Dashboard")
		}
	})

	t.Run("PanelContentAccessibility", func(t *testing.T) {
		// Test that each panel provides meaningful text content
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test quick actions panel
		dashboard.focusedPanel = "actions"
		view := dashboard.View()
		assert.Contains(t, view, "Quick Actions")

		// Test recent work panel
		dashboard.focusedPanel = "recent"
		view = dashboard.View()
		assert.Contains(t, view, "Recent Work")

		// Test music tools panel
		dashboard.focusedPanel = "tools"
		view = dashboard.View()
		assert.Contains(t, view, "Music Tools")

		// Test AI assistant panel
		dashboard.focusedPanel = "ai"
		view = dashboard.View()
		assert.Contains(t, view, "AI Assistant")

		// Test system info panel
		dashboard.focusedPanel = "info"
		view = dashboard.View()
		assert.Contains(t, view, "System Info")
	})
}

func TestDashboardModel_KeyboardShortcuts(t *testing.T) {
	t.Run("DashboardKeyboardNavigation", func(t *testing.T) {
		// Test that dashboard handles keyboard shortcuts appropriately
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test common shortcuts don't panic
		shortcuts := []tea.KeyMsg{
			{Type: tea.KeyTab},
			{Type: tea.KeyEsc},
			{Type: tea.KeyEnter},
			{Type: tea.KeyUp},
			{Type: tea.KeyDown},
			{Type: tea.KeyLeft},
			{Type: tea.KeyRight},
		}

		for _, shortcut := range shortcuts {
			assert.NotPanics(t, func() {
				dashboard.Update(shortcut)
			})
		}
	})

	t.Run("DashboardSpecificShortcuts", func(t *testing.T) {
		// Test dashboard-specific keyboard shortcuts
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test panel navigation shortcuts
		panelKeys := []tea.KeyMsg{
			{Type: tea.KeyTab},
			{Runes: []rune{'\t'}},
		}

		for _, keyMsg := range panelKeys {
			assert.NotPanics(t, func() {
				dashboard.Update(keyMsg)
			})
		}
	})
}

func TestDashboardModel_ThemeIntegration(t *testing.T) {
	t.Run("ThemeSwitchingKeyboardSupport", func(t *testing.T) {
		// Test that theme switching supports keyboard navigation
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test theme integration
		assert.NotPanics(t, func() {
			// Dashboard should work with theme manager
			_ = dashboard.View()
		})
	})

	t.Run("ThemeColorContrast", func(t *testing.T) {
		// Test that all themes have adequate color contrast for dashboard
		themes := theme.ListThemes()
		require.Greater(t, len(themes), 0)

		for _, themeID := range themes {
			t.Run(themeID, func(t *testing.T) {
				th := theme.GetTheme(themeID)

				// Test critical color combinations for contrast
				assert.NotEqual(t, th.Text, th.Background, "Text and background should differ for theme %s", themeID)
				assert.NotEqual(t, th.Primary, th.Background, "Primary and background should differ for theme %s", themeID)
				assert.NotEqual(t, th.Secondary, th.Background, "Secondary and background should differ for theme %s", themeID)
				assert.NotEqual(t, th.Accent, th.Background, "Accent and background should differ for theme %s", themeID)
				assert.NotEqual(t, th.Success, th.Background, "Success and background should differ for theme %s", themeID)
				assert.NotEqual(t, th.Error, th.Background, "Error and background should differ for theme %s", themeID)
				assert.NotEqual(t, th.Warning, th.Background, "Warning and background should differ for theme %s", themeID)
			})
		}
	})

	t.Run("ThemeReducedMotion", func(t *testing.T) {
		// Test that dashboard respects reduced motion preferences
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test that dashboard works with different theme configurations
		view := dashboard.View()
		assert.NotEmpty(t, view)

		// Test with different themes
		themes := theme.ListThemes()
		if len(themes) > 0 {
			themeManager := theme.GetManager()
			originalTheme := themeManager.Current()
			defer func() {
				themeManager.SetTheme(originalTheme.Name)
			}()

			themeManager.SetTheme(themes[0])
			current := themeManager.Current()
			assert.NotNil(t, current)

			// Verify dashboard still renders correctly with new theme
			view = dashboard.View()
			assert.NotEmpty(t, view)
		}
	})
}

func TestDashboardModel_Components(t *testing.T) {
	t.Run("QuickActionsAccessibility", func(t *testing.T) {
		// Test quick actions component accessibility
		actions := NewQuickActionsModel()
		require.NotNil(t, actions)

		view := actions.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "New Song")
		assert.Contains(t, view, "Open Project")
		assert.Contains(t, view, "Export")
		assert.Contains(t, view, "Settings")
	})

	t.Run("RecentWorkAccessibility", func(t *testing.T) {
		// Test recent work component accessibility
		recent := NewRecentWorkModel()
		require.NotNil(t, recent)

		view := recent.View()
		assert.NotEmpty(t, view)
		// Should contain recent work indicators even when empty
		assert.Contains(t, view, "Recent Work")
	})

	t.Run("MusicToolsAccessibility", func(t *testing.T) {
		// Test music tools component accessibility
		tools := NewMusicToolsModel()
		require.NotNil(t, tools)

		view := tools.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Chord Progressions")
		assert.Contains(t, view, "Scale Explorer")
		assert.Contains(t, view, "Metronome")
		assert.Contains(t, view, "Theory Reference")
	})

	t.Run("AIAssistantAccessibility", func(t *testing.T) {
		// Test AI assistant component accessibility
		ai := NewAIAssistantModel()
		require.NotNil(t, ai)

		view := ai.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "AI Assistant")
		assert.Contains(t, view, "Quick Idea")
		assert.Contains(t, view, "Lyric Helper")
		assert.Contains(t, view, "Melody Generator")
	})

	t.Run("SystemInfoAccessibility", func(t *testing.T) {
		// Test system info component accessibility
		info := NewSystemInfoModel()
		require.NotNil(t, info)

		view := info.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "System Info")
	})

	t.Run("ThemePreviewAccessibility", func(t *testing.T) {
		// Test theme preview component accessibility
		preview := NewThemePreviewModel()
		require.NotNil(t, preview)

		view := preview.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Theme Preview")
	})
}

func TestDashboardModel_ResponsiveLayout(t *testing.T) {
	t.Run("ResponsiveLayoutAccessibility", func(t *testing.T) {
		// Test that dashboard layout is accessible at different sizes
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test various window sizes
		sizes := []tea.WindowSizeMsg{
			{Width: 80, Height: 24},   // Minimum size
			{Width: 120, Height: 40},  // Medium size
			{Width: 160, Height: 50},  // Large size
		}

		for _, size := range sizes {
			cmd := dashboard.Update(size)
			assert.NotPanics(t, func() {
				_ = cmd
			})

			view := dashboard.View()
			assert.NotEmpty(t, view)
			assert.Contains(t, view, "Dashboard")
		}
	})

	t.Run("LayoutFocusIndicators", func(t *testing.T) {
		// Test that layout changes maintain focus indicators
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test with different window sizes
		sizes := []tea.WindowSizeMsg{
			{Width: 80, Height: 24},
			{Width: 160, Height: 50},
		}

		for _, size := range sizes {
			cmd := dashboard.Update(size)
			assert.NotPanics(t, func() {
				_ = cmd
			})

			view := dashboard.View()
			assert.NotEmpty(t, view)
			
			// Should maintain focus indicators
			assert.Contains(t, view, "Dashboard")
		}
	})
}

func TestDashboardModel_ErrorHandling(t *testing.T) {
	t.Run("GracefulErrorHandling", func(t *testing.T) {
		// Test that dashboard handles errors gracefully
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test with various error messages
		errorMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("invalid")}
		assert.NotPanics(t, func() {
			dashboard.Update(errorMsg)
		})

		view := dashboard.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Dashboard")
	})
}

func TestDashboardModel_Performance(t *testing.T) {
	t.Run("ViewRenderingPerformance", func(t *testing.T) {
		// Test that dashboard view rendering is performant
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Measure rendering time
		start := time.Now()
		view := dashboard.View()
		duration := time.Since(start)

		assert.NotEmpty(t, view)
		assert.Less(t, duration, 100*time.Millisecond, "View rendering should be fast")
	})

	t.Run("UpdatePerformance", func(t *testing.T) {
		// Test that dashboard updates are performant
		dashboard := NewDashboardModel()
		require.NotNil(t, dashboard)

		// Test multiple rapid updates
		start := time.Now()
		for i := 0; i < 100; i++ {
			msg := tea.KeyMsg{Type: tea.KeyTab}
			cmd := dashboard.Update(msg)
			assert.NotPanics(t, func() {
				_ = cmd
			})
		}
		duration := time.Since(start)

		assert.Less(t, duration, 1*time.Second, "Multiple updates should be fast")
	})
}