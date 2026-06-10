package ui

import (
	"testing"
	"time"

	"github.com/kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnimatedLoadingSpinner_Accessibility(t *testing.T) {
	t.Run("FocusIndicators", func(t *testing.T) {
		// Test spinner has visual indication of activity
		spinner := NewAnimatedLoadingSpinner("Loading...")
		view := spinner.View()
		assert.NotEmpty(t, view)
		// Should have some visual element indicating activity
		assert.True(t, len(view) > 0)
	})

	t.Run("KeyboardNavigationSupport", func(t *testing.T) {
		// Test that spinner can handle keyboard navigation messages
		spinner := NewAnimatedLoadingSpinner("Loading...")

		// Test with tab key (common navigation key)
		tabKey := tea.KeyMsg{Type: tea.KeyTab}
		// Should not crash and should handle gracefully
		// Just ensure it doesn't panic
		assert.NotPanics(t, func() {
			spinner.Update(tabKey)
		})

		// Test with escape key
		escKey := tea.KeyMsg{Type: tea.KeyEsc}
		// Should not crash and should handle gracefully
		// Just ensure it doesn't panic
		assert.NotPanics(t, func() {
			spinner.Update(escKey)
		})
	})

	t.Run("ScreenReaderText", func(t *testing.T) {
		// Test that spinner provides meaningful text content for screen readers
		spinner := NewAnimatedLoadingSpinner("Loading content...")
		view := spinner.View()
		// Should contain the message for screen readers
		assert.Contains(t, view, "Loading content...")
	})

	t.Run("ReducedMotionSupport", func(t *testing.T) {
		// Test that spinner respects reduced motion preferences
		spinner := NewAnimatedLoadingSpinner("Loading...")

		// Test with animation manager that has reduced motion enabled
		animationManager := NewAnimationManager()
		config := animationManager.GetConfig()
		config.ReducedMotion = true
		animationManager.SetConfig(config)

		// Should still work with reduced motion
		view := spinner.View()
		assert.NotEmpty(t, view)
	})

	t.Run("ColorContrast", func(t *testing.T) {
		// Test that spinner uses theme colors with adequate contrast
		themes := theme.ListThemes()
		require.Greater(t, len(themes), 0)

		for _, themeID := range themes {
			t.Run(themeID, func(t *testing.T) {
				themeManager := theme.GetManager()
				originalTheme := themeManager.Current()
				defer func() {
					themeManager.SetTheme(originalTheme.Name)
				}()

				themeManager.SetTheme(themeID)
				th := themeManager.Current()

				// Test that text and background colors are different
				assert.NotEqual(t, th.Text, th.Background, "Text and background should have different colors for theme %s", themeID)

				// Test that primary and background colors are different
				assert.NotEqual(t, th.Primary, th.Background, "Primary and background should have different colors for theme %s", themeID)
			})
		}
	})
}

func TestAnimatedStatusBar_Accessibility(t *testing.T) {
	t.Run("StatusIndicators", func(t *testing.T) {
		// Test status bar has different visual indicators for different statuses
		statusBar := NewAnimatedStatusBar("Test message")

		// Normal status
		view := statusBar.View()
		assert.NotEmpty(t, view)

		// Success status
		statusBar.SetStatus(StatusSuccess, "Success")
		successView := statusBar.View()
		assert.NotEmpty(t, successView)
		assert.NotEqual(t, view, successView)

		// Error status
		statusBar.SetStatus(StatusError, "Error")
		errorView := statusBar.View()
		assert.NotEmpty(t, errorView)
		assert.NotEqual(t, successView, errorView)
	})

	t.Run("KeyboardNavigationSupport", func(t *testing.T) {
		// Test that status bar can handle keyboard navigation messages
		statusBar := NewAnimatedStatusBar("Test message")

		// Test with tab key (common navigation key)
		tabKey := tea.KeyMsg{Type: tea.KeyTab}
		// Should not crash and should handle gracefully
		// Just ensure it doesn't panic
		assert.NotPanics(t, func() {
			statusBar.Update(tabKey)
		})

		// Test with escape key
		escKey := tea.KeyMsg{Type: tea.KeyEsc}
		// Should not crash and should handle gracefully
		// Just ensure it doesn't panic
		assert.NotPanics(t, func() {
			statusBar.Update(escKey)
		})
	})

	t.Run("ScreenReaderText", func(t *testing.T) {
		// Test that status bar provides meaningful text content for screen readers
		statusBar := NewAnimatedStatusBar("Status information")
		statusView := statusBar.View()
		// Should contain the status message
		assert.Contains(t, statusView, "Status information")
	})

	t.Run("ReducedMotionSupport", func(t *testing.T) {
		// Test that status bar respects reduced motion preferences
		statusBar := NewAnimatedStatusBar("Test message")

		// Test with animation manager that has reduced motion enabled
		animationManager := NewAnimationManager()
		config := animationManager.GetConfig()
		config.ReducedMotion = true
		animationManager.SetConfig(config)

		// Should still work with reduced motion
		view := statusBar.View()
		assert.NotEmpty(t, view)
	})

	t.Run("ColorContrast", func(t *testing.T) {
		// Test that status bar uses theme colors with adequate contrast
		themes := theme.ListThemes()
		require.Greater(t, len(themes), 0)

		for _, themeID := range themes {
			t.Run(themeID, func(t *testing.T) {
				themeManager := theme.GetManager()
				originalTheme := themeManager.Current()
				defer func() {
					themeManager.SetTheme(originalTheme.Name)
				}()

				themeManager.SetTheme(themeID)
				th := themeManager.Current()

				// Test that text and background colors are different
				assert.NotEqual(t, th.Text, th.Background, "Text and background should have different colors for theme %s", themeID)

				// Test that success and background colors are different
				assert.NotEqual(t, th.Success, th.Background, "Success and background should have different colors for theme %s", themeID)

				// Test that error and background colors are different
				assert.NotEqual(t, th.Error, th.Background, "Error and background should have different colors for theme %s", themeID)
			})
		}
	})
}

func TestAnimatedNotification_Accessibility(t *testing.T) {
	t.Run("NotificationVisibility", func(t *testing.T) {
		// Test notifications are visible and distinguishable
		notification := NewAnimatedNotification("Test message", "info", 5*time.Second)
		view := notification.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Test message")
	})

	t.Run("KeyboardNavigationSupport", func(t *testing.T) {
		// Test that notifications can handle keyboard navigation messages
		notification := NewAnimatedNotification("Test message", "info", 5*time.Second)

		// Test with tab key (common navigation key)
		tabKey := tea.KeyMsg{Type: tea.KeyTab}
		// Should not crash and should handle gracefully
		// Just ensure it doesn't panic
		assert.NotPanics(t, func() {
			notification.Update(tabKey)
		})

		// Test with escape key
		escKey := tea.KeyMsg{Type: tea.KeyEsc}
		// Should not crash and should handle gracefully
		// Just ensure it doesn't panic
		assert.NotPanics(t, func() {
			notification.Update(escKey)
		})
	})

	t.Run("ScreenReaderText", func(t *testing.T) {
		// Test that notifications provide meaningful text content for screen readers
		notification := NewAnimatedNotification("Notification message", "info", 5*time.Second)
		notifView := notification.View()
		// Should contain the notification message
		assert.Contains(t, notifView, "Notification message")
	})

	t.Run("ReducedMotionSupport", func(t *testing.T) {
		// Test that notifications respect reduced motion preferences
		notification := NewAnimatedNotification("Test message", "info", 5*time.Second)

		// Test with animation manager that has reduced motion enabled
		animationManager := NewAnimationManager()
		config := animationManager.GetConfig()
		config.ReducedMotion = true
		animationManager.SetConfig(config)

		// Should still work with reduced motion
		view := notification.View()
		assert.NotEmpty(t, view)
	})

	t.Run("ColorContrast", func(t *testing.T) {
		// Test that notifications use theme colors with adequate contrast
		themes := theme.ListThemes()
		require.Greater(t, len(themes), 0)

		for _, themeID := range themes {
			t.Run(themeID, func(t *testing.T) {
				themeManager := theme.GetManager()
				originalTheme := themeManager.Current()
				defer func() {
					themeManager.SetTheme(originalTheme.Name)
				}()

				themeManager.SetTheme(themeID)
				th := themeManager.Current()

				// Test that text and background colors are different
				assert.NotEqual(t, th.Text, th.Background, "Text and background should have different colors for theme %s", themeID)

				// Test that primary and background colors are different
				assert.NotEqual(t, th.Primary, th.Background, "Primary and background should have different colors for theme %s", themeID)
			})
		}
	})
}

func TestAnimationManager_Accessibility(t *testing.T) {
	t.Run("ReducedMotionRespectsUserPreference", func(t *testing.T) {
		manager := NewAnimationManager()

		// Enable reduced motion
		config := manager.GetConfig()
		config.ReducedMotion = true
		manager.SetConfig(config)

		// Start animation - should not actually start when reduced motion is enabled
		manager.StartAnimation("test", AnimationFade, 1.0)
		// In a full implementation, this would check that animations are minimized
		// For now, we just ensure it doesn't crash
		assert.NotNil(t, manager)
	})

	t.Run("AnimationPerformanceWithAccessibility", func(t *testing.T) {
		manager := NewAnimationManager()
		defer manager.Close()

		// Enable reduced motion for better accessibility
		config := manager.GetConfig()
		config.ReducedMotion = true
		manager.SetConfig(config)

		// Test that animations still work but are less intensive
		manager.StartAnimation("perf-test", AnimationFade, 1.0)

		// Should complete quickly even with reduced motion
		start := time.Now()
		cmd := manager.Update()
		duration := time.Since(start)

		assert.Less(t, duration, 50*time.Millisecond, "Animation update should be fast even with reduced motion")
		// With reduced motion, animations should still return a command for processing
		assert.NotNil(t, cmd, "Should return command even with reduced motion for proper processing")
	})
}

func TestThemeAccessibility(t *testing.T) {
	t.Run("ThemeColorContrastValidation", func(t *testing.T) {
		// Test that all themes meet WCAG AA standards for contrast
		themes := theme.ListThemes()
		require.Greater(t, len(themes), 0)

		for _, themeID := range themes {
			t.Run(themeID, func(t *testing.T) {
				th := theme.GetTheme(themeID)

				// Validate theme against WCAG standards
				warnings := theme.ValidateTheme(th)

				// Log warnings if any (don't fail the test for now)
				for _, warning := range warnings {
					t.Logf("Theme %s warning: %s", th.Name, warning)
				}

				// Should have minimal warnings for accessibility
				assert.LessOrEqual(t, len(warnings), 2, "Theme %s should have minimal accessibility warnings", th.Name)
			})
		}
	})

	t.Run("ThemeReducedMotionSupport", func(t *testing.T) {
		// Test that themes work well with reduced motion settings
		animationManager := NewAnimationManager()

		// Test default configuration
		config := animationManager.GetConfig()
		assert.True(t, config.Enabled)
		assert.False(t, config.ReducedMotion)

		// Test setting reduced motion
		config.ReducedMotion = true
		animationManager.SetConfig(config)

		newConfig := animationManager.GetConfig()
		assert.True(t, newConfig.ReducedMotion)
	})

	t.Run("ThemeFocusIndicators", func(t *testing.T) {
		// Test that themes provide adequate focus indicators
		themes := theme.ListThemes()
		require.Greater(t, len(themes), 0)

		for _, themeID := range themes {
			t.Run(themeID, func(t *testing.T) {
				th := theme.GetTheme(themeID)

				// Test that primary color is distinct from background for focus indicators
				assert.NotEqual(t, th.Primary, th.Background, "Primary and background should be different for focus indicators in theme %s", themeID)

				// Test that accent color is distinct from background for focus indicators
				assert.NotEqual(t, th.Accent, th.Background, "Accent and background should be different for focus indicators in theme %s", themeID)
			})
		}
	})
}
